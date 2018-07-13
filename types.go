package main

import ()

type Options uint

type Command struct {
	Options Options
}

// Changeset represents all of the changes in an environment and their changes
type Changeset struct {
	ChangeFiles map[string]ChangeFile
	Signature   string
}

// ChangeFile represents both a physical file on the local file system
// along with the entry in the changefile database
type ChangeFile struct {
	ID          int64
	DateCreated int64
	Hash        string
	Name        string
	DVCSetID    int64
	IsRun       bool
	IsDeleted   bool
	Content     string
	FullPath    string
	Ordinal     int
}

// Config contains a set of configuration values used throughout the application
type Config struct {
	Host          string `toml:"host"`
	DatabaseName  string `toml:"databaseName"`
	Username      string `toml:"username"`
	Password      string `toml:"password"`
	ChangeSetPath string `toml:"changesetPath"`
	DatabaseType  string `toml:"databaseType"`
}

// SortedColumns is a slice of Column objects
type SortedColumns []*Column

// Len is part of sort.Interface.
func (c SortedColumns) Len() int {
	return len(c)
}

// Swap is part of sort.Interface.
func (c SortedColumns) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (c SortedColumns) Less(i, j int) bool {
	return c[i].Position < c[j].Position
}
