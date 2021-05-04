package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/macinnir/dvc/lib"
)

// Insert inserts data into the database
func (c *Cmd) Insert(args []string) {

	if len(args) > 0 && args[0] == "help" {
		helpInsert()
		return
	}

	database := c.loadDatabase()

	reader := bufio.NewReader(os.Stdin)

	tableName := ""

	if len(args) > 0 {
		tableName = args[0]
	} else {
		tableName = lib.ReadCliInput(reader, "Table:")
	}

	if len(tableName) == 0 {
		fmt.Println("No table specified")
		return
	}

	if _, ok := database.Tables[tableName]; !ok {
		fmt.Printf("Table `%s` not found.\n", tableName)
		return
	}

	table := database.Tables[tableName]

	columns := table.ToSortedColumns()

	sql := fmt.Sprintf("INSERT INTO `%s` (\n", tableName)

	columnNames := []string{}
	values := []string{}

	for k := range columns {

		if columns[k].ColumnKey == "PRI" {
			continue
		}

		if columns[k].Name == "IsDeleted" {
			continue
		}

		columnNames = append(columnNames, fmt.Sprintf("`%s`", columns[k].Name))

		value := "?"
		if columns[k].Name == "DateCreated" {
			value = fmt.Sprintf("%d", time.Now().UnixNano()/1000000)
		} else {
			value = lib.ReadCliInput(reader, columns[k].Name+" ("+columns[k].DataType+"):")
		}

		if lib.IsString(columns[k]) {
			value = "'" + value + "'"
		} else {
			if len(value) == 0 {
				value = "0"
			}
		}

		values = append(values, value)
	}

	sql += "\t" + strings.Join(columnNames, ",\n\t")
	sql += "\n) VALUES (\n"
	sql += "\t" + strings.Join(values, ",\n\t")
	sql += "\n)\n"

	fmt.Println(sql)

	doInsertYN := lib.ReadCliInput(reader, "Run above SQL (Y/n")

	if doInsertYN != "Y" {
		return
	}

	connector, _ := connectorFactory(c.Config.DatabaseType, c.Config)
	x := lib.NewExecutor(c.Config, connector)
	x.RunSQL(sql)
}

func helpInsert() {
	fmt.Println(`
	insert Builds a SQL insert query based upon field by field input and executes it
	`)
}
