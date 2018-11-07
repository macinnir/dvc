package lib

import (
	"fmt"
	"sort"
	"strings"
)

// Query contains functionality for generating sql queries
type Query struct{}

// CreateChangeSQL generates sql statements based off of comparing two database objects
// localSchema is authority, remoteSchema will be upgraded to match localSchema
func (q *Query) CreateChangeSQL(localSchema *Database, remoteSchema *Database) (sql string, e error) {

	query := ""

	// What tables are in local that aren't in remote?
	for tableName, table := range localSchema.Tables {

		// Table does not exist on remote schema
		if _, ok := remoteSchema.Tables[tableName]; !ok {

			// fmt.Printf("Local table %s is not in remote\n", table.Name)
			query, e = q.CreateTable(table)
			// fmt.Printf("Running Query: %s\n", query)
			sql += query + "\n"
		} else {
			remoteTable := remoteSchema.Tables[tableName]
			query, e = q.CreateTableChangeSQL(table, remoteTable)
			if len(query) > 0 {
				sql += query + "\n"
			}
		}
	}

	// What tables are in remote that aren't in local?
	for _, table := range remoteSchema.Tables {

		// Table does not exist on local schema
		if _, ok := localSchema.Tables[table.Name]; !ok {
			query, e = q.DropTable(table)
			sql += query + "\n"
		}
	}

	return
}

// CreateTableChangeSQL returns a set of statements that alter a table's structure if and only if there is a difference between
// the local and remote tables
// If no change is found, an empty string is returned.
func (q *Query) CreateTableChangeSQL(localTable *Table, remoteTable *Table) (sql string, e error) {

	var query string

	for _, column := range localTable.Columns {

		// Column does not exist remotely
		if _, ok := remoteTable.Columns[column.Name]; !ok {
			query, e = q.AlterTableCreateColumn(localTable, column)
			if e != nil {
				return
			}

			if len(query) > 0 {
				sql += query + "\n"
			}

		} else {

			remoteColumn := remoteTable.Columns[column.Name]

			query, e = q.ChangeColumn(localTable, column, remoteColumn)

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
			query, e = q.AlterTableDropColumn(localTable, column)
			if e != nil {
				return
			}

			sql += query + "\n"
		}
	}

	return
}

// CreateTable returns a create table sql statement
func (q *Query) CreateTable(table *Table) (sql string, e error) {

	// colLen := len(table.Columns)
	idx := 1

	// Primary Key?
	primaryKey := ""

	cols := []string{}

	// Unique Keys
	uniqueKeyColumns := []*Column{}

	// Regular Keys (allows for multiple entries)
	multiKeyColumns := []*Column{}

	sortedColumns := make(SortedColumns, 0, len(table.Columns))

	for _, column := range table.Columns {
		sortedColumns = append(sortedColumns, column)
	}

	sort.Sort(sortedColumns)

	for _, column := range sortedColumns {

		colQuery := ""
		colQuery, e = q.CreateColumn(column)
		col := colQuery

		idx++

		switch column.ColumnKey {
		case "PRI":
			primaryKey = column.Name
		case "UNI":
			uniqueKeyColumns = append(uniqueKeyColumns, column)
		case "MUL":
			multiKeyColumns = append(multiKeyColumns, column)
		}
		cols = append(cols, col)
	}

	if len(primaryKey) > 0 {
		cols = append(cols, fmt.Sprintf("PRIMARY KEY(`%s`)", primaryKey))
	}

	sql = fmt.Sprintf("CREATE TABLE `%s` (\n\t%s\n) ENGINE = %s;", table.Name, strings.Join(cols, ",\n\t"), table.Engine)

	if len(uniqueKeyColumns) > 0 {
		sql += "\n"
		for _, uniqueKeyColumn := range uniqueKeyColumns {
			t, _ := q.AddUniqueIndex(table, uniqueKeyColumn)
			sql += t + "\n"
		}
	}

	if len(multiKeyColumns) > 0 {
		sql += "\n"
		for _, multiKeyColumn := range multiKeyColumns {
			t, _ := q.AddIndex(table, multiKeyColumn)
			sql += t + "\n"
		}
	}

	return
}

