package routes

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

func GenTSRoutes(controllers []*lib.Controller, config *lib.Config) error {

	lib.EnsureDir(config.TypescriptRoutesPath)

	// Clean out any old files
	files, e := ioutil.ReadDir(config.TypescriptRoutesPath)

	if e != nil {
		return e
	}

	for k := range files {
		os.Remove(path.Join(config.TypescriptRoutesPath, files[k].Name()))
	}

	for k := range controllers {

		g := NewTSRouteGenerator(controllers[k])

		routes, e := g.genTSRoutesFromController(controllers[k])
		if e != nil {
			return e
		}

		filePath := path.Join(config.TypescriptRoutesPath, controllers[k].Name+".ts")
		e = ioutil.WriteFile(
			filePath,
			[]byte(routes),
			0777,
		)

		if e != nil {
			log.Fatalf("Error writing file %s: %s", filePath, e.Error())
		}
	}

	return nil
}

type TSRouteGenerator struct {
	imports    map[string]struct{}
	controller *lib.Controller
	rootRoute  *lib.ControllerRoute
	itemRoute  *lib.ControllerRoute
}

func NewTSRouteGenerator(c *lib.Controller) *TSRouteGenerator {
	return &TSRouteGenerator{
		imports:    map[string]struct{}{},
		controller: c,
	}
}

// Generate typescript code for a set of routes in a controller
func (t *TSRouteGenerator) genTSRoutesFromController(controller *lib.Controller) (string, error) {

	// Identify root route (if any)
	for k := range controller.Routes {
		if strings.Count(controller.Routes[k].Path, "/") == 1 && controller.Routes[k].Method == "GET" {
			t.rootRoute = controller.Routes[k]
		}

		if strings.Count(controller.Routes[k].Path, "/") == 2 && controller.Routes[k].Method == "GET" && len(controller.Routes[k].Params) > 0 {
			t.itemRoute = controller.Routes[k]
		}
	}

	// fmt.Println("Generating TSRoute from ", controller.Name)
	var s strings.Builder

	s.WriteString(`/**
 * Generated Code; DO NOT EDIT
 *
 * ` + strings.Title(controller.Package) + "." + controller.Name + `
 */
 `)

	s.WriteString(`
import axios from 'axios'; 
import { useQuery, useMutation, queryCache } from 'react-query';
`)

	var rest strings.Builder

	for _, route := range controller.Routes {

		rest.WriteString(`
// ` + route.Name + ` ` + route.Description + `
// ` + route.Method + ` ` + route.Raw + `
`)

		if len(route.Permission) > 0 {
			rest.WriteString(`// @permission ` + route.Permission + `
`)
		}
		rest.WriteString(t.genTSRoute(controller, route))

		if route.Method == "GET" {
			rest.WriteString(`
` + t.genUseTSRoute(route))
		}
		if route.Method == "PUT" || route.Method == "POST" || route.Method == "DELETE" {
			rest.WriteString(`
` + t.genUseMutationTSRoute(route))
		}
		rest.WriteString(`
`)

	}

	var imports = []string{}
	for k := range t.imports {
		if len(k) == 0 {
			continue
		}
		imports = append(imports, k)
	}

	sort.Strings(imports)

	for k := range imports {

		importTypeDir := "models"
		if len(imports[k]) > 3 && imports[k][len(imports[k])-3:] == "DTO" {
			importTypeDir = "dtos"
		}

		if len(imports[k]) > 9 && imports[k][len(imports[k])-9:] == "Aggregate" {
			importTypeDir = "aggregates"
		}

		fmt.Fprintf(&s, "import { %s } from 'gen/%s/%s';\n", imports[k], importTypeDir, imports[k])
	}

	return s.String() + rest.String(), nil

}

