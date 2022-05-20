package gen

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"unicode"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

// GenerateTypescriptModels returns a string for a typscript types file
func GenerateTypescriptModels(config *lib.Config, routes *lib.RoutesJSONContainer) error {

	lib.EnsureDir(config.TypescriptModelsPath)

	var e error
	var files []os.FileInfo

	// fmt.Println("Checking for files in ", config.TypescriptModelsPath)

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

		tsFilePath := path.Join(config.TypescriptModelsPath, files[k].Name())
		// fmt.Println("Removing", tsFilePath)
		os.Remove(tsFilePath)
	}

	var str string

	for name := range routes.Models {

		model := routes.Models[name]

		if str, e = GenerateTypescriptModel(name, model); e != nil {
			return e
		}

		fullFilePath := path.Join(config.TypescriptModelsPath, name+".ts")

		if e = ioutil.WriteFile(fullFilePath, []byte(str), 0777); e != nil {
			return e
		}
	}

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

func GenerateTypesriptDTOs(config *lib.Config, routes *lib.RoutesJSONContainer) error {

	lib.EnsureDir(config.TypescriptDTOsPath)

	var e error
	var files []os.FileInfo

	if files, e = ioutil.ReadDir(config.TypescriptDTOsPath); e != nil {
		return e
	}

	for k := range files {

		if files[k].IsDir() {
			continue
		}

		if files[k].Name()[0:1] == "." {
			continue
		}

		tsFilePath := path.Join(config.TypescriptDTOsPath, files[k].Name())
		// fmt.Println("Removing", tsFilePath)
		os.Remove(tsFilePath)
	}

	for name := range routes.DTOs {
		str, _ := GenerateTypescriptDTO(name, routes.DTOs[name])
		fullFilePath := path.Join(config.TypescriptDTOsPath, name+".ts")
		// fmt.Println("Generating DTO", name, " => ", fullFilePath)
		ioutil.WriteFile(fullFilePath, []byte(str), 0777)
	}

	return nil
}

// GenerateTypescriptType returns a string for a type in typescript
// TODO need a map of all types so that import paths can be used for struct and array types
// TODO test for struct types (apart from array types)
func GenerateTypescriptDTO(name string, columns map[string]string) (string, error) {

	// ps, _ := lib.ParseStruct2(filePath)

	var sb strings.Builder
	TSFileHeader(&sb, name)
	ImportStrings(&sb, columns)

	columnNames := ColumnMapToNames(columns)

	sb.WriteString(`
export type ` + name + ` = {

`)

	for k := range columnNames {

		dataType := columns[columnNames[k]]

		// if filePath == "app/definitions/dtos/UpdateQuoteDTO.go" && (columnNames[k] == "Sales" || columnNames[k] == "Customers" || columnNames[k] == "Item") {
		// 	// fmt.Println(filePath)
		// 	fmt.Println(columnNames[k], " ==> ", ps.Fields[columnNames[k]])
		// }

		// TODO if the field type is a struct (or an array of structs) it needs to be imported
		fieldType := schema.GoTypeToTypescriptString(dataType)
		// fmt.Println("FieldType: ", fieldType, ps.Fields[columnNames[k]])

		sb.WriteString("\t// " + columnNames[k] + " " + dataType + "\n")
		sb.WriteString("\t" + columnNames[k] + ": " + fieldType + ";\n\n")
	}

	sb.WriteString("}\n\n")

	sb.WriteString("// new" + name + " is a factory method for new " + name + " objects\n")
	sb.WriteString("export const new" + name + " = () : " + name + " => ({\n")
	for k := range columnNames {
		dataType := columns[columnNames[k]]
		defaultVal := schema.GoTypeToTypescriptDefault(dataType)
		sb.WriteString("\t" + columnNames[k] + ": " + defaultVal + ",\n")
	}
	sb.WriteString("});\n\n")

	return sb.String(), nil
}

func NewTypescriptGenerator(config *lib.Config, routes *lib.RoutesJSONContainer) *TypescriptGenerator {
	return &TypescriptGenerator{
		routes,
		config,
	}
}

func (tg *TypescriptGenerator) GenerateTypesriptAggregates(config *lib.Config) error {

	lib.EnsureDir(tg.config.TypescriptAggregatesPath)

	var e error
	var files []os.FileInfo

	if files, e = ioutil.ReadDir(config.TypescriptAggregatesPath); e != nil {
		return e
	}

	for k := range files {

		if files[k].IsDir() {
			continue
		}

		if files[k].Name()[0:1] == "." {
			continue
		}

		tsFilePath := path.Join(config.TypescriptAggregatesPath, files[k].Name())
		// fmt.Println("Removing", tsFilePath)
		os.Remove(tsFilePath)
	}

	for name := range tg.routes.Aggregates {
		str, e := tg.GenerateTypescriptAggregate(name)
		if e != nil {
			fmt.Println("ERROR:", e.Error())
			return e
		}
		dest := path.Join(tg.config.TypescriptAggregatesPath, name+".ts")
		ioutil.WriteFile(dest, []byte(str), 0777)
	}

	return nil
}

type TypescriptGenerator struct {
	routes *lib.RoutesJSONContainer
	config *lib.Config
}

