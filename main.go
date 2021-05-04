package main

import (
	"fmt"
	"os"

	"github.com/macinnir/dvc/commands"
)

var (
	configFilePath = "dvc.toml"
)

func main() {

	cmd := &commands.Cmd{}
	e := cmd.Run(os.Args)

	if e != nil {
		fmt.Printf("ERROR: %s", e.Error())
		return
	}
}
