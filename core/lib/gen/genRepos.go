package gen

import (
	"bytes"
	"io/ioutil"
	"path"
	"sort"
	"sync"
	"text/template"
	"time"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

var RepoTemplate = template.Must(template.New("template-repo-file").Funcs(template.FuncMap{
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

package repos

import ( 
	"{{ .BasePackage }}/gen/definitions/collections" 
	"{{ .BasePackage }}/gen/definitions/models" 
	"{{ .BasePackage }}/gen/definitions/caches" 
	"{{ .BasePackage }}/gen/definitions/dal" 
	"{{ .BasePackage }}/core/components/config" 
	{{ if .CacheConfig.HasHashID }}"{{ .BasePackage }}/core/utils/hashids"{{end}}{{ if gt (len .CacheConfig.Location) 0}}
	"{{ .BasePackage }}/{{ .CacheConfig.Location }}"{{end}}{{if gt (len .CacheConfig.Search) 0}}
	"github.com/macinnir/dvc/core/lib/utils/query"{{end}}

	"fmt"
)


// {{.Table.Name}}Repo is a repo for {{.Table.Name}} objects
type {{.Table.Name}}Repo struct {
	config *config.Config
	{{.Table.Name | toArgName}}Cache caches.I{{.Table.Name}}Cache
	{{.Table.Name | toArgName}}DAL dal.I{{.Table.Name}}DAL
	{{ if .CacheConfig.HasHashID }}idHasher             *hashids.IDHasher{{end}}{{ if gt (len .CacheConfig.Properties) 0}}{{range $agg := .CacheConfig.Properties}}
	{{$agg.Aggregate.Table | toArgName}}Repo *{{$agg.Aggregate.Table}}Repo{{end}}{{end}}
}

// New{{.Table.Name}}Repo returns a new instance of {{.Table.Name}}Repo
func New{{.Table.Name}}Repo(
	config *config.Config, 
	{{.Table.Name | toArgName}}Cache caches.I{{.Table.Name}}Cache,
	{{.Table.Name | toArgName}}DAL dal.I{{.Table.Name}}DAL, 
	{{ if .CacheConfig.HasHashID }}idHasher             *hashids.IDHasher,{{end}}{{ if gt (len .CacheConfig.Properties) 0}}{{range $agg := .CacheConfig.Properties}}
	{{$agg.Aggregate.Table | toArgName}}Repo *{{$agg.Aggregate.Table}}Repo,{{end}}{{end}}
) *{{.Table.Name}}Repo {
	return &{{.Table.Name}}Repo{
		config, 
		{{.Table.Name | toArgName}}Cache,
		{{.Table.Name | toArgName}}DAL, 
		{{ if .CacheConfig.HasHashID }}idHasher,{{end}}{{ if gt (len .CacheConfig.Properties) 0}}{{range $agg := .CacheConfig.Properties}}
		{{$agg.Aggregate.Table | toArgName}}Repo,{{end}}{{end}}
	}
}

// FromID retrieves a {{.Table.Name}} model by its primary key
func (r *{{.Table.Name}}Repo) FromID(id int64, mustExist bool) (*models.{{.Table.Name}}, error) { 
	
	var model *models.{{.Table.Name}}
	var e error 
	
	if model, e = r.{{.Table.Name | toArgName}}Cache.FromID(id); e != nil { 
		return nil, e 
	}

	if model != nil { 
		return model, nil 
	}


	if model, e = r.{{.Table.Name | toArgName}}DAL.FromID(config.DEFAULT_SHARD, id, mustExist); e != nil { 
		return nil, e 
	}

	if model != nil { 
		r.{{.Table.Name | toArgName}}Cache.Save(model)
	}

	return model, e 
}

// Reset deletes a {{.Table.Name}} object from the cache and resaves it (if it exists).
func (r *{{.Table.Name}}Repo) Reset({{.PrimaryKey | toArgName}} int64) error { 

	var e error 
	var model *models.{{.Table.Name}}

	if model, e = r.{{.Table.Name | toArgName}}DAL.FromID(config.DEFAULT_SHARD, {{.PrimaryKey | toArgName}}, false); e != nil { 
		return e 
	}

	// No model in the database, so delete it from the cache
	if model == nil { 
		r.{{.Table.Name | toArgName}}Cache.Delete({{.PrimaryKey | toArgName}})
		return nil 
	}

	// Reset the cache 
	r.{{.Table.Name | toArgName}}Cache.Save(model) 
	return nil 
}

{{ if .CacheConfig.HasHashID }}// FromHashID returns a {{.Table.Name}} object based on its unique hashID
func (r *{{.Table.Name}}Repo) FromHashID(hashID string, mustExist bool) (*models.{{.Table.Name}}, error) { 

	var id = r.idHasher.FromHashID(hashID)

	return r.FromID(id, mustExist) 
}
{{ end }}{{ if gt (len .CacheConfig.Properties) 0 }}
// AggregateFromID returns a {{.Table.Name}}Aggregate object 
func (r *{{.Table.Name}}Repo) AggregateFromID(id int64, mustExist bool) (*aggregates.{{.Table.Name}}Aggregate, error) { 

	var model *models.{{.Table.Name}}
	var e error 

	if model, e = r.FromID(id, true); e != nil { 
		return nil, e 
	}

	var agg = &aggregates.{{.Table.Name}}Aggregate{ 
		{{.Table.Name}}: model, 
	}
	{{range $agg := .CacheConfig.Properties}}
	if agg.{{$agg.Aggregate.Property}}, e =  r.{{$agg.Aggregate.Table | toArgName}}Repo.From{{ if eq $agg.Aggregate.Type "Many"}}{{$agg.Aggregate.On}}{{ else }}ID{{ end }}(model.{{$agg.Aggregate.On}}, {{ if eq $agg.Aggregate.Type "Many" }}0, 0{{ else }}true{{end}}); e != nil { 
		return nil, e 
	}{{end}}

	return agg, nil
}
{{ end }}
// Create creates or updates an existing {{.Table.Name}} model
func (r *{{.Table.Name}}Repo) Create(model *models.{{.Table.Name}}) error { 
	
	var e error 

	if e = r.{{.Table.Name | toArgName}}DAL.Create(config.DEFAULT_SHARD, model); e != nil { 
		return e 
	}

	r.{{.Table.Name | toArgName}}Cache.Save(model) 
	
	return e 
}

// CreateMany creates a slice of {{.Table.Name}} objects 
func (r *{{.Table.Name}}Repo) CreateMany(modelSlice []*models.{{.Table.Name}}) error { 
	
	var  e error 

	if e = r.{{.Table.Name | toArgName}}DAL.CreateMany(config.DEFAULT_SHARD, modelSlice); e != nil { 
		return e 
	}

	for k := range modelSlice { 
		r.{{.Table.Name | toArgName}}Cache.Save(modelSlice[k]) 
	}

	return nil 
}

// Update updates an existing {{.Table.Name}} model
func (r *{{.Table.Name}}Repo) Update(model *models.{{.Table.Name}}) error { 

	var e = r.{{.Table.Name | toArgName}}DAL.Update(config.DEFAULT_SHARD, model) 
	
	if e == nil { 
		r.{{.Table.Name | toArgName}}Cache.Save(model) 
	}
	
	return e 
}

// UpdateMany updates a slice of {{.Table.Name}} objects 
func (r *{{.Table.Name}}Repo) UpdateMany(modelSlice []*models.{{.Table.Name}}) error { 

	for k := range modelSlice { 
		r.{{.Table.Name | toArgName}}Cache.Save(modelSlice[k]) 
	}

	return r.{{.Table.Name | toArgName}}DAL.UpdateMany(config.DEFAULT_SHARD, modelSlice)
}

// Delete removes a {{.Table.Name}} object from the cache
func (r *{{.Table.Name}}Repo) Delete(id int64) error { 
	r.{{.Table.Name | toArgName}}Cache.Delete(id) 
	return r.{{.Table.Name | toArgName}}DAL.Delete(config.DEFAULT_SHARD, id)
}

// DeleteMany deletes a slice of {{.Table.Name}} objects 
func (r *{{.Table.Name}}Repo) DeleteMany(modelSlice []*models.{{.Table.Name}}) error { 

	for k := range modelSlice { 
		r.{{.Table.Name | toArgName}}Cache.Delete(modelSlice[k].{{.PrimaryKey}}) 
	}

	return r.{{.Table.Name | toArgName}}DAL.DeleteMany(config.DEFAULT_SHARD, modelSlice)
}

// All returns a slice of {{.Table.Name}} objects 
func (r *{{.Table.Name}}Repo) All(page, limit int64) ([]*models.{{.Table.Name}}, error) { 
	
	var e error 
	var items []*models.{{.Table.Name}}

	if items, e = r.{{$.Table.Name | toArgName}}Cache.All(page, limit); e != nil { 
		return nil, e 
	}

	if len(items) == 0 { 
		
		if items, e = r.{{$.Table.Name | toArgName}}DAL.ManyPaged(config.DEFAULT_SHARD, limit, page*limit, "", ""); e != nil {
			return nil, fmt.Errorf("{{$.Table.Name}}Repo::All() -> {{$.Table.Name}}DAL.ManyPaged(): %w", e)
		}

		if len(items) > 0 { 
			for k := range items { 
				r.{{$.Table.Name | toArgName}}Cache.Save(items[k])
			}
		}
	}

	return items, nil 
}

// AllAsCollection returns a collection of {{.Table.Name}} objects 
func (r *{{.Table.Name}}Repo) AllAsCollection(page, limit int64) (*collections.{{.Table.Name}}Collection, error) { 
	
	var e error 
	var collection = &collections.{{.Table.Name}}Collection{}

	if collection.Data, e = r.{{$.Table.Name | toArgName}}Cache.All(page, limit); e != nil { 
		return nil, e 
	}

	if len(collection.Data) == 0 { 
		
		if collection.Data, e = r.{{$.Table.Name | toArgName}}DAL.ManyPaged(config.DEFAULT_SHARD, limit, page*limit, "", ""); e != nil {
			return nil, fmt.Errorf("{{$.Table.Name}}Repo::AllAsCollection() -> {{$.Table.Name}}DAL.ManyPaged(): %w", e)
		}

		if len(collection.Data) > 0 { 
			for k := range collection.Data { 
				r.{{$.Table.Name | toArgName}}Cache.Save(collection.Data[k])
			}
		}

		if collection.Count, e = r.{{$.Table.Name | toArgName}}DAL.Count(config.DEFAULT_SHARD).Run(); e != nil { 
			return nil, fmt.Errorf("{{$.Table.Name}}Repo::AllAsCollection() -> {{$.Table.Name}}DAL.Count(): %w", e)
		}
	}

	if collection.Count, e = r.{{$.Table.Name | toArgName}}Cache.Count(); e != nil { 
		return nil, e 
	}


	return collection, nil 
}
{{range $index := .CacheConfig.Indices}}
{{if not $index.Index.Unique}}// From{{$index.Columns | columnsToMethodName}} returns a collection of {{$.Table.Name}} objects by their indexed field '{{$index.Index.Field}}'
func (r *{{$.Table.Name}}Repo) From{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodParams}}, page, limit int64) ([]*models.{{$.Table.Name}}, error) { 

	var e error 
	var items = []*models.{{$.Table.Name}}{} 

	if items, e = r.{{$.Table.Name | toArgName}}Cache.From{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodArgs}}, page, limit); e != nil { 
		return nil, e 
	}

	if len(items) == 0 { 
		
		if items, e = r.{{$.Table.Name | toArgName}}DAL.ManyFrom{{$index.Columns | columnsToMethodName}}(config.DEFAULT_SHARD, {{$index.Columns | columnsToMethodArgs}}, limit, page*limit, "", ""); e != nil {
			return nil, fmt.Errorf("{{$.Table.Name}}Repo::All() -> {{$.Table.Name}}DAL.ManyPaged(): %w", e)
		}

		if len(items) > 0 { 
			for k := range items { 
				r.{{$.Table.Name | toArgName}}Cache.Save(items[k])
			}
		}
	}

	return items, nil
} 

