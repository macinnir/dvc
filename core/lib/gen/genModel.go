package gen

import (
	"bufio"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/cache"
	"github.com/macinnir/dvc/core/lib/schema"
)

// GenModels generates models
func GenModels(config *lib.Config, force bool, clean bool) error {

	start := time.Now()
	var schemaList *schema.SchemaList
	var e error

	schemaList, e = schema.LoadLocalSchemas()

	if e != nil {
		return e
	}

	fmt.Println("Generating models...")

	if clean {

		tableMap := map[string]struct{}{}
		for k := range schemaList.Schemas {
			schema := schemaList.Schemas[k]
			for l := range schema.Tables {
				tableMap[schema.Tables[l].Name] = struct{}{}
			}
		}

		CleanGoModels(lib.ModelsGenDir, tableMap)

	}

	var tablesCache cache.TablesCache
	tablesCache, e = cache.LoadTableCache()

	if e != nil {
		return e
	}

	generatedModelCount := 0

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

			generatedModelCount++

			// Update the models cache
			tablesCache.Models[tableKey] = tableHash

			// TODO verbose flag
			fmt.Printf("Generating model `%s`\n", table.Name)
			e = GenerateGoModel(config.BasePackage, "gen/definitions/models", schemaName, table)
			if e != nil {
				return e
			}
		}
	}

	cache.SaveTableCache(tablesCache)

	fmt.Printf("Generated %d models in %f seconds.\n", generatedModelCount, time.Since(start).Seconds())

	return nil
}

// GenerateGoModel returns a string for a model in golang
func GenerateGoModel(packageName, dir string, schemaName string, table *schema.Table) (e error) {

	lib.EnsureDir(dir)

	fullPath := path.Join(dir, table.Name+".go")

	// if lib.FileExists(p) {
	// 	lib.Infof("Updating `%s`", g.Options, table.Name)
	// 	e = g.updateGoModel(p, table)
	// 	return
	// }

	fmt.Println("Creating", table.Name)
	// lib.Infof("Creating `%s`", g.Options, table.Name)
	e = buildGoModel(packageName, fullPath, schemaName, table)
	return
}

// InspectFile inspects a file
func InspectFile(filePath string) (s *lib.GoStruct, e error) {

	fileBytes, e := ioutil.ReadFile(filePath)
	if e != nil {
		panic(e)
	}

	s, e = buildGoStructFromFile(fileBytes)
	if e != nil {
		fmt.Println("ERROR: ", filePath)
		panic(e)
	}

	return

}

// func (g *Gen) updateGoModel(p string, table *schema.Table) (e error) {

// 	var modelNode *lib.GoStruct
// 	var outFile []byte

// 	fileBytes, e := ioutil.ReadFile(p)
// 	modelNode, e = buildGoStructFromFile(fileBytes)
// 	resolveTableToModel(modelNode, table)

// 	if e != nil {
// 		return
// 	}

// 	outFile, e = buildFileFromModelNode(modelNode)
// 	if e != nil {
// 		return
// 	}

// 	e = ioutil.WriteFile(p, outFile, 0644)

// 	return
// }

func buildGoModel(packageName, fullPath string, schemaName string, table *schema.Table) (e error) {
	var modelNode *lib.GoStruct
	var outFile []byte

	fmt.Println("buildModelNodeFromTable")
	modelNode, e = buildModelNodeFromTable(packageName, table)
	if e != nil {
		fmt.Println("ERROR Building Model Node From Table ", table)
		return
	}

	fmt.Println("buildFileFromModelNode")
	outFile, e = buildFileFromModelNode(schemaName, table, modelNode)
	if e != nil {
		fmt.Println("ERROR Building File From Model Node ", table)
		return
	}

	fmt.Println("Writing file to", fullPath)
	e = ioutil.WriteFile(fullPath, outFile, 0644)
	return
}

