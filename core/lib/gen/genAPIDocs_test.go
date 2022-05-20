package gen

import (
	"testing"

	"github.com/macinnir/dvc/core/lib"
	"github.com/stretchr/testify/assert"
)

func TestBaseType(t *testing.T) {

	tests := []struct {
		in       string
		out      string
		baseType string
	}{
		{"aggregates.AccountAggregate", "AccountAggregate", "aggregate"},
		{"*aggregates.AccountAggregate", "AccountAggregate", "aggregate"},
		{"[]*aggregates.AccountAggregate", "AccountAggregate", "aggregate"},

		{"*models.Account", "Account", "model"},
		{"models.Account", "Account", "model"},
		{"[]*models.Account", "Account", "model"},

		{"dtos.CreateAccountDTO", "CreateAccountDTO", "dto"},
		{"*dtos.CreateAccountDTO", "CreateAccountDTO", "dto"},
		{"[]*dtos.CreateAccountDTO", "CreateAccountDTO", "dto"},
	}

	for k := range tests {
		result := cleanObject(tests[k].in)
		assert.Equal(t, tests[k].out, result, "`%s` should be type `%s` but is actually `%s`", tests[k].in, tests[k].out, result)
		result2 := getBaseType(tests[k].in)
		assert.Equal(t, tests[k].baseType, result2, "`%s` should have a base type of `%s` but is actually `%s`", tests[k].in, tests[k].baseType, result2)
	}

}

func TestObjectToString(t *testing.T) {

}

func TestRoutePathHTML(t *testing.T) {

	route := &lib.ControllerRoute{
		Name:        "GetAccountUsers",
		Description: "gets a collection of users across all accounts",
		Raw:         "/accounts/{accountID:[0-9]+}/users?page={page:[0-9]+}\u0026limit={limit:[0-9]+}",
		Method:      "GET",
		Params: []lib.ControllerRouteParam{
			{
				Name:    "accountID",
				Pattern: "[0-9]+",
				Type:    "int64",
			},
		},
		Queries: []lib.ControllerRouteQuery{
			{
				Name:         "page",
				Pattern:      "[0-9]+",
				Type:         "int64",
				VariableName: "page",
				ValueRaw:     "{page:[0-9]+}",
			},
			{
				Name:         "limit",
				Pattern:      "[0-9]+",
				Type:         "int64",
				VariableName: "limit",
				ValueRaw:     "{limit:[0-9]+}",
			},
		},
		IsAuth:         true,
		BodyType:       "",
		BodyFormat:     "",
		HasBody:        false,
		ResponseType:   "aggregates.UserCollectionAggregate",
		ResponseFormat: "JSON",
		ResponseCode:   200,
		Permission:     "Accounts_GetAccountUsers",
		ControllerName: "Accounts",
		FileName:       "core/api/base/AccountsController.go",
		LineNo:         0,
	}

	out := buildRoutePathHTML(route)
	assert.Equal(t, `/accounts/<span class="endpoint-query-var">{accountID}</span>/users?<span class="endpoint-query-var">page</span>=[0-9]+&<span class="endpoint-query-var">limit</span>=[0-9]+`, out)

}
