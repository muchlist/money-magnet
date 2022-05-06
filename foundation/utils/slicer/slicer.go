package slicer

import (
	"fmt"
	"strings"
)

// In do return true if value(1) available on list(2)
func In[T comparable](value T, list []T) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

// RemoveFrom do return new list with removed value(1) on list(2)
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

// RemoveAtIndex do return new list with removed item index(2) on list(1)
func RemoveAtIndex[T comparable](list []T, index int) []T {
	return append(list[:index], list[index+1:]...)
}

// RemoveAtIndexNotSorted do return new list with removed item index(2) on list(1)
// new list is not sorted as input but faster execution
func RemoveAtIndexNotSorted[T comparable](list []T, index int) []T {
	list[index] = list[len(list)-1]
	return list[:len(list)-1]
}

// ToStringSlice do convert interface{} to slice string if not error
func ToStringSlice(i any) ([]string, error) {
	var a []string

	switch v := i.(type) {
	case []any:
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
