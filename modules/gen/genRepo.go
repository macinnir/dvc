package gen

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/macinnir/dvc/lib"
)

// GenerateGoRepoFile generates a repo file in golang
func (g *Gen) GenerateGoRepoFile(dir string, table *lib.Table) (e error) {

	fileHead := ""
	fileFoot := ""
	imports := []string{}

	g.EnsureDir(dir)

	outFile := path.Join(dir, table.Name)

	lib.Debugf("Generating go repo file for table %s at path %s", g.Options, table.Name, outFile)

	if fileHead, fileFoot, imports, e = g.scanFileParts(outFile, true); e != nil {
		return
	}

	e = g.GenerateGoRepo(table, fileHead, fileFoot, imports, dir)
	return
}

// GenerateGoRepoFiles generates go repository files based on the database schema
func (g *Gen) GenerateGoRepoFiles(dir string, database *lib.Database) (e error) {

	for _, table := range database.Tables {

		lib.Debugf("Generating repo %s", g.Options, table.Name)
		e = g.GenerateGoRepoFile(dir, table)
		if e != nil {
			return
		}
	}

	return
}

// GenerateRepoInterfaces generates a go interfaces file for use by the services directory
func (g *Gen) GenerateRepoInterfaces(database *lib.Database, dir string) (e error) {

	var data = struct {
		BasePackage string
		Imports     []string
		Tables      map[string]*lib.Table
	}{
		BasePackage: g.Config.BasePackage,
		Imports: []string{
			fmt.Sprintf("%s/definitions/models", g.Config.BasePackage),
			"github.com/macinnir/dvc/modules/query",
		},
		Tables: database.Tables,
	}

	t := template.New("repo-interfaces")

	t.Funcs(template.FuncMap{"primaryKey": func(table *lib.Table) string {
		primaryKey := ""
		idType := "int64"
		for _, column := range table.Columns {
			if column.ColumnKey == "PRI" {
				primaryKey = column.Name
			}
		}

		return primaryKey + " " + idType
	}})

	tpl := `
// Package definitions outlines objects and functionality used in the {{.BasePackage}} application
package definitions
import ({{range .Imports}}
	"{{.}}"{{end}}
)

// Repos defines the container for all repository layer structs
type Repos struct {
	{{range .Tables}}
	{{.Name}} I{{.Name}}Repo{{end}}
}

{{range .Tables}}
// I{{.Name}}Repo outlines the repository methods on a {{.Name}} object 
type I{{.Name}}Repo interface {
	Create(model *models.{{.Name}}) (e error) 
	Update(model *models.{{.Name}}) (e error) 
	{{if .Columns.IsDeleted}}Delete(model *models.{{.Name}}) (e error){{end}}
	HardDelete(model *models.{{.Name}}) (e error) 
	GetByID({{. | primaryKey}}) (model *models.{{.Name}}, e error) 
	Run(q *query.SelectQuery) (collection []*models.{{.Name}}, e error) 
}
{{end}}
`

	p := path.Join(dir, "repos.go")
	f, _ := os.Create(p)
	t, _ = t.Parse(tpl)
	e = t.Execute(f, data)
	if e != nil {
		fmt.Println("Execute Error: ", e.Error())
	}
	f.Close()
	g.FmtGoCode(p)
	return
}

