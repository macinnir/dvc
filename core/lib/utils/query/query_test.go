package query_test

import (
	"testing"

	"github.com/macinnir/dvc/core/lib/utils/query"
	"github.com/macinnir/dvc/core/lib/utils/query/testassets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuerySelect(t *testing.T) {
	q := query.Select(&testassets.Comment{})
	var e error

	sql, e := q.String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t`", sql)

	sql, e = query.Select(&testassets.Comment{}).
		Where(
			query.GT("DateCreated", 2),
			query.Or(),
			query.EQ("Content", "foo's"),
			query.Or(),
			query.EQ("Name", "bar"),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`DateCreated` > 2 OR `t`.`Content` = '"+`foo\'`+"s' OR `t`.`Name` = 'bar'", sql)

	sql, e = query.Select(&testassets.Comment{}).
		Where(
			query.GT("DateCreated", 2),
			query.Or(),
			query.EQ("Content", "foo"),
			query.And(
				query.GTOE("DateCreated", 1),
				query.Or(),
				query.LTOE("DateCreated", 2),
				query.Or(),
				query.LT("DateCreated", 3),
			),
		).
		String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`DateCreated` > 2 OR `t`.`Content` = 'foo' AND ( `t`.`DateCreated` >= 1 OR `t`.`DateCreated` <= 2 OR `t`.`DateCreated` < 3 )", sql)

	sql, e = query.Select(&testassets.Comment{}).
		Where(
			query.WhereAll(),
			query.And(
				query.GT("DateCreated", 2),
				query.Or(),
				query.EQ("Content", "foo"),
			),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE 1=1 AND ( `t`.`DateCreated` > 2 OR `t`.`Content` = 'foo' )", sql)

	sql, e = query.Select(&testassets.Comment{}).
		Where(
			query.WhereAll(),
			query.Or(
				query.GT("DateCreated", 2),
				query.And(),
				query.EQ("Content", "foo"),
			),
			query.And(
				query.Between("ObjectID", 1, 2),
			),
			query.And(),
			query.IN("Content", "foo", "bar", "baz"),
			query.And(),
			query.NE("Content", "quux"),
			query.And(),
			query.NE("ObjectID", "5"),
		).
		OrderBy("Content", query.OrderDirASC).
		Limit(1, 2).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE 1=1 OR ( `t`.`DateCreated` > 2 AND `t`.`Content` = 'foo' ) AND ( `t`.`ObjectID` BETWEEN 1 AND 2 ) AND `t`.`Content` IN ( 'foo', 'bar', 'baz' ) AND `t`.`Content` != 'quux' AND `t`.`ObjectID` != 5 ORDER BY `t`.`Content` ASC LIMIT 1 OFFSET 2", sql)
}

func TestQuerySelect_LimitPage(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).LimitPage(10, 5).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` LIMIT 10 OFFSET 50", sql, "LimitPage() should have an offset that multiplies the limit by the page")
}

func TestMultipleOrderBy(t *testing.T) {
	q, e := query.Select(&testassets.Comment{}).OrderBy("CommentID", query.OrderDirASC).OrderBy("DateCreated", query.OrderDirDESC).String()
	assert.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` ORDER BY `t`.`CommentID` ASC, `t`.`DateCreated` DESC", q)
}

func TestQuerySelect_InvalidOrderByColumn(t *testing.T) {

	q, e := query.Select(&testassets.Comment{}).OrderBy("CommentID", query.OrderDirASC).String()
	assert.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` ORDER BY `t`.`CommentID` ASC", q)

	q, e = query.Select(&testassets.Comment{}).OrderBy("foo", query.OrderDirASC).String()
	assert.Equal(t, "Invalid Column Name at ORDER BY in model `Comment` -- foo", e.Error())
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` ORDER BY `t`.`foo` ASC", q)
}

func TestUnion(t *testing.T) {
	var e error
	sql, e := query.Union(
		query.Select(&testassets.Comment{}).Where(query.EQ("Content", "bar")),
		query.Select(&testassets.Comment{}).Where(query.EQ("Content", "baz")),
	)
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Content` = 'bar' UNION ALL SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Content` = 'baz'", sql)
}

