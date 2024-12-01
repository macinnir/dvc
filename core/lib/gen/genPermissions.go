package gen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"text/template"
	"time"
	"unicode"

	"github.com/macinnir/dvc/core/lib"
)

func FetchAllPermissionsFromControllers(controllers []*lib.Controller) (map[string]string, error) {

	permissionMap := LoadPermissionsFromJSON()

	for k := range controllers {

		controller := controllers[k]
		// Extract the permissions from the controller
		permissionMap[controller.Name+"_View"] = "View " + controller.Name

		for k := range controller.Routes {
			if len(controller.Routes[k].Permission) == 0 {
				continue
			}
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

// 0.014938 seconds
// 0.005944
// 0.000568
func GenTSPerms(config *lib.Config, permissions []PermissionTplType) (e error) {

	var start = time.Now()
	lib.EnsureDir(config.TypescriptPermissionsPath)
	var tsPermissionsPath = path.Join(config.TypescriptPermissionsPath, "permissions.ts")
	BuildTypescriptPermissions(permissions, tsPermissionsPath)
	lib.LogAdd(start, "%d ts permissions to %s", len(permissions), tsPermissionsPath)

	return
}

// 0.018100
// 0.000900
func GenGoPerms(config *lib.Config, permissions []PermissionTplType) (e error) {
	var start = time.Now()
	var permissionsFilePath = path.Join(lib.GoPermissionsDir, "permissions.go")
	BuildPermissionsGoFile(permissions, permissionsFilePath)
	lib.LogAdd(start, "%d go permissions to %s", len(permissions), permissionsFilePath)
	return
}

var goPermissionsFileTemplate = template.Must(template.New("go-permissions-file").Parse(`// Generated Code; DO NOT EDIT.

package permissions

import (
	"github.com/macinnir/dvc/core/lib/utils"
)

const (
	{{range .Permissions}}
	// {{.Title}} permission grants the ability of "{{.Description}}"
	{{.Title}} utils.Permission = "{{.Name}}"
	{{end}}
)

// Permissions returns a slice of permissions 
func Permissions() map[utils.Permission]string {

	return map[utils.Permission]string {
		{{range .Permissions}} 
		{{.Title}}: "{{.Description}}",{{end}}
	}

}
`))

type PermissionTplType struct {
	Title       string
	Description string
	Name        string
}

func BuildTplPermissions(permissionMap map[string]string) []PermissionTplType {

	var perms = make([]PermissionTplType, len(permissionMap))

	var k = 0
	for permissionName := range permissionMap {
		perms[k] = PermissionTplType{
			Title:       string(unicode.ToUpper(rune(permissionName[0]))) + permissionName[1:],
			Description: permissionMap[permissionName],
			Name:        permissionName,
		}
		k++
	}

	sort.Slice(perms, func(i, j int) bool {
		return perms[i].Name < perms[j].Name
	})

	return perms
}

func BuildPermissionsGoFile(permissions []PermissionTplType, permissionsFilePath string) error {

	var e error

	var tplVals = struct {
		Permissions []PermissionTplType
	}{
		Permissions: permissions,
	}

	var buffer = bytes.Buffer{}

	goPermissionsFileTemplate.Execute(&buffer, tplVals)

	var formatted []byte
	if formatted, e = format.Source(buffer.Bytes()); e != nil {
		fmt.Println("Format Error:", e.Error())
		return e
	}

	if e = ioutil.WriteFile(permissionsFilePath, formatted, lib.DefaultFileMode); e != nil {
		fmt.Println("Write file error: ", e.Error())
		return e
	}

	return nil
	// 	for k := range permissions {

	// 		// fmt.Println("Permission: ", k, permissions[k])
	// 		permTitle := string(unicode.ToUpper(rune(permissions[k][0]))) + permissions[k][1:]
	// 		b.WriteString(`	// ` + permTitle + ` permission grants the ability of "` + permissionMap[permissions[k]] + `"
	// 	` + permTitle + ` utils.Permission = "` + permissions[k] + `"
	// `)
	// 	}

	// 	b.WriteString(`)

	// 	// Permissions returns a slice of permissions
	// 	func Permissions() map[utils.Permission]string {
	// 		return map[utils.Permission]string {
	// 	`)

	// 	for k := range permissions {
	// 		permTitle := string(unicode.ToUpper(rune(permissions[k][0]))) + permissions[k][1:]
	// 		b.WriteString(`		` + permTitle + `: "` + permissionMap[permissions[k]] + `",
	// `)
	// 	}

	// 	b.WriteString(`	}
	// 	}

	// 	`)
	// 	// permissionsFilePath := path.Join("core", "definitions", "constants", "permissions", "permissions.go")
	// 	var permissionsFileBytes []byte

	// 	var e error
	// 	permissionsFileBytes, e = lib.FormatCode(b.String())
	// 	if e != nil {
	// 		panic(e)
	// 	}
	// 	ioutil.WriteFile(permissionsFilePath, permissionsFileBytes, 0777)
}

var tsPermissionsFileTemplate = template.Must(template.New("ts-permissions-file").Parse(`// Generated Code; DO NOT EDIT.
{{range .Permissions}}
// {{.Title}} -- {{.Description}}
export const {{.Title}}Permission = "{{.Name}}";
{{end}}`))

// BuildTypescriptPermissions returns a formatted typescript file of permission constants
func BuildTypescriptPermissions(permissions []PermissionTplType, permissionsFilePath string) error {

	var e error
	var tplVals = struct {
		Permissions []PermissionTplType
	}{
		Permissions: permissions,
	}
	var buffer = bytes.Buffer{}

	tsPermissionsFileTemplate.Execute(&buffer, tplVals)

	if e = ioutil.WriteFile(permissionsFilePath, buffer.Bytes(), lib.DefaultFileMode); e != nil {
		fmt.Println("Write file error: ", e.Error())
		return e
	}

	return nil

	// permissionsFile := "// Generated Code; DO NOT EDIT.\n\n"
	// for k := range permissions {
	// 	permission := strings.TrimSpace(permissions[k])
	// 	permTitle := string(unicode.ToUpper(rune(permission[0]))) + permission[1:]
	// 	permissionsFile += "// " + permTitle + " -- " + permissionMap[permission] + "\n"
	// 	permissionsFile += "export const " + permTitle + "Permission = \"" + permission + "\";\n"
	// }

	// return permissionsFile
}
