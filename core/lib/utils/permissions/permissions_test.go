package permissions_test

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/macinnir/dvc/core/lib/utils/permissions"
	"github.com/stretchr/testify/assert"
)

func TestBuildPerm(t *testing.T) {
	p := permissions.BuildPerm(1, 2, 3)
	assert.Equal(t, 1030201, p)
}

func TestFetchSectionFromPerm(t *testing.T) {
	p := 1010101
	s := permissions.FetchSectionFromPerm(p)
	assert.Equal(t, 1, s, "should return section # 1")
}

func TestFetchGroupFromPerm(t *testing.T) {
	p := 1010101
	s := permissions.FetchGroupFromPerm(p)
	assert.Equal(t, 1, s)

	p2 := 1030402
	s2 := permissions.FetchGroupFromPerm(p2)
	assert.Equal(t, 4, s2)
}

func TestFetchBasePermFromPerm(t *testing.T) {
	p := 1130402
	s := permissions.FetchBasePermFromPerm(p)
	assert.Equal(t, 13, s)
}

func TestHasPerm(t *testing.T) {
	perms := 1 + 2 + 4 + 8
	perm := 4
	assert.True(t, permissions.HasPerm(perm, perms))
}

func TestToBWPart(t *testing.T) {
	sectionBW := permissions.ToBWPart(3)
	assert.Equal(t, 8, sectionBW)
}

func TestUserPermission_AddPerm(t *testing.T) {

	section := 10
	group := 7

	up := permissions.NewUserPermissions()
	up.AddPerm(1090710)

	sectionBW := int(math.Pow(2, float64(section)))

	assert.Equal(t, sectionBW, up.Sections&sectionBW)

	if _, ok := up.Permissions[sectionBW]; !ok {
		assert.FailNow(t, fmt.Sprintf("Section #%d not found", sectionBW))
	}

	if _, ok := up.Permissions[sectionBW][group]; !ok {
		assert.FailNow(t, fmt.Sprintf("Group #%d not found", group))
	}

	assert.Equal(t, up.Permissions[sectionBW][group], 512)

}

func TestHasPermInSection(t *testing.T) {

	// perm is the permission that the method we are requesting has
	// var perm int = 1001

	// userPerms is the slice of permissions the user has
	var userPerms = permissions.NewUserPermissions([]int{1000910, 1040810, 1060408}...)

	assert.Equal(t, 1024, userPerms.Sections)
	assert.Equal(t, 1, userPerms.Permissions[1024][9])
	assert.True(t, userPerms.HasPerm(1000910))
	assert.True(t, userPerms.HasPerm(1040810))
	assert.True(t, userPerms.HasPerm(1060408))

	b, _ := json.Marshal(userPerms)
	fmt.Println(string(b))

	// assert.Equal(t, 8796093022208, userPerms.Permissions[128])

	// assert.True(t, userPerms.HasPerm(1001))
	// assert.True(t, userPerms.HasPerm(2001))
	// assert.True(t, userPerms.HasPerm(1002))
	// assert.True(t, userPerms.HasPerm(2002))
	// assert.True(t, userPerms.HasPerm(8796093022208008))

	// assert.False(t, userPerms.HasPerm(4002))
}
