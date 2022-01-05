package assert_test

import (
	"errors"
	"math"
	"os"
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

type roundedFloat float64

func (r roundedFloat) Equal(o roundedFloat) bool {
	return math.Round(float64(r)) == math.Round(float64(o))
}

func TestEqualC(t *testing.T) {
	assert.EqualC[roundedFloat](t, roundedFloat(1.1), roundedFloat(1.2))
}

func TestEqualCFail(t *testing.T) {
	r := newRecorder(t)
	defer r.AssertFatalCalled()

	assert.EqualC[roundedFloat](r, roundedFloat(1.1), roundedFloat(2.1))
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

func TestError(t *testing.T) {
	var err error = errors.New("")
	assert.Error(t, err)
}

func TestErrorFail(t *testing.T) {
	r := newRecorder(t)
	defer r.AssertFatalCalled()

	var err error
	assert.Error(r, err)
}

func TestNilError(t *testing.T) {
	var err error
	assert.NilError(t, err)
}

func TestNilErrorFail(t *testing.T) {
	r := newRecorder(t)
	defer r.AssertFatalCalled()

	var err error = errors.New("")
	assert.NilError(r, err)
}

func TestErrorIs(t *testing.T) {
	var err error = os.ErrClosed
	assert.ErrorIs(t, err, os.ErrClosed)
}

func TestErrorIsFail(t *testing.T) {
	r := newRecorder(t)
	defer r.AssertFatalCalled()

	var err error = errors.New("")
	assert.ErrorIs(r, err, os.ErrClosed)
}

func TestErrorAs(t *testing.T) {
	var err error = os.ErrClosed
	assert.ErrorAs(t, err, &os.ErrClosed)
}

func TestErrorAsFail(t *testing.T) {
	r := newRecorder(t)
	defer r.AssertFatalCalled()

	var err error
	assert.ErrorAs(r, err, &os.ErrClosed)
}
