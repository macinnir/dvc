package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	// "strconv"
	// "fmt"
	"crypto/sha1"
	"encoding/base64"
	"github.com/macinnir/dvc/query"
	"github.com/macinnir/dvc/types"
	// "io/ioutil"
	// "log"
	// "path"
	// "sort"
)

// Command Names
const (
	CommandImport        = "import"
	CommandGen           = "gen"
	CommandGenSchema     = "schema"
	CommandGenRepos      = "repos"
	CommandGenModels     = "models"
	CommandGenAll        = "all"
	CommandGenTypescript = "typescript"
	CommandCompare       = "compare"
	CommandHelp          = "help"
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
		Config: config,
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
	serverService      *serverService     // ServerService is the injected server manager
	Databases          map[string]*types.Database
}

func (d *DVC) initCommand() (server *types.Server) {

	var e error
	server, e = d.serverService.ConnectToServer(d.Config.Host, d.Config.Username, d.Config.Password)

	if e != nil {
		panic(e)
	}

	d.serverService.UseDatabase(server, d.Config.DatabaseName)

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
	database.Tables, e = d.serverService.FetchDatabaseTables(server, d.Config.DatabaseName)
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
func (d *DVC) CompareSchema(schemaFile string, reverse bool) (sql string, e error) {

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
	remoteSchema.Tables, e = d.serverService.FetchDatabaseTables(server, d.Config.DatabaseName)

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

	if reverse == true {
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

// Run prints help documentation to stdout
func (d *DVC) Run() (e error) {

	fmt.Println("Running Run()")

	fmt.Println("1. Fetching local changeset list")
	if d.LocalSQLPaths, e = d.Files.FetchLocalChangesetList(d.Config.ChangeSetPath); e != nil {
		return
	}

	for _, p := range d.LocalSQLPaths {
		fmt.Printf("\tPath: %s\n", p)
	}

	// fmt.Println("2. Building change files")
	// if d.LocalChangeFiles, e = d.Files.BuildChangeFiles(d.LocalSQLPaths); e != nil {
	// 	return
	// }

	// fmt.Println("3. Building changeset signature")
	// if d.ChangesetSignature, e = HashFileMd5(d.Config.ChangeSetPath + "/changesets.json"); e != nil {
	// 	return
	// }

	return
}

func (d *DVC) verifyChangesetFile() (e error) {
	// fmt.Printf("Looking for changeset file at path %s\n", d.Config.ChangeSetPath)
	// if _, e = os.Stat(d.Config.ChangeSetPath); os.IsNotExist(e) {
	// 	e = errors.New("changeset path does not exist")
	// 	return
	// }

	dt := ""
	for _, t := range DatabaseTypes {
		if d.Config.DatabaseType == t {
			dt = d.Config.DatabaseType
			break
		}
	}

	if dt == "" {
		e = errors.New("invalid database type")
		return
	}

	return
}

// // Run runs the dvc
// func (d *DVC) Run() {

// 	var e error
// 	var server Server
// 	var localChangeLogs LocalChangeLogs

// 	// 0. Look for changelog json file

// 	// 1. LocalChangeLogs
// 	// Find the local changeset files
// 	log.Println("1. LocalChangeLogs...")
// 	localChangeLogs, e = NewLocalChangeLogs(changesetsPath)
// 	if e != nil {
// 		log.Fatal("Local changelog error ", e)
// 	}
// 	e = localChangeLogs.FetchLocalChangesetFiles()
// 	if e != nil {
// 		log.Fatal("FetchLocalChangesetFiles", e)
// 	}
// 	log.Printf("\t Found %d changesets\n", len(localChangeLogs.ChangeSets))

// 	// for changesetName, changeset := range localChangeLogs.ChangeSets {
// 	// 	log.Printf("\t Changeset Dir: %s\n", changesetName)

// 	// 	for _, k := range changeset.Files {
// 	// 		log.Printf("\t\t Change File: %s\n", k)
// 	// 	}
// 	// }

// 	// 2. Server
// 	log.Println("2. Server")
// 	server = NewServer(dbHost, dbUser, dbPass)
// 	server.FetchDatabases()
// 	log.Printf("\t Found %d databases\n", len(server.Databases))

// 	// 3. Database
// 	log.Printf("3. Database `%s`\n", dbName)
// 	database, ok := server.Databases[dbName]
// 	if !ok {
// 		log.Printf("Database `%s` does not exist. Trying to create...", dbName)
// 		database, e = server.CreateDatabase(dbName)
// 		if e != nil {
// 			log.Fatal(e)
// 		}
// 	}

// 	e = database.Connect(dbUser, dbPass)
// 	if e != nil {
// 		log.Fatal(e)
// 	}

// 	// e = database.DropBaseTables()
// 	// if e != nil {
// 	// 	log.Fatal(e)
// 	// }

// 	e = database.FetchTables()
// 	if e != nil {
// 		database.logFatal("FetchTables(): " + e.Error())
// 	}

// 	log.Println("3.1 Database: Create Version Tables...")
// 	e = database.CreateBaseTablesIfNotExists()
// 	if e != nil {
// 		database.logFatal("CreateBaseTablesIfNotExists " + e.Error())
// 	}

// 	e = database.Start()
// 	log.Printf("Starting run #%d", database.runID)
// 	if e != nil {
// 		database.logFatal("Start()" + e.Error())
// 	}

// 	// 4. Database Changesets
// 	log.Printf("3.2 Fetching changes in `%s`\n", dbName)
// 	e = database.FetchSets()
// 	if e != nil {
// 		database.logFatal("FetchSets()" + e.Error())
// 	}

// 	// Resolve
// 	log.Println("3.3 Importing files...")

// 	// Delete sets not found
// 	for _, changesetName := range database.sortedSetKeys {

// 		dbChangeset := database.sets[changesetName]

// 		if _, ok := localChangeLogs.ChangeSets[changesetName]; !ok {
// 			database.DeleteSet(changesetName)
// 			continue
// 		}

// 		localChangeSet := localChangeLogs.ChangeSets[changesetName]

// 		for _, fileName := range dbChangeset.SortedFileKeys {

// 			if len(fileName) == 0 {
// 				continue
// 			}

// 			// fmt.Printf("Looking for %s\n", fileName)
// 			dbFile, _ := dbChangeset.Files[fileName]

// 			found := false

// 			for _, localFileName := range localChangeSet.Files {
// 				if localFileName == fileName {
// 					// found
// 					found = true
// 					break
// 				}
// 			}
// 			if found == false {
// 				database.log(fmt.Sprintf("File %s not found. Deleting...", fileName))
// 				database.DeleteFile(dbChangeset, dbFile)
// 			}
// 		}
// 	}

// 	for _, changesetName := range localChangeLogs.SortedChangesetKeys {

// 		changeset := localChangeLogs.ChangeSets[changesetName]

// 		if !database.SetExists(changesetName) {
// 			// log.Printf("Set Not Exist: %s", changesetName)
// 			if e = database.CreateSet(changesetName); e != nil {
// 				database.logFatal("CreateSet() " + e.Error())
// 			}
// 		}

// 		dbChangeset := database.sets[changesetName]

// 		for _, fileName := range changeset.Files {

// 			// fmt.Printf("File: %s\nOrdinal: %s\nAction: %s\nTarget: %s\n", fileName, fileOrdinal, fileAction, fileTarget)

// 			filePath := path.Join(changesetsPath, changesetName, fileName)
// 			if !database.FileExists(changesetName, fileName) {
// 				hash, _ := HashFileMd5(filePath)
// 				// log.Printf("File Not Exists: %s/%s", changesetName, fileName)

// 				var contentBytes []byte
// 				contentBytes, e = ioutil.ReadFile(filePath)
// 				if e != nil {
// 					database.logFatal("ReadFile() " + filePath + " " + e.Error())
// 				}

// 				content := string(contentBytes)

// 				_, e = database.CreateFile(dbChangeset, fileName, hash, content)
// 				if e != nil {
// 					database.logFatal("CreateFile() " + e.Error())
// 				}

// 			}
// 		}
// 	}

// 	log.Println("3.4 Applying database changes...")

// 	notRunChangeFiles := map[string]*DVCFile{}
// 	notRunChangeFileNames := []string{}

// 	for _, changesetName := range database.sortedSetKeys {

// 		for _, fileName := range database.sets[changesetName].SortedFileKeys {
// 			f := database.sets[changesetName].Files[fileName]
// 			if !f.IsRun {
// 				notRunChangeFileNames = append(notRunChangeFileNames, f.FullPath)
// 				notRunChangeFiles[f.FullPath] = f
// 			}
// 		}
// 	}

// 	sort.Strings(notRunChangeFileNames)

// 	database.log(fmt.Sprintf("Found %d changes to run \n", len(notRunChangeFiles)))

// 	for _, notRunChangeFileName := range notRunChangeFileNames {
// 		// fmt.Printf("\t\t %s\n", notRunChangeFileName)
// 		fileBytes, readFileErr := ioutil.ReadFile(path.Join(changesetsPath, notRunChangeFileName))
// 		if readFileErr != nil {
// 			log.Fatal(readFileErr)
// 		}
// 		fileString := string(fileBytes)
// 		e = database.RunChange(fileString)
// 		if e != nil {
// 			database.logFatal("RunChange() File: " + notRunChangeFileName + " - " + e.Error())
// 		}

// 		database.log(notRunChangeFiles[notRunChangeFileName].ToLogString())

// 		e = database.SetFileAsRun(notRunChangeFiles[notRunChangeFileName].ID)
// 		if e != nil {
// 			log.Fatal(e)
// 		}
// 	}
// 	database.finish()
// }
