package test

import (
	"errors"
	"strings"

	"github.com/macinnir/dvc/core/lib"
	"go.uber.org/zap"
)

var (
	ErrMissingEndpoint = errors.New("missing endpoint")
	ErrEndpointNoExist = errors.New("endpoint does not exist")
)

const CommandName = "test"

// Test tests an endpoint
func Cmd(log *zap.Logger, config *lib.Config, args []string) error {

	return nil
	// if len(args) < 1 {
	// 	log.Error("Missing command endpoint")
	// 	return ErrMissingEndpoint
	// }

	// endpoint := args[0]

	// log.Info("Testing", zap.String("endpoint", endpoint))

	// // load routes
	// fileBytes, e := ioutil.ReadFile(lib.RoutesFilePath)
	// if e != nil {
	// 	log.Error("Error reading routes file", zap.String("path", lib.RoutesFilePath), zap.Error(e))
	// 	return e
	// }

	// j := &gen.RoutesJSONContainer{}
	// if e = json.Unmarshal(fileBytes, j); e != nil {
	// 	log.Error("Error unmarshalling routes file", zap.Error(e))
	// 	return e
	// }

	// if _, ok := j.Routes[endpoint]; !ok {
	// 	log.Error(ErrEndpointNoExist.Error(), zap.String("endpoint", endpoint))
	// 	return ErrEndpointNoExist
	// }

	// route := j.Routes[endpoint]

	// fmt.Printf("%s => %s\n", route.Method, route.Raw)

	// if len(route.BodyType) > 0 {

	// 	bodyBaseType, bodyTypeName := getTypeFullKey(route.BodyType)
	// 	body := map[string]string{}
	// 	switch bodyBaseType {
	// 	case "dtos":
	// 		body = j.DTOs[bodyTypeName]
	// 	case "aggregates":
	// 		body = j.Aggregates[bodyTypeName]
	// 	case "models":
	// 		body = j.Models[bodyTypeName]
	// 	}

	// 	fmt.Printf("Body: %s\n", bodyTypeName)
	// 	for k := range body {
	// 		fmt.Printf("\t%s: %s\n", k, body[k])
	// 	}

	// }

	// returnBaseType, returnTypeName := getTypeFullKey(route.ResponseType)
	// response := map[string]string{}
	// switch returnBaseType {
	// case "dtos":
	// 	response = j.DTOs[returnTypeName]
	// case "aggregates":
	// 	response = j.Aggregates[returnTypeName]
	// case "models":
	// 	response = j.Models[returnTypeName]
	// }

	// log.Infof("Response: %s\n", route.ResponseType)
	// for k := range response {
	// 	fmt.Printf("\t%s: %s\n", k, response[k])
	// }

	// // var e error
	// // var sql string

	// // cmp := c.initCompare()

	// // if sql, e = cmp.ExportSchemaToSQL(); e != nil {
	// // 	lib.Error(e.Error(), c.Options)
	// // 	os.Exit(1)
	// // }

	// // fmt.Println(sql)

	// return nil
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
