package gen

import (
	"bytes"
	"io/ioutil"
	"path"
	"sort"
	"text/template"

	"github.com/macinnir/dvc/core/lib"
)

type RepoBootstrapConfig struct {
	Name   string
	Config *lib.CacheConfig
}

type ScoredRepoItem struct {
	Name  string
	Score int
}

func GenRepoBootstrapFile(basePackage string, cache map[string]*lib.CacheConfig) error {
	var e error

	var p = path.Join(lib.RepoGenDir, "bootstrap.go")
	var buf bytes.Buffer

	var caches = make([]*RepoBootstrapConfig, len(cache))
	var n = 0

	var repoMap = map[string]struct{}{}
	var scores = make([]*ScoredRepoItem, len(cache))
	var scoreMap = make(map[string]int, len(cache))

	for tableName := range cache {

		scores[n] = &ScoredRepoItem{
			Score: 0,
			Name:  tableName,
		}

		scoreMap[tableName] = n
		n++
	}

	// sort.Slice(scores, func(i, j int) bool {
	// 	return scores[i].Name < scores[j].Name
	// })

	for k := range cache {

		if cache[k].Aggregate == nil || len(cache[k].Aggregate.Properties) == 0 {
			continue
		}

		for l := range cache[k].Aggregate.Properties {
			var agg = cache[k].Aggregate.Properties[l]
			if _, ok := repoMap[agg.Table]; !ok {

				// The dependent repo hasn't been defined yet, so it should be bubbled up in the list
				scores[scoreMap[agg.Table]].Score--

				// The parent repo needs to be pushed down the list
				scores[scoreMap[k]].Score++

				repoMap[agg.Table] = struct{}{}
			}
		}
	}

	// Sort the items by their scores
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score < scores[j].Score
	})

	for k := range scores {

		var score = scores[k]

		caches[k] = &RepoBootstrapConfig{
			Config: cache[score.Name],
			Name:   score.Name,
		}
	}

	var data = struct {
		BasePackage string
		Caches      []*RepoBootstrapConfig
	}{
		BasePackage: basePackage,
		Caches:      caches,
	}

	if e = repoBootstrapFileTpl.Execute(&buf, data); e != nil {
		panic(e)
	}

	e = ioutil.WriteFile(p, buf.Bytes(), lib.DefaultFileMode)

	return e

}

var repoBootstrapFileTpl = template.Must(template.New("repo-bootstrap").Funcs(template.FuncMap{
	"toArgName": toArgName,
}).Parse(`// Generated Code; DO NOT EDIT.
package repos

import (
	"{{ .BasePackage }}/gen/caches" 
	"{{ .BasePackage }}/gen/definitions"
	"{{ .BasePackage }}/core/components/config"
	"{{ .BasePackage }}/core/utils/hashids"
)

// Cache is a container for all cache providers
type Repos struct { {{range $cache := .Caches}}
	{{$cache.Name}} *{{$cache.Name}}Repo{{end}}
}

// BootstrapRepos bootstraps the repos
func BootstrapRepos(
	caches *caches.Caches,
	dal *definitions.DAL, 
	config *config.Config, 
	idHasher *hashids.IDHasher, 
) *Repos {

	{{range $cache := .Caches}} 
	var {{$cache.Name | toArgName}}Repo = New{{$cache.Name}}Repo(
		config, 
		caches.{{$cache.Name}}, 
		dal.{{$cache.Name}},
		{{ if $cache.Config.HasHashID }}idHasher,{{end}}{{if $cache.Config.Aggregate }}{{range $agg := $cache.Config.Aggregate.Properties}}
		{{$agg.Table | toArgName}}Repo,{{end}}{{end}}
	)
	{{end}}

	return &Repos{ {{range $cache := .Caches}}
		{{$cache.Name}}: {{$cache.Name | toArgName}}Repo, {{end}}
	}
}`))
