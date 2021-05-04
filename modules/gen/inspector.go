package gen

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"sort"

	"github.com/fatih/structtag"
	"github.com/macinnir/dvc/lib"
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
			modelNode.Name = s.Name.Name
		}

		if s, ok := node.(*ast.StructType); !ok {

			return true

		} else {

			for _, field := range s.Fields.List {

				fieldType := src[field.Type.Pos()-1 : field.Type.End()-1]
				nodeField := &lib.GoStructField{
					Name:     field.Names[0].Name,
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

// buildModelNodeFromFile builds a node representation of a struct from a file
func buildModelNodeFromTable(table *lib.Table) (modelNode *lib.GoStruct, e error) {

	modelNode = lib.NewGoStruct()
	modelNode.Package = "models"
	modelNode.Name = table.Name
	modelNode.Comments = fmt.Sprintf("%s is a `%s` data model\n", table.Name, table.Name)

	hasNull := false

	sortedColumns := make(lib.SortedColumns, 0, len(table.Columns))

	for _, column := range table.Columns {
		sortedColumns = append(sortedColumns, column)
	}

	sort.Sort(sortedColumns)

	for _, col := range sortedColumns {
		fieldType := lib.DataTypeToGoTypeString(col)
		if fieldDataTypeIsNull(fieldType) {
			hasNull = true
		}
		modelNode.Fields.Append(&lib.GoStructField{
			Name:     col.Name,
			DataType: fieldType,
			Tags: []*lib.GoStructFieldTag{
				{Name: "db", Value: col.Name, Options: []string{}},
				{Name: "json", Value: col.Name, Options: []string{}},
			},
			Comments: "",
		})
	}

	if hasNull {
		modelNode.Imports.Append(NullPackage)
	}

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

func buildFileFromModelNode(modelNode *lib.GoStruct) (file []byte, e error) {

	fileString := "// Generated Code; DO NOT EDIT.\n\npackage " + modelNode.Package + "\n\n"
	if modelNode.Imports.Len() > 0 {
		fileString += modelNode.Imports.ToString() + "\n"
	}

	if len(modelNode.Comments) > 0 {
		fileString += "// " + modelNode.Comments
	}

	fileString += "type " + modelNode.Name + " struct {\n"

	for _, f := range *modelNode.Fields {
		fileString += "\t" + f.ToString()
	}

	fileString += "}\n"

	file = []byte(fileString)

	file, e = format.Source(file)
	if e != nil {
		log.Fatalf("FORMAT ERROR: File: %s; Error: %s\n%s", modelNode.Name, e.Error(), fileString)
	}
	return
}

func resolveTableToModel(modelNode *lib.GoStruct, table *lib.Table) {

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
				DataType: lib.DataTypeToGoTypeString(col),
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
			colDataType := lib.DataTypeToGoTypeString(col)

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
