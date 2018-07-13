package main

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
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

// WriteFile writes text to a file
func WriteFile(filePath string, contents string) (newFilePath string, e error) {

	// paths := []string{}

	// paths, e = d.Files.ScanChangesetDir(d.Config.ChangeSetPath)

	// if e != nil {
	// 	return
	// }

	// ordinalInt := 0

	// for _, p := range paths {

	// 	if len(p) < 11 {
	// 		continue
	// 	}

	// 	ordinal := p[0:6]

	// 	ordinalInt, e = strconv.Atoi(ordinal)

	// }

	// ordinalInt++

	// nextFile := fmt.Sprintf("%06d", ordinalInt) + ".sql"
	// newFilePath = d.Config.ChangeSetPath + "/changes.sql"

	e = ioutil.WriteFile(filePath, []byte(contents), 0644)
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

// scanFileParts scans a file for template parts, header, footer and import statements and returns those parts
func scanFileParts(filePath string, trackImports bool) (fileHead string, fileFoot string, imports []string, e error) {

	lineStart := -1
	lineEnd := -1
	var fileBytes []byte

	fileHead = ""
	fileFoot = ""
	imports = []string{}

	// Check if file exists
	if _, e = os.Stat(filePath); os.IsNotExist(e) {
		e = nil
		return
	}

	fileBytes, e = ioutil.ReadFile(filePath)

	if e != nil {
		fatal(e.Error())
	}

	fileString := string(fileBytes)
	fileLines := strings.Split(fileString, "\n")

	isImports := false

	for lineNum, line := range fileLines {

		line = strings.Trim(line, " ")

		if trackImports == true {

			if line == "import (" {
				isImports = true
				continue
			}

			if isImports == true {
				if line == ")" {
					isImports = false
					continue
				}

				imports = append(imports, line[2:len(line)-1])
				continue
			}

		}

		if line == "// #genStart" {
			lineStart = lineNum
			continue
		}

		if line == "// #genEnd" {
			lineEnd = lineNum
			continue
		}

		if lineStart == -1 {
			fileHead += line + "\n"
			continue
		}

		if lineEnd > -1 {
			fileFoot += line + "\n"
		}
	}

	if lineStart == -1 || lineEnd == -1 {
		e = fmt.Errorf("No gen tags found in outFile at path %s", filePath)
	}

	return
}

// WriteGoCodeToFile writes a string of golang code to a file and then formats it with `go fmt`
func WriteGoCodeToFile(goCode string, filePath string) (e error) {
	// outFile := "./repos/repos.go"

	e = ioutil.WriteFile(filePath, []byte(goCode), 0644)
	if e != nil {
		return
	}
	cmd := exec.Command("go", "fmt", filePath)
	e = cmd.Run()
	return
}

// scanStringForFuncSignature scans a string (a line of goCode) and returns matches if it is a golang function signature that matches
// signatureRegexp
func scanStringForFuncSignature(str string, signatureRegexp string) (matches []string) {

	lines := strings.Split(str, "\n")

	var validSignature = regexp.MustCompile(signatureRegexp)

	matches = []string{}

	for _, line := range lines {
		if validSignature.Match([]byte(line)) {
			matches = append(matches, line)
		}
	}

	return
}

// ReadSchemaFromFile Unmarshal's database json to a Database object
func ReadSchemaFromFile(filePath string) (database *Database, e error) {

	fileBytes, e := ioutil.ReadFile(filePath)

	if e != nil {
		return
	}

	database = &Database{}

	e = json.Unmarshal(fileBytes, database)
	if e != nil {
		return
	}
	return
}

// loadConfigFromFile loads a config file
func loadConfigFromFile(configFilePath string) (config *Config, e error) {

	// fmt.Printf("Looking for config at path %s\n", configFilePath)
	if _, e = os.Stat(configFilePath); os.IsNotExist(e) {
		e = fmt.Errorf("Config file `%s` not found", configFilePath)
		return
	}

	config = &Config{}
	_, e = toml.DecodeFile(configFilePath, config)
	return
}
