package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	// Environment variables for logging configuration.
	// The log path can be set to a file path, if none is provided it defaults to "stdout".
	constEnvLogPath = "GEN_ACI_TF_LOG_PATH"
	// The log level can be set to TRACE, DEBUG, INFO, WARN, ERROR, or FATAL.
	constEnvLogLevel = "GEN_ACI_TF_LOG_LEVEL"
	// The log flags for the formatting of the log messages.
	constLogFlags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile
	// The log file flags for writing to the log file.
	// The log file is created if it does not exist, and it is opened in append mode.
	constLogFileFlags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	// The log file permissions for the log file.
	// This is set to 0666 to allow all users to read and write to the log file.
	constLogFilePermissions = 0666
	// The default log level used when neither the environment variable nor SetLogLevel is called.
	constDefaultLogLevel = "INFO"
)

// The available log levels for the logger.
var logLevels = map[string]int{
	"FATAL": 1,
	"ERROR": 2,
	"WARN":  3,
	"INFO":  4,
	"DEBUG": 5,
	"TRACE": 6,
}

// The logger variable is a singleton instance of the Logger struct.
// Chose to use the standard log package to avoid adding a dependency on a third-party logging library.
var logger *Logger

// Logger is a wrapper for the standard log package.
type Logger struct {
	// Log is the standard logger.
	log *log.Logger
	// LogFile is the file where the logs are written to.
	// This is exposed in the struct to allow closing the file from everywhere.
	logFile *os.File
	// LogLevel is the current log level for the logger.
	logLevel int
	// exitFunc is called by Fatal/Fatalf. Defaults to os.Exit; overridden in tests.
	exitFunc func(int)
}

// Initialize a singleton logger instance if one does not exist yet else return existing logger.
func InitializeLogger() *Logger {
	if logger == nil {
		// Create a new logger instance with default settings.
		logger = &Logger{
			log:      log.New(os.Stdout, "", constLogFlags),
			exitFunc: os.Exit,
		}
	}

	// Check if the log path is set in the environment variables.
	// If it is set, open the file and set it as the output for the logger.
	if envLogPath := os.Getenv(constEnvLogPath); envLogPath != "" {
		logger.SetLogFile(envLogPath)
	}

	// Check if the log level is set in the environment variables.
	// If it is set, set the log level for the logger.
	// If it is not set, use the default log level.
	level := constDefaultLogLevel
	if envLogLevel := os.Getenv(constEnvLogLevel); envLogLevel != "" {
		level = envLogLevel
	}

	// Set the log level for the logger.
	logger.SetLogLevel(level)

	return logger
}

// ResetForTest clears the singleton (closing any open log file) so the next
// InitializeLogger call constructs a fresh instance. Intended for tests only.
func ResetForTest() {
	if logger != nil && logger.logFile != nil {
		_ = logger.logFile.Close()
	}
	logger = nil
}

// SetExitFunc overrides the function called by Fatal/Fatalf. Intended for tests.
func (l *Logger) SetExitFunc(f func(int)) {
	l.exitFunc = f
}

// Sets the log file for the logger from a path.
func (l *Logger) SetLogFile(logPath string) {
	file, err := os.OpenFile(logPath, constLogFileFlags, constLogFilePermissions)
	if err != nil {
		l.Fatalf("Failed to open log file: %s", err.Error())
		return
	}
	l.Tracef("Logging to file: %s", logPath)
	l.log.SetOutput(file)
	l.logFile = file
}

// Closes the log file if it is open and resets the logger to log to stdout.
func (l *Logger) CloseLogFile() {
	if l.logFile != nil {
		l.logFile.Close()
		l.logFile = nil
		l.log.SetOutput(os.Stdout)
		l.Trace("Log file closed, logging to stdout")
	}
}

// SetOutputForTesting sets the output writer for the logger.
// This is intended for use in tests to capture and validate log output.
func (l *Logger) SetOutputForTesting(w io.Writer) {
	l.log.SetOutput(w)
}

