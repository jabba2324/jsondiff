# JSON Diff Tool

A simple command-line tool to compare two JSON files and show their differences.

## Features

- Validates JSON files and pretty-prints them
- Performs exact match comparison between two JSON files
- Shows detailed differences including:
  - Missing/extra keys
  - Value mismatches
  - Array length differences
  - Type mismatches
- Flexible comparison options:
  - Case-insensitive key comparison
  - Case-insensitive string value comparison
  - Type-agnostic numeric comparison (1 == "1" == "1.0" == 1.0)
  - Type-agnostic boolean comparison (true == "true")
  - Null value comparison ("Harry Potter" == null)
  - Regex pattern matching for specific keys
  - Levenshtein distance fuzzy matching for specific keys
- Comprehensive unit tests

## Installation

```bash
# Clone the repository
git clone https://github.com/chrissewell/jsondiff.git
cd jsondiff

# Build the tool
go build -o jsondiff
```

## Usage

```bash
./jsondiff [options] file1.json file2.json
```

### Options

- `-concise`: Show concise output (suppresses validation messages)
- `-quiet`: Only show if files differ via exit code (0 for identical, 1 for different)
- `-keys-only`: Only compare keys/structure, ignore values
- `-ignore-case`: Ignore case when comparing keys
- `-ignore-case-values`: Ignore case when comparing string values
- `-ignore-numeric-type`: Ignore numeric types (e.g., 1 == "1" == "1.0" == 1.0)
- `-ignore-boolean-type`: Ignore boolean types (e.g., true == "true")
- `-ignore-null`: Ignore null values (e.g., "Harry Potter" == null)
- `-regex-match`: Use regex matching on specific key (format: key:pattern), can be specified multiple times
- `-levenshtein-key`: Apply Levenshtein distance matching on specific key, can be specified multiple times
- `-levenshtein-threshold`: Maximum Levenshtein distance to consider strings as equal (default: 3)

## Examples

### Basic Comparison

```bash
./jsondiff examples/example1.json examples/example2.json
```

### Using Case-Insensitive Comparison

```bash
./jsondiff -ignore-case-values examples/example1.json examples/example7.json
```

This will ignore case differences in string values (e.g., "John" equals "JOHN").

Output:
```
Validated JSON from examples/example1.json
Validated JSON from examples/example7.json
The JSON files are identical.
```

### Using Levenshtein Distance Matching

```bash
./jsondiff -levenshtein-key name -levenshtein-key "education.university" -levenshtein-threshold 2 examples/example17.json examples/example18.json
```

This will use Levenshtein distance matching on the "name" and "education.university" fields, considering values equal if their Levenshtein distance is less than or equal to 2. This is useful for comparing strings that might have small typos or variations (e.g., "John Smith" vs "John Smyth" or "Massachusetts Institute of Technology" vs "Massachusets Institute of Technology").

Output:
```
Validated JSON from examples/example17.json
Validated JSON from examples/example18.json
The JSON files are different.

Differences found:
description: value mismatch
- Software Engineer with 5 years of experience
+ Software Engineer with 6 years of experience
location: value mismatch
- New York City
+ New York
```

### Using Regex Pattern Matching

```bash
./jsondiff -regex-match "id:[A-Z]+-\d+-[A-Z]+" examples/example11.json examples/example12.json
```

This will use regex pattern matching on the "id" field, considering values equal if they both match the pattern (e.g., "ABC-123-XYZ" and "DEF-456-UVW" both match the pattern "[A-Z]+-\d+-[A-Z]+").

You can specify multiple regex matches for different keys:

```bash
./jsondiff -regex-match "id:[A-Z]+-\d+-[A-Z]+" -regex-match "timestamp:\d{4}-\d{2}-\d{2}" examples/example11.json examples/example12.json
```

Output:
```
Validated JSON from examples/example11.json
Validated JSON from examples/example12.json
The JSON files are identical.
```

### Using Null Value Comparison

```bash
./jsondiff -ignore-null examples/example1.json examples/example10.json
```

This will treat null values as equal to any value (e.g., "Harry Potter" equals null).

Output:
```
Validated JSON from examples/example1.json
Validated JSON from examples/example10.json
The JSON files are identical.
```

### Using Type-Agnostic Comparison

```bash
./jsondiff -ignore-numeric-type -ignore-boolean-type examples/example1.json examples/example8.json
```

This will treat numeric and boolean values as equal regardless of their type (e.g., `1` equals `"1"` and `true` equals `"true"`).

Output:
```
Validated JSON from examples/example1.json
Validated JSON from examples/example8.json
The JSON files are identical.
```

### Output Example for Basic Comparison

```
Validated JSON from examples/example1.json
Validated JSON from examples/example2.json
The JSON files are different.

Differences found:
name: value mismatch
- John
+ Jane
age: value mismatch
- 30
+ 31
address.city: value mismatch
- New York
+ Boston
hobbies[1]: value mismatch
- cycling
+ swimming
```

## Testing

To run the unit tests:

```bash
go test -v
```

The tests cover all major functionality including:
- Finding differences between JSON files
- Comparing only keys/structure
- Type-agnostic comparisons
- Levenshtein distance matching
- Regex pattern matching

## License

MIT