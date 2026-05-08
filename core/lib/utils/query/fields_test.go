package query_test

import (
	"testing"

	"github.com/macinnir/dvc/core/lib/utils/query"
	"github.com/macinnir/dvc/core/lib/utils/query/testassets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelectFields(t *testing.T) {
	sql, e := query.Select(&testassets.Job{}).
		Count("JobID", "ProjectsQuoted").
		Sum("TotalPrice", "SalesVolume").
		Sum("GrossProfit", "GM").
		// Field("COALESCE(SUM(TotalPrice), 0)", "SalesVolume").
		// Field("COALESCE(SUM(GrossProfit), 0)", "GM").
		Where(
			query.EQ("IsDeleted", 0),
			query.And(),
			query.Between("AwardDate", 1, 2),
		).
		String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT COUNT(`t`.`JobID`) AS `ProjectsQuoted`, COALESCE(SUM(`t`.`TotalPrice`), 0) AS `SalesVolume`, COALESCE(SUM(`t`.`GrossProfit`), 0) AS `GM` FROM `Job` `t` WHERE `t`.`IsDeleted` = 0 AND `t`.`AwardDate` BETWEEN 1 AND 2", sql)
}

func TestFields(t *testing.T) {
	sql, e := query.Select(&testassets.Job{}).
		Fields(
			query.NewField(query.FieldTypeBasic, "JobID"),
			query.NewField(query.FieldTypeBasic, "JobID", "foo"),
		).
		Where(
			query.EQ("IsDeleted", 0),
			query.And(),
			query.Between("AwardDate", 1, 2),
		).
		String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.`JobID`, `t`.`JobID` AS `foo` FROM `Job` `t` WHERE `t`.`IsDeleted` = 0 AND `t`.`AwardDate` BETWEEN 1 AND 2", sql)
}

func TestFieldAs(t *testing.T) {
	sql, e := query.Select(&testassets.Job{}).
		Field("JobID").
		FieldAs("JobID", "foo").
		Where(
			query.EQ("IsDeleted", 0),
			query.And(),
			query.Between("AwardDate", 1, 2),
		).
		String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.`JobID`, `t`.`JobID` AS `foo` FROM `Job` `t` WHERE `t`.`IsDeleted` = 0 AND `t`.`AwardDate` BETWEEN 1 AND 2", sql)
}

func TestField_InvalidField(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).
		Field("Foo").
		String()

	assert.Equal(t, "SELECT `t`.`Foo` FROM `Comment` `t`", sql)
	require.NotNil(t, e)
	assert.Equal(t, "Invalid Column Name at SELECT...Field in model `Comment` -- Foo", e.Error())
}

func TestMin(t *testing.T) {
	q1, e1 := query.Select(&testassets.FiscalYear{}).Min(testassets.FiscalYear_Column_Year, "MinYear").Where(query.EQ(testassets.FiscalYear_Column_IsDeleted, 0)).String()

	assert.Nil(t, e1)
	assert.Equal(t, "SELECT COALESCE(MIN(`t`.`Year`), 0) AS `MinYear` FROM `FiscalYear` `t` WHERE `t`.`IsDeleted` = 0", q1)
}

func TestMax(t *testing.T) {
	q1, e1 := query.Select(&testassets.FiscalYear{}).Max(testassets.FiscalYear_Column_Year, "MaxYear").Where(query.EQ(testassets.FiscalYear_Column_IsDeleted, 0)).String()

	assert.Nil(t, e1)
	assert.Equal(t, "SELECT COALESCE(MAX(`t`.`Year`), 0) AS `MaxYear` FROM `FiscalYear` `t` WHERE `t`.`IsDeleted` = 0", q1)
}

func TestDistinct(t *testing.T) {
	q1, e1 := query.Select(&testassets.FiscalYear{}).Distinct(testassets.FiscalYear_Column_Year).Where(query.EQ(testassets.FiscalYear_Column_IsDeleted, 0)).String()

	assert.Nil(t, e1)
	assert.Equal(t, "SELECT DISTINCT `t`.`Year` FROM `FiscalYear` `t` WHERE `t`.`IsDeleted` = 0", q1)
}

func TestCountAlias(t *testing.T) {
	sql, e := query.Select(&testassets.Job{}).
		Count("JobID", "ProjectsQuoted").
		Where(
			query.EQ("IsDeleted", 0),
			query.And(),
			query.Between("AwardDate", 1, 2),
		).
		Alias("j").
		String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT COUNT(`j`.`JobID`) AS `ProjectsQuoted` FROM `Job` `j` WHERE `j`.`IsDeleted` = 0 AND `j`.`AwardDate` BETWEEN 1 AND 2", sql)
}

func TestSum(t *testing.T) {
	sql, e := query.Select(&testassets.Job{}).
		Sum("TotalPrice", "SalesVolume").
		Where(
			query.EQ("IsDeleted", 0),
			query.And(),
			query.Between("AwardDate", 1, 2),
		).
		Alias("j").
		String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT COALESCE(SUM(`j`.`TotalPrice`), 0) AS `SalesVolume` FROM `Job` `j` WHERE `j`.`IsDeleted` = 0 AND `j`.`AwardDate` BETWEEN 1 AND 2", sql)
}

func TestSum_InvalidField(t *testing.T) {
	_, e := query.Select(&testassets.Job{}).
		Sum("Foo", "Foo").String()
	require.NotNil(t, e)
	assert.Equal(t, "Invalid Column Name at SELECT...Sum() in model `Job` -- Foo", e.Error())
	// assert.Equal(t, "SELECT SUM(`t`.`Foo`) AS `Foo` FROM `Job` `t`", sql)
}
