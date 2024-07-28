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

func GetObjAsMap[T any](val any) (map[string]T, error) {
	t, ok := val.(map[string]T)
	if !ok {
		log.Infof("unable to convert %v to type %T", val, t)
		return nil, fmt.Errorf("unable to convert %v to type %T", val, t)
	} else {
		return t, nil
	}
}
