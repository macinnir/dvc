package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"io/ioutil"
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

	p := path.Join(dir, table.Name+"DAL.go")

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
		fmt.Sprintf("%s/%s/models", g.Config.BasePackage, g.Config.Dirs.Definitions),
		fmt.Sprintf("%s/%s/integrations", g.Config.BasePackage, g.Config.Dirs.Definitions),
		"database/sql",
		// "github.com/macinnir/dvc/modules/utils",
		// "github.com/macinnir/dvc/modules/query",
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
// Package dal is the Data Access Layer
package dal

import ({{range .Imports}}
	"{{.}}"{{end}}
)

// {{.Table.Name}}DAL is a data repository for {{.Table.Name}} objects 
type {{.Table.Name}}DAL struct {
	db  integrations.IDB
	log integrations.ILog
}

// New{{.Table.Name}}DAL returns a new instance of {{.Table.Name}}Repo
func New{{.Table.Name}}DAL(db integrations.IDB, log integrations.ILog) *{{.Table.Name}}DAL {
	return &{{.Table.Name}}DAL{db, log}
}

// Create creates a new {{.Table.Name}} entry in the database 
func (r {{.Table.Name}}DAL) Create(model *models.{{.Table.Name}}) (e error) {
	
	var result sql.Result 
	result, e = r.db.NamedExec({{.Table.Name}}DALInsertSQL, model)
	if e != nil {
		r.log.Errorf("{{.Table.Name}}DAL.Insert > %s", e.Error())
		return 
	}

	model.{{.PrimaryKey}}, e = result.LastInsertId()

	r.log.Debugf("{{.Table.Name}}DAL.Insert(%d)", model.{{.PrimaryKey}})
	return 
}

// Update updates an existing {{.Table.Name}} entry in the database 
func (r *{{.Table.Name}}DAL) Update(model *models.{{.Table.Name}}) (e error) {
	_, e = r.db.NamedExec({{.Table.Name}}DALUpdateSQL, model)
	if e != nil {
		r.log.Errorf("{{.Table.Name}}DAL.Update(%d) > %s", model.{{.PrimaryKey}}, e.Error())
	} else {
		r.log.Debugf("{{.Table.Name}}DAL.Update(%d)", model.{{.PrimaryKey}})
	}
	return 
}{{if .IsDeleted}}

// Delete marks an existing {{.Table.Name}} entry in the database as deleted
func (r *{{.Table.Name}}DAL) Delete(model *models.{{.Table.Name}}) (e error) {
	_, e = r.db.NamedExec("UPDATE ` + "`{{.Table.Name}}` SET `IsDeleted`" + ` = 1 WHERE {{.PrimaryKey}} = :{{.PrimaryKey}}", model)
	if e != nil {
		r.log.Errorf("{{.Table.Name}}DAL.Delete(%d) > %s", model.{{.PrimaryKey}}, e.Error())
	} else {
		r.log.Debugf("{{.Table.Name}}DAL.Delete(%d)", model.{{.PrimaryKey}})
	}
	return 
}{{end}} 

// HardDelete performs a SQL DELETE operation on a {{.Table.Name}} entry in the database
func (r *{{.Table.Name}}DAL) HardDelete(model *models.{{.Table.Name}}) (e error) {
	_, e = r.db.NamedExec("DELETE FROM ` + "`{{.Table.Name}}`" + ` WHERE {{.PrimaryKey}} = :{{.PrimaryKey}}", model) 
	if e != nil {
		r.log.Errorf("{{.Table.Name}}DAL.HardDelete(%d) > %s", model.{{.PrimaryKey}}, e.Error())
	} else {
		r.log.Debugf("{{.Table.Name}}DAL.HardDelete(%d)", model.{{.PrimaryKey}})
	}
	return 
}

// FromID gets a single {{.Table.Name}} object by its Primary Key
func (r *{{.Table.Name}}DAL) FromID({{.PrimaryKey}} {{.IDType}}) (model *models.{{.Table.Name}}, e error) {
	
	model = &models.{{.Table.Name}}{}
	
	e = r.db.Get(model, "SELECT * FROM ` + "`{{.Table.Name}}` WHERE `{{.PrimaryKey}}` = ?" + `", {{.PrimaryKey}})
	
	if e == nil {
		r.log.Debugf("{{.Table.Name}}DAL.FromID(%d)", model.{{.PrimaryKey}})
	} else if e == sql.ErrNoRows {
		e = nil 
		model = nil 
		r.log.Debugf("{{.Table.Name}}DAL.FromID(%d) > NOT FOUND", model.{{.PrimaryKey}})
	} else {
		r.log.Errorf("{{.Table.Name}}DAL.FromID(%d) > %s", model.{{.PrimaryKey}}, e.Error())
	}
	
	return 
}
{{.FileFoot}}`

	t := template.New("dal-" + table.Name)
	t.Funcs(template.FuncMap{
		"insertFields": fetchTableInsertFieldsString,
		"insertValues": fetchTableInsertValuesString,
		"updateFields": fetchTableUpdateFieldsString,
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

// GenerateDALSQL generates a constants file filled with sql statements
func (g *Gen) GenerateDALSQL(dir string, database *lib.Database) (e error) {

	var contents string
	var formatted []byte

	g.EnsureDir(dir)

	contents, e = generateDALSQL("dal", database)

	if e != nil {
		return
	}

	formatted, e = format.Source([]byte(contents))

	if e != nil {
		fmt.Println(contents)
		return
	}

	p := path.Join(dir, "sql.go")
	e = ioutil.WriteFile(p, formatted, 0644)
	return
}

func generateDALSQL(basePackage string, database *lib.Database) (out string, e error) {

	var sb strings.Builder

	sb.WriteString("package " + basePackage + "\n")

	sortedTables := make(lib.SortedTables, 0, len(database.Tables))

	// Find the primary key
	for _, table := range database.Tables {
		sortedTables = append(sortedTables, table)
	}

	sort.Sort(sortedTables)

	for _, table := range sortedTables {
		var outTable string
		outTable, e = generateTableInsertAndUpdateFields(table)
		if e != nil {
			return
		}

		sb.WriteString("\n" + outTable + "\n")
	}

	out = sb.String()
	return
}

// generateTableInsertAndUpdateFields generates insert and update fields as a string for use in their
// respective SQL queries
func generateTableInsertAndUpdateFields(table *lib.Table) (fields string, e error) {
	data := struct {
		Table          *lib.Table
		PrimaryKey     string
		PrimaryKeyType string
		IDType         string
		Columns        lib.SortedColumns
		IsDeleted      bool
	}{
		Table: table,
	}
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
	t := template.New("dal-fields")
	t.Funcs(
		template.FuncMap{
			"primaryKey":   fetchTablePrimaryKey,
			"insertFields": fetchTableInsertFieldsString,
			"insertValues": fetchTableInsertValuesString,
			"updateFields": fetchTableUpdateFieldsString,
		},
	)

	tpl := `// {{.Table.Name}}DAL SQL
