# Ember Language Specification

## 1. Lexical Structure

### 1.1 Keywords

- `let`: Variable declaration
- `mut`: Mutability modifier
- `fn`: Function definition
- `if`, `else`: Control flow
- `return`: Return statement
- `true`, `false`: Boolean literals
- `while`, `for`: Loop constructs

### 1.2 Operators

- Arithmetic: `+`, `-`, `*`, `/`
- Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
- Logical: `!`
- Assignment: `=`

### 1.3 Delimiters

- Brackets: `()`, `{}`, `[]`
- Others: `,`, `;`

## 2. Syntax

### 2.1 Variable Declaration

```
let [mut] <identifier> = <expression>;
```

Variables are immutable by default. The `mut` keyword makes a variable mutable, allowing reassignment.

Examples:

```
let x = 5;           // Immutable variable
let mut y = 10;      // Mutable variable

y = 20;              // Valid - y is mutable
x = 30;              // Error - x is immutable
```

### 2.2 Assignment Expression

```
<identifier> = <expression>
```

Assignment is only valid for mutable variables. Attempting to assign to an immutable variable results in a runtime error.

Example:

```
let mut counter = 0;
counter = counter + 1;  // Valid - counter is mutable

let value = 5;
value = 10;            // Error - value is immutable
```

### 2.3 Function Definition

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

Function parameters are immutable by default. They cannot be reassigned within the function body.

Example:

```
let increment = fn(x) {
    x = x + 1;  // Error - function parameters are immutable
    return x;
};
```

### 2.4 Control Flow

#### If Statements

The language supports standard control flow if statements:

```
if (<condition>) {
    <consequence>
} else {
    <alternative>
}
```

#### Loops

Two types of loops are supported:

```
while (<condition>) {
    <body>
}

for (let [mut] <var> = <init>; <condition>; <increment>) {
    <body>
}
```

Example:

```typescript
// While loop with mutable counter
let mut i = 0;
while (i < 5) {
  i = i + 1;
}

// Equivalent for loop
for (let i = 0; i < 5; i++) {
  // loop body
}
```

### 2.7 Pointers

Ember supports pointers for referencing variables. Pointers are created using the `&` operator and dereferenced using the `*` operator.

#### 2.7.1 Creating Pointers

The address-of operator `&` creates a pointer to a variable:

```
let x = 10;
let p = &x;  // p is a pointer to x
```

#### 2.7.2 Dereferencing Pointers

The dereference operator `*` accesses the value a pointer points to:

```
let x = 10;
let p = &x;
let y = *p;  // y is 10 (the value of x)
```

#### 2.7.3 Modifying Values Through Pointers

Pointers can be used to modify the value of the variable they point to:

```
let mut x = 10;
let p = &x;
*p = 20;     // Changes x to 20
```

Note that the target variable must be mutable for modification through a pointer to work.

#### 2.7.4 Pointers to Complex Data Structures

Pointers can reference arrays and objects:

```
let mut arr = [1, 2, 3];
let p = &arr;
(*p)[0] = 42;  // Changes first element of arr to 42

let mut obj = {"name": "Alice", "age": 30};
let p = &obj;
(*p)["age"] = 31;  // Changes age to 31
```

#### 2.7.5 Pointers in Functions

Pointers are useful for modifying variables from within functions:

```
let modifyValue = fn(ptr) {
  *ptr = *ptr * 2;
};

let mut x = 10;
modifyValue(&x);  // x is now 20
```

#### 2.7.6 Pointer Safety

- Null pointer checks are recommended before dereferencing
- Pointer arithmetic is not supported to prevent memory safety issues
- Pointers cannot be created to arbitrary memory addresses

#### 2.7.7 Limitations

- Pointers must point to valid variables
- Dereferencing a null pointer causes a runtime error
- Pointer arithmetic (adding/subtracting from pointer addresses) is not supported

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
7. `=` - Assignment

## 5. Scoping and Mutability

- Lexical scoping
- Functions create new scopes
- Closures capture their environment
- Variables must be declared before use
- Variables are immutable by default
- The `mut` keyword makes variables mutable
- Function parameters are immutable
- Variables can be shadowed in nested scopes

### 5.1 Variable Shadowing

Variables can be shadowed in nested scopes. The inner variable is distinct from the outer variable, even if they have the same name.

Example:

```
let x = 5;
if (true) {
    let x = 10;  // Shadows outer x
    print(x);    // Prints 10
}
print(x);        // Prints 5
```

### 5.2 Mutability Rules

1. Variables are immutable by default
2. The `mut` keyword makes a variable mutable
3. Function parameters are immutable
4. Assignment is only valid for mutable variables
5. Mutability is checked at runtime

## 6. Error Handling

The interpreter reports:

- Syntax errors during parsing
- Type errors during evaluation
- Undefined variable references
- Invalid operator usage
- Mutability violations (attempting to assign to immutable variables)

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

// Mutability example
let mut counter = 0;
while (counter < 5) {
    counter = counter + 1;
    print(counter);
}
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
- Mutability violations (attempting to assign to immutable variables)

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

## Comments

Ember supports single-line comments using the `//` syntax. Everything after `//` until the end of the line is considered a comment and is ignored by the interpreter.
