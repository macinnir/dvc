package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/macinnir/dvc/core/lib"
)

// GenInterface takes makes the interface into a byte array
func GenInterface(comment, pkgName, ifaceName, ifaceComment string, methods []string, imports []string) ([]byte, error) {

	output := []string{"// " + comment, "", "package " + pkgName, "import ("}
	output = append(output, imports...)
	output = append(output, ")", "")
	if len(ifaceComment) > 0 {
		output = append(output, fmt.Sprintf("// %s", strings.Replace(ifaceComment, "\n", "\n// ", -1)))
	}
	output = append(output, fmt.Sprintf("type %s interface {", ifaceName))
	output = append(output, methods...)
	output = append(output, "}")
	code := strings.Join(output, "\n")
	return lib.FormatCode(code)
}

func fetchExistingInterfaceFiles(srcDir string) ([]string, error) {

	// fmt.Println("Fetching files for dir ", srcDir)
	filePaths := []string{}
	var files []os.FileInfo
	var e error

	// DAL
	if files, e = ioutil.ReadDir(srcDir); e != nil {
		return filePaths, e
	}

	for k := range files {

		f := files[k]

		filePath := path.Join(srcDir, files[k].Name())

		if files[k].IsDir() {
			var subFiles []string
			if subFiles, e = fetchExistingInterfaceFiles(filePath); e != nil {
				return filePaths, e
			}
			filePaths = append(filePaths, subFiles...)
			continue
		}

		// Filter out files that don't have upper case first letter names
		if !unicode.IsUpper([]rune(f.Name())[0]) {
			continue
		}

		// Verify this is a go file
		if f.Name()[len(f.Name())-3:] != ".go" {
			continue
		}

		filePaths = append(filePaths, filePath)
	}

	return filePaths, nil

}

func fetchSrcFilesForInterfaces(srcDir string) ([]string, error) {

	// fmt.Println("Fetching files for dir ", srcDir)
	filePaths := []string{}
	var files []os.FileInfo
	var e error

	// DAL
	if files, e = ioutil.ReadDir(srcDir); e != nil {
		return filePaths, e
	}

	for k := range files {

		f := files[k]

		filePath := path.Join(srcDir, files[k].Name())

		if files[k].IsDir() {
			var subFiles []string
			if subFiles, e = fetchSrcFilesForInterfaces(filePath); e != nil {
				return filePaths, e
			}
			filePaths = append(filePaths, subFiles...)
			continue
		}

		// Filter out files that don't have upper case first letter names
		if !unicode.IsUpper([]rune(f.Name())[0]) {
			continue
		}

		// Verify this is a go file
		if f.Name()[len(f.Name())-3:] != ".go" {
			continue
		}

		// Remove the go extension
		baseName := f.Name()[0 : len(f.Name())-3]

		// Skip over EXT files
		if baseName[len(baseName)-3:] == "Ext" {
			continue
		}

		// Skip over test files
		if baseName[len(baseName)-5:] == "_test" {
			continue
		}

		filePaths = append(filePaths, filePath)
	}

	return filePaths, nil

}

// services/services.go
func GenServicesBootstrap(config *lib.Config) error {

	lib.EnsureDir("app/definitions")
	lib.EnsureDir("core/definitions")

	var files []os.FileInfo
	var e error
	files, e = ioutil.ReadDir("core/services")
	packages := []string{}
	for k := range files {
		if files[k].IsDir() {
			packages = append(packages, path.Join("core/services", files[k].Name()))
		}
	}

	files, e = ioutil.ReadDir("app/services")
	for k := range files {
		if files[k].IsDir() {
			packages = append(packages, path.Join("app/services", files[k].Name()))
		}
	}

	// 	tpl := `// DO NOT EDIT; Auto generated
	// package services

	// import (
	// 	"` + path.Join(config.BasePackage, config.Dirs.Integrations) + `"
	// 	"` + path.Join(config.BasePackage, "gen/dal") + `"
	// 	"` + path.Join(config.BasePackage, "core/definitions") + `"
	// 	"` + path.Join(config.BasePackage, "core") + `"
	// 	"` + path.Join(config.BasePackage, "core/repos") + `"
	// 	"` + path.Join(config.BasePackage, "core/cache") + `"
	// `
	// 	for k := range packages {
	// 		tpl += "\t\"" + path.Join(config.BasePackage, packages[k]) + "\"\n"
	// 	}
	// 	tpl += `
	// )

	// func bootstrapServices(
	// 	d *dal.DAL,
	// 	config *core.Config,
	// 	i *integrations.Integrations,
	// 	a *AuthLog,
	// 	r *repos.Repos,
	// 	c *cache.Cache,
	// ) *definitions.Services {
	// 	return &definitions.Services {
	// `

	// 	for k := range packages {
	// 		packageName := path.Base(packages[k])
	// 		tpl += "\t\t" + strings.ToUpper(packageName[0:1]) + packageName[1:] + ": " + packageName + ".NewServices(d, config, i, a, r, c),\n"
	// 	}

	// 	tpl += `	}
	// }
	// `

	// 	ioutil.WriteFile(path.Join(config.Dirs.Services, "services.go"), []byte(tpl), 0777)

	// Write Definitions file
	tpl := `// DO NOT EDIT; Auto generated
package definitions 

import (
	"log"
	"` + path.Join(config.BasePackage, "core/app") + `" 
`

	for k := range packages {
		tpl += "\t\"" + path.Join(config.BasePackage, packages[k]) + "\"\n"
	}

	tpl += `)

// App is a container for the services layer down
type App struct { 
	*app.BaseApp 
	Services *Services 
}

// Services is a container for all services 
type Services struct {
`
	for k := range packages {
		packageName := path.Base(packages[k])
		tpl += "\t" + strings.ToUpper(packageName[0:1]) + packageName[1:] + " *" + packageName + ".Services\n"
	}
	tpl += `}

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

	app.Services = &Services {`
	for k := range packages {
		packageName := path.Base(packages[k])
		tpl += "\n\t\t" + strings.ToUpper(packageName[0:1]) + packageName[1:] + ": " + packageName + ".NewServices(app.DAL, app.Config, app.Integrations, authLog, coreRepos, app.Cache),"
	}
	tpl += `
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
`
	ioutil.WriteFile("gen/definitions/app.go", []byte(tpl), 0777)

	return e
}

