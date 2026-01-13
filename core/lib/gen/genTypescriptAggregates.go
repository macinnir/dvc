package gen

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
	"unicode"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

type TypescriptGenerator struct {
	routes *lib.RoutesJSONContainer
	config *lib.Config
}

func NewTypescriptGenerator(config *lib.Config, routes *lib.RoutesJSONContainer) *TypescriptGenerator {
	return &TypescriptGenerator{
		routes,
		config,
	}
}

func (tg *TypescriptGenerator) CleanTypescriptAggregates() error {

	var start = time.Now()
	var e error
	var files []os.FileInfo
	var removedCount = 0

	if files, e = ioutil.ReadDir(tg.config.TypescriptAggregatesPath); e != nil {
		return e
	}

	for k := range files {

		if files[k].IsDir() {
			continue
		}

		if files[k].Name()[0:1] == "." {
			continue
		}

		aggName := files[k].Name()[0 : len(files[k].Name())-3]

		if _, ok := tg.routes.Aggregates[aggName]; !ok {
			tsFilePath := path.Join(tg.config.TypescriptAggregatesPath, files[k].Name())
			os.Remove(tsFilePath)
			removedCount++
		}

	}

	fmt.Printf("Removed %d typescript Aggregates from `%s` in %f seconds\n", removedCount, tg.config.TypescriptAggregatesPath, time.Since(start).Seconds())

	return nil
}

// 0.008535
func (tg *TypescriptGenerator) GenerateTypescriptAggregates() error {

	var start = time.Now()
	var generatedCount = 0

	lib.EnsureDir(tg.config.TypescriptAggregatesPath)

	for name := range tg.routes.Aggregates {
		fmt.Println("Generating Typescript Aggregate:", name)
		tsDTOBytes, e := tg.GenerateTypescriptAggregate(name)
		if e != nil {
			fmt.Println("ERROR:", e.Error())
			return e
		}
		dest := path.Join(tg.config.TypescriptAggregatesPath, name+".ts")
		err := os.WriteFile(dest, tsDTOBytes, lib.DefaultFileMode)
		// f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, lib.DefaultFileMode)
		if err != nil {
			panic(err)
		}
		// if _, err = f.Write(tsDTOBytes); err != nil {
		// 	panic(err)
		// }
		// f.Close()
		// ioutil.WriteFile(dest, []byte(str), 0777)
		generatedCount++
	}

	fmt.Printf("Generated %d typescript Aggregates from `%s` in %f seconds\n", generatedCount, tg.config.TypescriptAggregatesPath, time.Since(start).Seconds())

	return nil
}

// GenerateTypescriptAggregate returns a string for a type in typescript
func (tg *TypescriptGenerator) GenerateTypescriptAggregate(name string) ([]byte, error) {

	columns := tg.routes.Aggregates[name]
	var buf bytes.Buffer

	TSFileHeader(&buf, name)
	ImportStrings(&buf, columns)

	buf.WriteString(`
export type ` + name + ` = `)

	inherits := InheritStrings(&buf, columns)

	if len(inherits) > 0 {
		buf.WriteString(strings.Join(inherits, " & ") + " & ")
	}

	buf.WriteString(`{

`)

	tg.GenerateTypescriptFields(&buf, name)

	buf.WriteString("}\n\n")

	buf.WriteString("// new" + name + " is a factory method for new " + name + " objects\n")
	buf.WriteString("export const new" + name + " = () : " + name + " => ({\n")

	tg.GenerateTypescriptDefaults(&buf, name)

	buf.WriteString("});\n\n")

	return buf.Bytes(), nil
}

func (tg *TypescriptGenerator) ExtractColumns(goType string) map[string]string {

	if goType[0:1] == "*" {
		goType = goType[1:]
	}

	// fmt.Println("GoType", goType)
	// if goType == "CustomerAggregate" {
	// 	fmt.Println(goType, " --> ", goType[len(goType)-9:])
	// }

	if len(goType) > 3 && goType[len(goType)-3:] == "DTO" {
		return tg.routes.DTOs[goType]
	}

	if len(goType) > 9 && goType[len(goType)-9:] == "Aggregate" {
		cols := tg.routes.Aggregates[goType]
		// fmt.Println("Got here", goType, len(cols))
		// if !ok {
		// 	fmt.Println("Nope")
		// }
		return cols

	}

	return tg.routes.Models[goType]
}

func (tg *TypescriptGenerator) GenerateTypescriptFields(sb io.Writer, objectName string) {

	columns := tg.ExtractColumns(objectName)

	columnNames := ColumnMapToNames(columns)

	for k := range columnNames {

		name := columnNames[k]

		// Uppercase fields only
		if !unicode.IsUpper(rune(name[0])) {
			continue
		}

		goType := columns[columnNames[k]]
		fieldType := schema.GoTypeToTypescriptString(goType)

		if len(name) > 9 && name[0:9] == "#embedded" {
			// tg.GenerateTypescriptFields(sb, fieldType)
			continue
		}

		fmt.Fprintf(sb, "\t// %s %s\n", columnNames[k], columns[columnNames[k]])
		fmt.Fprintf(sb, "\t%s: %s;\n\n", columnNames[k], fieldType)
	}
}

func (tg *TypescriptGenerator) GenerateTypescriptDefaults(sb io.Writer, objectName string) {

	columns := tg.ExtractColumns(objectName)
	columnNames := ColumnMapToNames(columns)

	for k := range columnNames {
		name := columnNames[k]
		goType := columns[columnNames[k]]
		fieldType := schema.GoTypeToTypescriptString(goType)

		if len(name) > 9 && name[0:9] == "#embedded" {

			typescriptDefault := schema.GoTypeToTypescriptDefault(fieldType)
			// if objectName == "QuestionAggregate" {
			// 	fmt.Println("\tEmbedded default for ", objectName, " field ", name, " type ", fieldType, " default ", typescriptDefault)
			// }

			fmt.Fprintf(sb, "\t..."+typescriptDefault+",\n")

			// tg.GenerateTypescriptDefaults(sb, fieldType)
			continue
		}

		// Public properties only
		if !unicode.IsUpper(rune(name[0])) {
			continue
		}

		defaultVal := schema.GoTypeToTypescriptDefault(columns[columnNames[k]])
		fmt.Fprintf(sb, "\t"+columnNames[k]+": "+defaultVal+",\n")
	}

}
