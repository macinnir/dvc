package lib

// DatabaseType is the type of database to be used.
type DatabaseType string

const (
	// DatabaseTypeMysql is the MySQL flavor of database
	DatabaseTypeMysql DatabaseType = "mysql"
	// DatabaseTypeSqlite is the Sqlite flavor of database
	DatabaseTypeSqlite DatabaseType = "sqlite"
)

// Options are the available runtime flags
type Options uint

// Command is the command line functionality
type Command struct {
	Options Options
}

// Changeset represents all of the changes in an environment and their changes
type Changeset struct {
	ChangeFiles map[string]ChangeFile
	Signature   string
}

// ChangeFile represents both a physical file on the local file system
// along with the entry in the changefile database
type ChangeFile struct {
	ID          int64
	DateCreated int64
	Hash        string
	Name        string
	DVCSetID    int64
	IsRun       bool
	IsDeleted   bool
	Content     string
	FullPath    string
	Ordinal     int
}

const (
	// OptLogInfo triggers verbose logging
	OptLogInfo = 1 << iota
	// OptLogDebug triggers extremely verbose logging
	OptLogDebug
	// OptSilent suppresses all logging
	OptSilent
	// OptReverse reverses the function
	OptReverse
	// OptSummary shows a summary of the actions instead of a raw stdout dump
	OptSummary
	// OptClean cleans
	OptClean
	// OptForce forces
	OptForce
)
