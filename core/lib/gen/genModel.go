package gen

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

// GenerateGoModel returns a string for a model in golang
func (g *Gen) GenerateGoModel(dir string, table *schema.Table) (e error) {

	lib.EnsureDir(dir)

	p := path.Join(dir, table.Name+".go")

	// if lib.FileExists(p) {
	// 	lib.Infof("Updating `%s`", g.Options, table.Name)
	// 	e = g.updateGoModel(p, table)
	// 	return
	// }

	// lib.Infof("Creating `%s`", g.Options, table.Name)s
	e = g.buildGoModel(p, table)
	return
}

// InspectFile inspects a file
func InspectFile(filePath string) (s *lib.GoStruct, e error) {

	fileBytes, e := ioutil.ReadFile(filePath)
	if e != nil {
		panic(e)
	}

	s, e = buildGoStructFromFile(fileBytes)
	if e != nil {
		fmt.Println("ERROR: ", filePath)
		panic(e)
	}

	return

}

// func (g *Gen) updateGoModel(p string, table *schema.Table) (e error) {

// 	var modelNode *lib.GoStruct
// 	var outFile []byte

// 	fileBytes, e := ioutil.ReadFile(p)
// 	modelNode, e = buildGoStructFromFile(fileBytes)
// 	resolveTableToModel(modelNode, table)

// 	if e != nil {
// 		return
// 	}

// 	outFile, e = buildFileFromModelNode(modelNode)
// 	if e != nil {
// 		return
// 	}

// 	e = ioutil.WriteFile(p, outFile, 0644)

// 	return
// }

func (g *Gen) buildGoModel(p string, table *schema.Table) (e error) {
	var modelNode *lib.GoStruct
	var outFile []byte

	// lib.Debugf("Generating model for table %s", g.Options, table.Name)

	modelNode, e = buildModelNodeFromTable(table)
	if e != nil {
		return
	}

	outFile, e = buildFileFromModelNode(table, modelNode)
	if e != nil {
		return
	}

	e = ioutil.WriteFile(p, outFile, 0644)
	return
}

// func (g *Gen) buildGoModelOld(p string, table *schema.Table) (e error) {

// 	oneToMany := g.Config.OneToMany[table.Name]
// 	oneToOne := g.Config.OneToOne[table.Name]

// 	lib.Debugf("Generating model for table %s", g.Options, table.Name)

// 	type Column struct {
// 		Name         string
// 		Type         string
// 		IsPrimaryKey bool
// 		IsNull       bool
// 		Ordinal      int
// 	}

// 	data := struct {
// 		OneToMany          string
// 		OneToOne           string
// 		Name               string
// 		IncludeNullPackage bool
// 		Columns            []Column
// 	}{
// 		OneToMany:          oneToMany,
// 		OneToOne:           oneToOne,
// 		Name:               table.Name,
// 		IncludeNullPackage: false,
// 		Columns:            []Column{},
// 	}

// 	tpl := `// Generated Code; DO NOT EDIT.

// package models
// {{if .IncludeNullPackage}}
// import (
// 	"gopkg.in/guregu/null.v3"
// )
// {{end}}

// // {{.Name}} represents a {{.Name}} domain object
// type {{.Name}} struct {
// 	{{range .Columns}}
// {{.Name}} {{.Type}} ` + "`db:\"{{.Name}}\" json:\"{{.Name}}\"`" + `{{end}}
// }`
// 	var sortedColumns = make(lib.SortedColumns, 0, len(table.Columns))

// 	for _, column := range table.Columns {
// 		sortedColumns = append(sortedColumns, column)
// 	}

// 	sort.Sort(sortedColumns)

// 	includeNullPackage := false

// 	for _, column := range sortedColumns {

// 		fieldType := schema.DataTypeToGoTypeString(column)

// 		if column.IsNullable {
// 			includeNullPackage = true
// 		}

// 		data.Columns = append(data.Columns, Column{
// 			Name:         column.Name,
// 			Type:         fieldType,
// 			IsPrimaryKey: column.ColumnKey == "PRI",
// 			IsNull:       column.IsNullable,
// 			Ordinal:      len(data.Columns),
// 		})
// 	}
// 	data.IncludeNullPackage = includeNullPackage

// 	t := template.Must(template.New("model-" + table.Name).Parse(tpl))
// 	f, err := os.Create(p)
// 	if err != nil {
// 		fmt.Println("ERROR: ", err.Error())
// 		return
// 	}

// 	err = t.Execute(f, data)
// 	if err != nil {
// 		fmt.Println("Execute ERROR: ", err.Error())
// 		return
// 	}

// 	f.Close()
// 	lib.FmtGoCode(p)

// 	return
// }

// CleanGoModels removes model files that are not found in the database.Tables map
func (g *Gen) CleanGoModels(dir string, database *schema.Schema) (e error) {

	lib.EnsureDir(dir)

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
	reader := bufio.NewReader(os.Stdin)

	for _, name := range dirFileNames {
		fileNameNoExt := name[0 : len(name)-3]
		if _, ok := database.Tables[fileNameNoExt]; !ok {
			fullFilePath := path.Join(dir, name)
			// log.Printf("Removing %s\n", fullFilePath)
			result := lib.ReadCliInput(reader, fmt.Sprintf("Delete unused model `%s` (Y/n)?", fileNameNoExt))
			if result == "Y" {
				fmt.Printf("Deleting model `%s` at path `%s`...\n", fileNameNoExt, fullFilePath)
				os.Remove(fullFilePath)
			}
		}
	}
	return
}