// buildModelNodeFromFile builds a node representation of a struct from a file
func buildModelNodeFromTable(packageName string, table *schema.Table) (modelNode *lib.GoStruct, e error) {

	modelNode = lib.NewGoStruct()
	modelNode.Package = "models"
	modelNode.Name = table.Name
	modelNode.Comments = fmt.Sprintf("%s is a `%s` data model\n", table.Name, table.Name)
	modelNode.Imports.Append("\"github.com/macinnir/dvc/core/lib/utils/query\"")
	modelNode.Imports.Append("\"github.com/macinnir/dvc/core/lib/utils/db\"")
	modelNode.Imports.Append("\"encoding/json\"")
	modelNode.Imports.Append("\"fmt\"")
	modelNode.Imports.Append("\"database/sql\"")

	hasNull := false

	sortedColumns := make(schema.SortedColumns, 0, len(table.Columns))

	for _, column := range table.Columns {
		sortedColumns = append(sortedColumns, column)
	}

	sort.Sort(sortedColumns)

	for _, col := range sortedColumns {
		fieldType := schema.DataTypeToGoTypeString(col)
		if fieldDataTypeIsNull(fieldType) {
			hasNull = true
		}
		modelNode.Fields.Append(&lib.GoStructField{
			Name:     col.Name,
			DataType: fieldType,
			Tags: []*lib.GoStructFieldTag{
				{Name: "db", Value: col.Name, Options: []string{}},
				{Name: "json", Value: col.Name, Options: []string{}},
			},
			Comments: "",
		})
	}

	if hasNull {
		modelNode.Imports.Append(NullPackage)
	}

	return
}

