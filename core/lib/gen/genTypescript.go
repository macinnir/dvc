package gen

import (
	"strings"

	"github.com/macinnir/dvc/core/lib/schema"
)

// GenerateTypescriptTypes returns a string for a typscript types file
func GenerateTypescriptTypes(database *schema.Schema) (string, error) {

	var e error
	var sb strings.Builder

	sb.WriteString("// #genStart \n\ndeclare namespace Models {\n\n")
	for k := range database.Tables {

		table := database.Tables[k]

		var str string
		if str, e = GenerateTypescriptType(table); e != nil {
			return "", e
		}

		sb.WriteString(str)
	}

	sb.WriteString("}\n")
	sb.WriteString("// #genEnd\n")

	return sb.String(), nil
}

// GenerateTypescriptType returns a string for a type in typescript
func GenerateTypescriptType(table *schema.Table) (string, error) {

	var sb strings.Builder

	sb.WriteString("\t/**\n\t * " + table.Name + "\n\t */\n")
	sb.WriteString("\texport interface " + table.Name + "{\n")
	for _, column := range table.Columns {

		fieldType := schema.DataTypeToTypescriptString(column)

		sb.WriteString("\t\t" + column.Name + ": " + fieldType + ";\n")
	}

	sb.WriteString("\t}\n\n")

	return sb.String(), nil
}
