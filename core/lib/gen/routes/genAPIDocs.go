package routes

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strings"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/gen"
	"github.com/macinnir/dvc/core/lib/schema"
)

func cleanObject(obj string) string {

	if len(obj) == 0 {
		return ""
	}

	obj = schema.ExtractBaseGoType(obj)

	if len(obj) > 8 && obj[0:8] == "appdtos." {
		obj = obj[8:]
	}

	if len(obj) > 9 && obj[0:9] == "*appdtos." {
		obj = obj[9:]
	}

	if len(obj) > 5 && obj[0:5] == "dtos." {
		obj = obj[5:]
	}

	if len(obj) > 11 && obj[0:11] == "aggregates." {
		obj = obj[11:]
	}

	if len(obj) > 12 && obj[0:12] == "*aggregates." {
		obj = obj[12:]
	}

	if len(obj) > 7 && obj[0:7] == "models." {
		obj = obj[7:]
	}

	if len(obj) > 8 && obj[0:8] == "*models." {
		obj = obj[8:]
	}

	return obj

}

func isArray(obj string) bool {
	return len(obj) > 2 && obj[0:2] == "[]"
}

func isComplexObject(obj string) bool {

	if (len(obj) > 5 && obj[0:5] == "dtos.") || (len(obj) > 6 && obj[0:6] == "*dtos.") || (len(obj) > 8 && obj[0:8] == "appdtos.") || (len(obj) > 9 && obj[0:9] == "*appdtos.") {
		return true
	}

	if (len(obj) > 11 && obj[0:11] == "aggregates.") || (len(obj) > 12 && obj[0:12] == "*aggregates.") {
		return true
	}

	if (len(obj) > 7 && obj[0:7] == "models.") || (len(obj) > 8 && obj[0:8] == "*models.") {
		return true
	}

	if len(obj) > 9 && obj[len(obj)-9:] == "Aggregate" {
		return true
	}

	if len(obj) > 3 && obj[len(obj)-3:] == "DTO" {
		return true
	}

	return false
}

func getBaseType(obj string) string {

	obj = schema.ExtractBaseGoType(obj)

	if obj == "int" || obj == "int64" {
		return "int"
	}

	if obj == "float64" {
		return "float"
	}

	if obj == "string" || obj == "null.String" {
		return "string"
	}

	if (len(obj) > 5 && obj[0:5] == "dtos.") || (len(obj) > 6 && obj[0:6] == "*dtos.") || (len(obj) > 8 && obj[0:8] == "appdtos.") || (len(obj) > 9 && obj[0:9] == "*appdtos.") {
		return "dto"
	}

	if (len(obj) > 11 && obj[0:11] == "aggregates.") || (len(obj) > 12 && obj[0:12] == "*aggregates.") {
		return "aggregate"
	}

	if (len(obj) > 7 && obj[0:7] == "models.") || (len(obj) > 8 && obj[0:8] == "*models.") {
		return "model"
	}

	if len(obj) > 9 && obj[len(obj)-9:] == "Aggregate" {
		return "aggregate"
	}

	if len(obj) > 3 && obj[len(obj)-3:] == "DTO" {
		return "dto"
	}

	return "model"
}

func defaultDataType(t string) string {
	switch t {
	case "int", "int64":
		return "0"
	case "string", "null.String":
		return `"string"`
	case "float64":
		return "0.0"
	default:
		return t

	}
}

func defaultDataTypeCSSClass(t string) string {
	switch t {
	case "int", "int64":
		return "datatype-number"
	case "string", "null.String":
		return `datatype-string`
	case "float64":
		return "datatype-number"
	default:
		return t

	}
}

