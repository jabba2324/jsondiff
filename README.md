# JSON Diff Tool

A simple command-line tool to compare two JSON files and show their differences.

## Features

- Validates JSON files and pretty-prints them
- Performs exact match comparison between two JSON files
- Shows detailed differences including:
  - Missing/extra keys
  - Value mismatches
  - Array length differences
  - Line-by-line comparison
- Modular code organization with separate files for different functionalities
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
- `-no-detail`: Skip detailed line-by-line comparison
- `-quiet`: Only show if files differ via exit code (0 for identical, 1 for different)
- `-keys-only`: Only compare keys/structure, ignore values
- `-ignore-case`: Ignore case when comparing keys

## Example

```bash
./jsondiff example1.json example2.json
```

Output:
```
Validated JSON from example1.json
Validated JSON from example2.json
The JSON files are different.

Differences found:
name: value mismatch - John vs Jane
age: value mismatch - 30 vs 31
address.city: value mismatch - New York vs Boston
hobbies[1]: value mismatch - cycling vs swimming

Detailed line-by-line comparison:
Line 3:
  - "name": "John",
  + "name": "Jane",
Line 4:
  - "age": 30,
  + "age": 31,
Line 7:
  - "city": "New York",
  + "city": "Boston",
Line 12:
  - "cycling",
  + "swimming",
```

## Testing

To run the unit tests:

```bash
go test -v
```

The tests cover all major functionality including:
- Finding differences between JSON files
- Comparing only keys/structure
- Line-by-line comparison

## License

MIT