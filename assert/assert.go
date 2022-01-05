package assert

import "testing"

func Equal[T comparable](t testing.TB, result, expected T) {
	if result != expected {
		t.Fatalf("%v != %v", result, expected)
	}
}

func Contains[T comparable](t testing.TB, s []T, elem T) {
	for _, v := range s {
		if v == elem {
			return
		}
	}

	t.Fatalf("%v not present in %v", elem, s)
}

func ContainsKey[T comparable, A any](t testing.TB, m map[T]A, key T) {
	if _, ok := m[key]; !ok {
		t.Fatalf("key %v not present in %v", key, m)
	}
}

func ContainsVal[A, T comparable](t testing.TB, m map[A]T, val T) {
	for _, v := range m {
		if v == val {
			return
		}
	}

	t.Fatalf("val %v not present in %v", val, m)
}
