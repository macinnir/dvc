package gen

import (
	"bufio"
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
	"time"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/cache"
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

// CleanGoDALs removes any repo files that aren't in the database.Tables map
func CleanGoDALs(dir string, database *schema.Schema) (e error) {

	// EXT
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

	reader := bufio.NewReader(os.Stdin)

	for _, name := range dirFileNames {

		// Skip anything that doesn't have the go extension
		if len(name) < 4 || name[len(name)-3:] != ".go" {
			continue
		}

		// Remove the extension
		modelName := name[0 : len(name)-3]

		// Skip tests
		if len(modelName) > 5 && modelName[len(modelName)-5:] == "_test" {
			continue
		}

		if modelName == "bootstrap" {
			continue
		}

		// DALExt
		if len(modelName) > 6 && modelName[len(modelName)-6:] == "DALExt" {
			fmt.Println("Ext file: ", name)
			modelName = modelName[:len(modelName)-6] // DALExt
		} else {
			modelName = modelName[:len(modelName)-3] // DAL
		}

		if _, ok := database.Tables[modelName]; !ok {
			if modelName != "Config" {
				fullFilePath := path.Join(dir, name)
				if result := lib.ReadCliInput(reader, fmt.Sprintf("Delete unused dal `%s`(Y/n)?", name)); result == "Y" {
					fmt.Printf("Removing %s\n", fullFilePath)
					os.Remove(fullFilePath)
				}
			}
		}
	}
	return
}

func toArgName(field string) string {
	return strings.ToLower(field[:1]) + field[1:]
}

func GenDALs(dalsDir string, config *lib.Config, force, clean bool) error {

	lib.EnsureDir(dalsDir)

	start := time.Now()
	var schemaList *schema.SchemaList
	var e error

	schemaList, e = schema.LoadLocalSchemas()

	if e != nil {
		return e
	}

	fmt.Println("Generating models...")

	if clean {
		for k := range schemaList.Schemas {
			CleanGoDALs(dalsDir, schemaList.Schemas[k])
		}
	}

	var tablesCache cache.TablesCache
	tablesCache, e = cache.LoadTableCache()

	if e != nil {
		return e
	}

	generatedDALCount := 0

	for k := range schemaList.Schemas {

		schemaName := schemaList.Schemas[k].Name

		for l := range schemaList.Schemas[k].Tables {

			table := schemaList.Schemas[k].Tables[l]
			tableKey := schemaName + "_" + table.Name

			var tableHash string
			tableHash, e = cache.HashTable(table)

			// If the model has been hashed before...
			if _, ok := tablesCache.Models[tableKey]; ok {

				// And the hash hasn't changed, skip...
				if tableHash == tablesCache.Models[tableKey] && !force {
					// fmt.Printf("Table `%s` hasn't changed! Skipping...\n", table.Name)
					continue
				}
			}

			generatedDALCount++

			// Update the dals cache
			if e = GenerateGoDAL(config, table, dalsDir); e != nil {
				return e
			}

			tablesCache.Models[tableKey] = tableHash
		}
	}

	cache.SaveTableCache(tablesCache)

	fmt.Println("Generating bootstrap file...")
	if e = GenerateDALsBootstrapFile(config, dalsDir, schemaList); e != nil {
		return e
	}

	fmt.Printf("Generated %d dals in %f seconds.\n", generatedDALCount, time.Since(start).Seconds())

	return nil

}

var dalTPL *template.Template

// GenerateGoDAL returns a string for a repo in golang
func GenerateGoDAL(config *lib.Config, table *schema.Table, dir string) (e error) {

	p := path.Join(dir, table.Name+"DAL.go")
	// TODO verbose flag
	// start := time.Now()
	// TODO verbose flag
	lib.EnsureDir(dir)

	imports := []string{}

	// lib.Debugf("Generating go dal file for table %s at path %s", g.Options, table.Name, p)

	data := struct {
		Table             *schema.Table
		Columns           schema.SortedColumns
		UpdateColumns     []*schema.Column
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
	}{
		Table:             table,
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
	}

	// if data.FileHead, data.FileFoot, imports, e = g.scanFileParts(p, true); e != nil {
	// 	lib.Errorf("ERROR: %s", g.Options, e.Error())
	// 	return
	// }

	// funcSig := fmt.Sprintf(`^func \(r \*%sRepo\) [A-Z].*$`, table.Name)
	// footMatches := g.scanStringForFuncSignature(fileFoot, funcSig)

	sortedColumns := make(schema.SortedColumns, 0, len(table.Columns))

	hasNull := false

	// Find the primary key
	for _, column := range table.Columns {
		if column.ColumnKey == "PRI" {
			data.PrimaryKey = column.Name
			data.PrimaryKeyType = column.DataType
		}

		goDataType := schema.DataTypeToGoTypeString(column)
		if len(goDataType) > 5 && goDataType[0:5] == "null." {
			hasNull = true
		}

		sortedColumns = append(sortedColumns, column)
	}

	sort.Sort(sortedColumns)
	data.Columns = sortedColumns

	insertColumnNames := []string{}
	insertColumnVals := []string{}
	insertColumnArgs := []string{}

	insertColumns := fetchInsertColumns(sortedColumns)

	for _, col := range insertColumns {
		insertColumnNames = append(insertColumnNames, fmt.Sprintf("`%s`", col.Name))
		insertColumnVals = append(insertColumnVals, "?")
		insertColumnArgs = append(insertColumnArgs, fmt.Sprintf("model.%s", col.Name))
	}

	data.InsertArgs = strings.Join(insertColumnArgs, ",")
	data.InsertSQL = fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", data.Table.Name, strings.Join(insertColumnNames, ","), strings.Join(insertColumnVals, ","))

	data.UpdateColumns = fetchUpdateColumns(sortedColumns)
	updateColumnNames := []string{}
	updateColumnArgs := []string{}
	for _, col := range data.UpdateColumns {
		updateColumnNames = append(updateColumnNames, fmt.Sprintf("`%s` = ?", col.Name))
		updateColumnArgs = append(updateColumnArgs, fmt.Sprintf("model.%s", col.Name))
	}

	updateColumnArgs = append(updateColumnArgs, fmt.Sprintf("model.%s", data.PrimaryKey))

	data.UpdateSQL = fmt.Sprintf("UPDATE `%s` SET %s WHERE %s = ?", data.Table.Name, strings.Join(updateColumnNames, ","), data.PrimaryKey)
	data.UpdateArgs = strings.Join(updateColumnArgs, ",")

	_, data.IsDeleted = table.Columns["IsDeleted"]
	_, data.IsDateCreated = table.Columns["DateCreated"]
	_, data.IsLastUpdated = table.Columns["LastUpdated"]

	switch data.PrimaryKeyType {
	case "varchar":
		data.IDType = "string"
	}

	defaultImports := []string{
		fmt.Sprintf("%s/%s", config.BasePackage, "gen/definitions/models"),
		"github.com/macinnir/dvc/core/lib/utils/db",
		"github.com/macinnir/dvc/core/lib/utils/log",
		"github.com/macinnir/dvc/core/lib/utils/errors",
		"database/sql",
		"context",
		"fmt",
		"strings",
	}

	if hasNull {
		imports = append(imports, "gopkg.in/guregu/null.v3")
	}

	// If either of the fields "DateCreated" or "LastUpdated" exist on this model,
	// the `time` package is needed
	if data.IsDateCreated || data.IsLastUpdated {
		imports = append(imports, "time")
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

	// var {{.Table.Name}}DALFields = []string{
	// 	{{range $col := .Columns}}"{{$col.Name}}",
	// 	{{end}}
	// }

	tpl := `// Generated Code; DO NOT EDIT.

package dal

import ({{range .Imports}}
	"{{.}}"{{end}}
)

// {{.Table.Name}}DAL is a data repository for {{.Table.Name}} objects
type {{.Table.Name}}DAL struct {
	db  []db.IDB
	log log.ILog
}

// New{{.Table.Name}}DAL returns a new instance of {{.Table.Name}}Repo
func New{{.Table.Name}}DAL(db []db.IDB, log log.ILog) *{{.Table.Name}}DAL {
	return &{{.Table.Name}}DAL{db, log}
}

func (r *{{.Table.Name}}DAL) Raw(shard int64, q string, args ...interface{}) ([]*models.{{.Table.Name}}, error) { 
	return (&models.{{.Table.Name}}{}).Raw(r.db[shard], fmt.Sprintf(q, args...)) 
}

func (r *{{.Table.Name}}DAL) Select(shard int64) *models.{{.Table.Name}}DALSelector { 
	return (&models.{{.Table.Name}}{}).Select(r.db[shard])
}

func (r *{{.Table.Name}}DAL) Count(shard int64) *models.{{.Table.Name}}DALCounter { 
	return (&models.{{.Table.Name}}{}).Count(r.db[shard])
}

func (r *{{.Table.Name}}DAL) Get(shard int64) *models.{{.Table.Name}}DALGetter { 
	return (&models.{{.Table.Name}}{}).Get(r.db[shard])
}

// Create creates a new {{.Table.Name}} entry in the database
func (r *{{.Table.Name}}DAL) Create(shard int64, model *models.{{.Table.Name}}) error { {{if .IsDateCreated}}
	
	model.DateCreated = time.Now().UnixNano() / 1000000{{end}}
	{{if .IsLastUpdated}}
	model.LastUpdated = time.Now().UnixNano() / 1000000
	{{end}}
	e := model.Create(r.db[shard])
	if e != nil {
		r.log.Errorf("{{.Table.Name}}DAL.Insert > %s", e.Error())
		return e		
	}

	r.log.Debugf("{{.Table.Name}}DAL.Insert(%d)", model.{{.PrimaryKey}})

	return nil
}

// CreateMany creates {{.Table.Name}} objects in chunks
func (r *{{.Table.Name}}DAL) CreateMany(shard int64, modelSlice []*models.{{.Table.Name}}) (e error) {

	// No records 
	if len(modelSlice) == 0 {
		return 
	}

	// Don't use a transaction if only a single value
	if len(modelSlice) == 1 {
		e = r.Create(shard, modelSlice[0])
		return
	}

	chunkSize := 25
	chunks := [][]*models.{{.Table.Name}}{}

	for i := 0; i < len(modelSlice); i += chunkSize {
		end := i + chunkSize
		if end > len(modelSlice) {
			end = len(modelSlice)
		}
		chunks = append(chunks, modelSlice[i:end])
	}

	for chunkID, chunk := range chunks {

		var tx *sql.Tx
		ctx := context.Background()
		tx, e = r.db[shard].BeginTx(ctx, nil)
		if e != nil {
			return
		}

		for insertID, model := range chunk {

			{{if .IsDateCreated}}
			model.DateCreated = time.Now().UnixNano() / 1000000{{end}}
			{{if .IsLastUpdated}}
			model.LastUpdated = time.Now().UnixNano() / 1000000{{end}}

			_, e = tx.ExecContext(ctx, "{{.InsertSQL}}", {{.InsertArgs}})
			if e != nil {
				r.log.Errorf("{{.Table.Name}}.CreateMany([](%d)) (Chunk %d.%d) > %s", len(modelSlice), chunkID, insertID, e.Error())
				break
			} else {
				r.log.Debugf("{{.Table.Name}}.CreateMany([](%d)) (Chunk %d.%d)", len(modelSlice), chunkID, insertID)
			}
		}

		if e != nil {
			return
		}

		e = tx.Commit()
	}

	return

}

// Update updates an existing {{.Table.Name}} entry in the database
func (r *{{.Table.Name}}DAL) Update(shard int64, model *models.{{.Table.Name}}) (e error) {
{{if .IsLastUpdated}}
	model.LastUpdated = time.Now().UnixNano() / 1000000{{end}}
	_, e = r.db[shard].Exec("{{.UpdateSQL}}", {{.UpdateArgs}})
	if e != nil {
		r.log.Errorf("{{.Table.Name}}DAL.Update(%d) > %s", model.{{.PrimaryKey}}, e.Error())
	} else {
		r.log.Debugf("{{.Table.Name}}DAL.Update(%d)", model.{{.PrimaryKey}})
	}
	return
}

// UpdateMany updates a slice of {{.Table.Name}} objects in chunks
func (r {{.Table.Name}}DAL) UpdateMany(shard int64, modelSlice []*models.{{.Table.Name}}) (e error) {

	// No records 
	if len(modelSlice) == 0 {
		return 
	}

	// Don't use a transaction if only a single value
	if len(modelSlice) == 1 {
		e = r.Update(shard, modelSlice[0])
		return
	}

	chunkSize := 25
	chunks := [][]*models.{{.Table.Name}}{}

	for i := 0; i < len(modelSlice); i += chunkSize {
		end := i + chunkSize
		if end > len(modelSlice) {
			end = len(modelSlice)
		}
		chunks = append(chunks, modelSlice[i:end])
	}

	for chunkID, chunk := range chunks {

		var tx *sql.Tx
		ctx := context.Background()
		tx, e = r.db[shard].BeginTx(ctx, nil)
		if e != nil {
			return
		}

		for updateID, model := range chunk {
{{if .IsLastUpdated}}
			model.LastUpdated = time.Now().UnixNano() / 1000000{{end}}

			_, e = tx.ExecContext(ctx, "{{.UpdateSQL}}", {{.UpdateArgs}})
			if e != nil {
				r.log.Errorf("{{.Table.Name}}.UpdateMany([](%d)) (Chunk %d.%d) > %s", len(modelSlice), chunkID, updateID, e.Error())
				break
			} else {
				r.log.Debugf("{{.Table.Name}}.UpdateMany([](%d)) (Chunk %d.%d)", len(modelSlice), chunkID, updateID)
			}
		}

		if e != nil {
			return
		}

		e = tx.Commit()
	}

	return

}{{if .IsDeleted}}

// Delete marks an existing {{.Table.Name}} entry in the database as deleted
func (r *{{.Table.Name}}DAL) Delete(shard int64, {{.PrimaryKey | toArgName}} {{.IDType}}) (e error) {
	_, e = r.db[shard].Exec("UPDATE ` + "`{{.Table.Name}}` SET `IsDeleted` = 1 WHERE `{{.PrimaryKey}}` = ?" + `", {{.PrimaryKey | toArgName}})
	if e != nil {
		r.log.Errorf("{{.Table.Name}}DAL.Delete(%d) > %s", {{.PrimaryKey | toArgName}}, e.Error())
	} else {
		r.log.Debugf("{{.Table.Name}}DAL.Delete(%d)", {{.PrimaryKey | toArgName}})
	}
	return
}

// DeleteMany marks {{.Table.Name}} objects in chunks as deleted
func (r {{.Table.Name}}DAL) DeleteMany(shard int64, modelSlice []*models.{{.Table.Name}}) (e error) {

	// No records 
	if len(modelSlice) == 0 {
		return 
	}

	// Don't use a transaction if only a single value
	if len(modelSlice) == 1 {
		e = r.Delete(shard, modelSlice[0].{{.PrimaryKey}})
		return
	}

	chunkSize := 25
	chunks := [][]*models.{{.Table.Name}}{}

	for i := 0; i < len(modelSlice); i += chunkSize {
		end := i + chunkSize
		if end > len(modelSlice) {
			end = len(modelSlice)
		}
		chunks = append(chunks, modelSlice[i:end])
	}

	for chunkID, chunk := range chunks {

		var tx *sql.Tx
		ctx := context.Background()
		tx, e = r.db[shard].BeginTx(ctx, nil)
		if e != nil {
			return
		}

		for deleteID, model := range chunk {
{{if .IsLastUpdated}}
			model.LastUpdated = time.Now().UnixNano() / 1000000{{end}}
			_, e = tx.ExecContext(ctx, "UPDATE ` + "`{{.Table.Name}}` SET `IsDeleted`= 1 WHERE `{{.PrimaryKey}}` = ?" + `", model.{{.PrimaryKey}})
			if e != nil {
				r.log.Errorf("{{.Table.Name}}.DeleteMany([](%d)) (Chunk %d.%d) > %s", len(modelSlice), chunkID, deleteID, e.Error())
				break
			} else {
				r.log.Debugf("{{.Table.Name}}.DeleteMany([](%d)) (Chunk %d.%d)", len(modelSlice), chunkID, deleteID)
			}
		}

		if e != nil {
			return
		}

		e = tx.Commit()
	}

	return

}{{end}}

// DeleteHard performs a SQL DELETE operation on a {{.Table.Name}} entry in the database
func (r *{{.Table.Name}}DAL) DeleteHard(shard int64, {{.PrimaryKey | toArgName}} {{.IDType}}) (e error) {
	_, e = r.db[shard].Exec("DELETE FROM ` + "`{{.Table.Name}}`" + ` WHERE {{.PrimaryKey}} = ?", {{.PrimaryKey | toArgName}})
	if e != nil {
		r.log.Errorf("{{.Table.Name}}DAL.HardDelete(%d) > %s", {{.PrimaryKey | toArgName}}, e.Error())
	} else {
		r.log.Debugf("{{.Table.Name}}DAL.HardDelete(%d)", {{.PrimaryKey | toArgName}})
	}
	return
}

// DeleteManyHard deletes {{.Table.Name}} objects in chunks
func (r {{.Table.Name}}DAL) DeleteManyHard(shard int64, modelSlice []models.{{.Table.Name}}) (e error) {

	// No records 
	if len(modelSlice) == 0 {
		return 
	}

	// Don't use a transaction if only a single value
	if len(modelSlice) == 1 {
		e = r.DeleteHard(shard, modelSlice[0].{{.PrimaryKey}})
		return
	}

	chunkSize := 25
	chunks := [][]models.{{.Table.Name}}{}

	for i := 0; i < len(modelSlice); i += chunkSize {
		end := i + chunkSize
		if end > len(modelSlice) {
			end = len(modelSlice)
		}
		chunks = append(chunks, modelSlice[i:end])
	}

	for chunkID, chunk := range chunks {

		var tx *sql.Tx
		ctx := context.Background()
		tx, e = r.db[shard].BeginTx(ctx, nil)
		if e != nil {
			return
		}

		for deleteID, model := range chunk {

			_, e = tx.ExecContext(ctx, "DELETE FROM ` + "`{{.Table.Name}}` WHERE `{{.PrimaryKey}}` = ?" + `", model.{{.PrimaryKey}})
			if e != nil {
				r.log.Errorf("{{.Table.Name}}.DeleteManyHard([](%d)) (Chunk %d.%d) > %s", len(modelSlice), chunkID, deleteID, e.Error())
				break
			} else {
				r.log.Debugf("{{.Table.Name}}.DeleteManyHard([](%d)) (Chunk %d.%d)", len(modelSlice), chunkID, deleteID)
			}
		}

		if e != nil {
			return
		}

		e = tx.Commit()
	}

	return
}

// FromID gets a single {{.Table.Name}} object by its Primary Key
func (r *{{.Table.Name}}DAL) FromID(shard int64, {{.PrimaryKey | toArgName}} {{.IDType}}, mustExist bool) (model *models.{{.Table.Name}}, e error) {

	model = &models.{{.Table.Name}}{}

	e = r.db[shard].Get(model, "SELECT * FROM ` + "`{{.Table.Name}}` WHERE `{{.PrimaryKey}}` = ?" + `", {{.PrimaryKey | toArgName}})

	if e == nil {
		
		{{ if .IsDeleted}}if model.IsDeleted == 1 && mustExist {
			model = nil
			e = errors.NewRecordNotFoundError()
			return
		}{{end}}

		r.log.Debugf("{{.Table.Name}}DAL.FromID(%d)", model.{{.PrimaryKey}})
		
	} else if e == sql.ErrNoRows {
		r.log.Debugf("{{.Table.Name}}DAL.FromID(%d) > NOT FOUND", {{.PrimaryKey | toArgName}})

		if mustExist {
			e = errors.NewRecordNotFoundError()
			return
		}

		e = nil
		model = nil
	} else {
		
		r.log.Errorf("{{.Table.Name}}DAL.FromID(%d) > %s", {{.PrimaryKey | toArgName}}, e.Error())
	}

	return
}

// FromIDs returns a slice of {{.Table.Name}} objects by a set of primary keys
func (r *{{.Table.Name}}DAL) FromIDs(shard int64, {{.PrimaryKey | toArgName}}s []{{.IDType}}) (model []*models.{{.Table.Name}}, e error) {

	model = []*models.{{.Table.Name}}{}
	
	// No records 
	if len({{.PrimaryKey | toArgName}}s) == 0 {
		return 
	}

	ids := []string{}
	for _, id := range {{.PrimaryKey | toArgName}}s {
		ids = append(ids, fmt.Sprintf("%d", id))
	}

	query := fmt.Sprintf("SELECT * FROM ` + "`{{.Table.Name}}` WHERE `{{.PrimaryKey}}` IN (%s){{if .IsDeleted}} AND IsDeleted = 0{{end}}" + `", strings.Join(ids, ","))

	e = r.db[shard].Select(&model, query)

	if e == nil {
		r.log.Debugf("{{.Table.Name}}DAL.FromIDs(%v)", {{.PrimaryKey | toArgName}}s)
	} else {
		r.log.Errorf("{{.Table.Name}}DAL.FromIDs(%v) > %s", {{.PrimaryKey | toArgName}}s, e.Error())
	}

	return
}

{{range $col := .UpdateColumns}}
// Set{{$col.Name}} sets the {{$col.Name}} column on a {{$.Table.Name}} object
func (r *{{$.Table.Name}}DAL) Set{{$col.Name}}(shard int64, {{$.PrimaryKey | toArgName}} {{$.IDType}}, {{$col.Name | toArgName}} {{$col | dataTypeToGoTypeString}}) (e error) {
	_, e = r.db[shard].Exec("UPDATE ` + "`{{$.Table.Name}}` SET `{{$col.Name}}` = ? WHERE `{{$.PrimaryKey}}` = ?" + `", {{$col.Name | toArgName}}, {{$.PrimaryKey | toArgName}})
	if e != nil {
		r.log.Errorf("{{$.Table.Name}}DAL.Set{{$col.Name}}(%d, %v) > %s", {{$.PrimaryKey | toArgName}}, {{$col.Name | toArgName}}, e.Error())
	} else {
		r.log.Debugf("{{$.Table.Name}}DAL.Set{{$col.Name}}(%d, %v)", {{$.PrimaryKey | toArgName}}, {{$col.Name | toArgName}})
	}
	return
}

// ManyFrom{{$col.Name}} returns a slice of {{$.Table.Name}} models from {{$col.Name}}
func (r *{{$.Table.Name}}DAL) ManyFrom{{$col.Name}}(shard int64, {{$col.Name | toArgName}} {{$col | dataTypeToGoTypeString}}, limit, offset int64, orderBy, orderDir string) (collection []*models.{{$.Table.Name}}, e error) {
	
	collection = []*models.{{$.Table.Name}}{}

	orderDirString := "ASC"
	if orderDir == "DESC" {
		orderDirString = "DESC"
	}

	query := "SELECT * FROM ` + "`{{$.Table.Name}}` WHERE `{{$col.Name}}` = ? " + `"
	{{if $.IsDeleted}}
	query += "` + " AND `IsDeleted` = 0" + `"{{end}}

	// Optional Order By 
	if len(orderBy) > 0 {
		query += fmt.Sprintf(" ORDER BY ` + "`%s` %s" + `", orderBy, orderDirString)
	}

	// Optional Limit 
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	e = r.db[shard].Select(&collection, query, {{$col.Name | toArgName}})

	if e != nil {
		r.log.Errorf("{{$.Table.Name}}DAL.ManyFrom{{$col.Name}}({{$col | dataTypeToFormatString}}, %d, %d, %s, %s) > %s", {{$col.Name | toArgName}}, limit, offset, orderBy, orderDir, e.Error())
	} else {
		r.log.Debugf("{{$.Table.Name}}DAL.ManyFrom{{$col.Name}}(%d, %d, %s, %s)", limit, offset, orderBy, orderDir)
	}
	return
}

{{if not $col.IsString}} // This is a string??{{end}}

{{if not $col.IsString}}
// ManyFrom{{$col.Name}}s returns a slice of {{$.Table.Name}} models from {{$col.Name}}s
func (r *{{$.Table.Name}}DAL) ManyFrom{{$col.Name}}s(shard int64, {{$col.Name | toArgName}}s []{{$col | dataTypeToGoTypeString}}, limit, offset int64, orderBy, orderDir string) (collection []*models.{{$.Table.Name}}, e error) {
	
	collection = []*models.{{$.Table.Name}}{}
	
	// No records 
	if len({{$col.Name | toArgName}}s) == 0 {
		return 
	}

	idGate := map[{{$col | dataTypeToGoTypeString}}]bool{} 
	ids := []string{}
	for _, id := range {{$col.Name | toArgName}}s {
		if _, ok := idGate[id]; !ok { 
			idGate[id] = true 
			ids = append(ids, fmt.Sprintf("{{$col.FmtType}}", id))
		}
	}

	orderDirString := "ASC"
	if orderDir == "DESC" {
		orderDirString = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM ` + "`{{$.Table.Name}}` WHERE `{{$col.Name}}` IN (%s) {{if $.IsDeleted}} AND `IsDeleted` = 0{{end}}" + `", strings.Join(ids, ",")) 

	// Optional Order By 
	if len(orderBy) > 0 {
		query += fmt.Sprintf(" ORDER BY ` + "`%s` %s" + `", orderBy, orderDirString)
	}

	// Optional Limit 
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	e = r.db[shard].Select(&collection, query)

	if e != nil {
		r.log.Errorf("{{$.Table.Name}}DAL.ManyFrom{{$col.Name}}s(%v, %d, %d, %s, %s) > %s", {{$col.Name | toArgName}}s, limit, offset, orderBy, orderDir, e.Error())
	} else {
		r.log.Debugf("{{$.Table.Name}}DAL.ManyFrom{{$col.Name}}s(%d, %d, %s, %s)", limit, offset, orderBy, orderDir)
	}
	return
}
{{end}}

// CountFrom{{$col.Name}} returns the number of {{$.Table.Name}} records from {{$col.Name}}
func (r *{{$.Table.Name}}DAL) CountFrom{{$col.Name}}(shard int64, {{$col.Name | toArgName}} {{$col | dataTypeToGoTypeString}}) (count int64, e error) {
	
	count = 0

	query := "SELECT COUNT(*) FROM ` + "`{{$.Table.Name}}`" + ` WHERE {{$col.Name}} = ?"

	{{if $.IsDeleted}}
	query += "` + " AND `IsDeleted` = 0" + `"{{end}}

	e = r.db[shard].Get(&count, query, {{$col.Name | toArgName}})
	
	if e != nil {
		r.log.Errorf("{{$.Table.Name}}DAL.CountFrom{{$col.Name}}({{$col | dataTypeToFormatString}}) > %s", {{$col.Name | toArgName}}, e.Error())
	} else {
		r.log.Debugf("{{$.Table.Name}}DAL.CountFrom{{$col.Name}}({{$col | dataTypeToFormatString}})", {{$col.Name | toArgName}})
	}

	return
}

// SingleFrom{{$col.Name}} returns a single {{$.Table.Name}} record by its {{$col.Name}}
func (r *{{$.Table.Name}}DAL) SingleFrom{{$col.Name}}(shard int64, {{$col.Name | toArgName}} {{$col | dataTypeToGoTypeString}}, mustExist bool) (model *models.{{$.Table.Name}}, e error) {

	model = &models.{{$.Table.Name}}{}

	query := "SELECT * FROM ` + "`{{$.Table.Name}}` WHERE `{{$col.Name}}` = ?" + `"
	{{if $.IsDeleted}}
	query += "` + " AND `IsDeleted` = 0" + `"{{end}}
	
	e = r.db[shard].Get(model, query, {{$col.Name | toArgName}})

	if e == nil {

		r.log.Debugf("{{$.Table.Name}}DAL.SingleFrom{{$col.Name}}({{$col | dataTypeToFormatString}})", {{$col.Name | toArgName}})

	} else if e == sql.ErrNoRows {

		r.log.Debugf("{{$.Table.Name}}DAL.SingleFrom{{$col.Name}}({{$col | dataTypeToFormatString}}) > NOT FOUND", {{$col.Name | toArgName}})

		if mustExist {
			e = errors.NewRecordNotFoundError()
			return
		}

		e = nil
		model = nil

	} else {
		r.log.Errorf("{{$.Table.Name}}DAL.SingleFrom{{$col.Name}}({{$col | dataTypeToFormatString}}) > %s", {{$col.Name | toArgName}}, e.Error())
	}

	return
}


{{end}}

// ManyPaged returns a slice of {{.Table.Name}} models
func (r *{{.Table.Name}}DAL) ManyPaged(shard int64, limit, offset int64, orderBy, orderDir string) (collection []*models.{{.Table.Name}}, e error) {
	
	collection = []*models.{{.Table.Name}}{}
	
	orderDirString := "ASC"
	if orderDir == "DESC" {
		orderDirString = "DESC"
	}

	query := "SELECT * FROM ` + "`{{.Table.Name}}` WHERE 1=1" + `"

	{{if $.IsDeleted}}
	query += "` + " AND `IsDeleted` = 0" + `"{{end}}
	
	// Optional Order By 
	if len(orderBy) > 0 {
		query += fmt.Sprintf(" ORDER BY ` + "`%s` %s" + `", orderBy, orderDirString)
	}

	// Optional Limit 
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	e = r.db[shard].Select(&collection, query)

	if e != nil {
		r.log.Errorf("{{.Table.Name}}DAL.GetMany(%d, %d, %s, %s) > %s", limit, offset, orderBy, orderDir, e.Error())
		e = fmt.Errorf("ManyPaged: Select %w", e) 
	} else {
		r.log.Debugf("{{.Table.Name}}DAL.GetMany(%d, %d, %s, %s)", limit, offset, orderBy, orderDir)
	}
	return
}`

	if dalTPL == nil {

		dalTPL = template.New("dal")

		dalTPL.Funcs(template.FuncMap{
			"insertFields":           fetchTableInsertFieldsString,
			"insertValues":           fetchTableInsertValuesString,
			"updateFields":           fetchTableUpdateFieldsString,
			"dataTypeToGoTypeString": schema.DataTypeToGoTypeString,
			"dataTypeToFormatString": schema.DataTypeToFormatString,
			"toArgName":              toArgName,
		})

		dalTPL, e = dalTPL.Parse(tpl)
		if e != nil {
			panic(e)
		}
	}

	f, err := os.Create(p)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}

	err = dalTPL.Execute(f, data)
	if err != nil {
		fmt.Println("Execute Error: ", err.Error())
		return
	}

	f.Close()

	// if e = lib.FmtGoCode(p); e != nil {
	// 	return e
	// 	// lib.Warn(e.Error(), g.Options)
	// }

	// TODO verbose flag
	// fmt.Printf("%f seconds\n", time.Since(start).Seconds())

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

func fetchTablePrimaryKey(table *schema.Table) string {
	primaryKey := ""
	idType := "int64"
	for _, column := range table.Columns {
		if column.ColumnKey == "PRI" {
			primaryKey = column.Name
		}
	}

	return primaryKey + " " + idType
}

func fetchTablePrimaryKeyName(table *schema.Table) string {
	primaryKey := ""
	for _, column := range table.Columns {
		if column.ColumnKey == "PRI" {
			primaryKey = column.Name
		}
	}

	return primaryKey
}

func fetchTableInsertFieldsString(columns schema.SortedColumns) string {

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

func fetchTableInsertValuesString(columns schema.SortedColumns) string {
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

func isInsertColumn(column *schema.Column) bool {
	if column.ColumnKey == "PRI" {
		return false
	}

	if column.Name == "IsDeleted" {
		return false
	}

	return true
}

func fetchInsertColumns(columns schema.SortedColumns) []*schema.Column {

	insertColumns := []*schema.Column{}

	for _, column := range columns {
		if !isInsertColumn(column) {
			continue
		}

		insertColumns = append(insertColumns, column)
	}

	return insertColumns
}

func isUpdateColumn(column *schema.Column) bool {
	if column.ColumnKey == "PRI" {
		return false
	}

	// if column.Name == "IsDeleted" {
	// 	return false
	// }

	if column.Name == "DateCreated" {
		return false
	}

	return true
}

func fetchUpdateColumns(columns schema.SortedColumns) []*schema.Column {

	updateColumns := []*schema.Column{}

	for _, column := range columns {
		if !isUpdateColumn(column) {
			continue
		}

		updateColumns = append(updateColumns, column)
	}

	return updateColumns
}

func fetchTableUpdateFieldsString(columns schema.SortedColumns) string {
	fields := []string{}
	for _, field := range columns {

		if !isUpdateColumn(field) {
			continue
		}

		fields = append(fields, "`"+field.Name+"` = :"+field.Name)
	}

	return strings.Join(fields, ",")
}

// GenerateDALsBootstrapFile generates a dal bootstrap file in golang
func GenerateDALsBootstrapFile(config *lib.Config, dir string, schemaList *schema.SchemaList) (e error) {

	tables := map[string]*schema.Table{}

	for k := range schemaList.Schemas {
		for l := range schemaList.Schemas[k].Tables {
			tables[l] = schemaList.Schemas[k].Tables[l]
		}
	}

	// Make the dir if it does not exist.
	lib.EnsureDir(dir)

	data := struct {
		Tables        map[string]*schema.Table
		BasePackage   string
		DBPackage     string
		LogPackage    string
		ModelsPackage string
	}{
		BasePackage:   config.BasePackage,
		Tables:        tables,
		DBPackage:     "github.com/macinnir/dvc/core/lib/utils/db",
		LogPackage:    "github.com/macinnir/dvc/core/lib/utils/log",
		ModelsPackage: fmt.Sprintf("%s/%s", config.BasePackage, "gen/definitions/models"),
	}

	tpl := `
// Package dal is the Data Access Layer
package dal

import (
	"{{ .DBPackage }}"
	"{{ .LogPackage }}"
	"{{ .ModelsPackage }}"
)

// DAL is a container for all dal structs
type DAL struct {
	{{range .Tables}}
	{{.Name}} *{{.Name}}DAL{{end}}
}

// BootstrapDAL bootstraps all of the DAL methods
func BootstrapDAL(db map[string][]db.IDB, log log.ILog) *DAL {

	d := &DAL{}
	{{range .Tables}}
	d.{{.Name}} = New{{.Name}}DAL(db[models.{{.Name}}_SchemaName], log){{end}}

	return d
}`

	p := path.Join(dir, "bootstrap.go")
	// lib.Debugf("Generating dal bootstrap file at path %s", g.Options, p)
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
