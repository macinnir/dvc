package selectcmd

import "fmt"

func Help() {
	fmt.Println(`
	select [tableName] 
	
		Selects rows from the table [tableName]
	`)
}
