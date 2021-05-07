package commands

import (
	"fmt"

	"github.com/macinnir/dvc/core/lib"
	"go.uber.org/zap"
)

const HelpCommandName = "help"

func PrintHelp(log *zap.Logger, config *lib.Config, args []string) error {
	fmt.Println("Commands:")
	fmt.Println("----------------------------")
	for commandName := range commands {
		fmt.Println("\t" + commandName)
	}
	fmt.Println("----------------------------")
	return nil
}

// PrintHelp prints help information
// func (c *Cmd) PrintHelp(args []string) {
// 	help := `usage: dvc [OPTIONS] [COMMAND] [ARGS]

// OPTIONS:

// 	-h, --help 		Show help
// 	-v, --verbose 	Show verbose logging
// 	-vv, --debug 	Show very verbose (debug) logging
// 	-s, --silent 	Disable all output

// COMMANDS:

// 	help 	This output

// `
// 	fmt.Printf(help)
// }
