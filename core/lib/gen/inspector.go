package gen

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"sort"
	"strings"

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
func buildModelNodeFromTable(table *schema.Table) (modelNode *lib.GoStruct, e error) {

	modelNode = lib.NewGoStruct()
	modelNode.Package = "models"
	modelNode.Name = table.Name
	modelNode.Comments = fmt.Sprintf("%s is a `%s` data model\n", table.Name, table.Name)

	hasNull := false

	sortedColumns := make(schema.SortedColumns, 0, len(table.Columns))

	for _, column := range table.Columns {
		sortedColumns = append(sortedColumns, column)
	}

	sort.Sort(sortedColumns)

	for _, col := range sortedColumns {
		fieldType := schema.DataTypeToGoTypeString(col)
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

func buildFileFromModelNode(table *schema.Table, modelNode *lib.GoStruct) (file []byte, e error) {

	insertColumns := fetchInsertColumns(table.ToSortedColumns())
	updateColumns := fetchUpdateColumns(table.ToSortedColumns())
	primaryKey := fetchTablePrimaryKeyName(table)

	var b strings.Builder
	b.WriteString("// Generated Code; DO NOT EDIT.\n\npackage " + modelNode.Package + "\n\n")
	if modelNode.Imports.Len() > 0 {
		b.WriteString(modelNode.Imports.ToString() + "\n")
	}

	if len(modelNode.Comments) > 0 {
		b.WriteString("// " + modelNode.Comments)
	}

	b.WriteString("type " + modelNode.Name + " struct {\n")

	for _, f := range *modelNode.Fields {
		b.WriteString("\t" + f.ToString())
	}

	b.WriteString("}\n")

	b.WriteString(`
// ` + modelNode.Name + `_Column is the type used for ` + modelNode.Name + ` columns
type ` + modelNode.Name + `_Column string

// ` + modelNode.Name + `_Columns specifies the columns in the ` + modelNode.Name + ` model
var ` + modelNode.Name + `_Columns = struct {
`)
	for _, f := range *modelNode.Fields {
		b.WriteString("\t" + f.Name + " string\n")
	}
	b.WriteString("}{\n")
	for _, f := range *modelNode.Fields {
		b.WriteString("\t" + f.Name + ": \"" + f.Name + "\",\n")
	}
	b.WriteString("}\n\n")

	// All Columns
	b.WriteString("var (\n")
	b.WriteString("\t// " + modelNode.Name + "_TableName is the name of the table \n")
	b.WriteString("\t" + modelNode.Name + "_TableName = \"" + modelNode.Name + "\"\n")
	b.WriteString("\t// " + modelNode.Name + "_PrimaryKey is the name of the table's primary key\n")
	b.WriteString("\t" + modelNode.Name + "_PrimaryKey = \"" + primaryKey + "\"\n")
	b.WriteString("\t// " + modelNode.Name + "_AllColumns is a list of all the columns\n")
	b.WriteString("\t" + modelNode.Name + "_AllColumns = []string{")
	for k, f := range *modelNode.Fields {
		b.WriteString("\"" + f.Name + "\"")
		if k < len(*modelNode.Fields)-1 {
			b.WriteByte(',')
		}
	}
	b.WriteString("}")
	b.WriteByte('\n')

	// Insert columns
	b.WriteString("\t// " + modelNode.Name + "_InsertColumns is a list of all insert columns for this model\n")
	b.WriteString("\t" + modelNode.Name + "_InsertColumns = []string{")
	for k := range insertColumns {
		col := insertColumns[k]
		b.WriteString("\"" + col.Name + "\"")
		if k < len(insertColumns)-1 {
			b.WriteByte(',')
		}
	}
	b.WriteString("}")
	b.WriteByte('\n')

	// Update columns
	b.WriteString("\t// " + modelNode.Name + "_UpdateColumns is a list of all update columns for this model\n")
	b.WriteString("\t" + modelNode.Name + "_UpdateColumns = []string{")
	for k := range updateColumns {
		col := updateColumns[k]
		b.WriteString("\"" + col.Name + "\"")
		if k < len(updateColumns)-1 {
			b.WriteByte(',')
		}
	}
	b.WriteString("}\n)")

	file = []byte(b.String())

	file, e = format.Source(file)
	if e != nil {
		log.Fatalf("FORMAT ERROR: File: %s; Error: %s\n%s", modelNode.Name, e.Error(), b.String())
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
