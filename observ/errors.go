package observ

import (
	"fmt"
	"runtime"
)

type StackTracer interface {
	Stacktrace() string
}

func Wrap(internalErr error, publicErr error) error {
	return wrapErrorWithStack(internalErr, publicErr)
}

type stacktraceError struct {
	err        error
	publicErr  error
	stacktrace string
}

func (e *stacktraceError) Unwrap() []error {
	return []error{e.publicErr, e.err}
}

func (e *stacktraceError) Error() string {
	return fmt.Sprintf("%s: %s", e.publicErr.Error(), e.err.Error())
}

func (e *stacktraceError) Stacktrace() string {
	return e.stacktrace
}

func wrapErrorWithStack(err error, publicErr error) error {
	stackBuf := make([]byte, 1024*4) // 5KB buffer
	stackLen := runtime.Stack(stackBuf, false)

	return &stacktraceError{
		err:        err,
		publicErr:  publicErr,
		stacktrace: string(stackBuf[:stackLen]),
	}
}
