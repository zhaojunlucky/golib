package collection

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func CreateTwoDimArray(row, column int) [][]int {
	var twoDimArr = make([][]int, row)

	for i := range twoDimArr {
		twoDimArr[i] = make([]int, column)
	}

	return twoDimArr
}

// GetObjAsSlice return a slice if the given any type object is a slice.
// if the given object is not a slice, then return error.
func GetObjAsSlice[T any](val any) ([]T, error) {
	t, ok := val.(T)
	if !ok {
		log.Infof("unable to convert %v to type %T", val, t)
	} else {
		return []T{t}, nil
	}

	tArr, ok := val.([]T)
	if !ok {
		log.Infof("unable to convert %v to type %T", val, tArr)
		return nil, fmt.Errorf("unable to convert %v to type %T", val, tArr)
	} else {
		return tArr, nil
	}
}

// GetObjAsMap return a map if the given any type object is a map.
// if the given object is not a map, then return error.
func GetObjAsMap[T any](val any) (map[string]T, error) {
	t, ok := val.(map[string]T)
	if !ok {
		log.Infof("unable to convert %v to type %T", val, t)
		return nil, fmt.Errorf("unable to convert %v to type %T", val, t)
	} else {
		return t, nil
	}
}
