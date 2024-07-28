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

func (m *MapWrapper) Get(key string) (val any, err error) {
	var ok bool
	val, ok = m.data[key]
	if !ok {
		err = fmt.Errorf("key %s not found in map", key)
	}
	return
}

func (m *MapWrapper) GetValueType(key string) (reflect.Type, error) {
	val, ok := m.data[key]
	if !ok {
		return nil, fmt.Errorf("key %s not found in map", key)
	}

	return reflect.TypeOf(val), nil
}

func (m *MapWrapper) GetString(key string) (string, error) {
	val, ok := m.data[key]
	if !ok {
		return "", fmt.Errorf("key %s not found in map", key)
	}
	strVal, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("key %s is not a string", key)
	}

	return strVal, nil
}

func (m *MapWrapper) GetBool(key string) (bool, error) {
	val, ok := m.data[key]
	if !ok {
		return false, fmt.Errorf("key %s not found in map", key)
	}
	boolVal, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("key %s is not a bool", key)
	}

	return boolVal, nil
}

func (m *MapWrapper) GetSlice(key string) ([]any, error) {
	val, ok := m.data[key]
	if !ok {
		return nil, fmt.Errorf("key %s not found in map", key)
	}
	listVal, ok := val.([]any)
	if !ok {
		return nil, fmt.Errorf("key %s is not a list", key)
	}
	return listVal, nil
}

func (m *MapWrapper) GetStringSlice(key string) ([]string, error) {
	val, ok := m.data[key]
	if !ok {
		return nil, fmt.Errorf("key %s not found in map", key)
	}
	listVal, ok := val.([]string)
	if !ok {
		return nil, fmt.Errorf("key %s is not a list", key)
	}
	return listVal, nil
}

func (m *MapWrapper) GetMap(key string) (map[string]any, error) {
	val, ok := m.data[key]
	if !ok {
		return nil, fmt.Errorf("key %s not found in map", key)
	}
	mapVal, ok := val.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("key %s is not a map", key)
	}
	return mapVal, nil
}

func (m *MapWrapper) GetStringMap(key string) (map[string]string, error) {
	val, ok := m.data[key]
	if !ok {
		return nil, fmt.Errorf("key %s not found in map", key)
	}
	mapVal, ok := val.(map[string]string)
	if !ok {
		return nil, fmt.Errorf("key %s is not a map", key)
	}
	return mapVal, nil
}

func (m *MapWrapper) GetMapWrapper(key string) (mapWrapper *MapWrapper, err error) {
	var data map[string]any
	data, err = m.GetMap(key)
	if err != nil {
		return
	}

	mapWrapper = NewMapWrapper(data)
	return
}

func NewMapWrapper(data map[string]any) *MapWrapper {
	return &MapWrapper{
		data: data,
	}
}
