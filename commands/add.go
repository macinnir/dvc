package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/macinnir/dvc/lib"
)

// Add adds an object to the database
func (c *Cmd) Add(args []string) {

	if len(args) > 0 && args[0] == "help" {
		helpAdd()
		return
	}

	database := c.loadDatabase()

	reader := bufio.NewReader(os.Stdin)

	sqlParts := []string{}
	sql := ""

	tableName := ""

	if len(args) > 0 {
		tableName = args[0]
	} else {
		tableName = lib.ReadCliInput(reader, "Table Name:")
	}

	// Creating a new table
	if _, ok := database.Tables[tableName]; !ok {

		fmt.Printf("Create table `%s`\n", tableName)

		start := fmt.Sprintf("CREATE TABLE `%s` (", tableName)
		sql += start
		sqlParts = append(sqlParts, start)
		cols := []string{
			fmt.Sprintf("%sID INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT", tableName),
			fmt.Sprintf("IsDeleted TINYINT UNSIGNED NOT NULL DEFAULT 0"),
			fmt.Sprintf("DateCreated BIGINT UNSIGNED NOT NULL DEFAULT 0"),
			fmt.Sprintf("LastUpdated BIGINT UNSIGNED NOT NULL DEFAULT 0"),
		}

		for k := range cols {
			fmt.Printf("`%s`.Column(%d): %s\n", tableName, k+1, cols[k])
		}

		includeColumns := lib.ReadCliInput(reader, fmt.Sprintf("Include the above columns (Y/n)?"))

		n := 1
		if includeColumns == "Y" {
			n = len(cols) + 1
		} else {
			cols = []string{}
		}

		for {

			col := lib.ReadCliInput(reader, fmt.Sprintf("`%s`.Column(%d):", tableName, n))
			if len(col) == 0 {
				break
			}

			// Remove trailing space
			col = strings.Trim(col, " ")

			// Remove trailing comma if exists
			if col[len(col)-1:] == "," {
				col = col[0 : len(col)-1]
			}

			colParts := strings.Split(col, " ")

			// Validate the sql part
			if len(colParts) < 2 {
				fmt.Println("Invalid sql part. Try again...")
				continue
			}

			// Validate the datatype
			dataType := strings.ToLower(colParts[1])

			if strings.Contains(dataType, "(") {
				dataType = strings.Split(dataType, "(")[0]
			}

			if lib.IsValidSQLType(dataType) == false {
				fmt.Printf("Invalid SQL type: %s. Try gain...\n", dataType)
				continue
			}

			n++

			for k := range colParts {

				if k == 0 {
					continue
				}

				colParts[k] = strings.ToUpper(colParts[k])
			}

			col = "`" + colParts[0] + "` " + strings.Join(colParts[1:], " ")
			sqlParts = append(sqlParts, col)
			cols = append(cols, col)
		}

		sql += strings.Join(cols, ", ") + ")"

		sqlParts = append(sqlParts, ")")

	} else {

		col := lib.ReadCliInput(reader, fmt.Sprintf("`%s`.Column:", tableName))
		if len(col) == 0 {
			fmt.Println("Column cannot be empty")
			return
		}

		colParts := strings.Split(col, " ")

		if _, ok := database.Tables[tableName].Columns[colParts[0]]; ok {
			fmt.Printf("Column `%s`.`%s` already exists", tableName, colParts[0])
			return
		}

		for k := range colParts {

			if k == 0 {
				continue
			}

			colParts[k] = strings.ToUpper(colParts[k])
		}

		col = "`" + colParts[0] + "` " + strings.Join(colParts[1:], " ")
		sql = fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s", tableName, col)
		sqlParts = append(sqlParts, sql)
	}

	fmt.Print("\n------------------------------------------------------\n")
	fmt.Print("--------------------- REVIEW -------------------------\n")
	fmt.Print("------------------------------------------------------\n")

	fmt.Println(sql)
	// for k := range sqlParts {
	// 	// Not first or last line
	// 	if k != 0 && k != len(sqlParts)-1 {
	// 		fmt.Print("\t")
	// 	}
	// 	fmt.Printf("%s", sqlParts[k])
	// 	// Not first or last line
	// 	if k != 0 && k != len(sqlParts)-1 {
	// 		fmt.Print(",")
	// 	}

	// 	fmt.Print("\n")
	// }

	fmt.Print("\n------------------------------------------------------\n")

	if lib.ReadCliInput(reader, "Are you sure want to execute the above SQL (Y/n)?") == "Y" {
		connector, _ := connectorFactory(c.Config.DatabaseType, c.Config)
		x := lib.NewExecutor(c.Config, connector)
		x.RunSQL(sql)

		c.Refresh([]string{})
	}
}
