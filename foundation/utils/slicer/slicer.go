package slicer

import (
	"fmt"
	"strings"
)

func In[T comparable](value T, list []T) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

func RemoveFrom[T comparable](value T, list []T) []T {
	copyList := make([]T, len(list))
	copy(copyList, list)

	for i, v := range copyList {
		if v == value {
			copyList = RemoveAtIndex(copyList, i)
		}
	}
	return copyList
}

func RemoveAtIndex[T comparable](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func RemoveAtIndexNotSorted[T comparable](slice []T, index int) []T {
	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func ToStringSlice(i interface{}) ([]string, error) {
	var a []string

	switch v := i.(type) {
	case []interface{}:
		for _, u := range v {
			str, ok := u.(string)
			if !ok {
				return a, fmt.Errorf("unable to cast %#v of type %T to string", u, u)
			}
			a = append(a, str)
		}
		return a, nil
	case []string:
		return v, nil
	case string:
		return strings.Fields(v), nil
	default:
		return a, fmt.Errorf("unable to cast %#v of type %T to []string", i, i)
	}
}
