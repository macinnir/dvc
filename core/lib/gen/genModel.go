package gen

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/cache"
	"github.com/macinnir/dvc/core/lib/schema"
)

// GenModels generates models
func GenModels(modelsDir string, force bool, clean bool) error {

	start := time.Now()
	var schemaList *schema.SchemaList
	var e error

	schemaList, e = schema.LoadLocalSchemas()

	if e != nil {
		return e
	}

	fmt.Println("Generating models...")

	if clean {
		for k := range schemaList.Schemas {
			CleanGoModels(modelsDir, schemaList.Schemas[k])
		}
	}

	var tablesCache cache.TablesCache
	tablesCache, e = cache.LoadTableCache()

	if e != nil {
		return e
	}

	generatedModelCount := 0

	for k := range schemaList.Schemas {

		schemaName := schemaList.Schemas[k].Name

		for l := range schemaList.Schemas[k].Tables {

			table := schemaList.Schemas[k].Tables[l]
			tableKey := schemaName + "_" + table.Name

			var tableHash string
			tableHash, e = cache.HashTable(table)

			// If the model has been hashed before...
			if _, ok := tablesCache.Models[tableKey]; ok {

				// And the hash hasn't changed, skip...
				if tableHash == tablesCache.Models[tableKey] && !force {
					// fmt.Printf("Table `%s` hasn't changed! Skipping...\n", table.Name)
					continue
				}
			}

			generatedModelCount++

			// Update the models cache
			tablesCache.Models[tableKey] = tableHash

			fmt.Printf("Generating model `%s`\n", table.Name)
			e = GenerateGoModel(modelsDir, table)
			if e != nil {
				return e
			}
		}
	}

	cache.SaveTableCache(tablesCache)

	fmt.Printf("Generated %d models in %f seconds.\n", generatedModelCount, time.Since(start).Seconds())

	return nil
}

// GenerateGoModel returns a string for a model in golang
func GenerateGoModel(dir string, table *schema.Table) (e error) {

	lib.EnsureDir(dir)

	p := path.Join(dir, table.Name+".go")

	// if lib.FileExists(p) {
	// 	lib.Infof("Updating `%s`", g.Options, table.Name)
	// 	e = g.updateGoModel(p, table)
	// 	return
	// }

	// lib.Infof("Creating `%s`", g.Options, table.Name)s
	e = buildGoModel(p, table)
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

func buildGoModel(p string, table *schema.Table) (e error) {
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
func CleanGoModels(dir string, database *schema.Schema) (e error) {

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
