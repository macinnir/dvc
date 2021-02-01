package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/macinnir/dvc/lib"
	"github.com/macinnir/dvc/modules/gen"
)

// Test tests an endpoint
func (c *Cmd) Test(args []string) {

	if len(args) > 0 && args[0] == "help" {
		helpTest()
		return
	}

	if len(args) < 1 {
		lib.Error("Missing command endpoint", c.Options)
		os.Exit(1)
	}

	endpoint := args[0]

	fmt.Printf("Testing %s\n", endpoint)

	// load routes
	fileBytes, e := ioutil.ReadFile("meta/routes.json")
	if e != nil {
		return
	}

	j := &gen.RoutesJSONContainer{}
	if e = json.Unmarshal(fileBytes, j); e != nil {
		return
	}

	if _, ok := j.Routes[endpoint]; !ok {
		lib.Errorf("Route `%s` does not exist", c.Options, endpoint)
		os.Exit(1)
	}

	route := j.Routes[endpoint]

	fmt.Printf("%s => %s\n", route.Method, route.Raw)

	if len(route.BodyType) > 0 {

		bodyBaseType, bodyTypeName := getTypeFullKey(route.BodyType)
		body := map[string]string{}
		switch bodyBaseType {
		case "dtos":
			body = j.DTOs[bodyTypeName]
		case "aggregates":
			body = j.Aggregates[bodyTypeName]
		case "models":
			body = j.Models[bodyTypeName]
		}

		fmt.Printf("Body: %s\n", bodyTypeName)
		for k := range body {
			fmt.Printf("\t%s: %s\n", k, body[k])
		}

	}

	returnBaseType, returnTypeName := getTypeFullKey(route.ResponseType)
	response := map[string]string{}
	switch returnBaseType {
	case "dtos":
		response = j.DTOs[returnTypeName]
	case "aggregates":
		response = j.Aggregates[returnTypeName]
	case "models":
		response = j.Models[returnTypeName]
	}

	fmt.Printf("Response: %s\n", route.ResponseType)
	for k := range response {
		fmt.Printf("\t%s: %s\n", k, response[k])
	}

	// var e error
	// var sql string

	// cmp := c.initCompare()

	// if sql, e = cmp.ExportSchemaToSQL(); e != nil {
	// 	lib.Error(e.Error(), c.Options)
	// 	os.Exit(1)
	// }

	// fmt.Println(sql)
}

func helpTest() {
	fmt.Println(`
	test [endpoint name] test an endpoint name 
	`)
}

func getTypeFullKey(name string) (baseType string, typeName string) {

	if len(name) == 0 {
		return
	}

	for {
		if len(name) == 0 {
			return
		}
		if name[0:1] == "*" {
			name = name[1:]
			continue
		}

		if name[0:2] == "[]" {
			name = name[2:]
			continue
		}
		break
	}

	if strings.Contains(name, ".") {

		parts := strings.Split(name, ".")
		baseType = parts[0]
		typeName = parts[1]

	} else {
		baseType = "aggregates"
		typeName = name
	}

	return
}
