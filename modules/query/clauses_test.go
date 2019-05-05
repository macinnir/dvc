package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEquals(t *testing.T) {
	b := Equals{"id": 1}
	sql, args := b.ToSQL()

	expectedSQL := "`id` = ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestEqualsIsNull(t *testing.T) {
	b := Equals{"id": nil}
	sql, args := b.ToSQL()

	expectedSQL := "`id` IS NULL"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{}
	assert.Equal(t, expectedArgs, args)
}

func TestNotEquals(t *testing.T) {
	b := NotEquals{"id": 1}
	sql, args := b.ToSQL()

	expectedSQL := "`id` != ?"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestNotEqualsIsNull(t *testing.T) {
	b := NotEquals{"id": nil}
	sql, args := b.ToSQL()

	expectedSQL := "`id` IS NOT NULL"
	assert.Equal(t, expectedSQL, sql)

	expectedArgs := []interface{}{}
	assert.Equal(t, expectedArgs, args)
}

func TestLessThan(t *testing.T) {
	b := LessThan{"id": 1}
	sql, args := b.ToSQL()

	expectedSQL := "`id` < ?"
	assert.Equal(t, expectedSQL, sql)
	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestGreaterThan(t *testing.T) {
	b := GreaterThan{"id": 1}
	sql, args := b.ToSQL()

	expectedSQL := "`id` > ?"
	assert.Equal(t, expectedSQL, sql)
	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestLessThanOrEqual(t *testing.T) {
	b := LessThanOrEqual{"id": 1}
	sql, args := b.ToSQL()

	expectedSQL := "`id` <= ?"
	assert.Equal(t, expectedSQL, sql)
	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestGreaterThanOrEqual(t *testing.T) {
	b := GreaterThanOrEqual{"id": 1}
	sql, args := b.ToSQL()

	expectedSQL := "`id` >= ?"
	assert.Equal(t, expectedSQL, sql)
	expectedArgs := []interface{}{1}
	assert.Equal(t, expectedArgs, args)
}

func TestIn(t *testing.T) {
	b := In{"id": {1, 2, 3}}
	expectedSQL := "`id` IN (?,?,?)"
	sql, args := b.ToSQL()
	assert.Equal(t, expectedSQL, sql)
	assert.Equal(t, 3, len(args))
	assert.Equal(t, 1, args[0])
}

func TestNotIn(t *testing.T) {
	b := NotIn{"id": {1, 2, 3}}
	expectedSQL := "`id` NOT IN (?,?,?)"
	sql, args := b.ToSQL()
	assert.Equal(t, expectedSQL, sql)
	assert.Equal(t, 3, len(args))
	assert.Equal(t, 1, args[0])
}

func TestOr(t *testing.T) {
	b := Or{}
	sql, args := b.ToSQL()
	assert.Equal(t, "OR", sql)
	assert.Equal(t, 0, len(args))
}

func TestAnd(t *testing.T) {
	b := And{}
	sql, args := b.ToSQL()
	assert.Equal(t, "AND", sql)
	assert.Equal(t, 0, len(args))
}

func TestOpenParenthesis(t *testing.T) {
	b := OpenParenthesis{}
	sql, args := b.ToSQL()
	assert.Equal(t, "(", sql)
	assert.Equal(t, 0, len(args))
}

func TestCloseParenthesis(t *testing.T) {
	b := CloseParenthesis{}
	sql, args := b.ToSQL()
	assert.Equal(t, ")", sql)
	assert.Equal(t, 0, len(args))
}

func TestDistinct(t *testing.T) {
	b := Distinct{}
	sql, args := b.ToSQL()
	assert.Equal(t, "DISTINCT", sql)
	assert.Equal(t, 0, len(args))
}

func TestOrderBy_Single(t *testing.T) {
	b := OrderBy{Field: "Foo", Dir: ASC}
	sql, _ := b.ToSQL()
	assert.Equal(t, "`Foo` ASC", sql)
}

func TestSet(t *testing.T) {
	b := Set{"Foo": "bar", "Baz": 123}
	sql, args := b.ToSQL()
	assert.Equal(t, "`Baz` = ?, `Foo` = ?", sql)
	assert.Equal(t, 2, len(args))
	assert.Equal(t, 123, args[0])
	assert.Equal(t, "bar", args[1])
}

func TestValues(t *testing.T) {
	b := Values{"A": 1, "B": 2}
	sql, args := b.ToSQL()
	assert.Equal(t, "(`A`, `B`) VALUES (?, ?)", sql)
	assert.Equal(t, 2, len(args))
	assert.Equal(t, 1, args[0])
	assert.Equal(t, 2, args[1])
}
