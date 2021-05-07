package add

import "fmt"

func Help() {
	fmt.Println(`
	add  [table]
		
		If the table name already exists, add a column to the existing table.
		Add an object to the database and then (by calling the import command) the local schema.
	`)
}
