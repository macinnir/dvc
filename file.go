package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/macinnir/dvc/types"
)

// FetchFile fetches a file contents
func FetchFile(filePath string) (changesetFileString string, e error) {

	// filePath := d.Config.ChangeSetPath + "/changes.sql"

	if _, e = os.Stat(filePath); os.IsNotExist(e) {
		return
	}

	var fileBytes []byte

	fileBytes, e = ioutil.ReadFile(filePath)
	if e != nil {
		return
	}

	changesetFileString = string(fileBytes)
	return
}

// ScanDir scans the changeset directory for any file with an extension
func (f *Files) ScanDir(rootPath string, extension string) (paths []string, e error) {

	var fileInfos []os.FileInfo
	paths = []string{}

	fileInfos, e = ioutil.ReadDir(rootPath)

	for _, f := range fileInfos {
		if f.IsDir() || filepath.Ext(f.Name()) != extension {
			continue
		}
		paths = append(paths, f.Name())
	}

	return

}

// ReadSchemaFromFile Unmarshal's database json to a Database object
func ReadSchemaFromFile(filePath string) (database *types.Database, e error) {

	fileBytes, e := ioutil.ReadFile(filePath)

	if e != nil {
		return
	}

	database = &types.Database{}

	e = json.Unmarshal(fileBytes, database)
	if e != nil {
		return
	}
	return
}

// loadConfigFromFile loads a config file
func loadConfigFromFile(configFilePath string) (config *types.Config, e error) {

	// fmt.Printf("Looking for config at path %s\n", configFilePath)
	if _, e = os.Stat(configFilePath); os.IsNotExist(e) {
		e = fmt.Errorf("Config file `%s` not found", configFilePath)
		return
	}

	config = &types.Config{}
	_, e = toml.DecodeFile(configFilePath, config)
	return
}