const (
	{{.Table.Name}}DALInsertSQL = "INSERT INTO ` + "`{{.Table.Name}}`" + ` ({{.Columns | insertFields}}) VALUES ({{.Columns | insertValues}})" 
	{{.Table.Name}}DALUpdateSQL = "UPDATE ` + "`{{.Table.Name}}`" + ` SET {{.Columns | updateFields}} WHERE {{.PrimaryKey}} = :{{.PrimaryKey}}"
)`

	t, e = t.Parse(tpl)
	if e != nil {
		panic(e)
	}

	outBytes := []byte{}
	out := bytes.NewBuffer(outBytes)

	e = t.Execute(out, data)

	if e != nil {
		return
	}

	fields = out.String()
	return
}

func fetchTablePrimaryKey(table *lib.Table) string {
	primaryKey := ""
	idType := "int64"
	for _, column := range table.Columns {
		if column.ColumnKey == "PRI" {
			primaryKey = column.Name
		}
	}

	return primaryKey + " " + idType
}

func fetchTableInsertFieldsString(columns lib.SortedColumns) string {

	fields := []string{}

	for _, field := range columns {
		if field.ColumnKey == "PRI" {
			continue
		}

		if field.Name == "IsDeleted" {
			continue
		}

		fields = append(fields, "`"+field.Name+"`")
	}

	return strings.Join(fields, ",")
}

func fetchTableInsertValuesString(columns lib.SortedColumns) string {
	fields := []string{}
	for _, field := range columns {

		if field.ColumnKey == "PRI" {
			continue
		}

		if field.Name == "IsDeleted" {
			continue
		}

		fields = append(fields, ":"+field.Name)
	}

	return strings.Join(fields, ",")
}

func fetchTableUpdateFieldsString(columns lib.SortedColumns) string {
	fields := []string{}
	for _, field := range columns {

		if field.ColumnKey == "PRI" {
			continue
		}

		if field.Name == "IsDeleted" {
			continue
		}

		if field.Name == "DateCreated" {
			continue
		}

		fields = append(fields, "`"+field.Name+"` = :"+field.Name)
	}

	return strings.Join(fields, ",")
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
	t.Funcs(template.FuncMap{"primaryKey": fetchTablePrimaryKey})

	tpl := `
// Package definitions outlines objects and functionality used in the {{.BasePackage}} application
package definitions
import ({{range .Imports}}
	"{{.}}"{{end}}
)	

// DAL defines the container for all data access layer structs
type DAL struct {
	{{range .Tables}}
	{{.Name}} I{{.Name}}DAL{{end}}
}

{{range .Tables}}
// I{{.Name}}DAL outlines the repository methods on a {{.Name}} object 
type I{{.Name}}DAL interface {
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
		Tables              map[string]*lib.Table
		BasePackage         string
		IntegrationsPackage string
	}{
		BasePackage:         g.Config.BasePackage,
		Tables:              database.Tables,
		IntegrationsPackage: fmt.Sprintf("%s/%s/integrations", g.Config.BasePackage, g.Config.Dirs.Definitions),
	}

	tpl := `
// Package dal is the Data Access Layer
package dal

import (
	"{{ .IntegrationsPackage }}"
)

// DAL is a container for all dal structs
type DAL struct {
	{{range .Tables}}
	{{.Name}} *{{.Name}}DAL{{end}}
}

// BootstrapDAL bootstraps all of the DAL methods
func BootstrapDAL(db integrations.IDB, log integrations.ILog) *DAL {

	d := &DAL{} 
	{{range .Tables}}
	d.{{.Name}} = New{{.Name}}DAL(db, log){{end}}

	return d
}`

	p := path.Join(dir, "bootstrap.go")
	lib.Debugf("Generating dal bootstrap file at path %s", g.Options, p)
	t := template.Must(template.New("repos-bootstrap").Parse(tpl))
	buffer := bytes.Buffer{}

	e = t.Execute(&buffer, data)
	if e != nil {
		fmt.Println("Template Error: ", e.Error())
		return
	}

	var formatted []byte
	formatted, e = format.Source(buffer.Bytes())

	if e = ioutil.WriteFile(p, formatted, 0644); e != nil {
		fmt.Println("Write file error: ", e.Error())
	}

	// fmt.Println(string(buffer.Bytes()))

	return
}
