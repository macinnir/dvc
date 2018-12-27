package gen

import (
	"fmt"
	"html/template"
	"os"
	"path"
)

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
