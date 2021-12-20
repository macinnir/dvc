package main

import (
	"fmt"
	"os"

	"github.com/macinnir/dvc/core/commands"
)

var (
	Version = "v0.0.0"
)

func main() {

	cmd := &commands.Cmd{}

	if len(os.Args) > 1 && os.Args[1] == "-v" {
		fmt.Println(Version)
		os.Exit(0)
	}

	e := cmd.Run(os.Args)

	if e != nil {
		fmt.Printf("ERROR: %s\n", e.Error())
		return
	}
}
