package logger

import (
	"fmt"
	"os"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

const (
	logFilePath = "./loggertest.log"
)

func initializeLogTest(t *testing.T) *Logger {
	test.InitializeTest(t, 1)
	genLogger := InitalizeLogger()

	if genLogger == nil {
		t.Errorf("Expected logger instance to be initialized, but got nil")
	}

	return genLogger
}

func cleanLogTest() {
	os.Remove(logFilePath)
}

func TestInitalizeLogger(t *testing.T) {
	genLogger := initializeLogTest(t)

	assert.Nil(t, genLogger.logFile, fmt.Sprintf("Expected log file to be nil, but got %v", genLogger.logFile))
	assert.NotEqual(t, genLogger.logLevel, nil, fmt.Sprintf("Expected log level to be %d, but got %d", logLevels["INFO"], genLogger.logLevel))

	cleanLogTest()
}

func TestLogFile(t *testing.T) {
	t.Setenv(envLogPath, logFilePath)
	genLogger := initializeLogTest(t)

	assert.Equal(t, genLogger.logFile.Name(), logFilePath, fmt.Sprintf("Expected log file to be %s, but got %v", logFilePath, genLogger.logFile.Name()))
	assert.NotEqualf(t, genLogger.logFile, nil, "Expected log file to be nil after closing, but got %v", genLogger.logFile)

	cleanLogTest()
}

func TestSetLogLevel(t *testing.T) {
	t.Setenv(envLogPath, logFilePath)
	t.Setenv(envLogLevel, "TRACE")
	genLogger := initializeLogTest(t)
	defer genLogger.CloseLogFile()

	for logLevel, logLevelInt := range logLevels {
		genLogger.SetLogLevel(logLevel)

		assert.Equal(t, genLogger.logLevel, logLevelInt, fmt.Sprintf("Expected log level to be %d, but got %d", logLevelInt, genLogger.logLevel))

		logMessage := fmt.Sprintf("%v message", logLevelInt)
		switch logLevel {
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
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logFileContent := string(bytes)
	for logLevel, logLevelInt := range logLevels {
		if logLevel != "FATAL" {
			logEntry := fmt.Sprintf("%s: %v message", logLevel, logLevelInt)
			assert.Contains(t, logFileContent, logEntry, fmt.Sprintf("Expected log file to contain '%s'", logEntry))
		} else {
			assert.NotContains(t, logFileContent, "FATAL message", "Expected log file to not contain 'FATAL message'")
		}
	}

	cleanLogTest()
}
