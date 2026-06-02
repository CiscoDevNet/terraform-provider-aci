package logger

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

// newTestLogger resets the package singleton, constructs a fresh logger,
// installs a no-op exitFunc so Fatal/Fatalf don't terminate the test process,
// and points output at the provided buffer. The cleanup closes any open
// log file and resets the singleton again.
func newTestLogger(t *testing.T, outputBuffer *bytes.Buffer) *Logger {
	t.Helper()
	test.InitializeTest(t)
	ResetForTest()
	t.Cleanup(ResetForTest)

	testLogger := InitializeLogger()
	testLogger.SetExitFunc(func(int) {}) // make Fatal/Fatalf non-terminating
	if outputBuffer != nil {
		testLogger.SetOutputForTesting(outputBuffer)
	}
	return testLogger
}

func TestInitializeLogger_FreshState(t *testing.T) {
	testLogger := newTestLogger(t, nil)

	assert.NotNil(t, testLogger)
	assert.Nil(t, testLogger.logFile)
	assert.Equal(t, logLevels[constDefaultLogLevel], testLogger.logLevel)
}

func TestInitializeLogger_FromEnv(t *testing.T) {
	logFilePath := filepath.Join(t.TempDir(), "env.log")
	t.Setenv(constEnvLogPath, logFilePath)
	t.Setenv(constEnvLogLevel, "TRACE")

	testLogger := newTestLogger(t, nil)

	assert.NotNil(t, testLogger.logFile)
	assert.Equal(t, logFilePath, testLogger.logFile.Name())
	assert.Equal(t, logLevels["TRACE"], testLogger.logLevel)
}

func TestSetLogFile_OpenFailureCallsExitFunc(t *testing.T) {
	testLogger := newTestLogger(t, &bytes.Buffer{})

	var exitCode int
	called := false
	testLogger.SetExitFunc(func(code int) {
		exitCode = code
		called = true
	})

	// A path under a non-existent directory cannot be opened.
	testLogger.SetLogFile(filepath.Join(t.TempDir(), "nope", "no-such-dir", "x.log"))

	assert.True(t, called, "expected exitFunc to be invoked on open failure")
	assert.Equal(t, 1, exitCode)
	assert.Nil(t, testLogger.logFile, "logFile must not be set when open fails")
}

func TestSetLogLevel_InvalidCallsExitFunc(t *testing.T) {
	testLogger := newTestLogger(t, &bytes.Buffer{})

	var exitCode int
	called := false
	testLogger.SetExitFunc(func(code int) {
		exitCode = code
		called = true
	})

	testLogger.SetLogLevel("NOT_A_LEVEL")

	assert.True(t, called, "expected exitFunc to be invoked on invalid level")
	assert.Equal(t, 1, exitCode)
}

func TestFatalAndFatalf_RouteThroughExitFunc(t *testing.T) {
	var outputBuffer bytes.Buffer
	testLogger := newTestLogger(t, &outputBuffer)

	exitCodes := []int{}
	testLogger.SetExitFunc(func(code int) { exitCodes = append(exitCodes, code) })

	testLogger.Fatal("plain fatal")
	testLogger.Fatalf("formatted fatal: %d", 42)

	assert.Equal(t, []int{1, 1}, exitCodes)
	out := outputBuffer.String()
	assert.Contains(t, out, "FATAL: plain fatal")
	assert.Contains(t, out, "FATAL: formatted fatal: 42")
}

func TestLogLevelFiltering(t *testing.T) {
	// Use a fresh buffer per case so assertions are independent.
	cases := []struct {
		level      string
		expectShow []string
		expectHide []string
	}{
		{
			level:      "WARN",
			expectShow: []string{"FATAL: f", "ERROR: e", "WARN: w"},
			expectHide: []string{"INFO: i", "DEBUG: d", "TRACE: t"},
		},
		{
			level:      "INFO",
			expectShow: []string{"FATAL: f", "ERROR: e", "WARN: w", "INFO: i"},
			expectHide: []string{"DEBUG: d", "TRACE: t"},
		},
		{
			level:      "TRACE",
			expectShow: []string{"FATAL: f", "ERROR: e", "WARN: w", "INFO: i", "DEBUG: d", "TRACE: t"},
			expectHide: nil,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.level, func(t *testing.T) {
			var outputBuffer bytes.Buffer
			testLogger := newTestLogger(t, &outputBuffer)
			testLogger.SetLogLevel(testCase.level)

			testLogger.Fatal("f")
			testLogger.Error("e")
			testLogger.Warn("w")
			testLogger.Info("i")
			testLogger.Debug("d")
			testLogger.Trace("t")

			out := outputBuffer.String()
			for _, want := range testCase.expectShow {
				assert.Contains(t, out, want, test.MessageContains(out, want, testCase.level))
			}
			for _, hide := range testCase.expectHide {
				assert.NotContains(t, out, hide, test.MessageNotContains(out, hide, testCase.level))
			}
		})
	}
}

