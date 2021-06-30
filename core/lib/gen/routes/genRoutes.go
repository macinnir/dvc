package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/fetcher"
	"github.com/macinnir/dvc/core/lib/gen"
)

// GenRoutes generates a list of routes from a directory of controller files
func GenRoutesAndPermissions(config *lib.Config) error {

	imports := []string{
		path.Join(config.BasePackage, config.Dirs.Controllers),
		path.Join(config.BasePackage, config.Dirs.IntegrationInterfaces),
		path.Join(config.BasePackage, config.Dirs.Aggregates),
		"net/http",
		"github.com/gorilla/mux",
	}

	// fmt.Println(imports)

	code := ""

	rest := ""
	// controllerCalls := []string{}

	hasBodyImports := false
	packageUsesPermission := false

	cf := fetcher.NewControllerFetcher()
	controllers, e := cf.Fetch(config.Dirs.Controllers)

	if e != nil {
		return e
	}

	for k := range controllers {

		controller := controllers[k]
		// fmt.Println("ControllerName:", controllerName)

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
		imports = append(imports, config.BasePackage+"/core/definitions/dtos")
	}

	imports = append(imports, "github.com/macinnir/dvc/core/lib/utils/request")

	if packageUsesPermission {
		imports = append(imports, "github.com/macinnir/dvc/core/lib/utils")
		imports = append(imports, path.Join(config.BasePackage, config.Dirs.Permissions))
	}

	final := `// Generated Code; DO NOT EDIT.

package api

import (
`

	for _, i := range imports {
		final += fmt.Sprintf("\t\"%s\"\n", i)
	}

	final += `)

// MapRoutesToControllers maps the routes to the controllers
func MapRoutesToControllers(r *mux.Router, auth integrations.IAuth, c *controllers.Controllers, res request.IResponseLogger, log integrations.ILog) {

	`
	final += code

	ioutil.WriteFile("core/api/routes.go", []byte(final), 0777)

	routesContainer := &lib.RoutesJSONContainer{
		Routes:     map[string]*lib.ControllerRoute{},
		DTOs:       genDTOSMap(),
		Models:     genModelsMap(),
		Aggregates: genAggregatesMap(),
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

		// fmt.Println("Route: " + route.Name)

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
			s = append(s, fmt.Sprintf("\tr.Handle(\"%s\", auth.AuthMiddleware(func(w http.ResponseWriter, currentUser *aggregates.UserAggregate, req *request.Request) {\n", route.Path))
			// s = append(s, fmt.Sprintf("\t\tcurrentUser := auth.GetCurrentUser(r)\n"))
			// args = append(args, "currentUser")
		} else {
			s = append(s, fmt.Sprintf("\tr.Handle(\"%s\", auth.AnonMiddleware(func(w http.ResponseWriter, req *request.Request) {\n", route.Path))
		}

		s = append(s, fmt.Sprintf("\n\t\tlog.Debug(\"ROUTE: %s %s => %s\")\n\n", route.Method, route.Path, route.Name))

		// Permission
		if route.IsAuth && len(route.Permission) > 0 {

			// ucFirst
			// permission := string(unicode.ToUpper(rune(route.Permission[0]))) + route.Permission[1:]

			s = append(s, `		
		if !utils.HasPerm(req, currentUser, permissions.`+route.Permission+`) {
			res.Forbidden(req, w)
			return
		}`)

		}

		if len(route.BodyType) > 0 {

			if route.BodyType[0:1] == "*" {
				route.BodyType = route.BodyType[1:]
			}

			s = append(s, fmt.Sprintf("\t\tbody := &%s{}", route.BodyType))
			s = append(s, "\t\treq.BodyJSON(body)\n")
		}

		if len(route.Params) > 0 {
			for _, param := range route.Params {
				s = append(s, fmt.Sprintf("\t\t// URL Param %s", param.Name))
				if param.Type == "int64" {
					s = append(s, fmt.Sprintf("\t\t%s := req.ArgInt64(\"%s\", 0)\n", param.Name, param.Name))
				} else {
					s = append(s, fmt.Sprintf("\t\t%s := req.Arg(\"%s\", \"\")\n", param.Name, param.Name))
				}

				args = append(args, param.Name)
			}
		}

		if len(route.Queries) > 0 {
			for _, query := range route.Queries {
				s = append(s, fmt.Sprintf("\t\t// Query Arg %s", query.VariableName))
				if query.Type == "int64" {
					s = append(s, fmt.Sprintf("\t\t%s := req.ArgInt64(\"%s\", 0)\n", query.VariableName, query.VariableName))
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

func genDTOSMap() map[string]map[string]string {

	dtosDir := "core/definitions/dtos"

	dirHandle, err := os.Open(dtosDir)

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

		fullPath := path.Join(dtosDir, name)

		model, e := gen.InspectFile(fullPath)
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

	modelsDir := "core/definitions/models"

	dirHandle, err := os.Open(modelsDir)

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
		fullPath := path.Join(modelsDir, name)
		// fmt.Println(fullPath)

		model, e := gen.InspectFile(fullPath)
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

func genAggregatesMap() map[string]map[string]string {

	modelsDir := "core/definitions/aggregates"

	dirHandle, err := os.Open(modelsDir)

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
		fullPath := path.Join(modelsDir, name)
		// fmt.Println(fullPath)

		fileBytes, e := ioutil.ReadFile(fullPath)
		if e != nil {
			panic(e)
		}

		contents := string(fileBytes)

		re := regexp.MustCompile("^type [a-zA-Z0-9]+ struct {$")
		contentLines := strings.Split(contents, "\n")
		currentStruct := ""
		for k := range contentLines {

			if re.Match([]byte(contentLines[k])) {
				structName := contentLines[k][5 : len(contentLines[k])-9]
				// fmt.Println(k, structName)
				result[structName] = map[string]string{}
				currentStruct = structName
				continue
			}

			if len(currentStruct) > 0 {

				contentLines[k] = strings.TrimSpace(contentLines[k])

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
						sParts := strings.Split(parts[0], ".")
						fieldName = sParts[1]
						fieldType = parts[0]
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

func genConstantsMap() map[string][]string {

	modelsDir := "core/definitions/constants"

	files, err := ioutil.ReadDir(modelsDir)

	if err != nil {
		panic(err)
	}

	// defer dirHandle.Close()

	// var dirFileNames []string
	// dirFileNames, err = dirHandle.Readdirnames(-1)

	// reader := bufio.NewReader(os.Stdin)

	result := map[string][]string{}

	for _, file := range files {

		if file.Name() == ".DS_Store" {
			continue
		}

		if file.IsDir() {
			continue
		}

		// fileNameNoExt := name[0 : len(name)-3]
		fullPath := path.Join(modelsDir, file.Name())
		// fmt.Println(fullPath)

		fileBytes, e := ioutil.ReadFile(fullPath)

		if e != nil {
			panic(e)
		}

		contents := string(fileBytes)

		re := regexp.MustCompile("^type [a-zA-Z0-9]+ [a-zA-Z0-9]+$")
		contentLines := strings.Split(contents, "\n")
		currentStruct := ""
		isConsts := false
		for k := range contentLines {
			// fmt.Println(k, contentLines[k])
			if re.Match([]byte(contentLines[k])) {
				structName := contentLines[k][5:]
				structName = strings.Split(structName, " ")[0]
				// fmt.Println(k, structName)
				result[structName] = []string{}
				currentStruct = structName
				continue
			}

			if contentLines[k] == "const (" {
				isConsts = true
				continue
			}

			if isConsts == true {
				contentLines[k] = strings.TrimSpace(contentLines[k])
				if contentLines[k] == ")" {
					break
				}

				if len(contentLines[k]) > 2 && contentLines[k][0:2] == "//" {
					continue
				}

				parts := strings.Split(contentLines[k], " ")
				result[currentStruct] = append(result[currentStruct], parts[0])
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
