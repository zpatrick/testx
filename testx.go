package testx

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

type Wrapper func(tb testing.TB) testing.TB

func Wrap(tb testing.TB, wrappers ...Wrapper) testing.TB {
	for _, w := range wrappers {
		tb = w(tb)
	}
	return tb
}

type LoggingMethodSwitch struct {
	Error  bool
	Errorf bool
	Fatal  bool
	Fatalf bool
	Log    bool
	Logf   bool
}

type LoggingOptions struct {
	DisableMethods LoggingMethodSwitch
}

type logger struct {
	testing.TB

	opts    LoggingOptions
	format  func(method func(...any), args ...any)
	formatf func(method func(string, ...any), format string, args ...any)
}

func (l *logger) Log(args ...any) {
	if l.opts.DisableMethods.Log {
		l.TB.Log(args...)
		return
	}

	l.format(l.TB.Log, args...)
}

func (l *logger) Logf(format string, args ...any) {
	if l.opts.DisableMethods.Logf {
		l.TB.Logf(format, args...)
		return
	}

	l.formatf(l.TB.Logf, format, args...)
}

func (l *logger) Error(args ...any) {
	if l.opts.DisableMethods.Error {
		l.TB.Error(args...)
		return
	}

	l.format(l.TB.Error, args...)
}

func (l *logger) Errorf(format string, args ...any) {
	if l.opts.DisableMethods.Errorf {
		l.TB.Errorf(format, args...)
		return
	}

	l.formatf(l.TB.Errorf, format, args...)
}

func (l *logger) Fatal(args ...any) {
	if l.opts.DisableMethods.Fatal {
		l.TB.Fatal(args...)
		return
	}

	l.format(l.TB.Fatal, args...)
}

func (l *logger) Fatalf(format string, args ...any) {
	if l.opts.DisableMethods.Fatalf {
		l.TB.Fatalf(format, args...)
		return
	}

	l.formatf(l.TB.Fatalf, format, args...)
}

type PrefixLoggingOptions struct {
	DisableMethods LoggingMethodSwitch
	Prefix         func() string
}

func WithPrefixLogging(opts PrefixLoggingOptions) Wrapper {
	return func(tb testing.TB) testing.TB {
		return newPrefixLogger(tb, opts)
	}
}

func Prefix(tb testing.TB, prefix string) testing.TB {
	return newPrefixLogger(tb, PrefixLoggingOptions{
		Prefix: func() string { return prefix },
	})
}

func Prefixf(tb testing.TB, format string, args ...any) testing.TB {
	return Prefix(tb, fmt.Sprintf(format, args...))
}

func newPrefixLogger(tb testing.TB, opts PrefixLoggingOptions) *logger {
	return &logger{
		TB:   tb,
		opts: LoggingOptions{DisableMethods: opts.DisableMethods},
		format: func(method func(...any), args ...any) {
			method(append([]any{opts.Prefix()}, args...)...)
		},
		formatf: func(method func(string, ...any), format string, args ...any) {
			method(opts.Prefix() + " " + fmt.Sprintf(format, args...))
		},
	}
}

type StackTraceLoggingOptions struct {
	DisableMethods LoggingMethodSwitch
	BufferSize     int
	SkipLines      int
	All            bool
}

func WithStackTraceLogging(opts StackTraceLoggingOptions) Wrapper {
	return func(tb testing.TB) testing.TB {
		return newStackTraceLogger(tb, opts)
	}
}

func newStackTraceLogger(tb testing.TB, opts StackTraceLoggingOptions) *logger {
	stack := func() string {
		buf := make([]byte, opts.BufferSize)
		runtime.Stack(buf, opts.All)

		if skip := opts.SkipLines; skip > 0 {
			lines := strings.Split(string(buf), "\n")
			if skip >= len(lines) {
				skip = len(lines) - 1
			}

			lines = lines[skip:]
			buf = []byte(strings.Join(lines, "\n"))
		}

		return fmt.Sprintf("\nstack: %s", string(buf))
	}

	return &logger{
		TB:   tb,
		opts: LoggingOptions{DisableMethods: opts.DisableMethods},
		format: func(method func(...any), args ...any) {
			method(append(args, stack())...)
		},
		formatf: func(method func(string, ...any), format string, args ...any) {
			method(fmt.Sprintf(format, args...) + stack())
		},
	}
}