func buildFileFromModelNode(schemaName string, table *schema.Table, modelNode *lib.GoStruct) (file []byte, e error) {

	insertColumns := fetchInsertColumns(table.ToSortedColumns())
	updateColumns := fetchUpdateColumns(table.ToSortedColumns())
	primaryKey := fetchTablePrimaryKeyName(table)

	var b strings.Builder
	b.WriteString("// Generated Code; DO NOT EDIT.\n\npackage " + modelNode.Package + "\n\n")
	if modelNode.Imports.Len() > 0 {
		b.WriteString(modelNode.Imports.ToString() + "\n")
	}

	hasAccountID := false
	hasUserID := false

	b.WriteString(`
const (

	// ` + modelNode.Name + `_SchemaName is the name of the schema group this model is in
	` + modelNode.Name + `_SchemaName = "` + schemaName + `"
	
	// ` + modelNode.Name + `_TableName is the name of the table 
	` + modelNode.Name + `_TableName query.TableName = "` + modelNode.Name + `"

	// Columns 
`)
	for _, f := range *modelNode.Fields {
		b.WriteString("\t" + modelNode.Name + "_Column_" + f.Name + " query.Column = \"" + f.Name + "\"\n")
		if f.Name == "AccountID" {
			hasAccountID = true
		}
		if f.Name == "UserID" {
			hasUserID = true
		}
	}

	b.WriteString(`
)

var (
	// ` + modelNode.Name + `_Columns is a list of all the columns
	` + modelNode.Name + `_Columns = []query.Column{
`)

	for k, f := range *modelNode.Fields {
		b.WriteString(modelNode.Name + "_Column_" + f.Name)
		if k < len(*modelNode.Fields)-1 {
			b.WriteByte(',')
		}
	}
	b.WriteString(`	}

	// ` + modelNode.Name + `_Column_Types maps columns to their string types
	` + modelNode.Name + `_Column_Types = map[query.Column]string{
`)

	// Column Types
	for k, f := range *modelNode.Fields {
		b.WriteString(modelNode.Name + "_Column_" + f.Name + ": \"" + schema.GoTypeFormatString(f.DataType) + "\"")
		if k < len(*modelNode.Fields)-1 {
			b.WriteByte(',')
		}
	}
	b.WriteString("}\n")

	// Update columns
	b.WriteString("\t// " + modelNode.Name + "_UpdateColumns is a list of all update columns for this model\n")
	b.WriteString("\t" + modelNode.Name + "_UpdateColumns = []query.Column{")
	for k := range updateColumns {
		col := updateColumns[k]
		b.WriteString(modelNode.Name + "_Column_" + col.Name)
		if k < len(updateColumns)-1 {
			b.WriteByte(',')
		}
	}
	b.WriteString("}\n")

	// Insert columns
	b.WriteString("\t// " + modelNode.Name + "_InsertColumns is a list of all insert columns for this model\n")
	b.WriteString("\t" + modelNode.Name + "_InsertColumns = []query.Column{")
	for k := range insertColumns {
		col := insertColumns[k]
		b.WriteString(modelNode.Name + "_Column_" + col.Name)
		if k < len(insertColumns)-1 {
			b.WriteByte(',')
		}
	}
	b.WriteString("}\n")

	// Primary Key
	b.WriteString("\t// " + modelNode.Name + "_PrimaryKey is the name of the table's primary key\n")
	b.WriteString("\t" + modelNode.Name + "_PrimaryKey query.Column = \"" + primaryKey + "\"\n)")

	// Model
	if len(modelNode.Comments) > 0 {
		b.WriteString("\n// " + modelNode.Comments)
	}
	b.WriteString("type " + modelNode.Name + " struct {\n")
	for _, f := range *modelNode.Fields {
		b.WriteString("\t" + f.ToString())
	}
	b.WriteString("}\n")

	if hasAccountID {
		b.WriteString(
			`
// Account satisifies the IAccountable interface
func (c *` + modelNode.Name + `) Account() int64 { 
	return c.AccountID
}
`)
	}
	if hasUserID {
		b.WriteString(
			`
// User satisifies the IUserable interface
func (c *` + modelNode.Name + `) User() int64 { 
	return c.UserID
}
`)
	}
	b.WriteString(`

// ` + modelNode.Name + `_TableName is the name of the table
func (c *` + modelNode.Name + `) Table_Name() query.TableName {
	return ` + modelNode.Name + `_TableName
}

func (c *` + modelNode.Name + `) Table_Columns() []query.Column {
	return ` + modelNode.Name + `_Columns
}

// Table_ColumnTypes returns a map of tableColumn names with their fmt string types
func (c *` + modelNode.Name + `) Table_Column_Types() map[query.Column]string {
	return ` + modelNode.Name + `_Column_Types
}

// Table_PrimaryKey returns the name of this table's primary key 
func (c *` + modelNode.Name + `) Table_PrimaryKey() query.Column {
	return ` + modelNode.Name + `_PrimaryKey
}

// Table_PrimaryKey_Value returns the value of this table's primary key
func (c *` + modelNode.Name + `) Table_PrimaryKey_Value() int64 {
	return c.` + primaryKey + `
}

// Table_InsertColumns is a list of all insert columns for this model
func (c *` + modelNode.Name + `) Table_InsertColumns() []query.Column {
	return ` + modelNode.Name + `_InsertColumns
}

// Table_UpdateColumns is a list of all update columns for this model
func (c *` + modelNode.Name + `) Table_UpdateColumns() []query.Column {
	return ` + modelNode.Name + `_UpdateColumns
}

// ` + modelNode.Name + `_SchemaName is the name of this table's schema
func (c *` + modelNode.Name + `) Table_SchemaName() string {
	return ` + modelNode.Name + `_SchemaName
}

// FromID returns a FromID query statement
func (c *` + modelNode.Name + `) FromID(db db.IDB, id int64) (query.IModel, error) {
	
	sel := query.Select(c)

	`)
	b.WriteString(`	sel.Fields(
`)
	for _, f := range *modelNode.Fields {
		b.WriteString(`	query.NewField(query.FieldTypeBasic, ` + modelNode.Name + `_Column_` + f.Name + `),
`)
	}
	b.WriteString(`
	)
	q, e := sel.String()
	if e != nil {
		return nil, fmt.Errorf("` + modelNode.Name + `.FromID.Query.String(): %w", e)
	}

	row := db.QueryRow(q)

	switch e = row.Scan(
`)

	for _, f := range *modelNode.Fields {
		b.WriteString(`		&c.` + f.Name + `,
`)
	}

	b.WriteString(`	); e { 
	case sql.ErrNoRows: 
		return nil, nil 
	case nil: 
		fmt.Printf("` + modelNode.Name + `DALGetter.Get(%s).Run()\n", q)
		return c, nil 
	default: 
		return nil, fmt.Errorf("` + modelNode.Name + `DALGetter(%s).Run(): %w", q, e)
	}
}

// String returns a json marshalled string of the object 
func (c *` + modelNode.Name + `) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

// Update updates a ` + modelNode.Name + ` record
func (c *` + modelNode.Name + `) Update(db db.IDB) error {
	var e error 
	var ql string 
	ql, _ = query.Update(c).
`)
	for k := range updateColumns {
		col := updateColumns[k]

		var value string
		if col.GoType == "null.String" {
			value = "c." + col.Name + ".String"
		} else {
			value = "c." + col.Name
		}

		b.WriteString("\t\tSet(" + modelNode.Name + "_Column_" + col.Name + ", " + value + ").\n")
	}
	b.WriteString(`
		Where(query.EQ(` + modelNode.Name + "_Column_" + primaryKey + `, c.` + primaryKey + `)).
	String()

	_, e = db.Exec(ql) 
	if e != nil {
		return fmt.Errorf("` + modelNode.Name + `.Update(): %w", e)
	}

	return e 
}

// Create inserts a ` + modelNode.Name + ` record
func (c *` + modelNode.Name + `) Create(db db.IDB) error { 

	var e error 
	q := query.Insert(c)

	if c.` + primaryKey + ` > 0 { 
		q.Set(` + modelNode.Name + "_Column_" + primaryKey + `, c.` + primaryKey + `)
	}
`)

	for k := range insertColumns {
		col := insertColumns[k]

		var value string
		if col.GoType == "null.String" {
			value = "c." + col.Name + ".String"
		} else {
			value = "c." + col.Name
		}

		b.WriteString("\tq.Set(" + modelNode.Name + "_Column_" + col.Name + ", " + value + ")\n")
	}
	b.WriteString(`

	ql, _ := q.String()
	var result sql.Result 
	result, e = db.Exec(ql) 
	if e != nil {
		return fmt.Errorf("` + modelNode.Name + `.Create(): %w", e)
	}

	// Assumes auto-increment 
	if c.` + primaryKey + ` == 0 {
		c.` + primaryKey + `, e = result.LastInsertId()
	}

	return e 
} 
	`)
	b.WriteString(`

// Destroy deletes a ` + modelNode.Name + ` record
func (c *` + modelNode.Name + `) Delete(db db.IDB) error {
	var e error 
	ql, _ := query.Delete(c).
		Where(
			query.EQ(` + modelNode.Name + "_Column_" + primaryKey + `, c.` + primaryKey + `),
		).String()

	_, e = db.Exec(ql)
	if e != nil {
		return fmt.Errorf("` + modelNode.Name + `.Delete(): %w", e)
	}

	return e
}

func (r *` + modelNode.Name + `) Raw(db db.IDB, queryRaw string) ([]*` + modelNode.Name + `, error) {

	var e error
	model := []*` + modelNode.Name + `{}
	sel := query.Select(r) 

`)

	b.WriteString(`	sel.Fields(
`)
	for _, f := range *modelNode.Fields {
		b.WriteString(`	query.NewField(query.FieldTypeBasic, ` + modelNode.Name + `_Column_` + f.Name + `),
`)
	}
	b.WriteString(`
	)

	q, e := sel.String()
	if e != nil {
		return nil, fmt.Errorf("` + modelNode.Name + `DAL.Raw.String(): %w", e)
	}

	var rows *sql.Rows 
	rows, e = db.Query(q) 

	if e != nil {
		if e == sql.ErrNoRows { 
			return nil, nil 
		}
		return nil, fmt.Errorf("` + modelNode.Name + `DAL.Raw.Run(%s): %w", q, e)
	}

	defer rows.Close() 
	for rows.Next() { 
		m := &` + modelNode.Name + `{}
		if e = rows.Scan(
`)

	for _, f := range *modelNode.Fields {
		b.WriteString(`			&m.` + f.Name + `,
`)
	}

	b.WriteString(`		); e != nil { 
			return nil, fmt.Errorf("` + modelNode.Name + `DALRaw(%s).Run(): %w", q, e)
		}

		model = append(model, m)
	}

	fmt.Printf("` + modelNode.Name + `DAL.Raw(%s).Run()\n", q)

	return model, nil
}

type ` + modelNode.Name + `DALSelector struct {
	db    	 db.IDB
	q     	 *query.Q
	isSingle bool 
}

func (r *` + modelNode.Name + `) Select(db db.IDB) *` + modelNode.Name + `DALSelector {
	return &` + modelNode.Name + `DALSelector{
		db:    db,
		q:     query.Select(r),
	}
}

func (r *` + modelNode.Name + `DALSelector) Alias(alias string) *` + modelNode.Name + `DALSelector { 
	r.q.Alias(alias) 
	return r
}

func (r *` + modelNode.Name + `DALSelector) Where(whereParts ...*query.WherePart) *` + modelNode.Name + `DALSelector {
	r.q.Where(whereParts...)
	return r
}

func (r *` + modelNode.Name + `DALSelector) Limit(limit, offset int64) *` + modelNode.Name + `DALSelector {
	r.q = r.q.Limit(limit, offset)
	return r
}

func (r *` + modelNode.Name + `DALSelector) OrderBy(col query.Column, dir query.OrderDir) *` + modelNode.Name + `DALSelector {
	r.q = r.q.OrderBy(col, dir)
	return r
}

func (r *` + modelNode.Name + `DALSelector) Run() ([]*` + modelNode.Name + `, error) {
	model := []*` + modelNode.Name + `{}
`)

	b.WriteString(`	r.q.Fields(
`)
	for _, f := range *modelNode.Fields {
		b.WriteString(`	query.NewField(query.FieldTypeBasic, ` + modelNode.Name + `_Column_` + f.Name + `),
`)
	}
	b.WriteString(`
	)

	q, e := r.q.String()
	if e != nil {
		return nil, fmt.Errorf("` + modelNode.Name + `DAL.Query.String(): %w", e)
	}

	var rows *sql.Rows 
	rows, e = r.db.Query(q) 

	if e != nil {
		if e == sql.ErrNoRows { 
			return nil, nil 
		}
		return nil, fmt.Errorf("` + modelNode.Name + `DALSelector.Run(%s): %w", q, e)
	}

	defer rows.Close() 
	for rows.Next() { 
		m := &` + modelNode.Name + `{}
		if e = rows.Scan(
`)

	for _, f := range *modelNode.Fields {
		b.WriteString(`			&m.` + f.Name + `,
`)
	}

	b.WriteString(`		); e != nil { 
			return nil, fmt.Errorf("` + modelNode.Name + `DALSelector(%s).Run(): %w", q, e)
		}

		model = append(model, m)
	}

	fmt.Printf("` + modelNode.Name + `DALSelector(%s).Run()\n", q)

	return model, nil
}

// Counter 
type ` + modelNode.Name + `DALCounter struct {
	db    db.IDB
	q     *query.Q
}

func (r *` + modelNode.Name + `) Count(db db.IDB) *` + modelNode.Name + `DALCounter {
	return &` + modelNode.Name + `DALCounter{
		db:    db,
		q:     query.Select(r).Count(r.Table_PrimaryKey(), "c"),
	}
}

func (r *` + modelNode.Name + `DALCounter) Alias(alias string) *` + modelNode.Name + `DALCounter { 
	r.q.Alias(alias) 
	return r
}

func (ds *` + modelNode.Name + `DALCounter) Where(whereParts ...*query.WherePart) *` + modelNode.Name + `DALCounter {
	ds.q.Where(whereParts...)
	return ds
}

func (ds *` + modelNode.Name + `DALCounter) Run() (int64, error) {

	count := int64(0)
	q, e := ds.q.String()
	if e != nil {
		return 0, fmt.Errorf("` + modelNode.Name + `DALCounter.Query.String(): %w", e)
	}

	row := ds.db.QueryRow(q)

	switch e = row.Scan(&count); e { 
	case sql.ErrNoRows: 
		return 0, nil 
	case nil: 
		fmt.Printf("` + modelNode.Name + `DALCounter.QueryRow(%s).Run()\n", q)
		return count, nil 
	default: 
		return 0, fmt.Errorf("` + modelNode.Name + `DALCounter.QueryRow(%s).Run(): %w", q, e)
	}
}

// Summer
type ` + modelNode.Name + `DALSummer struct {
	db    db.IDB
	q     *query.Q
}

func (r *` + modelNode.Name + `) Sum(db db.IDB, col query.Column) *` + modelNode.Name + `DALSummer {
	return &` + modelNode.Name + `DALSummer{
		db:    db,
		q:     query.Select(r).Sum(col, "c"),
	}
}

func (ds *` + modelNode.Name + `DALSummer) Where(whereParts ...*query.WherePart) *` + modelNode.Name + `DALSummer {
	ds.q.Where(whereParts...)
	return ds
}

func (ds *` + modelNode.Name + `DALSummer) Run() (float64, error) {

	sum := float64(0)
	q, e := ds.q.String()
	if e != nil {
		return 0, fmt.Errorf("` + modelNode.Name + `DALSummer.Query.String(): %w", e)
	}

	row := ds.db.QueryRow(q)

	switch e = row.Scan(&sum); e { 
	case sql.ErrNoRows: 
		return 0, nil 
	case nil: 
		fmt.Printf("` + modelNode.Name + `DALSummer.QueryRow(%s).Run()\n", q)
		return sum, nil 
	default: 
		return 0, fmt.Errorf("` + modelNode.Name + `DALSummer.QueryRow(%s).Run(): %w", q, e)
	}
}

// Minner
type ` + modelNode.Name + `DALMinner struct {
	db    db.IDB
	q     *query.Q
}

func (r *` + modelNode.Name + `) Min(db db.IDB, col query.Column) *` + modelNode.Name + `DALMinner {
	return &` + modelNode.Name + `DALMinner{
		db:    db,
		q:     query.Select(r).Min(col, "c"),
	}
}

func (ds *` + modelNode.Name + `DALMinner) Where(whereParts ...*query.WherePart) *` + modelNode.Name + `DALMinner {
	ds.q.Where(whereParts...)
	return ds
}

func (ds *` + modelNode.Name + `DALMinner) Run() (float64, error) {

	min := float64(0)
	q, e := ds.q.String()
	if e != nil {
		return 0, fmt.Errorf("` + modelNode.Name + `DALMinner.Query.String(): %w", e)
	}

	row := ds.db.QueryRow(q)

	switch e = row.Scan(&min); e { 
	case sql.ErrNoRows: 
		return 0, nil 
	case nil: 
		fmt.Printf("` + modelNode.Name + `DALMinner.QueryRow(%s).Run()\n", q)
		return min, nil 
	default: 
		return 0, fmt.Errorf("` + modelNode.Name + `DALMinner.QueryRow(%s).Run(): %w", q, e)
	}
}

// Maxer
type ` + modelNode.Name + `DALMaxer struct {
	db    db.IDB
	q     *query.Q
}

func (r *` + modelNode.Name + `) Max(db db.IDB, col query.Column) *` + modelNode.Name + `DALMaxer {
	return &` + modelNode.Name + `DALMaxer{
		db:    db,
		q:     query.Select(r).Max(col, "c"),
	}
}

func (ds *` + modelNode.Name + `DALMaxer) Where(whereParts ...*query.WherePart) *` + modelNode.Name + `DALMaxer {
	ds.q.Where(whereParts...)
	return ds
}

func (ds *` + modelNode.Name + `DALMaxer) Run() (float64, error) {

	max := float64(0)
	q, e := ds.q.String()
	if e != nil {
		return 0, fmt.Errorf("` + modelNode.Name + `DALMaxer.Query.String(): %w", e)
	}

	row := ds.db.QueryRow(q)

	switch e = row.Scan(&max); e { 
	case sql.ErrNoRows: 
		return 0, nil 
	case nil: 
		fmt.Printf("` + modelNode.Name + `DALMaxer.QueryRow(%s).Run()\n", q)
		return max, nil 
	default: 
		return 0, fmt.Errorf("` + modelNode.Name + `DALMaxer.QueryRow(%s).Run(): %w", q, e)
	}
}



type ` + modelNode.Name + `DALGetter struct {
	db    	 db.IDB
	q     	 *query.Q
}

func (r *` + modelNode.Name + `) Get(db db.IDB) *` + modelNode.Name + `DALGetter {
	return &` + modelNode.Name + `DALGetter{
		db:    db,
		q:     query.Select(r),
	}
}

func (r *` + modelNode.Name + `DALGetter) Alias(alias string) *` + modelNode.Name + `DALGetter { 
	r.q.Alias(alias) 
	return r
}

func (ds *` + modelNode.Name + `DALGetter) Where(whereParts ...*query.WherePart) *` + modelNode.Name + `DALGetter {
	ds.q.Where(whereParts...)
	return ds
}

func (ds *` + modelNode.Name + `DALGetter) OrderBy(col query.Column, dir query.OrderDir) *` + modelNode.Name + `DALGetter {
	ds.q = ds.q.OrderBy(col, dir)
	return ds
}

func (ds *` + modelNode.Name + `DALGetter) Run() (*` + modelNode.Name + `, error) {

	model := &` + modelNode.Name + `{}

	`)
	b.WriteString(`	ds.q.Fields(
`)
	for _, f := range *modelNode.Fields {
		b.WriteString(`	query.NewField(query.FieldTypeBasic, ` + modelNode.Name + `_Column_` + f.Name + `),
`)
	}
	b.WriteString(`
	)
	q, e := ds.q.String()
	if e != nil {
		return nil, fmt.Errorf("` + modelNode.Name + `DALGetter.Query.String(): %w", e)
	}

	row := ds.db.QueryRow(q)

	switch e = row.Scan(
`)

	for _, f := range *modelNode.Fields {
		b.WriteString(`		&model.` + f.Name + `,
`)
	}

	b.WriteString(`	); e { 
	case sql.ErrNoRows: 
		return nil, nil 
	case nil: 
		fmt.Printf("` + modelNode.Name + `DALGetter.Get(%s).Run()\n", q)
		return model, nil 
	default: 
		return nil, fmt.Errorf("` + modelNode.Name + `DALGetter(%s).Run(): %w", q, e)
	}
}
`)

	// Write the file

	file = []byte(b.String())

	file, e = format.Source(file)
	if e != nil {
		log.Fatalf("FORMAT ERROR: File: %s; Error: %s\n%s", modelNode.Name, e.Error(), b.String())
	}

	file, e = lib.FormatCode(string(file))
	return
}

