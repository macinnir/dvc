package permissions

import (
	"fmt"
	"math"
)

const (
	sectionDivisor = 100
	groupDivisor   = 10000
	totalDivisor   = 1000000
)

func BuildPerm(section, group, perm int) int {

	return totalDivisor + ((perm) * groupDivisor) + ((group) * sectionDivisor) + (section)
}

func HasPerm(perm int, perms int) bool {
	return perm&perms == perm
}

// FetchSectionFromPerm returns the section from the permission
func FetchSectionFromPerm(perm int) int {
	return perm % sectionDivisor
}

// FetchGroupFromPerm returns the group from the permission
func FetchGroupFromPerm(perm int) int {
	section := FetchSectionFromPerm(perm)
	return ((perm % groupDivisor) - section) / sectionDivisor
}

// FetchBasePermFromPerm returns the base permission number from the permission
func FetchBasePermFromPerm(perm int) int {
	return (perm % totalDivisor) / groupDivisor
}

type UserPermissions struct {
	// Sections is a number representing the sections that are utilized for this user's permissions
	Sections    int                 `json:"S"`
	Permissions map[int]map[int]int `json:"P"`
}

func NewUserPermissions(perms ...int) *UserPermissions {

	up := &UserPermissions{
		Sections:    0,
		Permissions: map[int]map[int]int{},
	}

	up.AddPerm(perms...)

	return up
}

func (up *UserPermissions) AddPerm(fullPerm ...int) {

	for k := range fullPerm {
		up.addPerm(fullPerm[k])
	}

}

func (up *UserPermissions) addPerm(fullPerm int) {

	var section = FetchSectionFromPerm(fullPerm)
	var group = FetchGroupFromPerm(fullPerm)
	var perm = FetchBasePermFromPerm(fullPerm)

	var sectionBW = ToBWPart(section)
	var permBW = ToBWPart(perm)

	// Check if the section exists
	fmt.Println("Sections: ", section, "(", sectionBW, ")", "&", up.Sections, " == 0 (", sectionBW&up.Sections, ")")
	// fmt.Println(section&up.Sections, section, up.Sections)
	if sectionBW&up.Sections == 0 {
		up.Sections += sectionBW
		fmt.Println("Adding section", sectionBW, "==>", up.Sections)
	}

	if _, ok := up.Permissions[sectionBW]; !ok {
		up.Permissions[sectionBW] = map[int]int{}
	}

	// Add if it does not exist
	if permBW&up.Permissions[sectionBW][group] == 0 {
		up.Permissions[sectionBW][group] += permBW
	}
}

// ToBWPart returns the bitwise value for the base10 part number
func ToBWPart(part int) int {
	return int(math.Pow(2, float64(part)))
}

func (up *UserPermissions) RemovePerm(fullPerm int) {

	var section = FetchSectionFromPerm(fullPerm)
	var group = FetchGroupFromPerm(fullPerm)
	var perm = FetchBasePermFromPerm(fullPerm)

	var sectionBW = ToBWPart(section)
	var permBW = ToBWPart(perm)

	// Check if the section exists
	fmt.Println("Sections: ", sectionBW, "&", up.Sections, " == 0 (", sectionBW&up.Sections, ")")
	// fmt.Println(section&up.Sections, section, up.Sections)
	if sectionBW&up.Sections == 0 {
		up.Sections += sectionBW
		fmt.Println("Adding section", sectionBW, "==>", up.Sections)
	}

	if _, ok := up.Permissions[sectionBW]; !ok {
		up.Permissions[sectionBW] = map[int]int{}
	}

	// Remove if it exists
	if permBW&up.Permissions[sectionBW][group] == permBW {
		up.Permissions[sectionBW][group] -= permBW
	}

}

func (up *UserPermissions) HasPerm(fullPerm int) bool {

	var section = ToBWPart(FetchSectionFromPerm(fullPerm))
	var group = FetchGroupFromPerm(fullPerm)
	var perm = ToBWPart(FetchBasePermFromPerm(fullPerm))

	fmt.Println("HasPerm; Section:", section, "Group:", group, "Perm:", perm)

	if section&up.Sections != section {
		fmt.Println("No Section", section)
		return false
	}

	if _, ok := up.Permissions[section]; !ok {
		fmt.Println("No permissions for section", section)
		return false
	}

	if _, ok := up.Permissions[section][group]; !ok {
		fmt.Println("No permissions for group", group)
		return false
	}

	fmt.Println(fullPerm, "==>", perm, "&", up.Permissions[section][group], "==", perm)
	return perm&up.Permissions[section][group] == perm
}
