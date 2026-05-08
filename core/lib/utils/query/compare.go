package query

import (
	"fmt"
)

// #region Equals

// EQ is an equals statement between a table column and a value
func EQ(fieldName Column, value interface{}) *WherePart {
	return newWherePart(
		WhereTypeEquals,
		string(fieldName),
		[]interface{}{
			value,
		},
	)
}

// EQInt is an equals statement between a table column and an int value
func EQInt(fieldName Column, value int) *WherePart {
	return EQ(fieldName, value)
}

// EQInt64 is an equals statement between a table column and an int64 value
func EQInt64(fieldName Column, value int64) *WherePart {
	return EQ(fieldName, value)
}

// EQString is an equals statement between a table column and a string value
func EQString(fieldName Column, value string) *WherePart {
	return EQ(fieldName, value)
}

// EQFloat is an equals statement between a table column and a float value
func EQFloat(fieldName Column, value float64) *WherePart {
	return EQ(fieldName, value)
}

// EQF allows for one column to be equal to another
// Example for a subselect
//
// query.Select(&models.UserGroupUser{}).Alias("ugu").FieldRaw("1", "n").Where(
//
//	query.EQF("UserID", "`u`.`UserID`"),
//	query.And(),
//	query.EQ("UserGroupID", groupID),
//	query.And(),
//	query.EQ("IsDeleted", 0),
//
// ),
func EQF(fieldName1, fieldName2 string) *WherePart {
	return newWherePart(
		WhereTypeEqualsField,
		fieldName1,
		[]interface{}{fieldName2},
	)
}

// #endregion

// #region Not Equals

// NE is a not equals statement between a table column and a value
func NE(fieldName Column, value interface{}) *WherePart {
	return newWherePart(
		WhereTypeNotEquals,
		string(fieldName),
		[]interface{}{value},
	)
}

// NEInt is a not equals statement between a table column and an int value
func NEInt(fieldName Column, value int) *WherePart {
	return NE(fieldName, value)
}

// NEInt64 is a not equals statement between a table column and an int64 value
func NEInt64(fieldName Column, value int64) *WherePart {
	return NE(fieldName, value)
}

// NEString is a not equals statement between a table column and a string value
func NEString(fieldName Column, value string) *WherePart {
	return NE(fieldName, value)
}

// NEFloat is a not equals statement between a table column and a float value
func NEFloat(fieldName Column, value float64) *WherePart {
	return NE(fieldName, value)
}

// NEF allows for one column to be not equal to another
// Example for a subselect
//
// query.Select(&models.UserGroupUser{}).Alias("ugu").FieldRaw("1", "n").Where(
//
//	query.NEF("UserID", "`u`.`UserID`"),
//	query.And(),
//	query.EQ("UserGroupID", groupID),
//	query.And(),
//	query.EQ("IsDeleted", 0),
//
// ),
func NEF(fieldName1, fieldName2 string) *WherePart {
	return newWherePart(
		WhereTypeNotEqualsField,
		fieldName1,
		[]interface{}{fieldName2},
	)
}

// #endregion

// #region Less Than

// LT is a less than statement between a table column and a value
// LT('foo', 1) => WHERE `t`.`foo` < 1
func LT(fieldName Column, value interface{}) *WherePart {
	return newWherePart(
		WhereTypeLessThan,
		string(fieldName),
		[]interface{}{value},
	)
}

func LTInt(fieldName Column, value int) *WherePart {
	return LT(fieldName, value)
}

func LTInt64(fieldName Column, value int64) *WherePart {
	return LT(fieldName, value)
}

func LTFloat(fieldName Column, value float64) *WherePart {
	return LT(fieldName, value)
}

// #endregion

// #region Greater Than

// GT is a greater than statement between a table column and a value
func GT(fieldName Column, value interface{}) *WherePart {
	return newWherePart(
		WhereTypeGreaterThan,
		string(fieldName),
		[]interface{}{value},
	)
}

