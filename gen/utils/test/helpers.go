package test

import (
	"fmt"
	"testing"
)

// InitializeTest is a helper function that logs the test name.
// It should be called at the beginning of each test.
func InitializeTest(t *testing.T) {
	t.Helper()
	t.Log("Executing:", t.Name())
}

// MessageUnexpectedError returns a formatted error message for unexpected errors.
func MessageUnexpectedError(err error) string {
	return fmt.Sprintf("Expected no error, but got '%s'", err)
}

// MessageEqual returns a formatted message for assert.Equal comparisons.
func MessageEqual(expected, actual any, caseName string) string {
	return fmt.Sprintf("Expected '%v', but got '%v' for case '%s'", expected, actual, caseName)
}

// MessageContains returns a formatted message for assert.Contains comparisons.
func MessageContains(collection, element any, caseName string) string {
	return fmt.Sprintf("Expected '%v' to contain '%v' for case '%s'", collection, element, caseName)
}

// MessageNotContains returns a formatted message for assert.NotContains comparisons.
func MessageNotContains(collection, element any, caseName string) string {
	return fmt.Sprintf("Expected '%v' to not contain '%v' for case '%s'", collection, element, caseName)
}

// MessageNotEmpty returns a formatted message for require.NotEmpty comparisons.
func MessageNotEmpty(object any, caseName string) string {
	return fmt.Sprintf("Expected '%v' to not be empty for case '%s'", object, caseName)
}

// MessageEmpty returns a formatted message for assert.Empty comparisons.
func MessageEmpty(object any, caseName string) string {
	return fmt.Sprintf("Expected '%v' to be empty for case '%s'", object, caseName)
}

// TestCase is a generic struct for table-driven tests.
// It provides a consistent structure for defining test cases with inputs and expected outputs.
//
// Fields:
//   - Name: A descriptive identifier for the test case, used with t.Run() to create named subtests.
//   - Input: The input data for the test. Can be any type (string, struct, etc.) and requires
//     type assertion in the test loop to access the concrete type.
//   - Expected: The expected result to assert against. Can be any type (bool, string, struct, etc.)
//     and requires type assertion in the test loop.
//
// For tests with multiple inputs or complex expected values, define helper structs
// (e.g., myTestInput, myTestExpected) and use them as the Input and Expected types.
type TestCase struct {
	Name     string
	Input    any
	Expected any
}
