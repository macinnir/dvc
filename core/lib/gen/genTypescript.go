package gen

import (
	"fmt"
	"io"
	"sort"
	"unicode"

	"github.com/macinnir/dvc/core/lib/schema"
)

// Additional helper functions for Typescript generation
// Don't require a constructor if it is an array, as we can just use the array literal syntax in that case
func requireConstructor(fullName string) bool {
	return fullName[0:1] != "["
}

func ImportString(sb io.Writer, fullName, objectName, importPath string, doIncludeConstructor bool) {

	fmt.Fprintf(sb, "import { %s", objectName)

	// If it starts as an array, do not include the import
	if doIncludeConstructor {
		fmt.Fprintf(sb, ", new%s", objectName)
	}

	fmt.Fprintf(sb, " } from '%s';\n", importPath)

}

type TypeToImport struct {
	FullName     string
	BaseTypeName string
	ObjectName   string
	ImportPath   string
	RequiresNew  bool
}

func getObjectNameAndImportPath(typeName string) (string, string) {
	if len(typeName) > 7 && typeName[0:7] == "models." {
		return typeName[7:], "gen/models/" + typeName[7:]
	} else if len(typeName) > 5 && typeName[0:5] == "dtos." {
		return typeName[5:], "gen/dtos/" + typeName[5:]
	} else if len(typeName) > 11 && typeName[0:11] == "aggregates." {
		return typeName[11:], "gen/aggregates/" + typeName[11:]
	} else {
		return typeName, "./" + typeName
	}
}

func NewTypeToImport(fullName string, baseType string) TypeToImport {

	t := TypeToImport{
		FullName:     fullName,
		BaseTypeName: baseType,
		ImportPath:   "",
		RequiresNew:  false,
	}

	t.ObjectName, t.ImportPath = getObjectNameAndImportPath(baseType)

	return t
}

func isImportable(typeName string) bool {
	return unicode.IsUpper(rune(typeName[0])) || typeName[0:1] == "#"
}

func isConstant(typeName string) bool {
	return len(typeName) > 10 && typeName[0:10] == "constants."
}

func ImportStrings(sb io.Writer, columns map[string]string) {

	// Build types to import
	// var n = 0
	var imported = map[string]TypeToImport{}
	var importNames = []string{}

	for name := range columns {

		if !isImportable(name) {
			continue
		}
		var fullName = columns[name]
		var baseType = schema.ExtractBaseGoType(fullName)

		if schema.IsGoTypeBaseType(baseType) {
			continue
		}

		if isConstant(baseType) {
			continue
		}

		var importObj TypeToImport
		if _, ok := imported[baseType]; ok {
			importObj = imported[baseType]
		} else {
			importObj = NewTypeToImport(fullName, baseType)
			importNames = append(importNames, baseType)
		}

		if requireConstructor(fullName) {
			importObj.RequiresNew = true
		}

		imported[baseType] = importObj
	}

	for _, name := range importNames {
		importType := imported[name]
		ImportString(sb, importType.FullName, importType.ObjectName, importType.ImportPath, importType.RequiresNew)
	}

	// for name := range columns {

	// 	n++
	// 	dataType := columns[name]

	// 	baseType := schema.ExtractBaseGoType(dataType)

	// 	if debug {
	// 		fmt.Printf("%d. ImportStrings checking column: %s(%s)\n", n, name, dataType)
	// 		fmt.Printf("\tbaseType: %s\n", baseType)
	// 	}
	// 	if !schema.IsGoTypeBaseType(baseType) {

	// 		if debug {
	// 			fmt.Printf("\tImporting baseType: %s\n", baseType)
	// 		}
	// 		ImportString(sb, fullName, importObj.ObjectName, importObj.ImportPath, debug)
	// 	}
	// }
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
