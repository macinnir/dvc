package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

//
// Commands
//

// GenerateReposBootstrapFile generates a repos bootstrap file in golang
func GenerateReposBootstrapFile(database *Database) (e error) {

	outFile := "./repos/repos.go"
	goCode, e := GenerateReposBootstrapGoCodeFromDatabase(database)

	if e != nil {
		return
	}

	e = WriteGoCodeToFile(goCode, outFile)

	return
}

// GenerateGoRepoFile generates a repo file in golang
func GenerateGoRepoFile(table *Table) (e error) {

	fileHead := ""
	fileFoot := ""
	goCode := ""
	imports := []string{}

	outFile := fmt.Sprintf("./repos/%s.go", table.Name)

	if fileHead, fileFoot, imports, e = scanFileParts(outFile, true); e != nil {
		return
	}

	goCode, e = GenerateGoRepo(table, fileFoot, imports)
	if e != nil {
		return
	}
	outFileContent := fileHead + goCode + fileFoot

	e = WriteGoCodeToFile(outFileContent, outFile)
	return
}

// GenerateGoSchemaFile generates a schema file in golang
func GenerateGoSchemaFile(database *Database) (e error) {

	var fileHead, fileFoot, goCode string

	outFile := fmt.Sprintf("./schema/schema.go")

	if fileHead, fileFoot, _, e = scanFileParts(outFile, false); e != nil {
		return
	}

	goCode, e = GenerateGoSchema(database)
	if e != nil {
		return
	}

	outFileContent := fileHead + goCode + fileFoot

	e = WriteGoCodeToFile(outFileContent, outFile)
	return
}

// GenerateGoModelFile generates a model file in golang
func GenerateGoModelFile(table *Table) (e error) {

	var fileHead, fileFoot, goCode string
	var imports []string

	outFile := fmt.Sprintf("./models/%s.go", table.Name)

	if fileHead, fileFoot, imports, e = scanFileParts(outFile, false); e != nil {
		return
	}

	goCode, e = GenerateGoModel(table, imports)
	if e != nil {
		return
	}

	goCode = fileHead + goCode + fileFoot

	e = WriteGoCodeToFile(goCode, outFile)
	return
}

// GenerateTypescriptTypesFile generates a typescript type file
func GenerateTypescriptTypesFile(database *Database) (e error) {

	var goCode string

	outFile := "./src/types/types.d.ts"

	goCode, e = GenerateTypescriptTypes(database)
	if e != nil {
		return
	}

	ioutil.WriteFile(outFile, []byte(goCode), 0644)

	return

}

//
// String Generators
//

// GenerateGoModels generates models for golang
func GenerateGoModels(database *Database) (goCode string, e error) {

	goCode = "// #genStart \n\n"

	imports := []string{}

	for _, table := range database.Tables {
		code := ""
		if code, e = GenerateGoModel(table, imports); e != nil {
			return
		}

		goCode += code
	}

	goCode += "// #genEnd\n"

	return
}

// GenerateReposBootstrapGoCodeFromDatabase generates golang code for a Repo Bootstrap file from
// a database object
func GenerateReposBootstrapGoCodeFromDatabase(database *Database) (goCode string, e error) {

	props := ""
	defs := ""

	for _, table := range database.Tables {
		props += fmt.Sprintf("\t%s I%sRepo\n",
			table.Name,
			table.Name,
		)

		defs += fmt.Sprintf("\t// %s\n", table.Name)
		defs += fmt.Sprintf("\trepo%s := new(%sRepo)\n", table.Name, table.Name)
		defs += fmt.Sprintf("\trepo%s.Dal = schema\n", table.Name)
		defs += fmt.Sprintf("\tr.%s = repo%s\n\n", table.Name, table.Name)
	}

	goCode = "package repos"
	goCode += "\n\nimport("
	goCode += "\n\tdal \"github.com/macinnir/go-dal\""
	goCode += "\n)"
	goCode += "\n\n// Repos is a collection of repositories"
	goCode += "\ntype Repos struct {"
	goCode += "\n\tBase\tdal.ISchema"
	goCode += fmt.Sprintf("\n%s", props)
	goCode += "\n}"
	goCode += "\n\n//Bootstrap bootstraps all of the repositories into a single repository object"
	goCode += "\nfunc Bootstrap(schema dal.ISchema) *Repos {"
	goCode += "\n\n\tr := new(Repos)"
	goCode += "\n\n\tr.Base = schema"
	goCode += "\n\n\t// Repos"
	goCode += fmt.Sprintf("\n\n%s", defs)
	goCode += "\n\n\treturn r"
	goCode += "\n}"

	return
}

