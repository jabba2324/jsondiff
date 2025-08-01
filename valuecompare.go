// Copyright (c) 2023 Chris Sewell
// Licensed under the MIT License

package main

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/agnivade/levenshtein"
)

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

// matchesRegex checks if both values match the given regex pattern
// Returns true if both values are strings and match the pattern
func matchesRegex(val1, val2 interface{}, pattern string) (bool, error) {
	// Check if both values are strings
	str1, isStr1 := val1.(string)
	str2, isStr2 := val2.(string)
	
	if !isStr1 || !isStr2 {
		return false, nil // Not comparing strings
	}
	
	// Compile the regex pattern
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}
	
	// Check if both strings match the pattern
	return re.MatchString(str1) && re.MatchString(str2), nil
}

// compareLevenshtein checks if two values are similar using Levenshtein distance
// Returns true if both values are strings and their Levenshtein distance is within the threshold
func compareLevenshtein(val1, val2 interface{}, threshold int) bool {
	// Check if both values are strings
	str1, isStr1 := val1.(string)
	str2, isStr2 := val2.(string)
	
	if !isStr1 || !isStr2 {
		return false // Not comparing strings
	}
	
	// Calculate Levenshtein distance
	distance := levenshtein.ComputeDistance(str1, str2)
	
	// Return true if distance is within threshold
	return distance <= threshold
}