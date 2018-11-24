package gen

import (
	"fmt"
	"github.com/macinnir/dvc/lib"
	"os"
	"sort"
	"strings"
)

// GenerateGoRepoFile generates a repo file in golang
func (g *Gen) GenerateGoRepoFile(dir string, table *lib.Table) (e error) {

	fileHead := ""
	fileFoot := ""
	goCode := ""
	imports := []string{}

	g.EnsureDir(dir)

	outFile := fmt.Sprintf("%s/%s.go", dir, table.Name)

	lib.Debugf("Generating go repo file for table %s at path %s", g.Options, table.Name, outFile)

	if fileHead, fileFoot, imports, e = g.scanFileParts(outFile, true); e != nil {
		return
	}

	goCode, e = g.GenerateGoRepo(table, fileFoot, imports)
	if e != nil {
		return
	}
	outFileContent := fileHead + goCode + fileFoot

	e = g.WriteGoCodeToFile(outFileContent, outFile)
	return
}

// GenerateGoRepoFiles generates go repository files based on the database schema
func (g *Gen) GenerateGoRepoFiles(reposDir string, database *lib.Database) (e error) {

	var files []string

	files, e = lib.FetchNonDirFileNames(reposDir)

	// clean out files that don't belong
	for _, file := range files {

		if file == "repos.go" {
			continue
		}

		existsInDatabase := false

		for _, table := range database.Tables {
			if file == table.Name+".go" {
				existsInDatabase = true
				break
			}
		}

		if !existsInDatabase {
			lib.Infof("Removing repo %s", g.Options, file)
			os.Remove(fmt.Sprintf("./%s/%s", reposDir, file))
		}
	}

	for _, table := range database.Tables {

		lib.Debugf("Generating repo %s", g.Options, table.Name)
		e = g.GenerateGoRepoFile(reposDir, table)
		if e != nil {
			return
		}
	}

	e = g.GenerateReposBootstrapFile(reposDir, database)
	return
}

// GenerateGoRepo returns a string for a repo in golang
func (g *Gen) GenerateGoRepo(table *lib.Table, fileFoot string, imports []string) (goCode string, e error) {

	primaryKey := ""
	primaryKeyType := ""

	funcSig := fmt.Sprintf(`^func \(r \*%sRepo\) [A-Z].*$`, table.Name)

	footMatches := g.scanStringForFuncSignature(fileFoot, funcSig)

	sortedColumns := make(lib.SortedColumns, 0, len(table.Columns))

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
		fmt.Sprintf("%s/models", g.Config.BasePackage),
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

// GenerateReposBootstrapGoCodeFromDatabase generates golang code for a Repo Bootstrap file from
// a database object
func (g *Gen) GenerateReposBootstrapGoCodeFromDatabase(database *lib.Database) (goCode string, e error) {

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

// GenerateReposBootstrapFile generates a repos bootstrap file in golang
func (g *Gen) GenerateReposBootstrapFile(dir string, database *lib.Database) (e error) {

	// Make the repos dir if it does not exist.
	g.EnsureDir(dir)

	outFile := fmt.Sprintf("%s/repos.go", dir)
	goCode, e := g.GenerateReposBootstrapGoCodeFromDatabase(database)
	lib.Debugf("Generating go Repos bootstrap file at path %s", g.Options, outFile)
	if e != nil {
		return
	}

	e = g.WriteGoCodeToFile(goCode, outFile)

	return
}
