package gen

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

	"github.com/macinnir/dvc/lib"
)

// GenRoutes generates a list of routes from a directory of controller files
func (g *Gen) GenRoutes(dir string) (e error) {

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

		controller, _ := g.BuildControllerObjFromController(path.Join(dir, filePath.Name()), src)

		// Include imports for dtos and response if necessary for JSON http body
		if controller.HasDTOsImport == true {
			hasBodyImports = true
		}

		rest += "\n" + g.BuildRoutesCodeFromController(controller) + "\n"

		controllerCalls = append(
			controllerCalls,
			"map"+extractNameFromFile(filePath.Name())+"Routes(r, auth, c)",
		)
	}

	code += strings.Join(controllerCalls, "\n\t")
	code += "\n\n}\n"
	code += rest

	if hasBodyImports {
		imports = append(imports, g.Config.BasePackage+"/core/utils/response")
		imports = append(imports, g.Config.BasePackage+"/core/definitions/dtos")
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
func mapRoutesToControllers(r *mux.Router, auth integrations.IAuth, c *controllers.Controllers) {

	`
	final += code

	ioutil.WriteFile("services/api/routes.go", []byte(final), 0777)

	return
}

// Controller represents a REST controller
type Controller struct {
	Name              string
	Path              string
	Routes            []ControllerRoute
	HasDTOsImport     bool
	HasResponseImport bool
}

// ControllerRoute represents a route inside a REST controller
type ControllerRoute struct {
	Name           string
	Description    string
	Raw            string
	Path           string
	Method         string
	Params         []ControllerRouteParam
	Queries        []ControllerRouteQuery
	IsAuth         bool
	BodyType       string
	BodyFormat     string
	HasBody        bool
	ResponseType   string
	ResponseFormat string
	ResponseCode   int
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

// ExtractRoutesFromController parses a file and extracts all of its @route comments
func (g *Gen) ExtractRoutesFromController(filePath string, src []byte) (routes []ControllerRoute, e error) {

	routes = []ControllerRoute{}

	// Get the controller name
	controllerName := extractNameFromFile(filePath)
	methods, _, _ := lib.ParseStruct(src, controllerName, true, true, "controllers")

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

			if doc == "// @auth" {
				route.IsAuth = true
				continue
			}

			if len(doc) > 8 && doc[0:9] == "// @body " {
				bodyComment := strings.Split(strings.Trim(doc[9:], " "), " ")
				route.BodyFormat = bodyComment[0]

				if len(bodyComment) > 1 {
					route.BodyType = bodyComment[1]
				}

				route.HasBody = true
				continue
			}

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

							if o.Pattern == "[0-9]" || o.Pattern == "[0-9]+" {
								o.Type = "int64"
							} else {
								o.Type = "string"
							}
						} else {
							o.VariableName = o.Name
							o.Type = "string"
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

		routes = append(routes, route)
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

// BuildControllerObjFromController builds a controller object from a controller file
func (g *Gen) BuildControllerObjFromController(filePath string, src []byte) (controller *Controller, e error) {

	controller = &Controller{
		Name: extractNameFromFile(filePath),
		Path: filePath,
	}

	controller.Routes, e = g.ExtractRoutesFromController(filePath, src)

	for _, r := range controller.Routes {
		if r.HasBody {
			controller.HasDTOsImport = true
			controller.HasResponseImport = true
			break
		}
	}

	return
}

// BuildRoutesCodeFromController builds controller code based on a route
func (g *Gen) BuildRoutesCodeFromController(controller *Controller) (out string) {

	s := []string{
		fmt.Sprintf("// map%sRoutes maps all of the routes for %s", controller.Name, controller.Name),
		fmt.Sprintf("func map%sRoutes(r *mux.Router, auth integrations.IAuth, c *controllers.Controllers) {\n", controller.Name),
	}

	for _, route := range controller.Routes {

		// Method comments
		s = append(s, fmt.Sprintf("\t// %s", route.Name))
		s = append(s, fmt.Sprintf("\t// %s %s", route.Method, route.Raw))

		// Method args
		args := []string{
			"w", // http.ResponseWriter
			"r", // *http.Request
		}

		if route.IsAuth {
			s = append(s, fmt.Sprintf("\tr.Handle(\"%s\", auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {\n", route.Path))
			s = append(s, fmt.Sprintf("\t\tcurrentUser := auth.GetCurrentUser(r)\n"))
			args = append(args, "currentUser")
		} else {
			s = append(s, fmt.Sprintf("\tr.HandleFunc(\"%s\", func(w http.ResponseWriter, r *http.Request) {\n", route.Path))
		}

		if len(route.BodyType) > 0 {
			s = append(s, fmt.Sprintf("\t\tbody := &%s{}", route.BodyType))
			s = append(s, "\t\te := request.GetBodyJSON(r, body)\n")
			s = append(s, "\t\tif e != nil {")
			s = append(s, "\t\t\tresponse.BadRequest(r, w, e)")
			s = append(s, "\t\t\treturn")
			s = append(s, "\t\t}\n")
		}

		if len(route.Params) > 0 {
			for _, param := range route.Params {
				s = append(s, fmt.Sprintf("\t\t// URL Param %s", param.Name))
				if param.Type == "int64" {
					s = append(s, fmt.Sprintf("\t\t%s := request.URLParamInt64(r, \"%s\", 0)\n", param.Name, param.Name))
				} else {
					s = append(s, fmt.Sprintf("\t\t%s := request.URLParamString(r, \"%s\", \"\")\n", param.Name, param.Name))
				}

				args = append(args, param.Name)
			}
		}

		if len(route.Queries) > 0 {
			for _, query := range route.Queries {
				s = append(s, fmt.Sprintf("\t\t// Query Arg %s", query.VariableName))
				if query.Type == "int64" {
					s = append(s, fmt.Sprintf("\t\t%s := request.QueryArgInt64(r, \"%s\", 0)\n", query.VariableName, query.VariableName))
				} else {
					s = append(s, fmt.Sprintf("\t\t%s := request.QueryArgString(r, \"%s\", \"\")\n", query.VariableName, query.VariableName))
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
