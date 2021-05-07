package clone

import (
	"github.com/macinnir/dvc/core/lib"
	"go.uber.org/zap"
)

const CommandName = "clone"

type Clone struct{}

// Clone clones a database
func Cmd(logger *zap.Logger, config *lib.Config, args []string) error {
	return nil
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
