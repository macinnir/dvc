package testassets

import (
	"encoding/json"

	"gopkg.in/guregu/null.v3"
)

// Job is a `Job` data model
type Job struct {
	JobID                  int64       `db:"JobID" json:"JobID"`
	DateCreated            int64       `db:"DateCreated" json:"DateCreated"`
	LastUpdated            int64       `db:"LastUpdated" json:"LastUpdated"`
	IsDeleted              int         `db:"IsDeleted" json:"IsDeleted"`
	QuoteNumberID          int64       `db:"QuoteNumberID" json:"QuoteNumberID"`
	CustomerID             int64       `db:"CustomerID" json:"CustomerID"`
	AwardDate              int64       `db:"AwardDate" json:"AwardDate"`
	AwardDateString        string      `db:"AwardDateString" json:"AwardDateString"`
	QuoteID                int64       `db:"QuoteID" json:"QuoteID"`
	Description            string      `db:"Description" json:"Description"`
	CustomerContactID      int64       `db:"CustomerContactID" json:"CustomerContactID"`
	Notes                  null.String `db:"Notes" json:"Notes"`
	CustomerPO2SentTo      int64       `db:"CustomerPO2SentTo" json:"CustomerPO2SentTo"`
	ThirdPartyName         string      `db:"ThirdPartyName" json:"ThirdPartyName"`
	CustomerPO2Number      string      `db:"CustomerPO2Number" json:"CustomerPO2Number"`
	BillingAddressZip      string      `db:"BillingAddressZip" json:"BillingAddressZip"`
	BillingAddress         string      `db:"BillingAddress" json:"BillingAddress"`
	BillingAddressState    string      `db:"BillingAddressState" json:"BillingAddressState"`
	CustomerPO1SentTo      int64       `db:"CustomerPO1SentTo" json:"CustomerPO1SentTo"`
	ProjectCost            float64     `db:"ProjectCost" json:"ProjectCost"`
	CustomerPO1Number      string      `db:"CustomerPO1Number" json:"CustomerPO1Number"`
	ThirdPartySplitPercent float64     `db:"ThirdPartySplitPercent" json:"ThirdPartySplitPercent"`
	CommissionTypeID       int64       `db:"CommissionTypeID" json:"CommissionTypeID"`
	GrossMarginPercent     float64     `db:"GrossMarginPercent" json:"GrossMarginPercent"`
	TotalPrice             float64     `db:"TotalPrice" json:"TotalPrice"`
	IsThirdPartySplit      int         `db:"IsThirdPartySplit" json:"IsThirdPartySplit"`
	GrossProfit            float64     `db:"GrossProfit" json:"GrossProfit"`
	BillingAddressCity     string      `db:"BillingAddressCity" json:"BillingAddressCity"`
	RemainingGrossProfit   float64     `db:"RemainingGrossProfit" json:"RemainingGrossProfit"`
	IsAddFreight           int         `db:"IsAddFreight" json:"IsAddFreight"`
	JobNumberString        string      `db:"JobNumberString" json:"JobNumberString"`
	ThirdPartyCommission   float64     `db:"ThirdPartyCommission" json:"ThirdPartyCommission"`
	ShippingAddressZip     string      `db:"ShippingAddressZip" json:"ShippingAddressZip"`
	ShippingAddress        string      `db:"ShippingAddress" json:"ShippingAddress"`
	ShippingAddressCity    string      `db:"ShippingAddressCity" json:"ShippingAddressCity"`
	ShippingAddressState   string      `db:"ShippingAddressState" json:"ShippingAddressState"`
	Sales1                 string      `db:"Sales1" json:"Sales1"`
	BidTypeID              int64       `db:"BidTypeID" json:"BidTypeID"`
	Vendor1ID              int64       `db:"Vendor1ID" json:"Vendor1ID"`
	MarketID               int64       `db:"MarketID" json:"MarketID"`
	Vendor2ID              int64       `db:"Vendor2ID" json:"Vendor2ID"`
	Sales2                 string      `db:"Sales2" json:"Sales2"`
	JEFDate                int64       `db:"JEFDate" json:"JEFDate"`
	JEFDateString          string      `db:"JEFDateString" json:"JEFDateString"`
}

// Comment_TableName is the name of the table
func (c *Job) Table_Name() string {
	return "Job"
}

func (c *Job) Table_Columns() []string {
	return []string{
		"JobID",
		"DateCreated",
		"LastUpdated",
		"IsDeleted",
		"QuoteNumberID",
		"CustomerID",
		"AwardDate",
		"AwardDateString",
		"QuoteID",
		"Description",
		"CustomerContactID",
		"Notes",
		"CustomerPO2SentTo",
		"ThirdPartyName",
		"CustomerPO2Number",
		"BillingAddressZip",
		"BillingAddress",
		"BillingAddressState",
		"CustomerPO1SentTo",
		"ProjectCost",
		"CustomerPO1Number",
		"ThirdPartySplitPercent",
		"CommissionTypeID",
		"GrossMarginPercent",
		"TotalPrice",
		"IsThirdPartySplit",
		"GrossProfit",
		"BillingAddressCity",
		"RemainingGrossProfit",
		"IsAddFreight",
		"JobNumberString",
		"ThirdPartyCommission",
		"ShippingAddressZip",
		"ShippingAddress",
		"ShippingAddressCity",
		"ShippingAddressState",
		"Sales1",
		"BidTypeID",
		"Vendor1ID",
		"MarketID",
		"Vendor2ID",
		"Sales2",
		"JEFDate",
		"JEFDateString",
	}
}

