package main

import (
	"bufio"
	"encoding/json"
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
	CommandInspect       Command = "inspect"
	CommandRm            Command = "rm"
	CommandInit          Command = "init"
	CommandLs            Command = "ls"
	CommandImport        Command = "import"
	CommandExport        Command = "export"
	CommandGen           Command = "gen"
	CommandGenApp        Command = "app"
	CommandGenCLI        Command = "cli"
	CommandGenAPI        Command = "api"
	CommandGenDal        Command = "dal"
	CommandGenDals       Command = "dals"
	CommandGenRepos      Command = "repos"
	CommandGenModels     Command = "models"
	CommandGenInterfaces Command = "interfaces"
	CommandGenTests      Command = "tests"
	CommandGenModel      Command = "model"
	CommandGenServices   Command = "services"
	CommandGenRoutes     Command = "routes"
	CommandGenTypescript Command = "typescript"
	CommandCompare       Command = "compare"
	CommandHelp          Command = "help"
	CommandRefresh       Command = "refresh"
	CommandInstall       Command = "install"
	CommandInsert        Command = "insert"
	CommandSelect        Command = "select"
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

		if len(args) > 0 && args[0] == "help" {
			helpRefresh()
			return
		}

		c.CommandRefresh(args)

	case CommandInsert:
		if len(args) > 0 && args[0] == "help" {
			helpInsert()
			return
		}

		c.CommandInsert(args)

	case CommandSelect:
		if len(args) > 0 && args[0] == "help" {
			helpSelect()
			return
		}
		c.CommandSelect(args)
	case CommandImport:
		if len(args) > 0 && args[0] == "help" {
			helpImport()
			return
		}
		c.CommandImport(args)
	case CommandExport:
		if len(args) > 0 && args[0] == "help" {
			helpExport()
			return
		}
		c.CommandExport(args)
	case CommandCompare:
		if len(args) > 0 && args[0] == "help" {
			helpCompare()
			return
		}
		c.CommandCompare(args)
	case CommandGen:
		if len(args) > 0 && args[0] == "help" {
			helpGen()
			return
		}
		c.CommandGen(args)
	case CommandHelp:
		helpCommandNames()
	case CommandInit:
		if len(args) > 0 && args[0] == "help" {
			helpInit()
			return
		}
		c.CommandInit(args)
	case CommandLs:
		if len(args) > 0 && args[0] == "help" {
			helpLs()
			return
		}
		c.CommandLs(args)
	case CommandAdd:
		if len(args) > 0 && args[0] == "help" {
			helpAdd()
			return
		}
		c.CommandAdd(args)
	case CommandInspect:
		c.CommandInspect(args)
	case CommandRm:
		if len(args) > 0 && args[0] == "help" {
			helpRm()
			return
		}
		c.CommandRm(args)
	default:
		fmt.Printf("Invalid command `%s`\n", cmd)
		helpCommandNames()
	}

	os.Exit(0)
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

