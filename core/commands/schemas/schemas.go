package schemas

import (
	"fmt"

	"github.com/macinnir/dvc/core/lib"
	"go.uber.org/zap"
)

const CommandName = "schemas"

// Import fetches the sql schema from the target database (specified in dvc.toml)
// and from that generates the json representation at `[schema name].schema.json`
func Cmd(log *zap.Logger, config *lib.Config, args []string) error {

	schemaOrdered := []string{}
	schemaMap := map[string][]string{}

	for k := range config.Databases {
		schemaName := lib.ExtractRootNameFromKey(config.Databases[k].Key)
		if _, ok := schemaMap[schemaName]; !ok {
			schemaMap[schemaName] = []string{}
			schemaOrdered = append(schemaOrdered, schemaName)
		}

		schemaMap[schemaName] = append(schemaMap[schemaName], config.Databases[k].Key+" -> "+config.Databases[k].Host+"/"+config.Databases[k].Name)
	}

	for k := range schemaOrdered {
		fmt.Println(" +  " + schemaOrdered[k])
		for l := range schemaMap[schemaOrdered[k]] {
			fmt.Println("    - ", schemaMap[schemaOrdered[k]][l])
		}
	}

	return nil
}
