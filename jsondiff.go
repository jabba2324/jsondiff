// Copyright (c) 2023 Chris Sewell
// Licensed under the MIT License

package main

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// DiffType represents the type of difference between JSON objects
type DiffType int

// Enum values for DiffType
const (
	ValueMismatch DiffType = iota
	KeyOnlyInFirst
	KeyOnlyInSecond
	ArrayLength
	TypeMismatch
)

// String returns the string representation of a DiffType
func (dt DiffType) String() string {
	switch dt {
	case ValueMismatch:
		return "value_mismatch"
	case KeyOnlyInFirst:
		return "key_only_in_first"
	case KeyOnlyInSecond:
		return "key_only_in_second"
	case ArrayLength:
		return "array_length"
	case TypeMismatch:
		return "type_mismatch"
	default:
		return "unknown"
	}
}

// Diff represents a difference between two JSON objects
type Diff struct {
	Path   string      // Path to the key where the difference was found
	Type   DiffType    // Type of difference
	Value1 interface{} // Value from the first object
	Value2 interface{} // Value from the second object
}

// FindDifferences recursively compares two JSON objects and returns a list of differences
// Options control whether to ignore case and whether to compare only keys
func FindDifferences(obj1, obj2 interface{}, path string, ignoreCase, ignoreCaseValues, ignoreNumericType, ignoreBooleanType, ignoreNullValues, keysOnly bool, regexMatches map[string]string, levenshteinKeys map[string]bool, levenshteinThreshold int) []Diff {
	return findDifferencesWithOptions(obj1, obj2, path, CompareOptions{
		IgnoreCase:           ignoreCase,
		IgnoreCaseValues:     ignoreCaseValues,
		IgnoreNumericType:    ignoreNumericType,
		IgnoreBooleanType:    ignoreBooleanType,
		IgnoreNullValues:     ignoreNullValues,
		KeysOnly:             keysOnly,
		RegexMatches:         regexMatches,
		LevenshteinKeys:      levenshteinKeys,
		LevenshteinThreshold: levenshteinThreshold,
	})
}

// compareValues compares two values with all the special handling options
// Returns true if the values are considered equal according to the options
func compareValues(val1, val2 interface{}, path string, options CompareOptions) bool {
	// Special handling for strings when IgnoreCaseValues is true
	if options.IgnoreCaseValues && !options.KeysOnly {
		str1, isStr1 := val1.(string)
		str2, isStr2 := val2.(string)
		if isStr1 && isStr2 && strings.EqualFold(str1, str2) {
			// Strings are equal when ignoring case
			return true
		}
	}

	// Special handling for null values
	if options.IgnoreNullValues && !options.KeysOnly {
		if val1 == nil || val2 == nil {
			// If either value is null, consider them equal
			return true
		}
	}

	// Special handling for regex matching
	if !options.KeysOnly && len(options.RegexMatches) > 0 {
		// Check if this key path has a regex pattern
		if pattern, ok := options.RegexMatches[path]; ok {
			// Check if both values match the pattern
			matches, err := matchesRegex(val1, val2, pattern)
			if err == nil && matches {
				// Both values match the pattern, consider them equal
				return true
			}
		}
	}

	// Special handling for Levenshtein distance
	if !options.KeysOnly && len(options.LevenshteinKeys) > 0 && options.LevenshteinThreshold > 0 {
		// Check if this key path should use Levenshtein distance
		if _, ok := options.LevenshteinKeys[path]; ok {
			// Check if strings are similar using Levenshtein distance
			if compareLevenshtein(val1, val2, options.LevenshteinThreshold) {
				// Strings are similar enough, consider them equal
				return true
			}
		}
	}

	// Special handling for boolean types
	if options.IgnoreBooleanType && !options.KeysOnly {
		if equal, ok := compareBooleanValues(val1, val2); ok && equal {
			// Values are equal when compared as booleans
			return true
		}
	}

	// Special handling for numeric types
	if options.IgnoreNumericType && !options.KeysOnly {
		if compareNumericValues(val1, val2) {
			// Values are equal when compared as numbers
			return true
		}
	}

	// Standard comparison
	return reflect.DeepEqual(val1, val2)
}

// findDifferencesWithOptions is the internal implementation that handles all comparison options
func findDifferencesWithOptions(obj1, obj2 interface{}, path string, options CompareOptions) []Diff {
	differences := []Diff{}

	// If types are different, that's a difference
	type1 := reflect.TypeOf(obj1)
	type2 := reflect.TypeOf(obj2)
	if type1 != type2 {
		differences = append(differences, Diff{
			Path:   path,
			Type:   TypeMismatch,
			Value1: type1,
			Value2: type2,
		})
		return differences
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
				differences = append(differences, Diff{
					Path:   newPath,
					Type:   KeyOnlyInSecond,
					Value1: nil,
					Value2: val2,
				})
			} else if !ok2 {
				differences = append(differences, Diff{
					Path:   newPath,
					Type:   KeyOnlyInFirst,
					Value1: val1,
					Value2: nil,
				})
			} else {
				// Compare values using all the special handling options
				if options.KeysOnly {
					// In keys-only mode, only check structure of complex objects
					if isComplex(val1) {
						differences = append(differences, findDifferencesWithOptions(val1, val2, newPath, options)...)
					}
				} else {
					// Check if values are equal according to the options
					if !compareValues(val1, val2, newPath, options) {
						if isComplex(val1) {
							// Recursively compare nested structures
							differences = append(differences, findDifferencesWithOptions(val1, val2, newPath, options)...)
						} else {
							// For primitive types, just compare values
							differences = append(differences, Diff{
								Path:   newPath,
								Type:   ValueMismatch,
								Value1: val1,
								Value2: val2,
							})
						}
					}
				}
			}
		}

	case []interface{}:
		// Compare arrays
		arr1 := obj1.([]interface{})
		arr2 := obj2.([]interface{})

		// Check array lengths
		if len(arr1) != len(arr2) {
			differences = append(differences, Diff{
				Path:   path,
				Type:   ArrayLength,
				Value1: len(arr1),
				Value2: len(arr2),
			})
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

			// Compare values using all the special handling options
			if options.KeysOnly {
				// In keys-only mode, only check structure of complex objects
				if isComplex(val1) {
					differences = append(differences, findDifferencesWithOptions(val1, val2, newPath, options)...)
				}
			} else {
				// Check if values are equal according to the options
				if !compareValues(val1, val2, newPath, options) {
					if isComplex(val1) {
						// Recursively compare nested structures
						differences = append(differences, findDifferencesWithOptions(val1, val2, newPath, options)...)
					} else {
						// For primitive types, just compare values
						differences = append(differences, Diff{
							Path:   newPath,
							Type:   ValueMismatch,
							Value1: val1,
							Value2: val2,
						})
					}
				}
			}
		}

	default:
		// For primitive types, just compare values if not in keys-only mode
		if !options.KeysOnly && !compareValues(obj1, obj2, path, options) {
			differences = append(differences, Diff{
				Path:   path,
				Type:   ValueMismatch,
				Value1: obj1,
				Value2: obj2,
			})
		}
	}

	return differences
}