// func (g *Gen) buildGoModelOld(p string, table *schema.Table) (e error) {

// 	oneToMany := g.Config.OneToMany[table.Name]
// 	oneToOne := g.Config.OneToOne[table.Name]

// 	lib.Debugf("Generating model for table %s", g.Options, table.Name)

// 	type Column struct {
// 		Name         string
// 		Type         string
// 		IsPrimaryKey bool
// 		IsNull       bool
// 		Ordinal      int
// 	}

// 	data := struct {
// 		OneToMany          string
// 		OneToOne           string
// 		Name               string
// 		IncludeNullPackage bool
// 		Columns            []Column
// 	}{
// 		OneToMany:          oneToMany,
// 		OneToOne:           oneToOne,
// 		Name:               table.Name,
// 		IncludeNullPackage: false,
// 		Columns:            []Column{},
// 	}

// 	tpl := `// Generated Code; DO NOT EDIT.

// package models
// {{if .IncludeNullPackage}}
// import (
// 	"gopkg.in/guregu/null.v3"
// )
// {{end}}

// // {{.Name}} represents a {{.Name}} domain object
// type {{.Name}} struct {
// 	{{range .Columns}}
// {{.Name}} {{.Type}} ` + "`db:\"{{.Name}}\" json:\"{{.Name}}\"`" + `{{end}}
// }`
// 	var sortedColumns = make(lib.SortedColumns, 0, len(table.Columns))