// CollectionFrom{{$index.Columns | columnsToMethodName}} returns a collection of {{$.Table.Name}} objects by their indexed field '{{$index.Index.Field}}'
func (r *{{$.Table.Name}}Repo) CollectionFrom{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodParams}}, page, limit int64) (*collections.{{$.Table.Name}}Collection, error) { 

	var e error 
	var collection = &collections.{{$.Table.Name}}Collection{} 

	if collection.Data, e = r.From{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodArgs}}, page, limit); e != nil { 
		return nil, e 
	}

	if collection.Count, e = r.{{$.Table.Name | toArgName}}Cache.CountFrom{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodArgs}}); e != nil { 
		return nil, e 
	}

	return collection, nil
}
{{else}}// From{{$index.Columns | columnsToMethodName}} returns a single {{$.Table.Name}} by its unique field(s) '{{$index.Index.Field}}'
func (r *{{$.Table.Name}}Repo) From{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodParams}}, mustExist bool) (*models.{{$.Table.Name}}, error) { 
	
	var model, e = r.{{$.Table.Name | toArgName}}Cache.From{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodArgs}}) 
	
	if e != nil || model == nil { 
	
		if model, e = r.{{$.Table.Name | toArgName}}DAL.SingleFrom{{$index.Columns | columnsToMethodName}}(config.DEFAULT_SHARD, {{$index.Columns | columnsToMethodArgs}}, mustExist); e != nil {
			return nil, e 
		}

		if model != nil { 
			r.{{$.Table.Name | toArgName}}Cache.Save(model) 
		}
	}

	return model, e
}
{{ if gt (len $.CacheConfig.Properties) 0 }}
// AggregateFrom{{$index.Columns | columnsToMethodName}} returns a {{$.Table.Name}}Aggregate object by its unique field(s) '{{$index.Index.Field}}'
func (r *{{$.Table.Name}}Repo) AggregateFrom{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodParams}}, mustExist bool) (*aggregates.{{$.Table.Name}}Aggregate, error) { 

	var model, e = r.From{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodArgs}}, mustExist) 
	if e != nil { 
		return nil, e 
	}
	
	var agg = &aggregates.{{$.Table.Name}}Aggregate{ 
		{{$.Table.Name}}: model, 
	}
	{{range $agg := $.CacheConfig.Properties}}
	if agg.{{$agg.Aggregate.Property}}, e =  r.{{$agg.Aggregate.Table | toArgName}}Repo.From{{ if eq $agg.Aggregate.Type "Many"}}{{$agg.Aggregate.On}}{{ else }}ID{{ end }}(model.{{$agg.Aggregate.On}}, {{ if eq $agg.Aggregate.Type "Many" }}0, 0{{ else }}true{{end}}); e != nil { 
		return nil, e 
	}
	{{end}}
	return agg, nil
}{{end}}{{end}}{{end}}
{{range $index := .CacheConfig.Indices}}{{ if not $index.Index.Unique }}
// AddIndex_{{$index.Columns | columnsToMethodName}} adds an index on the {{$index.Index.Field}} field(s)
func (r *{{$.Table.Name}}Repo) AddIndex_{{$index.Columns | columnsToMethodName}}(model *models.{{$.Table.Name}}) { 
	r.{{$.Table.Name | toArgName}}Cache.AddIndex_{{$index.Columns | columnsToMethodName}}(model)
}

