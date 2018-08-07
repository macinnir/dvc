package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/macinnir/dvc/gen"
	"github.com/macinnir/dvc/logger"
	"github.com/macinnir/dvc/types"
)

// Command is a type that represents the possible commands passed in at run time
type Command string

// Command Names
const (
	CommandImport        Command = "import"
	CommandGen           Command = "gen"
	CommandGenSchema     Command = "schema"
	CommandGenRepos      Command = "repos"
	CommandGenModels     Command = "models"
	CommandGenAll        Command = "all"
	CommandGenTypescript Command = "typescript"
	CommandCompare       Command = "compare"
	CommandHelp          Command = "help"
)

// Cmd is a container for handling commands
type Cmd struct {
	Options types.Options
	// errLog  *log.Logger
	cmd    string
	dvc    *DVC
	config *types.Config
}

// Main is the main function that handles commands arguments and routes them to their correct options and functions
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
			c.Options &^= types.OptLogDebug | types.OptSilent
			c.Options |= types.OptLogInfo
		case "-vv", "--debug":
			c.Options &^= types.OptLogInfo | types.OptSilent
			c.Options |= types.OptLogDebug
		case "-s", "--silent":
			c.Options &^= types.OptLogDebug | types.OptLogInfo
			c.Options |= types.OptSilent
		default:
			logger.Debug(fmt.Sprintf("Arg: %s", arg), c.Options)
			if len(arg) > 0 && arg[0] == '-' {
				logger.Errorf("Unrecognized option '%s'. Try the --help option for more information\n", c.Options, arg)
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
		logger.Error("No command provided", c.Options)
		return
	}
	logger.Debugf("cmd: %s, %v\n", c.Options, cmd, args)
	// fmt.Printf("cmd: %s, %v\n", cmd, args)

	switch cmd {
	case CommandImport:
		c.CommandImport(args)
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

// CommandImport is the `import` command
func (c *Cmd) CommandImport(args []string) {

	var e error
	schemaFile := c.dvc.Config.DatabaseName + ".schema.json"

	if e = c.dvc.ImportSchema(schemaFile); e != nil {
		logger.Error(e.Error(), c.Options)
		os.Exit(1)
	}

	logger.Infof("Schema `%s`.`%s` imported to %s.", c.Options, c.dvc.Config.Host, c.dvc.Config.DatabaseName, schemaFile)
}

// CommandCompare handles the `compare` command
func (c *Cmd) CommandCompare(args []string) {

	var e error

	cmd := "print"
	sql := ""
	schemaFile := c.dvc.Config.DatabaseName + ".schema.json"
	outfile := ""

	// logger.Debugf("Args: %v", c.Options, args)
	if len(args) == 0 {
		goto Main
	}

	for len(args) > 0 {

		switch args[0] {
		case "-r", "--reverse":
			c.Options |= types.OptReverse
		case "-u", "--summary":
			c.Options |= types.OptSummary
		case "apply":
			cmd = "apply"
		default:

			if len(args[0]) > len("-o=") && args[0][0:len("-o=")] == "-o=" {
				outfile = args[0][len("-o="):]
				if len(outfile) == 0 {
					logger.Error("Outfile argument cannot be empty", c.Options)
					os.Exit(1)
				}
				cmd = "write"
			} else if args[0][0] == '-' {
				logger.Errorf("Unrecognized option '%s'. Try the --help option for more information\n", c.Options, args[0])
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
	if sql, e = c.dvc.CompareSchema(schemaFile, c.Options); e != nil {
		logger.Error(e.Error(), c.Options)
		os.Exit(1)
	}

	if len(sql) == 0 {
		logger.Info("No changes found.", c.Options)
		os.Exit(0)
	}

	switch cmd {
	case "write":

		if len(args) == 0 {
			logger.Error("Missing file path for `dvc compare -o=[filePath]`", c.Options)
			os.Exit(1)
		}

		filePath := args[0]

		logger.Debugf("Writing changeset sql to path `%s`", c.Options, filePath)
		e = ioutil.WriteFile(filePath, []byte(sql), 0644)
		if e != nil {
			logger.Error(e.Error(), c.Options)
			os.Exit(1)
		}
	case "apply":
		e = c.dvc.ApplyChangeset(sql)
		if e != nil {
			logger.Error(e.Error(), c.Options)
			os.Exit(1)
		}

	case "print":
		// Print to stdout
		fmt.Printf("%s", sql)
	default:
		logger.Errorf("Unknown argument: `%s`", c.Options, cmd)
		os.Exit(1)
	}

}

// Write SQL to STDOUT

// CommandGen handles the `gen` command
func (c *Cmd) CommandGen(args []string) {

	var e error
	var database *types.Database
	fmt.Printf("Args: %v", args)
	if len(args) < 1 {
		logger.Error("Missing gen type [schema | model | repo]", c.Options)
		os.Exit(1)
	}
	subCmd := Command(args[0])

	if len(args) > 1 {
		args = args[1:]
	}

	for len(args) > 0 {
		switch args[0] {
		case "-c", "--clean":
			c.Options |= types.OptClean
		}
		args = args[1:]
	}

	logger.Debugf("Gen Subcommand: %s", c.Options, subCmd)

	// Load the schema
	schemaFile := c.dvc.Config.DatabaseName + ".schema.json"
	database, e = ReadSchemaFromFile(schemaFile)
	if e != nil {
		logger.Error(e.Error(), c.Options)
		os.Exit(1)
	}

	g := &gen.Gen{
		Options: c.Options,
	}

	switch subCmd {
	case CommandGenSchema:
		e = g.GenerateGoSchemaFile(database)
		if e != nil {
			logger.Error(e.Error(), c.Options)
			os.Exit(1)
		}

	case CommandGenRepos:

		e = g.GenerateGoRepoFiles(database)

		if e != nil {
			logger.Error(e.Error(), c.Options)
			os.Exit(1)
		}

	case "repo":
		if len(args) < 4 {
			logger.Error("Missing repo name", c.Options)
			os.Exit(1)

		}

		var t *types.Table

		if t, e = database.FindTableByName(args[3]); e != nil {
			logger.Error(e.Error(), c.Options)
			os.Exit(1)
		}

		if e = g.GenerateGoRepoFile(t); e != nil {
			logger.Error(e.Error(), c.Options)
			os.Exit(1)
		}

	case CommandGenModels:

		for _, table := range database.Tables {
			e = g.GenerateGoModelFile(table)
			if e != nil {
				logger.Error(e.Error(), c.Options)
				os.Exit(1)
			}
		}

	case "model":
		if len(args) < 4 {
			logger.Error("missing model name", c.Options)
			os.Exit(1)
		}

		var t *types.Table

		if t, e = database.FindTableByName(args[3]); e != nil {
			logger.Error(e.Error(), c.Options)
			os.Exit(1)
		}

		if e = g.GenerateGoModelFile(t); e != nil {
			logger.Error(e.Error(), c.Options)
		}
	case CommandGenAll:

		// Generate schema
		logger.Debug("Generating schema...", c.Options)
		e = g.GenerateGoSchemaFile(database)
		if e != nil {
			logger.Error(e.Error(), c.Options)
			os.Exit(1)
		}
		logger.Debug("done\n", c.Options)

		// Generate repos
		logger.Debugf("Generating %d repos...", c.Options, len(database.Tables))
		e = g.GenerateGoRepoFiles(database)
		if e != nil {
			logger.Error(e.Error(), c.Options)
			os.Exit(1)
		}
		logger.Debug("done", c.Options)

		// Generate models
		logger.Debugf("Generating %d models...", c.Options, len(database.Tables))
		for _, table := range database.Tables {
			e = g.GenerateGoModelFile(table)
			if e != nil {
				logger.Error(e.Error(), c.Options)
				os.Exit(1)
			}
		}
		logger.Debug("done\n", c.Options)

	case "typescript":
		g.GenerateTypescriptTypesFile(database)
	default:
		logger.Errorf("Unknown output type: `%s`", c.Options, subCmd)
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
