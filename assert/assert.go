package assert

import (
	"errors"
	"testing"
)

// Equal calls t.Fatal if result != expected.
func Equal[T comparable](t testing.TB, result, expected T) {
	if result != expected {
		t.Fatalf("%v != %v", result, expected)
	}
}

// A Comparable can be compared to other instance of the same type.
type Comparable[T any] interface {
	// Equal should return true if t is equal to the receiver.
	Equal(t T) bool
}

// EqualC calls t.Fatal if result != expected.
func EqualC[T any](t testing.TB, result Comparable[T], expected T) {
	if !result.Equal(expected) {
		t.Fatalf("%v != %v", result, expected)
	}
}

// Contains calls t.Fatal if elem is not present in s.
func Contains[T comparable](t testing.TB, s []T, elem T) {
	for _, v := range s {
		if v == elem {
			return
		}
	}

	t.Fatalf("%v not present in %v", elem, s)
}

// Equal calls t.Fatal if the specified key is not present in m.
func ContainsKey[T comparable, A any](t testing.TB, m map[T]A, key T) {
	if _, ok := m[key]; !ok {
		t.Fatalf("key %v not present in %v", key, m)
	}
}

// Equal calls t.Fatal if the specified val is not present in m.
func ContainsVal[A, T comparable](t testing.TB, m map[A]T, val T) {
	for _, v := range m {
		if v == val {
			return
		}
	}

	t.Fatalf("val %v not present in %v", val, m)
}

// Error calls t.Fatal if err is nil.
func Error(t testing.TB, err error) {
	if err == nil {
		t.Fatalf("error is nil")
	}
}

// NilError calls t.Fatal if err is not nil.
func NilError(t testing.TB, err error) {
	if err == nil {
		return
	}

	t.Fatalf("error is not nil: %v", err)
}

// ErrorsIs calls t.Fatal if errors.Is(err, target) fails.
func ErrorIs(t testing.TB, err, target error) {
	if !errors.Is(err, target) {
		t.Fatalf("error %v is not %v", err, target)
	}
}
