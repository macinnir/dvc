package gen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/macinnir/dvc/core/lib"
)

// GenAPITests generates api tests
func (g *Gen) GenAPITests() (e error) {

	fileBytes, e := ioutil.ReadFile(lib.RoutesFilePath)
	if e != nil {
		return
	}

	c := &lib.RoutesJSONContainer{}
	if e = json.Unmarshal(fileBytes, c); e != nil {
		return
	}

	for k := range c.Routes {

		fmt.Printf("Testing %s \n", k)
		fmt.Printf("Body: ")

	}

	// fmt.Println(c.Aggregates["ConversationAggregate"])

	// for k := range c.Aggregates {
	// 	fmt.Println(k)

	// }

	// b := false
	// for k := range c.Routes {
	// 	controller := c.Routes[k]
	// 	for l := range controller.Routes {
	// 		// if controller.Routes[l].Name != "CreateBusiness" {
	// 		// 	continue
	// 		// }
	// 		// b = true

	// 		bodyCat, bodyType := getTypeFullKey(controller.Routes[l].BodyType)

	// 		responseCat, responseType := getTypeFullKey(controller.Routes[l].ResponseType)

	// 		fmt.Println(
	// 			controller.Routes[l].Name,
	// 			controller.Routes[l].Method,
	// 			controller.Routes[l].Raw,
	// 			controller.Routes[l].BodyType,
	// 			controller.Routes[l].ResponseType,
	// 			fmt.Sprintf("%s.%s", bodyCat, bodyType),
	// 			fmt.Sprintf("%s.%s", responseCat, responseType),
	// 		)

	// 		// break
	// 	}

	// 	if b {
	// 		break
	// 	}
	// }

	return

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
