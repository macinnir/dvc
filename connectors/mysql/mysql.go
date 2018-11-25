/**
 * MySQL
 * @implements IConnector
 */

package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/macinnir/dvc/lib"
)

// MySQL contains functionality for interacting with a server
type MySQL struct {
	Config *lib.Config
}

// Connect connects to a server and returns a new server object
func (ss *MySQL) Connect() (server *lib.Server, e error) {
	server = &lib.Server{Host: ss.Config.Connection.Host}
	var connectionString = ss.Config.Connection.Username + ":" + ss.Config.Connection.Password + "@tcp(" + ss.Config.Connection.Host + ")/" + ss.Config.Connection.DatabaseName + "?charset=utf8"
	server.Connection, e = sql.Open("mysql", connectionString)
	return
}

// FetchDatabases fetches a set of database names from the target server
// populating the Databases property with a map of Database objects
func (ss *MySQL) FetchDatabases(server *lib.Server) (databases map[string]*lib.Database, e error) {

	var rows *sql.Rows
	databases = map[string]*lib.Database{}

	if rows, e = server.Connection.Query("SHOW DATABASES"); e != nil {
		return
	}

	if rows != nil {
		defer rows.Close()
	}

	for rows.Next() {
		databaseName := ""
		rows.Scan(&databaseName)
		databases[databaseName] = &lib.Database{Name: databaseName, Host: server.Host}
	}

	return
}

// UseDatabase switches the connection context to the passed in database
func (ss *MySQL) UseDatabase(server *lib.Server, databaseName string) (e error) {

	if server.CurrentDatabase == databaseName {
		return
	}

	_, e = server.Connection.Exec(fmt.Sprintf("USE %s", databaseName))
	if e == nil {
		server.CurrentDatabase = databaseName
	}
	return
}

// FetchDatabaseTables fetches the complete set of tables from this database
func (ss *MySQL) FetchDatabaseTables(server *lib.Server, databaseName string) (tables map[string]*lib.Table, e error) {

	var rows *sql.Rows
	query := "select `TABLE_NAME`, `ENGINE`, `VERSION`, `ROW_FORMAT`, `TABLE_ROWS`, `DATA_LENGTH`, `TABLE_COLLATION`, `AUTO_INCREMENT` FROM information_schema.tables WHERE TABLE_SCHEMA = '" + databaseName + "'"
	// fmt.Printf("Query: %s\n", query)
	if rows, e = server.Connection.Query(query); e != nil {
		return
	}

	if rows != nil {
		defer rows.Close()
	}

	tables = map[string]*lib.Table{}

	for rows.Next() {

		table := &lib.Table{}

		rows.Scan(
			&table.Name,
			&table.Engine,
			&table.Version,
			&table.RowFormat,
			&table.Rows,
			&table.DataLength,
			&table.Collation,
			&table.AutoIncrement,
		)

		table.Columns, e = ss.FetchTableColumns(server, databaseName, table.Name)

		if e != nil {
			log.Fatalf("ERROR: %s", e.Error())
			return
		}

		tables[table.Name] = table
	}

	return
}

// FetchTableColumns lists all of the columns in a table
func (ss *MySQL) FetchTableColumns(server *lib.Server, databaseName string, tableName string) (columns map[string]*lib.Column, e error) {

	var rows *sql.Rows

	query := fmt.Sprintf(`
		SELECT 
			COLUMN_NAME, 
			ORDINAL_POSITION, 
			COALESCE(COLUMN_DEFAULT, '') as COLUMN_DEFAULT, 
			CASE IS_NULLABLE 
				WHEN 'YES' THEN 1
				ELSE 0
			END AS IS_NULLABLE,
			DATA_TYPE, 
			COALESCE(CHARACTER_MAXIMUM_LENGTH, 0) as MaxLength, 
			COALESCE(NUMERIC_PRECISION, 0) as NumericPrecision, 
			COALESCE(CHARACTER_SET_NAME, '') AS CharSet, 
			COLUMN_TYPE,
			COLUMN_KEY,
			EXTRA,
			COALESCE(NUMERIC_SCALE, 0) as NumericScale 
		FROM information_schema.COLUMNS 
		WHERE 
			TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'
	`, databaseName, tableName)

	if rows, e = server.Connection.Query(query); e != nil {
		return
	}

	if rows != nil {
		defer rows.Close()
	}

	columns = map[string]*lib.Column{}

	for rows.Next() {
		column := lib.Column{}
		if e = rows.Scan(
			&column.Name,
			&column.Position,
			&column.Default,
			&column.IsNullable,
			&column.DataType,
			&column.MaxLength,
			&column.Precision,
			&column.CharSet,
			&column.Type,
			&column.ColumnKey,
			&column.Extra,
			&column.NumericScale,
		); e != nil {
			return
		}
		column.IsUnsigned = strings.Contains(strings.ToLower(column.Type), " unsigned")
		columns[column.Name] = &column
	}

	// fmt.Printf("Fetching columns database: %s, table: %s - columns: %d\n", databaseName, tableName, len(columns))

	return
}

