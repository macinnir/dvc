package utils

import (
	"testing"

	"github.com/macinnir/dvc/core/lib/utils/request"
	"github.com/stretchr/testify/assert"
)

type testUser struct {
	id          int64
	activated   bool
	disabled    bool
	locked      bool
	permissions []string
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

func TestHasPerm(t *testing.T) {
	req := &request.Request{
		Method: "GET",
		Path:   "/foo/2/bar",
		Headers: map[string]string{
			RequestPathUserIDArgName: "1",
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
		permissions: []string{"somePermission"},
	}

	assert.True(t, HasPerm(req, tu, "somePermission"))
}

func TestHasPermAsOwner(t *testing.T) {
	req := &request.Request{
		Method:  "GET",
		Path:    "/foo/2/bar",
		Headers: map[string]string{},
		Params: map[string]string{
			RequestPathUserIDArgName: "2",
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
		permissions: []string{"somePermissionAsOwner"},
	}

	assert.True(t, HasPerm(req, tu, "somePermissionAsOwner"))
}

func TestHasPermAsNotOwner(t *testing.T) {
	req := &request.Request{
		Method:  "GET",
		Path:    "/foo/3/bar",
		Headers: map[string]string{},
		Params: map[string]string{
			RequestPathUserIDArgName: "3",
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
		permissions: []string{"somePermissionAsOwner"},
	}

	assert.False(t, HasPerm(req, tu, "somePermissionAsOwner"))
}

// 11.9 ns/op 0 B/op 0 allocs/op
func BenchmarkHasPerm(b *testing.B) {
	req := &request.Request{
		Method: "GET",
		Path:   "/foo/2/bar",
		Headers: map[string]string{
			RequestPathUserIDArgName: "1",
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
		permissions: []string{"somePermission"},
	}

	for n := 0; n < b.N; n++ {
		// fmt.Println("running", n)
		HasPerm(req, tu, "somePermission")
	}
}

// 32.6 ns/op  0 B/op 0 allocs/op
func BenchmarkHasPermAsOwner(b *testing.B) {
	req := &request.Request{
		Method:  "GET",
		Path:    "/foo/2/bar",
		Headers: map[string]string{},
		Params: map[string]string{
			RequestPathUserIDArgName: "2",
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
		permissions: []string{"somePermissionAsOwner"},
	}

	for n := 0; n < b.N; n++ {
		HasPerm(req, tu, "somePermissionAsOwner")
	}
}

// 31.3 ns/op 0 B/op 0 allocs/op
// go test -benchmem -run=^$ -bench ^(BenchmarkHasPermAsNotOwner)$ github.com/macinnir/dvc/core/lib/utils -v
func BenchmarkHasPermAsNotOwner(b *testing.B) {
	req := &request.Request{
		Method:  "GET",
		Path:    "/foo/3/bar",
		Headers: map[string]string{},
		Params: map[string]string{
			RequestPathUserIDArgName: "3",
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
		permissions: []string{"somePermissionAsOwner"},
	}

	for n := 0; n < b.N; n++ {
		HasPerm(req, tu, "somePermissionAsOwner")
	}

}
