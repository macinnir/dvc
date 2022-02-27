package utils_test

import (
	"testing"

	"github.com/macinnir/dvc/core/lib/utils"
	"github.com/macinnir/dvc/core/lib/utils/request"
	"github.com/macinnir/dvc/core/lib/utils/types"
	"github.com/stretchr/testify/assert"
)

type testUser struct {
	id          int64
	activated   bool
	disabled    bool
	locked      bool
	permissions []string
	account     int64
}

func (tu *testUser) ID() int64 {
	return tu.id
}

func (tu *testUser) Activated() bool {
	return tu.activated
}

func (tu *testUser) Disabled() bool {
	return tu.disabled
}

func (tu *testUser) Locked() bool {
	return tu.locked
}

func (tu *testUser) Permissions() []string {
	return tu.permissions
}

func (tu *testUser) Account() int64 {
	return tu.account
}

func (tu *testUser) SettingMgr() *types.SettingsManager {
	return nil
}

func TestHasPerm(t *testing.T) {
	req := &request.Request{
		Method: "GET",
		Path:   "/foo/2/bar",
		Headers: map[string]string{
			utils.RequestPathUserIDArgName: "1",
		},
		Params:     map[string]string{},
		Body:       "",
		ActionType: 0,
	}

	tu := &testUser{
		id:          2,
		activated:   true,
		disabled:    false,
		locked:      false,
		permissions: []string{"someFeature_somePermission"},
	}

	assert.True(t, utils.HasPerm(req, tu, "someFeature_somePermission"))
}

func TestHasPermAsSuperuser(t *testing.T) {

	req := &request.Request{
		Method: "GET",
		Path:   "/some/feature",
	}

	tu := &testUser{
		id:        utils.SuperUserID,
		activated: true,
	}

	assert.True(t, utils.HasPerm(req, tu, "someFeature_somePermission"), "Super user requires no permissions")
}

func TestHasPermByFeature(t *testing.T) {
	req := &request.Request{
		Method: "GET",
		Path:   "/some/feature",
	}

	tu := &testUser{
		activated:   true,
		permissions: []string{"someFeature_*"},
	}

	assert.True(t, utils.HasPerm(req, tu, "someFeature_somePermission"))
}

func TestHasPermAsOwner(t *testing.T) {
	req := &request.Request{
		Method:  "GET",
		Path:    "/foo/2/bar",
		Headers: map[string]string{},
		Params: map[string]string{
			utils.RequestPathUserIDArgName: "2",
		},
		UserID:     2,
		Body:       "",
		ActionType: 0,
	}

	tu := &testUser{
		id:          2,
		activated:   true,
		disabled:    false,
		locked:      false,
		permissions: []string{"someFeature_somePermissionAsOwner"},
	}

	assert.True(t, utils.HasPerm(req, tu, "someFeature_somePermissionAsOwner"))
}

func TestHasPermAsNotOwner(t *testing.T) {
	req := &request.Request{
		Method:  "GET",
		Path:    "/foo/3/bar",
		Headers: map[string]string{},
		Params: map[string]string{
			utils.RequestPathUserIDArgName: "3",
		},
		UserID:     2,
		Body:       "",
		ActionType: 0,
	}

	tu := &testUser{
		id:          2,
		activated:   true,
		disabled:    false,
		locked:      false,
		permissions: []string{"someFeature_somePermissionAsOwner"},
	}

	assert.False(t, utils.HasPerm(req, tu, "someFeature_somePermissionAsOwner"))
}

// 11.9 ns/op 0 B/op 0 allocs/op
func BenchmarkHasPerm(b *testing.B) {
	req := &request.Request{
		Method: "GET",
		Path:   "/foo/2/bar",
		Headers: map[string]string{
			utils.RequestPathUserIDArgName: "1",
		},
		Params:     map[string]string{},
		Body:       "",
		ActionType: 0,
	}

	tu := &testUser{
		id:          2,
		activated:   true,
		disabled:    false,
		locked:      false,
		permissions: []string{"someFeature_somePermission"},
	}

	for n := 0; n < b.N; n++ {
		// fmt.Println("running", n)
		utils.HasPerm(req, tu, "someFeature_somePermission")
	}
}

// 32.6 ns/op  0 B/op 0 allocs/op
func BenchmarkHasPermAsOwner(b *testing.B) {
	req := &request.Request{
		Method:  "GET",
		Path:    "/foo/2/bar",
		Headers: map[string]string{},
		Params: map[string]string{
			utils.RequestPathUserIDArgName: "2",
		},
		UserID:     2,
		Body:       "",
		ActionType: 0,
	}

	tu := &testUser{
		id:          2,
		activated:   true,
		disabled:    false,
		locked:      false,
		permissions: []string{"someFeature_somePermissionAsOwner"},
	}

	for n := 0; n < b.N; n++ {
		utils.HasPerm(req, tu, "someFeature_somePermissionAsOwner")
	}
}

// 31.3 ns/op 0 B/op 0 allocs/op
// 108  ns/op 1 B/op 1 allocs/op
// 38.2 ns/op 0 B/op 0 allocs/op -- Remove strings.Split
// go test -benchmem -run=^$ -bench ^(BenchmarkHasPermAsNotOwner)$ github.com/macinnir/dvc/core/lib/utils -v
func BenchmarkHasPermAsNotOwner(b *testing.B) {
	req := &request.Request{
		Method:  "GET",
		Path:    "/foo/3/bar",
		Headers: map[string]string{},
		Params: map[string]string{
			utils.RequestPathUserIDArgName: "3",
		},
		UserID:     2,
		Body:       "",
		ActionType: 0,
	}

	tu := &testUser{
		id:          2,
		activated:   true,
		disabled:    false,
		locked:      false,
		permissions: []string{"someFeature_somePermissionAsOwner"},
	}

	for n := 0; n < b.N; n++ {
		utils.HasPerm(req, tu, "someFeature_somePermissionAsOwner")
	}

}
