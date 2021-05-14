package lib

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"
)

// FetchNonDirFileNames returns a list of files in a directory
// that are only regular files
func FetchNonDirFileNames(dirPath string) (files []string, e error) {

	var filesInfo []os.FileInfo
	files = []string{}

	if filesInfo, e = ioutil.ReadDir(dirPath); e != nil {
		return
	}

	for _, f := range filesInfo {
		fileName := f.Name()
		if fileName[0:1] == "." || f.Mode().IsDir() {
			continue
		}

		files = append(files, f.Name())
	}

	return
}

// FetchDirFileNames fetches returns a slice of strings of the name of the directories in a directory (non-recursive)
func FetchDirFileNames(dirPath string) (fileNames []string, e error) {

	var filesInfo []os.FileInfo
	fileNames = []string{}

	if filesInfo, e = ioutil.ReadDir(dirPath); e != nil {
		return
	}

	for _, f := range filesInfo {

		fileName := f.Name()

		if fileName[0:1] == "." || !f.Mode().IsDir() {
			continue
		}

		fileNames = append(fileNames, f.Name())
	}

	return
}

// HashStringMd5 genereates an MD5 hash of a string
func HashStringMd5(s string) string {
	data := []byte(s)
	return fmt.Sprintf("%x", md5.Sum(data))
}

// HashFileMd5 returns an MD5 checksum of the file at `filePath`
func HashFileMd5(filePath string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil
}

// RunCommand runs a system command
func RunCommand(name string, args ...string) (stdout string, stderr string, exitCode int) {
	// log.Println("run command:", name, args)
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout = outbuf.String()
	stderr = errbuf.String()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			// log.Printf("Could not get exit code for failed program: %v, %v", name, args)
			exitCode = 1
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	// log.Printf("command result, stdout: %v, stderr: %v, exitCode: %v", stdout, stderr, exitCode)
	return
}

// EnsureDir creates a new dir if the dir is not found
func EnsureDir(dir string) (e error) {

	// fmt.Println("Ensuring directory: ", dir)
	// lib.Debugf("Ensuring directory: %s", g.Options, dir)

	if _, e = os.Stat(dir); os.IsNotExist(e) {

		e = os.MkdirAll(dir, 0777)

		if e != nil {
			log.Fatalf("Could not created dir at path: %s", dir)
		}

	}
	return
}

// FmtGoCode formats a go file
func FmtGoCode(filePath string) (e error) {
	_, stdError, exitCode := RunCommand("go", "fmt", filePath)

	if exitCode > 0 {
		e = fmt.Errorf("fmt error: %s", stdError)
	}
	return
}

// WriteGoCodeToFile writes a string of golang code to a file and then formats it with `go fmt`
func WriteGoCodeToFile(goCode string, filePath string) (e error) {
	// outFile := "./repos/repos.go"

	e = ioutil.WriteFile(filePath, []byte(goCode), 0644)
	if e != nil {
		return
	}

	FmtGoCode(filePath)
	// cmd := exec.Command("go", "fmt", filePath)
	// e = cmd.Run()
	// fmt.Printf("WriteCodeToFile: %s\n", e.Error())
	return
}

func DirExists(dirPath string) bool {
	if _, e := os.Stat(dirPath); os.IsNotExist(e) {
		return false
	}

	return true
}

func DirIsEmpty(dirPath string) bool {

	f, e := os.Open(dirPath)
	if e != nil {
		return false
	}

	defer f.Close()

	_, e = f.Readdirnames(1)
	if e == io.EOF {
		return true
	}

	return false
}

func FileExists(filePath string) bool {
	// Check if file exists
	if _, e := os.Stat(filePath); os.IsNotExist(e) {
		// fmt.Printf("File %s does not exist", filePath)
		return false
	}

	return true

}
