# Ember Language Specification

## 1. Lexical Structure

### 1.1 Keywords

- `let`: Variable declaration
- `fn`: Function definition
- `if`, `else`: Control flow
- `return`: Return statement
- `true`, `false`: Boolean literals

### 1.2 Operators

- Arithmetic: `+`, `-`, `*`, `/`
- Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
- Logical: `!`

### 1.3 Delimiters

- Brackets: `()`, `{}`, `[]`
- Others: `,`, `;`

## 2. Syntax

### 2.1 Variable Declaration

```
let <identifier> = <expression>;
```

### 2.2 Function Definition

```
let <name> = fn(<params>) {
    <body>
};
```

Functions are first-class values and support closures:

```
let makeAdder = fn(x) {
    fn(y) { x + y };  // Captures x from outer scope
};
```

### 2.3 Control Flow

```
if (<condition>) {
    <consequence>
} else {
    <alternative>
}
```

## 3. Type System

Currently supported types:

- Integers: Whole numbers (`5`, `10`, `-3`)
- Booleans: `true` or `false`
- Functions: First-class closures
- Null: Represents absence of value

### 3.1 Type Coercion

- No implicit type coercion
- Boolean conditions must evaluate to boolean values
- Arithmetic operations require integer operands

## 4. Operator Precedence

From highest to lowest:

1. `()` - Grouping
2. Function calls
3. `*`, `/` - Multiplication, Division
4. `+`, `-` - Addition, Subtraction
5. `>`, `<`, `>=`, `<=` - Comparison
6. `==`, `!=` - Equality

## 5. Scoping

- Lexical scoping
- Functions create new scopes
- Closures capture their environment
- Variables must be declared before use
- No variable shadowing in the same scope

## 6. Error Handling

The interpreter reports:

- Syntax errors during parsing
- Type errors during evaluation
- Undefined variable references
- Invalid operator usage

## Built-in Functions

### Array Operations

- `len(array)`: Returns length of array or string
- `push(array, item)`: Adds item to array, returns new array
- `map(array, fn)`: Applies function to each element
- `reduce(array, fn, initial)`: Reduces array to single value

### Arithmetic Functions

- `add(x, y)`: Adds two integers
- `sub(x, y)`: Subtracts two integers
- `mul(x, y)`: Multiplies two integers
- `div(x, y)`: Divides two integers

### Utility Functions

- `print(...args)`: Prints arguments to stdout
- `len(arg)`: Returns length of strings or arrays

### Examples

```typescript
// Array operations
let nums = [1, 2, 3, 4, 5];
let doubled = map(nums, fn(x) { x * 2 });
let sum = reduce(nums, add, 0);

// Function composition
let numbers = [1, 2, 3, 4, 5];
let square = fn(x) { x * x };
let squares = map(numbers, square);
let sumOfSquares = reduce(squares, add, 0);

// Built-in arithmetic
let result = add(mul(5, 2), div(10, 2));  // 15
```

## Type System

Currently supported types:

- Integers: Whole numbers (`5`, `10`, `-3`)
- Booleans: `true` or `false`
- Arrays: Ordered collections (`[1, 2, 3]`)
- Functions: First-class closures
- Null: Represents absence of value

### Array Type

Arrays are:

- Zero-indexed
- Homogeneous (same type recommended)
- Immutable (operations return new arrays)
- Support built-in operations (map, reduce, push)

## Error Handling

The interpreter reports:

- Syntax errors during parsing
- Type errors during evaluation
- Undefined variable references
- Invalid operator usage
- Invalid function arguments
- Array operation errors

## Command Line Interface

The Ember interpreter can be used in two modes:

### REPL Mode

```bash
ember
```

Starts an interactive shell where you can type and evaluate Ember code directly.

### File Execution Mode

```bash
ember filename.em
```

Executes an Ember source file. Files must have the `.em` extension.

## File Format

Ember source files:

- Must have the `.em` extension
- Are UTF-8 encoded
- Use semicolons (;) as statement terminators
- Support both single-line and multi-line comments

Example:

```typescript
// Single line comment
/* Multi-line
   comment */

let main = fn() {
    // Your code here
};

main();
```

## Development Tools

### Debug Mode

The interpreter supports a debug mode that provides insight into the compilation and execution process:

```bash
DEBUG=1 ember program.em
```

This shows:

1. **Source Code**: The original program text
2. **Tokens**: The stream of tokens from lexical analysis
3. **AST**: A visual representation of the Abstract Syntax Tree
4. **Result**: The final evaluation result

### AST Visualization

The Abstract Syntax Tree shows the hierarchical structure of the program. Each node is color-coded:

- Purple: Keywords and control flow (let, if, return)
- Blue: Function-related nodes
- Cyan: Numbers and identifiers
- Orange: Strings and operators
- Green: Boolean values
- White: General expressions
- Gray: Tree structure

Example:

```
└── Let Statement
    ├── Identifier: x
    └── Infix: +
        ├── Integer: 5
        └── Integer: 10
```

The tree depth indicates:

1. Scope nesting (functions, blocks)
2. Expression composition
3. Operator precedence
4. Control flow structure
