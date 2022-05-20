package gen

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenControllerBootstrap(t *testing.T) {

	// var result = GenControllerBootstrap("foo", []string{ "bar", "baz", "quux"})
	// assert.Equal(`

	// `)

}

func TestGetConstantsFromGoFile(t *testing.T) {

	var buffer bytes.Buffer
	buffer.WriteString(`package constants

// ProjectType is a type of project
type ProjectType int

const (
	// ProjectTypeUnknown is an unknown project type
	ProjectTypeUnknown ProjectType = 1
	// ProjectTypeNew is a project on a new application/service
	ProjectTypeNew ProjectType = 2
	// ProjectTypeSupport is a project on an existing application/service
	ProjectTypeSupport ProjectType = 3
)

func (p ProjectType) String() string {
	switch p {
	case ProjectTypeUnknown:
		return "Unknown"
	case ProjectTypeNew:
		return "New"
	case ProjectTypeSupport:
		return "Support"
	default:
		return "Unknown"
	}
}

// Int64 returns the int64 representation of the constant
func (p ProjectType) Int64() int64 {
	return int64(p)
}`)

	var name, constants = getConstantsFromGoFile(&buffer)

	assert.Equal(t, "ProjectType", name)
	require.Equal(t, 3, len(constants))
	assert.Equal(t, "ProjectTypeUnknown", constants[0])
	assert.Equal(t, "ProjectTypeNew", constants[1])
	assert.Equal(t, "ProjectTypeSupport", constants[2])
}

