package lib

import (
	"fmt"
	"sort"
)

// GenerateGoSchemaFile generates a schema file in golang
func (g *Gen) GenerateGoSchemaFile(dir string, database *Database) (e error) {

	g.EnsureDir(dir)

	var fileHead, fileFoot, goCode string

	outFile := fmt.Sprintf("%s/schema.go", dir)

	Debugf("Generating go schema file at path %s", g.Options, outFile)

	if fileHead, fileFoot, _, e = g.scanFileParts(outFile, false); e != nil {
		return
	}

	goCode, e = g.GenerateGoSchema(database)

	if e != nil {
		return
	}

	// Add package statement
	if !g.fileExists(outFile) {
		fileHead = "package schema\n\n"
	}

	outFileContent := fileHead + goCode + fileFoot

	e = g.WriteGoCodeToFile(outFileContent, outFile)
	return
}

// GenerateGoSchema generates golang code for a schema file
func (g *Gen) GenerateGoSchema(database *Database) (goCode string, e error) {

	goCode = "\n// #genStart"
	goCode += "\n// Schema defines the data access layer schema"
	goCode += "\ntype Schema struct {"
	goCode += "\n\tSchema dal.ISchema"
	goCode += "\n}"
	goCode += "\n// Init initializes the DAL Schema"
	goCode += "\nfunc (s *Schema) Init() {"

	for _, table := range database.Tables {
		cols := ""

		sortedColumns := make(SortedColumns, 0, len(table.Columns))

		for _, column := range table.Columns {
			sortedColumns = append(sortedColumns, column)
		}

		sort.Sort(sortedColumns)

		for _, column := range sortedColumns {
			cols += fmt.Sprintf("\t\t\t\"%s\",\n", column.Name)
		}

		goCode += fmt.Sprintf("\n\n\t// %s", table.Name)
		goCode += fmt.Sprintf("\n\ts.Schema.AddTable(")
		goCode += fmt.Sprintf("\n\t\t\"%s\",", table.Name)
		goCode += fmt.Sprintf("\n\t\t[]string{\n%s\t\t})", cols)
	}

	goCode += "\n}\n\n// #genEnd\n"

	return
}
