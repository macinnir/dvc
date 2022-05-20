package gen

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/macinnir/dvc/core/lib"
)

// GenControllerBootstrap generates the bootstrap file for the applications controllers
func GenControllerBootstrap(basePackageName string, dirs []string) string {

	var sb strings.Builder

	sb.WriteString(`// DO NOT EDIT: Auto generated
package gen

import (

`)
	for k := range dirs {
		sb.WriteString("\t\"" + path.Join(basePackageName, dirs[k]) + "\"\n")
	}
	sb.WriteString(`	"` + basePackageName + `/gen/definitions"

	"github.com/macinnir/dvc/core/lib/utils/request"
)

// Controllers is the main container for all of the controller modules 
type Controllers struct {
`)

	for k := range dirs {
		packageName := path.Base(dirs[k])
		sb.WriteString("\t" + strings.ToUpper(packageName[0:1]) + packageName[1:] + " *" + packageName + ".Controllers\n")
	}

	sb.WriteString(`}

// NewControllers bootstraps all of the controller modules 
func NewControllers(s *definitions.App, r request.IResponseLogger) *Controllers { 
	return &Controllers{
`)

	for k := range dirs {
		packageName := path.Base(dirs[k])
		sb.WriteString("\t\t" + strings.ToUpper(packageName[0:1]) + packageName[1:] + ": " + packageName + ".NewControllers(s, r),\n")
	}

	sb.WriteString(`	}
}`)

	return sb.String()
}

func LoadRoutes(config *lib.Config) (*lib.RoutesJSONContainer, error) {

	var e error

	if _, e = os.Stat(lib.RoutesFilePath); os.IsNotExist(e) {
		return nil, errors.New("Routes file does not exist at path" + lib.RoutesFilePath)
	}

	routes := &lib.RoutesJSONContainer{}

	var fileBytes []byte

	if fileBytes, e = ioutil.ReadFile(lib.RoutesFilePath); e != nil {
		return nil, e
	}

	if e = json.Unmarshal(fileBytes, routes); e != nil {
		return nil, e
	}

	return routes, nil
}

// GenRoutes generates a list of routes from a directory of controller files
func GenRoutesAndPermissions(controllers []*lib.Controller, dirs []string, config *lib.Config) error {

	var e error
	imports := []string{
		path.Join(config.BasePackage, config.Dirs.IntegrationInterfaces),
		// path.Join(config.BasePackage, config.Dirs.Aggregates),
		"net/http",
		"github.com/gorilla/mux",
	}

	// fmt.Println(imports)

	code := ""

	rest := ""
	// controllerCalls := []string{}

	hasBodyImports := false
	packageUsesPermission := false

	if e != nil {
		return e
	}

	for k := range controllers {

		controller := controllers[k]

		if controller.PermCount > 0 {
			packageUsesPermission = true
		}

		// Documentation routes
		controllers = append(controllers, controller)

		// Include imports for dtos and response if necessary for JSON http body
		if controller.HasDTOsImport {
			hasBodyImports = true
		}

		var routesString string

		routesString, e = buildRoutesCodeFromController(controller)

		if e != nil {
			return e
		}

		rest += "\n" + routesString + "\n"

		// controllerCalls = append(
		// 	controllerCalls,
		// 	"map"+strings.Title(controller.Package)+controller.Name+"Routes(res, r, auth, c, log)",
		// )
	}

	// code += strings.Join(controllerCalls, "\n\t")
	code += rest
	code += "\n\n}\n"

	if hasBodyImports {
		// imports = append(imports, g.Config.BasePackage+"/core/utils/response")
		imports = append(imports, path.Join(config.BasePackage, lib.CoreDTOsDir))
	}

	imports = append(imports, lib.LibRequests)

	if packageUsesPermission {
		imports = append(imports, lib.LibUtils)
		imports = append(imports, path.Join(config.BasePackage, lib.GoPermissionsPath))
	}

	final := `// Generated Code; DO NOT EDIT.

package gen

import (
`

	for _, i := range imports {
		final += fmt.Sprintf("\t\"%s\"\n", i)
	}

	final += `
	appdtos "` + path.Join(config.BasePackage, lib.AppDTOsDir) + `"
)

// MapRoutesToControllers maps the routes to the controllers
func MapRoutesToControllers(r *mux.Router, auth integrations.IAuth, c *Controllers, res request.IResponseLogger, log integrations.ILog) {

	`
	final += code

	ioutil.WriteFile("gen/routes.go", []byte(final), 0777)

	// DTOS
	dtos := genDTOSMap(lib.CoreDTOsDir)
	appDTOs := genDTOSMap(lib.AppDTOsDir)
	for dtoName := range appDTOs {
		dtos[dtoName] = appDTOs[dtoName]
	}

	// Aggregates
	aggregates := genAggregatesMap("core/definitions/aggregates")
	appAggregates := genAggregatesMap("app/definitions/aggregates")
	for aggregateName := range appAggregates {
		aggregates[aggregateName] = appAggregates[aggregateName]
	}

	routesContainer := &lib.RoutesJSONContainer{
		Routes:     map[string]*lib.ControllerRoute{},
		DTOs:       dtos,
		Models:     genModelsMap(),
		Aggregates: aggregates,
		Constants:  genConstantsMap(),
	}

	for k := range controllers {
		for i := range controllers[k].Routes {
			key := controllers[k].Routes[i].Name
			routesContainer.Routes[key] = controllers[k].Routes[i]
		}
	}

	routesJSON, _ := json.MarshalIndent(routesContainer, "  ", "    ")
	// fmt.Println("Writing Routes JSON to path", lib.RoutesFilePath)
	ioutil.WriteFile(lib.RoutesFilePath, routesJSON, 0777)

	var controllerBootstrapFile = GenControllerBootstrap(config.BasePackage, dirs)

	ioutil.WriteFile(lib.ControllersBootstrapGenFile, []byte(controllerBootstrapFile), 0777)

	return nil
}

