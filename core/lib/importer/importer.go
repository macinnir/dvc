package importer

import (
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
	rootNames := []string{}
	configs := map[string]*lib.ConfigDatabase{}
	rootNameGate := map[string]bool{}

	for k := range config.Databases {
		rootName := lib.ExtractRootNameFromKey(config.Databases[k].Key)
		if _, ok := rootNameGate[rootName]; !ok {
			rootNameGate[rootName] = true
			rootNames = append(rootNames, rootName)
			configs[rootName] = config.Databases[k]
		}
	}

	// 2. Loop through each database and get a copy its schema
	for k := range rootNames {

		rootName := rootNames[k]
		connector, e := connectors.DBConnectorFactory(configs[rootName])
		if e != nil {
			return nil, e
		}
		executor := executor.NewExecutor(
			configs[rootName],
			connector,
		).Connect()

		s := &schema.Schema{
			Name: rootName,
		}
		s.Tables, e = connector.FetchDatabaseTables(
			executor,
			configs[rootName].Name,
		)

		if e != nil {
			return nil, e
		}

		schemaList.Schemas = append(schemaList.Schemas, s)
	}
	// database.Enums = c.connector.FetchEnums(server)

	return schemaList, nil
}

// FetchAllSchemas pulls in the schemas for all databases
func FetchAllSchemas(config *lib.Config) (map[string]*schema.Schema, error) {

	schemas := map[string]*schema.Schema{}

	for k := range config.Databases {
		config := config.Databases[k]

		rootName := lib.ExtractRootNameFromKey(config.Key)
		connector, e := connectors.DBConnectorFactory(config)
		if e != nil {
			return nil, e
		}
		executor := executor.NewExecutor(
			config,
			connector,
		).Connect()

		s := &schema.Schema{
			Name: rootName,
		}
		s.Tables, e = connector.FetchDatabaseTables(
			executor,
			config.Name,
		)

		if e != nil {
			return nil, e
		}

		schemas[config.Key] = s
	}
	// database.Enums = c.connector.FetchEnums(server)

	return schemas, nil
}
