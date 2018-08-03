package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/macinnir/dvc/connectors/mysql"
	"github.com/macinnir/dvc/query"
	"github.com/macinnir/dvc/types"
)

// NewDVC creates a new DVC instance
// Can be called 2 ways:
// 	1. NewDvc(filePath)
//  2. NewDvc(host, databaseName, username, password, changesetPath, databaseType)
func NewDVC(args ...string) (dvc *DVC, e error) {

	var config *types.Config
	var configFilePath string

	if len(args) == 1 {

		// load config from file path
		configFilePath = args[0]
		config, e = loadConfigFromFile(configFilePath)

		if e != nil {
			return
		}

	} else {

		if len(args) < 6 {
			e = errors.New("not enough arguments")
			return
		}

		// build config from arguments
		config = &types.Config{
			Host:          args[0],
			DatabaseName:  args[1],
			Username:      args[2],
			Password:      args[3],
			ChangeSetPath: args[4],
			DatabaseType:  args[5],
		}
	}

	dvc = &DVC{
		Config:    config,
		connector: &mysql.MySQL{},
	}

	return
}

// DVC is the core object for running Database Version Control
type DVC struct {
	Config             *types.Config      // Config is the config object
	LocalSQLPaths      []string           // LocalSQLPaths is a list of paths pulled from the changesets.json file
	ChangesetSignature string             // ChangesetSignature is a SHA signature for the changesets.json file
	LocalChangeFiles   []types.ChangeFile // LocalChangeFiles is a list of paths to local change files
	Files              *Files             // Files is the injected file manager
	connector          types.IConnector   // IConnector is the injected server manager
	Databases          map[string]*types.Database
}

func (d *DVC) initCommand() (server *types.Server) {

	var e error
	server, e = d.connector.ConnectToServer(d.Config.Host, d.Config.Username, d.Config.Password)

	if e != nil {
		panic(e)
	}

	e = d.connector.UseDatabase(server, d.Config.DatabaseName)

	return

}

// ImportSchema calles `FetchSchema` and then marshal's it into a JSON object, writing it to the default schema.json file
// @command import
func (d *DVC) ImportSchema(fileName string) (e error) {

	server := d.initCommand()

	database := &types.Database{
		Host: server.Host,
		Name: d.Config.DatabaseName,
	}
	database.Tables, e = d.connector.FetchDatabaseTables(server, d.Config.DatabaseName)
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
func (d *DVC) CompareSchema(schemaFile string, options types.Options) (sql string, e error) {

	var localSchema *types.Database
	var remoteSchema *types.Database

	localSchema, e = ReadSchemaFromFile(schemaFile)
	if e != nil {
		return
	}

	server := d.initCommand()

	remoteSchema = &types.Database{}

	remoteSchema.Host = server.Host
	remoteSchema.Name = d.Config.DatabaseName
	remoteSchema.Tables, e = d.connector.FetchDatabaseTables(server, d.Config.DatabaseName)

	if e != nil {
		return
	}

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

	if localSha == remoteSha {
		// fmt.Println("They are the same...")
		return
	}

	sql = ""

	query := &query.Query{}

	if options&types.OptReverse == types.OptReverse {
		sql, e = query.CreateChangeSQL(remoteSchema, localSchema)
	} else {
		sql, e = query.CreateChangeSQL(localSchema, remoteSchema)
	}

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
		fmt.Printf("Running sql: %s", sql)

		server.Connection.Exec(sql)
		if e != nil {
			return
		}
	}

	return
}
