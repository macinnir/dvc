package transfer

import "fmt"

func Help() {
	fmt.Println(`
	transfer [fromDatabase] [toDatabase] [[table_name]]
	
		Transfer schema and data from one database to another (on the same server).

		If a table_name is provided, only product sql (or run) for that specific table. 

		[-r|--run] 

		Default behavior is to print the sql statements, without performing any further action. 

			-r, --run 		Run the sql. 


	`)

}