// GenerateGoRepo returns a string for a repo in golang
func (g *Gen) GenerateGoRepo(table *lib.Table, fileHead string, fileFoot string, imports []string, dir string) (e error) {

	primaryKey := ""
	primaryKeyType := ""

	funcSig := fmt.Sprintf(`^func \(r \*%sRepo\) [A-Z].*$`, table.Name)

	footMatches := g.scanStringForFuncSignature(fileFoot, funcSig)

	sortedColumns := make(lib.SortedColumns, 0, len(table.Columns))

	oneToMany := g.Config.OneToMany[table.Name]
	oneToOne := g.Config.OneToOne[table.Name]

	// fmt.Println("OneToMany", table.Name, " ==> ", oneToMany)

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
		fmt.Sprintf("%s/definitions/models", g.Config.BasePackage),
		fmt.Sprintf("%s/definitions", g.Config.BasePackage),
		"github.com/macinnir/dvc/modules/utils",
		"github.com/macinnir/dvc/modules/query",
		"fmt",
		"log",
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

	footMatchCode := []string{}
	if len(footMatches) > 0 {
		footMatchPrefixLen := len(fmt.Sprintf("func (r *%sRepo) ", table.Name))
		for _, footMatch := range footMatches {
			footMatch = strings.Trim(footMatch, " ")
			footMatchLen := len(footMatch)
			footMatchCode = append(footMatchCode, footMatch[footMatchPrefixLen:footMatchLen-1])
		}
	}

	data := struct {
		BasePackage    string
		OneToMany      string
		OneToOne       string
		Imports        []string
		Name           string
		IsDeleted      bool
		PrimaryKey     string
		PrimaryKeyType string
		FootMatches    []string
		FileHead       string
		FileFoot       string
	}{
		BasePackage:    g.Config.BasePackage,
		OneToMany:      oneToMany,
		OneToOne:       oneToOne,
		Imports:        imports,
		Name:           table.Name,
		IsDeleted:      isDeleted,
		PrimaryKey:     primaryKey,
		PrimaryKeyType: idType,
		FootMatches:    footMatchCode,
		FileHead:       fileHead,
		FileFoot:       fileFoot,
	}

	tpl := `
{{.FileHead}}
// #genStart 
package repos

import ({{range .Imports}}
	"{{.}}"{{end}} 
)

// {{.Name}}Repo is a repository for {{.Name}} objects 
type {{.Name}}Repo struct {
	config *models.Config
	dal *definitions.Dal 
	store utils.IStore
}

// New{{.Name}}Repo returns a new instance of the {{.Name}}Repo
func New{{.Name}}Repo(config *models.Config, dal *definitions.Dal, store utils.IStore) *{{.Name}}Repo {
	return &{{.Name}}Repo{config, dal, store}
}

// Create creates a new {{.Name}} entry
func (r *{{.Name}}Repo) Create(model *models.{{.Name}}) (e error) {
	e = r.dal.{{.Name}}.Create(model) 
	if e != nil {
		log.Printf("ERR {{.Name}}Repo.Create > %s", e.Error())
	} else {
		log.Printf("INF {{.Name}}Repo.Create > #%d", model.{{.PrimaryKey}})
		r.store.Set(fmt.Sprintf("{{.BasePackage}}_{{.Name}}_%d", model.{{.PrimaryKey}}), model) 
	}
	return 
}

// Update updates an existing {{.Name}} entry
func (r *{{.Name}}Repo) Update(model *models.{{.Name}}) (e error) {
	e = r.dal.{{.Name}}.Update(model)
	if e != nil {
		log.Printf("ERR {{.Name}}Repo.Update > %s", e.Error()) 
	} else {
		log.Printf("INF {{.Name}}Repo.Update > #%d", model.{{.PrimaryKey}})
		r.store.Set(fmt.Sprintf("{{.BasePackage}}_{{.Name}}_%d", model.{{.PrimaryKey}}), model) 
	}
	return 
}
{{if .IsDeleted}}
// Delete marks an existing {{.Name}} object as deleted
func (r *{{.Name}}Repo) Delete(model *models.{{.Name}}) (e error) {
	e = r.dal.{{.Name}}.Delete(model)
	if e == nil {
		log.Printf("INF {{.Name}}Repo.Delete > #%d", model.{{.PrimaryKey}}) 
		r.store.Delete(fmt.Sprintf("{{.BasePackage}}_{{.Name}}_%d", model.{{.PrimaryKey}}))
	} else {
		log.Printf("ERR {{.Name}}Repo.Delete > %s", e.Error()) 
	}
	return 
}
{{end}}

// HardDelete performs a SQL DELETE operation on a {{.Name}} entry in the database
func (r *{{.Name}}Repo) HardDelete(model *models.{{.Name}}) (e error) {
	e = r.dal.{{.Name}}.HardDelete(model)
	if e != nil {
		log.Printf("ERR {{.Name}}Repo.HardDelete > %s", e.Error()) 
	} else {
		log.Printf("INF {{.Name}}Repo.HardDelete > #%d", model.{{.PrimaryKey}}) 
		r.store.Delete(fmt.Sprintf("{{.BasePackage}}_{{.Name}}_%d", model.{{.PrimaryKey}}))
	}

	return
}

// GetByID gets a single {{.Name}} object by a Primary Key
func (r {{.Name}}Repo) GetByID({{.PrimaryKey}} {{.PrimaryKeyType}})(model *models.{{.Name}}, e error) {
	model = &models.{{.Name}}{}
	
	e = r.store.Get(fmt.Sprintf("{{.BasePackage}}_{{.Name}}_%d", {{.PrimaryKey}}), model)
	
	if e == nil {
		log.Printf("INF {{.Name}}Repo.GetByID > #%d [From Cache]", model.{{.PrimaryKey}})
	} else {
		if model, e = r.dal.{{.Name}}.GetByID({{.PrimaryKey}}); e == nil {
			log.Printf("INF {{.Name}}Repo.GetByID > #%d", model.{{.PrimaryKey}})
			r.store.Set(fmt.Sprintf("{{.BasePackage}}_{{.Name}}_%d", {{.PrimaryKey}}), model)
		} else {
			log.Printf("ERR {{.Name}}Repo.GetByID > #%d > %s", model.{{.PrimaryKey}}, e.Error())
			return 
		}
	} 

	{{if ne .OneToMany ""}}
	model.{{.OneToMany}}s, _ = r.dal.{{.OneToMany}}.Run(query.Select().Where(query.Equals{"{{.Name}}ID": {{.Name}}ID}))
	{{end}}

	{{if ne .OneToOne ""}}
	model.{{.OneToOne}}, _ = r.dal.{{.OneToOne}}.GetByID(model.{{.OneToOne}}ID)
	{{end}} 
	return 
}

// Run runs the select query (from Select()) for {{.Name}} objects
func (r *{{.Name}}Repo) Run(q *query.SelectQuery) (collection []*models.{{.Name}}, e error) {
	log.Println("INF {{.Name}}Repo.GetMany")
	return r.dal.{{.Name}}.Run(q)
}

// #genEnd
{{.FileFoot}} 
`

	// Store    utils.IStore
	p := path.Join(dir, table.Name+".go")
	t := template.Must(template.New("repo-" + table.Name).Parse(tpl))
	f, err := os.Create(p)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}

	err = t.Execute(f, data)
	if err != nil {
		fmt.Println("Execute Error: ", err.Error())
		return
	}

	f.Close()

	g.FmtGoCode(p)

	return
}

