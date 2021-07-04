package testassets

import (
	"encoding/json"

	"github.com/macinnir/dvc/core/lib/utils/query"
)

const (

	// TaskBatchSchedule_SchemaName is the name of the schema group this model is in
	TaskBatchSchedule_SchemaName = "dvc"

	// TaskBatchSchedule_TableName is the name of the table
	TaskBatchSchedule_TableName query.TableName = "TaskBatchSchedule"

	// Columns
	TaskBatchSchedule_Column_TaskBatchScheduleID query.Column = "TaskBatchScheduleID"
	TaskBatchSchedule_Column_TaskBatchID         query.Column = "TaskBatchID"
	TaskBatchSchedule_Column_DOM                 query.Column = "DOM"
	TaskBatchSchedule_Column_DOW                 query.Column = "DOW"
	TaskBatchSchedule_Column_HOD                 query.Column = "HOD"
	TaskBatchSchedule_Column_MOH                 query.Column = "MOH"
	TaskBatchSchedule_Column_DateCreated         query.Column = "DateCreated"
	TaskBatchSchedule_Column_IsDeleted           query.Column = "IsDeleted"
	TaskBatchSchedule_Column_IsActive            query.Column = "IsActive"
	TaskBatchSchedule_Column_CurrentStatus       query.Column = "CurrentStatus"
	TaskBatchSchedule_Column_LastRunDate         query.Column = "LastRunDate"
	TaskBatchSchedule_Column_LastSuccessDate     query.Column = "LastSuccessDate"
	TaskBatchSchedule_Column_LastErrorDate       query.Column = "LastErrorDate"
	TaskBatchSchedule_Column_NumFailed           query.Column = "NumFailed"
	TaskBatchSchedule_Column_NumFinished         query.Column = "NumFinished"
	TaskBatchSchedule_Column_NumAttempted        query.Column = "NumAttempted"
)

var (
	// TaskBatchSchedule_Columns is a list of all the columns
	TaskBatchSchedule_Columns = []query.Column{
		TaskBatchSchedule_Column_TaskBatchScheduleID, TaskBatchSchedule_Column_TaskBatchID, TaskBatchSchedule_Column_DOM, TaskBatchSchedule_Column_DOW, TaskBatchSchedule_Column_HOD, TaskBatchSchedule_Column_MOH, TaskBatchSchedule_Column_DateCreated, TaskBatchSchedule_Column_IsDeleted, TaskBatchSchedule_Column_IsActive, TaskBatchSchedule_Column_CurrentStatus, TaskBatchSchedule_Column_LastRunDate, TaskBatchSchedule_Column_LastSuccessDate, TaskBatchSchedule_Column_LastErrorDate, TaskBatchSchedule_Column_NumFailed, TaskBatchSchedule_Column_NumFinished, TaskBatchSchedule_Column_NumAttempted}

	// TaskBatchSchedule_Column_Types maps columns to their string types
	TaskBatchSchedule_Column_Types = map[query.Column]string{
		TaskBatchSchedule_Column_TaskBatchScheduleID: "%d", TaskBatchSchedule_Column_TaskBatchID: "%d", TaskBatchSchedule_Column_DOM: "%d", TaskBatchSchedule_Column_DOW: "%d", TaskBatchSchedule_Column_HOD: "%d", TaskBatchSchedule_Column_MOH: "%d", TaskBatchSchedule_Column_DateCreated: "%d", TaskBatchSchedule_Column_IsDeleted: "%d", TaskBatchSchedule_Column_IsActive: "%d", TaskBatchSchedule_Column_CurrentStatus: "%d", TaskBatchSchedule_Column_LastRunDate: "%d", TaskBatchSchedule_Column_LastSuccessDate: "%d", TaskBatchSchedule_Column_LastErrorDate: "%d", TaskBatchSchedule_Column_NumFailed: "%d", TaskBatchSchedule_Column_NumFinished: "%d", TaskBatchSchedule_Column_NumAttempted: "%d"}
	// TaskBatchSchedule_UpdateColumns is a list of all update columns for this model
	TaskBatchSchedule_UpdateColumns = []query.Column{TaskBatchSchedule_Column_TaskBatchID, TaskBatchSchedule_Column_DOM, TaskBatchSchedule_Column_DOW, TaskBatchSchedule_Column_HOD, TaskBatchSchedule_Column_MOH, TaskBatchSchedule_Column_IsDeleted, TaskBatchSchedule_Column_IsActive, TaskBatchSchedule_Column_CurrentStatus, TaskBatchSchedule_Column_LastRunDate, TaskBatchSchedule_Column_LastSuccessDate, TaskBatchSchedule_Column_LastErrorDate, TaskBatchSchedule_Column_NumFailed, TaskBatchSchedule_Column_NumFinished, TaskBatchSchedule_Column_NumAttempted}
	// TaskBatchSchedule_InsertColumns is a list of all insert columns for this model
	TaskBatchSchedule_InsertColumns = []query.Column{TaskBatchSchedule_Column_TaskBatchID, TaskBatchSchedule_Column_DOM, TaskBatchSchedule_Column_DOW, TaskBatchSchedule_Column_HOD, TaskBatchSchedule_Column_MOH, TaskBatchSchedule_Column_DateCreated, TaskBatchSchedule_Column_IsActive, TaskBatchSchedule_Column_CurrentStatus, TaskBatchSchedule_Column_LastRunDate, TaskBatchSchedule_Column_LastSuccessDate, TaskBatchSchedule_Column_LastErrorDate, TaskBatchSchedule_Column_NumFailed, TaskBatchSchedule_Column_NumFinished, TaskBatchSchedule_Column_NumAttempted}
	// TaskBatchSchedule_PrimaryKey is the name of the table's primary key
	TaskBatchSchedule_PrimaryKey query.Column = "TaskBatchScheduleID"
)

