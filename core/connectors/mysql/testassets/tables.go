package testassets

import "github.com/macinnir/dvc/core/lib/schema"

func TablesNoChange() []*schema.Schema {
	return []*schema.Schema{
		{
			Tables: map[string]*schema.Table{
				"Foo": {
					Name: "Foo",
					Columns: map[string]*schema.Column{
						"FooID":       columnFooID,
						"DateCreated": columnDateCreated,
						"IsDeleted":   columnIsDeleted,
					},
				},
			},
		},
		{
			Tables: map[string]*schema.Table{
				"Foo": {
					Name: "Foo",
					Columns: map[string]*schema.Column{
						"FooID":       columnFooID,
						"DateCreated": columnDateCreated,
						"IsDeleted":   columnIsDeleted,
					},
				},
			},
		},
	}
}

func TablesDropColumn() []*schema.Schema {
	return []*schema.Schema{
		{
			Tables: map[string]*schema.Table{
				"Foo": {
					Name: "Foo",
					Columns: map[string]*schema.Column{
						"FooID":       columnFooID,
						"DateCreated": columnDateCreated,
						"IsDeleted":   columnIsDeleted,
					},
				},
			},
		},
		{
			Tables: map[string]*schema.Table{
				"Foo": {
					Name: "Foo",
					Columns: map[string]*schema.Column{
						"FooID":       columnFooID,
						"Name":        columnName,
						"DateCreated": columnDateCreated,
						"IsDeleted":   columnIsDeleted,
					},
				},
			},
		},
	}
}

func TablesAddColumn() []*schema.Schema {
	return []*schema.Schema{
		{
			Tables: map[string]*schema.Table{
				"Foo": {
					Name: "Foo",
					Columns: map[string]*schema.Column{
						"FooID":       columnFooID,
						"Name":        columnName,
						"DateCreated": columnDateCreated,
						"IsDeleted":   columnIsDeleted,
					},
				},
			},
		},
		{
			Tables: map[string]*schema.Table{
				"Foo": {
					Name: "Foo",
					Columns: map[string]*schema.Column{
						"FooID":       columnFooID,
						"DateCreated": columnDateCreated,
						"IsDeleted":   columnIsDeleted,
					},
				},
			},
		},
	}
}

func TablesChangeVarcharColumnSize() []*schema.Schema {
	return []*schema.Schema{
		{
			Tables: map[string]*schema.Table{
				"Foo": {
					Name: "Foo",
					Columns: map[string]*schema.Column{
						"FooID":       columnFooID,
						"Name":        columnName,
						"DateCreated": columnDateCreated,
						"IsDeleted":   columnIsDeleted,
					},
				},
			},
		},
		{
			Tables: map[string]*schema.Table{
				"Foo": {
					Name: "Foo",
					Columns: map[string]*schema.Column{
						"FooID":       columnFooID,
						"Name":        columnNameSize100,
						"DateCreated": columnDateCreated,
						"IsDeleted":   columnIsDeleted,
					},
				},
			},
		},
	}
}

func TablesAddColumnWithUniqueIndex() []*schema.Schema {
	return []*schema.Schema{
		{
			Tables: map[string]*schema.Table{
				"Foo": {
					Name: "Foo",
					Columns: map[string]*schema.Column{
						"FooID":       columnFooID,
						"Name":        columnNameWithUniqueIndex,
						"DateCreated": columnDateCreated,
						"IsDeleted":   columnIsDeleted,
					},
				},
			},
		},
		{
			Tables: map[string]*schema.Table{
				"Foo": {
					Name: "Foo",
					Columns: map[string]*schema.Column{
						"FooID":       columnFooID,
						"DateCreated": columnDateCreated,
						"IsDeleted":   columnIsDeleted,
					},
				},
			},
		},
	}
}

func TablesAddColumnWithIndex() []*schema.Schema {
	return []*schema.Schema{
		{
			Tables: map[string]*schema.Table{
				"Foo": {
					Name: "Foo",
					Columns: map[string]*schema.Column{
						"FooID":       columnFooID,
						"Name":        columnNameWithIndex,
						"DateCreated": columnDateCreated,
						"IsDeleted":   columnIsDeleted,
					},
				},
			},
		},
		{
			Tables: map[string]*schema.Table{
				"Foo": {
					Name: "Foo",
					Columns: map[string]*schema.Column{
						"FooID":       columnFooID,
						"DateCreated": columnDateCreated,
						"IsDeleted":   columnIsDeleted,
					},
				},
			},
		},
	}
}

func TablesDropAutoIncrement() []*schema.Schema {
	return []*schema.Schema{
		{
			Tables: map[string]*schema.Table{
				"Foo": {
					Name: "Foo",
					Columns: map[string]*schema.Column{
						"FooID":       columnFooIDBigIntNoAI,
						"Name":        columnName,
						"DateCreated": columnDateCreated,
						"IsDeleted":   columnIsDeleted,
					},
				},
			},
		},
		{
			Tables: map[string]*schema.Table{
				"Foo": {
					Name: "Foo",
					Columns: map[string]*schema.Column{
						"FooID":       columnFooID,
						"Name":        columnName,
						"DateCreated": columnDateCreated,
						"IsDeleted":   columnIsDeleted,
					},
				},
			},
		},
	}
}

func TablesAddTable() []*schema.Schema {
	return []*schema.Schema{
		{
			Tables: map[string]*schema.Table{
				"Foo": {
					Name: "Foo",
					Columns: map[string]*schema.Column{
						"FooID":       columnFooID,
						"DateCreated": columnDateCreated,
						"IsDeleted":   columnIsDeleted,
					},
				},
			},
		},
		{},
	}
}

func TablesDropTable() []*schema.Schema {
	return []*schema.Schema{
		{},
		{
			Tables: map[string]*schema.Table{
				"Foo": {
					Name: "Foo",
					Columns: map[string]*schema.Column{
						"FooID":       columnFooID,
						"DateCreated": columnDateCreated,
						"IsDeleted":   columnIsDeleted,
					},
				},
			},
		},
	}
}
