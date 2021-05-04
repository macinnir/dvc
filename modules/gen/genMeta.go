package gen

import (
	"html/template"
	"os"
	"path"

	"github.com/macinnir/dvc/lib"
)

// GenMeta returns a string for a model in golang
func (g *Gen) GenMeta(dir string, database *lib.Database) (e error) {

	lib.EnsureDir(dir)

	p := path.Join(dir, "models.go")

	lib.Infof("Generating `%s`", g.Options, p)

	tpl := `// Generated Code; DO NOT EDIT.

package meta 

// Tables returns metaData about Tables 
func Tables() map[string]map[string][]string {
	return map[string]map[string][]string {
	{{range .Tables}}
		"{{.Name}}": { 
			{{range .Columns}}
			"{{.Name}}": {"{{.FmtType}}", "{{.DataType}}", "{{.GoType}}"}, 
			{{end}}
		}, 
	{{end}}
	}
}
	
	`

	t := template.Must(template.New("meta-tables").Parse(tpl))
	f, e := os.Create(p)
	if e != nil {
		panic(e)
	}

	e = t.Execute(f, database)
	if e != nil {
		panic(e)
	}

	f.Close()
	lib.FmtGoCode(p)
	return
}
