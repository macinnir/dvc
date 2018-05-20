package main

import (
	"fmt"
	"sort"
	"strings"
)

func QueryCreateChangeSQL(localSchema *Database, remoteSchema *Database) (sql string, e error) {

	query := ""

	// What tables are in local that aren't in remote?
	for tableName, table := range localSchema.Tables {

		// Table does not exist on remote schema
		if _, ok := remoteSchema.Tables[tableName]; !ok {

			// fmt.Printf("Local table %s is not in remote\n", table.Name)
			query, e = QueryCreateTable(table)
			// fmt.Printf("Running Query: %s\n", query)
			sql += query + "\n"
		} else {
			remoteTable := remoteSchema.Tables[tableName]
			query, e = QueryCreateTableChangeSQL(table, remoteTable)
			if len(query) > 0 {
				sql += query + "\n"
			}
		}
	}

	// What tables are in remote that aren't in local?
	for _, table := range remoteSchema.Tables {

		// Table does not exist on local schema
		if _, ok := localSchema.Tables[table.Name]; !ok {
			query, e = QueryDropTable(table)
			sql += query + "\n"
		}
	}

	return
}

// QueryCreateTableChangeSQL returns a set of statements that alter a table's structure if and only if there is a difference between
// the local and remote tables
// If no change is found, an empty string is returned.
func QueryCreateTableChangeSQL(localTable *Table, remoteTable *Table) (sql string, e error) {

	var query string

	for _, column := range localTable.Columns {

		// Column does not exist remotely
		if _, ok := remoteTable.Columns[column.Name]; !ok {
			query, e = QueryAlterTableCreateColumn(localTable, column)
			if e != nil {
				return
			}

			if len(query) > 0 {
				sql += query + "\n"
			}

		} else {

			remoteColumn := remoteTable.Columns[column.Name]

			query, e = QueryAlterTableChangeColumn(localTable, column, remoteColumn)

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
			query, e = QueryAlterTableDropColumn(localTable, column)
			if e != nil {
				return
			}

			sql += query + "\n"
		}
	}

	return
}

// QueryAlterTableChangeColumn returns an alter table sql statement that adds or removes an index from a column
// if and only if the one (e.g. local) has a column and the other (e.g. remote) does not
// Truth table
// 		Remote 	| 	Local 	| 	Result
// ---------------------------------------------------------
// 1. 	MUL		| 	none 	| 	Drop index
// 2. 	UNI		| 	none 	| 	Drop unique index
// 3. 	none 	| 	MUL 	|  	Create index
// 4. 	none 	| 	UNI 	| 	Create unique index
// 5. 	MUL		| 	UNI 	| 	Drop index; Create unique index
// 6. 	UNI 	| 	MUL 	| 	Drop unique index; Create index
// 7. 	none	| 	none	| 	Do nothing
// 8. 	MUL		| 	MUL		| 	Do nothing
// 9. 	UNI		|   UNI		| 	Do nothing
func QueryAlterTableChangeColumn(table *Table, localColumn *Column, remoteColumn *Column) (sql string, e error) {

	t := ""
	query := ""

	// 7,8,9
	if localColumn.ColumnKey == remoteColumn.ColumnKey {
		return
	}

	// <7
	if localColumn.ColumnKey != remoteColumn.ColumnKey {

		// 1,2: There is no indexing on the local schema
		if localColumn.ColumnKey == "" {
			switch remoteColumn.ColumnKey {
			// 1
			case "MUL":
				t, _ = QueryAlterTableDropIndex(table, localColumn)
				query += t + "\n"
			// 2
			case "UNI":
				t, _ = QueryAlterTableDropUniqueIndex(table, localColumn)
				query += t + "\n"
			}
		}

		// 3, 4: There is no indexing on the remote schema
		if remoteColumn.ColumnKey == "" {
			switch localColumn.ColumnKey {
			// 3
			case "MUL":
				t, _ = QueryAlterTableAddIndex(table, localColumn)
				query += t + "\n"
			// 4
			case "UNI":
				t, _ = QueryAlterTableAddUniqueIndex(table, localColumn)
				query += t + "\n"
			}
		}

		// 5
		if remoteColumn.ColumnKey == "MUL" && localColumn.ColumnKey == "UNI" {
			t, _ = QueryAlterTableDropIndex(table, localColumn)
			query += t + "\n"
			t, _ = QueryAlterTableAddUniqueIndex(table, localColumn)
			query += t + "\n"
		}

		// 6
		if remoteColumn.ColumnKey == "UNI" && localColumn.ColumnKey == "MUL" {
			t, _ = QueryAlterTableDropUniqueIndex(table, localColumn)
			query += t + "\n"
			t, _ = QueryAlterTableAddIndex(table, localColumn)
			query += t + "\n"
		}
	}

	sql = query
	return

}

