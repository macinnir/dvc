package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelect_ShouldBuildObject(t *testing.T) {
	s := Select()
	s.Object = NewTestObject()
	assert.IsType(t, &TestObject{}, s.Object)
	assert.Equal(t, "TestObject", s.Object.GetName())
}

func TestSelect_Where_ShouldAddQueryParts(t *testing.T) {
	s := Select()
	s.Object = NewTestObject()
	s.Where(Equals{"id": 1})
	sql, args := s.ToSQL()
	assert.Equal(t, "SELECT `ID`,`Name`,`DateCreated` FROM `TestObject` WHERE ( `id` = ? ) ", sql)
	assert.Equal(t, 1, len(args))
}

func TestSelect_Where_ShouldAddMultipleWhereParts(t *testing.T) {
	s := Select()
	s.Object = NewTestObject()
	s.Where(Equals{"id": 1, "foo": "bar"})
	s.Where(Equals{"id": 2, "foo": "baz"})
	sql, args := s.ToSQL()
	assert.Equal(t, "SELECT `ID`,`Name`,`DateCreated` FROM `TestObject` WHERE ( `foo` = ? AND `id` = ? ) AND ( `foo` = ? AND `id` = ? ) ", sql)
	assert.Equal(t, 4, len(args))
}

func TestSelect_Or_ShouldAddQueryParts(t *testing.T) {
	s := Select()
	s.Object = NewTestObject()
	s.Or(Equals{"id": 1, "foo": "bar"})
	s.Or(Equals{"id": 2, "foo": "baz"})
	sql, args := s.ToSQL()
	assert.Equal(t, "SELECT `ID`,`Name`,`DateCreated` FROM `TestObject` WHERE ( `foo` = ? AND `id` = ? ) OR ( `foo` = ? AND `id` = ? ) ", sql)
	assert.Equal(t, 4, len(args))
}

func TestSelect_OrderBy(t *testing.T) {
	s := Select().OrderBy("a", DESC).OrderBy("b", ASC)
	s.Object = NewTestObject()
	sql, args := s.ToSQL()
	assert.Equal(t, 0, len(args))
	assert.Equal(t, "SELECT `ID`,`Name`,`DateCreated` FROM `TestObject` ORDER BY `a` DESC, `b` ASC ", sql)
}

func TestSelect_Limit(t *testing.T) {
	s := Select().Limit(10000, 1234)
	s.Object = NewTestObject()
	sql, args := s.ToSQL()
	assert.Equal(t, 2, len(args))
	assert.Equal(t, "SELECT `ID`,`Name`,`DateCreated` FROM `TestObject` LIMIT ?,?", sql)
}

func TestSelect_AllTheThings(t *testing.T) {
	s := Select().
		Distinct().
		Where(
			Equals{"A": 1},
			Or{},
			Equals{"C": 2},
		).
		Or(
			NotEquals{"B": 2},
			And{},
			LessThan{"C": 123},
			And{},
			GreaterThan{"D": 234},
		).
		Limit(123, 456).
		OrderBy("B", DESC).OrderBy("A", ASC)
	s.Object = NewTestObject()
	sql, args := s.ToSQL()
	assert.Equal(t, "SELECT DISTINCT `ID`,`Name`,`DateCreated` FROM `TestObject` WHERE ( `A` = ? OR `C` = ? ) OR ( `B` != ? AND `C` < ? AND `D` > ? ) ORDER BY `B` DESC, `A` ASC LIMIT ?,?", sql)
	assert.Equal(t, 7, len(args))
}
