package query_test

import (
	"fmt"
	"testing"

	"github.com/macinnir/dvc/core/lib/utils/query"
	"github.com/macinnir/dvc/core/lib/utils/query/testassets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// #region EQ
func TestEQ(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.EQ("Content", "foo's"),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Content` = 'foo\\'s'", sql)
}

func TestEQInt(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.EQInt("CommentID", 123),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`CommentID` = 123", sql)
}

func TestEQInt64(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.EQInt64("CommentID", 1234567890123),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`CommentID` = 1234567890123", sql)
}

func TestEQFloat(t *testing.T) {
	sql, e := query.Select(&testassets.JobSales{}).
		Where(
			query.EQFloat("CommissionPercent", 4.5),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `JobSales` `t` WHERE `t`.`CommissionPercent` = 4.5", sql)
}

func TestEQString(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.EQString("Content", "foo's"),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Content` = 'foo\\'s'", sql)
}

func TestEQInvalidField(t *testing.T) {
	_, e := query.Select(&testassets.Comment{}).
		Where(
			query.EQ("NonExistentField", "foo"),
		).String()
	require.NotNil(t, e)
	assert.Equal(t, "Invalid Column Name at WHERE... in model `Comment` -- NonExistentField", e.Error())
}

// #endregion EQ

// #region NEQ
func TestNEQ(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.NE("Content", "foo's"),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Content` != 'foo\\'s'", sql)
}

func TestNEQInt(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.NEInt("CommentID", 123),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`CommentID` != 123", sql)
}

func TestNEQInt64(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.NEInt64("CommentID", 1234567890123),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`CommentID` != 1234567890123", sql)
}

func TestNEQFloat(t *testing.T) {
	sql, e := query.Select(&testassets.JobSales{}).
		Where(
			query.NEFloat("CommissionPercent", 4.5),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `JobSales` `t` WHERE `t`.`CommissionPercent` != 4.5", sql)
}

func TestNEQString(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.NEString("Content", "foo's"),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Content` != 'foo\\'s'", sql)
}

func TestNEQInvalidField(t *testing.T) {
	_, e := query.Select(&testassets.Comment{}).
		Where(
			query.NE("NonExistentField", "foo"),
		).String()
	require.NotNil(t, e)
	assert.Equal(t, "Invalid Column Name at WHERE... in model `Comment` -- NonExistentField", e.Error())
}

// #endregion NEQ

// #region Less Than
func TestLT(t *testing.T) {
	sql, e := query.Select(&testassets.JobSales{}).
		Where(
			query.LT("CommissionPercent", 4.5),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `JobSales` `t` WHERE `t`.`CommissionPercent` < 4.5", sql)
}

func TestLTInt(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.LTInt("CommentID", 123),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`CommentID` < 123", sql)
}

func TestLTInt64(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.LTInt64("CommentID", 1234567890123),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`CommentID` < 1234567890123", sql)
}

func TestLTFloat(t *testing.T) {
	sql, e := query.Select(&testassets.JobSales{}).
		Where(
			query.LTFloat("CommissionPercent", 4.5),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `JobSales` `t` WHERE `t`.`CommissionPercent` < 4.5", sql)
}

func TestLTInvalidField(t *testing.T) {
	_, e := query.Select(&testassets.Comment{}).
		Where(
			query.LT("NonExistentField", 4.5),
		).String()
	require.NotNil(t, e)
	assert.Equal(t, "Invalid Column Name at WHERE... in model `Comment` -- NonExistentField", e.Error())
}

// #endregion Less Than

// #region Equal Field
func TestEQField(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.EQF("Content", "`t`.`Name`"),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Content` = `t`.`Name`", sql)
}

func TestEQFieldInvalidField(t *testing.T) {
	_, e := query.Select(&testassets.Comment{}).
		Where(
			query.EQF("NonExistentField", "`t`.`Name`"),
		).String()
	require.NotNil(t, e)
	assert.Equal(t, "Invalid Column Name at WHERE... in model `Comment` -- NonExistentField", e.Error())
	// assert.Equal(t, "Invalid field name `NonExistentField` at WHERE...EQField in model `Comment`", e.Error())
}

// #endregion Equal Field

// #region Not Equal Field

func TestNEF(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.NEF("Content", "`t`.`Name`"),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Content` != `t`.`Name`", sql)
}

func TestNEFieldInvalidField(t *testing.T) {
	_, e := query.Select(&testassets.Comment{}).
		Where(
			query.NEF("NonExistentField", "`t`.`Name`"),
		).String()
	require.NotNil(t, e)
	assert.Equal(t, "Invalid Column Name at WHERE... in model `Comment` -- NonExistentField", e.Error())
}

// #endregion

// #region IN
func TestIN(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.IN("Content", "foo", "bar", "baz"),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Content` IN ( 'foo', 'bar', 'baz' )", sql)
}

func TestINString(t *testing.T) {

	args := []string{"foo", "bar", "baz"}

	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.INString(
				"Content",
				args...,
			),
		).String()

	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Content` IN ( 'foo', 'bar', 'baz' )", sql)

}

func TestINInt64(t *testing.T) {

	args := []int64{1, 2, 3}

	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.INInt64(
				"CommentID",
				args...,
			),
		).String()

	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`CommentID` IN ( 1, 2, 3 )", sql)

}

func TestINInt(t *testing.T) {

	args := []int{1, 2, 3}

	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.INInt(
				"CommentID",
				args...,
			),
		).String()

	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`CommentID` IN ( 1, 2, 3 )", sql)

}

// #endregion
// #region NOT IN
func TestNotIN(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.NOTIN("Content", "foo", "bar", "baz"),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Content` NOT IN ( 'foo', 'bar', 'baz' )", sql)
}

func TestNotINInt(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.NOTINInt("CommentID", 1, 2, 3),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`CommentID` NOT IN ( 1, 2, 3 )", sql)
}
func TestNotINInt64(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.NOTINInt64("CommentID", 1, 2, 3),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`CommentID` NOT IN ( 1, 2, 3 )", sql)
}

func TestNotINFloat(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).
		Where(
			query.NOTINFloat("CommentID", 1.1, 2.2, 3.3),
		).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`CommentID` NOT IN ( 1.1, 2.2, 3.3 )", sql)
}

// #endregion

func TestWhereLike(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).Where(query.Like("Name", "Foo%")).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Name` LIKE 'Foo%'", sql)
}

func TestWhereLike_InvalidValue(t *testing.T) {
	_, e := query.Select(&testassets.Comment{}).Where(query.Like("CommentID", "Foo%")).String()
	require.NotNil(t, e)
	assert.Equal(t, "Invalid value at WHERE...LIKE in model `Comment` -- `%d` value: Foo%", e.Error())
}

func TestWhereNotLike(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).Where(query.NotLike("Name", "Foo%")).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Name` NOT LIKE 'Foo%'", sql)
}

func TestWhereNotLike_InvalidValue(t *testing.T) {
	_, e := query.Select(&testassets.Comment{}).Where(query.NotLike("CommentID", "Foo%")).String()
	require.NotNil(t, e)
	assert.Equal(t, "Invalid value at WHERE...NOT LIKE in model `Comment` -- `%d` value: Foo%", e.Error())
}

func TestWhereIsNull(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).Where(query.IsNull("Name")).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Name` IS NULL", sql)
}

