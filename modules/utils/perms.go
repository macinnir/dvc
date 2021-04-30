package utils

// Permission is the name of a permission
type Permission string

type IUserContainer interface {
	ID() int64
	Activated() bool
	Disabled() bool
	Locked() bool
	Permissions() []string
}

func HasPerm(user IUserContainer, perm Permission) bool {

	// System user
	if user.ID() == 1 {
		return true
	}

	if !user.Activated() || user.Disabled() || user.Locked() {
		return false
	}

	permissions := user.Permissions()
	for k := range permissions {
		if permissions[k] == string(perm) {
			return true
		}
	}

	return false
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