// CommandLs lists database information
// TODO search fields
// TODO select from tables
// TODO show row counts in a table
func (c *Cmd) CommandLs(args []string) {

	database := c.loadDatabase()

	// Options
	// ls 							Show all Tables
	// ls [name] 					Show all Columns in table [name] is found
	// ls [partialName] 			Show all tables with name containing [partialName]
	// ls .[fieldPartialName] 		Show all columns with name containing [fieldPartialName]

	// ls - Show all tables
	if len(args) == 0 {

		results := database.ToSortedTables()
		t := lib.NewCLITable([]string{"Table", "Columns", "Collation", "Engine", "Row Format"})

		for _, table := range results {
			t.Row()
			t.Col(table.Name)
			t.Colf("%d", len(table.Columns))
			t.Col(table.Collation)
			t.Col(table.Engine)
			t.Col(table.RowFormat)
		}

		fmt.Println(t.String())
		return
	}

	// ls .[fieldPartialName] - Show all columns with names containing [fieldPartialName]
	if strings.HasPrefix(args[0], ".") {

		if args[0] == "." {
			return
		}

		fieldSearch := strings.Trim(strings.ToLower(args[0][1:]), " ")

		fmt.Printf("Searching all fields for `%s`\n", fieldSearch)

		sortedTables := database.ToSortedTables()

		columns := []*lib.ColumnWithTable{}

		for j := range sortedTables {
			for k := range sortedTables[j].Columns {
				if strings.Contains(strings.ToLower(sortedTables[j].Columns[k].Name), fieldSearch) {
					column := &lib.ColumnWithTable{
						Column:    sortedTables[j].Columns[k],
						TableName: sortedTables[j].Name,
					}
					columns = append(columns, column)
				}
			}
		}

		t := lib.NewCLITable([]string{"Table", "Name", "Type", "MaxLength", "Null", "Default", "Extra", "Key"})

		for _, col := range columns {

			t.Row()
			t.Col(col.TableName)
			t.Col(col.Name)
			t.Col(col.DataType)
			t.Colf("%d", col.MaxLength)

			if col.IsNullable {
				t.Col("YES")
			} else {
				t.Col("NO")
			}

			t.Col(col.Default)
			t.Col(col.Extra)
			t.Col(col.ColumnKey)
		}

		fmt.Println(t.String())

		return
	}

	tableName := strings.Trim(strings.ToLower(args[0]), " ")

	// findTable
	for k := range database.Tables {
		if strings.ToLower(database.Tables[k].Name) == tableName {
			tableName = database.Tables[k].Name
			// Found
			break
		}
	}

	// ls [tableName] -  Show all column in the table [tableName]
	if _, ok := database.Tables[tableName]; ok {

		fmt.Printf("Listing all columns for table `%s`\n", tableName)

		table := database.Tables[tableName]

		sortedColumns := table.ToSortedColumns()

		t := lib.NewCLITable([]string{"Name", "Type", "MaxLength", "Null", "Default", "Extra", "Key"})

		for _, col := range sortedColumns {

			t.Row()

			t.Col(col.Name)
			t.Col(col.DataType)
			t.Colf("%d", col.MaxLength)

			if col.IsNullable {
				t.Col("YES")
			} else {
				t.Col("NO")
			}

			t.Col(col.Default)
			t.Col(col.Extra)
			t.Col(col.ColumnKey)
		}

		fmt.Println(t.String())
		return
	}

	// ls [partialTableName] - Show all tables with name containing [partialTableName]
	sortedTables := database.ToSortedTables()
	tableSearch := strings.Trim(strings.ToLower(args[0]), " ")

	t := lib.NewCLITable([]string{"Table", "Columns", "Collation", "Engine", "Row Format"})
	for _, table := range sortedTables {

		if !strings.Contains(strings.ToLower(table.Name), tableSearch) {
			continue
		}

		t.Row()
		t.Col(table.Name)
		t.Colf("%d", len(table.Columns))
		t.Col(table.Collation)
		t.Col(table.Engine)
		t.Col(table.RowFormat)
	}

	fmt.Println(t.String())
}

// CommandRefresh is the refresh command
func (c *Cmd) CommandRefresh(args []string) {
	totalTime := time.Now()

	// Import
	start := time.Now()
	c.CommandImport(args)
	fmt.Printf("Import: %f seconds\n", time.Since(start).Seconds())

	// Gen Models
	start = time.Now()
	c.CommandGen([]string{"models"})
	fmt.Printf("Models: %f seconds\n", time.Since(start).Seconds())

	start = time.Now()
	c.CommandGen([]string{"dals"})
	fmt.Printf("DALs: %f seconds\n", time.Since(start).Seconds())

	start = time.Now()
	c.CommandGen([]string{"interfaces"})
	fmt.Printf("Interfaces: %f seconds\n", time.Since(start).Seconds())

	start = time.Now()
	c.CommandGen([]string{"routes"})
	fmt.Printf("Routes: %f seconds\n", time.Since(start).Seconds())

	fmt.Printf("Total: %f seconds\n", time.Since(totalTime).Seconds())
}

