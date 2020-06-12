package gen

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
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

	code := `package main 
import (
	"` + g.Config.BasePackage + `/core/controllers"
	"` + g.Config.BasePackage + `/core/utils"

	"github.com/gorilla/mux"
)

// mapRoutesToControllers maps the routes to the controllers
func mapRoutesToControllers(r *mux.Router, auth *utils.Auth, c *controllers.Controllers) {

	`
	rest := ""
	controllerCalls := []string{}

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

		log.Println(path.Join(dir, filePath.Name()))
		src, e = ioutil.ReadFile(path.Join(dir, filePath.Name()))

		if e != nil {
			log.Println("Error with ", path.Join(dir, filePath.Name()))
			return
		}

		controller, _ := g.BuildControllerObjFromController(path.Join(dir, filePath.Name()), src)
		rest += "\n" + g.BuildRoutesCodeFromController(controller) + "\n"

		controllerCalls = append(
			controllerCalls,
			"map"+extractNameFromFile(filePath.Name())+"Routes(r, auth, c)",
		)
	}

	code += strings.Join(controllerCalls, "\n\t")
	code += "\n\n}\n"

	code += rest

	ioutil.WriteFile("services/api/routes.go", []byte(code), 0777)

	return
}

// Controller represents a REST controller
type Controller struct {
	Name   string
	Path   string
	Routes []ControllerRoute
}

// ControllerRoute represents a route inside a REST controller
type ControllerRoute struct {
	Name        string
	Description string
	Raw         string
	Path        string
	Method      string
	Params      []ControllerRouteParam
	Queries     []ControllerRouteQuery
	IsAuth      bool
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
							o.Pattern = queryValueParts[1][0 : len(queryValueParts[1])-1]

							if o.Pattern == "[0-9]" || o.Pattern == "[0-9]+" {
								o.Type = "int64"
							} else {
								o.Type = "string"
							}
						}

						route.Queries = append(route.Queries, o)
					}

				} else {
					route.Path = route.Raw
				}

				// Params
				if strings.Contains(route.Path, "{") {

					routeParts := strings.Split(route.Path, "{")

					for _, p := range routeParts[1:] {

						if !strings.HasSuffix(p, "}") || !strings.Contains(p, ":") {
							continue
						}

						paramParts := strings.Split(p, ":")

						param := ControllerRouteParam{
							Name:    paramParts[0],
							Pattern: paramParts[1][0 : len(paramParts[1])-1],
						}

						if param.Pattern == "[0-9]" || param.Pattern == "[0-9]+" {
							param.Type = "int64"
						} else {
							param.Type = "string"
						}

						route.Params = append(route.Params, param)
					}
				}

			} else {
				route.Description += " " + doc[3:]
			}

		}

		routes = append(routes, route)
	}

	return
}

// BuildControllerObjFromController builds a controller object from a controller file
func (g *Gen) BuildControllerObjFromController(filePath string, src []byte) (controller *Controller, e error) {

	controller = &Controller{
		Name: extractNameFromFile(filePath),
		Path: filePath,
	}

	controller.Routes, e = g.ExtractRoutesFromController(filePath, src)

	return
}

// BuildRoutesCodeFromController builds controller code based on a route
func (g *Gen) BuildRoutesCodeFromController(controller *Controller) (out string) {

	s := []string{
		fmt.Sprintf("// map%sRoutes maps all of the routes for %s", controller.Name, controller.Name),
		fmt.Sprintf("func map%sRoutes(r *mux.Router, auth *utils.Auth, c *controllers.Controllers) {", controller.Name),
	}

	for _, route := range controller.Routes {

		handleFunc := ""
		if route.IsAuth {
			handleFunc = fmt.Sprintf("r.Handle(\"%s\", auth.AuthMiddleware(c.%s.%s)).", route.Path, controller.Name, route.Name)
		} else {
			handleFunc = fmt.Sprintf("r.HandleFunc(\"%s\", c.%s.%s).", route.Path, controller.Name, route.Name)
		}

		r := []string{
			fmt.Sprintf("// %s", route.Name),
			handleFunc,
			fmt.Sprintf("\tMethods(\"%s\").", route.Method),
		}

		if len(route.Queries) > 0 {
			r = append(r, "\tQueries(")
			for _, query := range route.Queries {
				r = append(r, fmt.Sprintf("\t\t\"%s\", \"%s\",", query.Name, query.ValueRaw))
			}
			r = append(r, "\t).")

		}

		r = append(r, fmt.Sprintf("\tName(\"%s\")", route.Name))

		s = append(s, "\n\t"+strings.Join(r, "\n\t"))
	}

	out = strings.Join(s, "\n") + "\n}"

	return

}

func extractNameFromFile(fileName string) (name string) {
	baseName := filepath.Base(fileName)
	return baseName[0 : len(baseName)-3]
}
