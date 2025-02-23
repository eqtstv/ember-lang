# Ember Programming Language

Ember is an interpreted programming language implemented in Go, designed to be simple yet powerful with a focus on readability and expressiveness.

## Features

- C-like syntax with modern conveniences
- First-class functions and closures
- Integer arithmetic and boolean operations
- Variables with `let` keyword
- Control flow (`if`/`else`)
- Return statements
- Operator precedence parsing

## Quick Start

1. Clone the repository and navigate to the project directory:

```bash
cd ember_lang
```

2. Build and start the REPL:

```bash
make build
make run
```

3. Try some examples in the REPL:

```ember
⟶ let x = 10;
⟶ let y = 5;
⟶ let add = fn(a, b) { a + b; };
⟶ add(x, y);
15
```

## Project Structure

```
ember_lang/
├── cmd/ember/         # Command-line interface
├── ember_lang/        # Implementation
│   ├── lexer/         # Tokenization
│   ├── parser/        # Syntax analysis
│   ├── ast/           # Abstract Syntax Tree
│   ├── token/         # Token definitions
│   └── repl/          # Interactive shell
└── docs/              # Documentation
└── examples/          # Example code
```

## Development

Requirements:

- Go 1.21 or later
- Make

Common tasks:

```bash
make build    # Build the ember binary
make test     # Run tests
make lint     # Run linter
make run      # Run the REPL
```

## Language Syntax

Basic example:

```
let fib = fn(n) {
    if (n <= 1) {
        return n;
    }
    return fib(n - 1) + fib(n - 2);
};

let result = fib(10);  // Calculate 10th Fibonacci number
```
