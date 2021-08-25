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
	database := localSchemaList.Schemas[0]
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
	fmt.Println("args", args)
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
		return nil
	}

	// ls .[fieldPartialName] - Show all columns with names containing [fieldPartialName]
	if strings.HasPrefix(args[0], ".") {

		if args[0] == "." {
			return nil
		}

		fieldSearch := strings.Trim(strings.ToLower(args[0][1:]), " ")

		fmt.Printf("Searching all fields for `%s`\n", fieldSearch)

		sortedTables := database.ToSortedTables()

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

		return nil
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
		return nil
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
	return nil
}
