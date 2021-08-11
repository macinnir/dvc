package typescript

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

// GenerateTypescriptTypes returns a string for a typscript types file
func GenerateTypescriptTypes(config *lib.Config) error {

	lib.EnsureDir(config.TypescriptModelsPath)

	var e error
	var files []os.FileInfo

	fmt.Println("Checking for files in ", config.TypescriptModelsPath)

	if files, e = ioutil.ReadDir(config.TypescriptModelsPath); e != nil {
		return e
	}

	for k := range files {
		if files[k].IsDir() {
			continue
		}

		if files[k].Name()[0:1] == "." {
			continue
		}

		os.Remove(path.Join(config.TypescriptModelsPath, files[k].Name()))
	}

	var schemaList *schema.SchemaList

	if schemaList, e = schema.LoadLocalSchemas(); e != nil {
		return e
	}

	for k := range schemaList.Schemas {

		database := schemaList.Schemas[k]

		for k := range database.Tables {

			table := database.Tables[k]

			var str string

			if str, e = GenerateTypescriptType(table); e != nil {
				return e
			}

			if e = ioutil.WriteFile(path.Join(config.TypescriptModelsPath, table.Name+".ts"), []byte(str), 0777); e != nil {
				return e
			}

		}
	}

	return nil
}

// GenerateTypescriptType returns a string for a type in typescript
func GenerateTypescriptType(table *schema.Table) (string, error) {

	var sb strings.Builder
	sb.WriteString(`/**
 * Generated Code; DO NOT EDIT
 * 
 * ` + table.Name + `
 */
export type ` + table.Name + ` = {

`)

	columnNames := make([]string, len(table.Columns))
	n := 0
	for columnName := range table.Columns {
		columnNames[n] = columnName
		n++
	}

	sort.Strings(columnNames)

	for k := range columnNames {
		fieldType := schema.DataTypeToTypescriptString(table.Columns[columnNames[k]])
		if table.Columns[columnNames[k]].MaxLength > 0 {
			sb.WriteString("\t// " + columnNames[k] + " " + table.Columns[columnNames[k]].DataType + "(" + fmt.Sprint(table.Columns[columnNames[k]].MaxLength) + ")\n")
		} else {
			sb.WriteString("\t// " + columnNames[k] + " " + table.Columns[columnNames[k]].DataType + "\n")
		}
		sb.WriteString("\t" + columnNames[k] + ": " + fieldType + ";\n\n")
	}

	sb.WriteString("}\n\n")

	sb.WriteString("// new" + table.Name + " is a factory method for new " + table.Name + " objects\n")
	sb.WriteString("export const new" + table.Name + " = () : " + table.Name + " => ({\n")
	for k := range columnNames {
		defaultVal := schema.DataTypeToTypescriptDefault(table.Columns[columnNames[k]])
		sb.WriteString("\t" + columnNames[k] + ": " + defaultVal + ",\n")
	}
	sb.WriteString("});\n\n")

	return sb.String(), nil
}
