package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMin_ShouldBuildObject(t *testing.T) {
	s := Min("foo")
	s.Object = NewTestObject()
	assert.IsType(t, &TestObject{}, s.Object)
	assert.Equal(t, "TestObject", s.Object.GetName())
}

func TestMin_Where_ShouldAddQueryParts(t *testing.T) {
	s := Min("foo")
	s.Object = NewTestObject()
	s.Where(Equals{"id": 1})
	sql, args := s.ToSQL()
	assert.Equal(t, "SELECT MIN(`foo`) FROM `TestObject` WHERE ( `id` = ? ) ", sql)
	assert.Equal(t, 1, len(args))
}

func TestMin_Where_ShouldAddMultipleWhereParts(t *testing.T) {
	s := Min("foo")
	s.Object = NewTestObject()
	s.Where(Equals{"id": 1, "foo": "bar"})
	s.Where(Equals{"id": 2, "foo": "baz"})
	sql, args := s.ToSQL()
	assert.Equal(t, "SELECT MIN(`foo`) FROM `TestObject` WHERE ( `foo` = ? AND `id` = ? ) AND ( `foo` = ? AND `id` = ? ) ", sql)
	assert.Equal(t, 4, len(args))
}

func TestMin_Or_ShouldAddQueryParts(t *testing.T) {
	s := Min("foo")
	s.Object = NewTestObject()
	s.Or(Equals{"id": 1, "foo": "bar"})
	s.Or(Equals{"id": 2, "foo": "baz"})
	sql, args := s.ToSQL()
	assert.Equal(t, "SELECT MIN(`foo`) FROM `TestObject` WHERE ( `foo` = ? AND `id` = ? ) OR ( `foo` = ? AND `id` = ? ) ", sql)
	assert.Equal(t, 4, len(args))
}

func TestMin_AllTheThings(t *testing.T) {
	s := Min("foo").
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
		)
	s.Object = NewTestObject()
	sql, args := s.ToSQL()
	assert.Equal(t, "SELECT MIN(`foo`) FROM `TestObject` WHERE ( `A` = ? OR `C` = ? ) OR ( `B` != ? AND `C` < ? AND `D` > ? ) ", sql)
	assert.Equal(t, 5, len(args))
}
