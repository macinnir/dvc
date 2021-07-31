package mysql

import (
	"testing"

	"github.com/macinnir/dvc/core/connectors/mysql/testassets"
	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDropIndex(t *testing.T) {

	table := new(schema.Table)
	table.Name = "foo"

	column := new(schema.Column)
	column.Name = "bar"

	result := dropIndex(table, column)
	assert.Equal(t, result.Type, schema.DropIndex)
	assert.Equal(t, "ALTER TABLE `foo` DROP INDEX `i_foo_bar`;", result.SQL)
}

func TestDropUniqueIndex(t *testing.T) {

	table := new(schema.Table)
	table.Name = "foo"

	column := new(schema.Column)
	column.Name = "bar"

	result := dropUniqueIndex(table, column)
	assert.Equal(
		t,
		"ALTER TABLE `foo` DROP INDEX `ui_foo_bar`;",
		result.SQL)
	assert.Equal(t, schema.DropIndex, result.Type)
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
	column := new(schema.Column)
	column.Name = "foo"
	column.DataType = "enum"
	column.Type = "ENUM('bar','baz')"
	column.Default = "bar"
	column.IsNullable = false
	result := createColumnSegment(column)
	assert.Equal(t, "`foo` ENUM('bar','baz') NOT NULL DEFAULT 'bar'", result)
}

func TestCreateChangeSQL_NoChanges(t *testing.T) {

	s := NewMySQL(&lib.ConfigDatabase{})

	tables := testassets.TablesNoChange()

	comparison := s.CreateChangeSQL(tables[0], tables[1], "Foo")

	assert.Equal(t, 0, comparison.Additions)
	assert.Equal(t, 0, comparison.Deletions)
	assert.Equal(t, 0, comparison.Alterations)
	assert.Equal(t, 0, len(comparison.Changes))

	if len(comparison.Changes) > 0 {
		for k := range comparison.Changes {
			t.Log(comparison.Changes[k].SQL)
		}
	}

}

func TestCreateChangeSQL_DropColumn(t *testing.T) {

	s := NewMySQL(&lib.ConfigDatabase{})

	tables := testassets.TablesDropColumn()

	comparison := s.CreateChangeSQL(tables[0], tables[1], "Foo")

	assert.Equal(t, 0, comparison.Additions)
	assert.Equal(t, 1, comparison.Deletions)
	assert.Equal(t, 0, comparison.Alterations)
	require.Equal(t, 1, len(comparison.Changes))

	assert.Equal(t, schema.DropColumn, comparison.Changes[0].Type)
	assert.Equal(t, "ALTER TABLE `Foo` DROP COLUMN `Name`;", comparison.Changes[0].SQL)

}

func TestCreateChangeSQL_ChangeVarcharColumnSize(t *testing.T) {

	s := NewMySQL(&lib.ConfigDatabase{})

	tables := testassets.TablesChangeVarcharColumnSize()

	comparison := s.CreateChangeSQL(tables[0], tables[1], "Foo")

	assert.Equal(t, 0, comparison.Additions)
	assert.Equal(t, 0, comparison.Deletions)
	assert.Equal(t, 1, comparison.Alterations)
	require.Equal(t, 1, len(comparison.Changes))

	assert.Equal(t, schema.ChangeColumn, comparison.Changes[0].Type)
	assert.Equal(t, "ALTER TABLE `Foo` CHANGE `Name` `Name` varchar(200) NOT NULL DEFAULT '';", comparison.Changes[0].SQL)

}

func TestCreateChangeSQL_AddColumn(t *testing.T) {

	s := NewMySQL(&lib.ConfigDatabase{})

	tables := testassets.TablesAddColumn()

	comparison := s.CreateChangeSQL(tables[0], tables[1], "Foo")

	assert.Equal(t, 1, comparison.Additions)
	assert.Equal(t, 0, comparison.Deletions)
	assert.Equal(t, 0, comparison.Alterations)
	require.Equal(t, 1, len(comparison.Changes))

	assert.Equal(t, schema.AddColumn, comparison.Changes[0].Type)
	assert.Equal(t, "ALTER TABLE `Foo` ADD COLUMN `Name` varchar(200) NOT NULL DEFAULT '';", comparison.Changes[0].SQL)
}

func TestCreateChangeSQL_AddColumnWithUniqueIndex(t *testing.T) {

	s := NewMySQL(&lib.ConfigDatabase{})

	tables := testassets.TablesAddColumnWithUniqueIndex()

	comparison := s.CreateChangeSQL(tables[0], tables[1], "Foo")

	assert.Equal(t, 2, comparison.Additions)
	assert.Equal(t, 0, comparison.Deletions)
	assert.Equal(t, 0, comparison.Alterations)
	require.Equal(t, 2, len(comparison.Changes))

	assert.Equal(t, schema.AddColumn, comparison.Changes[0].Type)
	assert.Equal(t, "ALTER TABLE `Foo` ADD COLUMN `Name` varchar(200) NOT NULL DEFAULT '';", comparison.Changes[0].SQL)
	assert.Equal(t, schema.AddIndex, comparison.Changes[1].Type)
	assert.Equal(t, "ALTER TABLE `Foo` ADD UNIQUE INDEX `ui_Foo_Name` (`Name`);", comparison.Changes[1].SQL)
}

func TestCreateChangeSQL_DropColumnWithUniqueIndex(t *testing.T) {

	s := NewMySQL(&lib.ConfigDatabase{})

	tables := testassets.TablesAddColumnWithUniqueIndex()

	comparison := s.CreateChangeSQL(tables[1], tables[0], "Foo")

	assert.Equal(t, 0, comparison.Additions)
	assert.Equal(t, 2, comparison.Deletions)
	assert.Equal(t, 0, comparison.Alterations)
	require.Equal(t, 2, len(comparison.Changes))

	assert.Equal(t, schema.DropIndex, comparison.Changes[0].Type)
	assert.Equal(t, "ALTER TABLE `Foo` DROP INDEX `ui_Foo_Name`;", comparison.Changes[0].SQL)
	assert.Equal(t, schema.DropColumn, comparison.Changes[1].Type)
	assert.Equal(t, "ALTER TABLE `Foo` DROP COLUMN `Name`;", comparison.Changes[1].SQL)
}

func TestCreateChangeSQL_AddColumnWithIndex(t *testing.T) {

	s := NewMySQL(&lib.ConfigDatabase{})

	tables := testassets.TablesAddColumnWithIndex()

	comparison := s.CreateChangeSQL(tables[0], tables[1], "Foo")

	assert.Equal(t, 2, comparison.Additions)
	assert.Equal(t, 0, comparison.Deletions)
	assert.Equal(t, 0, comparison.Alterations)
	require.Equal(t, 2, len(comparison.Changes))

	assert.Equal(t, schema.AddColumn, comparison.Changes[0].Type)
	assert.Equal(t, "ALTER TABLE `Foo` ADD COLUMN `Name` varchar(200) NOT NULL DEFAULT '';", comparison.Changes[0].SQL)
	assert.Equal(t, schema.AddIndex, comparison.Changes[1].Type)
	assert.Equal(t, "ALTER TABLE `Foo` ADD INDEX `i_Foo_Name` (`Name`);", comparison.Changes[1].SQL)
}

func TestCreateChangeSQL_DropColumnWithIndex(t *testing.T) {

	s := NewMySQL(&lib.ConfigDatabase{})

	tables := testassets.TablesAddColumnWithIndex()

	comparison := s.CreateChangeSQL(tables[1], tables[0], "Foo")

	assert.Equal(t, 0, comparison.Additions)
	assert.Equal(t, 2, comparison.Deletions)
	assert.Equal(t, 0, comparison.Alterations)
	require.Equal(t, 2, len(comparison.Changes))

	assert.Equal(t, schema.DropIndex, comparison.Changes[0].Type)
	assert.Equal(t, "ALTER TABLE `Foo` DROP INDEX `i_Foo_Name`;", comparison.Changes[0].SQL)
	assert.Equal(t, schema.DropColumn, comparison.Changes[1].Type)
	assert.Equal(t, "ALTER TABLE `Foo` DROP COLUMN `Name`;", comparison.Changes[1].SQL)
}

func TestCreateChangeSQL_DropAutoIncrement(t *testing.T) {

	s := NewMySQL(&lib.ConfigDatabase{})

	tables := testassets.TablesDropAutoIncrement()

	comparison := s.CreateChangeSQL(tables[0], tables[1], "Foo")

	assert.Equal(t, 0, comparison.Additions)
	assert.Equal(t, 0, comparison.Deletions)
	assert.Equal(t, 1, comparison.Alterations)
	require.Equal(t, 1, len(comparison.Changes))

	assert.Equal(t, schema.ChangeColumn, comparison.Changes[0].Type)
	assert.Equal(t, "ALTER TABLE `Foo` CHANGE `FooID` `FooID` bigint(20) UNSIGNED NOT NULL;", comparison.Changes[0].SQL)
}

func TestCreateTable(t *testing.T) {

	s := NewMySQL(&lib.ConfigDatabase{})

	tables := testassets.TablesAddTable()

	comparison := s.CreateChangeSQL(tables[0], tables[1], "Foo")
	assert.Equal(t, 1, comparison.Additions)
	assert.Equal(t, 0, comparison.Deletions)
	assert.Equal(t, 0, comparison.Alterations)
	require.Equal(t, 1, len(comparison.Changes))

	assert.Equal(t, schema.CreateTable, comparison.Changes[0].Type)
	assert.Equal(t, "CREATE TABLE `Foo` (\n\t`FooID` int(10) UNSIGNED NOT NULL auto_increment,\n\t`IsDeleted` tinyint(4) SIGNED NOT NULL DEFAULT 0,\n\t`DateCreated` bigint(20) UNSIGNED NOT NULL DEFAULT 0,\n\tPRIMARY KEY(`FooID`)\n);", comparison.Changes[0].SQL)

	// for k := range comparison.Changes {
	// 	t.Log(comparison.Changes[k].SQL)
	// }

}

func TestDropTable(t *testing.T) {

	s := NewMySQL(&lib.ConfigDatabase{})

	tables := testassets.TablesDropTable()

	comparison := s.CreateChangeSQL(tables[0], tables[1], "Foo")
	assert.Equal(t, 0, comparison.Additions)
	assert.Equal(t, 1, comparison.Deletions)
	assert.Equal(t, 0, comparison.Alterations)
	require.Equal(t, 1, len(comparison.Changes))

	assert.Equal(t, schema.DropTable, comparison.Changes[0].Type)
	assert.Equal(t, "DROP TABLE `Foo`;", comparison.Changes[0].SQL)

	for k := range comparison.Changes {
		t.Log(comparison.Changes[k].SQL)
	}

}
