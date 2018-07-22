package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDropTable(t *testing.T) {

	var e error
	var result string
	q := &Query{}
	table := &Table{Name: "Foo"}
	result, e = q.DropTable(table)

	assert.Nil(t, e)
	assert.Equal(t, "DROP TABLE `Foo`;", result)
}

func TestCreateColumn(t *testing.T) {
	var e error
	var result string
	q := &Query{}

	// int
	column := &Column{Name: "Foo", IsNullable: false, DataType: "int", Type: "int(10) unsigned", Extra: "auto_increment"}
	result, e = q.CreateColumn(column)

	assert.Nil(t, e)
	assert.Equal(t, "`Foo` int(10) unsigned NOT NULL auto_increment", result)

	// varchar
	column.DataType = "varchar"
	column.Type = "varchar(32)"
	column.Default = ""
	column.Extra = ""
	column.ColumnKey = ""

	result, e = q.CreateColumn(column)
	assert.Nil(t, e)
	assert.Equal(t, "`Foo` varchar(32) NOT NULL DEFAULT ''", result)

	// enum
	column.DataType = "enum"
	column.Type = "ENUM('Foo', 'Bar', 'Baz')"
	column.Default = "Foo"
	result, e = q.CreateColumn(column)
	assert.Nil(t, e)
	assert.Equal(t, "`Foo` ENUM('Foo', 'Bar', 'Baz') NOT NULL DEFAULT 'Foo'", result)

	//
}

func TestCreateColumnWithDefault(t *testing.T) {}
func TestCreateColumnWithExtra(t *testing.T) {

}
