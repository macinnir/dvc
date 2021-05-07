package main

import (
	"fmt"
	"os"

	"github.com/macinnir/dvc/core/commands"
)

func main() {

	cmd := &commands.Cmd{}
	e := cmd.Run(os.Args)

	if e != nil {
		fmt.Printf("ERROR: %s\n", e.Error())
		return
	}
}