// CommandInsert inserts data into the database
func (c *Cmd) CommandInsert(args []string) {

	if len(args) > 0 && args[0] == "help" {
		helpInsert()
		return
	}

	database := c.loadDatabase()

	reader := bufio.NewReader(os.Stdin)

	tableName := ""

	if len(args) > 0 {
		tableName = args[0]
	} else {
		tableName = lib.ReadCliInput(reader, "Table:")
	}

	if len(tableName) == 0 {
		fmt.Println("No table specified")
		return
	}

	if _, ok := database.Tables[tableName]; !ok {
		fmt.Printf("Table `%s` not found.\n", tableName)
		return
	}

	table := database.Tables[tableName]

	columns := table.ToSortedColumns()

	sql := fmt.Sprintf("INSERT INTO `%s` (\n", tableName)

	columnNames := []string{}
	values := []string{}

	for k := range columns {

		if columns[k].ColumnKey == "PRI" {
			continue
		}

		if columns[k].Name == "IsDeleted" {
			continue
		}

		columnNames = append(columnNames, fmt.Sprintf("`%s`", columns[k].Name))

		value := "?"
		if columns[k].Name == "DateCreated" {
			value = fmt.Sprintf("%d", time.Now().UnixNano()/1000000)
		} else {
			value = lib.ReadCliInput(reader, columns[k].Name+" ("+columns[k].DataType+"):")
		}

		if lib.IsString(columns[k]) {
			value = "'" + value + "'"
		} else {
			if len(value) == 0 {
				value = "0"
			}
		}

		values = append(values, value)
	}

	sql += "\t" + strings.Join(columnNames, ",\n\t")
	sql += "\n) VALUES (\n"
	sql += "\t" + strings.Join(values, ",\n\t")
	sql += "\n)\n"

	fmt.Println(sql)

	doInsertYN := lib.ReadCliInput(reader, "Run above SQL (Y/n")

	if doInsertYN != "Y" {
		return
	}

	connector, _ := connectorFactory(c.Config.DatabaseType, c.Config)
	x := lib.NewExecutor(c.Config, connector)
	x.RunSQL(sql)
}

