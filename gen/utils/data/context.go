package data

import (
	"fmt"
	"strings"
)

// Context carries shared state across the generation pipeline.
type Context struct {
	Diagnostics *Diagnostics
}

func NewContext() *Context {
	return &Context{
		Diagnostics: NewDiagnostics(),
	}
}

// Diagnostics accumulates non-fatal errors and warnings during generation.
// Hard failures (I/O, JSON) still return error immediately; definition mistakes
// (missing role, unresolved placeholders) accumulate here for a single summary.
type Diagnostics struct {
	errors []string
}

func NewDiagnostics() *Diagnostics {
	return &Diagnostics{}
}

func (d *Diagnostics) AddError(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	d.errors = append(d.errors, msg)
	// Debug-level only: the final summary printed by ctx.Diagnostics.Error()
	// surfaces every accumulated diagnostic; logging at Errorf here would
	// double-report every issue.
	genLogger.Debugf(format, v...)
}

func (d *Diagnostics) Error() error {
	if len(d.errors) == 0 {
		return nil
	}
	return fmt.Errorf("generation encountered %d error(s):\n  - %s", len(d.errors), strings.Join(d.errors, "\n  - "))
}
