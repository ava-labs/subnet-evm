// Package testutils contains test utilities ONLY to be used outside plugin/evm.
// The aim is to reduce changes in geth tests by using the utilities defined here.
// This package MUST NOT be imported by non-test packages.

package testutils

import (
	"runtime"
	"strings"
)

// panicIfCallsFromNonTest should be added at the top of every function defined in this package
// to enforce this package to be used only by tests.
func panicIfCallsFromNonTest() {
	pc := make([]uintptr, 64)
	runtime.Callers(0, pc)
	frames := runtime.CallersFrames(pc)
	for {
		f, more := frames.Next()
		if strings.HasPrefix(f.File, "/testing/") || strings.HasSuffix(f.File, "_test.go") {
			return
		}
		if !more {
			panic("no test file in call stack")
		}
	}
}
