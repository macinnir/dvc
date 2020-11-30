package commands

import (
	"fmt"

	"github.com/macinnir/dvc/modules/data"
)

// Data is the base command for managing static data sets
func (c *Cmd) Data(args []string) {

	if len(args) > 0 && args[0] == "help" {
		helpData()
		return
	}

	database := c.loadDatabase()
	connector, _ := connectorFactory(c.Config.DatabaseType, c.Config)

	d, e := data.NewData(c.Config, c.Options, database, connector)

	if e != nil {
		panic(e)
	}

	if len(args) == 0 {
		fmt.Println("Missing argument [help | import | apply]")
		return
	}

	cmd := args[0]
	args = args[1:]

	switch cmd {
	case "import":
		d.Import(args)
	case "apply":
		d.Apply(args)
	case "rm":
		d.Remove(args)
	default:
		fmt.Println("Unknown command")
	}
}

func helpData() {
	fmt.Println(`
	data [import|apply|rm] [[table]]

		import - Import data from database tables to static files 

				Without any arguments, import looks for any files in the meta/data and imports data from the tables with the same names as the file 

				import [table] imports data for a single table, creating the table file if it does not exist. 

		apply - Apply data to the database from a static file 

				Without any arguments, import looks for any files in the meta/data folder and applies the contents of those files to the database tables with the same name as the files respectively

				import [table] applies data for a single table (if it exists) to the database. 
		
		rm [table] - Remove a data file 

		help - This output 
	
	`)
}
