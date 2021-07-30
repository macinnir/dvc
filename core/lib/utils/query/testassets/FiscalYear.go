package testassets

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/macinnir/dvc/core/lib/utils/db"
	"github.com/macinnir/dvc/core/lib/utils/query"
)

const (

	// FiscalYear_SchemaName is the name of the schema group this model is in
	FiscalYear_SchemaName = "rfq"

	// FiscalYear_TableName is the name of the table
	FiscalYear_TableName query.TableName = "FiscalYear"

	// Columns
	FiscalYear_Column_FiscalYearID    query.Column = "FiscalYearID"
	FiscalYear_Column_FiscalYearKey   query.Column = "FiscalYearKey"
	FiscalYear_Column_Year            query.Column = "Year"
	FiscalYear_Column_DateFrom        query.Column = "DateFrom"
	FiscalYear_Column_DateTo          query.Column = "DateTo"
	FiscalYear_Column_DateCreated     query.Column = "DateCreated"
	FiscalYear_Column_IsDeleted       query.Column = "IsDeleted"
	FiscalYear_Column_TotalFiscalDays query.Column = "TotalFiscalDays"
	FiscalYear_Column_IsLocked        query.Column = "IsLocked"
)

var (
	// FiscalYear_Columns is a list of all the columns
	FiscalYear_Columns = []query.Column{
		FiscalYear_Column_FiscalYearID, FiscalYear_Column_FiscalYearKey, FiscalYear_Column_Year, FiscalYear_Column_DateFrom, FiscalYear_Column_DateTo, FiscalYear_Column_DateCreated, FiscalYear_Column_IsDeleted, FiscalYear_Column_TotalFiscalDays, FiscalYear_Column_IsLocked}

	// FiscalYear_Column_Types maps columns to their string types
	FiscalYear_Column_Types = map[query.Column]string{
		FiscalYear_Column_FiscalYearID: "%d", FiscalYear_Column_FiscalYearKey: "%d", FiscalYear_Column_Year: "%d", FiscalYear_Column_DateFrom: "%d", FiscalYear_Column_DateTo: "%d", FiscalYear_Column_DateCreated: "%d", FiscalYear_Column_IsDeleted: "%d", FiscalYear_Column_TotalFiscalDays: "%d", FiscalYear_Column_IsLocked: "%d"}
	// FiscalYear_UpdateColumns is a list of all update columns for this model
	FiscalYear_UpdateColumns = []query.Column{FiscalYear_Column_FiscalYearKey, FiscalYear_Column_Year, FiscalYear_Column_DateFrom, FiscalYear_Column_DateTo, FiscalYear_Column_IsDeleted, FiscalYear_Column_TotalFiscalDays, FiscalYear_Column_IsLocked}
	// FiscalYear_InsertColumns is a list of all insert columns for this model
	FiscalYear_InsertColumns = []query.Column{FiscalYear_Column_FiscalYearKey, FiscalYear_Column_Year, FiscalYear_Column_DateFrom, FiscalYear_Column_DateTo, FiscalYear_Column_DateCreated, FiscalYear_Column_TotalFiscalDays, FiscalYear_Column_IsLocked}
	// FiscalYear_PrimaryKey is the name of the table's primary key
	FiscalYear_PrimaryKey query.Column = "FiscalYearID"
)

// FiscalYear is a `FiscalYear` data model
type FiscalYear struct {
	FiscalYearID    int64 `db:"FiscalYearID" json:"FiscalYearID"`
	FiscalYearKey   int64 `db:"FiscalYearKey" json:"FiscalYearKey"`
	Year            int64 `db:"Year" json:"Year"`
	DateFrom        int64 `db:"DateFrom" json:"DateFrom"`
	DateTo          int64 `db:"DateTo" json:"DateTo"`
	DateCreated     int64 `db:"DateCreated" json:"DateCreated"`
	IsDeleted       int   `db:"IsDeleted" json:"IsDeleted"`
	TotalFiscalDays int64 `db:"TotalFiscalDays" json:"TotalFiscalDays"`
	IsLocked        int64 `db:"IsLocked" json:"IsLocked"`
}

// FiscalYear_TableName is the name of the table
func (c *FiscalYear) Table_Name() query.TableName {
	return FiscalYear_TableName
}

func (c *FiscalYear) Table_Columns() []query.Column {
	return FiscalYear_Columns
}

// Table_ColumnTypes returns a map of tableColumn names with their fmt string types
func (c *FiscalYear) Table_Column_Types() map[query.Column]string {
	return FiscalYear_Column_Types
}

// Table_PrimaryKey returns the name of this table's primary key
func (c *FiscalYear) Table_PrimaryKey() query.Column {
	return FiscalYear_PrimaryKey
}

// Table_PrimaryKey_Value returns the value of this table's primary key
func (c *FiscalYear) Table_PrimaryKey_Value() int64 {
	return c.FiscalYearID
}

