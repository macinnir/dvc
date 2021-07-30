package testassets

import (
	"encoding/json"

	"github.com/macinnir/dvc/core/lib/utils/db"
	"github.com/macinnir/dvc/core/lib/utils/query"
	"gopkg.in/guregu/null.v3"
)

const (
	Job_TableName query.TableName = "Job"

	// Columns
	Job_Column_JobID                  query.Column = "JobID"
	Job_Column_DateCreated            query.Column = "DateCreated"
	Job_Column_LastUpdated            query.Column = "LastUpdated"
	Job_Column_IsDeleted              query.Column = "IsDeleted"
	Job_Column_QuoteNumberID          query.Column = "QuoteNumberID"
	Job_Column_CustomerID             query.Column = "CustomerID"
	Job_Column_AwardDate              query.Column = "AwardDate"
	Job_Column_AwardDateString        query.Column = "AwardDateString"
	Job_Column_QuoteID                query.Column = "QuoteID"
	Job_Column_Description            query.Column = "Description"
	Job_Column_CustomerContactID      query.Column = "CustomerContactID"
	Job_Column_Notes                  query.Column = "Notes"
	Job_Column_CustomerPO2SentTo      query.Column = "CustomerPO2SentTo"
	Job_Column_ThirdPartyName         query.Column = "ThirdPartyName"
	Job_Column_CustomerPO2Number      query.Column = "CustomerPO2Number"
	Job_Column_BillingAddressZip      query.Column = "BillingAddressZip"
	Job_Column_BillingAddress         query.Column = "BillingAddress"
	Job_Column_BillingAddressState    query.Column = "BillingAddressState"
	Job_Column_CustomerPO1SentTo      query.Column = "CustomerPO1SentTo"
	Job_Column_ProjectCost            query.Column = "ProjectCost"
	Job_Column_CustomerPO1Number      query.Column = "CustomerPO1Number"
	Job_Column_ThirdPartySplitPercent query.Column = "ThirdPartySplitPercent"
	Job_Column_CommissionTypeID       query.Column = "CommissionTypeID"
	Job_Column_GrossMarginPercent     query.Column = "GrossMarginPercent"
	Job_Column_TotalPrice             query.Column = "TotalPrice"
	Job_Column_IsThirdPartySplit      query.Column = "IsThirdPartySplit"
	Job_Column_GrossProfit            query.Column = "GrossProfit"
	Job_Column_BillingAddressCity     query.Column = "BillingAddressCity"
	Job_Column_RemainingGrossProfit   query.Column = "RemainingGrossProfit"
	Job_Column_IsAddFreight           query.Column = "IsAddFreight"
	Job_Column_JobNumberString        query.Column = "JobNumberString"
	Job_Column_ThirdPartyCommission   query.Column = "ThirdPartyCommission"
	Job_Column_ShippingAddressZip     query.Column = "ShippingAddressZip"
	Job_Column_ShippingAddress        query.Column = "ShippingAddress"
	Job_Column_ShippingAddressCity    query.Column = "ShippingAddressCity"
	Job_Column_ShippingAddressState   query.Column = "ShippingAddressState"
	Job_Column_Sales1                 query.Column = "Sales1"
	Job_Column_BidTypeID              query.Column = "BidTypeID"
	Job_Column_Vendor1ID              query.Column = "Vendor1ID"
	Job_Column_MarketID               query.Column = "MarketID"
	Job_Column_Vendor2ID              query.Column = "Vendor2ID"
	Job_Column_Sales2                 query.Column = "Sales2"
	Job_Column_JEFDate                query.Column = "JEFDate"
	Job_Column_JEFDateString          query.Column = "JEFDateString"
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
func (c *Job) Table_Name() query.TableName {
	return Job_TableName
}

func (c *Job) Table_Columns() []query.Column {
	return []query.Column{
		Job_Column_JobID,
		Job_Column_DateCreated,
		Job_Column_LastUpdated,
		Job_Column_IsDeleted,
		Job_Column_QuoteNumberID,
		Job_Column_CustomerID,
		Job_Column_AwardDate,
		Job_Column_AwardDateString,
		Job_Column_QuoteID,
		Job_Column_Description,
		Job_Column_CustomerContactID,
		Job_Column_Notes,
		Job_Column_CustomerPO2SentTo,
		Job_Column_ThirdPartyName,
		Job_Column_CustomerPO2Number,
		Job_Column_BillingAddressZip,
		Job_Column_BillingAddress,
		Job_Column_BillingAddressState,
		Job_Column_CustomerPO1SentTo,
		Job_Column_ProjectCost,
		Job_Column_CustomerPO1Number,
		Job_Column_ThirdPartySplitPercent,
		Job_Column_CommissionTypeID,
		Job_Column_GrossMarginPercent,
		Job_Column_TotalPrice,
		Job_Column_IsThirdPartySplit,
		Job_Column_GrossProfit,
		Job_Column_BillingAddressCity,
		Job_Column_RemainingGrossProfit,
		Job_Column_IsAddFreight,
		Job_Column_JobNumberString,
		Job_Column_ThirdPartyCommission,
		Job_Column_ShippingAddressZip,
		Job_Column_ShippingAddress,
		Job_Column_ShippingAddressCity,
		Job_Column_ShippingAddressState,
		Job_Column_Sales1,
		Job_Column_BidTypeID,
		Job_Column_Vendor1ID,
		Job_Column_MarketID,
		Job_Column_Vendor2ID,
		Job_Column_Sales2,
		Job_Column_JEFDate,
		Job_Column_JEFDateString,
	}
}

func (c *Job) Table_Column_Types() map[query.Column]string {
	return map[query.Column]string{
		Job_Column_JobID:                  "%d",
		Job_Column_DateCreated:            "%d",
		Job_Column_LastUpdated:            "%d",
		Job_Column_IsDeleted:              "%d",
		Job_Column_QuoteNumberID:          "%d",
		Job_Column_CustomerID:             "%d",
		Job_Column_AwardDate:              "%d",
		Job_Column_AwardDateString:        "%s",
		Job_Column_QuoteID:                "%d",
		Job_Column_Description:            "%s",
		Job_Column_CustomerContactID:      "%d",
		Job_Column_Notes:                  "%s",
		Job_Column_ThirdPartyName:         "%s",
		Job_Column_CustomerPO2Number:      "%s",
		Job_Column_BillingAddressZip:      "%s",
		Job_Column_BillingAddress:         "%s",
		Job_Column_BillingAddressState:    "%s",
		Job_Column_CustomerPO1Number:      "%s",
		Job_Column_BillingAddressCity:     "%s",
		Job_Column_JobNumberString:        "%s",
		Job_Column_ShippingAddressZip:     "%s",
		Job_Column_ShippingAddress:        "%s",
		Job_Column_ShippingAddressCity:    "%s",
		Job_Column_ShippingAddressState:   "%s",
		Job_Column_Sales1:                 "%s",
		Job_Column_ProjectCost:            "%f",
		Job_Column_ThirdPartySplitPercent: "%f",
		Job_Column_GrossMarginPercent:     "%f",
		Job_Column_TotalPrice:             "%f",
		Job_Column_GrossProfit:            "%f",
		Job_Column_RemainingGrossProfit:   "%f",
		Job_Column_ThirdPartyCommission:   "%f",
		Job_Column_CustomerPO2SentTo:      "%d",
		Job_Column_CustomerPO1SentTo:      "%d",
		Job_Column_CommissionTypeID:       "%d",
		Job_Column_IsAddFreight:           "%d",
		Job_Column_IsThirdPartySplit:      "%d",
		Job_Column_BidTypeID:              "%d",
		Job_Column_Vendor1ID:              "%d",
		Job_Column_MarketID:               "%d",
		Job_Column_Vendor2ID:              "%d",
		Job_Column_JEFDate:                "%d",
		Job_Column_Sales2:                 "%s",
		Job_Column_JEFDateString:          "%s",
	}
}

func (c *Job) Table_Column_Values() map[query.Column]interface{} {
	return map[query.Column]interface{}{
		Job_Column_JobID:                  c.JobID,
		Job_Column_DateCreated:            c.DateCreated,
		Job_Column_LastUpdated:            c.LastUpdated,
		Job_Column_IsDeleted:              c.IsDeleted,
		Job_Column_QuoteNumberID:          c.QuoteNumberID,
		Job_Column_CustomerID:             c.CustomerID,
		Job_Column_AwardDate:              c.AwardDate,
		Job_Column_AwardDateString:        c.AwardDateString,
		Job_Column_QuoteID:                c.QuoteID,
		Job_Column_Description:            c.Description,
		Job_Column_CustomerContactID:      c.CustomerContactID,
		Job_Column_Notes:                  c.Notes,
		Job_Column_ThirdPartyName:         c.ThirdPartyName,
		Job_Column_CustomerPO2Number:      c.CustomerPO2Number,
		Job_Column_BillingAddressZip:      c.BillingAddressZip,
		Job_Column_BillingAddress:         c.BillingAddress,
		Job_Column_BillingAddressState:    c.BillingAddressState,
		Job_Column_CustomerPO1Number:      c.CustomerPO1Number,
		Job_Column_BillingAddressCity:     c.BillingAddressCity,
		Job_Column_JobNumberString:        c.JobNumberString,
		Job_Column_ShippingAddressZip:     c.ShippingAddressZip,
		Job_Column_ShippingAddress:        c.ShippingAddress,
		Job_Column_ShippingAddressCity:    c.ShippingAddressCity,
		Job_Column_ShippingAddressState:   c.ShippingAddressState,
		Job_Column_Sales1:                 c.Sales1,
		Job_Column_ProjectCost:            c.ProjectCost,
		Job_Column_ThirdPartySplitPercent: c.ThirdPartySplitPercent,
		Job_Column_GrossMarginPercent:     c.GrossMarginPercent,
		Job_Column_TotalPrice:             c.TotalPrice,
		Job_Column_GrossProfit:            c.GrossProfit,
		Job_Column_RemainingGrossProfit:   c.RemainingGrossProfit,
		Job_Column_ThirdPartyCommission:   c.ThirdPartyCommission,
		Job_Column_CustomerPO2SentTo:      c.CustomerPO2SentTo,
		Job_Column_CustomerPO1SentTo:      c.CustomerPO1SentTo,
		Job_Column_CommissionTypeID:       c.CommissionTypeID,
		Job_Column_IsAddFreight:           c.IsAddFreight,
		Job_Column_IsThirdPartySplit:      c.IsThirdPartySplit,
		Job_Column_BidTypeID:              c.BidTypeID,
		Job_Column_Vendor1ID:              c.Vendor1ID,
		Job_Column_MarketID:               c.MarketID,
		Job_Column_Vendor2ID:              c.Vendor2ID,
		Job_Column_JEFDate:                c.JEFDate,
		Job_Column_Sales2:                 c.Sales2,
		Job_Column_JEFDateString:          c.JEFDateString,
	}
}

// Comment_PrimaryKey is the name of the table's primary key
func (c *Job) Table_PrimaryKey() query.Column {
	return Job_Column_JobID
}

func (c *Job) Table_PrimaryKey_Value() int64 {
	return c.JobID
}

// Comment_InsertColumns is a list of all insert columns for this model
// TODO
func (c *Job) Table_InsertColumns() []query.Column {
	return []query.Column{}
}

// Comment_UpdateColumns is a list of all update columns for this model
// TODO
func (c *Job) Table_UpdateColumns() []query.Column {
	return []query.Column{}
}

func (c *Job) String() string {
	str, _ := json.Marshal(c)
	return string(str)
}

func (c *Job) Create(db db.IDB) error {
	return nil
}

func (c *Job) Update(db db.IDB) error {
	return nil
}

func (c *Job) Delete(db db.IDB) error {
	return nil
}

func (c *Job) FromID(db db.IDB, id int64) (query.IModel, error) {
	return nil, nil
}
