package transfer

import (
	"errors"
	"fmt"
	"os"

	"github.com/macinnir/dvc/core/connectors"
	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/executor"
	"github.com/macinnir/dvc/core/lib/importer"
	"github.com/macinnir/dvc/core/lib/schema"
	"go.uber.org/zap"
)

const CommandName = "transfer"

type Transfer struct{}

// Clone clones a database
func Cmd(logger *zap.Logger, config *lib.Config, args []string) error {

	var e error
	var remoteSchemas map[string]*schema.Schema
	doRun := false

	fromDatabase := ""
	toDatabase := ""
	tableName := ""

	if len(args) < 2 {
		fmt.Println("Usage: dvc transfer [source_database] [destination_database] [[table_name]] [[-r|--run]]")
		return errors.New("Insufficient arguments...")
	}
	fromDatabase = args[0]
	toDatabase = args[1]

	if fromDatabase == toDatabase {
		return errors.New("Source and destination databases must be different")
	}

	if len(args) > 2 {
		doRun = args[2] == "-r" || args[2] == "--run"
		if !doRun {
			tableName = args[2]
		}
	}

	if len(args) > 3 {
		doRun = args[3] == "-r" || args[3] == "--run"
	}

	fmt.Println("Args: ", len(args))
	fmt.Printf("From: %s; To: %s; Table: %s\n", fromDatabase, toDatabase, tableName)

	if remoteSchemas, e = importer.FetchAllSchemas(config); e != nil {
		return e
	}

	if _, ok := remoteSchemas[fromDatabase]; !ok {
		return errors.New("Unknown database key: " + fromDatabase)
	}

	sql := []string{}
	tables := []string{}

	configs := map[string]*lib.ConfigDatabase{}
	for k := range config.Databases {
		configs[config.Databases[k].Key] = config.Databases[k]
	}

	// Migrate single table
	if len(tableName) > 0 {
		if _, ok := remoteSchemas[fromDatabase].Tables[tableName]; !ok {
			return errors.New("Table `" + tableName + "` not found in source schema")
		}

		tables = []string{tableName}
	} else {
		for tableName := range remoteSchemas[fromDatabase].Tables {
			tables = append(tables, tableName)
		}
	}

	sourceDatabaseConfig := configs[fromDatabase]
	sourceDatabaseName := sourceDatabaseConfig.Name

	for k := range tables {
		sql = append(sql, transferSQL(sourceDatabaseName, toDatabase, tables[k])...)
	}

	if doRun {
		var connector connectors.IConnector
		if connector, e = connectors.DBConnectorFactory(sourceDatabaseConfig); e != nil {
			return e
		}

		server := executor.NewExecutor(
			sourceDatabaseConfig,
			connector,
		).Connect()

		for k := range sql {
			fmt.Println("RUNNING: " + sql[k])
			_, e = server.Connection.Exec(sql[k])
			if e != nil {
				fmt.Printf("SQL ERROR: %s\n", e.Error())
				fmt.Printf("SQL: %s\n", sql[k])
				os.Exit(1)
			}
		}
	} else {

		for k := range sql {
			fmt.Println(sql[k])
		}
	}
	// connector, e := connectors.DBConnectorFactory(config)

	// for k := range remoteSchemas {
	// 	fmt.Println(k)
	// }

	return nil
}

func transferSQL(fromDatabase, toDatabase, tableName string) []string {
	return []string{
		fmt.Sprintf("CREATE TABLE `%s`.`%s` LIKE `%s`.`%s`", toDatabase, tableName, fromDatabase, tableName),
		fmt.Sprintf("INSERT INTO `%s`.`%s` SELECT * FROM `%s`.`%s`", toDatabase, tableName, fromDatabase, tableName),
		fmt.Sprintf("DROP TABLE `%s`.`%s`", fromDatabase, tableName),
	}
}

// 	config := &lib.Config{}

// 	if len(args) == 0 {
// 		lib.Error("Missing target database name", c.Options)
// 		return
// 	}

// 	targetDatabase := args[0]

// 	config.DatabaseType = c.Config.DatabaseType
// 	config.Connection.Username = c.Config.Connection.Username
// 	config.Connection.Password = c.Config.Connection.Password
// 	config.Connection.Host = c.Config.Connection.Host
// 	config.Connection.DatabaseName = targetDatabase

// 	connector, _ := connectorFactory(c.Config.DatabaseType, config)
// 	cmp, _ := compare.NewCompare(config, c.Options, connector)

// 	// Do the comparison
// 	schemaFile := c.Config.Connection.DatabaseName + ".schema.json"
// 	var sql string
// 	var e error

// 	if sql, e = cmp.CompareSchema(schemaFile); e != nil {
// 		lib.Error(e.Error(), c.Options)
// 		os.Exit(1)
// 	}

// 	if len(sql) == 0 {
// 		fmt.Printf("No schema changes found with target %s\n", config.Connection.DatabaseName)
// 		lib.Info("No changes found.", c.Options)
// 		os.Exit(0)
// 	} else {
// 		writeSQLToLog(sql)
// 		e = cmp.ApplyChangeset(sql)
// 		if e != nil {
// 			lib.Error(e.Error(), c.Options)
// 			os.Exit(1)
// 		}
// 	}

// 	// database := c.loadDatabase()

// 	// reader := bufio.NewReader(os.Stdin)

// 	// dbName := ""

// 	// if len(args) > 0 {
// 	// 	dbName = args[0]
// 	// } else {
// 	// 	dbName = lib.ReadCliInput(reader, "Destination Database:")
// 	// }

// 	// if len(dbName) == 0 {
// 	// 	fmt.Println("No target database specified")
// 	// 	return
// 	// }

// 	// fmt.Println(sql)

// 	// doInsertYN := lib.ReadCliInput(reader, "Run above SQL (Y/n")

// 	// if doInsertYN != "Y" {
// 	// 	return
// 	// }

// 	// connector, _ := connectorFactory(c.Config.DatabaseType, c.Config)
// 	// x := lib.NewExecutor(c.Config, connector)
// 	// x.RunSQL(sql)
// }
