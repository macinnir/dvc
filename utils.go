package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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

func fatal(msg string) {
	fmt.Printf("ERROR: %s\n", msg)
	os.Exit(1)
}
