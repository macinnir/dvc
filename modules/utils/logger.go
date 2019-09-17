package utils

import (
	"log"
	"os"
)

// NewLogger returns a new Logger object
func NewLogger(logFilePath string) *Logger {

	f, e := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	// log.SetFlags()

	if e != nil {
		log.Fatalf("error opening log file %s", logFilePath)
	}

	logger := &Logger{
		file: f,
	}

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(logger)

	return logger
}

// Logger is a logging utility
type Logger struct {
	file *os.File
	n    int64
}

func (l *Logger) Write(p []byte) (n int, e error) {
	n, e = l.file.Write(p)
	return
}

// Finish closes the Logger
func (l *Logger) Finish() {
	l.file.Close()
}