func GenInterfaces(srcDir, destDir string) error {

	existingInterfaces, _ := fetchExistingInterfaceFiles(destDir)

	existingInterfaceMap := map[string]bool{}
	for k := range existingInterfaces {
		baseName := path.Base(existingInterfaces[k])
		s := baseName[1 : len(baseName)-3]
		existingInterfaceMap[s] = true
		// fmt.Println(path.Base(existingInterfaces[k][1 : len(existingInterfaces[k])-3])
	}

	start := time.Now()

	lib.EnsureDir(destDir)

	var e error
	generatedInterfaces := 0
	var files []string
	if files, e = fetchSrcFilesForInterfaces(srcDir); e != nil {
		return e
	}

	// removeInterface := []string{}
	addedInterfaceMap := map[string]bool{}
	for k := range files {

		srcFile := files[k]
		baseName := filepath.Base(srcFile)
		structName := baseName[0 : len(baseName)-3]
		addedInterfaceMap[structName] = true
		interfaceName := "I" + structName
		packageName := filepath.Base(filepath.Dir(srcFile))
		destDirName := filepath.Base(destDir)
		// TODO verbose flag?
		// fmt.Println("DestDir: ", destDirName)

		subDestDir := destDir

		if packageName != destDirName {
			subDestDir = path.Join(destDir, packageName)
		}

		// fmt.Println("SubDestDir: ", subDestDir)
		destFile := path.Join(subDestDir, interfaceName+".go")
		srcFiles := []string{srcFile}

		// fmt.Println("Generating interface: " + srcFile + " ==> " + destFile)

		// Check if EXT file exists
		extFile := srcFile[0:len(srcFile)-3] + "Ext.go"
		if _, e = os.Stat(extFile); e == nil {
			srcFiles = append(srcFiles, extFile)
		}

		var i []byte
		i, e = genInterfaces(
			srcFiles,
			structName,
			"Generated Code; DO NOT EDIT.",
			packageName,
			interfaceName,
			fmt.Sprintf("%s describes the %s", interfaceName, baseName),
			true,
			true,
		)

		if e != nil {
			fmt.Println("ERROR: " + srcFile + " => " + e.Error())
			return e
		}

		lib.EnsureDir(subDestDir)

		// TODO verbose flag
		// fmt.Printf("Generating interface %s...\n", destFile)
		ioutil.WriteFile(destFile, i, 0644)

		generatedInterfaces++

	}

	// for k := range existingInterfaceMap {
	// 	if _, ok := addedInterfaceMap[k]; !ok {
	// 		fullPath := path.Join(destDir, "I"+k+".go")
	// 		// if e = os.Remove(fullPath); e != nil {
	// 		// 	return fmt.Errorf("remove file %s: %e", fullPath, e)
	// 		// }
	// 		fmt.Println("Removed interface at ", fullPath)
	// 	}
	// }

	fmt.Printf("Generated %d interfaces to %s in %f seconds\n", generatedInterfaces, destDir, time.Since(start).Seconds())

	return nil

}

// GenInterfaces runs GenInterface on all the target files
func genInterfaces(
	files []string,
	structType,
	comment,
	pkgName,
	ifaceName,
	ifaceComment string,
	copyDocuments,
	copyTypeDoc bool,
) (result []byte, e error) {

	// fmt.Printf("Generating interface: %s\n", files[0])
	allMethods := []string{}
	allImports := []string{}

	mset := make(map[string]struct{})
	iset := make(map[string]struct{})

	var typeDoc string

	for _, f := range files {

		var src []byte
		if src, e = ioutil.ReadFile(f); e != nil {
			return
		}

		methods, imports, parsedTypeDoc := lib.ParseStruct(src, structType, copyDocuments, copyTypeDoc, pkgName)
		for _, m := range methods {
			if _, ok := mset[m.Code]; !ok {
				allMethods = append(allMethods, m.Lines()...)
				mset[m.Code] = struct{}{}
			}
		}
		for _, i := range imports {
			if _, ok := iset[i]; !ok {
				allImports = append(allImports, i)
				iset[i] = struct{}{}
			}
		}
		if typeDoc == "" {
			typeDoc = parsedTypeDoc
		}
	}

	if typeDoc != "" {
		ifaceComment += "\n" + typeDoc
	}

	result, e = GenInterface(comment, pkgName, ifaceName, ifaceComment, allMethods, allImports)

	return
}
