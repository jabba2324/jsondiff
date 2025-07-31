// Copyright (c) 2023 Chris Sewell
// Licensed under the MIT License

package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestCompareLines(t *testing.T) {
	tests := []struct {
		name     string
		str1     string
		str2     string
		expected string
	}{
		{
			name: "Different lines",
			str1: "line1\nline2\nline3",
			str2: "line1\nmodified\nline3",
			expected: "Line 2",  // We expect output to contain "Line 2"
		},
		{
			name: "Different number of lines",
			str1: "line1\nline2",
			str2: "line1\nline2\nline3",
			expected: "Line 3",  // We expect output to contain "Line 3"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			CompareLines(tt.str1, tt.str2)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if !strings.Contains(output, tt.expected) {
				t.Errorf("CompareLines() output = %v, should contain %v", output, tt.expected)
			}
		})
	}
}