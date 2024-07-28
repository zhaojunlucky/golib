package collection

import (
	"reflect"
	"testing"
)

func Test_GetObjAsSlice(t *testing.T) {
	s := "hello"

	arr, err := GetObjAsSlice[string](interface{}(s))

	if err != nil {
		t.Fatal(err)
	}

	if len(arr) != 1 {
		t.Fatalf("expected 1, got %d", len(arr))
	}

	if arr[0] != s {
		t.Fatalf("failed to convert %v to type %T", s, arr)
	}
}

func Test_GetObjAsSliceArr(t *testing.T) {
	s := []string{"hello", "world", "1"}

	arr, err := GetObjAsSlice[string](interface{}(s))

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(arr, s) {
		t.Fatalf("failed to convert %v to type %T", s, arr)
	}
}