func buildRoutesCodeFromController(controller *lib.Controller) (out string, e error) {

	s := []string{
		"",
		"\t////",
		"\t// " + strings.Title(controller.Package) + "." + controller.Name,
		"\t////",
		"",
	}
	// 	fmt.Sprintf("// map%sRoutes maps all of the routes for %s", controller.Name, controller.Name),
	// 	fmt.Sprintf("func map%s%sRoutes(res request.IResponseLogger, r *mux.Router, auth integrations.IAuth, c *controllers.Controllers.%s, log integrations.ILog) {\n", strings.Title(controller.Package), controller.Name, strings.Title(controller.Package)),
	// }

	for _, route := range controller.Routes {

		// Method comments
		s = append(s, fmt.Sprintf("\t// %s.%s.%s", strings.Title(controller.Package), controller.Name, route.Name))
		s = append(s, fmt.Sprintf("\t// %s %s", route.Method, route.Raw))
		if !route.IsAuth {
			s = append(s, "\t// @anonymous")
		}

		if route.IsAuth && len(route.Permission) == 0 {
			s = append(s, "\t// @anyone")
		}

		// Method args
		args := []string{
			"w",   // http.ResponseWriter
			"req", // *http.Request
		}

		if route.IsAuth {
			s = append(s, fmt.Sprintf("\tr.Handle(\"%s\", auth.AuthMiddleware(func(w http.ResponseWriter, req *request.Request) {\n", route.Path))
			// s = append(s, fmt.Sprintf("\t\tcurrentUser := auth.GetCurrentUser(r)\n"))
			// args = append(args, "currentUser")
		} else {
			s = append(s, fmt.Sprintf("\tr.Handle(\"%s\", auth.AnonMiddleware(func(w http.ResponseWriter, req *request.Request) {\n", route.Path))
		}

		s = append(s, fmt.Sprintf("\t\tlog.Debug(\"ROUTE: %s %s => %s\")\n", route.Method, route.Path, route.Name))

		// Permission
		if route.IsAuth && len(route.Permission) > 0 {

			// ucFirst
			// permission := string(unicode.ToUpper(rune(route.Permission[0]))) + route.Permission[1:]

			s = append(s, `		// Requires permission `+route.Permission+`
		if !utils.HasPerm(req, req.User, permissions.`+route.Permission+`) {
			res.Forbidden(req, w)
			return
		}
`)

		}

		if len(route.BodyType) > 0 {

			if route.BodyType[0:1] == "*" {
				route.BodyType = route.BodyType[1:]
			}
			s = append(s, "\t\t// Parse the body of type "+route.BodyType)
			s = append(s, fmt.Sprintf("\t\tbody := &%s{}", route.BodyType))
			s = append(s, "\t\treq.BodyJSON(body)\n")
		}

		if len(route.Params) > 0 {
			for _, param := range route.Params {
				s = append(s, fmt.Sprintf("\t\t// URL Param %s", param.Name))
				s = append(s, fmt.Sprintf("\t\t// %s", param.Pattern))
				if param.Type == "int64" {
					defaultValue := "0"
					if param.Pattern == "-?[0-9]+" {
						defaultValue = "-1"
					}
					s = append(s, fmt.Sprintf("\t\t%s := req.ArgInt64(\"%s\", %s)\n", param.Name, param.Name, defaultValue))
				} else {
					s = append(s, fmt.Sprintf("\t\t%s := req.Arg(\"%s\", \"\")\n", param.Name, param.Name))
				}

				args = append(args, param.Name)
			}
		}

		if len(route.Queries) > 0 {
			for _, query := range route.Queries {
				s = append(s, fmt.Sprintf("\t\t// Query Arg %s", query.VariableName))
				s = append(s, fmt.Sprintf("\t\t// %s", query.Pattern))

				if query.Type == "int64" {
					defaultValue := "0"
					if query.Pattern == "-?[0-9]+" {
						defaultValue = "-1"
					}
					s = append(s, fmt.Sprintf("\t\t%s := req.ArgInt64(\"%s\", %s)\n", query.VariableName, query.VariableName, defaultValue))
				} else {
					s = append(s, fmt.Sprintf("\t\t%s := req.Arg(\"%s\", \"\")\n", query.VariableName, query.VariableName))
				}

				args = append(args, query.VariableName)
			}
		}

		// Add the body as the last argument
		if len(route.BodyType) > 0 {
			args = append(args, "body")
		}

		s = append(s, fmt.Sprintf("\t\tc.%s.%s.%s(", strings.Title(controller.Package), controller.Name, route.Name)+strings.Join(args, ", ")+")\n")

		s = append(s, "\t})).")

		s = append(s, fmt.Sprintf("\t\tMethods(\"%s\").", route.Method))

		if len(route.Queries) > 0 {
			s = append(s, "\t\tQueries(")
			for _, query := range route.Queries {
				s = append(s, fmt.Sprintf("\t\t\t\"%s\", \"%s\",", query.Name, query.ValueRaw))
			}
			s = append(s, "\t\t).")
		}

		s = append(s, fmt.Sprintf("\t\tName(\"%s\")\n", route.Name))
	}

	out = strings.Join(s, "\n") // + "\n}"

	return

}

