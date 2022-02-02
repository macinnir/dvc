package gen

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"unicode"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/fetcher"
)

func fetchAllPermissionsFromControllers(controllersDir string) (map[string]string, error) {

	cf := fetcher.NewControllerFetcher()
	controllers, _, e := cf.FetchAll()

	if e != nil {
		return nil, e
	}

	permissionMap := LoadPermissionsFromJSON()

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
	var fileBytes []byte

	// Core permissions
	if _, e := os.Stat(lib.CorePermissionsFile); !os.IsNotExist(e) {
		fileBytes, _ = ioutil.ReadFile(lib.CorePermissionsFile)
		json.Unmarshal(fileBytes, &permissionMap)
	}

	// User permissions
	if _, e := os.Stat(lib.PermissionsFile); !os.IsNotExist(e) {
		fileBytes, _ = ioutil.ReadFile(lib.PermissionsFile)
		userPermissions := map[string]string{}
		json.Unmarshal(fileBytes, &userPermissions)

		for k := range userPermissions {
			permissionMap[k] = userPermissions[k]
		}
	}

	return permissionMap
}

func GenTSPerms(config *lib.Config) (e error) {

	lib.EnsureDir(config.TypescriptPermissionsPath)

	var permissionMap map[string]string
	permissionMap, e = fetchAllPermissionsFromControllers(config.Dirs.Controllers)
	if e != nil {
		return
	}
	str := BuildTypescriptPermissions(permissionMap)
	e = ioutil.WriteFile(path.Join(config.TypescriptPermissionsPath, "permissions.ts"), []byte(str), 0777)
	return
}

func GenGoPerms(config *lib.Config) (e error) {

	var permissionMap map[string]string
	permissionMap, e = fetchAllPermissionsFromControllers(config.Dirs.Controllers)
	if e != nil {
		return
	}

	BuildPermissionsGoFile(permissionMap)
	return
}

func BuildSettingsGoFile(permissionMap map[string]string) {

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
	lib.EnsureDir(lib.GoPermissionsPath)
	// permissionsFilePath := path.Join("core", "definitions", "constants", "permissions", "permissions.go")
	permissionsFilePath := path.Join(lib.GoPermissionsPath, "permissions.go")
	var permissionsFileBytes []byte

	var e error
	permissionsFileBytes, e = lib.FormatCode(b.String())
	if e != nil {
		panic(e)
	}
	ioutil.WriteFile(permissionsFilePath, permissionsFileBytes, 0777)
}

// BuildTypescriptSettings returns a formatted typescript file of setting constants
func BuildTypescriptSettings(settingsMap map[string]string) string {

	settings := make([]string, 0, len(settingsMap))
	for k := range settingsMap {
		if len(k) == 0 {
			continue
		}
		settings = append(settings, k)
	}

	sort.Strings(settings)

	settingsFile := "// Generated Code; DO NOT EDIT.\n\n"
	for k := range settings {
		setting := strings.TrimSpace(settings[k])
		settingTitle := string(unicode.ToUpper(rune(setting[0]))) + setting[1:]
		settingsFile += "// " + settingTitle + " -- " + settingsMap[setting] + "\n"
		settingsFile += "export const " + settingTitle + "Setting = \"" + setting + "\";\n"
	}

	return settingsFile
}
