package connections

import (
	"fmt"

	"github.com/macinnir/dvc/core/lib"
	"go.uber.org/zap"
)

const CommandName = "connections"

// Import fetches the sql schema from the target database (specified in dvc.toml)
// and from that generates the json representation at `[schema name].schema.json`
func Cmd(log *zap.Logger, config *lib.Config, args []string) error {

	for k := range config.Databases {
		fmt.Println(" - ", config.Databases[k].Key, "(Schema: "+lib.ExtractRootNameFromKey(config.Databases[k].Key)+") => "+config.Databases[k].Host)
	}

	return nil
}
