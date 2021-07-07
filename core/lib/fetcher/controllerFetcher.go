package fetcher

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/macinnir/dvc/core/lib"
)

const (
	RouteArgTypeString = "string"
	RouteArgTypeNumber = "int64"

	// Route Tags
	RouteTagAnonymous = "@anonymous"
	RouteTagAnyone    = "@anyone"
	RouteTagBody      = "@body"
	RouteTagResponse  = "@response"
	RouteTagRoute     = "@route"
)

type ControllerFetcher struct {
	routeMap map[string]bool
}

func NewControllerFetcher() *ControllerFetcher {
	return &ControllerFetcher{
		routeMap: map[string]bool{},
	}
}

func (cf *ControllerFetcher) FetchAll() (controllers []*lib.Controller, dirs []string, e error) {

	controllers, dirs, e = cf.Fetch("core/api")

	if e != nil {
		return
	}

	appControllers, appDirs, e := cf.Fetch("app/api")
	if e != nil {
		return
	}
	controllers = append(controllers, appControllers...)
	dirs = append(dirs, appDirs...)
	return
}

func (cf *ControllerFetcher) Fetch(dir string) (controllers []*lib.Controller, dirs []string, e error) {

	controllers = []*lib.Controller{}
	dirs = []string{}

	var files []os.FileInfo
	files, e = ioutil.ReadDir(dir)

	if e != nil {
		log.Println("ERROR: Fetch Controllers - ", dir, e.Error())
		return
	}

	for k := range files {

		filePath := path.Join(dir, files[k].Name())

		if files[k].IsDir() {
			dirs = append(dirs, filePath)
			var subControllers []*lib.Controller
			if subControllers, _, e = cf.Fetch(filePath); e != nil {
				return
			}
			controllers = append(controllers, subControllers...)
			continue
		}

		// Build a controller object from the controller file
		var controller *lib.Controller
		if controller, e = cf.BuildControllerObjFromControllerFile(filePath); e != nil {
			return
		}

		if controller != nil {
			controllers = append(controllers, controller)
		}
	}

	return

}

// BuildControllerObjFromControllerFile parses a file and extracts all of its @route comments
func (cf *ControllerFetcher) BuildControllerObjFromControllerFile(filePath string) (controller *lib.Controller, e error) {

	pkgName := filepath.Base(filepath.Dir(filePath))

	controllerName := extractControllerNameFromFileName(filePath)

	if controllerName == "" {
		return nil, nil
	}

	var src []byte

	src, e = ioutil.ReadFile(filePath)

	if e != nil {
		log.Println("Error with ", filePath)
		return
	}

	controller = &lib.Controller{
		Name:    controllerName,
		Path:    filePath,
		Routes:  []*lib.ControllerRoute{},
		Package: pkgName,
	}
	controllerFullName := controller.Name + "Controller"

	// Get the controller name
	var methods []lib.Method
	methods, _, controller.Description = lib.ParseStruct(src, controllerFullName, true, true, "controllers")

	// Remove the name of the controller from the description
	controller.Description = strings.TrimPrefix(controller.Description, controller.Name)

	for k := range methods {

		route := parseControllerMethod(methods[k], controller)

		routeSignature := route.Method + " " + route.Path
		if _, ok := cf.routeMap[routeSignature]; ok {
			e = fmt.Errorf("Duplicate route signature `%s` for method `%s`.`%s`", routeSignature, controller.Name, route.Name)
			return
		}

		controller.Routes = append(controller.Routes, route)
	}

	return
}

func parseControllerMethod(
	method lib.Method,
	controller *lib.Controller,
) *lib.ControllerRoute {

	var e error

	route := &lib.ControllerRoute{
		Queries: []lib.ControllerRouteQuery{},
		Params:  []lib.ControllerRouteParam{},
	}

	isAnyone := false

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
		if len(doc) > 12 && doc[0:13] == "// "+RouteTagAnonymous {
			route.IsAuth = false
			continue
		}

		// @anyone
		if len(doc) > 9 && doc[0:10] == "// "+RouteTagAnyone {
			isAnyone = true
			continue
		}

		// @body
		if len(doc) > 9 && doc[0:9] == fmt.Sprintf("// %s ", RouteTagBody) {

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

		// @response (last line)
		if len(doc) > 13 && doc[0:13] == fmt.Sprintf("// %s ", RouteTagResponse) {

			responseComment := strings.Split(strings.Trim(doc[13:], " "), " ")

			if route.ResponseCode, e = strconv.Atoi(responseComment[0]); e != nil {
				log.Fatalf("Invalid @response comment: %s at %s.%s", doc, controller.Name, route.Name)
			}

			if len(responseComment) > 1 {
				route.ResponseFormat = responseComment[1]
			}

			if len(responseComment) > 2 {
				route.ResponseType = responseComment[2]
			}

			continue
		}

		var e error
		// @route
		if len(doc) > 9 && doc[0:9] == "// "+RouteTagRoute {

			if e = parseRouteString(route, doc); e != nil {
				log.Fatalf("method `%s.%s`: `%s`", controller.Name, route.Name, e.Error())
			}

		} else {
			route.Description += " " + doc[3:]
		}
	}

	if route.IsAuth && !isAnyone {
		controller.PermCount++
		route.Permission = controller.Name + "_" + route.Name
	}

	return route
}

func extractControllerNameFromFileName(path string) string {

	fileName := filepath.Base(path)

	// Must be aleast 14 chars (e.g. AController.go)

	// fmt.Println("Extracting controller name", fileName)

	if
	// 14 chars
	len(fileName) < 14 ||
		// .go extension
		fileName[len(fileName)-3:] != ".go" ||
		// Uppercase first letter
		!unicode.IsUpper([]rune(fileName)[0]) ||
		// Not a test file
		fileName[len(fileName)-8:] == "_test.go" {
		return ""
	}

	return fileName[:len(fileName)-13]
}

func parseRouteString(route *lib.ControllerRoute, routeString string) error {

	lineParts := strings.Split(routeString, " ")
	if len(lineParts) < 4 {
		return fmt.Errorf("invalid route comment `%s`", routeString)
	}
	route.Method = lineParts[2]
	route.Raw = lineParts[3]

	route.Path, route.Queries = extractQueriesFromRoutePath(route.Raw)

	params, _ := extractParamsFromRoutePath(route.Path)
	route.Params = append(route.Params, params...)
	return nil
}

func extractParamsFromRoutePath(routePath string) (params []lib.ControllerRouteParam, e error) {

	params = []lib.ControllerRouteParam{}

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

func extractParamFromString(paramString string) (param lib.ControllerRouteParam) {

	// Incase there are parts after the param, split on the closing bracket
	pParts := strings.Split(paramString, "}")
	paramString = pParts[0]

	paramParts := strings.Split(paramString, ":")

	param = lib.ControllerRouteParam{
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

func extractQueriesFromRoutePath(routePath string) (string, []lib.ControllerRouteQuery) {

	queries := []lib.ControllerRouteQuery{}

	if !strings.Contains(routePath, "?") {
		return routePath, queries
	}

	subParts := strings.Split(routePath, "?")
	routePath = subParts[0]
	queryStrings := strings.Split(subParts[1], "&")

	for _, query := range queryStrings {

		if !strings.Contains(query, "=") {
			continue
		}

		queryParts := strings.Split(query, "=")

		o := lib.ControllerRouteQuery{
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

		queries = append(queries, o)
	}

	return routePath, queries
}
