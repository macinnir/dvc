package importcmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/importer"
	"github.com/macinnir/dvc/core/lib/schema"
	"go.uber.org/zap"
)

const CommandName = "import"

// Import fetches the sql schema from the target database (specified in dvc.toml)
// and from that generates the json representation at `[schema name].schema.json`
func Cmd(log *zap.Logger, config *lib.Config, args []string) error {

	if len(args) > 0 {
		return ImportSingleSchema(config, args)
	}

	return ImportAll(log, config)
}

func ImportSingleSchema(config *lib.Config, args []string) error {
	schemaName := ""
	connectionName := ""

	if len(args) > 0 {
		schemaName = args[0]
	}

	if len(args) > 1 {
		connectionName = args[1]
	}

	schemaType := "app"
	isCore := importer.IsCoreSchemaName(schemaName)
	if isCore {
		schemaType = "core"
	}

	fmt.Printf("Importing %s schema `%s` from `%s`\n", schemaType, schemaName, connectionName)

	var e error
	var remoteSchema *schema.Schema
	if remoteSchema, e = importer.FetchSchema(config, schemaName, connectionName); e != nil {
		return fmt.Errorf("Import: FetchSchema: %w", e)
	}

	srcFile := ""

	if isCore {
		// Write to core schema
		srcFile = lib.CoreSchemasFilePath
		// fmt.Println("Writing to ", lib.CoreSchemasFilePath)
	} else {
		// Write to app schema
		srcFile = lib.SchemasFilePath
		// fmt.Println("Writing to ", lib.SchemasFilePath)
	}

	localSchemas := &schema.SchemaList{}
	srcBytes, _ := ioutil.ReadFile(srcFile)
	if e = json.Unmarshal(srcBytes, localSchemas); e != nil {
		return fmt.Errorf("unmarshal schema: %w", e)
	}

	targetSchemaKey := -1

	for k := range localSchemas.Schemas {
		if localSchemas.Schemas[k].Name == schemaName {
			targetSchemaKey = k
			break
		}
	}

	// Remote schema does not exist on the local set of schemas
	// Add it
	if targetSchemaKey == -1 {
		localSchemas.Schemas = append(localSchemas.Schemas, remoteSchema)
	} else {
		localSchemas.Schemas[targetSchemaKey] = remoteSchema
	}

	var dbBytes []byte
	dbBytes, _ = json.MarshalIndent(localSchemas, " ", "    ")

	return os.WriteFile(srcFile, dbBytes, 0777)
}
func ImportAll(log *zap.Logger, config *lib.Config) error {

	var start = time.Now()

	var e error
	var allSchemas *schema.SchemaList

	if allSchemas, e = importer.FetchAllUniqueSchemas(config); e != nil {
		log.Error("Error importing all schemas", zap.Error(e))
		return e
	}
	coreSchemaList := &schema.SchemaList{
		Schemas: []*schema.Schema{},
	}

	appSchemaList := &schema.SchemaList{
		Schemas: []*schema.Schema{},
	}

	// for k := range allSchemas.Schemas {
	// 	fmt.Println(k, "Schema: "+allSchemas.Schemas[k].Name)

	// }

	// os.Exit(1)
	for k := range allSchemas.Schemas {
		if importer.IsCoreSchemaName(allSchemas.Schemas[k].Name) {
			coreSchemaList.Schemas = append(coreSchemaList.Schemas, allSchemas.Schemas[k])
		} else {
			appSchemaList.Schemas = append(appSchemaList.Schemas, allSchemas.Schemas[k])
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Core Schema List
	go func() {
		defer wg.Done()
		dbBytes, _ := json.MarshalIndent(coreSchemaList, " ", "    ")
		if e = os.WriteFile(lib.CoreSchemasFilePath, dbBytes, 0644); e != nil {
			log.Fatal(fmt.Sprintf("Error writing core schemas: %s", e.Error()))
		}
	}()

	// App Schema List
	go func() {
		defer wg.Done()
		dbBytes, _ := json.MarshalIndent(appSchemaList, " ", "    ")
		if e = os.WriteFile(lib.SchemasFilePath, dbBytes, 0644); e != nil {
			log.Fatal(fmt.Sprintf("Error writing app schemas: %s", e.Error()))
		}
	}()

	lib.LogAdd(start, "imported all schemas")

	return nil
}