// QueryAlterTableCreateColumn returns an alter table sql statement that adds a column
func QueryAlterTableCreateColumn(table *Table, column *Column) (sql string, e error) {

	query := ""

	query, e = QueryCreateColumn(column)
	sql = fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s;", table.Name, query)

	return
}

// QueryAlterTableDropColumn returns an alter table sql statement that drops a column
func QueryAlterTableDropColumn(table *Table, column *Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`;", table.Name, column.Name)
	return
}

// QueryAlterTableAddIndex returns an alter table sql statement that adds an index to a table
func QueryAlterTableAddIndex(table *Table, column *Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` ADD INDEX `i_%s` (`%s`)", table.Name, column.Name, column.Name)
	return
}

// QueryAlterTableAddUniqueIndex returns an alter table sql statement that adds a unique index to a table
func QueryAlterTableAddUniqueIndex(table *Table, column *Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` ADD UNIQUE INDEX `ui_%s` (`%s`)", table.Name, column.Name, column.Name)
	return
}

// QueryAlterTableDropIndex returns an alter table sql statement that drops an index
func QueryAlterTableDropIndex(table *Table, column *Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `i_%s`", table.Name, column.Name)
	return
}

// QueryAlterTableDropUniqueIndex returns an alter table sql statement that drops a unique index
func QueryAlterTableDropUniqueIndex(table *Table, column *Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `ui_%s`", table.Name, column.Name)
	return
}

// QueryDropTable returns a drop table sql statement
func QueryDropTable(table *Table) (sql string, e error) {
	sql = fmt.Sprintf("DROP TABLE `%s`;", table.Name)
	return
}

// QueryCreateTable returns a create table sql statement
func QueryCreateTable(table *Table) (sql string, e error) {

	// colLen := len(table.Columns)
	idx := 1

	// Primary Key?
	primaryKey := ""

	cols := []string{}

	// Unique Keys
	uniqueKeys := []string{}

	// Regular Keys (allows for multiple entries)
	multiKeys := []string{}

	sortedColumns := make(SortedColumns, 0, len(table.Columns))

	for _, column := range table.Columns {
		sortedColumns = append(sortedColumns, column)
	}

	sort.Sort(sortedColumns)

	for _, column := range sortedColumns {

		colQuery := ""
		colQuery, e = QueryCreateColumn(column)
		col := colQuery

		idx++

		switch column.ColumnKey {
		case "PRI":
			primaryKey = column.Name
		case "UNI":
			uniqueKeys = append(uniqueKeys, column.Name)
		case "MUL":
			multiKeys = append(multiKeys, column.Name)
		}
		cols = append(cols, col)
	}

	if len(primaryKey) > 0 {
		cols = append(cols, fmt.Sprintf("PRIMARY KEY(`%s`)", primaryKey))
	}

	if len(uniqueKeys) > 0 {
		for _, uniqueKey := range uniqueKeys {
			cols = append(cols, fmt.Sprintf("UNIQUE KEY (`%s`)", uniqueKey))
		}
	}

	if len(multiKeys) > 0 {
		for _, multiKey := range multiKeys {
			cols = append(cols, fmt.Sprintf("KEY (`%s`)", multiKey))
		}
	}

	sql = fmt.Sprintf("CREATE TABLE `%s` (\n\t%s\n);", table.Name, strings.Join(cols, ",\n\t"))

	return
}

// QueryCreateColumn returns a create table column sql statement
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
