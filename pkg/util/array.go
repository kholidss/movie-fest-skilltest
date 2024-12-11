package util

import (
	"reflect"
)

// InArray check if an element is exist in the array
func InArray(val interface{}, array interface{}) bool {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				return true
			}
		}
	}
	return false
}

// IsDuplicateArray check if any duplicate value in array
func IsDuplicateArray(val any) bool {
	slice, ok := val.([]string)
	if !ok {
		return false
	}

	seen := make(map[string]struct{})
	for _, v := range slice {
		if _, exists := seen[v]; exists {
			return true
		}
		seen[v] = struct{}{}
	}

	return false
}