var controllerFilePath = "/Users/robertmacinnis/go/src/axis-api/core/controllers/AssignmentsController.go"
var controller = `package controllers

import (
	"axis-api/core/definitions/aggregates"
	"axis-api/core/definitions/dtos"
	"axis-api/core/definitions/models"
	iservices "axis-api/core/definitions/services"
	"axis-api/core/utils"
	"axis-api/core/utils/errors"
	"axis-api/core/utils/request"
	"axis-api/core/utils/response"
	"net/http"
)

// AssignmentsController is a container for user methods
type AssignmentsController struct {
	assignmentService iservices.IAssignmentService
}

// NewAssignmentsController returns a new AssignmentsController
func NewAssignmentsController(assignmentService iservices.IAssignmentService) *AssignmentsController {
	return &AssignmentsController{assignmentService}
}

// CreateAssignment creates an assignment
// @route POST /assignments
// @body JSON dtos.CreateAssignmentRequestDTO
// @response 200 JSON *aggregates.AssignmentAggregate
func (a *AssignmentsController) CreateAssignment(w http.ResponseWriter, r *http.Request, body *dtos.CreateAssignmentRequestDTO) {

	currentUser := auth.GetCurrentUser(r)
	var e error

	dto := &dtos.CreateAssignmentRequestDTO{}
	e = request.GetBodyJSON(r, dto)

	if e != nil {
		response.BadRequest(r, w, e)
		return
	}

	var model *models.Assignment
	model, e = a.assignmentService.CreateAssignment(currentUser, dto)

	agg := aggregates.NewAssignmentAggregate(model)

	if e != nil {
		response.HandleError(r, w, e)
		return
	}

	response.JSON(r, w, agg)
}

// UpdateAssignment updates an assignment
// @route PUT /assignments/{assignmentID:[0-9]+}
// @response 200 JSON *models.Assignment
func (a *AssignmentsController) UpdateAssignment(w http.ResponseWriter, r *http.Request) {

	var e error

	dto := &dtos.UpdateAssignmentRequestDTO{}
	e = request.GetBodyJSON(r, dto)

	if e != nil {
		response.BadRequest(r, w, e)
		return
	}

	assignmentID := request.URLParamInt64(r, "assignmentID", 0)

	if assignmentID == 0 {
		response.BadRequest(r, w, errors.NewError("Invalid assignment id"))
		return
	}

	currentUser := auth.GetCurrentUser(r)

	var assignment *models.Assignment
	assignment, e = a.assignmentService.UpdateAssignment(currentUser, assignmentID, dto)

	if e != nil {
		response.HandleError(r, w, e)
		return
	}

	response.JSON(r, w, assignment)
}

// GetAssignmentsByCurrentUser returns a slice of assignments by the current user
// @route GET /assignments?classID={fooID:[0-9a-zA-Z-]+}&barID={barID:[0-9]+}&any={any}
// @response 200 JSON []models.Assignment
// @auth
func (a *AssignmentsController) GetAssignmentsByCurrentUser(w http.ResponseWriter, r *http.Request) {

	currentUser := auth.GetCurrentUser(r)

	var e error
	var assignments []models.Assignment

	assignments, e = a.assignmentService.GetMyAssignments(currentUser)

	if e != nil {
		response.HandleError(r, w, e)
		return
	}

	response.JSON(r, w, assignments)
}

// GetFooByAssignment returns a slice of foo by an assignment
// @route GET /assignments/{assignmentID:[0-9]+}/foo
// @response 200 JSON []models.Foo
// @auth
func (a *AssignmentsController) GetFooByAssignment(w http.ResponseWriter, r *http.Request, assignmentID int64) {

	currentUser := auth.GetCurrentUser(r)

	var e error
	var foo []models.Foo

	foo, e = a.assignmentService.GetFooByAssignment(currentUser, assignmentID)

	if e != nil {
		response.HandleError(r, w, e)
		return
	}

	response.JSON(r, w, foo)
}

// GetAssignmentByID returns a assignment by its unique ID
// @route GET /assignments/{assignmentID:[0-9]+}
// @response 200 JSON *aggregates.AssignmentAggregate
func (a *AssignmentsController) GetAssignmentByID(w http.ResponseWriter, r *http.Request) {

	var e error
	var agg *aggregates.AssignmentAggregate

	assignmentID := request.URLParamInt64(r, "assignmentID", 0)
	if assignmentID == 0 {
		response.BadRequest(r, w, errors.NewError("Invalid memoID"))
		return
	}

	currentUser := auth.GetCurrentUser(r)
	agg, e = a.assignmentService.GetAssignmentByID(currentUser, assignmentID)

	if e != nil {
		response.HandleError(r, w, e)
		return
	}

	response.JSON(r, w, agg)
}

// DeleteAssignment deletes an assignment
// @route DELETE /assignments/{assignmentID:[0-9]+}
// @response 200
// @auth
func (a *AssignmentsController) DeleteAssignment(w http.ResponseWriter, r *http.Request) {

	var e error

	assignmentID := request.URLParamInt64(r, "assignmentID", 0)
	if assignmentID == 0 {
		response.NotFound(r, w)
		return
	}

	currentUser := auth.GetCurrentUser(r)
	e = a.assignmentService.DeleteAssignment(currentUser, assignmentID)

	if e != nil {
		response.HandleError(r, w, e)
		return
	}

	response.OK(r, w)
}
`

