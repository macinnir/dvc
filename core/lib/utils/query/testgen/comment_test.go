package testgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

func TestQuerySave_Insert(t *testing.T) {
	sql, e := (&Comment{
		DateCreated: 1620919850194,
		Content:     null.StringFrom("here is some test content"),
	}).Save().String()
	assert.Nil(t, e)
	assert.Equal(t, "INSERT INTO `Comment` ( `DateCreated`, `Content`, `ObjectType`, `ObjectID` ) VALUES ( 1620919850194, 'here is some test content', 0, 0 )", sql)
}

func TestQuerySave_Update(t *testing.T) {
	sql, e := (&Comment{
		DateCreated: 1620919850194,
		CommentID:   123,
		ObjectType:  1,
		ObjectID:    2,
		Content:     null.StringFrom("here is some test content"),
	}).Save().String()
	assert.Nil(t, e)
	assert.Equal(t, "UPDATE `Comment` SET `Content` = 'here is some test content', `ObjectType` = 1, `ObjectID` = 2 WHERE `CommentID` = 123", sql)
}

func TestComment_ToString(t *testing.T) {
	str := (&Comment{
		DateCreated: 1620919850194,
		CommentID:   123,
		ObjectType:  1,
		ObjectID:    2,
		Content:     null.StringFrom("here is some test content"),
	}).String()

	assert.Equal(t, `{"CommentID":123,"DateCreated":1620919850194,"IsDeleted":0,"Content":"here is some test content","ObjectType":1,"ObjectID":2}`, str)
}

// String() 3219, Save() 1023

// 2261 ns/op
// 2017 ns/op  		1472 B/op 		26 allocs/op
// 1959 ns/op 		1456 B/op 		26 allocs/op
func BenchmarkQuerySave_Insert(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.ReportAllocs()
		(&Comment{
			DateCreated: 1620919850194,
			Content:     null.StringFrom("here is some test content"),
			ObjectType:  1,
			ObjectID:    2,
		}).Save().String()
	}
	// assert.Nil(t, e)
	// assert.Equal(t, "INSERT INTO `Comment` ( `DateCreated`, `Content`, `ObjectType`, `ObjectID` ) VALUES ( 0, 'here is some test content', 0, 0 )", sql)
}

// 3219 ns/op
// 2965 ns/op
// 2267 ns/op
// 2088 ns/op 		1696 B/op 		28 allocs/op
// 1810 ns/op 		1552 B/op		27 allocs/op
//
func BenchmarkQuerySave_Update(b *testing.B) {
	for n := 0; n < b.N; n++ {
		(&Comment{
			DateCreated: 1620919850194,
			CommentID:   123,
			ObjectType:  1,
			ObjectID:    2,
			Content:     null.StringFrom("here is some test content"),
		}).Save().String()
	}
}
