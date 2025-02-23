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
- Comparison: `==`, `!=`, `<`, `>`
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

- Integers
- Booleans
- Functions
- Null

## 4. Operator Precedence

From highest to lowest:

1. `()` - Grouping
2. `*`, `/` - Multiplication, Division
3. `+`, `-` - Addition, Subtraction
4. `>`, `<` - Comparison
5. `==`, `!=` - Equality
