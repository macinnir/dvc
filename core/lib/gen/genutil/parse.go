package genutil

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"

	"github.com/fatih/structtag"
	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

// InspectFile inspects a file
func ParseFileToGoStruct(filePath string) (*lib.GoStruct, error) {

	var s *lib.GoStruct
	var e error

	fileBytes, e := ioutil.ReadFile(filePath)
	if e != nil {
		return nil, e
	}

	s, e = ParseStringToGoStruct(fileBytes)
	if e != nil {
		fmt.Println("ERROR: ", filePath)
		return nil, e
	}

	return s, nil

}

// buildModelNodeFromFile builds a node representation of a struct from a file
func ParseStringToGoStruct(src []byte) (*lib.GoStruct, error) {

	var e error
	var modelNode = lib.NewGoStruct()
	var tree *ast.File

	var srcString = string(src)
	_, tree, e = ParseFileToAST(src)

	if e != nil {
		return nil, e
	}

	// typeDecl := tree.Decls[0].(*ast.GenDecl)
	// structDecl := typeDecl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType)
	// fields := structDecl.Fields.List

	// for k := range fields {
	// 	typeExpr := fields[k].Type
	// 	start := typeExpr.Pos() - 1
	// 	end := typeExpr.End() - 1

	// 	typeInSource := src[start:end]

	// 	fmt.Println(typeInSource)
	// }

	ast.Inspect(tree, func(node ast.Node) bool {

		// Check if this is a package
		if s, ok := node.(*ast.File); ok {

			modelNode.Package = s.Name.Name
			if len(s.Comments) > 0 {
				modelNode.Comments = s.Comments[0].Text()
			}
			modelNode.Imports = &lib.GoFileImports{}

			for _, i := range s.Imports {
				// This is a named import
				if i.Name != nil {
					modelNode.Imports.Append(i.Name.Name + " " + i.Path.Value)
				} else {
					modelNode.Imports.Append(i.Path.Value)
				}
			}

			// for _, d := range s.Decls {
			// 	GetReceiverTypeName
			// }

		}

		// Declaration of our struct
		if s, ok := node.(*ast.TypeSpec); ok {
			if len(modelNode.Name) == 0 {
				// fmt.Println("Type Name: ", s.Name.Name)
				modelNode.Name = s.Name.Name
			}
		}

		if s, ok := node.(*ast.StructType); !ok {

			return true

		} else {

			for _, field := range s.Fields.List {

				fieldName := field.Names[0].Name

				if fieldName == "db" || fieldName == "isSingle" || fieldName == "q" {
					continue
				}

				fieldType := srcString[field.Type.Pos()-1 : field.Type.End()-1]
				nodeField := &lib.GoStructField{
					Name:     fieldName,
					Tags:     []*lib.GoStructFieldTag{},
					DataType: fieldType,
					Comments: field.Comment.Text(),
				}
				if field.Tag != nil {
					tagString := field.Tag.Value[1 : len(field.Tag.Value)-1]
					// fmt.Printf("Tag: %s\n", tagString)
					tags, e := structtag.Parse(tagString)
					if e != nil {
						log.Fatal(e)
					}
					for _, tag := range tags.Tags() {
						nodeField.Tags = append(nodeField.Tags, &lib.GoStructFieldTag{
							Name:    tag.Key,
							Value:   tag.Name,
							Options: tag.Options,
						})
					}
				}

				modelNode.Fields.Append(nodeField)
			}
		}

		return false
	})

	return modelNode, nil
}

// ParseFileToAST takes a file path and parses the contents of that file into
// an AST representation
func ParseFileToAST(fileBytes []byte) (*token.FileSet, *ast.File, error) {

	var fileSet = token.NewFileSet()

	var tree, e = parser.ParseFile(fileSet, "", fileBytes, parser.ParseComments)
	if e != nil {
		return nil, nil, e
	}

	return fileSet, tree, nil
}

// Deprecated
func ResolveTableToModel(modelNode *lib.GoStruct, table *schema.Table) {

	fieldMap := map[string]int{}
	modelFields := &lib.GoStructFields{}

	nullImportIndex := -1
	hasNullField := false

	for k, i := range *modelNode.Imports {
		if i == lib.NullPackage {
			nullImportIndex = k
			break
		}
	}

	i := 0
	for _, m := range *modelNode.Fields {

		// Skip any fields not in the database
		if _, ok := table.Columns[m.Name]; !ok {
			continue
		}

		fieldMap[m.Name] = i
		modelFields.Append(m)
		i++
	}

	for name, col := range table.Columns {

		fieldIndex, ok := fieldMap[name]

		// Add any fields not in the model
		if !ok {
			modelFields.Append(&lib.GoStructField{
				Name:     col.Name,
				DataType: schema.DataTypeToGoTypeString(col),
				Tags: []*lib.GoStructFieldTag{
					{
						Name:    "db",
						Value:   col.Name,
						Options: []string{},
					},
					{
						Name:    "json",
						Value:   col.Name,
						Options: []string{},
					},
				},
			})
		} else {

			// Check that the datatype hasn't changed
			colDataType := schema.DataTypeToGoTypeString(col)

			// log.Println(colDataType, fieldIndex, name)

			if colDataType != (*modelFields)[fieldIndex].DataType {
				(*modelFields)[fieldIndex].DataType = colDataType
			}
		}
	}

	// Finally check for nullables
	for _, f := range *modelFields {

		if schema.IsNull(f.DataType) {
			hasNullField = true
		}
	}

	// If the package needs null, and it hasn't been added, add it
	if hasNullField && nullImportIndex == -1 {
		modelNode.Imports.Append(lib.NullPackage)
	}

	// If no null import is needed, but one exists, remove it
	if !hasNullField && nullImportIndex > -1 {
		*modelNode.Imports = append((*modelNode.Imports)[:nullImportIndex], (*modelNode.Imports)[nullImportIndex+1:]...)
	}

	modelNode.Fields = modelFields
	return
}

func ParseFileNameToModelName(fileName, prefix, suffix string) string {

	// Remove .go or .ts extension
	var rootName = fileName[0 : len(fileName)-3]
	var modelName = rootName

	if len(prefix) > 0 {
		modelName = modelName[len(prefix):]
	}

	if len(modelName) > 5 && modelName[len(modelName)-5:] == "_test" {
		modelName = modelName[0 : len(modelName)-5]
	}

	if len(suffix) > 0 {
		modelName = modelName[0 : len(modelName)-len(suffix)]
	}

	return modelName
}
