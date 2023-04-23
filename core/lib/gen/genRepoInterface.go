package gen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"text/template"
	"time"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

var RepoInterfaceTemplate = template.Must(template.New("template-repo-interface-file").Funcs(template.FuncMap{
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
	"{{ .BasePackage }}/gen/definitions/models" 
	"{{ .BasePackage }}/gen/definitions/collections"{{ if gt (len .CacheConfig.Location) 0}}
	"{{ .BasePackage }}/{{ .CacheConfig.Location }}"{{end}}
)

// I{{.Table.Name}}Repo is an interface for the {{.Table.Name}} repo
type I{{.Table.Name}}Repo interface {
	FromID(id int64, mustExist bool) (*models.{{.Table.Name}}, error){{if .CacheConfig.HasHashID}}
	FromHashID(hashID string, mustExist bool) (*models.{{.Table.Name}}, error){{end}}{{ if gt (len .CacheConfig.Properties) 0 }}
	AggregateFromID(id int64, mustExist bool) (*aggregates.{{.Table.Name}}Aggregate, error){{end}}
	Reset({{.PrimaryKey | toArgName}} int64) error
	Create(model *models.{{.Table.Name}}) error
	CreateMany(modelSlice []*models.{{.Table.Name}}) error
	Update(model *models.{{.Table.Name}}) error
	UpdateMany(modelSlice []*models.{{.Table.Name}}) error
	Delete(id int64) error
	DeleteMany(modelSlice []*models.{{.Table.Name}}) error
	All(page, limit int64) ([]*models.{{.Table.Name}}, error)
	AllAsCollection(page, limit int64) (*collections.{{.Table.Name}}Collection, error)

	{{range $index := .CacheConfig.Indices}}{{ if not $index.Index.Unique}}From{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodParams}}, page, limit int64) ([]*models.{{$.Table.Name}}, error)
	CollectionFrom{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodParams}}, page, limit int64) (*collections.{{$.Table.Name}}Collection, error){{else}}From{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodParams}}, mustExist bool) (*models.{{$.Table.Name}}, error){{if gt (len $.CacheConfig.Properties) 0}}
	AggregateFrom{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodParams}}, mustExist bool) (*aggregates.{{$.Table.Name}}Aggregate, error){{end}}{{end}}
	{{end}}{{range $index := .CacheConfig.Indices}}{{ if not $index.Index.Unique }}
	AddIndex_{{$index.Columns | columnsToMethodName}}(model *models.{{$.Table.Name}}) 
	RemoveIndex_{{$index.Columns | columnsToMethodName}}(model *models.{{$.Table.Name}}) {{else}}
	AddUniqueIndex_{{$index.Columns | columnsToMethodName}}(model *models.{{$.Table.Name}}) 
	RemoveUniqueIndex_{{$index.Columns | columnsToMethodName}}(model *models.{{$.Table.Name}}){{end}}{{end}}
	{{if gt (len .CacheConfig.Search) 0}}{{range $search := .CacheConfig.Search}}
	// Search{{$search.SearchColumns | columnsToMethodName}} searches the {{range $col := $search.SearchColumns}}{{$col.Name}}{{end}} 
	// leftOrRightOrCenter is 2 == Center, 1 == Right, 0 (default) == Left
	Search{{$search.SearchColumns | columnsToMethodName}}({{$search.ConditionColumns | columnsToMethodParams}}{{if gt (len $search.ConditionColumns) 0}},{{end}}{{ $search.SearchColumns | columnsToMethodParams }}, leftOrRightOrCenter, page, limit int64) ([]*models.{{$.Table.Name}}, error){{end}}{{end}} 
}

`))

func GenRepoInterfaces(basePackage string, tables []*schema.Table, cache map[string]*lib.CacheConfig) error {

	var tableMap = map[string]int{}
	for k := range tables {
		tableMap[tables[k].Name] = k
	}

	start := time.Now()
	var generatedCacheCount = 0
	lib.EnsureDir(lib.RepoInterfaceGenDir)

	for tableName := range cache {

		var table = tables[tableMap[tableName]]

		if e := GenerateGoRepoInterface(basePackage, cache[tableName], table, lib.RepoInterfaceGenDir); e != nil {
			return e
		}
		generatedCacheCount++
	}
	fmt.Printf("Generated %d repo interfaces in %f seconds.\n", generatedCacheCount, time.Since(start).Seconds())

	return nil
}

// var dalTPL *template.Template

// GenerateGoDAL returns a string for a repo in golang
func GenerateGoRepoInterface(basePackage string, cacheConfig *lib.CacheConfig, table *schema.Table, dir string) (e error) {

	p := path.Join(dir, "I"+table.Name+"Repo.go")
	fmt.Println("Generating Cache file to path: ", p)

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
	}

	data.Columns = sortedColumns

	_, data.IsDeleted = table.Columns["IsDeleted"]
	_, data.IsDateCreated = table.Columns["DateCreated"]
	_, data.IsLastUpdated = table.Columns["LastUpdated"]

	switch data.PrimaryKeyType {
	case "varchar":
		data.IDType = "string"
	}

	var buf bytes.Buffer
	if e = RepoInterfaceTemplate.Execute(&buf, data); e != nil {
		panic(e)
		return
	}

	var bufBytes = buf.Bytes()
	// fmt.Println("Cache file: ", string(bufBytes))

	e = ioutil.WriteFile(p, bufBytes, lib.DefaultFileMode)
	if e != nil {
		panic(e)
	}

	return
}
