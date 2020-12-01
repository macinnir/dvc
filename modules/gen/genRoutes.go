package gen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/macinnir/dvc/lib"
)

// GenRoutes generates a list of routes from a directory of controller files
func (g *Gen) GenRoutes() (e error) {

	dir := "core/controllers"

	var files []os.FileInfo
	files, e = ioutil.ReadDir(dir)

	if e != nil {
		log.Println("Error with path ", dir, e.Error())
		return
	}

	imports := []string{
		g.Config.BasePackage + "/core/controllers",
		g.Config.BasePackage + "/core/utils/request",
		g.Config.BasePackage + "/core/definitions/integrations",
		"net/http",
		"github.com/gorilla/mux",
	}

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
		src, e = ioutil.ReadFile(path.Join(dir, filePath.Name()))

		if e != nil {
			log.Println("Error with ", path.Join(dir, filePath.Name()))
			return
		}

		// Build a controller object from the controller file
		controller, usesPerms, _ := g.BuildControllerObjFromControllerFile(path.Join(dir, filePath.Name()), src)

		if usesPerms == true {
			usesPermissions = true
		}

		// Documentation routes
		controllers = append(controllers, controller)

		// Include imports for dtos and response if necessary for JSON http body
		if controller.HasDTOsImport == true {
			hasBodyImports = true
		}

		routesString, perms := g.BuildRoutesCodeFromController(controller)

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
		imports = append(imports, g.Config.BasePackage+"/core/definitions/dtos")
	}

	if usesPermissions {
		imports = append(imports, "github.com/macinnir/dvc/modules/utils")
		imports = append(imports, g.Config.BasePackage+"/core/definitions/constants/permissions")
	}

	final := `// Generated Code; DO NOT EDIT.

package main

import (
`

	for _, i := range imports {
		final += fmt.Sprintf("\t\"%s\"\n", i)
	}

	final += `)

// mapRoutesToControllers maps the routes to the controllers
func mapRoutesToControllers(r *mux.Router, auth integrations.IAuth, c *controllers.Controllers, res integrations.IResponseLogger, log integrations.ILog) {

	`
	final += code

	ioutil.WriteFile("services/api/routes.go", []byte(final), 0777)

	routesContainer := &RoutesJSONContainer{
		Routes: controllers,
		DTOs:   genDTOSMap(),
		Models: genModelsMap(),
	}

	if e = lib.EnsureDir("meta"); e != nil {
		return
	}

	routesJSON, _ := json.MarshalIndent(routesContainer, "  ", "    ")
	routesJSONFilePath := "meta/routes.json"
	fmt.Println("Writing Routes JSON to to path", routesJSONFilePath)
	ioutil.WriteFile(routesJSONFilePath, routesJSON, 0777)

	g.buildPermissions()

	return
}

// RoutesJSONContainer is a container for JSON Routes
type RoutesJSONContainer struct {
	Routes []*Controller                `json:"routes"`
	DTOs   map[string]map[string]string `json:"dtos"`
	Models map[string]map[string]string `json:"models"`
}

