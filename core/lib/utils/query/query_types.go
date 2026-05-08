package query

import "strings"

type QueryType int

const (
	QueryTypeSelect QueryType = iota
	QueryTypeInsert
	QueryTypeUpdate
	QueryTypeDelete
	QueryTypeRaw
)

func Select(model IModel) *Q {
	q := Query(model)
	q.queryType = QueryTypeSelect
	return q
}

func Raw(model IModel, query string) *Q {
	q := Query(model)
	q.queryType = QueryTypeRaw
	q.raw = query
	return q
}

func Update(model IModel) *Q {
	q := Query(model)
	q.queryType = QueryTypeUpdate
	return q
}

func Delete(model IModel) *Q {
	q := Query(model)
	q.queryType = QueryTypeDelete
	return q
}

func Insert(model IModel) *Q {
	q := Query(model)
	q.queryType = QueryTypeInsert
	return q
}

func Union(queries ...*Q) (string, error) {

	sqls := []string{}
	for k := range queries {
		query, e := queries[k].String()
		if e != nil {
			return "", e
		}
		sqls = append(sqls, query)
	}

	return strings.Join(sqls, " UNION ALL "), nil
}

// Exists is a where clause for the SQL EXISTS statement
func Exists(clause *Q) *WherePart {
	clauseString, e := clause.String()

	w := newWherePart(
		WhereTypeExists,
		"",
		[]interface{}{clauseString},
	)
	if e != nil {
		w.e = e
	}
	return w
}

// Exists is a where clause for the SQL EXISTS statement
func NotExists(clause *Q) *WherePart {
	clauseString, e := clause.String()

	w := newWherePart(
		WhereTypeNotExists,
		"",
		[]interface{}{clauseString},
	)
	if e != nil {
		w.e = e
	}
	return w
}
