package data

import "fmt"

func Help() {
	fmt.Println(`
	data [import|apply|rm] [[table]]

		import - Import data from database tables to static files 

				Without any arguments, import looks for any files in the meta/data and imports data from the tables with the same names as the file 

				import [table] imports data for a single table, creating the table file if it does not exist. 

		apply - Apply data to the database from a static file 

				Without any arguments, import looks for any files in the meta/data folder and applies the contents of those files to the database tables with the same name as the files respectively

				import [table] applies data for a single table (if it exists) to the database. 
		
		rm [table] - Remove a data file 

		help - This output 
	
	`)
}
