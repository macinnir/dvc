package version

import (
	"errors"
	"fmt"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/versioner"
	"go.uber.org/zap"
)

const CommandName = "version"

// Compare handles the `compare` command
func Cmd(log *zap.Logger, config *lib.Config, args []string) error {

	if len(args) == 0 {
		return errors.New("usage: dvc version [version] [[version type]]")
	}

	version := args[0]
	versionType := "patch"

	if len(args) > 1 {
		versionType = args[1]
	}

	nextVersion, e := versioner.NextVersion(version, versionType)

	if e != nil {
		return e
	}

	fmt.Print(nextVersion)
	return nil
}
