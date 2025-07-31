// Copyright (c) 2023 Chris Sewell
// Licensed under the MIT License

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// JSONFile represents a parsed JSON file
type JSONFile struct {
	Data       interface{}
	PrettyJSON string
}

// CompareOptions contains options for JSON comparison
type CompareOptions struct {
	IgnoreCase        bool // If true, key comparisons will be case-insensitive
	IgnoreCaseValues  bool // If true, string value comparisons will be case-insensitive
	IgnoreNumericType bool // If true, numeric types are compared by value, not type (e.g., 1 == "1" == "1.0")
	IgnoreBooleanType bool // If true, boolean types are compared by value, not type (e.g., true == "true")
	IgnoreNullValues  bool // If true, null values are considered equal to any value
	KeysOnly          bool // If true, only compare keys/structure, not values
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
func FindDifferences(obj1, obj2 interface{}, path string, ignoreCase, ignoreCaseValues, ignoreNumericType, ignoreBooleanType, ignoreNullValues, keysOnly bool) []string {
	return findDifferencesWithOptions(obj1, obj2, path, CompareOptions{
		IgnoreCase:        ignoreCase,
		IgnoreCaseValues:  ignoreCaseValues,
		IgnoreNumericType: ignoreNumericType,
		IgnoreBooleanType: ignoreBooleanType,
		IgnoreNullValues:  ignoreNullValues,
		KeysOnly:          keysOnly,
	})
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

// compareBooleanValues compares two values as booleans, ignoring their original types
// Returns true if both values can be converted to booleans and are equal
func compareBooleanValues(val1, val2 interface{}) (bool, bool) {
	// Extract boolean values
	bool1, isBool1 := val1.(bool)
	boolStr1, isBoolStr1 := val1.(string)
	bool2, isBool2 := val2.(bool)
	boolStr2, isBoolStr2 := val2.(string)
	
	// Check if we're comparing boolean values
	if !isBool1 && !isBoolStr1 && !isBool2 && !isBoolStr2 {
		// Not comparing booleans
		return false, false
	}
	
	// Convert to boolean values
	var b1, b2 bool
	var ok1, ok2 bool
	
	if isBool1 {
		b1 = bool1
		ok1 = true
	} else if isBoolStr1 {
		lowerStr := strings.ToLower(boolStr1)
		if lowerStr == "true" {
			b1 = true
			ok1 = true
		} else if lowerStr == "false" {
			b1 = false
			ok1 = true
		}
	}
	
	if isBool2 {
		b2 = bool2
		ok2 = true
	} else if isBoolStr2 {
		lowerStr := strings.ToLower(boolStr2)
		if lowerStr == "true" {
			b2 = true
			ok2 = true
		} else if lowerStr == "false" {
			b2 = false
			ok2 = true
		}
	}
	
	// If both values could be converted to booleans, compare them
	if ok1 && ok2 {
		return b1 == b2, true
	}
	
	// Couldn't convert both values to booleans
	return false, false
}

// convertToFloat64 attempts to convert a value to float64
// Returns the converted value and a boolean indicating success
func convertToFloat64(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	case string:
		// Try to parse string as number
		f, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return f, true
		}
	}
	return 0, false
}

// compareNumericValues compares two values as numbers, ignoring their original types
// Returns true if both values can be converted to numbers and are equal
func compareNumericValues(val1, val2 interface{}) bool {
	// Try to convert both values to float64
	num1, ok1 := convertToFloat64(val1)
	num2, ok2 := convertToFloat64(val2)

	// Compare the numeric values if both conversions succeeded
	if ok1 && ok2 {
		return num1 == num2
	}
	
	// Values couldn't be compared as numbers
	return false
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
			} else {
				// Special handling for strings when IgnoreCaseValues is true
				if options.IgnoreCaseValues && !options.KeysOnly {
					str1, isStr1 := val1.(string)
					str2, isStr2 := val2.(string)
					if isStr1 && isStr2 && strings.EqualFold(str1, str2) {
						// Strings are equal when ignoring case
						continue
					}
				}
				
				// Special handling for null values
				if options.IgnoreNullValues && !options.KeysOnly {
					if val1 == nil || val2 == nil {
						// If either value is null, consider them equal
						continue
					}
				}
				
				// Special handling for boolean types
				if options.IgnoreBooleanType && !options.KeysOnly {
					if equal, ok := compareBooleanValues(val1, val2); ok && equal {
						// Values are equal when compared as booleans
						continue
					}
				}
				
				// Special handling for numeric types
				if options.IgnoreNumericType && !options.KeysOnly {
					if compareNumericValues(val1, val2) {
						// Values are equal when compared as numbers
						continue
					}
				}
				
				if !options.KeysOnly && !reflect.DeepEqual(val1, val2) {
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
			
			// Special handling for strings when IgnoreCaseValues is true
			if options.IgnoreCaseValues && !options.KeysOnly {
				str1, isStr1 := val1.(string)
				str2, isStr2 := val2.(string)
				if isStr1 && isStr2 && strings.EqualFold(str1, str2) {
					// Strings are equal when ignoring case, continue to next element
					continue
				}
			}
			
			// Special handling for null values
			if options.IgnoreNullValues && !options.KeysOnly {
				if val1 == nil || val2 == nil {
					// If either value is null, consider them equal
					continue
				}
			}
			
			// Special handling for boolean types
			if options.IgnoreBooleanType && !options.KeysOnly {
				if equal, ok := compareBooleanValues(val1, val2); ok && equal {
					// Values are equal when compared as booleans
					continue
				}
			}
			
			// Special handling for numeric types
			if options.IgnoreNumericType && !options.KeysOnly {
				if compareNumericValues(val1, val2) {
					// Values are equal when compared as numbers
					continue
				}
			}
			
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
		if !options.KeysOnly {
			// Special handling for null values
			if options.IgnoreNullValues {
				if obj1 == nil || obj2 == nil {
					// If either value is null, consider them equal
					return differences
				}
			}
			
			// Special handling for strings when IgnoreCaseValues is true
			if options.IgnoreCaseValues {
				str1, isStr1 := obj1.(string)
				str2, isStr2 := obj2.(string)
				if isStr1 && isStr2 && strings.EqualFold(str1, str2) {
					// Strings are equal when ignoring case
					return differences
				}
			}
			
			// Special handling for boolean types
			if options.IgnoreBooleanType {
				if equal, ok := compareBooleanValues(obj1, obj2); ok && equal {
					// Values are equal when compared as booleans
					return differences
				}
			}
			
			// Special handling for numeric types
			if options.IgnoreNumericType {
				if compareNumericValues(obj1, obj2) {
					// Values are equal when compared as numbers
					return differences
				}
			}
			
			// Standard comparison for non-strings or when not ignoring case
			if !reflect.DeepEqual(obj1, obj2) {
				differences = append(differences, fmt.Sprintf("%s: value mismatch - %v vs %v", path, obj1, obj2))
			}
		}
	}

	return differences
}

