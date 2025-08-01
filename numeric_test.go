// Copyright (c) 2023 Chris Sewell
// Licensed under the MIT License

package main

import (
	"testing"
)

func TestNumericTypeComparison(t *testing.T) {
	// Load test files with different numeric representations
	file15, err := ReadAndValidateJSON("examples/example15.json", true)
	if err != nil {
		t.Fatalf("Failed to read examples/example15.json: %v", err)
	}

	file16, err := ReadAndValidateJSON("examples/example16.json", true)
	if err != nil {
		t.Fatalf("Failed to read examples/example16.json: %v", err)
	}

	// Test with ignore-numeric-type flag OFF
	diffsWithoutFlag := FindDifferences(file15.Data, file16.Data, "", false, false, false, false, false, false, nil, nil, 0)
	
	// Should find differences for all fields
	expectedDiffCount := 11 // All fields should be different
	if len(diffsWithoutFlag) != expectedDiffCount {
		t.Errorf("Expected %d differences with ignore-numeric-type OFF, got %d", expectedDiffCount, len(diffsWithoutFlag))
	}

	// Test with ignore-numeric-type flag ON
	diffsWithFlag := FindDifferences(file15.Data, file16.Data, "", false, false, true, false, false, false, nil, nil, 0)
	
	// Should find no differences
	if len(diffsWithFlag) != 0 {
		t.Errorf("Expected 0 differences with ignore-numeric-type ON, got %d: %v", len(diffsWithFlag), diffsWithFlag)
	}

	// Test specific numeric representations
	testCases := []struct {
		name   string
		val1   interface{}
		val2   interface{}
		equal  bool
	}{
		{"Integer vs String", 1, "1", true},
		{"Float vs String", 1.0, "1", true},
		{"Float vs String with decimal", 1.0, "1.0", true},
		{"Integer vs Float", 1, 1.0, true},
		{"Scientific notation", 1e2, "100", true},
		{"Scientific notation string", "1e2", 100, true},
		{"Zero representations", 0, "0", true},
		{"Zero float vs string", 0.0, "0", true},
		{"Negative numbers", -5, "-5", true},
		{"Different values", 42, "43", false},
		{"Non-numeric string", 42, "42a", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := compareNumericValues(tc.val1, tc.val2)
			if result != tc.equal {
				t.Errorf("compareNumericValues(%v, %v) = %v, want %v", 
					tc.val1, tc.val2, result, tc.equal)
			}
		})
	}
}