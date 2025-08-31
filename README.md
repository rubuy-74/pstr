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

## 🛠 Tech Stack

- [Go](https://golang.org/) (1.24+)

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
    go run cmd/main.go
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

## 📁 Project Structure

```text
├── cmd/
│   └── main.go              # Application entry point (Interactive CLI)
├── internal/
│   ├── parser/
│   │   ├── parser.go        # Regex string to token parsing
│   │   ├── state_machine.go # Token to NFA conversion
│   │   ├── matching.go      # NFA-based string matching logic
│   │   └── parser_test.go   # Tests for the parser
│   └── state_machine/
│       └── state_machine.go # NFA data structures and implementation
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
