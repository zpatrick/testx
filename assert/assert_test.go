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

func TestEqualSlices(t *testing.T) {
	assert.EqualSlices(t, []int{1, 2}, []int{1, 2})
	assert.EqualSlices(t, []string{"a", "b"}, []string{"a", "b"})
	assert.EqualSlices[int](t, nil, nil)
}

func TestEqualSlicesFail(t *testing.T) {
	testCases := []struct {
		Name string
		A    []int
		B    []int
	}{
		{"nil", []int{1, 2, 3}, nil},
		{"missing element", []int{1, 2, 3}, []int{1, 2}},
		{"extra element", []int{1, 2, 3}, []int{1, 2, 3, 4}},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			r := newRecorder(t)
			defer r.AssertFatalCalled()

			assert.EqualSlices(r, tc.A, tc.B)
		})
	}
}

func TestEqualMaps(t *testing.T) {
	assert.EqualMaps(t, map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1, "b": 2})
	assert.EqualMaps[int, bool](t, nil, nil)
}

func TestEqualMapsFail(t *testing.T) {
	testCases := []struct {
		Name string
		A    map[string]int
		B    map[string]int
	}{
		{"nil", map[string]int{"a": 1, "b": 2}, nil},
		{"missing element", map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1}},
		{"extra element", map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1, "b": 2, "c": 3}},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			r := newRecorder(t)
			defer r.AssertFatalCalled()

			assert.EqualMaps(r, tc.A, tc.B)
		})
	}
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
	assert.Contains(t, []int{1, 2, 3, 4, 5}, 1, 2, 3)
	assert.Contains(t, []int{1, 2, 3, 4, 5}, 2, 5)
}

func TestContainsFail(t *testing.T) {
	testCases := []struct {
		Name string
		S    []int
		V    []int
	}{
		{"below", []int{1, 2, 3}, []int{0}},
		{"above", []int{1, 2, 3}, []int{4}},
		{"all out", []int{1, 2, 3}, []int{5, 6}},
		{"mixed", []int{1, 2, 3}, []int{3, 4}},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			r := newRecorder(t)
			defer r.AssertFatalCalled()

			assert.Contains(r, tc.S, tc.V...)
		})
	}
}

func TestContainsKeys(t *testing.T) {
	assert.ContainsKeys(t, map[int]string{1: "", 2: ""}, 1)
	assert.ContainsKeys(t, map[int]string{1: "", 2: ""}, 2)
	assert.ContainsKeys(t, map[int]string{1: "", 2: ""}, 1, 2)
	assert.ContainsKeys(t, map[string]int{"a": 0, "b": 0}, "a")
	assert.ContainsKeys(t, map[string]int{"a": 0, "b": 0}, "b")
	assert.ContainsKeys(t, map[string]int{"a": 0, "b": 0}, "a", "b")
}

func TestContainsKeysFail(t *testing.T) {
	testCases := []struct {
		Name string
		M    map[string]int
		Keys []string
	}{
		{"out", map[string]int{"a": 1, "b": 2}, []string{"c"}},
		{"mixed", map[string]int{"a": 1, "b": 2}, []string{"b", "c"}},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			r := newRecorder(t)
			defer r.AssertFatalCalled()

			assert.ContainsKeys(r, tc.M, tc.Keys...)
		})
	}
}

func TestContainsVals(t *testing.T) {
	assert.ContainsVals(t, map[int]string{0: "a", 1: "b"}, "a")
	assert.ContainsVals(t, map[int]string{0: "a", 1: "b"}, "b")
	assert.ContainsVals(t, map[int]string{0: "a", 1: "b"}, "a", "b")
	assert.ContainsVals(t, map[string]int{"a": 1, "b": 2}, 1)
	assert.ContainsVals(t, map[string]int{"a": 1, "b": 2}, 2)
	assert.ContainsVals(t, map[string]int{"a": 1, "b": 2}, 1, 2)
}

func TestContainsValsFail(t *testing.T) {
	testCases := []struct {
		Name string
		M    map[string]int
		Vals []int
	}{
		{"out", map[string]int{"a": 1, "b": 2}, []int{3}},
		{"mixed", map[string]int{"a": 1, "b": 2}, []int{2, 3}},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			r := newRecorder(t)
			defer r.AssertFatalCalled()

			assert.ContainsVals(r, tc.M, tc.Vals...)
		})
	}
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
