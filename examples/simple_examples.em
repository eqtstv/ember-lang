// Simple arithmetic
5 + 5 * 10

// Variables
let x = 5;
let y = 10;
print(x + y);

// Functions
let add = fn(a, b) {
    return a + b;
};
print(add(x, y));

// Loops
let mut i = 0;
while (i < 10) {
    i = i++;
}
print(i);

// For loop
for (let i = 0; i < 10; i++) {
    let x = i;
}
print(x);

// Arrays
let mut numbers = [1, 2, 3];
numbers[0] = 10;
print(numbers);

// Maps
let mut mapping = {"a": 1, "b": 2, "c": 3};
mapping["a"] = 10;
print(mapping);