package importcmd

import "fmt"

func Help() {
	fmt.Println(`
	import	[[schema_name]] [[connection_name]] 
	
		Build a schema definition file based on the target database.
		This will overwrite any existing schema definition file.

		If schema_name and connection_name arguments are provided, the import process will only apply to the that schema
		and will replace it inside of whatever file it exists (core or app) without affecting any other schemas. 
	`)
}
