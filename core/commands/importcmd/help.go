package importcmd

import "fmt"

func Help() {
	fmt.Println(`
	import	Build a schema definition file based on the target database.
			This will overwrite any existing schema definition file.
	`)
}