// RemoveIndex_{{$index.Columns | columnsToMethodName}} removes an index on the {{$index.Index.Field}} field
func (r *{{$.Table.Name}}Repo) RemoveIndex_{{$index.Columns | columnsToMethodName}}(model *models.{{$.Table.Name}}) { 
	r.{{$.Table.Name | toArgName}}Cache.RemoveIndex_{{$index.Columns | columnsToMethodName}}(model)
}{{else}}

// AddUniqueIndex_{{$index.Columns | columnsToMethodName}} adds an index on the {{$index.Index.Field}} field
func (r *{{$.Table.Name}}Repo) AddUniqueIndex_{{$index.Columns | columnsToMethodName}}(model *models.{{$.Table.Name}}) { 
	r.{{$.Table.Name | toArgName}}Cache.AddUniqueIndex_{{$index.Columns | columnsToMethodName}}(model)
}

// RemoveUniqueIndex_{{$index.Columns | columnsToMethodName}} removes an index on the {{$index.Index.Field}} field
func (r *{{$.Table.Name}}Repo) RemoveUniqueIndex_{{$index.Columns | columnsToMethodName}}(model *models.{{$.Table.Name}}) { 
	r.{{$.Table.Name | toArgName}}Cache.RemoveUniqueIndex_{{$index.Columns | columnsToMethodName}}(model)
}{{end}}{{end}}
{{if and (.CacheConfig.Search) (gt (len .CacheConfig.Search) 0)}}{{range $search := .CacheConfig.Search}}
// Search{{$search.SearchColumns | columnsToMethodName}} searches the {{range $col := $search.SearchColumns}}{{$col.Name}},{{end}} column(s) 
// leftOrRightOrCenter is 2 == Center, 1 == Right, 0 (default) == Left
func (r *{{$.Table.Name}}Repo)Search{{$search.SearchColumns | columnsToMethodName}}({{$search.ConditionColumns | columnsToMethodParams}}{{if gt (len $search.ConditionColumns) 0}},{{end}}{{ $search.SearchColumns | columnsToMethodParams }}, leftOrRightOrCenter, page, limit int64) ([]*models.{{$.Table.Name}}, error) { 
	
	var e error
	var model []*models.{{$.Table.Name}}

	q := r.{{$.Table.Name | toArgName}}DAL.Select(config.DEFAULT_SHARD).Where(
		query.EQ(models.{{$.Table.Name}}_Column_IsDeleted, 0), 
		query.And(), 
	)

	switch leftOrRightOrCenter { 
		case 2: 
			// Search both 
			{{range $col := $search.SearchColumns}}{{$col.Name | toArgName}} = "%" + {{$col.Name | toArgName}} + "%"
			{{end}}
		case 1: 
			// Search right 
			{{range $col := $search.SearchColumns}}{{$col.Name | toArgName}} = "%" + {{$col.Name | toArgName}}
			{{end}}
		default:  
			// Search left (0)
			{{range $col := $search.SearchColumns}}{{$col.Name | toArgName}} = {{$col.Name | toArgName}} + "%"
			{{end}}
	} 
	{{if gt (len $search.ConditionColumns) 0}}
	// Conditions
	q.Where(
		query.Ands({{range $col := $search.ConditionColumns}}
			query.EQ(models.{{$.Table.Name}}_Column_{{$col.Name}}, {{ $col.Name | toArgName }}),{{end}}
		),
	)
	{{end}}
	// Search Fields
	q.Where(
		query.Ors({{range $col := $search.SearchColumns}}
			query.Like(models.{{$.Table.Name}}_Column_{{$col.Name}}, {{ $col.Name | toArgName }}),{{end}}
		),
	)
	
	if limit > 0 { 
		q.Limit(limit, limit*page)
	}

	if model, e = q.Run(); e != nil { 
		return nil, e 
	}

	return model, nil 
}
{{end}}{{end}}

