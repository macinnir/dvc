package compare

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"github.com/macinnir/dvc/lib"
	"io/ioutil"
	"strings"
)

func schemasAreSame(localSchema *lib.Database, remoteSchema *lib.Database) (same bool, e error) {

	// Remote Signature
	var localBytes []byte
	var remoteBytes []byte

	localBytes, e = json.Marshal(localSchema)
	if e != nil {
		return
	}

	remoteBytes, e = json.Marshal(remoteSchema)
	if e != nil {
		return
	}

	// Local signature
	localHasher := sha1.New()
	_, e = localHasher.Write(localBytes)
	if e != nil {
		return
	}
	localSha := base64.URLEncoding.EncodeToString(localHasher.Sum(nil))
	// fmt.Printf("Local SHA %s\n", localSha)

	// Remote signature
	remoteHasher := sha1.New()
	_, e = remoteHasher.Write(remoteBytes)
	if e != nil {
		return
	}
	remoteSha := base64.URLEncoding.EncodeToString(remoteHasher.Sum(nil))
	// fmt.Printf("Remote SHA %s\n", remoteSha)

	same = localSha == remoteSha
	return
}

// NewCompare creates a new Compare instance
func NewCompare(config *lib.Config, options lib.Options) (compare *Compare, e error) {

	compare = &Compare{
		Config:  config,
		Options: options,
	}

	return
}

type Compare struct {
	Config             *lib.Config              // Config is the config object
	LocalSQLPaths      []string                 // LocalSQLPaths is a list of paths pulled from the changesets.json file
	ChangesetSignature string                   // ChangesetSignature is a SHA signature for the changesets.json file
	LocalChangeFiles   []lib.ChangeFile         // LocalChangeFiles is a list of paths to local change files
	Files              *lib.Files               // Files is the injected file manager
	Connector          lib.IConnector           // IConnector is the injected server manager
	Databases          map[string]*lib.Database // A map of databases
	Options            lib.Options
}

func (c *Compare) initCommand() (server *lib.Server) {

	var e error
	server, e = c.Connector.Connect()

	if e != nil {
		panic(e)
	}

	e = c.Connector.UseDatabase(server, c.Config.Connection.DatabaseName)

	return

}

// ExportSchemaToSQL exports the current schema to sql
func (c *Compare) ExportSchemaToSQL() (sql string, e error) {

	schemaFile := c.Config.Connection.DatabaseName + ".schema.json"

	var localSchema *lib.Database
	emptySchema := &lib.Database{}

	if localSchema, e = lib.ReadSchemaFromFile(schemaFile); e != nil {
		return
	}

	sql, e = c.Connector.CreateChangeSQL(localSchema, emptySchema)
	return
}

// func (c *Compare) ExportDataToSQL(tableName string)
// allow for choosing individual tables
// func (c *Compare) ExportDataToJSON(tableName string)
// func (c *Compare) ImportDataFromJSON(tableName string)

// ApplyChangeset runs the sql produced by the `CompareSchema` command against the target database
// @command compare [reverse] apply
func (c *Compare) ApplyChangeset(changeset string) (e error) {

	server := c.initCommand()

	statements := strings.Split(changeset, ";")

	for _, s := range statements {
		sql := strings.Trim(strings.Trim(s, " "), "\n")
		if len(sql) == 0 {
			continue
		}
		lib.Debugf("Running sql: \n%s\n", c.Options, sql)

		_, e = server.Connection.Exec(sql)
		if e != nil {
			return
		}
	}

	return
}

// CompareSchema returns a string that contains a new line (`\n`) separated list of sql statements
// This comparison assumes the local `schemaFile` is the authority and the remote database is the
// schema to be updated
// @param reverse bool If true, the remote and local schema comparison is flipped in that the remote schema is treated as the authority
// 		and the local schema is treated as the schema to be updated.
// @command compare [reverse]
func (c *Compare) CompareSchema(schemaFile string) (sql string, e error) {

	var localSchema *lib.Database
	var remoteSchema *lib.Database

	localSchema, e = lib.ReadSchemaFromFile(schemaFile)
	if e != nil {
		return
	}

	server := c.initCommand()

	if remoteSchema, e = c.buildRemoteSchema(server); e != nil {
		return
	}

	sql = ""
	same := false

	if same, e = schemasAreSame(localSchema, remoteSchema); e != nil {
		return
	}

	if same {
		return
	}

	if c.Options&lib.OptReverse == lib.OptReverse {
		sql, e = c.Connector.CreateChangeSQL(remoteSchema, localSchema)
	} else {
		sql, e = c.Connector.CreateChangeSQL(localSchema, remoteSchema)
	}

	return
}

// ImportSchema calles `FetchSchema` and then marshal's it into a JSON object, writing it to the default schema.json file
// @command import
func (c *Compare) ImportSchema(fileName string) (e error) {

	server := c.initCommand()

	database := &lib.Database{
		Host: server.Host,
		Name: c.Config.Connection.DatabaseName,
	}
	database.Tables, e = c.Connector.FetchDatabaseTables(server, c.Config.Connection.DatabaseName)
	if e != nil {
		return
	}
	filePath := "./" + fileName
	dbBytes := []byte{}
	dbBytes, e = json.MarshalIndent(database, " ", "    ")
	e = ioutil.WriteFile(filePath, dbBytes, 0644)
	return
}

func (c *Compare) buildRemoteSchema(server *lib.Server) (remoteSchema *lib.Database, e error) {

	remoteSchema = &lib.Database{}

	remoteSchema.Host = server.Host
	remoteSchema.Name = c.Config.Connection.DatabaseName
	remoteSchema.Tables, e = c.Connector.FetchDatabaseTables(server, c.Config.Connection.DatabaseName)

	return
}
