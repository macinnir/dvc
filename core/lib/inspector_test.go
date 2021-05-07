package lib

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	src = []byte(`
		package main

		import (
		    "fmt"
		)

		// Person ...
		type Person struct {
		    name string
		    age int
		    telephone string
		}

		// Name ...
		func (p *Person) Name() string {
		    return p.name
		}

		// SetName ...
		func (p *Person) SetName(name string) {
		    p.name = name
		}

		// Age ...
		func (p *Person) Age() int {
		    return p.Age
		}

		// Age ...
		func (p *Person) SetAge(age int) {
		    p.Age = age
		}

		// AgeAndName ...
		func (p *Person) AgeAndName() (int, string) {
		    return p.age, p.name
		}

		func (p *Person) SetAgeAndName(name string, age int) {
		    p.name = name
		    p.age = age
		}

		// TelephoneAndName ...
		func (p *Person) GetNameAndTelephone() (name, telephone string) {
		    telephone = p.telephone
		    name = p.name 
		    return
		}

		func (p *Person) SetNameAndTelephone(name, telephone string) {
		    p.name = name
		    p.telephone = telephone
		}

		func SomeFunction() string {
		    return "Something"
		}`)
)

func TestMethod_GetLines(t *testing.T) {
	docs := []string{`// TestMethod is great`}
	code := `func TestMethod() string {return "I am great"}`

	method := Method{Code: code, Documents: docs}
	lines := method.Lines()

	assert.Equal(t, "// TestMethod is great", lines[0])
	assert.Equal(t, "func TestMethod() string {return \"I am great\"}", lines[1])
}

func TestParseStruct(t *testing.T) {
	methods, imports, typeDoc := ParseStruct(src, "Person", true, true, "")

	assert.Equal(t, "Name() (string)", methods[0].Code)

	imp := imports[0]
	trimmedImp := strings.TrimSpace(imp)

	assert.Equal(t, `"fmt"`, trimmedImp)
	assert.Equal(t, "Person ...", typeDoc)
}

func TestGetReceiverTypeName(t *testing.T) {
	fset := token.NewFileSet()
	a, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	assert.Nil(t, err, "ParseFile returned an error")

	hasPersonFuncDecl := false
	for _, d := range a.Decls {
		typeName, fd := GetReceiverTypeName(src, d)
		if typeName == "" {
			continue
		}
		switch typeName {
		case "Person":
			assert.NotNil(t, fd, "receiver type with name %s had a nil func decl")
			// OK
			hasPersonFuncDecl = true
		}
	}

	assert.True(t, hasPersonFuncDecl, "Never registered a func decl with the `Person` receiver type")
}

func TestFormatFieldList(t *testing.T) {
	fset := token.NewFileSet()
	a, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	assert.Nil(t, err, "ParseFile returned an error")

	for _, d := range a.Decls {
		if a, fd := GetReceiverTypeName(src, d); a == "Person" {
			methodName := fd.Name.String()
			params := FormatFieldList(src, fd.Type.Params, "")
			results := FormatFieldList(src, fd.Type.Results, "")

			var expectedParams []string
			var expectedResults []string
			switch methodName {
			case "Name":
				expectedResults = []string{"string"}
			case "SetName":
				expectedParams = []string{"name string"}
			case "Age":
				expectedResults = []string{"int"}
			case "SetAge":
				expectedParams = []string{"age int"}
			case "AgeAndName":
				expectedResults = []string{"int", "string"}
			case "SetAgeAndName":
				expectedParams = []string{"name string", "age int"}
			case "GetNameAndTelephone":
				expectedResults = []string{"name, telephone string"}
			case "SetNameAndTelephone":
				expectedParams = []string{"name, telephone string"}
			}
			assert.Equalf(t, expectedParams, params, "%s must have the expected params", methodName)
			assert.Equalf(t, expectedResults, results, "%s must have the expected results", methodName)
		}
	}
}

func TestParseStruct_NoCopyTypeDocs(t *testing.T) {
	_, _, typeDoc := ParseStruct(src, "Person", true, false, "")
	assert.Equal(t, "", typeDoc)
}
