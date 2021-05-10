package compare

import (
	"fmt"
	"strings"

	"github.com/macinnir/dvc/core/connectors"
	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/compare"
	"github.com/macinnir/dvc/core/lib/executor"
	"github.com/macinnir/dvc/core/lib/importer"
	"github.com/macinnir/dvc/core/lib/schema"
	"go.uber.org/zap"
)

const CommandName = "compare"

// Compare handles the `compare` command
func Cmd(log *zap.Logger, config *lib.Config, args []string) error {

	summarize := false
	apply := false
	safeMode := false

	for k := range args {
		switch args[k] {
		case "-u", "--summarize":
			summarize = true
		case "-a", "--apply":
			apply = true
		case "-s", "--safe-mode":
			safeMode = true
		}
	}

	var e error
	var remoteSchemas map[string]*schema.Schema
	remoteSchemas, e = importer.FetchAllSchemas(config)
	if e != nil {
		log.Error("Error importing schemas", zap.Error(e))
		return e
	}

	var localSchemaList *schema.SchemaList
	localSchemaList, _ = schema.LoadLocalSchemas()

	comparisons := compare.CompareSchemas(config, localSchemaList, remoteSchemas)

	if summarize {
		compare.PrintComparisonSummary(comparisons)
	} else {
		compare.PrintComparisons(comparisons)
	}

	if apply {

		configs := map[string]*lib.ConfigDatabase{}
		for k := range config.Databases {
			configs[config.Databases[k].Key] = config.Databases[k]
		}

		for k := range comparisons {

			config := configs[comparisons[k].DatabaseKey]
			connector, e := connectors.DBConnectorFactory(config)
			if e != nil {
				panic(e)
			}
			server := executor.NewExecutor(
				config,
				connector,
			).Connect()

			changes := comparisons[k].Changes

			for l := range changes {

				if safeMode && changes[l].IsDestructive {
					fmt.Println("IS_DESTRUCTIVE!!!!")
				}

				change := strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(changes[l].SQL), "\n", ""), "\t", "")
				if len(change) > 40 {
					change = change[0:40] + "..."
				}
				fmt.Printf("Applying query %s...%s\n", changes[l].Type, change)
				_, e = server.Connection.Exec(changes[l].SQL)
				if e != nil {
					panic(e)
				}
			}
		}
	}

	return nil
	// cmd := "print"
	// sql := ""
	// outfile := ""

	// lib.Debugf("Args: %v", c.Options, args)
	// if len(args) > 0 {

	// 	for len(args) > 0 {

	// 		switch args[0] {
	// 		case "-r", "--reverse":
	// 			c.Options |= lib.OptReverse
	// 		case "-u", "--summary":
	// 			c.Options |= lib.OptSummary
	// 		case "print":
	// 			cmd = "print"
	// 		case "apply":
	// 			cmd = "apply"
	// 		default:

	// 			if len(args[0]) > len("-o=") && args[0][0:len("-o=")] == "-o=" {
	// 				outfile = args[0][len("-o="):]
	// 				if len(outfile) == 0 {
	// 					lib.Error("Outfile argument cannot be empty", c.Options)
	// 					os.Exit(1)
	// 				}
	// 				cmd = "write"
	// 			} else if args[0][0] == '-' {
	// 				lib.Errorf("Unrecognized option '%s'. Try the --help option for more information\n", c.Options, args[0])
	// 				os.Exit(1)
	// 				// c.errLog.Fatalf("Unrecognized option '%s'. Try the --help option for more information\n", arg)
	// 			}

	// 			// Check if outfile argument is non-empty

	// 			break
	// 		}
	// 		args = args[1:]
	// 	}

	// }

	// Do the comparison
	// TODO pass all options (e.g. verbose)
	// TODO -reverse | -r as option
	// reverse := false
	// if sql, e = cmp.CompareSchema(config.Databases[0], reverse); e != nil {
	// 	log.Sugar().Errorf(e.Error())
	// 	return e
	// }

	// if len(sql) == 0 {
	// 	log.Info("No changes found")
	// 	return nil
	// }

	// switch cmd {
	// case "apply":
	// 	writeSQLToLog(sql)
	// 	e = cmp.ApplyChangeset(sql)
	// 	if e != nil {
	// 		lib.Error(e.Error(), c.Options)
	// 		os.Exit(1)
	// 	}

	// case "print":
	// 	// Print to stdout
	// 	fmt.Printf("%s", sql)
	// default:
	// 	lib.Errorf("Unknown argument: `%s`", c.Options, cmd)
	// 	os.Exit(1)
	// }
}
