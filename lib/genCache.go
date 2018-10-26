package lib

import (
	"fmt"
	"os"
	"sort"
)

// GenerateGoCacheFile generates a repo file in golang
func (g *Gen) GenerateGoCacheFile(dir string, table *Table) (e error) {

	fileHead := ""
	fileFoot := ""
	goCode := ""
	imports := []string{}

	g.EnsureDir(dir)

	outFile := fmt.Sprintf("%s/%s.go", dir, table.Name)

	Debugf("Generating go cache file for table %s at path %s", g.Options, table.Name, outFile)

	if fileHead, fileFoot, imports, e = g.scanFileParts(outFile, true); e != nil {
		return
	}

	goCode, e = g.GenerateGoCache(table, fileFoot, imports)
	if e != nil {
		return
	}
	outFileContent := fileHead + goCode + fileFoot

	e = g.WriteGoCodeToFile(outFileContent, outFile)
	return
}

// GenerateGoCacheFiles generates go repository files based on the database schema
func (g *Gen) GenerateGoCacheFiles(reposDir string, database *Database) (e error) {

	var files []string

	files, e = FetchNonDirFileNames(reposDir)

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
			Infof("Removing repo %s", g.Options, file)
			os.Remove(fmt.Sprintf("./%s/%s", reposDir, file))
		}
	}

	fmt.Printf("CacheDir: %s %d\n", reposDir, len(files))

	for _, table := range database.Tables {

		Infof("Generating cache %s", g.Options, table.Name)
		e = g.GenerateGoCacheFile(reposDir, table)
		if e != nil {
			return
		}
	}

	e = g.GenerateReposBootstrapFile(reposDir, database)
	return
}

// GenerateGoCache returns a string for a repo in golang
func (g *Gen) GenerateGoCache(table *Table, fileFoot string, imports []string) (goCode string, e error) {

	primaryKey := ""
	primaryKeyType := ""

	// funcSig := fmt.Sprintf(`^func \(r \*%sCache\) [A-Z].*$`, table.Name)

	// footMatches := g.scanStringForFuncSignature(fileFoot, funcSig)

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
		"fmt",
		"goalgopher/models",
		"goalgopher/utils",
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
	goCode += "package cache\n\n"
	goCode += "import (\n"
	for _, i := range imports {
		goCode += "\t\"" + i + "\"\n"
	}
	goCode += ")\n\n"

	// Struct
	goCode += fmt.Sprintf("// %sCache is a data repository for %s objects\n",
		table.Name,
		table.Name,
	)

	goCode += fmt.Sprintf("type %sCache struct {\n", table.Name)
	goCode += fmt.Sprintf("\trepo *%sRepo\n", table.Name)
	goCode += fmt.Sprintf("\tcache utils.IStore\n")
	goCode += "}\n\n"

	// Create
	goCode += fmt.Sprintf("// Create creates a new %s entry in the database\n", table.Name)
	goCode += fmt.Sprintf("func (r *%sCache) Create(model *models.%s) (e error) {\n",
		table.Name,
		table.Name,
	)

	hasPrimaryKey := len(primaryKey) > 0 && idType == "int64"

	goCode += "\tif e = r.repo.Create(model); e != nil {\n"
	goCode += "\t\treturn\n"
	goCode += "\t}\n"

	if hasPrimaryKey {
		goCode += "\tr.cache.Set(fmt.Sprintf(\"" + table.Name + "_%d\", model." + primaryKey + "), model)\n"
	}

	goCode += "\treturn\n"
	goCode += "}\n\n"

	// Update
	goCode += fmt.Sprintf("// Update updates an existing %s entry in the database\n", table.Name)
	goCode += fmt.Sprintf("func (r *%sCache) Update(model *models.%s) (e error) {\n\n",
		table.Name,
		table.Name,
	)
	goCode += "\tif e = r.repo.Update(model); e != nil {\n"
	goCode += "\t\treturn\n"
	goCode += "\t}\n"

	if hasPrimaryKey {
		goCode += "\tr.cache.Set(fmt.Sprintf(\"" + table.Name + "_%d\", model." + primaryKey + "), model)\n"
	}

	goCode += "\treturn\n"
	goCode += "}\n\n"

	// Delete

	if isDeleted {

		goCode += fmt.Sprintf("// Delete marks an existing %s entry in the database as deleted\n", table.Name)
		goCode += fmt.Sprintf("func (r *%sRepo) Delete(model *models.%s) (e error) {\n\n",
			table.Name,
			table.Name,
		)

		goCode += "\tif e = r.repo.Delete(model); e != nil {\n"
		goCode += "\t\treturn\n"
		goCode += "\t}\n"
		goCode += "\tr.cache.Delete(fmt.Sprintf(\"" + table.Name + "_%d\", model." + primaryKey + "))\n"
		goCode += "\treturn\n"
		goCode += "}\n\n"
	}

	goCode += fmt.Sprintf("// HardDelete performs a SQL DELETE operation on a %s entry in the database\n", table.Name)
	goCode += fmt.Sprintf("func (r *%sRepo) HardDelete(model *models.%s) (e error) {\n\n",
		table.Name,
		table.Name,
	)

	goCode += "\tif e = r.repo.HardDelete(model); e != nil {\n"
	goCode += "\t\treturn\n"
	goCode += "\t}\n"
	goCode += "\tr.cache.Delete(fmt.Sprintf(\"" + table.Name + "_%d\", model." + primaryKey + "))\n"
	goCode += "\treturn\n"
	goCode += "}\n\n"

	// GetByID

	goCode += fmt.Sprintf("// GetByID gets a single %s object by a Primary Key\n", table.Name)
	goCode += fmt.Sprintf("func (r *%sRepo) GetByID(%s %s) (model *models.%s, e error) {\n",
		table.Name,
		primaryKey,
		idType,
		table.Name,
	)

	goCode += "\tmodel = &models.Account{}\n"
	goCode += "\tkey := fmt.Sprintf(\"" + table.Name + "_%d\", " + primaryKey + ")\n"
	goCode += "\tif e = r.cache.Get(key, model); e != nil {\n"
	goCode += "\t\t// redis.Nil -- no rows\n"
	goCode += "\t\tmodel, e = r.repo.GetByID(AccountID)\n"
	goCode += "\t\tr.cache.Set(key, model)\n"
	goCode += "\t}\n"
	goCode += "\treturn\n"
	goCode += "}\n"

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
	collection, e = r.repo.GetMany(limit, offset, args...)
	return
	`
	goCode += "\n}\n\n"

	// Single
	goCode += fmt.Sprintf("// GetSingle gets one %s object\n",
		table.Name,
	)

	goCode += fmt.Sprintf("func (r *%sRepo) GetSingle(args ...string) (model *models.%s, e error) {\n",
		table.Name,
		table.Name,
	)

	goCode += `
	model, e = r.repo.GetSingle(args...)
	return
	`
	goCode += "\n}\n\n"

	goCode += "// #genEnd\n"

	return
}
