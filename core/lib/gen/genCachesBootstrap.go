package gen

import (
	"bytes"
	"io/ioutil"
	"path"
	"sort"
	"text/template"

	"github.com/macinnir/dvc/core/lib"
)

func GenerateCacheBootstrapFile(basePackage string, cache map[string]*lib.CacheConfig) error {
	var e error

	var p = path.Join(lib.CacheGenDir, "bootstrap.go")
	var buf bytes.Buffer

	var caches = make([]string, len(cache))
	var n = 0
	for k := range cache {
		caches[n] = k
		n++
	}

	sort.Strings(caches)

	var data = struct {
		BasePackage string
		Caches      []string
	}{
		BasePackage: basePackage,
		Caches:      caches,
	}

	if e = cacheBootstrapFileTpl.Execute(&buf, data); e != nil {
		panic(e)
	}

	e = ioutil.WriteFile(p, buf.Bytes(), lib.DefaultFileMode)

	return e

}

var cacheBootstrapFileTpl = template.Must(template.New("cache-bootstrap").Parse(`// Generated Code; DO NOT EDIT.
package caches

import (
	"{{ .BasePackage }}/core/components/redis" 
)

// Cache is a container for all cache providers
type Caches struct {
	{{range $name := .Caches}}
	{{$name}} *{{$name}}Cache{{end}}
}

// BootstrapCache bootstraps the cache
func BootstrapCaches(cache redis.IRedis) *Caches {

	return &Caches{ {{range $index := .Caches}}
		{{$index}}: New{{$index}}Cache(cache),{{end}}
	}
}`))
