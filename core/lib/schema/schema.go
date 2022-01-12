package schema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

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
	case "char", "varchar", "tinytext", "mediumtext", "text", "longtext", "enum", "set":
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

func GoTypeToTypescriptString(goDataType string) (fieldType string) {

	fieldType = "number"

	isArray := len(goDataType) > 0 && goDataType[0:1] == "["

	// if isArray {
	// 	fmt.Println("goDataType: ", goDataType[1:])
	// }

	switch goDataType {
	case "string", "null.String":
		fieldType = "string"
	case "bool":
		fieldType = "boolean"
	}

	if isArray {
		fieldType += "[]"
	}

	return
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
func GoTypeToTypescriptDefault(goDataType string) (fieldType string) {

	fieldType = "0"

	if len(goDataType) > 0 && goDataType[0:1] == "[" {
		return "[]"
	}

	switch goDataType {
	case "string", "null.String":
		fieldType = "''"
	case "bool":
		fieldType = "false"
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