// TODO Search with conditions (e.g. AccountID, UserID)
// TODO Search OR (e.g. Title OR Description)
// TODO Search left/all/right search
`))

func GenRepos(basePackage string, tables []*schema.Table, cache map[string]*lib.CacheConfig) error {

	var tableMap = map[string]int{}
	for k := range tables {
		tableMap[tables[k].Name] = k
	}

	start := time.Now()
	var generatedRepoCount = 0
	lib.EnsureDir(lib.RepoGenDir)
	lib.EnsureDir(lib.CollectionGenDir)

	var wg sync.WaitGroup
	mutex := &sync.Mutex{}

	for tableName := range cache {

		wg.Add(1)

		go func(tableName string) {

			defer wg.Done()

			var table = tables[tableMap[tableName]]

			if e := GenerateGoRepo(basePackage, cache[tableName], table, lib.RepoGenDir); e != nil {
				panic(e)
			}

			GenerateRepoCollection(basePackage, tableName)
			GenerateRepoCollectionItem(basePackage, tableName)
			mutex.Lock()
			generatedRepoCount++
			mutex.Unlock()
		}(tableName)
	}

	wg.Wait()

	lib.LogAdd(start, "%d repos", generatedRepoCount)

	return nil
}

// GenerateGoDAL returns a string for a repo in golang
func GenerateGoRepo(basePackage string, cacheConfig *lib.CacheConfig, table *schema.Table, dir string) (e error) {

	p := path.Join(dir, table.Name+"Repo.go")
	// fmt.Println("Generating Repo file to path: ", p)

	data := struct {
		BasePackage   string
		Table         *schema.Table
		Columns       schema.SortedColumns
		StringColumns []*schema.Column

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
		Table:             table,
		StringColumns:     []*schema.Column{},
		HasNull:           false,
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

		// TODO this is all helper code repeated for genDAL, etc.
		var column = table.Columns[k]

		if column.ColumnKey == "PRI" {
			data.PrimaryKey = column.Name
			data.PrimaryKeyType = column.DataType
		}

		if schema.DataTypeToGoTypeString(column) == "string" {
			data.StringColumns = append(data.StringColumns, column)
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
	if e = RepoTemplate.Execute(&buf, data); e != nil {
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
