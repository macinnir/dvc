package lib

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
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

// FetchDirFileNames fetches returns a slice of strings of the name of the files in a directory (non-recursive)
func FetchDirFileNames(dirPath string) (dirs []string, e error) {

	var filesInfo []os.FileInfo
	dirs = []string{}

	if filesInfo, e = ioutil.ReadDir(dirPath); e != nil {
		return
	}

	for _, f := range filesInfo {

		fileName := f.Name()

		if fileName[0:1] == "." || !f.Mode().IsDir() {
			continue
		}

		dirs = append(dirs, f.Name())
	}

	return
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