// CreateChangeSQL generates sql statements based off of comparing two database objects
// localSchema is authority, remoteSchema will be upgraded to match localSchema
func (ss *MySQL) CreateChangeSQL(localSchema *lib.Database, remoteSchema *lib.Database) (sql string, e error) {

	query := ""

	dropTableStatements := map[string]string{}
	createTableStatements := map[string]string{}

	// What tables are in local that aren't in remote?
	for tableName, table := range localSchema.Tables {

		// Table does not exist on remote schema
		if _, ok := remoteSchema.Tables[tableName]; !ok {
			query, e = createTable(table)
			createTableStatements[tableName] = query
		} else {
			remoteTable := remoteSchema.Tables[tableName]
			query, e = createTableChangeSQL(table, remoteTable)
			if len(query) > 0 {
				sql += query + "\n"
			}
		}
	}

	// What tables are in remote that aren't in local?
	for _, table := range remoteSchema.Tables {

		// Table does not exist on local schema
		if _, ok := localSchema.Tables[table.Name]; !ok {
			query, e = dropTable(table)
			dropTableStatements[table.Name] = query
		}
	}

	// Rename Table

	if len(dropTableStatements) > 0 && len(createTableStatements) > 0 {

		for dropTableName := range dropTableStatements {

			for createTableName := range createTableStatements {

				localTable := localSchema.Tables[createTableName]
				remoteTable := remoteSchema.Tables[dropTableName]

				// # of columns
				if len(localTable.Columns) == len(remoteTable.Columns) {

					same := true

					// Same column names
					for localColumnName := range localTable.Columns {
						if _, ok := remoteTable.Columns[localColumnName]; !ok {
							same = false
							break
						}
					}

					if same == true {
						delete(dropTableStatements, dropTableName)
						delete(createTableStatements, createTableName)
						sql += fmt.Sprintf("RENAME TABLE `%s` TO `%s`;\n", dropTableName, createTableName)
						break
					}
				}
			}
		}
	}

	if len(dropTableStatements) > 0 {
		for _, s := range dropTableStatements {
			sql += s + "\n"
		}
	}

	if len(createTableStatements) > 0 {
		for _, s := range createTableStatements {
			sql += s + "\n"
		}
	}

	if len(ss.Config.Enums) > 0 {
		sql += ss.compareEnums(remoteSchema, localSchema)
	}

	return
}

func (ss *MySQL) compareEnums(remoteSchema *lib.Database, localSchema *lib.Database) (sql string) {

	sql += ""

	// Compare remote to local
	for tableName, localTable := range localSchema.Enums {

		tableSQL := ""

		// remoteTable := remoteSchema.Enums[tableName]
		localTableSchema := localSchema.Tables[tableName]

		fieldMap := []string{}

		for fieldName := range localTableSchema.Columns {
			// fmt.Println("fieldName: ", fieldName)
			fieldMap = append(fieldMap, fieldName)
		}

		for _, localRow := range localTable {

			// If out of range with remote table, create a new entry
			fields := []string{}
			values := []string{}

			for _, fieldName := range fieldMap {

				// fmt.Println("fieldName:", fieldName)

				column := localTableSchema.Columns[fieldName]
				fields = append(fields, fmt.Sprintf("`%s`", fieldName))
				dataType := column.DataType
				value, valueExists := localRow[fieldName]

				if !valueExists {
					panic(fmt.Sprintf("Value for field `%s`.`%s` does not exist in enumerations. Please add this field to enumerations before continuing.", tableName, fieldName))
				}

				if isFloatingPointType(dataType) || isFixedPointType(dataType) {
					values = append(values, fmt.Sprintf("%f", value))
				} else if isString(dataType) {
					values = append(values, "'"+strings.Replace(fmt.Sprintf("%s", value), "'", "\\'", -1)+"'")
				} else if isInt(dataType) {
					values = append(values, fmt.Sprintf("%.0f", value))
				}
			}
			tableSQL += fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s);\n", tableName, strings.Join(fields, ","), strings.Join(values, ","))
		}
		if len(tableSQL) > 0 {
			sql += fmt.Sprintf("DELETE FROM `%s`;\n", tableName) + tableSQL
		}
	}

	return
}

