package gen

import (
	"fmt"
	"github.com/macinnir/dvc/lib"
	"os"
	"sort"
	"strings"
)

// GenerateGoServiceFile generates a repo file in golang
func (g *Gen) GenerateGoServiceFile(dir string, table *lib.Table) (e error) {

	fileHead := ""
	fileFoot := ""
	goCode := ""
	imports := []string{}

	g.EnsureDir(dir)

	outFile := fmt.Sprintf("%s/%s.go", dir, table.Name)

	lib.Debugf("Generating go service file for table %s at path %s", g.Options, table.Name, outFile)

	if fileHead, fileFoot, imports, e = g.scanFileParts(outFile, true); e != nil {
		return
	}

	goCode, e = g.GenerateGoService(table, fileFoot, imports)
	if e != nil {
		return
	}
	outFileContent := fileHead + goCode + fileFoot

	e = g.WriteGoCodeToFile(outFileContent, outFile)
	return
}

// GenerateGoServiceFiles generates go repository files based on the database schema
func (g *Gen) GenerateGoServiceFiles(reposDir string, database *lib.Database) (e error) {

	var files []string

	files, e = lib.FetchNonDirFileNames(reposDir)

	// clean out files that don't belong
	for _, file := range files {

		if file == "services.go" {
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
			lib.Infof("Removing service %s", g.Options, file)
			os.Remove(fmt.Sprintf("./%s/%s", reposDir, file))
		}
	}

	for _, table := range database.Tables {

		lib.Debugf("Generating service %s", g.Options, table.Name)
		e = g.GenerateGoServiceFile(reposDir, table)
		if e != nil {
			return
		}
	}

	e = g.GenerateServicesBootstrapFile(reposDir, database)
	return
}

// GenerateGoService returns a string for a repo in golang
func (g *Gen) GenerateGoService(table *lib.Table, fileFoot string, imports []string) (goCode string, e error) {

	primaryKey := ""
	primaryKeyType := ""

	funcSig := fmt.Sprintf(`^func \(r \*%sService\) [A-Z].*$`, table.Name)

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
		fmt.Sprintf("%s/repos", g.Config.BasePackage),
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
	goCode += "package services\n\n"
	goCode += "import (\n"
	for _, i := range imports {
		goCode += "\t\"" + i + "\"\n"
	}
	goCode += ")\n\n"
	goCode += fmt.Sprintf("// I%sService outlines the service methods for %s objects\n", table.Name, table.Name)
	goCode += fmt.Sprintf("type I%sService interface {\n", table.Name)
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
		footMatchPrefixLen := len(fmt.Sprintf("func (r *%sService) ", table.Name))
		for _, footMatch := range footMatches {
			footMatch = strings.Trim(footMatch, " ")
			footMatchLen := len(footMatch)
			goCode += "\t" + footMatch[footMatchPrefixLen:footMatchLen-1] + "\n"
		}
	}

	goCode += "}\n\n"

	// Struct
	goCode += fmt.Sprintf("// %sService is a service for %s objects\n",
		table.Name,
		table.Name,
	)

	goCode += fmt.Sprintf("type %sService struct {\n", table.Name)
	goCode += "\tConfig *models.Config\n"
	goCode += "\tRepos *repos.Repos\n"
	goCode += "}\n\n"

	// Create

	goCode += fmt.Sprintf("// Create creates a new %s entry\n", table.Name)

	goCode += fmt.Sprintf("func (r *%sService) Create(model *models.%s) (e error) {\n",
		table.Name,
		table.Name,
	)

	goCode += fmt.Sprintf("\treturn r.Repos.%sRepo.Create(model)\n", table.Name)
	goCode += "}\n\n"

	// Update
	goCode += fmt.Sprintf("// Update updates an existing %s entry\n", table.Name)
	goCode += fmt.Sprintf("func (r *%sService) Update(model *models.%s) (e error) {\n",
		table.Name,
		table.Name,
	)
	goCode += fmt.Sprintf("\te = r.Repos.%sRepo.Update(model)\n", table.Name)
	goCode += "\treturn\n"
	goCode += "}\n\n"

	// Delete
	if isDeleted {
		goCode += fmt.Sprintf("// Delete marks an existing %s object as deleted\n", table.Name)
		goCode += fmt.Sprintf("func (r *%sService) Delete(model *models.%s) (e error) {\n",
			table.Name,
			table.Name,
		)

		goCode += fmt.Sprintf("\te = r.Repos.%sRepo.Delete(model)\n", table.Name)
		goCode += "\treturn\n"
		goCode += "}\n\n"
	}

	goCode += fmt.Sprintf("// HardDelete performs a SQL DELETE operation on a %s entry in the database\n", table.Name)
	goCode += fmt.Sprintf("func (r *%sService) HardDelete(model *models.%s) (e error) {\n",
		table.Name,
		table.Name,
	)

	goCode += fmt.Sprintf("\te = r.Repos.%sRepo.HardDelete(model)\n", table.Name)
	goCode += "\treturn\n"
	goCode += "}\n\n"

	// SelectByID
	goCode += fmt.Sprintf("// GetByID gets a single %s object by a Primary Key\n", table.Name)
	goCode += fmt.Sprintf("func (r *%sService) GetByID(%s %s) (model *models.%s, e error) {\n",
		table.Name,
		primaryKey,
		idType,
		table.Name,
	)

	goCode += fmt.Sprintf("\treturn r.Repos.%sRepo.GetByID(%s)\n", table.Name, primaryKey)
	goCode += "}\n\n"

	// Select
	goCode += fmt.Sprintf("// GetMany gets %s objects\n",
		table.Name,
	)

	goCode += fmt.Sprintf("func (r *%sService) GetMany(limit int, offset int, args ...string) (collection []*models.%s, e error) {\n",
		table.Name,
		table.Name,
	)

	goCode += fmt.Sprintf("\treturn r.Repos.%sRepo.GetMany(limit, offset, args...)\n", table.Name)
	goCode += "}\n\n"

	// Single
	goCode += fmt.Sprintf("// GetSingle gets one %s object\n",
		table.Name,
	)

	goCode += fmt.Sprintf("func (r *%sService) GetSingle(args ...string) (model *models.%s, e error) {\n",
		table.Name,
		table.Name,
	)

	goCode += fmt.Sprintf("\treturn r.Repos.%sRepo.GetSingle(args...)\n", table.Name)
	goCode += "}\n\n"
	goCode += "// #genEnd\n"

	return
}