// 	for _, column := range table.Columns {
// 		sortedColumns = append(sortedColumns, column)
// 	}

// 	sort.Sort(sortedColumns)

// 	includeNullPackage := false

// 	for _, column := range sortedColumns {

// 		fieldType := schema.DataTypeToGoTypeString(column)

// 		if column.IsNullable {
// 			includeNullPackage = true
// 		}

// 		data.Columns = append(data.Columns, Column{
// 			Name:         column.Name,
// 			Type:         fieldType,
// 			IsPrimaryKey: column.ColumnKey == "PRI",
// 			IsNull:       column.IsNullable,
// 			Ordinal:      len(data.Columns),
// 		})
// 	}
// 	data.IncludeNullPackage = includeNullPackage

// 	t := template.Must(template.New("model-" + table.Name).Parse(tpl))
// 	f, err := os.Create(p)
// 	if err != nil {
// 		fmt.Println("ERROR: ", err.Error())
// 		return
// 	}

// 	err = t.Execute(f, data)
// 	if err != nil {
// 		fmt.Println("Execute ERROR: ", err.Error())
// 		return
// 	}

// 	f.Close()
// 	lib.FmtGoCode(p)

// 	return
// }

// CleanGoModels removes model files that are not found in the database.Tables map
func CleanGoModels(dir string, tables map[string]struct{}) (e error) {

	lib.EnsureDir(dir)

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
	reader := bufio.NewReader(os.Stdin)

	for _, name := range dirFileNames {
		fileNameNoExt := name[0 : len(name)-3]
		if _, ok := tables[fileNameNoExt]; !ok {
			fullFilePath := path.Join(dir, name)
			// log.Printf("Removing %s\n", fullFilePath)
			result := lib.ReadCliInput(reader, fmt.Sprintf("Delete unused model `%s` (Y/n)?", fileNameNoExt))
			if result == "Y" {
				fmt.Printf("Deleting model `%s` at path `%s`...\n", fileNameNoExt, fullFilePath)
				os.Remove(fullFilePath)
			}
		}
	}
	return
}
