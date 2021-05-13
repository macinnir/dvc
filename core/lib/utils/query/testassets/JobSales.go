package testassets

import "encoding/json"

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
func (c *JobSales) Table_Name() string {
	return "JobSales"
}

func (c *JobSales) Table_Columns() []string {
	return []string{
		"JobSalesID",
		"JobID",
		"UserID",
		"CommissionPercent",
		"DateCreated",
		"IsDeleted",
		"CommissionDollars",
		"IsHouse",
	}
}

func (c *JobSales) Table_Column_Types() map[string]string {
	return map[string]string{
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

func (c *JobSales) Table_Column_Values() map[string]interface{} {
	return map[string]interface{}{
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
func (c *JobSales) Table_PrimaryKey() string {
	return "JobSalesID"
}

func (c *JobSales) Table_PrimaryKey_Value() int64 {
	return c.JobSalesID
}

// Comment_InsertColumns is a list of all insert columns for this model
func (c *JobSales) Table_InsertColumns() []string {
	return []string{"DateCreated", "Content", "ObjectType", "ObjectID"}
}

// Comment_UpdateColumns is a list of all update columns for this model
func (c *JobSales) Table_UpdateColumns() []string {
	return []string{"Content", "ObjectType", "ObjectID"}
}

func (c *JobSales) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func (c *JobSales) Destroy() string {
	return ""
}

func (c *JobSales) Save() string {
	return ""
}
