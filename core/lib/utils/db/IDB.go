package db

import (
	"context"
	"database/sql"
)

// IDB is a container for db interactions
type IDB interface {
	Close()
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	ExecMany(stmts []string, chunkSize int) (e error)
	Host() string // The host name (from config)
	Name() string // The name of the database (from config)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}
