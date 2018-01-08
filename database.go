package dvc

import (
	"database/sql"
	"fmt"
	"log"
	"path"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type DVCLog struct {
	LogType    string
	LogMessage string
}

type DVCSet struct {
	ID             int64
	DateCreated    int64
	Name           string
	Files          map[string]*DVCFile
	IsDeleted      bool
	SortedFileKeys []string
}

type DVCFile struct {
	ID          int64
	DateCreated int64
	FileHash    string
	Name        string
	DVCSetID    int64
	IsRun       bool
	IsDeleted   bool
	Content     string
	FullPath    string
}

// ToLogString parses the file name into a human readable string for logging the action
func (d *DVCFile) ToLogString() (logString string) {

	// 001_{action}_{target}
	fileNameParts := strings.Split(d.Name, "_")

	// action
	fileAction := fileNameParts[1]

	// target
	fileTarget := strings.Join(fileNameParts[2:], "_")

	// Remove the `.sql` extension
	fileTarget = fileTarget[0 : len(fileTarget)-4]
	switch fileAction {
	case "createTable":
		logString = fmt.Sprintf("Creating table `%s`", fileTarget)
	case "alterTable":
		fileTargetParts := strings.Split(fileTarget, "__")
		fileActionParts := strings.Split(fileTargetParts[1], "_")
		logString = fmt.Sprintf("Altering table `%s` - %s %s", fileTargetParts[0], fileActionParts[0], fileActionParts[1])
	case "dropTable":
		logString = fmt.Sprintf("Dropping table `%s`", fileTarget)
	case "createView":
		logString = fmt.Sprintf("Creating view `%s`", fileTarget)
	case "alterView":
		logString = fmt.Sprintf("Altering view `%s`", fileTarget)
	case "dropView":
		logString = fmt.Sprintf("Dropping view `%s`", fileTarget)
	case "insert":
		logString = fmt.Sprintf("Inserting data into `%s`", fileTarget)
	}

	return
}

// Database is a collection of methods for managing the dvc tables
// in the target database `database`
type Database struct {
	runID         int64
	name          string
	conn          *sql.DB
	host          string
	tran          *sql.Tx
	sets          map[string]*DVCSet
	sortedSetKeys []string
	tables        map[string]Table
	logs          []DVCLog
}

// Connect connects to a database
// Todo "named" connections
func (d *Database) Connect(user string, pass string) (e error) {
	var connectionString = user + ":" + pass + "@tcp(" + d.host + ")/" + d.name + "?charset=utf8"
	d.conn, e = sql.Open("mysql", connectionString)
	return
}

// Start creates a new dvc run
func (d *Database) Start() (e error) {
	var result sql.Result

	result, e = d.conn.Exec("INSERT INTO dvcRun (`dateCreated`) VALUES (UNIX_TIMESTAMP())")

	if e != nil {
		return
	}

	d.runID, e = result.LastInsertId()
	return
}

// CreateBaseTablesIfNotExists creates the required base tables on the target database if they don't exist
func (d *Database) CreateBaseTablesIfNotExists() (e error) {

	if !d.TableExists("dvcSet") {
		_, e = d.conn.Exec("create table if not exists dvcSet (`id` int unsigned not null primary key auto_increment, `dateCreated` int unsigned not null default 0, `name` varchar(64) not null default '', isDeleted tinyint unsigned not null default 0)")
		if e != nil {
			return
		}
	}

	if !d.TableExists("dvcFile") {
		_, e = d.conn.Exec("create table if not exists dvcFile(`id` int unsigned not null primary key auto_increment, `dateCreated` int unsigned not null default 0, `fileHash` varchar(255) not null default '', `name` varchar(255) not null default '', `dvcSetId` int unsigned not null default 0, `isRun` tinyint unsigned not null default 0, isDeleted tinyint unsigned not null default 0, content TEXT NOT NULL)")
		if e != nil {
			return
		}
	}

	if !d.TableExists("dvcRun") {
		_, e = d.conn.Exec("create table if not exists dvcRun(`id` int unsigned not null primary key auto_increment, `dateCreated` int not null default 0)")
		if e != nil {
			return
		}
	}

	if !d.TableExists("dvcLog") {
		_, e = d.conn.Exec("create table if not exists dvcLog(`id` int unsigned not null primary key auto_increment, `dateCreated` int unsigned not null default 0, `dvcRunId` int unsigned not null default 0, `logType` ENUM('error', 'info') not null default 'info', `logMessage` text not null)")
	}

	return
}

// DropBaseTables drops the tables used by DVC
func (d *Database) DropBaseTables() (e error) {

	_, e = d.conn.Exec("DROP TABLE `dvcSet`")
	if e != nil {
		return
	}

	_, e = d.conn.Exec("DROP TABLE `dvcFile`")
	if e != nil {
		return
	}

	_, e = d.conn.Exec("DROP TABLE `dvcRun`")
	if e != nil {
		return
	}

	_, e = d.conn.Exec("DROP TABLE `dvcLog`")
	return
}

// FetchFilesBySet fetches a list of files associated with a changeset
// from the target database database
func (d *Database) FetchFilesBySet(dvcSetID int64, dvcSetName string) (dvcFiles map[string]*DVCFile, sortedFileKeys []string, e error) {

	var rows *sql.Rows
	sortedFileKeys = []string{}

	if rows, e = d.conn.Query("SELECT * FROM dvcFile WHERE dvcSetId = ? and isDeleted = 0", dvcSetID); e != nil {
		d.logFatal("FetchFilesBySet() " + e.Error())
		return
	}

	if rows != nil {
		defer rows.Close()
	}

	dvcFiles = map[string]*DVCFile{}

	for rows.Next() {
		dvcFile := DVCFile{}
		e = rows.Scan(
			&dvcFile.ID,
			&dvcFile.DateCreated,
			&dvcFile.FileHash,
			&dvcFile.Name,
			&dvcFile.DVCSetID,
			&dvcFile.IsRun,
			&dvcFile.IsDeleted,
			&dvcFile.Content,
		)
		if e != nil {
			return
		}
		sortedFileKeys = append(sortedFileKeys, dvcFile.Name)
		dvcFile.FullPath = path.Join(dvcSetName, dvcFile.Name)
		dvcFiles[dvcFile.Name] = &dvcFile
	}

	sort.Strings(sortedFileKeys)

	return
}

// FetchSets fetches a list of files associated with a changeset
// from the target database database
func (d *Database) FetchSets() (e error) {

	var rows *sql.Rows

	if rows, e = d.conn.Query("SELECT * FROM dvcSet WHERE isDeleted = 0"); e != nil {
		return
	}

	if rows != nil {
		defer rows.Close()
	}

	d.sets = map[string]*DVCSet{}
	d.sortedSetKeys = []string{}

	for rows.Next() {
		dvcSet := DVCSet{}
		e = rows.Scan(
			&dvcSet.ID,
			&dvcSet.DateCreated,
			&dvcSet.Name,
			&dvcSet.IsDeleted,
		)
		if e != nil {
			log.Fatal(e)
			return
		}
		d.sortedSetKeys = append(d.sortedSetKeys, dvcSet.Name)
		dvcSet.Files, dvcSet.SortedFileKeys, e = d.FetchFilesBySet(dvcSet.ID, dvcSet.Name)
		if e != nil {
			log.Fatal(e)
		}
		d.sets[dvcSet.Name] = &dvcSet
	}
	sort.Strings(d.sortedSetKeys)
	return
}

// CreateSet creates a new set and adds it to database sets collection
func (d *Database) CreateSet(name string) (e error) {
	var stmt *sql.Stmt
	var result sql.Result

	if stmt, e = d.conn.Prepare("insert into dvcSet(`dateCreated`, `name`) values (UNIX_TIMESTAMP(), ?)"); e != nil {
		return
	}

	if result, e = stmt.Exec(name); e != nil {
		return
	}

	setID, e := result.LastInsertId()
	if e != nil {
		return
	}

	d.sets[name] = &DVCSet{ID: setID, Name: name, Files: map[string]*DVCFile{}}

	d.sortedSetKeys = append(d.sortedSetKeys, name)
	sort.Strings(d.sortedSetKeys)
	return
}

// DeleteSet marks a set as deleted
func (d *Database) DeleteSet(setName string) (e error) {

	var set *DVCSet
	var ok bool

	if set, ok = d.sets[setName]; !ok {
		e = fmt.Errorf("Set %s not found", setName)
		return
	}

	// Delete from database
	var stmt *sql.Stmt
	stmt, e = d.conn.Prepare("UPDATE dvcSet SET isDeleted = 1 WHERE id = ?")

	if e != nil {
		return
	}
	_, e = stmt.Exec(set.ID)
	if e != nil {
		return
	}

	// Delete from sortedName cache
	var deletedSetIdx int
	for idx, s := range d.sortedSetKeys {
		if setName == s {
			deletedSetIdx = idx
			break
		}
	}
	copy(d.sortedSetKeys[deletedSetIdx:], d.sortedSetKeys[deletedSetIdx+1:])
	d.sortedSetKeys[len(d.sortedSetKeys)-1] = ""
	d.sortedSetKeys = d.sortedSetKeys[:len(d.sortedSetKeys)-1]

	// Delete from map
	delete(d.sets, setName)

	return
}

// CreateFile creates a new changeset file entry in the database
func (d *Database) CreateFile(set *DVCSet, name string, fileHash string, content string) (f *DVCFile, e error) {
	var stmt *sql.Stmt
	var result sql.Result
	var fileID int64

	if stmt, e = d.conn.Prepare("insert into dvcFile(`dateCreated`, `dvcSetId`, `name`, `fileHash`, `content`) values (UNIX_TIMESTAMP(), ?, ?, ?, ?)"); e != nil {
		return
	}

	if result, e = stmt.Exec(set.ID, name, fileHash, content); e != nil {
		return
	}

	fileID, e = result.LastInsertId()
	if e != nil {
		return
	}

	fullPath := path.Join(set.Name, name)
	f = &DVCFile{ID: fileID, DVCSetID: set.ID, Name: name, FileHash: fileHash, IsRun: false, FullPath: fullPath}
	set.Files[name] = f
	set.SortedFileKeys = append(set.SortedFileKeys, name)
	sort.Strings(set.SortedFileKeys)
	return
}

// DeleteFile marks a file as deleted
func (d *Database) DeleteFile(set *DVCSet, file *DVCFile) (e error) {

	// Delete from database
	var stmt *sql.Stmt
	stmt, e = d.conn.Prepare("UPDATE dvcFile SET isDeleted = 1 WHERE id = ?")

	if e != nil {
		return
	}
	_, e = stmt.Exec(file.ID)

	if e != nil {
		return
	}

	// Delete from sortedName cache
	var deletedFileIdx int
	for idx, f := range set.SortedFileKeys {
		if file.Name == f {
			deletedFileIdx = idx
			break
		}
	}
	copy(set.SortedFileKeys[deletedFileIdx:], set.SortedFileKeys[deletedFileIdx+1:])
	set.SortedFileKeys[len(set.SortedFileKeys)-1] = ""
	set.SortedFileKeys = set.SortedFileKeys[:len(set.SortedFileKeys)-1]

	// Delete from map
	delete(set.Files, file.Name)

	return
}

// FetchTables fetches the complete set of tables from this database
// and populates the Tables map with a collection of Table objects
func (d *Database) FetchTables() (e error) {

	var rows *sql.Rows

	if rows, e = d.conn.Query("SHOW TABLES"); e != nil {
		return
	}

	if rows != nil {
		defer rows.Close()
	}

	d.tables = map[string]Table{}

	for rows.Next() {
		table := Table{}
		rows.Scan(
			&table.Name,
		)
		d.tables[table.Name] = table
		// fmt.Printf("Table: %s\n", table.Name)
	}

	return
}

func (d *Database) SetExists(changesetName string) (ok bool) {
	_, ok = d.sets[changesetName]
	return
}

func (d *Database) FileExists(changesetName string, fileName string) (ok bool) {
	if !d.SetExists(changesetName) {
		ok = false
		return
	}

	_, ok = d.sets[changesetName].Files[fileName]
	return
}

func (d *Database) TableExists(tableName string) (ok bool) {
	_, ok = d.tables[tableName]
	return
}

func (d *Database) RunChange(sqlString string) (e error) {
	_, e = d.conn.Exec(sqlString)
	return
}

func (d *Database) SetFileAsRun(dvcFileId int64) (e error) {
	var stmt *sql.Stmt

	stmt, e = d.conn.Prepare("UPDATE dvcFile SET isRun = 1 WHERE id = ?")
	if e != nil {
		return
	}

	_, e = stmt.Exec(dvcFileId)
	return
}

func (d *Database) logFatal(msg string) {
	d.logs = append(d.logs, DVCLog{LogMessage: msg, LogType: "error"})
	d.writeLogs()
	log.Fatal(msg)
}

func (d *Database) log(l string) {
	log.Println(l)
	d.logs = append(d.logs, DVCLog{LogMessage: l, LogType: "info"})
}

func (d *Database) writeLogs() {
	if len(d.logs) > 0 {
		tx, e := d.conn.Begin()
		if e != nil {
			log.Fatal(e)
		}

		for _, l := range d.logs {
			stmt, e := tx.Prepare("INSERT INTO dvcLog(dateCreated, logType, logMessage, dvcRunId) VALUES (UNIX_TIMESTAMP(), ?, ?, ?)")
			if e != nil {
				log.Fatal(e)
			}

			_, e = stmt.Exec(l.LogType, l.LogMessage, d.runID)
			if e != nil {
				log.Fatal(e)
			}
		}

		tx.Commit()
	}
}

func (d *Database) finish() {
	d.writeLogs()
}
