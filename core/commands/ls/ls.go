package ls

import (
	"fmt"
	"strings"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
	"go.uber.org/zap"
)

const CommandName = "ls"

// Ls lists database information
// TODO search fields
// TODO select from tables
// TODO show row counts in a table

// dvc ls - Show all tables in all schemas
// dvc ls [name] - Show all columns in table [name] is found
// dvc ls [partialName] - Show all tables with name containing [partialName]
// dvc ls .[fieldPartialName] - Show all columns with name containing [fieldPartialName]
func Cmd(logger *zap.Logger, config *lib.Config, args []string) error {

	localSchemaList, _ := schema.LoadLocalSchemas()

	var showList = false
	var showSchemaName = ""
	var showTableName = ""

	for len(args) > 0 {
		if args[0] == "-l" || args[0] == "--list" {
			showList = true
			args = args[1:]
			continue
		}

		if args[0] == "-s" || args[0] == "--schema" {
			showSchemaName = args[1]
			args = args[2:]
			continue
		}

		if args[0] == "-t" || args[0] == "--table" {
			showTableName = args[1]
			args = args[2:]
			continue
		}
	}

	if showList {
		if len(showTableName) > 0 {
			for k := range localSchemaList.Schemas {
				localSchema := localSchemaList.Schemas[k]
				for _, table := range localSchema.ToSortedTables() {
					if table.Name == showTableName {
						for _, column := range table.ToSortedColumns() {
							fmt.Printf("%s\n", column.Name)
						}
						return nil
					}
				}
			}
			return nil
		} else {
			for k := range localSchemaList.Schemas {
				if showSchemaName != "" && localSchemaList.Schemas[k].Name != showSchemaName {
					continue
				}
				localSchema := localSchemaList.Schemas[k]
				for _, table := range localSchema.ToSortedTables() {
					fmt.Printf("%s\n", table.Name)
				}
				// fmt.Println(localSchema.Name)
			}
			return nil
		}
	}

	if len(showTableName) == 0 {

		t := lib.NewCLITable([]string{"Schema", "Table", "Columns", "Collation", "Engine", "Row Format"})

		for k := range localSchemaList.Schemas {
			localSchema := localSchemaList.Schemas[k]
			if showSchemaName != "" && localSchema.Name != showSchemaName {
				continue
			}
			results := localSchema.ToSortedTables()

			for _, table := range results {
				t.Row()
				t.Col(localSchema.Name)
				t.Col(table.Name)
				t.Colf("%d", len(table.Columns))
				t.Col(table.Collation)
				t.Col(table.Engine)
				t.Col(table.RowFormat)
			}
		}
		fmt.Println(t.String())

		return nil
	}

	// ls [tableName] -  Show all column in the table [tableName]
	var table *schema.Table
	var tableName = strings.Trim(strings.ToLower(showTableName), " ")
	var schemaName = ""
	{
		for k := range localSchemaList.Schemas {

			localSchema := localSchemaList.Schemas[k]

			// findTable
			for k := range localSchema.Tables {
				if strings.ToLower(localSchema.Tables[k].Name) == tableName {
					table = localSchema.Tables[k]
					schemaName = localSchema.Name
					// Found
					break
				}
			}
		}
	}

	if table != nil {

		t := lib.NewCLITable([]string{"Schema", "Name", "Type", "MaxLength", "Null", "Default", "Extra", "Key"})

		fmt.Printf("Listing all columns for table `%s`\n", tableName)

		sortedColumns := table.ToSortedColumns()

		for _, col := range sortedColumns {

			t.Row()

			t.Col(schemaName)
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
		return nil
	}

	// ls [partialTableName] - Show all tables with name containing [partialTableName]
	t := lib.NewCLITable([]string{"Schema", "Table", "Columns", "Collation", "Engine", "Row Format"})

	for k := range localSchemaList.Schemas {
		localSchema := localSchemaList.Schemas[k]
		sortedTables := localSchema.ToSortedTables()
		tableSearch := strings.Trim(strings.ToLower(showTableName), " ")

		for _, table := range sortedTables {

			if !strings.Contains(strings.ToLower(table.Name), tableSearch) {
				continue
			}

			t.Row()
			t.Col(localSchema.Name)
			t.Col(table.Name)
			t.Colf("%d", len(table.Columns))
			t.Col(table.Collation)
			t.Col(table.Engine)
			t.Col(table.RowFormat)
		}
	}

	fmt.Println(t.String())
	return nil
}
