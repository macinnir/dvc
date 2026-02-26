package gen

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/macinnir/dvc/core/lib"
)

func GenAPIData(config *lib.Config, routes *lib.RoutesJSONContainer) error {

	var b strings.Builder
	var e error
	var jsonBytes []byte
	if jsonBytes, e = json.Marshal(routes); e != nil {
		return e
	}

	b.WriteString(`package apidocs

func ShowAPIData() string { 
	return ` + "`")

	b.WriteString(string(jsonBytes))
	b.WriteString("`\n")
	b.WriteString("}\n")

	if e = os.WriteFile("gen/apidocs/apiData.go", []byte(b.String()), 0644); e != nil {
		return e
	}

	return nil

}
