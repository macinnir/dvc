package query

import "sort"

// DomainObject is an interface that outlines functionality for every domain object
type DomainObject interface {
	GetFields() map[string]*DomainObjectField
	GetFieldsOrdered() []*DomainObjectField
	GetName() string
	Build()
}

// DomainObjectField is meta data about a field in a domain object
type DomainObjectField struct {
	Name         string
	Type         string
	IsPrimaryKey bool
	IsNull       bool
	Ordinal      int
}

// FieldList is a slice of fields in a domain object
type FieldList []*DomainObjectField

func (fieldList FieldList) Len() int { return len(fieldList) }

// Less is part of sort.Interface. We use count as the value to sort by
func (fieldList FieldList) Less(i, j int) bool { return fieldList[i].Ordinal < fieldList[j].Ordinal }

// Swap is part of sort.Interface.
func (fieldList FieldList) Swap(i, j int) { fieldList[i], fieldList[j] = fieldList[j], fieldList[i] }

// BaseDomainObject is the base object inherited by all domain objects
// Note: These fields include json tags for omission when the inheriting object attempts to marshal (e.g. for a JSON response from an API)
type BaseDomainObject struct {
	FieldList map[string]*DomainObjectField `json:"-"`
	TableName string                        `json:"-"`
}

// GetName returns the name of the domain object
func (b *BaseDomainObject) GetName() string {
	return b.TableName
}

// GetFields returns a map of the fields in the domain object
func (b *BaseDomainObject) GetFields() map[string]*DomainObjectField {
	return b.FieldList
}

// GetFieldsOrdered returns a slice of fields in the domain object, ordered by their ordinal
func (b *BaseDomainObject) GetFieldsOrdered() []*DomainObjectField {
	orderedFieldList := make(FieldList, 0, len(b.FieldList))
	for _, field := range b.FieldList {
		orderedFieldList = append(orderedFieldList, field)
	}
	sort.Sort(orderedFieldList)
	return orderedFieldList
}
