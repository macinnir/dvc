package gen

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"unicode"

	"github.com/macinnir/dvc/core/lib/schema"
)

func ImportString(sb io.Writer, a, b, c string) {

	fmt.Fprint(sb, "import { "+b)

	// If it starts as an array, do not include the import
	if a[0:1] != "[" {
		fmt.Fprint(sb, ", new"+b)
	}

	fmt.Fprint(sb, " } from '"+c+"';\n")

}

func ImportStrings(sb io.Writer, columns map[string]string, typeName string) {

	imported := map[string]struct{}{}

	var sortedCols = make([]string, 0, len(columns))
	for name := range columns {
		sortedCols = append(sortedCols, name)
	}

	sort.Strings(sortedCols)

	// imports := [][]string{}
	for _, name := range sortedCols {

		// Golang properties are only public if they start with an uppercase letter
		// Embedded (inherited) data types start with #embedded (e.g. [#embedded1] => "*models.User")
		if !unicode.IsUpper(rune(name[0])) && name[0:1] != "#" {
			continue
		}

		dataType := columns[name]
		baseType := schema.ExtractBaseGoType(dataType)

		if typeName == "QuestionAggregate" {
			fmt.Println("Column " + name + " Data Type " + dataType + " Base Type " + baseType)
		}

		if !schema.IsGoTypeBaseType(baseType) {

			if typeName == "QuestionAggregate" {
				fmt.Println("Not BaseType " + baseType)
			}

			// Constants are not imported
			if len(baseType) > 10 && baseType[0:10] == "constants." {
				continue
			}

			// Already imported
			if _, ok := imported[baseType]; ok {
				continue
			}

			imported[baseType] = struct{}{}

			if len(baseType) > 7 && baseType[0:7] == "models." {
				ImportString(sb, dataType, baseType[7:], "gen/models/"+baseType[7:])
			} else if len(baseType) > 5 && baseType[0:5] == "dtos." {
				ImportString(sb, dataType, baseType[5:], "gen/dtos/"+baseType[5:])
			} else {
				ImportString(sb, dataType, baseType, "./"+baseType)
			}
		} else {
			if typeName == "QuestionAggregate" {
				fmt.Println("BaseType " + baseType)
			}
		}
	}
}

func ImportStrings2(columns map[string]string) string {

	var buf bytes.Buffer

	imported := map[string]struct{}{}

	// imports := [][]string{}
	for name := range columns {

		// Golang properties are only public if they start with an uppercase letter
		// Embedded (inherited) properties start with #embedded
		if !unicode.IsUpper(rune(name[0])) && name[0:1] != "#" {
			continue
		}

		dataType := columns[name]

		baseType := schema.ExtractBaseGoType(dataType)

		if !schema.IsGoTypeBaseType(baseType) {

			if len(baseType) > 10 && baseType[0:10] == "constants." {
				continue
			}

			// Already imported
			if _, ok := imported[baseType]; ok {
				continue
			}

			imported[baseType] = struct{}{}

			if len(baseType) > 7 && baseType[0:7] == "models." {
				ImportString(&buf, dataType, baseType[7:], "gen/models/"+baseType[7:])
			} else if len(baseType) > 5 && baseType[0:5] == "dtos." {
				ImportString(&buf, dataType, baseType[5:], "gen/dtos/"+baseType[5:])
			} else {
				ImportString(&buf, dataType, baseType, "./"+baseType)
			}
		}
	}

	return buf.String()
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
	fmt.Fprint(sb, `/**
 * Generated Code; DO NOT EDIT
 *
 * `+name+`
 */
`)
}
