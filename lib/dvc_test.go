package lib

import (
	"testing"
)

func TestDvcConfigValues(t *testing.T) {

	config := &Config{}
	config.ChangeSetPath = "test/resources/changes"
	config.DatabaseType = "mysql"
	config.Connection.DatabaseName = "dbTest"
	config.Connection.Host = "127.0.0.1:3307"
	config.Connection.Username = "root"
	config.Connection.Password = "root"
	config.BasePackage = "testbp"

	d, e := NewDVC(config)

	if e != nil {
		t.Errorf("dvc should not have returned an error: %s", e.Error())
		return
	}

	if d.Config.Connection.DatabaseName != "dbTest" {
		t.Errorf("wrong database name `%s` should be `%s`", d.Config.Connection.DatabaseName, "dbTest")
	}

	if d.Config.Connection.Host != "127.0.0.1:3307" {
		t.Errorf("wrong host `%s` should be `%s`", d.Config.Connection.Host, "127.0.0.1:3307")
	}

	if d.Config.Connection.Username != "root" {
		t.Errorf("Wrong username `%s` should be `%s`", d.Config.Connection.Username, "root")
	}

	if d.Config.Connection.Password != "root" {
		t.Errorf("Wrong password `%s` should be `%s`", d.Config.Connection.Password, "root")
	}

	if d.Config.ChangeSetPath != "test/resources/changes" {
		t.Errorf("Wrong changeSetPath `%s` should be `%s`", d.Config.ChangeSetPath, "test/resources/changes")
	}

	if d.Config.DatabaseType != "mysql" {
		t.Errorf("Wrong database type `%s` should be `%s`", d.Config.DatabaseType, "mysql")
	}

	if d.Config.BasePackage != "testbp" {
		t.Errorf("Wrong basePackage `%s` should be `%s`", d.Config.BasePackage, "testbp")
	}
}
func TestDvcMissingArguments(t *testing.T) {

	_, e := NewDVC(&Config{})

	if e == nil {
		t.Error("No error was thrown with no arguments")
		return
	}

	if e.Error() != "not enough arguments" {
		t.Error("NewDVC should have returned a `not enough arguments` error")
	}
}

// func TestCompareSchema(t *testing.T) {
// }

// func TestCompareSchema(t *testing.T) {}
