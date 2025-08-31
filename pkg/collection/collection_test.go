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

func Test_GetObjAsMap(t *testing.T) {
	m := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	result, err := GetObjAsMap[string](interface{}(m))

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(result, m) {
		t.Fatalf("failed to convert %v to type %T", m, result)
	}
}

func Test_GetObjAsMap_IntValues(t *testing.T) {
	m := map[string]int{
		"count": 42,
		"total": 100,
	}

	result, err := GetObjAsMap[int](interface{}(m))

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(result, m) {
		t.Fatalf("failed to convert %v to type %T", m, result)
	}
}

func Test_GetObjAsMap_Error(t *testing.T) {
	// Test with non-map type
	s := "not a map"

	_, err := GetObjAsMap[string](interface{}(s))

	if err == nil {
		t.Fatal("expected error when converting non-map to map, but got nil")
	}
}

func Test_GetObjAsMap_WrongValueType(t *testing.T) {
	// Test with map but wrong value type
	m := map[string]int{
		"key1": 1,
		"key2": 2,
	}

	_, err := GetObjAsMap[string](interface{}(m))

	if err == nil {
		t.Fatal("expected error when converting map[string]int to map[string]string, but got nil")
	}
}
