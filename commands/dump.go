package commands

import (
	"fmt"
	"strings"

	"github.com/macinnir/dvc/lib"
)

func dbFieldCleanString(val string) string {

	return strings.Replace(
		strings.Replace(
			val,
			"\\", "\\\\", -1,
		),
		"'", `\'`, -1,
	)
}

// Dump produces sql insert statements for the database
func (c *Cmd) Dump(args []string) {

	if len(args) > 0 && args[0] == "help" {
		helpDump()
		return
	}

	database := c.loadDatabase()

	// reader := bufio.NewReader(os.Stdin)

	// tableName := ""

	// if len(args) > 0 {
	// 	tableName = args[0]
	// } else {
	// 	tableName = lib.ReadCliInput(reader, "Table:")
	// }

	tables := database.ToSortedTables()
	connector, _ := connectorFactory(c.Config.DatabaseType, c.Config)
	server := lib.NewExecutor(c.Config, connector).Connect()

	fmt.Println("-- DVC Table Dump")
	fmt.Printf("-- Host: %s\t Database: %s\n", c.Config.Connection.Host, c.Config.Connection.DatabaseName)

	insertBatchMax := 100

	for k := range tables {

		insertBatchCounter := 0
		table := tables[k]

		fmt.Println("")
		fmt.Println("--")
		fmt.Printf("-- Dumping data for table `%s`\n", table.Name)
		fmt.Println("--")
		fmt.Println("")

		// fmt.Printf("DROP TABLE IF EXISTS `%s`;\n", table.Name)

		tableData := connector.FetchEnum(server, table.Name)

		cols := table.ToSortedColumns()

		colNames := []string{}
		for j := range cols {
			colNames = append(colNames, fmt.Sprintf("`%s`", cols[j].Name))
		}

		insertStart := fmt.Sprintf("\nINSERT INTO `%s` (%s) VALUES ", table.Name, strings.Join(colNames, ","))

		if len(tableData) > 0 {
			fmt.Printf("LOCK TABLES `%s` WRITE;\n", table.Name)
			fmt.Print(insertStart)
			for rowNum := range tableData {

				if insertBatchCounter > insertBatchMax {
					fmt.Println("; \n", insertStart)
					insertBatchCounter = 0
				}

				if insertBatchCounter > 0 {
					fmt.Print(",")
				}

				insertBatchCounter++

				// fields := []string{}
				values := []string{}
				n := 0

				for colNum := range cols {

					colName := cols[colNum].Name

					// l := table.Columns[colName]
					// m := tableData[rowNum][colName]

					// fields = append(fields, fmt.Sprintf("`%s`", l.Name))

					// fmt.Printf("Parsing value %s.%s => %s\n", table.Name, l.Name, table.Columns[colName].DataType)

					// t := table.Columns[l].FmtType
					val := ""

					switch table.Columns[colName].DataType {
					case "char", "varchar", "text", "date", "datetime", "enum":

						val = fmt.Sprintf("'%s'", dbFieldCleanString(fmt.Sprintf("%s", tableData[rowNum][colName])))
					default:
						val = fmt.Sprintf("%v", tableData[rowNum][colName])
					}
					// switch table.Columns[colName].DataType {
					// case "int":
					// 	if table.Columns[colName].IsUnsigned {
					// 		val = fmt.Sprintf("%d", int64(*(m.(*uint32))))
					// 	} else {
					// 		val = fmt.Sprintf("%d", int64(*(m.(*int32))))
					// 	}
					// case "bigint":
					// 	val = fmt.Sprintf("%d", int64(*(m.(*uint64))))
					// case "tinyint":
					// 	if table.Columns[colName].IsUnsigned {
					// 		val = fmt.Sprintf("%d", int64(*(m.(*uint8))))
					// 	} else {
					// 		val = fmt.Sprintf("%d", int64(*(m.(*int8))))
					// 	}
					// case "char", "varchar", "text":
					// 	val = fmt.Sprintf("'%s'", *(m.(*string)))
					// case "decimal":
					// 	val = fmt.Sprintf("%f", float64(*(m.(*float64))))
					// }

					values = append(values, val)

					n++
				}

				query := fmt.Sprintf("(%s)", strings.Join(values, ","))

				// queries = append(queries, query)
				fmt.Print(query)

				// }

				// sql := fmt.Sprintf("INSERT INTO `%s` (\n", table.Name)

				// columnNames := []string{}
				// values := []string{}

				// for k := range columns {

				// 	if columns[k].ColumnKey == "PRI" {
				// 		continue
				// 	}

				// 	if columns[k].Name == "IsDeleted" {
				// 		continue
				// 	}

				// 	columnNames = append(columnNames, fmt.Sprintf("`%s`", columns[k].Name))

				// 	value := "?"
				// 	if columns[k].Name == "DateCreated" {
				// 		value = fmt.Sprintf("%d", time.Now().UnixNano()/1000000)
				// 	} else {
				// 		// value = lib.ReadCliInput(reader, columns[k].Name+" ("+columns[k].DataType+"):")
				// 	}

				// 	if lib.IsString(columns[k]) {
				// 		value = "'" + value + "'"
				// 	} else {
				// 		if len(value) == 0 {
				// 			value = "0"
				// 		}
				// 	}

				// 	values = append(values, value)
				// }

				// sql += "\t" + strings.Join(columnNames, ",\n\t")
				// sql += "\n) VALUES (\n"
				// sql += "\t" + strings.Join(values, ",\n\t")
				// sql += "\n)\n"
				// fmt.Println(sql)
			}
			fmt.Print(";\nUNLOCK TABLES;\n\n")
		}
	}
	// doInsertYN := lib.ReadCliInput(reader, "Run above SQL (Y/n")

	// if doInsertYN != "Y" {
	// 	return
	// }

	// x := lib.NewExecutor(c.Config, connector)
	// x.RunSQL(sql)
}

func helpDump() {
	fmt.Println(`
	dumps Prints SQL insert statements for the database 
	`)
}
