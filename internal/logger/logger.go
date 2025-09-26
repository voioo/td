// Package logger provides structured logging functionality for the td application.
package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Level represents the logging level.
type Level int

const (
	// LevelDebug is for detailed debug information.
	LevelDebug Level = iota
	// LevelInfo is for general information.
	LevelInfo
	// LevelWarn is for warning messages.
	LevelWarn
	// LevelError is for error messages.
	LevelError
	// LevelFatal is for fatal errors that cause program exit.
	LevelFatal
)

// String returns the string representation of the level.
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging functionality.
type Logger struct {
	level  Level
	writer io.Writer
	json   bool
}

// NewLogger creates a new logger with the specified level and output destination.
func NewLogger(level Level, writer io.Writer, json bool) *Logger {
	return &Logger{
		level:  level,
		writer: writer,
		json:   json,
	}
}

// NewDefaultLogger creates a logger with default settings (INFO level, stderr, non-JSON).
func NewDefaultLogger() *Logger {
	return NewLogger(LevelInfo, os.Stderr, false)
}

// SetLevel sets the logging level.
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// SetJSON sets whether to output in JSON format.
func (l *Logger) SetJSON(json bool) {
	l.json = json
}

// log writes a log entry if the level is enabled.
func (l *Logger) log(level Level, message string, fields map[string]interface{}) {
	if level < l.level {
		return
	}

	entry := LogEntry{
		Time:    time.Now().UTC(),
		Level:   level.String(),
		Message: message,
		Fields:  fields,
	}

	var output string
	if l.json {
		if jsonBytes, err := json.Marshal(entry); err == nil {
			output = string(jsonBytes)
		} else {
			output = fmt.Sprintf("{\"error\":\"failed to marshal log entry\",\"message\":\"%s\"}", message)
		}
	} else {
		output = l.formatText(entry)
	}

	fmt.Fprintln(l.writer, output)

	if level == LevelFatal {
		os.Exit(1)
	}
}

// formatText formats a log entry as human-readable text.
func (l *Logger) formatText(entry LogEntry) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s [%s] %s",
		entry.Time.Format("2006-01-02T15:04:05Z"),
		entry.Level,
		entry.Message))

	if len(entry.Fields) > 0 {
		sb.WriteString(" |")
		for k, v := range entry.Fields {
			sb.WriteString(fmt.Sprintf(" %s=%v", k, v))
		}
	}

	return sb.String()
}

// Debug logs a debug message.
func (l *Logger) Debug(message string, fields ...Field) {
	l.log(LevelDebug, message, fieldsToMap(fields))
}

// Info logs an info message.
func (l *Logger) Info(message string, fields ...Field) {
	l.log(LevelInfo, message, fieldsToMap(fields))
}

// Warn logs a warning message.
func (l *Logger) Warn(message string, fields ...Field) {
	l.log(LevelWarn, message, fieldsToMap(fields))
}

// Error logs an error message.
func (l *Logger) Error(message string, fields ...Field) {
	l.log(LevelError, message, fieldsToMap(fields))
}

// Fatal logs a fatal message and exits.
func (l *Logger) Fatal(message string, fields ...Field) {
	l.log(LevelFatal, message, fieldsToMap(fields))
}

// Field represents a structured logging field.
type Field struct {
	Key   string
	Value interface{}
}

// F creates a new field.
func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// fieldsToMap converts fields to a map.
func fieldsToMap(fields []Field) map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}

	m := make(map[string]interface{}, len(fields))
	for _, f := range fields {
		m[f.Key] = f.Value
	}
	return m
}

// LogEntry represents a structured log entry.
type LogEntry struct {
	Time    time.Time              `json:"time"`
	Level   string                 `json:"level"`
	Message string                 `json:"message"`
	Fields  map[string]interface{} `json:"fields,omitempty"`
}

// Global logger instance
var defaultLogger = NewDefaultLogger()

// SetDefaultLogger sets the default logger.
func SetDefaultLogger(l *Logger) {
	defaultLogger = l
}

// Debug logs a debug message using the default logger.
func Debug(message string, fields ...Field) {
	defaultLogger.Debug(message, fields...)
}

// Info logs an info message using the default logger.
func Info(message string, fields ...Field) {
	defaultLogger.Info(message, fields...)
}

// Warn logs a warning message using the default logger.
func Warn(message string, fields ...Field) {
	defaultLogger.Warn(message, fields...)
}

// Error logs an error message using the default logger.
func Error(message string, fields ...Field) {
	defaultLogger.Error(message, fields...)
}

// Fatal logs a fatal message using the default logger.
func Fatal(message string, fields ...Field) {
	defaultLogger.Fatal(message, fields...)
}