// FetchEnums fetches enum data for all enums listed in config
func (ss *MySQL) FetchEnums(server *lib.Server) (enums map[string][]map[string]interface{}) {

	enums = make(map[string][]map[string]interface{})
	for _, enum := range ss.Config.Enums {
		// fmt.Printf("Building Enum: %s\n", enum)
		enums[enum] = fetchEnum(server, enum)
	}

	return
}

func fetchEnum(server *lib.Server, enum string) (objects []map[string]interface{}) {

	var e error
	var rows *sql.Rows

	if rows, e = server.Connection.Query(fmt.Sprintf("SELECT * FROM `%s`", enum)); e != nil {
		return
	}

	defer rows.Close()
	columnNames, _ := rows.Columns()
	columnTypes, _ := rows.ColumnTypes()

	for rows.Next() {

		values := make([]interface{}, len(columnNames))
		object := map[string]interface{}{}
		for i, column := range columnTypes {
			switch column.ScanType().Name() {
			case "RawBytes":
				if isFloatingPointType(column.DatabaseTypeName()) || isFixedPointType(column.DatabaseTypeName()) {
					v := 0.0
					object[column.Name()] = &v
				}

				if isString(column.DatabaseTypeName()) {
					v := ""
					object[column.Name()] = &v
				}

				if isInt(column.DatabaseTypeName()) {
					v := 0
					object[column.Name()] = &v
				}
			default:
				object[column.Name()] = reflect.New(column.ScanType()).Interface()
			}

			values[i] = object[column.Name()]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(values...); err != nil {
			panic(err)
		}

		objects = append(objects, object)
	}

	return
}

func getBytes(src interface{}) []byte {

	if a, ok := src.([]uint8); ok {
		return a
	}
	return nil
}

// createTableChangeSQL returns a set of statements that alter a table's structure if and only if there is a difference between
// the local and remote tables
// If no change is found, an empty string is returned.
func createTableChangeSQL(localTable *lib.Table, remoteTable *lib.Table) (sql string, e error) {

	var query string

	createColumnStatements := map[string]string{}
	dropColumnStatements := map[string]string{}

	for _, column := range localTable.Columns {

		// Column does not exist remotely
		if _, ok := remoteTable.Columns[column.Name]; !ok {
			query, e = alterTableCreateColumn(localTable, column)
			if e != nil {
				return
			}

			if len(query) > 0 {
				createColumnStatements[column.Name] += query
			}

		} else {

			remoteColumn := remoteTable.Columns[column.Name]

			query, e = changeColumn(localTable, column, remoteColumn)

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
			query, e = alterTableDropColumn(localTable, column)
			if e != nil {
				return
			}

			dropColumnStatements[column.Name] = query
		}
	}

	if len(dropColumnStatements) > 0 && len(createColumnStatements) > 0 {

		for dropColumnName := range dropColumnStatements {

			for createColumnName := range createColumnStatements {

				if remoteTable.Columns[dropColumnName].DataType == localTable.Columns[createColumnName].DataType {
					sql += alterTableRenameColumn(localTable, localTable.Columns[createColumnName], dropColumnName) + "\n"
					delete(dropColumnStatements, dropColumnName)
					delete(createColumnStatements, createColumnName)
					break
				}
			}

		}

	}

	if len(dropColumnStatements) > 0 {
		for _, s := range dropColumnStatements {
			sql += s + "\n"
		}
	}

	if len(createColumnStatements) > 0 {
		for _, s := range createColumnStatements {
			sql += s + "\n"
		}
	}

	return
}

// createTable returns a create table sql statement
func createTable(table *lib.Table) (sql string, e error) {

	// colLen := len(table.Columns)
	idx := 1

	// Primary Key?
	primaryKey := ""

	cols := []string{}

	// Unique Keys
	uniqueKeyColumns := []*lib.Column{}

	// Regular Keys (allows for multiple entries)
	multiKeyColumns := []*lib.Column{}

	sortedColumns := make(lib.SortedColumns, 0, len(table.Columns))

	for _, column := range table.Columns {
		sortedColumns = append(sortedColumns, column)
	}

	sort.Sort(sortedColumns)

	for _, column := range sortedColumns {

		colQuery := ""
		colQuery, e = createColumnSegment(column)
		col := colQuery

		idx++

		switch column.ColumnKey {
		case KeyPRI:
			primaryKey = column.Name
		case KeyUNI:
			uniqueKeyColumns = append(uniqueKeyColumns, column)
		case KeyMUL:
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
			t, _ := addUniqueIndex(table, uniqueKeyColumn)
			sql += t + "\n"
		}
	}

	if len(multiKeyColumns) > 0 {
		sql += "\n"
		for _, multiKeyColumn := range multiKeyColumns {
			t, _ := addIndex(table, multiKeyColumn)
			sql += t + "\n"
		}
	}

	return
}

// dropTable returns a drop table sql statement
func dropTable(table *lib.Table) (sql string, e error) {
	sql = fmt.Sprintf("DROP TABLE `%s`;", table.Name)
	return
}

// changeColumn returns an alter table sql statement that adds or removes an index from a column
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
func changeColumn(table *lib.Table, localColumn *lib.Column, remoteColumn *lib.Column) (sql string, e error) {

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
			case KeyMUL:
				t, _ = dropIndex(table, localColumn)
				query += t + "\n"
			// 2
			case KeyUNI:
				t, _ = dropUniqueIndex(table, localColumn)
				query += t + "\n"
			}
		}

		// 3, 4: There is no indexing on the remote schema
		if remoteColumn.ColumnKey == "" {
			switch localColumn.ColumnKey {
			// 3
			case KeyMUL:
				t, _ = addIndex(table, localColumn)
				query += t + "\n"
			// 4
			case KeyUNI:
				t, _ = addUniqueIndex(table, localColumn)
				query += t + "\n"
			}
		}

		// 5
		if remoteColumn.ColumnKey == KeyMUL && localColumn.ColumnKey == KeyUNI {
			t, _ = dropIndex(table, localColumn)
			query += t + "\n"
			t, _ = addUniqueIndex(table, localColumn)
			query += t + "\n"
		}

		// 6
		if remoteColumn.ColumnKey == KeyUNI && localColumn.ColumnKey == KeyMUL {
			t, _ = dropUniqueIndex(table, localColumn)
			query += t + "\n"
			t, _ = addIndex(table, localColumn)
			query += t + "\n"
		}
	}

	sql = query
	return

}

