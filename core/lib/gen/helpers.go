package gen

import (
	"fmt"
	"strings"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

func columnsToMethodName(columns []*schema.Column) string {

	if len(columns) == 0 {
		return ""
	}

	if len(columns) == 1 {
		return columns[0].Name
	}

	var sb strings.Builder
	for k := range columns {
		sb.WriteString(columns[k].Name)
		if k < len(columns)-1 {
			sb.WriteString("And")
		}
	}

	return sb.String()
}

func columnsToMethodArgs(columns []*schema.Column) string {

	if len(columns) == 0 {
		return ""
	}

	if len(columns) == 1 {
		return toArgName(columns[0].Name)
	}

	var sb strings.Builder
	for k := range columns {
		sb.WriteString(toArgName(columns[k].Name))
		if k < len(columns)-1 {
			sb.WriteString(",")
		}
	}

	return sb.String()
}

func columnsToMethodParams(columns []*schema.Column) string {

	if len(columns) == 0 {
		return ""
	}

	if len(columns) == 1 {
		return toArgName(columns[0].Name) + " " + schema.DataTypeToGoTypeString(columns[0])
	}

	var sb strings.Builder
	for k := range columns {
		sb.WriteString(toArgName(columns[k].Name) + " " + schema.DataTypeToGoTypeString(columns[k]))
		if k < len(columns)-1 {
			sb.WriteString(",")
		}
	}

	return sb.String()
}

func columnsToKey(columns []*schema.Column) string {

	if len(columns) == 0 {
		return ""
	}

	if len(columns) == 1 {
		return columns[0].Name
	}

	var sb strings.Builder
	for k := range columns {
		sb.WriteString(columns[k].Name)
		if k < len(columns)-1 {
			sb.WriteString("_")
		}
	}

	return sb.String()
}

func columnModelValuesToKey(columns []*schema.Column) string {

	if len(columns) == 0 {
		return ""
	}

	if len(columns) == 1 {
		return `fmt.Sprint(model.` + columns[0].Name + `)`
	}

	var sb strings.Builder
	for k := range columns {
		sb.WriteString(`fmt.Sprint(model.` + columns[k].Name + `)`)
		if k < len(columns)-1 {
			sb.WriteString(` + "_" + `)
		}
	}

	return sb.String()
}

func columnValuesToKey(columns []*schema.Column) string {

	if len(columns) == 0 {
		return ""
	}

	if len(columns) == 1 {
		return `fmt.Sprint(` + toArgName(columns[0].Name) + `)`
	}

	var sb strings.Builder
	for k := range columns {
		sb.WriteString(`fmt.Sprint(` + toArgName(columns[k].Name) + `)`)
		if k < len(columns)-1 {
			sb.WriteString(` + "_" + `)
		}
	}

	return sb.String()
}

type IndexColumn struct {
	Columns []*schema.Column
	Index   *lib.CacheConfigIndex
}

type AggregateColumn struct {
	Column    *schema.Column
	Aggregate *lib.CacheConfigAggregateProperty
}

type SearchColumn struct {
	ConditionColumns []*schema.Column
	SearchColumns    []*schema.Column
	Search           *lib.CacheConfigSearch
}

type CacheData struct {
	Indices    []*IndexColumn
	Properties []*AggregateColumn
	HasHashID  bool
	Location   string
	Search     []*SearchColumn
}

func ParseIndices(cacheConfig *lib.CacheConfig, table *schema.Table) *CacheData {

	var aggCount = 0
	if cacheConfig.Aggregate != nil {
		aggCount = len(cacheConfig.Aggregate.Properties)
	}

	var searchCount = 0
	if cacheConfig.Search != nil {
		searchCount = len(cacheConfig.Search)
	}

	var data = &CacheData{
		Indices:    []*IndexColumn{},
		Properties: make([]*AggregateColumn, aggCount),
		HasHashID:  cacheConfig.HasHashID,
		Search:     make([]*SearchColumn, searchCount),
	}

	if cacheConfig.Aggregate != nil {
		data.Location = cacheConfig.Aggregate.Location
	}

	for k := range cacheConfig.Indices {

		var index = cacheConfig.Indices[k]
		var fields = []string{strings.TrimSpace(index.Field)}

		if strings.Contains(index.Field, ",") {
			fields = strings.Split(index.Field, ",")
		}
		var columns = make([]*schema.Column, len(fields))

		for l := range fields {
			var field = strings.TrimSpace(fields[l])

			if _, ok := table.Columns[field]; !ok {
				panic(fmt.Sprintf("Index: Table %s has no column named `%s`", table.Name, field))
			}

			columns[l] = table.Columns[field]
			fields[l] = field
		}

		// fmt.Println("@@@@@@@@@@ Index: ", cacheConfig.Indices[k].Field)
		cacheConfig.Indices[k].Field = strings.Join(fields, ", ")

		data.Indices = append(data.Indices, &IndexColumn{
			Columns: columns,
			Index:   cacheConfig.Indices[k],
		})
	}

	if aggCount > 0 {
		for k := range cacheConfig.Aggregate.Properties {

			if _, ok := table.Columns[cacheConfig.Aggregate.Properties[k].On]; !ok {
				panic(fmt.Sprintf("Aggregate.On: Table `%s` has no column named `%s`", table.Name, cacheConfig.Aggregate.Properties[k].On))
			}

			var agg = &AggregateColumn{
				Aggregate: cacheConfig.Aggregate.Properties[k],
				Column:    table.Columns[cacheConfig.Aggregate.Properties[k].On],
			}

			data.Properties[k] = agg
		}
	}

	if searchCount > 0 {
		for k := range cacheConfig.Search {

			var searchColumns = []*schema.Column{}
			var conditionColumns = []*schema.Column{}
			for l := range cacheConfig.Search[k].Fields {

				if _, ok := table.Columns[cacheConfig.Search[k].Fields[l]]; !ok {
					panic(fmt.Sprintf("Invalid Search column `%s` on table `%s`", cacheConfig.Search[k].Fields[l], table.Name))
				}

				searchColumns = append(searchColumns, table.Columns[cacheConfig.Search[k].Fields[l]])
			}

			for l := range cacheConfig.Search[k].Conditions {

				if _, ok := table.Columns[cacheConfig.Search[k].Conditions[l]]; !ok {
					panic(fmt.Sprintf("Invalid Search condition column `%s` on table `%s`", cacheConfig.Search[k].Conditions[l], table.Name))
				}

				conditionColumns = append(conditionColumns, table.Columns[cacheConfig.Search[k].Conditions[l]])
			}

			data.Search[k] = &SearchColumn{
				ConditionColumns: conditionColumns,
				SearchColumns:    searchColumns,
				Search:           cacheConfig.Search[k],
			}
		}
	}

	return data
}

func toArgName(field string) string {
	if len(field) == 0 {
		return ""
	}
	return strings.ToLower(field[:1]) + field[1:]
}
