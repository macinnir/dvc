package gen

import (
	"bytes"
	"io/ioutil"
	"path"
	"sort"
	"text/template"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

var CacheTemplate = template.Must(template.New("template-cache-file").Funcs(template.FuncMap{
	"dataTypeToGoTypeString": schema.DataTypeToGoTypeString,
	"dataTypeToFormatString": schema.DataTypeToFormatString,
	"toArgName":              toArgName,
	"columnsToMethodName":    columnsToMethodName,
	"columnsToMethodArgs":    columnsToMethodArgs,
	"columnsToMethodParams":  columnsToMethodParams,
	"columnsToKey":           columnsToKey,
	"columnValuesToKey":      columnValuesToKey,
	"columnModelValuesToKey": columnModelValuesToKey,
}).Parse(`// Generated Code; DO NOT EDIT.

package caches

import ( 
	"{{ .BasePackage }}/core/components/redis" 
	"{{ .BasePackage }}/gen/definitions/models" 

	"fmt"
)

const ( 
	// {{.Table.Name}}_Index_ID_Key is the string key for the primary key ({{.PrimaryKey}})
	{{.Table.Name}}_Index_ID_Key = "{{.Table.Name | toArgName}}_idx_id"
{{range $index := .CacheConfig.Indices}}{{if $index.Index.Unique}}
	// {{$.Table.Name}}_Index_{{$index.Columns | columnsToKey}}_Key is the string key of a unique index (hash map) for field {{$index.Index.Field}}
	{{$.Table.Name}}_Index_{{$index.Columns | columnsToKey}}_Key = "{{$.Table.Name | toArgName}}_idx_{{$index.Columns | columnsToKey}}"
{{end}}{{end}})

// {{.Table.Name}}_Key returns the string key for the unique primary key 
func {{.Table.Name}}_Key(id int64) string { 
	return fmt.Sprintf("{{ .Table.Name | toArgName }}_%d", id)
}
{{range $index := .CacheConfig.Indices}}{{if not $index.Index.Unique}}

// {{$.Table.Name}}_Index_{{range $column := $index.Columns}}{{$column.Name}}_{{end}}Key is the string key of a non-unique index for fields {{range $column := $index.Columns}}{{$column.Name}}, {{end}}
func {{$.Table.Name}}_Index_{{range $column := $index.Columns}}{{$column.Name}}_{{end}}Key({{range $column := $index.Columns}}
	{{$column.Name | toArgName}} {{$column | dataTypeToGoTypeString}},{{end}}
) string { 
	return fmt.Sprintf("{{$.Table.Name | toArgName}}_idx{{range $column := $index.Columns}}_{{$column.Name | toArgName}}_{{$column | dataTypeToFormatString}}{{end}}",{{range $column := $index.Columns}}
		{{$column.Name | toArgName}},{{end}}	
	)
}{{end}}{{end}}

// {{.Table.Name}}Cache is a cache for {{.Table.Name}} objects
type {{.Table.Name}}Cache struct {
	cache redis.IRedis
}

// New{{.Table.Name}}Cache returns a new instance of {{.Table.Name}}Cache
func New{{.Table.Name}}Cache(
	cache redis.IRedis, 
) *{{.Table.Name}}Cache {
	return &{{.Table.Name}}Cache{
		cache,
	}
}

// FromID retrieves a {{.Table.Name}} model by its primary key
func (r *{{.Table.Name}}Cache) FromID(id int64) (*models.{{.Table.Name}}, error) { 
	
	var e error 

	var model = &models.{{.Table.Name}}{} 

	// Get the item by key 
	e = r.cache.Get({{.Table.Name}}_Key(id), model) 

	return model, e 
}

// Save creates or updates an existing {{.Table.Name}} model
func (r *{{.Table.Name}}Cache) Save(model *models.{{.Table.Name}}) { 
	
	r.cache.Set({{.Table.Name}}_Key(model.{{.PrimaryKey}}), model)

	// Primary Key 
	r.cache.ZAdd({{.Table.Name}}_Index_ID_Key, 0, {{.Table.Name}}_Key(model.{{.PrimaryKey}}))
	{{range $index := .CacheConfig.Indices}}{{ if not $index.Index.Unique }}
	// Index: {{$index.Index.Field}}
	r.AddIndex_{{$index.Columns | columnsToMethodName}}(model)
	{{else}}
	// Unique Index: {{$index.Index.Field}}
	r.AddUniqueIndex_{{$index.Columns | columnsToMethodName}}(model)
	{{end}}{{end}}
}

{{range $index := .CacheConfig.Indices}}{{ if not $index.Index.Unique }}
// AddIndex_{{$index.Columns | columnsToMethodName}} adds an index on the {{$index.Index.Field}} field(s)
func (r *{{$.Table.Name}}Cache) AddIndex_{{$index.Columns | columnsToMethodName}}(model *models.{{$.Table.Name}}) { 
	r.cache.ZAdd(
		{{$.Table.Name}}_Index_{{range $column := $index.Columns}}{{$column.Name}}_{{end}}Key({{range $column := $index.Columns}}
			model.{{$column.Name}},{{end}}
		), 
		0,
		{{$.Table.Name}}_Key(model.{{$.PrimaryKey}}),
	)
}

// RemoveIndex_{{$index.Columns | columnsToMethodName}} removes an index on the {{$index.Index.Field}} field
func (r *{{$.Table.Name}}Cache) RemoveIndex_{{$index.Columns | columnsToMethodName}}(model *models.{{$.Table.Name}}) { 
	r.cache.ZRem(
		{{$.Table.Name}}_Index_{{range $column := $index.Columns}}{{$column.Name}}_{{end}}Key({{range $column := $index.Columns}}
			model.{{$column.Name}},{{end}}
		), 
		{{$.Table.Name}}_Key(model.{{$.PrimaryKey}}),
	)
}
{{else}}
// AddUniqueIndex_{{$index.Columns | columnsToMethodName}} adds an index on the {{$index.Index.Field}} field
func (r *{{$.Table.Name}}Cache) AddUniqueIndex_{{$index.Columns | columnsToMethodName}}(model *models.{{$.Table.Name}}) { 
	r.cache.HSet(
		{{$.Table.Name}}_Index_{{$index.Columns | columnsToKey}}_Key,
		{{$index.Columns | columnModelValuesToKey}},
		{{$.Table.Name}}_Key(model.{{$.PrimaryKey}}),
	)
}

// RemoveUniqueIndex_{{$index.Columns | columnsToMethodName}} removes an index on the {{$index.Index.Field}} field
func (r *{{$.Table.Name}}Cache) RemoveUniqueIndex_{{$index.Columns | columnsToMethodName}}(model *models.{{$.Table.Name}}) { 
	r.cache.HDel(
		{{$.Table.Name}}_Index_{{$index.Columns | columnsToKey}}_Key, 
		{{$index.Columns | columnModelValuesToKey}}, 
	)
}
{{end}}{{end}}

// Delete removes a {{.Table.Name}} object from the cache
func (r *{{.Table.Name}}Cache) Delete(id int64) { 
	
	var e error 
	var model = &models.{{.Table.Name}}{} 

	if e = r.cache.Get({{.Table.Name}}_Key(id), model); e == nil { 

		// Delete the key 
		r.cache.Del({{.Table.Name}}_Key(id))

		// Primary Key 
		r.cache.ZRem({{.Table.Name}}_Index_ID_Key, {{.Table.Name}}_Key(id))
		{{range $index := .CacheConfig.Indices}}{{ if not $index.Index.Unique}}
		// Index: {{$index.Index.Field}}
		r.RemoveIndex_{{$index.Columns | columnsToMethodName}}(model)
		{{else}}
		// Unique Index: {{$index.Index.Field}}
		r.RemoveUniqueIndex_{{$index.Columns | columnsToMethodName}}(model)
		{{end}}{{end}}
	}
}

// All returns a slice of {{.Table.Name}} objects 
func (r *{{.Table.Name}}Cache) All(page, limit int64) ([]*models.{{.Table.Name}}, error) { 

	var e error 
	var keys []string 
	var collection = []*models.{{.Table.Name}}{} 

	if keys, e = r.cache.ZRangeByScore({{.Table.Name}}_Index_ID_Key, 0, 0, page*limit, limit); e == nil { 
		for k := range keys { 
			var model = &models.{{.Table.Name}}{} 
			if e = r.cache.Get(keys[k], model); e == nil { 
				collection = append(collection, model) 
			}
		}
	}

	return collection, e 
}

// Count returns a count of all {{.Table.Name}} objects 
func (r *{{.Table.Name}}Cache) Count() (int64, error) { 
	return r.cache.ZCard({{.Table.Name}}_Index_ID_Key)
}

{{range $index := .CacheConfig.Indices}}
{{if not $index.Index.Unique}}// From{{$index.Columns | columnsToMethodName}} returns a slice of {{$.Table.Name}} objects by their indexed field '{{$index.Index.Field}}'
func (r *{{$.Table.Name}}Cache) From{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodParams}}, page, limit int64) ([]*models.{{$.Table.Name}}, error) { 

	var e error 
	var keys []string 
	var collection = []*models.{{$.Table.Name}}{} 

	if keys, e = r.cache.ZRangeByScore({{$.Table.Name}}_Index_{{range $column := $index.Columns}}{{$column.Name}}_{{end}}Key({{range $column := $index.Columns}}{{$column.Name | toArgName}},{{end}}), 0, 0, limit*page, limit); e == nil { 
		for k := range keys { 
			var model = &models.{{$.Table.Name}}{} 
			if e = r.cache.Get(keys[k], model); e == nil { 
				collection = append(collection, model) 
			}
		}
	}

	return collection, e 
} 

// CountFrom{{$index.Columns | columnsToMethodName}} returns a count of items by {{$index.Index.Field}}
func (r *{{$.Table.Name}}Cache) CountFrom{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodParams}}) (int64, error) { 
	return r.cache.ZCard({{$.Table.Name}}_Index_{{range $column := $index.Columns}}{{$column.Name}}_{{end}}Key({{$index.Columns | columnsToMethodArgs}}))
}
{{else}}// From{{range $column := $index.Columns}}{{$column.Name}}{{end}} returns a single {{$.Table.Name}} by its unique field(s) '{{$index.Index.Field}}'
func (r *{{$.Table.Name}}Cache) From{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodParams}}) (*models.{{$.Table.Name}}, error) { 
	var key = r.cache.HGet({{$.Table.Name}}_Index_{{range $column := $index.Columns}}{{$column.Name}}_{{end}}Key, {{$index.Columns | columnValuesToKey}})

	if len(key) > 0 { 
		var model = &models.{{$.Table.Name}}{}
		r.cache.Get(key, model)
		return model, nil 
	}

	return nil, nil 
}
{{end}}
{{end}}
`))

func GenCaches(tables []*schema.Table, basePackage string, config map[string]*lib.CacheConfig) error {

	var tableMap = map[string]int{}
	for k := range tables {
		tableMap[tables[k].Name] = k
	}

	// start := time.Now()
	var generatedCacheCount = 0
	lib.EnsureDir(lib.CacheGenDir)
	// TODO Verbose flag
	// fmt.Println("Generating caches", config)
	for k := range config {

		var table = tables[tableMap[k]]

		if e := GenerateGoCache(basePackage, config[k], table, lib.CacheGenDir); e != nil {
			return e
		}
		generatedCacheCount++
	}
	// TODO Verbose flag
	// fmt.Printf("Generated %d caches in %f seconds.\n", generatedCacheCount, time.Since(start).Seconds())

	return nil
}

// GenerateGoDAL returns a string for a repo in golang
func GenerateGoCache(basePackage string, cacheConfig *lib.CacheConfig, table *schema.Table, dir string) (e error) {

	p := path.Join(dir, table.Name+"Cache.go")
	// TODO Verbose mode
	// fmt.Println("Generating Cache file to path: ", p)

	data := struct {
		BasePackage string
		Table       *schema.Table
		Columns     schema.SortedColumns

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
		CacheConfig       *CacheData
	}{
		BasePackage:       basePackage,
		HasNull:           false,
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
		CacheConfig:       ParseIndices(cacheConfig, table),
	}

	sortedColumns := make(schema.SortedColumns, 0, len(table.Columns))

	// var insertColumns = []*schema.Column{}

	// Find the primary key
	for k := range table.Columns {

		// fmt.Println("Column:", table.Columns[k].Name)

		var column = table.Columns[k]

		if column.ColumnKey == "PRI" {
			data.PrimaryKey = column.Name
			data.PrimaryKeyType = column.DataType
		}

		goDataType := schema.DataTypeToGoTypeString(column)
		if len(goDataType) > 5 && goDataType[0:5] == "null." {
			data.HasNull = true
		}

		sortedColumns = append(sortedColumns, column)

		// if isInsertColumn(column) {
		// 	insertColumns = append(insertColumns, column)
		// }

		sortedColumns = append(sortedColumns, column)
	}

	sort.Sort(sortedColumns)
	data.Columns = sortedColumns

	_, data.IsDeleted = table.Columns["IsDeleted"]
	_, data.IsDateCreated = table.Columns["DateCreated"]
	_, data.IsLastUpdated = table.Columns["LastUpdated"]

	switch data.PrimaryKeyType {
	case "varchar":
		data.IDType = "string"
	}

	var buf bytes.Buffer
	if e = CacheTemplate.Execute(&buf, data); e != nil {
		panic(e)
		return
	}

	var bufBytes = buf.Bytes()

	e = ioutil.WriteFile(p, bufBytes, lib.DefaultFileMode)
	if e != nil {
		panic(e)
	}

	return
}
