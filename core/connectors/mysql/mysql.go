/**
 * MySQL
 * @implements IConnector
 */

// TODO Add/remove primary key
// ALTER TABLE `Profile` ADD PRIMARY KEY (`ProfileID`);
// ALTER TABLE `Profile` DROP PRIMARY KEY
// ALTER TABLE `Persons` ADD CONSTRAINT PK_Person PRIMARY KEY (ID,LastName);
// ALTER TABLE Persons DROP PRIMARY KEY;
// ALTER TABLE `Persons` DROP CONSTRAINT `PK_Person`

package mysql

import (
	"database/sql"
	"fmt"

	// "gopkg.in/guregu/null.v3"
	"log"
	"sort"
	"strings"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

// MySQL contains functionality for interacting with a server
type MySQL struct {
	config *lib.ConfigDatabase
}

func NewMySQL(config *lib.ConfigDatabase) *MySQL {
	return &MySQL{
		config: config,
	}
}

// Connect connects to a server and returns a new server object
func (ss *MySQL) Connect() (server *schema.Server, e error) {
	// fmt.Println("Connecting to ", ss.config.Name)
	server = &schema.Server{Host: ss.config.Host}
	var connectionString = ss.config.User + ":" + ss.config.Pass + "@tcp(" + ss.config.Host + ")/" + ss.config.Name + "?charset=utf8"
	server.Connection, e = sql.Open("mysql", connectionString)
	return
}

// FetchDatabases fetches a set of database names from the target server
// populating the Databases property with a map of Database objects
func (ss *MySQL) FetchDatabases(server *schema.Server) (map[string]*schema.Database, error) {

	var e error
	var databases = map[string]*schema.Database{}
	var rows *sql.Rows

	if rows, e = server.Connection.Query("SELECT SCHEMA_NAME, DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA"); e != nil {
		return nil, e
	}

	if rows != nil {
		defer rows.Close()
	}

	for rows.Next() {

		databaseName := ""
		characterSet := ""
		collation := ""
		rows.Scan(
			&databaseName,
			&characterSet,
			&collation,
		)
		databases[databaseName] = &schema.Database{
			Name:                databaseName,
			DefaultCharacterSet: characterSet,
			DefaultCollation:    collation,
		}

		fmt.Println("CharacterSet", characterSet)
		fmt.Println("Collation", collation)
	}

	return databases, nil
}

// FetchDatabases fetches a set of database names from the target server
// populating the Databases property with a map of Database objects
func (ss *MySQL) fetchDatabaseMeta(server *schema.Server, schema *schema.Database) (e error) {

	var rows *sql.Rows

	if rows, e = server.Connection.Query(fmt.Sprintf("SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = '%s'", schema.Name)); e != nil {
		return
	}

	if e == sql.ErrNoRows {
		fmt.Println("SCHEMATA not found for database ", schema.Name)
		return nil
	}

	if rows != nil {
		defer rows.Close()
	}

	for rows.Next() {
		rows.Scan(
			&schema.DefaultCharacterSet,
			&schema.DefaultCollation,
		)
	}

	return
}

// UseDatabase switches the connection context to the passed in database
func (ss *MySQL) UseDatabase(server *schema.Server, databaseName string) (e error) {

	if server.CurrentDatabase == databaseName {
		return
	}

	_, e = server.Connection.Exec(fmt.Sprintf("USE %s", databaseName))
	if e == nil {
		server.CurrentDatabase = databaseName
	}
	return
}

func (ss *MySQL) FetchDatabase(server *schema.Server, databaseName string) (*schema.Database, error) {

	var e error
	database := &schema.Database{
		Name: databaseName,
	}

	if e = ss.fetchDatabaseMeta(server, database); e != nil {
		return nil, e
	}

	if database.Tables, e = ss.fetchDatabaseTables(server, databaseName); e != nil {
		return nil, e
	}

	return database, nil

}

// fetchDatabaseTables fetches the complete set of tables from this database
func (ss *MySQL) fetchDatabaseTables(server *schema.Server, databaseName string) (tables map[string]*schema.Table, e error) {

	var rows *sql.Rows
	query := "select t.`TABLE_NAME`, t.`ENGINE`, t.`VERSION`, t.`ROW_FORMAT`, t.`TABLE_ROWS`, t.`DATA_LENGTH`, t.`TABLE_COLLATION`, COALESCE(t.`AUTO_INCREMENT`, 0) AS `AUTO_INCREMENT`, ccsa.CHARACTER_SET_NAME FROM information_schema.tables t JOIN information_schema.`COLLATION_CHARACTER_SET_APPLICABILITY` ccsa ON ccsa.`COLLATION_NAME` = t.`TABLE_COLLATION` WHERE TABLE_SCHEMA = '" + databaseName + "'"

	if rows, e = server.Connection.Query(query); e != nil {
		return
	}

	if rows != nil {
		defer rows.Close()
	}

	tables = map[string]*schema.Table{}

	for rows.Next() {

		table := &schema.Table{}

		rows.Scan(
			&table.Name,
			&table.Engine,
			&table.Version,
			&table.RowFormat,
			&table.Rows,
			&table.DataLength,
			&table.Collation,
			&table.AutoIncrement,
			&table.CharacterSet,
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
func (ss *MySQL) FetchTableColumns(server *schema.Server, databaseName string, tableName string) (columns map[string]*schema.Column, e error) {

	var rows *sql.Rows

	query := fmt.Sprintf(`
		SELECT
			COLUMN_NAME,
			-- ORDINAL_POSITION,
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
			COALESCE(NUMERIC_SCALE, 0) as NumericScale,
			COALESCE(COLLATION_NAME, '') as Collation
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

	columns = map[string]*schema.Column{}

	for rows.Next() {
		column := schema.Column{}
		if e = rows.Scan(
			&column.Name,
			// &column.Position,
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
			&column.Collation,
		); e != nil {
			return
		}

		column.IsUnsigned = strings.Contains(strings.ToLower(column.Type), " unsigned")
		column.FmtType = schema.DataTypeToFormatString(&column)
		column.GoType = schema.DataTypeToGoTypeString(&column)

		if column.Default == "''" {
			column.Default = ""
		}
		column.IsString = isString(column.DataType)
		columns[column.Name] = &column
	}

	// fmt("Fetching columns database: %s, table: %s - columns: %d\n", databaseName, tableName, len(columns))

	return
}

// CreateChangeSQL generates sql statements based off of comparing two database objects
// localSchema is authority, remoteSchema will be upgraded to match localSchema
func (ss *MySQL) CreateChangeSQL(localSchema *schema.Schema, remoteSchema *schema.Schema, databaseName string) *schema.SchemaComparison {

	comparison := &schema.SchemaComparison{
		Database: "",
		Changes:  []*schema.SchemaChange{},
	}

	createTableStatements := map[string][]*schema.SchemaChange{}
	dropTableStatements := map[string][]*schema.SchemaChange{}

	// Character Encoding
	if len(localSchema.DefaultCharacterSet) > 0 && (localSchema.DefaultCharacterSet != remoteSchema.DefaultCharacterSet ||
		localSchema.DefaultCollation != remoteSchema.DefaultCollation) {

		comparison.Changes = append(comparison.Changes, alterDatabaseCharacterSet(databaseName, localSchema.DefaultCharacterSet, localSchema.DefaultCollation))
		comparison.Alterations++
	}

	// What tables are in local that aren't in remote?
	for tableName, table := range localSchema.Tables {

		// Table does not exist on remote schema
		if _, ok := remoteSchema.Tables[tableName]; !ok {
			createTableStatements[tableName] = createTable(table)
		} else {
			remoteTable := remoteSchema.Tables[tableName]

			createTableChangeSQL(comparison, table, remoteTable)
		}
	}

	// What tables are in remote that aren't in local?
	for _, table := range remoteSchema.Tables {

		// Table does not exist on local schema
		if _, ok := localSchema.Tables[table.Name]; !ok {
			dropTableStatements[table.Name] = dropTable(table)
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

					if same {
						comparison.Changes = append(comparison.Changes, &schema.SchemaChange{
							Type: schema.RenameTable,
							SQL:  fmt.Sprintf("RENAME TABLE `%s` TO `%s`;\n", dropTableName, createTableName),
						})

						delete(dropTableStatements, dropTableName)
						delete(createTableStatements, createTableName)
						break
					}
				}
			}
		}
	}

	if len(dropTableStatements) > 0 {
		for k := range dropTableStatements {
			for l := range dropTableStatements[k] {
				comparison.Deletions++
				comparison.Changes = append(comparison.Changes, dropTableStatements[k][l])
			}
		}
	}

	if len(createTableStatements) > 0 {
		for k := range createTableStatements {
			for l := range createTableStatements[k] {
				comparison.Additions++
				comparison.Changes = append(comparison.Changes, createTableStatements[k][l])
			}
		}
	}

	return comparison
}

// CompareEnums returns a set of sql statements based on the difference between local (authority) and remote
func (ss *MySQL) CompareEnums(
	remoteSchema *schema.Schema,
	localSchema *schema.Schema,
	tableName string,
) (sql string) {

	sql += ""

	localTable := localSchema.Enums[tableName]

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

	return
}

// FetchEnums fetches enum data for all enums listed in config
// func (ss *MySQL) FetchEnums(server *schema.Server) (enums map[string][]map[string]interface{}) {

// 	enums = make(map[string][]map[string]interface{})
// 	for _, enum := range ss.Config.Enums {
// 		// fmt.Printf("Building Enum: %s\n", enum)
// 		enums[enum] = ss.FetchEnum(server, enum)
// 	}

// 	return
// }

func (ss *MySQL) FetchEnum(server *schema.Server, tableName string) (objects []map[string]interface{}) {

	var e error
	var rows *sql.Rows

	if rows, e = server.Connection.Query(fmt.Sprintf("SELECT * FROM `%s`", tableName)); e != nil {
		return
	}

	defer rows.Close()
	columnNames, _ := rows.Columns()
	// columnTypes, _ := rows.ColumnTypes()

	values := make([]interface{}, len(columnNames))
	valuePtrs := make([]interface{}, len(columnNames))

	for rows.Next() {

		for i := range columnNames {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		object := map[string]interface{}{}
		for i, col := range columnNames {
			val := values[i]
			b, ok := val.([]byte)
			var v interface{}
			if ok {
				v = string(b)
			} else {
				v = val
			}

			// fmt.Println(col, v)
			object[col] = v
		}
		objects = append(objects, object)

		// object := map[string]interface{}{}
		// // types := map[string]string{}

		// for i, column := range columnTypes {

		// 	fmt.Printf("%s.%s > %s > ScanTypeName: %s\n", tableName, column.Name(), column.DatabaseTypeName(), column.ScanType().Name())

		// 	switch column.ScanType().Name() {
		// 	case "RawBytes":

		// 		// fmt.Println(column.Name(), "RawBytes", column.DatabaseTypeName())

		// 		if isFloatingPointType(column.DatabaseTypeName()) || isFixedPointType(column.DatabaseTypeName()) {
		// 			v := 0.0
		// 			object[column.Name()] = &v
		// 		}

		// 		if isString(column.DatabaseTypeName()) {
		// 			nullable, _ := column.Nullable()
		// 			if nullable {
		// 				fmt.Printf("Nullable String: %s.%s\n", tableName, column.Name())
		// 				str := sql.NullString{}
		// 				object[column.Name()] = &str
		// 			} else {
		// 				v := ""
		// 				object[column.Name()] = &v
		// 			}
		// 		}

		// 		if isInt(column.DatabaseTypeName()) {
		// 			v := 0
		// 			object[column.Name()] = &v
		// 		}

		// 	case "NullTime":
		// 		v := ""
		// 		// fmt.Println(column.Name(), "NullTime", column.DatabaseTypeName())
		// 		object[column.Name()] = &v
		// 		// fmt.Println("!!", column.ScanType().Name(), column.DatabaseTypeName())
		// 	default:
		// 		// fmt.Println(column.Name(), "Default", column.DatabaseTypeName())
		// 		object[column.Name()] = reflect.New(column.ScanType()).Interface()
		// 	}

		// 	values[i] = object[column.Name()]
		// }

		// // Scan the result into the column pointers...
		// if err := rows.Scan(values...); err != nil {
		// 	panic(err)
		// }

		// for k :=

		// objects = append(objects, object)
	}

	return
}

// func mapType() {
// 	row := map[string]interface{}{}
// 	for j, l := range enums[k] {
// 		// fmt.Printf("%s %v\n", j, l)

// 		switch table.Columns[j].DataType {
// 		case "int":
// 			row[j] = *(l.(*uint32))
// 			// fmt.Println(j, *(l.(*uint32)))
// 		case "bigint":
// 			row[j] = *(l.(*uint64))
// 			// fmt.Println(j, *(l.(*uint64)))
// 		case "tinyint":
// 			if table.Columns[j].IsUnsigned {
// 				row[j] = *(l.(*uint8))
// 				// fmt.Println(j, *(l.(*uint8)))
// 			} else {
// 				row[j] = *(l.(*int8))
// 				// fmt.Println(j, *(l.(*int8)))
// 			}
// 		// case "int":
// 		// 	fmt.Println(j, l.(int8))
// 		// case "float64":
// 		// 	fmt.Println(j, l.(float64))
// 		case "varchar", "char", "text":
// 			row[j] = *(l.(*string))
// 			// fmt.Println(j, *(l.(*string)))
// 		default:
// 			fmt.Println(j, "???")
// 		}

// 		// val := reflect.ValueOf(l)
// 		// fmt.Println(j, val)
// 	}
// }

// createTableChangeSQL returns a set of statements that alter a table's structure if and only if there is a difference between
// the local and remote tables
// If no change is found, an empty string is returned.
func createTableChangeSQL(comparison *schema.SchemaComparison, localTable *schema.Table, remoteTable *schema.Table) {

	if localTable.CharacterSet != remoteTable.CharacterSet ||
		localTable.Collation != remoteTable.Collation {
		fmt.Printf("%s CHANGE CHARSET FROM: %s/%s => %s/%s\n", remoteTable.Name, remoteTable.CharacterSet, remoteTable.Collation, localTable.CharacterSet, localTable.Collation)
		comparison.Changes = append(comparison.Changes, alterTableCharacterSet(localTable.Name, localTable.CharacterSet, localTable.Collation))
		comparison.Alterations++
	}

	createColumnStatements := map[string]*schema.SchemaChange{}
	dropColumnStatements := map[string]*schema.SchemaChange{}
	addIndexStatements := map[string]*schema.SchemaChange{}
	dropIndexStatements := map[string]*schema.SchemaChange{}

	for _, column := range localTable.Columns {

		// Column does not exist remotely
		if _, ok := remoteTable.Columns[column.Name]; !ok {

			createColumnStatements[column.Name] = alterTableCreateColumn(localTable, column)

			// If a unique index is added to this column, include it

			if column.ColumnKey == KeyUNI {
				addIndexStatements[column.Name] = addUniqueIndex(localTable, column)
			}
			if column.ColumnKey == KeyMUL {
				addIndexStatements[column.Name] = addIndex(localTable, column)
			}

		} else {

			remoteColumn := remoteTable.Columns[column.Name]
			changeColumn(comparison, localTable, column, remoteColumn)

		}
	}

	for _, column := range remoteTable.Columns {

		// Column does not exist locally
		if _, ok := localTable.Columns[column.Name]; !ok {
			dropColumn := alterTableDropColumn(localTable, column)
			dropColumnStatements[column.Name] = dropColumn

			// If this column had a unique index, make sure it's dropped
			if column.ColumnKey == KeyUNI {
				dropIndexStatements[column.Name] = dropUniqueIndex(remoteTable, column)
			}

			if column.ColumnKey == KeyMUL {
				dropIndexStatements[column.Name] = dropIndex(remoteTable, column)
			}
		}
	}

	if len(dropColumnStatements) > 0 && len(createColumnStatements) > 0 {

		for dropColumnName := range dropColumnStatements {

			for createColumnName := range createColumnStatements {

				// The dataTypes are the same, possibly a renamed column. Simply rename the column
				if remoteTable.Columns[dropColumnName].DataType == localTable.Columns[createColumnName].DataType {
					comparison.Changes = append(comparison.Changes, alterTableChangeColumn(localTable, localTable.Columns[createColumnName], dropColumnName))
					comparison.Alterations++
					delete(dropColumnStatements, dropColumnName)
					delete(createColumnStatements, createColumnName)
					break
				}
			}

		}

	}

	if len(dropColumnStatements) > 0 {
		for k := range dropColumnStatements {

			// Include the index change if it exists
			if _, ok := dropIndexStatements[k]; ok {
				comparison.Changes = append(comparison.Changes, dropIndexStatements[k])
				comparison.Deletions++
			}

			comparison.Changes = append(comparison.Changes, dropColumnStatements[k])
			comparison.Deletions++
		}
	}

	if len(createColumnStatements) > 0 {

		for k := range createColumnStatements {

			comparison.Changes = append(comparison.Changes, createColumnStatements[k])
			comparison.Additions++

			// Include the index change if it exists
			if _, ok := addIndexStatements[k]; ok {
				comparison.Changes = append(comparison.Changes, addIndexStatements[k])
				comparison.Additions++
			}
		}
	}

	return
}

// createTable returns a create table sql statement
func createTable(table *schema.Table) []*schema.SchemaChange {

	changes := []*schema.SchemaChange{}

	// colLen := len(table.Columns)
	idx := 1

	// Primary Key?
	primaryKey := ""

	cols := []string{}

	// Unique Keys
	uniqueKeyColumns := []*schema.Column{}

	// Regular Keys (allows for multiple entries)
	multiKeyColumns := []*schema.Column{}

	sortedColumns := make(schema.SortedColumns, 0, len(table.Columns))

	for _, column := range table.Columns {
		sortedColumns = append(sortedColumns, column)
	}

	sort.Sort(sortedColumns)

	for _, column := range sortedColumns {

		colQuery := ""
		colQuery = createColumnSegment(column)
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

	sql := fmt.Sprintf("CREATE TABLE `%s` (\n\t%s\n)", table.Name, strings.Join(cols, ",\n\t"))

	if len(table.Engine) > 0 {
		sql += fmt.Sprintf(" ENGINE = %s", table.Engine)
	}

	if len(table.CharacterSet) > 0 {
		sql += fmt.Sprintf(" CHARACTER SET %s", table.CharacterSet)
	}

	if len(table.Collation) > 0 {
		sql += fmt.Sprintf(" COLLATE %s", table.Collation)
	}

	sql += ";"

	// Create table
	changes = append(changes, &schema.SchemaChange{
		Type: schema.CreateTable,
		SQL:  sql,
	})

	if len(uniqueKeyColumns) > 0 {
		for k := range uniqueKeyColumns {
			changes = append(changes, addUniqueIndex(table, uniqueKeyColumns[k]))
		}
	}

	if len(multiKeyColumns) > 0 {
		for k := range multiKeyColumns {
			changes = append(changes, addIndex(table, multiKeyColumns[k]))
		}
	}

	return changes
}

// dropTable returns a drop table sql statement
func dropTable(table *schema.Table) []*schema.SchemaChange {
	return []*schema.SchemaChange{
		{
			Type:          schema.DropTable,
			SQL:           fmt.Sprintf("DROP TABLE `%s`;", table.Name),
			IsDestructive: true,
		},
	}
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
func changeColumn(comparison *schema.SchemaComparison, table *schema.Table, localColumn *schema.Column, remoteColumn *schema.Column) {

	// 7,8,9
	// if localColumn.ColumnKey == remoteColumn.ColumnKey {
	// 	return
	// }

	// <7
	// The key for this column has been added/removed.
	if localColumn.ColumnKey != remoteColumn.ColumnKey {

		// 1,2: There is no indexing on the local schema
		if localColumn.ColumnKey == "" {
			switch remoteColumn.ColumnKey {
			// 1
			case KeyMUL:
				comparison.Changes = append(comparison.Changes, dropIndex(table, localColumn))
				comparison.Deletions++
			// 2
			case KeyUNI:
				comparison.Changes = append(comparison.Changes, dropUniqueIndex(table, localColumn))
				comparison.Deletions++
			case KeyPRI:
				comparison.Changes = append(comparison.Changes, dropPrimaryKey(table, localColumn))
				comparison.Deletions++
			}
		}

		// 3, 4: There is no indexing on the remote schema
		if remoteColumn.ColumnKey == "" {
			switch localColumn.ColumnKey {
			// 3
			case KeyMUL:
				comparison.Changes = append(comparison.Changes, addIndex(table, localColumn))
				comparison.Additions++
			// 4
			case KeyUNI:
				comparison.Changes = append(comparison.Changes, addUniqueIndex(table, localColumn))
				comparison.Additions++
			case KeyPRI:
				comparison.Changes = append(comparison.Changes, addPrimaryKey(table, localColumn))
				comparison.Additions++
			}
		}

		// 5 Drop multi key, add unique key
		if remoteColumn.ColumnKey == KeyMUL && localColumn.ColumnKey == KeyUNI {
			comparison.Changes = append(comparison.Changes, dropIndex(table, localColumn))
			comparison.Deletions++

			comparison.Changes = append(comparison.Changes, addUniqueIndex(table, localColumn))
			comparison.Additions++
		}

		// 6 Drop unique key, add remote key
		if remoteColumn.ColumnKey == KeyUNI && localColumn.ColumnKey == KeyMUL {
			comparison.Changes = append(comparison.Changes, dropUniqueIndex(table, localColumn))
			comparison.Deletions++

			comparison.Changes = append(comparison.Changes, addIndex(table, localColumn))
			comparison.Additions++
		}
	}

	if localColumn.DataType != remoteColumn.DataType ||
		localColumn.CharSet != remoteColumn.CharSet ||
		localColumn.Collation != remoteColumn.Collation ||
		localColumn.MaxLength != remoteColumn.MaxLength {
		comparison.Changes = append(comparison.Changes, alterTableChangeColumn(table, localColumn, localColumn.Name))
		comparison.Alterations++
	}
}

func alterTableChangeColumn(table *schema.Table, newColumn *schema.Column, oldColumnName string) *schema.SchemaChange {
	query := createColumnSegment(newColumn)

	sql := fmt.Sprintf("ALTER TABLE `%s` CHANGE `%s` %s", table.Name, oldColumnName, query)

	return &schema.SchemaChange{
		Type: schema.ChangeColumn,
		SQL:  sql + ";",
	}
}

// alterTableCreateColumn returns an alter table sql statement that adds a column
func alterTableCreateColumn(table *schema.Table, column *schema.Column) *schema.SchemaChange {
	query := createColumnSegment(column)

	sql := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s", table.Name, query)

	return &schema.SchemaChange{
		Type: schema.AddColumn,
		SQL:  sql + ";",
	}
}

// alterTableDropColumn returns an alter table sql statement that drops a column
func alterTableDropColumn(table *schema.Table, column *schema.Column) *schema.SchemaChange {
	return &schema.SchemaChange{
		Type:          schema.DropColumn,
		SQL:           fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`;", table.Name, column.Name),
		IsDestructive: true,
	}
}

func addIndex(table *schema.Table, column *schema.Column) *schema.SchemaChange {
	return &schema.SchemaChange{
		Type: schema.AddIndex,
		SQL:  fmt.Sprintf("ALTER TABLE `%s` ADD INDEX `i_%s_%s` (`%s`);", table.Name, table.Name, column.Name, column.Name),
	}
}

func addUniqueIndex(table *schema.Table, column *schema.Column) *schema.SchemaChange {
	return &schema.SchemaChange{
		Type: schema.AddIndex,
		SQL:  fmt.Sprintf("ALTER TABLE `%s` ADD UNIQUE INDEX `ui_%s_%s` (`%s`);", table.Name, table.Name, column.Name, column.Name),
	}
}

// alternative: https://www.techonthenet.com/mysql/primary_keys.php
func addPrimaryKey(table *schema.Table, column *schema.Column) *schema.SchemaChange {
	return &schema.SchemaChange{
		Type: schema.AddIndex,
		SQL:  fmt.Sprintf("ALTER TABLE `%s` ADD PRIMARY KEY (`%s`);", table.Name, column.Name),
	}
}

func alterDatabaseCharacterSet(databaseName, characterSet, collation string) *schema.SchemaChange {
	return &schema.SchemaChange{
		Type: schema.ChangeCharacterSet,
		SQL:  fmt.Sprintf("ALTER DATABASE `%s` CHARACTER SET %s COLLATE %s;", databaseName, characterSet, collation),
	}
}

func alterTableCharacterSet(tableName, characterSet, collation string) *schema.SchemaChange {
	return &schema.SchemaChange{
		Type: schema.ChangeCharacterSet,
		SQL:  fmt.Sprintf("ALTER TABLE `%s` CONVERT TO CHARACTER SET %s COLLATE %s;", tableName, characterSet, collation),
	}
}

// dropIndex returns an alter table sql statement that drops an index
func dropIndex(table *schema.Table, column *schema.Column) *schema.SchemaChange {
	return &schema.SchemaChange{
		Type:          schema.DropIndex,
		SQL:           fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `i_%s_%s`;", table.Name, table.Name, column.Name),
		IsDestructive: true,
	}
}

// dropUniqueIndex returns an alter table sql statement that drops a unique index
func dropUniqueIndex(table *schema.Table, column *schema.Column) *schema.SchemaChange {
	return &schema.SchemaChange{
		Type:          schema.DropIndex,
		SQL:           fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `ui_%s_%s`;", table.Name, table.Name, column.Name),
		IsDestructive: true,
	}
}

// dropPrimaryKey returns an alter table sql statement that drops a primary key
func dropPrimaryKey(table *schema.Table, column *schema.Column) *schema.SchemaChange {
	return &schema.SchemaChange{
		Type:          schema.DropIndex,
		SQL:           fmt.Sprintf("ALTER TABLE `%s` DROP PRIMARY KEY;", table.Name),
		IsDestructive: true,
	}
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
	case ColTypeVarchar, ColTypeEnum, ColTypeChar, ColTypeTinyText, ColTypeMediumText, ColTypeText, ColTypeLongText, ColTypeDate, ColTypeDateTime:
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

// createColumnSegment returns a table column sql segment
// Data Types
// INT SIGNED 	11 columns
func createColumnSegment(column *schema.Column) (sql string) {

	if isInt(column.DataType) {

		sql = fmt.Sprintf("`%s` %s(%d)", column.Name, column.DataType, intColLength(column.DataType, column.IsUnsigned))

		if column.IsUnsigned {
			sql += fmt.Sprintf(" %s", SignedUnsigned)
		} else {
			sql += fmt.Sprintf(" %s", SignedSigned)
		}

	} else if isFixedPointType(column.DataType) {

		sql = fmt.Sprintf("`%s` %s(%d,%d)", column.Name, column.DataType, column.Precision, column.NumericScale)

		if column.IsUnsigned {
			sql += fmt.Sprintf(" %s", SignedUnsigned)
		} else {
			sql += fmt.Sprintf(" %s", SignedSigned)
		}
	} else if isFloatingPointType(column.DataType) {
		sql = fmt.Sprintf("`%s` %s(%d,%d)", column.Name, column.DataType, column.Precision, column.NumericScale)
		if column.IsUnsigned {
			sql += fmt.Sprintf(" %s", SignedUnsigned)
		} else {
			sql += fmt.Sprintf(" %s", SignedSigned)
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

	if len(column.CharSet) > 0 {
		sql += " CHARACTER SET " + column.CharSet
	}

	if len(column.Collation) > 0 {
		sql += " COLLATE " + column.Collation
	}

	if !column.IsNullable {
		sql += " NOT"
	}
	sql += " NULL"

	// Add single quotes to string default
	if hasDefaultString(column.DataType) {

		defaultString := ""
		// Just use the NULL default (instead of a null string) if the field is nullable and the default is NULL
		// Sometimes default strings include their own single quotes
		if column.IsNullable && column.Default == "NULL" {
			// Don't add quotes
			defaultString = column.Default
		} else if len(column.Default) > 0 && column.Default[0:1] == "'" {
			// Don't add quotes
			defaultString = column.Default
		} else {
			defaultString = "'" + column.Default + "'"
		}
		// fmt.Println("string datatype: ", column.DataType, defaultString)
		sql += fmt.Sprintf(" DEFAULT %s", defaultString)
	} else if len(column.Default) > 0 {
		sql += fmt.Sprintf(" DEFAULT %s", column.Default)
	}

	if len(column.Extra) > 0 {
		sql += " " + column.Extra
	}

	return

}
