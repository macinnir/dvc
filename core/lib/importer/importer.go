package importer

import (
	"errors"
	"fmt"

	"github.com/macinnir/dvc/core/connectors"
	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/executor"
	"github.com/macinnir/dvc/core/lib/schema"
)

// FetchAllUniqueSchemas pulls in the schemas for all databases with a unique root name
// Datasebases A_0, A_1, B_0, C_0, C_1 would pull in schemas for A_0, B_0, C_0 (using the first of the set as the reference point )
func FetchAllUniqueSchemas(config *lib.Config) (*schema.SchemaList, error) {

	schemaList := &schema.SchemaList{
		Schemas: []*schema.Schema{},
	}

	// 1. Get unique names for each schema
	schemaNames := []string{}
	configs := map[string]*lib.ConfigDatabase{}
	rootNameGate := map[string]bool{}

	for k := range config.Databases {
		// fmt.Println(k, config.Databases[k].Key)
		schemaName := lib.ExtractRootNameFromKey(config.Databases[k].Key)

		// fmt.Println(k, rootName)

		if _, ok := rootNameGate[schemaName]; !ok {
			rootNameGate[schemaName] = true
			schemaNames = append(schemaNames, schemaName)
			configs[schemaName] = config.Databases[k]
		}
	}

	// 2. Loop through each database and get a copy its schema
	for k := range schemaNames {

		schemaName := schemaNames[k]

		connector, e := connectors.DBConnectorFactory(configs[schemaName])
		if e != nil {
			return nil, e
		}
		executor := executor.NewExecutor(
			configs[schemaName],
			connector,
		).Connect()

		var s *schema.Database
		s, e = connector.FetchDatabase(
			executor,
			configs[schemaName].Name,
		)

		if e != nil {
			return nil, e
		}

		schemaList.Schemas = append(schemaList.Schemas, s.ToSchema(schemaName))
	}
	// database.Enums = c.connector.FetchEnums(server)

	return schemaList, nil
}

func FetchSchema(config *lib.Config, schemaName, connectionName string) (*schema.Schema, error) {

	if len(schemaName) == 0 {
		return nil, errors.New("Schema Name cannot be empty")
	}

	if len(connectionName) == 0 {
		return nil, errors.New("Connection name cannot empty")
	}

	schemaToConnectionMap := map[string]map[string]bool{}

	for k := range config.Databases {

		thisSchemaName := lib.ExtractRootNameFromKey(config.Databases[k].Key)

		if _, ok := schemaToConnectionMap[thisSchemaName]; !ok {
			schemaToConnectionMap[thisSchemaName] = map[string]bool{}
		}

		schemaToConnectionMap[thisSchemaName][config.Databases[k].Key] = true
	}

	if _, ok := schemaToConnectionMap[schemaName]; !ok {
		return nil, fmt.Errorf("Unknown schema name `%s`", schemaName)
	}

	if _, ok := schemaToConnectionMap[schemaName][connectionName]; !ok {
		return nil, fmt.Errorf("Unknown connection name `%s` for schema `%s`", connectionName, schemaName)
	}

	var connectionConfig *lib.ConfigDatabase

	for k := range config.Databases {
		if config.Databases[k].Key == connectionName {
			connectionConfig = config.Databases[k]
			break
		}
	}

	if connectionConfig == nil {
		return nil, fmt.Errorf("Connection `%s` not found", connectionName)
	}

	if lib.ExtractRootNameFromKey(connectionConfig.Key) != schemaName {
		return nil, fmt.Errorf("Invalid schema name `%s` for connection `%s`", schemaName, connectionName)
	}

	connector, e := connectors.DBConnectorFactory(connectionConfig)
	if e != nil {
		return nil, e
	}

	executor := executor.NewExecutor(
		connectionConfig,
		connector,
	).Connect()

	var s *schema.Database
	s, e = connector.FetchDatabase(
		executor,
		connectionConfig.Name,
	)

	if e != nil {
		return nil, e
	}

	return s.ToSchema(schemaName), nil
}

// FetchAllSchemas pulls in the schemas for all databases
func FetchAllSchemas(config *lib.Config) (map[string]*schema.Schema, error) {

	schemas := map[string]*schema.Schema{}

	for k := range config.Databases {
		config := config.Databases[k]

		schemaName := lib.ExtractRootNameFromKey(config.Key)
		connector, e := connectors.DBConnectorFactory(config)
		if e != nil {
			return nil, e
		}
		executor := executor.NewExecutor(
			config,
			connector,
		).Connect()

		var database *schema.Database
		database, e = connector.FetchDatabase(
			executor,
			config.Name,
		)

		if e != nil {
			return nil, e
		}

		schemas[config.Key] = database.ToSchema(schemaName)
	}
	// database.Enums = c.connector.FetchEnums(server)

	return schemas, nil
}

func IsCoreSchemaName(schemaName string) bool {
	return schemaName == lib.CoreSchemasLogName || schemaName == lib.CoreSchemasName
}
