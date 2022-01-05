package assert_test

import (
	"testing"

	"github.com/zpatrick/testx/assert"
)

type recorder struct {
	testing.TB
	fatalCalled bool
}

func newRecorder(t *testing.T) *recorder {
	return &recorder{TB: t}
}

func (r *recorder) Fatal(args ...any) {
	r.fatalCalled = true
}

func (r *recorder) Fatalf(format string, args ...any) {
	r.Fatal()
}

func (r *recorder) AssertFatalCalled() {
	if !r.fatalCalled {
		r.TB.Fatal("fatal was not called")
	}
}

func TestEqual(t *testing.T) {
	assert.Equal(t, 1, 1)
	assert.Equal(t, "a", "a")
}

func TestEqualFail(t *testing.T) {
	r := newRecorder(t)
	defer r.AssertFatalCalled()

	assert.Equal(r, 1, 2)
}

func TestContains(t *testing.T) {
	assert.Contains(t, []int{1, 2, 3, 4, 5}, 1)
	assert.Contains(t, []int{1, 2, 3, 4, 5}, 3)
	assert.Contains(t, []int{1, 2, 3, 4, 5}, 5)
}

func TestContainsFail(t *testing.T) {
	r := newRecorder(t)
	defer r.AssertFatalCalled()

	assert.Contains(r, []int{1, 2, 3, 4, 5}, 6)
}

func TestContainsKey(t *testing.T) {
	assert.ContainsKey(t, map[int]string{1: "", 2: ""}, 1)
	assert.ContainsKey(t, map[int]string{1: "", 2: ""}, 2)
	assert.ContainsKey(t, map[string]int{"a": 0, "b": 0}, "a")
	assert.ContainsKey(t, map[string]int{"a": 0, "b": 0}, "b")
}

func TestContainsKeyFail(t *testing.T) {
	r := newRecorder(t)
	defer r.AssertFatalCalled()

	assert.ContainsKey(r, map[int]string{1: "", 2: ""}, 3)
}

func TestContainsVal(t *testing.T) {
	assert.ContainsVal(t, map[int]string{0: "a", 1: "b"}, "a")
	assert.ContainsVal(t, map[int]string{0: "a", 1: "b"}, "b")
	assert.ContainsVal(t, map[string]int{"a": 1, "b": 2}, 1)
	assert.ContainsVal(t, map[string]int{"a": 1, "b": 2}, 2)
}

func TestContainsValFail(t *testing.T) {
	r := newRecorder(t)
	defer r.AssertFatalCalled()

	assert.ContainsVal(r, map[int]string{0: "a", 1: "b"}, "c")
}
