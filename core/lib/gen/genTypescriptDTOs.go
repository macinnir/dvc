package gen

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

func CleanTypescriptDTOs(config *lib.Config, routes *lib.RoutesJSONContainer) error {

	// var start = time.Now()
	var e error
	var files []os.FileInfo

	if files, e = ioutil.ReadDir(config.TypescriptDTOsPath); e != nil {
		return e
	}

	var removedCount = 0

	for k := range files {

		if files[k].IsDir() {
			continue
		}

		if files[k].Name()[0:1] == "." {
			continue
		}

		dtoName := files[k].Name()[0 : len(files[k].Name())-3]

		if _, ok := routes.DTOs[dtoName]; !ok {
			tsFilePath := path.Join(config.TypescriptDTOsPath, files[k].Name())
			// fmt.Println("Removing", tsFilePath)
			os.Remove(tsFilePath)
			removedCount++
		}

	}

	// TODO Verbose flag
	// fmt.Printf("Removed %d typescript DTOs from `%s` in %f seconds\n", removedCount, config.TypescriptDTOsPath, time.Since(start).Seconds())

	return nil
}

func GenerateTypesriptDTOs(config *lib.Config, routes *lib.RoutesJSONContainer) error {

	// var start = time.Now()
	var generatedCount = 0

	lib.EnsureDir(config.TypescriptDTOsPath)

	for name := range routes.DTOs {
		tsDTOBytes, _ := GenerateTypescriptDTO(name, routes.DTOs[name])
		fullFilePath := path.Join(config.TypescriptDTOsPath, name+".ts")
		err := os.WriteFile(fullFilePath, tsDTOBytes, lib.DefaultFileMode)
		if err != nil {
			panic(err)
		}
		generatedCount++
	}

	// TODO Verbose
	// fmt.Printf("Generated %d typescript DTOs from `%s` in %f seconds\n", generatedCount, config.TypescriptDTOsPath, time.Since(start).Seconds())

	return nil
}

var typescriptDTOTemplate = template.Must(template.New("typescript-dto-file").Parse(`/**
* Generated Code; DO NOT EDIT
*
* {{.Name}}
*/

{{.Imports}}

export type {{.Name}} = { 
	{{range .Columns}}
		// {{.Name}} {{.DataType}}
		{{.Name}}: {{.FieldType}};{{end}}
}

// new {{.Name}} is a factory method for new {{.Name}} objects
export const new{{.Name}} = () : {{.Name}} => ({ 
	{{range .Columns}}
	{{.Name}}: {{.DefaultValue}},{{end}}
});

 `))

// 0.010107

// GenerateTypescriptType returns a string for a type in typescript
// TODO need a map of all types so that import paths can be used for struct and array types
// TODO test for struct types (apart from array types)
func GenerateTypescriptDTO(name string, columns map[string]string) ([]byte, error) {

	var buf bytes.Buffer

	var data = struct {
		Name    string
		Imports string
		Columns []struct {
			Name         string
			FieldType    string
			DataType     string
			DefaultValue string
		}
	}{
		Imports: GenDTOImportStrings(columns),
		Name:    name,
	}

	// ps, _ := lib.ParseStruct2(filePath)

	// TSFileHeader(&buf, name)

	columnNames := ColumnMapToNames(columns)

	// 	sb.WriteString(`
	// export type ` + name + ` = {

	// `)

	for k := range columnNames {

		dataType := columns[columnNames[k]]

		// if filePath == "app/definitions/dtos/UpdateQuoteDTO.go" && (columnNames[k] == "Sales" || columnNames[k] == "Customers" || columnNames[k] == "Item") {
		// 	// fmt.Println(filePath)
		// 	fmt.Println(columnNames[k], " ==> ", ps.Fields[columnNames[k]])
		// }

		// TODO if the field type is a struct (or an array of structs) it needs to be imported
		fieldType := schema.GoTypeToTypescriptString(dataType)
		// fmt.Println("FieldType: ", fieldType, ps.Fields[columnNames[k]])

		data.Columns = append(data.Columns, struct {
			Name         string
			FieldType    string
			DataType     string
			DefaultValue string
		}{
			Name:         columnNames[k],
			DataType:     dataType,
			FieldType:    fieldType,
			DefaultValue: schema.GoTypeToTypescriptDefault(dataType),
		})

		// sb.WriteString("\t// " + columnNames[k] + " " + dataType + "\n")
		// sb.WriteString("\t" + columnNames[k] + ": " + fieldType + ";\n\n")
	}

	// sb.WriteString("}\n\n")

	// sb.WriteString("// new" + name + " is a factory method for new " + name + " objects\n")
	// sb.WriteString("export const new" + name + " = () : " + name + " => ({\n")
	// for k := range columnNames {
	// 	dataType := columns[columnNames[k]]
	// 	defaultVal := schema.GoTypeToTypescriptDefault(dataType)
	// 	sb.WriteString("\t" + columnNames[k] + ": " + defaultVal + ",\n")
	// }
	// sb.WriteString("});\n\n")

	typescriptDTOTemplate.Execute(&buf, data)

	return buf.Bytes(), nil
}

func GenDTOImportStrings(columns map[string]string) string {

	var buf bytes.Buffer
	ImportStrings(&buf, columns)
	return buf.String()

	imported := map[string]struct{}{}

	// imports := [][]string{}
	for name := range columns {

		if !isImportable(name) {
			continue
		}

		dataType := columns[name]

		baseType := schema.ExtractBaseGoType(dataType)

		if !schema.IsGoTypeBaseType(baseType) {

			if isConstant(baseType) {
				continue
			}

			// Already imported
			if _, ok := imported[baseType]; ok {
				continue
			}

			imported[baseType] = struct{}{}

			if len(baseType) > 7 && baseType[0:7] == "models." {
				ImportString(&buf, dataType, baseType[7:], "gen/models/"+baseType[7:], false)
			} else if len(baseType) > 5 && baseType[0:5] == "dtos." {
				ImportString(&buf, dataType, baseType[5:], "gen/dtos/"+baseType[5:], false)
			} else {
				ImportString(&buf, dataType, baseType, "./"+baseType, false)
			}
		}
	}

	return buf.String()
}
