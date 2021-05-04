package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Init creates a default dvc.toml file in the CWD
func (c *Cmd) Init(args []string) {

	if len(args) > 0 && args[0] == "help" {
		helpInit()
		return
	}

	var e error

	if _, e = os.Stat("./dvc.toml"); os.IsNotExist(e) {

		reader := bufio.NewReader(os.Stdin)

		// https://tutorialedge.net/golang/reading-console-input-golang/
		// BasePackage
		fmt.Print("> Base Package:")
		basePackage, _ := reader.ReadString('\n')
		basePackage = strings.Replace(basePackage, "\n", "", -1)

		fmt.Print("> Base directory (leave blank for current):")
		baseDir, _ := reader.ReadString('\n')
		baseDir = strings.Replace(baseDir, "\n", "", -1)

		// Host
		fmt.Print("> Database Host:")
		host, _ := reader.ReadString('\n')
		host = strings.Replace(host, "\n", "", -1)

		// databaseName
		fmt.Print("> Database Name:")
		databaseName, _ := reader.ReadString('\n')
		databaseName = strings.Replace(databaseName, "\n", "", -1)

		// databaseUser
		fmt.Print("> Database User:")
		databaseUser, _ := reader.ReadString('\n')
		databaseUser = strings.Replace(databaseUser, "\n", "", -1)

		// databasePass
		fmt.Print("> Database Password:")
		databasePass, _ := reader.ReadString('\n')
		databasePass = strings.Replace(databasePass, "\n", "", -1)

		content := "databaseType = \"mysql\"\nbasePackage = \"" + basePackage + "\"\n\nenums = []\n\n"
		content += "[connection]\nhost = \"" + host + "\"\ndatabaseName = \"" + databaseName + "\"\nusername = \"" + databaseUser + "\"\npassword = \"" + databasePass + "\"\n\n"

		packages := []string{
			"repos",
			"models",
			"typescript",
			"services",
			"dal",
			"definitions",
		}

		content += "[packages]\n"
		for _, p := range packages {
			if p == "typescript" {
				continue
			}

			content += fmt.Sprintf("%s = \"%s\"\n", p, path.Join(basePackage, p))
		}

		// content += "[packages]\ncache = \"myPackage/cache\"\nmodels = \"myPackage/models\"\nschema = \"myPackage/schema\"\nrepos = \"myPackage/repos\"\n\n"

		content += "[dirs]\n"

		for _, p := range packages {

			if baseDir != "" {
				content += fmt.Sprintf("%s = \"%s\"\n", p, path.Join(baseDir, p))
			} else {
				content += fmt.Sprintf("%s = \"%s\"\n", p, p)
			}
		}

		// content += "[dirs]\nrepos = \"repos\"\ncache = \"cache\"\nmodels = \"models\"\nschema = \"schema\"\ntypescript = \"ts\""

		ioutil.WriteFile("./dvc.toml", []byte(content), 0644)

	} else {
		fmt.Println("dvc.toml already exists in this directory")
	}
}
