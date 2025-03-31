// Variables and arithmetic
let age = 25;
let temperature = 18 + 5;
let isHot = temperature > 20;

// Mutability
let x = 5;         // Immutable by default
// x = 10;         // Error: Cannot assign to immutable variable: x

let mut y = 5;     // Explicitly mutable
y = 10;            // Works fine

// For loop
for (let i = 0; i < 5; i++) {
    let forIndex = i;
}
print(forIndex);

// While loop
let mut i = 0;
while (i < 5) {
    i = i + 1;
}
print(i);

// Arrays
let numbers = [1, 2, 3, 4, 5];
let doubled = map(numbers, fn(x) { x * 2 });  // [2, 4, 6, 8, 10]
let sum = reduce(numbers, add, 0);  // 15
let numbers = push(numbers, 6);  // [1, 2, 3, 4, 5, 6]

// Hashes
let person = {
    "name": "John",
    "age": 30,
    "city": "New York"
};
let name = person["name"];  // "John"
print(name);

// Functions
let greet = fn(name) {
    return "Hello, " + name + "!";
};
let greeting = greet("John");  // Returns "Hello, John!"
print(greeting);

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
let mut result = addFive(10);   // Returns 15
print(result);
result = addFive(20);   // Returns 25
print(result);

// Recursive functions
let fib = fn(n) {
    if (n <= 1) {
        return n;
    }
    return fib(n - 1) + fib(n - 2);
};

let mut result = fib(10);  // Calculate 10th Fibonacci number
print(result); // 55