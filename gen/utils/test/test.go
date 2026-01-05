package test

import (
	"testing"
)

// InitializeTest is a helper function that logs the test name.
// It should be called at the beginning of each test.
func InitializeTest(t *testing.T) {
	t.Helper()
	t.Log("Executing:", t.Name())
}
