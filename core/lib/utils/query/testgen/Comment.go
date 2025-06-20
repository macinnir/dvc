package testgen

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/macinnir/dvc/core/lib/utils/db"
	"github.com/macinnir/dvc/core/lib/utils/query"
	"gopkg.in/guregu/null.v3"
)

const (

	// Comment_SchemaName is the name of the schema group this model is in
	Comment_SchemaName = "core"

	// Comment_TableName is the name of the table
	Comment_TableName query.TableName = "Comment"

	// Columns

	Comment_Column_CommentID   query.Column = "CommentID"
	Comment_Column_Content     query.Column = "Content"
	Comment_Column_DateCreated query.Column = "DateCreated"
	Comment_Column_IsDeleted   query.Column = "IsDeleted"
	Comment_Column_ObjectID    query.Column = "ObjectID"
	Comment_Column_ObjectType  query.Column = "ObjectType"
)

var (
	// Comment_Columns is a list of all the columns
	Comment_Columns = []query.Column{
		Comment_Column_CommentID,
		Comment_Column_Content,
		Comment_Column_DateCreated,
		Comment_Column_IsDeleted,
		Comment_Column_ObjectID,
		Comment_Column_ObjectType,
	}

	// Comment_Column_Types maps columns to their string types
	Comment_Column_Types = map[query.Column]string{
		Comment_Column_CommentID:   "%d",
		Comment_Column_Content:     "%s",
		Comment_Column_DateCreated: "%d",
		Comment_Column_IsDeleted:   "%d",
		Comment_Column_ObjectID:    "%d",
		Comment_Column_ObjectType:  "%d",
	}

	// Update columns
	// Comment_UpdateColumns is a list of all update columns for this model
	Comment_UpdateColumns = []query.Column{
		Comment_Column_Content,
		Comment_Column_ObjectID,
		Comment_Column_ObjectType,
	}

	// Insert columns
	// Comment_InsertColumns is a list of all insert columns for this model
	Comment_InsertColumns = []query.Column{
		Comment_Column_Content,
		Comment_Column_DateCreated,
		Comment_Column_ObjectID,
		Comment_Column_ObjectType,
	}

	// Primary Key
	// Comment_PrimaryKey is the name of the table's primary key
	Comment_PrimaryKey query.Column = "CommentID"
)

// Comment is a data model
type Comment struct {
	CommentID   int64       `db:"CommentID" json:"CommentID"`
	Content     null.String `db:"Content" json:"Content"`
	DateCreated int64       `db:"DateCreated" json:"DateCreated"`
	IsDeleted   int         `db:"IsDeleted" json:"IsDeleted"`
	ObjectID    int64       `db:"ObjectID" json:"ObjectID"`
	ObjectType  int64       `db:"ObjectType" json:"ObjectType"`
}

// 63

// 68

// Comment_TableName is the name of the table
func (c *Comment) Table_Name() query.TableName {
	return Comment_TableName
}

func (c *Comment) Table_Columns() []query.Column {
	return Comment_Columns
}

// Table_ColumnTypes returns a map of tableColumn names with their fmt string types
func (c *Comment) Table_Column_Types() map[query.Column]string {
	return Comment_Column_Types
}

// Table_PrimaryKey returns the name of this table's primary key
func (c *Comment) Table_PrimaryKey() query.Column {
	return Comment_PrimaryKey
}

// Table_PrimaryKey_Value returns the value of this table's primary key
func (c *Comment) Table_PrimaryKey_Value() int64 {
	return c.CommentID
}

// Table_InsertColumns is a list of all insert columns for this model
func (c *Comment) Table_InsertColumns() []query.Column {
	return Comment_InsertColumns
}

// Table_UpdateColumns is a list of all update columns for this model
func (c *Comment) Table_UpdateColumns() []query.Column { // 100
	return Comment_UpdateColumns
}

// Comment_SchemaName is the name of this table's schema
func (c *Comment) Table_SchemaName() string {
	return Comment_SchemaName
}

