package commands

import (
	"bufio"
	"fmt"
	"os"

	"github.com/macinnir/dvc/lib"
)

// Rm removes an object from the database
// dvc rm [table]
func (c *Cmd) Rm(args []string) {

	if len(args) > 0 && args[0] == "help" {
		helpRm()
		return
	}

	database := c.loadDatabase()

	reader := bufio.NewReader(os.Stdin)

	sql := ""

	tableName := ""
	if len(args) > 0 {
		tableName = args[0]
	} else {
		tableName = lib.ReadCliInput(reader, "Table Name:")
		if _, ok := database.Tables[tableName]; !ok {
			fmt.Printf("Table `%s` does not exist.", tableName)
			return
		}
	}

	tableOrColumn := lib.ReadCliInput(reader, fmt.Sprintf("Drop the (t)able `%s` or select a (c)olumn?", tableName))
	if tableOrColumn == "t" {
		sql = fmt.Sprintf("DROP TABLE `%s`", tableName)
	} else if tableOrColumn == "c" {
		columnName := lib.ReadCliInput(reader, fmt.Sprintf("`%s`.Column:", tableName))
		if _, ok := database.Tables[tableName].Columns[columnName]; !ok {
			fmt.Printf("Column `%s`.`%s` doesn't exist.", tableName, columnName)
			return
		}

		sql = fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`", tableName, columnName)
	} else {
		fmt.Println("Invalid entry")
		return
	}

	fmt.Print("\n------------------------------------------------------\n")
	fmt.Print("--------------------- REVIEW -------------------------\n")
	fmt.Print("------------------------------------------------------\n")

	fmt.Println(sql)

	fmt.Print("\n------------------------------------------------------\n")

	if lib.ReadCliInput(reader, "Are you sure want to execute the above SQL (Y/n)?") == "Y" {

		// Apply the change
		connector, _ := connectorFactory(c.Config.DatabaseType, c.Config)
		x := lib.NewExecutor(c.Config, connector)
		x.RunSQL(sql)

		// Import the schema
		c.Import([]string{})
	}
}
