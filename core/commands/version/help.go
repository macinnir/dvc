package version

import "fmt"

func Help() {
	fmt.Println(`
	version [version] [[type]]
		
		Generate a version string from an input string

		Default behavior (no arguments) is to generate a new patch version from the input version. 

			type 			[major|minor|patch]
					
	`)

}
