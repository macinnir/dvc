package gen

import (
	"fmt"
	"github.com/macinnir/dvc/lib"
	"io/ioutil"
)

// GenerateTypescriptTypes returns a string for a typscript types file
func (g *Gen) GenerateTypescriptTypes(database *lib.Database) (goCode string, e error) {
	goCode = "// #genStart \n\n"
	for _, table := range database.Tables {

		str := ""

		if str, e = g.GenerateTypescriptType(table); e != nil {
			return
		}

		goCode += str
	}

	goCode += "// #genEnd\n"

	return
}

// GenerateTypescriptType returns a string for a type in typescript
func (g *Gen) GenerateTypescriptType(table *lib.Table) (goCode string, e error) {

	goCode += fmt.Sprintf("/**\n * %s\n */\n", table.Name)
	goCode += fmt.Sprintf("declare interface %s {\n", table.Name)
	for _, column := range table.Columns {

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

		goCode += fmt.Sprintf("\t%s: %s;\n", column.Name, fieldType)
	}

	goCode += "}\n\n"

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