// DropTable returns a drop table sql statement
func (q *Query) DropTable(table *Table) (sql string, e error) {
	sql = fmt.Sprintf("DROP TABLE `%s`;", table.Name)
	return
}

// CreateColumn returns a table column sql segment
func (q *Query) CreateColumn(column *Column) (sql string, e error) {

	sql = fmt.Sprintf("`%s` %s", column.Name, column.Type)
	if !column.IsNullable {
		sql += " NOT"
	}
	sql += " NULL"

	if column.DataType == "char" || column.DataType == "varchar" || column.DataType == "enum" {
		sql += fmt.Sprintf(" DEFAULT '%s'", column.Default)
	} else if len(column.Default) > 0 {
		sql += fmt.Sprintf(" DEFAULT %s", column.Default)
	}

	if len(column.Extra) > 0 {
		sql += " " + column.Extra
	}

	return

}

// AlterTableDropColumn returns an alter table sql statement that drops a column
func (q *Query) AlterTableDropColumn(table *Table, column *Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`;", table.Name, column.Name)
	return
}

// ChangeColumn returns an alter table sql statement that adds or removes an index from a column
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
func (q *Query) ChangeColumn(table *Table, localColumn *Column, remoteColumn *Column) (sql string, e error) {

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
				t, _ = q.DropIndex(table, localColumn)
				query += t + "\n"
			// 2
			case "UNI":
				t, _ = q.DropUniqueIndex(table, localColumn)
				query += t + "\n"
			}
		}

		// 3, 4: There is no indexing on the remote schema
		if remoteColumn.ColumnKey == "" {
			switch localColumn.ColumnKey {
			// 3
			case "MUL":
				t, _ = q.AddIndex(table, localColumn)
				query += t + "\n"
			// 4
			case "UNI":
				t, _ = q.AddUniqueIndex(table, localColumn)
				query += t + "\n"
			}
		}

		// 5
		if remoteColumn.ColumnKey == "MUL" && localColumn.ColumnKey == "UNI" {
			t, _ = q.DropIndex(table, localColumn)
			query += t + "\n"
			t, _ = q.AddUniqueIndex(table, localColumn)
			query += t + "\n"
		}

		// 6
		if remoteColumn.ColumnKey == "UNI" && localColumn.ColumnKey == "MUL" {
			t, _ = q.DropUniqueIndex(table, localColumn)
			query += t + "\n"
			t, _ = q.AddIndex(table, localColumn)
			query += t + "\n"
		}
	}

	sql = query
	return

}

// AlterTableCreateColumn returns an alter table sql statement that adds a column
func (q *Query) AlterTableCreateColumn(table *Table, column *Column) (sql string, e error) {

	query := ""

	query, e = q.CreateColumn(column)
	sql = fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s;", table.Name, query)

	return
}

// AddIndex returns an alter table sql statement that adds an index to a table
func (q *Query) AddIndex(table *Table, column *Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` ADD INDEX `i_%s` (`%s`);", table.Name, column.Name, column.Name)
	return
}

// AddUniqueIndex returns an alter table sql statement that adds a unique index to a table
func (q *Query) AddUniqueIndex(table *Table, column *Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` ADD UNIQUE INDEX `ui_%s` (`%s`);", table.Name, column.Name, column.Name)
	return
}

// DropIndex returns an alter table sql statement that drops an index
func (q *Query) DropIndex(table *Table, column *Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `i_%s`;", table.Name, column.Name)
	return
}

// DropUniqueIndex returns an alter table sql statement that drops a unique index
func (q *Query) DropUniqueIndex(table *Table, column *Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `ui_%s`;", table.Name, column.Name)
	return
}
