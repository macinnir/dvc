package genutil

import (
	"strings"

	"github.com/macinnir/dvc/core/lib/schema"
)

func IsInsertColumn(column *schema.Column) bool {
	if column.ColumnKey == "PRI" {
		return false
	}

	if column.Name == "IsDeleted" {
		return false
	}

	if column.DataType == "vector" {
		return false
	}

	return true
}

func FetchInsertColumns(columns schema.SortedColumns) []*schema.Column {

	insertColumns := []*schema.Column{}

	for _, column := range columns {
		if !IsInsertColumn(column) {
			continue
		}

		insertColumns = append(insertColumns, column)
	}

	return insertColumns
}

func IsUpdateColumn(column *schema.Column) bool {
	if column.ColumnKey == "PRI" {
		return false
	}

	if column.Name == "IsDeleted" {
		return false
	}

	if column.Name == "DateCreated" {
		return false
	}

	if column.DataType == "vector" {
		return false
	}

	return true
}

func IsSpecialColumn(column *schema.Column) bool {
	if column.DataType == "vector" {
		return true
	}
	return false
}

func FetchUpdateColumns(columns schema.SortedColumns) []*schema.Column {

	updateColumns := []*schema.Column{}

	for _, column := range columns {
		if !IsUpdateColumn(column) {
			continue
		}

		updateColumns = append(updateColumns, column)
	}

	return updateColumns
}

func FetchTableUpdateFieldsString(columns schema.SortedColumns) string {
	fields := []string{}
	for _, field := range columns {

		if !IsUpdateColumn(field) {
			continue
		}

		fields = append(fields, "`"+field.Name+"` = :"+field.Name)
	}

	return strings.Join(fields, ",")
}

func FetchTablePrimaryKey(table *schema.Table) string {
	primaryKey := ""
	idType := "int64"
	for _, column := range table.Columns {
		if column.ColumnKey == "PRI" {
			primaryKey = column.Name
		}
	}

	return primaryKey + " " + idType
}

func FetchTablePrimaryKeyName(table *schema.Table) string {
	primaryKey := ""
	for _, column := range table.Columns {
		if column.ColumnKey == "PRI" {
			primaryKey = column.Name
		}
	}

	return primaryKey
}

func FetchTableInsertFieldsString(columns schema.SortedColumns) string {

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

func FetchTableInsertValuesString(columns schema.SortedColumns) string {
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
