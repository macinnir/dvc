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
		case "all":

			// Generate schema
			fmt.Print("Generating schema...")
			e = GenerateGoSchemaFile(database)
			if e != nil {
				fatal(e.Error())
			}
			fmt.Print("done\n")

			// Generate repos
			fmt.Printf("Generating %d repos...", len(database.Tables))
			for _, table := range database.Tables {
				e = GenerateGoRepoFile(table)
				if e != nil {
					fatal(e.Error())
				}
			}
			fmt.Print("done\n")

			fmt.Print("Generating repo bootstrap file...")
			GenerateReposBootstrapFile(database)
			fmt.Print("done\n")

			// Generate models
			fmt.Printf("Generating %d models...", len(database.Tables))
			for _, table := range database.Tables {
				e = GenerateGoModelFile(table)
				if e != nil {
					fatal(e.Error())
				}
			}
			fmt.Print("done\n")

		case "typescript":
			GenerateTypescriptTypesFile(database)
		default:
			fatal(fmt.Sprintf("Unknown output type: `%s`", subCmd))
		}

	case "compare":

		sql := ""

		nextArgIdx := 2

		schemaFile := dvc.Config.DatabaseName + ".schema.json"

		isCompareFlipped := false
		if argLen > 2 && args[nextArgIdx] == "reverse" {
			isCompareFlipped = true
			nextArgIdx++
		}

		if sql, e = dvc.CompareSchema(schemaFile, isCompareFlipped); e != nil {
			fatal(e.Error())
		}

		if len(sql) == 0 {
			log.Printf("No changes found.")
			os.Exit(0)
		}

		if argLen > nextArgIdx {
			// "--file="

			switch args[nextArgIdx] {
			case "write":

				nextArgIdx++

				if argLen > nextArgIdx {

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
		fmt.Println("\timport")
		fmt.Println("\t\tBuild a schema definition file based on the target database. This will overwrite any existing schema definition file.")
		fmt.Println("\tgen ( models | repos | schema | all)")
		fmt.Println("\t\tall\t Generate all options below")
		fmt.Println("\t\tmodels\t Generate models based on imported schema information. Will fail if no imported schema file exists.")
		fmt.Println("\t\trepos\t Generate repositories based on imported schema information. Will fail if no imported schema file exists.")
		fmt.Println("\t\tschema\t Generate go-dal schema bootstrap code based on imported schema information. Will fail if no imported schema file exists.")
		fmt.Println("\tcompare [reverse] [ ( write <path> | apply ) ]")
		fmt.Println("\t\t Default behavior (no arguments) is to compare local schema as authority against remote database as target and write the resulting sql to stdout.")
		fmt.Println("\t\t reverse\tOptional reverse command will swith the roles of the schemas, making the remote database the authority and the local schema the target for updating.")
		fmt.Println("\t\t write\tAfter performing the comparison, the resulting sequel statements will be written to a filepath <path> (required).")
		fmt.Println("\t\t apply\tAfter performing the comparison, the resulting sequel statements will be immediately applied to the target database.")
		fmt.Println("\timport")

	default:
		fatal(fmt.Sprintf("Unknown command: `%s`", cmd))
	}

}
