package gen

import (
	"fmt"
	"github.com/macinnir/dvc/lib"
	"html/template"
	"os"
	"path"
	"sort"
	"strings"
)

// GenerateGoServiceFile generates a repo file in golang
func (g *Gen) GenerateGoServiceFile(dir string, table *lib.Table) (e error) {

	fileHead := ""
	fileFoot := ""
	imports := []string{}

	g.EnsureDir(dir)

	outFile := path.Join(dir, table.Name)

	lib.Debugf("Generating go service file for table %s at path %s", g.Options, table.Name, outFile)

	if fileHead, fileFoot, imports, e = g.scanFileParts(outFile, true); e != nil {
		return
	}

	e = g.GenerateGoService(table, fileHead, fileFoot, imports, dir)
	return
}

// GenerateGoServiceFiles generates go repository files based on the database schema
func (g *Gen) GenerateGoServiceFiles(reposDir string, database *lib.Database) (e error) {

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
func (g *Gen) GenerateGoService(table *lib.Table, fileHead string, fileFoot string, imports []string, dir string) (e error) {

	primaryKey := ""
	primaryKeyType := ""

	funcSig := fmt.Sprintf(`^func \(r \*%sService\) [A-Z].*$`, table.Name)

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
		fmt.Sprintf("%s/models", g.Config.BasePackage),
		"github.com/macinnir/dvc/modules/utils",
		"fmt",
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
		footMatchPrefixLen := len(fmt.Sprintf("func (r *%sService) ", table.Name))
		for _, footMatch := range footMatches {
			footMatch = strings.Trim(footMatch, " ")
			footMatchLen := len(footMatch)
			footMatchCode = append(footMatchCode, footMatch[footMatchPrefixLen:footMatchLen-1])
		}
	}

	data := struct {
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
package services 

import ({{range .Imports}}
	"{{.}}"{{end}} 
)

// I{{.Name}}Service outlines the service methods for {{.Name}} objects 
type I{{.Name}}Service interface {
	Create(model *models.{{.Name}}) (e error)
	Update(model *models.{{.Name}}) (e error) 
{{if .IsDeleted }}	Delete(model *models.{{.Name}}) (e error){{end}}
	HardDelete(model *models.{{.Name}}) (e error) 
	GetMany(args map[string]interface{}, orderBy map[string]string, limit []int64) (collection []*models.{{.Name}}, e error) 
	GetByID({{.PrimaryKey}} {{.PrimaryKeyType}})(model *models.{{.Name}}, e error) 
	{{range .FootMatches}}
	{{.}}
	{{end}}
}

// {{.Name}}Service is a service for '{{.Name}}' objects 
type {{.Name}}Service struct {
	config *models.Config
	repos *Repos 
	store utils.IStore
}

// New{{.Name}}Service returns a new instance of the {{.Name}}Service
func New{{.Name}}Service(config *models.Config, repos *Repos, store utils.IStore) *{{.Name}}Service {
	return &{{.Name}}Service{config, repos, store}
}

// Create creates a new {{.Name}} entry
func (r *{{.Name}}Service) Create(model *models.{{.Name}}) (e error) {
	e = r.repos.{{.Name}}.Create(model) 
	if e != nil {
		r.store.Set(fmt.Sprintf("{{.Name}}_%d", model.{{.PrimaryKey}}), model) 
	}
	return 
}

// Update updates an existing {{.Name}} entry
func (r *{{.Name}}Service) Update(model *models.{{.Name}}) (e error) {
	// Has Permission? -- CreatedBy? InGroup?
	e = r.repos.{{.Name}}.Update(model)
	if e == nil {
		r.store.Set(fmt.Sprintf("{{.Name}}_%d", model.{{.PrimaryKey}}), model) 
	}
	return 
}
{{if .IsDeleted}}
// Delete marks an existing {{.Name}} object as deleted
func (r *{{.Name}}Service) Delete(model *models.{{.Name}}) (e error) {
	e = r.repos.{{.Name}}.Delete(model)
	r.store.Delete(fmt.Sprintf("{{.Name}}_%d", model.{{.PrimaryKey}}))
	return 
}
{{end}}

// HardDelete performs a SQL DELETE operation on a {{.Name}} entry in the database
func (r *{{.Name}}Service) HardDelete(model *models.{{.Name}}) (e error) {
	e = r.repos.{{.Name}}.HardDelete(model)
	r.store.Delete(fmt.Sprintf("{{.Name}}_%d", model.{{.PrimaryKey}}))
	return
}

// GetByID gets a single {{.Name}} object by a Primary Key
func (r {{.Name}}Service) GetByID({{.PrimaryKey}} {{.PrimaryKeyType}})(model *models.{{.Name}}, e error) {
	model = &models.{{.Name}}{}
	if e = r.store.Get(fmt.Sprintf("{{.Name}}_%d", {{.PrimaryKey}}), model); e != nil {
		fmt.Println("Getting model from database...")
		if model, e = r.repos.{{.Name}}.GetByID({{.PrimaryKey}}); e == nil {
			r.store.Set(fmt.Sprintf("{{.Name}}_%d", {{.PrimaryKey}}), model)
		} else {
			return 
		}
	} 

	{{if ne .OneToMany ""}}
	model.{{.OneToMany}}s, _ = r.repos.{{.OneToMany}}.GetMany(map[string]interface{}{ "{{.Name}}ID": {{.Name}}ID }, map[string]string{}, []int64{})
	{{end}}

	{{if ne .OneToOne ""}}
	model.{{.OneToOne}}, _ = r.repos.{{.OneToOne}}.GetByID(model.{{.OneToOne}}ID)
	{{end}} 
	return 
}

// GetMany gets {{.Name}} objects\n",
func (r *{{.Name}}Service) GetMany(args map[string]interface{}, orderBy map[string]string, limit []int64) (collection []*models.{{.Name}}, e error) {
	return r.repos.{{.Name}}.GetMany(args, orderBy, limit)
}

// #genEnd
{{.FileFoot}} 
`

	// Store    utils.IStore
	p := path.Join(dir, table.Name+".go")
	t := template.Must(template.New("service-" + table.Name).Parse(tpl))
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

// GenerateServicesBootstrapGoCodeFromDatabase generates golang code for a Repo Bootstrap file from
// a database object
func (g *Gen) GenerateServicesBootstrapGoCodeFromDatabase(database *lib.Database, dir string) (e error) {

	data := struct {
		Tables      map[string]*lib.Table
		BasePackage string
	}{
		Tables:      database.Tables,
		BasePackage: g.Config.BasePackage,
	}

	tpl := `
package services 

import (
	"{{.BasePackage}}/models"
	"github.com/macinnir/dvc/modules/utils" 
)


// Services is a collection of services 
type Services struct {
	{{range .Tables}}
	{{.Name}} I{{.Name}}Service 
	{{end}} 
}

// NewServices bootstraps all of the services into a single service object 
func NewServices (config *models.Config, repos *Repos, store utils.IStore) *Services {
	r := new(Services) 
	{{range .Tables}}
	r.{{.Name}} = New{{.Name}}Service(config, repos, store)
	{{end}} 
	return r
}
`

	// Store    utils.IStore
	p := path.Join(dir, "services.go")
	t := template.Must(template.New("service-services").Parse(tpl))
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

// GenerateServicesBootstrapFile generates a repos bootstrap file in golang
func (g *Gen) GenerateServicesBootstrapFile(dir string, database *lib.Database) (e error) {

	// Make the repos dir if it does not exist.
	g.EnsureDir(dir)
	e = g.GenerateServicesBootstrapGoCodeFromDatabase(database, dir)
	return
}

// GetOrphanedServices returns a slice of service files that are not in the database.Tables map
func (g *Gen) GetOrphanedServices(dir string, database *lib.Database) []string {
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
func (g *Gen) CleanGoServices(dir string, database *lib.Database) (e error) {
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
