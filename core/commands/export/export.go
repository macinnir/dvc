package export

import (
	"github.com/macinnir/dvc/core/lib"
	"go.uber.org/zap"
)

const CommandName = "export"

// Export export SQL create statements to standard out
func Cmd(logger *zap.Logger, config *lib.Config, args []string) error {

	return nil
	// var e error
	// var sql string

	// cmp := c.initCompare()

	// if sql, e = cmp.ExportSchemaToSQL(); e != nil {
	// 	lib.Error(e.Error(), c.Options)
	// 	os.Exit(1)
	// }

	// fmt.Println(sql)
}