func TestUpdate(t *testing.T) {
	sql, e := query.Update(&testassets.Comment{}).
		Set("Content", "bar").
		Set("ObjectID", 1).
		Where(query.EQ("CommentID", 123)).String()
	require.Nil(t, e)
	assert.Equal(t, "UPDATE `Comment` SET `Content` = 'bar', `ObjectID` = 1 WHERE `CommentID` = 123", sql)
}

func TestUpdate_InvalidField(t *testing.T) {
	sql, e := query.Update(&testassets.Comment{}).
		Set("Foo", "bar").
		Set("ObjectID", 1).
		Where(query.EQ("CommentID", 123)).String()
	require.NotNil(t, e)
	assert.Equal(t, "UPDATE `Comment` SET `Foo` = 'bar', `ObjectID` = 1 WHERE `CommentID` = 123", sql)
}

func TestDelete(t *testing.T) {
	sql, e := query.Delete(&testassets.Comment{}).
		Where(query.EQ("CommentID", 123)).String()
	require.Nil(t, e)
	assert.Equal(t, "DELETE FROM `Comment` WHERE `CommentID` = 123", sql)
}

func TestAlias(t *testing.T) {
	sql, e := query.Select(&testassets.Job{}).
		Alias("j").
		Count("JobID", "ProjectsQuoted").
		// Field("COALESCE(SUM(TotalPrice), 0)", "SalesVolume").
		// Field("COALESCE(SUM(GrossProfit), 0)", "GM").
		Where(
			query.EQ("IsDeleted", 0),
			query.And(),
			query.Between("AwardDate", 1, 2),
		).
		String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT COUNT(`j`.`JobID`) AS `ProjectsQuoted` FROM `Job` `j` WHERE `j`.`IsDeleted` = 0 AND `j`.`AwardDate` BETWEEN 1 AND 2", sql)
}

// func TestQuerySave_Insert(t *testing.T) {
// 	sql, e := (&testassets.Comment{
// 		Content: null.StringFrom("here is some test content"),
// 	}).Save()
// 	assert.Nil(t, e)
// 	assert.Equal(t, "INSERT INTO `Comment` ( `DateCreated`, `Content`, `ObjectType`, `ObjectID` ) VALUES ( 0, 'here is some test content', 0, 0 )", sql)
// }

// func TestQuerySave_Update(t *testing.T) {
// 	sql, e := Save(&testassets.Comment{
// 		CommentID:  123,
// 		ObjectType: 1,
// 		ObjectID:   2,
// 		Content:    null.StringFrom("here is some test content"),
// 	}).String()
// 	assert.Nil(t, e)
// 	assert.Equal(t, "UPDATE `Comment` SET `Content` = 'here is some test content', `ObjectType` = 1, `ObjectID` = 2 WHERE `CommentID` = 123", sql)
// }

// func TestWhereSelect(t *testing.T) {
// 	expected := "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`CommentID` < (SELECT QuoteNumberFullInt FROM QuoteNumber WHERE QuoteNumberID = %d) AND `t`.`IsDeleted` = 0 AND `t`.`JobID` > 0 ORDER BY `t`.`QuoteNumberFullInt` DESC LIMIT 1"

// 	q, e := query.Select(&testassets.Comment{}).Where(
// 		query.LT(testassets.Comment_Column_CommentID, )
// 	).String()

// 	assert.Nil(t, e)
// 	assert.Equal(t, expected, q)
// }

func TestAvgAndCountsAndSums(t *testing.T) {

	q, e := query.Select(&testassets.FiscalYear{}).
		Count(testassets.FiscalYear_Column_FiscalYearID, "Overall").
		Sum(testassets.FiscalYear_Column_FiscalYearID, "ChangesFound").
		Avg(testassets.FiscalYear_Column_FiscalYearID, "AverageChangesFound").String()

	assert.Nil(t, e)
	assert.Equal(t, "SELECT COUNT(`t`.`FiscalYearID`) AS `Overall`, COALESCE(SUM(`t`.`FiscalYearID`), 0) AS `ChangesFound`, COALESCE(AVG(`t`.`FiscalYearID`), 0) AS `AverageChangesFound` FROM `FiscalYear` `t`", q)

}

func TestOrderDirFromString(t *testing.T) {
	assert.Equal(t, query.OrderDirASC, query.OrderDirFromString("asc"))
	assert.Equal(t, query.OrderDirDESC, query.OrderDirFromString("desc"))
	assert.Equal(t, query.OrderDirASC, query.OrderDirFromString("foo"))
}