func genDTOSMap(dir string) map[string]map[string]string {

	dirHandle, err := os.Open(dir)

	if err != nil {
		panic(err)
	}

	defer dirHandle.Close()

	var dirFileNames []string
	dirFileNames, err = dirHandle.Readdirnames(-1)

	if err != nil {
		panic(err)
	}
	// reader := bufio.NewReader(os.Stdin)

	result := map[string]map[string]string{}

	for _, name := range dirFileNames {

		if name == ".DS_Store" {
			continue
		}

		fullPath := path.Join(dir, name)

		model, e := ParseFileToGoStruct(fullPath)
		if e != nil {
			panic(e)
		}
		k := 0

		result[model.Name] = map[string]string{}

		for k < model.Fields.Len() {
			result[model.Name][model.Fields.Get(k).Name] = model.Fields.Get(k).DataType
			k++
		}
	}

	return result
}

func genModelsMap() map[string]map[string]string {

	dirHandle, err := os.Open(lib.ModelsGenDir)

	if err != nil {
		panic(err)
	}

	defer dirHandle.Close()

	var dirFileNames []string
	dirFileNames, err = dirHandle.Readdirnames(-1)

	if err != nil {
		panic(err)
	}
	// reader := bufio.NewReader(os.Stdin)

	result := map[string]map[string]string{}

	for _, name := range dirFileNames {

		if name == ".DS_Store" {
			continue
		}

		modelName := name[0 : len(name)-3]
		fullPath := path.Join(lib.ModelsGenDir, name)

		model, e := ParseFileToGoStruct(fullPath)
		if e != nil {
			panic(e)
		}
		k := 0

		result[modelName] = map[string]string{}

		for k < model.Fields.Len() {
			result[modelName][model.Fields.Get(k).Name] = model.Fields.Get(k).DataType
			k++
		}
	}

	return result
}

