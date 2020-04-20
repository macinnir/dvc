package gen

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"

	"github.com/macinnir/dvc/lib"
)

// GenerateGoModel returns a string for a model in golang
func (g *Gen) GenerateGoModel(dir string, table *lib.Table) (e error) {

	g.EnsureDir(dir)

	p := path.Join(dir, table.Name+".go")

	if g.fileExists(p) {
		lib.Infof("Updating `%s`", g.Options, table.Name)
		e = g.updateGoModel(p, table)
		return
	}

	e = g.buildGoModel(p, table)
	return
}

func (g *Gen) updateGoModel(p string, table *lib.Table) (e error) {

	var modelNode *lib.GoStruct
	var outFile []byte

	fileBytes, e := ioutil.ReadFile(p)
	modelNode, e = buildModelNodeFromFile(fileBytes)
	resolveTableToModel(modelNode, table)

	if e != nil {
		return
	}

	outFile, e = buildFileFromModelNode(modelNode)
	if e != nil {
		return
	}

	e = ioutil.WriteFile(p, outFile, 0644)

	return
}

func (g *Gen) buildGoModel(p string, table *lib.Table) (e error) {
	var modelNode *lib.GoStruct
	var outFile []byte

	lib.Debugf("Generating model for table %s", g.Options, table.Name)

	modelNode, e = buildModelNodeFromTable(table)
	if e != nil {
		return
	}

	outFile, e = buildFileFromModelNode(modelNode)
	if e != nil {
		return
	}

	e = ioutil.WriteFile(p, outFile, 0644)
	return
}

func (g *Gen) buildGoModelOld(p string, table *lib.Table) (e error) {

	var fileHead, fileFoot string
	oneToMany := g.Config.OneToMany[table.Name]
	oneToOne := g.Config.OneToOne[table.Name]

	lib.Debugf("Generating model for table %s", g.Options, table.Name)
	if fileHead, fileFoot, _, e = g.scanFileParts(p, false); e != nil {
		return
	}

	type Column struct {
		Name         string
		Type         string
		IsPrimaryKey bool
		IsNull       bool
		Ordinal      int
	}

	data := struct {
		OneToMany          string
		OneToOne           string
		Name               string
		IncludeNullPackage bool
		Columns            []Column
		FileHead           string
		FileFoot           string
	}{
		OneToMany:          oneToMany,
		OneToOne:           oneToOne,
		Name:               table.Name,
		IncludeNullPackage: false,
		Columns:            []Column{},
		FileHead:           fileHead,
		FileFoot:           fileFoot,
	}

	tpl := `
package models
{{if .IncludeNullPackage}}
import (
	"gopkg.in/guregu/null.v3"
)
{{end}}

// {{.Name}} represents a {{.Name}} domain object 
type {{.Name}} struct { 
	{{range .Columns}}
{{.Name}} {{.Type}} ` + "`db:\"{{.Name}}\" json:\"{{.Name}}\"`" + `{{end}}
}
`
	var sortedColumns = make(lib.SortedColumns, 0, len(table.Columns))

	for _, column := range table.Columns {
		sortedColumns = append(sortedColumns, column)
	}

	sort.Sort(sortedColumns)

	includeNullPackage := false

	for _, column := range sortedColumns {

		fieldType := dataTypeToGoTypeString(column)

		if column.IsNullable {
			includeNullPackage = true
		}

		data.Columns = append(data.Columns, Column{
			Name:         column.Name,
			Type:         fieldType,
			IsPrimaryKey: column.ColumnKey == "PRI",
			IsNull:       column.IsNullable,
			Ordinal:      len(data.Columns),
		})
	}
	data.IncludeNullPackage = includeNullPackage

	t := template.Must(template.New("model-" + table.Name).Parse(tpl))
	f, err := os.Create(p)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}

	err = t.Execute(f, data)
	if err != nil {
		fmt.Println("Execute ERROR: ", err.Error())
		return
	}

	f.Close()
	g.FmtGoCode(p)

	return
}

// CleanGoModels removes model files that are not found in the database.Tables map
func (g *Gen) CleanGoModels(dir string, database *lib.Database) (e error) {

	g.EnsureDir(dir)

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

	for _, name := range dirFileNames {
		fileNameNoExt := name[0 : len(name)-3]
		if _, ok := database.Tables[fileNameNoExt]; !ok {
			fullFilePath := path.Join(dir, name)
			log.Printf("Removing %s\n", fullFilePath)
			os.Remove(fullFilePath)
		}
	}
	return
}