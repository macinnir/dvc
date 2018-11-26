package main

import (
	"fmt"
	"os"
)

var (
	configFilePath = "dvc.toml"
)

func main() {

	cmd := &Cmd{}
	e := cmd.Main(os.Args)

	if e != nil {
		fmt.Printf("ERROR: %s", e.Error())
		return
	}
}
