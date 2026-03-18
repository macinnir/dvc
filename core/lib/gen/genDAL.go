package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/gen/genutil"
	"github.com/macinnir/dvc/core/lib/schema"
)

// GetOrphanedDals gets repo files that aren't in the database.Tables map
func (g *Gen) GetOrphanedDals(dir string, database *schema.Schema) []string {
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
		if (len(name) > 8 && name[len(name)-8:] == "_test.go") || name == "repos.go" {
			continue
		}

		fileNameNoExt := name[0 : len(name)-3]
		if _, ok := database.Tables[fileNameNoExt]; !ok {
			orphans = append(orphans, name)
		}
	}

	return orphans
}

func GenDALs(tables []*schema.Table, config *lib.Config) error {

	// start := time.Now()
	generatedDALCount := 0
	lib.EnsureDir(lib.DalsGenDir)

	for k := range tables {
		if e := GenerateGoDAL(config, tables[k], lib.DalsGenDir); e != nil {
			return fmt.Errorf("GenDALs(%s): %w", tables[k].Name, e)
		}
		generatedDALCount++
	}

	// TODO Verbose flag
	// fmt.Printf("Generated %d dals in %f seconds.\n", generatedDALCount, time.Since(start).Seconds())
	return nil
}

// var dalTPL *template.Template

// GenerateGoDAL returns a string for a repo in golang
func GenerateGoDAL(config *lib.Config, table *schema.Table, dir string) (e error) {

	p := path.Join(dir, table.Name+"DAL.go")
	// TODO verbose flag
	// start := time.Now()
	// TODO verbose flag
	// lib.EnsureDir(dir)

	// imports := []string{}

	// // lib.Debugf("Generating go dal file for table %s at path %s", g.Options, table.Name, p)

	var data = struct {
		BasePackage       string
		Table             *schema.Table
		Columns           schema.SortedColumns
		UpdateColumns     []*schema.Column
		StringColumns     []*schema.Column
		SpecialColumns    []*schema.Column
		InsertSQL         string
		InsertArgs        string
		UpdateSQL         string
		UpdateArgs        string
		PrimaryKey        string
		PrimaryKeyType    string
		PrimaryKeyArgName string
		IDType            string
		IsDeleted         bool
		IsDateCreated     bool
		IsLastUpdated     bool
		Imports           []string
		FileHead          string
		FileFoot          string
		HasNull           bool
		HasSpecialColumns bool
	}{
		BasePackage:       config.BasePackage,
		Table:             table,
		UpdateColumns:     []*schema.Column{},
		StringColumns:     []*schema.Column{},
		SpecialColumns:    []*schema.Column{},
		InsertSQL:         "",
		InsertArgs:        "",
		UpdateSQL:         "",
		UpdateArgs:        "",
		PrimaryKey:        "",
		PrimaryKeyType:    "",
		PrimaryKeyArgName: "",
		IDType:            "int64",
		IsDeleted:         false,
		IsDateCreated:     false,
		IsLastUpdated:     false,
		Imports:           []string{},
		FileHead:          "",
		FileFoot:          "",
		HasNull:           false,
		HasSpecialColumns: false,
	}

	// if data.FileHead, data.FileFoot, imports, e = g.scanFileParts(p, true); e != nil {
	// 	lib.Errorf("ERROR: %s", g.Options, e.Error())
	// 	return
	// }

	// funcSig := fmt.Sprintf(`^func \(r \*%sRepo\) [A-Z].*$`, table.Name)
	// footMatches := g.scanStringForFuncSignature(fileFoot, funcSig)

	sortedColumns := make(schema.SortedColumns, 0, len(table.Columns))

	var insertColumns = []*schema.Column{}

	data.PrimaryKey = "Foo"

	// Find the primary key
	for k := range table.Columns {

		// fmt.Println("Column:", table.Columns[k].Name)

		var column = table.Columns[k]

		if column.ColumnKey == "PRI" {
			data.PrimaryKey = column.Name
			data.PrimaryKeyType = column.DataType
		}

		if schema.DataTypeToGoTypeString(column) == "string" {
			data.StringColumns = append(data.StringColumns, column)
		}

		sort.Slice(data.StringColumns, func(i, j int) bool {
			return data.StringColumns[i].Name < data.StringColumns[j].Name
		})

		goDataType := schema.DataTypeToGoTypeString(column)
		if len(goDataType) > 5 && goDataType[0:5] == "null." {
			data.HasNull = true
		}

		sortedColumns = append(sortedColumns, column)

		if genutil.IsInsertColumn(column) {
			insertColumns = append(insertColumns, column)
		}

		if genutil.IsUpdateColumn(column) {
			data.UpdateColumns = append(data.UpdateColumns, column)
		}

		if genutil.IsSpecialColumn(column) {
			data.SpecialColumns = append(data.SpecialColumns, column)
			data.HasSpecialColumns = true
		}

		sortedColumns = append(sortedColumns, column)
	}

	sort.Sort(sortedColumns)
	data.Columns = sortedColumns

	var insertColumnNames bytes.Buffer
	var insertColumnVals bytes.Buffer
	var insertColumnArgs bytes.Buffer

	sort.Slice(insertColumns, func(a, b int) bool { return insertColumns[a].Name < insertColumns[b].Name })
	sort.Slice(data.UpdateColumns, func(a, b int) bool { return data.UpdateColumns[a].Name < data.UpdateColumns[b].Name })

	for k, col := range insertColumns {

		insertColumnVals.WriteString("?")
		insertColumnNames.WriteString("`" + col.Name + "`")
		insertColumnArgs.WriteString("model." + col.Name)

		if k < len(insertColumns)-1 {
			insertColumnNames.WriteString(", ")
			insertColumnArgs.WriteString(", ")
			insertColumnVals.WriteString(", ")
		}
	}

	data.InsertArgs = insertColumnArgs.String()
	data.InsertSQL = "INSERT INTO `" + data.Table.Name + "` (" + insertColumnNames.String() + ") VALUES (" + insertColumnVals.String() + ")"

	var updateColumnNames bytes.Buffer
	var updateColumnArgs bytes.Buffer

	for k, col := range data.UpdateColumns {

		updateColumnNames.WriteString("`" + col.Name + "` = ?")
		updateColumnArgs.WriteString("model." + col.Name)
		updateColumnArgs.WriteString(", ")

		if k < len(data.UpdateColumns)-1 {
			updateColumnNames.WriteString(", ")
		}
	}

	updateColumnArgs.WriteString("model." + data.PrimaryKey)

	data.UpdateArgs = updateColumnArgs.String()
	data.UpdateSQL = "UPDATE `" + data.Table.Name + "` SET " + updateColumnNames.String() + " WHERE `" + data.PrimaryKey + "` = ?"

	_, data.IsDeleted = table.Columns["IsDeleted"]
	_, data.IsDateCreated = table.Columns["DateCreated"]
	_, data.IsLastUpdated = table.Columns["LastUpdated"]

	switch data.PrimaryKeyType {
	case "varchar":
		data.IDType = "string"
	}

	var buf bytes.Buffer
	if e = DALTemplate.Execute(&buf, data); e != nil {
		return
	}

	ioutil.WriteFile(p, buf.Bytes(), lib.DefaultFileMode)

	return
}

