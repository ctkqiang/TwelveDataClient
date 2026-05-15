package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
	"twelve_data_client/internal/color"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var levelNames = map[LogLevel]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	FatalLevel: "FATAL",
}

var levelColors = map[LogLevel]func(string, ...interface{}) string{
	DebugLevel: color.Dimf,
	InfoLevel:  color.Cyanf,
	WarnLevel:  color.Yellowf,
	ErrorLevel: color.Redf,
	FatalLevel: color.Redf,
}

// Logger is a structured logger with levels and timestamps
type Logger struct {
	level        LogLevel
	output       io.Writer
	prefix       string
	useTimestamp bool
	useColor     bool
}

// Config holds logger configuration
type Config struct {
	Level        LogLevel
	Output       io.Writer
	Prefix       string
	UseTimestamp bool
	UseColor     bool
}

// NewLogger creates a new logger with default configuration
func NewLogger(config *Config) *Logger {
	if config == nil {
		config = &Config{}
	}

	if config.Output == nil {
		config.Output = os.Stdout
	}

	if config.Level == 0 {
		config.Level = InfoLevel
	}

	return &Logger{
		level:        config.Level,
		output:       config.Output,
		prefix:       config.Prefix,
		useTimestamp: config.UseTimestamp,
		useColor:     config.UseColor,
	}
}

// NewDefaultLogger creates a logger with sensible defaults
func NewDefaultLogger() *Logger {
	return NewLogger(&Config{
		Level:        InfoLevel,
		Output:       os.Stdout,
		UseTimestamp: true,
		UseColor:     true,
	})
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetPrefix sets the logger prefix
func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

// log is the internal logging function
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	// Build the log message
	var message strings.Builder

	// Add timestamp if enabled
	if l.useTimestamp {
		timestamp := time.Now().Format("15:04:05")
		if l.useColor {
			message.WriteString(color.Dimf("[%s]", timestamp))
		} else {
			message.WriteString(fmt.Sprintf("[%s]", timestamp))
		}
		message.WriteString(" ")
	}

	// Add log level
	levelStr := levelNames[level]
	if l.useColor {
		colorFunc := levelColors[level]
		message.WriteString(colorFunc("[%s]", levelStr))
	} else {
		message.WriteString(fmt.Sprintf("[%s]", levelStr))
	}
	message.WriteString(" ")

	// Add prefix if set
	if l.prefix != "" {
		if l.useColor {
			message.WriteString(color.Blue(l.prefix))
		} else {
			message.WriteString(l.prefix)
		}
		message.WriteString(" ")
	}

	// Add the actual message
	if len(args) > 0 {
		message.WriteString(fmt.Sprintf(format, args...))
	} else {
		message.WriteString(format)
	}

	// Write to output
	logger := log.New(l.output, "", 0)
	logger.Println(message.String())

	// Handle fatal level
	if level == FatalLevel {
		os.Exit(1)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DebugLevel, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(InfoLevel, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WarnLevel, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ErrorLevel, format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FatalLevel, format, args...)
}

// LogError logs an error with context
func (l *Logger) LogError(operation string, err error, context ...string) {
	if err == nil {
		return
	}

	msg := fmt.Sprintf("%s failed: %v", operation, err)
	if len(context) > 0 {
		msg += " (" + strings.Join(context, ", ") + ")"
	}

	l.Error(msg)
}

// LogSuccess logs a successful operation
func (l *Logger) LogSuccess(operation string, details ...string) {
	msg := fmt.Sprintf("%s successful", operation)
	if len(details) > 0 {
		msg += ": " + strings.Join(details, ", ")
	}

	if l.useColor {
		l.Info(color.Green(msg))
	} else {
		l.Info(msg)
	}
}

// LogRequest logs an HTTP request
func (l *Logger) LogRequest(method string, url string) {
	if l.useColor {
		l.Debug("Request: %s %s", color.Cyan(method), color.White(url))
	} else {
		l.Debug("Request: %s %s", method, url)
	}
}

// LogResponse logs an HTTP response
func (l *Logger) LogResponse(statusCode int, message string) {
	statusStr := fmt.Sprintf("%d", statusCode)
	if statusCode >= 200 && statusCode < 300 {
		if l.useColor {
			statusStr = color.Green(statusStr)
		}
	} else if statusCode >= 400 {
		if l.useColor {
			statusStr = color.Red(statusStr)
		}
	}

	l.Debug("Response: %s %s", statusStr, message)
}

// LogDuration logs an operation's duration
func (l *Logger) LogDuration(operation string, duration time.Duration) {
	l.Debug("%s completed in %v", operation, duration)
}

// Global logger instance
var defaultLogger *Logger

func init() {
	defaultLogger = NewDefaultLogger()
}

// Global convenience functions

// SetLevel sets the minimum log level for the default logger
func SetLevel(level LogLevel) {
	defaultLogger.SetLevel(level)
}

// SetPrefix sets the prefix for the default logger
func SetPrefix(prefix string) {
	defaultLogger.SetPrefix(prefix)
}

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

// Fatal logs a fatal message and exits using the default logger
func Fatal(format string, args ...interface{}) {
	defaultLogger.Fatal(format, args...)
}

// LogError logs an error with context using the default logger
func LogError(operation string, err error, context ...string) {
	defaultLogger.LogError(operation, err, context...)
}

// LogSuccess logs a successful operation using the default logger
func LogSuccess(operation string, details ...string) {
	defaultLogger.LogSuccess(operation, details...)
}

// LogRequest logs an HTTP request using the default logger
func LogRequest(method string, url string) {
	defaultLogger.LogRequest(method, url)
}

// LogResponse logs an HTTP response using the default logger
func LogResponse(statusCode int, message string) {
	defaultLogger.LogResponse(statusCode, message)
}

// LogDuration logs an operation's duration using the default logger
func LogDuration(operation string, duration time.Duration) {
	defaultLogger.LogDuration(operation, duration)
}