func TestWhereIsNotNull(t *testing.T) {
	sql, e := query.Select(&testassets.Comment{}).Where(query.IsNotNull("Name")).String()
	require.Nil(t, e)
	assert.Equal(t, "SELECT `t`.* FROM `Comment` `t` WHERE `t`.`Name` IS NOT NULL", sql)
}

func TestMod(t *testing.T) {
	q, e := query.Select(&testassets.Job{}).Where(query.Mod("IsDeleted", 1, 0)).String()
	expected := "SELECT `t`.* FROM `Job` `t` WHERE MOD(`t`.`IsDeleted`, 1) = 0"

	assert.Nil(t, e)
	assert.Equal(t, expected, q)
}

func TestModF(t *testing.T) {
	q, e := query.Select(&testassets.Job{}).Where(query.Modf(1, "IsDeleted", 0)).String()
	expected := "SELECT `t`.* FROM `Job` `t` WHERE MOD(1, `t`.`IsDeleted`) = 0"

	assert.Nil(t, e)
	assert.Equal(t, expected, q)
}

func ExampleMod() {
	query.Mod("foo", 2, 1)
	fmt.Println("MOD(`t`.`foo`, 2) = 1")

	// Output: MOD(`t`.`foo`, 2) = 1
}

func TestBitAnd(t *testing.T) {
	q, e := query.Select(&testassets.Job{}).Where(query.BitAnd("IsDeleted", 1, 0)).String()
	expected := "SELECT `t`.* FROM `Job` `t` WHERE `t`.`IsDeleted` & 1 = 0"

	assert.Nil(t, e)
	assert.Equal(t, expected, q)
}

