package testgen

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/macinnir/dvc/core/lib/utils/query"
	"gopkg.in/guregu/null.v3"
)

type Time struct {
	time.Time
}

func (t *Time) String() string {
	return fmt.Sprint(t.Time.UnixNano() / 1000000)
}

// MarshalJSON implements the json.Marshaller interface
func (t *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.UnixNano() / 1000000)
}

// UnmarshalJSON implements the json.Unmarshaller interface
func (t *Time) UnmarshalJSON(data []byte) error {
	var i int64
	if e := json.Unmarshal(data, &i); e != nil {
		return e
	}

	t.Time = time.Unix(0, i*1000000)
	return nil
}

// // MarshalText implements encoding.TextMarshaler.
// // It will encode a blank string when this String is null.
// func (t *Time) MarshalText() ([]byte, error) {
// 	return []byte(fmt.Sprint(t.Time.))
// 	// if !s.Valid {
// 	// 	return []byte{}, nil
// 	// }
// 	// return []byte(s.String), nil
// }

// // UnmarshalText implements encoding.TextUnmarshaler.
// // It will unmarshal to a null String if the input is a blank string.
// func (s *String) UnmarshalText(text []byte) error {
// 	s.String = string(text)
// 	s.Valid = s.String != ""
// 	return nil
// }

func NewTime(t time.Time) Time {
	return Time{
		Time: t,
	}
}

type NullString struct {
	sql.NullString
}

func (n *NullString) String() string {
	return n.NullString.String
}

func NewNullString(str string) NullString {
	return NullString{
		NullString: sql.NullString{
			String: str,
			Valid:  true,
		},
	}
}

func (n *NullString) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String)
}

func (n *NullString) UnMarshalJSON(data []byte) error {
	var s string
	if e := json.Unmarshal(data, &s); e != nil {
		return e
	}
	n.NullString = sql.NullString{String: s, Valid: true}
	return nil
}

var (
	Comment_TableName query.TableName = "Comment"

	Comment_Column_CommentID   query.Column = "CommentID"
	Comment_Column_DateCreated query.Column = "DateCreated"
	Comment_Column_IsDeleted   query.Column = "IsDeleted"
	Comment_Column_Content     query.Column = "Content"
	Comment_Column_ObjectType  query.Column = "ObjectType"
	Comment_Column_ObjectID    query.Column = "ObjectID"

	Comment_Columns                    = []query.Column{Comment_Column_CommentID, Comment_Column_DateCreated, Comment_Column_IsDeleted, Comment_Column_Content, Comment_Column_ObjectType, Comment_Column_ObjectID}
	Comment_Column_Types               = map[query.Column]string{"CommentID": "%d", "DateCreated": "%d", "IsDeleted": "%d", "Content": "%s", "ObjectType": "%d", "ObjectID": "%d"}
	Comment_UpdateColumns              = []query.Column{Comment_Column_Content, Comment_Column_ObjectType, Comment_Column_ObjectID}
	Comment_InsertColumns              = []query.Column{Comment_Column_DateCreated, Comment_Column_Content, Comment_Column_ObjectType, Comment_Column_ObjectID}
	Comment_PrimaryKey    query.Column = "CommentID"
)

// Comment is a `Comment` data model
type Comment struct {
	CommentID   int64       `db:"CommentID" json:"CommentID"`
	DateCreated int64       `db:"DateCreated" json:"DateCreated"`
	IsDeleted   int         `db:"IsDeleted" json:"IsDeleted"`
	Content     null.String `db:"Content" json:"Content"`
	ObjectType  int64       `db:"ObjectType" json:"ObjectType"`
	ObjectID    int64       `db:"ObjectID" json:"ObjectID"`
}

// Comment_TableName is the name of the table
func (c *Comment) Table_Name() query.TableName {
	return Comment_TableName
}

func (c *Comment) Table_Columns() []query.Column {
	return Comment_Columns
}

func (c *Comment) Table_Column_Types() map[query.Column]string {
	return Comment_Column_Types
}

func (c *Comment) Table_PrimaryKey() query.Column {
	return Comment_PrimaryKey
}

// Comment_PrimaryKey is the name of the table's primary key
func (c *Comment) Table_PrimaryKey_Value() int64 {
	return c.CommentID
}

// Comment_InsertColumns is a list of all insert columns for this model
func (c *Comment) Table_InsertColumns() []query.Column {
	return Comment_InsertColumns
}

// Comment_UpdateColumns is a list of all update columns for this model
func (c *Comment) Table_UpdateColumns() []query.Column {
	return Comment_UpdateColumns
}

func (c *Comment) Select() *query.Q {
	return query.Select(c)
}

func (c *Comment) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

func (c *Comment) Create() string {

	var sql string
	sql, _ = query.Insert(c).
		Set("CommentID", c.CommentID).
		Set("DateCreated", c.DateCreated).
		Set("Content", c.Content.String).
		Set("ObjectType", c.ObjectType).
		Set("ObjectID", c.ObjectID).
		String()
	return sql
}

func (c *Comment) Update() string {
	var sql string
	sql, _ = query.Update(c).
		Set("Content", c.Content.String).
		Set("ObjectType", c.ObjectType).
		Set("ObjectID", c.ObjectID).
		Where(query.EQ("CommentID", c.CommentID)).
		String()
	return sql
}

func (c *Comment) Destroy() string {
	sql, _ := query.Delete(c).
		Where(
			query.EQ("CommentID", c.CommentID),
		).String()
	return sql
}

func (c *Comment) FromID(id int64) string {
	return ""
}
