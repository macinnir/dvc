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
func Cmd(logger *zap.Logger, config *lib.Config, args []string) error {

	localSchemaList, _ := schema.LoadLocalSchemas()
	// database, e := lib.LoadSchema(config.Databases[0])

	// if e != nil {
	// 	return e
	// }

	// Options
	// ls 							Show all Tables
	// ls [name] 					Show all Columns in table [name] is found
	// ls [partialName] 			Show all tables with name containing [partialName]
	// ls .[fieldPartialName] 		Show all columns with name containing [fieldPartialName]
	// ls - Show all tables
	// fmt.Println("args", args)
	if len(args) == 0 {

		t := lib.NewCLITable([]string{"Schema", "Table", "Columns", "Collation", "Engine", "Row Format"})

		for k := range localSchemaList.Schemas {

			localSchema := localSchemaList.Schemas[k]
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

	// ls .[fieldPartialName] - Show all columns with names containing [fieldPartialName]
	if strings.HasPrefix(args[0], ".") {

		if args[0] == "." {
			return nil
		}

		fieldSearch := strings.Trim(strings.ToLower(args[0][1:]), " ")

		fmt.Printf("Searching all fields for `%s`\n", fieldSearch)

		t := lib.NewCLITable([]string{"Schema", "Table", "Name", "Type", "MaxLength", "Null", "Default", "Extra", "Key"})

		for k := range localSchemaList.Schemas {

			localSchema := localSchemaList.Schemas[k]
			sortedTables := localSchema.ToSortedTables()

			columns := []*schema.ColumnWithTable{}

			for j := range sortedTables {
				for k := range sortedTables[j].Columns {
					if strings.Contains(strings.ToLower(sortedTables[j].Columns[k].Name), fieldSearch) {
						column := &schema.ColumnWithTable{
							Column:    sortedTables[j].Columns[k],
							TableName: sortedTables[j].Name,
						}
						columns = append(columns, column)
					}
				}
			}

			for _, col := range columns {

				t.Row()
				t.Col(localSchema.Name)
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
		}

		fmt.Println(t.String())

		return nil
	}

	// ls [tableName] -  Show all column in the table [tableName]
	var table *schema.Table
	tableName := strings.Trim(strings.ToLower(args[0]), " ")
	schemaName := ""

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
		tableSearch := strings.Trim(strings.ToLower(args[0]), " ")

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
