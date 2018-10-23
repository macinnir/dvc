package lib

import (
	"testing"
)

func TestDvcBadConfigFilePath(t *testing.T) {

	_, e := NewDVC("badConfigPath")

	if e == nil {
		t.Error("Dvc should have returned an error when given a bad config file path...")
		return
	}
}

// func TestDvcInvalidDatabaseType(t *testing.T) {
// 	_, e := NewDVC("./test_resources/config2.toml")

// 	if e == nil {
// 		t.Error("Dvc should have thrown an error when the config2.toml has a bad database type.")
// 		return
// 	}

// 	err := e.Error()

// 	fmt.Printf("Error: %s\n", err)

// 	if err != "invalid database type" {
// 		t.Error("Dvc should have returned a config error indicating a bad value for the database type")
// 		return
// 	}
// }

func TestDvcGoodConfigFilePath(t *testing.T) {
	d, e := NewDVC("../test/dvc.toml")

	if e != nil {
		t.Errorf("dvc should not have returned an error: %s", e.Error())
		return
	}

	if d.Config.DatabaseName != "dbTest" {
		t.Errorf("wrong database name `%s` should be `%s`", d.Config.DatabaseName, "dbTest")
	}

	if d.Config.Host != "127.0.0.1:3307" {
		t.Errorf("wrong host `%s` should be `%s`", d.Config.Host, "127.0.0.1:3307")
	}

	if d.Config.Username != "root" {
		t.Errorf("Wrong username `%s` should be `%s`", d.Config.Username, "root")
	}

	if d.Config.Password != "root" {
		t.Errorf("Wrong password `%s` should be `%s`", d.Config.Password, "root")
	}

	if d.Config.ChangeSetPath != "test/resources/changes" {
		t.Errorf("Wrong changeSetPath `%s` should be `%s`", d.Config.ChangeSetPath, "changes")
	}

	if d.Config.DatabaseType != "mysql" {
		t.Errorf("Wrong database type `%s` should be `%s`", d.Config.DatabaseType, "mysql")
	}

}
func TestDvcMissingArguments(t *testing.T) {

	_, e := NewDVC()

	if e == nil {
		t.Error("No error was thrown with no arguments")
		return
	}

	if e.Error() != "not enough arguments" {
		t.Error("NewDVC should have returned a `not enough arguments` error")
	}
}
func TestDvcMultipleArguments(t *testing.T) {
	d, e := NewDVC("host", "name", "user", "pass", "test_resources/changes", "mysql")

	if e != nil {
		t.Error("should not have thrown an error")
	}

	if d.Config.Host != "host" || d.Config.DatabaseName != "name" || d.Config.Username != "user" || d.Config.Password != "pass" || d.Config.ChangeSetPath != "test_resources/changes" || d.Config.DatabaseType != "mysql" {
		t.Error("config value not correctly set")
	}
}

// func TestDvcImportSchema(t *testing.T) {
// 	var dvcTest *DVC
// 	var e error
// 	if dvcTest, e = NewDVC("config.toml"); e != nil {
// 		t.Errorf("should not have thrown an error: %s", e.Error())
// 		return
// 	}

// 	if e = dvcTest.ImportSchema("config."); e != nil {
// 		t.Errorf("should not have thrown an error: %s", e.Error())
// 	}
// }

// func TestDvcCompareSchema(t *testing.T) {
// 	var dvcTest *DVC
// 	var e error
// 	if dvcTest, e = NewDVC("config.toml"); e != nil {
// 		t.Errorf("should not have thrown an error: %s", e.Error())
// 		return
// 	}

// 	sql := ""

// 	if sql, e = dvcTest.CompareSchema(); e != nil {
// 		t.Errorf("should not have thrown an error: %s", e.Error())
// 	}

// 	fmt.Printf(sql)
// }