// Sets the log level for the logger.
// The log level can be set to TRACE, DEBUG, INFO, WARN, ERROR, or FATAL.
func (l *Logger) SetLogLevel(logLevel string) {
	if level, ok := logLevels[logLevel]; ok {
		l.logLevel = level
		l.Tracef("Log level set to: %s", logLevel)
		return
	}

	// If the log level is not valid, log a fatal message and exit the program.
	allowedLevels := []string{}
	for allowedLevel := range logLevels {
		allowedLevels = append(allowedLevels, allowedLevel)
	}
	l.Fatalf("Invalid log level: %s. Allowed levels are: %s", logLevel, allowedLevels)
}

// The calldepth passed to log.Output so that Llongfile/Lshortfile resolve to
// the original caller (the user's code), not a frame inside this package.
//
// Stack at the moment runtime.Caller is invoked from inside log.Output:
//
//	0: log.Output
//	1: (*Logger).println / (*Logger).printf / (*Logger).fatal
//	2: (*Logger).Debug / Debugf / Info / ... / Fatal / Fatalf
//	3: caller of the public method   ← what we want printed
const constLogCallDepth = 3

// Logs a plain message when the log level is valid to be logged.
// Used by Trace/Debug/Info/Warn/Error.
func (l *Logger) println(errorLevel, message string) {
	if logLevels[errorLevel] <= l.logLevel {
		// calldepth=constLogCallDepth so Llongfile points at the caller of the public level method.
		l.log.Output(constLogCallDepth, errorLevel+": "+message)
	}
}

// Logs a formatted message when the log level is valid to be logged.
// Used by Tracef/Debugf/Infof/Warnf/Errorf. Arguments are handled in the manner of fmt.Printf.
func (l *Logger) printf(errorLevel, format string, v ...any) {
	if logLevels[errorLevel] <= l.logLevel {
		// calldepth=constLogCallDepth so Llongfile points at the caller of the public level method.
		l.log.Output(constLogCallDepth, errorLevel+": "+fmt.Sprintf(format, v...))
	}
}

// Logs a fatal message and exits the program via exitFunc. Shared by Fatal and
// Fatalf so both share the same calldepth as the other levels.
func (l *Logger) fatal(message string) {
	// calldepth=constLogCallDepth so Llongfile points at the caller of Fatal/Fatalf.
	l.log.Output(constLogCallDepth, "FATAL: "+message)
	l.exitFunc(1)
}

// Fatal logs a fatal message and exits the program.
func (l *Logger) Fatal(message string) {
	l.fatal(message)
}

// Fatalf logs a formatted fatal message and exits the program.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Fatalf(format string, v ...any) {
	l.fatal(fmt.Sprintf(format, v...))
}

// Error logs an error message.
func (l *Logger) Error(message string) {
	l.println("ERROR", message)
}

// Errorf logs a formatted error message.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Errorf(format string, v ...any) {
	l.printf("ERROR", format, v...)
}

// Warn logs a warning message.
func (l *Logger) Warn(message string) {
	l.println("WARN", message)
}

// Warnf logs a formatted warning message.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Warnf(format string, v ...any) {
	l.printf("WARN", format, v...)
}

// Info logs an info message.
func (l *Logger) Info(message string) {
	l.println("INFO", message)
}

// Infof logs a formatted info message.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Infof(format string, v ...any) {
	l.printf("INFO", format, v...)
}

// Debug logs a debug message.
func (l *Logger) Debug(message string) {
	l.println("DEBUG", message)
}

// Debugf logs a formatted debug message.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Debugf(format string, v ...any) {
	l.printf("DEBUG", format, v...)
}

// Trace logs a trace message.
func (l *Logger) Trace(message string) {
	l.println("TRACE", message)
}

// Tracef logs a formatted trace message.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Tracef(format string, v ...any) {
	l.printf("TRACE", format, v...)
}
