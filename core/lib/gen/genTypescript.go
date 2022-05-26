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

	fmt.Fprintf(sb, "import { %s", b)

	// If it starts as an array, do not include the import
	if a[0:1] != "[" {
		fmt.Fprintf(sb, ", new%s", b)
	}

	fmt.Fprintf(sb, " } from '%s';\n", c)

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
		}
	}
}

func ImportStrings2(columns map[string]string) string {

	var buf bytes.Buffer

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
	fmt.Fprintf(sb, `/**
 * Generated Code; DO NOT EDIT
 *
 * `+name+`
 */
`)
}