// Table_InsertColumns is a list of all insert columns for this model
func (c *FiscalYear) Table_InsertColumns() []query.Column {
	return FiscalYear_InsertColumns
}

// Table_UpdateColumns is a list of all update columns for this model
func (c *FiscalYear) Table_UpdateColumns() []query.Column {
	return FiscalYear_UpdateColumns
}

// FiscalYear_SchemaName is the name of this table's schema
func (c *FiscalYear) Table_SchemaName() string {
	return FiscalYear_SchemaName
}

// FromID returns a FromID query statement
func (c *FiscalYear) FromID(db db.IDB, id int64) (query.IModel, error) {
	q, _ := query.Select(c).Where(query.EQ(FiscalYear_Column_FiscalYearID, id)).String()
	var e error

	model := &FiscalYear{}

	if e = db.Get(model, q); e != nil {
		if e == sql.ErrNoRows {
			return nil, nil
		}
		return nil, e
	}

	return model, nil
}

// String returns a json marshalled string of the object
func (c *FiscalYear) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

// Update updates a FiscalYear record
func (c *FiscalYear) Update(db db.IDB) error {
	var e error
	var ql string
	ql, _ = query.Update(c).
		Set(FiscalYear_Column_FiscalYearKey, c.FiscalYearKey).
		Set(FiscalYear_Column_Year, c.Year).
		Set(FiscalYear_Column_DateFrom, c.DateFrom).
		Set(FiscalYear_Column_DateTo, c.DateTo).
		Set(FiscalYear_Column_IsDeleted, c.IsDeleted).
		Set(FiscalYear_Column_TotalFiscalDays, c.TotalFiscalDays).
		Set(FiscalYear_Column_IsLocked, c.IsLocked).
		Where(query.EQ(FiscalYear_Column_FiscalYearID, c.FiscalYearID)).
		String()

	_, e = db.Exec(ql)
	if e != nil {
		return fmt.Errorf("FiscalYear.Update(): %w", e)
	}

	return e
}

// Create inserts a FiscalYear record
func (c *FiscalYear) Create(db db.IDB) error {

	var e error
	q := query.Insert(c)

	if c.FiscalYearID > 0 {
		q.Set(FiscalYear_Column_FiscalYearID, c.FiscalYearID)
	}
	q.Set(FiscalYear_Column_FiscalYearKey, c.FiscalYearKey)
	q.Set(FiscalYear_Column_Year, c.Year)
	q.Set(FiscalYear_Column_DateFrom, c.DateFrom)
	q.Set(FiscalYear_Column_DateTo, c.DateTo)
	q.Set(FiscalYear_Column_DateCreated, c.DateCreated)
	q.Set(FiscalYear_Column_TotalFiscalDays, c.TotalFiscalDays)
	q.Set(FiscalYear_Column_IsLocked, c.IsLocked)

	ql, _ := q.String()
	var result sql.Result
	result, e = db.Exec(ql)
	if e != nil {
		return fmt.Errorf("FiscalYear.Create(): %w", e)
	}

	// Assumes auto-increment
	if c.FiscalYearID == 0 {
		c.FiscalYearID, e = result.LastInsertId()
	}

	return e
}

// Destroy deletes a FiscalYear record
func (c *FiscalYear) Delete(db db.IDB) error {
	var e error
	ql, _ := query.Delete(c).
		Where(
			query.EQ(FiscalYear_Column_FiscalYearID, c.FiscalYearID),
		).String()

	_, e = db.Exec(ql)
	if e != nil {
		return fmt.Errorf("FiscalYear.Delete(): %w", e)
	}

	return e
}

func (r *FiscalYear) Raw(db db.IDB, queryRaw string) ([]*FiscalYear, error) {

	var e error
	model := []*FiscalYear{}
	e = db.Select(&model, queryRaw)

	if e != nil {
		return nil, fmt.Errorf("FiscalYear.Query(%s).Run(): %w", queryRaw, e)
	}

	fmt.Printf("FiscalYear.Query(%s).Run()\n", queryRaw)

	return model, nil
}

type FiscalYearDALSelector struct {
	db       db.IDB
	q        *query.Q
	isSingle bool
}

func (r *FiscalYear) Select(db db.IDB) *FiscalYearDALSelector {
	return &FiscalYearDALSelector{
		db: db,
		q:  query.Select(r),
	}
}

func (r *FiscalYearDALSelector) Alias(alias string) *FiscalYearDALSelector {
	r.q.Alias(alias)
	return r
}

func (r *FiscalYearDALSelector) Sum(fieldName query.Column, fieldAlias string) *FiscalYearDALSelector {
	r.q.Sum(fieldName, fieldAlias)
	return r
}

