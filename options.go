// Copyright (c) 2023 Chris Sewell
// Licensed under the MIT License

package main

// CompareOptions contains options for JSON comparison
type CompareOptions struct {
	IgnoreCase           bool              // If true, key comparisons will be case-insensitive
	IgnoreCaseValues     bool              // If true, string value comparisons will be case-insensitive
	IgnoreNumericType    bool              // If true, numeric types are compared by value, not type (e.g., 1 == "1" == "1.0")
	IgnoreBooleanType    bool              // If true, boolean types are compared by value, not type (e.g., true == "true")
	IgnoreNullValues     bool              // If true, null values are considered equal to any value
	KeysOnly             bool              // If true, only compare keys/structure, not values
	RegexMatches         map[string]string // Map of key paths to regex patterns for value matching
	LevenshteinKeys      map[string]bool   // Map of key paths to apply Levenshtein distance matching
	LevenshteinThreshold int               // Maximum Levenshtein distance to consider strings as equal
}