// createColumnSegment returns a table column sql segment
// Data Types
// INT SIGNED 	11 columns
func createColumnSegment(column *lib.Column) (sql string, e error) {

	if isInt(column.DataType) {

		sql = fmt.Sprintf("`%s` %s(%d)", column.Name, column.DataType, intColLength(column.DataType, column.IsUnsigned))

		if column.IsUnsigned == true {
			sql += fmt.Sprintf(" %s ", SignedUnsigned)
		} else {
			sql += fmt.Sprintf(" %s ", SignedSigned)
		}

	} else if isFixedPointType(column.DataType) {

		sql = fmt.Sprintf("`%s` %s(%d,%d)", column.Name, column.DataType, column.Precision, column.NumericScale)

		if column.IsUnsigned == true {
			sql += fmt.Sprintf(" %s ", SignedUnsigned)
		} else {
			sql += fmt.Sprintf(" %s ", SignedSigned)
		}
	} else if isFloatingPointType(column.DataType) {
		sql = fmt.Sprintf("`%s` %s(%d,%d)", column.Name, column.DataType, column.Precision, column.NumericScale)
		if column.IsUnsigned == true {
			sql += fmt.Sprintf(" %s ", SignedUnsigned)
		} else {
			sql += fmt.Sprintf(" %s ", SignedSigned)
		}
	} else if isString(column.DataType) {
		// Use the text from the `Type` field (the `COLUMN_TYPE` column) directly
		if strings.ToLower(column.DataType) == ColTypeEnum {
			sql = fmt.Sprintf("`%s` %s", column.Name, column.Type)
		} else if stringHasLength(column.DataType) {
			sql = fmt.Sprintf("`%s` %s(%d)", column.Name, column.DataType, column.MaxLength)
		} else {
			sql = fmt.Sprintf("`%s` %s", column.Name, column.DataType)
		}

	} else {
		sql = fmt.Sprintf("`%s` %s", column.Name, column.DataType)
	}

	if !column.IsNullable {
		sql += " NOT"
	}
	sql += " NULL"

	// Add single quotes to string default
	if hasDefaultString(column.DataType) {
		sql += fmt.Sprintf(" DEFAULT '%s'", column.Default)
	} else if len(column.Default) > 0 {
		sql += fmt.Sprintf(" DEFAULT %s", column.Default)
	}

	if len(column.Extra) > 0 {
		sql += " " + column.Extra
	}

	return

}

