package query

import (
	"context"
	"database/sql"
)

// DBInterface is a container for db interactions
type DBInterface interface {
	Close()
	Exec(query string, args ...interface{}) (sql.Result, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	ExecMany(stmts []string, chunkSize int) (e error)
	Host() string // The host name (from config)
	Name() string // The name of the database (from config)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
