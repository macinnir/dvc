package gen

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/macinnir/dvc/lib"
)

// GenInterface takes makes the interface into a byte array
func GenInterface(comment, pkgName, ifaceName, ifaceComment string, methods []string, imports []string) ([]byte, error) {
	output := []string{"// " + comment, "", "package " + pkgName, "import ("}
	output = append(output, imports...)
	output = append(output, ")", "")
	if len(ifaceComment) > 0 {
		output = append(output, fmt.Sprintf("// %s", strings.Replace(ifaceComment, "\n", "\n// ", -1)))
	}
	output = append(output, fmt.Sprintf("type %s interface {", ifaceName))
	output = append(output, methods...)
	output = append(output, "}")
	code := strings.Join(output, "\n")
	return lib.FormatCode(code)
}

// GenInterfaces runs GenInterface on all the target files
func GenInterfaces(files []string, structType, comment, pkgName, ifaceName, ifaceComment string, copyDocuments, copyTypeDoc bool) (result []byte, e error) {

	allMethods := []string{}
	allImports := []string{}

	mset := make(map[string]struct{})
	iset := make(map[string]struct{})

	var typeDoc string

	for _, f := range files {

		var src []byte
		if src, e = ioutil.ReadFile(f); e != nil {
			return
		}

		methods, imports, parsedTypeDoc := lib.ParseStruct(src, structType, copyDocuments, copyTypeDoc, pkgName)
		for _, m := range methods {
			if _, ok := mset[m.Code]; !ok {
				allMethods = append(allMethods, m.Lines()...)
				mset[m.Code] = struct{}{}
			}
		}
		for _, i := range imports {
			if _, ok := iset[i]; !ok {
				allImports = append(allImports, i)
				iset[i] = struct{}{}
			}
		}
		if typeDoc == "" {
			typeDoc = parsedTypeDoc
		}
	}

	if typeDoc != "" {
		ifaceComment += "\n" + typeDoc
	}

	result, e = GenInterface(comment, pkgName, ifaceName, ifaceComment, allMethods, allImports)

	return
}
