package main

import (
	"fmt"
	"os"
)

var (
	configFilePath = "dvc.toml"
)

func main() {

	var e error

	cmd := &Cmd{}
	if cmd.dvc, e = NewDVC(configFilePath); e != nil {
		fmt.Printf("ERROR: %s", e.Error())
	}
	e = cmd.Main(os.Args)

	if e != nil {
		fmt.Printf("ERROR: %s", e.Error())
		return
	}
}