func alterTableRenameColumn(table *lib.Table, newColumn *lib.Column, oldColumnName string) (sql string) {
	query, _ := createColumnSegment(newColumn)
	sql = fmt.Sprintf("ALTER TABLE `%s` CHANGE `%s` %s;", table.Name, oldColumnName, query)
	return
}

// alterTableCreateColumn returns an alter table sql statement that adds a column
func alterTableCreateColumn(table *lib.Table, column *lib.Column) (sql string, e error) {
	query := ""
	query, e = createColumnSegment(column)
	sql = fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s;", table.Name, query)
	return
}

// alterTableDropColumn returns an alter table sql statement that drops a column
func alterTableDropColumn(table *lib.Table, column *lib.Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`;", table.Name, column.Name)
	return
}

// addIndex returns an alter table sql statement that adds an index to a table
func addIndex(table *lib.Table, column *lib.Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` ADD INDEX `i_%s` (`%s`);", table.Name, column.Name, column.Name)
	return
}

// addUniqueIndex returns an alter table sql statement that adds a unique index to a table
func addUniqueIndex(table *lib.Table, column *lib.Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` ADD UNIQUE INDEX `ui_%s` (`%s`);", table.Name, column.Name, column.Name)
	return
}

// dropIndex returns an alter table sql statement that drops an index
func dropIndex(table *lib.Table, column *lib.Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `i_%s_%s`;", table.Name, table.Name, column.Name)
	return
}

// dropUniqueIndex returns an alter table sql statement that drops a unique index
func dropUniqueIndex(table *lib.Table, column *lib.Column) (sql string, e error) {
	sql = fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `ui_%s_%s`;", table.Name, table.Name, column.Name)
	return
}

// stringHasLength
func stringHasLength(dataType string) bool {
	switch strings.ToLower(dataType) {
	case ColTypeVarchar, ColTypeChar:
		return true
	}
	return false
}

// hasDefaultString
func hasDefaultString(dataType string) bool {
	switch strings.ToLower(dataType) {
	case ColTypeVarchar, ColTypeChar, ColTypeEnum:
		return true
	}
	return false
}

// isString
// String Types: https://dev.mysql.com/doc/refman/8.0/en/string-types.html
func isString(dataType string) bool {
	switch strings.ToLower(dataType) {
	case ColTypeVarchar, ColTypeEnum, ColTypeChar, ColTypeTinyText, ColTypeMediumText, ColTypeText, ColTypeLongText:
		return true
	}
	return false
}

// isInt
// Integer DataTypes: https://dev.mysql.com/doc/refman/8.0/en/integer-types.html
func isInt(dataType string) bool {
	switch strings.ToLower(dataType) {
	case ColTypeTinyint:
		return true
	case ColTypeSmallint:
		return true
	case ColTypeMediumint:
		return true
	case ColTypeInt:
		return true
	case ColTypeBigint:
		return true
	}
	return false
}

// Fixed Point Types
// https://dev.mysql.com/doc/refman/8.0/en/fixed-point-types.html
func isFixedPointType(dataType string) bool {
	switch strings.ToLower(dataType) {
	case ColTypeDecimal:
		return true
	case ColTypeNumeric:
		return true
	}
	return false
}

// Floating Point Types
// https://dev.mysql.com/doc/refman/8.0/en/floating-point-types.html
func isFloatingPointType(dataType string) bool {
	switch strings.ToLower(dataType) {
	case ColTypeFloat:
		return true
	case ColTypeDouble:
		return true
	}

	return false
}

func intColLength(dataType string, isUnsigned bool) int {
	switch dataType {
	case ColTypeTinyint:
		if isUnsigned {
			return 3
		}
		return 4
	case ColTypeSmallint:
		if isUnsigned {
			return 5
		}
		return 6
	case ColTypeMediumint:
		if isUnsigned {
			return 8
		}
		return 9
	case ColTypeInt:
		if isUnsigned {
			return 10
		}
		return 11
	case ColTypeBigint:
		return 20
	}

	return 0
}