func (t *TSRouteGenerator) genTSRoute(controller *lib.Controller, route *lib.ControllerRoute) string {

	var str strings.Builder

	var tsRouteName = strings.ToLower(route.Name[0:1]) + route.Name[1:]

	hasBody := false

	str.WriteString(`export const ` + tsRouteName + ` = async (`)

	argIndex := 0

	args := []string{}

	if route.Name == "GetUserFromID" {
		fmt.Println("GetUserFromID: Params", len(route.Params))
	}

	// Start with Params
	if len(route.Params) > 0 {
		for k := range route.Params {
			argType := schema.GoTypeToTypescriptString(route.Params[k].Type)
			args = append(args, route.Params[k].Name+" : "+argType)
			t.AddImport(route.Params[k].Type)
			argIndex++
		}
	}

	// Next Queries
	if len(route.Queries) > 0 {
		for k := range route.Queries {
			argType := schema.GoTypeToTypescriptString(route.Queries[k].Type)
			args = append(args, route.Queries[k].Name+" : "+argType)
			t.AddImport(route.Queries[k].Type)
			argIndex++
		}
	}

	isFormData := false

	// Body is always the last argument
	if (route.Method == "POST" || route.Method == "PUT") && route.HasBody {
		if route.BodyFormat == "FormData" {
			hasBody = true
			isFormData = true
			args = append(args, `body : FormData`)
		}

		if len(route.BodyType) > 0 {
			hasBody = true
			bodyType := schema.GoTypeToTypescriptString(route.BodyType)
			args = append(args, `body : `+bodyType)
			t.AddImport(route.BodyType)
		}
	}

	if len(args) > 0 {
		if route.Method == "POST" || route.Method == "PUT" || route.Method == "DELETE" {
			str.WriteString(`args : { ` + strings.Join(args, ", ") + ` }`)
		} else {
			str.WriteString(strings.Join(args, ", "))
		}
	}

	str.WriteString(`) => await axios.` + strings.ToLower(route.Method))

	if route.Name == "ExportPolicies" {
		fmt.Println("ExportPolicies: ", route.ResponseType, "Format:", route.ResponseFormat)
	}

	var responseType = "any"
	var isBlob = false
	if len(route.ResponseType) > 0 {
		responseType = schema.GoTypeToTypescriptString(route.ResponseType)
	}

	if route.ResponseFormat == "BLOB" {
		responseType = ""
		isBlob = true
	}

	if len(responseType) > 0 {
		if tsRouteName == "getUserSession" {

			fmt.Println("getUserSession: ", route.ResponseType, responseType)
		}

		str.WriteString(`<` + responseType + `>`)

		t.AddImport(route.ResponseType)
	}

	str.WriteString("(")
	var routePath = route.Path

	// Arguments are passed inside an `args` object
	var argsPrefix = "args."
	if !(route.Method == "POST" || route.Method == "PUT" || route.Method == "DELETE") {
		argsPrefix = ""
	}
	// Replace params
	if len(route.Params) > 0 {
		for k := range route.Params {
			routePath = strings.Replace(routePath, "{"+route.Params[k].Name+":"+route.Params[k].Pattern+"}", "${"+argsPrefix+route.Params[k].Name+"}", 1)
		}
	}

	if len(route.Queries) > 0 {

		routePath += "?"

		for k := range route.Queries {

			routePath += route.Queries[k].Name + "=${" + argsPrefix + route.Queries[k].Name + "}"

			if k < len(route.Queries)-1 {
				routePath += "&"
			}
			// routePath = strings.Replace(routePath, "{"+route.Queries[k].Name+":"+route.Queries[k].Pattern+"}", "${"+route.Queries[k].Name+"}", 1)
		}
	}

	str.WriteString("`" + routePath + "`")

	if hasBody {
		str.WriteString(", " + argsPrefix + "body")
	} else {
		// No body, but should have one for put and post
		if route.Method == "POST" || route.Method == "PUT" {
			str.WriteString(", {}")
		}
	}

	if isBlob {
		str.WriteString(", { responseType: 'blob' }")
	}

	if isFormData {
		str.WriteString(`, { headers: { "content-type": "multipart/form-data" } }`)
	}

	str.WriteString(");")
	return str.String()
}

