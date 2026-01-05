package logger

import (
	"fmt"
	"os"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	logFilePath = "./loggertest.log"
)

func initializeLogTest(t *testing.T) *Logger {
	t.Helper()
	test.InitializeTest(t)
	genLogger := InitializeLogger()
	require.NotNil(t, genLogger, "logger must be initialized")
	return genLogger
}

func TestInitializeLogger(t *testing.T) {
	t.Parallel()
	genLogger := initializeLogTest(t)

	assert.Nil(t, genLogger.logFile)
	assert.NotNil(t, genLogger.logLevel)
}

func TestLogFile(t *testing.T) {
	t.Setenv(constEnvLogPath, logFilePath)
	t.Cleanup(func() {
		os.Remove(logFilePath)
	})

	genLogger := initializeLogTest(t)

	require.NotNil(t, genLogger.logFile)
	assert.Equal(t, logFilePath, genLogger.logFile.Name())
}

func TestSetLogLevel(t *testing.T) {
	// Note: This test cannot use t.Parallel() because it shares state via genLogger and log file
	t.Setenv(constEnvLogPath, logFilePath)
	t.Setenv(constEnvLogLevel, "TRACE")
	t.Cleanup(func() {
		os.Remove(logFilePath)
	})

	genLogger := initializeLogTest(t)
	t.Cleanup(func() {
		genLogger.CloseLogFile()
	})

	// Define ordered log levels to ensure deterministic execution
	// Using a slice instead of ranging over map to ensure consistent order
	orderedLogLevels := []struct {
		name  string
		level int
	}{
		{"TRACE", logLevels["TRACE"]},
		{"DEBUG", logLevels["DEBUG"]},
		{"INFO", logLevels["INFO"]},
		{"WARN", logLevels["WARN"]},
		{"ERROR", logLevels["ERROR"]},
	}

	// Sequential execution - subtests share genLogger state so cannot be parallel
	for _, ll := range orderedLogLevels {
		genLogger.SetLogLevel(ll.name)
		assert.Equal(t, ll.level, genLogger.logLevel, "log level mismatch for %s", ll.name)

		logMessage := fmt.Sprintf("%v message", ll.level)
		switch ll.name {
		case "TRACE":
			genLogger.Trace(logMessage)
		case "DEBUG":
			genLogger.Debug(logMessage)
		case "INFO":
			genLogger.Info(logMessage)
		case "WARN":
			genLogger.Warn(logMessage)
		case "ERROR":
			genLogger.Error(logMessage)
		}
	}

	bytes, err := os.ReadFile(logFilePath)
	require.NoError(t, err, "failed to read log file")

	logFileContent := string(bytes)
	for _, ll := range orderedLogLevels {
		logEntry := fmt.Sprintf("%s: %v message", ll.name, ll.level)
		assert.Contains(t, logFileContent, logEntry)
	}
	assert.NotContains(t, logFileContent, "FATAL message")

	genLogger.SetLogLevel("INFO")
}
