package gen

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/macinnir/dvc/lib"
)

// GetOrphanedDals gets repo files that aren't in the database.Tables map
func (g *Gen) GetOrphanedDals(dir string, database *lib.Database) []string {
	dirHandle, err := os.Open(dir)
	if err != nil {
		log.Fatalf("Directory not found: %s", dir)
		// panic(err)
	}

	defer dirHandle.Close()
	var dirFileNames []string
	dirFileNames, err = dirHandle.Readdirnames(-1)
	if err != nil {
		panic(err)
	}

	orphans := []string{}

	for _, name := range dirFileNames {

		// Skip tests
		if (len(name) > 8 && name[len(name)-8:len(name)] == "_test.go") || name == "repos.go" {
			continue
		}

		fileNameNoExt := name[0 : len(name)-3]
		if _, ok := database.Tables[fileNameNoExt]; !ok {
			orphans = append(orphans, name)
		}
	}

	return orphans
}

// CleanGoDALs removes any repo files that aren't in the database.Tables map
func (g *Gen) CleanGoDALs(dir string, database *lib.Database) (e error) {
	dirHandle, err := os.Open(dir)
	if err != nil {
		log.Fatalf("Directory not found: %s", dir)
	}

	defer dirHandle.Close()
	var dirFileNames []string
	dirFileNames, err = dirHandle.Readdirnames(-1)
	if err != nil {
		panic(err)
	}

	for _, name := range dirFileNames {

		// Skip tests
		if (len(name) > 8 && name[len(name)-8:len(name)] == "_test.go") || name == "repos.go" {
			continue
		}

		fileNameNoExt := name[0 : len(name)-3]
		if _, ok := database.Tables[fileNameNoExt]; !ok {
			if fileNameNoExt != "Config" {
				fullFilePath := path.Join(dir, name)
				fmt.Printf("Removing %s\n", fullFilePath)
				os.Remove(fullFilePath)
			}
		}
	}
	return
}

