package gen

import (
	"testing"

	"github.com/macinnir/dvc/core/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenTSRoutes(t *testing.T) {

	ctl := &lib.Controller{
		Name:        "TestController",
		Description: "does a thing",
		Package:     "foo",
		Routes: []*lib.ControllerRoute{
			{
				Name:           "ActivateSMSNotifications",
				Description:    "activates sms notifications for the current user",
				Path:           "/self/notifications/sms",
				Method:         "POST",
				Params:         []lib.ControllerRouteParam{},
				Queries:        []lib.ControllerRouteQuery{},
				IsAuth:         true,
				BodyType:       "dtos.ActivateSMSNotificationsDTO",
				BodyFormat:     "JSON",
				HasBody:        true,
				ResponseType:   "*aggregates.UserAggregate",
				ResponseFormat: "JSON",
				ResponseCode:   200,
				Permission:     "",
				ControllerName: "Self",
				FileName:       "core/api/base/SelfController.go",
				LineNo:         0,
			},
		},
	}

	var e error
	assert.NotNil(t, ctl)
	// _, e := genTSRoutesFromController(ctl)

	require.Nil(t, e)

}

func TestGenTSRoute(t *testing.T) {

	tests := []struct {
		route    *lib.ControllerRoute
		expected string
	}{
		{
			&lib.ControllerRoute{
				Name:           "CreateUser",
				Description:    "Creates a user",
				Path:           "/users",
				Method:         "POST",
				Params:         []lib.ControllerRouteParam{},
				Queries:        []lib.ControllerRouteQuery{},
				ResponseType:   "*aggregates.UserAggregate",
				ResponseFormat: "JSON",
				BodyType:       "dtos.CreateUserDTO",
				BodyFormat:     "JSON",
				HasBody:        true,
			},
			"export const createUser = async (body : CreateUserDTO) => await axios.post<UserAggregate>(`/users`, body);",
		},
		{
			&lib.ControllerRoute{
				Name:        "CreateVendorItem",
				Description: "creates a vendor item",
				Path:        "/vendors/{vendorID:[0-9]+}/items",
				Method:      "POST",
				Params: []lib.ControllerRouteParam{
					{
						Name:    "vendorID",
						Pattern: "[0-9]+",
						Type:    "int64",
					},
				},
				Queries:        []lib.ControllerRouteQuery{},
				IsAuth:         true,
				BodyType:       "appdtos.CreateVendorItemDTO",
				BodyFormat:     "JSON",
				HasBody:        true,
				ResponseType:   "*models.VendorItem",
				ResponseFormat: "JSON",
				ResponseCode:   200,
				Permission:     "Vendors_CreateVendorItem",
				ControllerName: "Vendors",
				FileName:       "app/api/joc/VendorsController.go",
				LineNo:         0,
			},
			"export const createVendorItem = async (vendorID : number, body : CreateVendorItemDTO) => await axios.post<VendorItem>(`/vendors/${vendorID}/items`, body);",
		},
		{
			&lib.ControllerRoute{
				Name:        "DownloadSalesActivityReportExcel",
				Description: "downloads a SalesActivityReport excel document",
				Path:        "/reports/sales-activity/download/excel?year={year:[0-9]+}",
				Method:      "GET",
				Params:      []lib.ControllerRouteParam{},
				Queries: []lib.ControllerRouteQuery{
					{
						Name:         "year",
						Pattern:      "[0-9]+",
						Type:         "int64",
						VariableName: "year",
						ValueRaw:     "{year:[0-9]+}",
					},
				},
				IsAuth:         true,
				BodyType:       "",
				BodyFormat:     "",
				HasBody:        false,
				ResponseType:   "#blob",
				ResponseFormat: "",
				ResponseCode:   200,
				Permission:     "Reports_DownloadSalesActivityReportExcel",
				ControllerName: "Reports",
				FileName:       "app/api/joc/ReportsController.go",
				LineNo:         0,
			},
			"export const downloadSalesActivityReportExcel = async (year : number) => await axios.get(`/reports/sales-activity/download/excel?year=${year}`, { responseType: 'blob' });",
		},
		{
			&lib.ControllerRoute{
				Name:        "EnableUser",
				Description: "enables a user",
				Path:        "/users/{userID:[0-9]+}/enable",
				Method:      "PUT",
				Params: []lib.ControllerRouteParam{
					{
						Name:    "userID",
						Pattern: "[0-9]+",
						Type:    "int64",
					},
				},
				Queries:        []lib.ControllerRouteQuery{},
				IsAuth:         true,
				BodyType:       "",
				BodyFormat:     "",
				HasBody:        false,
				ResponseType:   "",
				ResponseFormat: "",
				ResponseCode:   200,
				Permission:     "Users_EnableUser",
				ControllerName: "Users",
				FileName:       "core/api/base/UsersController.go",
				LineNo:         0,
			},
			"export const enableUser = async (userID : number) => await axios.put<any>(`/users/${userID}/enable`, {});",
		},
		{
			&lib.ControllerRoute{
				Name:        "DeleteCustomer",
				Description: "deletes a company",
				Path:        "/customers/{customerID:[0-9]+}",
				Method:      "DELETE",
				Params: []lib.ControllerRouteParam{
					{
						Name:    "customerID",
						Pattern: "[0-9]+",
						Type:    "int64",
					},
				},
				Queries:        []lib.ControllerRouteQuery{},
				IsAuth:         true,
				BodyType:       "",
				BodyFormat:     "",
				HasBody:        false,
				ResponseType:   "",
				ResponseFormat: "",
				ResponseCode:   200,
				Permission:     "Customers_DeleteCustomer",
				ControllerName: "Customers",
				FileName:       "app/api/joc/CustomersController.go",
				LineNo:         0,
			},
			"export const deleteCustomer = async (customerID : number) => await axios.delete<any>(`/customers/${customerID}`);",
		},
	}

	for k := range tests {
		str := ""
		// str := genTSRoute(tests[k].route)
		assert.Equal(t, tests[k].expected, str)
	}

}
