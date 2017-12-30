package dvc

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sort"
)

// ChangeSet represents a logical (and physical) collection of change files
type ChangeSet struct {
	Name  string   `json:"name"`
	Files []string `json:"files"`
}

// LocalChangeLogs represents local sql change scripts, organized by changesets
type LocalChangeLogs struct {
	RootDir             string `json:"rootDir"`
	ChangeSets          map[string]ChangeSet
	SortedChangesetKeys []string
}

// FetchLocalChangesetFiles fetches a list of changeset directories and their child sql files
// and populates the ChangeSets property with a map of ChangeSet objects indexed by their name
func (l *LocalChangeLogs) FetchLocalChangesetFiles() (e error) {

	l.SortedChangesetKeys, e = FetchDirFileNames(l.RootDir)
	sort.Strings(l.SortedChangesetKeys)

	if e != nil {
		return
	}

	l.ChangeSets = map[string]ChangeSet{}

	for _, changeDir := range l.SortedChangesetKeys {

		c := ChangeSet{Name: changeDir, Files: []string{}}

		fullPath := l.RootDir + string(filepath.Separator) + changeDir
		changeFiles, e := FetchNonDirFileNames(fullPath)

		if e != nil {
			log.Fatal(e)
		}

		sort.Strings(changeFiles)
		c.Files = changeFiles
		l.ChangeSets[changeDir] = c

	}
	return
}

func NewLocalChangeLogs(rootDir string) (localChangeLogs LocalChangeLogs, e error) {

	if _, e = os.Stat(rootDir); os.IsNotExist(e) {
		e = errors.New("changesets dir not found")
		return
	}

	localChangeLogs = LocalChangeLogs{RootDir: rootDir}

	return
}
