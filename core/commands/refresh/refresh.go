package refresh

import (
	"fmt"
	"time"

	"github.com/macinnir/dvc/core/commands/gen"
	"github.com/macinnir/dvc/core/commands/importcmd"
	"github.com/macinnir/dvc/core/lib"
	"go.uber.org/zap"
)

const CommandName = "refresh"

// Refresh is the refresh command
func Cmd(logger *zap.Logger, config *lib.Config, args []string) error {

	var genArgs = []string{"all", "-c"}
	for k := range args {
		if args[k] == "-f" {
			genArgs = append(genArgs, "-f")
		}
	}
	var start = time.Now()
	importcmd.Cmd(logger, config, []string{})
	gen.Cmd(logger, config, genArgs)
	fmt.Printf("Finished in %f seconds\n", time.Since(start).Seconds())

	return nil

	// totalTime := time.Now()

	// // Import
	// start := time.Now()
	// c.Import(args)
	// fmt.Printf("Import: %f seconds\n", time.Since(start).Seconds())

	// // Gen Models
	// start = time.Now()
	// c.Gen([]string{"models"})
	// fmt.Printf("Models: %f seconds\n", time.Since(start).Seconds())

	// // Gen DALs
	// start = time.Now()
	// c.Gen([]string{"dals"})
	// fmt.Printf("DALs: %f seconds\n", time.Since(start).Seconds())

	// // Gen Interfaces
	// start = time.Now()
	// c.Gen([]string{"interfaces"})
	// fmt.Printf("Interfaces: %f seconds\n", time.Since(start).Seconds())

	// // Gen routes
	// start = time.Now()
	// c.Gen([]string{"routes"})
	// fmt.Printf("Routes: %f seconds\n", time.Since(start).Seconds())

	// fmt.Printf("Total: %f seconds\n", time.Since(totalTime).Seconds())
}

func helpRefresh() {
	fmt.Println(`
	refresh 
	
	FLAGS

		-c 	Clean 	Delete any previously generated files that have been orphaned due to a change in the Database Schema (e.g. a dropped table)
		-f  Force 	Force file regeneration, regardless of whether or not a change has been detected for that schema. 

	Alias for running all of the following commands (in order):

		dvc import 
		dvc gen models -c 
		dvc gen dals 
		dvc gen interfaces 
		dvc gen goperms
		dvc gen tsperms 
		dvc gen ts 
		dvc gen routes 
	`)
}
