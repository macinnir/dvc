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
