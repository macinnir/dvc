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

	columnFooIDBigIntNoAI = &schema.Column{
		Name:       "FooID",
		IsNullable: false,
		IsUnsigned: true,
		DataType:   "bigint",
		MaxLength:  0,
		Precision:  10,
		Type:       "bigint(20) unsigned",
		ColumnKey:  "PRI",
		Extra:      "",
	}

	columnDateCreated = &schema.Column{
		Name:         "DateCreated",
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

	columnNameSize100 = &schema.Column{
		Name:         "Name",
		Default:      "",
		IsNullable:   false,
		IsUnsigned:   false,
		DataType:     "varchar",
		MaxLength:    100,
		Precision:    3,
		CharSet:      "",
		Type:         "varchar(100)",
		ColumnKey:    "",
		NumericScale: 0,
		Extra:        "",
		FmtType:      "%s",
		GoType:       "string",
		IsString:     true,
	}

	columnNameWithUniqueIndex = &schema.Column{
		Name:         "Name",
		Default:      "",
		IsNullable:   false,
		IsUnsigned:   false,
		DataType:     "varchar",
		MaxLength:    200,
		Precision:    3,
		CharSet:      "",
		Type:         "varchar(200)",
		ColumnKey:    "UNI",
		NumericScale: 0,
		Extra:        "",
		FmtType:      "%s",
		GoType:       "string",
		IsString:     true,
	}

	columnNameWithIndex = &schema.Column{
		Name:         "Name",
		Default:      "",
		IsNullable:   false,
		IsUnsigned:   false,
		DataType:     "varchar",
		MaxLength:    200,
		Precision:    3,
		CharSet:      "",
		Type:         "varchar(200)",
		ColumnKey:    "MUL",
		NumericScale: 0,
		Extra:        "",
		FmtType:      "%s",
		GoType:       "string",
		IsString:     true,
	}
)
