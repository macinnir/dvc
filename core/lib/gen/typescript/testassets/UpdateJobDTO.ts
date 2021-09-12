import { JobSales } from "gen/models/JobSales";
import { UpdateJobItemDTO } from "./UpdateJobItemDTO";

/**
 * Generated Code; DO NOT EDIT
 *
 * UpdateJobDTO
 */
export type UpdateJobDTO = {

	// AwardDate string
	AwardDate: string;

	// BillingAddress string
	BillingAddress: string;

	// BillingAddressCity string
	BillingAddressCity: string;

	// BillingAddressState string
	BillingAddressState: string;

	// BillingAddressZip string
	BillingAddressZip: string;

	// CommissionTypeID int64
	CommissionTypeID: number;

	// CustomerPO1Number string
	CustomerPO1Number: string;

	// CustomerPO1SentTo int64
	CustomerPO1SentTo: number;

	// CustomerPO2Number string
	CustomerPO2Number: string;

	// CustomerPO2SentTo int64
	CustomerPO2SentTo: number;

	// GrossMarginPercent float64
	GrossMarginPercent: number;

	// GrossProfit float64
	GrossProfit: number;

	// IsAddFreight int
	IsAddFreight: number;

	// IsThirdPartySplit int
	IsThirdPartySplit: number;

	// Items []&{1975 UpdateJobItemDTO}
	Items: UpdateJobItemDTO[];

	// JEFDate string
	JEFDate: string;

	// Notes string
	Notes: string;

	// ProjectCost float64
	ProjectCost: number;

	// RemainingGrossProfit float64
	RemainingGrossProfit: number;

	// Sales []&{2034 0xc00031a000}
	Sales: JobSales[];

	// ShippingAddress string
	ShippingAddress: string;

	// ShippingAddressCity string
	ShippingAddressCity: string;

	// ShippingAddressState string
	ShippingAddressState: string;

	// ShippingAddressZip string
	ShippingAddressZip: string;

	// ThirdPartyCommission float64
	ThirdPartyCommission: number;

	// ThirdPartyName string
	ThirdPartyName: string;

	// ThirdPartySplitPercent float64
	ThirdPartySplitPercent: number;

	// TotalPrice float64
	TotalPrice: number;

}

// newUpdateJobDTO is a factory method for new UpdateJobDTO objects
export const newUpdateJobDTO = () : UpdateJobDTO => ({
	AwardDate: '',
	BillingAddress: '',
	BillingAddressCity: '',
	BillingAddressState: '',
	BillingAddressZip: '',
	CommissionTypeID: 0,
	CustomerPO1Number: '',
	CustomerPO1SentTo: 0,
	CustomerPO2Number: '',
	CustomerPO2SentTo: 0,
	GrossMarginPercent: 0,
	GrossProfit: 0,
	IsAddFreight: 0,
	IsThirdPartySplit: 0,
	Items: [],
	JEFDate: '',
	Notes: '',
	ProjectCost: 0,
	RemainingGrossProfit: 0,
	Sales: [],
	ShippingAddress: '',
	ShippingAddressCity: '',
	ShippingAddressState: '',
	ShippingAddressZip: '',
	ThirdPartyCommission: 0,
	ThirdPartyName: '',
	ThirdPartySplitPercent: 0,
	TotalPrice: 0,
});

