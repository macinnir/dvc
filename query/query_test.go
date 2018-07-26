package query

import (
	"github.com/macinnir/dvc/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDropTable(t *testing.T) {

	var e error
	var result string
	q := &Query{}
	table := &types.Table{Name: "Foo"}
	result, e = q.DropTable(table)

	assert.Nil(t, e)
	assert.Equal(t, "DROP TABLE `Foo`;", result)
}

func TestCreateColumn(t *testing.T) {
	var e error
	var result string
	q := &Query{}

	// int
	column := &types.Column{Name: "Foo", IsNullable: false, DataType: "int", Type: "int(10) unsigned", Extra: "auto_increment"}
	result, e = q.CreateColumn(column)
	assert.Nil(t, e)
	assert.Equal(t, "`Foo` int(10) unsigned NOT NULL auto_increment", result)

	// int w/ default
	column.Extra = ""
	column.Default = "0"
	result, e = q.CreateColumn(column)
	assert.Nil(t, e)
	assert.Equal(t, "`Foo` int(10) unsigned NOT NULL DEFAULT 0", result)

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
}

func TestAlterTableDropColumn(t *testing.T) {
	var e error
	var result string

	table := &types.Table{Name: "Foo"}
	column := &types.Column{Name: "Bar"}

	q := &Query{}
	result, e = q.AlterTableDropColumn(table, column)

	assert.Nil(t, e)
	assert.Equal(t, "ALTER TABLE `Foo` DROP COLUMN `Bar`;", result)
}

func TestAlterTableCreateColumn(t *testing.T) {
	var e error
	var result string

	table := &types.Table{Name: "Foo"}
	column := &types.Column{Name: "Bar", IsNullable: false, DataType: "int", Type: "int(10) unsigned", Extra: "auto_increment"}
	q := &Query{}
	result, e = q.AlterTableCreateColumn(table, column)

	assert.Nil(t, e)
	assert.Equal(t, "ALTER TABLE `Foo` ADD COLUMN `Bar` int(10) unsigned NOT NULL auto_increment;", result)
}

func TestAddIndex(t *testing.T) {
	var e error
	var result string

	table := &types.Table{Name: "Foo"}
	column := &types.Column{Name: "Bar"}

	q := &Query{}

	result, e = q.AddIndex(table, column)

	assert.Nil(t, e)
	assert.Equal(t, "ALTER TABLE `Foo` ADD INDEX `i_Bar` (`Bar`);", result)
}

func TestAddUniqueIndex(t *testing.T) {
	var e error
	var result string

	table := &types.Table{Name: "Foo"}
	column := &types.Column{Name: "Bar"}

	q := &Query{}

	result, e = q.AddUniqueIndex(table, column)

	assert.Nil(t, e)
	assert.Equal(t, "ALTER TABLE `Foo` ADD UNIQUE INDEX `ui_Bar` (`Bar`);", result)
}

func TestDropIndex(t *testing.T) {

	var e error
	var result string

	table := &types.Table{Name: "Foo"}
	column := &types.Column{Name: "Bar"}

	q := &Query{}

	result, e = q.DropIndex(table, column)

	assert.Nil(t, e)
	assert.Equal(t, "ALTER TABLE `Foo` DROP INDEX `i_Bar`;", result)
}

func TestDropUniqueIndex(t *testing.T) {

	var e error
	var result string

	table := &types.Table{Name: "Foo"}
	column := &types.Column{Name: "Bar"}

	q := &Query{}

	result, e = q.DropUniqueIndex(table, column)

	assert.Nil(t, e)
	assert.Equal(t, "ALTER TABLE `Foo` DROP INDEX `ui_Bar`;", result)
}