// FromID returns a FromID query statement
func (c *Comment) FromID(db db.IDB, id int64) (query.IModel, error) {

	sel := query.Select(c)
	sel.Fields(
		query.NewField(query.FieldTypeBasic, Comment_Column_CommentID),
		query.NewField(query.FieldTypeBasic, Comment_Column_Content),
		query.NewField(query.FieldTypeBasic, Comment_Column_DateCreated),
		query.NewField(query.FieldTypeBasic, Comment_Column_IsDeleted),
		query.NewField(query.FieldTypeBasic, Comment_Column_ObjectID),
		query.NewField(query.FieldTypeBasic, Comment_Column_ObjectType),
	)
	q, e := sel.String()
	if e != nil {
		return nil, fmt.Errorf("Comment.FromID.Query.String(): %w", e)
	}

	row := db.QueryRow(q)

	switch e = row.Scan(
		&c.CommentID,
		&c.Content,
		&c.DateCreated,
		&c.IsDeleted,
		&c.ObjectID,
		&c.ObjectType,
	); e {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		// fmt.Printf("CommentDALGetter.Get(%s).Run()\n", q)
		return c, nil
	default:
		return nil, fmt.Errorf("CommentDALGetter(%s).Run(): %w", q, e)
	}
}

// String returns a json marshalled string of the object
func (c *Comment) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

// Update updates a Comment record
func (c *Comment) Update(db db.IDB) error {
	var e error
	var ql string
	ql, _ = query.Update(c).
		Set(Comment_Column_Content, c.Content.String).
		Set(Comment_Column_ObjectID, c.ObjectID).
		Set(Comment_Column_ObjectType, c.ObjectType).
		Where(query.EQ(Comment_Column_CommentID, c.CommentID)).
		String()
	_, e = db.Exec(ql)
	if e != nil {
		return fmt.Errorf("Comment.Update(): %w", e)
	}

	return e
}

// Create inserts a Comment record
func (c *Comment) Create(db db.IDB) error {

	var e error

	q := query.Insert(c)

	if c.CommentID > 0 {
		q.Set(Comment_Column_CommentID, c.CommentID)
	}

	q.Set(Comment_Column_Content, c.Content.String)
	q.Set(Comment_Column_DateCreated, c.DateCreated)
	q.Set(Comment_Column_ObjectID, c.ObjectID)
	q.Set(Comment_Column_ObjectType, c.ObjectType)

	ql, _ := q.String()
	var result sql.Result
	result, e = db.Exec(ql)

	if e != nil {
		return fmt.Errorf("Comment.Create(): %w", e) // 177
	}

	// Assumes auto-increment
	if c.CommentID == 0 {
		c.CommentID, e = result.LastInsertId()
	}

	return e
}

// Destroy deletes a Comment record
func (c *Comment) Delete(db db.IDB) error {
	var e error
	ql, _ := query.Delete(c).
		Where(
			query.EQ(Comment_Column_CommentID, c.CommentID),
		).String()

	_, e = db.Exec(ql)
	if e != nil {
		return fmt.Errorf("Comment.Delete(): %w", e)
	}

	return e
}

func (r *Comment) Raw(db db.IDB, queryRaw string) ([]*Comment, error) {

	var e error
	model := []*Comment{}
	sel := query.Select(r)
	sel.Fields(
		query.NewField(query.FieldTypeBasic, Comment_Column_CommentID),
		query.NewField(query.FieldTypeBasic, Comment_Column_Content),
		query.NewField(query.FieldTypeBasic, Comment_Column_DateCreated),
		query.NewField(query.FieldTypeBasic, Comment_Column_IsDeleted),
		query.NewField(query.FieldTypeBasic, Comment_Column_ObjectID),
		query.NewField(query.FieldTypeBasic, Comment_Column_ObjectType),
	)

	q, e := sel.String()
	if e != nil {
		return nil, fmt.Errorf("CommentDAL.Raw.String(): %w", e)
	}

	var rows *sql.Rows
	rows, e = db.Query(q)

	if e != nil {
		if e == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("CommentDAL.Raw.Run(%s): %w", q, e)
	}

	defer rows.Close()
	for rows.Next() {
		m := &Comment{}
		if e = rows.Scan(
			&m.CommentID,
			&m.Content,
			&m.DateCreated,
			&m.IsDeleted,
			&m.ObjectID,
			&m.ObjectType,
		); e != nil {
			return nil, fmt.Errorf("CommentDALRaw(%s).Run(): %w", q, e)
		}
		model = append(model, m)
	}

	// fmt.Printf("CommentDAL.Raw(%s).Run()\n", q)

	return model, nil
}

type ICommentDALSelector interface {
	Select(db db.IDB) ICommentDALSelector
}

type CommentDALSelector struct {
	db       db.IDB
	q        *query.Q
	isSingle bool
}

