// Copyright (c) 2023 Chris Sewell
// Licensed under the MIT License

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

// stringSliceFlag is a custom flag type that allows multiple values
type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	// Define flags
	concisePtr := flag.Bool("concise", false, "Show concise output")
	quietPtr := flag.Bool("quiet", false, "Only show if files differ, no details")
	outputJSONPtr := flag.String("output-json", "", "Write differences to a JSON file")
	keysOnlyPtr := flag.Bool("keys-only", false, "Only compare keys, ignore values")
	ignoreCasePtr := flag.Bool("ignore-case", false, "Ignore case when comparing keys")
	ignoreCaseValuesPtr := flag.Bool("ignore-case-values", false, "Ignore case when comparing string values")
	ignoreNumericTypePtr := flag.Bool("ignore-numeric-type", false, "Ignore numeric types (e.g., 1 == \"1\" == \"1.0\" == 1.0)")
	ignoreBooleanTypePtr := flag.Bool("ignore-boolean-type", false, "Ignore boolean types (e.g., true == \"true\")")
	ignoreNullValuesPtr := flag.Bool("ignore-null", false, "Ignore null values (e.g., \"Harry Potter\" == null)")
	var regexMatchList stringSliceFlag
	flag.Var(&regexMatchList, "regex-match", "Use regex matching on specific key (format: key:pattern), can be specified multiple times")
	var levenshteinKeyList stringSliceFlag
	flag.Var(&levenshteinKeyList, "levenshtein-key", "Apply Levenshtein distance matching on specific key, can be specified multiple times")
	levenshteinThresholdPtr := flag.Int("levenshtein-threshold", 3, "Maximum Levenshtein distance to consider strings as equal (default: 3)")

	// Parse flags
	flag.Parse()

	// Check if we have exactly two arguments after flags
	args := flag.Args()
	if len(args) != 2 {
		fmt.Println("Usage: jsondiff [options] <file1.json> <file2.json>")
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	file1Path := args[0]
	file2Path := args[1]

	// Read and validate first JSON file
	jsonFile1, err := ReadAndValidateJSON(file1Path, *concisePtr)
	if err != nil {
		fmt.Printf("Error with first file: %v\n", err)
		os.Exit(1)
	}

	// Read and validate second JSON file
	jsonFile2, err := ReadAndValidateJSON(file2Path, *concisePtr)
	if err != nil {
		fmt.Printf("Error with second file: %v\n", err)
		os.Exit(1)
	}

	// Parse regex match options
	regexMatches := make(map[string]string)
	for _, regexMatch := range regexMatchList {
		// Split the regex match option into key and pattern
		parts := strings.SplitN(regexMatch, ":", 2)
		if len(parts) == 2 {
			key := parts[0]
			pattern := parts[1]

			// Check for duplicate keys
			if _, exists := regexMatches[key]; exists {
				fmt.Printf("Warning: Duplicate regex match key '%s'. Only the last pattern will be used.\n", key)
			}

			regexMatches[key] = pattern
		} else {
			fmt.Println("Invalid regex match format. Expected format: key:pattern")
			os.Exit(1)
		}
	}

	// Parse Levenshtein keys
	levenshteinKeys := make(map[string]bool)
	for _, key := range levenshteinKeyList {
		levenshteinKeys[key] = true
	}

	// Get differences based on options
	differences := FindDifferences(jsonFile1.Data, jsonFile2.Data, "", *ignoreCasePtr, *ignoreCaseValuesPtr, *ignoreNumericTypePtr, *ignoreBooleanTypePtr, *ignoreNullValuesPtr, *keysOnlyPtr, regexMatches, levenshteinKeys, *levenshteinThresholdPtr)
	
	// Write differences to JSON file if requested
	if *outputJSONPtr != "" {
		outputJSON, err := json.MarshalIndent(differences, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling differences to JSON: %v\n", err)
			os.Exit(1)
		}
		
		err = os.WriteFile(*outputJSONPtr, outputJSON, 0644)
		if err != nil {
			fmt.Printf("Error writing differences to file: %v\n", err)
			os.Exit(1)
		}
		
		if !*quietPtr {
			fmt.Printf("Differences written to %s\n", *outputJSONPtr)
		}
	}

	// Check if files are identical
	if len(differences) == 0 {
		if !*quietPtr {
			fmt.Println("The JSON files are identical.")
		}
		os.Exit(0)
	} else {
		if !*quietPtr {
			fmt.Println("The JSON files are different.")

			// Show the differences
			fmt.Println("\nDifferences found:")
			for _, diff := range differences {
				switch diff.Type {
				case ValueMismatch:
					fmt.Printf("%s: value mismatch\n- %v\n+ %v\n", diff.Path, diff.Value1, diff.Value2)
				case KeyOnlyInFirst:
					fmt.Printf("%s: key exists only in first file\n", diff.Path)
				case KeyOnlyInSecond:
					fmt.Printf("%s: key exists only in second file\n", diff.Path)
				case ArrayLength:
					fmt.Printf("%s: array length mismatch\n- %v\n+ %v\n", diff.Path, diff.Value1, diff.Value2)
				case TypeMismatch:
					fmt.Printf("%s: type mismatch\n- %v\n+ %v\n", diff.Path, diff.Value1, diff.Value2)
				}
			}
		}
		os.Exit(1) // Exit with non-zero status if files differ
	}
}
