package importcmd

import (
	"encoding/json"
	"io/ioutil"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/importer"
	"github.com/macinnir/dvc/core/lib/schema"
	"go.uber.org/zap"
)

const CommandName = "import"

// Import fetches the sql schema from the target database (specified in dvc.toml)
// and from that generates the json representation at `[schema name].schema.json`
func Cmd(log *zap.Logger, config *lib.Config, args []string) error {

	var e error
	var schemaList *schema.SchemaList

	schemaList, e = importer.FetchAllUniqueSchemas(config)
	if e != nil {
		log.Error("Error importing all schemas", zap.Error(e))
		return e
	}

	var dbBytes []byte
	dbBytes, _ = json.MarshalIndent(schemaList, " ", "    ")
	e = ioutil.WriteFile(lib.SchemasFilePath, dbBytes, 0644)

	return nil
}
