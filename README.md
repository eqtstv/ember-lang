# Ember Programming Language

Ember is an interpreted programming language implemented in Go. It's designed to be simple yet powerful, featuring a C-like syntax.

## Features

- C-like syntax
- Variables with `let` keyword
- First-class functions
- Integer arithmetic
- Boolean operations
- Comparison operators
- Control flow (`if`/`else`)
- Return statements

## Example Code

```
let add = fn(a, b) {
    a + b;
};

add(5, 10);

let x = 10;

if (x > 5) {
    return true;
} else {
    return false;
}
```

## Project Structure

- `/lexer` - Tokenizes source code into tokens
- `/token` - Defines token types and structures
- `/repl` - Interactive shell for testing the language

## Running the REPL

To start the interactive shell:

```bash
go run main.go
```

You'll see a prompt `->` where you can type Ember code and see it tokenized.
