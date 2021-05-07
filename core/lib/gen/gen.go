package gen

import (
	"regexp"
	"strings"

	"github.com/macinnir/dvc/core/lib"
)

// Gen conntains all of the generator functionality
type Gen struct {
	Config *lib.Config
}

//
// String Generators
//

// scanStringForFuncSignature scans a string (a line of goCode) and returns matches if it is a golang function signature that matches
// signatureRegexp
func (g *Gen) scanStringForFuncSignature(str string, signatureRegexp string) (matches []string) {

	lines := strings.Split(str, "\n")

	var validSignature = regexp.MustCompile(signatureRegexp)

	matches = []string{}

	for _, line := range lines {
		if validSignature.Match([]byte(line)) {
			matches = append(matches, line)
		}
	}

	return
}
