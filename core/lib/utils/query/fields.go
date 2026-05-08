package query

type FieldType int

const (
	FieldTypeBasic FieldType = iota
	FieldTypeRaw
	FieldTypeCount
	FieldTypeSum
	FieldTypeAvg
	FieldTypeMin
	FieldTypeMax
	FieldTypeDistinct
)

type Field struct {
	FieldType FieldType
	Name      Column
	As        string
	Raw       string
}

// NewField creates a new field.
//
//	NewField(FieldTypeBasic, "Foo")
//	NewField(FieldTypeBasic, "Foo", "Bar") <-- `Foo` AS `Bar`
func NewField(fieldType FieldType, column Column, opts ...string) *Field {

	as := ""

	if len(opts) > 0 {
		as = opts[0]
	}

	return &Field{
		FieldType: fieldType,
		Name:      column,
		As:        as,
		Raw:       "",
	}
}

// NewRawField creates a new field.
//
//	NewRawField("`t`.`Foo` AS `Bar`)
func NewRawField(raw string) *Field {

	return &Field{
		FieldType: FieldTypeRaw,
		Name:      "",
		As:        "",
		Raw:       raw,
	}
}

// Fields injects fields as raw strings into the field clause of the query
//
//	sql, e := query.Select(&testassets.Job{}).
//		Fields(
//			NewField(FieldTypeBasic, "JobID"),
//			NewField(FieldTypeBasic, "Name", "Foo"),
//		)
func (q *Q) Fields(fields ...*Field) *Q {
	q.fields = fields
	return q
}

func (q *Q) Raw(query string) *Q {
	q.raw = query
	return q
}

// Field includes a specific field in the columns to be returned by a result set
func (q *Q) Field(name Column) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "SELECT...Field", string(name))
	}

	q.fields = append(q.fields, NewField(FieldTypeBasic, name))

	return q
}

// FieldAs includes a specific field in the columns to be returned by a set aliased by `as`
func (q *Q) FieldAs(name Column, as string) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "SELECT...Field...as", string(name))
	}

	q.fields = append(q.fields, NewField(FieldTypeBasic, name, as))

	return q
}

// FieldRaw allows for an arbitrary string (e.g. "NOW()") to be included in the select columns
func (q *Q) FieldRaw(fieldStr, as string) *Q {
	q.fields = append(q.fields, NewRawField(fieldStr+" AS "+"`"+as+"`"))

	return q
}

// Count creates a count statement
//
//	q.Count(query.Column("Foo"), "FooCounted")
//	COALESCE(COUNT(`t`.`Foo`), 0) AS `FooCounted`
func (q *Q) Count(name Column, as string) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "SELECT...COUNT()", string(name))
		return q
	}

	q.fields = append(q.fields, NewField(FieldTypeCount, name, as))
	return q

	// return q.FieldRaw("COUNT(`"+q.alias+"`.`"+string(name)+"`)", as)
}

// Sum creates a sum statement
//
//	q.Sum(query.Column("Foo"), "FooSummed")
//	COALESCE(SUM(`t`.`Foo`), 0) AS `FooSummed`
func (q *Q) Sum(name Column, as string) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "SELECT...Sum()", string(name))
		return q
	}

	q.fields = append(q.fields, NewField(FieldTypeSum, name, as))
	return q

	// return q.FieldRaw("COALESCE(SUM(`"+q.alias+"`.`"+string(name)+"`), 0)", as)
}

// Avg creates an Avg statement
//
//	q.Avg(query.Column("Foo"), "FooAveraged")
//	COALESCE(AVG(`t`.`Foo`), 0) AS `FooAveraged`
func (q *Q) Avg(name Column, as string) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "SELECT...Sum()", string(name))
		return q
	}

	q.fields = append(q.fields, NewField(FieldTypeAvg, name, as))
	return q

	// return q.FieldRaw("COALESCE(SUM(`"+q.alias+"`.`"+string(name)+"`), 0)", as)
}

// Min creates a min statement
//
//	q.Min(query.Column("Foo"), "MinFoo")
//	COALESCE(MIN(`t`.`Foo`), 0) AS `MinFoo`
func (q *Q) Min(name Column, as string) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "SELECT...Min()", string(name))
		return q
	}

	q.fields = append(q.fields, NewField(FieldTypeMin, name, as))
	return q

	// return q.FieldRaw("COALESCE(MIN(`"+q.alias+"`.`"+string(name)+"`), 0)", as)
}

// Distinct creates a distinct statement
//
//	q.Distinct(query.Column("Foo"), "DISTINCT")
//	DISTINCT `t`.`Foo`
func (q *Q) Distinct(name Column) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "SELECT...Distinct()", string(name))
		return q
	}

	q.fields = append(q.fields, NewField(FieldTypeDistinct, name))
	return q
}

// Max creates a max statement
//
//	q.Max(query.Column("Foo"), "MaxFoo")
//	COALESCE(MAX(`t`.`Foo`), 0) AS `MaxFoo`
func (q *Q) Max(name Column, as string) *Q {

	if _, ok := q.columnTypes[name]; !ok {
		q.errorInvalidColumn(QUERY_ERROR_INVALID_COLUMN, "SELECT...Max()", string(name))
		return q
	}

	q.fields = append(q.fields, NewField(FieldTypeMax, name, as))
	return q

	// return q.FieldRaw("COALESCE(MAX(`"+q.alias+"`.`"+string(name)+"`), 0)", as)
}

func (q *Q) FromFieldToString(field *Field) string {

	as := ""

	if len(field.As) > 0 {
		as = " AS `" + field.As + "`"
	}

	switch field.FieldType {
	case FieldTypeCount:
		return "COUNT(`" + q.alias + "`.`" + string(field.Name) + "`)" + as
	case FieldTypeSum:
		return "COALESCE(SUM(`" + q.alias + "`.`" + string(field.Name) + "`), 0)" + as
	case FieldTypeAvg:
		return "COALESCE(AVG(`" + q.alias + "`.`" + string(field.Name) + "`), 0)" + as
	case FieldTypeMin:
		return "COALESCE(MIN(`" + q.alias + "`.`" + string(field.Name) + "`), 0)" + as
	case FieldTypeMax:
		return "COALESCE(MAX(`" + q.alias + "`.`" + string(field.Name) + "`), 0)" + as
	case FieldTypeDistinct:
		return "DISTINCT `" + q.alias + "`.`" + string(field.Name) + "`" + as
	case FieldTypeRaw:
		return field.Raw
	// FieldTypeBasic
	default:
		return "`" + q.alias + "`.`" + string(field.Name) + "`" + as
	}
}
