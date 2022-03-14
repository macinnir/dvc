package db

import (
	"context"
	"fmt"

	"github.com/macinnir/dvc/core/lib/utils/log"

	// Mysql Package
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host string
	Name string
	User string
	Pass string
}

// MySQL is mysql
type MySQL struct {
	config *Config
	db     *sqlx.DB
	log    log.ILog
}

// New returns a new MySQL object
func New(config *Config, log log.ILog) IDB {
	m := &MySQL{
		config: config,
		log:    log,
	}

	m.connect()
	return m
}

func (m *MySQL) Host() string {
	return m.config.Host
}

func (m *MySQL) Name() string {
	return m.config.Name
}

// connect connects to the database
func (m *MySQL) connect() {

	var e error
	m.log.Printf(
		"MARIADB: Connecting to %s/%s with user %s",
		m.config.Host,
		m.config.Name,
		m.config.User,
	)

	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", m.config.User, m.config.Pass, m.config.Host, m.config.Name)

	if m.db, e = sqlx.Connect("mysql", dbConnectionString); e != nil {
		m.log.Fatalf("ERROR: Database Connection: %s", e.Error())
		return
	}

	m.log.Println("MARIADB: CONNECTED")

	return
}

// NamedExec using this DB.
// Any named placeholder parameters are replaced with fields from arg.
func (m *MySQL) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return m.db.NamedExec(query, arg)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (m *MySQL) Exec(query string, args ...interface{}) (sql.Result, error) {
	return m.db.Exec(query, args...)
}

// Get using this DB.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (m *MySQL) Get(dest interface{}, query string, args ...interface{}) error {
	return m.db.Get(dest, query, args...)
}

// Select using this DB.
// Any placeholder parameters are replaced with supplied args.
func (m *MySQL) Select(dest interface{}, query string, args ...interface{}) error {
	return m.db.Select(dest, query, args...)
}

// Close closes the database and prevents new queries from starting.
// Close then waits for all queries that have started processing on the server
// to finish.
//
// It is rare to Close a DB, as the DB handle is meant to be
// long-lived and shared between many goroutines.
func (m *MySQL) Close() {
	m.db.Close()
}

// BeginTx starts a transaction.
//
// The provided context is used until the transaction is committed or rolled back.
// If the context is canceled, the sql package will roll back
// the transaction. Tx.Commit will return an error if the context provided to
// BeginTx is canceled.
//
// The provided TxOptions is optional and may be nil if defaults should be used.
// If a non-default isolation level is used that the driver doesn't support,
// an error will be returned.
func (m *MySQL) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return m.db.BeginTx(ctx, opts)
}

func buildExecManyChunks(stmts []string, chunkSize int) [][]string {

	chunks := [][]string{}

	// No records
	if len(stmts) == 0 {
		chunks = [][]string{{}}
		return chunks
	}

	// Don't use a transaction if only a single value
	if len(stmts) == 1 {
		chunks = append(chunks, []string{stmts[0]})
		return chunks
	}

	for i := 0; i < len(stmts); i += chunkSize {
		end := i + chunkSize
		if end > len(stmts) {
			end = len(stmts)
		}
		chunks = append(chunks, stmts[i:end])
	}

	return chunks
}

// UpdateMany updates a slice of User objects in chunks
func (m *MySQL) ExecMany(stmts []string, chunkSize int) (e error) {

	chunks := buildExecManyChunks(stmts, chunkSize)

	for k := range chunks {

		var tx *sql.Tx
		ctx := context.Background()
		tx, e = m.BeginTx(ctx, nil)

		if e != nil {
			return
		}

		for l := range chunks[k] {

			if _, e = tx.ExecContext(ctx, chunks[k][l]); e != nil {
				tx.Rollback()
				return
			}

		}

		if e != nil {
			return
		}

		e = tx.Commit()
	}

	return

}
