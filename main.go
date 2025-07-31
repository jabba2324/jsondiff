// Copyright (c) 2023 Chris Sewell
// Licensed under the MIT License

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Define flags
	concisePtr := flag.Bool("concise", false, "Show concise output")
	noDetailPtr := flag.Bool("no-detail", false, "Skip detailed line-by-line comparison")
	quietPtr := flag.Bool("quiet", false, "Only show if files differ, no details")
	keysOnlyPtr := flag.Bool("keys-only", false, "Only compare keys, ignore values")
	ignoreCasePtr := flag.Bool("ignore-case", false, "Ignore case when comparing keys")
	
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

	// Get differences based on options
	differences := FindDifferences(jsonFile1.Data, jsonFile2.Data, "", *ignoreCasePtr, *keysOnlyPtr)
	
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
				fmt.Println(diff)
			}
			
			// Compare pretty-printed JSON strings line by line for visual diff
			if !*noDetailPtr {
				fmt.Println("\nDetailed line-by-line comparison:")
				CompareLines(jsonFile1.PrettyJSON, jsonFile2.PrettyJSON)
			}
		}
		os.Exit(1) // Exit with non-zero status if files differ
	}
}