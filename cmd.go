package main

import (
	"fmt"
	"github.com/macinnir/dvc/modules/compare"
	"io/ioutil"
	"os"

	"github.com/macinnir/dvc/connectors/mysql"
	"github.com/macinnir/dvc/connectors/sqlite"
	"github.com/macinnir/dvc/lib"
	"github.com/macinnir/dvc/modules/gen"
)

// Command is a type that represents the possible commands passed in at run time
type Command string

// Command Names
const (
	CommandImport        Command = "import"
	CommandExport        Command = "export"
	CommandGen           Command = "gen"
	CommandGenSchema     Command = "schema"
	CommandGenRepos      Command = "repos"
	CommandGenCaches     Command = "cache"
	CommandGenModels     Command = "models"
	CommandGenAll        Command = "all"
	CommandGenTypescript Command = "typescript"
	CommandCompare       Command = "compare"
	CommandHelp          Command = "help"
)

// Cmd is a container for handling commands
type Cmd struct {
	Options lib.Options

	// errLog  *log.Logger
	cmd    string
	Config *lib.Config
}

// Main is the main function that handles commands arguments
// and routes them to their correct options and functions
func (c *Cmd) Main(args []string) (err error) {

	args = args[1:]
	var cmd Command
	var arg string

	for len(args) > 0 {
		arg = args[0]
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
		default:
			lib.Debug(fmt.Sprintf("Arg: %s", arg), c.Options)
			if len(arg) > 0 && arg[0] == '-' {
				lib.Errorf("Unrecognized option '%s'. Try the --help option for more information\n", c.Options, arg)
				// c.errLog.Fatalf("Unrecognized option '%s'. Try the --help option for more information\n", arg)
				return nil
			}

			cmd = Command(arg)
		}

		args = args[1:]

		if len(cmd) > 0 {
			break
		}

	}

	if len(cmd) == 0 {
		lib.Error("No command provided", c.Options)
		return
	}
	lib.Debugf("cmd: %s, %v\n", c.Options, cmd, args)
	// fmt.Printf("cmd: %s, %v\n", cmd, args)

	c.dvc.Options = c.Options

	switch cmd {
	case CommandImport:
		c.CommandImport(args)
	case CommandExport:
		c.CommandExport(args)
	case CommandCompare:
		c.CommandCompare(args)
	case CommandGen:
		c.CommandGen(args)
	case CommandHelp:
		c.PrintHelp(args)
	}

	os.Exit(0)
	return
}

func (c *Cmd) initCompare() *compare.Compare {
	cmp, _ := compare.NewCompare(c.Config, c.Options)
	cmp.Connector, _ = connectorFactory(c.Config.DatabaseType, c.Config)
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

// CommandImport is the `import` command
func (c *Cmd) CommandImport(args []string) {

	var e error
	cmp := c.initCompare()

	if e = cmp.ImportSchema("./" + c.Config.Connection.DatabaseName + ".schema.json"); e != nil {
		lib.Error(e.Error(), c.Options)
		os.Exit(1)
	}

	lib.Infof("Schema `%s`.`%s` imported to %s.", c.Options, c.Config.Connection.Host, c.Config.Connection.DatabaseName, c.Config.Connection.DatabaseName+".schema.json")
}

func (c *Cmd) CommandExport(args []string) {
	var e error
	var sql string

	cmp := c.initCompare()

	if sql, e = cmp.ExportSchemaToSQL(); e != nil {
		lib.Error(e.Error(), c.Options)
		os.Exit(1)
	}

	fmt.Println(sql)
}

// CommandCompare handles the `compare` command
func (c *Cmd) CommandCompare(args []string) {

	var e error
	cmp := c.initCompare()

	cmd := "print"
	sql := ""
	schemaFile := c.Config.Connection.DatabaseName + ".schema.json"
	outfile := ""

	// lib.Debugf("Args: %v", c.Options, args)
	if len(args) == 0 {
		goto Main
	}

	for len(args) > 0 {

		switch args[0] {
		case "-r", "--reverse":
			c.Options |= lib.OptReverse
		case "-u", "--summary":
			c.Options |= lib.OptSummary
		case "print":
			cmd = "print"
		case "apply":
			cmd = "apply"
		default:

			if len(args[0]) > len("-o=") && args[0][0:len("-o=")] == "-o=" {
				outfile = args[0][len("-o="):]
				if len(outfile) == 0 {
					lib.Error("Outfile argument cannot be empty", c.Options)
					os.Exit(1)
				}
				cmd = "write"
			} else if args[0][0] == '-' {
				lib.Errorf("Unrecognized option '%s'. Try the --help option for more information\n", c.Options, args[0])
				os.Exit(1)
				// c.errLog.Fatalf("Unrecognized option '%s'. Try the --help option for more information\n", arg)
			}

			// Check if outfile argument is non-empty

			goto Main
		}
		args = args[1:]
	}

Main:

	// Do the comparison
	// TODO pass all options (e.g. verbose)
	if sql, e = cmp.CompareSchema(schemaFile); e != nil {
		lib.Error(e.Error(), c.Options)
		os.Exit(1)
	}

	if len(sql) == 0 {
		lib.Info("No changes found.", c.Options)
		os.Exit(0)
	}

	switch cmd {
	case "write":

		if len(args) == 0 {
			lib.Error("Missing file path for `dvc compare -o=[filePath]`", c.Options)
			os.Exit(1)
		}

		filePath := args[0]

		lib.Debugf("Writing changeset sql to path `%s`", c.Options, filePath)
		e = ioutil.WriteFile(filePath, []byte(sql), 0644)
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}
	case "apply":
		e = cmp.ApplyChangeset(sql)
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}

	case "print":
		// Print to stdout
		fmt.Printf("%s", sql)
	default:
		lib.Errorf("Unknown argument: `%s`", c.Options, cmd)
		os.Exit(1)
	}

}

