package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
func (a *AssignmentsController) CreateAssignment(w http.ResponseWriter, r *http.Request) {

	currentUser := utils.GetCurrentUser(r)
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

	currentUser := utils.GetCurrentUser(r)

	var assignment *models.Assignment
	assignment, e = a.assignmentService.UpdateAssignment(currentUser, assignmentID, dto)

	if e != nil {
		response.HandleError(r, w, e)
		return
	}

	response.JSON(r, w, assignment)
}

// GetAssignmentsByCurrentUser returns a slice of assignments by the current user
// @route GET /assignments?classID={fooID:[0-9a-zA-Z-]+}&barID={barID:[0-9]+}
func (a *AssignmentsController) GetAssignmentsByCurrentUser(w http.ResponseWriter, r *http.Request) {

	currentUser := utils.GetCurrentUser(r)

	var e error
	var assignments []models.Assignment

	assignments, e = a.assignmentService.GetMyAssignments(currentUser)

	if e != nil {
		response.HandleError(r, w, e)
		return
	}

	response.JSON(r, w, assignments)
}

// GetAssignmentByID returns a assignment by its unique ID
// @route GET /assignments/{assignmentID:[0-9]+}
func (a *AssignmentsController) GetAssignmentByID(w http.ResponseWriter, r *http.Request) {

	var e error
	var agg *aggregates.AssignmentAggregate

	assignmentID := request.URLParamInt64(r, "assignmentID", 0)
	if assignmentID == 0 {
		response.BadRequest(r, w, errors.NewError("Invalid memoID"))
		return
	}

	currentUser := utils.GetCurrentUser(r)
	agg, e = a.assignmentService.GetAssignmentByID(currentUser, assignmentID)

	if e != nil {
		response.HandleError(r, w, e)
		return
	}

	response.JSON(r, w, agg)
}

// DeleteAssignment deletes an assignment
// @route DELETE /assignments/{assignmentID:[0-9]+}
// @auth
func (a *AssignmentsController) DeleteAssignment(w http.ResponseWriter, r *http.Request) {

	var e error

	assignmentID := request.URLParamInt64(r, "assignmentID", 0)
	if assignmentID == 0 {
		response.NotFound(r, w)
		return
	}

	currentUser := utils.GetCurrentUser(r)
	e = a.assignmentService.DeleteAssignment(currentUser, assignmentID)

	if e != nil {
		response.HandleError(r, w, e)
		return
	}

	response.OK(r, w)
}
`

var routesCode = `// mapAssignmentsControllerRoutes maps all of the routes for AssignmentsController
func mapAssignmentsControllerRoutes(r *mux.Router, auth *utils.Auth, c *controllers.Controllers) {

	// CreateAssignment
	r.HandleFunc("/assignments", c.AssignmentsController.CreateAssignment).
		Methods("POST").
		Name("CreateAssignment")

	// UpdateAssignment
	r.HandleFunc("/assignments/{assignmentID:[0-9]+}", c.AssignmentsController.UpdateAssignment).
		Methods("PUT").
		Name("UpdateAssignment")

	// GetAssignmentsByCurrentUser
	r.HandleFunc("/assignments", c.AssignmentsController.GetAssignmentsByCurrentUser).
		Methods("GET").
		Queries(
			"classID", "{fooID:[0-9a-zA-Z-]+}",
			"barID", "{barID:[0-9]+}",
		).
		Name("GetAssignmentsByCurrentUser")

	// GetAssignmentByID
	r.HandleFunc("/assignments/{assignmentID:[0-9]+}", c.AssignmentsController.GetAssignmentByID).
		Methods("GET").
		Name("GetAssignmentByID")

	// DeleteAssignment
	r.Handle("/assignments/{assignmentID:[0-9]+}", auth.AuthMiddleware(c.AssignmentsController.DeleteAssignment)).
		Methods("DELETE").
		Name("DeleteAssignment")
}`

func TestExtractRoutesFromController(t *testing.T) {

	// Arrange
	g := &Gen{}

	// Act
	result, e := g.ExtractRoutesFromController(controllerFilePath, []byte(controller))

	// Assert
	require.Nil(t, e)
	require.Less(t, 0, len(result))

	assert.Equal(t, "POST", result[0].Method)
	assert.Equal(t, "/assignments", result[0].Path)
	assert.Equal(t, "CreateAssignment", result[0].Name)
	assert.Equal(t, "creates an assignment", result[0].Description)

	assert.Equal(t, "PUT", result[1].Method)
	assert.Equal(t, "/assignments/{assignmentID:[0-9]+}", result[1].Path)
	assert.Equal(t, "/assignments/{assignmentID:[0-9]+}", result[1].Raw)
	assert.Equal(t, "UpdateAssignment", result[1].Name)
	assert.Equal(t, "updates an assignment", result[1].Description)

	// Params
	require.Equal(t, 1, len(result[1].Params))
	assert.Equal(t, "assignmentID", result[1].Params[0].Name)
	assert.Equal(t, "[0-9]+", result[1].Params[0].Pattern)
	assert.Equal(t, "int64", result[1].Params[0].Type)

	assert.Equal(t, "GET", result[2].Method)
	assert.Equal(t, "/assignments", result[2].Path)
	assert.Equal(t, "GetAssignmentsByCurrentUser", result[2].Name)
	assert.Equal(t, "returns a slice of assignments by the current user", result[2].Description)

	require.Equal(t, 2, len(result[2].Queries))

	assert.Equal(t, "classID", result[2].Queries[0].Name)
	assert.Equal(t, "fooID", result[2].Queries[0].VariableName)
	assert.Equal(t, "[0-9a-zA-Z-]+", result[2].Queries[0].Pattern)
	assert.Equal(t, "{fooID:[0-9a-zA-Z-]+}", result[2].Queries[0].ValueRaw)
	assert.Equal(t, "string", result[2].Queries[0].Type)

	assert.Equal(t, "barID", result[2].Queries[1].Name)
	assert.Equal(t, "barID", result[2].Queries[1].VariableName)
	assert.Equal(t, "[0-9]+", result[2].Queries[1].Pattern)
	assert.Equal(t, "{barID:[0-9]+}", result[2].Queries[1].ValueRaw)
	assert.Equal(t, "int64", result[2].Queries[1].Type)

	assert.Equal(t, "GET", result[3].Method)
	assert.Equal(t, "/assignments/{assignmentID:[0-9]+}", result[3].Path)
	assert.Equal(t, "GetAssignmentByID", result[3].Name)
	assert.Equal(t, "returns a assignment by its unique ID", result[3].Description)

	assert.Equal(t, "DELETE", result[4].Method)
	assert.Equal(t, "/assignments/{assignmentID:[0-9]+}", result[4].Path)
	assert.Equal(t, "DeleteAssignment", result[4].Name)
	assert.Equal(t, "deletes an assignment", result[4].Description)
}

func TestGetControllerName(t *testing.T) {

	result := extractNameFromFile(controllerFilePath)
	assert.Equal(t, "AssignmentsController", result)
}

func TestBuildRoutesCodeFromController(t *testing.T) {
	g := &Gen{}

	ctl, e := g.BuildControllerObjFromController(controllerFilePath, []byte(controller))
	require.Nil(t, e)

	code := g.BuildRoutesCodeFromController(ctl)

	assert.Equal(t, routesCode, code)
}