// GenerateDALSQL generates a constants file filled with sql statements
func (g *Gen) GenerateDALSQL(dir string, database *schema.Schema) (e error) {

	var contents string
	var formatted []byte

	lib.EnsureDir(dir)

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

func generateDALSQL(basePackage string, database *schema.Schema) (out string, e error) {

	var sb strings.Builder

	sb.WriteString("package " + basePackage + "\n")

	sortedTables := database.ToSortedTables()

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
func generateTableInsertAndUpdateFields(table *schema.Table) (fields string, e error) {
	data := struct {
		Table          *schema.Table
		PrimaryKey     string
		PrimaryKeyType string
		IDType         string
		Columns        schema.SortedColumns
		IsDeleted      bool
	}{
		Table: table,
	}
	sortedColumns := make(schema.SortedColumns, 0, len(table.Columns))

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
			"primaryKey":   genutil.FetchTablePrimaryKey,
			"insertFields": genutil.FetchTableInsertFieldsString,
			"insertValues": genutil.FetchTableInsertValuesString,
			"updateFields": genutil.FetchTableUpdateFieldsString,
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

var bootstrapFileTpl = template.Must(template.New("repos-bootstrap").Parse(`// Generated Code; DO NOT EDIT.
package definitions

import (
	"{{ .DBPackage }}"
	"{{ .LogPackage }}"
	"{{ .ModelsPackage }}"
	"{{ .DALPackage }}"
)

// DAL is a container for all dal structs
type DAL struct {
	{{range .Tables}}
	{{.Name}} *dal.{{.Name}}DAL{{end}}
}

// BootstrapDAL bootstraps all of the DAL methods
func BootstrapDAL(db map[string][]db.IDB, log log.ILog) *DAL {

	d := &DAL{}
	{{range .Tables}}
	d.{{.Name}} = dal.New{{.Name}}DAL(db[models.{{.Name}}_SchemaName], log){{end}}

	return d
}`))

// GenerateDALsBootstrapFile generates a dal bootstrap file in golang
func GenerateDALsBootstrapFile(config *lib.Config, schemaList *schema.SchemaList) error {

	// var start = time.Now()
	var e error

	tables := map[string]*schema.Table{}

	for k := range schemaList.Schemas {
		for l := range schemaList.Schemas[k].Tables {
			tables[l] = schemaList.Schemas[k].Tables[l]
		}
	}

	data := struct {
		Tables        map[string]*schema.Table
		BasePackage   string
		DBPackage     string
		LogPackage    string
		ModelsPackage string
		DALPackage    string
	}{
		BasePackage:   config.BasePackage,
		Tables:        tables,
		DBPackage:     "github.com/macinnir/dvc/core/lib/utils/db",
		LogPackage:    "github.com/macinnir/dvc/core/lib/utils/log",
		ModelsPackage: fmt.Sprintf("%s/%s", config.BasePackage, "gen/definitions/models"),
		DALPackage:    fmt.Sprintf("%s/%s", config.BasePackage, "gen/dal"),
	}

	// // lib.Debugf("Generating dal bootstrap file at path %s", g.Options, p)
	buffer := bytes.Buffer{}

	e = bootstrapFileTpl.Execute(&buffer, data)
	if e != nil {
		fmt.Println("Template Error: ", e.Error())
		return e
	}

	var formatted []byte
	if formatted, e = format.Source(buffer.Bytes()); e != nil {
		fmt.Println("Format Error:", e.Error())
		return e
	}

	if e = ioutil.WriteFile(lib.DALBootstrapFile, formatted, lib.DefaultFileMode); e != nil {
		fmt.Println("Write file error: ", e.Error())
		return e
	}

	// TODO verbose flag
	// fmt.Printf("Generated dal bootstrap file to %s in %f seconds\n", lib.DALBootstrapFile, time.Since(start).Seconds())

	return nil
}
