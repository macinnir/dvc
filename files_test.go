package main

import (
	"testing"
)

func TestFetchLocalChangesetListChangesetFileNotFound(t *testing.T) {

	f := Files{}

	_, e := f.FetchLocalChangesetList("missing_changesets.json")

	if e == nil {
		t.Error("Files should have returned an error when it couldn't find the changesets file")
		return
	}
}

func TestFetchLocalChangesetListSuccess(t *testing.T) {

	changesetPath := "test_resources/changes"

	f := &Files{}

	sqlPaths, e := f.FetchLocalChangesetList(changesetPath)

	if e != nil {
		t.Error("should not have returned an error (should have found the changesets file)")
		return
	}

	if len(sqlPaths) != 3 {
		t.Error("should have 3 changesets from changeset file")
		return
	}

	want := "test_resources/changes/0000/foo.sql"

	if sqlPaths[0] != want {
		t.Errorf("should have had `%s` in the first changeset item (found: %s)", want, sqlPaths[0])
	}
}
func TestBuildChangeFiles(t *testing.T) {

	changesetPath := "test_resources/changes"

	f := &Files{}

	sqlPaths, e := f.FetchLocalChangesetList(changesetPath)

	if e != nil {
		t.Error("fetchLocalChangesetList should have not thrown an error")
	}

	var changeFiles []ChangeFile

	changeFiles, e = f.BuildChangeFiles(sqlPaths)

	if e != nil {
		t.Errorf("should not have thrown an error: %s", e.Error())
	}

	if len(changeFiles) != 3 {
		t.Error("should have 3 changesets")
		return
	}

	if changeFiles[0].Name != "test_resources/changes/0000/foo.sql" {
		t.Error("should have correct name to changeset sql file")
	}

	if changeFiles[0].Content != "create table foo (id int);" {
		t.Errorf("should have correct sql content: %s", changeFiles[0].Content)
	}

	if changeFiles[0].Hash != "0b900d2fe13c1f66c90de3bc51d42c05" {
		t.Errorf("should have generated the correct hash.")
	}
}
