package commands

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/macinnir/dvc/core/commands/add"
	"github.com/macinnir/dvc/core/commands/clone"
	"github.com/macinnir/dvc/core/commands/compare"
	"github.com/macinnir/dvc/core/commands/data"
	"github.com/macinnir/dvc/core/commands/dump"
	"github.com/macinnir/dvc/core/commands/export"
	"github.com/macinnir/dvc/core/commands/gen"
	"github.com/macinnir/dvc/core/commands/importcmd"
	"github.com/macinnir/dvc/core/commands/insert"
	"github.com/macinnir/dvc/core/commands/inspect"
	"github.com/macinnir/dvc/core/commands/ls"
	"github.com/macinnir/dvc/core/commands/refresh"
	"github.com/macinnir/dvc/core/commands/rm"
	"github.com/macinnir/dvc/core/commands/selectcmd"
	"github.com/macinnir/dvc/core/commands/test"

	"github.com/macinnir/dvc/core/lib"
)

var (
	ErrNoCommandSpecified = errors.New("No command specified")
	ErrInvalidCommand     = errors.New("Invalid command")
)

// Command is a type that represents the possible commands passed in at run time
type Command func(log *zap.Logger, config *lib.Config, args []string) (e error)
type Help func()

var commands map[string]Command
var helps map[string]Help

func makeCommands() {
	commands = map[string]Command{
		add.CommandName:       add.Cmd,
		clone.CommandName:     clone.Cmd,
		compare.CommandName:   compare.Cmd,
		data.CommandName:      data.Cmd,
		dump.CommandName:      dump.Cmd,
		export.CommandName:    export.Cmd,
		gen.CommandName:       gen.Cmd,
		HelpCommandName:       PrintHelp,
		importcmd.CommandName: importcmd.Cmd,
		insert.CommandName:    insert.Cmd,
		inspect.CommandName:   inspect.Cmd,
		ls.CommandName:        ls.Cmd,
		refresh.CommandName:   refresh.Cmd,
		rm.CommandName:        rm.Cmd,
		selectcmd.CommandName: selectcmd.Cmd,
		test.CommandName:      test.Cmd,
	}
}

func makeHelp() {
	helps = map[string]Help{
		add.CommandName:       add.Help,
		clone.CommandName:     clone.Help,
		compare.CommandName:   compare.Help,
		data.CommandName:      data.Help,
		dump.CommandName:      dump.Help,
		export.CommandName:    export.Help,
		gen.CommandName:       gen.Help,
		importcmd.CommandName: importcmd.Help,
		insert.CommandName:    insert.Help,
		inspect.CommandName:   inspect.Help,
		ls.CommandName:        ls.Help,
		refresh.CommandName:   refresh.Help,
		rm.CommandName:        rm.Help,
		selectcmd.CommandName: selectcmd.Help,
		test.CommandName:      test.Help,
	}
}

// Cmd is a container for handling commands
type Cmd struct {
	Options   lib.Options
	cmd       string
	Config    *lib.Config
	newModels map[string]string
}

// Run is the main function that handles commands arguments
// and routes them to their correct options and functions
func (c *Cmd) Run(args []string) error {

	// Ensure that the meta directory exists
	lib.EnsureDir(lib.MetaDirectory)

	// Logger
	logger, e := zap.Config{
		Level:    zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding: "console",
	}.Build()

	if e != nil {
		return e
	}

	// Config
	var config *lib.Config
	config, e = lib.LoadConfig()
	logger.Info("Loading config!")
	if e != nil {
		return e
	}

	// TODO flags
	// Remove initial file path
	args = args[1:]

	if len(args) == 0 {
		return ErrNoCommandSpecified
	}

	command := args[0]

	args = args[1:]

	if len(args) > 0 && args[0] == "help" {
		makeHelp()
		if _, ok := helps[command]; !ok {
			return ErrInvalidCommand
		}

		helps[command]()
		return nil
	}

	makeCommands()

	if _, ok := commands[command]; !ok {
		return ErrInvalidCommand
	}

	e = commands[command](logger, config, args)

	return e

	// var cmd Command
	// var arg string
	// args := []string{}

	// for len(inputArgs) > 0 {

	// 	arg = inputArgs[0]

	// 	switch arg {
	// 	case "-h", "--help":
	// 		cmd = CommandHelp
	// 	case "-v", "--verbose":
	// 		c.Options &^= lib.OptLogDebug | lib.OptSilent
	// 		c.Options |= lib.OptLogInfo
	// 	case "-vv", "--debug":
	// 		c.Options &^= lib.OptLogInfo | lib.OptSilent
	// 		c.Options |= lib.OptLogDebug
	// 	case "-s", "--silent":
	// 		c.Options &^= lib.OptLogDebug | lib.OptLogInfo
	// 		c.Options |= lib.OptSilent
	// 	case "-f", "--force":
	// 		c.Options |= lib.OptForce
	// 	case "-c", "--clean":
	// 		c.Options |= lib.OptClean
	// 	default:

	// 		lib.Debug(fmt.Sprintf("Arg: %s", arg), c.Options)

	// 		if len(arg) > 0 && arg[0] == '-' {
	// 			lib.Errorf("Unrecognized option '%s'. Try the --help option for more information\n", c.Options, arg)
	// 			// c.errLog.Fatalf("Unrecognized option '%s'. Try the --help option for more information\n", arg)
	// 			return nil
	// 		}

	// 		if isCommand(arg) && len(cmd) == 0 {
	// 			cmd = Command(arg)
	// 		} else {
	// 			args = append(args, arg)
	// 		}
	// 	}

	// 	inputArgs = inputArgs[1:]
	// }

	// if len(cmd) == 0 {
	// 	lib.Error("No command provided", c.Options)
	// 	return
	// }
	// lib.Debugf("cmd: %s, %v\n", c.Options, cmd, args)

	// if cmd != CommandInit {
	// 	var e error
}

// SearchType is the type of search
type SearchType int

const (
	// SearchTypeWildcard is a wildcard
	SearchTypeWildcard SearchType = iota
	// SearchTypeStartingWildcard is a wildcard at the start (e.g. `*foo`)
	SearchTypeStartingWildcard
	// SearchTypeEndingWildcard is a wildcard at the end (e.g. `foo*`)
	SearchTypeEndingWildcard
	// SearchTypeBoth is a wildcard on both the start and the end (e.g. `*foo*`)
	SearchTypeBoth
)
