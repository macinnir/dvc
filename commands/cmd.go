package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/macinnir/dvc/modules/compare"

	"github.com/macinnir/dvc/connectors/mysql"
	"github.com/macinnir/dvc/connectors/sqlite"
	"github.com/macinnir/dvc/lib"
)

// Command is a type that represents the possible commands passed in at run time
type Command string

// TablesCache stores an md5 hash of the JSON representation of a table in the schema.json file
// These hashes are used to skip unchanged models for DAL and Model generation
type TablesCache struct {
	Dals   map[string]string
	Models map[string]string
}

// NewTablesCache is a factory method for TablesCache
func NewTablesCache() TablesCache {
	return TablesCache{
		Dals:   map[string]string{},
		Models: map[string]string{},
	}
}

// Command Names
const (
	CommandAdd           Command = "add"
	CommandClone         Command = "clone"
	CommandCompare       Command = "compare"
	CommandData          Command = "data"
	CommandDump          Command = "dump"
	CommandExport        Command = "export"
	CommandGen           Command = "gen"
	CommandGenAPI        Command = "api"
	CommandGenAPITests   Command = "apitests"
	CommandGenApp        Command = "app"
	CommandGenCLI        Command = "cli"
	CommandGenDal        Command = "dal"
	CommandGenDals       Command = "dals"
	CommandGenInterfaces Command = "interfaces"
	CommandGenModel      Command = "model"
	CommandGenModels     Command = "models"
	CommandGenRepos      Command = "repos"
	CommandGenRoutes     Command = "routes"
	CommandGenServices   Command = "services"
	CommandGenTests      Command = "tests"
	CommandGenTSPerms    Command = "tsperms"
	CommandGenTypescript Command = "typescript"
	CommandHelp          Command = "help"
	CommandImport        Command = "import"
	CommandInit          Command = "init"
	CommandInsert        Command = "insert"
	CommandInspect       Command = "inspect"
	CommandInstall       Command = "install"
	CommandLs            Command = "ls"
	CommandRefresh       Command = "refresh"
	CommandRm            Command = "rm"
	CommandSelect        Command = "select"
	CommandTest          Command = "test"
)

// Cmd is a container for handling commands
type Cmd struct {
	Options lib.Options

	// errLog  *log.Logger
	cmd            string
	Config         *lib.Config
	existingModels TablesCache
	newModels      map[string]string
}

func isCommand(cmd string) bool {
	commands := map[string]bool{
		"add":     true,
		"clone":   true,
		"compare": true,
		"data":    true,
		"dump":    true,
		"export":  true,
		"gen":     true,
		"help":    true,
		"import":  true,
		"init":    true,
		"insert":  true,
		"inspect": true,
		"install": true,
		"ls":      true,
		"refresh": true,
		"rm":      true,
		"select":  true,
		"test":    true,
	}

	_, ok := commands[cmd]
	return ok
}

// Run is the main function that handles commands arguments
// and routes them to their correct options and functions
func (c *Cmd) Run(inputArgs []string) (err error) {

	inputArgs = inputArgs[1:]
	var cmd Command
	var arg string
	args := []string{}

	for len(inputArgs) > 0 {

		arg = inputArgs[0]

		switch arg {
		case "-h", "--help":
			cmd = CommandHelp
		case "-v", "--verbose":
			c.Options &^= lib.OptLogDebug | lib.OptSilent
			c.Options |= lib.OptLogInfo
		case "-vv", "--debug":
			c.Options &^= lib.OptLogInfo | lib.OptSilent
			c.Options |= lib.OptLogDebug
		case "-s", "--silent":
			c.Options &^= lib.OptLogDebug | lib.OptLogInfo
			c.Options |= lib.OptSilent
		case "-f", "--force":
			c.Options |= lib.OptForce
		case "-c", "--clean":
			c.Options |= lib.OptClean
		default:

			lib.Debug(fmt.Sprintf("Arg: %s", arg), c.Options)

			if len(arg) > 0 && arg[0] == '-' {
				lib.Errorf("Unrecognized option '%s'. Try the --help option for more information\n", c.Options, arg)
				// c.errLog.Fatalf("Unrecognized option '%s'. Try the --help option for more information\n", arg)
				return nil
			}

			if isCommand(arg) && len(cmd) == 0 {
				cmd = Command(arg)
			} else {
				args = append(args, arg)
			}
		}

		inputArgs = inputArgs[1:]
	}

	if len(cmd) == 0 {
		lib.Error("No command provided", c.Options)
		return
	}
	lib.Debugf("cmd: %s, %v\n", c.Options, cmd, args)

	if cmd != CommandInit {
		var e error

		// Find the file config file
		maxDepth := 5
		curDepth := 0
		configPath := "dvc.toml"
		configFound := false

		for {
			if curDepth == maxDepth {
				break
			}

			curDepth++

			lib.Debugf("Looking for %s\n", c.Options, configPath)
			_, e = os.Stat(configPath)

			if os.IsNotExist(e) {
				configPath = "../" + configPath
			} else {
				configFound = true
				break
			}
		}

		if configFound {
			lib.Debugf("Found config at %s\n", c.Options, configPath)
		} else {
			lib.Errorf("Could not find config file (dvc.toml).", c.Options, configPath)
			return
		}

		configDir := path.Dir(configPath)

		if configDir != "." {

			lib.Debugf("Changing CWD to %s\n", c.Options, configDir)
			cwd := ""

			if e = os.Chdir(configDir); e != nil {
				lib.Errorf("Error (Chdir()): %s\n", c.Options, e.Error())
				return
			}

			if cwd, e = os.Getwd(); e != nil {
				lib.Errorf("Error (Getwd()): %s\n", c.Options, e.Error())
				return
			}

			lib.Debugf("CWD now %s\n", c.Options, cwd)
		}

		c.Config, e = loadConfigFromFile("./dvc.toml")
		if e != nil {
			lib.Error("Error loading config (./dvc.toml)", c.Options)
		}

	}

	switch cmd {
	case CommandRefresh:
		c.Refresh(args)
	case CommandClone:
		c.Clone(args)
	case CommandInsert:
		c.Insert(args)
	case CommandDump:
		c.Dump(args)
	case CommandSelect:
		c.CommandSelect(args)
	case CommandImport:
		c.Import(args)
	case CommandExport:
		c.Export(args)
	case CommandTest:
		c.Test(args)
	case CommandCompare:
		c.Compare(args)
	case CommandGen:
		c.Gen(args)
	case CommandHelp:
		helpCommandNames()
	case CommandInit:
		c.Init(args)
	case CommandLs:
		c.Ls(args)
	case CommandAdd:
		c.Add(args)
	case CommandInspect:
		c.CommandInspect(args)
	case CommandRm:
		c.Rm(args)
	case CommandData:
		c.Data(args)
	default:
		fmt.Printf("Invalid command `%s`\n", cmd)
		helpCommandNames()
	}

	os.Exit(0)
	return
}