func (r *FiscalYearDALSelector) Count(fieldName query.Column, fieldAlias string) *FiscalYearDALSelector {
	r.q.Count(fieldName, fieldAlias)
	return r
}

func (r *FiscalYearDALSelector) Where(whereParts ...query.WherePart) *FiscalYearDALSelector {
	r.q.Where(whereParts...)
	return r
}

func (r *FiscalYearDALSelector) Limit(limit, offset int64) *FiscalYearDALSelector {
	r.q = r.q.Limit(limit, offset)
	return r
}

func (r *FiscalYearDALSelector) OrderBy(col query.Column, dir query.OrderDir) *FiscalYearDALSelector {
	r.q = r.q.OrderBy(col, dir)
	return r
}

func (r *FiscalYearDALSelector) Run() ([]*FiscalYear, error) {
	model := []*FiscalYear{}
	q, e := r.q.String()
	if e != nil {
		return nil, fmt.Errorf("FiscalYearDAL.Query.String(): %w", e)
	}

	e = r.db.Select(&model, q)

	if e != nil {
		return nil, fmt.Errorf("FiscalYearDAL.Query(%s).Run(): %w", q, e)
	}

	fmt.Printf("FiscalYearDAL.Query(%s).Run()\n", q)

	return model, nil
}

// Counter
type FiscalYearDALCounter struct {
	db db.IDB
	q  *query.Q
}

func (r *FiscalYear) Count(db db.IDB) *FiscalYearDALCounter {
	return &FiscalYearDALCounter{
		db: db,
		q:  query.Select(r).Count(r.Table_PrimaryKey(), "c"),
	}
}

func (r *FiscalYearDALCounter) Alias(alias string) *FiscalYearDALCounter {
	r.q.Alias(alias)
	return r
}

func (ds *FiscalYearDALCounter) Where(whereParts ...query.WherePart) *FiscalYearDALCounter {
	ds.q.Where(whereParts...)
	return ds
}

func (ds *FiscalYearDALCounter) Run() (int64, error) {

	count := int64(0)
	q, e := ds.q.String()
	if e != nil {
		return 0, fmt.Errorf("FiscalYearDALCounter.Query.String(): %w", e)
	}

	e = ds.db.Get(&count, q)

	if e != nil {
		return 0, fmt.Errorf("FiscalYearDALCounter.Query(%s).Run(): %w", q, e)
	}

	fmt.Printf("FiscalYearDALCounter.Query(%s).Run()\n", q)

	return count, nil
}

// Summer
type FiscalYearDALSummer struct {
	db db.IDB
	q  *query.Q
}

func (r *FiscalYear) Sum(db db.IDB, col query.Column) *FiscalYearDALSummer {
	return &FiscalYearDALSummer{
		db: db,
		q:  query.Select(r).Sum(col, "c"),
	}
}

func (ds *FiscalYearDALSummer) Where(whereParts ...query.WherePart) *FiscalYearDALSummer {
	ds.q.Where(whereParts...)
	return ds
}

func (ds *FiscalYearDALSummer) Run() (float64, error) {

	sum := float64(0)
	q, e := ds.q.String()
	if e != nil {
		return 0, fmt.Errorf("FiscalYearDALSummer.Query.String(): %w", e)
	}

	e = ds.db.Get(&sum, q)

	if e != nil {
		return 0, fmt.Errorf("FiscalYearDALSummer.Query(%s).Run(): %w", q, e)
	}

	fmt.Printf("FiscalYearDALSummer.Query(%s).Run()\n", q)

	return sum, nil
}

type FiscalYearDALGetter struct {
	db db.IDB
	q  *query.Q
}

func (r *FiscalYear) Get(db db.IDB) *FiscalYearDALGetter {
	return &FiscalYearDALGetter{
		db: db,
		q:  query.Select(r),
	}
}

func (r *FiscalYearDALGetter) Alias(alias string) *FiscalYearDALGetter {
	r.q.Alias(alias)
	return r
}

func (ds *FiscalYearDALGetter) Where(whereParts ...query.WherePart) *FiscalYearDALGetter {
	ds.q.Where(whereParts...)
	return ds
}

func (ds *FiscalYearDALGetter) OrderBy(col query.Column, dir query.OrderDir) *FiscalYearDALGetter {
	ds.q = ds.q.OrderBy(col, dir)
	return ds
}

func (ds *FiscalYearDALGetter) Run() (*FiscalYear, error) {

	model := &FiscalYear{}
	q, e := ds.q.String()
	if e != nil {
		return nil, fmt.Errorf("FiscalYearDALGetter.Query.String(): %w", e)
	}

	e = ds.db.Get(model, q)

	if e != nil {
		if e == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("FiscalYearDALGetter.Run(%s): %w", q, e)
	}

	fmt.Printf("FiscalYearDALGetter.Get(%s).Run()\n", q)

	return model, nil
}
