// Copyright (c) 2023 Chris Sewell
// Licensed under the MIT License

package main

import (
	"reflect"
	"testing"
)

func TestCaseSensitivity(t *testing.T) {
	// Load test files
	file1, err := ReadAndValidateJSON("example1.json", true)
	if err != nil {
		t.Fatalf("Failed to read example1.json: %v", err)
	}

	file6, err := ReadAndValidateJSON("example6.json", true)
	if err != nil {
		t.Fatalf("Failed to read example6.json: %v", err)
	}

	// Test key case sensitivity
	diffs := FindDifferences(file1.Data, file6.Data, "", false, false, false, false, false, false, nil)
	if len(diffs) == 0 {
		t.Error("Expected differences due to case-sensitive keys, but found none")
	}

	// Verify specific differences
	expectedKeyDiffs := []string{
		"name: key exists only in first file",
		"Name: key exists only in second file",
		"age: key exists only in first file",
		"Age: key exists only in second file",
		"address: key exists only in first file",
		"ADDRESS: key exists only in second file",
	}

	for _, expected := range expectedKeyDiffs {
		found := false
		for _, diff := range diffs {
			if contains(diff, expected) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find diff containing '%s', but didn't", expected)
		}
	}

	// Test case-insensitive comparison
	ignoreCaseDiffs := FindDifferences(file1.Data, file6.Data, "", true, false, false, false, false, false, nil)
	
	// Should only find differences in values, not in keys
	for _, diff := range ignoreCaseDiffs {
		for _, keyDiff := range expectedKeyDiffs {
			if contains(diff, keyDiff) {
				t.Errorf("Found unexpected key difference in case-insensitive mode: %s", diff)
			}
		}
	}

	// Test case-insensitive key-only comparison
	keyDiffs := FindDifferences(file1.Data, file6.Data, "", true, false, false, false, false, true, nil)
	if len(keyDiffs) > 0 {
		t.Errorf("Expected structures to be equal in case-insensitive mode, but found differences: %v", keyDiffs)
	}
}

