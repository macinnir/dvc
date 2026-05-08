package query_test

import (
	"testing"

	"github.com/macinnir/dvc/core/lib/utils/query"
	"github.com/macinnir/dvc/core/lib/utils/query/testassets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParen(t *testing.T) {
	expected := "SELECT `t`.* FROM `FiscalYear` `t` WHERE `t`.`IsDeleted` = 0 AND ( ( `t`.`DateFrom` BETWEEN 1 AND 2 ) OR ( `t`.`DateTo` BETWEEN 3 AND 4 ) )"
	actual, e := query.Select(&testassets.FiscalYear{}).
		Where(
			query.EQ("IsDeleted", 0),
			query.And(
				query.Paren(query.Between("DateFrom", 1, 2)),
				query.Or(),
				query.Paren(query.Between("DateTo", 3, 4)),
			),
		).
		String()

	assert.Nil(t, e)
	assert.Equal(t, expected, actual)
}

func TestAnds(t *testing.T) {
	q1, e1 := query.Select(&testassets.FiscalYear{}).
		Where(
			query.Ands(
				query.EQ(testassets.FiscalYear_Column_IsDeleted, 0),
				query.Between(testassets.FiscalYear_Column_DateFrom, 1, 2),
				query.EQ(testassets.FiscalYear_Column_IsLocked, 0),
				query.EQ(testassets.FiscalYear_Column_Year, 2021),
			),
		).String()

	q2, e2 := query.Select(&testassets.FiscalYear{}).
		Where(
			query.Ands(
				query.EQ(testassets.FiscalYear_Column_IsDeleted, 0),
				query.Between(testassets.FiscalYear_Column_DateFrom, 1, 2),
				query.EQ(testassets.FiscalYear_Column_IsLocked, 0),
				query.EQ(testassets.FiscalYear_Column_Year, 2021),
				nil,
			),
		).String()

	q3, e3 := query.Select(&testassets.FiscalYear{}).
		Where(
			query.Ands(
				query.EQ(testassets.FiscalYear_Column_IsDeleted, 0),
				query.Between(testassets.FiscalYear_Column_DateFrom, 1, 2),
				query.EQ(testassets.FiscalYear_Column_IsLocked, 0),
				nil,
				query.EQ(testassets.FiscalYear_Column_Year, 2021),
			),
		).String()
	q4, e4 := query.Select(&testassets.FiscalYear{}).
		Where(
			query.Ands(
				nil,
				query.EQ(testassets.FiscalYear_Column_IsDeleted, 0),
				query.Between(testassets.FiscalYear_Column_DateFrom, 1, 2),
				query.EQ(testassets.FiscalYear_Column_IsLocked, 0),
				query.EQ(testassets.FiscalYear_Column_Year, 2021),
			),
		).String()

	assert.Nil(t, e1)
	assert.Equal(t, "SELECT `t`.* FROM `FiscalYear` `t` WHERE `t`.`IsDeleted` = 0 AND `t`.`DateFrom` BETWEEN 1 AND 2 AND `t`.`IsLocked` = 0 AND `t`.`Year` = 2021", q1)
	assert.Nil(t, e2)
	assert.Equal(t, "SELECT `t`.* FROM `FiscalYear` `t` WHERE `t`.`IsDeleted` = 0 AND `t`.`DateFrom` BETWEEN 1 AND 2 AND `t`.`IsLocked` = 0 AND `t`.`Year` = 2021", q2)
	assert.Nil(t, e3)
	assert.Equal(t, "SELECT `t`.* FROM `FiscalYear` `t` WHERE `t`.`IsDeleted` = 0 AND `t`.`DateFrom` BETWEEN 1 AND 2 AND `t`.`IsLocked` = 0 AND `t`.`Year` = 2021", q3)
	assert.Nil(t, e4)
	assert.Equal(t, "SELECT `t`.* FROM `FiscalYear` `t` WHERE `t`.`IsDeleted` = 0 AND `t`.`DateFrom` BETWEEN 1 AND 2 AND `t`.`IsLocked` = 0 AND `t`.`Year` = 2021", q4)
}

