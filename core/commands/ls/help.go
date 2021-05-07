package ls

import "fmt"

func Help() {
	fmt.Println(`
		ls 		List all tables in the schema.
	
				[full table name] 	List information about an object in the database (e.g. columns of a table)
	
				[partial table name] List all tables with a name containing [partial table name]
	
				.[partial or full column name] List all columns with a name containing [partial or full column name]
		
		`)
}
