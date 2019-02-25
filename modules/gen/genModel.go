package gen

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"sort"

	"github.com/macinnir/dvc/lib"
)

// GenerateGoModel returns a string for a model in golang
func (g *Gen) GenerateGoModel(dir string, table *lib.Table) (e error) {

	g.EnsureDir(dir)

	var fileHead, fileFoot string
	oneToMany := g.Config.OneToMany[table.Name]
	oneToOne := g.Config.OneToOne[table.Name]

	p := path.Join(dir, table.Name+".go")
	lib.Debugf("Generating model for table %s at path %s", g.Options, table.Name, p)
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
{{.FileHead}} 
// #genStart

package models
import (
	q "github.com/macinnir/dvc/modules/query" 
	{{if .IncludeNullPackage}}"gopkg.in/guregu/null.v3"{{end}}
)

// {{.Name}} represents a {{.Name}} domain object 
type {{.Name}} struct { 
	q.BaseDomainObject
	{{range .Columns}}
{{.Name}} {{.Type}} ` + "`db:\"{{.Name}}\" json:\"{{.Name}}\"`" + `{{end}}
{{if ne .OneToMany ""}}
	{{.OneToMany}}s []*{{.OneToMany}}{{end}}
{{if ne .OneToOne ""}}
	{{.OneToOne}} *{{.OneToOne}}{{end}}
}

// New{{.Name}} returns a new {{.Name}} domain object
func New{{.Name}}() *{{.Name}} {
	o := &{{.Name}}{}
	o.Build()
	return o
}

// Build builds the internal meta data for a {{.Name}} domain object
func (o *{{.Name}}) Build() {
	o.BaseDomainObject.TableName = "{{.Name}}"
	o.BaseDomainObject.FieldList = map[string]*q.DomainObjectField{
		{{range.Columns}}
		"{{.Name}}": {
			Name: "{{.Name}}", 
			Type: "{{.Type}}", 
			IsPrimaryKey: {{.IsPrimaryKey}}, 
			IsNull: {{.IsNull}}, 
			Ordinal: {{.Ordinal}}, 
		}, 
		{{end}}
	}
}
// #genEnd
{{.FileFoot}}
`
	var sortedColumns = make(lib.SortedColumns, 0, len(table.Columns))

	for _, column := range table.Columns {
		sortedColumns = append(sortedColumns, column)
	}

	sort.Sort(sortedColumns)

	includeNullPackage := false

	for _, column := range sortedColumns {

		fieldType := "int64"
		switch column.DataType {
		case "varchar":
			fieldType = "string"
		case "enum":
			fieldType = "string"
		case "text":
			fieldType = "string"
		case "date":
			fieldType = "string"
		case "datetime":
			fieldType = "string"
		case "char":
			fieldType = "string"
		case "decimal":
			fieldType = "float64"
		}

		if column.IsNullable == true {
			includeNullPackage = true
			switch fieldType {
			case "string":
				// fieldType = "sql.NullString"
				fieldType = "null.String"
			case "int64":
				// fieldType = "sql.NullInt64"
				fieldType = "null.Int"
			case "float64":
				// fieldType = "sql.NullFloat64"
				fieldType = "null.Float"
			}
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

// GenerateDefaultConfigModelFile generates a default config model file
func (g *Gen) GenerateDefaultConfigModelFile(dir string) {

	model := &lib.Table{
		Columns: map[string]*lib.Column{},
		Name:    "Config",
	}

	model.Columns["AppName"] = &lib.Column{Name: "AppName", DataType: "varchar"}
	model.Columns["DBName"] = &lib.Column{Name: "DBName", DataType: "varchar"}
	model.Columns["DBHost"] = &lib.Column{Name: "DBHost", DataType: "varchar"}
	model.Columns["DBUser"] = &lib.Column{Name: "DBUser", DataType: "varchar"}
	model.Columns["DBPass"] = &lib.Column{Name: "DBPass", DataType: "varchar"}
	model.Columns["Domain"] = &lib.Column{Name: "Domain", DataType: "varchar"}
	model.Columns["Port"] = &lib.Column{Name: "Port", DataType: "varchar"}
	model.Columns["HTTPS"] = &lib.Column{Name: "HTTPS", DataType: "varchar"}
	model.Columns["URLVersionPrefix"] = &lib.Column{Name: "URLVersionPrefix", DataType: "varchar"}
	model.Columns["TokenExpiryMinute"] = &lib.Column{Name: "TokenExpiryMinute", DataType: "int"}
	model.Columns["TokenIssuerName"] = &lib.Column{Name: "TokenIssuerName", DataType: "varchar"}
	model.Columns["PublicDomain"] = &lib.Column{Name: "PublicDomain", DataType: "varchar"}
	model.Columns["RedisHost"] = &lib.Column{Name: "RedisHost", DataType: "varchar"}
	model.Columns["RedisPassword"] = &lib.Column{Name: "RedisPassword", DataType: "varchar"}
	model.Columns["RedisDB"] = &lib.Column{Name: "RedisDB", DataType: "int"}
	model.Columns["WSURL"] = &lib.Column{Name: "WSURL", DataType: "varchar"}
	model.Columns["HTMLURL"] = &lib.Column{Name: "HTMLURL", DataType: "varchar"}

	g.GenerateGoModel(dir, model)
}

// GenerateDefaultConfigJsonFile generates a default json configuration file
// in the root of the project directory
func (g *Gen) GenerateDefaultConfigJsonFile(dir string) {

	data := map[string]string{
		"DatabaseName": g.Config.Connection.DatabaseName,
		"Host":         g.Config.Connection.Host,
		"Username":     g.Config.Connection.Username,
		"Password":     g.Config.Connection.Password,
		"AppName":      g.Config.BasePackage,
	}
	tpl := `{
	"AppName": "{{ .AppName }}",
	"DBName": "{{ .DatabaseName }}", 
	"DBHost": "{{ .Host }}", 
	"DBUser": "{{ .Username }}", 
	"DBPass": "{{ .Password }}", 
	"URLVersionPrefix": "v1",
	"TokenExpiryMinute": 1440,
	"TokenIssuerName": "{{ .AppName }}",
	"Domain": "0.0.0.0",
	"Port": "8081",
	"HTTPS": "http",
	"PublicDomain": "localhost",
	"RedisHost": "127.0.0.1:6379",
	"RedisPassword": "",
	"RedisDB": 0,
	"WSURL": ":8082",
	"HTMLURL": ":8080"
}
`
	t := template.Must(template.New("config-json").Parse(tpl))
	p := path.Join(dir, "config.json")
	f, err := os.Create(p)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return
	}

	err = t.Execute(f, data)
	if err != nil {
		fmt.Println("Execute Error: ", err.Error())
		return
	}

	f.Close()
}

// CleanGoModels removes model files that are not found in the database.Tables map
func (g *Gen) CleanGoModels(dir string, database *lib.Database) (e error) {
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
			if fileNameNoExt != "Config" {
				fullFilePath := path.Join(dir, name)
				fmt.Printf("TEST: Removing %s\n", fullFilePath)
				os.Remove(fullFilePath)
			}
		}
	}
	return
}
