package gen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strings"
	"unicode"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/fetcher"
)

func fetchAllPermissions(controllersDir string) (map[string]string, error) {

	cf := fetcher.NewControllerFetcher()
	permissionMap := LoadPermissionsFromJSON()
	controllers, _, e := cf.FetchAll()
	if e != nil {
		return nil, e
	}

	for k := range controllers {

		controller := controllers[k]
		// Extract the permissions from the controller
		permissionMap[controller.Name+"_View"] = "View " + controller.Name

		for k := range controller.Routes {
			permissionMap[controller.Routes[k].Permission] = controller.Routes[k].Description
		}
	}

	return permissionMap, nil
}

// LoadPermissionsFromJSON loads a set of permissions from a JSON file
func LoadPermissionsFromJSON() map[string]string {
	permissionMap := map[string]string{}
	fileBytes, e := ioutil.ReadFile(lib.PermissionsFile)
	if e != nil {
		panic(e)
	}
	json.Unmarshal(fileBytes, &permissionMap)
	return permissionMap
}

func GenTSPerms(config *lib.Config) (e error) {
	var permissionMap map[string]string
	permissionMap, e = fetchAllPermissions(config.Dirs.Controllers)
	if e != nil {
		return
	}
	str := BuildTypescriptPermissions(permissionMap)
	fmt.Println(str)
	return
}

func GenGoPerms(config *lib.Config) (e error) {

	var permissionMap map[string]string
	permissionMap, e = fetchAllPermissions(config.Dirs.Controllers)
	if e != nil {
		return
	}

	BuildPermissionsGoFile(permissionMap)
	return
}

func BuildPermissionsGoFile(permissionMap map[string]string) {

	permissions := make([]string, 0, len(permissionMap))

	for k := range permissionMap {
		if len(k) == 0 {
			continue
		}
		permissions = append(permissions, k)
	}

	sort.Strings(permissions)

	// for k := range permissions {

	// 	permission := permissions[k]
	// 	description := permissionMap[permissions[k]]

	// 	// fmt.Println(permission + ": " + description)

	// }

	var b strings.Builder

	b.WriteString(`// Generated Code; DO NOT EDIT.

	package permissions

	import (
		"github.com/macinnir/dvc/core/lib/utils"
	)
	
	const (
	`)
	for k := range permissions {

		// fmt.Println("Permission: ", k, permissions[k])
		permTitle := string(unicode.ToUpper(rune(permissions[k][0]))) + permissions[k][1:]
		b.WriteString(`	// ` + permTitle + ` permission grants the ability of "` + permissionMap[permissions[k]] + `"
	` + permTitle + ` utils.Permission = "` + permissions[k] + `"
`)
	}

	b.WriteString(`)
	
	// Permissions returns a slice of permissions 
	func Permissions() map[utils.Permission]string {
		return map[utils.Permission]string {
	`)

	for k := range permissions {
		permTitle := string(unicode.ToUpper(rune(permissions[k][0]))) + permissions[k][1:]
		b.WriteString(`		` + permTitle + `: "` + permissionMap[permissions[k]] + `",
`)
	}

	b.WriteString(`	}
	}	
	
	`)
	permissionsFilePath := path.Join("core", "definitions", "constants", "permissions", "permissions.go")
	var permissionsFileBytes []byte

	var e error
	permissionsFileBytes, e = lib.FormatCode(b.String())
	if e != nil {
		panic(e)
	}
	ioutil.WriteFile(permissionsFilePath, permissionsFileBytes, 0777)
}

// BuildTypescriptPermissions returns a formatted typescript file of permission constants
func BuildTypescriptPermissions(permissionMap map[string]string) string {

	permissions := make([]string, 0, len(permissionMap))
	for k := range permissionMap {
		if len(k) == 0 {
			continue
		}
		permissions = append(permissions, k)
	}

	sort.Strings(permissions)

	permissionsFile := "// Generated Code; DO NOT EDIT.\n\n"
	for k := range permissions {
		permission := strings.TrimSpace(permissions[k])
		permTitle := string(unicode.ToUpper(rune(permission[0]))) + permission[1:]
		permissionsFile += "// " + permTitle + " -- " + permissionMap[permission] + "\n"
		permissionsFile += "export const " + permTitle + "Permission = \"" + permission + "\";\n"
	}

	return permissionsFile
}
