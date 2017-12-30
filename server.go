package dvc

import (
	"database/sql"
	"fmt"
	"log"
)

type Server struct {
	Name      string `json:"name"`
	Databases map[string]*Database
	conn      *sql.DB
}

// FetchDatabases fetches a set of database names from the target server
// populating the Databases property with a map of Database objects
func (s *Server) FetchDatabases() (e error) {

	var rows *sql.Rows
	var databases = map[string]*Database{}

	if rows, e = s.conn.Query("SHOW DATABASES"); e != nil {
		return
	}

	if rows != nil {
		defer rows.Close()
	}

	for rows.Next() {
		databaseName := ""
		rows.Scan(&databaseName)
		databases[databaseName] = &Database{name: databaseName, host: s.Name}
	}

	s.Databases = databases

	return
}

// Connect connects to a server
func (s *Server) Connect(user string, pass string) (e error) {
	var connectionString = user + ":" + pass + "@tcp(" + s.Name + ":3306)/?charset=utf8"
	log.Printf("Connecting to server: %s", connectionString)
	s.conn, e = sql.Open("mysql", connectionString)
	return
}

// CreateDatabase creates a new databse
func (s *Server) CreateDatabase(databaseName string) (database *Database, e error) {
	_, e = s.conn.Exec(fmt.Sprintf("CREATE DATABASE `%s`", databaseName))
	if e != nil {
		return
	}
	s.Databases[databaseName] = &Database{name: databaseName, host: s.Name}
	database = s.Databases[databaseName]
	return
}

// NewServer creates a new instance of a Server object
func NewServer(host string, username string, password string) (s Server) {
	s = Server{Name: host}
	e := s.Connect(username, password)
	if e != nil {
		log.Fatal(e)
	}

	return
}
