package lib

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

// NewDVC creates a new DVC instance
// Can be called 2 ways:
// 	1. NewDvc(filePath)
//  2. NewDvc(host, databaseName, username, password, changesetPath, databaseType)
func NewDVC(config *Config) (dvc *DVC, e error) {

	dvc = &DVC{
		Config: config,
	}

	return
}

// DVC is the core object for running Database Version Control
type DVC struct {
	Config             *Config              // Config is the config object
	LocalSQLPaths      []string             // LocalSQLPaths is a list of paths pulled from the changesets.json file
	ChangesetSignature string               // ChangesetSignature is a SHA signature for the changesets.json file
	LocalChangeFiles   []ChangeFile         // LocalChangeFiles is a list of paths to local change files
	Files              *Files               // Files is the injected file manager
	Connector          IConnector           // IConnector is the injected server manager
	Databases          map[string]*Database // A map of databases
}

func (d *DVC) initCommand() (server *Server) {

	var e error
	server, e = d.Connector.Connect()

	if e != nil {
		panic(e)
	}

	e = d.Connector.UseDatabase(server, d.Config.Connection.DatabaseName)

	return

}

// ImportSchema calles `FetchSchema` and then marshal's it into a JSON object, writing it to the default schema.json file
// @command import
func (d *DVC) ImportSchema(fileName string) (e error) {

	server := d.initCommand()

	database := &Database{
		Host: server.Host,
		Name: d.Config.Connection.DatabaseName,
	}
	database.Tables, e = d.Connector.FetchDatabaseTables(server, d.Config.Connection.DatabaseName)
	if e != nil {
		return
	}
	filePath := "./" + fileName
	dbBytes := []byte{}
	dbBytes, e = json.MarshalIndent(database, " ", "    ")
	e = ioutil.WriteFile(filePath, dbBytes, 0644)
	return
}

// CompareSchema returns a string that contains a new line (`\n`) separated list of sql statements
// This comparison assumes the local `schemaFile` is the authority and the remote database is the
// schema to be updated
// @param reverse bool If true, the remote and local schema comparison is flipped in that the remote schema is treated as the authority
// 		and the local schema is treated as the schema to be updated.
// @command compare [reverse]
func (d *DVC) CompareSchema(schemaFile string, options Options) (sql string, e error) {

	var localSchema *Database
	var remoteSchema *Database

	localSchema, e = ReadSchemaFromFile(schemaFile)
	if e != nil {
		return
	}

	server := d.initCommand()

	if remoteSchema, e = d.buildRemoteSchema(server); e != nil {
		return
	}

	sql = ""
	same := false

	if same, e = d.schemasAreSame(localSchema, remoteSchema); e != nil {
		return
	}

	if same {
		return
	}

	if options&OptReverse == OptReverse {
		sql, e = d.Connector.CreateChangeSQL(remoteSchema, localSchema)
	} else {
		sql, e = d.Connector.CreateChangeSQL(localSchema, remoteSchema)
	}

	return
}

func (d *DVC) buildRemoteSchema(server *Server) (remoteSchema *Database, e error) {

	remoteSchema = &Database{}

	remoteSchema.Host = server.Host
	remoteSchema.Name = d.Config.Connection.DatabaseName
	remoteSchema.Tables, e = d.Connector.FetchDatabaseTables(server, d.Config.Connection.DatabaseName)

	return
}

func (d *DVC) schemasAreSame(localSchema *Database, remoteSchema *Database) (same bool, e error) {

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

func (d *DVC) ExportSchemaToSQL(options Options) (sql string, e error) {

	schemaFile := d.Config.Connection.DatabaseName + ".schema.json"

	var localSchema *Database
	emptySchema := &Database{}

	if localSchema, e = ReadSchemaFromFile(schemaFile); e != nil {
		return
	}

	sql, e = d.Connector.CreateChangeSQL(localSchema, emptySchema)
	return
}

// ApplyChangeset runs the sql produced by the `CompareSchema` command against the target database
// @command compare [reverse] apply
func (d *DVC) ApplyChangeset(changeset string) (e error) {

	server := d.initCommand()

	statements := strings.Split(changeset, ";")

	for _, s := range statements {
		sql := strings.Trim(strings.Trim(s, " "), "\n")
		if len(sql) == 0 {
			continue
		}
		fmt.Printf("Running sql: \n%s\n", sql)

		_, e = server.Connection.Exec(sql)
		if e != nil {
			return
		}
	}

	return
}
