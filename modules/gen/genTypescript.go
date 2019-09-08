package gen

import (
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/macinnir/dvc/lib"
)

// GenerateTypescriptTypes returns a string for a typscript types file
func (g *Gen) GenerateTypescriptTypes(database *lib.Database) (goCode string, e error) {
	goCode = "// #genStart \n\ndeclare namespace Models {\n\n"

	tableKeys := []string{}
	for key := range database.Tables {
		tableKeys = append(tableKeys, key)
	}

	sort.Strings(tableKeys)

	for _, key := range tableKeys {

		table := database.Tables[key]

		str := ""

		if str, e = g.GenerateTypescriptType(table); e != nil {
			return
		}

		goCode += str
	}

	goCode += "}\n"
	goCode += "// #genEnd\n"

	return
}

// GenerateTypescriptType returns a string for a type in typescript
func (g *Gen) GenerateTypescriptType(table *lib.Table) (goCode string, e error) {

	goCode += fmt.Sprintf("\t/**\n\t * %s\n\t */\n", table.Name)
	goCode += fmt.Sprintf("\texport interface %s {\n", table.Name)

	columnKeys := []string{}

	for key := range table.Columns {
		columnKeys = append(columnKeys, key)
	}

	sort.Strings(columnKeys)

	for _, key := range columnKeys {

		column := table.Columns[key]

		fieldType := "number"
		switch column.DataType {
		case "varchar":
			fieldType = "string"
		case "enum":
			fieldType = "string"
		case "text":
			fieldType = "string"
		case "date":
			fieldType = "string"
		case "datetime":
			fieldType = "string"
		case "char":
			fieldType = "string"
		}
		// decimal
		//

		goCode += fmt.Sprintf("\t\t%s: %s;\n", column.Name, fieldType)
	}

	goCode += "\t}\n\n"

	return

}

// GenerateTypescriptTypesFile generates a typescript type file
func (g *Gen) GenerateTypescriptTypesFile(dir string, database *lib.Database) (e error) {

	g.EnsureDir(dir)

	var goCode string

	outFile := fmt.Sprintf("%s/types.d.ts", dir)
	lib.Debugf("Generating typescript types file at path %s", g.Options, outFile)
	goCode, e = g.GenerateTypescriptTypes(database)
	if e != nil {
		return
	}

	ioutil.WriteFile(outFile, []byte(goCode), 0644)

	return

}