// Controller represents a REST controller
type Controller struct {
	Name              string            `json:"Name"`
	Description       string            `json:"Description"`
	Path              string            `json:"-"`
	Routes            []ControllerRoute `json:"Routes"`
	HasDTOsImport     bool              `json:"-"`
	HasResponseImport bool              `json:"-"`
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
func (g *Gen) BuildControllerObjFromControllerFile(filePath string, src []byte) (controller *Controller, usesPerms bool, e error) {

	controller = &Controller{
		Name:   extractNameFromFile(filePath),
		Path:   filePath,
		Routes: []ControllerRoute{},
	}

	// Get the controller name
	controllerName := extractNameFromFile(filePath)
	var methods []lib.Method
	methods, _, controller.Description = lib.ParseStruct(src, controllerName, true, true, "controllers")

	// Remove the name of the controller from the description
	controller.Description = strings.TrimPrefix(controller.Description, controller.Name)

	for _, method := range methods {

		route := ControllerRoute{
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

func (g *Gen) buildPermissions() {

	permissionMap := map[string]string{}
	fileBytes, e := ioutil.ReadFile("meta/permissions.json")
	if e != nil {
		panic(e)
	}
	json.Unmarshal(fileBytes, &permissionMap)

	permissions := make([]string, 0, len(permissionMap))
	for k := range permissionMap {
		permissions = append(permissions, k)
	}

	sort.Strings(permissions)

	// for k := range permissions {

	// 	permission := permissions[k]
	// 	description := permissionMap[permissions[k]]

	// 	// fmt.Println(permission + ": " + description)

	// }

	permissionsFile := `// Generated Code; DO NOT EDIT.

	package permissions

	import (
		"github.com/macinnir/dvc/modules/utils"
	)
	
	const (
	`
	for k := range permissions {
		permTitle := string(unicode.ToUpper(rune(permissions[k][0]))) + permissions[k][1:]
		permissionsFile += "\t// " + permTitle + "Permission is the `" + permissions[k] + "` permission\n"
		permissionsFile += "\t" + permTitle + " utils.Permission = \"" + permissions[k] + "\"\n"
	}

	permissionsFile += `)
	
	// Permissions returns a slice of permissions 
	func Permissions() map[utils.Permission]string {
		return map[utils.Permission]string {
	`

	for k := range permissions {
		permTitle := string(unicode.ToUpper(rune(permissions[k][0]))) + permissions[k][1:]
		permissionsFile += "\t\t" + permTitle + ": \"" + permissionMap[permissions[k]] + "\",\n"
	}

	permissionsFile += `	}
	}	
	
	`
	permissionsFilePath := path.Join("core", "definitions", "constants", "permissions", "permissions.go")
	fmt.Println("Writing the permissions file to path ", permissionsFilePath)
	var permissionsFileBytes []byte
	permissionsFileBytes, e = lib.FormatCode(permissionsFile)
	if e != nil {
		panic(e)
	}
	ioutil.WriteFile(permissionsFilePath, []byte(permissionsFileBytes), 0777)
}

// BuildTypescriptPermissions returns a formatted typescript file of permission constants
func (g *Gen) BuildTypescriptPermissions() string {

	permissionMap := map[string]string{}
	fileBytes, e := ioutil.ReadFile("meta/permissions.json")
	if e != nil {
		panic(e)
	}
	json.Unmarshal(fileBytes, &permissionMap)

	permissions := make([]string, 0, len(permissionMap))
	for k := range permissionMap {
		permissions = append(permissions, k)
	}

	sort.Strings(permissions)

	permissionsFile := "// Generated Code; DO NOT EDIT.\n\n"
	for k := range permissions {
		permTitle := string(unicode.ToUpper(rune(permissions[k][0]))) + permissions[k][1:]
		permissionsFile += "// " + permTitle + " -- " + permissionMap[permissions[k]] + "\n"
		permissionsFile += "export const " + permTitle + "Permission = \"" + permissions[k] + "\";\n"
	}

	return permissionsFile
}

// BuildRoutesCodeFromController builds controller code based on a route
func (g *Gen) BuildRoutesCodeFromController(controller *Controller) (out string, perms []string) {

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

			// for k := range route.Perms {
			// s = append(s, "\t\tif !utils.HasPerm(currentUser, \""+route.Perms[k]+"\") {")
			s = append(s, "\t\tif !utils.HasPerm(currentUser.UserID, currentUser.PermissionNames, permissions."+permission+") {")
			s = append(s, "\t\t\tres.Forbidden(req, w)")
			s = append(s, "\t\t\treturn")
			s = append(s, "\t\t}\n")
			// }
		}

		if len(route.BodyType) > 0 {
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

		if route.IsAuth {
			s = append(s, "\t})).")
		} else {
			s = append(s, "\t}).")
		}

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