var routesCode = `// mapAssignmentsControllerRoutes maps all of the routes for AssignmentsController
func mapAssignmentsControllerRoutes(r *mux.Router, auth integrations.IAuth, c *controllers.Controllers) {

	// CreateAssignment
	// POST /assignments
	r.HandleFunc("/assignments", func(w http.ResponseWriter, r *http.Request) {

		body := &dtos.CreateAssignmentRequestDTO{}
		e := request.GetBodyJSON(r, body)

		if e != nil {
			response.BadRequest(r, w, e)
			return
		}

		c.AssignmentsController.CreateAssignment(w, r, body)

	}).
		Methods("POST").
		Name("CreateAssignment")

	// UpdateAssignment
	// PUT /assignments/{assignmentID:[0-9]+}
	r.HandleFunc("/assignments/{assignmentID:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {

		// URL Param assignmentID
		assignmentID := request.URLParamInt64(r, "assignmentID", 0)

		c.AssignmentsController.UpdateAssignment(w, r, assignmentID)

	}).
		Methods("PUT").
		Name("UpdateAssignment")

	// GetAssignmentsByCurrentUser
	// GET /assignments?classID={fooID:[0-9a-zA-Z-]+}&barID={barID:[0-9]+}&any={any}
	r.Handle("/assignments", auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {

		currentUser := auth.GetCurrentUser(r)

		// Query Arg fooID
		fooID := request.QueryArgString(r, "fooID", "")

		// Query Arg barID
		barID := request.QueryArgInt64(r, "barID", 0)

		// Query Arg any
		any := request.QueryArgString(r, "any", "")

		c.AssignmentsController.GetAssignmentsByCurrentUser(w, r, currentUser, fooID, barID, any)

	})).
		Methods("GET").
		Queries(
			"classID", "{fooID:[0-9a-zA-Z-]+}",
			"barID", "{barID:[0-9]+}",
			"any", "{any}",
		).
		Name("GetAssignmentsByCurrentUser")

	// GetFooByAssignment
	// GET /assignments/{assignmentID:[0-9]+}/foo
	r.Handle("/assignments/{assignmentID:[0-9]+}/foo", auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {

		currentUser := auth.GetCurrentUser(r)

		// URL Param assignmentID
		assignmentID := request.URLParamInt64(r, "assignmentID", 0)

		c.AssignmentsController.GetFooByAssignment(w, r, currentUser, assignmentID)

	})).
		Methods("GET").
		Name("GetFooByAssignment")

	// GetAssignmentByID
	// GET /assignments/{assignmentID:[0-9]+}
	r.HandleFunc("/assignments/{assignmentID:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {

		// URL Param assignmentID
		assignmentID := request.URLParamInt64(r, "assignmentID", 0)

		c.AssignmentsController.GetAssignmentByID(w, r, assignmentID)

	}).
		Methods("GET").
		Name("GetAssignmentByID")

	// DeleteAssignment
	// DELETE /assignments/{assignmentID:[0-9]+}
	r.Handle("/assignments/{assignmentID:[0-9]+}", auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {

		currentUser := auth.GetCurrentUser(r)

		// URL Param assignmentID
		assignmentID := request.URLParamInt64(r, "assignmentID", 0)

		c.AssignmentsController.DeleteAssignment(w, r, currentUser, assignmentID)

	})).
		Methods("DELETE").
		Name("DeleteAssignment")

}`

// func TestExtractRoutesFromController(t *testing.T) {

// 	// Arrange
// 	g := &Gen{}

// 	// Act
// 	result, e := g.ExtractRoutesFromController(controllerFilePath, []byte(controller))

// 	// Assert
// 	require.Nil(t, e)
// 	require.Less(t, 0, len(result))

// 	// CreateAssignment
// 	assert.Equal(t, "POST", result[0].Method)
// 	assert.Equal(t, "/assignments", result[0].Path)
// 	assert.Equal(t, "CreateAssignment", result[0].Name)
// 	assert.Equal(t, "creates an assignment", result[0].Description)
// 	assert.Equal(t, "JSON", result[0].BodyFormat)
// 	assert.Equal(t, "dtos.CreateAssignmentRequestDTO", result[0].BodyType)
// 	assert.Equal(t, int(200), result[0].ResponseCode)
// 	assert.Equal(t, "JSON", result[0].ResponseFormat)
// 	assert.Equal(t, "*aggregates.AssignmentAggregate", result[0].ResponseType)

// 	// UpdateAssignment
// 	assert.Equal(t, "PUT", result[1].Method)
// 	assert.Equal(t, "/assignments/{assignmentID:[0-9]+}", result[1].Path)
// 	assert.Equal(t, "/assignments/{assignmentID:[0-9]+}", result[1].Raw)
// 	assert.Equal(t, "UpdateAssignment", result[1].Name)
// 	assert.Equal(t, "updates an assignment", result[1].Description)
// 	assert.Equal(t, int(200), result[1].ResponseCode)
// 	assert.Equal(t, "JSON", result[1].ResponseFormat)
// 	assert.Equal(t, "*models.Assignment", result[1].ResponseType)

