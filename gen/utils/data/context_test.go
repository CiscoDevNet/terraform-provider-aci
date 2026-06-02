package data

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiagnosticsAddError(t *testing.T) {
	t.Parallel()
	d := NewDiagnostics()

	assert.NoError(t, d.Error(), "empty Diagnostics should return nil error")

	d.AddError("first error for class '%s'", "fvTenant")
	d.AddError("second error for class '%s'", "fvBD")

	err := d.Error()
	assert.Error(t, err, "Diagnostics with errors should return non-nil")
	assert.True(t, strings.Contains(err.Error(), "first error for class 'fvTenant'"), "summary should contain first error")
	assert.True(t, strings.Contains(err.Error(), "second error for class 'fvBD'"), "summary should contain second error")
	assert.True(t, strings.Contains(err.Error(), "2 error(s)"), "summary should include count")
}