// TaskBatchSchedule is a `TaskBatchSchedule` data model
type TaskBatchSchedule struct {
	TaskBatchScheduleID int64 `db:"TaskBatchScheduleID" json:"TaskBatchScheduleID"`
	TaskBatchID         int64 `db:"TaskBatchID" json:"TaskBatchID"`
	DOM                 int64 `db:"DOM" json:"DOM"`
	DOW                 int64 `db:"DOW" json:"DOW"`
	HOD                 int64 `db:"HOD" json:"HOD"`
	MOH                 int64 `db:"MOH" json:"MOH"`
	DateCreated         int64 `db:"DateCreated" json:"DateCreated"`
	IsDeleted           int   `db:"IsDeleted" json:"IsDeleted"`
	IsActive            int   `db:"IsActive" json:"IsActive"`
	CurrentStatus       int   `db:"CurrentStatus" json:"CurrentStatus"`
	LastRunDate         int64 `db:"LastRunDate" json:"LastRunDate"`
	LastSuccessDate     int64 `db:"LastSuccessDate" json:"LastSuccessDate"`
	LastErrorDate       int64 `db:"LastErrorDate" json:"LastErrorDate"`
	NumFailed           int64 `db:"NumFailed" json:"NumFailed"`
	NumFinished         int64 `db:"NumFinished" json:"NumFinished"`
	NumAttempted        int64 `db:"NumAttempted" json:"NumAttempted"`
}

// TaskBatchSchedule_TableName is the name of the table
func (c *TaskBatchSchedule) Table_Name() query.TableName {
	return TaskBatchSchedule_TableName
}

func (c *TaskBatchSchedule) Table_Columns() []query.Column {
	return TaskBatchSchedule_Columns
}

// Table_ColumnTypes returns a map of tableColumn names with their fmt string types
func (c *TaskBatchSchedule) Table_Column_Types() map[query.Column]string {
	return TaskBatchSchedule_Column_Types
}

// Table_PrimaryKey returns the name of this table's primary key
func (c *TaskBatchSchedule) Table_PrimaryKey() query.Column {
	return TaskBatchSchedule_PrimaryKey
}

// Table_PrimaryKey_Value returns the value of this table's primary key
func (c *TaskBatchSchedule) Table_PrimaryKey_Value() int64 {
	return c.TaskBatchScheduleID
}

// Table_InsertColumns is a list of all insert columns for this model
func (c *TaskBatchSchedule) Table_InsertColumns() []query.Column {
	return TaskBatchSchedule_InsertColumns
}

