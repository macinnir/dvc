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

	start = time.Now()
	c.Gen([]string{"dals"})
	fmt.Printf("DALs: %f seconds\n", time.Since(start).Seconds())

	start = time.Now()
	c.Gen([]string{"interfaces"})
	fmt.Printf("Interfaces: %f seconds\n", time.Since(start).Seconds())

	start = time.Now()
	c.Gen([]string{"routes"})
	fmt.Printf("Routes: %f seconds\n", time.Since(start).Seconds())

	fmt.Printf("Total: %f seconds\n", time.Since(totalTime).Seconds())
}