func baseObjectToString(routes *lib.RoutesJSONContainer, obj map[string]string, embedded bool, level int) string {

	var names = []string{}
	for name := range obj {
		names = append(names, name)
	}

	sort.Strings(names)

	var sb strings.Builder

	if !embedded {
		sb.WriteString(`{`)
	}

	var n = 0
	for _, name := range names {

		value := obj[name]

		if len(name) > 9 && name[0:9] == "#embedded" {
			sb.WriteString(objectToString(routes, value, true, level))
		} else if len(value) == 17 && value == "map[string]string" {
			sb.WriteString("\n" + strings.Repeat("\t", level+1) + `"` + name + `": { "string": "string" }`)
		} else if (len(value) == 7 && value == "[]int64") || (len(value) == 6 && value == "[]int") {
			sb.WriteString("\n" + strings.Repeat("\t", level+1) + `"` + name + `": [0]`)
		} else if len(value) == 9 && value == "[]float64" {
			sb.WriteString("\n" + strings.Repeat("\t", level+1) + `"` + name + `": [0.0]`)
		} else if len(value) == 8 && value == "[]string" {
			sb.WriteString("\n" + strings.Repeat("\t", level+1) + `"` + name + `": ["string"]`)
		} else if len(value) > 2 && value[0:2] == "[]" {
			sb.WriteString("\n" + strings.Repeat("\t", level+1) + `"` + name + `" : [` + "\n" + strings.Repeat("\t", level+2) + "{" + objectToString(routes, value, true, level+2) + "\n" + strings.Repeat("\t", level+2) + "}\n" + strings.Repeat("\t", level+1) + "]")
		} else if isComplexObject(value) {
			sb.WriteString("\n" +
				strings.Repeat("\t", level+1) + `"` + name + `" : {` +
				objectToString(routes, value, true, level+1) + "\n" +
				strings.Repeat("\t", level+1) + "}")

		} else {

			dataType := defaultDataType(value)
			class := defaultDataTypeCSSClass(value)
			sb.WriteString("\n" + strings.Repeat("\t", level+1) + `"` + name + `": <span class="` + class + `">` + dataType + `</span>`)
		}

		if n < len(obj)-1 {
			sb.WriteString(`,`)
		}
		n++
	}
	if !embedded {
		sb.WriteString(`
}`)
	}

	return sb.String()
}

func baseObjectToTable(routes *lib.RoutesJSONContainer, obj map[string]string) string {

	var names = []string{}
	for name := range obj {
		names = append(names, name)
	}

	sort.Strings(names)

	var sb strings.Builder

	for k := range names {

		var name = names[k]
		value := obj[name]
		simpleType := schema.ExtractBaseGoType(value)

		if len(name) > 9 && name[0:9] == "#embedded" {
			sb.WriteString(objectToTable(routes, value))
			continue
		}

		sb.WriteString(`
			<tr><td style="width: 300px;">` + name + `</td><td><span class="field-type">`)
		if isComplexObject(simpleType) {
			sb.WriteString(`<a href="#type-` + cleanObject(simpleType) + `">` + value + `</a>`)
		} else {
			sb.WriteString(value)
		}

		sb.WriteString(`</span></td></tr>`)
	}

	return sb.String()
}

func objectToString(routes *lib.RoutesJSONContainer, obj string, embedded bool, level int) string {

	var baseType = getBaseType(obj)
	var objName = cleanObject(obj)

	switch baseType {
	case "dto":
		if _, ok := routes.DTOs[objName]; !ok {

			return ""
		}

		return baseObjectToString(routes, routes.DTOs[objName], embedded, level)

	case "model":
		if _, ok := routes.Models[objName]; !ok {
			fmt.Println("Can't find model:", obj, "("+objName+")")
			return ""
		}

		return baseObjectToString(routes, routes.Models[objName], embedded, level)
	case "aggregate":
		if _, ok := routes.Aggregates[objName]; !ok {
			return ""
		}

		return baseObjectToString(routes, routes.Aggregates[objName], embedded, level)
	case "int":
		return "0"
	case "string":

	}

	return ""

}

func objectToTable(routes *lib.RoutesJSONContainer, obj string) string {

	var baseType = getBaseType(obj)
	var objName = cleanObject(obj)
	var objType = map[string]string{}

	switch baseType {
	case "dto":
		if _, ok := routes.DTOs[objName]; !ok {
			return ""
		}
		objType = routes.DTOs[objName]

	case "model":
		if _, ok := routes.Models[objName]; !ok {
			return ""
		}
		objType = routes.Models[objName]

	case "aggregate":
		if _, ok := routes.Aggregates[objName]; !ok {
			return ""
		}
		objType = routes.Aggregates[objName]
	}

	return baseObjectToTable(routes, objType)

	// if (len(obj) > 5 && obj[0:5] == "dtos.") || (len(obj) > 6 && obj[0:6] == "*dtos.") || (len(obj) > 7 && obj[0:7] == "models.") || (len(obj) > 8 && obj[0:8] == "*models.") {

	// }

}

func goTypeToSwiftType(goType string) string {

	switch goType {
	case "int":
		return "Int"
	case "int64":
		return "Int64"
	case "string":
		return "String"
	case "null.String":
		return "String?"
	case "float64", "null.Float":
		return "Float64"
	}

	isArray := false
	if len(goType) > 2 && goType[0:2] == "[]" {
		isArray = true
	}

	goType = cleanObject(goType)

	if isArray {
		return "[" + goType + "]"
	}

	return goType
}

