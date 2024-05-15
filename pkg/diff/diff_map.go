package diff

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// DiffJSON generete 2 map[string]interface from any input struct who use json tag
// Usually updatedData is nil version of existingData
// Nil value in updatedData wil be ignore (value in existing data considered not change)
func DiffJSON(existingData, updatedData interface{}) (before, after map[string]interface{}, err error) {
	json1, err := json.Marshal(existingData)
	if err != nil {
		return nil, nil, err
	}
	json2, err := json.Marshal(updatedData)
	if err != nil {
		return nil, nil, err
	}

	var data1, data2 map[string]interface{}
	if err := json.Unmarshal(json1, &data1); err != nil {
		return nil, nil, err
	}
	if err := json.Unmarshal(json2, &data2); err != nil {
		return nil, nil, err
	}

	before, after = FindDifference(data1, data2)
	return before, after, nil
}

// FindDiference generete before after map[string]interface from any map[string] interface input
// usually input2 is nil version of input1
// example :
// input1 address.street = "hawaii", address.postal = 8000
// input2 address.street = "himalayan", address.postal = nil
// result before & after will record address.street has changed and ignore address.postal
func FindDifference(input1, input2 map[string]interface{}) (before map[string]interface{}, after map[string]interface{}) {

	beforeTemp := map[string]interface{}{}
	afterTemp := map[string]interface{}{}

	findDifferences(input1, input2, beforeTemp, afterTemp, "")
	return beforeTemp, afterTemp
}

// findDifferences is internal helper using recursive for detect all different between 2 map
func findDifferences(data1, data2 map[string]interface{}, before, after map[string]interface{}, prefix string) {
	for key, val2 := range data2 {
		if val2 == nil {
			continue
		}
		if val1, ok := data1[key]; ok {
			if val1 != nil && reflect.TypeOf(val1).Kind() == reflect.Map {
				// Recursive call untuk nested maps
				findDifferences(val1.(map[string]interface{}), val2.(map[string]interface{}), before, after, prefix+key+".")
			} else if val1 != nil && reflect.TypeOf(val1).Kind() == reflect.Slice {
				// Handle slice of maps
				handleSliceDifferences(val1, val2, before, after, prefix+key)
			} else if !reflect.DeepEqual(val1, val2) {
				before[prefix+key] = val1
				after[prefix+key] = val2
			}
		} else {
			before[prefix+key] = ""
			after[prefix+key] = val2
		}
	}
}

func handleSliceDifferences(val1, val2 interface{}, before, after map[string]interface{}, prefix string) {
	slice1, ok1 := val1.([]interface{})
	slice2, ok2 := val2.([]interface{})

	if !ok1 || !ok2 {
		before[prefix] = val1
		after[prefix] = val2
		return
	}

	// Iterasi terhadap slice dan komparasi element
	for i := 0; i < len(slice2); i++ {
		prefixedKey := fmt.Sprintf("%s[%d]", prefix, i)
		if i < len(slice1) {
			if reflect.TypeOf(slice1[i]).Kind() == reflect.Map && reflect.TypeOf(slice2[i]).Kind() == reflect.Map {
				// Recursive call untuk nested maps di dalam slice / array
				findDifferences(slice1[i].(map[string]interface{}), slice2[i].(map[string]interface{}), before, after, prefixedKey+".")
			} else if !reflect.DeepEqual(slice1[i], slice2[i]) {
				before[prefixedKey] = slice1[i]
				after[prefixedKey] = slice2[i]
			}
		} else {
			// Element exists di slice2 tapi tidak di slice1
			before[prefixedKey] = ""
			after[prefixedKey] = slice2[i]
		}
	}

	// Handle yang ada di slice1 tetapi tidak ada di slice2
	for i := len(slice2); i < len(slice1); i++ {
		prefixedKey := fmt.Sprintf("%s[%d]", prefix, i)
		before[prefixedKey] = slice1[i]
		after[prefixedKey] = ""
	}
}

// ConvertToNestedMap transforms a map with dot notation keys into a nested map
func ConvertToNestedMap(input map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range input {
		parts := strings.Split(key, ".")
		n := len(parts)
		current := result

		for i, part := range parts {
			// Check for array notation like "locations[0]"
			if openIndex := strings.Index(part, "["); openIndex != -1 {
				closeIndex := strings.Index(part, "]")
				if closeIndex > openIndex {
					// Extract array index
					index := part[openIndex+1 : closeIndex]
					arrayKey := part[:openIndex]

					// Ensure the array exists
					if _, exists := current[arrayKey]; !exists {
						current[arrayKey] = make([]interface{}, 0)
					}

					// Ensure the array is of the correct length
					array := current[arrayKey].([]interface{})
					for len(array) <= int(index[0]-'0') {
						array = append(array, make(map[string]interface{}))
					}
					current[arrayKey] = array

					// Descend into the correct element
					if i == n-1 {
						array[int(index[0]-'0')] = value
					} else {
						current = array[int(index[0]-'0')].(map[string]interface{})
					}
				}
			} else {
				if i == n-1 {
					current[part] = value
				} else {
					if _, exists := current[part]; !exists {
						current[part] = make(map[string]interface{})
					}
					current = current[part].(map[string]interface{})
				}
			}
		}
	}
	return result
}
