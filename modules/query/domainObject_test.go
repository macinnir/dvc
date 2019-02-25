package query

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestObject is a test domain object
type TestObject struct {
	BaseDomainObject
	ID          int64
	Name        string
	DateCreated time.Time
}

func NewTestObject() *TestObject {
	o := &TestObject{}
	o.Build()
	return o
}

func (o *TestObject) Build() {
	o.BaseDomainObject.TableName = "TestObject"
	o.BaseDomainObject.FieldList = map[string]*DomainObjectField{
		"ID": {
			Name:         "ID",
			Type:         "int64",
			IsPrimaryKey: true,
			IsNull:       false,
			Ordinal:      0,
		},
		"Name": {
			Name:    "Name",
			Type:    "string",
			IsNull:  false,
			Ordinal: 1,
		},
		"DateCreated": {
			Name:    "DateCreated",
			Type:    "time",
			IsNull:  true,
			Ordinal: 2,
		},
	}
}

func TestDomainObject_GetName(t *testing.T) {
	d := NewTestObject()
	assert.Equal(t, "TestObject", d.GetName())
}

func TestDomainObject_GetFields(t *testing.T) {
	d := NewTestObject()

	fields := d.GetFields()

	assert.Equal(t, 3, len(fields))
	field, ok := fields["ID"]

	assert.True(t, ok)
	assert.IsType(t, &DomainObjectField{}, field)
	assert.Equal(t, "int64", field.Type)
	assert.True(t, field.IsPrimaryKey)
	assert.False(t, field.IsNull)
}

func TestFieldList_swap(t *testing.T) {
	objA := &DomainObjectField{
		Name:    "a",
		Ordinal: 0,
	}
	objB := &DomainObjectField{
		Name:    "b",
		Ordinal: 1,
	}
	fieldList := FieldList{
		objA,
		objB,
	}

	fieldList.Swap(0, 1)

	assert.Equal(t, "b", fieldList[0].Name)
	assert.Equal(t, "a", fieldList[1].Name)
}
