package commands

import (
	"fmt"
	"os"

	"github.com/macinnir/dvc/lib"
)

// Export export SQL create statements to standard out
func (c *Cmd) Export(args []string) {

	if len(args) > 0 && args[0] == "help" {
		helpExport()
		return
	}

	var e error
	var sql string

	cmp := c.initCompare()

	if sql, e = cmp.ExportSchemaToSQL(); e != nil {
		lib.Error(e.Error(), c.Options)
		os.Exit(1)
	}

	fmt.Println(sql)
}
