package apitest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const url = "http://someurl.com"

var at *APITest

func TestMain(t *testing.T) {
	at = NewAPITest(t, "http://someurl.com")
}

func TestGetProfile_ShouldReturnDefault(t *testing.T) {
	profile := at.GetActiveProfile()
	assert.Equal(t, defaultProfileName, profile.name)
	e := at.NewProfile("profile_test_123")
	assert.Nil(t, e)

	at.SetActiveProfile("profile_test_123")
	profile = at.GetActiveProfile()
	assert.Equal(t, "profile_test_123", profile.name)

	at.ResetProfileToDefault()
	at.DestroyProfile("profile_test_123")
}

func TestDestroyProfile(t *testing.T) {
	at.NewProfile("profile_test_to_be_destroyed")
	assert.Equal(t, 2, len(at.GetProfileNames()))

	at.DestroyProfile("profile_test_to_be_destroyed")
	assert.Equal(t, 1, len(at.GetProfileNames()))
}

func TestDestroyProfile_ShouldNotAllowDeletingDefaultProfile(t *testing.T) {
	e := at.DestroyProfile(defaultProfileName)
	assert.NotNil(t, e)
	assert.Equal(t, 1, len(at.GetProfileNames()))
}

func TestDestroyProfile_ShouldNotAllowDeletingCurrentProfile(t *testing.T) {
	at.NewProfile("new_profile")
	at.SetActiveProfile("new_profile")
	e := at.DestroyProfile("new_profile")
	assert.NotNil(t, e)
	assert.Equal(t, 2, len(at.GetProfileNames()))
	at.ResetProfileToDefault()
	e = at.DestroyProfile("new_profile")
	assert.Nil(t, e)
	assert.Equal(t, 1, len(at.GetProfileNames()))
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

func TestSetAuthKey(t *testing.T) {
	at.SetAuthKey("12345")
	assert.Equal(t, "12345", at.GetAuthKey())
}

func TestCreateProfile(t *testing.T) {
	e := at.NewProfile("user1")
	assert.Nil(t, e)

	profileNames := at.GetProfileNames()
	assert.Equal(t, 2, len(profileNames))
}

func TestSetString_ShouldBeIsolatedToProfile(t *testing.T) {
	at.SetString("user0_string", "siloedString")
	at.SetActiveProfile("user1")
	val := at.GetString("user0_string")
	assert.Equal(t, "", val)
}

func TestSetInt_ShouldBeIsolatedToProfile(t *testing.T) {
	at.ResetProfileToDefault()
	at.SetInt("user0_int", 123)
	at.SetActiveProfile("user1")
	val := at.GetInt("user0_int")
	assert.Equal(t, int64(-1), val)
}

func TestSetGlobalString(t *testing.T) {

	val := at.SetGlobalString("global_string", "global_string_val")
	assert.Equal(t, "global_string_val", val)

	profileVal := at.GetString("global_string")
	assert.Equal(t, "", profileVal)

	val = at.GetGlobalString("global_string")
	assert.Equal(t, "global_string_val", val)
}

func TestSetIntGlobal_ShouldBeSiloedFromProfiles(t *testing.T) {

	val := at.SetIntGlobal("global_int", 123)
	assert.Equal(t, int64(123), val)

	profileVal := at.GetInt("global_int")
	assert.Equal(t, int64(-1), profileVal)

	val = at.GetIntGlobal("global_int")
	assert.Equal(t, int64(123), val)
}

func TestGetGlobalString_ShouldReturnEmptyStringIfNotSet(t *testing.T) {
	val := at.GetGlobalString("not set")
	assert.Equal(t, "", val)
}

func TestGetGlobalInt_ShouldReturnEmptyStringIfNotSet(t *testing.T) {
	val := at.GetIntGlobal("not set")
	assert.Equal(t, int64(-1), val)
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

func TestIncrementGlobal(t *testing.T) {
	val := at.SetIntGlobal("counter", 1)
	val = at.IncrementGlobal("counter")
	assert.Equal(t, int64(2), val)
	at.IncrementGlobal("counter")
	val = at.GetIntGlobal("counter")
	assert.Equal(t, int64(3), val)
}

func TestIncrementGlobal_ShouldReturn1IfNotExist(t *testing.T) {
	val := at.IncrementGlobal("counter3_does_not_exist")
	assert.Equal(t, int64(1), val)
	val = at.GetIntGlobal("counter3_does_not_exist")
	assert.Equal(t, int64(1), val)
}

func TestDecrementGlobal(t *testing.T) {
	val := at.SetIntGlobal("counter", 1)
	val = at.DecrementGlobal("counter")
	assert.Equal(t, int64(0), val)
	at.DecrementGlobal("counter")
	val = at.GetIntGlobal("counter")
	assert.Equal(t, int64(-1), val)
}

func TestDecrementGlobal_ShouldReturnNegativeIfNotExist(t *testing.T) {
	val := at.DecrementGlobal("counter4_does_not_exist")
	assert.Equal(t, int64(-1), val)
	val = at.GetIntGlobal("counter4_does_not_exist")
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

func TestSetObjectGlobal(t *testing.T) {
	obj := map[string]string{
		"foo": "bar",
	}

	at.SetObjectGlobal("testObj1", obj)

	objRet := at.GetObjectGlobal("testObj1").(map[string]string)
	assert.Equal(t, objRet["foo"], "bar")
}
