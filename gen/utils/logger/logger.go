package logger

import (
	"fmt"
	"log"
	"os"
)

const (
	// Environment variables for logging configuration.
	// The log path can be set to a file path, if none is provided it defaults to "stdout".
	envLogPath = "GEN_LOG_PATH"
	// The log level can be set to TRACE, DEBUG, INFO, WARN, ERROR, or FATAL.
	envLogLevel = "GEN_LOG_LEVEL"
	// The log flags for the formatting of the log messages.
	logFlags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile
	// The log file flags for writing to the log file.
	logFileFlags       = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	logFilePermissions = 0666
)

// The logger variable is a singleton instance of the Logger struct.
// Chose to use the standard log package to avoid adding a dependency on a third-party logging library.
var logger *Logger

// The default log level for the logger.
var logLevel = "INFO"

// The available log levels for the logger.
var logLevels = map[string]int{
	"FATAL": 1,
	"ERROR": 2,
	"WARN":  3,
	"INFO":  4,
	"DEBUG": 5,
	"TRACE": 6,
}

// Logger is a wrapper for the standard log package.
type Logger struct {
	// Log is the standard logger.
	log *log.Logger
	// LogFile is the file where the logs are written to.
	// This is exposed in the struct to allow closing the file from everywhere.
	logFile *os.File
	// LogLevel is the current log level for the logger.
	logLevel int
}

// Initialize a singleton logger instance if one does not exist yet else return existing logger.
func InitalizeLogger() *Logger {
	if logger == nil {
		// Create a new logger instance with default settings.
		logger = &Logger{log: log.New(os.Stdout, "", logFlags)}
	}

	// Check if the log path is set in the environment variables.
	// If it is set, open the file and set it as the output for the logger.
	if logPath := os.Getenv(envLogPath); logPath != "" {
		logger.SetLogFile(logPath)
	}

	// Check if the log level is set in the environment variables.
	// If it is set, set the log level for the logger.
	// If it is not set, use the default log level.
	if envLogLevel := os.Getenv(envLogLevel); envLogLevel != "" {
		logLevel = envLogLevel
	}

	// Set the log level for the logger.
	logger.SetLogLevel(logLevel)

	return logger
}

// Sets the log file for the logger from a path.
func (l *Logger) SetLogFile(logPath string) {
	file, err := os.OpenFile(logPath, logFileFlags, logFilePermissions)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to open log file: %s", err.Error()))
	}
	logger.Info(fmt.Sprintf("Logging to file: %s", logPath))
	logger.log.SetOutput(file)
	logger.logFile = file
}

// Closes the log file if it is open and resets the logger to log to stdout.
func (l *Logger) CloseLogFile() {
	if l.logFile != nil {
		l.logFile.Close()
		l.logFile = nil
		l.log.SetOutput(os.Stdout)
		l.Info("Log file closed, logging to stdout")
	}
}

// Sets the log level for the logger.
// The log level can be set to TRACE, DEBUG, INFO, WARN, ERROR, or FATAL.
func (l *Logger) SetLogLevel(logLevel string) {
	if level, ok := logLevels[logLevel]; ok {
		l.logLevel = level
		logger.Info(fmt.Sprintf("Log level set to: %s", logLevel))
		return
	}

	// If the log level is not valid, log a fatal message and exit the program.
	allowedLevels := []string{}
	for allowedLevel := range logLevels {
		allowedLevels = append(allowedLevels, allowedLevel)
	}
	logger.Fatal(fmt.Sprintf("Invalid log level: %s. Allowed levels are: %s", logLevel, allowedLevels))

}

// Logs a message when the log level is valid to be logged.
func (l *Logger) println(errorLevel, message string) {
	if logLevels[errorLevel] <= l.logLevel {
		l.log.Printf("%s: %s", errorLevel, message)
	}
}

// Fatal logs a fatal message and exits the program.
func (l *Logger) Fatal(message string) {
	l.log.Fatalf("FATAL: %s", message)
}

// Error logs an error message
func (l *Logger) Error(message string) {
	l.println("ERROR", message)
}

// Warn logs a warning message.
func (l *Logger) Warn(message string) {
	l.println("WARN", message)
}

// Info logs an info message.
func (l *Logger) Info(message string) {
	l.println("INFO", message)
}

// Debug logs a debug message.
func (l *Logger) Debug(message string) {
	l.println("DEBUG", message)
}

// Trace logs a trace message.
func (l *Logger) Trace(message string) {
	l.println("TRACE", message)
}
