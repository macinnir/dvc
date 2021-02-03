package commands

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/macinnir/dvc/lib"
)

// Compare handles the `compare` command
func (c *Cmd) Compare(args []string) {

	if len(args) > 0 && args[0] == "help" {
		helpCompare()
		return
	}

	var e error
	cmp := c.initCompare()

	cmd := "print"
	sql := ""
	schemaFile := c.Config.Connection.DatabaseName + ".schema.json"
	outfile := ""

	// lib.Debugf("Args: %v", c.Options, args)
	if len(args) > 0 {

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

				break
			}
			args = args[1:]
		}

	}

	cmp.SetOptions(c.Options)

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
