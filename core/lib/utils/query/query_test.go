package query

import (
	"testing"

	"github.com/macinnir/dvc/core/lib/utils/query/testassets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuerySelect(t *testing.T) {
	q := Select(&testassets.Comment{})
	var e error

	sql, e := q.String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t`", sql)

	sql, e = Select(&testassets.Comment{}).
		Where(
			GT("DateCreated", 2),
			Or(),
			EQ("Content", "foo"),
			Or(),
			EQ("Name", "bar"),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`DateCreated` > 2 OR `t`.`Content` = 'foo' OR `t`.`Name` = 'bar'", sql)

	sql, e = Select(&testassets.Comment{}).
		Where(
			GT("DateCreated", 2),
			Or(),
			EQ("Content", "foo"),
			And(
				GTOE("DateCreated", 1),
				Or(),
				LTOE("DateCreated", 2),
				Or(),
				LT("DateCreated", 3),
			),
		).
		String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`DateCreated` > 2 OR `t`.`Content` = 'foo' AND ( `t`.`DateCreated` >= 1 OR `t`.`DateCreated` <= 2 OR `t`.`DateCreated` < 3 )", sql)

	sql, e = Select(&testassets.Comment{}).
		Where(
			WhereAll(),
			And(
				GT("DateCreated", 2),
				Or(),
				EQ("Content", "foo"),
			),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE 1=1 AND ( `t`.`DateCreated` > 2 OR `t`.`Content` = 'foo' )", sql)

	sql, e = Select(&testassets.Comment{}).
		Where(
			WhereAll(),
			Or(
				GT("DateCreated", 2),
				And(),
				EQ("Content", "foo"),
			),
			And(
				Between("ObjectID", 1, 2),
			),
			And(),
			IN("Content", "foo", "bar", "baz"),
			And(),
			NE("Content", "quux"),
			And(),
			NE("ObjectID", "5"),
		).
		OrderBy("Content", "asc").
		Limit(1, 2).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE 1=1 OR ( `t`.`DateCreated` > 2 AND `t`.`Content` = 'foo' ) AND ( `t`.`ObjectID` BETWEEN 1 AND 2 ) AND `t`.`Content` IN ( 'foo', 'bar', 'baz' ) AND `t`.`Content` <> 'quux' AND `t`.`ObjectID` <> 5 ORDER BY `t`.`Content` ASC LIMIT 1 OFFSET 2", sql)
}

func TestQuerySelect_InvalidFieldName(t *testing.T) {

	sql, e := Select(&testassets.Comment{}).
		Where(
			EQ("Foo", "Bar"),
		).String()

	require.NotNil(t, e)
	assert.Equal(t, "WHERE(): INVALID COLUMN: `Comment`.`Foo`", e.Error())
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Foo` = Bar", sql)

}

func TestQuery_INString(t *testing.T) {

	args := []string{"foo", "bar", "baz"}

	sql, e := Select(&testassets.Comment{}).
		Where(
			INString(
				"Content",
				args,
			),
		).String()

	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Content` IN ( 'foo', 'bar', 'baz' )", sql)

}

func TestQuery_INInt64(t *testing.T) {

	args := []int64{1, 2, 3}

	sql, e := Select(&testassets.Comment{}).
		Where(
			INInt64(
				"CommentID",
				args,
			),
		).String()

	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`CommentID` IN ( 1, 2, 3 )", sql)

}

func TestQuerySelect_EmptyWhereClause(t *testing.T) {

	q := Select(&testassets.Comment{})
	// TODO extra where clause
	sql, e := q.Where().String()
	require.NotNil(t, e)
	assert.Equal(t, "EMPTY_WHERE_CLAUSE: `Comment`", e.Error())
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE ", sql)

	q = Select(&testassets.Comment{})
	sql, e = q.Where(EQ("CommentID", 1)).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`CommentID` = 1", sql)
}

func TestQuerySelect_InvalidField(t *testing.T) {
	sql, e := Select(&testassets.Comment{}).
		Field("Foo").
		String()

	assert.Equal(t, "SELECT `t`.`Foo` FROM `Comment` `t`", sql)
	require.NotNil(t, e)
	assert.Equal(t, "SELECT: INVALID COLUMN: `Comment`.`Foo`", e.Error())
}

func TestWhereLike(t *testing.T) {
	sql, e := Select(&testassets.Comment{}).Where(Like("Name", "Foo%")).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Name` LIKE 'Foo%'", sql)
}

