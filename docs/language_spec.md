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
