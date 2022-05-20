package gen

import (
	"testing"

	"github.com/macinnir/dvc/core/lib/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fooObj = &schema.Table{
	Name: "Foo",
	Columns: map[string]*schema.Column{
		"FooID": {Name: "FooID", DataType: "int"},
		"Bar":   {Name: "Bar", DataType: "varchar", MaxLength: 200},
		"Baz":   {Name: "Baz", DataType: "datetime"},
	},
}

var fooColumns = map[string]string{
	"FooID": "int",
	"Bar":   "string",
	"Baz":   "string",
}

var tsFoo = `/**
 * Generated Code; DO NOT EDIT
 *
 * Foo
 */

export type Foo = {

	// Bar string
	Bar: string;

	// Baz string
	Baz: string;

	// FooID int
	FooID: number;

}

// newFoo is a factory method for creating new Foo objects
export const newFoo = () : Foo => ({ 
	Bar: '',
	Baz: '',
	FooID: 0,
});

`

func TestGenerateTypescriptType(t *testing.T) {

	g, e := GenerateTypescriptModel("Foo", fooColumns)

	assert.Nil(t, e)
	assert.Equal(t, tsFoo, g)
}

var barObj = &schema.Table{
	Name: "Bar",
	Columns: map[string]*schema.Column{
		"FooID2": {Name: "FooID2", DataType: "int"},
		"Bar2":   {Name: "Bar2", DataType: "varchar", MaxLength: 200},
		"Baz2":   {Name: "Baz2", DataType: "datetime"},
	},
}

var tsBar = `/**
* Bar
*/
export type Bar = {

   // Bar2 varchar(200)
   Bar2: string;

   // Baz2 datetime
   Baz2: string;

   // FooID2 int
   FooID2: number;

}

// newFoo is a factory method for new Foo objects
export const newBar = () : Bar => ({
   Bar2: '',
   Baz2: '',
   FooID2: 0,
});

`

func TestGenerateTypescriptDTO(t *testing.T) {
	str, _ := GenerateTypescriptDTO("CreateWebsiteDTO", map[string]string{
		"Title":                   "string",
		"BaseURL":                 "string",
		"ScheduleIntervalMinutes": "int64",
		"IsActive":                "int",
		"IDs":                     "[]int64",
	})
	expected := `/**
 * Generated Code; DO NOT EDIT
 *
 * CreateWebsiteDTO
 */

export type CreateWebsiteDTO = {

	// BaseURL string
	BaseURL: string;

	// IDs []int64
	IDs: number[];

	// IsActive int
	IsActive: number;

	// ScheduleIntervalMinutes int64
	ScheduleIntervalMinutes: number;

	// Title string
	Title: string;

}

// newCreateWebsiteDTO is a factory method for new CreateWebsiteDTO objects
export const newCreateWebsiteDTO = () : CreateWebsiteDTO => ({
	BaseURL: '',
	IDs: [],
	IsActive: 0,
	ScheduleIntervalMinutes: 0,
	Title: '',
});

`

	assert.Equal(t, expected, str)
}

func TestGenerateTypescriptDTO_WithModelsAndDTOs(t *testing.T) {
	var str, e = GenerateTypescriptDTO("UpdateJobDTO", map[string]string{
		"AwardDate":           "string",
		"BillingAddress":      "string",
		"BillingAddressCity":  "string",
		"BillingAddressState": "string",
		"BillingAddressZip":   "string",
		"CustomerPO1Number":   "string",
		"CustomerPO1SentTo":   "int64",
		"CustomerPO2Number":   "string",
		"CustomerPO2SentTo":   "int64",
		"GrossMarginPercent":  "float64",
		"GrossProfit":         "float64",
		"IsAddFreight":        "int",
		"IsThirdPartySplit":   "int",
		"Items":               "[]*UpdateJobItemDTO",
		"JEFDate":             "string",
		"Notes":               "string",
		"ProjectCost":         "float64",
		"Sales":               "[]*models.JobSales",
	})

	var expected = `/**
 * Generated Code; DO NOT EDIT
 *
 * UpdateJobDTO
 */
import { UpdateJobItemDTO } from './UpdateJobItemDTO';
import { JobSales } from 'gen/models/JobSales';

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

	// Items []*UpdateJobItemDTO
	Items: UpdateJobItemDTO[];

	// JEFDate string
	JEFDate: string;

	// Notes string
	Notes: string;

	// ProjectCost float64
	ProjectCost: number;

	// Sales []*models.JobSales
	Sales: JobSales[];

}

// newUpdateJobDTO is a factory method for new UpdateJobDTO objects
export const newUpdateJobDTO = () : UpdateJobDTO => ({
	AwardDate: '',
	BillingAddress: '',
	BillingAddressCity: '',
	BillingAddressState: '',
	BillingAddressZip: '',
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
	Sales: [],
});

`

	require.Nil(t, e)
	assert.Equal(t, expected, str)
}