func (r *Comment) Select(db db.IDB) *CommentDALSelector {
	return &CommentDALSelector{
		db: db,
		q:  query.Select(r),
	}
}

func (r *CommentDALSelector) Alias(alias string) *CommentDALSelector {
	r.q.Alias(alias)
	return r
}

func (r *CommentDALSelector) Where(whereParts ...*query.WherePart) *CommentDALSelector {
	r.q.Where(whereParts...)
	return r
}

func (r *CommentDALSelector) Limit(limit, offset int64) *CommentDALSelector {
	r.q = r.q.Limit(limit, offset)
	return r
}

func (r *CommentDALSelector) OrderBy(col query.Column, dir query.OrderDir) *CommentDALSelector {
	r.q = r.q.OrderBy(col, dir)
	return r
}

func (r *CommentDALSelector) Run() ([]*Comment, error) {
	model := []*Comment{}
	r.q.Fields(
		query.NewField(query.FieldTypeBasic, Comment_Column_CommentID),
		query.NewField(query.FieldTypeBasic, Comment_Column_Content),
		query.NewField(query.FieldTypeBasic, Comment_Column_DateCreated),
		query.NewField(query.FieldTypeBasic, Comment_Column_IsDeleted),
		query.NewField(query.FieldTypeBasic, Comment_Column_ObjectID),
		query.NewField(query.FieldTypeBasic, Comment_Column_ObjectType),
	)

	q, e := r.q.String()
	if e != nil {
		return nil, fmt.Errorf("CommentDAL.Query.String(): %w", e)
	}

	var rows *sql.Rows
	rows, e = r.db.Query(q)

	if e != nil {
		if e == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("CommentDALSelector.Run(%s): %w", q, e)
	}

	defer rows.Close()
	for rows.Next() {
		m := &Comment{}
		if e = rows.Scan(
			&m.CommentID,
			&m.Content,
			&m.DateCreated,
			&m.IsDeleted,
			&m.ObjectID,
			&m.ObjectType,
		); e != nil {
			return nil, fmt.Errorf("CommentDALSelector(%s).Run(): %w", q, e)
		}

		model = append(model, m)
	}

	// fmt.Printf("CommentDALSelector(%s).Run()\n", q)

	return model, nil
}

// Counter
type CommentDALCounter struct {
	db db.IDB
	q  *query.Q
}

func (r *Comment) Count(db db.IDB) *CommentDALCounter {
	return &CommentDALCounter{
		db: db,
		q:  query.Select(r).Count(r.Table_PrimaryKey(), "c"),
	}
}

func (r *CommentDALCounter) Alias(alias string) *CommentDALCounter {
	r.q.Alias(alias)
	return r
}

func (ds *CommentDALCounter) Where(whereParts ...*query.WherePart) *CommentDALCounter {
	ds.q.Where(whereParts...)
	return ds
}

func (ds *CommentDALCounter) Run() (int64, error) {

	count := int64(0)
	q, e := ds.q.String()
	if e != nil {
		return 0, fmt.Errorf("CommentDALCounter.Query.String(): %w", e)
	}

	row := ds.db.QueryRow(q)

	switch e = row.Scan(&count); e {
	case sql.ErrNoRows:
		return 0, nil
	case nil:
		// fmt.Printf("CommentDALCounter.QueryRow(%s).Run()\n", q)
		return count, nil
	default:
		return 0, fmt.Errorf("CommentDALCounter.QueryRow(%s).Run(): %w", q, e)
	}
}

// Summer
type CommentDALSummer struct {
	db db.IDB
	q  *query.Q
}

func (r *Comment) Sum(db db.IDB, col query.Column) *CommentDALSummer {
	return &CommentDALSummer{
		db: db,
		q:  query.Select(r).Sum(col, "c"),
	}
}

func (ds *CommentDALSummer) Where(whereParts ...*query.WherePart) *CommentDALSummer {
	ds.q.Where(whereParts...)
	return ds
}

func (ds *CommentDALSummer) Run() (float64, error) {

	sum := float64(0)
	q, e := ds.q.String()
	if e != nil {
		return 0, fmt.Errorf("CommentDALSummer.Query.String(): %w", e)
	}

	row := ds.db.QueryRow(q)

	switch e = row.Scan(&sum); e {
	case sql.ErrNoRows:
		return 0, nil
	case nil:
		// fmt.Printf("CommentDALSummer.QueryRow(%s).Run()\n", q)
		return sum, nil
	default:
		return 0, fmt.Errorf("CommentDALSummer.QueryRow(%s).Run(): %w", q, e)
	}
}

