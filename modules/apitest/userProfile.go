package apitest

import (
	"errors"
	"fmt"
	"log"
)

const (
	defaultProfileName = "profile0"
)

// UserProfile contains a user-contextual data structure for use within the testing process
type UserProfile struct {
	authKey string
	name    string
	id      int
	strings map[string]string
	ints    map[string]int64
	objects map[string]interface{}
}

// NewUserProfile returns a new UserProfile
func NewUserProfile(id int, name string) *UserProfile {
	userProfile := &UserProfile{
		id:      id,
		name:    name,
		strings: map[string]string{},
		ints:    map[string]int64{},
		objects: map[string]interface{}{},
	}
	return userProfile
}

// NewProfile creates a user profile indexed by `name`
func (a *APITest) NewProfile(name string) (e error) {

	if _, ok := a.userProfiles[name]; ok {
		e = fmt.Errorf("Duplicate profile name `%s`", name)
		return
	}

	a.userProfiles[name] = NewUserProfile(len(a.userProfiles), name)
	return
}

// DestroyProfile destroys a profile
// - Cannot destroy the default profile
// - Cannot destroy the profile if it is the currently active profile
func (a *APITest) DestroyProfile(profileName string) (e error) {
	if _, ok := a.userProfiles[profileName]; !ok {
		e = fmt.Errorf("DestroyProfile: profile `%s` does not exist", profileName)
		return
	}

	if profileName == defaultProfileName {
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

// SetAuthKey sets the auth key for the active user profile
func (a *APITest) SetAuthKey(authKey string) {
	log.Printf("Setting auth key to %s", authKey)
	a.userProfiles[a.activeProfile].authKey = authKey
}

// GetAuthKey gets the auth key for the active user profile
func (a *APITest) GetAuthKey() string {
	return a.userProfiles[a.activeProfile].authKey
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

// SetActiveProfile sets the active profile
func (a *APITest) SetActiveProfile(profileName string) {
	profile, e := a.GetProfile(profileName)
	if e != nil {
		panic(fmt.Sprintf("SetActiveProfile: Unknown profile `%s`", profileName))
	}
	a.activeProfile = profile.name
}

// ResetProfileToDefault sets the activeProfile to the default profile
func (a *APITest) ResetProfileToDefault() {
	a.activeProfile = defaultProfileName
}

// GetActiveProfile gets the active profile by its name
func (a *APITest) GetActiveProfile() *UserProfile {
	profile, _ := a.GetProfile(a.activeProfile)
	return profile
}
