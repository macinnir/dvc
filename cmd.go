package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
	"unicode"

	"github.com/BurntSushi/toml"
	"github.com/macinnir/dvc/modules/compare"
	"github.com/macinnir/dvc/modules/gen"
	"github.com/macinnir/dvc/modules/gen/interfaces"

	"github.com/macinnir/dvc/connectors/mysql"
	"github.com/macinnir/dvc/connectors/sqlite"
	"github.com/macinnir/dvc/lib"
)

// Command is a type that represents the possible commands passed in at run time
type Command string

// Command Names
const (
	CommandInit          Command = "init"
	CommandImport        Command = "import"
	CommandExport        Command = "export"
	CommandGen           Command = "gen"
	CommandGenApp        Command = "app"
	CommandGenCLI        Command = "cli"
	CommandGenAPI        Command = "api"
	CommandGenDal        Command = "dal"
	CommandGenRepos      Command = "repos"
	CommandGenModels     Command = "models"
	CommandGenInterfaces Command = "interfaces"
	CommandGenModel      Command = "model"
	CommandGenServices   Command = "services"
	CommandGenTypescript Command = "typescript"
	CommandCompare       Command = "compare"
	CommandHelp          Command = "help"
	CommandRefresh       Command = "refresh"
	CommandInstall       Command = "install"
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

	if cmd != CommandInit {
		var e error
		c.Config, e = loadConfigFromFile("./dvc.toml")

		if e != nil {
			fmt.Println("Could not load config file.")
			// fmt.Printf("ERROR: %s\n", e.Error())
			// os.Exit(1)
			return
		}

	}

	switch cmd {
	// case CommandRefresh:
	// 	c.CommandImport(args)
	// 	c.CommandGen([]string{"models"})
	// 	c.CommandGen([]string{"dal"})
	// 	c.CommandGen([]string{"repos"})
	// 	c.CommandGen([]string{"services"})
	// 	c.CommandGen([]string{"api"})
	// case CommandInstall:
	// 	c.CommandImport(args)
	// 	c.CommandGen([]string{"app"})
	// 	c.CommandGen([]string{"models"})
	// 	c.CommandGen([]string{"dal"})
	// 	c.CommandGen([]string{"repos"})
	// 	c.CommandGen([]string{"services"})
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
	case CommandInit:
		c.CommandInit(args)
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

// CommandInit creates a default dvc.toml file in the CWD
func (c *Cmd) CommandInit(args []string) {

	var e error

	if _, e = os.Stat("./dvc.toml"); os.IsNotExist(e) {

		reader := bufio.NewReader(os.Stdin)

		// https://tutorialedge.net/golang/reading-console-input-golang/
		// BasePackage
		fmt.Print("> Base Package:")
		basePackage, _ := reader.ReadString('\n')
		basePackage = strings.Replace(basePackage, "\n", "", -1)

		fmt.Print("> Base directory (leave blank for current):")
		baseDir, _ := reader.ReadString('\n')
		baseDir = strings.Replace(baseDir, "\n", "", -1)

		// Host
		fmt.Print("> Database Host:")
		host, _ := reader.ReadString('\n')
		host = strings.Replace(host, "\n", "", -1)

		// databaseName
		fmt.Print("> Database Name:")
		databaseName, _ := reader.ReadString('\n')
		databaseName = strings.Replace(databaseName, "\n", "", -1)

		// databaseUser
		fmt.Print("> Database User:")
		databaseUser, _ := reader.ReadString('\n')
		databaseUser = strings.Replace(databaseUser, "\n", "", -1)

		// databasePass
		fmt.Print("> Database Password:")
		databasePass, _ := reader.ReadString('\n')
		databasePass = strings.Replace(databasePass, "\n", "", -1)

		content := "databaseType = \"mysql\"\nbasePackage = \"" + basePackage + "\"\n\nenums = []\n\n"
		content += "[connection]\nhost = \"" + host + "\"\ndatabaseName = \"" + databaseName + "\"\nusername = \"" + databaseUser + "\"\npassword = \"" + databasePass + "\"\n\n"

		packages := []string{
			"repos",
			"models",
			"typescript",
			"services",
			"dal",
			"definitions",
		}

		content += "[packages]\n"
		for _, p := range packages {
			if p == "typescript" {
				continue
			}

			content += fmt.Sprintf("%s = \"%s\"\n", p, path.Join(basePackage, p))
		}

		// content += "[packages]\ncache = \"myPackage/cache\"\nmodels = \"myPackage/models\"\nschema = \"myPackage/schema\"\nrepos = \"myPackage/repos\"\n\n"

		content += "[dirs]\n"

		for _, p := range packages {

			if baseDir != "" {
				content += fmt.Sprintf("%s = \"%s\"\n", p, path.Join(baseDir, p))
			} else {
				content += fmt.Sprintf("%s = \"%s\"\n", p, p)
			}
		}

		// content += "[dirs]\nrepos = \"repos\"\ncache = \"cache\"\nmodels = \"models\"\nschema = \"schema\"\ntypescript = \"ts\""

		ioutil.WriteFile("./dvc.toml", []byte(content), 0644)

	} else {
		fmt.Println("dvc.toml already exists in this directory")
	}

}

// CommandImport fetches the sql schema from the target database (specified in dvc.toml)
// and from that generates the json representation at `[schema name].schema.json`
func (c *Cmd) CommandImport(args []string) {

	var e error
	cmp := c.initCompare()

	if e = cmp.ImportSchema("./" + c.Config.Connection.DatabaseName + ".schema.json"); e != nil {
		lib.Error(e.Error(), c.Options)
		os.Exit(1)
	}

	lib.Infof("Schema `%s`.`%s` imported to %s.", c.Options, c.Config.Connection.Host, c.Config.Connection.DatabaseName, c.Config.Connection.DatabaseName+".schema.json")
}

// CommandExport export SQL create statements to standard out
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

	cmp.Options = c.Options

	// Do the comparison
	// TODO pass all options (e.g. verbose)
	if sql, e = cmp.CompareSchema(schemaFile); e != nil {
		lib.Error(e.Error(), c.Options)
		os.Exit(1)
	}

	if len(sql) == 0 {
		fmt.Println("No changes found")
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
		writeSQLToLog(sql)
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

func writeSQLToLog(sql string) {

	sqlLog := time.Now().Format("20060102150405") + "\n"
	sqlLog += sql

	filePath := "./dvc-changes.log"

	if _, e := os.Stat(filePath); os.IsNotExist(e) {
		ioutil.WriteFile(filePath, []byte(sqlLog), 0600)
	} else {
		f, err := os.OpenFile("./dvc-changes.log", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if _, err = f.WriteString(sqlLog); err != nil {
			panic(err)
		}
	}

}

// CommandGen handles the `gen` command
func (c *Cmd) CommandGen(args []string) {

	var e error
	var database *lib.Database
	// fmt.Printf("Args: %v", args)
	if len(args) < 1 {
		lib.Error("Missing gen type [schema | models | repos | caches | ts]", c.Options)
		os.Exit(1)
	}
	subCmd := Command(args[0])
	cwd, _ := os.Getwd()

	if len(args) > 0 {
		args = args[1:]
	}

	argLen := len(args)
	n := 0

	// dvc gen models -c
	for n < argLen {
		switch args[n] {
		case "-c":
			c.Options |= lib.OptClean
		}
		n++
	}

	// for len(args) > 0 {
	// 	switch args[0] {
	// 	case "-c", "--clean":
	// 		c.Options |= lib.OptClean
	// 	}
	// 	args = args[1:]
	// }

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
	// case CommandGenSchema:
	// 	e = g.GenerateGoSchemaFile(c.Config.Dirs.Schema, database)
	// 	if e != nil {
	// 		lib.Error(e.Error(), c.Options)
	// 		os.Exit(1)
	// 	}
	// case CommandGenCaches:
	// 	fmt.Println("CommandGenCaches")
	// 	e = g.GenerateGoCacheFiles(c.Config.Dirs.Cache, database)
	// 	if e != nil {
	// 		lib.Error(e.Error(), c.Options)
	// 		os.Exit(1)
	// 	}
	case CommandGenRepos:

		if c.Options&lib.OptClean == lib.OptClean {
			g.CleanGoRepos(c.Config.Dirs.Repos, database)
		}

		e = g.GenerateGoRepoFiles(c.Config.Dirs.Repos, database)
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}

		e = g.GenerateReposBootstrapFile(c.Config.Dirs.Repos, database)
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}

		lib.Debug("Generating repo interfaces at "+c.Config.Dirs.Definitions, c.Options)
		g.EnsureDir(c.Config.Dirs.Definitions)
		e = g.GenerateRepoInterfaces(database, c.Config.Dirs.Definitions)
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}
	case CommandGenDal:

		if argLen == 0 {
			lib.Error("Missing dal name", c.Options)
			os.Exit(1)
		}

		// lib.Error(fmt.Sprintf("Args: %s", args[0]), c.Options)
		table, e := database.FindTableByName(args[0])
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}

		e = g.GenerateGoDAL(table, c.Config.Dirs.Dal)
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}

		// if c.Options&lib.OptClean == lib.OptClean {
		// 	g.CleanGoDALs(c.Config.Dirs.Dal, database)
		// }

		// for _, table := range database.Tables {

		// 	lib.Debugf("Generating dal %s", g.Options, table.Name)
		// 	e = g.GenerateGoDAL(table, c.Config.Dirs.Dal)
		// 	if e != nil {
		// 		return
		// 	}
		// }

		// if e != nil {
		// 	lib.Error(e.Error(), c.Options)
		// 	os.Exit(1)
		// }

		// Create the dal bootstrap file in the dal repo
		e = g.GenerateDALsBootstrapFile(c.Config.Dirs.Dal, database)
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}

		e = g.GenerateDALSQL(c.Config.Dirs.Dal, database)
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}

		// lib.Debug("Generating dal interfaces at "+c.Config.Dirs.Definitions, c.Options)
		// g.EnsureDir(c.Config.Dirs.Definitions)
		// e = g.GenerateDALInterfaces(database, c.Config.Dirs.Definitions)
		// if e != nil {
		// 	lib.Error(e.Error(), c.Options)
		// 	os.Exit(1)
		// }

	case CommandGenInterfaces:

		genInterfaces := func(srcDir, srcType string) (e error) {

			var files []os.FileInfo
			// DAL
			if files, e = ioutil.ReadDir(srcDir); e != nil {
				return
			}
			for _, f := range files {

				if !unicode.IsUpper([]rune(f.Name())[0]) {
					continue
				}

				srcFile := path.Join(srcDir, f.Name())
				baseName := f.Name()[0 : len(f.Name())-3]

				// srcFile := path.Join(c.Config.Dirs.Dal, baseName + ".go")
				destFile := path.Join(c.Config.Dirs.Definitions, srcType, "I"+baseName+".go")
				interfaceName := "I" + baseName
				packageName := srcType

				var i []byte
				i, e = interfaces.Run(
					[]string{srcFile},
					baseName,
					"Generated Code; DO NOT EDIT.",
					packageName,
					interfaceName,
					fmt.Sprintf("%s describes the %s", interfaceName, baseName),
					true,
					true,
				)
				if e != nil {
					fmt.Println("ERROR", e.Error())
					return
				}

				// fmt.Println("Generating ", destFile)
				// fmt.Println("Writing to: ", destFile)

				ioutil.WriteFile(destFile, i, 0644)

				// fmt.Println("Name: ", baseName, "Path: ", srcFile)

			}

			return
		}

		e = genInterfaces(c.Config.Dirs.Dal, "dal")
		if e != nil {
			fmt.Println("ERROR", e.Error())
			os.Exit(1)
		}

		e = genInterfaces(c.Config.Dirs.Services, "services")
		if e != nil {
			fmt.Println("ERROR", e.Error())
			os.Exit(1)
		}

		os.Exit(0)

		// result, err := interfaces.Make(files, args.StructType, args.Comment, args.PkgName, args.IfaceName, args.IfaceComment, args.copyDocs, args.CopyTypeDoc)

	case CommandGenModels:

		modelsDir := path.Join(c.Config.Dirs.Definitions, "models")
		if c.Options&lib.OptClean == lib.OptClean {
			g.CleanGoModels(modelsDir, database)
		}

		for _, table := range database.Tables {
			e = g.GenerateGoModel(modelsDir, table)
			if e != nil {
				lib.Error(e.Error(), c.Options)
				os.Exit(1)
			}
		}

		e = g.GenerateDALSQL(c.Config.Dirs.Dal, database)
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}

		// // Config.go
		// if _, e = os.Stat(path.Join(modelsDir, "Config.go")); os.IsNotExist(e) {
		// 	lib.Debugf("Generating default Config.go file at %s", c.Options, path.Join(modelsDir, "Config.go"))
		// 	g.GenerateDefaultConfigModelFile(modelsDir)
		// }

		// config.json
		// if _, e = os.Stat(path.Join(cwd, "config.json")); os.IsNotExist(e) {
		// 	lib.Debugf("Generating default config.json file at %s", c.Options, path.Join(cwd, "config.json"))
		// 	g.GenerateDefaultConfigJsonFile(cwd)
		// }

	case CommandGenServices:
		g.GenerateServiceInterfaces(c.Config.Dirs.Definitions, c.Config.Dirs.Services)
		g.GenerateServiceBootstrapFile(c.Config.Dirs.Services)

	case CommandGenApp:
		g.GenerateGoApp(cwd)
	// case CommandGenCLI:
	// g.GenerateGoCLI(cwd)
	case CommandGenAPI:
		// g.GenerateGoAPI(cwd)
		g.GenerateAPIRoutes(c.Config.Dirs.API)
	case "ts":
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
	
	init 	Initialize a dvc.toml configuration file. 

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
