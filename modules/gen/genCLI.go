package gen

import (
	"fmt"
	"html/template"
	"os"
	"path"
)

// GenerateGoCLI generates a boilerplate cli app
func (g *Gen) GenerateGoCLI(dir string) (e error) {

	tpl := `
package main

import (
	"github.com/macinnir/dvc/modules/utils"
	base "{{ .BasePackage }}"
	"{{ .BasePackage }}/models"
	"gopkg.in/guregu/null.v3"
)

func main() {
	app := base.NewApp("{{ .BasePackage }}_cli")
	app.InitConfig()
	app.InitLogging()
	app.InitRepos()
	app.InitStore()
	app.InitServices()

}
`
	p := path.Join(dir, "cli/main.go")
	t := template.Must(template.New("cli").Parse(tpl))
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

	g.FmtGoCode(p)

	return
}
