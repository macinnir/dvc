package typescript

import (
	"testing"

	"github.com/macinnir/dvc/core/lib/schema"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTypescriptType(t *testing.T) {

	g, e := GenerateTypescriptType(fooObj)

	assert.Nil(t, e)
	assert.Equal(t, tsFoo, g)
}

var fooObj = &schema.Table{
	Name: "Foo",
	Columns: map[string]*schema.Column{
		"FooID": {Name: "FooID", DataType: "int"},
		"Bar":   {Name: "Bar", DataType: "varchar", MaxLength: 200},
		"Baz":   {Name: "Baz", DataType: "datetime"},
	},
}

var tsFoo = `/**
 * Generated Code; DO NOT EDIT
 * 
 * Foo
 */
export type Foo = {

	// Bar varchar(200)
	Bar: string;

	// Baz datetime
	Baz: string;

	// FooID int
	FooID: number;

}

// newFoo is a factory method for new Foo objects
export const newFoo = () : Foo => ({
	Bar: '',
	Baz: '',
	FooID: 0,
});

`

var barObj = &schema.Table{
	Name: "Bar",
	Columns: map[string]*schema.Column{
		"FooID2": {Name: "FooID2", DataType: "int"},
		"Bar2":   {Name: "Bar2", DataType: "varchar", MaxLength: 200},
		"Baz2":   {Name: "Baz2", DataType: "datetime"},
	},
}

var tsBar = `/**
* Bar
*/
export type Bar = {

   // Bar2 varchar(200)
   Bar2: string;

   // Baz2 datetime
   Baz2: string;

   // FooID2 int
   FooID2: number;

}

// newFoo is a factory method for new Foo objects
export const newBar = () : Bar => ({
   Bar2: '',
   Baz2: '',
   FooID2: 0,
});

`
