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
)

func loadPermissions() map[string]string {
	permissionMap := map[string]string{}
	fileBytes, e := ioutil.ReadFile(lib.PermissionsFile)
	if e != nil {
		panic(e)
	}
	json.Unmarshal(fileBytes, &permissionMap)
	return permissionMap
}

func GenPermissionsGoFile() {

	permissionMap := loadPermissions()
	permissions := make([]string, 0, len(permissionMap))
	for k := range permissionMap {
		permissions = append(permissions, k)
	}

	sort.Strings(permissions)

	// for k := range permissions {

	// 	permission := permissions[k]
	// 	description := permissionMap[permissions[k]]

	// 	// fmt.Println(permission + ": " + description)

	// }

	permissionsFile := `// Generated Code; DO NOT EDIT.

	package permissions

	import (
		"github.com/macinnir/dvc/core/lib/utils"
	)
	
	const (
	`
	for k := range permissions {
		permTitle := string(unicode.ToUpper(rune(permissions[k][0]))) + permissions[k][1:]
		permissionsFile += "\t// " + permTitle + " Permission is the `" + permissions[k] + "` permission\n"
		permissionsFile += "\t" + permTitle + " utils.Permission = \"" + permissions[k] + "\"\n"
	}

	permissionsFile += `)
	
	// Permissions returns a slice of permissions 
	func Permissions() map[utils.Permission]string {
		return map[utils.Permission]string {
	`

	for k := range permissions {
		permTitle := string(unicode.ToUpper(rune(permissions[k][0]))) + permissions[k][1:]
		permissionsFile += "\t\t" + permTitle + ": \"" + permissionMap[permissions[k]] + "\",\n"
	}

	permissionsFile += `	}
	}	
	
	`
	permissionsFilePath := path.Join("core", "definitions", "constants", "permissions", "permissions.go")
	fmt.Println("Writing the permissions file to path ", permissionsFilePath)
	var permissionsFileBytes []byte
	var e error
	permissionsFileBytes, e = lib.FormatCode(permissionsFile)
	if e != nil {
		panic(e)
	}
	ioutil.WriteFile(permissionsFilePath, []byte(permissionsFileBytes), 0777)
}

// BuildTypescriptPermissions returns a formatted typescript file of permission constants
func BuildTypescriptPermissions() string {

	permissionMap := map[string]string{}
	fileBytes, e := ioutil.ReadFile(lib.PermissionsFile)
	if e != nil {
		panic(e)
	}
	json.Unmarshal(fileBytes, &permissionMap)

	permissions := make([]string, 0, len(permissionMap))
	for k := range permissionMap {
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
