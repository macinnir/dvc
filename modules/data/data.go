package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/macinnir/dvc/lib"
)

// APPLY: Apply data in files to database
// 1. Verify that the data folder exists
// 2. Find the data files (named by tables)
// 3. Validate those files
// 4. Connect to the database
// 5. Verify that those tables exist in the database. If not, skip them.
// 6. TRUNCATE all data from those tables
// 7. Build SQL INSERT statements based on the data file
// 8. Run the SQL

// IMPORT: Import data to files from database
// 1. Verify that the data folder exists
// 2. Find the data files (named by tables)
// 3. Validate those files
// 4. Connect to the database
// 5. Confirm whether or not those tables exist in the database. If not, skip them.
// 6. Query all data from the table and convert to data file contents
// 7. Write content to data file

// NewData creates a new Data instance
func NewData(config *lib.Config, options lib.Options, database *lib.Database, connector lib.IConnector) (data *Data, e error) {
	data = &Data{
		config:    config,
		options:   options,
		database:  database,
		connector: connector,
	}
	return
}

type Data struct {
	config    *lib.Config
	options   lib.Options
	connector lib.IConnector
	database  *lib.Database
}

// func (d *Data) Apply() {

// 	x := lib.NewExecutor(d.config, d.connector)
// 	x.RunSQL(sql)

// }
// Import imports table data from database to data files
func (d *Data) Import(args []string) {

	files := []string{}

	if len(args) > 0 {
		files = []string{fmt.Sprintf("%s.json", args[0])}
	} else {
		files = d.getFiles()
	}

	if len(files) == 0 {
		fmt.Println("No files listed. Try importing a single table via `dvc import [table name]`")
		return
	}

	server := lib.NewExecutor(d.config, d.connector).Connect()

	for k := range files {

		// Remove .json file extension
		tableName := files[k][0 : len(files[k])-5]

		fmt.Printf("Looking for %s\n", tableName)

		if _, ok := d.database.Tables[tableName]; !ok {
			fmt.Printf("Cannot import data from table %s -- does not exist in database.\n", tableName)
			continue
		}

		fileData := []map[string]interface{}{}

		table := d.database.Tables[tableName]

		enums := d.connector.FetchEnum(server, tableName)

		for k := range enums {

			row := map[string]interface{}{}
			for j, l := range enums[k] {
				// fmt.Printf("%s %v\n", j, l)

				switch table.Columns[j].DataType {
				case "int":
					row[j] = *(l.(*uint32))
					// fmt.Println(j, *(l.(*uint32)))
				case "bigint":
					row[j] = *(l.(*uint64))
					// fmt.Println(j, *(l.(*uint64)))
				case "tinyint":
					if table.Columns[j].IsUnsigned {
						row[j] = *(l.(*uint8))
						// fmt.Println(j, *(l.(*uint8)))
					} else {
						row[j] = *(l.(*int8))
						// fmt.Println(j, *(l.(*int8)))
					}
				// case "int":
				// 	fmt.Println(j, l.(int8))
				// case "float64":
				// 	fmt.Println(j, l.(float64))
				case "varchar", "char", "text":
					row[j] = *(l.(*string))
					// fmt.Println(j, *(l.(*string)))
				default:
					fmt.Println(j, "???")
				}

				// val := reflect.ValueOf(l)
				// fmt.Println(j, val)
			}

			fileData = append(fileData, row)
		}

		filePath := fmt.Sprintf("meta/data/%s", files[k])
		// if !lib.FileExists(filePath) {
		// 	fmt.Printf("%s does not exist!\n", filePath)
		// } else {
		// 	fmt.Printf("%s DOES exist\n", filePath)
		// }
		jsonContent, _ := json.MarshalIndent(fileData, "", "    ")

		ioutil.WriteFile(filePath, []byte(jsonContent), 0644)
	}

}

func (d *Data) getFiles() []string {

	fmt.Println("Getting files?")
	dataPath := "meta/data"

	lib.EnsureDir(dataPath)

	files, e := lib.FetchNonDirFileNames(dataPath)
	if e != nil {
		panic(e)
	}

	return files

}

// Apply applies the data files to the database
func (d *Data) Apply(args []string) {

	files := []string{}

	if len(args) > 0 {
		files = []string{fmt.Sprintf("%s.json", args[0])}
	} else {
		files = d.getFiles()
	}

	fmt.Println("Files: ", len(files))

	x := lib.NewExecutor(d.config, d.connector)

	queries := []string{}
	for k := range files {
		queries = append(queries, d.genInsertQueriesForTable(files[k])...)
	}

	allQueryString := strings.Join(queries, "\n")
	fmt.Println(allQueryString)

	x.RunSQL(allQueryString)
	// for k := range queries {
	// 	fmt.Println(queries[k])
	// }
}

func (d *Data) genInsertQueriesForTable(file string) []string {

	filePath := fmt.Sprintf("meta/data/%s", file)
	tableName := file[0 : len(file)-5]

	if _, ok := d.database.Tables[tableName]; !ok {
		fmt.Printf("Cannot apply data to table %s -- does not exist in database.\n", tableName)
		return []string{}
	}

	table := d.database.Tables[tableName]

	if !lib.FileExists(filePath) {
		fmt.Printf("Cannot apply data to table %s -- no data file exists. Run `dvc data import %s` to store data for this table.\n", tableName, tableName)
		return []string{}
	}

	tableData := []map[string]interface{}{}

	content, _ := ioutil.ReadFile(filePath)
	json.Unmarshal(content, &tableData)

	cols := []string{}
	for l := range tableData[0] {
		cols = append(cols, l)
	}

	sort.Strings(cols)
	queries := []string{fmt.Sprintf("DELETE FROM `%s`;", tableName)}

	for j := range tableData {

		fields := []string{}
		values := []string{}
		n := 0

		for i := range cols {

			l := cols[i]
			m := tableData[j][l]

			fields = append(fields, fmt.Sprintf("`%s`", l))

			// t := table.Columns[l].FmtType
			val := ""
			switch table.Columns[l].DataType {
			case "int", "bigint", "tinyint":
				val = fmt.Sprintf("%d", int64(m.(float64)))
			case "char", "varchar", "text":
				val = fmt.Sprintf("'%s'", m)
			case "decimal":
				val = fmt.Sprintf("%f", m)
			}

			values = append(values, val)

			n++
		}

		query := fmt.Sprintf("REPLACE INTO `%s` ( %s ) VALUES ( %s );", tableName, strings.Join(fields, ","), strings.Join(values, ","))

		queries = append(queries, query)

	}
	return queries
}

// Remove removes the file data locally
func (d *Data) Remove(args []string) {

	if len(args) == 0 {
		fmt.Println("Usage: dvc data rm [table name]")
		return
	}

	tableName := args[0]

	filePath := fmt.Sprintf("meta/data/%s.json", tableName)

	if !lib.FileExists(filePath) {
		fmt.Printf("Data file `%s` does not exist.\n", filePath)
	}

	os.Remove(filePath)
}
