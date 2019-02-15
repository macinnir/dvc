package gen

import (
	"fmt"
	"github.com/macinnir/dvc/lib"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// Gen conntains all of the generator functionality
type Gen struct {
	Options lib.Options
	Config  *lib.Config
}

// EnsureDir creates a new dir if the dir is not found
func (g *Gen) EnsureDir(dir string) (e error) {

	if _, e = os.Stat(dir); os.IsNotExist(e) {
		os.Mkdir(dir, 0777)
	}
	return
}

// FmtGoCode formats a go file
func (g *Gen) FmtGoCode(filePath string) {
	_, stdError, exitCode := lib.RunCommand("go", "fmt", filePath)

	if exitCode > 0 {
		lib.Warnf("fmt error: %s", g.Options, stdError)
	}
}

// WriteGoCodeToFile writes a string of golang code to a file and then formats it with `go fmt`
func (g *Gen) WriteGoCodeToFile(goCode string, filePath string) (e error) {
	// outFile := "./repos/repos.go"

	e = ioutil.WriteFile(filePath, []byte(goCode), 0644)
	if e != nil {
		return
	}

	g.FmtGoCode(filePath)
	// cmd := exec.Command("go", "fmt", filePath)
	// e = cmd.Run()
	// fmt.Printf("WriteCodeToFile: %s\n", e.Error())
	return
}

func (g *Gen) fileExists(filePath string) bool {
	// Check if file exists
	if _, e := os.Stat(filePath); os.IsNotExist(e) {
		// fmt.Printf("File %s does not exist", filePath)
		return false
	}

	return true

}

// scanFileParts scans a file for template parts, header, footer and import statements and returns those parts
func (g *Gen) scanFileParts(filePath string, trackImports bool) (fileHead string, fileFoot string, imports []string, e error) {

	lineStart := -1
	lineEnd := -1
	var fileBytes []byte

	fileHead = ""
	fileFoot = ""
	imports = []string{}

	if !g.fileExists(filePath) {
		return
	}

	fileBytes, e = ioutil.ReadFile(filePath)

	if e != nil {
		lib.Error(e.Error(), g.Options)
		return
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

//
// String Generators
//

// scanStringForFuncSignature scans a string (a line of goCode) and returns matches if it is a golang function signature that matches
// signatureRegexp
func (g *Gen) scanStringForFuncSignature(str string, signatureRegexp string) (matches []string) {

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
