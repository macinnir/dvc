package gen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"github.com/fatih/structtag"
	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

// NullPackage is the package name used for handling nulls
const NullPackage = "\"gopkg.in/guregu/null.v3\""

// buildModelNodeFromFile builds a node representation of a struct from a file
func buildGoStructFromFile(fileBytes []byte) (modelNode *lib.GoStruct, e error) {

	src := string(fileBytes)
	modelNode = lib.NewGoStruct()
	var tree *ast.File

	_, tree, e = parseFileToAST(fileBytes)

	if e != nil {
		return
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

				fieldType := src[field.Type.Pos()-1 : field.Type.End()-1]
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

	return
}

func fieldDataTypeIsNull(fieldType string) bool {
	return len(fieldType) > 5 && fieldType[0:5] == "null."
}

// parseFileToAST takes a file path and parses the contents of that file into
// an AST representation
func parseFileToAST(fileBytes []byte) (fileSet *token.FileSet, tree *ast.File, e error) {

	fileSet = token.NewFileSet()

	tree, e = parser.ParseFile(fileSet, "", fileBytes, parser.ParseComments)
	if e != nil {
		return
	}

	return
}

func resolveTableToModel(modelNode *lib.GoStruct, table *schema.Table) {

	fieldMap := map[string]int{}
	modelFields := &lib.GoStructFields{}

	nullImportIndex := -1
	hasNullField := false

	for k, i := range *modelNode.Imports {
		if i == NullPackage {
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

		if fieldDataTypeIsNull(f.DataType) {
			hasNullField = true
		}
	}

	// If the package needs null, and it hasn't been added, add it
	if hasNullField && nullImportIndex == -1 {
		modelNode.Imports.Append(NullPackage)
	}

	// If no null import is needed, but one exists, remove it
	if !hasNullField && nullImportIndex > -1 {
		*modelNode.Imports = append((*modelNode.Imports)[:nullImportIndex], (*modelNode.Imports)[nullImportIndex+1:]...)
	}

	modelNode.Fields = modelFields
	return
}