func TestCloseLogFile(t *testing.T) {
	logFilePath := filepath.Join(t.TempDir(), "close.log")
	testLogger := newTestLogger(t, nil)
	testLogger.SetLogFile(logFilePath)
	assert.NotNil(t, testLogger.logFile)

	testLogger.CloseLogFile()

	assert.Nil(t, testLogger.logFile, "logFile should be nil after close")
	// Second close must be a no-op.
	testLogger.CloseLogFile()
	assert.Nil(t, testLogger.logFile)
}

func TestSetOutputForTesting(t *testing.T) {
	var outputBuffer bytes.Buffer
	testLogger := newTestLogger(t, &outputBuffer)
	testLogger.SetLogLevel("WARN")

	expected := "Test warning message for SetOutputForTesting"
	testLogger.Warn(expected)

	assert.Contains(t, outputBuffer.String(), fmt.Sprintf("WARN: %s", expected))
}

func TestFormattedMethods(t *testing.T) {
	var outputBuffer bytes.Buffer
	testLogger := newTestLogger(t, &outputBuffer)
	testLogger.SetLogLevel("TRACE")

	testLogger.Tracef("trace count=%d name=%s", 1, "foo")
	testLogger.Debugf("debug count=%d name=%s", 2, "bar")
	testLogger.Infof("info count=%d name=%s", 3, "baz")
	testLogger.Warnf("warn count=%d name=%s", 4, "qux")
	testLogger.Errorf("error count=%d name=%s", 5, "quux")

	out := outputBuffer.String()
	for _, entry := range []string{
		"TRACE: trace count=1 name=foo",
		"DEBUG: debug count=2 name=bar",
		"INFO: info count=3 name=baz",
		"WARN: warn count=4 name=qux",
		"ERROR: error count=5 name=quux",
	} {
		assert.Contains(t, out, entry, test.MessageContains(out, entry, t.Name()))
	}
}

// TestLogCallDepth verifies that Llongfile/Lshortfile in the log output points
// at the caller of the public level method (this test file), not at logger.go.
// This locks in the calldepth chosen by constLogCallDepth.
func TestLogCallDepth(t *testing.T) {
	var outputBuffer bytes.Buffer
	testLogger := newTestLogger(t, &outputBuffer)
	testLogger.SetLogLevel("TRACE")
	// SetLogLevel itself emits a Tracef from logger.go; reset so the assertion
	// only inspects output produced by the calls under test below.
	outputBuffer.Reset()

	testLogger.Debug("calldepth-debug")
	testLogger.Debugf("calldepth-debugf=%d", 1)
	testLogger.Info("calldepth-info")
	testLogger.Infof("calldepth-infof=%d", 2)

	out := outputBuffer.String()
	assert.Contains(t, out, "logger_test.go:", test.MessageContains(out, "logger_test.go:", t.Name()))
	assert.NotContains(t, out, "logger.go:", test.MessageNotContains(out, "logger.go:", t.Name()))
}

func TestSetLogFile_WritesToFile(t *testing.T) {
	logFilePath := filepath.Join(t.TempDir(), "writes.log")
	testLogger := newTestLogger(t, nil)
	testLogger.SetLogFile(logFilePath)
	t.Cleanup(testLogger.CloseLogFile)

	testLogger.SetLogLevel("INFO")
	testLogger.Info("hello-file")

	contents, err := os.ReadFile(logFilePath)
	assert.NoError(t, err)
	assert.Contains(t, string(contents), "INFO: hello-file")
}

// TestInitializeLogger_Idempotent verifies that calling InitializeLogger
// twice returns the same singleton instance (production behavior relied on
// by per-package var genLogger = logger.InitializeLogger()).
func TestInitializeLogger_Idempotent(t *testing.T) {
	test.InitializeTest(t)
	ResetForTest()
	t.Cleanup(ResetForTest)

	a := InitializeLogger()
	b := InitializeLogger()
	assert.Same(t, a, b, "InitializeLogger must return the same singleton on subsequent calls")
}
