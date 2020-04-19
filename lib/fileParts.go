package lib

import (
	"strings"
)

// GoFileImports is a collection of imports for a model file
type GoFileImports []string

// Len returns the length of the collection of imports
func (i *GoFileImports) Len() int {
	return len(*i)
}

// Get returns an import at an index
func (i *GoFileImports) Get(idx int) string {
	return (*i)[idx]
}

// Append adds to the collection of imports
func (i *GoFileImports) Append(m string) {
	*i = append(*i, m)
}

// ToString returns a string representation of imports on a go file
func (i *GoFileImports) ToString() string {
	ret := ""
	if len(*i) > 0 {
		ret += "import ("
		for _, m := range *i {
			ret += "\n\t" + m
		}
		ret += "\n)\n"
	}

	return ret
}

// GoStruct represents a model struct
type GoStruct struct {
	Package  string
	Name     string
	Fields   *GoStructFields
	Comments string
	Imports  *GoFileImports
}

// NewGoStruct returns a new GoStruct
func NewGoStruct() *GoStruct {
	return &GoStruct{
		Fields:  &GoStructFields{},
		Imports: &GoFileImports{},
	}
}

// GoStructFields is a collection of a go struct fields
type GoStructFields []*GoStructField

// Len returns the length of the collection of imports
func (i *GoStructFields) Len() int {
	return len(*i)
}

// Get returns an import at an index
func (i *GoStructFields) Get(idx int) *GoStructField {
	return (*i)[idx]
}

// Append adds to the collection of imports
func (i *GoStructFields) Append(m *GoStructField) {
	*i = append(*i, m)
}

// GoStructField is a field on a Model struct
type GoStructField struct {
	Name     string
	DataType string
	Tags     []*GoStructFieldTag
	Comments string
}

// ToString returns a tring representation of a go struct field
func (m *GoStructField) ToString() string {

	str := m.Name +
		" " +
		m.DataType +
		" "
	if len(m.Tags) > 0 {
		str += "`"
		tags := []string{}
		for _, t := range m.Tags {
			tags = append(tags, t.ToString())
		}
		str += strings.Join(tags, " ")
		str += "`"
	}

	if len(m.Comments) > 0 {
		str += " // " + m.Comments
	} else {
		str += "\n"
	}

	return str

}

// GoStructFieldTag is a tag on a field on a model struct
type GoStructFieldTag struct {
	Name    string
	Value   string
	Options []string
}

// ToString returns a string representation of a go struct field tag
func (t *GoStructFieldTag) ToString() string {
	str := t.Name + ":\"" + t.Value
	if len(t.Options) > 0 {
		str += "," + strings.Join(t.Options, ",")
	}
	str += "\""
	return str
}
