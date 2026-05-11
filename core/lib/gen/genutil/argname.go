package genutil

import "strings"

func ToArgName(field string) string {
	if len(field) == 0 {
		return ""
	}
	return strings.ToLower(field[:1]) + field[1:]
}
