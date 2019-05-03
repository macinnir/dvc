package apitest

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

const (
	DefaultProfileName = "profile0"
	AuthKey            = "authKey"
	RefreshTokenKey    = "refreshToken"
)

// UserProfile contains a user-contextual data structure for use within the testing process
type UserProfile struct {
	sync.Mutex
	baseURL  string
	Name     string
	ID       int
	strings  map[string]string
	ints     map[string]int64
	objects  map[string]interface{}
	Requests *Requests
	Logger   *Logger
}

// GetInt gets the int value at profile `profileName` and `key` without
// changing the current profile
func (u *UserProfile) GetInt(key string) (val int64) {

	var ok bool

	if val, ok = u.ints[key]; !ok {
		val = 0
	}

	return val
}

// SetString sets a string value at `key` within the active profile
func (u *UserProfile) SetString(key, value string) string {
	// u.Logger.Debug(fmt.Sprintf("About to set %s ==> %s, %v \n", key, value, u))
	u.strings[key] = value
	return value
}

// SetAuthKey sets the authentication key for network requests
func (u *UserProfile) SetAuthKey(value string) string {
	u.strings[AuthKey] = value
	return value
}

// GetString returns a string value in profile `profileName` at `key` without changing the active profile
func (u *UserProfile) GetString(key string) string {

	if _, ok := u.strings[key]; !ok {
		return ""
	}

	return u.strings[key]
}

// SetInt sets an int value at `name` for profile `profileName`
func (u *UserProfile) SetInt(name string, val int64) int64 {
	u.ints[name] = val
	return val
}

// GetObject gets an object at `key` for profile `profileName`
func (u *UserProfile) GetObject(key string) (val interface{}) {

	var ok bool

	if val, ok = u.objects[key]; !ok {
		return nil
		// a.logger.Error(fmt.Sprintf("Invalid object name `%s` in profile `%s`", key, a.activeProfile))
	}

	return
}

// SetObject sets an object at `name` for the active user profile
func (u *UserProfile) SetObject(key string, obj interface{}) {
	u.objects[key] = obj
}

// Decrement sets an integer in the number cache if it doesn't exist
// And then decrements it by 1
// The resulting value is returned
func (u *UserProfile) Decrement(key string) int64 {

	val, ok := u.ints[key]

	if !ok {
		val = 0
	}

	val--
	u.ints[key] = val
	return val
}

// Increment sets an integer in the number cache for profile `profileName` if it doesn't exist
// And then increments it by 1
// The resulting value is returned
func (u *UserProfile) Increment(key string) int64 {

	val, ok := u.ints[key]
	if !ok {
		val = 0
	}

	val++

	u.ints[key] = val
	return val
}

// NewUserProfile returns a new UserProfile
func NewUserProfile(id int, name string, logger *Logger, t *testing.T) *UserProfile {
	userProfile := &UserProfile{
		ID:      id,
		Name:    name,
		strings: map[string]string{},
		ints:    map[string]int64{},
		objects: map[string]interface{}{},
		Logger:  logger,
	}

	userProfile.Requests = NewRequests(userProfile, logger, t)
	return userProfile
}

// NewProfile creates a user profile indexed by `name`
func (a *APITest) NewProfile(name string) (userProfile *UserProfile, e error) {

	if _, ok := a.userProfiles[name]; ok {
		e = fmt.Errorf("Duplicate profile name `%s`", name)
		return
	}

	a.userProfiles[name] = NewUserProfile(len(a.userProfiles), name, a.logger, a.t)
	a.userProfiles[name].baseURL = a.baseURL

	return a.userProfiles[name], nil
}

// DestroyProfile destroys a profile
// - Cannot destroy the default profile
// - Cannot destroy the profile if it is the currently active profile
func (a *APITest) DestroyProfile(profileName string) (e error) {
	if _, ok := a.userProfiles[profileName]; !ok {
		e = fmt.Errorf("DestroyProfile: profile `%s` does not exist", profileName)
		return
	}

	if profileName == DefaultProfileName {
		e = errors.New("DestroyProfile: cannot delete default profile")
		return
	}

	if profileName == a.activeProfile {
		e = fmt.Errorf("DestroyProfile: cannot delete the currently active profile")
		return
	}

	delete(a.userProfiles, profileName)

	return
}

// GetProfileNames returns a list of profile names
func (a *APITest) GetProfileNames() []string {

	profileNames := []string{}

	for profileName := range a.userProfiles {
		profileNames = append(profileNames, profileName)
	}

	return profileNames
}

// GetProfile gets a profile by its name
func (a *APITest) GetProfile(profileName string) (profile *UserProfile, e error) {

	var ok bool

	if profile, ok = a.userProfiles[profileName]; !ok {
		e = fmt.Errorf("Undefined profile name `%s`", profileName)
		return
	}

	return
}

// CheckProfile checks a profile
func (a *APITest) CheckProfile(profileName string) {
	if _, ok := a.userProfiles[profileName]; !ok {
		panic(fmt.Sprintf("Undefined profile name `%s`", profileName))
	}
}

// SetActiveProfile sets the active profile
func (a *APITest) SetActiveProfile(profileName string) {
	profile, e := a.GetProfile(profileName)
	if e != nil {
		panic(fmt.Sprintf("SetActiveProfile: Unknown profile `%s`", profileName))
	}
	a.activeProfile = profile.Name
}

// ResetProfileToDefault sets the activeProfile to the default profile
func (a *APITest) ResetProfileToDefault() {
	a.activeProfile = DefaultProfileName
}
