package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/macinnir/dvc/core/lib"
)

// var goAppBootstrapFileTemplate = template.Must(template.New("go-app-bootstrap-file").Parse(`

// `))

// GenAppBootstrapFile generates the services bootstrap file
func GenAppBootstrapFile(basePackage string) error {

	var start = time.Now()

	lib.EnsureDir(lib.AppServicesDir)
	lib.EnsureDir(lib.CoreServicesDir)

	var files []os.FileInfo
	var e error
	files, _ = ioutil.ReadDir(lib.CoreServicesDir)
	packages := []string{}
	for k := range files {
		if files[k].IsDir() {
			packages = append(packages, path.Join(lib.CoreServicesDir, files[k].Name()))
		}
	}

	files, _ = ioutil.ReadDir(lib.AppServicesDir)
	for k := range files {
		if files[k].IsDir() {
			packages = append(packages, path.Join(lib.AppServicesDir, files[k].Name()))
		}
	}

	// Write Definitions file
	var sb strings.Builder

	sb.WriteString(`// DO NOT EDIT; Auto generated
package services

import (
	"log"
	"` + path.Join(basePackage, "core/app") + `" 
`)

	for k := range packages {
		sb.WriteString("\t\"" + path.Join(basePackage, packages[k]) + "\"\n")
	}

	sb.WriteString(`)

// App is a container for the services layer down
type App struct { 
	*app.BaseApp 
	Services *Services 
}

// Services is a container for all services 
type Services struct {
`)
	for k := range packages {
		packageName := path.Base(packages[k])
		sb.WriteString("\t" + strings.ToUpper(packageName[0:1]) + packageName[1:] + " *" + packageName + ".Services\n")
	}
	sb.WriteString(`}

// InitAppFromCLI initializes the application (presumably from the command line)
func InitAppFromCLI(
	configFilePath, 
	appName, 
	version, 
	commitHash, 
	buildDate, 
	clientVersion string,
) *App { 
	
	if len(appName) == 0 { 
		log.Fatal("App name cannot be empty") 
	}

	baseApp, coreRepos, authLog := app.NewBaseApp(configFilePath, appName, version, commitHash, buildDate, clientVersion) 

	app := &App { 
		BaseApp: baseApp, 
	}

	app.Services = &Services {`)
	for k := range packages {
		packageName := path.Base(packages[k])
		sb.WriteString("\n\t\t" + strings.ToUpper(packageName[0:1]) + packageName[1:] + ": " + packageName + ".NewServices(app.DAL, app.Config, app.Integrations, authLog, coreRepos, app.Cache),")
	}
	sb.WriteString(`
	}

	return app
} 

// Finish cleans up any connections from the app
func (a *App) Finish() {
	for schemaName := range a.Integrations.DB {
		for k := range a.Integrations.DB[schemaName] {
			a.Integrations.DB[schemaName][k].Close()
		}
	}
}
`)
	var bootstrapDir = filepath.Dir(lib.AppBootstrapFile)
	lib.EnsureDir(bootstrapDir)

	ioutil.WriteFile(lib.AppBootstrapFile, []byte(sb.String()), 0777)

	fmt.Printf("Generated app bootstrap file to %s in %f seconds\n", lib.AppBootstrapFile, time.Since(start).Seconds())

	return e
}
