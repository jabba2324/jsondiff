package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"sort"
	"strings"
)

// JSONFile represents a parsed JSON file
type JSONFile struct {
	Data       interface{}
	PrettyJSON string
}

// CompareOptions contains options for JSON comparison
type CompareOptions struct {
	IgnoreCase bool // If true, key comparisons will be case-insensitive
	KeysOnly   bool // If true, only compare keys/structure, not values
}

// ReadAndValidateJSON reads a JSON file, validates it, and returns the parsed object and pretty-printed string
func ReadAndValidateJSON(filePath string, concise bool) (*JSONFile, error) {
	// Read file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Parse JSON
	var jsonObj interface{}
	err = json.Unmarshal(data, &jsonObj)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}

	// Pretty print (for validation)
	prettyJSON, err := json.MarshalIndent(jsonObj, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to pretty print: %v", err)
	}

	if !concise {
		fmt.Printf("Validated JSON from %s\n", filePath)
	}
	
	return &JSONFile{
		Data:       jsonObj,
		PrettyJSON: string(prettyJSON),
	}, nil
}

// FindDifferences recursively compares two JSON objects and returns a list of differences
// Options control whether to ignore case and whether to compare only keys
func FindDifferences(obj1, obj2 interface{}, path string, ignoreCase, keysOnly bool) []string {
	return findDifferencesWithOptions(obj1, obj2, path, CompareOptions{IgnoreCase: ignoreCase, KeysOnly: keysOnly})
}

// isComplex checks if a value is a complex type (map or array)
func isComplex(val interface{}) bool {
	switch val.(type) {
	case map[string]interface{}, []interface{}:
		return true
	default:
		return false
	}
}

// findDifferencesWithOptions is the internal implementation that handles all comparison options
func findDifferencesWithOptions(obj1, obj2 interface{}, path string, options CompareOptions) []string {
	differences := []string{}

	// If types are different, that's a difference
	type1 := reflect.TypeOf(obj1)
	type2 := reflect.TypeOf(obj2)
	if type1 != type2 {
		return append(differences, fmt.Sprintf("%s: type mismatch - %v vs %v", path, type1, type2))
	}

	// Handle different types
	switch obj1.(type) {
	case map[string]interface{}:
		// Compare maps
		map1 := obj1.(map[string]interface{})
		map2 := obj2.(map[string]interface{})

		// Get all keys from both maps
		allKeys := make(map[string]bool)
		
		// If case-insensitive, create case-insensitive maps for lookup
		var lookupMap1, lookupMap2 map[string]interface{}
		var keyMap1, keyMap2 map[string]string
		
		if options.IgnoreCase {
			lookupMap1 = make(map[string]interface{})
			lookupMap2 = make(map[string]interface{})
			keyMap1 = make(map[string]string)
			keyMap2 = make(map[string]string)
			
			// Create case-insensitive lookup maps
			for key, val := range map1 {
				lKey := strings.ToLower(key)
				lookupMap1[lKey] = val
				keyMap1[lKey] = key
				allKeys[lKey] = true
			}
			
			for key, val := range map2 {
				lKey := strings.ToLower(key)
				lookupMap2[lKey] = val
				keyMap2[lKey] = key
				allKeys[lKey] = true
			}
		} else {
			// Standard case-sensitive comparison
			for k := range map1 {
				allKeys[k] = true
			}
			for k := range map2 {
				allKeys[k] = true
			}
		}

		// Sort keys for consistent output
		keys := make([]string, 0, len(allKeys))
		for k := range allKeys {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// Check each key
		for _, key := range keys {
			var newPath, originalKey1, originalKey2 string
			var val1, val2 interface{}
			var ok1, ok2 bool
			
			if options.IgnoreCase {
				// For case-insensitive, key is already lowercase
				originalKey1, ok1 = keyMap1[key]
				originalKey2, ok2 = keyMap2[key]
				
				if ok1 {
					val1 = lookupMap1[key]
					newPath = originalKey1
				} else if ok2 {
					newPath = originalKey2
				} else {
					newPath = key // Shouldn't happen, but just in case
				}
				
				if ok2 {
					val2 = lookupMap2[key]
				}
			} else {
				// Standard case-sensitive comparison
				newPath = key
				val1, ok1 = map1[key]
				val2, ok2 = map2[key]
			}
			
			if path != "" {
				newPath = path + "." + newPath
			}

			if !ok1 {
				differences = append(differences, fmt.Sprintf("%s: key exists only in second file", newPath))
			} else if !ok2 {
				differences = append(differences, fmt.Sprintf("%s: key exists only in first file", newPath))
			} else if !options.KeysOnly && !reflect.DeepEqual(val1, val2) {
				// Only compare values if not in keys-only mode
				if isComplex(val1) {
					// Recursively compare nested structures
					differences = append(differences, findDifferencesWithOptions(val1, val2, newPath, options)...)
				} else {
					// For primitive types, just compare values
					differences = append(differences, fmt.Sprintf("%s: value mismatch - %v vs %v", newPath, val1, val2))
				}
			} else if options.KeysOnly && isComplex(val1) {
				// In keys-only mode, still check structure of nested objects
				differences = append(differences, findDifferencesWithOptions(val1, val2, newPath, options)...)
			}
		}

	case []interface{}:
		// Compare arrays
		arr1 := obj1.([]interface{})
		arr2 := obj2.([]interface{})

		// Check array lengths
		if len(arr1) != len(arr2) {
			differences = append(differences, fmt.Sprintf("%s: array length mismatch - %d vs %d", path, len(arr1), len(arr2)))
		}

		// Compare array elements
		minLen := len(arr1)
		if len(arr2) < minLen {
			minLen = len(arr2)
		}

		for i := 0; i < minLen; i++ {
			newPath := fmt.Sprintf("%s[%d]", path, i)
			val1 := arr1[i]
			val2 := arr2[i]
			
			if !options.KeysOnly && !reflect.DeepEqual(val1, val2) {
				// Only compare values if not in keys-only mode
				if isComplex(val1) {
					differences = append(differences, findDifferencesWithOptions(val1, val2, newPath, options)...)
				} else {
					differences = append(differences, fmt.Sprintf("%s: value mismatch - %v vs %v", newPath, val1, val2))
				}
			} else if options.KeysOnly && isComplex(val1) {
				// In keys-only mode, still check structure of nested objects
				differences = append(differences, findDifferencesWithOptions(val1, val2, newPath, options)...)
			}
		}

	default:
		// For primitive types, just compare values if not in keys-only mode
		if !options.KeysOnly && !reflect.DeepEqual(obj1, obj2) {
			differences = append(differences, fmt.Sprintf("%s: value mismatch - %v vs %v", path, obj1, obj2))
		}
	}

	return differences
}

