package commands

import (
	"fmt"
	"sort"

	"github.com/macinnir/dvc/core/lib"
	"go.uber.org/zap"
)

const HelpCommandName = "help"

func PrintHelp(log *zap.Logger, config *lib.Config, args []string) error {

	commandList := []string{}
	for commandName := range commands {
		commandList = append(commandList, commandName)
	}

	sort.Strings(commandList)

	fmt.Println("Commands:")
	fmt.Println("----------------------------")
	for k := range commandList {
		fmt.Println("\t" + commandList[k])
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
