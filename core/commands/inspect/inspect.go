package inspect

import (
	"github.com/macinnir/dvc/core/lib"
	"go.uber.org/zap"
)

const CommandName = "inspect"

// CommandInspect is the inspect command
func Cmd(logger *zap.Logger, config *lib.Config, args []string) error {
	return nil
}
