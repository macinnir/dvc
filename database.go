package main

import (
	"database/sql"
	"fmt"
	// "database/sql"
	// "fmt"
	// "log"
	// "path"
	// "sort"
)

// Database Types
const (
	DatabaseTypeMysql     = "mysql"
	DatabaseTypeSQLServer = "sqlserver"
	DatabaseTypePostGRES  = "postgres"
	DatabaseTypeSQLite    = "sqlite"
)

var DatabaseTypes = []string{
	DatabaseTypeMysql,
	DatabaseTypeSQLServer,
	DatabaseTypePostGRES,
	DatabaseTypeSQLite,
}

// IRepo defines the type of methods needed for a database integration
type IRepo interface {
	Connect() (e error)
	CreateRun() (runID int64, e error)
	FetchAllTableNames() (tableNames map[string]string, e error)
	CreateBaseTablesIfNotExists() (e error)
	DropBaseTables() (e error)
	// FetchFilesBySet(dvcSetID int64, dvcSetName string) (dvcFiles map[string]*ChangeFile, sortedFileKeys []string, e error)
	// FetchSets() (e error)
	// CreateSet(name string) (e error)
	// DeleteSet(setName string) (e error)
	// CreateFile(set *ChangeSet, name string, fileHash string, content string) (f *ChangeFile, e error)
	// DeleteFile(set *ChangeSet, file *ChangeFile) (e error)
	// FetchTables() (e error)
	// SetExists(changesetName string) (ok bool)
	// FileExists(changesetName string, fileName string) (ok bool)
	// TableExists(tableName string) (ok bool)
	// RunChange(sqlString string) (e error)
	// SetFileAsRun(dvcFileId int64) (e error)
}

// func NewDatabaseMgr(config *Config) (e error, d *DatabaseMgr) {

// 	d = &DatabaseMgr{
// 		config: config,
// 	}

// 	return
// }

type DatabaseMgr struct {
	databaseName string
	connection   *sql.DB
}

// ListTables fetches the complete set of tables from this database
func (d *DatabaseMgr) ListTables(databaseName string) (tables map[string]*Table, e error) {

	var rows *sql.Rows
	query := "select `TABLE_NAME`, `ENGINE`, `VERSION`, `ROW_FORMAT`, `TABLE_ROWS`, `DATA_LENGTH`, `TABLE_COLLATION`, `AUTO_INCREMENT` FROM information_schema.tables WHERE TABLE_SCHEMA = '" + databaseName + "'"
	// fmt.Printf("Query: %s\n", query)
	if rows, e = d.connection.Query(query); e != nil {
		return
	}

	if rows != nil {
		defer rows.Close()
	}

	tables = map[string]*Table{}

	for rows.Next() {
		table := &Table{}
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
		tables[table.Name] = table
	}

	return
}

