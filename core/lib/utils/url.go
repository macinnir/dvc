package utils

import (
	"strconv"
	"strings"
)

func QueryIDs(query string) []int64 {

	idsString := QueryStrings(query)

	ids := []int64{}
	for k := 0; k < len(idsString); k++ {
		n := int64(0)
		var e error
		if n, e = strconv.ParseInt(idsString[k], 10, 64); e != nil {
			continue
		}

		ids = append(ids, n)
	}

	return ids
}

func QueryStrings(query string) []string {

	var idsString []string

	if strings.Contains(query, ",") {
		idsString = strings.Split(query, ",")
	} else {
		idsString = []string{query}
	}

	return idsString

}
