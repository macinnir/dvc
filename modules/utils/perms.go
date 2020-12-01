package utils

// Permission is the name of a permission
type Permission string

// HasPerm verifies that a permission exists in a userProfile's permissions
func HasPerm(userID int64, perms []string, permName Permission) bool {

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