// Table_UpdateColumns is a list of all update columns for this model
func (c *TaskBatchSchedule) Table_UpdateColumns() []query.Column {
	return TaskBatchSchedule_UpdateColumns
}

// TaskBatchSchedule_SchemaName is the name of this table's schema
func (c *TaskBatchSchedule) Table_SchemaName() string {
	return TaskBatchSchedule_SchemaName
}

// Select starts a select statement
func (c *TaskBatchSchedule) FromID(id int64) string {
	q, _ := query.Select(c).Where(query.EQ(TaskBatchSchedule_Column_TaskBatchScheduleID, id)).String()
	return q
}

// String returns a json marshalled string of the object
func (c *TaskBatchSchedule) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

// Update updates a TaskBatchSchedule record
func (c *TaskBatchSchedule) Update() string {

	var sql string
	sql, _ = query.Update(c).
		Set(TaskBatchSchedule_Column_TaskBatchID, c.TaskBatchID).
		Set(TaskBatchSchedule_Column_DOM, c.DOM).
		Set(TaskBatchSchedule_Column_DOW, c.DOW).
		Set(TaskBatchSchedule_Column_HOD, c.HOD).
		Set(TaskBatchSchedule_Column_MOH, c.MOH).
		Set(TaskBatchSchedule_Column_IsDeleted, c.IsDeleted).
		Set(TaskBatchSchedule_Column_IsActive, c.IsActive).
		Set(TaskBatchSchedule_Column_CurrentStatus, c.CurrentStatus).
		Set(TaskBatchSchedule_Column_LastRunDate, c.LastRunDate).
		Set(TaskBatchSchedule_Column_LastSuccessDate, c.LastSuccessDate).
		Set(TaskBatchSchedule_Column_LastErrorDate, c.LastErrorDate).
		Set(TaskBatchSchedule_Column_NumFailed, c.NumFailed).
		Set(TaskBatchSchedule_Column_NumFinished, c.NumFinished).
		Set(TaskBatchSchedule_Column_NumAttempted, c.NumAttempted).
		Where(query.EQ(TaskBatchSchedule_Column_TaskBatchScheduleID, c.TaskBatchScheduleID)).
		String()

	return sql
}

// Create inserts a TaskBatchSchedule record
func (c *TaskBatchSchedule) Create() string {

	var sql string
	q := query.Insert(c)

	if c.TaskBatchScheduleID > 0 {
		q.Set(TaskBatchSchedule_Column_TaskBatchScheduleID, c.TaskBatchScheduleID)
	}
	q.Set(TaskBatchSchedule_Column_TaskBatchID, c.TaskBatchID)
	q.Set(TaskBatchSchedule_Column_DOM, c.DOM)
	q.Set(TaskBatchSchedule_Column_DOW, c.DOW)
	q.Set(TaskBatchSchedule_Column_HOD, c.HOD)
	q.Set(TaskBatchSchedule_Column_MOH, c.MOH)
	q.Set(TaskBatchSchedule_Column_DateCreated, c.DateCreated)
	q.Set(TaskBatchSchedule_Column_IsActive, c.IsActive)
	q.Set(TaskBatchSchedule_Column_CurrentStatus, c.CurrentStatus)
	q.Set(TaskBatchSchedule_Column_LastRunDate, c.LastRunDate)
	q.Set(TaskBatchSchedule_Column_LastSuccessDate, c.LastSuccessDate)
	q.Set(TaskBatchSchedule_Column_LastErrorDate, c.LastErrorDate)
	q.Set(TaskBatchSchedule_Column_NumFailed, c.NumFailed)
	q.Set(TaskBatchSchedule_Column_NumFinished, c.NumFinished)
	q.Set(TaskBatchSchedule_Column_NumAttempted, c.NumAttempted)

	sql, _ = q.String()
	return sql
}

// Destroy deletes a TaskBatchSchedule record
func (c *TaskBatchSchedule) Destroy() string {
	sql, _ := query.Delete(c).
		Where(
			query.EQ(TaskBatchSchedule_Column_TaskBatchScheduleID, c.TaskBatchScheduleID),
		).String()
	return sql
}
