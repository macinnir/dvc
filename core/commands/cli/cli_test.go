package cli

import (
	"testing"

	"github.com/macinnir/dvc/core/lib"
	"github.com/stretchr/testify/assert"
)

var anonRoute = &lib.ControllerRoute{
	Path:   "/join/{userID:[0-9]+}/activate?activationCode={activationCode:[0-9a-zA-Z-]+}\u0026nonceID={nonceID:[0-9]+}",
	Method: "GET",
	Params: []lib.ControllerRouteParam{
		{
			Name:    "userID",
			Pattern: "[0-9]+",
			Type:    "int64",
		},
	},
	Queries: []lib.ControllerRouteQuery{
		{
			Name:         "activationCode",
			Pattern:      "[0-9a-zA-Z-]+",
			Type:         "string",
			VariableName: "activationCode",
			ValueRaw:     "{activationCode:[0-9a-zA-Z-]+}",
		},
		{
			Name:         "nonceID",
			Pattern:      "[0-9]+",
			Type:         "int64",
			VariableName: "nonceID",
			ValueRaw:     "{nonceID:[0-9]+}",
		},
	},
	IsAuth:         false,
	BodyType:       "",
	BodyFormat:     "",
	HasBody:        false,
	ResponseType:   "",
	ResponseFormat: "",
	ResponseCode:   200,
	Permission:     "",
}

var getCurrentUserRoute = &lib.ControllerRoute{
	Name:           "GetCurrentUser",
	Description:    "gets the current user",
	Path:           "/users/self",
	Method:         "GET",
	Params:         []lib.ControllerRouteParam{},
	Queries:        []lib.ControllerRouteQuery{},
	IsAuth:         true,
	BodyType:       "",
	BodyFormat:     "",
	HasBody:        false,
	ResponseType:   "*aggregates.UserAggregate",
	ResponseFormat: "JSON",
	ResponseCode:   200,
	Permission:     "Identity_GetCurrentUser",
}

func TestExtractBasePath(t *testing.T) {

	path := extractBasePath("http://localhost:8080/api/", anonRoute)

	// assert.Equal(t, "http://localhost:8080/api/join/{userID:[0-9]+}/activate", path)
	assert.Equal(t, "http://localhost:8080/api/join/{userID:[0-9]+}/activate?activationCode={activationCode:[0-9a-zA-Z-]+}&nonceID={nonceID:[0-9]+}", path)
}

func TestExtractBasePath_CurrentUser(t *testing.T) {
	path := extractBasePath("http://localhost:8080/api", getCurrentUserRoute)
	assert.Equal(t, "http://localhost:8080/api/users/self", path)
}

func TestApplyParam(t *testing.T) {
	url := "http://localhost:8080/api/join/{userID:[0-9]+}/activate"
	result := applyParam(url, lib.ControllerRouteParam{
		Name: "userID",
	}, "1234")

	assert.Equal(t, "http://localhost:8080/api/join/1234/activate", result)
}

func TestApplyQuery(t *testing.T) {
	url := "http://localhost:8080/api/join/{userID:[0-9]+}/activate?activationCode={activationCode:[0-9a-zA-Z-]+}\u0026nonceID={nonceID:[0-9]+}"
	result := applyQuery(url, lib.ControllerRouteQuery{
		Name: "nonceID",
	}, "11111")
	assert.Equal(t, "http://localhost:8080/api/join/{userID:[0-9]+}/activate?activationCode={activationCode:[0-9a-zA-Z-]+}\u0026nonceID=11111", result)
}
