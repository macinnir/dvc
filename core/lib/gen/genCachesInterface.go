package gen

import (
	"bytes"
	"io/ioutil"
	"path"
	"text/template"
	"time"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

var CacheInterfaceTemplate = template.Must(template.New("template-cache-interface-file").Funcs(template.FuncMap{
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
	"{{ .BasePackage }}/gen/definitions/models" 
)

// I{{.Table.Name}}Cache is an interface for the {{.Table.Name}} cache
type I{{.Table.Name}}Cache interface {
	FromID(id int64) (*models.{{.Table.Name}}, error) 
	Save(model *models.{{.Table.Name}}){{range $index := .CacheConfig.Indices}}{{ if not $index.Index.Unique }}
	AddIndex_{{$index.Columns | columnsToMethodName}}(model *models.{{$.Table.Name}}) 
	RemoveIndex_{{$index.Columns | columnsToMethodName}}(model *models.{{$.Table.Name}}) {{else}}
	AddUniqueIndex_{{$index.Columns | columnsToMethodName}}(model *models.{{$.Table.Name}}) 
	RemoveUniqueIndex_{{$index.Columns | columnsToMethodName}}(model *models.{{$.Table.Name}}){{end}}{{end}}
	Delete(id int64) 
	All(page, limit int64) ([]*models.{{.Table.Name}}, error)
	Count() (int64, error) 
	{{range $index := .CacheConfig.Indices}}{{ if not $index.Index.Unique}}From{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodParams}}, page, limit int64) ([]*models.{{$.Table.Name}}, error)
	CountFrom{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodParams}}) (int64, error)
	{{else}}From{{$index.Columns | columnsToMethodName}}({{$index.Columns | columnsToMethodParams}}) (*models.{{$.Table.Name}}, error)
	{{end}}{{end}}
}

`))

func GenCacheInterfaces(basePackage string, tables []*schema.Table, cache map[string]*lib.CacheConfig) error {

	var tableMap = map[string]int{}
	for k := range tables {
		tableMap[tables[k].Name] = k
	}

	start := time.Now()
	var generatedCacheCount = 0
	lib.EnsureDir(lib.CacheInterfaceGenDir)

	for tableName := range cache {

		var table = tables[tableMap[tableName]]

		if e := GenerateGoCacheInterface(basePackage, cache[tableName], table, lib.CacheInterfaceGenDir); e != nil {
			return e
		}
		generatedCacheCount++
	}

	lib.LogAdd(start, "%d cache interfaces", generatedCacheCount)

	return nil
}

// var dalTPL *template.Template

// GenerateGoDAL returns a string for a repo in golang
func GenerateGoCacheInterface(basePackage string, cacheConfig *lib.CacheConfig, table *schema.Table, dir string) (e error) {

	p := path.Join(dir, "I"+table.Name+"Cache.go")
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
	if e = CacheInterfaceTemplate.Execute(&buf, data); e != nil {
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
