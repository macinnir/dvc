package gen

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/macinnir/dvc/lib"
)

// GenerateServiceInterfaces scans the services directory and outputs 2 files
// 	1. A services bootstrap file in the services directory
//  2. A services definition file in the definitions directory
func (g *Gen) GenerateServiceInterfaces(definitionsDir string, servicesDir string) (e error) {

	g.EnsureDir(servicesDir)

	var data = struct {
		BasePackage string
		Imports     []string
		Services    map[string][]string
	}{
		BasePackage: g.Config.BasePackage,
		Imports: []string{
			fmt.Sprintf("%s/definitions/models", g.Config.BasePackage),
			"github.com/macinnir/dvc/modules/utils",
		},
		Services: map[string][]string{},
	}

	viewmodelsPath := fmt.Sprintf("%s/viewmodels", definitionsDir)

	if g.dirExists(viewmodelsPath) && !g.dirIsEmpty(viewmodelsPath) {
		data.Imports = append(data.Imports, fmt.Sprintf("%s/definitions/viewmodels", g.Config.BasePackage))
	}

	// fmt.Sprintf("%s/definitions/viewmodels", g.Config.BasePackage),
	// g.EnsureDir("definitions/viewmodels")

	var serviceNames []string

	serviceNames, e = g.getServiceNames(servicesDir)
	if e != nil {
		panic(e)
	}

	var fileBytes []byte

	funcPrefix := "func (s *Services) "

	for _, serviceName := range serviceNames {

		data.Services[serviceName] = []string{}

		// Get the service file
		serviceFilePath := path.Join(servicesDir, serviceName+".go")
		fmt.Printf("ServiceFilePath: %s\n", serviceFilePath)
		fileBytes, e = ioutil.ReadFile(serviceFilePath)

		if e != nil {
			lib.Error(e.Error(), g.Options)
			return
		}

		fileString := string(fileBytes)
		fileLines := strings.Split(fileString, "\n")

		var validSignature = regexp.MustCompile(`^func \(s \*Services\) [A-Z].*$`)

		for _, line := range fileLines {
			if validSignature.Match([]byte(line)) {
				// Remove the prefix and the ending space and open bracket
				funcLine := line[len(funcPrefix) : len(line)-2]
				data.Services[serviceName] = append(data.Services[serviceName], funcLine)
			}
		}
	}

	for serviceName, services := range data.Services {
		fmt.Printf("Service: %s\n", serviceName)
		for _, method := range services {
			fmt.Printf("Service: %s.%s\n", serviceName, method)
		}
	}

	t := template.New("service-interfaces")

	// {{range $serviceName, $service := .Services}}	{{$serviceName}} I{{$serviceName}}Service
	// // I{{$serviceName}}Service outlines the service methods for the {{$serviceName}} service
	// type I{{$serviceName}}Service interface {
	// {{end}}

	tpl := `
// Package definitions outlines objects and functionality used in the {{.BasePackage}} application
package definitions
import ({{range .Imports}}
	"{{.}}"{{end}}
)

// Services defines the container for all service layer structs
type IServices interface {

	{{range $serviceName, $service := .Services}}
	// {{$serviceName}}
	{{range $service}}	{{.}}
	{{end}}
{{end}}}
// #genEnd
`

	p := path.Join(definitionsDir, "services.go")
	f, _ := os.Create(p)
	t, _ = t.Parse(tpl)
	e = t.Execute(f, data)
	if e != nil {
		fmt.Println("Execute Error: ", e.Error())
	}
	f.Close()
	g.FmtGoCode(p)
	return
}

func (g *Gen) GenerateServiceBootstrapFile(servicesDir string) (e error) {

	t := template.New("service-bootstrap")

	var serviceNames []string

	serviceNames, e = g.getServiceNames(servicesDir)
	if e != nil {
		return
	}

	// {{range .Services}}services.{{.}} = &{{.}}{config, repos, store}
	// {{end}}

	var data = struct {
		BasePackage string
		Services    []string
	}{
		BasePackage: g.Config.BasePackage,
		Services:    serviceNames,
	}
	tpl := `
// Package services provides the service methods objects and functionality used in the {{.BasePackage}} application
package services 
import (
	"{{.BasePackage}}/definitions/models"
	"{{.BasePackage}}/definitions"
	"github.com/macinnir/dvc/modules/utils" 
)

type Services struct {
	config *models.Config
	repos *definitions.Repos 
	store utils.IStore
}

// Bootstrap instantiates a new Services instance and all of its members 
func Bootstrap(config *models.Config, repos *definitions.Repos, store utils.IStore) *Services {
	services := &Services{config, repos, store} 
	return services
}
// #genEnd 
	`

	p := path.Join(servicesDir, "bootstrap.go")
	f, _ := os.Create(p)
	t, _ = t.Parse(tpl)
	e = t.Execute(f, data)
	if e != nil {
		fmt.Println("Execute Error: ", e.Error())
	}
	f.Close()
	g.FmtGoCode(p)
	return
}

// GetServiceNames gets a list of services in the services directory
func (g *Gen) getServiceNames(dir string) (serviceNames []string, e error) {

	serviceNames = []string{}
	dirFileNames := []string{}
	var dirHandle *os.File

	dirHandle, e = os.Open(dir)
	if e != nil {
		return
	}

	defer dirHandle.Close()
	dirFileNames, e = dirHandle.Readdirnames(-1)
	if e != nil {
		return
	}

	for _, fileName := range dirFileNames {

		var fileInfo os.FileInfo

		p := path.Join(dir, fileName)

		// Skip directories
		if fileInfo, e = os.Stat(p); e != nil {
			fmt.Printf("CODEGEN.Services> SKIP: File Not Found: `%s`\n ", p)
			continue
		}

		if fileInfo.IsDir() {
			fmt.Printf("CODEGEN.Services> SKIP: File is directory: `%s`\n", p)
			continue
		}

		if !isGeneratableFile(fileName) {
			continue
		}

		fmt.Printf("CODEGEN.Services> GENERATE: %s\n", fileName[0:len(fileName)-3])
		serviceNames = append(serviceNames, fileName[0:len(fileName)-3])
	}
	return
}

func isGeneratableFile(fileName string) bool {

	fileLen := len(fileName)

	if fileLen < 4 {
		fmt.Printf("CODEGEN.Services> SKIP: FileName too short: `%s`\n", fileName)
		return false
	}

	if fileName[fileLen-3:] != ".go" {
		fmt.Printf("CODEGEN.Services> SKIP: Not a go file: `%s`\n", fileName)
		return false
	}

	if fileLen > 8 && fileName[fileLen-8:] == "_test.go" {
		fmt.Printf("CODEGEN.Services> SKIP: Test file: `%s`\n", fileName)
		return false
	}

	// https://yourbasic.org/golang/regexp-cheat-sheet/
	// https://regex-golang.appspot.com/assets/html/index.html
	r, e := regexp.MatchString("^[A-Z]{1}.+", fileName)
	if e != nil {
		panic(e)
	}

	// Skip bootstrap file, test files and anything not a go file
	if !r {
		fmt.Printf("CODEGEN.Services> SKIP: Invalid file format: %s\n", fileName)
		return false
	}

	return true
}
