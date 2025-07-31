// Copyright (c) 2023 Chris Sewell
// Licensed under the MIT License

package main

import (
	"fmt"
	"strings"
)

// CompareLines compares two strings line by line and prints differences
func CompareLines(str1, str2 string) {
	lines1 := strings.Split(str1, "\n")
	lines2 := strings.Split(str2, "\n")

	maxLines := len(lines1)
	if len(lines2) > maxLines {
		maxLines = len(lines2)
	}

	for i := 0; i < maxLines; i++ {
		var line1, line2 string
		if i < len(lines1) {
			line1 = lines1[i]
		}
		if i < len(lines2) {
			line2 = lines2[i]
		}

		if line1 != line2 {
			fmt.Printf("Line %d:\n  - %s\n  + %s\n", i+1, line1, line2)
		}
	}
}