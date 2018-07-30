package logger

import (
	"fmt"
	"github.com/macinnir/dvc/types"
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
func Log(content string, logLevel LogLevel, options types.Options) {

	logLevelOption := LogLevelWarning
	logLevelName := "WARN"

	// If silent, only show errors
	if options&types.OptSilent == types.OptSilent {
		logLevelOption = LogLevelError
	} else if options&types.OptLogInfo == types.OptLogInfo {
		logLevelOption = LogLevelInfo
	} else if options&types.OptLogDebug == types.OptLogDebug {
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

func Error(content string, options types.Options) {
	Log(content, LogLevelError, options)
}

func Errorf(content string, options types.Options, a ...interface{}) {
	Error(fmt.Sprintf(content, a...), options)
}

func Warn(content string, options types.Options) {
	Log(content, LogLevelWarning, options)
}

func Warnf(content string, options types.Options, a ...interface{}) {
	Warn(fmt.Sprintf(content, a...), options)
}

func Info(content string, options types.Options) {
	Log(content, LogLevelInfo, options)
}

func Infof(content string, options types.Options, a ...interface{}) {
	Info(fmt.Sprintf(content, a...), options)
}

func Debug(content string, options types.Options) {
	Log(content, LogLevelDebug, options)
}

func Debugf(content string, options types.Options, a ...interface{}) {
	Debug(fmt.Sprintf(content, a...), options)
}
