package main

import (
	"database/sql"
	"fmt"
	// "fmt"
	// "log"
)

// ServerMgr manages a server
type ServerMgr struct {
	Config          *Config
	conn            *sql.DB
	currentDatabase string
}

// Connect connects to a server
func (s *ServerMgr) Connect() (e error) {

	if s.conn == nil {
		var connectionString = s.Config.Username + ":" + s.Config.Password + "@tcp(" + s.Config.Host + ")/?charset=utf8"
		// log.Printf("Connecting to server: %s", connectionString)
		s.conn, e = sql.Open("mysql", connectionString)
	}
	return
}

// FetchDatabases fetches a set of database names from the target server
// populating the Databases property with a map of Database objects
func (s *ServerMgr) FetchDatabases() (databases map[string]*Database, e error) {

	var rows *sql.Rows
	databases = map[string]*Database{}

	if rows, e = s.conn.Query("SHOW DATABASES"); e != nil {
		return
	}

	if rows != nil {
		defer rows.Close()
	}

	for rows.Next() {
		databaseName := ""
		rows.Scan(&databaseName)
		databases[databaseName] = &Database{Name: databaseName, Host: s.Config.Host}
	}

	return
}

// UseDatabase switches the connection context to the passed in database
func (s *ServerMgr) UseDatabase(databaseName string) (e error) {

	if s.currentDatabase == databaseName {
		return
	}

	_, e = s.conn.Exec(fmt.Sprintf("USE %s", databaseName))
	return
}

// // CreateDatabase creates a new databse
// func (s *Server) CreateDatabase(databaseName string) (database *Database, e error) {
// 	_, e = s.conn.Exec(fmt.Sprintf("CREATE DATABASE `%s`", databaseName))
// 	if e != nil {
// 		return
// 	}
// 	s.Databases[databaseName] = &Database{name: databaseName, host: s.Name}
// 	database = s.Databases[databaseName]
// 	return
// }

// // NewServer creates a new instance of a Server object
// func NewServer(host string, username string, password string) (s Server) {
// 	s = Server{Name: host}
// 	e := s.Connect(username, password)
// 	if e != nil {
// 		log.Fatal(e)
// 	}

// 	return
// }
