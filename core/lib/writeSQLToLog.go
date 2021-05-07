package lib

import (
	"io/ioutil"
	"os"
	"time"
)

func WriteSQLToLog(sql string) error {

	EnsureDir(MetaDirectory)

	sqlLog := time.Now().Format("20060102150405") + "\n"
	sqlLog += sql

	if _, e := os.Stat(ChangeFilePath); os.IsNotExist(e) {
		ioutil.WriteFile(ChangeFilePath, []byte(sqlLog), 0600)
	} else {
		f, err := os.OpenFile(ChangeFilePath, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err = f.WriteString(sqlLog); err != nil {
			return err
		}
	}

	return nil
}
