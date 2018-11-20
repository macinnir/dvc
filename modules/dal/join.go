package dal

// Join represents a join clause
type Join struct {
	Table  *Table
	Fields []JoinField
}
