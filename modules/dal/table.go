package dal

import "fmt"

// NewTable creates a new table object
func NewTable(name string) *Table {
	t := new(Table)
	t.Name = name
	t.Fields = map[string]*Field{}
	return t
}

// Table represents a database table
type Table struct {
	Name      string
	Fields    map[string]*Field
	Alias     string
	FieldKeys []string
}

// AddFields adds a collection of fields
func (t *Table) AddFields(fieldNames []string) {
	for _, fieldName := range fieldNames {
		t.AddField(fieldName)
	}
}

// AddField adds a new field to the database object
func (t *Table) AddField(fieldName string) (e error) {

	var ok bool

	if _, ok = t.Fields[fieldName]; ok {
		e = fmt.Errorf("the field `%s`.`%s` has already been defined", t.Name, fieldName)
		return
	}

	newField := new(Field)
	newField.Name = fieldName
	t.Fields[fieldName] = newField
	t.FieldKeys = append(t.FieldKeys, fieldName)

	return
}

// Field gets a field from the table by its name
func (t *Table) Field(fieldName string) *Field {

	var field *Field
	var ok bool

	if field, ok = t.Fields[fieldName]; !ok {
		panic(fmt.Sprintf("Field `%s`.`%s` not found", t.Name, fieldName))
	}

	return field
}
