package gen

import (
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

// GenerateTypescriptModels returns a string for a typscript types file
func GenerateTypescriptModels(config *lib.Config, routes *lib.RoutesJSONContainer) error {

	var start = time.Now()

	lib.EnsureDir(config.TypescriptModelsPath)

	var generatedCount = 0

	var wg sync.WaitGroup
	var mutex = sync.Mutex{}

	var modelNames = []string{}
	for name := range routes.Models {
		modelNames = append(modelNames, name)
	}

	sort.Strings(modelNames)

	for _, n := range modelNames {

		wg.Add(1)

		go func(name string, modelProps map[string]string) {

			var e error
			var str string

			defer wg.Done()

			if str, e = GenerateTypescriptModel(name, modelProps); e != nil {
				log.Fatalf("Error generating typescript model %s: %s", name, e.Error())
			}

			var fullFilePath = path.Join(config.TypescriptModelsPath, name+".ts")
			// log.Printf("Generating typescript model %s => %s", name, fullFilePath)

			if e = os.WriteFile(fullFilePath, []byte(str), 0777); e != nil {
				log.Fatalf("Error writing to file %s: %s", fullFilePath, e.Error())
			}

			mutex.Lock()
			generatedCount++
			mutex.Unlock()

		}(n, routes.Models[n])
	}

	wg.Wait()
	lib.LogAdd(start, "%d typescript models to %s", generatedCount, config.TypescriptModelsPath)
	return nil
}

// GenerateTypescriptType returns a string for a type in typescript
func GenerateTypescriptModel(name string, columns map[string]string) (string, error) {

	var sb strings.Builder
	TSFileHeader(&sb, name)
	ImportStrings(&sb, columns, name)
	columnNames := ColumnMapToNames(columns)

	sb.WriteString(`
export type ` + name + ` = {

`)

	for k := range columnNames {
		dataType := columns[columnNames[k]]
		fieldType := schema.GoBaseTypeToBaseTypescriptType(dataType)
		sb.WriteString("\t// " + columnNames[k] + " " + dataType + "\n")
		sb.WriteString("\t" + columnNames[k] + ": " + fieldType + ";\n\n")
	}

	sb.WriteString(`}

// new` + name + ` is a factory method for creating new ` + name + ` objects
export const new` + name + ` = () : ` + name + ` => ({ 
`)

	for k := range columnNames {
		defaultVal := schema.GoBaseTypeToBaseTypescriptDefault(columns[columnNames[k]])
		sb.WriteString("\t" + columnNames[k] + ": " + defaultVal + ",\n")
	}
	sb.WriteString("});\n\n")

	return sb.String(), nil
}
