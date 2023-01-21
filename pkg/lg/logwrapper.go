package lg

import "log"

// LogFunc is a logging function
type LogFunc func(fmt string, args ...any)

// Info logs an info-level log message
var Info LogFunc

// Warning logs a warning-level log message
var Warning LogFunc

// Error logs an error-level log message
var Error LogFunc

// Debug logs a debug-level log message
var Debug LogFunc

var defaultLogger defaultLogWrapper

func init() {
	Debug = defaultLogger.Debug
	Info = defaultLogger.Info
	Warning = defaultLogger.Warning
	Error = defaultLogger.Error
}

type defaultLogWrapper struct {
}

// Info prints an info message to the logger
func (l *defaultLogWrapper) Info(fmt string, args ...any) {
	log.Printf(fmt, args...)
}

// Warning prints a warning message to the logger
func (l *defaultLogWrapper) Warning(fmt string, args ...any) {
	log.Printf(fmt, args...)
}

func (l *defaultLogWrapper) Error(fmt string, args ...any) {
	log.Printf(fmt, args...)
}

func (l *defaultLogWrapper) Debug(fmt string, args ...any) {
	log.Printf(fmt, args...)
}
