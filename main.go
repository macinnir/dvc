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

	argLen := len(args)

	if argLen < 2 {
		fmt.Println("Missing command...")
		os.Exit(1)
	}

	cmd := args[1]

	if dvc, e = NewDVC(configFilePath); e != nil {
		log.Printf("ERROR: %s", e.Error())
		os.Exit(1)
	}

	switch cmd {
	case "import":

		schemaFile := dvc.Config.DatabaseName + ".schema.json"

		if e = dvc.ImportSchema(schemaFile); e != nil {
			go log.Printf("Error: %s", e.Error())
			os.Exit(1)
		}

		log.Printf("Schema `%s`.`%s` imported to %s.", dvc.Config.Host, dvc.Config.DatabaseName, schemaFile)
		os.Exit(0)

	case "compare":

		sql := ""

		schemaFile := dvc.Config.DatabaseName + ".schema.json"

		if sql, e = dvc.CompareSchema(schemaFile); e != nil {
			log.Printf("ERROR: %s", e.Error())
			os.Exit(1)
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
						log.Printf("ERROR: %s", e.Error())
						os.Exit(1)
					}

					os.Exit(0)

				} else {
					log.Printf("Missing path to write sql to")
					os.Exit(1)
				}
			case "apply":

				e = dvc.ApplyChangeset(sql)
				if e != nil {
					log.Printf("ERROR: %s", e.Error())
					os.Exit(1)
				}

			default:
				fmt.Printf("Unknown argument: `%s`\n", args[2])
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
		log.Printf("Unknown command: `%s`", cmd)
		os.Exit(1)
	}

}
