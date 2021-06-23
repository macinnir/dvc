package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractRootNameFromKey(t *testing.T) {
	tests := [][]string{
		{"foo", "foo"},
		{"foo_", "foo"},
		{"foo_123", "foo"},
		{"foo_bar_123", "foo_bar"},
		{"foo_bar_123_456", "foo_bar_123"},
	}

	for k := range tests {
		actual := ExtractRootNameFromKey(tests[k][0])
		assert.Equal(t, tests[k][1], actual)
	}
}

func TestGetDatabasesByRootName(t *testing.T) {
	config := &Config{
		Databases: []*ConfigDatabase{
			{Key: "Foo_0", Name: "Foo", Host: "1.0.0.0"},
			{Key: "Foo_1", Name: "Foo", Host: "1.0.0.1"},
			{Key: "Foo_2", Name: "Foo", Host: "1.0.0.2"},
			{Key: "Foo_3", Name: "Foo", Host: "1.0.0.3"},
			{Key: "Bar_0", Name: "Bar", Host: "1.0.0.0"},
			{Key: "Bar_1", Name: "Bar", Host: "1.0.0.1"},
			{Key: "Bar_2", Name: "Bar", Host: "1.0.0.2"},
		},
	}

	fooDatabases := GetDatabasesByRootName("Foo", config)

	assert.Len(t, fooDatabases, 4)
	assert.Equal(t, "Foo_0", fooDatabases[0].Key)
	assert.Equal(t, "Foo_1", fooDatabases[1].Key)
	assert.Equal(t, "Foo_2", fooDatabases[2].Key)
	assert.Equal(t, "Foo_3", fooDatabases[3].Key)

}
