# Ember Programming Language

Ember is an interpreted programming language implemented in Go, designed to be simple yet powerful with a focus on readability and expressiveness.

## Features

- C-like syntax with modern conveniences
- First-class functions and closures
- Dynamic typing with integers, booleans, arrays, hashes, and functions
- Lexical scoping and proper closures
- Built-in integer arithmetic and boolean operations
- Array operations (map, reduce, push)
- Built-in functions for common operations
- Variables with `let` keyword
- Control flow (`if`/`else`)
- Return statements
- Operator precedence parsing
- REPL with error reporting

## Language Syntax

Basic examples:

```typescript
// Variables and arithmetic
let age = 25;
let temperature = 18 + 5;
let isHot = temperature > 20;

// Arrays
let numbers = [1, 2, 3, 4, 5];
let doubled = numbers.map(fn(x) { x * 2 });  // [2, 4, 6, 8, 10]
let sum = numbers.reduce(fn(acc, x) { acc + x }, 0);  // 15
numbers.push(6);  // [1, 2, 3, 4, 5, 6]

// Hashes
let person = {
    "name": "John",
    "age": 30,
    "city": "New York"
};
let name = person["name"];  // "John"
person["country"] = "USA";  // Add new key-value pair

// Functions
let greet = fn(name) {
    return "Hello, " + name + "!";
};
greet("John");  // Returns "Hello, John!"

// Conditionals
let max = fn(a, b) {
    if (a > b) {
        return a;
    } else {
        return b;
    }
};

// Functions and closures
let makeAdder = fn(x) {
    return fn(y) {
        return x + y;
    };
};

let addFive = makeAdder(5);
addFive(10);   // Returns 15
addFive(20);   // Returns 25

// Recursive functions
let fib = fn(n) {
    if (n <= 1) {
        return n;
    }
    return fib(n - 1) + fib(n - 2);
};

let result = fib(10);  // Calculate 10th Fibonacci number
```

## Getting Started

You can use Ember in two ways:

### Option 1: Quick Start (REPL only)

If you just want to try Ember without installing:

```bash
cd ember_lang
make run
```

This will start the REPL (Read-Eval-Print Loop) interactive shell.

### Option 2: System Installation

For full installation that allows running Ember files from anywhere:

1. Clone the repository:

```bash
cd ember_lang
```

2. Install the Ember binary:

```bash
make install
```

This will install the `ember` command to `/usr/local/bin/ember`.

## Usage

There are two ways to use Ember:

### 1. REPL (Interactive Mode)

Start the REPL by running:

```bash
ember
```

You'll see:

```
Ember Programming Language v0.0.1 (prototype)
Type "help" for more information.
⟶
```

### 2. File Execution

Run Ember files (with .em extension):

```bash
ember fibonacci.em
```

### Example Program

Create a file `hello.em`:

```typescript
let greet = fn(name) {
    return "Hello, " + name + "!";
};

print(greet("World"));
```

Run it:

```bash
ember hello.em
```

## Project Structure

```
ember_lang/
├── cmd/ember/         # Command-line interface
├── ember_lang/        # Implementation
│   ├── lexer/        # Tokenization
│   ├── parser/       # Syntax analysis
│   ├── ast/          # Abstract Syntax Tree
│   ├── token/        # Token definitions
│   ├── object/       # Runtime object system
│   ├── evaluator/    # Expression evaluation
│   └── repl/         # Interactive shell
└── docs/             # Documentation
└── examples/         # Example code
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
make run      # Start REPL (during development)
make install  # Install ember to /usr/local/bin
```

### Debugging

You can enable debug output by setting the DEBUG environment variable:

```bash
DEBUG=1 ember fibonacci.em
```

This will show:

1. Source code
2. Token stream (lexical analysis)
3. Abstract Syntax Tree (AST)
4. Final result

The AST visualization shows the hierarchical structure of your code:

```
└── Program
    └── Let Statement
        ├── Identifier: fibonacci
        └── Function: fn(n)
            └── Block Statement
                ├── If Expression
                │   ├── Infix: <=
                │   │   ├── Identifier: n
                │   │   └── Integer: 1
```

The tree depth indicates:

- Code blocks and scopes
- Expression nesting
- Operator precedence
- Control flow structure
