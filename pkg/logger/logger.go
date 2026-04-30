package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Level represents the log level
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelNames = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

// Logger provides structured logging with levels
type Logger struct {
	level  Level
	logger *log.Logger
}

// New creates a new logger
func New(level Level) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", 0),
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level <= DEBUG {
		l.log(DEBUG, format, args...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	if l.level <= INFO {
		l.log(INFO, format, args...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level <= WARN {
		l.log(WARN, format, args...)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	if l.level <= ERROR {
		l.log(ERROR, format, args...)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
	os.Exit(1)
}

// log writes a formatted log message
func (l *Logger) log(level Level, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := levelNames[level]
	message := fmt.Sprintf(format, args...)
	l.logger.Printf("[%s] %s: %s", timestamp, levelStr, message)
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// Default logger instance
var defaultLogger = New(INFO)

// Debug logs a debug message using the default logger
func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

// Info logs an info message using the default logger
func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

// Warn logs a warning message using the default logger
func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

// Error logs an error message using the default logger
func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

// Fatal logs a fatal message using the default logger and exits
func Fatal(format string, args ...interface{}) {
	defaultLogger.Fatal(format, args...)
}

// SetLevel sets the logging level for the default logger
func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}