// Minner
type CommentDALMinner struct {
	db db.IDB
	q  *query.Q
}

func (r *Comment) Min(db db.IDB, col query.Column) *CommentDALMinner {
	return &CommentDALMinner{
		db: db,
		q:  query.Select(r).Min(col, "c"),
	}
}

func (ds *CommentDALMinner) Where(whereParts ...*query.WherePart) *CommentDALMinner {
	ds.q.Where(whereParts...)
	return ds
}

func (ds *CommentDALMinner) Run() (float64, error) {

	min := float64(0)
	q, e := ds.q.String()
	if e != nil {
		return 0, fmt.Errorf("CommentDALMinner.Query.String(): %w", e)
	}

	row := ds.db.QueryRow(q)

	switch e = row.Scan(&min); e {
	case sql.ErrNoRows:
		return 0, nil
	case nil:
		// fmt.Printf("CommentDALMinner.QueryRow(%s).Run()\n", q)
		return min, nil
	default:
		return 0, fmt.Errorf("CommentDALMinner.QueryRow(%s).Run(): %w", q, e)
	}
}

// Maxer
type CommentDALMaxer struct {
	db db.IDB
	q  *query.Q
}

func (r *Comment) Max(db db.IDB, col query.Column) *CommentDALMaxer {
	return &CommentDALMaxer{
		db: db,
		q:  query.Select(r).Max(col, "c"),
	}
}

func (ds *CommentDALMaxer) Where(whereParts ...*query.WherePart) *CommentDALMaxer {
	ds.q.Where(whereParts...)
	return ds
}

func (ds *CommentDALMaxer) Run() (float64, error) {

	max := float64(0)
	q, e := ds.q.String()
	if e != nil {
		return 0, fmt.Errorf("CommentDALMaxer.Query.String(): %w", e)
	}

	row := ds.db.QueryRow(q)

	switch e = row.Scan(&max); e {
	case sql.ErrNoRows:
		return 0, nil
	case nil:
		// fmt.Printf("CommentDALMaxer.QueryRow(%s).Run()\n", q)
		return max, nil
	default:
		return 0, fmt.Errorf("CommentDALMaxer.QueryRow(%s).Run(): %w", q, e)
	}
}

type CommentDALGetter struct {
	db db.IDB
	q  *query.Q
}

func (r *Comment) Get(db db.IDB) *CommentDALGetter {
	return &CommentDALGetter{
		db: db,
		q:  query.Select(r),
	}
}

func (r *CommentDALGetter) Alias(alias string) *CommentDALGetter {
	r.q.Alias(alias)
	return r
}

func (ds *CommentDALGetter) Where(whereParts ...*query.WherePart) *CommentDALGetter {
	ds.q.Where(whereParts...)
	return ds
}

func (ds *CommentDALGetter) OrderBy(col query.Column, dir query.OrderDir) *CommentDALGetter {
	ds.q = ds.q.OrderBy(col, dir)
	return ds
}

func (ds *CommentDALGetter) Run() (*Comment, error) {

	model := &Comment{}

	ds.q.Fields(
		query.NewField(query.FieldTypeBasic, Comment_Column_CommentID),
		query.NewField(query.FieldTypeBasic, Comment_Column_Content),
		query.NewField(query.FieldTypeBasic, Comment_Column_DateCreated),
		query.NewField(query.FieldTypeBasic, Comment_Column_IsDeleted),
		query.NewField(query.FieldTypeBasic, Comment_Column_ObjectID),
		query.NewField(query.FieldTypeBasic, Comment_Column_ObjectType),
	)
	q, e := ds.q.String()
	if e != nil {
		return nil, fmt.Errorf("CommentDALGetter.Query.String(): %w", e)
	}

	row := ds.db.QueryRow(q)

	switch e = row.Scan(
		&model.CommentID,
		&model.Content,
		&model.DateCreated,
		&model.IsDeleted,
		&model.ObjectID,
		&model.ObjectType,
	); e {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		// fmt.Printf("CommentDALGetter.Get(%s).Run()\n", q)
		return model, nil
	default:
		return nil, fmt.Errorf("CommentDALGetter(%s).Run(): %w", q, e)
	}
}
