package testassets

import (
	"encoding/json"

	"gopkg.in/guregu/null.v3"
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
func (c *Comment) Table_Name() string {
	return "Comment"
}

func (c *Comment) Table_Columns() []string {
	return []string{
		"CommentID",
		"DateCreated",
		"IsDeleted",
		"Content",
		"ObjectType",
		"ObjectID",
	}
}

func (c *Comment) Table_Column_Types() map[string]string {
	return map[string]string{
		"CommentID":   "%d",
		"DateCreated": "%d",
		"IsDeleted":   "%d",
		"Content":     "%s",
		"ObjectType":  "%d",
		"ObjectID":    "%d",
	}
}

func (c *Comment) Table_Column_Values() map[string]interface{} {
	return map[string]interface{}{
		"CommentID":   c.CommentID,
		"DateCreated": c.DateCreated,
		"IsDeleted":   c.IsDeleted,
		"Content":     c.Content.String,
		"ObjectType":  c.ObjectType,
		"ObjectID":    c.ObjectID,
	}
}

// Comment_PrimaryKey is the name of the table's primary key
func (c *Comment) Table_PrimaryKey() string {
	return "CommentID"
}

func (c *Comment) Table_PrimaryKey_Value() int64 {
	return c.CommentID
}

// Comment_InsertColumns is a list of all insert columns for this model
func (c *Comment) Table_InsertColumns() []string {
	return []string{"DateCreated", "Content", "ObjectType", "ObjectID"}
}

// Comment_UpdateColumns is a list of all update columns for this model
func (c *Comment) Table_UpdateColumns() []string {
	return []string{"Content", "ObjectType", "ObjectID"}
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

func (c *Comment) FromID() string {
	return ""
}
