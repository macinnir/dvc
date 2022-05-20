package schema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/macinnir/dvc/core/lib"
)

const (
	SchemaTypeMySQL      = "mysql"
	SchemaTypePostgreSQL = "postgresql"
	SchemaTypeSQLServer  = "sqlserver"
)

// loadDatabase loads a database from configuration
func LoadLocalSchemas() (*SchemaList, error) {

	var e error
	var appSchema *SchemaList
	var coreSchema *SchemaList

	if appSchema, e = loadSchema(lib.SchemasFilePath); e != nil {
		return nil, e
	}

	if coreSchema, e = loadSchema(lib.CoreSchemasFilePath); e != nil {
		return nil, e
	}

	appSchema.Schemas = append(appSchema.Schemas, coreSchema.Schemas...)

	return appSchema, nil
}

func loadSchema(filePath string) (*SchemaList, error) {

	fmt.Println("Load schema from ", filePath)

	var e error
	var fileBytes []byte

	if fileBytes, e = ioutil.ReadFile(filePath); e != nil {
		return nil, e
	}

	schemaList := &SchemaList{
		Schemas: []*Schema{},
	}

	if e = json.Unmarshal(fileBytes, schemaList); e != nil {
		return nil, e
	}

	return schemaList, nil
}

// DataTypeToGoTypeString converts a database type to its equivalent golang datatype
func IsValidSQLType(str string) bool {
	switch str {
	case "int":
		return true
	case "varchar":
		return true
	case "enum":
		return true
	case "text":
		return true
	case "date":
		return true
	case "datetime":
		return true
	case "bigint":
		return true
	case "tinyint":
		return true
	case "char":
		return true
	case "decimal":
		return true
	}
	return false
}

// DataTypeToFormatString converts a database type to its equivalent golang datatype
func GoTypeFormatString(goType string) (fieldType string) {

	fieldType = "%s"

	switch goType {
	case "int", "int64":
		fieldType = "%d"
	case "string":
		fieldType = "%s"
	case "float", "float64":
		fieldType = "%f"
	}

	return
}

// DataTypeToFormatString converts a database type to its equivalent golang datatype
// TODO move to mysql specific
// TODO make column types constants in mysql
func DataTypeToFormatString(column *Column) (fieldType string) {

	fieldType = "%s"

	switch column.DataType {
	case "int", "bigint", "tinyint":
		fieldType = "%d"
	case "varchar", "enum", "text", "date", "datetime", "char":
		fieldType = "%s"
	case "decimal":
		fieldType = "%f"
	}

	return
}

// DataTypeToGoTypeString converts a database type to its equivalent golang datatype
func DataTypeToGoTypeString(column *Column) (fieldType string) {
	fieldType = "int64"

	switch column.DataType {
	case "int", "bigint":
		fieldType = "int64"
	case "tinyint":
		fieldType = "int"
	case "char", "varchar", "tinytext", "mediumtext", "text", "longtext", "enum", "set", "date", "datetime":
		fieldType = "string"
	case "decimal":
		fieldType = "float64"
	}

	if column.IsNullable {
		switch fieldType {
		case "string":
			// fieldType = "sql.NullString"
			fieldType = "null.String"
		case "int64":
			// fieldType = "sql.NullInt64"
			fieldType = "null.Int"
		case "float64":
			// fieldType = "sql.NullFloat64"
			fieldType = "null.Float"
		}
	}
	return
}

// TODO move to mysql specific
// TODO make column types constants in mysql
func DataTypeToTypescriptString(dbDataType string) (fieldType string) {

	fieldType = "number"

	switch dbDataType {
	case "char", "varchar", "tinytext", "mediumtext", "text", "longtext", "enum", "set", "datetime", "date", "time":
		fieldType = "string"
	}

	return
}

