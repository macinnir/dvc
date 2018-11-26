package mysql

import (
	"github.com/macinnir/dvc/lib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDropIndex(t *testing.T) {

	table := new(lib.Table)
	table.Name = "foo"

	column := new(lib.Column)
	column.Name = "bar"

	result, e := dropIndex(table, column)
	assert.Nil(t, e)
	assert.Equal(t, "ALTER TABLE `foo` DROP INDEX `i_foo_bar`;", result)
}

func TestDropUniqueIndex(t *testing.T) {

	table := new(lib.Table)
	table.Name = "foo"

	column := new(lib.Column)
	column.Name = "bar"

	result, e := dropUniqueIndex(table, column)
	assert.Nil(t, e)
	assert.Equal(t, "ALTER TABLE `foo` DROP INDEX `ui_foo_bar`;", result)
}

func TestHasDefaultString(t *testing.T) {
	assert.True(t, hasDefaultString("varchar"))
	assert.True(t, hasDefaultString("char"))
	assert.True(t, hasDefaultString("enum"))
	assert.False(t, hasDefaultString("int"))
}

func TestIsString(t *testing.T) {
	assert.True(t, isString("varchar"))
	assert.True(t, isString("char"))
	assert.True(t, isString("ENUM"))
	assert.False(t, isString("int"))
}

func TestIsInt(t *testing.T) {
	assert.True(t, isInt("tinyint"))
	assert.True(t, isInt("smallint"))
	assert.True(t, isInt("mediumint"))
	assert.True(t, isInt("int"))
	assert.True(t, isInt("bigint"))
	assert.False(t, isInt("char"))
}

func TestIsFixedPointType(t *testing.T) {
	assert.True(t, isFixedPointType("decimal"))
	assert.True(t, isFixedPointType("numeric"))
	assert.False(t, isFixedPointType("char"))
}

func TestIsFloatingPointType(t *testing.T) {
	assert.True(t, isFloatingPointType("float"))
	assert.True(t, isFloatingPointType("double"))
	assert.False(t, isFloatingPointType("char"))
}

func TestCreateColumnEnum(t *testing.T) {
	column := new(lib.Column)
	column.Name = "foo"
	column.DataType = "enum"
	column.Type = "ENUM('bar','baz')"
	column.Default = "bar"
	column.IsNullable = false
	result, e := createColumnSegment(column)
	assert.Nil(t, e)
	assert.Equal(t, "`foo` ENUM('bar','baz') NOT NULL DEFAULT 'bar'", result)
}

// func TestCreateColumn(t *testing.T) {
// 	m := new(MySQL)
// }
