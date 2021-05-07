package executor

import (
	"strings"

	"github.com/macinnir/dvc/core/connectors"
	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
	"go.uber.org/zap"
)

// Executor executes commands
type Executor struct {
	config    *lib.ConfigDatabase   // Config is the config object
	connector connectors.IConnector // IConnector is the injected server manager
	server    *schema.Server        // Server
	log       *zap.Logger           // Logger
}

// NewExecutor returns a new Executor
func NewExecutor(config *lib.ConfigDatabase, connector connectors.IConnector) *Executor {
	return &Executor{
		config:    config,
		connector: connector,
	}
}

// Connect connects
func (x *Executor) Connect() (server *schema.Server) {

	var e error
	server, e = x.connector.Connect()

	// TODO log
	if e != nil {
		panic(e)
	}

	e = x.connector.UseDatabase(server, x.config.Name)

	// TODO log
	if e != nil {
		panic(e)
	}

	return

}

func parseSQLStatements(str string) []string {

	stmts := []string{}

	if strings.Contains(str, ";") {
		stmts = strings.Split(str, ";")
	}

	nonEmptyStatements := []string{}
	for k := range stmts {
		l := strings.TrimSpace(stmts[k])
		if len(l) == 0 {
			continue
		}

		nonEmptyStatements = append(nonEmptyStatements, l)
	}

	return nonEmptyStatements
}

// RunSQL runs the sql produced by the `CompareSchema` command against the target database
// @command compare [reverse] apply
func (x *Executor) RunSQL(sql string) (e error) {

	server := x.Connect()
	defer server.Connection.Close()
	tx, _ := server.Connection.Begin()

	stmts := parseSQLStatements(sql)

	stmtLen := len(stmts)

	for k := range stmts {
		x.log.Sugar().Infof("Running %d of %d: %s", k+1, stmtLen, sql)
		_, e = tx.Exec(sql)
		if e != nil {
			tx.Rollback()
			return
		}
	}

	x.log.Info("Finished")
	e = tx.Commit()
	if e != nil {
		panic(e)
	}

	return
}
