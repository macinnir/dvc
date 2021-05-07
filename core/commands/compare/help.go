package compare

import "fmt"

func Help() {
	fmt.Println(`
	compare 
		
		Compare two schemas and output the difference.

		[-r|--reverse] [ ( write <path> | apply ) ]

		Default behavior (no arguments) is to compare local schema as authority against
		remote database as target and write the resulting sql to stdout.

			-r, --reverse 	Switches the roles of the schemas. The remote database becomes the authority
							and the local schema the target for updating.

							Use this option when attempting to generate sql alter statements against a database that
							matches the structure of your local schema, in order to make it match a database with the structure
							of the remote.

			write			After performing the comparison, the resulting sequel statements will be written to a filepath <path> (required).

							Example: dvc compare write path/to/changeset.sql

			apply 			After performing the comparison, apply the the resulting sql statements directly to the target database.

							E.g. dvc compare apply
	`)

}