func baseObjectToSwiftSource(routes *lib.RoutesJSONContainer, obj map[string]string) string {

	var names = []string{}
	for name := range obj {
		names = append(names, name)
	}

	sort.Strings(names)

	var sb strings.Builder

	for k := range names {

		var name = names[k]
		value := obj[name]
		// simpleType := schema.ExtractBaseGoType(value)

		if len(name) > 9 && name[0:9] == "#embedded" {
			sb.WriteString(objectToSwiftSource(routes, value))
			continue
		}

		sb.WriteString(`
	var ` + name + ` : ` + goTypeToSwiftType(value))
		// if isComplexObject(simpleType) {
		// 	sb.WriteString(`<a href="#type-` + cleanObject(simpleType) + `">` + value + `</a>`)
		// } else {
		// 	sb.WriteString(value)
		// }

		// sb.WriteString(`</span></td></tr>`)
	}

	return sb.String()
}

func objectToSwiftSource(routes *lib.RoutesJSONContainer, obj string) string {

	// struct Account : Decodable {
	//    var AccountID: Int64
	//    var CreatedBy: Int64
	//    var DateCreated: Int64
	//    var IsDeleted: Int
	//    var MaxUserCount: Int64
	//    var OwnerID: Int64
	//    var Title: String
	//    var UserCount: Int64
	// }

	var baseType = getBaseType(obj)
	var objName = cleanObject(obj)

	var objType = map[string]string{}

	switch baseType {
	case "dto":
		if _, ok := routes.DTOs[objName]; !ok {
			return ""
		}
		objType = routes.DTOs[objName]

	case "model":
		if _, ok := routes.Models[objName]; !ok {
			return ""
		}
		objType = routes.Models[objName]

	case "aggregate":
		if _, ok := routes.Aggregates[objName]; !ok {
			return ""
		}
		objType = routes.Aggregates[objName]
	}

	return baseObjectToSwiftSource(routes, objType)

}

func buildRoutePathHTML(route *lib.ControllerRoute) string {

	path := route.Raw

	if len(route.Params) > 0 {

		for _, param := range route.Params {

			pattern := "{" + param.Name + ":" + param.Pattern + "}"

			path = strings.Replace(path, pattern, `<span class="endpoint-query-var">{`+param.Name+`}</span>`, 1)

		}

	}

	if len(route.Queries) > 0 {

		for _, query := range route.Queries {

			pattern := query.Name + "={" + query.Name + ":" + query.Pattern + "}"
			path = strings.Replace(path, pattern, `<span class="endpoint-query-var">`+query.Name+`</span>=`+query.Pattern, 1)

		}

	}

	return path

	// if strings.Contains(`/accounts/<span class="endpoint-query-var">{accountID}</span>/users?<span class="endpoint-query-var">page</span>=[0-9]+&<span class="endpoint-query-var">limit</span>=[0-9]+</span>
}