// GenerateGoDAL returns a string for a repo in golang
func (g *Gen) GenerateGoDAL(table *lib.Table, dir string) (e error) {

	imports := []string{}

	g.EnsureDir(dir)

	p := path.Join(dir, table.Name+".go")

	lib.Debugf("Generating go dal file for table %s at path %s", g.Options, table.Name, p)

	data := struct {
		Table          *lib.Table
		Columns        lib.SortedColumns
		PrimaryKey     string
		PrimaryKeyType string
		IDType         string
		IsDeleted      bool
		Imports        []string
		FileHead       string
		FileFoot       string
	}{
		Table:          table,
		PrimaryKey:     "",
		PrimaryKeyType: "",
		IDType:         "int64",
		IsDeleted:      false,
		Imports:        []string{},
		FileHead:       "",
		FileFoot:       "",
	}

	if data.FileHead, data.FileFoot, imports, e = g.scanFileParts(p, true); e != nil {
		lib.Errorf("ERROR: %s", g.Options, e.Error())
		return
	}

	// funcSig := fmt.Sprintf(`^func \(r \*%sRepo\) [A-Z].*$`, table.Name)
	// footMatches := g.scanStringForFuncSignature(fileFoot, funcSig)

	sortedColumns := make(lib.SortedColumns, 0, len(table.Columns))

	// Find the primary key
	for _, column := range table.Columns {
		if column.ColumnKey == "PRI" {
			data.PrimaryKey = column.Name
			data.PrimaryKeyType = column.DataType
		}
		sortedColumns = append(sortedColumns, column)
	}

	sort.Sort(sortedColumns)

	data.Columns = sortedColumns

	_, data.IsDeleted = table.Columns["IsDeleted"]

	switch data.PrimaryKeyType {
	case "varchar":
		data.IDType = "string"
	}

	defaultImports := []string{
		fmt.Sprintf("%s/definitions/models", g.Config.BasePackage),
		"database/sql",
		"github.com/jmoiron/sqlx",
		"log",
		"github.com/macinnir/dvc/modules/utils",
		"github.com/macinnir/dvc/modules/query",
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

	data.Imports = imports

	tpl := `{{.FileHead}}
// #genStart
// Package dal is the Data Access Layer
package dal

import ({{range .Imports}}
	"{{.}}"{{end}}
)

// {{.Table.Name}}Dal is a data repository for {{.Table.Name}} objects 
type {{.Table.Name}}Dal struct {
	db *sqlx.DB
}

// New{{.Table.Name}}Dal returns a new instance of {{.Table.Name}}Repo
func New{{.Table.Name}}Dal(db *sqlx.DB) *{{.Table.Name}}Dal {
	return &{{.Table.Name}}Dal{db}
}

// Create creates a new {{.Table.Name}} entry in the database 
func (r {{.Table.Name}}Dal) Create(model *models.{{.Table.Name}}) (e error) {
	
	var result sql.Result 
	result, e = r.db.NamedExec("INSERT INTO ` + "`{{.Table.Name}}`" + ` ({{.Columns | insertFields}}) VALUES ({{.Columns | insertValues}})", model)

	if e != nil {
		log.Printf("ERR {{.Table.Name}}Dal.Insert > %s", e.Error())
		return 
	}

	model.{{.PrimaryKey}}, e = result.LastInsertId()

	log.Printf("INF {{.Table.Name}}Dal.Insert > #%d", model.{{.PrimaryKey}})
	return 
}

// Update updates an existing {{.Table.Name}} entry in the database 
func (r *{{.Table.Name}}Dal) Update(model *models.{{.Table.Name}}) (e error) {
	_, e = r.db.NamedExec("UPDATE ` + "`{{.Table.Name}}`" + ` SET {{.Columns | updateFields}} WHERE {{.PrimaryKey}} = :{{.PrimaryKey}}", model)
	if e != nil {
		log.Printf("ERR {{.Table.Name}}Dal.Update > %s", e.Error())
	} else {
		log.Printf("INF {{.Table.Name}}Dal.Update > #%d", model.{{.PrimaryKey}})
	}
	return 
}{{if .IsDeleted}}

// Delete marks an existing {{.Table.Name}} entry in the database as deleted
func (r *{{.Table.Name}}Dal) Delete(model *models.{{.Table.Name}}) (e error) {
	_, e = r.db.NamedExec("UPDATE ` + "`{{.Table.Name}}` SET `IsDeleted`" + ` = 1 WHERE {{.PrimaryKey}} = :{{.PrimaryKey}}", model)
	if e != nil {
		log.Printf("ERR {{.Table.Name}}Dal.Delete > %s", e.Error())
	} else {
		log.Printf("INF {{.Table.Name}}Dal.Delete > #%d", model.{{.PrimaryKey}})
	}
	return 
}{{end}} 

// HardDelete performs a SQL DELETE operation on a {{.Table.Name}} entry in the database
func (r *{{.Table.Name}}Dal) HardDelete(model *models.{{.Table.Name}}) (e error) {
	_, e = r.db.NamedExec("DELETE FROM ` + "`{{.Table.Name}}`" + ` WHERE {{.PrimaryKey}} = :{{.PrimaryKey}}", model) 
	if e != nil {
		log.Printf("ERR {{.Table.Name}}Dal.HardDelete > %s", e.Error())
	} else {
		log.Printf("INF {{.Table.Name}}Dal.HardDelete > #%d", model.{{.PrimaryKey}})
	}
	return 
}

// GetByID gets a single {{.Table.Name}} object by a Primary Key
func (r *{{.Table.Name}}Dal) GetByID({{.PrimaryKey}} {{.IDType}}) (model *models.{{.Table.Name}}, e error) {
	model = &models.{{.Table.Name}}{}
	if e = r.db.Get(model, "SELECT * FROM ` + "`{{.Table.Name}}` WHERE `{{.PrimaryKey}}` = ?" + `", {{.PrimaryKey}}); e == sql.ErrNoRows {
		e = utils.NewRecordNotFoundError()
	}

	if e != nil {
		log.Printf("ERR {{.Table.Name}}Dal.GetByID > %s", e.Error())
	} else {
		log.Printf("INF {{.Table.Name}}Dal.GetByID > #%d", model.{{.PrimaryKey}})
	}

	return 
}

func (r *{{.Table.Name}}Dal) Run(q *query.SelectQuery) (collection []*models.{{.Table.Name}}, e error) {

	collection = []*models.{{.Table.Name}}{}
	q.Object = models.New{{.Table.Name}}()
	sql, args := q.ToSQL()

	e = r.db.Select(&collection, sql, args...)

	if e != nil {
		log.Printf("ERR {{.Table.Name}}Dal.GetMany > %s", e.Error())
	} else {
		log.Println("INF {{.Table.Name}}Dal.GetMany")
	}

	return
}

func (r *{{.Table.Name}}Dal) Count(q *query.CountQuery) (count int64, e error) {

	count = 0
	q.Object = models.New{{.Table.Name}}()
	sql, args := q.ToSQL()

	e = r.db.Get(&count, sql, args...)

	if e != nil {
		log.Printf("ERR {{.Table.Name}}Dal.Count > %s", e.Error())
	} else {
		log.Println("INF {{.Table.Name}}Dal.Count")
	}

	return
}
// #genEnd
{{.FileFoot}}`

	t := template.New("dal-" + table.Name)
	t.Funcs(template.FuncMap{
		"insertFields": func(columns lib.SortedColumns) string {
			fields := []string{}

			for _, field := range columns {
				if field.ColumnKey == "PRI" {
					continue
				}

				fields = append(fields, "`"+field.Name+"`")
			}

			return strings.Join(fields, ",")
		},
		"insertValues": func(columns lib.SortedColumns) string {
			fields := []string{}
			for _, field := range columns {
				if field.ColumnKey == "PRI" {
					continue
				}

				fields = append(fields, ":"+field.Name)
			}

			return strings.Join(fields, ",")
		},
		"updateFields": func(columns lib.SortedColumns) string {
			fields := []string{}
			for _, field := range columns {
				if field.ColumnKey == "PRI" {
					continue
				}
				fields = append(fields, "`"+field.Name+"` = :"+field.Name)
			}

			return strings.Join(fields, ",")
		},
	})

	t, e = t.Parse(tpl)
	if e != nil {
		panic(e)
	}

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

func (g *Gen) GenerateDALInterfaces(database *lib.Database, dir string) (e error) {

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

	t := template.New("dal-interface")
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

// Dal defines the container for all data access layer structs
type Dal struct {
	{{range .Tables}}
	{{.Name}} I{{.Name}}Dal{{end}}
}

{{range .Tables}}
// I{{.Name}}Dal outlines the repository methods on a {{.Name}} object 
type I{{.Name}}Dal interface {
	Create(model *models.{{.Name}}) (e error) 
	Update(model *models.{{.Name}}) (e error) 
	{{if .Columns.IsDeleted}}Delete(model *models.{{.Name}}) (e error){{end}}
	HardDelete(model *models.{{.Name}}) (e error) 
	GetByID({{. | primaryKey}}) (model *models.{{.Name}}, e error) 
	Run(q *query.SelectQuery) (collection []*models.{{.Name}}, e error) 
	Count(q *query.CountQuery) (count int64, e error) 
}
{{end}}
`
	// {{if .Columns.IsDeleted}}Delete(model *models.{{.Name}}) (e error){{end}}
	p := path.Join(dir, "dal.go")
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

// GenerateDALsBootstrapFile generates a dal bootstrap file in golang
func (g *Gen) GenerateDALsBootstrapFile(dir string, database *lib.Database) (e error) {

	// Make the repos dir if it does not exist.
	g.EnsureDir(dir)

	data := struct {
		Tables      map[string]*lib.Table
		BasePackage string
	}{
		BasePackage: g.Config.BasePackage,
		Tables:      database.Tables,
	}

	tpl := `
// Package dal is the Data Access Layer
package dal
import (
	"log" 
	"github.com/jmoiron/sqlx"
	"{{.BasePackage}}/definitions"
)

// Bootstrap instantiates a collection of repositories
func Bootstrap(db *sqlx.DB) *definitions.Dal {
	log.Println("INF Bootstrapping DAL") 
	r := new(definitions.Dal) 
	{{range .Tables}}
	// {{.Name}}
	r.{{.Name}} = New{{.Name}}Dal(db)
	{{end}}
	return r
}`

	p := path.Join(dir, "bootstrap.go")
	t := template.Must(template.New("repos-bootstrap").Parse(tpl))
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
