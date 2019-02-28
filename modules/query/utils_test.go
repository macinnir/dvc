package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSortedKeys(t *testing.T) {
	exp := map[string]interface{}{"b": 1, "a": 2}
	keys := getSortedKeys(exp)
	assert.Equal(t, keys[0], "a")
	assert.Equal(t, keys[1], "b")
}

func TestGetSortedInKeys(t *testing.T) {
	exp := map[string][]interface{}{
		"b": []interface{}{},
		"a": []interface{}{},
	}
	keys := getSortedInKeys(exp)

	assert.Equal(t, keys[0], "a")
	assert.Equal(t, keys[1], "b")
}

func TestBuildWhereClause_ShouldReturnSameStatementIfSingle(t *testing.T) {
	clauses := []string{"a = ?"}
	result := buildWhereClause(clauses)
	assert.Equal(t, "a = ?", result)
}

func TestBuildWHereClause_ShouldReturnClausesSeparatedByAnd(t *testing.T) {
	clauses := []string{"a = ?", "b = ?"}
	result := buildWhereClause(clauses)
	assert.Equal(t, "a = ? AND b = ?", result)
}

func TestEscapeField(t *testing.T) {
	field := "a"
	result := escapeField(field)
	assert.Equal(t, "`a`", result)
}
