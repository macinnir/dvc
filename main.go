package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	configFilePath = "dvc.toml"
)

func fatal(msg string) {
	fmt.Printf("ERROR: %s\n", msg)
	os.Exit(1)
}

func main() {

	var e error
	var dvc *DVC

	args := os.Args

	argLen := len(args)

	if argLen < 2 {
		fatal("Missing command...")
	}

	cmd := args[1]

	if dvc, e = NewDVC(configFilePath); e != nil {
		fatal(e.Error())
	}

	switch cmd {
	case "import":

		schemaFile := dvc.Config.DatabaseName + ".schema.json"

		if e = dvc.ImportSchema(schemaFile); e != nil {
			fatal(e.Error())
		}

		log.Printf("Schema `%s`.`%s` imported to %s.", dvc.Config.Host, dvc.Config.DatabaseName, schemaFile)
		os.Exit(0)

	case "gen":

		var database *Database

		if argLen < 3 {
			fatal("Missing gen type [schema | model | repo]")
		}
		subCmd := args[2]

		// Load the schema
		schemaFile := dvc.Config.DatabaseName + ".schema.json"
		database, e = dvc.ReadSchemaFromFile(schemaFile)
		if e != nil {
			fatal(e.Error())
		}

		switch subCmd {
		case "schema":
			e = GenerateGoSchemaFile(database)
			if e != nil {
				fatal(e.Error())
			}
		case "repos":

			for _, table := range database.Tables {
				e = GenerateGoRepoFile(table)
				if e != nil {
					fatal(e.Error())
				}
			}

			GenerateReposBootstrapFile(database)

		case "repo":
			if argLen < 4 {
				fatal("Missing repo name")
			}

			var t *Table

			if t, e = FindTableByName(database, args[3]); e != nil {
				fatal(e.Error())
			}

			if e = GenerateGoRepoFile(t); e != nil {
				fatal(e.Error())
			}

		case "models":

			for _, table := range database.Tables {
				e = GenerateGoModelFile(table)
				if e != nil {
					fatal(e.Error())
				}
			}

		case "model":
			if argLen < 4 {
				fatal("missing model name")
			}

			var t *Table

			if t, e = FindTableByName(database, args[3]); e != nil {
				fatal(e.Error())
			}

			if e = GenerateGoModelFile(t); e != nil {
				fatal(e.Error())
			}

		case "typescript":
			GenerateTypescriptTypesFile(database)
		default:
			fatal(fmt.Sprintf("Unknown output type: `%s`", subCmd))
		}

	case "compare":

		sql := ""

		schemaFile := dvc.Config.DatabaseName + ".schema.json"

		if sql, e = dvc.CompareSchema(schemaFile); e != nil {
			fatal(e.Error())
		}

		if len(sql) == 0 {
			log.Printf("No changes found.")
			os.Exit(0)
		}

		if argLen > 2 {
			// "--file="

			switch args[2] {
			case "write":
				if argLen > 3 {

					filePath := args[3]

					log.Printf("Writing sql to path `%s`", filePath)
					e = ioutil.WriteFile(filePath, []byte(sql), 0644)
					if e != nil {
						fatal(e.Error())
					}
				} else {
					fatal("Missing path to write sql to")
				}
			case "apply":
				e = dvc.ApplyChangeset(sql)
				if e != nil {
					fatal(e.Error())
				}

			default:
				fatal(fmt.Sprintf("Unknown argument: `%s`\n", args[2]))
				os.Exit(1)
			}

		}

		// Write SQL to STDOUT
		fmt.Printf("%s", sql)
		os.Exit(0)

	case "help":
		fmt.Println("DVC Help")
		fmt.Println("\tcompare")
		fmt.Println("\t\t[default] - Write the change sql to stdout")
		fmt.Println("\t\twrite [path] - Write the sql to a file path")
		fmt.Println("\t\tapply - Apply the changes after the comparison sql has been generated.")
		fmt.Println("\timport")

	default:
		fatal(fmt.Sprintf("Unknown command: `%s`", cmd))
	}

}
