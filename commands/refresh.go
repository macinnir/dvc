package commands

import (
	"fmt"
	"time"
)

// Refresh is the refresh command
func (c *Cmd) Refresh(args []string) {

	if len(args) > 0 && args[0] == "help" {
		helpRefresh()
		return
	}

	totalTime := time.Now()

	// Import
	start := time.Now()
	c.Import(args)
	fmt.Printf("Import: %f seconds\n", time.Since(start).Seconds())

	// Gen Models
	start = time.Now()
	c.Gen([]string{"models"})
	fmt.Printf("Models: %f seconds\n", time.Since(start).Seconds())

	// Gen DALs
	start = time.Now()
	c.Gen([]string{"dals"})
	fmt.Printf("DALs: %f seconds\n", time.Since(start).Seconds())

	// Gen Interfaces
	start = time.Now()
	c.Gen([]string{"interfaces"})
	fmt.Printf("Interfaces: %f seconds\n", time.Since(start).Seconds())

	// Gen routes
	start = time.Now()
	c.Gen([]string{"routes"})
	fmt.Printf("Routes: %f seconds\n", time.Since(start).Seconds())

	fmt.Printf("Total: %f seconds\n", time.Since(totalTime).Seconds())
}

func helpRefresh() {
	fmt.Println(`
	refresh 
	
	FLAGS

		-c 	Clean 	Delete any previously generated files that have been orphaned due to a change in the Database Schema (e.g. a dropped table)
		-f  Force 	Force file regeneration, regardless of whether or not a change has been detected for that schema. 

	Alias for running all of the following commands (in order):

		1. import
		2. gen models
		3. gen dals
		4. gen interfaces
		5. gen routes
	`)
}
