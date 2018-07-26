package main

import (
	"fmt"
	"github.com/macinnir/dvc/gen"
	"github.com/macinnir/dvc/types"
	"io/ioutil"
	"log"
	"os"
)

// Cmd is a container for handling commands
type Cmd struct {
	Options types.Options
	errLog  *log.Logger
	cmd     string
	dvc     *DVC
	config  *types.Config
}

// NewCmd returns a new Cmd instance
func NewCmd() *Cmd {

	var e error
	var dvc *DVC

	if dvc, e = NewDVC(configFilePath); e != nil {
		fatal(e.Error())
	}

	return &Cmd{
		errLog: log.New(os.Stderr, "", 0),
		cmd:    "",
		dvc:    dvc,
	}
}

// Main is the main function that handles commands arguments and routes them to their correct options and functions
func (c *Cmd) Main(args []string) (err error) {

	cmd := ""

	for len(args) > 0 {

		switch args[0] {
		case "-h", "--help":

		case "-v", "--verbose":
			c.Options &^= OptVerboseDebug | OptSilent
			c.Options |= OptVerbose
		case "-s", "--silent":
			c.Options &^= OptVerboseDebug | OptVerbose
			c.Options |= OptSilent
		case "-vv", "--debug":
			c.Options &^= OptVerbose | OptSilent
			c.Options |= OptVerboseDebug
		default:
			arg := args[0]
			if len(arg) > 0 && arg[0] == '-' {
				c.errLog.Fatalf("Unrecognized option '%s'. Try the --help option for more information\n", arg)
				return nil
			}

			cmd = arg
			break

		}

		args = args[1:]
	}

	if len(cmd) == 0 {
		fmt.Printf("No command provided")
		return
	}

	switch cmd {
	case CommandImport:
		c.CommandImport(args)
	case CommandCompare:
		c.CommandCompare(args)
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
		fatal(e.Error())
	}

	log.Printf("Schema `%s`.`%s` imported to %s.", c.dvc.Config.Host, c.dvc.Config.DatabaseName, schemaFile)
}

// CommandCompare handles the `compare` command
func (c *Cmd) CommandCompare(args []string) {

	var e error
	sql := ""

	nextArgIdx := 1

	schemaFile := c.dvc.Config.DatabaseName + ".schema.json"

	isCompareFlipped := false
	if len(args) > 1 && args[nextArgIdx] == "reverse" {
		isCompareFlipped = true
		nextArgIdx++
	}

	if sql, e = c.dvc.CompareSchema(schemaFile, isCompareFlipped); e != nil {
		fatal(e.Error())
	}

	if len(sql) == 0 {
		log.Printf("No changes found.")
		os.Exit(0)
	}

	if len(args) > nextArgIdx {
		// "--file="

		switch args[nextArgIdx] {
		case "write":

			nextArgIdx++

			if len(args) > nextArgIdx {

				filePath := args[nextArgIdx]

				log.Printf("Writing sql to path `%s`", filePath)
				e = ioutil.WriteFile(filePath, []byte(sql), 0644)
				if e != nil {
					fatal(e.Error())
				}
			} else {
				fatal("Missing path to write sql to")
			}
		case "apply":
			e = c.dvc.ApplyChangeset(sql)
			if e != nil {
				fatal(e.Error())
			}

		default:
			fatal(fmt.Sprintf("Unknown argument: `%s`\n", args[2]))
		}

	}

	// Write SQL to STDOUT
	fmt.Printf("%s", sql)
}

// CommandGen handles the `gen` command
func (c *Cmd) CommandGen(args []string) {

	var e error
	var database *types.Database

	if len(args) < 3 {
		fatal("Missing gen type [schema | model | repo]")
	}
	subCmd := args[2]

	// Load the schema
	schemaFile := c.dvc.Config.DatabaseName + ".schema.json"
	database, e = ReadSchemaFromFile(schemaFile)
	if e != nil {
		fatal(e.Error())
	}

	switch subCmd {
	case CommandGenSchema:
		e = gen.GenerateGoSchemaFile(database)
		if e != nil {
			fatal(e.Error())
		}

	case CommandGenRepos:

		for _, table := range database.Tables {
			e = gen.GenerateGoRepoFile(table)
			if e != nil {
				fatal(e.Error())
			}
		}

		gen.GenerateReposBootstrapFile(database)

	case "repo":
		if len(args) < 4 {
			fatal("Missing repo name")
		}

		var t *types.Table

		if t, e = database.FindTableByName(args[3]); e != nil {
			fatal(e.Error())
		}

		if e = gen.GenerateGoRepoFile(t); e != nil {
			fatal(e.Error())
		}

	case CommandGenModels:

		for _, table := range database.Tables {
			e = gen.GenerateGoModelFile(table)
			if e != nil {
				fatal(e.Error())
			}
		}

	case "model":
		if len(args) < 4 {
			fatal("missing model name")
		}

		var t *types.Table

		if t, e = database.FindTableByName(args[3]); e != nil {
			fatal(e.Error())
		}

		if e = gen.GenerateGoModelFile(t); e != nil {
			fatal(e.Error())
		}
	case CommandGenAll:

		// Generate schema
		fmt.Print("Generating schema...")
		e = gen.GenerateGoSchemaFile(database)
		if e != nil {
			fatal(e.Error())
		}
		fmt.Print("done\n")

		// Generate repos
		fmt.Printf("Generating %d repos...", len(database.Tables))
		for _, table := range database.Tables {
			e = gen.GenerateGoRepoFile(table)
			if e != nil {
				fatal(e.Error())
			}
		}
		fmt.Print("done\n")

		fmt.Print("Generating repo bootstrap file...")
		gen.GenerateReposBootstrapFile(database)
		fmt.Print("done\n")

		// Generate models
		fmt.Printf("Generating %d models...", len(database.Tables))
		for _, table := range database.Tables {
			e = gen.GenerateGoModelFile(table)
			if e != nil {
				fatal(e.Error())
			}
		}
		fmt.Print("done\n")

	case "typescript":
		gen.GenerateTypescriptTypesFile(database)
	default:
		fatal(fmt.Sprintf("Unknown output type: `%s`", subCmd))
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
		
		compare [reverse] [ ( write <path> | apply ) ]
			
			Default behavior (no arguments) is to compare local schema as authority against 
			remote database as target and write the resulting sql to stdout.

			reverse 	Switches the roles of the schemas. The remote database becomes the authority 
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

const (
	// OptVerbose triggers verbose logging
	OptVerbose = 1 << iota
	// OptVerboseDebug triggers extremely verbose logging
	OptVerboseDebug
	// OptSilent suppresses all logging
	OptSilent
)
