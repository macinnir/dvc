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

	var sb strings.Builder
	sb.WriteString(`// ` + comment + `

package ` + pkgName + `

import (
`)
	for k := range imports {
		sb.WriteString("\t" + imports[k] + "\n")
	}

	sb.WriteString(`)

`)

	// output := []string{"// " + comment, "", "package " + pkgName, "import ("}
	// output = append(output, imports...)
	// output = append(output, ")", "")

	if len(ifaceComment) > 0 {
		sb.WriteString(fmt.Sprintf("// %s", strings.Replace(ifaceComment, "\n", "\n// ", -1)))
	}
	sb.WriteString("\ntype " + ifaceName + " interface {\n")
	for k := range methods {
		sb.WriteString("\t" + methods[k] + "\n")
	}

	sb.WriteString("}\n")
	return lib.FormatCode(sb.String())
	// return []byte(sb.String()), nil

	// output = append(output, fmt.Sprintf("type %s interface {", ifaceName))
	// output = append(output, methods...)
	// output = append(output, "}")
	// code := strings.Join(output, "\n")
	// return lib.FormatCode(code)
}

func fetchExistingInterfaceFiles(srcDir string) ([]string, error) {

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
			if subFiles, e = fetchExistingInterfaceFiles(filePath); e != nil {
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

		filePaths = append(filePaths, filePath)
	}

	return filePaths, nil

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
		if len(baseName) > 5 && baseName[len(baseName)-5:] == "_test" {
			continue
		}

		filePaths = append(filePaths, filePath)
	}

	return filePaths, nil

}

func GenInterfaces(srcDir, destDir string) error {

	fmt.Println("Generating interface files from", srcDir, " => ", destDir)

	lib.EnsureDir(destDir)

	var start = time.Now()
	var e error
	var generatedInterfaces = 0
	var files []string
	if files, e = fetchSrcFilesForInterfaces(srcDir); e != nil {
		return e
	}

	// var interfaceMap = map[string]struct{}{}
	// for k := range files {
	// 	fmt.Println("Source File:", filepath.Dir(files[k]), filepath.Base(files[k]))
	// 	// var baseFile = filepath.Base(files[k])
	// 	// interfaceMap[path.Join(destDir, "I"+baseFile[0:len(baseFile)-3]+".go")] = struct{}{}
	// }

	// var existingInterfaceFiles []string
	// if existingInterfaceFiles, e = fetchSrcFilesForInterfaces(destDir); e != nil {
	// 	return e
	// }

	// for k := range existingInterfaceFiles {
	// 	fmt.Println("Existing Interface:", existingInterfaceFiles[k])
	// }

	var genInterfaceMap = map[string]struct{}{}
	for k := range files {
		var srcFile = files[k]
		baseName := filepath.Base(srcFile)
		var structName = baseName[0 : len(baseName)-3]

		interfaceName := "I" + structName
		packageName := filepath.Base(filepath.Dir(srcFile))
		destDirName := filepath.Base(destDir)
		destSubDir := destDir

		if packageName != destDirName {
			destSubDir = path.Join(destDir, packageName)
		}

		destFile := path.Join(destSubDir, interfaceName+".go")
		genInterfaceMap[destFile] = struct{}{}
		// fmt.Println("Generating", destFile)
		GenInterface2(structName, srcFile, packageName, interfaceName, destSubDir, destFile)
		generatedInterfaces++
	}

	// for k := range existingInterfaceFiles {
	// 	if _, ok := genInterfaceMap[existingInterfaceFiles[k]]; !ok {
	// 		fmt.Println("Removing interface at path", existingInterfaceFiles[k])
	// 	}
	// }

	fmt.Printf("Generated %d interfaces to %s in %f seconds\n", generatedInterfaces, destDir, time.Since(start).Seconds())

	return nil

}

func GenInterface2(structName, srcFile, packageName, interfaceName, destDir, destFile string) error {

	var e error

	var i []byte

	var src []byte
	if src, e = ioutil.ReadFile(srcFile); e != nil {
		return e
	}

	i, e = genInterface(
		src,
		structName,
		"Generated Code; DO NOT EDIT.",
		packageName,
		interfaceName,
		fmt.Sprintf("%s describes %s", interfaceName, structName),
		true,
		true,
	)

	if e != nil {
		fmt.Println("ERROR: " + srcFile + " => " + e.Error())
		return e
	}

	lib.EnsureDir(destDir)

	// TODO verbose flag
	// fmt.Printf("Generating interface %s...\n", destFile)
	ioutil.WriteFile(destFile, i, 0644)

	return nil
}

// genInterface generates an interface for the given file
func genInterface(
	src []byte,
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

	if typeDoc != "" {
		ifaceComment += "\n" + typeDoc
	}

	result, e = GenInterface(comment, pkgName, ifaceName, ifaceComment, allMethods, allImports)

	return
}
