// Copyright (c) 2023 Chris Sewell
// Licensed under the MIT License

package main

import (
	"testing"
)

func TestLevenshteinComparison(t *testing.T) {
	// Load test files with similar but not identical text values
	file17, err := ReadAndValidateJSON("examples/example17.json", true)
	if err != nil {
		t.Fatalf("Failed to read examples/example17.json: %v", err)
	}

	file18, err := ReadAndValidateJSON("examples/example18.json", true)
	if err != nil {
		t.Fatalf("Failed to read examples/example18.json: %v", err)
	}

	// Test without Levenshtein distance
	diffsWithoutLevenshtein := FindDifferences(file17.Data, file18.Data, "", false, false, false, false, false, false, nil, nil, 0)
	
	// Should find differences for name, description, location, and university
	expectedDiffCount := 4
	if len(diffsWithoutLevenshtein) != expectedDiffCount {
		t.Errorf("Expected %d differences without Levenshtein, got %d", expectedDiffCount, len(diffsWithoutLevenshtein))
	}

	// Test with Levenshtein distance on name and university
	levenshteinKeys := map[string]bool{
		"name": true,
		"education.university": true,
	}
	
	diffsWithLevenshtein := FindDifferences(file17.Data, file18.Data, "", false, false, false, false, false, false, nil, levenshteinKeys, 2)
	
	// Should find differences only for description and location
	expectedDiffCount = 2
	if len(diffsWithLevenshtein) != expectedDiffCount {
		t.Errorf("Expected %d differences with Levenshtein, got %d: %v", expectedDiffCount, len(diffsWithLevenshtein), diffsWithLevenshtein)
	}

	// Test specific Levenshtein comparisons
	testCases := []struct {
		name      string
		val1      interface{}
		val2      interface{}
		threshold int
		equal     bool
	}{
		{"Small edit distance", "hello", "hallo", 1, true},
		{"Medium edit distance", "hello", "hola", 3, true},
		{"Large edit distance", "hello", "goodbye", 5, false},
		{"Empty string", "", "hello", 3, false},
		{"Case difference", "Hello", "hello", 1, true},
		{"Transposition", "hello", "helol", 2, true},
		{"Insertion", "hello", "hellox", 1, true},
		{"Deletion", "hello", "hell", 1, true},
		{"Non-string values", 42, "42", 0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := compareLevenshtein(tc.val1, tc.val2, tc.threshold)
			if result != tc.equal {
				t.Errorf("compareLevenshtein(%v, %v, %d) = %v, want %v", 
					tc.val1, tc.val2, tc.threshold, result, tc.equal)
			}
		})
	}
}