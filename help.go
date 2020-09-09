package main

import "fmt"

func helpCommandNames() {
	cmds := []string{
		"add [table]",
		"compare",
		"export",
		"gen (dal (table)|dals|models|interfaces|routes)",
		"help",
		"import",
		"init",
		"insert",
		"ls",
		"refresh",
		"rm",
	}
	fmt.Println("Commands:")
	fmt.Println("----------------------------")
	for k := range cmds {
		fmt.Println("\t" + cmds[k])
	}
	fmt.Println("----------------------------")
}

func helpAdd() {
	fmt.Println(`
	add  [table]
		
		If the table name already exists, add a column to the existing table.
		Add an object to the database and then (by calling the import command) the local schema.
	`)
}

func helpCompare() {
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

func helpExport() {
	fmt.Println(`
	export 	Export the database to stdout
	`)
}

func helpGen() {
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

func helpImport() {
	fmt.Println(`
	import	Build a schema definition file based on the target database.
			This will overwrite any existing schema definition file.
	`)
}

func helpInsert() {
	fmt.Println(`
	insert Builds a SQL insert query based upon field by field input and executes it
	`)
}

func helpSelect() {
	fmt.Println(`
	select [tableName] 
	
		Selects rows from the table [tableName]
	`)
}

func helpInit() {
	fmt.Println(`
	init 	Initialize a dvc.toml configuration file in the CWD
	`)
}

func helpLs() {
	fmt.Println(`
	ls 		List all tables in the schema.

			[full table name] 	List information about an object in the database (e.g. columns of a table)

			[partial table name] List all tables with a name containing [partial table name]

			.[partial or full column name] List all columns with a name containing [partial or full column name]
	
	`)
}

func helpRefresh() {
	fmt.Println(`
	refresh Alias for running all of the following commands (in order):

		1. import
		2. gen models
		3. gen dals
		4. gen interfaces
		5. gen routes
	`)
}

func helpRm() {
	fmt.Println(`
	rm Removes an object from the database schema (with prompts)
	`)
}
