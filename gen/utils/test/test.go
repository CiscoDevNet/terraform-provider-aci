package test

import (
	"fmt"
	"testing"
)

func InitializeTest(t *testing.T) {
	println(fmt.Sprintf("Executing: %s", t.Name()))
}
