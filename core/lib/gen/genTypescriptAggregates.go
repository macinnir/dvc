package gen

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
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
	var files []os.DirEntry
	var removedCount = 0

	if files, e = os.ReadDir(tg.config.TypescriptAggregatesPath); e != nil {
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

	lib.LogRemove(start, "%d typescript Aggregates from `%s`", removedCount, tg.config.TypescriptAggregatesPath)

	return nil
}

// 0.008535
func (tg *TypescriptGenerator) GenerateTypescriptAggregates() error {

	var start = time.Now()
	var generatedCount = 0

	lib.EnsureDir(tg.config.TypescriptAggregatesPath)

	var aggregateNames = []string{}
	for name := range tg.routes.Aggregates {
		aggregateNames = append(aggregateNames, name)
	}
	sort.Strings(aggregateNames)

	var wg sync.WaitGroup
	var mutex = sync.Mutex{}

	for _, name := range aggregateNames {

		wg.Add(1)
		go func(name string) {

			defer wg.Done()

			tsDTOBytes, e := tg.GenerateTypescriptAggregate(name)
			if e != nil {
				log.Fatalf("Error generating typescript aggregate %s: %s", name, e.Error())
			}
			dest := path.Join(tg.config.TypescriptAggregatesPath, name+".ts")
			// log.Printf("Generating Aggregate %s => %s\n", name, dest)
			f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, lib.DefaultFileMode)
			if err != nil {
				log.Fatalf("Error opening file %s: %s", dest, err.Error())
			}
			if _, err = f.Write(tsDTOBytes); err != nil {
				log.Fatalf("Error writing to file %s: %s", dest, err.Error())
			}
			f.Close()

			mutex.Lock()
			generatedCount++
			mutex.Unlock()

		}(name)
	}

	wg.Wait()
	lib.LogAdd(start, "%d typescript Aggregates from %s", generatedCount, tg.config.TypescriptAggregatesPath)
	return nil
}

// GenerateTypescriptAggregate returns a string for a type in typescript
func (tg *TypescriptGenerator) GenerateTypescriptAggregate(name string) ([]byte, error) {

	columns := tg.routes.Aggregates[name]

	var buf bytes.Buffer

	TSFileHeader(&buf, name)

	// fmt.Printf("Aggregate: %s\n", name)

	ImportStrings(&buf, columns, name)

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

	// Remove pointer symbol
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

		goType := columns[name]
		fieldType := schema.GoTypeToTypescriptString(goType)

		if len(name) > 9 && name[0:9] == "#embedded" {
			// tg.GenerateTypescriptFields(sb, fieldType)
			continue
		}

		fmt.Fprint(sb, "\t// "+name+" "+columns[name]+"\n")
		fmt.Fprint(sb, "\t"+name+": "+fieldType+";\n\n")
	}
}

func (tg *TypescriptGenerator) GenerateTypescriptDefaults(sb io.Writer, objectName string) {

	var doLog = objectName == "QuestionAggregate"
	if doLog {
		fmt.Println("GenerateTypescriptDefaults", objectName)
	}

	columns := tg.ExtractColumns(objectName)

	columnNames := ColumnMapToNames(columns)

	for k := range columnNames {
		name := columnNames[k]
		goType := columns[columnNames[k]]
		fieldType := schema.GoTypeToTypescriptString(goType)

		if len(name) > 9 && name[0:9] == "#embedded" {

			if doLog {
				fmt.Println("Embedded", fieldType)
			}
			fmt.Fprint(sb, "\t..."+schema.GoTypeToTypescriptDefault(fieldType)+",\n")
			// tg.GenerateTypescriptDefaults(sb, fieldType)
			continue
		}

		if !unicode.IsUpper(rune(name[0])) {
			continue
		}

		defaultVal := schema.GoTypeToTypescriptDefault(columns[columnNames[k]])
		fmt.Fprint(sb, "\t"+columnNames[k]+": "+defaultVal+",\n")
	}

}