func ExtractBaseGoType(goDataType string) string {

	if len(goDataType) == 0 {
		return ""
	}

	for {

		if len(goDataType) == 0 {
			break
		}

		if goDataType[0:1] == "*" {
			goDataType = goDataType[1:]
			continue
		}

		if len(goDataType) > 4 && goDataType[0:4] == "map[" {
			goDataType = goDataType[strings.Index(goDataType, "]")+1:]
			continue
		}

		if len(goDataType) > 2 && goDataType[0:2] == "[]" {
			goDataType = goDataType[2:]
			continue
		}

		// if len(goDataType) > 7 && goDataType[0:7] == "models." {
		// 	goDataType = goDataType[7:]
		// 	continue
		// }

		break
	}

	return goDataType
}

func IsGoTypeBaseType(goDataType string) bool {
	switch goDataType {
	case "byte", "interface{}", "bytes.Buffer", "string", "null.String", "null.Float", "int", "int64", "float64", "bool":
		return true
	default:
		return false
	}
}

func GoBaseTypeToBaseTypescriptType(goDataType string) string {

	if strings.Contains(goDataType, ".") && goDataType != "null.String" && goDataType != "null.Float" {
		parts := strings.Split(goDataType, ".")
		goDataType = parts[len(parts)-1]
	}

	switch goDataType {
	case "string", "null.String":
		return "string"
	case "int", "int64", "float64", "null.Float":
		return "number"
	case "bool":
		return "boolean"
	default:
		if goDataType[0:1] == "*" {
			goDataType = goDataType[1:]
		}

		return goDataType
	}
}

func GoBaseTypeToBaseTypescriptDefault(goDataType string) string {
	switch goDataType {
	case "[]byte", "interface{}", "bytes.Buffer":
		return "null"
	case "string", "null.String":
		return "''"
	case "int", "int64", "float64", "null.Float":
		return "0"
	case "bool":
		return "false"
	default:
		return "null"
	}
}

func ParseMapTypeToTypescriptString(goDataType string) string {

	// Remove the map[ prefix
	goDataType = goDataType[4:]

	// key type
	keyType := goDataType[0:strings.Index(goDataType, "]")]

	// Remove the key
	goDataType = goDataType[len(keyType)+1:]

	return fmt.Sprintf("{ [key: %s]: %s }", GoBaseTypeToBaseTypescriptType(keyType), GoBaseTypeToBaseTypescriptType(goDataType))
}

func GoTypeToTypescriptString(goDataType string) string {

	if len(goDataType) > 4 && goDataType[0:4] == "map[" {
		return ParseMapTypeToTypescriptString(goDataType)
	}

	if len(goDataType) > 6 && goDataType[0:6] == "[]map[" {
		return ParseMapTypeToTypescriptString(goDataType[2:]) + "[]"
	}

	// DTOs
	if len(goDataType) > 5 && goDataType[0:5] == "dtos." {
		return goDataType[5:]
	}

	if len(goDataType) > 6 && goDataType[0:6] == "*dtos." {
		return goDataType[6:]
	}

	if len(goDataType) > 7 && goDataType[0:7] == "[]dtos." {
		return goDataType[7:] + "[]"
	}

	if len(goDataType) > 8 && goDataType[0:8] == "[]*dtos." {
		return goDataType[8:] + "[]"
	}

	// APPDTOs
	if len(goDataType) > 8 && goDataType[0:8] == "appdtos." {
		return goDataType[8:]
	}

	if len(goDataType) > 9 && goDataType[0:9] == "*appdtos." {
		return goDataType[9:]
	}

	if len(goDataType) > 10 && goDataType[0:10] == "[]appdtos." {
		return goDataType[10:] + "[]"
	}

	if len(goDataType) > 11 && goDataType[0:11] == "[]*appdtos." {
		return goDataType[11:] + "[]"
	}

	// Models
	if len(goDataType) > 7 && goDataType[0:7] == "models." {
		return goDataType[7:]
	}

	if len(goDataType) > 8 && goDataType[0:8] == "*models." {
		return goDataType[8:]
	}

	if len(goDataType) > 9 && goDataType[0:9] == "[]models." {
		return goDataType[9:] + "[]"
	}

	if len(goDataType) > 10 && goDataType[0:10] == "[]*models." {
		return goDataType[10:] + "[]"
	}

	// Aggregates
	if len(goDataType) > 11 && goDataType[0:11] == "aggregates." {
		return goDataType[11:]
	}

	if len(goDataType) > 12 && goDataType[0:12] == "*aggregates." {
		return goDataType[12:]
	}

	if len(goDataType) > 13 && goDataType[0:13] == "[]aggregates." {
		return goDataType[13:] + "[]"
	}

	if len(goDataType) > 14 && goDataType[0:14] == "[]*aggregates." {
		return goDataType[14:] + "[]"
	}

	if (len(goDataType) >= 13 && goDataType[0:13] == "*bytes.Buffer") ||
		(len(goDataType) >= 12 && goDataType[0:12] == "bytes.Buffer") ||
		(len(goDataType) >= 11 && goDataType[0:11] == "interface{}") ||
		(len(goDataType) >= 6 && goDataType[0:6] == "[]byte") {
		return "any"
	}

	if len(goDataType) > 4 && goDataType[:4] == "[][]" {
		return GoBaseTypeToBaseTypescriptType(goDataType[4:]) + "[][]"
	}

	if len(goDataType) > 2 && goDataType[0:2] == "[]" {
		return GoBaseTypeToBaseTypescriptType(goDataType[2:]) + "[]"
	}

	return GoBaseTypeToBaseTypescriptType(goDataType)
}

