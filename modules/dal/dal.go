package dal

import (
	"database/sql"
)

// Dal manages the connection to and the querying of the database
type Dal struct {
	Connection *sql.DB
	Schema     *Schema
}

// Connect connects to the database
func (d *Dal) Connect(dbConnection *sql.DB) (e error) {
	d.Connection = dbConnection
	return
}

// NewDal assigns a new Dal instance to the package global DB
func NewDal(dbConnection *sql.DB) *Dal {
	dal := new(Dal)
	dal.Connect(dbConnection)
	dal.Schema = NewSchema(dal, "test")
	return dal
}