// Write SQL to STDOUT

// CommandGen handles the `gen` command
func (c *Cmd) CommandGen(args []string) {

	var e error
	var database *lib.Database
	// fmt.Printf("Args: %v", args)
	if len(args) < 1 {
		lib.Error("Missing gen type [schema | model | repo]", c.Options)
		os.Exit(1)
	}
	subCmd := Command(args[0])

	if len(args) > 1 {
		args = args[1:]
	}

	for len(args) > 0 {
		switch args[0] {
		case "-c", "--clean":
			c.Options |= lib.OptClean
		}
		args = args[1:]
	}

	lib.Debugf("Gen Subcommand: %s", c.Options, subCmd)

	// Load the schema
	schemaFile := c.Config.Connection.DatabaseName + ".schema.json"
	database, e = lib.ReadSchemaFromFile(schemaFile)
	if e != nil {
		lib.Error(e.Error(), c.Options)
		os.Exit(1)
	}

	g := &gen.Gen{
		Options: c.Options,
		Config:  c.Config,
	}

	switch subCmd {
	case CommandGenSchema:
		e = g.GenerateGoSchemaFile(c.Config.Dirs.Schema, database)
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}

	case CommandGenRepos:

		fmt.Println("CommandGenRepos")
		e = g.GenerateGoRepoFiles(c.Config.Dirs.Repos, database)

		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}
	case CommandGenCaches:
		fmt.Println("CommandGenCaches")
		e = g.GenerateGoCacheFiles(c.Config.Dirs.Cache, database)

		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}
	case "repo":
		if len(args) < 4 {
			lib.Error("Missing repo name", c.Options)
			os.Exit(1)

		}

		var t *lib.Table

		if t, e = database.FindTableByName(args[3]); e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}

		if e = g.GenerateGoRepoFile(c.Config.Dirs.Repos, t); e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}

	case CommandGenModels:

		for _, table := range database.Tables {
			e = g.GenerateGoModelFile(c.Config.Dirs.Models, table)
			if e != nil {
				lib.Error(e.Error(), c.Options)
				os.Exit(1)
			}
		}

	case "model":
		if len(args) < 4 {
			lib.Error("missing model name", c.Options)
			os.Exit(1)
		}

		var t *lib.Table

		if t, e = database.FindTableByName(args[3]); e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}

		if e = g.GenerateGoModelFile(c.Config.Dirs.Models, t); e != nil {
			lib.Error(e.Error(), c.Options)
		}
	case CommandGenAll:

		// Schema
		lib.Debug("Generating schema...", c.Options)
		e = g.GenerateGoSchemaFile(c.Config.Dirs.Schema, database)
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}

		// Repos
		lib.Infof("Generating %d repos...", c.Options, len(database.Tables))
		e = g.GenerateGoRepoFiles(c.Config.Dirs.Repos, database)
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}

		// Models
		lib.Infof("Generating %d models", c.Options, len(database.Tables))
		for _, table := range database.Tables {
			e = g.GenerateGoModelFile(c.Config.Dirs.Models, table)
			if e != nil {
				lib.Error(e.Error(), c.Options)
				os.Exit(1)
			}
		}

		// Cache
		// lib.Info("Generating cache file", c.Options)
		// e = g.GenerateGoCacheFiles(c.dvc.Config.Dirs.Cache, database)

		// if e != nil {
		// 	lib.Error(e.Error(), c.Options)
		// 	os.Exit(1)
		// }

		// Typescript
		// lib.Info("Generating typescript types file", c.Options)
		// g.GenerateTypescriptTypesFile(c.dvc.Config.Dirs.Typescript, database)

	case "typescript":
		g.GenerateTypescriptTypesFile(c.Config.Dirs.Typescript, database)
	default:
		lib.Errorf("Unknown output type: `%s`", c.Options, subCmd)
		os.Exit(1)
	}
}

// PrintHelp prints help information
func (c *Cmd) PrintHelp(args []string) {
	help := `usage: dvc [OPTIONS] [COMMAND] [ARGS]
	
OPTIONS: 
	-h, --help 		Show help 
	-v, --verbose 	Show verbose logging 
	-vv, --debug 	Show very verbose (debug) logging
	-s, --silent 	Disable all output 

COMMANDS: 
	
	import	Build a schema definition file based on the target database. 
			This will overwrite any existing schema definition file.
	
	gen 	Generate go code based on the local schema json file. 
			Will fail if no imported schema file json file exists. 
			Requires one (and only one) of the following sub-commands: 

		models 		Generate models. 
		repos 		Generate repos 
		schema 		Generate go-dal schema bootstrap code based on imported schema information. 
		all 		Generate all above 
	
	compare [-r|--reverse] [ ( write <path> | apply ) ]
		
		Default behavior (no arguments) is to compare local schema as authority against 
		remote database as target and write the resulting sql to stdout.

		-r, --reverse 	Switches the roles of the schemas. The remote database becomes the authority 
						and the local schema the target for updating.

						Use this option when attempting to generate sql alter statements against a database that 
						matches the structure of your local schema, in order to make it match a database with the structure 
						of the remote.

		write		After performing the comparison, the resulting sequel statements will be written to a filepath <path> (required).

					Example: dts compare write path/to/changeset.sql
		
		apply 		After performing the comparison, apply the the resulting sql statements directly to the target database. 

					E.g. dts compare apply 
	
		import	[[path/to/local/schema.json]] 
				
				Generate a local schema json file based on the remote target database. 
				
				If no path is provided, the default path of ./[databaseName].json will be used. 
				This overwrites any existing json schema file. 

	`
	fmt.Printf(help)
}