func GoTypeToTypescriptDefault(goDataType string) (fieldType string) {

	if len(goDataType) > 3 && goDataType[0:4] == "map[" {

		if goDataType[0:11] == "map[string]" {
			return "{}"
		}

		return "[]"
	}

	if len(goDataType) >= 6 && goDataType[0:6] == "[]byte" {
		return "null"
	}

	if len(goDataType) > 2 && goDataType[0:2] == "[]" {
		return "[]"
	}

	baseType := ExtractBaseGoType(goDataType)

	if len(baseType) >= 13 && baseType[0:13] == "*bytes.Buffer" {
		return "null"
	}

	if len(baseType) >= 12 && baseType[0:12] == "bytes.Buffer" {
		return "null"
	}

	if len(baseType) > 7 && baseType[0:7] == "models." {
		return "new" + baseType[7:] + "()"
	}

	if len(baseType) > 5 && baseType[0:5] == "dtos." {
		return "new" + baseType[5:] + "()"
	}

	if len(baseType) > 6 && baseType[0:6] == "*dtos." {
		return "new" + baseType[6:] + "()"
	}

	if len(baseType) > 8 && baseType[0:8] == "appdtos." {
		return "new" + baseType[8:] + "()"
	}

	if len(baseType) > 9 && baseType[0:9] == "*appdtos." {
		return "new" + baseType[9:] + "()"
	}

	if !IsGoTypeBaseType(baseType) {
		return "new" + baseType + "()"
	}

	return GoBaseTypeToBaseTypescriptDefault(goDataType)
}

// TODO move to mysql specific
// TODO make column types constants in mysql
func DataTypeToTypescriptDefault(dataType string) (fieldType string) {

	fieldType = "0"

	switch dataType {
	case "char", "varchar", "tinytext", "mediumtext", "text", "longtext", "enum", "set", "datetime", "date", "time":
		fieldType = "''"
	}

	return
}

// TODO move to mysql specific
// TODO make column types constants in mysql
func IsString(column *Column) bool {

	switch column.DataType {
	case "varchar", "enum", "text", "longtext", "tinytext", "date", "datetime", "char":
		return true
	default:
		return false
	}

}

// IsNull returns true if the field uses the `null` package
func IsNull(fieldType string) bool {
	return len(fieldType) > 5 && fieldType[0:5] == "null."
}

type SchemaComparison struct {
	Database    string
	DatabaseKey string
	Additions   int
	Alterations int
	Deletions   int
	Changes     []*SchemaChange
}

type SchemaChange struct {
	Type          string
	SQL           string
	IsDestructive bool
}
