package compare

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"

	"github.com/macinnir/dvc/core/connectors"
	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

func objectsAreSame(local interface{}, remote interface{}) (same bool, e error) {

	// Remote Signature
	var localBytes []byte
	var remoteBytes []byte

	localBytes, e = json.Marshal(local)
	if e != nil {
		return
	}

	remoteBytes, e = json.Marshal(remote)
	if e != nil {
		return
	}

	// Local signature
	localHasher := sha1.New()
	_, e = localHasher.Write(localBytes)
	if e != nil {
		return
	}
	localSha := base64.URLEncoding.EncodeToString(localHasher.Sum(nil))
	// fmt.Printf("Local SHA %s\n", localSha)

	// Remote signature
	remoteHasher := sha1.New()
	_, e = remoteHasher.Write(remoteBytes)
	if e != nil {
		return
	}
	remoteSha := base64.URLEncoding.EncodeToString(remoteHasher.Sum(nil))
	// fmt.Printf("Remote SHA %s\n", remoteSha)

	same = localSha == remoteSha
	return
}

// NewCompare creates a new Compare instance
func NewCompare(config *lib.ConfigDatabase, connector connectors.IConnector) (compare *Compare, e error) {

	compare = &Compare{
		config:    config,
		connector: connector,
	}

	return
}

type Compare struct {
	config             *lib.ConfigDatabase       // Config is the config object
	connector          connectors.IConnector     // IConnector is the injected server manager
	LocalSQLPaths      []string                  // LocalSQLPaths is a list of paths pulled from the changesets.json file
	ChangesetSignature string                    // ChangesetSignature is a SHA signature for the changesets.json file
	LocalChangeFiles   []lib.ChangeFile          // LocalChangeFiles is a list of paths to local change files
	Files              *lib.Files                // Files is the injected file manager
	Databases          map[string]*schema.Schema // A map of databases
}

// ExportSchemaToSQL exports the current schema to sql
// func (c *Compare) ExportSchemaToSQL() (sql string, e error) {

// 	schemaFile := c.config.Name + ".schema.json"

// 	var localSchema *schema.Schema
// 	emptySchema := &schema.Schema{}

// 	if localSchema, e = schema.ReadSchemaFromFile(schemaFile); e != nil {
// 		return
// 	}

// 	sql, e = c.connector.CreateChangeSQL(localSchema, emptySchema)
// 	return
// }

// func (c *Compare) ExportDataToSQL(tableName string)
// allow for choosing individual tables
// func (c *Compare) ExportDataToJSON(tableName string)
// func (c *Compare) ImportDataFromJSON(tableName string)

// ApplyChangeset runs the sql produced by the `CompareSchema` command against the target database
// @command compare [reverse] apply
// func (c *Compare) ApplyChangeset(changeset string) (e error) {

// 	e = lib.NewExecutor(c.config, c.connector).RunSQL(changeset)
// 	return
// server := c.initCommand()

// statements := strings.Split(changeset, ";")

// defer server.Connection.Close()

// tx, _ := server.Connection.Begin()

// nonEmptyStatements := []string{}
// for _, s := range statements {
// 	if len(strings.Trim(strings.Trim(s, " "), "\n")) == 0 {
// 		continue
// 	}

// 	nonEmptyStatements = append(nonEmptyStatements, s)
// }

// for i, s := range nonEmptyStatements {
// 	sql := strings.Trim(strings.Trim(s, " "), "\n")
// 	if len(sql) == 0 {
// 		continue
// 	}
// 	// fmt.Printf("\rRunning %d of %d sql statements...", i+1, len(nonEmptyStatements))
// 	fmt.Printf("Running %d of %d: \n%s\n", i+1, len(nonEmptyStatements), sql)
// 	// lib.Debugf("Running sql: \n%s\n", c.Options, sql)

// 	_, e = tx.Exec(sql)
// 	if e != nil {
// 		tx.Rollback()
// 		return
// 	}
// }
// fmt.Print("Finished\n")
// e = tx.Commit()
// if e != nil {
// 	panic(e)
// }

// return
// }

// CompareSchema returns a string that contains a new line (`\n`) separated list of sql statements
// This comparison assumes the local `schemaFile` is the authority and the remote database is the
// schema to be updated
// @param reverse bool If true, the remote and local schema comparison is flipped in that the remote schema is treated as the authority
// 		and the local schema is treated as the schema to be updated.
// @command compare [reverse]
func CompareSchemas(
	config *lib.Config,
	localSchemaList *schema.SchemaList,
	remoteSchemas map[string]*schema.Schema,
) (string, error) {

	var e error

	// var same bool
	// if same, e = objectsAreSame(localSchemaList, remoteSchemaList); e != nil || same {
	// 	return "", nil
	// }

	configMap := map[string]*lib.ConfigDatabase{}
	for k := range config.Databases {
		configMap[config.Databases[k].Key] = config.Databases[k]
	}

	sql := ""

	localSchemaMap := map[string]int{}
	for k := range localSchemaList.Schemas {
		localSchemaMap[localSchemaList.Schemas[k].Name] = k
	}

	for key := range remoteSchemas {
		remoteSchema := remoteSchemas[key]
		rootName := lib.ExtractRootNameFromKey(key)

		var connector connectors.IConnector
		connector, e = connectors.DBConnectorFactory(configMap[key])

		localSchema := localSchemaList.Schemas[localSchemaMap[rootName]]
		sql += "\n\n"
		sql += "-- \n"
		sql += "-- Change for " + key + " (" + configMap[key].Host + "/" + configMap[key].Name + ")\n"
		sql += "-- \n\n"
		changeSQL := ""
		changeSQL, e = connector.CreateChangeSQL(localSchema, remoteSchema)
		if e != nil {
			return "", e
		}

		sql += changeSQL
		sql += "-- End for " + key + "\n"
	}

	// if len(localSchema.Enums) > 0 {
	// 	// Compare remote to local
	// 	for tableName := range localSchema.Enums {
	// 		if same, _ = objectsAreSame(localSchema.Enums[tableName], remoteSchema.Enums[tableName]); !same {
	// 			sql += c.connector.CompareEnums(remoteSchema, localSchema, tableName)
	// 		}
	// 	}
	// }

	// if len(sql) == 0 {
	// 	fmt.Printf("The schema objects were not the same, but no change sql was generated.\n\n Something strange is afoot...\n")
	// }

	return sql, nil
}
