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
	assert.Equal(t, logFilePath, genLogger.logFile.Name(), test.MessageEqual(logFilePath, genLogger.logFile.Name(), t.Name()))
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
	testCases := []test.TestCase{
		{Name: "test_TRACE", Input: "TRACE", Expected: logLevels["TRACE"]},
		{Name: "test_DEBUG", Input: "DEBUG", Expected: logLevels["DEBUG"]},
		{Name: "test_INFO", Input: "INFO", Expected: logLevels["INFO"]},
		{Name: "test_WARN", Input: "WARN", Expected: logLevels["WARN"]},
		{Name: "test_ERROR", Input: "ERROR", Expected: logLevels["ERROR"]},
	}

	// Sequential execution - subtests share genLogger state so cannot be parallel
	for _, testCase := range testCases {
		genLogger.SetLogLevel(testCase.Input.(string))
		assert.Equal(t, testCase.Expected.(int), genLogger.logLevel, test.MessageEqual(testCase.Expected.(int), genLogger.logLevel, testCase.Name))

		logMessage := fmt.Sprintf("%v message", testCase.Expected.(int))
		switch testCase.Input.(string) {
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
	require.NoError(t, err, test.MessageUnexpectedError(err))

	logFileContent := string(bytes)
	for _, testCase := range testCases {
		logEntry := fmt.Sprintf("%s: %v message", testCase.Input.(string), testCase.Expected.(int))
		assert.Contains(t, logFileContent, logEntry, test.MessageContains(logFileContent, logEntry, testCase.Name))
	}
	assert.NotContains(t, logFileContent, "FATAL message", test.MessageNotContains(logFileContent, "FATAL message", t.Name()))

	genLogger.SetLogLevel("INFO")
}
