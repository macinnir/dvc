package model

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/gen/genutil"
	"github.com/macinnir/dvc/core/lib/schema"
)

// 0.607244

// GenModels generates models
func GenModels(tables []*schema.Table, config *lib.Config) error {

	// start := time.Now()

	generatedModelCount := 0

	lib.EnsureDir(lib.ModelsGenDir)

	for k := range tables {

		var table = tables[k]

		fullPath := path.Join(lib.ModelsGenDir, table.Name+".go")
		if e := buildGoModel(config.BasePackage, fullPath, table); e != nil {
			return e
		}
		generatedModelCount++
	}

	// TODO Verbose flag
	// fmt.Printf("Generated %d models in %f seconds.\n", generatedModelCount, time.Since(start).Seconds())
	return nil
}

func buildGoModel(packageName, fullPath string, table *schema.Table) (e error) {
	// var modelNode *lib.GoStruct
	var outFile []byte
	outFile, e = buildFileFromModelNode(table)
	if e != nil {
		fmt.Println("ERROR Building File From Model Node ", table, e.Error())
		return
	}

	os.WriteFile(fullPath, outFile, lib.DefaultFileMode)
	return
}

// buildModelNodeFromFile builds a node representation of a struct from a file
func buildModelNodeFromTable(table *schema.Table) (*lib.GoStruct, error) {

	var modelNode = lib.NewGoStruct()
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

		if schema.IsNull(fieldType) {
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
		modelNode.Imports.Append(lib.NullPackage)
	}

	return modelNode, nil
}

type GoModelTemplateVals struct {
	Name          string
	Schema        string
	HasNull       bool
	HasAccountID  bool
	HasUserID     bool
	UpdateColumns []GoModelTemplateFieldVal
	InsertColumns []GoModelTemplateFieldVal
	PrimaryKey    string
	Fields        []GoModelTemplateFieldVal
	SelectFields  []GoModelTemplateFieldVal
}

// schema.GoTypeFormatString
type GoModelTemplateFieldVal struct {
	Name       string
	Type       string
	DBType     string
	GoType     string
	FormatType string
}

func buildFileFromModelNode(table *schema.Table) ([]byte, error) {

	var vals = GoModelTemplateVals{
		Name:         table.Name,
		Schema:       table.SchemaName,
		Fields:       make([]GoModelTemplateFieldVal, len(table.Columns)),
		SelectFields: []GoModelTemplateFieldVal{},
	}

	var sortedColumns = make([]string, len(table.Columns))

	var k = 0
	for colName := range table.Columns {
		sortedColumns[k] = table.Columns[colName].Name
		k++
	}

	sort.Strings(sortedColumns)

	for l := range sortedColumns {

		col := table.Columns[sortedColumns[l]]

		if col.Name == "AccountID" {
			vals.HasAccountID = true
		}

		if col.Name == "UserID" {
			vals.HasUserID = true
		}

		if col.ColumnKey == "PRI" {
			vals.PrimaryKey = col.Name
		}

		field := GoModelTemplateFieldVal{
			Name:   sortedColumns[l],
			Type:   col.Type,
			DBType: col.DataType,
			GoType: schema.DataTypeToGoTypeString(col),
		}

		if field.GoType == "null.String" || field.GoType == "null.Float" {
			vals.HasNull = true
		}

		field.FormatType = schema.GoTypeFormatString(field.GoType)

		vals.Fields[l] = field

		if genutil.IsInsertColumn(col) {
			vals.InsertColumns = append(vals.InsertColumns, field)
		}

		if genutil.IsUpdateColumn(col) {
			vals.UpdateColumns = append(vals.UpdateColumns, field)
		}

		if !genutil.IsSpecialColumn(col) {
			vals.SelectFields = append(vals.SelectFields, field)
		}

	}

	var e error
	var buf bytes.Buffer

	if e = goModelTemplate.Execute(&buf, vals); e != nil {
		return nil, e
	}
	// 0.025663
	return buf.Bytes(), nil

}
