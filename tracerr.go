package tracerr

import (
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
)

type TraceableError interface {
	// TraceableError is a no-op function used inside the Wrap function to
	// distinguish this type of error from any other.
	TraceableError()
}

type traceableError struct {
	err   error
	stack []*stackFrame
}

// Max depth for catured stack traces
const maxStackDepth = 64

// Wrap takes an existing error and attaches the current stack trace to it.
// The passed in error will not be wrapped if it is nil or it already has a
// stack trace attached to it (i.e. it has already been wrapped).
func Wrap(err error) error {
	// Handle nil errors
	if err == nil {
		return err
	}

	// Check if err is already implementing TraceError
	if _, ok := err.(TraceableError); ok {
		return err
	}

	// Capture stack trace and wrap err
	return &traceableError{
		err:   err,
		stack: captureStack(2, maxStackDepth),
	}
}

func (t *traceableError) TraceableError() {}

// Error returns the original error message plus the stack trace captured
// at the time the error was first wrapped.
func (t *traceableError) Error() string {
	str := t.err.Error()
	for _, frame := range t.stack {
		str += fmt.Sprintf("\n  at %s", frame.string())
	}
	return str
}

// stackFrame represents a particular function in the call stack.
type stackFrame struct {
	file     string
	line     int
	function string
}

// string converts a given stack frame to a formated string.
func (s *stackFrame) string() string {
	return fmt.Sprintf("%s (%s:%d)", s.function, s.file, s.line)
}

// newStackFrame returns a new stack frame initialized from the passed
// in program counter.
func newStackFrame(pc uintptr) *stackFrame {
	fn := runtime.FuncForPC(pc)
	file, line := fn.FileLine(pc)
	packagePath, funcSignature := parseFuncName(fn.Name())
	_, fileName := filepath.Split(file)

	return &stackFrame{
		file:     filepath.Join(packagePath, fileName),
		line:     line,
		function: funcSignature,
	}
}

// captureStack returns a slice of stack frames representing the stack
// of the calling go routine.
func captureStack(skip, maxDepth int) []*stackFrame {
	pcs := make([]uintptr, maxDepth)
	count := runtime.Callers(skip+1, pcs)

	frames := make([]*stackFrame, count)
	for i, pc := range pcs[0:count] {
		frames[i] = newStackFrame(pc)
	}

	return frames
}

// parseFuncName returns the package path and function signature for a
// give Func name.
func parseFuncName(fnName string) (packagePath, signature string) {
	regEx := regexp.MustCompile("([^\\(]*)\\.(.*)")
	parts := regEx.FindStringSubmatch(fnName)
	packagePath = parts[1]
	signature = parts[2]
	return
}
