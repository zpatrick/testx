package assert

import (
	"errors"
	"testing"
)

// Equal calls t.Fatalf if result != expected.
func Equal[T comparable](t testing.TB, result, expected T) {
	t.Helper()

	if result != expected {
		t.Fatalf("%v != %v", result, expected)
	}
}

// EqualSlices calls t.Fatalf if result expected do not contain the same elements in the same order.
func EqualSlices[T comparable, TS ~[]T](t testing.TB, result, expected TS) {
	t.Helper()

	if resLen, expLen := len(result), len(expected); resLen != expLen {
		t.Fatalf("%v != %v: slices are not the same length: %d != %d", result, expected, resLen, expLen)
	}

	for i := 0; i < len(expected)-1; i++ {
		if res, exp := result[i], expected[i]; res != exp {
			t.Fatalf("%v != %v: unequal elements at index %d: %v != %v", result, expected, i, res, exp)
		}
	}
}

// EqualMaps calls t.Fatalf if result expected do not contain the same elements.
func EqualMaps[K, V comparable](t testing.TB, result, expected map[K]V) {
	t.Helper()

	if resLen, expLen := len(result), len(expected); resLen != expLen {
		t.Fatalf("%v != %v: maps are not the same length: %d != %d", result, expected, resLen, expLen)
	}

	for expKey, expVal := range expected {
		resVal, ok := result[expKey]
		if !ok {
			t.Fatalf("%v != %v: result did not contain key %v", result, expected, expKey)
		}

		if expVal != resVal {
			t.Fatalf("%v != %v: unequal elements at key %v: %v != %v", result, expected, expKey, resVal, expVal)
		}
	}
}

// A Comparable can be compared to other instances of the same type.
type Comparable[T any] interface {
	// Equal should return true if t is equal to the receiver.
	Equal(t T) bool
}

// EqualC calls t.Fatalf if result != expected.
func EqualC[T any](t testing.TB, result Comparable[T], expected T) {
	t.Helper()

	if !result.Equal(expected) {
		t.Fatalf("%v != %v", result, expected)
	}
}

// Contains calls t.Fatalf if any value v is not present in s.
func Contains[T comparable](t testing.TB, s []T, v ...T) {
	t.Helper()

	for _, expected := range v {
		var found bool
		for _, actual := range s {
			if expected == actual {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("%v not present in %v", expected, s)
		}
	}
}

// ContainsKeys calls t.Fatalf if any of the specified keys are not present in m.
func ContainsKeys[T comparable, A any](t testing.TB, m map[T]A, keys ...T) {
	t.Helper()

	for _, key := range keys {
		if _, ok := m[key]; !ok {
			t.Fatalf("key %v not present in %v", key, m)
		}
	}
}

// ContainsVals calls t.Fatalf if any of the specified vals are not present in m.
func ContainsVals[A, T comparable](t testing.TB, m map[A]T, vals ...T) {
	t.Helper()

	for _, expected := range vals {
		var found bool
		for _, actual := range m {
			if expected == actual {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("val %v not present in %v", expected, m)
		}
	}
}

// Error calls t.Fatalf if err is nil.
func Error(t testing.TB, err error) {
	t.Helper()

	if err == nil {
		t.Fatalf("error is nil")
	}
}

// NilError calls t.Fatalf if err is not nil.
func NilError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("error is not nil: %v", err)
	}
}

// ErrorsIs calls t.Fatalf if errors.Is(err, target) fails.
func ErrorIs(t testing.TB, err, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("error.Is check failed for %v (target: %v)", err, target)
	}
}

// ErrorsAs calls t.Fatalf if errors.As(err, target) fails.
func ErrorAs(t testing.TB, err error, target any) {
	t.Helper()

	if !errors.As(err, target) {
		t.Fatalf("error.As check failed for %v (target: %v)", err, target)
	}
}
