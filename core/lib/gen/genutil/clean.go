package genutil

import (
	"fmt"
	"os"
	"path"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

// CleanFiles removes model files that are not found in the database.Tables map
func CleanFiles(name string, dir string, schemaList *schema.SchemaList, prefix, suffix string) error {

	lib.EnsureDir(dir)

	// var start = time.Now()
	var e error
	var dirHandle *os.File

	dirHandle, e = os.Open(dir)
	if e != nil {
		return e
	}

	defer dirHandle.Close()
	var modelFiles []string
	modelFiles, e = dirHandle.Readdirnames(-1)
	if e != nil {
		return e
	}

	var removedCount = 0

	for k := range modelFiles {

		var fileName = modelFiles[k]
		var modelName = ParseFileNameToModelName(fileName, prefix, suffix)

		// TODO This needs to be configurable
		if fileName == "mocks_test.go" {
			continue
		}

		// go, ts
		if _, ok := schemaList.TableMap[modelName]; !ok {
			fullFilePath := path.Join(dir, fileName)
			// TODO Verbose flag
			fmt.Printf("Deleting `%s` (`%s`)\n", fullFilePath, modelName)
			os.Remove(fullFilePath)
			removedCount++
		}
	}

	// TODO Verbose flag
	// fmt.Printf("Removed %d %s from `%s` in %f seconds\n", removedCount, name, dir, time.Since(start).Seconds())

	return nil
}
