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

func main() {

	var e error
	var dvc *DVC
	args := os.Args

	if len(args) < 2 {
		log.Fatal("Missing command")
		return
	}

	cmd := args[1]

	if dvc, e = NewDVC(configFilePath); e != nil {
		fatal(e.Error())
	}

	switch cmd {
	case CommandImport:

		schemaFile := dvc.Config.DatabaseName + ".schema.json"

		if e = dvc.ImportSchema(schemaFile); e != nil {
			fatal(e.Error())
		}

		log.Printf("Schema `%s`.`%s` imported to %s.", dvc.Config.Host, dvc.Config.DatabaseName, schemaFile)
		os.Exit(0)

	case CommandGen:

		var database *Database

		if len(args) < 3 {
			fatal("Missing gen type [schema | model | repo]")
		}
		subCmd := args[2]

		// Load the schema
		schemaFile := dvc.Config.DatabaseName + ".schema.json"
		database, e = ReadSchemaFromFile(schemaFile)
		if e != nil {
			fatal(e.Error())
		}

		switch subCmd {
		case CommandGenSchema:
			e = GenerateGoSchemaFile(database)
			if e != nil {
				fatal(e.Error())
			}

		case CommandGenRepos:

			for _, table := range database.Tables {
				e = GenerateGoRepoFile(table)
				if e != nil {
					fatal(e.Error())
				}
			}

			GenerateReposBootstrapFile(database)

		case "repo":
			if len(args) < 4 {
				fatal("Missing repo name")
			}

			var t *Table

			if t, e = database.FindTableByName(args[3]); e != nil {
				fatal(e.Error())
			}

			if e = GenerateGoRepoFile(t); e != nil {
				fatal(e.Error())
			}

		case CommandGenModels:

			for _, table := range database.Tables {
				e = GenerateGoModelFile(table)
				if e != nil {
					fatal(e.Error())
				}
			}

		case "model":
			if len(args) < 4 {
				fatal("missing model name")
			}

			var t *Table

			if t, e = database.FindTableByName(args[3]); e != nil {
				fatal(e.Error())
			}

			if e = GenerateGoModelFile(t); e != nil {
				fatal(e.Error())
			}
		case CommandGenAll:

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

	case CommandCompare:

		sql := ""

		nextArgIdx := 2

		schemaFile := dvc.Config.DatabaseName + ".schema.json"

		isCompareFlipped := false
		if len(args) > 2 && args[nextArgIdx] == "reverse" {
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
				e = dvc.ApplyChangeset(sql)
				if e != nil {
					fatal(e.Error())
				}

			default:
				fatal(fmt.Sprintf("Unknown argument: `%s`\n", args[2]))
			}

		}

		// Write SQL to STDOUT
		fmt.Printf("%s", sql)
		os.Exit(0)

	case CommandHelp:
		dvc.PrintHelp()

	default:
		fatal(fmt.Sprintf("Unknown command: `%s`", cmd))
	}

}
