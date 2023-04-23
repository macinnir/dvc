package gen

import (
	"bytes"
	"io/ioutil"
	"path"
	"text/template"

	"github.com/macinnir/dvc/core/lib"
)

var collectionItemFileTpl = template.Must(template.New("repo-collection-item").Parse(`// Generated Code; DO NOT EDIT.
package collections

import (
	"{{ .BasePackage }}/gen/definitions/models" 
)

type {{.TableName}}CollectionItem struct { 
	*models.{{ .TableName }}
	ID string 
}
`))

func GenerateRepoCollectionItem(basePackage string, tableName string) {
	var p = path.Join(lib.CollectionGenDir, tableName+"CollectionItem.go")
	var buf bytes.Buffer

	var data = struct {
		TableName   string
		BasePackage string
	}{
		TableName:   tableName,
		BasePackage: basePackage,
	}

	collectionItemFileTpl.Execute(&buf, data)
	ioutil.WriteFile(p, buf.Bytes(), lib.DefaultFileMode)
}
