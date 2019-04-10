package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate_ShouldBuildObject(t *testing.T) {
	s := Update(NewTestObject())
	assert.IsType(t, &TestObject{}, s.Object)
	assert.Equal(t, "TestObject", s.Object.GetName())
}

func TestUpdate_Where_ShouldAddQueryParts(t *testing.T) {
	s := Update(NewTestObject()).
		Set("A", 1).
		Where(Equals{"id": 1})
	sql, args := s.ToSQL()
	assert.Equal(t, "UPDATE `TestObject` SET `A` = ? WHERE ( `id` = ? ) ", sql)
	assert.Equal(t, 2, len(args))
}

func TestUpdate_Where_ShouldAddMultipleWhereParts(t *testing.T) {
	s := Update(NewTestObject()).
		Set("A", 1).
		Set("B", 2).
		Where(Equals{"id": 1, "foo": "bar"}).
		Where(Equals{"id": 2, "foo": "baz"})
	sql, args := s.ToSQL()
	assert.Equal(t, "UPDATE `TestObject` SET `A` = ?, `B` = ? WHERE ( `foo` = ? AND `id` = ? ) AND ( `foo` = ? AND `id` = ? ) ", sql)
	assert.Equal(t, 6, len(args))
}

func TestUpdate_Or_ShouldAddQueryParts(t *testing.T) {
	s := Update(NewTestObject()).
		Set("A", 1).
		Set("B", 2).
		Or(Equals{"id": 1, "foo": "bar"}).
		Or(Equals{"id": 2, "foo": "baz"})
	sql, args := s.ToSQL()
	assert.Equal(t, "UPDATE `TestObject` SET `A` = ?, `B` = ? WHERE ( `foo` = ? AND `id` = ? ) OR ( `foo` = ? AND `id` = ? ) ", sql)
	assert.Equal(t, 6, len(args))
}

func TestUpdate_AllTheThings(t *testing.T) {
	s := Update(NewTestObject()).
		Set("A", "a").
		Set("B", 123).
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
	sql, args := s.ToSQL()
	assert.Equal(t, "UPDATE `TestObject` SET `A` = ?, `B` = ? WHERE ( `A` = ? OR `C` = ? ) OR ( `B` != ? AND `C` < ? AND `D` > ? ) ", sql)
	assert.Equal(t, 7, len(args))
}
