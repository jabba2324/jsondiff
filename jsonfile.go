// Copyright (c) 2023 Chris Sewell
// Licensed under the MIT License

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// JSONFile represents a parsed JSON file
type JSONFile struct {
	Data interface{}
}

// ReadAndValidateJSON reads a JSON file, validates it, and returns the parsed object
func ReadAndValidateJSON(filePath string, concise bool) (*JSONFile, error) {
	// Read file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Parse JSON
	var jsonObj interface{}
	err = json.Unmarshal(data, &jsonObj)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}

	if !concise {
		fmt.Printf("Validated JSON from %s\n", filePath)
	}
	
	return &JSONFile{
		Data: jsonObj,
	}, nil
}