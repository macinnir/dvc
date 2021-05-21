package utils

import (
	"github.com/macinnir/dvc/core/lib/utils/request"
)

// Permission is the name of a permission
type Permission string

const (
	SuperUserID              = int64(1)
	RequestPathUserIDArgName = "userID"
	AsOwnerSuffix            = "AsOwner"
)

type IUserContainer interface {
	ID() int64
	Activated() bool
	Disabled() bool
	Locked() bool
	Permissions() []string
}

func HasPerm(req *request.Request, user IUserContainer, perm Permission) bool {

	// System user
	if user.ID() == SuperUserID {
		return true
	}

	if !user.Activated() || user.Disabled() || user.Locked() {
		return false
	}

	hasPerm := false

	permissions := user.Permissions()
	for k := range permissions {
		if permissions[k] == string(perm) {
			hasPerm = true
		}
	}

	if !hasPerm {
		return false
	}

	// Check suffix
	if len(perm) <= len(AsOwnerSuffix) {
		return hasPerm
	}

	suffix := string(perm[len(perm)-len(AsOwnerSuffix):])

	// Check if this permission is "AsOwner"
	if suffix == AsOwnerSuffix {

		// Check if the request contains a "userID" argument and that it matches the current user
		return req.ArgInt64(RequestPathUserIDArgName, 0) == req.UserID
	}

	return hasPerm
}

// HasPerm verifies that a permission exists in a userProfile's permissions
func HasPermOld(userID int64, perms []string, permName Permission) bool {

	// Superuser
	if userID == 1 {
		return true
	}

	for k := range perms {
		if perms[k] == string(permName) {
			return true
		}
	}

	return false

	// // Check if the device has been registered
	// if user.Device == nil || user.Device.DateRegistered == 0 {
	// 	return false
	// }

	// if len(user.UserProfile.Permissions.String) > 0 {

	// 	// God mode catch-all
	// 	if user.UserProfile.Permissions.String == "*" {
	// 		return true
	// 	}

	// 	if strings.Contains(user.UserProfile.Permissions.String, "#"+string(permName)+"#") {
	// 		return true
	// 	}

	// }

	// // Iterate through all roles to see if they have the permission
	// for k := range user.UserProfile.Roles {

	// 	if strings.Contains(user.UserProfile.Roles[k].Permissions.String, "#"+string(permName)+"#") {
	// 		return true
	// 	}

	// }

	// return false
}
