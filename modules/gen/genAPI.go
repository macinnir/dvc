package gen

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

type apiPart struct {
	baseRoute string
	route     string
	verb      string
	args      map[string]string
	method    string
}

type apiPartGroup struct {
	baseRoute string
	routes    map[string]apiPart
	routeKeys sort.StringSlice
}

func (g *Gen) GenerateAPIRoutes(apiDir string) {

	// Find the directory where the api exists
	if _, e := os.Stat(apiDir); os.IsNotExist(e) {
		panic(fmt.Errorf("Can't generate API routes: No API directory exists at path %s", apiDir))
	}

	apiFiles, _ := g.getServiceNames(apiDir)
	routes := map[string]apiPartGroup{}
	groupKeys := []string{}
	currentRoutes := []string{}
	// for _, apiName := range apiFiles {
	// 	objName := strings.ToLower(apiName[0:1]) + apiName[1:]
	// 	file += fmt.Sprintf("\t%s := &%s{app}\n", objName, apiName)
	// }

	validRouteSig := regexp.MustCompile(`^// @route.*$`)

	for _, apiName := range apiFiles {

		// objName := strings.ToLower(apiName[0:1]) + apiName[1:]
		apiFilePath := path.Join(apiDir, apiName+".go")

		fileBytes, _ := ioutil.ReadFile(apiFilePath)
		fileString := string(fileBytes)
		fileLines := strings.Split(fileString, "\n")

		for _, line := range fileLines {

			// This is a line that starts with `// @route`
			if validRouteSig.Match([]byte(line)) {
				currentRoutes = append(currentRoutes, line)
				continue
			}

			// Allow for other comments
			if len(line) < 7 || line[0:5] != "func " {
				continue
			}

			// This is the line below @route line(s)
			if len(currentRoutes) > 0 {
				for _, currentRoute := range currentRoutes {

					args := map[string]string{}
					lineParts := strings.Split(line, " ")
					currentRoute = currentRoute[10:]
					currentRouteParts := strings.Split(currentRoute, " ")
					baseRoute := currentRouteParts[1]
					route := baseRoute

					queryPart := ""

					if strings.Contains(route, "?") {
						routeParts := strings.Split(route, "?")
						queryPart = routeParts[1]
						baseRoute = routeParts[0]

						keyValues := []string{routeParts[1]}
						if strings.Contains(routeParts[1], "&") {
							keyValues = strings.Split(routeParts[1], "&")
						}

						for _, keyValue := range keyValues {
							if strings.Contains(keyValue, "=") {
								kvParts := strings.Split(keyValue, "=")
								args[kvParts[0]] = kvParts[1]
							}
						}

					}

					verb := currentRouteParts[0]

					subRouteKey := fmt.Sprintf("%s_%d_%s", verb, len(queryPart), queryPart)
					// if len(queryPart) > 0 {
					// 	subRouteKey += "?" + queryPart + " "
					// }
					// subRouteKey += verb

					p := apiPart{
						verb:      verb,
						baseRoute: baseRoute,
						route:     route,
						method:    fmt.Sprintf("routes.%s", lineParts[3][:len(lineParts[3])-2]),
						args:      args,
					}

					// Keyed by the base route path
					if _, ok := routes[baseRoute]; !ok {
						routes[baseRoute] = apiPartGroup{
							baseRoute: baseRoute,
							routes:    map[string]apiPart{},
							routeKeys: sort.StringSlice{},
						}
						groupKeys = append(groupKeys, baseRoute)
					}

					routes[baseRoute].routes[subRouteKey] = p
					routeBaseRoute := routes[baseRoute]
					routeBaseRoute.routeKeys = append(routeBaseRoute.routeKeys, subRouteKey)
					routes[baseRoute] = routeBaseRoute

					currentRoutes = []string{}
					continue
				}
			}
		}
	}

	file := `package main

	func (routes *Routes) init() {`
	file += "\n\n"
	file += "\n\troutes.routes = []Route{\n"

	sort.Strings(groupKeys)

	for _, routeGroupKey := range groupKeys {

		log.Printf("RouteKey: %s\n", routeGroupKey)
		routeGroup := routes[routeGroupKey]

		sort.Sort(sort.Reverse(routeGroup.routeKeys))

		file += fmt.Sprintf("\t\t// Group: %s\n", routeGroupKey)

		for _, routeKey := range routeGroup.routeKeys {

			argsString := ""

			route := routeGroup.routes[routeKey]

			if len(route.args) > 0 {

				argsParts := []string{}

				// Reassemble into URL query string
				// argsString += "?"
				// for k, v := range route.args {
				// 	argsParts = append(argsParts, fmt.Sprintf("%s=%s", k, v))
				// }
				// argsString += strings.Join(argsParts, "&")

				for k, v := range route.args {
					argsParts = append(argsParts, fmt.Sprintf("\"%s\": \"%s\"", k, v))
				}

				argsString = strings.Join(argsParts, ", ")
			}

			file += fmt.Sprintf("\t\t{ \"%s\", \"%s\", map[string]string{%s}, %s }, \n", route.baseRoute, route.verb, argsString, route.method)
		}
	}

	file += "\t}\n\n\treturn \n}"

	ioutil.WriteFile(path.Join(g.Config.Dirs.API, "routes.go"), []byte(file), 0644)

	// List all of the api files
}

func (g *Gen) GenerateGoAPI(dir string) {

	var tpl = `
package main

import (
    "compress/gzip"
    baseApp "{{ .BasePackage }}"
	"{{ .BasePackage }}/repos"
	"{{ .BasePackage }}/services"
	"github.com/macinnir/dvc/modules/utils"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var err error

// This is the main function that is executed at startup
func main() {

	app := baseApp.NewApp("{{ .BasePackage }}_api")
	app.InitRepos()
	app.InitStore()
	// app.InitCache()
	app.InitServices()

	// Router
	r := mux.NewRouter()
	r.HandleFunc("/", utils.NotImplementedHandler)

	// Routes
	Bootstrap(r.PathPrefix("/"+utils.Config.URLVersionPrefix).Subrouter(), app.Repos, app.Services, app.Store)

	httpProtocolString := "http"

	url := utils.Config.Domain + ":" + utils.Config.Port

	err = http.ListenAndServe(
		url,
		handlers.LoggingHandler(
			os.Stdout,
			utils.CORSHandler(
				routes.AuthHandler(
					app.Repos.User,
					handlers.CompressHandlerLevel(
						r,
						gzip.BestSpeed,
					),
				),
			),
		),
	)

	if err != nil {
		app.Finish()
		log.Fatal("ListenAndServe: ", err)
	}

	app.Finish()
}
// Bootstrap bootstraps all of the routes
func Bootstrap(r *mux.Router, re *repos.Repos, se *services.Services, store utils.IStore) {

}
`
	var e error

	if _, e = os.Stat(path.Join(dir, "api")); os.IsNotExist(e) {
		e = os.Mkdir(path.Join(dir, "api"), 0777)
		if e != nil {
			fmt.Println("ERROR: ", e.Error())
		}
	}

	p := path.Join(dir, "api", "main.go")
	fmt.Println("Generating API To path:", p)
	t := template.Must(template.New("app").Parse(tpl))
	f, err := os.Create(p)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}
	err = t.Execute(f, g.Config)
	if err != nil {
		fmt.Println("Execute Error: ", err.Error())
		return
	}

	f.Close()

	g.FmtGoCode(p)

}
