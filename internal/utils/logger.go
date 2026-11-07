package utils

import (
	"log"
	"os"
	"sync"
)

// Simple logger wrapper. Uses the standard library logger but provides
// convenient level functions. The implementation is deliberately lightweight
// – the project does not yet depend on a structured logger.

var (
	logger   *log.Logger
	loggerMu sync.Mutex
)

func init() {
	// Default logger writes to stdout with standard flags.
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmsgprefix)
}

// SetOutput allows tests to redirect logger output.
func SetOutput(out *os.File) {
	loggerMu.Lock()
	defer loggerMu.Unlock()
	logger.SetOutput(out)
}

// Info logs an informational message.
func Info(v ...interface{}) {
	loggerMu.Lock()
	defer loggerMu.Unlock()
	logger.SetPrefix("INFO: ")
	logger.Print(v...)
}

// Debug logs a debug message.
func Debug(v ...interface{}) {
	loggerMu.Lock()
	defer loggerMu.Unlock()
	logger.SetPrefix("DEBUG: ")
	logger.Print(v...)
}

// Error logs an error message.
func Error(v ...interface{}) {
	loggerMu.Lock()
	defer loggerMu.Unlock()
	logger.SetPrefix("ERROR: ")
	logger.Print(v...)
}
