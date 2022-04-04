package fetcher

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractControllerNameFromFileName(t *testing.T) {
	var tests = []struct {
		path string
		name string
	}{
		{"AController.go", "A"},
		{"AController_test.go", ""},
		{"ALongControllerName", ""},
		{"foo/bar/BazController.go", "Baz"},
	}

	for k := range tests {
		assert.Equal(t, tests[k].name, extractControllerNameFromFileName(tests[k].path))
	}
}

func TestRegepx(t *testing.T) {
	r := regexp.MustCompile("-?[0-9]+")

	tests := []string{
		"-1",
		"12345",
	}

	for k := range tests {
		if !r.Match([]byte(tests[k])) {
			t.FailNow()
		}
	}
}