func TestWhereLike_InvalidValue(t *testing.T) {
	_, e := Select(&testassets.Comment{}).Where(Like("CommentID", "Foo%")).String()
	require.NotNil(t, e)
	assert.Equal(t, "LIKE: INVALID VALUE: `Comment`.`%d` => Foo%", e.Error())
}

func TestWhereNotLike(t *testing.T) {
	sql, e := Select(&testassets.Comment{}).Where(NotLike("Name", "Foo%")).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Name` NOT LIKE 'Foo%'", sql)
}

func TestWhereNotLike_InvalidValue(t *testing.T) {
	_, e := Select(&testassets.Comment{}).Where(NotLike("CommentID", "Foo%")).String()
	require.NotNil(t, e)
	assert.Equal(t, "NOT LIKE: INVALID VALUE: `Comment`.`%d` => Foo%", e.Error())
}

func TestUnion(t *testing.T) {
	var e error
	sql, e := Union(
		Select(&testassets.Comment{}).Where(EQ("Content", "bar")),
		Select(&testassets.Comment{}).Where(EQ("Content", "baz")),
	)
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Content` = 'bar' UNION ALL SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Content` = 'baz'", sql)
}

func TestUpdate(t *testing.T) {
	sql, e := Update(&testassets.Comment{}).
		Set("Content", "bar").
		Set("ObjectID", 1).
		Where(EQ("CommentID", 123)).String()
	require.Nil(t, e)
	assert.Equal(t, "UPDATE `Comment` SET `Content` = 'bar', `ObjectID` = 1 WHERE `CommentID` = 123", sql)
}

func TestUpdate_InvalidField(t *testing.T) {
	sql, e := Update(&testassets.Comment{}).
		Set("Foo", "bar").
		Set("ObjectID", 1).
		Where(EQ("CommentID", 123)).String()
	require.NotNil(t, e)
	assert.Equal(t, "UPDATE `Comment` SET `Foo` = 'bar', `ObjectID` = 1 WHERE `CommentID` = 123", sql)
}

func TestDelete(t *testing.T) {
	sql, e := Delete(&testassets.Comment{}).
		Where(EQ("CommentID", 123)).String()
	require.Nil(t, e)
	assert.Equal(t, "DELETE FROM `Comment` WHERE `CommentID` = 123", sql)
}

func TestInsert(t *testing.T) {
	sql, e := Insert(&testassets.Comment{}).
		Set("DateCreated", 1).
		Set("Content", "foo").
		Set("ObjectType", 2).
		Set("ObjectID", 3).
		String()
	require.Nil(t, e)
	assert.Equal(t, "INSERT INTO `Comment` ( `DateCreated`, `Content`, `ObjectType`, `ObjectID` ) VALUES ( 1, 'foo', 2, 3 )", sql)
}

func TestInsert_InvalidFieldName(t *testing.T) {

	sql, e := Insert(&testassets.Comment{}).
		Set("Foo", "Bar").String()

	require.NotNil(t, e)
	assert.Equal(t, "INSERT(): INVALID COLUMN: `Comment`.`Foo`", e.Error())
	assert.Equal(t, "INSERT INTO `Comment` ( `Foo` ) VALUES ( Bar )", sql)

}

func TestSelectFields(t *testing.T) {
	sql, e := Select(&testassets.Job{}).
		Count("JobID", "ProjectsQuoted").
		Sum("TotalPrice", "SalesVolume").
		Sum("GrossProfit", "GM").
		// Field("COALESCE(SUM(TotalPrice), 0)", "SalesVolume").
		// Field("COALESCE(SUM(GrossProfit), 0)", "GM").
		Where(
			EQ("IsDeleted", 0),
			And(),
			Between("AwardDate", 1, 2),
		).
		String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT COUNT(`t`.`JobID`) AS `ProjectsQuoted`, COALESCE(SUM(`t`.`TotalPrice`), 0) AS `SalesVolume`, COALESCE(SUM(`t`.`GrossProfit`), 0) AS `GM` FROM `Job` `t` WHERE `t`.`IsDeleted` = 0 AND `t`.`AwardDate` BETWEEN 1 AND 2", sql)
}

func TestSum_InvalidField(t *testing.T) {
	_, e := Select(&testassets.Job{}).
		Sum("Foo", "Foo").String()
	require.NotNil(t, e)
	assert.Equal(t, "Sum(): INVALID COLUMN: `Job`.`Foo`", e.Error())
	// assert.Equal(t, "SELECT SUM(`t`.`Foo`) AS `Foo` FROM `Job` `t`", sql)
}

func TestSelectFields2(t *testing.T) {
	sql, e := Select(&testassets.Job{}).
		Field("JobID").
		FieldAs("JobID", "foo").
		Where(
			EQ("IsDeleted", 0),
			And(),
			Between("AwardDate", 1, 2),
		).
		String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.`JobID`, `t`.`JobID` AS `foo` FROM `Job` `t` WHERE `t`.`IsDeleted` = 0 AND `t`.`AwardDate` BETWEEN 1 AND 2", sql)
}

func TestSelectFields3(t *testing.T) {
	sql, e := Select(&testassets.Job{}).
		Fields(
			"`t`.`JobID`",
			"`t`.`JobID` AS `foo`",
		).
		Where(
			EQ("IsDeleted", 0),
			And(),
			Between("AwardDate", 1, 2),
		).
		String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.`JobID`, `t`.`JobID` AS `foo` FROM `Job` `t` WHERE `t`.`IsDeleted` = 0 AND `t`.`AwardDate` BETWEEN 1 AND 2", sql)
}

func TestSelectAlias(t *testing.T) {
	sql, e := Select(&testassets.Job{}).
		Alias("j").
		Count("JobID", "ProjectsQuoted").
		// Field("COALESCE(SUM(TotalPrice), 0)", "SalesVolume").
		// Field("COALESCE(SUM(GrossProfit), 0)", "GM").
		Where(
			EQ("IsDeleted", 0),
			And(),
			Between("AwardDate", 1, 2),
		).
		String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT COUNT(`j`.`JobID`) AS `ProjectsQuoted` FROM `Job` `j` WHERE `j`.`IsDeleted` = 0 AND `j`.`AwardDate` BETWEEN 1 AND 2", sql)
}

func TestSelectExists(t *testing.T) {

	actual, e := Select(&testassets.Job{}).
		Alias("j").
		Count("JobID", "ProjectsQuoted").
		Sum("TotalPrice", "SalesVolume").
		Where(
			Exists(
				Select(&testassets.JobSales{}).
					Alias("js").
					FieldRaw("1", "n").
					Where(
						EQF("JobID", "`j`.`JobID`"),
						And(),
						EQ("IsDeleted", 0),
						And(),
						EQ("UserID", 1),
					),
			),

			//"SELECT 1 FROM `JobSales` `js` WHERE `js`.`JobID` = `j`.`JobID` AND `js`.`IsDeleted` = 0 AND `js`.`UserID` = 1"
			And(),
			EQ("IsDeleted", 0),
			And(),
			Between("AwardDate", 1, 2),
		).String()

	require.Nil(t, e)
	expected := "SELECT COUNT(`j`.`JobID`) AS `ProjectsQuoted`, COALESCE(SUM(`j`.`TotalPrice`), 0) AS `SalesVolume` FROM `Job` `j` WHERE"
	expected += " EXISTS ( SELECT 1 AS `n` FROM `JobSales` `js` WHERE `js`.`JobID` = `j`.`JobID` AND `js`.`IsDeleted` = 0 AND `js`.`UserID` = 1 )"
	expected += " AND `j`.`IsDeleted` = 0 AND `j`.`AwardDate` BETWEEN 1 AND 2"

	assert.Equal(t, expected, actual)
}

func TestWhereTypeAll(t *testing.T) {
	q, e := Select(&testassets.Job{}).Where(WhereAll()).String()
	require.Nil(t, e)
	expected := "SELECT `t`.* FROM `Job` `t` WHERE 1=1"
	assert.Equal(t, expected, q)

	q, e = Select(&testassets.Job{}).Where(WhereAll(), And(), EQ("IsDeleted", 0)).String()
	require.Nil(t, e)
	expected = "SELECT `t`.* FROM `Job` `t` WHERE 1=1 AND `t`.`IsDeleted` = 0"
	assert.Equal(t, expected, q)
}

func TestWhere_MultiWheres(t *testing.T) {
	q := Select(&testassets.Job{}).Where(WhereAll())
	q.Where(And(), EQ("IsDeleted", 0))
	r, e := q.String()

	expected := "SELECT `t`.* FROM `Job` `t` WHERE 1=1 AND `t`.`IsDeleted` = 0"

	assert.Nil(t, e)
	assert.Equal(t, expected, r)
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
