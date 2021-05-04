package gen

import (
	"testing"

	"github.com/macinnir/dvc/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	TestGenTableExpected1String = `// Table1DAL SQL
const (
	Table1DALInsertSQL = "INSERT INTO ` + "`Table1` (`Foo`,`Bar`,`DateCreated`) VALUES (:Foo,:Bar,:DateCreated)\"" + ` 
	Table1DALUpdateSQL = "UPDATE ` + "`Table1` SET `Foo` = :Foo,`Bar` = :Bar WHERE Table1ID = :Table1ID\"" + `
)`
	TestGenTableExpected2String = `// Table2DAL SQL
const (
	Table2DALInsertSQL = "INSERT INTO ` + "`Table2` (`Baz`,`Quux`,`DateCreated`) VALUES (:Baz,:Quux,:DateCreated)\"" + ` 
	Table2DALUpdateSQL = "UPDATE ` + "`Table2` SET `Baz` = :Baz,`Quux` = :Quux WHERE Table2ID = :Table2ID\"" + `
)`
	TestGenDatabase1String = `package footest

` + TestGenTableExpected1String + `

` + TestGenTableExpected2String + `
`
)

var TestGenTableObj1 = &lib.Table{
	Name: "Table1",
	Columns: map[string]*lib.Column{
		"Table1ID":    {Name: "Table1ID", ColumnKey: "PRI", DataType: "int", Position: 1},
		"Foo":         {Name: "Foo", DataType: "varchar", Position: 2},
		"Bar":         {Name: "Bar", DataType: "varchar", Position: 3},
		"IsDeleted":   {Name: "IsDeleted", DataType: "tinyint", Position: 4},
		"DateCreated": {Name: "DateCreated", DataType: "bigint", Position: 5},
	},
}

var TestGenTableObj2 = &lib.Table{
	Name: "Table2",
	Columns: map[string]*lib.Column{
		"Table2ID":    {Name: "Table2ID", ColumnKey: "PRI", DataType: "int", Position: 1},
		"Baz":         {Name: "Baz", DataType: "varchar", Position: 2},
		"Quux":        {Name: "Quux", DataType: "varchar", Position: 3},
		"IsDeleted":   {Name: "IsDeleted", DataType: "tinyint", Position: 4},
		"DateCreated": {Name: "DateCreated", DataType: "bigint", Position: 5},
	},
}

var TestGenDatabase1 = &lib.Database{
	Tables: map[string]*lib.Table{
		"Table1": TestGenTableObj1,
		"Table2": TestGenTableObj2,
	},
}

func TestGenerateTableInsertAndUpdateFields(t *testing.T) {

	out, e := generateTableInsertAndUpdateFields(TestGenTableObj1)

	require.Nil(t, e)
	assert.Equal(t, TestGenTableExpected1String, out)

}

func TestGenerateDALSQL(t *testing.T) {
	out, e := generateDALSQL("footest", TestGenDatabase1)
	require.Nil(t, e)
	assert.Equal(t, TestGenDatabase1String, out)
}
