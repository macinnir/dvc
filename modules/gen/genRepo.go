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

// GetOrphanedRepos gets repo files that aren't in the database.Tables map
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

// CleanGoRepos removes any repo files that aren't in the database.Tables map
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

// GenerateGoRepo returns a string for a repo in golang
func (g *Gen) GenerateGoRepo(table *lib.Table, dir string) (e error) {

	imports := []string{}

	g.EnsureDir(dir)

	p := path.Join(dir, table.Name+".go")

	lib.Debugf("Generating go repo file for table %s at path %s", g.Options, table.Name, p)

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
		fmt.Sprintf("%s/models", g.Config.BasePackage),
		"database/sql",
		"github.com/jmoiron/sqlx",
		"log",
		"strings",
		"strconv",
		"github.com/macinnir/dvc/modules/utils",
		"fmt",
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

	data.Imports = imports

	tpl := `{{.FileHead}}
// #genStart
package repos 

import ({{range .Imports}}
	"{{.}}"{{end}}
)

// {{.Table.Name}}Repo is a data repository for {{.Table.Name}} objects 
type {{.Table.Name}}Repo struct {
	db *sqlx.DB
}

// New{{.Table.Name}}Repo returns a new instance of {{.Table.Name}}Repo
func New{{.Table.Name}}Repo(db *sqlx.DB) *{{.Table.Name}}Repo {
	return &{{.Table.Name}}Repo{db}
}

// Create creates a new {{.Table.Name}} entry in the database 
func (r {{.Table.Name}}Repo) Create(model *models.{{.Table.Name}}) (e error) {
	
	var result sql.Result 
	result, e = r.db.NamedExec("INSERT INTO ` + "`{{.Table.Name}}`" + ` ({{.Columns | insertFields}}) VALUES ({{.Columns | insertValues}})", model)

	if e != nil {
		return 
	}

	model.{{.PrimaryKey}}, e = result.LastInsertId()
	return 
}

	// if len(footMatches) > 0 {
	// 	footMatchPrefixLen := len(fmt.Sprintf("func (r *%sRepo) ", table.Name))
	// 	for _, footMatch := range footMatches {
	// 		footMatch = strings.Trim(footMatch, " ")
	// 		footMatchLen := len(footMatch)
	// 		goCode += "\t" + footMatch[footMatchPrefixLen:footMatchLen-1] + "\n"
	// 	}
	// }

// Update updates an existing {{.Table.Name}} entry in the database 
func (r *{{.Table.Name}}Repo) Update(model *models.{{.Table.Name}}) (e error) {
	_, e = r.db.NamedExec("UPDATE ` + "`{{.Table.Name}}`" + ` SET {{.Columns | updateFields}} WHERE {{.PrimaryKey}} = :{{.PrimaryKey}}", model)
	return 
}{{if .IsDeleted}}

// Delete marks an existing {{.Table.Name}} entry in the database as deleted
func (r *{{.Table.Name}}Repo) Delete(model *models.{{.Table.Name}}) (e error) {
	_, e = r.db.NamedExec("UPDATE ` + "`{{.Table.Name}}` SET `IsDeleted`" + ` = 1 WHERE {{.PrimaryKey}} = :{{.PrimaryKey}}", model)
	return 
}{{end}} 

// HardDelete performs a SQL DELETE operation on a {{.Table.Name}} entry in the database
func (r *{{.Table.Name}}Repo) HardDelete(model *models.{{.Table.Name}}) (e error) {
	_, e = r.db.NamedExec("DELETE FROM ` + "`{{.Table.Name}}`" + ` WHERE {{.PrimaryKey}} = :{{.PrimaryKey}}", model) 
	return 
}

// GetByID gets a single {{.Table.Name}} object by a Primary Key
func (r *{{.Table.Name}}Repo) GetByID({{.PrimaryKey}} {{.IDType}}) (model *models.{{.Table.Name}}, e error) {
	log.Println("Getting {{.Table.Name}} at ID ", {{.PrimaryKey}})
	model = &models.{{.Table.Name}}{}
	if e = r.db.Get(model, "SELECT * FROM ` + "`{{.Table.Name}}` WHERE `{{.PrimaryKey}}` = ?" + `", {{.PrimaryKey}}); e == sql.ErrNoRows {
		e = utils.NewRecordNotFoundError("{{.Table.Name}}", strconv.Itoa(int({{.PrimaryKey}})))
	}
	return 
}

// GetMany gets {{.Table.Name}} objects 
func (r *{{.Table.Name}}Repo) GetMany(args map[string]interface{}, orderBy map[string]string, limit []int64) (collection []*models.{{.Table.Name}}, e error) {
	collection = []*models.{{.Table.Name}}{}
	n := 1
	where := []string{"1=1"} 
	whereArgs := []interface{}{} 
	for field, val := range args {
		where = append(where, field + " = ?")
		whereArgs = append(whereArgs, val)
		n++
	}
	query := "SELECT * FROM ` + "`{{.Table.Name}}`" + ` WHERE " + strings.Join(where, " AND ") 
	orderBys := []string{} 
	if len(orderBy) > 0 {
		for col, dir := range orderBy {
			if dir != "ASC" && dir != "DESC" {
				e = errors.New("Invalid order by on table {{.Table.Name}}")
				return 
			}
			orderBys = append(orderBys, fmt.Sprintf("%s %s", col, dir)) 
		}
	}

	if len(orderBys) > 0 {
		query += " ORDER BY " + strings.Join(orderBys, ",")
	}

	if len(limit) > 0 {
		query += " LIMIT " + strconv.Itoa(int(limit[0]));
		if(len(limit) > 1) {
			 query += "," + strconv.Itoa(int(limit[1])); 
		} 
	}

	e = r.db.Select(&collection, query, whereArgs...) 

	return 
}
// #genEnd
{{.FileFoot}}`

	t := template.New("repo-" + table.Name)
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

func (g *Gen) GenerateRepoInterfaces(database *lib.Database, dir string) (e error) {

	var data = struct {
		Imports []string
		Tables  map[string]*lib.Table
	}{
		Imports: []string{
			fmt.Sprintf("%s/models", g.Config.BasePackage),
		},
		Tables: database.Tables,
	}

	t := template.New("repo-interface")
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
// #genStart
package services 
import ({{range .Imports}}
	"{{.}}"{{end}}
)	

// Repos is the container for all repositories
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
	GetMany(args map[string]interface{}, orderBy map[string]string, limit []int64) (collection []*models.{{.Name}}, e error) 
	GetByID({{. | primaryKey}}) (model *models.{{.Name}}, e error) 
}
{{end}}
`
	// {{if .Columns.IsDeleted}}Delete(model *models.{{.Name}}) (e error){{end}}
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

// GenerateReposBootstrapFile generates a repos bootstrap file in golang
func (g *Gen) GenerateReposBootstrapFile(dir string, database *lib.Database) (e error) {

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
package repos 
import (
	"github.com/jmoiron/sqlx"
	"{{.BasePackage}}/services"
)

// Bootstrap instantiates a collection of repositories
func Bootstrap(db *sqlx.DB) *services.Repos {
	r := new(services.Repos) 
	{{range .Tables}}
	// {{.Name}}
	r.{{.Name}} = New{{.Name}}Repo(db)
	{{end}}
	return r
}`

	p := path.Join(dir, "repos.go")
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
