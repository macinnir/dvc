package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelete_ShouldBuildObject(t *testing.T) {
	s := Delete()
	s.Object = NewTestObject()
	assert.IsType(t, &TestObject{}, s.Object)
	assert.Equal(t, "TestObject", s.Object.GetName())
}

func TestDelete_Where_ShouldAddQueryParts(t *testing.T) {
	s := Delete()
	s.Object = NewTestObject()
	s.Where(Equals{"id": 1})
	sql, args := s.ToSQL()
	assert.Equal(t, "DELETE FROM `TestObject` WHERE ( `id` = ? ) ", sql)
	assert.Equal(t, 1, len(args))
}

func TestDelete_Where_ShouldAddMultipleWhereParts(t *testing.T) {
	s := Delete()
	s.Object = NewTestObject()
	s.Where(Equals{"id": 1, "foo": "bar"})
	s.Where(Equals{"id": 2, "foo": "baz"})
	sql, args := s.ToSQL()
	assert.Equal(t, "DELETE FROM `TestObject` WHERE ( `foo` = ? AND `id` = ? ) AND ( `foo` = ? AND `id` = ? ) ", sql)
	assert.Equal(t, 4, len(args))
}

func TestDelete_Or_ShouldAddQueryParts(t *testing.T) {
	s := Delete()
	s.Object = NewTestObject()
	s.Or(Equals{"id": 1, "foo": "bar"})
	s.Or(Equals{"id": 2, "foo": "baz"})
	sql, args := s.ToSQL()
	assert.Equal(t, "DELETE FROM `TestObject` WHERE ( `foo` = ? AND `id` = ? ) OR ( `foo` = ? AND `id` = ? ) ", sql)
	assert.Equal(t, 4, len(args))
}

func TestDelete_Limit(t *testing.T) {
	s := Delete().Limit(10000)
	s.Object = NewTestObject()
	sql, args := s.ToSQL()
	assert.Equal(t, 1, len(args))
	assert.Equal(t, "DELETE FROM `TestObject` LIMIT ?", sql)
}

func TestDelete_AllTheThings(t *testing.T) {
	s := Delete().
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
		Limit(123)
	s.Object = NewTestObject()
	sql, args := s.ToSQL()
	assert.Equal(t, "DELETE FROM `TestObject` WHERE ( `A` = ? OR `C` = ? ) OR ( `B` != ? AND `C` < ? AND `D` > ? ) LIMIT ?", sql)
	assert.Equal(t, 6, len(args))
}
