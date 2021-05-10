package testassets

import "github.com/macinnir/dvc/core/lib/schema"

var (
	columnFooID = &schema.Column{
		Name:       "FooID",
		IsNullable: false,
		IsUnsigned: true,
		DataType:   "int",
		MaxLength:  0,
		Precision:  10,
		Type:       "int(10) unsigned",
		ColumnKey:  "PRI",
		Extra:      "auto_increment",
	}

	columnDateCreated = &schema.Column{
		Name:         "DateCreated",
		Position:     5,
		Default:      "0",
		IsNullable:   false,
		IsUnsigned:   true,
		DataType:     "bigint",
		MaxLength:    0,
		Precision:    20,
		CharSet:      "",
		Type:         "bigint(20) unsigned",
		ColumnKey:    "",
		NumericScale: 0,
		Extra:        "",
		FmtType:      "%d",
		GoType:       "int64",
		IsString:     false,
	}

	columnIsDeleted = &schema.Column{
		Name:         "IsDeleted",
		Position:     3,
		Default:      "0",
		IsNullable:   false,
		IsUnsigned:   false,
		DataType:     "tinyint",
		MaxLength:    0,
		Precision:    3,
		CharSet:      "",
		Type:         "tinyint(4)",
		ColumnKey:    "",
		NumericScale: 0,
		Extra:        "",
		FmtType:      "%d",
		GoType:       "int",
		IsString:     false,
	}

	columnName = &schema.Column{
		Name:         "Name",
		Position:     4,
		Default:      "",
		IsNullable:   false,
		IsUnsigned:   false,
		DataType:     "varchar",
		MaxLength:    200,
		Precision:    3,
		CharSet:      "",
		Type:         "varchar(200)",
		ColumnKey:    "",
		NumericScale: 0,
		Extra:        "",
		FmtType:      "%s",
		GoType:       "string",
		IsString:     true,
	}
)