// GenerateServicesBootstrapGoCodeFromDatabase generates golang code for a Repo Bootstrap file from
// a database object
func (g *Gen) GenerateServicesBootstrapGoCodeFromDatabase(database *lib.Database) (goCode string, e error) {

	props := ""
	defs := ""

	for _, table := range database.Tables {

		props += fmt.Sprintf("\t%s I%sService\n",
			table.Name,
			table.Name,
		)

		defs += fmt.Sprintf("\t// %s\n", table.Name)
		defs += fmt.Sprintf("\tservice%s := new(%sService)\n", table.Name, table.Name)
		defs += fmt.Sprintf("\tservice%s.Config = config\n", table.Name)
		defs += fmt.Sprintf("\tservice%s.Repos = repos\n", table.Name)
		defs += fmt.Sprintf("\tr.%s = service%s\n\n", table.Name, table.Name)
	}

	goCode = "package services"
	goCode += "\n\nimport("
	goCode += "\n\t\"" + g.Config.BasePackage + "/" + g.Config.Packages.Repos + "\""
	goCode += "\n)"
	goCode += "\n\n// Services is a collection of services"
	goCode += "\ntype Services struct {"
	goCode += "\n\tRepos " + g.Config.Packages.Repos + ".Repos"
	goCode += fmt.Sprintf("\n%s", props)
	goCode += "\n}"
	goCode += "\n\n//Bootstrap bootstraps all of the services into a single service object"
	goCode += "\nfunc Bootstrap(repos *Repos, config *models.Config) *Services {"
	goCode += "\n\n\tr := new(Services)"
	goCode += "\n\n\tr.Repos = repos"
	goCode += "\n\n\t// Services"
	goCode += fmt.Sprintf("\n\n%s", defs)
	goCode += "\n\n\treturn r"
	goCode += "\n}"

	return
}

// GenerateServicesBootstrapFile generates a repos bootstrap file in golang
func (g *Gen) GenerateServicesBootstrapFile(dir string, database *lib.Database) (e error) {

	// Make the repos dir if it does not exist.
	g.EnsureDir(dir)

	outFile := fmt.Sprintf("%s/services.go", dir)
	goCode, e := g.GenerateServicesBootstrapGoCodeFromDatabase(database)
	lib.Debugf("Generating go Repos bootstrap file at path %s", g.Options, outFile)
	if e != nil {
		return
	}

	e = g.WriteGoCodeToFile(goCode, outFile)

	return
}
