package lib

// DataTypeToGoTypeString converts a database type to its equivalent golang datatype
func DataTypeToGoTypeString(column *Column) (fieldType string) {
	fieldType = "int64"
	switch column.DataType {
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

	if column.IsNullable == true {
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
