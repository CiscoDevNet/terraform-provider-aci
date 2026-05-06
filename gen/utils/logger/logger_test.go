package logger

import (
	"bytes"
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
	test.InitializeTest(t)
	genLogger := InitializeLogger()
	assert.NotNil(t, genLogger, "logger must be initialized")
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

	assert.NotNil(t, genLogger.logFile)
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
	assert.NoError(t, err, test.MessageUnexpectedError(err))

	logFileContent := string(bytes)
	for _, testCase := range testCases {
		logEntry := fmt.Sprintf("%s: %v message", testCase.Input.(string), testCase.Expected.(int))
		assert.Contains(t, logFileContent, logEntry, test.MessageContains(logFileContent, logEntry, testCase.Name))
	}
	assert.NotContains(t, logFileContent, "FATAL message", test.MessageNotContains(logFileContent, "FATAL message", t.Name()))

	genLogger.SetLogLevel("INFO")
}

func TestSetOutputForTesting(t *testing.T) {
	genLogger := initializeLogTest(t)
	genLogger.SetLogLevel("WARN")

	// Capture log output using a buffer.
	var logBuffer bytes.Buffer
	genLogger.SetOutputForTesting(&logBuffer)

	// Restore original log output after test.
	defer func() {
		genLogger.SetOutputForTesting(os.Stdout)
	}()

	// Log a warning message.
	expectedMessage := "Test warning message for SetOutputForTesting"
	genLogger.Warn(expectedMessage)

	// Verify the warning was captured in the buffer.
	logOutput := logBuffer.String()
	expectedLogEntry := fmt.Sprintf("WARN: %s", expectedMessage)
	assert.Contains(t, logOutput, expectedLogEntry, test.MessageContains(logOutput, expectedLogEntry, t.Name()))
}

// TestFormattedMethods verifies that the *f variants format their arguments
// in the manner of fmt.Printf and emit the level prefix expected by callers.
func TestFormattedMethods(t *testing.T) {
	genLogger := initializeLogTest(t)
	genLogger.SetLogLevel("TRACE")
	t.Cleanup(func() {
		genLogger.SetLogLevel("INFO")
	})

	var logBuffer bytes.Buffer
	genLogger.SetOutputForTesting(&logBuffer)
	t.Cleanup(func() {
		genLogger.SetOutputForTesting(os.Stdout)
	})

	genLogger.Tracef("trace count=%d name=%s", 1, "foo")
	genLogger.Debugf("debug count=%d name=%s", 2, "bar")
	genLogger.Infof("info count=%d name=%s", 3, "baz")
	genLogger.Warnf("warn count=%d name=%s", 4, "qux")
	genLogger.Errorf("error count=%d name=%s", 5, "quux")

	logOutput := logBuffer.String()
	expectedEntries := []string{
		"TRACE: trace count=1 name=foo",
		"DEBUG: debug count=2 name=bar",
		"INFO: info count=3 name=baz",
		"WARN: warn count=4 name=qux",
		"ERROR: error count=5 name=quux",
	}
	for _, entry := range expectedEntries {
		assert.Contains(t, logOutput, entry, test.MessageContains(logOutput, entry, t.Name()))
	}
}

// TestLogCallDepth verifies that Llongfile/Lshortfile in the log output points
// at the caller of the public level method (this test file), not at logger.go.
// This locks in the calldepth chosen by constLogCallDepth.
func TestLogCallDepth(t *testing.T) {
	genLogger := initializeLogTest(t)
	genLogger.SetLogLevel("TRACE")
	t.Cleanup(func() {
		genLogger.SetLogLevel("INFO")
	})

	var logBuffer bytes.Buffer
	genLogger.SetOutputForTesting(&logBuffer)
	t.Cleanup(func() {
		genLogger.SetOutputForTesting(os.Stdout)
	})

	// Each of these calls should report logger_test.go as the source file.
	genLogger.Debug("calldepth-debug")
	genLogger.Debugf("calldepth-debugf=%d", 1)
	genLogger.Info("calldepth-info")
	genLogger.Infof("calldepth-infof=%d", 2)

	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, "logger_test.go:", test.MessageContains(logOutput, "logger_test.go:", t.Name()))
	assert.NotContains(t, logOutput, "logger.go:", test.MessageNotContains(logOutput, "logger.go:", t.Name()))
}
