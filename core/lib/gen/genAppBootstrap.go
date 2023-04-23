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

	var files []os.FileInfo
	var e error
	packages := []string{}

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
	"` + path.Join(basePackage, "core/services/base") + `"
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
	Base *base.Services
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

	var app = &App { 
		BaseApp: app.NewBaseApp(configFilePath, appName, version, commitHash, buildDate, clientVersion),
		Services: &Services{}, 
	}

	// Base Services 
	app.Services.Base = base.NewServices(app.DAL, app.Config, app.Integrations, app.Log, app.Repos, app.Cache)`)
	for k := range packages {
		packageName := path.Base(packages[k])
		sb.WriteString("\n\tapp.Services." + strings.ToUpper(packageName[0:1]) + packageName[1:] + " = " + packageName + ".NewServices(app.Services.Base, app.DAL, app.Config, app.Integrations, app.Log, app.BaseApp.Repos, app.Cache)")
	}
	sb.WriteString(`
	return app
} 

// Finish cleans up any connections from the app
func (a *App) Finish() {
	a.Finish()
}
`)
	var bootstrapDir = filepath.Dir(lib.AppBootstrapFile)
	lib.EnsureDir(bootstrapDir)

	ioutil.WriteFile(lib.AppBootstrapFile, []byte(sb.String()), 0777)

	fmt.Printf("Generated app bootstrap file to %s in %f seconds\n", lib.AppBootstrapFile, time.Since(start).Seconds())

	return e
}
