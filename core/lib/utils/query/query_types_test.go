package query_test

import (
	"testing"

	"github.com/macinnir/dvc/core/lib/utils/query"
	"github.com/macinnir/dvc/core/lib/utils/query/testassets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsert(t *testing.T) {
	sql, e := query.Insert(&testassets.Comment{}).
		Set("DateCreated", 1).
		Set("Content", "foo").
		Set("ObjectType", 2).
		Set("ObjectID", 3).
		String()
	require.Nil(t, e)
	assert.Equal(t, "INSERT INTO `Comment` ( `DateCreated`, `Content`, `ObjectType`, `ObjectID` ) VALUES ( 1, 'foo', 2, 3 )", sql)
}

func TestInsert_InvalidFieldName(t *testing.T) {

	sql, e := query.Insert(&testassets.Comment{}).
		Set("Foo", "Bar").String()

	require.NotNil(t, e)
	assert.Equal(t, "Invalid Column Name at INSERT...SET in model `Comment` -- Foo", e.Error())
	assert.Equal(t, "INSERT INTO `Comment` ( `Foo` ) VALUES ( Bar )", sql)

}

func TestRaw(t *testing.T) {

	q, e := query.Raw(&testassets.TaskBatchSchedule{}, "SELECT * FROM `TaskBatchSchedule` WHERE 1=1 ORDER BY `t`.`LastRunDate` DESC LIMIT 1 OFFSET 2").
		String()

	assert.Nil(t, e)
	assert.Equal(t, "SELECT * FROM `TaskBatchSchedule` WHERE 1=1 ORDER BY `t`.`LastRunDate` DESC LIMIT 1 OFFSET 2", q)
}

func TestSelectExists(t *testing.T) {

	actual, e := query.Select(&testassets.Job{}).
		Count("JobID", "ProjectsQuoted").
		Sum("TotalPrice", "SalesVolume").
		Where(
			query.Exists(
				query.Select(&testassets.JobSales{}).
					Alias("js").
					FieldRaw("1", "n").
					Where(
						query.EQF("JobID", "`t`.`JobID`"),
						query.And(),
						query.EQ("IsDeleted", 0),
						query.And(),
						query.EQ("UserID", 1),
					),
			),

			//"SELECT 1 FROM `JobSales` `js` WHERE `js`.`JobID` = `j`.`JobID` AND `js`.`IsDeleted` = 0 AND `js`.`UserID` = 1"
			query.And(),
			query.EQ("IsDeleted", 0),
			query.And(),
			query.Between("AwardDate", 1, 2),
		).String()

	require.Nil(t, e)
	expected := "SELECT COUNT(`t`.`JobID`) AS `ProjectsQuoted`, COALESCE(SUM(`t`.`TotalPrice`), 0) AS `SalesVolume` FROM `Job` `t` WHERE"
	expected += " EXISTS ( SELECT 1 AS `n` FROM `JobSales` `js` WHERE `js`.`JobID` = `t`.`JobID` AND `js`.`IsDeleted` = 0 AND `js`.`UserID` = 1 )"
	expected += " AND `t`.`IsDeleted` = 0 AND `t`.`AwardDate` BETWEEN 1 AND 2"

	assert.Equal(t, expected, actual)
}

func TestSelectNotExists(t *testing.T) {

	actual, e := query.Select(&testassets.Job{}).
		Count("JobID", "ProjectsQuoted").
		Sum("TotalPrice", "SalesVolume").
		Where(
			query.NotExists(
				query.Select(&testassets.JobSales{}).
					Alias("js").
					FieldRaw("1", "n").
					Where(
						query.EQF("JobID", "`t`.`JobID`"),
						query.And(),
						query.EQ("IsDeleted", 0),
						query.And(),
						query.EQ("UserID", 1),
					),
			),

			//"SELECT 1 FROM `JobSales` `js` WHERE `js`.`JobID` = `j`.`JobID` AND `js`.`IsDeleted` = 0 AND `js`.`UserID` = 1"
			query.And(),
			query.EQ("IsDeleted", 0),
			query.And(),
			query.Between("AwardDate", 1, 2),
		).String()

	require.Nil(t, e)
	expected := "SELECT COUNT(`t`.`JobID`) AS `ProjectsQuoted`, COALESCE(SUM(`t`.`TotalPrice`), 0) AS `SalesVolume` FROM `Job` `t` WHERE"
	expected += " NOT EXISTS ( SELECT 1 AS `n` FROM `JobSales` `js` WHERE `js`.`JobID` = `t`.`JobID` AND `js`.`IsDeleted` = 0 AND `js`.`UserID` = 1 )"
	expected += " AND `t`.`IsDeleted` = 0 AND `t`.`AwardDate` BETWEEN 1 AND 2"

	assert.Equal(t, expected, actual)
}
