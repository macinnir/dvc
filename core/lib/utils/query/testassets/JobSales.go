package testassets

import (
	"encoding/json"

	"github.com/macinnir/dvc/core/lib/utils/query"
)

const (
	JobSales_TableName query.TableName = "JobSales"

	JobSales_Column_JobSalesID        query.Column = "JobSalesID"
	JobSales_Column_JobID             query.Column = "JobID"
	JobSales_Column_UserID            query.Column = "UserID"
	JobSales_Column_CommissionPercent query.Column = "CommissionPercent"
	JobSales_Column_DateCreated       query.Column = "DateCreated"
	JobSales_Column_IsDeleted         query.Column = "IsDeleted"
	JobSales_Column_CommissionDollars query.Column = "CommissionDollars"
	JobSales_Column_IsHouse           query.Column = "IsHouse"
)

// JobSales is a `JobSales` data model
type JobSales struct {
	JobSalesID        int64   `db:"JobSalesID" json:"JobSalesID"`
	JobID             int64   `db:"JobID" json:"JobID"`
	UserID            int64   `db:"UserID" json:"UserID"`
	CommissionPercent float64 `db:"CommissionPercent" json:"CommissionPercent"`
	DateCreated       int64   `db:"DateCreated" json:"DateCreated"`
	IsDeleted         int     `db:"IsDeleted" json:"IsDeleted"`
	CommissionDollars float64 `db:"CommissionDollars" json:"CommissionDollars"`
	IsHouse           int     `db:"IsHouse" json:"IsHouse"`
}

// Comment_TableName is the name of the table
func (c *JobSales) Table_Name() query.TableName {
	return JobSales_TableName
}

func (c *JobSales) Table_Columns() []query.Column {
	return []query.Column{
		JobSales_Column_JobSalesID,
		JobSales_Column_JobID,
		JobSales_Column_UserID,
		JobSales_Column_CommissionPercent,
		JobSales_Column_DateCreated,
		JobSales_Column_IsDeleted,
		JobSales_Column_CommissionDollars,
		JobSales_Column_IsHouse,
	}
}

func (c *JobSales) Table_Column_Types() map[query.Column]string {
	return map[query.Column]string{
		"JobSalesID":        "%d",
		"JobID":             "%d",
		"UserID":            "%d",
		"CommissionPercent": "%d",
		"DateCreated":       "%d",
		"IsDeleted":         "%d",
		"CommissionDollars": "%d",
		"IsHouse":           "%d",
	}
}

func (c *JobSales) Table_Column_Values() map[query.Column]interface{} {
	return map[query.Column]interface{}{
		"JobSalesID":        c.JobSalesID,
		"JobID":             c.JobID,
		"UserID":            c.UserID,
		"CommissionPercent": c.CommissionPercent,
		"DateCreated":       c.DateCreated,
		"IsDeleted":         c.IsDeleted,
		"CommissionDollars": c.CommissionDollars,
		"IsHouse":           c.IsHouse,
	}
}

// Comment_PrimaryKey is the name of the table's primary key
func (c *JobSales) Table_PrimaryKey() query.Column {
	return "JobSalesID"
}

func (c *JobSales) Table_PrimaryKey_Value() int64 {
	return c.JobSalesID
}

// Comment_InsertColumns is a list of all insert columns for this model
func (c *JobSales) Table_InsertColumns() []query.Column {
	return []query.Column{"DateCreated", "Content", "ObjectType", "ObjectID"}
}

// Comment_UpdateColumns is a list of all update columns for this model
func (c *JobSales) Table_UpdateColumns() []query.Column {
	return []query.Column{"Content", "ObjectType", "ObjectID"}
}

func (c *JobSales) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func (c *JobSales) Destroy() string {
	return ""
}

func (c *JobSales) Update() string {
	return ""
}

func (c *JobSales) Create() string {
	return ""
}

func (c *JobSales) FromID(id int64) string {
	return ""
}
