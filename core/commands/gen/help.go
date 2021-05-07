package gen

import "fmt"

func Help() {
	fmt.Println(`
	gen 	Generate go code based on the local schema json file.
			Will fail if no imported schema file json file exists.
			Requires one (and only one) of the following sub-commands:

				dal [dalName] 	Generate dal [dalName]
				dals 			Generate dals
				interfaces 		Generate interfaces
				models 			Generate models.
				routes 			Generate routes
	`)
}
