package testx_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/zpatrick/testx"
	"github.com/zpatrick/testx/assert"
)

var (
	stackTraceLoggingOpts = testx.StackTraceLoggingOptions{
		BufferSize: 2000,
		SkipLines:  6,
		DisableMethods: testx.LoggingMethodSwitch{
			Fatal:  false,
			Fatalf: true,
			Log:    true,
			Logf:   true,
		},
	}
	prefixLoggingOpts = testx.PrefixLoggingOptions{
		PrefixFunc: func() string {
			return fmt.Sprintf("(%s)", time.Now().Format("15:04:05.9999"))
		},
	}
)

func setupTB(t testing.TB) testing.TB {
	return testx.Wrap(t,
		testx.WithStackTraceLogging(stackTraceLoggingOpts),
		testx.WithPrefixLogging(prefixLoggingOpts),
	)
}

func TestExample(t *testing.T) {
	tb := setupTB(t)
	tb.Log("test running")

	assert.Equal(tb, 1, 1)
	tb.Log("test is still running")
	assert.Equal(testx.Prefix(tb, "unexpected random number generated:"), rand.Int(), 1)
	tb.Log("test completed")
}
