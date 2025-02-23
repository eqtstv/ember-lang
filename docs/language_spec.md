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