func GTInt(fieldName Column, value int) *WherePart {
	return GT(fieldName, value)
}

func GTInt64(fieldName Column, value int64) *WherePart {
	return GT(fieldName, value)
}

func GTFloat(fieldName Column, value float64) *WherePart {
	return GT(fieldName, value)
}

// #endregion

// #region Less Than or Equal To

// LTOE is a less than or equals (<=) statement between a table column and a value
//
//	`t`.`Col` <= value
func LTOE(fieldName Column, value interface{}) *WherePart {
	return newWherePart(
		WhereTypeLessThanOrEqualTo,
		string(fieldName),
		[]interface{}{value},
	)
}

func LTOEInt(fieldName Column, value int) *WherePart {
	return LTOE(fieldName, value)
}

func LTOEInt64(fieldName Column, value int64) *WherePart {
	return LTOE(fieldName, value)
}

func LTOEFloat(fieldName Column, value float64) *WherePart {
	return LTOE(fieldName, value)
}

// #endregion

// #region Greater Than or Equal To

// GTOE is a greater than or equals statement (>=) between a table column and a value
//
//	`t`.`Col` >= value
func GTOE(fieldName Column, value interface{}) *WherePart {
	return newWherePart(
		WhereTypeGreaterThanOrEqualTo,
		string(fieldName),
		[]interface{}{value},
	)
}

func GTOEInt(fieldName Column, value int) *WherePart {
	return GTOE(fieldName, value)
}

func GTOEInt64(fieldName Column, value int64) *WherePart {
	return GTOE(fieldName, value)
}

func GTOEFloat(fieldName Column, value float64) *WherePart {
	return GTOE(fieldName, value)
}

// #endregion

// #region Modulo

// Mod is applies modulo operation on column and value testing if it equals remainder
//
//	MOD(`t`.`Field`, value) = remainder
func Mod(fieldName Column, value, remainder int64) *WherePart {
	return newWherePart(
		WhereTypeMod,
		string(fieldName),
		[]interface{}{value, remainder},
	)
}

// Modf MOD(value, `t`.`Field`) = remainder
// Example: query.Mod("foo", 2, 1) -> `t`.`Foo` % 2 = 1
func Modf(value int64, fieldName Column, remainder int64) *WherePart {
	return newWherePart(
		WhereTypeModF,
		string(fieldName),
		[]interface{}{value, remainder},
	)
}

// #endregion

// BitAnd `t`.`Field` & a = b
// Example: query.BitAnd("foo", 2, 1) -> `t`.`Foo` & 2 = 1
func BitAnd(fieldName Column, a, b int64) *WherePart {
	return newWherePart(
		WhereTypeBitAnd,
		string(fieldName),
		[]interface{}{a, b},
	)
}

// #region IN

// IN is an IN clause
// Example: query.IN("col1", "foo", "bar", "baz")
func IN(fieldName Column, values ...interface{}) *WherePart {
	return newWherePart(
		WhereTypeIN,
		string(fieldName),
		values,
	)
}

// INInt is a helper function for converting a slice of string arguments into
// a slice of interface arguments, passed into an IN clause and returned
func INInt(fieldName Column, values ...int) *WherePart {
	interfaces := make([]interface{}, len(values))

	for k := range values {
		interfaces[k] = values[k]
	}

	return IN(fieldName, interfaces...)
}

// INInt64 is a helper function for converting a slice of string arguments into
// a slice of interface arguments, passed into an IN clause and returned
func INInt64(fieldName Column, values ...int64) *WherePart {
	interfaces := make([]interface{}, len(values))

	for k := range values {
		interfaces[k] = values[k]
	}

	return IN(fieldName, interfaces...)
}

// INString is a helper function for converting a slice of string arguments into
// a slice of interface arguments, passed into an IN clause and returned
func INString(fieldName Column, values ...string) *WherePart {
	interfaces := make([]interface{}, len(values))

	for k := range values {
		interfaces[k] = values[k]
	}

	return IN(fieldName, interfaces...)
}

