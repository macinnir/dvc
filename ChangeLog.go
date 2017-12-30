package dvc

import "database/sql"

// ChangeLog represents a sql statement stored in a changelog file
type ChangeLog struct {
	ID            int64  `json:"id"`
	DBChangeLogID int64  `json:"dbChangeLogId"`
	DateCreated   int64  `json:"dateCreated"`
	FilePath      string `json:"filePath"`
}

func (c *ChangeLog) Apply(conn *sql.DB, databaseName string) {

}