// GenerateGoSchema generates golang code for a schema file
func GenerateGoSchema(database *Database) (goCode string, e error) {

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

		goCode += fmt.Sprintf("\n\n//%s", table.Name)
		goCode += fmt.Sprintf("\n\ts.Schema.AddTable(")
		goCode += fmt.Sprintf("\n\t\t\"%s\",", table.Name)
		goCode += fmt.Sprintf("[]string{%s\n\t})", cols)
	}

	goCode += "\n}\n\n// #genEnd\n"

	return
}

// GenerateGoRepo returns a string for a repo in golang
func GenerateGoRepo(table *Table, fileFoot string, imports []string) (goCode string, e error) {

	primaryKey := ""
	primaryKeyType := ""

	funcSig := fmt.Sprintf(`^func \(r \*%sRepo\) [A-Z].*$`, table.Name)

	footMatches := scanStringForFuncSignature(fileFoot, funcSig)

	sortedColumns := make(SortedColumns, 0, len(table.Columns))

	// Find the primary key
	for _, column := range table.Columns {
		if column.ColumnKey == "PRI" {
			primaryKey = column.Name
			primaryKeyType = column.DataType
		}

		sortedColumns = append(sortedColumns, column)
	}

	sort.Sort(sortedColumns)

	_, isDeleted := table.Columns["IsDeleted"]

	idType := "int64"
	switch primaryKeyType {
	case "varchar":
		idType = "string"
	}

	defaultImports := []string{
		"goalgopher/models",
		"database/sql",
		"github.com/macinnir/go-dal",
		"errors",
	}

	if len(imports) > 0 {

		for _, di := range defaultImports {

			exists := false

			for _, ii := range imports {
				if ii == di {
					exists = true
					break
				}
			}

			if !exists {
				imports = append(imports, di)
			}

		}

	} else {
		imports = defaultImports
	}

	goCode += "// #genStart\n\n"
	goCode += "package repos\n\n"
	goCode += "import (\n"
	for _, i := range imports {
		goCode += "\t\"" + i + "\"\n"
	}
	goCode += ")\n\n"
	goCode += fmt.Sprintf("// I%sRepo outlines the repository methods for %s objects\n", table.Name, table.Name)
	goCode += fmt.Sprintf("type I%sRepo interface {\n", table.Name)
	goCode += fmt.Sprintf("\tCreate(model *models.%s) (e error)\n", table.Name)
	goCode += fmt.Sprintf("\tUpdate(model *models.%s) (e error)\n", table.Name)

	if isDeleted {
		goCode += fmt.Sprintf("\tDelete(model *models.%s) (e error)\n", table.Name)
	}

	goCode += fmt.Sprintf("\tHardDelete(model *models.%s) (e error)\n", table.Name)
	goCode += fmt.Sprintf("\tGetMany(limit int, offset int, args ...string)(collection []*models.%s, e error)\n", table.Name)
	goCode += fmt.Sprintf("\tGetByID(%s %s)(model *models.%s, e error)\n", primaryKey, idType, table.Name)
	goCode += fmt.Sprintf("\tGetSingle(args ...string)(model *models.%s, e error)\n", table.Name)

	if len(footMatches) > 0 {
		footMatchPrefixLen := len(fmt.Sprintf("func (r *%sRepo) ", table.Name))
		for _, footMatch := range footMatches {
			footMatch = strings.Trim(footMatch, " ")
			footMatchLen := len(footMatch)
			goCode += "\t" + footMatch[footMatchPrefixLen:footMatchLen-1] + "\n"
		}
	}

	goCode += "}\n\n"

	// Struct
	goCode += fmt.Sprintf("// %sRepo is a data repository for %s objects\n",
		table.Name,
		table.Name,
	)

	goCode += fmt.Sprintf("type %sRepo struct {\n", table.Name)
	goCode += "\tDal dal.ISchema\n"
	goCode += "}\n\n"

	// Create

	goCode += fmt.Sprintf("// Create creates a new %s entry in the database\n", table.Name)

	goCode += fmt.Sprintf("func (r *%sRepo) Create(model *models.%s) (e error) {\n",
		table.Name,
		table.Name,
	)

	hasPrimaryKey := len(primaryKey) > 0 && idType == "int64"

	if hasPrimaryKey {
		goCode += "\n\tvar result sql.Result\n\n"
	}

	goCode += fmt.Sprintf("\tq := r.Dal.Insert(\"%s\")\n", table.Name)

	for _, column := range sortedColumns {
		if column.ColumnKey == "PRI" {
			continue
		}
		goCode += fmt.Sprintf("\tq.Set(\"%s\", model.%s)\n", column.Name, column.Name)
	}

	if hasPrimaryKey {
		goCode += fmt.Sprintf("\n\tresult, e = q.Exec()\n")
	} else {
		goCode += fmt.Sprintf("\n\t_, e = q.Exec()\n")
	}
	goCode += fmt.Sprintf("\n\tif e != nil {\n")
	goCode += "\t\treturn\n"
	goCode += "\t}\n"

	if hasPrimaryKey {
		goCode += fmt.Sprintf("\n\tmodel.%s, e = result.LastInsertId()\n", primaryKey)
	}

	goCode += "\n\treturn\n"

	goCode += "}\n\n"

	// Update
	goCode += fmt.Sprintf("// Update updates an existing %s entry in the database\n", table.Name)
	goCode += fmt.Sprintf("func (r *%sRepo) Update(model *models.%s) (e error) {\n\n",
		table.Name,
		table.Name,
	)

	goCode += fmt.Sprintf("\tq := r.Dal.Update(\"%s\")\n", table.Name)
	for _, column := range sortedColumns {
		if column.ColumnKey == "PRI" {
			primaryKey = column.Name
			continue
		}
		goCode += fmt.Sprintf("\tq.Set(\"%s\", model.%s)\n", column.Name, column.Name)
	}

	goCode += fmt.Sprintf("\n\tq.Where(\"%s\", model.%s)\n", primaryKey, primaryKey)

	goCode += fmt.Sprintf("\n\t_, e = q.Exec()\n")
	goCode += fmt.Sprintf("\n\treturn\n")

	goCode += "}\n\n"

	// Delete

	if isDeleted {

		goCode += fmt.Sprintf("// Delete marks an existing %s entry in the database as deleted\n", table.Name)
		goCode += fmt.Sprintf("func (r *%sRepo) Delete(model *models.%s) (e error) {\n\n",
			table.Name,
			table.Name,
		)

		goCode += fmt.Sprintf("\tq := r.Dal.Update(\"%s\")\n", table.Name)
		goCode += "\tq.Set(\"IsDeleted\", 1)\n"
		goCode += fmt.Sprintf("\tq.Where(\"%s\", model.%s)", primaryKey, primaryKey)

		goCode += fmt.Sprintf("\n\t_, e = q.Exec()\n")
		goCode += fmt.Sprintf("\n\treturn\n")

		goCode += "}\n\n"
	}

	goCode += fmt.Sprintf("// HardDelete performs a SQL DELETE operation on a %s entry in the database\n", table.Name)
	goCode += fmt.Sprintf("func (r *%sRepo) HardDelete(model *models.%s) (e error) {\n\n",
		table.Name,
		table.Name,
	)

	goCode += fmt.Sprintf("\tq := r.Dal.Delete(\"%s\")\n", table.Name)
	goCode += fmt.Sprintf("\tq.Where(\"%s\", model.%s)", primaryKey, primaryKey)

	goCode += fmt.Sprintf("\n\t_, e = q.Exec()\n")
	goCode += fmt.Sprintf("\n\treturn\n")

	goCode += "}\n\n"

	// SelectByID

	goCode += fmt.Sprintf("// GetByID gets a single %s object by a Primary Key\n", table.Name)
	goCode += fmt.Sprintf("func (r *%sRepo) GetByID(%s %s) (model *models.%s, e error) {\n",
		table.Name,
		primaryKey,
		idType,
		table.Name,
	)
	goCode += "\n\tvar rows *sql.Rows\n\n"
	goCode += fmt.Sprintf("\tq := r.Dal.Select(\"%s\")\n", table.Name)
	goCode += fmt.Sprintf("\tq.Where(\"%s\", %s)\n", primaryKey, primaryKey)
	goCode += `
	rows, e = q.Query()

	if rows != nil {
		defer rows.Close()
	}

	if e != nil {
		return
	}

	rows.Next()
	model, e = r.scanRow(rows)
	return`

	goCode += "\n}\n\n"

	// Select

	goCode += fmt.Sprintf("// GetMany gets %s objects\n",
		table.Name,
	)

	goCode += fmt.Sprintf("func (r *%sRepo) GetMany(limit int, offset int, args ...string) (collection []*models.%s, e error) {\n",
		table.Name,
		table.Name,
	)

	goCode += "\n\tvar rows *sql.Rows\n\n"
	goCode += fmt.Sprintf("\tq := r.Dal.Select(\"%s\")\n", table.Name)

	goCode += `
	argLen := len(args) 
	idx := 0
	idx1 := 1

	for argLen > 0 {
		q.Where(args[idx], args[idx1])
		argLen = argLen - 2
		idx = idx + 2
		idx1 = idx1 + 2
	}
	
	`
	goCode += "\tq.Limit(limit)\n"
	goCode += "\tq.Offset(offset)\n"
	goCode += fmt.Sprintf(`
	rows, e = q.Query()

	if rows != nil {
		defer rows.Close()
	}

	if e != nil {
		return
	}

	collection = []*models.%s{}

	for rows.Next() {

		var model *models.%s

		model, e = r.scanRow(rows)

		if e != nil {
			return 
		}

		collection = append(collection, model)
	}

	return`, table.Name, table.Name)

	goCode += "\n}\n\n"

	// Single
	goCode += fmt.Sprintf("// GetSingle gets one %s object\n",
		table.Name,
	)

	goCode += fmt.Sprintf("func (r *%sRepo) GetSingle(args ...string) (model *models.%s, e error) {\n",
		table.Name,
		table.Name,
	)

	goCode += fmt.Sprintf(`

	var collection []*models.%s
	if collection, e = r.GetMany(1, 0, args...); e != nil {
		return 
	}

	if len(collection) == 0 {
		e = errors.New("No rows")
		return 
	}

	model = collection[0]
	return 

	`, table.Name)

	goCode += "\n}\n\n"

	// Scan
	goCode += fmt.Sprintf("// scanRow scans all of the rows for %s models \n",
		table.Name,
	)
	goCode += fmt.Sprintf("func (r *%sRepo) scanRow(row *sql.Rows) (model *models.%s, e error) {\n\n",
		table.Name,
		table.Name,
	)

	goCode += fmt.Sprintf("\tmodel = new(models.%s)\n\n", table.Name)

	goCode += fmt.Sprintf("\te = row.Scan(\n")

	for _, column := range sortedColumns {
		goCode += fmt.Sprintf("\t\t&model.%s,\n", column.Name)
	}
	goCode += fmt.Sprintf("\t)\n\n")

	goCode += "\treturn\n"
	goCode += "}\n\n"
	goCode += "// #genEnd\n"

	return
}

// GenerateGoModel returns a string for a model in golang
func GenerateGoModel(table *Table, imports []string) (goCode string, e error) {

	goCode += "// #genStart\n\n"
	goCode += "package models\n\n"

	var sortedColumns = make(SortedColumns, 0, len(table.Columns))

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

// GenerateTypescriptTypes returns a string for a typscript types file
func GenerateTypescriptTypes(database *Database) (goCode string, e error) {
	goCode = "// #genStart \n\n"
	for _, table := range database.Tables {

		str := ""

		if str, e = GenerateTypescriptType(table); e != nil {
			return
		}

		goCode += str
	}

	goCode += "// #genEnd\n"

	return
}

// GenerateTypescriptType returns a string for a type in typescript
func GenerateTypescriptType(table *Table) (goCode string, e error) {

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
