package gen

import (
	"fmt"
	"github.com/macinnir/dvc/lib"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

// GenerateServiceInterfaces scans the services directory and outputs 2 files
// 	1. A services bootstrap file in the services directory
//  2. A services definition file in the definitions directory
func (g *Gen) GenerateServiceInterfaces(definitionsDir string, servicesDir string) (e error) {

	var data = struct {
		BasePackage string
		Imports     []string
		Services    map[string][]string
	}{
		BasePackage: g.Config.BasePackage,
		Imports: []string{
			fmt.Sprintf("%s/definitions/models", g.Config.BasePackage),
			fmt.Sprintf("%s/definitions/viewmodels", g.Config.BasePackage),
			"github.com/macinnir/dvc/modules/utils",
		},
		Services: map[string][]string{},
	}

	var serviceNames []string

	serviceNames, e = g.getServiceNames(servicesDir)
	if e != nil {
		return
	}

	var fileBytes []byte

	for _, serviceName := range serviceNames {

		data.Services[serviceName] = []string{}

		// Get the service file
		fileBytes, e = ioutil.ReadFile(path.Join(servicesDir, serviceName+".go"))

		if e != nil {
			lib.Error(e.Error(), g.Options)
			return
		}

		fileString := string(fileBytes)
		fileLines := strings.Split(fileString, "\n")

		funcSig := fmt.Sprintf(`^func \(s \*%s\) [A-Z].*$`, serviceName)
		var validSignature = regexp.MustCompile(funcSig)
		funcPrefix := fmt.Sprintf("func (s *%s) ", serviceName)

		for _, line := range fileLines {
			if validSignature.Match([]byte(line)) {
				// Remove the prefix and the ending space and open bracket
				funcLine := line[len(funcPrefix) : len(line)-2]
				data.Services[serviceName] = append(data.Services[serviceName], funcLine)
			}
		}
	}

	t := template.New("service-interfaces")

	tpl := `
// Package definitions outlines objects and functionality used in the {{.BasePackage}} application
package definitions
import ({{range .Imports}}
	"{{.}}"{{end}}
)

// Services defines the container for all service layer structs
type Services struct {
{{range $serviceName, $service := .Services}}	{{$serviceName}} I{{$serviceName}}Service
{{end}}}

{{range $serviceName, $service := .Services}}
// I{{$serviceName}}Service outlines the service methods for the {{$serviceName}} service 
type I{{$serviceName}}Service interface {
{{range $service}}	{{.}}
{{end}}}
{{end}}
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

// Bootstrap instantiates a new Services instance and all of its members 
func Bootstrap(config *models.Config, repos *definitions.Repos, store utils.IStore) *definitions.Services {
	services := &definitions.Services{} 
	{{range .Services}}services.{{.}} = &{{.}}{config, repos, store}
	{{end}}
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

	fmt.Println("Gettings service names")
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

	var fileInfo os.FileInfo

	for _, fileName := range dirFileNames {

		p := path.Join(dir, fileName)
		fmt.Printf("1. File %s\n", p)
		// Skip directories
		if fileInfo, e = os.Stat(p); e != nil {
			fmt.Printf("Could not find file %s. Skipping...\n ", p)
			continue
		}

		if fileInfo.IsDir() {
			fmt.Printf("File %s is a directory. Skipping...\n", p)
		}

		r, _ := regexp.MatchString("[A-Z]{1}.+", fileName)

		if !r {
			fmt.Printf("Filename: %s does have an uppercase first letter. Skipping...\n", fileName)
			continue
		}

		// Skip bootstrap file, test files and anything not a go file
		if (len(fileName) > 8 && fileName[len(fileName)-8:len(fileName)] == "_test.go") ||
			fileName == "bootstrap.go" ||
			(len(fileName) > 3 && fileName[len(fileName)-3:len(fileName)] != ".go") {
			fmt.Printf("Skipping %s\n", fileName)
			continue
		}

		serviceNames = append(serviceNames, fileName[0:len(fileName)-3])
	}
	return
}
