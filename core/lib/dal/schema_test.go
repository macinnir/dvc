package dal

import (
	"database/sql"
	"testing"
)

func TestSchema(t *testing.T) {

	conn := new(sql.DB)
	d := NewDal(conn)

	s := d.Schema

	s.AddTable("foo", []string{})
	if s.Table("foo").Alias != "f" {
		t.Errorf("Table alias should have been '%s'", "f")
	}
	s.Table("foo").AddField("id")
	s.Table("foo").AddField("name")
	if s.Table("foo").Field("id").Name != "id" {
		t.Errorf("Table foe should have a column name id but doesn't")
	}

	s.AddTable("foe", []string{
		"id",
		"name",
		"createdAt",
		"updatedAt",
		"fooId",
		"foo2Id",
	})
	if s.Table("foe").Alias != "f1" {
		t.Errorf("Table alias should have been '%s'", "f1")
	}

	s.AddTable("foo2", []string{
		"id",
		"name",
	})

	var q IQuery
	var selectString string

	q = s.Select("foe")
	q.Join("foo")
	q.OnField("fooId", "foe", "id")
	q.Join("foo2").OnField("id", "foe", "foo2Id")
	q.Where("name", "foo").Where("!fooId", 123).Or().Where("foodId", 111).ToSQL()

	selectString = q.ToSQL()

	expectedSelectString := "SELECT `f1`.`id`, `f1`.`name`, `f1`.`createdAt`, `f1`.`updatedAt`, `f1`.`fooId`, `f1`.`foo2Id` FROM `foe` `f1` JOIN `foo` `f` ON `f`.`fooId` = `f1`.`id` JOIN `foo2` `f2` ON `f2`.`id` = `f1`.`foo2Id` WHERE `f1`.`name` = ? AND `f1`.`fooId` != ? OR `f1`.`foodId` = ?"
	if selectString != expectedSelectString {
		t.Errorf("Actual Query: \n\n %s\n\nExpeted:\n\n %s\n\n", selectString, expectedSelectString)
	}

	expectedSelectStringWithJoinField := "SELECT `f1`.`id`, `f1`.`name`, `f1`.`createdAt`, `f1`.`updatedAt`, `f1`.`fooId`, `f1`.`foo2Id`, `f2`.`id` as `f2Id` FROM `foe` `f1` JOIN `foo` `f` ON `f`.`fooId` = `f1`.`id` JOIN `foo2` `f2` ON `f2`.`id` = `f1`.`foo2Id` WHERE `f1`.`name` = ? AND `f1`.`fooId` != ? OR `f1`.`foodId` = ?"

	q.SelectJoinField("foo2", "id", "fooId")

	if selectString != expectedSelectString {
		t.Errorf("Actual Query: \n\n %s\n\nExpeted:\n\n %s\n\n", selectString, expectedSelectStringWithJoinField)
	}
	expectedUpdateString := "UPDATE `foe` `f1` SET `f1`.`name` = ?, `f1`.`fooId` = ? WHERE `f1`.`id` = ?"
	q = s.Update("foe").Set("name", "Another name").Set("fooId", 123).Where("id", 123)
	updateString := q.ToSQL()

	if q.GetValues()[0] != "Another name" {
		t.Errorf("Incorrect values. Expected: %s Actual: %s", "Another name", q.GetValues()[0])
	}

	if updateString != expectedUpdateString {
		t.Errorf("Query %s not correct", updateString)
	}

	expectedDeleteString := "DELETE FROM `f1` USING `foe` AS `f1` WHERE `f1`.`id` = ?"
	deleteString := s.Delete("foe").Where("id", 1).ToSQL()
	if deleteString != expectedDeleteString {
		t.Errorf("Query %s not correct", deleteString)
	}

	expectedInsertString := "INSERT INTO `foe` (`name`,`fooId`) VALUES (?,?)"
	insertString := s.Insert("foe").Set("name", "foo").Set("fooId", 1).ToSQL()
	if insertString != expectedInsertString {
		t.Errorf("Query %s not correct", insertString)
	}

	countString := s.Count("foo").Where("fooId", 123).ToSQL()
	expectedCountString := "SELECT COUNT(DISTINCT `f`.`id`) FROM `foo` `f` WHERE `f`.`fooId` = ?"
	if countString != expectedCountString {
		t.Errorf("Query %s not correct", countString)
	}

	// q.Set
	// q.Update()
}
