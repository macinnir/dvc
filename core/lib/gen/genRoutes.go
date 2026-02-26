package gen

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

var goControllerBootstrapTemplate = template.Must(template.New("go-controllers-bootstrap-file").Parse(`// Generated Code; DO NOT EDIT.

package routes 

import(
	{{ range .Imports}}"{{.}}"
	{{end}}
)


// Controllers is the main container for all of the controller modules 
type Controllers struct {
	{{ range .Controllers }}
	{{.Title}} *{{.Name}}.Controllers{{end}}
}

// NewControllers bootstraps all of the controller modules 
func NewControllers(s *app.App, r request.IResponseLogger) *Controllers {
	return &Controllers {
		{{ range .Controllers }}
		{{.Title}}: {{.Name}}.NewControllers(s, r),{{ end }} 
	}
}

`))

type GoControllerBootstrapTemplateValues struct {
	Imports     []string
	Controllers []struct {
		Title string
		Name  string
	}
}

// 0.085753

// GenControllerBootstrap generates the bootstrap file for the applications controllers
func GenControllerBootstrap(basePackageName string, dirs []string) string {

	var vals = GoControllerBootstrapTemplateValues{
		Imports: make([]string, len(dirs)+2),
		Controllers: make([]struct {
			Title string
			Name  string
		}, len(dirs)),
	}

	for k := range dirs {
		packageName := path.Base(dirs[k])
		vals.Imports[k] = path.Join(basePackageName, dirs[k])
		vals.Controllers[k] = struct {
			Title string
			Name  string
		}{
			strings.ToUpper(packageName[0:1]) + packageName[1:],
			packageName,
		}
	}

	vals.Imports[len(dirs)] = path.Join(basePackageName, "core/app")
	vals.Imports[len(dirs)+1] = "github.com/macinnir/dvc/core/lib/utils/request"

	var buf bytes.Buffer

	goControllerBootstrapTemplate.Execute(&buf, vals)

	return buf.String()
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

var goRoutesTemplate = template.Must(
	template.New("go-routes-file").
		Funcs(template.FuncMap{
			"hasPrefix": strings.HasPrefix,
		}).
		Parse(`// Generated Code; DO NOT EDIT.

package routes

import (
	{{range .Imports}}"{{.}}"
	{{end}}
	appdtos "{{.AppDTOsPath}}"
)

// MapRoutesToControllers maps the routes to the controllers
func MapRoutesToControllers(r *mux.Router, app *app.App, res request.IResponseLogger) {

	var auth = app.Auth
	var log = app.Utils.Logger
	var c = NewControllers(app, res)

	{{range .Controllers}}
	////
	// {{.Package}}.{{.Name}}
	////

	{{range .Routes}}
	// {{.Package}}.{{.Controller}}.{{.Name}}
	// {{.Method}} {{.Raw}}{{if not .IsAuth}}
	// @anonymous{{else}}{{if eq (len .Permission) 0}}
	// @anyone{{else}}
	// @permission {{.Permission}}{{end}}{{end}}
	r.Handle("{{.Path}}", {{ if .IsAuth }}auth.AuthMiddleware({{ if gt (len .Permission) 0 }}permissions.{{.Permission}}{{ else }}""{{end}},{{ else }}auth.AnonMiddleware({{ end }}func(w http.ResponseWriter, req *request.Request) {

		log.Debug("ROUTE: {{.Method}} {{.Path}} => {{.Name}}")
{{if gt (len .BodyType) 0}}
		// Parse the body of type {{.BodyType}} 
		var body = &{{.BodyType}}{}
		req.BodyJSON(body){{end}}
		
		{{range .Params}}
		// URL Param {{.Name}}
		// {{.Pattern}}{{if eq .Type "int64"}}
		var {{.Name}} = req.ArgInt64("{{.Name}}", {{ if eq .Pattern "-?[0-9]+" }}-1{{else}}0{{end}})
		{{else}}
		var {{.Name}} = req.Arg("{{.Name}}", "")
		{{end}}{{end}}
		{{range .Queries}}
		// Query Arg {{.VariableName}}
		// {{.Pattern}}{{if eq .Type "int64"}}
		var {{.VariableName}} = req.ArgInt64("{{.VariableName}}",{{ if eq .Pattern "-?[0-9]+" }}-1{{else}}0{{end}})
	{{else}}
		{{if hasPrefix .Pattern "-?" }}var {{.VariableName}} = req.RootRequest.URL.Query().Get("{{.VariableName}}")
		{{else}}var {{.VariableName}} = req.Arg("{{.VariableName}}", "")
		{{end}}
	{{end}}{{end}}
		c.{{.Package}}.{{.Controller}}.{{.Name}}(w, req{{range .Params}}, {{.Name}}{{end}}{{range .Queries}}, {{.VariableName}}{{end}}{{if gt (len .BodyType) 0}}, body{{end}})
	})).
		Methods("{{.Method}}").
		{{if gt (len .RequiredQueries) 0}}
			Queries({{range .RequiredQueries}}
				"{{.VariableName}}", "{{.ValueRaw}}",{{end}}			
			).
		{{end}}
		Name("{{.Name}}")
	{{end}}
	{{end}}

}
`))

type RoutesTplValues struct {
	Imports     []string
	AppDTOsPath string
	Controllers []*lib.Controller
}

// GenRoutesAndPermissions generates a list of routes from a directory of controller files
// 0.0824
// 0.008664
func GenRoutesAndPermissions(schemaList *schema.SchemaList, controllers []*lib.Controller, dirs []string, config *lib.Config) (*lib.RoutesJSONContainer, error) {

	var buf bytes.Buffer
	// var start = time.Now()
	var e error

	var routesTplValues = RoutesTplValues{
		Imports: []string{
			// path.Join(config.BasePackage, config.Dirs.IntegrationInterfaces),
			// path.Join(config.BasePackage, config.Dirs.Aggregates),
			"net/http",
			"github.com/gorilla/mux",
			lib.LibRequests,
			path.Join(config.BasePackage, "core/app"),
			path.Join(config.BasePackage, lib.CoreDTOsDir),
			path.Join(config.BasePackage, lib.GoPermissionsDir),
		},
		AppDTOsPath: path.Join(config.BasePackage, lib.AppDTOsDir),
		Controllers: controllers,
	}
	goRoutesTemplate.Execute(&buf, routesTplValues)
	lib.EnsureDir(filepath.Dir(lib.RoutesBootstrapFile))
	ioutil.WriteFile(lib.RoutesBootstrapFile, buf.Bytes(), 0777)
	// TODO Verbose mode
	// fmt.Printf("Generated routes bootstrap file to `%s` in %f seconds\n", lib.RoutesBootstrapFile, time.Since(start).Seconds())

	// start = time.Now()
	// DTOS
	dtos := genDTOSMap(lib.CoreDTOsDir)
	appDTOs := genDTOSMap(lib.AppDTOsDir)
	for dtoName := range appDTOs {
		dtos[dtoName] = appDTOs[dtoName]
	}

	// Aggregates
	aggregates := genAggregatesMap(lib.CoreAggregatesDir)
	appAggregates := genAggregatesMap(lib.AppAggregatesDir)

	for aggregateName := range appAggregates {
		aggregates[aggregateName] = appAggregates[aggregateName]
	}

	var permissions = map[string]string{}
	permissions, e = FetchAllPermissionsFromControllers(controllers)

	if e != nil {
		return nil, e
	}

	var routesContainer = &lib.RoutesJSONContainer{
		Routes:      map[string]*lib.ControllerRoute{},
		DTOs:        dtos,
		Models:      genModelsMap(schemaList),
		Aggregates:  aggregates,
		Constants:   genConstantsMap(),
		Permissions: permissions,
	}

	for k := range controllers {
		for i := range controllers[k].Routes {
			key := controllers[k].Routes[i].Name
			routesContainer.Routes[key] = controllers[k].Routes[i]
		}
	}

	routesJSON, _ := json.MarshalIndent(routesContainer, "  ", "    ")
	ioutil.WriteFile(lib.RoutesFilePath, routesJSON, 0777)
	// fmt.Printf("Generated routes to `%s` in %f seconds\n", lib.RoutesFilePath, time.Since(start).Seconds())

	// start = time.Now()
	// Controller Bootstrap
	var controllerBootstrapFile = GenControllerBootstrap(config.BasePackage, dirs)
	lib.EnsureDir(filepath.Dir(lib.ControllersBootstrapGenFile))
	ioutil.WriteFile(lib.ControllersBootstrapGenFile, []byte(controllerBootstrapFile), 0777)
	// TODO Verbose mode
	// fmt.Printf("Generated ControllerBootstrapGenFile to `%s` in %f seconds\n", lib.ControllersBootstrapGenFile, time.Since(start).Seconds())

	return routesContainer, nil
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
				s = append(s, fmt.Sprintf("\t\t// `%s` - %s", query.Pattern, query.Pattern[0:2]))

				if query.Type == "int64" {
					defaultValue := "0"
					if query.Pattern == "-?[0-9]+" {
						defaultValue = "-1"
					}
					s = append(s, fmt.Sprintf("\t\t%s := req.ArgInt64(\"%s\", %s)\n", query.VariableName, query.VariableName, defaultValue))
				} else {
					if query.Pattern[0:2] == "-?" {
						s = append(s, fmt.Sprintf("\t\t%s := req.RootRequest.URL.Query().Get(\"%s\")", query.VariableName, query.VariableName))
					} else {
						s = append(s, fmt.Sprintf("\t\t%s := req.Arg(\"%s\", \"\")\n", query.VariableName, query.VariableName))
					}
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
			var maybeQueries = []string{}
			for _, query := range route.Queries {
				if query.Pattern[0:2] == "-?" {
					continue
				}
				maybeQueries = append(maybeQueries, fmt.Sprintf("\t\t\t\"%s\", \"%s\",", query.VariableName, query.ValueRaw))
			}
			if len(maybeQueries) > 0 {
				s = append(s, "\t\tQueries(")
				s = append(s, maybeQueries...)
				// for _, query := range route.Queries {
				// 	s = append(s, fmt.Sprintf("\t\t\t\"%s\", \"%s\",", query.Name, query.ValueRaw))
				// }
				s = append(s, "\t\t).")
			}
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

func genModelsMap(schemas *schema.SchemaList) map[string]map[string]string {

	var tableMap = map[string]map[string]string{}
	for k := range schemas.Schemas {

		var s = schemas.Schemas[k]

		for l := range s.Tables {
			var table = s.Tables[l]

			tableMap[table.Name] = map[string]string{}

			for m := range table.Columns {

				tableMap[table.Name][m] = schema.DataTypeToGoTypeString(table.Columns[m])

			}

		}

	}

	return tableMap

	// dirHandle, err := os.Open(lib.ModelsGenDir)

	// if err != nil {
	// 	panic(err)
	// }

	// defer dirHandle.Close()

	// var dirFileNames []string
	// dirFileNames, err = dirHandle.Readdirnames(-1)

	// if err != nil {
	// 	panic(err)
	// }
	// // reader := bufio.NewReader(os.Stdin)

	// result := map[string]map[string]string{}

	// for _, name := range dirFileNames {

	// 	if name == ".DS_Store" {
	// 		continue
	// 	}

	// 	modelName := name[0 : len(name)-3]
	// 	fullPath := path.Join(lib.ModelsGenDir, name)

	// 	model, e := ParseFileToGoStruct(fullPath)
	// 	if e != nil {
	// 		panic(e)
	// 	}
	// 	k := 0

	// 	result[modelName] = map[string]string{}

	// 	for k < model.Fields.Len() {
	// 		result[modelName][model.Fields.Get(k).Name] = model.Fields.Get(k).DataType
	// 		k++
	// 	}
	// }

	// return result
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

					var isSlice = len(fieldType) > 2 && fieldType[0:2] == "[]"
					var containsDotSeparator = strings.Contains(fieldType, ".")
					var isPointer = false
					if isSlice {
						fieldType = fieldType[2:]
						isPointer = len(fieldType) > 1 && fieldType[0:1] == "*"
						if isPointer && !containsDotSeparator {
							fieldType = "[]*aggregates." + fieldType[1:]
						} else {
							fieldType = "[]" + fieldType
						}
					} else {
						if !containsDotSeparator && isPointer {
							fieldType = "*aggregates." + fieldType[1:]
						}
					}

				} else {

					var containsDotSeparator = strings.Contains(parts[0], ".")
					var isPointer = len(parts[0]) > 1 && parts[0][0:1] == "*"
					// This is an embedded type (DTO or model)
					if containsDotSeparator || isPointer {
						// sParts := strings.Split(parts[0], ".")
						// fieldName = sParts[1]
						fieldType = parts[0]
						fieldName = "#embedded" + fmt.Sprint(k)

						// Has a pointer but no dot separator, so it's an embedded struct in the same `aggregates` package
						if !containsDotSeparator && isPointer {
							fieldType = "*aggregates." + fieldType[1:]
						}
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