func TestModAndBitwise(t *testing.T) {
	// now := time.Now()

	dayBW := int64(16)       // int64(math.Pow(2, float64(now.Day())))
	weekdayBW := int64(1)    // int64(math.Pow(2, float64(now.Weekday())))
	hourBW := int64(2048)    // int64(math.Pow(2, float64(now.Hour())))
	minute := int64(56)      // int64(time.Now().Minute())
	seconds := 1625414209657 // time.Now().UnixNano() / 1000000
	// SELECT * FROM TaskBatchSchedule tbs WHERE tbs.IsActive = 1 AND tbs.IsDeleted = 0 AND ( tbs.DOM = 0 OR 16 & tbs.DOM = 16 ) AND ( tbs.DOW = 0 OR 1 & tbs.DOW = 1 ) AND ( tbs.HOD = 0 OR 2048 & tbs.HOD = 2048 ) AND ( MOD(56, tbs.MOH ) = 0 AND ( tbs.LastRunDate + 60000 < 1625414209657 )

	q := "SELECT `t`.* FROM `TaskBatchSchedule` `t` WHERE `t`.`IsActive` = 1 AND `t`.`IsDeleted` = 0 "
	q += fmt.Sprintf("AND ( `t`.`DOM` = 0 OR `t`.`DOM` & %d = %d ) ", dayBW, dayBW)
	q += fmt.Sprintf("AND ( `t`.`DOW` = 0 OR `t`.`DOW` & %d = %d ) ", weekdayBW, weekdayBW)
	q += fmt.Sprintf("AND ( `t`.`HOD` = 0 OR `t`.`HOD` & %d = %d ) ", hourBW, hourBW)
	// Every minute of the hour
	q += fmt.Sprintf("AND ( MOD(%d, `t`.`MOH`) = 0 ) ", minute)
	// It has been atleast 60 seconds since the last run
	q += fmt.Sprintf("AND ( `t`.`LastRunDate` + 60000 < %d )", seconds)

	actual, e := query.Select(&testassets.TaskBatchSchedule{}).
		Where(
			query.EQ("IsActive", 1),
			query.And(),
			query.EQ("IsDeleted", 0),
			query.And(query.EQ("DOM", 0), query.Or(), query.BitAnd("DOM", dayBW, dayBW)),
			query.And(query.EQ("DOW", 0), query.Or(), query.BitAnd("DOW", weekdayBW, weekdayBW)),
			query.And(query.EQ("HOD", 0), query.Or(), query.BitAnd("HOD", hourBW, hourBW)),
			query.And(query.Modf(minute, "MOH", 0)),
			query.And(query.Rawf("`t`.`LastRunDate` + 60000 < %d", seconds)),
		).
		String()

	// fmt.Println(q)

	assert.Nil(t, e)
	assert.Equal(t, q, actual)

}
