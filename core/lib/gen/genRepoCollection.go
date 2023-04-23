package gen

import (
	"bytes"
	"io/ioutil"
	"path"
	"text/template"

	"github.com/macinnir/dvc/core/lib"
)

var collectionFileTpl = template.Must(template.New("repo-collection").Parse(`// Generated Code; DO NOT EDIT.
package collections

import (
	"{{ .BasePackage }}/gen/definitions/models" 
)

type {{.TableName}}Collection struct { 
	Count int64 
	Data []*models.{{ .TableName }}
}
`))

func GenerateRepoCollection(basePackage string, tableName string) {
	var p = path.Join(lib.CollectionGenDir, tableName+"Collection.go")
	var buf bytes.Buffer

	var data = struct {
		TableName   string
		BasePackage string
	}{
		TableName:   tableName,
		BasePackage: basePackage,
	}

	collectionFileTpl.Execute(&buf, data)
	ioutil.WriteFile(p, buf.Bytes(), lib.DefaultFileMode)
}
