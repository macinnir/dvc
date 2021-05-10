package gen

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/macinnir/dvc/core/lib"
)

// GenRoutes generates a list of routes from a directory of controller files
func GenRoutes(config *lib.Config) (e error) {

	permissionMap := loadPermissions()

	var files []os.FileInfo
	files, e = ioutil.ReadDir(config.Dirs.Controllers)

	if e != nil {
		log.Println("Error with path ", config.Dirs.Controllers, e.Error())
		return
	}

	imports := []string{
		path.Join(config.BasePackage, config.Dirs.API),
		path.Join(config.BasePackage, config.Dirs.Integrations),
		path.Join(config.BasePackage, config.Dirs.Aggregates),
		"net/http",
		"github.com/gorilla/mux",
	}

	// fmt.Println(imports)

	code := ""

	rest := ""
	controllerCalls := []string{}

	hasBodyImports := false
	usesPermissions := false

	allPerms := []string{}

	controllers := []*Controller{}

	for _, filePath := range files {

		fileName := filePath.Name()
		// Filter out files that don't have upper case first letter names
		if !unicode.IsUpper([]rune(fileName)[0]) {
			continue
		}

		// Skip non-go files
		if len(fileName) > 3 && fileName[len(fileName)-3:] != ".go" {
			continue
		}

		// Skip tests
		if len(fileName) > 8 && fileName[len(fileName)-8:] == "_test.go" {
			continue
		}

		var src []byte

		// log.Println(path.Join(dir, filePath.Name()))
		src, e = ioutil.ReadFile(path.Join(config.Dirs.Controllers, filePath.Name()))

		if e != nil {
			log.Println("Error with ", path.Join(config.Dirs.Controllers, filePath.Name()))
			return
		}

		// Build a controller object from the controller file
		controller, usesPerms, _ := BuildControllerObjFromControllerFile(path.Join(config.Dirs.Controllers, filePath.Name()), src)

		if usesPerms == true {
			usesPermissions = true
		}

		// Documentation routes
		controllers = append(controllers, controller)

		// Include imports for dtos and response if necessary for JSON http body
		if controller.HasDTOsImport == true {
			hasBodyImports = true
		}

		var routesString string
		var perms []string

		routesString, perms, e = BuildRoutesCodeFromController(permissionMap, controller)

		if e != nil {
			return
		}

		allPerms = append(allPerms, perms...)
		// allPerms = append(allPerms, "Manage"+controller.Name[0:len(controller.Name)-len("Controller")])

		rest += "\n" + routesString + "\n"

		controllerCalls = append(
			controllerCalls,
			"map"+extractNameFromFile(filePath.Name())+"Routes(res, r, auth, c, log)",
		)
	}

	code += strings.Join(controllerCalls, "\n\t")
	code += "\n\n}\n"
	code += rest

	if hasBodyImports {
		// imports = append(imports, g.Config.BasePackage+"/core/utils/response")
		imports = append(imports, config.BasePackage+"/core/definitions/dtos")
	}

	imports = append(imports, "github.com/macinnir/dvc/core/lib/utils/request")

	if usesPermissions {
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

	routesContainer := &RoutesJSONContainer{
		Routes:     map[string]*ControllerRoute{},
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

	if e = lib.EnsureDir("meta"); e != nil {
		return
	}

	routesJSON, _ := json.MarshalIndent(routesContainer, "  ", "    ")
	fmt.Println("Writing Routes JSON to to path", lib.RoutesFilePath)
	ioutil.WriteFile(lib.RoutesFilePath, routesJSON, 0777)

	return
}

// RoutesJSONContainer is a container for JSON Routes
type RoutesJSONContainer struct {
	Routes     map[string]*ControllerRoute  `json:"routes"`
	DTOs       map[string]map[string]string `json:"dtos"`
	Models     map[string]map[string]string `json:"models"`
	Aggregates map[string]map[string]string `json:"aggregates"`
	Constants  map[string][]string          `json:"constants"`
}

// Controller represents a REST controller
type Controller struct {
	Name              string             `json:"Name"`
	Description       string             `json:"Description"`
	Path              string             `json:"-"`
	Routes            []*ControllerRoute `json:"Routes"`
	HasDTOsImport     bool               `json:"-"`
	HasResponseImport bool               `json:"-"`
}

// ControllerRoute represents a route inside a REST controller
type ControllerRoute struct {
	Name           string                 `json:"Name"`
	Description    string                 `json:"Description"`
	Raw            string                 `json:"Path"`
	Path           string                 `json:"-"`
	Method         string                 `json:"Method"`
	Params         []ControllerRouteParam `json:"Params"`
	Queries        []ControllerRouteQuery `json:"Queries"`
	IsAuth         bool                   `json:"IsAuth"`
	BodyType       string                 `json:"BodyType"`
	BodyFormat     string                 `json:"BodyFormat"`
	HasBody        bool                   `json:"HasBody"`
	ResponseType   string                 `json:"ResponseType"`
	ResponseFormat string                 `json:"ResponseFormat"`
	ResponseCode   int                    `json:"ResponseCode"`
	Permission     string                 `json:"Permission"`
}

// ControllerRouteParam represents a param inside a controller route
type ControllerRouteParam struct {
	Name    string
	Pattern string
	Type    string
}

// ControllerRouteQuery represents a query inside a controller route
type ControllerRouteQuery struct {
	Name         string
	Pattern      string
	Type         string
	VariableName string
	ValueRaw     string
}

// DocRoute is a route in the documentation
type DocRoute struct {
	Name           string
	Description    string
	Method         string
	Path           string
	HasBody        bool
	BodyType       string
	BodyFormat     string
	ResponseType   string
	ResponseFormat string
	ResponseCode   int
}

func newDocRoute(route ControllerRoute) (docRoute *DocRoute) {
	docRoute = &DocRoute{
		Name:        route.Name,
		Description: route.Description,
		Method:      route.Method,
		Path:        route.Path,
		HasBody:     route.HasBody,
		BodyType:    route.BodyType,
	}

	return
}

// DocRouteParam represents a parameter inside a route for documentation
type DocRouteParam struct {
	Name    string
	Pattern string
	Type    string
}

// DocRouteQuery represents a query inside a route for documentation
type DocRouteQuery struct {
	Name    string
	Pattern string
	Type    string
}

// BuildControllerObjFromControllerFile parses a file and extracts all of its @route comments
func BuildControllerObjFromControllerFile(filePath string, src []byte) (controller *Controller, usesPerms bool, e error) {

	controller = &Controller{
		Name:   extractNameFromFile(filePath),
		Path:   filePath,
		Routes: []*ControllerRoute{},
	}

	// Get the controller name
	controllerName := extractNameFromFile(filePath)
	var methods []lib.Method
	methods, _, controller.Description = lib.ParseStruct(src, controllerName, true, true, "controllers")

	// Remove the name of the controller from the description
	controller.Description = strings.TrimPrefix(controller.Description, controller.Name)

	for _, method := range methods {

		route := &ControllerRoute{
			Queries: []ControllerRouteQuery{},
			Params:  []ControllerRouteParam{},
		}

		for line, doc := range method.Documents {

			// This is the title of the method
			if line == 0 {
				lineParts := strings.Split(doc, " ")
				route.Name = lineParts[1]
				route.Description = strings.Join(lineParts[2:], " ")
				continue
			}
			route.IsAuth = true

			// @anonymous
			if len(doc) > 12 && doc[0:13] == "// @anonymous" {
				route.IsAuth = false
				continue
			}

			// @body
			if len(doc) > 9 && doc[0:9] == "// @body " {

				bodyComment := strings.Split(strings.Trim(doc[9:], " "), " ")
				route.BodyFormat = bodyComment[0]

				if len(bodyComment) > 1 {
					route.BodyType = bodyComment[1]
				}

				controller.HasDTOsImport = true
				controller.HasResponseImport = true
				route.HasBody = true
				continue
			}

			// @response
			if len(doc) > 13 && doc[0:13] == "// @response " {

				responseComment := strings.Split(strings.Trim(doc[13:], " "), " ")

				if route.ResponseCode, e = strconv.Atoi(responseComment[0]); e != nil {
					log.Fatalf("Invalid @response comment: %s at %s.%s", doc, controllerName, route.Name)
				}

				if len(responseComment) > 1 {
					route.ResponseFormat = responseComment[1]
				}

				if len(responseComment) > 2 {
					route.ResponseType = responseComment[2]
				}

				continue
			}

			// @perm
			if len(doc) > 9 && doc[0:9] == "// @perm " {
				route.Permission = strings.TrimSpace(doc[9:])
				usesPerms = true
				continue
			}

			// @route
			if len(doc) > 9 && doc[0:9] == "// @route" {

				lineParts := strings.Split(doc, " ")
				if len(lineParts) < 4 {
					log.Fatalf("Invalid route comment `%s` for method `%s.%s`", doc, controllerName, route.Name)
				}
				route.Method = lineParts[2]
				route.Raw = lineParts[3]

				// Queries
				if strings.Contains(route.Raw, "?") {

					subParts := strings.Split(route.Raw, "?")
					route.Path = subParts[0]
					queries := strings.Split(subParts[1], "&")

					for _, query := range queries {
						if !strings.Contains(query, "=") {
							continue
						}

						queryParts := strings.Split(query, "=")

						o := ControllerRouteQuery{
							Name:     queryParts[0],
							ValueRaw: queryParts[1],
						}

						if strings.Contains(o.ValueRaw, ":") {
							queryValueParts := strings.Split(o.ValueRaw, ":")
							// Remove the starting "{"
							o.VariableName = queryValueParts[0][1:]

							// Remove the ending "}"
							o.Pattern = strings.Join(queryValueParts[1:], ":")
							o.Pattern = o.Pattern[0 : len(o.Pattern)-1]

							// Check if the value isn't a constant value
							if o.Pattern == "[0-9]" || o.Pattern == "[0-9]+" {
								o.Type = "int64"
							} else {
								o.Type = "string"
							}
						} else {
							// Try to parse the value as an int64
							// e.g. param=123

							o.VariableName = o.Name

							if _, e := strconv.ParseInt(o.ValueRaw, 10, 64); e != nil {
								o.Type = "string"
							} else {
								o.Type = "int64"
							}
						}

						route.Queries = append(route.Queries, o)
					}

				} else {
					route.Path = route.Raw
				}

				params, _ := extractParamsFromRoutePath(route.Path)

				route.Params = append(route.Params, params...)

			} else {
				route.Description += " " + doc[3:]
			}

		}

		controller.Routes = append(controller.Routes, route)
	}

	return
}

func extractParamsFromRoutePath(routePath string) (params []ControllerRouteParam, e error) {

	params = []ControllerRouteParam{}

	// Params
	if strings.Contains(routePath, "{") {

		routeParts := strings.Split(routePath, "{")

		for _, p := range routeParts[1:] {

			if !strings.Contains(p, "}") || !strings.Contains(p, ":") {
				continue
			}

			param := extractParamFromString(p)

			params = append(params, param)
		}
	}

	return
}

func extractParamFromString(paramString string) (param ControllerRouteParam) {

	// Incase there are parts after the param, split on the closing bracket
	pParts := strings.Split(paramString, "}")
	paramString = pParts[0]

	paramParts := strings.Split(paramString, ":")

	param = ControllerRouteParam{
		Name:    paramParts[0],
		Pattern: paramParts[1],
	}

	param.Type = matchPatternToDataType(param.Pattern)
	return
}

func matchPatternToDataType(pattern string) string {
	if pattern == "[0-9]" || pattern == "[0-9]+" {
		return "int64"
	}

	return "string"
}

// BuildRoutesCodeFromController builds controller code based on a route
func BuildRoutesCodeFromController(permissionMap map[string]string, controller *Controller) (out string, perms []string, e error) {

	permsGate := map[string]bool{}
	perms = []string{}

	s := []string{
		fmt.Sprintf("// map%sRoutes maps all of the routes for %s", controller.Name, controller.Name),
		fmt.Sprintf("func map%sRoutes(res request.IResponseLogger, r *mux.Router, auth integrations.IAuth, c *controllers.Controllers, log integrations.ILog) {\n", controller.Name),
	}

	for _, route := range controller.Routes {

		// Method comments
		s = append(s, fmt.Sprintf("\t// %s", route.Name))
		s = append(s, fmt.Sprintf("\t// %s %s", route.Method, route.Raw))
		if !route.IsAuth {
			s = append(s, "\t// @anonymous")
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
			permission := string(unicode.ToUpper(rune(route.Permission[0]))) + route.Permission[1:]
			if _, ok := permsGate[route.Permission]; !ok {
				// fmt.Println("adding permission " + route.Permission)
				permsGate[route.Permission] = true
				perms = append(perms, route.Permission)
			}

			if _, ok := permissionMap[permission]; !ok {
				e = errors.New("Permission " + permission + " does not exist")
				return
			}

			s = append(s, "\t\tif !utils.HasPerm(currentUser, permissions."+permission+") {")
			s = append(s, "\t\t\tres.Forbidden(req, w)")
			s = append(s, "\t\t\treturn")
			s = append(s, "\t\t}\n")
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

		s = append(s, fmt.Sprintf("\t\tc.%s.%s(", controller.Name, route.Name)+strings.Join(args, ", ")+")\n")

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

	out = strings.Join(s, "\n") + "\n}"

	return

}

func extractNameFromFile(fileName string) (name string) {
	baseName := filepath.Base(fileName)
	return baseName[0 : len(baseName)-3]
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

		// fileNameNoExt := name[0 : len(name)-3]
		fullPath := path.Join(dtosDir, name)
		// fmt.Println(fullPath)

		model, e := InspectFile(fullPath)
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

		model, e := InspectFile(fullPath)
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