// CommandAdd adds an object to the database
func (c *Cmd) CommandAdd(args []string) {

	database := c.loadDatabase()

	reader := bufio.NewReader(os.Stdin)

	sqlParts := []string{}
	sql := ""

	tableName := ""

	if len(args) > 0 {
		tableName = args[0]
	} else {
		tableName = lib.ReadCliInput(reader, "Table Name:")
	}

	// Creating a new table
	if _, ok := database.Tables[tableName]; !ok {

		fmt.Printf("Create table `%s`\n", tableName)

		start := fmt.Sprintf("CREATE TABLE `%s` (", tableName)
		sql += start
		sqlParts = append(sqlParts, start)
		cols := []string{
			fmt.Sprintf("%sID INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT", tableName),
			fmt.Sprintf("IsDeleted TINYINT UNSIGNED NOT NULL DEFAULT 0"),
			fmt.Sprintf("DateCreated BIGINT UNSIGNED NOT NULL DEFAULT 0"),
			fmt.Sprintf("LastUpdated BIGINT UNSIGNED NOT NULL DEFAULT 0"),
		}

		for k := range cols {
			fmt.Printf("`%s`.Column(%d): %s\n", tableName, k+1, cols[k])
		}

		includeColumns := lib.ReadCliInput(reader, fmt.Sprintf("Include the above columns (Y/n)?"))

		n := 1
		if includeColumns == "Y" {
			n = len(cols) + 1
		} else {
			cols = []string{}
		}

		for {

			col := lib.ReadCliInput(reader, fmt.Sprintf("`%s`.Column(%d):", tableName, n))
			if len(col) == 0 {
				break
			}

			// Remove trailing space
			col = strings.Trim(col, " ")

			// Remove trailing comma if exists
			if col[len(col)-1:] == "," {
				col = col[0 : len(col)-1]
			}

			colParts := strings.Split(col, " ")

			// Validate the sql part
			if len(colParts) < 2 {
				fmt.Println("Invalid sql part. Try again...")
				continue
			}

			// Validate the datatype
			dataType := strings.ToLower(colParts[1])

			if strings.Contains(dataType, "(") {
				dataType = strings.Split(dataType, "(")[0]
			}

			if lib.IsValidSQLType(dataType) == false {
				fmt.Printf("Invalid SQL type: %s. Try gain...\n", dataType)
				continue
			}

			n++

			for k := range colParts {

				if k == 0 {
					continue
				}

				colParts[k] = strings.ToUpper(colParts[k])
			}

			col = "`" + colParts[0] + "` " + strings.Join(colParts[1:], " ")
			sqlParts = append(sqlParts, col)
			cols = append(cols, col)
		}

		sql += strings.Join(cols, ", ") + ")"

		sqlParts = append(sqlParts, ")")

	} else {

		col := lib.ReadCliInput(reader, fmt.Sprintf("`%s`.Column:", tableName))
		if len(col) == 0 {
			fmt.Println("Column cannot be empty")
			return
		}

		colParts := strings.Split(col, " ")

		if _, ok := database.Tables[tableName].Columns[colParts[0]]; ok {
			fmt.Printf("Column `%s`.`%s` already exists", tableName, colParts[0])
			return
		}

		for k := range colParts {

			if k == 0 {
				continue
			}

			colParts[k] = strings.ToUpper(colParts[k])
		}

		col = "`" + colParts[0] + "` " + strings.Join(colParts[1:], " ")
		sql = fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s", tableName, col)
		sqlParts = append(sqlParts, sql)
	}

	fmt.Print("\n------------------------------------------------------\n")
	fmt.Print("--------------------- REVIEW -------------------------\n")
	fmt.Print("------------------------------------------------------\n")

	fmt.Println(sql)
	// for k := range sqlParts {
	// 	// Not first or last line
	// 	if k != 0 && k != len(sqlParts)-1 {
	// 		fmt.Print("\t")
	// 	}
	// 	fmt.Printf("%s", sqlParts[k])
	// 	// Not first or last line
	// 	if k != 0 && k != len(sqlParts)-1 {
	// 		fmt.Print(",")
	// 	}

	// 	fmt.Print("\n")
	// }

	fmt.Print("\n------------------------------------------------------\n")

	if lib.ReadCliInput(reader, "Are you sure want to execute the above SQL (Y/n)?") == "Y" {
		connector, _ := connectorFactory(c.Config.DatabaseType, c.Config)
		x := lib.NewExecutor(c.Config, connector)
		x.RunSQL(sql)

		c.CommandRefresh([]string{})
	}
}

// CommandInspect is the inspect command
func (c *Cmd) CommandInspect(args []string) {

	model, _ := gen.InspectFile(args[0])

	k := 0
	for k < model.Fields.Len() {
		fmt.Println(model.Fields.Get(k).Name + " > " + model.Fields.Get(k).DataType)
	}
}

// CommandSelect selects rows from the database
func (c *Cmd) CommandSelect(args []string) {

	// database := c.loadDatabase()

	// reader := bufio.NewReader(os.Stdin)

	// tableName := ""

	// if len(args) > 0 {
	// 	tableName = args[0]
	// 	args = args[1:]
	// } else {
	// 	tableName = lib.ReadCliInput(reader, "Table Name:")
	// }

	// if _, ok := database.Tables[tableName]; !ok {
	// 	fmt.Printf("Unknown table `%s`\n", tableName)
	// 	return
	// }

	// query := fmt.Sprintf("SELECT * FROM `%s` LIMIT 100\n", tableName)

}