func (t *TSRouteGenerator) genUseTSRoute(route *lib.ControllerRoute) string {

	var str strings.Builder

	var tsRouteName = "use" + route.Name

	str.WriteString(`export const ` + tsRouteName + ` = (`)

	argIndex := 0

	argNames := []string{}
	argNamesWithTypes := []string{}

	// Start with Params
	if len(route.Params) > 0 {
		for k := range route.Params {
			if argIndex > 0 {
				str.WriteString(", ")
			}
			argType := schema.GoTypeToTypescriptString(route.Params[k].Type)
			str.WriteString(route.Params[k].Name + " : " + argType)

			t.AddImport(route.Params[k].Type)

			argIndex++
			argNames = append(argNames, route.Params[k].Name)
			argNamesWithTypes = append(argNamesWithTypes, route.Params[k].Name+" : "+argType)
		}
	}

	// Next Queries
	if len(route.Queries) > 0 {
		for k := range route.Queries {

			if argIndex > 0 {
				str.WriteString(", ")
			}

			argType := schema.GoTypeToTypescriptString(route.Queries[k].Type)
			str.WriteString(route.Queries[k].Name + " : " + argType)

			argIndex++
			argNames = append(argNames, route.Queries[k].Name)
			argNamesWithTypes = append(argNamesWithTypes, route.Queries[k].Name+" : "+argType)
		}
	}

	str.WriteString(`) => useQuery(["` + strings.ToLower(route.Name[0:1]) + route.Name[1:] + `"`)

	if len(argNames) > 0 {
		str.WriteString(`, ` + strings.Join(argNames, ", "))
	}
	str.WriteString(`], (_ : string`)
	if len(argNames) > 0 {
		str.WriteString(`, ` + strings.Join(argNamesWithTypes, ", "))
	}
	str.WriteString(`) => ` + strings.ToLower(route.Name[0:1]) + route.Name[1:] + `(`)
	if len(argNames) > 0 {
		str.WriteString(strings.Join(argNames, ","))
	}
	str.WriteString("));")

	return str.String()
}

func (t *TSRouteGenerator) genUseMutationTSRoute(route *lib.ControllerRoute) string {

	var str strings.Builder

	var tsRouteName = "use" + route.Name

	str.WriteString(`export const ` + tsRouteName + ` = () => useMutation(` + strings.ToLower(route.Name[0:1]) + route.Name[1:])

	if t.rootRoute != nil {
		str.WriteString(`, {
	onSuccess: (data, variables) => { 
		queryCache.invalidateQueries(["` + strings.ToLower(t.rootRoute.Name[0:1]) + t.rootRoute.Name[1:] + `"]);`)

		if t.itemRoute != nil && len(route.Params) > 0 {
			// if route.Name == "UpdateUserGroup" {
			// 	fmt.Println("got here!", t.itemRoute.Params[0].Name, " ==> ", route.Params[0].Name)
			// }
			containsParam := false
			for k := range route.Params {
				if route.Params[k].Name == t.itemRoute.Params[0].Name {
					containsParam = true
					break
				}
			}

			if containsParam {
				str.WriteString(`
		queryCache.invalidateQueries(["` + strings.ToLower(t.itemRoute.Name[0:1]) + t.itemRoute.Name[1:] + `", variables.` + t.itemRoute.Params[0].Name + `]);`)
			}
		}
		str.WriteString(`
	},
	throwOnError: true, 
});
	`)
	} else {
		str.WriteString(`);`)
	}

	return str.String()
}

func (t *TSRouteGenerator) AddImport(importType string) {

	if len(importType) == 0 {
		return
	}

	if importType[0:2] == "[]" {
		importType = importType[2:]
	}

	// Double slice?
	if len(importType) > 2 && importType[0:2] == "[]" {
		importType = importType[2:]
	}

	if len(importType) > 11 && importType[0:11] == "map[string]" {
		importType = importType[11:]
	}

	tsType := schema.GoTypeToTypescriptString(importType)

	if len(tsType) > 0 && tsType != "any" && !schema.IsGoTypeBaseType(importType) {
		if tsType[len(tsType)-2:] == "[]" {
			tsType = tsType[0 : len(tsType)-2]
		}
		t.imports[tsType] = struct{}{}
	}
}
