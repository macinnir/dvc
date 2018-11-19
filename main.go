package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/macinnir/dvc/connectors/mysql"
	"github.com/macinnir/dvc/connectors/sqlite"
	"github.com/macinnir/dvc/lib"
	"os"
)

var (
	configFilePath = "dvc.toml"
)

func main() {

	config, e := loadConfigFromFile("./dvc.toml")
	if e != nil {
		fmt.Println("Could not load config file.")
		// fmt.Printf("ERROR: %s\n", e.Error())
		// os.Exit(1)
		return
	}
	// config = &Config{
	// 	Host:          args[0],
	// 	DatabaseName:  args[1],
	// 	Username:      args[2],
	// 	Password:      args[3],
	// 	ChangeSetPath: args[4],
	// 	DatabaseType:  args[5],
	// }

	cmd := &Cmd{
		Config: config,
	}

	if cmd.dvc, e = lib.NewDVC(config); e != nil {
		fmt.Printf("ERROR: %s", e.Error())
	}

	cmd.dvc.Connector = connectorFactory(config.DatabaseType, config)
	e = cmd.Main(os.Args)

	if e != nil {
		fmt.Printf("ERROR: %s", e.Error())
		return
	}
}

// loadConfigFromFile loads a config file
func loadConfigFromFile(configFilePath string) (config *lib.Config, e error) {

	// fmt.Printf("Looking for config at path %s\n", configFilePath)
	if _, e = os.Stat(configFilePath); os.IsNotExist(e) {
		e = fmt.Errorf("Config file `%s` not found", configFilePath)
		return
	}

	config = &lib.Config{}
	_, e = toml.DecodeFile(configFilePath, config)
	return
}

func connectorFactory(connectorName string, config *lib.Config) (connector lib.IConnector) {
	if connectorName == "mysql" {
		connector = &mysql.MySQL{
			Config: config,
		}
	}

	if connectorName == "sqlite" {
		connector = &sqlite.Sqlite{
			Config: config,
		}
	}

	return
}
