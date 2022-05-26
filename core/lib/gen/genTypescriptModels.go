package gen

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

// GenerateTypescriptModels returns a string for a typscript types file
func GenerateTypescriptModels(config *lib.Config, routes *lib.RoutesJSONContainer) error {

	var start = time.Now()

	lib.EnsureDir(config.TypescriptModelsPath)

	var e error
	var str string

	var generatedCount = 0

	for name := range routes.Models {

		model := routes.Models[name]

		if str, e = GenerateTypescriptModel(name, model); e != nil {
			return e
		}

		fullFilePath := path.Join(config.TypescriptModelsPath, name+".ts")

		if e = ioutil.WriteFile(fullFilePath, []byte(str), 0777); e != nil {
			return e
		}

		generatedCount++
	}

	fmt.Printf("Generated %d typescript models to %s in %f seconds\n", generatedCount, config.TypescriptModelsPath, time.Since(start).Seconds())

	return nil
}

// GenerateTypescriptType returns a string for a type in typescript
func GenerateTypescriptModel(name string, columns map[string]string) (string, error) {

	var sb strings.Builder
	TSFileHeader(&sb, name)
	ImportStrings(&sb, columns)
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
