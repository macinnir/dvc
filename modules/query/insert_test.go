package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert_ShouldBuildObject(t *testing.T) {
	s := Insert()
	s.Object = NewTestObject()
	assert.IsType(t, &TestObject{}, s.Object)
	assert.Equal(t, "TestObject", s.Object.GetName())
}
func TestInsert_AllTheThings(t *testing.T) {
	s := Insert().
		Value("A", "a").
		Value("B", 123)
	s.Object = NewTestObject()
	sql, args := s.ToSQL()
	assert.Equal(t, "INSERT INTO `TestObject` (`A`, `B`) VALUES (?, ?)", sql)
	assert.Equal(t, 2, len(args))
}