// CommandRm removes an object from the database
// dvc rm [table]
func (c *Cmd) CommandRm(args []string) {

	database := c.loadDatabase()

	reader := bufio.NewReader(os.Stdin)

	sql := ""

	tableName := ""
	if len(args) > 0 {
		tableName = args[0]
	} else {
		tableName = lib.ReadCliInput(reader, "Table Name:")
		if _, ok := database.Tables[tableName]; !ok {
			fmt.Printf("Table `%s` does not exist.", tableName)
			return
		}
	}

	tableOrColumn := lib.ReadCliInput(reader, fmt.Sprintf("Drop the (t)able `%s` or select a (c)olumn?", tableName))
	if tableOrColumn == "t" {
		sql = fmt.Sprintf("DROP TABLE `%s`", tableName)
	} else if tableOrColumn == "c" {
		columnName := lib.ReadCliInput(reader, fmt.Sprintf("`%s`.Column:", tableName))
		if _, ok := database.Tables[tableName].Columns[columnName]; !ok {
			fmt.Printf("Column `%s`.`%s` doesn't exist.", tableName, columnName)
			return
		}

		sql = fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`", tableName, columnName)
	} else {
		fmt.Println("Invalid entry")
		return
	}

	fmt.Print("\n------------------------------------------------------\n")
	fmt.Print("--------------------- REVIEW -------------------------\n")
	fmt.Print("------------------------------------------------------\n")

	fmt.Println(sql)

	fmt.Print("\n------------------------------------------------------\n")

	if lib.ReadCliInput(reader, "Are you sure want to execute the above SQL (Y/n)?") == "Y" {

		// Apply the change
		connector, _ := connectorFactory(c.Config.DatabaseType, c.Config)
		x := lib.NewExecutor(c.Config, connector)
		x.RunSQL(sql)

		// Import the schema
		c.CommandImport([]string{})
	}
}

// CommandImport fetches the sql schema from the target database (specified in dvc.toml)
// and from that generates the json representation at `[schema name].schema.json`
func (c *Cmd) CommandImport(args []string) {

	fmt.Println("Importing...")
	var e error
	cmp := c.initCompare()

	if e = cmp.ImportSchema("./" + c.Config.Connection.DatabaseName + ".schema.json"); e != nil {
		lib.Error(e.Error(), c.Options)
		os.Exit(1)
	}

	curDir, _ := os.Getwd()
	lib.Infof("Schema `%s`.`%s` imported to %s", c.Options, c.Config.Connection.Host, c.Config.Connection.DatabaseName, path.Join(curDir, c.Config.Connection.DatabaseName+".schema.json"))
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

// CommandGen handles the `gen` command
func (c *Cmd) CommandGen(args []string) {

	var e error

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

	g := &gen.Gen{
		Options: c.Options,
		Config:  c.Config,
	}

	database := c.loadDatabase()
	c.genTableCache(database)

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
		c.GenRepos(g, database)
	case CommandGenDals:
		c.GenDals(g, database)
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

	case CommandGenInterfaces:
		c.GenInterfaces(g)
		// result, err := interfaces.Make(files, args.StructType, args.Comment, args.PkgName, args.IfaceName, args.IfaceComment, args.copyDocs, args.CopyTypeDoc)
	case CommandGenRoutes:
		c.GenRoutes(g)
	case CommandGenTests:

		serviceSuffix := "Service"
		srcDir := c.Config.Dirs.Services
		fmt.Println(srcDir)
		var files []os.FileInfo
		// DAL
		if files, e = ioutil.ReadDir(srcDir); e != nil {
			fmt.Println("ERROR", e.Error())
			return
		}
		for _, f := range files {

			// Filter out files that don't have upper case first letter names
			if !unicode.IsUpper([]rune(f.Name())[0]) {
				continue
			}

			srcFile := path.Join(srcDir, f.Name())

			// Remove the go extension
			baseName := f.Name()[0 : len(f.Name())-3]

			// Skip over EXT files
			if baseName[len(baseName)-3:] == "Ext" {
				continue
			}

			// Skip over test files
			if baseName[len(baseName)-5:] == "_test" {
				continue
			}

			// fmt.Println(baseName)

			if baseName == "DesignationService" {
				e = g.GenServiceTest(baseName[0:len(baseName)-len(serviceSuffix)], srcFile)

				if e != nil {
					panic(e)
				}
			}
		}

	case CommandGenModels:

		c.GenModels(g, database)

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

	

	

	

	
	help 	This output

	
	

		

	
`
	fmt.Printf(help)
}

// GenRepos generates repos
func (c *Cmd) GenRepos(g *gen.Gen, database *lib.Database) {

	var e error

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
}

// GenModels generates models
func (c *Cmd) GenModels(g *gen.Gen, database *lib.Database) {

	fmt.Println("Generating models...")

	var e error

	modelsDir := path.Join(c.Config.Dirs.Definitions, "models")
	if c.Options&lib.OptClean == lib.OptClean {
		g.CleanGoModels(modelsDir, database)
	}

	for _, table := range database.Tables {

		// If the model has been hashed before...
		if _, ok := c.existingModels.Models[table.Name]; ok {

			// And the hash hasn't changed, skip...
			if c.newModels[table.Name] == c.existingModels.Models[table.Name] {

				// fmt.Printf("Table `%s` hasn't changed! Skipping...\n", table.Name)
				continue
			}
		}

		// Update the models cache
		c.existingModels.Models[table.Name] = c.newModels[table.Name]

		// fmt.Printf("Generating model `%s`\n", table.Name)
		e = g.GenerateGoModel(modelsDir, table)
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}
	}

	c.saveTableCache()

	// e = g.GenerateDALSQL(c.Config.Dirs.Dal, database)
	// if e != nil {
	// 	lib.Error(e.Error(), c.Options)
	// 	os.Exit(1)
	// }
}

