package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/macinnir/dvc/lib"
)

// Import fetches the sql schema from the target database (specified in dvc.toml)
// and from that generates the json representation at `[schema name].schema.json`
func (c *Cmd) Import(args []string) {
	if len(args) > 0 && args[0] == "help" {
		helpImport()
		return
	}
	var e error
	cmp := c.initCompare()

	if e = cmp.ImportSchema("./" + c.Config.Connection.DatabaseName + ".schema.json"); e != nil {
		lib.Error(e.Error(), c.Options)
		os.Exit(1)
	}

	fmt.Println("Importing...")
	curDir, _ := os.Getwd()
	lib.Infof("Schema `%s`.`%s` imported to %s", c.Options, c.Config.Connection.Host, c.Config.Connection.DatabaseName, path.Join(curDir, c.Config.Connection.DatabaseName+".schema.json"))
}
