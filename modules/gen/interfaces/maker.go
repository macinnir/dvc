package interfaces

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/tools/imports"
)

// Method describes code and Documents for a method
type Method struct {
	Code      string
	Documents []string
}

// Lines return a a slice of documentation and code
func (m *Method) Lines() []string {
	var lines []string
	lines = append(lines, m.Documents...)
	lines = append(lines, m.Code)
	return lines
}

// GetReceiverTypeName returns the name of the receiver type and the declaration
func GetReceiverTypeName(src []byte, fl interface{}) (name string, funcDef *ast.FuncDecl) {

	var ok bool
	var expr ast.Expr
	var e error

	if funcDef, ok = fl.(*ast.FuncDecl); !ok {
		return
	}

	if expr, e = GetReceiverType(funcDef); e != nil {
		return
	}

	name = string(src[expr.Pos()-1 : expr.End()-1])
	if len(name) > 0 && name[0] == '*' {
		name = name[1:]
	}
	return
}

// GetReceiverType returns the receiver type
func GetReceiverType(fd *ast.FuncDecl) (expr ast.Expr, e error) {

	if fd.Recv == nil {
		e = fmt.Errorf("fd is a function, not a method")
		return
	}
	expr = fd.Recv.List[0].Type
	return

}

// FormatFieldList formats the field list
func FormatFieldList(src []byte, fieldList *ast.FieldList, pkgName string) (fields []string) {

	if fieldList == nil {
		return
	}

	for _, l := range fieldList.List {

		names := make([]string, len(l.Names))

		for i, n := range l.Names {
			names[i] = n.Name
		}

		t := string(src[l.Type.Pos()-1 : l.Type.End()-1])

		t = strings.Replace(t, pkgName+".", "", -1)

		if len(names) > 0 {
			typeSharingArgs := strings.Join(names, ", ")
			fields = append(fields, fmt.Sprintf("%s %s", typeSharingArgs, t))
		} else {
			fields = append(fields, t)
		}
	}

	return fields
}

// FormatCode formats the code
func FormatCode(code string) ([]byte, error) {
	opts := &imports.Options{
		TabIndent: true,
		TabWidth:  2,
		Fragment:  true,
		Comments:  true,
	}
	return imports.Process("", []byte(code), opts)
}

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
	return FormatCode(code)
}

// ParseStruct parses a struct
func ParseStruct(src []byte, structName string, copyDocuments bool, copyTypeDocuments bool, pkgName string) (methods []Method, imports []string, typeDoc string) {

	var file *ast.File
	var e error

	fset := token.NewFileSet()

	if file, e = parser.ParseFile(fset, "", src, parser.ParseComments); e != nil {
		log.Fatal(e.Error())
	}

	for _, i := range file.Imports {
		if i.Name != nil {
			imports = append(imports, fmt.Sprintf("%s %s", i.Name.String(), i.Path.Value))
		} else {
			imports = append(imports, fmt.Sprintf("%s", i.Path.Value))
		}
	}

	for _, d := range file.Decls {

		if a, decl := GetReceiverTypeName(src, d); a == structName {

			if !decl.Name.IsExported() {
				continue
			}

			params := FormatFieldList(src, decl.Type.Params, pkgName)
			ret := FormatFieldList(src, decl.Type.Results, pkgName)
			method := fmt.Sprintf("%s(%s) (%s)", decl.Name.String(), strings.Join(params, ", "), strings.Join(ret, ", "))

			var Documents []string

			if decl.Doc != nil && copyDocuments {
				for _, d := range decl.Doc.List {
					Documents = append(Documents, string(src[d.Pos()-1:d.End()-1]))
				}
			}

			methods = append(methods, Method{
				Code:      method,
				Documents: Documents,
			})

		}
	}

	if copyTypeDocuments {
		pkg := &ast.Package{Files: map[string]*ast.File{"": file}}
		doc := doc.New(pkg, "", doc.AllDecls)
		for _, t := range doc.Types {
			if t.Name == structName {
				typeDoc = strings.TrimSuffix(t.Doc, "\n")
			}
		}
	}

	return
}

// Run runs GenInterface on all the target files
func Run(files []string, structType, comment, pkgName, ifaceName, ifaceComment string, copyDocuments, copyTypeDoc bool) (result []byte, e error) {

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

		methods, imports, parsedTypeDoc := ParseStruct(src, structType, copyDocuments, copyTypeDoc, pkgName)
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
