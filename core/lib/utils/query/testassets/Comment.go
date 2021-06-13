package testassets

import (
	"encoding/json"

	"github.com/macinnir/dvc/core/lib/utils/query"
	"gopkg.in/guregu/null.v3"
)

const (

	// ListItem_SchemaName is the name of the schema group this model is in
	ListItem_SchemaName = "rfq"

	// Comment_TableName
	Comment_TableName query.TableName = "Comment"

	// Columns
	Comment_Column_CommentID   query.Column = "CommentID"
	Comment_Column_DateCreated query.Column = "DateCreated"
	Comment_Column_IsDeleted   query.Column = "IsDeleted"
	Comment_Column_Content     query.Column = "Content"
	Comment_Column_ObjectType  query.Column = "ObjectType"
	Comment_Column_ObjectID    query.Column = "ObjectID"
	Comment_Column_Name        query.Column = "Name"
)

// Comment is a `Comment` data model
type Comment struct {
	CommentID   int64       `db:"CommentID" json:"CommentID"`
	DateCreated int64       `db:"DateCreated" json:"DateCreated"`
	IsDeleted   int         `db:"IsDeleted" json:"IsDeleted"`
	Content     null.String `db:"Content" json:"Content"`
	ObjectType  int64       `db:"ObjectType" json:"ObjectType"`
	ObjectID    int64       `db:"ObjectID" json:"ObjectID"`
	Name        string      `db:"Name" json:"Name"`
}

// Comment_TableName is the name of the table
func (c *Comment) Table_Name() query.TableName {
	return Comment_TableName
}

func (c *Comment) Table_Columns() []query.Column {
	return []query.Column{
		Comment_Column_CommentID,
		Comment_Column_DateCreated,
		Comment_Column_IsDeleted,
		Comment_Column_Content,
		Comment_Column_ObjectType,
		Comment_Column_ObjectID,
		Comment_Column_Name,
	}
}

func (c *Comment) Table_Column_Types() map[query.Column]string {
	return map[query.Column]string{
		"CommentID":   "%d",
		"DateCreated": "%d",
		"IsDeleted":   "%d",
		"Content":     "%s",
		"ObjectType":  "%d",
		"ObjectID":    "%d",
		"Name":        "%s",
	}
}

func (c *Comment) Table_Column_Values() map[query.Column]interface{} {
	return map[query.Column]interface{}{
		"CommentID":   c.CommentID,
		"DateCreated": c.DateCreated,
		"IsDeleted":   c.IsDeleted,
		"Content":     c.Content.String,
		"ObjectType":  c.ObjectType,
		"ObjectID":    c.ObjectID,
		"Name":        c.Name,
	}
}

// Comment_PrimaryKey is the name of the table's primary key
func (c *Comment) Table_PrimaryKey() query.Column {
	return Comment_Column_CommentID
}

func (c *Comment) Table_PrimaryKey_Value() int64 {
	return c.CommentID
}

// Comment_InsertColumns is a list of all insert columns for this model
func (c *Comment) Table_InsertColumns() []query.Column {
	return []query.Column{
		Comment_Column_DateCreated,
		Comment_Column_Content,
		Comment_Column_ObjectType,
		Comment_Column_ObjectID,
		Comment_Column_Name,
	}
}

// Comment_UpdateColumns is a list of all update columns for this model
func (c *Comment) Table_UpdateColumns() []query.Column {
	return []query.Column{
		Comment_Column_Content,
		Comment_Column_ObjectType,
		Comment_Column_ObjectID,
		Comment_Column_Name,
	}
}

func (c *Comment) String() string {
	str, _ := json.Marshal(c)
	return string(str)
}

func (c *Comment) Create() string {
	return ""
}

func (c *Comment) Update() string {
	return ""
}

func (c *Comment) Destroy() string {
	return ""
}

func (c *Comment) FromID(id int64) string {
	return ""
}