func GenAPIDocs(config *lib.Config, routes *lib.RoutesJSONContainer) {

	dir := "gen/apidocs"
	lib.EnsureDir(dir)

	outFile := path.Join(dir, "apidocs.go")

	var sb strings.Builder

	sb.WriteString(htmlHeader())

	var controllerNames = []string{}
	var controllers = map[string][]string{}
	var methodCount = len(routes.Routes)

	for k := range routes.Routes {

		if _, ok := controllers[routes.Routes[k].ControllerName]; !ok {
			controllers[routes.Routes[k].ControllerName] = []string{}
			controllerNames = append(controllerNames, routes.Routes[k].ControllerName)
		}

		controllers[routes.Routes[k].ControllerName] = append(controllers[routes.Routes[k].ControllerName], k)
	}

	sort.Strings(controllerNames)

	sb.WriteString(`
		<h3>Routes (` + fmt.Sprint(methodCount) + `)</h3>
		<ul>`)

	for k := range controllerNames {

		controllerName := controllerNames[k]

		sb.WriteString(`
					<li>
						<a href="#controller-` + controllerName + `" class="controller">` + controllerName + `</a>
						<ul>`)

		sort.Strings(controllers[controllerName])

		for routeName := range controllers[controllerName] {
			route := routes.Routes[controllers[controllerName][routeName]]
			sb.WriteString(`
							<li><a href="#route-` + route.Name + `" class="method">` + route.Name + `</a></li>`)
		}
		sb.WriteString(`
						</ul>
					</li>
		`)
	}

	typeNames := []string{}

	for aggregateName := range routes.Aggregates {
		typeNames = append(typeNames, aggregateName)
	}

	for dtoName := range routes.DTOs {
		typeNames = append(typeNames, dtoName)
	}

	for modelName := range routes.Models {
		typeNames = append(typeNames, modelName)
	}

	sb.WriteString(`
	</ul>
	<h3>Types (` + fmt.Sprint(len(typeNames)) + `)</h3>
	<ul>
`)

	sort.Strings(typeNames)

	for k := range typeNames {
		sb.WriteString(`
					<li><a href="#type-` + typeNames[k] + `">` + typeNames[k] + `</a></li>`)
	}

	sb.WriteString(`
				</ul>
			</div>
			<div id="main-right">
	`)

	var usages = map[string]map[string]struct{}{}

	for k := range controllerNames {
		sb.WriteString(`
				<div id="controller-` + controllerNames[k] + `" class="controller-container">
					<h2>` + controllerNames[k] + `</h2>
					<div class="inner" id="controller-` + controllerNames[k] + `-inner">`)

		for l := range controllers[controllerNames[k]] {

			route := routes.Routes[controllers[controllerNames[k]][l]]

			sb.WriteString(`
						<a name="route-` + route.Name + `"></a>
						<div id="route-` + route.Name + `" class="endpoint endpoint-` + strings.ToLower(route.Method) + `">
							<h3>
								<div class="endpoint-title-left">
									<span class="endpoint-method endpoint-` + strings.ToLower(route.Method) + `-method">` + route.Method + `</span>
									<span class="endpoint-title">` + route.Name + `</span>
									<span class="endpoint-path">` + buildRoutePathHTML(route) + `</span>
								</div>
								<div class="endpoint-title-right">
							`)
			if !route.IsAuth {
				sb.WriteString(`
									<div class="endpoint-perm endpoint-perm-anonymous"><i class="bi-unlock-fill"></i> Anonymous</div>`)
			} else {
				if len(route.Permission) > 0 {
					sb.WriteString(`
									<div class="endpoint-perm endpoint-perm-perm"><i class="bi-lock-fill"></i><a href="#permission-` + route.Permission + `">` + route.Permission + `</a></div>`)
				} else {
					sb.WriteString(`
									<div class="endpoint-perm endpoint-perm-anyone"><i class="bi-lock-fill"></i> Anyone</div>`)
				}
			}
			sb.WriteString(`
								</div>
							</h3>
							<div class="endpoint-description">` + route.Description + `</div>`)
			if len(route.Params) > 0 {
				sb.WriteString(`
							<h4>URL Parameters (` + fmt.Sprint(len(route.Params)) + `)</h4>
							<div class="endpoint-parameters">
								<table class="table">
									<thead>
										<tr>
											<th style="width: 300px;">Field</th>
											<th style="width: 300px;">Type</th>
											<th>Pattern</th>
										</tr>
									</thead>
									<tbody>`)

				for k := range route.Params {
					p := route.Params[k]
					sb.WriteString(`
										<tr><td><span class="endpoint-query-var">` + p.Name + `</a></td><td><span class="field-type">` + p.Type + `</span></td><td>` + p.Pattern + `</td></tr>
									`)
				}

				sb.WriteString(`
									</tbody>
								</table>

							</div>
				`)
			}

			if len(route.Queries) > 0 {
				sb.WriteString(`
							<h4>Query Parameters (` + fmt.Sprint(len(route.Queries)) + `)</h4>
							<div class="endpoint-parameters">
								<table class="table">
									<thead>
										<tr>
											<th style="width: 300px;">Field</th>
											<th style="width: 300px;">Type</th>
											<th>Pattern</th>
										</tr>
									</thead>
								<tbody>`)
				for k := range route.Queries {

					q := route.Queries[k]

					sb.WriteString(`
									<tr><td><span class="endpoint-query-var">` + q.Name + ` </a></td><td><span class="field-type">` + q.Type + `</span></td><td>` + q.Pattern + `</td></tr>
									`)
				}

				sb.WriteString(`
								</tbody>
							</table>
						</div>
				`)
			}

			if len(route.BodyType) > 0 {

				sb.WriteString(`
							<h4>Body</h4>
							<div class="endpoint-body">
	 							<div class="endpoint-body-inner source-code">`)

				sb.WriteString(objectToString(routes, route.BodyType, false, 0))

				bodyType := cleanObject(route.BodyType)

				if _, ok := usages[bodyType]; !ok {
					usages[bodyType] = map[string]struct{}{}
				}

				usages[bodyType][route.Name] = struct{}{}

				sb.WriteString(`</div>
								<div class="endpoint-body-footer"><a href="#type-` + bodyType + `">` + bodyType + `</a></div>
							</div>`)

			}
			if len(route.ResponseType) > 0 {

				var responseType = cleanObject(route.ResponseType)
				if _, ok := usages[responseType]; !ok {
					usages[responseType] = map[string]struct{}{}
				}

				usages[responseType][route.Name] = struct{}{}

				sb.WriteString(`
								<h4>Response</h4>
								<div class="endpoint-response">
									<div class="endpoint-response-inner source-code">`)
				sb.WriteString(objectToString(routes, route.ResponseType, false, 0))
				sb.WriteString(`</div>
									<div class="endpoint-response-footer">
										<a href="#type-` + responseType + `">` + responseType + `</a>
									</div>
								</div>`)
			}

			sb.WriteString(`
								<div class="endpoint-footer">
									
									<div class="endpoint-footer-left">
										<a href="javascript:void(0);" onclick="copyHashToClipboard('#route-` + route.Name + `')">
											Copy Link
										</a>
									</div>

									<div class="endpoint-footer-right">File: ` + route.FileName + `</div>

								</div>

							</div>
			`)
		}
		sb.WriteString(`
					</div>
				</div>`)
	}

	sb.WriteString(`<div id="types-container">`)
	for k := range typeNames {

		var rawName = typeNames[k]
		var name = cleanObject(rawName)

		sb.WriteString(`
						<div id="type-` + name + `" class="type">
							<h2>
								<div class="type-title-left">` + name + `</div>
								<div class="type-title-right"></div>
							</h2>
							<div class="type-inner" id="type-` + name + `-inner">
								<table class="field-table table">
									<thead>
										<tr>
											<th style="width: 300px;">Field</th>
											<th>Type</th>
										</tr>
									</thead>
									<tbody>
										` + objectToTable(routes, rawName) + `
									</tbody>
								</table>		
							</div>
							
							<div class="collapse" id="collapse-type-source-` + name + `-swift">
							<div class="source-code" id="collapse-type-source-` + name + `-swift-inner">struct ` + name + ` : Codable {` + objectToSwiftSource(routes, rawName) + `
}</div>
							</div>
							
							<div class="type-footer">
							
								<div class="type-footer-left">
									Swift: 
									<a data-bs-toggle="collapse" href="#collapse-type-source-` + name + `-swift">
										Show 
									</a> |
									<a href="javascript:void(0);" onclick="copyContentToClipboard('collapse-type-source-` + name + `-swift-inner')">
										Copy
									</a>
									&nbsp; | &nbsp;
									<a href="javascript:void(0);" onclick="copyHashToClipboard('#type-` + name + `')">
										Copy Link
									</a>
								</div>

								<div class="type-footer-right">File: </div>

							</div>`)

		if _, ok := usages[name]; ok {

			routeNames := []string{}

			for routeName := range usages[name] {
				routeNames = append(routeNames, routeName)
			}

			sort.Strings(routeNames)

			sb.WriteString(`<div class="type-used-in"><strong>Used In:</strong>&nbsp;`)

			for l := range routeNames {
				sb.WriteString(`<a href="#route-` + routeNames[l] + `">` + routeNames[l] + `</a>`)
				if l < len(routeNames)-1 {
					sb.WriteString("&nbsp;|&nbsp;")
				}
			}
			sb.WriteString(`</div>`)
		}
		sb.WriteString(`
						</div>

		`)
	}

	sb.WriteString(`</div><!-- /types-container -->
	<div id="permissions-container">
		<h2>Permissions</h2>
	`)

	permissionMap, _ := gen.FetchAllPermissionsFromControllers(config.Dirs.Controllers)

	permissionNames := make([]string, len(permissionMap))
	var n = 0
	for permissionName := range permissionMap {
		permissionNames[n] = permissionName
		n++
	}

	sort.Strings(permissionNames)

	sb.WriteString(`<table class="table">
		<thead>
			<tr>
				<th style="width: 300px;">Permission</th>
				<th>Description</th>
			</tr>
		</thead>
		<tbody>`)
	for k := range permissionNames {
		permissionName := permissionNames[k]

		sb.WriteString(`
			<tr class="permission">
				<td class="permission-title">
				<a name="permission-` + permissionName + `"></a>
				` + permissionName + `</td>
				<td class="permission-description">` + permissionMap[permissionName] + `</td>
			</tr>
		`)
	}

	sb.WriteString(`
		</tbody>
	</table>
	`)

	sb.WriteString(`
		</div><!-- /#permissions-container -->
	`)

	// 					<div id="type-WebsiteURLsAggregate" class="type">
	// 						<h2>
	// 							<div class="type-title-left">WebsiteURLsAggregate</div>
	// 							<div class="type-title-right"></div>
	// 						</h2>
	// 						<div class="type-inner" id="type-WebsiteURLsAggregate-inner">
	// 							<table class="field-table table"><thead><tr><th style="width: 300px;">Field</th><th>Type</th></tr></thead><tbody>
	// 								<tr><td style="width: 300px;">Count</td><td><span class="field-type">int64</span></td></tr>

	// 								<tr><td style="width: 300px;">Data</td><td><span class="field-type"><a href="#type-WebsiteURLAggregate">[]*WebsiteURLAggregate</a></span></td></tr>
	// 							</tbody></table>

	// 						</div>
	// 						<a class="btn btn-primary" data-bs-toggle="collapse" href="#collapse-type-source-WebsiteURLsAggregate-json">JSON</a> |
	// 						<a class="btn btn-primary" data-bs-toggle="collapse" href="#collapse-type-source-WebsiteURLsAggregate-kotlin">Kotlin</a> |
	// 						<script type="text/javascript">
	// 						<!--
	// 							<a class="btn" onclick="copyJSONSourceWebsiteURLsAggregate()">Copy</a>
	// 							function copyJSONSourceWebsiteURLsAggregate() {
	// 								var copyText = document.getElementById("collapse-type-source-WebsiteURLsAggregate-json-inner")
	// 								navigator.clipboard.writeText(copyText.innerText)
	// 							}

	// 							function copyKotlinSourceWebsiteURLsAggregate() {
	// 								var copyText = document.getElementById("collapse-type-source-WebsiteURLsAggregate-kotlin-inner")
	// 								navigator.clipboard.writeText(copyText.innerText)
	// 							}

	// 						-->
	// 						</script>
	// 						<div class="collapse" id="collapse-type-source-WebsiteURLsAggregate-json">
	// 							<a class="btn" onclick="copyJSONSourceWebsiteURLsAggregate()">Copy</a>
	// 							<div class="source-code" id="collapse-type-source-WebsiteURLsAggregate-json-inner">
	// {
	//    "Count": 0
	//    "Data":       [
	//          {
	//             "Website": "string",
	//             "FolderName": "string",
	//             "FolderCSSSelectorEnabled": 0,
	//             "FolderCSSSelector": "string"
	//             "PasswordSelector": "string",
	//             "Title": "string",
	//             "UsernameSelector": "string",
	//             "WebSiteURLID": 0,
	//             "Frequency": 0,
	//             "IsFlagged": 0,
	//             "LastCrawled": 0,
	//             "LastError": 0,
	//             "LoginURL": "string",
	//             "Notes": "string",
	//             "Password": "string",
	//             "WebsiteStatus": 0,
	//             "URL": "string",
	//             "Username": "string",
	//             "Versions": 0,
	//             "ImportedFromScrubberRunID": 0,
	//             "IsDeleted": 0,
	//             "LastChangeDetected": 0,
	//             "SeparateLoginFields": 0,
	//             "URLHash": "string",
	//             "AccountID": 0,
	//             "CSSSelector": "string",
	//             "IsPasswordProtected": 0,
	//             "ProtectedSelector": "string",
	//             "ClickOnAfterLogin": "string",
	//             "IsActive": 0,
	//             "LatestErrorCount": 0,
	//             "SeparateLoginFieldsNextButton": "string",
	//             "AssignedTo": 0,
	//             "CurrentScrubberRunID": 0,
	//             "SubmitSelector": "string",
	//             "WebsiteFolderID": 0,
	//             "DateImported": 0,
	//             "NextCrawl": 0,
	//             "IsCSSSelectorEnabled": 0,
	//             "PauseAfterLoginSeconds": 0,
	//             "RedirectAfterLogin": "string"
	//             "WebSiteID": 0

	//          }
	//       ]
	// }

	// 							</div>
	// 						</div>
	// 						<div class="collapse" id="collapse-type-source-WebsiteURLsAggregate-kotlin">
	// 							<a class="btn" onclick="copyKotlinSourceWebsiteURLsAggregate()">Copy</a>
	// 							<div class="source-code" id="collapse-type-source-WebsiteURLsAggregate-kotlin-inner">
	// struct WebsiteURLsAggregate : Decodable {
	//    var Count: Int64
	//    var Data:       [
	//          {
	//             var AccountID: Int64
	//             var AssignedTo: Int64
	//             var CSSSelector: String
	//             var ClickOnAfterLogin: String
	//             var CurrentScrubberRunID: Int64
	//             var DateImported: Int64
	//             var Frequency: Int
	//             var ImportedFromScrubberRunID: Int64
	//             var IsActive: Int
	//             var IsCSSSelectorEnabled: Int
	//             var IsDeleted: Int
	//             var IsFlagged: Int
	//             var IsPasswordProtected: Int
	//             var LastChangeDetected: Int64
	//             var LastCrawled: Int64
	//             var LastError: Int64
	//             var LatestErrorCount: Int64
	//             var LoginURL: String
	//             var NextCrawl: Int64
	//             var Notes: String
	//             var Password: String
	//             var PasswordSelector: String
	//             var PauseAfterLoginSeconds: Int64
	//             var ProtectedSelector: String
	//             var RedirectAfterLogin: String
	//             var SeparateLoginFields: Int
	//             var SeparateLoginFieldsNextButton: String
	//             var SubmitSelector: String
	//             var Title: String
	//             var URL: String
	//             var URLHash: String
	//             var Username: String
	//             var UsernameSelector: String
	//             var Versions: Int64
	//             var WebSiteID: Int64
	//             var WebSiteURLID: Int64
	//             var WebsiteFolderID: Int64
	//             var WebsiteStatus: Int

	//             var FolderCSSSelector: String
	//             var FolderCSSSelectorEnabled: Int
	//             var FolderName: String
	//             var Website: String
	//          }
	//       ]
	// }

	// 						</div>
	// 					</div>

	sb.WriteString(`
						

		<!-- /#main-right -->
		</div>

		<script type="application/javascript">

			function copyHashToClipboard(content) {
				content = location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '') + location.pathname + (location.search ? location.search : '') + content
				navigator.clipboard.writeText(content)
			}
	
			function copyContentToClipboard(id) { 
				var copyText = document.getElementById(id);
				navigator.clipboard.writeText(copyText.innerText); 
			}

			console.log('Running this!!'); 
			// window.scrollTo(0, 10000);
		</script>
	</body>

</html>` + "`" + `

}
`)

	ioutil.WriteFile(outFile, []byte(sb.String()), 0777)

}