// SearchType is the type of search
type SearchType int

const (
	// SearchTypeWildcard is a wildcard
	SearchTypeWildcard SearchType = iota
	// SearchTypeStartingWildcard is a wildcard at the start (e.g. `*foo`)
	SearchTypeStartingWildcard
	// SearchTypeEndingWildcard is a wildcard at the end (e.g. `foo*`)
	SearchTypeEndingWildcard
	// SearchTypeBoth is a wildcard on both the start and the end (e.g. `*foo*`)
	SearchTypeBoth
)

func writeSQLToLog(sql string) {

	lib.EnsureDir("meta")

	sqlLog := time.Now().Format("20060102150405") + "\n"
	sqlLog += sql

	filePath := "meta/dvc-changes.log"

	if _, e := os.Stat(filePath); os.IsNotExist(e) {
		ioutil.WriteFile(filePath, []byte(sqlLog), 0600)
	} else {
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if _, err = f.WriteString(sqlLog); err != nil {
			panic(err)
		}
	}
}

// loadDatabase loads a database
func (c *Cmd) loadDatabase() *lib.Database {
	// Load the schema
	schemaFile := c.Config.Connection.DatabaseName + ".schema.json"
	database, e := lib.ReadSchemaFromFile(schemaFile)
	if e != nil {
		lib.Error(e.Error(), c.Options)
		os.Exit(1)
	}

	return database
}

func (c *Cmd) genTableCache(database *lib.Database) {

	var e error

	lib.EnsureDir("meta")

	tableCachePath := "meta/.tables"

	c.existingModels = NewTablesCache()

	if _, e := os.Stat(tableCachePath); e == nil {
		if fileBytes, e := ioutil.ReadFile(tableCachePath); e == nil {
			if e = json.Unmarshal(fileBytes, &c.existingModels); e != nil {
				panic("Can't read table cache")
			}
		}
	}

	c.newModels = map[string]string{}
	for _, table := range database.Tables {

		marshalledTable := []byte{}

		if marshalledTable, e = json.Marshal(table); e != nil {
			panic(e.Error())
		}

		// Build the list of new model hashes to check against
		c.newModels[table.Name] = lib.HashStringMd5(string(marshalledTable))
	}
}

func (c *Cmd) saveTableCache() {

	lib.EnsureDir("meta")

	var e error
	newModelData := []byte{}
	if newModelData, e = json.Marshal(c.existingModels); e != nil {
		panic(fmt.Sprintf("Can't serialize newModels: %s", e.Error()))
	}
	if e = ioutil.WriteFile("meta/.tables", newModelData, 0777); e != nil {
		panic("Could not write table cache to file")
	}
}

// loadConfigFromFile loads a config file
func loadConfigFromFile(configFilePath string) (config *lib.Config, e error) {

	// fmt.Printf("Looking for config at path %s\n", configFilePath)
	if _, e = os.Stat(configFilePath); os.IsNotExist(e) {
		e = fmt.Errorf("Config file `%s` not found", configFilePath)
		return
	}

	config = &lib.Config{
		OneToMany: map[string]string{},
		ManyToOne: map[string]string{},
	}
	_, e = toml.DecodeFile(configFilePath, config)

	if e != nil {
		return
	}

	if len(config.OneToMany) > 0 {

		for col, subCol := range config.OneToMany {
			config.ManyToOne[subCol] = col
		}
	}

	return
}

func (c *Cmd) initCompare() *compare.Compare {
	connector, _ := connectorFactory(c.Config.DatabaseType, c.Config)
	cmp, _ := compare.NewCompare(c.Config, c.Options, connector)
	return cmp
}

func connectorFactory(databaseType string, config *lib.Config) (connector lib.IConnector, e error) {

	t := lib.DatabaseType(databaseType)

	switch t {
	case lib.DatabaseTypeMysql:
		connector = &mysql.MySQL{
			Config: config,
		}
	case lib.DatabaseTypeSqlite:
		connector = &sqlite.Sqlite{
			Config: config,
		}
	default:
		e = errors.New("Invalid database type")
	}

	return
}