func TestFindDifferences(t *testing.T) {
	// Load test files
	file1, err := ReadAndValidateJSON("example1.json", true)
	if err != nil {
		t.Fatalf("Failed to read example1.json: %v", err)
	}

	file2, err := ReadAndValidateJSON("example2.json", true)
	if err != nil {
		t.Fatalf("Failed to read example2.json: %v", err)
	}

	file3, err := ReadAndValidateJSON("example3.json", true)
	if err != nil {
		t.Fatalf("Failed to read example3.json: %v", err)
	}

	file4, err := ReadAndValidateJSON("example4.json", true)
	if err != nil {
		t.Fatalf("Failed to read example4.json: %v", err)
	}

	file5, err := ReadAndValidateJSON("example5.json", true)
	if err != nil {
		t.Fatalf("Failed to read example5.json: %v", err)
	}

	file7, err := ReadAndValidateJSON("example7.json", true)
	if err != nil {
		t.Fatalf("Failed to read example7.json: %v", err)
	}

	file8, err := ReadAndValidateJSON("example8.json", true)
	if err != nil {
		t.Fatalf("Failed to read example8.json: %v", err)
	}

	file9, err := ReadAndValidateJSON("example9.json", true)
	if err != nil {
		t.Fatalf("Failed to read example9.json: %v", err)
	}

	file10, err := ReadAndValidateJSON("example10.json", true)
	if err != nil {
		t.Fatalf("Failed to read example10.json: %v", err)
	}

	file11, err := ReadAndValidateJSON("example11.json", true)
	if err != nil {
		t.Fatalf("Failed to read example11.json: %v", err)
	}

	file12, err := ReadAndValidateJSON("example12.json", true)
	if err != nil {
		t.Fatalf("Failed to read example12.json: %v", err)
	}

	tests := []struct {
		name              string
		obj1              interface{}
		obj2              interface{}
		ignoreCase        bool
		ignoreCaseValues  bool
		ignoreNumericType bool
		ignoreBooleanType bool
		ignoreNullValues  bool
		regexMatches      map[string]string
		keysOnly          bool
		expectDiff        bool
		expectedDiffs     int
		expectedValues    []string
	}{
		// Full comparison tests
		{
			name:              "Identical files",
			obj1:              file1.Data,
			obj2:              file3.Data,
			ignoreCase:        false,
			ignoreCaseValues:  false,
			ignoreNumericType: false,
			ignoreBooleanType: false,
			ignoreNullValues:  false,
			regexMatches:      nil,
			keysOnly:          false,
			expectDiff:        false,
		},
		{
			name:              "Different values",
			obj1:              file1.Data,
			obj2:              file2.Data,
			ignoreCase:        false,
			ignoreCaseValues:  false,
			ignoreNumericType: false,
			ignoreBooleanType: false,
			ignoreNullValues:  false,
			regexMatches:      nil,
			keysOnly:          false,
			expectDiff:        true,
			expectedDiffs:     4,
			expectedValues: []string{
				"name: value mismatch",
				"age: value mismatch",
				"address.city: value mismatch",
				"hobbies[1]: value mismatch",
			},
		},
		{
			name:              "Different structure",
			obj1:              file1.Data,
			obj2:              file5.Data,
			ignoreCase:        false,
			ignoreCaseValues:  false,
			ignoreNumericType: false,
			ignoreBooleanType: false,
			ignoreNullValues:  false,
			regexMatches:      nil,
			keysOnly:          false,
			expectDiff:        true,
			expectedDiffs:     11,
			expectedValues: []string{
				"name: value mismatch",
				"age: value mismatch",
				"address.street: value mismatch",
				"address.city: value mismatch",
				"address.zip: key exists only in first file",
				"address.state: key exists only in second file",
				"address.country: key exists only in second file",
				"hobbies: array length mismatch",
				"hobbies[0]: value mismatch",
				"hobbies[1]: value mismatch",
				"email: key exists only in second file",
			},
		},
		// Case-insensitive value comparison test
		{
			name:              "Case-insensitive value comparison",
			obj1:              file1.Data,
			obj2:              file7.Data,
			ignoreCase:        false,
			ignoreCaseValues:  true,
			ignoreNumericType: false,
			ignoreBooleanType: false,
			ignoreNullValues:  false,
			regexMatches:      nil,
			keysOnly:          false,
			expectDiff:        false,
		},
		// Numeric type comparison test
		{
			name:              "Numeric type comparison",
			obj1:              file1.Data,
			obj2:              file8.Data,
			ignoreCase:        false,
			ignoreCaseValues:  false,
			ignoreNumericType: true,
			ignoreBooleanType: false,
			ignoreNullValues:  false,
			regexMatches:      nil,
			keysOnly:          false,
			expectDiff:        false,
		},
		// Boolean type comparison test
		{
			name:              "Boolean type comparison",
			obj1:              file1.Data,
			obj2:              file9.Data,
			ignoreCase:        false,
			ignoreCaseValues:  false,
			ignoreNumericType: false,
			ignoreBooleanType: true,
			ignoreNullValues:  false,
			regexMatches:      nil,
			keysOnly:          false,
			expectDiff:        false,
		},
		// Null value comparison test
		{
			name:              "Null value comparison",
			obj1:              file1.Data,
			obj2:              file10.Data,
			ignoreCase:        false,
			ignoreCaseValues:  false,
			ignoreNumericType: false,
			ignoreBooleanType: false,
			ignoreNullValues:  true,
			regexMatches:      nil,
			keysOnly:          false,
			expectDiff:        false,
		},
		// Keys-only comparison tests
		{
			name:              "Same structure different values",
			obj1:              file1.Data,
			obj2:              file4.Data,
			ignoreCase:        false,
			ignoreCaseValues:  false,
			ignoreNumericType: false,
			ignoreBooleanType: false,
			ignoreNullValues:  false,
			regexMatches:      nil,
			keysOnly:          true,
			expectDiff:        false,
		},
		{
			name:              "Different structure keys-only",
			obj1:              file1.Data,
			obj2:              file5.Data,
			ignoreCase:        false,
			ignoreCaseValues:  false,
			ignoreNumericType: false,
			ignoreBooleanType: false,
			ignoreNullValues:  false,
			regexMatches:      nil,
			keysOnly:          true,
			expectDiff:        true,
			expectedDiffs:     5,
			expectedValues: []string{
				"address.country: key exists only in second file",
				"address.state: key exists only in second file",
				"address.zip: key exists only in first file",
				"email: key exists only in second file",
				"hobbies: array length mismatch",
			},
		},
		// Regex match test
		{
			name:              "Regex match comparison",
			obj1:              file11.Data,
			obj2:              file12.Data,
			ignoreCase:        false,
			ignoreCaseValues:  false,
			ignoreNumericType: false,
			ignoreBooleanType: false,
			ignoreNullValues:  false,
			regexMatches:      map[string]string{"id": "[A-Z]+-\\d+-[A-Z]+"},
			keysOnly:          false,
			expectDiff:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diffs := FindDifferences(tt.obj1, tt.obj2, "", tt.ignoreCase, tt.ignoreCaseValues, tt.ignoreNumericType, tt.ignoreBooleanType, tt.ignoreNullValues, tt.keysOnly, tt.regexMatches)
			
			if (len(diffs) > 0) != tt.expectDiff {
				t.Errorf("Expected diff: %v, got: %v, differences: %v", tt.expectDiff, len(diffs) > 0, diffs)
			}

			if tt.expectedDiffs > 0 && len(diffs) != tt.expectedDiffs {
				t.Errorf("Expected %d diffs, got %d", tt.expectedDiffs, len(diffs))
			}

			if tt.expectedValues != nil {
				for _, expected := range tt.expectedValues {
					found := false
					for _, diff := range diffs {
						if contains(diff, expected) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected to find diff containing '%s', but didn't", expected)
					}
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return reflect.DeepEqual(s, substr) || len(s) >= len(substr) && s[:len(substr)] == substr
}