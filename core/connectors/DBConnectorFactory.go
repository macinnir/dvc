package connectors

import (
	"errors"

	"github.com/macinnir/dvc/core/connectors/mysql"
	"github.com/macinnir/dvc/core/connectors/sqlite"
	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

// DBConnectorFactory is a factory method for database connections
func DBConnectorFactory(config *lib.ConfigDatabase) (connector IConnector, e error) {

	switch config.Type {
	case schema.SchemaTypeMySQL:
		connector = mysql.NewMySQL(config)
	case schema.SchemaTypeSQLite:
		connector = sqlite.NewSqlite(config)
	default:
		e = errors.New("invalid database type")
	}

	return
}

// IConnector defines the shape of a connector to a database
type IConnector interface {
	Connect() (server *schema.Server, e error)
	FetchDatabases(server *schema.Server) (databases map[string]*schema.Schema, e error)
	// FetchEnums(server *Server) (enums map[string][]map[string]interface{})
	FetchEnum(server *schema.Server, tableName string) []map[string]interface{}
	UseDatabase(server *schema.Server, databaseName string) (e error)
	FetchDatabase(server *schema.Server, databaseName string) (schema *schema.Schema, e error)
	FetchTableColumns(server *schema.Server, databaseName string, tableName string) (columns map[string]*schema.Column, e error)
	CreateChangeSQL(localSchema *schema.Schema, remoteSchema *schema.Schema, databaseName string) (s *schema.SchemaComparison)
	// CompareEnums(remoteSchema *schema.Schema, localSchema *schema.Schema, tableName string) (sql string)
}
