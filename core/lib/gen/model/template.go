package model

import "text/template"

var goModelTemplate = template.Must(template.New("go-model-file").Parse(`// Generated Code; DO NOT EDIT.
package models 

import (
	"github.com/macinnir/dvc/core/lib/utils/query"
	"github.com/macinnir/dvc/core/lib/utils/db"
	"encoding/json"
	"fmt"
	"database/sql" {{ if .HasNull }}
	"gopkg.in/guregu/null.v3"{{ end }}
)

const (

	// {{ $.Name }}_SchemaName is the name of the schema group this model is in
	{{ $.Name }}_SchemaName = "{{ $.Schema }}"
	
	// {{ $.Name }}_TableName is the name of the table 
	{{ $.Name }}_TableName query.TableName = "{{ $.Name }}"

	// Columns 
	{{ range .Fields }}
	{{ $.Name }}_Column_{{ .Name }} query.Column = "{{ .Name }}"{{ end }}
)

var (
	// {{ $.Name }}_Columns is a list of all the columns
	{{ $.Name }}_Columns = []query.Column{ {{ range .Fields }}
		{{ $.Name }}_Column_{{ .Name }},{{ end }}
	}

	// {{ $.Name }}_Column_Types maps columns to their string types
	{{ $.Name }}_Column_Types = map[query.Column]string{ {{ range .Fields }}
		{{ $.Name }}_Column_{{ .Name }}: "{{ .FormatType }}",{{ end }} 
	}

	// Update columns
	// {{ $.Name }}_UpdateColumns is a list of all update columns for this model 
	{{ $.Name }}_UpdateColumns = []query.Column{ {{ range .UpdateColumns }}
		{{ $.Name }}_Column_{{ .Name }},{{ end }}
	}

	// Insert columns
	// {{ $.Name }}_InsertColumns is a list of all insert columns for this model 
	{{ $.Name }}_InsertColumns = []query.Column{ {{ range .InsertColumns}}
		{{ $.Name }}_Column_{{ .Name }},{{ end }}
	}

	// Primary Key
	// {{ $.Name }}_PrimaryKey is the name of the table's primary key 
	{{ $.Name }}_PrimaryKey query.Column = "{{.PrimaryKey}}"
)

// {{ $.Name }} is a data model
type {{ $.Name }} struct { 
	{{ range .Fields }}
	{{ .Name }} {{ .GoType }} ` + "`" + `db:"{{ .Name }}" json:"{{ .Name }}"` + "`" + `{{ end }}
}

{{ if .HasAccountID }}// Account satifies the IAccountable interface 
func (c *{{ $.Name }}) Account() int64 { 
	return c.AccountID
}{{ end }} // 63

{{ if .HasUserID }}// User satisifies the IUserable interface 
func (c *{{ $.Name }}) User() int64 { 
	return c.UserID
}{{ end }} // 68

// {{ $.Name }}_TableName is the name of the table
func (c *{{ $.Name }}) Table_Name() query.TableName {
	return {{ $.Name }}_TableName
}

func (c *{{ $.Name }}) Table_Columns() []query.Column {
	return {{ $.Name }}_Columns
}

// Table_ColumnTypes returns a map of tableColumn names with their fmt string types
func (c *{{ $.Name }}) Table_Column_Types() map[query.Column]string {
	return {{ $.Name }}_Column_Types
}

// Table_PrimaryKey returns the name of this table's primary key 
func (c *{{ $.Name }}) Table_PrimaryKey() query.Column {
	return {{ $.Name }}_PrimaryKey
}

// Table_PrimaryKey_Value returns the value of this table's primary key
func (c *{{ $.Name }}) Table_PrimaryKey_Value() int64 {
	return c.{{ $.PrimaryKey }}
}

// Table_InsertColumns is a list of all insert columns for this model
func (c *{{ $.Name }}) Table_InsertColumns() []query.Column {
	return {{ $.Name }}_InsertColumns
}

// Table_UpdateColumns is a list of all update columns for this model
func (c *{{ $.Name }}) Table_UpdateColumns() []query.Column { // 100
	return {{ $.Name }}_UpdateColumns
}

// {{ $.Name }}_SchemaName is the name of this table's schema
func (c *{{ $.Name }}) Table_SchemaName() string {
	return {{ $.Name }}_SchemaName
}


// String returns a json marshalled string of the object 
func (c *{{ $.Name }}) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

// Update updates a {{ $.Name }} record
func (c *{{ $.Name }}) Update(db db.IDB) error {
	var e error 
	var ql string 
	ql, _ = query.Update(c).{{ range .UpdateColumns }}
		Set({{ $.Name }}_Column_{{ .Name }}, c.{{ .Name }}{{ if eq .GoType "null.String" }}.String{{ end }}).{{ end }}
		Where(query.EQ({{ $.Name }}_Column_{{ $.PrimaryKey }}, c.{{ $.PrimaryKey }})).
		String()
	_, e = db.Exec(ql) 
	if e != nil {
		return fmt.Errorf("{{ $.Name }}.Update(): %w", e)
	}

	return e 
}

// Create inserts a {{ $.Name }} record
func (c *{{ $.Name }}) Create(db db.IDB) error {
	
	var e error 
	
	q := query.Insert(c)
	
	if c.{{ $.PrimaryKey }} > 0 { 
		q.Set({{ $.Name }}_Column_{{ $.PrimaryKey }}, c.{{ $.PrimaryKey }})
	}
	
	{{ range .InsertColumns }}
		q.Set({{ $.Name }}_Column_{{ .Name }}, c.{{ .Name }}{{ if eq .GoType "null.String" }}.String{{ end }}){{ end }}

	ql, _ := q.String()
	var result sql.Result
	result, e = db.Exec(ql)

	if e != nil {
		return fmt.Errorf("{{ $.Name }}.Create(): %w", e) // 177
	}

	// Assumes auto-increment
	if c.{{ $.PrimaryKey }} == 0 {
		c.{{ $.PrimaryKey }}, e = result.LastInsertId()
	}

	return e 
}


// Destroy deletes a {{ $.Name }} record
func (c *{{ $.Name }}) Delete(db db.IDB) error {
	var e error 
	ql, _ := query.Delete(c).
		Where(
			query.EQ({{ $.Name }}_Column_{{ $.PrimaryKey }}, c.{{ $.PrimaryKey }}),
		).String()

	_, e = db.Exec(ql)
	if e != nil {
		return fmt.Errorf("{{ $.Name }}.Delete(): %w", e)
	}

	return e
}

func (r *{{ $.Name }}) Raw(db db.IDB, queryRaw string) ([]*{{ $.Name }}, error) {

	var e error
	model := []*{{ $.Name }}{}
	sel := query.Select(r) 
	sel.Fields({{ range .SelectFields }}
		{{ if eq .DBType "vector" }}{{ else }}query.NewField(query.FieldTypeBasic, {{ $.Name }}_Column_{{ .Name }}),{{ end }}{{ end }}
	)

	q, e := sel.String()
	if e != nil {
		return nil, fmt.Errorf("{{ $.Name }}DAL.Raw.String(): %w", e)
	}

	var rows *sql.Rows 
	rows, e = db.Query(q) 

	if e != nil {
		if e == sql.ErrNoRows { 
			return nil, nil 
		}
		return nil, fmt.Errorf("{{ $.Name }}DAL.Raw.Run(%s): %w", q, e)
	}

	defer rows.Close() 
	for rows.Next() { 
		m := &{{ $.Name }}{}
		if e = rows.Scan({{ range .SelectFields }}
			&m.{{ .Name }},{{ end }} 
		); e != nil { 
			return nil, fmt.Errorf("{{ $.Name }}DALRaw(%s).Run(): %w", q, e)
		}
		model = append(model, m)
	}

	// fmt.Printf("{{ $.Name }}DAL.Raw(%s).Run()\n", q)

	return model, nil
}

type I{{ $.Name }}DALSelector interface { 
	Select(db db.IDB) I{{ $.Name }}DALSelector
}

type {{ $.Name }}DALSelector struct {
	db    	 db.IDB
	q     	 *query.Q
	isSingle bool 
}

func (r *{{ $.Name }}) Select(db db.IDB) *{{ $.Name }}DALSelector {
	return &{{ $.Name }}DALSelector{
		db:    db,
		q:     query.Select(r),
	}
}

func (r *{{ $.Name }}DALSelector) Alias(alias string) *{{ $.Name }}DALSelector { 
	r.q.Alias(alias) 
	return r
}

func (r *{{ $.Name }}DALSelector) Where(whereParts ...*query.WherePart) *{{ $.Name }}DALSelector {
	r.q.Where(whereParts...)
	return r
}

func (r *{{ $.Name }}DALSelector) Limit(limit, offset int64) *{{ $.Name }}DALSelector {
	r.q = r.q.Limit(limit, offset)
	return r
}

func (r *{{ $.Name }}DALSelector) OrderBy(col query.Column, dir query.OrderDir) *{{ $.Name }}DALSelector {
	r.q = r.q.OrderBy(col, dir)
	return r
}

func (r *{{ $.Name }}DALSelector) String() (string, error) { 
	
	r.q.Fields({{ range .SelectFields }}
		{{ if eq .DBType "vector" }}{{ else }}query.NewField(query.FieldTypeBasic, {{ $.Name }}_Column_{{ .Name }}),{{ end }}{{ end }}
	)

	q, e := r.q.String()
	if e != nil {
		return "", fmt.Errorf("{{ $.Name }}DAL.Query.String(): %w", e)
	}

	return q, nil 
}

func (r *{{ $.Name }}DALSelector) Run() ([]*{{ $.Name }}, error) {

	var e error 
	var q = "" 
	var model = []*{{ $.Name }}{}
	q, e = r.String()
	if e != nil { 
		return nil, fmt.Errorf("{{ $.Name }}DALSelector.Query.String(): %w", e)
	}

	var rows *sql.Rows 
	rows, e = r.db.Query(q) 

	if e != nil {
		if e == sql.ErrNoRows { 
			return nil, nil 
		}
		return nil, fmt.Errorf("{{ $.Name }}DALSelector.Run(%s): %w", q, e)
	}

	defer rows.Close() 
	for rows.Next() { 
		m := &{{ $.Name }}{}
		if e = rows.Scan({{ range .SelectFields }}
			&m.{{ .Name }},{{ end }}
		); e != nil { 
			return nil, fmt.Errorf("{{ $.Name }}DALSelector(%s).Run(): %w", q, e)
		}

		model = append(model, m)
	}

	// fmt.Printf("{{ $.Name }}DALSelector(%s).Run()\n", q)

	return model, nil
}

// Counter 
type {{ $.Name }}DALCounter struct {
	db    db.IDB
	q     *query.Q
}

func (r *{{ $.Name }}) Count(db db.IDB) *{{ $.Name }}DALCounter {
	return &{{ $.Name }}DALCounter{
		db:    db,
		q:     query.Select(r).Count(r.Table_PrimaryKey(), "c"),
	}
}

func (r *{{ $.Name }}DALCounter) Alias(alias string) *{{ $.Name }}DALCounter { 
	r.q.Alias(alias) 
	return r
}

func (ds *{{ $.Name }}DALCounter) Where(whereParts ...*query.WherePart) *{{ $.Name }}DALCounter {
	ds.q.Where(whereParts...)
	return ds
}

func (ds *{{ $.Name }}DALCounter) Run() (int64, error) {

	count := int64(0)
	q, e := ds.q.String()
	if e != nil {
		return 0, fmt.Errorf("{{ $.Name }}DALCounter.Query.String(): %w", e)
	}

	row := ds.db.QueryRow(q)

	switch e = row.Scan(&count); e { 
	case sql.ErrNoRows: 
		return 0, nil 
	case nil: 
		// fmt.Printf("{{ $.Name }}DALCounter.QueryRow(%s).Run()\n", q)
		return count, nil 
	default: 
		return 0, fmt.Errorf("{{ $.Name }}DALCounter.QueryRow(%s).Run(): %w", q, e)
	}
}

// Summer
type {{ $.Name }}DALSummer struct {
	db    db.IDB
	q     *query.Q
}

func (r *{{ $.Name }}) Sum(db db.IDB, col query.Column) *{{ $.Name }}DALSummer {
	return &{{ $.Name }}DALSummer{
		db:    db,
		q:     query.Select(r).Sum(col, "c"),
	}
}

func (ds *{{ $.Name }}DALSummer) Where(whereParts ...*query.WherePart) *{{ $.Name }}DALSummer {
	ds.q.Where(whereParts...)
	return ds
}

func (ds *{{ $.Name }}DALSummer) Run() (float64, error) {

	sum := float64(0)
	q, e := ds.q.String()
	if e != nil {
		return 0, fmt.Errorf("{{ $.Name }}DALSummer.Query.String(): %w", e)
	}

	row := ds.db.QueryRow(q)

	switch e = row.Scan(&sum); e { 
	case sql.ErrNoRows: 
		return 0, nil 
	case nil: 
		// fmt.Printf("{{ $.Name }}DALSummer.QueryRow(%s).Run()\n", q)
		return sum, nil 
	default: 
		return 0, fmt.Errorf("{{ $.Name }}DALSummer.QueryRow(%s).Run(): %w", q, e)
	}
}

// Minner
type {{ $.Name }}DALMinner struct {
	db    db.IDB
	q     *query.Q
}

func (r *{{ $.Name }}) Min(db db.IDB, col query.Column) *{{ $.Name }}DALMinner {
	return &{{ $.Name }}DALMinner{
		db:    db,
		q:     query.Select(r).Min(col, "c"),
	}
}

func (ds *{{ $.Name }}DALMinner) Where(whereParts ...*query.WherePart) *{{ $.Name }}DALMinner {
	ds.q.Where(whereParts...)
	return ds
}

func (ds *{{ $.Name }}DALMinner) Run() (float64, error) {

	min := float64(0)
	q, e := ds.q.String()
	if e != nil {
		return 0, fmt.Errorf("{{ $.Name }}DALMinner.Query.String(): %w", e)
	}

	row := ds.db.QueryRow(q)

	switch e = row.Scan(&min); e { 
	case sql.ErrNoRows: 
		return 0, nil 
	case nil: 
		// fmt.Printf("{{ $.Name }}DALMinner.QueryRow(%s).Run()\n", q)
		return min, nil 
	default: 
		return 0, fmt.Errorf("{{ $.Name }}DALMinner.QueryRow(%s).Run(): %w", q, e)
	}
}

// Maxer
type {{ $.Name }}DALMaxer struct {
	db    db.IDB
	q     *query.Q
}

func (r *{{ $.Name }}) Max(db db.IDB, col query.Column) *{{ $.Name }}DALMaxer {
	return &{{ $.Name }}DALMaxer{
		db:    db,
		q:     query.Select(r).Max(col, "c"),
	}
}

func (ds *{{ $.Name }}DALMaxer) Where(whereParts ...*query.WherePart) *{{ $.Name }}DALMaxer {
	ds.q.Where(whereParts...)
	return ds
}

func (ds *{{ $.Name }}DALMaxer) Run() (float64, error) {

	max := float64(0)
	q, e := ds.q.String()
	if e != nil {
		return 0, fmt.Errorf("{{ $.Name }}DALMaxer.Query.String(): %w", e)
	}

	row := ds.db.QueryRow(q)

	switch e = row.Scan(&max); e { 
	case sql.ErrNoRows: 
		return 0, nil 
	case nil: 
		// fmt.Printf("{{ $.Name }}DALMaxer.QueryRow(%s).Run()\n", q)
		return max, nil 
	default: 
		return 0, fmt.Errorf("{{ $.Name }}DALMaxer.QueryRow(%s).Run(): %w", q, e)
	}
}



type {{ $.Name }}DALGetter struct {
	db    	 db.IDB
	q     	 *query.Q
}

func (r *{{ $.Name }}) Get(db db.IDB) *{{ $.Name }}DALGetter {
	return &{{ $.Name }}DALGetter{
		db:    db,
		q:     query.Select(r),
	}
}

func (r *{{ $.Name }}DALGetter) Alias(alias string) *{{ $.Name }}DALGetter { 
	r.q.Alias(alias) 
	return r
}

func (ds *{{ $.Name }}DALGetter) Where(whereParts ...*query.WherePart) *{{ $.Name }}DALGetter {
	ds.q.Where(whereParts...)
	return ds
}

func (ds *{{ $.Name }}DALGetter) OrderBy(col query.Column, dir query.OrderDir) *{{ $.Name }}DALGetter {
	ds.q = ds.q.OrderBy(col, dir)
	return ds
}

func (ds *{{ $.Name }}DALGetter) Run() (*{{ $.Name }}, error) {

	model := &{{ $.Name }}{}

	ds.q.Fields({{ range $.Fields }}
		{{ if eq .DBType "vector" }}{{ else }}query.NewField(query.FieldTypeBasic, {{ $.Name }}_Column_{{ .Name }}),{{ end }}{{ end }}
	)
	q, e := ds.q.String()
	if e != nil {
		return nil, fmt.Errorf("{{ $.Name }}DALGetter.Query.String(): %w", e)
	}

	row := ds.db.QueryRow(q)

	switch e = row.Scan({{ range .SelectFields }}
		{{ if eq .DBType "vector" }}{{ else }}&model.{{ .Name }}, {{ end }}{{ end }}
	); e { 
	case sql.ErrNoRows: 
		return nil, nil 
	case nil: 
		// fmt.Printf("{{ $.Name }}DALGetter.Get(%s).Run()\n", q)
		return model, nil 
	default: 
		return nil, fmt.Errorf("{{ $.Name }}DALGetter(%s).Run(): %w", q, e)
	}
}
`))
