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
func buildModelNodeFromFile(fileBytes []byte) (modelNode *lib.GoStruct, e error) {

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
			modelNode.Comments = s.Comments[0].Text()
			modelNode.Imports = &lib.GoFileImports{}

			for _, i := range s.Imports {
				// This is a named import
				if i.Name != nil {
					modelNode.Imports.Append(i.Name.Name + " " + i.Path.Value)
				} else {
					modelNode.Imports.Append(i.Path.Value)
				}
			}

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
		fieldType := dataTypeToGoTypeString(col)
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

	fileString := "package " + modelNode.Package + "\n\n"
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

	hasNullImport := false
	for _, i := range *modelNode.Imports {
		if i == NullPackage {
			hasNullImport = true
			break
		}
	}

	for n, m := range *modelNode.Fields {
		fieldMap[m.Name] = n
		// Skip any fields not in the database
		if _, ok := table.Columns[m.Name]; !ok {
			continue
		}

		modelFields.Append(m)
	}

	for name, col := range table.Columns {

		_, ok := fieldMap[name]

		// Add any fields not in the model
		if !ok {
			modelFields.Append(&lib.GoStructField{
				Name:     col.Name,
				DataType: dataTypeToGoTypeString(col),
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
		}
	}

	// Finally check for nullables
	for _, f := range *modelFields {

		// If the package needs null, and it hasn't been added, add it
		if fieldDataTypeIsNull(f.DataType) && !hasNullImport {
			modelNode.Imports.Append(NullPackage)
			hasNullImport = true
		}
	}

	modelNode.Fields = modelFields
	return
}

// dataTypeToGoTypeString converts a database type to its equivalent golang datatype
func dataTypeToGoTypeString(column *lib.Column) (fieldType string) {
	fieldType = "int64"
	switch column.DataType {
	case "varchar":
		fieldType = "string"
	case "enum":
		fieldType = "string"
	case "text":
		fieldType = "string"
	case "date":
		fieldType = "string"
	case "datetime":
		fieldType = "string"
	case "char":
		fieldType = "string"
	case "decimal":
		fieldType = "float64"
	}

	if column.IsNullable == true {
		switch fieldType {
		case "string":
			// fieldType = "sql.NullString"
			fieldType = "null.String"
		case "int64":
			// fieldType = "sql.NullInt64"
			fieldType = "null.Int"
		case "float64":
			// fieldType = "sql.NullFloat64"
			fieldType = "null.Float"
		}
	}
	return
}
