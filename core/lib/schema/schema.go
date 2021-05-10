package schema

import (
	"encoding/json"
	"io/ioutil"

	"github.com/macinnir/dvc/core/lib"
)

const (
	SchemaTypeMySQL      = "mysql"
	SchemaTypeSQLite     = "sqlite"
	SchemaTypePostgreSQL = "postgresql"
	SchemaTypeSQLServer  = "sqlserver"
)

// loadDatabase loads a database from configuration
func LoadLocalSchemas() (*SchemaList, error) {

	fileBytes, e := ioutil.ReadFile(lib.SchemasFilePath)

	if e != nil {
		return nil, e
	}

	schemaList := &SchemaList{
		Schemas: []*Schema{},
	}

	e = json.Unmarshal(fileBytes, schemaList)
	if e != nil {
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
func DataTypeToFormatString(column *Column) (fieldType string) {

	fieldType = "%s"

	switch column.DataType {
	case "int":
		fieldType = "%d"
	case "bigint":
		fieldType = "%d"
	case "tinyint":
		fieldType = "%d"
	case "varchar":
		fieldType = "%s"
	case "enum":
		fieldType = "%s"
	case "text":
		fieldType = "%s"
	case "date":
		fieldType = "%s"
	case "datetime":
		fieldType = "%s"
	case "char":
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
	case "int":
		fieldType = "int64"
	case "tinyint":
		fieldType = "int"
	case "varchar":
		fieldType = "string"
	case "enum":
		fieldType = "string"
	case "text":
		fieldType = "string"
	case "date":
		fieldType = "string"
	case "datetime":
		fieldType = "string"
	case "char":
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

func IsString(column *Column) bool {

	switch column.DataType {
	case "tinyint":
		return false
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
	case "char":
		return true
	case "decimal":
		return false
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