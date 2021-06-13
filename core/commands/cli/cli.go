package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/TylerBrock/colorjson"
	"github.com/macinnir/dvc/core/lib"
	"go.uber.org/zap"
)

const CommandName = "cli"

// Cmd adds an object to the database
func Cmd(log *zap.Logger, config *lib.Config, args []string) error {

	fmt.Println("Loading the local schema...")

	// Load the schema
	fileBytes, e := ioutil.ReadFile(lib.RoutesFilePath)
	if e != nil {
		return e
	}

	c := &lib.RoutesJSONContainer{}
	if e = json.Unmarshal(fileBytes, c); e != nil {
		return e
	}

	reader := bufio.NewReader(os.Stdin)

	host := lib.ReadCliInput(reader, "Host(http://localhost:8080/api)> ")

	if len(strings.TrimSpace(host)) == 0 {
		host = "http://localhost:8080/api"
	}

	authHeader := lib.ReadCliInput(reader, "AuthHeader> ")

	if len(strings.TrimSpace(authHeader)) == 0 {
		authHeader = "HK3w2iQWijOFSSxH6nbsCdSb8HeoYqSwD3XUu62d"
		// return fmt.Errorf("auth header cannot be empty")
	}

	cmd := ""

	for {

		for {
			cmd = lib.ReadCliInput(reader, "Cmd> ")
			if len(cmd) == 0 {
				fmt.Println("Command cannot be empty, please provide a command")
				continue
			}

			if _, ok := c.Routes[cmd]; !ok {

				for k := range c.Routes {
					if len(c.Routes[k].Name) <= len(cmd) {
						continue
					}

					// Did you type the first part of an existing command?
					if c.Routes[k].Name[0:len(cmd)] != cmd {
						continue
					}

					fmt.Println("> " + c.Routes[k].Name)
				}

				fmt.Println("Invalid command...")
				continue
			}
			break
		}

		route := c.Routes[cmd]
		requestURL := extractBasePath(host, route)

		if len(route.Params) > 0 {
			// params := make([]string, len(route.Params))
			for k := range route.Params {
				paramValue := lib.ReadCliInput(reader, fmt.Sprintf("URL Param `%s` (%s):> ", route.Params[k].Name, route.Params[k].Type))
				requestURL = applyParam(requestURL, route.Params[k], paramValue)
			}
		}

		if len(route.Queries) > 0 {
			for k := range route.Queries {
				queryValue := lib.ReadCliInput(reader, fmt.Sprintf("URL Query `%s` (%s):> ", route.Queries[k].Name, route.Queries[k].Type))
				requestURL = applyQuery(requestURL, route.Queries[k], queryValue)
			}
		}

		fmt.Println("CMD: " + route.Name)
		fmt.Println("URL: ", requestURL)

		var response *http.Response

		if route.Method == "GET" {
			request, e := http.NewRequest("GET", requestURL, nil)
			if e != nil {
				fmt.Println("ERROR: ", e.Error())
				continue
			}

			request.Header.Set("Content-Type", "application/json")

			if route.IsAuth {
				request.Header.Set("Authorization", "Bearer "+authHeader)
			}

			client := &http.Client{}
			response, e = client.Do(request)
			if e != nil {
				fmt.Println("ERROR: ", e.Error())
				continue
			}

			defer response.Body.Close()
			var responseBody []byte
			responseBody, _ = ioutil.ReadAll(response.Body)
			fmt.Printf("Status Code: %d\n", response.StatusCode)
			// fmt.Println(string(responseBody))

			var obj map[string]interface{}
			json.Unmarshal(responseBody, &obj)

			f := colorjson.NewFormatter()
			f.Indent = 4
			s, _ := f.Marshal(obj)
			fmt.Println(string(s))

		}
	}
	// fmt.Println(path.Join(host, route.Path))

	// for k := range c.Routes {

	// 	fmt.Printf("Route %s \n", k)
	// 	fmt.Printf("Body: ")

	// }

	// fmt.Println("This is the CLI")
	return nil
}

func applyParam(url string, param lib.ControllerRouteParam, value string) string {
	r := regexp.MustCompile(`{` + param.Name + `:.*?}`)
	return r.ReplaceAllString(url, value)
}

func applyQuery(url string, param lib.ControllerRouteQuery, value string) string {
	r := regexp.MustCompile(`{` + param.Name + `:.*?}`)
	return r.ReplaceAllString(url, value)
}

func extractBasePath(host string, route *lib.ControllerRoute) string {

	fmt.Println("Name: ", route.Name)
	// Remove last slash from host
	if host[len(host)-1:] == "/" {
		host = host[0 : len(host)-1]
	}

	// Add slash if route doesn't have one
	path := ""

	if len(route.Raw) > 0 {
		path = route.Raw
	} else {
		path = route.Path
	}

	// b, _ := json.MarshalIndent(route, "", "   ")
	// fmt.Println(string(b))
	// fmt.Println("Route Name: ", route.Name)
	// fmt.Println("Route Path:", route.Path)
	if path[0:1] != "/" {
		path = "/" + path
	}

	// if strings.Contains(path, "?") {
	// 	path = strings.Split(path, "?")[0]
	// }

	return host + path
}