func (c *Job) Table_Column_Types() map[string]string {
	return map[string]string{
		"JobID":                  "%d",
		"DateCreated":            "%d",
		"LastUpdated":            "%d",
		"IsDeleted":              "%d",
		"QuoteNumberID":          "%d",
		"CustomerID":             "%d",
		"AwardDate":              "%d",
		"AwardDateString":        "%s",
		"QuoteID":                "%d",
		"Description":            "%s",
		"CustomerContactID":      "%d",
		"Notes":                  "%s",
		"ThirdPartyName":         "%s",
		"CustomerPO2Number":      "%s",
		"BillingAddressZip":      "%s",
		"BillingAddress":         "%s",
		"BillingAddressState":    "%s",
		"CustomerPO1Number":      "%s",
		"BillingAddressCity":     "%s",
		"JobNumberString":        "%s",
		"ShippingAddressZip":     "%s",
		"ShippingAddress":        "%s",
		"ShippingAddressCity":    "%s",
		"ShippingAddressState":   "%s",
		"Sales1":                 "%s",
		"ProjectCost":            "%f",
		"ThirdPartySplitPercent": "%f",
		"GrossMarginPercent":     "%f",
		"TotalPrice":             "%f",
		"GrossProfit":            "%f",
		"RemainingGrossProfit":   "%f",
		"ThirdPartyCommission":   "%f",
		"CustomerPO2SentTo":      "%d",
		"CustomerPO1SentTo":      "%d",
		"CommissionTypeID":       "%d",
		"IsAddFreight":           "%d",
		"IsThirdPartySplit":      "%d",
		"BidTypeID":              "%d",
		"Vendor1ID":              "%d",
		"MarketID":               "%d",
		"Vendor2ID":              "%d",
		"JEFDate":                "%d",
		"Sales2":                 "%s",
		"JEFDateString":          "%s",
	}
}

func (c *Job) Table_Column_Values() map[string]interface{} {
	return map[string]interface{}{
		"JobID":                  c.JobID,
		"DateCreated":            c.DateCreated,
		"LastUpdated":            c.LastUpdated,
		"IsDeleted":              c.IsDeleted,
		"QuoteNumberID":          c.QuoteNumberID,
		"CustomerID":             c.CustomerID,
		"AwardDate":              c.AwardDate,
		"AwardDateString":        c.AwardDateString,
		"QuoteID":                c.QuoteID,
		"Description":            c.Description,
		"CustomerContactID":      c.CustomerContactID,
		"Notes":                  c.Notes,
		"ThirdPartyName":         c.ThirdPartyName,
		"CustomerPO2Number":      c.CustomerPO2Number,
		"BillingAddressZip":      c.BillingAddressZip,
		"BillingAddress":         c.BillingAddress,
		"BillingAddressState":    c.BillingAddressState,
		"CustomerPO1Number":      c.CustomerPO1Number,
		"BillingAddressCity":     c.BillingAddressCity,
		"JobNumberString":        c.JobNumberString,
		"ShippingAddressZip":     c.ShippingAddressZip,
		"ShippingAddress":        c.ShippingAddress,
		"ShippingAddressCity":    c.ShippingAddressCity,
		"ShippingAddressState":   c.ShippingAddressState,
		"Sales1":                 c.Sales1,
		"ProjectCost":            c.ProjectCost,
		"ThirdPartySplitPercent": c.ThirdPartySplitPercent,
		"GrossMarginPercent":     c.GrossMarginPercent,
		"TotalPrice":             c.TotalPrice,
		"GrossProfit":            c.GrossProfit,
		"RemainingGrossProfit":   c.RemainingGrossProfit,
		"ThirdPartyCommission":   c.ThirdPartyCommission,
		"CustomerPO2SentTo":      c.CustomerPO2SentTo,
		"CustomerPO1SentTo":      c.CustomerPO1SentTo,
		"CommissionTypeID":       c.CommissionTypeID,
		"IsAddFreight":           c.IsAddFreight,
		"IsThirdPartySplit":      c.IsThirdPartySplit,
		"BidTypeID":              c.BidTypeID,
		"Vendor1ID":              c.Vendor1ID,
		"MarketID":               c.MarketID,
		"Vendor2ID":              c.Vendor2ID,
		"JEFDate":                c.JEFDate,
		"Sales2":                 c.Sales2,
		"JEFDateString":          c.JEFDateString,
	}
}

// Comment_PrimaryKey is the name of the table's primary key
func (c *Job) Table_PrimaryKey() string {
	return "JobID"
}

func (c *Job) Table_PrimaryKey_Value() int64 {
	return c.JobID
}

// Comment_InsertColumns is a list of all insert columns for this model
func (c *Job) Table_InsertColumns() []string {
	return []string{"DateCreated", "Content", "ObjectType", "ObjectID"}
}

// Comment_UpdateColumns is a list of all update columns for this model
func (c *Job) Table_UpdateColumns() []string {
	return []string{"Content", "ObjectType", "ObjectID"}
}

func (c *Job) String() string {
	str, _ := json.Marshal(c)
	return string(str)
}

func (c *Job) Destroy() string {
	return ""
}

func (c *Job) Create() string {
	return ""
}

func (c *Job) Update() string {
	return ""
}

func (c *Job) FromID(id int64) string {
	return ""
}