func htmlHeader() string {
	return `package apidocs 

	func ShowDocs() string { 
	return ` + "`" + `<html>
	<head>
		<title>API Docs</title>
		<!-- CSS only -->
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">
		<!-- JavaScript Bundle with Popper -->
		<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-ka7Sk0Gln4gmtz2MlQnikT1wXgYsOg+OMhuP+IlRH9sENBO0LRn5q+8nbTov4+1p" crossorigin="anonymous"></script>
		<!-- Bootstrap Icons --> 
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.7.2/font/bootstrap-icons.css">
		<style type="text/css">
		<!--
		* { font-family: arial; }

		body { 
			height: 100%; 
			width: 100%; 
		}

		.field-type { 
			font-style: italic; 
			color: #666; 
		}

		table.field-table { 
			margin-bottom: 0px !important; 
		}

		#container { 
			display: flex; 
			position: relative; 
		}
		#main-left { 
			position: fixed; 
			top: 0; 
			left: 0; 
			width: 400px; 
			border-right: solid 1px #909090;
			height: 100%; 
			overflow-y: auto; 
			padding-left: 20px; 
			min-width: 300px; 
		}

		#main-left h2 { 
			font-size: 18px; 
			margin-bottom: 20px; 
		}

		#main-left h3 { 

		}

		#main-left ul { 
			padding-left: 0px; 
			list-style-type: none; 
		}

		#main-left ul li { 
			padding-left: 0px; 
			display: block; 
			margin-bottom: 10px; 
		}

		#main-left ul li a {
			font-size: 14px; 
			color: #000; 
			text-decoration: none; 
			display: block; 
			border-bottom: solid 1px #999;  
			margin-bottom: 10px; 
		}

		#main-left ul li ul li a { 
			padding-left: 20px; 
			border-bottom: none; 
			font-size: 13px; 
			color: #444; 
		}

		#main-left .menu-type { 

		}

		#main-right { 
			margin-left: 400px; 
			padding: 20px; 
		}

		.controller-container { 
			border: solid 1px #ccc; 
			border-radius: 5px; 
			margin-bottom: 20px; 
		}

		.controller-container h2 { 
			padding: 10px; 
			background-color: #eee; 
			border-bottom: solid 1px #ccc; 
			border-top-left-radius: 5px; 
			border-top-right-radius: 5px; 
			margin-bottom: 0px; 
		}

		.controller-container .inner { 
			padding: 20px; 
		}

		.endpoint { 
			padding: 20px; 
			border: solid 1px #909090; 
			border-radius: 3px; 
			margin-bottom: 20px; 
		}

		.endpoint h4 { 
			font-size: 16px; 
			font-weight: 700; 
		}

		.endpoint-method { 
			border-radius: 3px; 
			color: #fff; 
			font-size: 14px; 
			font-weight: 700; 
			min-width: 80px; 
			padding: 6px 0; 
			text-align: center; 
			text-shadow: 0 1px 0 rgb(0, 0, 0 / 10%); 
			display: inline-block; 
		}

		.endpoint-title { 
			color: #3b4151; 
			font-size: 16px; 
			padding: 0 10px; 
		}

		.endpoint-query-var { 
			color: #FF8300;
			font-weight: 500; 
		}

		.endpoint-description, .endpoint-parameters, .endpoint-body, .endpoint-response { 
			background-color: #fff; 
			color: #000; 
			border-radius: 3px; 
			margin-top: 10px; 
			margin-bottom: 10px; 
			padding: 20px; 
		}

		.source-code { 
			background-color: rgb(51,51,51); 
			white-space: pre-wrap; 
			word-break: break-all; 
			color: #fff; 
			font-size: 12px; 
			padding: 10px; 
		}

		.endpoint-response-footer, .endpoint-body-footer { 
			font-size: 14px; 
			margin-top: 10px; 
		}

		.endpoint-parameters, endpoint-body { 
			border: solid 1px #999; 
		}


		.endpoint-path { 
			font-size: 13px; 
			color: #000; 
			font-style: italic; 
			padding: 0 10px; 
			font-weight: normal; 
			display: block; 
			margin-left: 88px; 
		}

		.endpoint-post { 
			border-color: #49cc90; 
			background: rgba(73, 204, 144, .1); 
		}

		.endpoint-post-method { 
			background: #49cc90; 
		}

		.endpoint-put { 
			border-color: #fca130; 
			background: rgba(252, 161, 48, .1); 
		}

		.endpoint-put-method { 
			background: #fca130;
		}

		.endpoint-delete { 
			border-color: #f93e3e; 
			background: rgba(249, 62, 62, .1); 
		}

		.endpoint-delete-method { 
			background: #f93e3e;
		}

		.endpoint-get { 
			background: rgba(97, 175, 254, .1); 
			border-color: #61affe; 
		}

		.endpoint-get-method { 
			background: #61affe; 
		}

		.endpoint h3 { 
			font-weight: bold; 
			margin-bottom: 10px; 
			font-size: 20px; 
			display: flex; 
		}

		.endpoint-title-left { 
			flex-grow: 1; 
		}

		.endpoint-title-right { 
			text-align: right; 
			align-content: right; 
		}

		.endpoint-perm { 
			font-size: 14px; 
			padding: 6px; 
			border-radius: 3px; 
			background-color: #eee; 
			display: inline-block; 
			font-weight: normal; 
			border: solid 1px #000; 
		}

		.endpoint-perm a { 
			color: #444; 
			text-decoration: none; 
		}

		.endpoint-perm-anonymous { 
			background-color: #efe; 
			border-color: #9f9; 
		}

		.endpoint-perm-anyone { 
			background-color: #eef; 
			border-color: #99f; 
		}

		.endpoint-perm-perm { 
			background-color: #fee; 
			border-color: #f99; 
		}

		.endpoint-footer { 
			display: flex; 
			font-size: 10px; 
		}
		.endpoint-footer-left { 
			flex-grow: 1; 
			text-align: left; 
		}

		.endpoint-footer-right { 
			text-align: right; 
		}

		#types-container { 
			margin-bottom: 20px; 
		}


		.type { 
			border: solid 1px #ccc; 
			border-radius: 5px; 
			margin-bottom: 20px; 
		}

		.type h2 { 
			padding: 10px; 
			background-color: #eee; 
			border-bottom: solid 1px #ccc; 
			border-top-left-radius: 5px; 
			border-top-right-radius: 5px; 
			margin-bottom: 0px; 
			display: flex; 
		}

		.type h2 .type-title-left { 
			flex-grow: 1; 
		}

		.type h2 .type-title-right { 
			text-align: right; 
			font-size: 12px; 
		}

		.type table thead tr td { 

		}

		.type-footer { 
			display: flex; 
			font-size: 10px; 
			padding: 10px;
		}
		.type-footer-left { 
			flex-grow: 1; 
			text-align: left; 
		}

		.type-footer-right { 
			text-align: right; 
		}

		.type-used-in { 
			font-size: 10px; 
			padding: 10px; 
		}


		.datatype-string { 
			color: rgb(162, 252, 162);
		}

		.datatype-number { 
			color: rgb(211, 99, 99); 
		}

		#permissions-container { 
			border: solid 1px #ccc; 
			border-radius: 5px; 
			margin-bottom: 20px; 
		}

		.permission {}
		.permission-title {}
		.permission-description {}

		-->
		</style>
	</head>
	<body>
	<div id="main-left">
		<h2>API Docs</h2>`
}