// GenDals generates dals
func (c *Cmd) GenDals(g *gen.Gen, database *lib.Database) {

	fmt.Println("Generating dals...")
	var e error

	// Loop through the schema's tables and build md5 hashes of each to check against
	for _, table := range database.Tables {

		// If the model has been hashed before...
		if _, ok := c.existingModels.Dals[table.Name]; ok {

			// And the hash hasn't changed, skip...
			if c.newModels[table.Name] == c.existingModels.Dals[table.Name] {

				// fmt.Printf("Table `%s` hasn't changed! Skipping...\n", table.Name)
				continue
			}
		}

		// Update the dals cache
		c.existingModels.Dals[table.Name] = c.newModels[table.Name]

		// fmt.Printf("Generating %sDAL...\n", table.Name)
		e = g.GenerateGoDAL(table, c.Config.Dirs.Dal)
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}
	}

	c.saveTableCache()

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
}

func (c *Cmd) GenInterfaces(g *gen.Gen) {

	fmt.Println("Generating interfaces...")
	var e error

	genInterfaces := func(srcDir, srcType string) (e error) {

		var files []os.FileInfo
		// DAL
		if files, e = ioutil.ReadDir(srcDir); e != nil {
			return
		}
		for _, f := range files {

			// Filter out files that don't have upper case first letter names
			if !unicode.IsUpper([]rune(f.Name())[0]) {
				continue
			}

			srcFile := path.Join(srcDir, f.Name())

			// Remove the go extension
			baseName := f.Name()[0 : len(f.Name())-3]

			// Skip over EXT files
			if baseName[len(baseName)-3:] == "Ext" {
				continue
			}

			// Skip over test files
			if baseName[len(baseName)-5:] == "_test" {
				continue
			}

			// srcFile := path.Join(c.Config.Dirs.Dal, baseName + ".go")
			destFile := path.Join(c.Config.Dirs.Definitions, srcType, "I"+baseName+".go")
			interfaceName := "I" + baseName
			packageName := srcType

			srcFiles := []string{srcFile}
			// var src []byte
			// if src, e = ioutil.ReadFile(srcFile); e != nil {
			// 	return
			// }

			// Check if EXT file exists
			extFile := srcFile[0:len(srcFile)-3] + "Ext.go"
			if _, e = os.Stat(extFile); e == nil {
				srcFiles = append(srcFiles, extFile)
				// concatenate the contents of the ext file with the contents of the regular file
				// var extSrc []byte
				// if extSrc, e = ioutil.ReadFile(extFile); e != nil {
				// 	return
				// }
				// src = append(src, extSrc...)
			}

			var i []byte
			i, e = gen.GenInterfaces(
				srcFiles,
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
}

func (c *Cmd) GenRoutes(g *gen.Gen) {

	fmt.Println("Generating routes...")
	var e error

	e = g.GenRoutes("core/controllers")
	if e != nil {
		lib.Error(e.Error(), c.Options)
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

	tableCachePath := ".tables"

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
	var e error
	newModelData := []byte{}
	if newModelData, e = json.Marshal(c.existingModels); e != nil {
		panic(fmt.Sprintf("Can't serialize newModels: %s", e.Error()))
	}
	if e = ioutil.WriteFile(".tables", newModelData, 0777); e != nil {
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
