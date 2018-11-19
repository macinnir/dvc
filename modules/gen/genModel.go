package gen

import (
	"fmt"
	"github.com/macinnir/dvc/lib"
	"sort"
)

// GenerateGoModel returns a string for a model in golang
func (g *Gen) GenerateGoModel(table *lib.Table, imports []string) (goCode string, e error) {

	goCode += "// #genStart\n\n"
	goCode += "package models\n\n"

	var sortedColumns = make(lib.SortedColumns, 0, len(table.Columns))

	for _, column := range table.Columns {
		sortedColumns = append(sortedColumns, column)
	}

	sort.Sort(sortedColumns)

	includeNullPackage := false

	fieldCode := ""

	for _, column := range sortedColumns {

		fieldType := "int64"
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
		case "decimal":
			fieldType = "float64"
		}

		if column.IsNullable == true {
			includeNullPackage = true
			switch fieldType {
			case "string":
				// fieldType = "sql.NullString"
				fieldType = "null.String"
			case "int64":
				// fieldType = "sql.NullInt64"
				fieldType = "null.Int"
			case "float64":
				// fieldType = "sql.NullFloat64"
				fieldType = "null.Float"
			}
		}

		fieldCode += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", column.Name, fieldType, column.Name)
	}

	if includeNullPackage == true {
		// goCode += "import \"database/sql\"\n"
		goCode += "import \"gopkg.in/guregu/null.v3\"\n\n"
	}

	goCode += fmt.Sprintf("// %s represents a %s model\n", table.Name, table.Name)
	goCode += fmt.Sprintf("type %s struct {\n", table.Name)
	goCode += fieldCode
	goCode += "}\n\n"
	goCode += "// #genEnd"
	return
}

// GenerateGoModels generates models for golang
func (g *Gen) GenerateGoModels(database *lib.Database) (goCode string, e error) {

	goCode = "// #genStart \n\n"

	imports := []string{}

	for _, table := range database.Tables {
		code := ""
		if code, e = g.GenerateGoModel(table, imports); e != nil {
			return
		}

		goCode += code
	}

	goCode += "// #genEnd\n"

	return
}

// GenerateGoModelFile generates a model file in golang
func (g *Gen) GenerateGoModelFile(dir string, table *lib.Table) (e error) {

	g.EnsureDir(dir)

	var fileHead, fileFoot, goCode string
	var imports []string

	outFile := fmt.Sprintf("%s/%s.go", dir, table.Name)
	lib.Debugf("Generating model for table %s at path %s", g.Options, table.Name, outFile)

	if fileHead, fileFoot, imports, e = g.scanFileParts(outFile, false); e != nil {
		return
	}

	goCode, e = g.GenerateGoModel(table, imports)
	if e != nil {
		return
	}

	goCode = fileHead + goCode + fileFoot

	e = g.WriteGoCodeToFile(goCode, outFile)
	return
}
