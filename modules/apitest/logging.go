package apitest

import (
	"log"
	"strings"

	"github.com/fatih/color"
)

// LogLevel specifies the log levels
type LogLevel int

const (
	// LogLevels
	LogLevelError LogLevel = 0
	LogLevelWarn  LogLevel = 1
	LogLevelInfo  LogLevel = 2
	LogLevelDebug LogLevel = 3
)

// var warnColor
var red func(a ...interface{}) string
var yellow func(a ...interface{}) string
var blue func(a ...interface{}) string
var white func(a ...interface{}) string
var heading func(a ...interface{}) string

type Logger struct {
	logLevel    LogLevel
	indentCount int
}

func (l *Logger) Indent() {
	l.indentCount++
}

func (l *Logger) UnIndent() {
	if l.indentCount > 0 {
		l.indentCount = l.indentCount - 1
	}

	// l.indentCount = 0
	// fmt.Printf("Unindent!!!! %d", l.indentCount)
}

func InitLogger(logLevel LogLevel) *Logger {
	red = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	blue = color.New(color.FgBlue).SprintFunc()
	white = color.New(color.FgWhite).SprintFunc()
	heading = color.New(color.FgBlack, color.Bold, color.BgHiWhite).SprintFunc()

	// Remove all formatting from the global logger
	log.SetFlags(0)

	return &Logger{logLevel, 0}
}

func (l *Logger) Heading(message string) {
	l.Log(heading(message))
	l.Indent()
}

// Log writes a log message to stdout
func (l *Logger) Log(message string) {
	message = strings.Repeat(" ", (l.indentCount*2)) + message
	log.Println(message)
}

func (l *Logger) Error(message string) {
	if l.logLevel >= LogLevelInfo {
		l.Log(red(message))
	}
}

func (l *Logger) Warn(message string) {
	if l.logLevel >= LogLevelInfo {
		l.Log(yellow(message))
	}
}
func (l *Logger) Info(message string) {
	if l.logLevel >= LogLevelInfo {
		l.Log(blue(message))
	}
}

func (l *Logger) Debug(message string) {
	if l.logLevel >= LogLevelDebug {
		l.Log(white(message))
	}
}
