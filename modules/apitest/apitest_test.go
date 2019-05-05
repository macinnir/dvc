package apitest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const url = "http://someurl.com"

var at *APITest

func TestMain(t *testing.T) {
	at = NewAPITest(t, "http://someurl.com", LogLevelError)
}

func TestGetProfile_ShouldReturnDefault(t *testing.T) {
	profile, e := at.GetProfile(DefaultProfileName)
	assert.Nil(t, e)
	assert.Equal(t, DefaultProfileName, profile.Name)
	_, e = at.NewProfile("profile_test_123")
	assert.Nil(t, e)
}

func TestDestroyProfile(t *testing.T) {
	numProfiles := len(at.GetProfileNames())
	at.NewProfile("profile_test_to_be_destroyed")
	assert.Equal(t, numProfiles+1, len(at.GetProfileNames()))

	at.DestroyProfile("profile_test_to_be_destroyed")
	assert.Equal(t, numProfiles, len(at.GetProfileNames()))
}

func TestDestroyProfile_ShouldNotAllowDeletingDefaultProfile(t *testing.T) {
	numProfiles := len(at.GetProfileNames())
	e := at.DestroyProfile(DefaultProfileName)
	assert.NotNil(t, e)
	assert.Equal(t, numProfiles, len(at.GetProfileNames()))
}

func TestDestroyProfile_ShouldNotAllowDeletingCurrentProfile(t *testing.T) {
	numProfiles := len(at.GetProfileNames())
	at.NewProfile("new_profile")
	at.SetActiveProfile("new_profile")
	e := at.DestroyProfile("new_profile")
	assert.NotNil(t, e)
	assert.Equal(t, numProfiles+1, len(at.GetProfileNames()))
	at.ResetProfileToDefault()
	e = at.DestroyProfile("new_profile")
	assert.Nil(t, e)
	assert.Equal(t, numProfiles, len(at.GetProfileNames()))
}

func TestSetString(t *testing.T) {
	val := at.SetString("foo", "bar")
	assert.Equal(t, "bar", val)
	val = at.GetString("foo")
	assert.Equal(t, "bar", val)
}

func TestSetInt(t *testing.T) {
	val := at.SetInt("foo", 1)
	assert.Equal(t, int64(1), val)
	val = at.GetInt("foo")
	assert.Equal(t, int64(1), val)
}

func TestCreateProfile(t *testing.T) {
	numProfiles := len(at.GetProfileNames())
	_, e := at.NewProfile("user1")
	assert.Nil(t, e)

	profileNames := at.GetProfileNames()
	assert.Equal(t, numProfiles+1, len(profileNames))
}

func TestGetString_ShouldReturnEmptyStringIfNotSet(t *testing.T) {
	val := at.GetString("not set")
	assert.Equal(t, "", val)
}

func TestGetInt_ShouldReturnZeroIfNotSet(t *testing.T) {
	val := at.GetInt("not set")
	assert.Equal(t, int64(0), val)
}

func TestIncrement(t *testing.T) {
	val := at.SetInt("counter", 1)
	val = at.Increment("counter")
	assert.Equal(t, int64(2), val)
	at.Increment("counter")
	val = at.GetInt("counter")
	assert.Equal(t, int64(3), val)
}

func TestIncrement_ShouldReturn1IfNotExist(t *testing.T) {
	val := at.Increment("counter1_does_not_exist")
	assert.Equal(t, int64(1), val)
	val = at.GetInt("counter1_does_not_exist")
	assert.Equal(t, int64(1), val)
}

func TestDecrement(t *testing.T) {
	val := at.SetInt("counter", 1)
	val = at.Decrement("counter")
	assert.Equal(t, int64(0), val)
	at.Decrement("counter")
	val = at.GetInt("counter")
	assert.Equal(t, int64(-1), val)
}

func TestDecrement_ShouldReturnNegativeIfNotExist(t *testing.T) {
	val := at.Decrement("counter2_does_not_exist")
	assert.Equal(t, int64(-1), val)
	val = at.GetInt("counter2_does_not_exist")
	assert.Equal(t, int64(-1), val)
}

func TestSetObject(t *testing.T) {
	obj := map[string]string{
		"foo": "bar",
	}

	at.SetObject("testObj1", obj)

	objRet := at.GetObject("testObj1").(map[string]string)
	assert.Equal(t, objRet["foo"], "bar")
}
