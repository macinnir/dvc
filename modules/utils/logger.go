package utils

import (
	"log"
	"os"
)

func NewLogger(appName string, logFilePath string) *Logger {

	f, e := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if e != nil {
		log.Fatalf("error opening log file %s", logFilePath)
	}

	return &Logger{
		file:    f,
		appName: appName,
	}
}

type Logger struct {
	appName string
	file    *os.File
}

func (l *Logger) Write(p []byte) (n int, e error) {

	s := l.appName + "> " + string(p)
	n, e = l.file.Write([]byte(s))
	return
}

func (l *Logger) Finish() {
	l.file.Close()
}