// GenerateReposBootstrapGoCodeFromDatabase generates golang code for a Repo Bootstrap file from
// a database object
func (g *Gen) GenerateReposBootstrapGoCodeFromDatabase(database *lib.Database, dir string) (e error) {

	data := struct {
		Tables      map[string]*lib.Table
		BasePackage string
	}{
		Tables:      database.Tables,
		BasePackage: g.Config.BasePackage,
	}

	tpl := `
package repos

import (
	"{{.BasePackage}}/definitions"
	"{{.BasePackage}}/definitions/models"
	"github.com/macinnir/dvc/modules/utils" 
)

// Bootstrap bootstraps all of the repos into a single repo object 
func Bootstrap (config *models.Config, dal *definitions.Dal, store utils.IStore) *definitions.Repos {
	r := new(definitions.Repos) 
	{{range .Tables}}
	r.{{.Name}} = New{{.Name}}Repo(config, dal, store)
	{{end}} 
	return r
}
`

	// Store    utils.IStore
	p := path.Join(dir, "bootstrap.go")
	t := template.Must(template.New("repo-repos").Parse(tpl))
	f, err := os.Create(p)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}

	err = t.Execute(f, data)
	if err != nil {
		fmt.Println("Execute Error: ", err.Error())
		return
	}

	f.Close()
	g.FmtGoCode(p)

	return

}

// GenerateReposBootstrapFile generates a repos bootstrap file in golang
func (g *Gen) GenerateReposBootstrapFile(dir string, database *lib.Database) (e error) {

	// Make the repos dir if it does not exist.
	g.EnsureDir(dir)
	e = g.GenerateReposBootstrapGoCodeFromDatabase(database, dir)
	return
}

// GetOrphanedRepos returns a slice of service files that are not in the database.Tables map
func (g *Gen) GetOrphanedRepos(dir string, database *lib.Database) []string {
	dirHandle, err := os.Open(dir)
	if err != nil {
		panic(err)
	}

	defer dirHandle.Close()
	var dirFileNames []string
	dirFileNames, err = dirHandle.Readdirnames(-1)
	if err != nil {
		panic(err)
	}

	orphans := []string{}

	for _, name := range dirFileNames {

		if fileInfo, e := os.Stat(name); e != nil || fileInfo.IsDir() {
			continue
		}
		// Skip tests, repo definitions and service definitions
		if (len(name) > 8 && name[len(name)-8:len(name)] == "_test.go") || name == "repos.go" || name == "services.go" {
			continue
		}

		fileNameNoExt := name[0 : len(name)-3]
		if _, ok := database.Tables[fileNameNoExt]; !ok {
			orphans = append(orphans, name)
		}
	}

	return orphans
}

// CleanGoServices removes service files not found in the database.Tables map
func (g *Gen) CleanGoRepos(dir string, database *lib.Database) (e error) {
	dirHandle, err := os.Open(dir)
	if err != nil {
		panic(err)
	}

	defer dirHandle.Close()
	var dirFileNames []string
	dirFileNames, err = dirHandle.Readdirnames(-1)
	if err != nil {
		panic(err)
	}

	for _, name := range dirFileNames {

		if fileInfo, e := os.Stat(name); e != nil || fileInfo.IsDir() {
			continue
		}
		// Skip tests, repo definitions and service definitions
		if (len(name) > 8 && name[len(name)-8:len(name)] == "_test.go") || name == "repos.go" || name == "services.go" {
			continue
		}

		fileNameNoExt := name[0 : len(name)-3]
		if _, ok := database.Tables[fileNameNoExt]; !ok {
			fullFilePath := path.Join(dir, name)
			fmt.Printf("TEST: Removing %s\n", fullFilePath)
			os.Remove(fullFilePath)
		}
	}
	return
}
