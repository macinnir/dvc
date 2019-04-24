package apitest

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// APITest tests your API
type APITest struct {
	// The injected testing object
	t *testing.T
	// The base URL used for all API calls
	baseURL string
	// UserProfiles
	userProfiles map[string]*UserProfile
	// The name of the active UserProfile
	activeProfile string
	// Map of global string values accessible across all profiles
	stringVals map[string]string
	// Map of global int64 values accessible across all profiles
	intVals map[string]int64
	// Map of global object values accessible across all profiles
	objectVals map[string]interface{}
	// Unique key for the current testing session
	sessionKey string
}

// NewAPITest returns a new APITest instance
func NewAPITest(t *testing.T, baseURL string) *APITest {

	apiTest := &APITest{
		t:            t,
		baseURL:      baseURL,
		userProfiles: map[string]*UserProfile{},
		stringVals:   map[string]string{},
		intVals:      map[string]int64{},
		objectVals:   map[string]interface{}{},
	}

	apiTest.init()

	return apiTest
}

// DoAs runs a function (fn) within the context of profile `profileName`
// then sets the context back to the previous profile
func (a *APITest) DoAs(profileName string, fn func()) {

	oldProfile := a.GetActiveProfile()
	a.SetActiveProfile(profileName)
	fn()
	a.SetActiveProfile(oldProfile.name)
}

func (a *APITest) init() {
	rand.Seed(time.Now().UnixNano())
	a.sessionKey = RandString(10)

	// Create the first user profile
	a.NewProfile(defaultProfileName)
	a.SetActiveProfile(defaultProfileName)
}

/**
 * Strings
 */

// SetString sets a string value at `key` within the active profile
func (a *APITest) SetString(key, value string) string {
	a.userProfiles[a.activeProfile].strings[key] = value
	return value
}

// SetGlobalString sets a string value at `key` globally
func (a *APITest) SetGlobalString(key, value string) string {
	a.stringVals[key] = value
	return value
}

// GetStringAt returns a string value in profile `profileName` at `key` without changing the active profile
func (a *APITest) GetStringAt(profileName, key string) string {
	return a.userProfiles[profileName].strings[key]
}

// GetString returns a string value at `key` within the active profile
func (a *APITest) GetString(key string) string {
	return a.GetStringAt(a.activeProfile, key)
}

// GetGlobalString gets the string value identified by `name`
func (a *APITest) GetGlobalString(key string) string {
	val, ok := a.stringVals[key]
	if !ok {
		return ""
	}

	return val
}

// SetRandString generates a random string of length `length`,
// sets it at `key` for the active user profile
func (a *APITest) SetRandString(key string, length int) string {
	randString := RandString(length)
	return a.SetString(key, randString)
}

// SetRandGlobalString generates a random string of length `length`,
// sets it at `key` globally
func (a *APITest) SetRandGlobalString(key string, length int) string {
	randString := RandString(length)
	return a.SetGlobalString(key, randString)
}

// RandEmail returns a random email
func (a *APITest) RandEmail() string {
	profile := a.GetActiveProfile()
	// return a.RandString(rand.Intn(10)) + "@" + a.RandString(rand.Intn(10)) + "." + a.RandDataString("domainSuffixes")
	return RandString(rand.Intn(10)) + fmt.Sprintf("@%s_prof_%d", a.sessionKey, profile.id) + a.RandDataString("domainSuffixes")
}

/**
 * Ints
 */

// SetIntGlobal sets an int value at `name` globally
func (a *APITest) SetIntGlobal(name string, v int64) int64 {
	a.intVals[name] = v
	return v
}

// SetInt sets an int value at `name` for the current profile
func (a *APITest) SetInt(name string, v int64) int64 {
	a.userProfiles[a.activeProfile].ints[name] = v
	return v
}

// GetIntGlobal gets the int value at `name` globally
func (a *APITest) GetIntGlobal(name string) int64 {
	val, ok := a.intVals[name]
	if !ok {
		return -1
	}

	return val
}

// GetInt gets the int value at `key` for the current profile
func (a *APITest) GetInt(key string) (val int64) {
	return a.GetIntAt(a.activeProfile, key)
}

// GetIntAt gets the int value at profile `profileName` and `key` without
// changing the current profile
func (a *APITest) GetIntAt(profileName string, key string) (val int64) {

	var ok bool

	if val, ok = a.userProfiles[profileName].ints[key]; !ok {
		return -1
	}

	return val
}

// SetRandInt generates a random number where 0 <= n <= max,
// sets it to the number cache and returns the number
func (a *APITest) SetRandInt(name string, max int) int64 {
	randNum := RandInt(max)
	a.SetInt(name, randNum)
	return randNum
}

// Increment sets an integer in the number cache if it doesn't exist
// And then increments it by 1
// The resulting value is returned
func (a *APITest) Increment(key string) int64 {

	val, ok := a.userProfiles[a.activeProfile].ints[key]
	if !ok {
		val = 0
	}

	val++

	a.userProfiles[a.activeProfile].ints[key] = val

	return val
}

// IncrementGlobal sets an integer in the number cache if it doesn't exist
// And then decrements it by 1
// The resulting value is returned
func (a *APITest) IncrementGlobal(key string) int64 {

	val, ok := a.intVals[key]

	if !ok {
		val = 0
	}

	val++

	a.intVals[key] = val

	return val
}

// Decrement sets an integer in the number cache if it doesn't exist
// And then decrements it by 1
// The resulting value is returned
func (a *APITest) Decrement(key string) int64 {

	val, ok := a.userProfiles[a.activeProfile].ints[key]

	if !ok {
		val = 0
	}

	val--
	a.userProfiles[a.activeProfile].ints[key] = val

	return val
}

// DecrementGlobal sets an integer in the number cache if it doesn't exist
// And then decrements it by 1
// The resulting value is returned
func (a *APITest) DecrementGlobal(key string) int64 {

	val, ok := a.intVals[key]

	if !ok {
		val = 0
	}

	val--

	a.intVals[key] = val

	return val
}

/**
 * Objects
 */

// SetObjectGlobal sets an object at `name` globally
func (a *APITest) SetObjectGlobal(key string, obj interface{}) {
	a.objectVals[key] = obj
}

// SetObject sets an object at `name` for the active user profile
func (a *APITest) SetObject(key string, obj interface{}) {
	a.userProfiles[a.activeProfile].objects[key] = obj
}

// GetObjectGlobal gets an object at `key` globally
func (a *APITest) GetObjectGlobal(name string) interface{} {
	val, ok := a.objectVals[name]
	if !ok {
		panic(fmt.Sprintf("Invalid object name `%s`", name))
	}

	return val
}

// GetObject gets an object at `key` for the active user profile
func (a *APITest) GetObject(key string) (val interface{}) {

	var ok bool

	if val, ok = a.userProfiles[a.activeProfile].objects[key]; !ok {
		panic(fmt.Sprintf("Invalid object name `%s` in profile `%s`", key, a.activeProfile))
	}

	return
}

// RandDataString returns a random data string from `name` set in dataStrings
func (a *APITest) RandDataString(name string) string {
	return dataStrings[name][rand.Intn(len(dataStrings[name]))]
}
