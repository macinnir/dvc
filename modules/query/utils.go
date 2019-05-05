package query

import (
	"fmt"
	"sort"
	"strings"
)

// IQueryPart describes the functionality for a query part
type IQueryPart interface {
	ToSQL() (sql string, args []interface{})
}

func getSortedKeys(exp map[string]interface{}) []string {
	sortedKeys := make([]string, 0, len(exp))
	for k := range exp {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

func getSortedInKeys(exp map[string][]interface{}) []string {
	sortedKeys := make([]string, 0, len(exp))
	for k := range exp {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

func buildWhereClause(clauses []string) string {
	if len(clauses) > 1 {
		return strings.Join(clauses, " AND ")
	}

	return clauses[0]
}

func escapeField(field string) string {
	return fmt.Sprintf("`%s`", field)
}

func buildWhereClauseString(parts []IQueryPart) (sql string, args []interface{}) {

	wheres := []string{}

	if len(parts) > 0 {
		for _, where := range parts {
			whereSQL, whereArgs := where.ToSQL()
			wheres = append(wheres, whereSQL)
			if whereArgs != nil {
				for _, whereArg := range whereArgs {
					args = append(args, whereArg)
				}
			}
		}
	}

	if len(wheres) > 0 {
		sql = "WHERE " + strings.Join(wheres, " ") + " "
	}

	return
}
