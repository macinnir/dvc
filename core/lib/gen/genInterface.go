package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/macinnir/dvc/core/lib"
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

func GenInterfaces(srcDir, destDir string) error {

	start := time.Now()

	lib.EnsureDir(destDir)

	var e error
	var files []os.FileInfo

	// DAL
	if files, e = ioutil.ReadDir(srcDir); e != nil {
		return e
	}

	generatedInterfaces := 0

	for _, f := range files {

		// Filter out files that don't have upper case first letter names
		if !unicode.IsUpper([]rune(f.Name())[0]) {
			continue
		}

		srcFile := path.Join(srcDir, f.Name())

		// Remove the go extension
		baseName := f.Name()[0 : len(f.Name())-3]

		// Skip over EXT files
		if baseName[len(baseName)-3:] == "Ext" {
			continue
		}

		// Skip over test files
		if baseName[len(baseName)-5:] == "_test" {
			continue
		}

		// srcFile := path.Join(c.Config.Dirs.Dal, baseName + ".go")
		interfaceName := "I" + baseName
		destFile := path.Join(destDir, interfaceName+".go")
		packageName := filepath.Base(destDir)

		srcFiles := []string{srcFile}
		// var src []byte
		// if src, e = ioutil.ReadFile(srcFile); e != nil {
		// 	return
		// }

		// Check if EXT file exists
		extFile := srcFile[0:len(srcFile)-3] + "Ext.go"
		if _, e = os.Stat(extFile); e == nil {
			srcFiles = append(srcFiles, extFile)
			// concatenate the contents of the ext file with the contents of the regular file
			// var extSrc []byte
			// if extSrc, e = ioutil.ReadFile(extFile); e != nil {
			// 	return
			// }
			// src = append(src, extSrc...)
		}

		var i []byte
		i, e = genInterfaces(
			srcFiles,
			baseName,
			"Generated Code; DO NOT EDIT.",
			packageName,
			interfaceName,
			fmt.Sprintf("%s describes the %s", interfaceName, baseName),
			true,
			true,
		)

		if e != nil {
			return e
		}

		// fmt.Println("Generating ", destFile)
		// fmt.Println("Writing to: ", destFile)

		ioutil.WriteFile(destFile, i, 0644)

		generatedInterfaces++
		// fmt.Println("Name: ", baseName, "Path: ", srcFile)

	}

	fmt.Printf("Generated %d interfaces to %s in %f seconds\n", generatedInterfaces, destDir, time.Since(start).Seconds())

	return nil

}

// GenInterfaces runs GenInterface on all the target files
func genInterfaces(
	files []string,
	structType,
	comment,
	pkgName,
	ifaceName,
	ifaceComment string,
	copyDocuments,
	copyTypeDoc bool,
) (result []byte, e error) {

	// fmt.Printf("Generating interface: %s\n", files[0])
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

	// fmt.Println(allImports)
	result, e = GenInterface(comment, pkgName, ifaceName, ifaceComment, allMethods, allImports)

	return
}