// GenerateTypescriptAggregate returns a string for a type in typescript
func (tg *TypescriptGenerator) GenerateTypescriptAggregate(name string) (string, error) {

	columns := tg.routes.Aggregates[name]

	var sb strings.Builder

	TSFileHeader(&sb, name)

	ImportStrings(&sb, columns)

	sb.WriteString(`
export type ` + name + ` = `)

	inherits := InheritStrings(&sb, columns)

	if len(inherits) > 0 {
		sb.WriteString(strings.Join(inherits, " & ") + " & ")
	}

	sb.WriteString(`{

`)

	tg.GenerateTypescriptFields(&sb, name)

	sb.WriteString("}\n\n")

	sb.WriteString("// new" + name + " is a factory method for new " + name + " objects\n")
	sb.WriteString("export const new" + name + " = () : " + name + " => ({\n")

	tg.GenerateTypescriptDefaults(&sb, name)

	sb.WriteString("});\n\n")

	return sb.String(), nil
}

func (tg *TypescriptGenerator) ExtractColumns(goType string) map[string]string {

	if goType[0:1] == "*" {
		goType = goType[1:]
	}

	// fmt.Println("GoType", goType)
	// if goType == "CustomerAggregate" {
	// 	fmt.Println(goType, " --> ", goType[len(goType)-9:])
	// }

	if len(goType) > 3 && goType[len(goType)-3:] == "DTO" {
		return tg.routes.DTOs[goType]
	}

	if len(goType) > 9 && goType[len(goType)-9:] == "Aggregate" {
		cols := tg.routes.Aggregates[goType]
		// fmt.Println("Got here", goType, len(cols))
		// if !ok {
		// 	fmt.Println("Nope")
		// }
		return cols

	}

	return tg.routes.Models[goType]
}

func (tg *TypescriptGenerator) GenerateTypescriptFields(sb io.Writer, objectName string) {

	columns := tg.ExtractColumns(objectName)

	columnNames := ColumnMapToNames(columns)

	for k := range columnNames {

		name := columnNames[k]

		// Uppercase fields only
		if !unicode.IsUpper(rune(name[0])) {
			continue
		}

		goType := columns[columnNames[k]]
		fieldType := schema.GoTypeToTypescriptString(goType)

		if len(name) > 9 && name[0:9] == "#embedded" {
			// tg.GenerateTypescriptFields(sb, fieldType)
			continue
		}

		fmt.Fprintf(sb, "\t// %s %s\n", columnNames[k], columns[columnNames[k]])
		fmt.Fprintf(sb, "\t%s: %s;\n\n", columnNames[k], fieldType)
	}
}

func (tg *TypescriptGenerator) GenerateTypescriptDefaults(sb io.Writer, objectName string) {

	columns := tg.ExtractColumns(objectName)

	columnNames := ColumnMapToNames(columns)

	for k := range columnNames {
		name := columnNames[k]
		goType := columns[columnNames[k]]
		fieldType := schema.GoTypeToTypescriptString(goType)

		if len(name) > 9 && name[0:9] == "#embedded" {
			fmt.Fprintf(sb, "\t..."+schema.GoTypeToTypescriptDefault(fieldType)+",\n")
			// tg.GenerateTypescriptDefaults(sb, fieldType)
			continue
		}

		if !unicode.IsUpper(rune(name[0])) {
			continue
		}

		defaultVal := schema.GoTypeToTypescriptDefault(columns[columnNames[k]])
		fmt.Fprintf(sb, "\t"+columnNames[k]+": "+defaultVal+",\n")
	}

}

func ImportString(sb io.Writer, parts []string) {

	fmt.Fprintf(sb, "import { %s", parts[1])

	// If it starts as an array, do not include the import
	if parts[0][0:1] != "[" {
		fmt.Fprintf(sb, ", new%s", parts[1])
	}

	fmt.Fprintf(sb, " } from '%s';\n", parts[2])

}

func ImportStrings(sb io.Writer, columns map[string]string) {

	imported := map[string]struct{}{}

	// imports := [][]string{}
	for name := range columns {

		if !unicode.IsUpper(rune(name[0])) && name[0:1] != "#" {
			continue
		}

		dataType := columns[name]

		baseType := schema.ExtractBaseGoType(dataType)

		if !schema.IsGoTypeBaseType(baseType) {

			if len(baseType) > 10 && baseType[0:10] == "constants." {
				continue
			}

			if _, ok := imported[baseType]; ok {
				continue
			}

			imported[baseType] = struct{}{}

			if len(baseType) > 7 && baseType[0:7] == "models." {
				ImportString(sb, []string{dataType, baseType[7:], "gen/models/" + baseType[7:]})
			} else if len(baseType) > 5 && baseType[0:5] == "dtos." {
				ImportString(sb, []string{dataType, baseType[5:], "gen/dtos/" + baseType[5:]})
			} else {
				ImportString(sb, []string{dataType, baseType, "./" + baseType})
			}
		}
	}
}

func InheritStrings(sb io.Writer, columns map[string]string) []string {

	imports := []string{}
	for name := range columns {

		if len(name) > 9 && name[0:9] == "#embedded" {
			imports = append(imports, schema.GoTypeToTypescriptString(columns[name]))
		}
	}

	return imports
}

func ColumnMapToNames(columns map[string]string) []string {
	k := 0
	columnNames := make([]string, len(columns))
	for columnName := range columns {
		columnNames[k] = columnName
		k++
	}

	sort.Strings(columnNames)

	return columnNames
}

func TSFileHeader(sb io.Writer, name string) {
	fmt.Fprintf(sb, `/**
 * Generated Code; DO NOT EDIT
 *
 * `+name+`
 */
`)
}
