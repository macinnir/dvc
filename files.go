package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type IFiles interface {
}

// Files
type Files struct {
}

// ScanChangesetDir scans the changeset directory for any sql files
func (f *Files) ScanChangesetDir(rootPath string) (paths []string, e error) {

	var fileInfos []os.FileInfo
	paths = []string{}

	fileInfos, e = ioutil.ReadDir(rootPath)

	for _, f := range fileInfos {
		if f.IsDir() || filepath.Ext(f.Name()) != ".sql" {
			continue
		}
		paths = append(paths, f.Name())
	}

	return

}

// FetchLocalChangesetList gets the contents of the changesets.json file
// and returns a collection of paths ([]string)
func (f *Files) FetchLocalChangesetList(rootPath string) (sqlPaths []string, e error) {

	changesetJSONPath := rootPath + "/changesets.json"

	fmt.Printf("Looking for changeset file at path %s\n", changesetJSONPath)
	if _, e = os.Stat(changesetJSONPath); os.IsNotExist(e) {
		return
	}

	var raw []byte

	raw, e = ioutil.ReadFile(changesetJSONPath)
	if e != nil {
		return
	}

	sqlPaths = []string{}

	json.Unmarshal(raw, &sqlPaths)

	for idx, sqlPath := range sqlPaths {
		sqlPaths[idx] = rootPath + "/" + sqlPath
	}

	return
}

// BuildChangeFiles returns a collection of ChangeFile objects based on
// a collection of sql paths
func (f *Files) BuildChangeFiles(sqlPaths []string) (changeFiles []ChangeFile, e error) {

	if len(sqlPaths) == 0 {
		return
	}

	changeFiles = []ChangeFile{}

	ordinal := 0

	for _, sqlPath := range sqlPaths {

		ordinal = ordinal + 1

		var changeFile *ChangeFile

		if changeFile, e = f.BuildChangeFile(sqlPath, ordinal); e != nil {
			return
		}

		changeFiles = append(changeFiles, *changeFile)
	}

	return
}

// BuildChangeFile builds a changeFile object based on a path and an ordinal
func (f *Files) BuildChangeFile(sqlPath string, ordinal int) (changeFile *ChangeFile, e error) {

	var content []byte

	if _, e = os.Stat(sqlPath); os.IsNotExist(e) {
		e = fmt.Errorf("Missing file in changeset: `%s`", sqlPath)
		return
	}

	if content, e = ioutil.ReadFile(sqlPath); e != nil {
		return
	}

	var hash string
	if hash, e = HashFileMd5(sqlPath); e != nil {
		return
	}

	changeFile = &ChangeFile{
		Name:    sqlPath,
		Content: string(content),
		Ordinal: ordinal,
		Hash:    hash,
	}

	return
}
