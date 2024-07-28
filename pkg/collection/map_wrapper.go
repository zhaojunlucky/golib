package collection

import (
	"fmt"
	"reflect"
)

type MapWrapper struct {
	data map[string]any
}

func (m *MapWrapper) Has(key string) bool {
	_, ok := m.data[key]
	return ok
}

func (m *MapWrapper) GetChild(key string) (*MapWrapper, error) {
	mapObj, ok := m.data[key]
	if !ok {
		return nil, fmt.Errorf("key %s not found in map", key)
	}
	return NewMapWrapper(mapObj.(map[string]any)), nil
}

func (m *MapWrapper) GetAny(key string) (any, error) {
	mapObj, ok := m.data[key]
	if !ok {
		return nil, fmt.Errorf("key %s not found in map", key)
	}
	return mapObj, nil

}

func (m *MapWrapper) Get(key string, val any) error {
	mapObj, ok := m.data[key]
	if !ok {
		return fmt.Errorf("key %s not found in map", key)
	}

	valType := reflect.ValueOf(val)

	if valType.Kind() != reflect.Pointer {
		return fmt.Errorf("value is not a pointer")
	}

	eleType := valType.Elem()

	switch eleType.Kind() {
	case reflect.Slice, reflect.Array:
		return m.getSlice(eleType, key, mapObj)
	case reflect.String:
		return m.getValue(eleType, key, mapObj)
	case reflect.Map:
		return m.getMap(eleType, key, mapObj)
	case reflect.Pointer, reflect.Struct, reflect.UnsafePointer, reflect.Chan, reflect.Interface,
		reflect.Uintptr, reflect.Complex64, reflect.Complex128, reflect.Func, reflect.Invalid:
		return fmt.Errorf("key %s is not a primitive type/slice/map/array", key)
	default:
		return m.getValue(eleType, key, mapObj)
	}
}

func (m *MapWrapper) getSlice(tgtVal reflect.Value, key string, srcObj any) error {

	srcType := reflect.ValueOf(srcObj)

	if tgtVal.Type().Kind() == reflect.Array {
		if tgtVal.Len() != srcType.Len() {
			return fmt.Errorf("key %s source length %d != target length %d", key, srcType.Len(), tgtVal.Len())
		}
	} else {
		tgtVal.Set(reflect.MakeSlice(tgtVal.Type(), srcType.Len(), srcType.Len()))
	}

	srcSlice := srcObj.([]any)

	for i := 0; i < len(srcSlice); i++ {
		to := tgtVal.Index(i)
		from := reflect.ValueOf(srcSlice[i])

		if from.Type().AssignableTo(to.Type()) {
			to.Set(from)
		} else if from.Type().ConvertibleTo(to.Type()) {
			to.Set(from.Convert(to.Type()))
		} else {
			return fmt.Errorf("key %s %d elem %v %s can't be convert to type %s", key, i, from, from.Type().Name(), to.Type().Name())
		}
	}
	return nil
}

func (m *MapWrapper) getValue(tgtVal reflect.Value, key string, srcObj any) error {

	srcType := reflect.ValueOf(srcObj)

	if srcType.Type().AssignableTo(tgtVal.Type()) {
		tgtVal.Set(srcType)
	} else if srcType.Type().ConvertibleTo(tgtVal.Type()) {
		tgtVal.Set(srcType.Convert(tgtVal.Type()))
	} else {
		return fmt.Errorf("key %s source %s != target %s", key, srcType.Type(), tgtVal.Type())
	}

	return nil
}

func (m *MapWrapper) getMap(tgtVal reflect.Value, key string, srcObj any) error {
	tgtKeyType := tgtVal.Type().Key().String()
	tgtValType := tgtVal.Type().Elem().String()

	srcType := reflect.ValueOf(srcObj)

	if srcType.Kind() != reflect.Map {
		return fmt.Errorf("key %s is not a map", key)
	}

	srcKeyType := srcType.Type().Key().String()
	srcValType := srcType.Type().Elem().String()

	if tgtKeyType != srcKeyType || srcKeyType != "string" {
		return fmt.Errorf("key %s source map[%s]%s != map[%s]%s", key, srcKeyType, srcValType, tgtKeyType, tgtValType)
	}

	if tgtValType == srcValType || srcType.Type().Elem().AssignableTo(tgtVal.Type().Elem()) || srcType.Type().Elem().ConvertibleTo(tgtVal.Type().Elem()) {
		tgtVal.Set(srcType)
	}

	mapData := srcObj.(map[string]any)
	tgtVal.Set(reflect.MakeMapWithSize(tgtVal.Type(), len(mapData)))
	for k, v := range mapData {
		vt := reflect.TypeOf(v)
		if vt.AssignableTo(tgtVal.Type().Elem()) {
			tgtVal.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
		} else if vt.ConvertibleTo(tgtVal.Type().Elem()) {
			tgtVal.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v).Convert(tgtVal.Type().Elem()))
		} else {
			return fmt.Errorf("key %s value %v type %s can't be convert to type %s", key, v, srcValType, tgtValType)
		}
	}

	return nil
}

func NewMapWrapper(data map[string]any) *MapWrapper {
	return &MapWrapper{
		data: data,
	}
}
