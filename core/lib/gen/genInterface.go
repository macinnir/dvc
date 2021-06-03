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

func fetchSrcFilesForInterfaces(srcDir string) ([]string, error) {

	// fmt.Println("Fetching files for dir ", srcDir)
	filePaths := []string{}
	var files []os.FileInfo
	var e error

	// DAL
	if files, e = ioutil.ReadDir(srcDir); e != nil {
		return filePaths, e
	}

	for k := range files {

		f := files[k]

		filePath := path.Join(srcDir, files[k].Name())

		if files[k].IsDir() {
			var subFiles []string
			if subFiles, e = fetchSrcFilesForInterfaces(filePath); e != nil {
				return filePaths, e
			}
			filePaths = append(filePaths, subFiles...)
			continue
		}

		// Filter out files that don't have upper case first letter names
		if !unicode.IsUpper([]rune(f.Name())[0]) {
			continue
		}

		// Verify this is a go file
		if f.Name()[len(f.Name())-3:] != ".go" {
			continue
		}

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

		filePaths = append(filePaths, filePath)
	}

	return filePaths, nil

}

func GenInterfaces(srcDir, destDir string) error {

	start := time.Now()

	lib.EnsureDir(destDir)

	var e error
	generatedInterfaces := 0
	var files []string
	if files, e = fetchSrcFilesForInterfaces(srcDir); e != nil {
		return e
	}

	for k := range files {

		srcFile := files[k]
		baseName := filepath.Base(srcFile)
		structName := baseName[0 : len(baseName)-3]
		interfaceName := "I" + structName
		packageName := filepath.Base(filepath.Dir(srcFile))
		subDestDir := path.Join(destDir, packageName)
		destFile := path.Join(subDestDir, interfaceName+".go")
		srcFiles := []string{srcFile}

		// fmt.Println("Generating interface: " + srcFile + " ==> " + destFile)

		// Check if EXT file exists
		extFile := srcFile[0:len(srcFile)-3] + "Ext.go"
		if _, e = os.Stat(extFile); e == nil {
			srcFiles = append(srcFiles, extFile)
		}

		var i []byte
		i, e = genInterfaces(
			srcFiles,
			structName,
			"Generated Code; DO NOT EDIT.",
			packageName,
			interfaceName,
			fmt.Sprintf("%s describes the %s", interfaceName, baseName),
			true,
			true,
		)

		if e != nil {
			fmt.Println("ERROR: " + srcFile + " => " + e.Error())
			return e
		}

		lib.EnsureDir(subDestDir)

		ioutil.WriteFile(destFile, i, 0644)

		generatedInterfaces++

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

	result, e = GenInterface(comment, pkgName, ifaceName, ifaceComment, allMethods, allImports)

	return
}