func TestOrs(t *testing.T) {
	q1, e1 := query.Select(&testassets.FiscalYear{}).
		Where(
			query.Ors(
				query.EQ(testassets.FiscalYear_Column_IsDeleted, 0),
				query.EQ(testassets.FiscalYear_Column_IsLocked, 0),
				query.EQ(testassets.FiscalYear_Column_Year, 2021),
			),
		).String()

	q2, e2 := query.Select(&testassets.FiscalYear{}).
		Where(
			query.Ors(
				query.EQ(testassets.FiscalYear_Column_IsDeleted, 0),
				query.EQ(testassets.FiscalYear_Column_IsLocked, 0),
				query.EQ(testassets.FiscalYear_Column_Year, 2021),
				nil,
			),
		).String()

	q3, e3 := query.Select(&testassets.FiscalYear{}).
		Where(
			query.Ors(
				query.EQ(testassets.FiscalYear_Column_IsDeleted, 0),
				query.EQ(testassets.FiscalYear_Column_IsLocked, 0),
				nil,
				query.EQ(testassets.FiscalYear_Column_Year, 2021),
			),
		).String()

	q4, e4 := query.Select(&testassets.FiscalYear{}).
		Where(
			query.Ors(
				nil,
				query.EQ(testassets.FiscalYear_Column_IsDeleted, 0),
				query.EQ(testassets.FiscalYear_Column_IsLocked, 0),
				query.EQ(testassets.FiscalYear_Column_Year, 2021),
			),
		).String()

	assert.Nil(t, e1)
	assert.Equal(t, "SELECT `t`.* FROM `FiscalYear` `t` WHERE `t`.`IsDeleted` = 0 OR `t`.`IsLocked` = 0 OR `t`.`Year` = 2021", q1)
	assert.Nil(t, e2)
	assert.Equal(t, "SELECT `t`.* FROM `FiscalYear` `t` WHERE `t`.`IsDeleted` = 0 OR `t`.`IsLocked` = 0 OR `t`.`Year` = 2021", q2)
	assert.Nil(t, e3)
	assert.Equal(t, "SELECT `t`.* FROM `FiscalYear` `t` WHERE `t`.`IsDeleted` = 0 OR `t`.`IsLocked` = 0 OR `t`.`Year` = 2021", q3)
	assert.Nil(t, e4)
	assert.Equal(t, "SELECT `t`.* FROM `FiscalYear` `t` WHERE `t`.`IsDeleted` = 0 OR `t`.`IsLocked` = 0 OR `t`.`Year` = 2021", q4)
}

func TestWhereTypeAll(t *testing.T) {
	q, e := query.Select(&testassets.Job{}).Where(query.WhereAll()).String()
	require.Nil(t, e)
	expected := "SELECT `t`.* FROM `Job` `t` WHERE 1=1"
	assert.Equal(t, expected, q)

	q, e = query.Select(&testassets.Job{}).Where(query.WhereAll(), query.And(), query.EQ("IsDeleted", 0)).String()
	require.Nil(t, e)
	expected = "SELECT `t`.* FROM `Job` `t` WHERE 1=1 AND `t`.`IsDeleted` = 0"
	assert.Equal(t, expected, q)
}

func TestWhere_MultiWheres(t *testing.T) {
	q := query.Select(&testassets.Job{}).Where(query.WhereAll())
	q.Where(query.And(), query.EQ("IsDeleted", 0))
	r, e := q.String()

	expected := "SELECT `t`.* FROM `Job` `t` WHERE 1=1 AND `t`.`IsDeleted` = 0"

	assert.Nil(t, e)
	assert.Equal(t, expected, r)
}

func TestWhere_InvalidFieldName(t *testing.T) {

	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.EQ("Foo", "Bar"),
		).String()

	require.NotNil(t, e)
	assert.Equal(t, "Invalid Column Name at WHERE... in model `Comment` -- Foo", e.Error())
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Foo` = Bar", sql)

}

func TestQuerySelect_MultipleWhere(t *testing.T) {
	q := query.Select(&testassets.Comment{}).
		Where(query.EQ("Content", "foo"))
	q.Where(query.And())
	q.Where(query.EQ("Name", "bar"))
	sql, e := q.String()

	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Content` = 'foo' AND `t`.`Name` = 'bar'", sql)
}

func TestQuerySelect_EmptyWhereClause(t *testing.T) {

	q := query.Select(&testassets.Comment{})
	// TODO extra where clause
	wheres := []*query.WherePart{}
	sql, e := q.Where(wheres...).String()
	require.Nil(t, e)
	// assert.Equal(t, "Empty where clause at WHERE in model `Comment` -- ", e.Error())
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t`", sql)

	q = query.Select(&testassets.Comment{})
	sql, e = q.Where(query.EQ("CommentID", 1)).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`CommentID` = 1", sql)
}
