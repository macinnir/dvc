package lib

import (
	"fmt"
	"log"
)

// LogLevel is a custom type for passing the level of logging specified by runtime flags
type LogLevel uint

const (

	// LogLevelAlways indicates that a log should always be printed regardless of flags
	LogLevelAlways LogLevel = 0

	// LogLevelError indicates that a log should
	LogLevelError LogLevel = 1

	// LogLevelWarning indicates a potential problem in the system.
	LogLevelWarning LogLevel = 2

	// LogLevelInfo indicates that a log should be shown if the OptVerbose flag is passed
	LogLevelInfo LogLevel = 3

	// LogLevelDebug indicates that a log should be written if the OptVerboseDebug flag is passed
	LogLevelDebug LogLevel = 4
)

// Log writes to a log based on its log level and the runtime options passed
func Log(content string, logLevel LogLevel, options Options) {

	logLevelOption := LogLevelWarning
	logLevelName := "WARN"

	// If silent, only show errors
	if options&OptSilent == OptSilent {
		logLevelOption = LogLevelError
	} else if options&OptLogInfo == OptLogInfo {
		logLevelOption = LogLevelInfo
	} else if options&OptLogDebug == OptLogDebug {
		logLevelOption = LogLevelDebug
	}

	switch logLevel {
	case LogLevelDebug:
		logLevelName = "DEBUG"
	case LogLevelInfo:
		logLevelName = "INFO"
	case LogLevelWarning:
		logLevelName = "WARNING"
	case LogLevelError:
		logLevelName = "ERROR"
	}

	if logLevel <= logLevelOption {
		log.Println(fmt.Sprintf("%s> %s", logLevelName, content))
	}
}

// Error writes an error level log
func Error(content string, options Options) {
	Log(content, LogLevelError, options)
}

// Errorf writes a formatted error level log
func Errorf(content string, options Options, a ...interface{}) {
	Error(fmt.Sprintf(content, a...), options)
}

// Warn writes a warn level log
func Warn(content string, options Options) {
	Log(content, LogLevelWarning, options)
}

// Warnf writes a formatted warn level log
func Warnf(content string, options Options, a ...interface{}) {
	Warn(fmt.Sprintf(content, a...), options)
}

// Info writes an info level log
func Info(content string, options Options) {
	Log(content, LogLevelInfo, options)
}

// Infof writes a formatted info level log
func Infof(content string, options Options, a ...interface{}) {
	Info(fmt.Sprintf(content, a...), options)
}

// Debug writes a debug level log
func Debug(content string, options Options) {
	Log(content, LogLevelDebug, options)
}

// Debugf writes a formatted debug level log
func Debugf(content string, options Options, a ...interface{}) {
	Debug(fmt.Sprintf(content, a...), options)
}
