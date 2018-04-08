package main

import (
	"fmt"
	"strings"
)

func QueryCreateChangeSQL(localSchema *Database, remoteSchema *Database) (sql string, e error) {

	query := ""

	// What tables are in local that aren't in remote?
	for tableName, table := range localSchema.Tables {

		// Table does not exist
		if _, ok := remoteSchema.Tables[tableName]; !ok {

			// fmt.Printf("Local table %s is not in remote\n", table.Name)
			query, e = QueryCreateTable(&table)
			// fmt.Printf("Running Query: %s\n", query)
			sql += query + "\n"
		} else {
			remoteTable := remoteSchema.Tables[tableName]
			query, e = QueryCreateTableChangeSQL(&table, &remoteTable)
			if len(query) > 0 {
				sql += query + "\n"
			}
		}
	}

	for _, table := range remoteSchema.Tables {
		if _, ok := localSchema.Tables[table.Name]; !ok {
			query, e = QueryDropTable(&table)
			sql += query + "\n"
		}
	}

	return
}

func QueryCreateTableChangeSQL(localTable *Table, remoteTable *Table) (sql string, e error) {

	var query string

	for _, column := range localTable.Columns {

		// Column does not exist remotely
		if _, ok := remoteTable.Columns[column.Name]; !ok {
			query, e = QueryAlterTableCreateColumn(localTable, &column)
			if e != nil {
				return
			}

			if len(query) > 0 {
				sql += query + "\n"
			}

		} else {

			remoteColumn := remoteTable.Columns[column.Name]

			query, e = QueryAlterTableChangeColumn(localTable, &column, &remoteColumn)

			if e != nil {
				return
			}

			if len(query) > 0 {
				sql += query + "\n"
			}
		}
	}

	for _, column := range remoteTable.Columns {

		// Column does not exist locally
		if _, ok := localTable.Columns[column.Name]; !ok {
			query, e = QueryAlterTableDropColumn(localTable, &column)
			if e != nil {
				return
			}

			sql += query + "\n"
		}
	}

	return
}

func QueryAlterTableChangeColumn(table *Table, localColumn *Column, remoteColumn *Column) (sql string, e error) {

	query := ""

	// if localColumn.

	sql = query
	return

}

func QueryAlterTableCreateColumn(table *Table, column *Column) (sql string, e error) {

	query := ""

	query, e = QueryCreateColumn(column)
	sql = fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s;", table.Name, query)

	return

}

func QueryAlterTableDropColumn(table *Table, column *Column) (sql string, e error) {

	sql = fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`;", table.Name, column.Name)

	return

}

func QueryDropTable(table *Table) (sql string, e error) {
	sql = fmt.Sprintf("DROP TABLE `%s`;", table.Name)
	return
}

func QueryCreateTable(table *Table) (sql string, e error) {

	// colLen := len(table.Columns)
	idx := 1

	// Primary Key?
	primaryKey := ""

	cols := []string{}

	for _, column := range table.Columns {

		colQuery := ""
		colQuery, e = QueryCreateColumn(&column)
		col := colQuery

		idx++

		if column.ColumnKey == "PRI" {
			primaryKey = column.Name
		}

		cols = append(cols, col)
	}
	if len(primaryKey) > 0 {
		cols = append(cols, fmt.Sprintf("PRIMARY KEY(`%s`)", primaryKey))
	}

	sql = fmt.Sprintf("CREATE TABLE `%s` (\n\t%s\n);", table.Name, strings.Join(cols, ",\n\t"))

	return
}

func QueryCreateColumn(column *Column) (sql string, e error) {

	sql = fmt.Sprintf("`%s` %s", column.Name, column.Type)
	if !column.IsNullable {
		sql += " NOT"
	}
	sql += " NULL"

	if len(column.Default) > 0 || column.DataType == "varchar" {

		columnDefault := column.Default

		switch column.DataType {
		case "varchar":
			columnDefault = fmt.Sprintf("'%s'", column.Default)
		case "enum":
			columnDefault = fmt.Sprintf("'%s'", column.Default)
		}

		sql += fmt.Sprintf(" DEFAULT %s", columnDefault)
	}

	if len(column.Extra) > 0 {
		sql += " " + column.Extra
	}

	return

}