var structRegex = regexp.MustCompile("^type [a-zA-Z0-9]+ struct {$")

func genAggregatesMap(dir string) map[string]map[string]string {

	dirHandle, err := os.Open(dir)

	if err != nil {
		panic(err)
	}

	defer dirHandle.Close()

	var dirFileNames []string
	dirFileNames, err = dirHandle.Readdirnames(-1)

	if err != nil {
		panic(err)
	}
	// reader := bufio.NewReader(os.Stdin)

	result := map[string]map[string]string{}

	for _, name := range dirFileNames {

		if name == ".DS_Store" {
			continue
		}

		// fileNameNoExt := name[0 : len(name)-3]
		fullPath := path.Join(dir, name)
		// fmt.Println(fullPath)

		fileBytes, e := ioutil.ReadFile(fullPath)
		if e != nil {
			panic(e)
		}

		contents := string(fileBytes)

		contentLines := strings.Split(contents, "\n")
		currentStruct := ""
		for k := range contentLines {

			contentLines[k] = strings.TrimSpace(contentLines[k])

			if len(contentLines[k]) == 0 {
				continue
			}

			if structRegex.Match([]byte(contentLines[k])) {
				structName := contentLines[k][5 : len(contentLines[k])-9]
				// fmt.Println(k, structName)
				result[structName] = map[string]string{}
				currentStruct = structName
				continue
			}

			if len(currentStruct) > 0 {

				if contentLines[k] == "}" {
					currentStruct = ""
					continue
				}

				parts := []string{}
				preParts := strings.Split(contentLines[k], " ")
				for l := range preParts {
					if len(preParts[l]) == 0 {
						continue
					}

					parts = append(parts, preParts[l])
				}

				fieldName := ""
				fieldType := ""
				if len(parts) > 1 {
					fieldName = parts[0]
					fieldType = parts[1]
				} else {
					// This is an embedded type
					if strings.Contains(parts[0], ".") {
						// sParts := strings.Split(parts[0], ".")
						// fieldName = sParts[1]
						fieldType = parts[0]
						fieldName = "#embedded" + fmt.Sprint(k)
					}

					// fmt.Println(">>>> " + strings.TrimSpace(parts[0]))
				}

				if len(fieldName) > 0 && len(fieldType) > 0 {
					result[currentStruct][fieldName] = fieldType
				}

			}

			// if len(contentLines[k]) > 5 && contentLines[k][0:5] == "type " {
			// }
		}
		// model, e := InspectFile(fullPath)
		// if e != nil {
		// 	panic(e)
		// }
		// k := 0

		// for k < model.Fields.Len() {
		// 	result[model.Name][model.Fields.Get(k).Name] = model.Fields.Get(k).DataType
		// 	k++
		// }
	}

	return result
}

var goConstantRegex = regexp.MustCompile("^type [a-zA-Z0-9]+ [a-zA-Z0-9]+$")

func genConstantsMap() map[string][]string {

	files, err := ioutil.ReadDir(lib.CoreConstantsDir)

	if err != nil {
		panic(err)
	}

	result := map[string][]string{}

	for _, file := range files {

		if file.Name() == ".DS_Store" {
			continue
		}

		if file.IsDir() {
			continue
		}

		fullPath := path.Join(lib.CoreConstantsDir, file.Name())

		file, _ := os.Open(fullPath)

		key, value := getConstantsFromGoFile(file)

		result[key] = value

		file.Close()
	}

	return result
}

func getConstantsFromGoFile(file io.Reader) (string, []string) {

	scanner := bufio.NewScanner(file)

	var structName = ""
	isConsts := false
	var constants = []string{}

	for scanner.Scan() {

		var line = scanner.Text()

		if goConstantRegex.Match([]byte(line)) {
			structName = strings.Split(line[5:], " ")[0]
			continue
		}

		if line == "const (" {
			isConsts = true
			continue
		}

		if isConsts {

			line = strings.TrimSpace(line)

			if line == ")" {
				break
			}

			if len(line) > 2 && line[0:2] == "//" {
				continue
			}

			parts := strings.Split(line, " ")

			constants = append(constants, parts[0])

		}
	}

	return structName, constants

}
