# pstr - A Regex Engine in Go

<p align="center">
  <img src="https://img.shields.io/badge/build-passing-brightgreen" alt="Build Status" />
  <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License" />
  <img src="https://img.shields.io/badge/go-1.24%2B-blue" alt="Go Version" />
</p>

<p align="center">
  <b>pstr</b> is a simple regex engine written in Go from scratch for educational purposes. It can parse basic regular expressions, convert them into a Nondeterministic Finite Automaton (NFA), and perform matching on input strings.
</p>

---

## 🚀 Features

- **Basic Regex Parsing**: Supports literals, `( )` groups, `[ ]` character classes, and quantifiers like `*`, `+`, `?`, and `{m,n}`.
- **NFA Engine**: Converts parsed regex tokens into an NFA state machine.
- **String Matching**: Checks if an input string is valid according to the generated NFA.
- **Interactive CLI**: A simple command-line interface to test regex patterns in real-time.
- **Exposed API**: An API endpoint to check regex patterns programmatically.

## 🛠 Tech Stack

- [Go](https://golang.org/) (1.24+)
- [Fiber](https://gofiber.io/)

---

## 🏗️ Getting Started

### Prerequisites

- Go 1.24 or newer

### ▶️ Running the CLI

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/rubuy-74/pstr
    cd pstr
    ```

2.  **Run the interactive CLI:**
    You can use the provided shell script:
    ```bash
    ./run.sh
    ```
    Or run it directly with Go:
    ```bash
    go run cmd/pstr/main.go
    ```

3.  **Test a pattern:**
    The application will prompt you to enter a regex and then a string to check against it.

    *Example Interaction:*
    ```
    Enter regex
    > (a|b)*c
    Enter string to check
    > ababc
    Congratulations, the string is VALID
    ```

### ▶️ Running the API

1.  **Run the API server:**
    ```bash
    go run cmd/pstr/main.go
    ```
    The server will start on port `3000`.

2.  **Send a POST request:**
    You can use `curl` or any API client to send a `POST` request to the `/check` endpoint.

    *Example with `curl`:*
    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"regex": "(a|b)*c", "string": "ababc"}' http://localhost:3000/check
    ```

    *Expected Response (Success):*
    ```json
    {
        "valid": true
    }
    ```

    *Expected Response (Error):*
    ```json
    {
        "error": "failed to parse regex",
        "message": "missing left operand for | operator at position 0"
    }
    ```

## 🧪 Testing

> **Note**: This testing section was created using Cursor AI to provide comprehensive test coverage and reliability verification.

The project includes extensive test suites to ensure reliability and prevent crashes on edge cases. All tests verify that the regex engine handles invalid inputs gracefully instead of crashing.

### ▶️ Running Tests

#### **Run All Tests (Basic)**
```bash
go test ./...
```
- Runs all tests in all packages
- Shows only pass/fail status
- Fast and clean output

#### **Run All Tests (Verbose)**
```bash
go test ./... -v
```
- Shows detailed output for each test
- Great for debugging and seeing what's being tested
- Shows individual test case results

#### **Run Tests for Specific Package**
```bash
go test ./internal/parser/
go test ./internal/state_machine/
go test ./internal/
```

#### **Run Specific Test Functions**
```bash
go test ./... -run="TestEmptyInputValidation"
go test ./... -run="TestProcessRepeatEdgeCases"
go test ./... -run="TestPanicRecovery"
```

#### **Run Tests with Coverage**
```bash
go test ./... -cover
```
- Shows test coverage percentage
- Helps identify untested code

#### **Run Tests with Race Detection**
```bash
go test ./... -race
```
- Detects race conditions in concurrent code
- Important for production code

#### **Run Tests with Benchmarking**
```bash
go test ./... -bench=.
```
- Runs benchmark tests (if you have any)
- Measures performance

#### **Generate Coverage Report**
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```
- Generates detailed coverage report
- Opens HTML report in browser

### 📊 Test Coverage

The project includes comprehensive test suites covering:

- **Integration Tests**: Complete pipeline testing with edge cases
- **Parser Tests**: Original functionality + reliability tests
- **State Machine Tests**: NFA creation and validation
- **Reliability Tests**: Edge cases that previously caused crashes
- **Panic Recovery Tests**: Ensures no crashes on invalid inputs
- **Memory Safety Tests**: Validates safe memory access

**Total Test Cases**: 50+ individual test cases covering all reliability fixes.

### 🎯 Test Categories

#### **Input Validation Tests**
- Empty string handling
- Whitespace-only inputs
- Invalid operator usage

#### **Array Bounds Safety Tests**
- `processRepeat()` with no preceding tokens
- `processOr()` with missing operands
- `processBrackets()` with empty/invalid content
- `processGroup()` with empty/unclosed groups
- `ToNFA()` with empty token lists

#### **Type Safety Tests**
- `Token.ToNFA()` with invalid type assertions
- Safe type assertion handling
- Panic recovery verification

#### **Error Handling Tests**
- Proper error propagation
- Meaningful error messages
- Graceful failure handling

#### **Edge Case Tests**
- Malformed regex patterns
- Unclosed brackets/groups/ranges
- Invalid range syntax
- Complex combinations

### 💡 Recommended Usage

- **Daily Development**: `go test ./...`
- **Debugging Issues**: `go test ./... -v`
- **Before Commits**: `go test ./... -race -cover`
- **CI/CD Pipeline**: `go test ./... -v -cover`

## 📁 Project Structure

```text
├── cmd/
│   └── pstr/
│       └── main.go          # API endpoint and CLI entry point
├── internal/
│   ├── models/
│   │   ├── state/
│   │   │   └── state.go       # NFA state data structures
│   │   ├── token/
│   │   │   └── token.go       # Regex token data structures
│   │   └── token_type/
│   │       └── token_type.go  # Enum for token types
│   ├── parser/
│   │   ├── parser.go        # Regex string to token parsing
│   │   ├── parser_test.go   # Tests for the parser
│   │   └── reliability_test.go # Reliability and edge case tests
│   ├── state_machine/
│   │   ├── state_machine.go # Token to NFA conversion and matching logic
│   │   └── state_machine_test.go # State machine tests
│   ├── integration_test.go  # End-to-end integration tests
│   └── utils/
│       └── utils.go         # Utility functions
├── go.mod                   # Go module definition
├── run.sh                   # Script to run the CLI
└── TODO.md                  # Project goals and references
```

---

## 📜 License

This project is licensed under the MIT License.

---

<p align="center">
  <sub>Made with ❤️ by <a href="https://github.com/rubuy-74">rubuy-74</a></sub>
</p>
