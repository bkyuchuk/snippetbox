package assert

import (
	"reflect"
	"testing"
)

func Equal[T any](t *testing.T, got, want T) {
	t.Helper()

	if !isEqual(got, want) {
		t.Errorf("got %v; want: %v", got, want)
	}
}

func isEqual[T any](got, want T) bool {
	if isNil(got) != isNil(want) {
		return false
	}

	return reflect.DeepEqual(got, want)
}

func isNil(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice, reflect.Chan, reflect.Interface, reflect.Pointer, reflect.UnsafePointer:
		return rv.IsNil()
	}

	return false
}

func Nil(t *testing.T, got any) {
	t.Helper()

	if !isNil(got) {
		t.Errorf("got %q; want: nil", got)
	}
}

func True(t *testing.T, got bool) {
	t.Helper()

	if !got {
		t.Errorf("got %t; want: true", got)
	}
}
