package lib

import (
	"fmt"
	"strings"
)

// Executor executes commands
type Executor struct {
	config    *Config    // Config is the config object
	connector IConnector // IConnector is the injected server manager
	server    *Server
}

// NewExecutor returns a new Executor
func NewExecutor(config *Config, connector IConnector) *Executor {
	return &Executor{
		config:    config,
		connector: connector,
	}
}

// Connect connects
func (x *Executor) Connect() (server *Server) {

	var e error
	server, e = x.connector.Connect()

	if e != nil {
		panic(e)
	}

	e = x.connector.UseDatabase(server, x.config.Connection.DatabaseName)

	return

}

// RunSQL runs the sql produced by the `CompareSchema` command against the target database
// @command compare [reverse] apply
func (x *Executor) RunSQL(sql string) (e error) {

	server := x.Connect()

	statements := strings.Split(sql, ";")

	defer server.Connection.Close()

	tx, _ := server.Connection.Begin()

	nonEmptyStatements := []string{}
	for _, s := range statements {
		if len(strings.Trim(strings.Trim(s, " "), "\n")) == 0 {
			continue
		}

		nonEmptyStatements = append(nonEmptyStatements, s)
	}

	for i, s := range nonEmptyStatements {
		sql := strings.Trim(strings.Trim(s, " "), "\n")
		if len(sql) == 0 {
			continue
		}
		// fmt.Printf("\rRunning %d of %d sql statements...", i+1, len(nonEmptyStatements))
		fmt.Printf("Running %d of %d: \n%s\n", i+1, len(nonEmptyStatements), sql)
		// lib.Debugf("Running sql: \n%s\n", c.Options, sql)

		_, e = tx.Exec(sql)
		if e != nil {
			tx.Rollback()
			return
		}
	}
	fmt.Print("Finished\n")
	e = tx.Commit()
	if e != nil {
		panic(e)
	}

	return
}