// 	require.Equal(t, 1, len(result[1].Params))
// 	assert.Equal(t, "assignmentID", result[1].Params[0].Name)
// 	assert.Equal(t, "[0-9]+", result[1].Params[0].Pattern)
// 	assert.Equal(t, "int64", result[1].Params[0].Type)

// 	// GetAssignmentsByCurrentUser
// 	assert.Equal(t, "GET", result[2].Method)
// 	assert.Equal(t, "/assignments", result[2].Path)
// 	assert.Equal(t, "GetAssignmentsByCurrentUser", result[2].Name)
// 	assert.Equal(t, "returns a slice of assignments by the current user", result[2].Description)
// 	assert.Equal(t, int(200), result[2].ResponseCode)
// 	assert.Equal(t, "JSON", result[2].ResponseFormat)
// 	assert.Equal(t, "[]models.Assignment", result[2].ResponseType)

// 	require.Equal(t, 3, len(result[2].Queries))

// 	assert.Equal(t, "classID", result[2].Queries[0].Name)
// 	assert.Equal(t, "fooID", result[2].Queries[0].VariableName)
// 	assert.Equal(t, "[0-9a-zA-Z-]+", result[2].Queries[0].Pattern)
// 	assert.Equal(t, "{fooID:[0-9a-zA-Z-]+}", result[2].Queries[0].ValueRaw)
// 	assert.Equal(t, "string", result[2].Queries[0].Type)

// 	assert.Equal(t, "barID", result[2].Queries[1].Name)
// 	assert.Equal(t, "barID", result[2].Queries[1].VariableName)
// 	assert.Equal(t, "[0-9]+", result[2].Queries[1].Pattern)
// 	assert.Equal(t, "{barID:[0-9]+}", result[2].Queries[1].ValueRaw)
// 	assert.Equal(t, "int64", result[2].Queries[1].Type)

// 	assert.Equal(t, "any", result[2].Queries[2].Name)
// 	assert.Equal(t, "any", result[2].Queries[2].VariableName)
// 	assert.Equal(t, "", result[2].Queries[2].Pattern)
// 	assert.Equal(t, "{any}", result[2].Queries[2].ValueRaw)
// 	assert.Equal(t, "string", result[2].Queries[2].Type)

// 	// GetAssignmentFoo
// 	assert.Equal(t, "GET", result[3].Method)
// 	assert.Equal(t, "/assignments/{assignmentID:[0-9]+}/foo", result[3].Path)
// 	assert.Equal(t, "GetFooByAssignment", result[3].Name)
// 	assert.Equal(t, "returns a slice of foo by an assignment", result[3].Description)
// 	require.Equal(t, 1, len(result[3].Params))
// 	assert.Equal(t, "assignmentID", result[3].Params[0].Name)
// 	assert.Equal(t, "int64", result[3].Params[0].Type)

// 	// GetAssignmentByID
// 	assert.Equal(t, "GET", result[4].Method)
// 	assert.Equal(t, "/assignments/{assignmentID:[0-9]+}", result[4].Path)
// 	assert.Equal(t, "GetAssignmentByID", result[4].Name)
// 	assert.Equal(t, "returns a assignment by its unique ID", result[4].Description)
// 	require.Equal(t, 1, len(result[4].Params))
// 	assert.Equal(t, "assignmentID", result[4].Params[0].Name)
// 	assert.Equal(t, "int64", result[4].Params[0].Type)
// 	assert.Equal(t, int(200), result[4].ResponseCode)
// 	assert.Equal(t, "JSON", result[4].ResponseFormat)
// 	assert.Equal(t, "*aggregates.AssignmentAggregate", result[4].ResponseType)

// 	// DeleteAssignment
// 	assert.Equal(t, "DELETE", result[5].Method)
// 	assert.Equal(t, "/assignments/{assignmentID:[0-9]+}", result[5].Path)
// 	assert.Equal(t, "DeleteAssignment", result[5].Name)
// 	assert.Equal(t, "deletes an assignment", result[5].Description)
// }
