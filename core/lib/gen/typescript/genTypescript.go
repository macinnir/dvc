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

// GenerateTypescriptModels returns a string for a typscript types file
func GenerateTypescriptModels(config *lib.Config) error {

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

			if str, e = GenerateTypescriptModel(table); e != nil {
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
func GenerateTypescriptModel(table *schema.Table) (string, error) {

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
		fieldType := schema.DataTypeToTypescriptString(table.Columns[columnNames[k]].DataType)
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
		defaultVal := schema.DataTypeToTypescriptDefault(table.Columns[columnNames[k]].DataType)
		sb.WriteString("\t" + columnNames[k] + ": " + defaultVal + ",\n")
	}
	sb.WriteString("});\n\n")

	return sb.String(), nil
}

func GenerateTypesriptDTOs(config *lib.Config) {

	dtoPaths := []string{}

	files, _ := ioutil.ReadDir("core/definitions/dtos")
	for k := range files {
		dtoPaths = append(dtoPaths, "core/definitions/dtos/"+files[k].Name())
	}

	files, _ = ioutil.ReadDir("app/definitions/dtos")
	for k := range files {
		dtoPaths = append(dtoPaths, "app/definitions/dtos/"+files[k].Name())
	}

	lib.EnsureDir(config.TypescriptDTOsPath)

	for k := range dtoPaths {
		str, _ := GenerateTypescriptDTO(dtoPaths[k])
		basePath := path.Base(dtoPaths[k])
		dest := path.Join(config.TypescriptDTOsPath, basePath[0:len(basePath)-3]+".ts")
		// fmt.Println("Generating: ", dtoPaths[k], " ==> ", dest)
		ioutil.WriteFile(dest, []byte(str), 0777)
	}
}

// GenerateTypescriptType returns a string for a type in typescript
// TODO need a map of all types so that import paths can be used for struct and array types
// TODO test for struct types (apart from array types)
func GenerateTypescriptDTO(filePath string) (string, error) {

	ps, _ := lib.ParseStruct2(filePath)

	var sb strings.Builder
	sb.WriteString(`/**
 * Generated Code; DO NOT EDIT
 *
 * ` + ps.Name + `
 */
export type ` + ps.Name + ` = {

`)

	columnNames := make([]string, len(ps.Fields))
	n := 0
	for columnName := range ps.Fields {
		columnNames[n] = columnName
		n++
	}

	sort.Strings(columnNames)

	for k := range columnNames {
		// if filePath == "app/definitions/dtos/UpdateQuoteDTO.go" && (columnNames[k] == "Sales" || columnNames[k] == "Customers" || columnNames[k] == "Item") {
		// 	// fmt.Println(filePath)
		// 	fmt.Println(columnNames[k], " ==> ", ps.Fields[columnNames[k]])
		// }

		// TODO if the field type is a struct (or an array of structs) it needs to be imported
		fieldType := schema.GoTypeToTypescriptString(ps.Fields[columnNames[k]])
		// fmt.Println("FieldType: ", fieldType, ps.Fields[columnNames[k]])

		sb.WriteString("\t// " + columnNames[k] + " " + ps.Fields[columnNames[k]] + "\n")
		sb.WriteString("\t" + columnNames[k] + ": " + fieldType + ";\n\n")
	}

	sb.WriteString("}\n\n")

	sb.WriteString("// new" + ps.Name + " is a factory method for new " + ps.Name + " objects\n")
	sb.WriteString("export const new" + ps.Name + " = () : " + ps.Name + " => ({\n")
	for k := range columnNames {
		defaultVal := schema.GoTypeToTypescriptDefault(ps.Fields[columnNames[k]])
		sb.WriteString("\t" + columnNames[k] + ": " + defaultVal + ",\n")
	}
	sb.WriteString("});\n\n")

	return sb.String(), nil
}