// #endregion

// #region NOT IN

// NOTIN is an NOT IN clause
// Example: query.NOTIN("col1", "foo", "bar", "baz")
// Example: queyr.NOTIN("col2", 1, 2, 3)
func NOTIN(fieldName Column, values ...interface{}) *WherePart {
	return newWherePart(
		WhereTypeNotIN,
		string(fieldName),
		values,
	)
}

// NOTINInt is a helper function for converting a slice of string arguments into
// a slice of interface arguments, passed into a NOT IN clause and returned
func NOTINInt(fieldName Column, values ...int) *WherePart {
	interfaces := make([]interface{}, len(values))

	for k := range values {
		interfaces[k] = values[k]
	}

	return NOTIN(fieldName, interfaces...)
}

// NOTINInt64 is a helper function for converting a slice of string arguments into
// a slice of interface arguments, passed into a NOT IN clause and returned
func NOTINInt64(fieldName Column, values ...int64) *WherePart {
	interfaces := make([]interface{}, len(values))

	for k := range values {
		interfaces[k] = values[k]
	}

	return NOTIN(fieldName, interfaces...)
}

// NOTINFloat is a helper function for converting a slice of string arguments into
// a slice of interface arguments, passed into a NOT IN clause and returned
func NOTINFloat(fieldName Column, values ...float64) *WherePart {
	interfaces := make([]interface{}, len(values))

	for k := range values {
		interfaces[k] = values[k]
	}

	return NOTIN(fieldName, interfaces...)
}

// NOTINString is a helper function for converting a slice of string arguments into
// a slice of interface arguments, passed into a NOT IN clause and returned
func NOTINString(fieldName Column, values ...string) *WherePart {
	interfaces := make([]interface{}, len(values))

	for k := range values {
		interfaces[k] = values[k]
	}

	return NOTIN(fieldName, interfaces...)
}

// #endregion

// Rawf is a raw SQL statement
// Example: query.Rawf("`t`.`LastRunDate` + 60000 < %d", seconds)),
func Rawf(str string, args ...interface{}) *WherePart {
	return newWherePart(
		WhereTypeRaw,
		"",
		[]interface{}{fmt.Sprintf(str, args...)},
	)
}

// #region Between

// Between is a BETWEEN statement
// Example: Between("")
func Between(fieldName Column, from, to interface{}) *WherePart {
	return newWherePart(
		WhereTypeBetween,
		string(fieldName),
		[]interface{}{from, to},
	)
}

func BetweenInts(fieldName Column, from, to int) *WherePart {
	return Between(fieldName, from, to)
}

func BetweenInt64s(fieldName Column, from, to int64) *WherePart {
	return Between(fieldName, from, to)
}

func BetweenFloats(fieldName Column, from, to float64) *WherePart {
	return Between(fieldName, from, to)
}

func BetweenStrings(fieldName Column, from, to string) *WherePart {
	return Between(fieldName, from, to)
}

// #endregion

// #region Like

func Like(fieldName Column, value string) *WherePart {
	return newWherePart(
		WhereTypeLike,
		string(fieldName),
		[]interface{}{value},
	)
}

// #endregion Like

// #region Not Like
func NotLike(fieldName Column, value string) *WherePart {
	return newWherePart(
		WhereTypeNotLike,
		string(fieldName),
		[]interface{}{value},
	)
}

// #endregion Not Like

// #region Is Null

func IsNull(fieldName Column) *WherePart {
	return newWherePart(
		WhereTypeIsNull,
		string(fieldName),
		[]interface{}{},
	)
}

// #endregion Is Null

// #region Is Not Null

func IsNotNull(fieldName Column) *WherePart {
	return newWherePart(
		WhereTypeIsNotNull,
		string(fieldName),
		[]interface{}{},
	)
}

// #endregion Is Not Null
