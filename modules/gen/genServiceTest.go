package gen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"

	"github.com/macinnir/dvc/lib"
	"golang.org/x/tools/imports"
)

// FormatCode formats the code
func FormatCode(code string) ([]byte, error) {
	opts := &imports.Options{
		TabIndent: true,
		TabWidth:  2,
		Fragment:  true,
		Comments:  true,
	}
	return imports.Process("", []byte(code), opts)
}

// GenServiceTest generates a test for a service file
func (g *Gen) GenServiceTest(serviceName, filePath string) (e error) {

	fmt.Println("ServiceName", serviceName, "FilePath", filePath)

	structName := fmt.Sprintf(serviceName + "Service")

	var src []byte

	src, e = ioutil.ReadFile(filePath)

	if e != nil {
		return
	}

	fset := token.NewFileSet()

	var file *ast.File

	if file, e = parser.ParseFile(fset, "", src, parser.ParseComments); e != nil {
		return
	}

	for _, d := range file.Decls {
		if a, decl := lib.GetReceiverTypeName(src, d); a == structName {
			if !decl.Name.IsExported() {
				continue
			}

			for _, l := range decl.Type.Params.List {
				for _, n := range l.Names {
					fmt.Println(n.Name)
				}
			}
		}
	}

	return

	imports := []string{
		"testing",
		"github.com/golang/mock/gomock",
		"github.com/stretchr/testify/assert",
		"github.com/stretchr/testify/require",
		fmt.Sprintf("%s/%s/models", g.Config.BasePackage, g.Config.Dirs.Definitions),
		fmt.Sprintf("%s/%s/integrations", g.Config.BasePackage, g.Config.Dirs.Definitions),
	}

	data := struct {
		Name    string
		Imports []string
	}{
		Name:    serviceName,
		Imports: imports,
	}

	tplString := `
package services 

import (
{{range $import := .Imports}}	"{{$import}}"
{{end}})

type {{.Name}}ServiceTestObj struct {
	service *{{.Name}}Service
}

func Setup{{.Name}}ServiceTests(t *testing.T) (testObj *{{.Name}}ServiceTestObj) {
	ctrl := gomock.NewController(t) 

	testObj = &{{.Name}}ServiceTestObj {
		
	}

	testObj.service = New{{.Name}}Service(

	)

	return 
}
	
	
`

	t := template.Must(template.New("service-test-" + serviceName).Parse(tplString))
	var out bytes.Buffer
	e = t.Execute(&out, data)
	if e != nil {
		return
	}

	fmt.Println(out.String())

	return
}
