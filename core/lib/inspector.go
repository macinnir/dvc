package lib

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"sort"
	"strings"

	"golang.org/x/tools/imports"
)

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

type ParsedStruct struct {
	Docs   []string
	Name   string
	Fields map[string]string
}

// Reflection: https://stackoverflow.com/questions/20170275/how-to-find-the-type-of-an-object-in-go
//

// func ParseStruct2(filePath string) (*ParsedStruct, error) {

// 	p := &ParsedStruct{
// 		Docs:   []string{},
// 		Fields: map[string]string{},
// 	}

// 	fileBytes, _ := ioutil.ReadFile(filePath)

// 	src := string(fileBytes)

// 	// Create the AST by parsing src.
// 	fset := token.NewFileSet() // positions are relative to fset
// 	f, err := parser.ParseFile(fset, path.Base(filePath), src, 0)
// 	if err != nil {
// 		return nil, err
// 	}

// 	fieldNo := 0

// 	for _, node := range f.Decls {
// 		switch node.(type) {

// 		case *ast.GenDecl:
// 			decl := node.(*ast.GenDecl)
// 			for _, spec := range decl.Specs {
// 				switch spec.(type) {
// 				case *ast.TypeSpec:
// 					typeSpec := spec.(*ast.TypeSpec)

// 					if decl.Doc != nil {
// 						for _, d := range decl.Doc.List {
// 							p.Docs = append(p.Docs, string(src[d.Pos()-1:d.End()-1]))
// 						}
// 					}

// 					p.Name = typeSpec.Name.Name

// 					switch typeSpec.Type.(type) {
// 					case *ast.StructType:
// 						structType := typeSpec.Type.(*ast.StructType)
// 						for _, field := range structType.Fields.List {

// 							fieldType := ""

// 							// switch field.Type.(type) {
// 							// case *ast.ArrayType:
// 							// fmt.Println(structType.Fields)
// 							// at := field.Type.(*ast.ArrayType)
// 							// fieldType = "[]" + fmt.Sprint(at.Elt)
// 							typeExpr := field.Type

// 							// Fetch Type from scource code
// 							typeInSource := src[typeExpr.Pos()-1 : typeExpr.End()-1]
// 							fmt.Println(typeExpr.Pos()-1, typeExpr.End()-1, typeInSource)
// 							fieldType = typeInSource
// 							// fmt.Println("ArrayType: ", at.Elt)
// 							// case *ast.Ident:
// 							// 	i := field.Type.(*ast.Ident)
// 							// 	fieldType = i.Name
// 							// }

// 							if len(field.Names) > 0 {
// 								for _, name := range field.Names {
// 									fmt.Printf("\tField: name=%s type=%s\n", name.Name, fieldType)
// 									p.Fields[name.Name] = fieldType

// 								}
// 							} else {
// 								p.Fields[fmt.Sprintf("#embedded%d", fieldNo)] = fieldType
// 							}

// 						}

// 					}
// 				}
// 			}
// 		}
// 	}

// 	return p, nil
// }

// ParseStruct parses a struct
func ParseStruct(src []byte, structName string, copyDocuments bool, copyTypeDocuments bool, pkgName string) ([]Method, []string, string) {

	var methods []Method
	var imports = []string{}
	var typeDoc string

	var file *ast.File
	var e error

	var fset = token.NewFileSet()

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
				Name:      decl.Name.String(),
				Code:      method,
				Documents: Documents,
				LineNo:    int(decl.Recv.Pos()),
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

	sort.Slice(methods, func(i, j int) bool { return methods[i].Name < methods[j].Name })

	return methods, imports, typeDoc
}

// Method describes code and Documents for a method
type Method struct {
	Name      string
	Code      string
	Documents []string
	LineNo    int
}

// Lines return a a slice of documentation and code
func (m *Method) Lines() []string {
	var lines []string
	lines = append(lines, m.Documents...)
	lines = append(lines, m.Code)
	return lines
}

// FormatFieldList returns a list of formatted list of fields
func FormatFieldList(src []byte, fieldList *ast.FieldList, pkgName string) (fields []string) {

	if fieldList == nil {
		return
	}

	// Loop through the list of fields in the field list
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

// GetFuncDeclarationReceiverType returns the receiver type of a function declaration
func GetFuncDeclarationReceiverType(fd *ast.FuncDecl) (expr ast.Expr, e error) {

	if fd.Recv == nil {
		e = fmt.Errorf("fd is a function, not a method")
		return
	}
	expr = fd.Recv.List[0].Type
	return

}

// GetReceiverTypeName returns the name of the receiver type and the declaration
func GetReceiverTypeName(src []byte, fl interface{}) (string, *ast.FuncDecl) {

	var name string
	var funcDef *ast.FuncDecl

	var ok bool
	var expr ast.Expr
	var e error

	if funcDef, ok = fl.(*ast.FuncDecl); !ok {
		return "", nil
	}

	if expr, e = GetFuncDeclarationReceiverType(funcDef); e != nil {
		return "", nil
	}

	name = string(src[expr.Pos()-1 : expr.End()-1])
	if len(name) > 0 && name[0] == '*' {
		name = name[1:]
	}
	return name, funcDef
}