func (d *DatabaseMgr) ListColumns(databaseName string, tableName string) (columns map[string]*Column, e error) {
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
			EXTRA
		FROM information_schema.COLUMNS 
		WHERE 
			TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'
	`, databaseName, tableName)

	if rows, e = d.connection.Query(query); e != nil {
		return
	}

	if rows != nil {
		defer rows.Close()
	}

	columns = map[string]*Column{}

	for rows.Next() {
		column := Column{}
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
		); e != nil {
			return
		}
		columns[column.Name] = &column
	}

	return
}

func (d *DatabaseMgr) BuildDatabase(databaseName string) (database *Database, e error) {

	database = &Database{
		Name: databaseName,
	}

	if database.Tables, e = d.ListTables(database.Name); e != nil {
		return
	}

	tableLen := len(database.Tables)

	if tableLen == 0 {
		// e = fmt.Errorf("no tables found in database %s", database.Name)
		return
	}

	// fmt.Printf("Found %d tables\n", tableLen)

	keys := make([]string, tableLen)
	idx := 0

	for tableName := range database.Tables {
		keys[idx] = tableName
		idx = idx + 1
	}

	// fmt.Printf("Keys %d\n", len(keys))

	for _, tableName := range keys {

		table := database.Tables[tableName]

		// fmt.Printf("\tTable: %s - %s\n", tableName, table.Name)
		table.Columns, e = d.ListColumns(database.Name, table.Name)

		// for columName, column := range columns {

		// 	fmt.Printf("\t\tColumn: %s - %s\n", columName, column.Type)
		// }

		database.Tables[tableName] = table
	}

	return
}

// func (d *DatabaseMgr) InitIntegration() (e error) {
// 	switch d.config.DatabaseType {
// 	case "mysql":
// 		d.db, e = NewMySQL(d.config)
// 	default:
// 		e = errors.New("invalid database type")
// 	}
// 	return
// }

// Database is a collection of methods for managing the dvc tables
// in the target database `database`

// Connect connects to a database
// Todo "named" connections
// func (d *Database) Connect(host string, name string, user string, pass string) (e error) {
// 	d.db.Connect(host, name, user, pass)
// }

// // StartRun creates a new dvc run
// func (d *Database) StartRun() (e error) {
// 	d.tableNames, e = d.db.FetchAllTableNames()
// 	if e != nil {
// 		return
// 	}
// 	d.runID, e = d.db.CreateRun()
// 	return
// }

// CreateBaseTablesIfNotExists creates the required base tables on the target database if they don't exist
// func (d *Database) InitDvcTables() (e error) {

// 	if !d.TableExists("dvcSet") {
// 		_, e = d.conn.Exec("create table if not exists dvcSet (`id` int unsigned not null primary key auto_increment, `dateCreated` int unsigned not null default 0, `name` varchar(64) not null default '', isDeleted tinyint unsigned not null default 0)")
// 		if e != nil {
// 			return
// 		}
// 	}

// 	if !d.TableExists("dvcFile") {
// 		_, e = d.conn.Exec("create table if not exists dvcFile(`id` int unsigned not null primary key auto_increment, `dateCreated` int unsigned not null default 0, `fileHash` varchar(255) not null default '', `name` varchar(255) not null default '', `dvcSetId` int unsigned not null default 0, `isRun` tinyint unsigned not null default 0, isDeleted tinyint unsigned not null default 0, content TEXT NOT NULL)")
// 		if e != nil {
// 			return
// 		}
// 	}

// 	if !d.TableExists("dvcRun") {
// 		_, e = d.conn.Exec("create table if not exists dvcRun(`id` int unsigned not null primary key auto_increment, `dateCreated` int not null default 0)")
// 		if e != nil {
// 			return
// 		}
// 	}

// 	if !d.TableExists("dvcLog") {
// 		_, e = d.conn.Exec("create table if not exists dvcLog(`id` int unsigned not null primary key auto_increment, `dateCreated` int unsigned not null default 0, `dvcRunId` int unsigned not null default 0, `logType` ENUM('error', 'info') not null default 'info', `logMessage` text not null)")
// 	}

// 	return
// }

// // DropBaseTables drops the tables used by DVC
// func (d *Database) DropBaseTables() (e error) {

// 	_, e = d.conn.Exec("DROP TABLE `dvcSet`")
// 	if e != nil {
// 		return
// 	}

// 	_, e = d.conn.Exec("DROP TABLE `dvcFile`")
// 	if e != nil {
// 		return
// 	}

// 	_, e = d.conn.Exec("DROP TABLE `dvcRun`")
// 	if e != nil {
// 		return
// 	}

// 	_, e = d.conn.Exec("DROP TABLE `dvcLog`")
// 	return
// }

// // FetchFilesBySet fetches a list of files associated with a changeset
// // from the target database database
// func (d *Database) FetchFilesBySet(dvcSetID int64, dvcSetName string) (dvcFiles map[string]*DVCFile, sortedFileKeys []string, e error) {

// 	var rows *sql.Rows
// 	sortedFileKeys = []string{}

// 	if rows, e = d.conn.Query("SELECT * FROM dvcFile WHERE dvcSetId = ? and isDeleted = 0", dvcSetID); e != nil {
// 		d.logFatal("FetchFilesBySet() " + e.Error())
// 		return
// 	}

// 	if rows != nil {
// 		defer rows.Close()
// 	}

// 	dvcFiles = map[string]*DVCFile{}

// 	for rows.Next() {
// 		dvcFile := DVCFile{}
// 		e = rows.Scan(
// 			&dvcFile.ID,
// 			&dvcFile.DateCreated,
// 			&dvcFile.FileHash,
// 			&dvcFile.Name,
// 			&dvcFile.DVCSetID,
// 			&dvcFile.IsRun,
// 			&dvcFile.IsDeleted,
// 			&dvcFile.Content,
// 		)
// 		if e != nil {
// 			return
// 		}
// 		sortedFileKeys = append(sortedFileKeys, dvcFile.Name)
// 		dvcFile.FullPath = path.Join(dvcSetName, dvcFile.Name)
// 		dvcFiles[dvcFile.Name] = &dvcFile
// 	}

// 	sort.Strings(sortedFileKeys)

// 	return
// }

// // FetchSets fetches a list of files associated with a changeset
// // from the target database database
// func (d *Database) FetchSets() (e error) {

// 	var rows *sql.Rows

// 	if rows, e = d.conn.Query("SELECT * FROM dvcSet WHERE isDeleted = 0"); e != nil {
// 		return
// 	}

// 	if rows != nil {
// 		defer rows.Close()
// 	}

// 	d.sets = map[string]*DVCSet{}
// 	d.sortedSetKeys = []string{}

// 	for rows.Next() {
// 		dvcSet := DVCSet{}
// 		e = rows.Scan(
// 			&dvcSet.ID,
// 			&dvcSet.DateCreated,
// 			&dvcSet.Name,
// 			&dvcSet.IsDeleted,
// 		)
// 		if e != nil {
// 			log.Fatal(e)
// 			return
// 		}
// 		d.sortedSetKeys = append(d.sortedSetKeys, dvcSet.Name)
// 		dvcSet.Files, dvcSet.SortedFileKeys, e = d.FetchFilesBySet(dvcSet.ID, dvcSet.Name)
// 		if e != nil {
// 			log.Fatal(e)
// 		}
// 		d.sets[dvcSet.Name] = &dvcSet
// 	}
// 	sort.Strings(d.sortedSetKeys)
// 	return
// }

// // CreateSet creates a new set and adds it to database sets collection
// func (d *Database) CreateSet(name string) (e error) {
// 	var stmt *sql.Stmt
// 	var result sql.Result

// 	if stmt, e = d.conn.Prepare("insert into dvcSet(`dateCreated`, `name`) values (UNIX_TIMESTAMP(), ?)"); e != nil {
// 		return
// 	}

// 	if result, e = stmt.Exec(name); e != nil {
// 		return
// 	}

// 	setID, e := result.LastInsertId()
// 	if e != nil {
// 		return
// 	}

// 	d.sets[name] = &DVCSet{ID: setID, Name: name, Files: map[string]*DVCFile{}}

// 	d.sortedSetKeys = append(d.sortedSetKeys, name)
// 	sort.Strings(d.sortedSetKeys)
// 	return
// }

// // DeleteSet marks a set as deleted
// func (d *Database) DeleteSet(setName string) (e error) {

// 	var set *DVCSet
// 	var ok bool

// 	if set, ok = d.sets[setName]; !ok {
// 		e = fmt.Errorf("Set %s not found", setName)
// 		return
// 	}

// 	// Delete from database
// 	var stmt *sql.Stmt
// 	stmt, e = d.conn.Prepare("UPDATE dvcSet SET isDeleted = 1 WHERE id = ?")

// 	if e != nil {
// 		return
// 	}
// 	_, e = stmt.Exec(set.ID)
// 	if e != nil {
// 		return
// 	}

// 	// Delete from sortedName cache
// 	var deletedSetIdx int
// 	for idx, s := range d.sortedSetKeys {
// 		if setName == s {
// 			deletedSetIdx = idx
// 			break
// 		}
// 	}
// 	copy(d.sortedSetKeys[deletedSetIdx:], d.sortedSetKeys[deletedSetIdx+1:])
// 	d.sortedSetKeys[len(d.sortedSetKeys)-1] = ""
// 	d.sortedSetKeys = d.sortedSetKeys[:len(d.sortedSetKeys)-1]

// 	// Delete from map
// 	delete(d.sets, setName)

// 	return
// }

// // CreateFile creates a new changeset file entry in the database
// func (d *Database) CreateFile(set *DVCSet, name string, fileHash string, content string) (f *DVCFile, e error) {

// 	var stmt *sql.Stmt
// 	var result sql.Result
// 	var fileID int64

// 	if stmt, e = d.conn.Prepare(`
// 		INSERT INTO dvcFile (
// 			dateCreated,
// 			dvcSetId,
// 			name,
// 			fileHash,
// 			content
// 		) values (
// 			UNIX_TIMESTAMP(),
// 			?,
// 			?,
// 			?,
// 			?
// 		)`); e != nil {
// 		return
// 	}

// 	if result, e = stmt.Exec(set.ID, name, fileHash, content); e != nil {
// 		return
// 	}

// 	fileID, e = result.LastInsertId()
// 	if e != nil {
// 		return
// 	}

// 	fullPath := path.Join(set.Name, name)
// 	f = &DVCFile{ID: fileID, DVCSetID: set.ID, Name: name, FileHash: fileHash, IsRun: false, FullPath: fullPath}
// 	set.Files[name] = f
// 	set.SortedFileKeys = append(set.SortedFileKeys, name)
// 	sort.Strings(set.SortedFileKeys)
// 	return
// }

// // DeleteFile marks a file as deleted
// func (d *Database) DeleteFile(set *DVCSet, file *DVCFile) (e error) {

// 	// Delete from database
// 	var stmt *sql.Stmt
// 	stmt, e = d.conn.Prepare("UPDATE dvcFile SET isDeleted = 1 WHERE id = ?")

// 	if e != nil {
// 		return
// 	}
// 	_, e = stmt.Exec(file.ID)

// 	if e != nil {
// 		return
// 	}

// 	// Delete from sortedName cache
// 	var deletedFileIdx int
// 	for idx, f := range set.SortedFileKeys {
// 		if file.Name == f {
// 			deletedFileIdx = idx
// 			break
// 		}
// 	}
// 	copy(set.SortedFileKeys[deletedFileIdx:], set.SortedFileKeys[deletedFileIdx+1:])
// 	set.SortedFileKeys[len(set.SortedFileKeys)-1] = ""
// 	set.SortedFileKeys = set.SortedFileKeys[:len(set.SortedFileKeys)-1]

// 	// Delete from map
// 	delete(set.Files, file.Name)

// 	return
// }

// // TableExists returns true if the table currently exists in the database, false if not
// func (d *Database) TableExists(tableName string) (exists bool) {

// }

// // SetExists verifies that a changeset exists
// func (d *Database) SetExists(changesetName string) (ok bool) {
// 	_, ok = d.sets[changesetName]
// 	return
// }

// // FileExists verifies that a file exists in the list of changesets
// func (d *Database) FileExists(changesetName string, fileName string) (ok bool) {
// 	if !d.SetExists(changesetName) {
// 		ok = false
// 		return
// 	}

// 	_, ok = d.sets[changesetName].Files[fileName]
// 	return
// }

// // TableExists verifies that a table exists within the colleciton of tables
// func (d *Database) TableExists(tableName string) (ok bool) {
// 	_, ok = d.tables[tableName]
// 	return
// }

// // RunChange runs a sql string
// func (d *Database) RunChange(sqlString string) (e error) {
// 	_, e = d.conn.Exec(sqlString)
// 	return
// }

// // SetFileAsRun sets a file as being run
// // Prevents a logged file in the database from running this file again
// func (d *Database) SetFileAsRun(dvcFileID int64) (e error) {
// 	var stmt *sql.Stmt

// 	stmt, e = d.conn.Prepare("UPDATE dvcFile SET isRun = 1 WHERE id = ?")
// 	if e != nil {
// 		return
// 	}

// 	_, e = stmt.Exec(dvcFileID)
// 	return
// }

// func (d *Database) finish() {
// 	d.writeLogs()
// }
