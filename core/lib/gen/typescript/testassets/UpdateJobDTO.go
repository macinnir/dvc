package dtos

import "joc-rfq-api/gen/definitions/models"

// UpdateJobDTO is the shape of an UpdateJob request body
type UpdateJobDTO struct {
	TotalPrice             float64             `json:"TotalPrice"`
	Notes                  string              `json:"Notes"`
	AwardDate              string              `json:"AwardDate"`
	JEFDate                string              `json:"JEFDate"`
	ProjectCost            float64             `json:"ProjectCost"`
	GrossProfit            float64             `json:"GrossProfit"`
	GrossMarginPercent     float64             `json:"GrossMarginPercent"`
	RemainingGrossProfit   float64             `json:"RemainingGrossProfit"`
	ThirdPartyCommission   float64             `json:"ThirdPartyCommission"`
	ThirdPartySplitPercent float64             `json:"ThirdPartySplitPercent"`
	BillingAddress         string              `json:"BillingAddress"`
	BillingAddressCity     string              `json:"BillingAddressCity"`
	BillingAddressState    string              `json:"BillingAddressState"`
	BillingAddressZip      string              `json:"BillingAddressZip"`
	ShippingAddress        string              `json:"ShippingAddress"`
	ShippingAddressCity    string              `json:"ShippingAddressCity"`
	ShippingAddressState   string              `json:"ShippingAddressState"`
	ShippingAddressZip     string              `json:"ShippingAddressZip"`
	CommissionTypeID       int64               `json:"CommissionTypeID"`
	CustomerPO1Number      string              `json:"CustomerPO1Number"`
	CustomerPO1SentTo      int64               `json:"CustomerPO1SentTo"`
	CustomerPO2Number      string              `json:"CustomerPO2Number"`
	CustomerPO2SentTo      int64               `json:"CustomerPO2SentTo"`
	IsAddFreight           int                 `json:"IsAddFreight"`
	IsThirdPartySplit      int                 `json:"IsThirdPartySplit"`
	ThirdPartyName         string              `json:"ThirdPartyName"`
	Items                  []*UpdateJobItemDTO `json:"Items"`
	Sales                  []*models.JobSales  `json:"Sales"`
}
