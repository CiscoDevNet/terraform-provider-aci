package test

import (
	"runtime"
	"strings"
	"testing"
)

func InitializeTest(t *testing.T, functionCounter int) {
	println("Executing:", GetTestName(t, functionCounter+1))
}

func GetTestName(t *testing.T, functionCounter int) string {
	counter, _, _, success := runtime.Caller(functionCounter + 1)
	if !success {
		t.Fatalf("Failed to get caller information: %v", success)
	}
	splittedName := strings.Split(runtime.FuncForPC(counter).Name(), ".")
	return splittedName[len(splittedName)-1]
}
