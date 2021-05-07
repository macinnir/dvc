package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSQLStatements(t *testing.T) {

	str := ";foo;bar;"

	stmts := parseSQLStatements(str)

	assert.Len(t, stmts, 2)
	assert.Equal(t, "foo", stmts[0])
	assert.Equal(t, "bar", stmts[1])

}
