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
	// TODO safeMode := false

	filteredArgs := []string{}

	for k := range args {
		switch args[k] {
		case "-u", "--summarize":
			summarize = true
		case "-a", "--apply":
			apply = true
		// case "-s", "--safe-mode":
		// 	safeMode = true
		default:
			filteredArgs = append(filteredArgs, args[k])
		}
	}

	if len(filteredArgs) > 0 {
		return CompareSingle(config, filteredArgs, summarize, apply)
	}

	return CompareAll(config, summarize, apply)
}

func CompareSingle(config *lib.Config, args []string, summarize, apply bool) error {

	localSchema := ""
	remoteConnectionName := ""
	if len(args) > 0 {
		localSchema = args[0]
	}

	if len(args) > 1 {
		remoteConnectionName = args[1]
	}

	if len(localSchema) == 0 {
		return fmt.Errorf("Missing local schema name")
	}

	if len(remoteConnectionName) == 0 {
		return fmt.Errorf("Missing remote connection name")
	}

	var e error
	var remoteSchema *schema.Schema
	remoteSchema, e = importer.FetchSchema(config, localSchema, remoteConnectionName)
	if e != nil {
		return fmt.Errorf("Error importing schemas: %w", e)
	}

	var localSchemaList *schema.SchemaList
	localSchemaList, _ = schema.LoadLocalSchemas()

	var targetLocalSchema *schema.Schema
	for k := range localSchemaList.Schemas {
		if localSchemaList.Schemas[k].Name == localSchema {
			targetLocalSchema = localSchemaList.Schemas[k]
		}
	}

	if targetLocalSchema == nil {
		return fmt.Errorf("Cannot find Target Local Schema `%s`", localSchema)
	}

	comparisons := compare.CompareSchemas(config, &schema.SchemaList{
		Schemas: []*schema.Schema{
			targetLocalSchema,
		},
	}, map[string]*schema.Schema{
		remoteConnectionName: remoteSchema,
	})

	if summarize {
		compare.PrintComparisonSummary(comparisons)
	}

	if !summarize && !apply {
		compare.PrintComparisons(comparisons)
	}

	if apply {

		configs := map[string]*lib.ConfigDatabase{}
		for k := range config.Databases {
			configs[config.Databases[k].Key] = config.Databases[k]
		}

		if e = applyChanges(
			configs,
			comparisons,
		); e != nil {
			fmt.Println("SQL ERROR: ", e.Error())
		}

	}

	return nil
}

func CompareAll(config *lib.Config, summarize, apply bool) error {

	var e error
	var remoteSchemas map[string]*schema.Schema
	remoteSchemas, e = importer.FetchAllSchemas(config)
	if e != nil {
		return fmt.Errorf("error importing schemas: %w", e)
	}

	var localSchemaList *schema.SchemaList
	localSchemaList, e = schema.LoadLocalSchemas()
	if e != nil {
		return fmt.Errorf("error loading local schemas: %w", e)
	}

	if e != nil {
		return fmt.Errorf("Error loading local schemas: %w", e)
	}

	comparisons := compare.CompareSchemas(config, localSchemaList, remoteSchemas)

	if summarize {
		compare.PrintComparisonSummary(comparisons)
	}

	if !summarize && !apply {
		compare.PrintComparisons(comparisons)
	}

	if apply {

		configs := map[string]*lib.ConfigDatabase{}
		for k := range config.Databases {
			configs[config.Databases[k].Key] = config.Databases[k]
		}

		if e = applyChanges(
			configs,
			comparisons,
		); e != nil {
			fmt.Println("SQL ERROR: ", e.Error())
		}

	}

	return nil
}

func applyChanges(
	configs map[string]*lib.ConfigDatabase,
	comparisons []*schema.SchemaComparison,
) error {

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

			// TODO if safeMode && changes[l].IsDestructive {
			// 	fmt.Println("IS_DESTRUCTIVE!!!!")
			// }

			change := strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(changes[l].SQL), "\n", ""), "\t", "")
			if len(change) > 80 {
				change = change[0:80] + "..."
			}
			fmt.Printf("Applying query %s.%d %s...%s\n", config.Key, l, changes[l].Type, change)
			_, e = server.Connection.Exec(changes[l].SQL)
			if e != nil {
				return fmt.Errorf("SQL ERROR on database %s (%s/%s):\n\n %s \n\n %w", config.Key, config.Host, config.Name, changes[l].SQL, e)
			}
		}
	}

	return nil
}
