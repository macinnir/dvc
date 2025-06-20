package testgen

import (
	"context"
	"database/sql"
	"testing"

	"github.com/macinnir/dvc/core/lib/utils/db"
	"github.com/macinnir/dvc/core/lib/utils/query"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

type MockDB struct{}

func (m *MockDB) Close() {}
func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}
func (m *MockDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return nil, nil
}
func (m *MockDB) ExecMany(stmts []string, chunkSize int) (e error) {
	return nil
}
func (m *MockDB) Host() string {
	return ""
} // The host name (from config) {}
func (m *MockDB) Name() string {
	return ""
} // The name of the database (from config) {}
func (m *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}
func (m *MockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return nil
}

func NewMockDB() db.IDB {
	return &MockDB{}
}

func TestQuerySave_Insert(t *testing.T) {

	sql := (&Comment{
		CommentID:   12345,
		DateCreated: 1620919850194,
		Content:     null.StringFrom("here is some test content"),
	}).Create(NewMockDB())
	assert.Equal(t, "INSERT INTO `Comment` ( `CommentID`, `DateCreated`, `Content`, `ObjectType`, `ObjectID` ) VALUES ( 12345, 1620919850194, 'here is some test content', 0, 0 )", sql)
}

func TestQuerySave_Update(t *testing.T) {
	sql, _ := query.Update(
		&Comment{
			DateCreated: 1620919850194,
			CommentID:   123,
			ObjectType:  1,
			ObjectID:    2,
			Content:     null.StringFrom("here is some test content"),
		},
	).
		Set(Comment_Column_Content, c.Content.String).
		Set(Comment_Column_ObjectID, c.ObjectID).
		Set(Comment_Column_ObjectType, c.ObjectType).
		Where(query.EQ(Comment_Column_CommentID, c.CommentID)).
		String()
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

	assert.Equal(t, `{"CommentID":123,"Content":"here is some test content","DateCreated":1620919850194,"IsDeleted":0,"ObjectID":2,"ObjectType":1}`, str)
}

// String() 3219, Save() 1023

// 2261 ns/op
// 2017 ns/op  		1472 B/op 		26 allocs/op
// 1959 ns/op 		1456 B/op 		26 allocs/op
// 1795 ns/op 		1512 B/op 		25 allocs/op
func BenchmarkQuerySave_Create(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.ReportAllocs()
		(&Comment{
			DateCreated: 1620919850194,
			Content:     null.StringFrom("here is some test content"),
			ObjectType:  1,
			ObjectID:    2,
		}).Create(NewMockDB())
	}
	// assert.Nil(t, e)
	// assert.Equal(t, "INSERT INTO `Comment` ( `DateCreated`, `Content`, `ObjectType`, `ObjectID` ) VALUES ( 0, 'here is some test content', 0, 0 )", sql)
}

// 3219 ns/op
// 2965 ns/op
// 2267 ns/op
// 2088 ns/op 		1696 B/op 		28 allocs/op
// 1810 ns/op 		1552 B/op		27 allocs/op
func BenchmarkQuerySave_Update(b *testing.B) {
	for n := 0; n < b.N; n++ {
		(&Comment{
			DateCreated: 1620919850194,
			CommentID:   123,
			ObjectType:  1,
			ObjectID:    2,
			Content:     null.StringFrom("here is some test content"),
		}).Update(NewMockDB())
	}
}
