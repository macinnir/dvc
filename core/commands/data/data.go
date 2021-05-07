package data

import (
	"github.com/macinnir/dvc/core/lib"
	"go.uber.org/zap"
)

const CommandName = "data"

// Data is the base command for managing static data sets
func Cmd(logger *zap.Logger, config *lib.Config, args []string) error {

	return nil

	// database := c.loadDatabase()
	// connector, _ := connectorFactory(c.Config.DatabaseType, c.Config)

	// d, e := data.NewData(c.Config, c.Options, database, connector)

	// if e != nil {
	// 	panic(e)
	// }

	// if len(args) == 0 {
	// 	fmt.Println("Missing argument [help | import | apply]")
	// 	return
	// }

	// cmd := args[0]
	// args = args[1:]

	// switch cmd {
	// case "import":
	// 	d.Import(args)
	// case "apply":
	// 	d.Apply(args)
	// case "rm":
	// 	d.Remove(args)
	// default:
	// 	fmt.Println("Unknown command")
	// }
}
