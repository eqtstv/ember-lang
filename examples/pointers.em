// Basic pointer example
let mut x = 42;
let mut x_ptr = &x;

print(x); // 42
print(*x_ptr); // 42

*x_ptr = 24;

print(x) // 24
print(*x_ptr) // 24

print("--------------------------------");

// Array pointers
let mut arr = [1, 2, 3];
let p = &arr;

print(arr); // [1, 2, 3]
print(*p); // [1, 2, 3]

*p = [4, 5, 6];

print(arr) // [4, 5, 6]
print(*p) // [4, 5, 6]

print("--------------------------------");

// Pointers in Functions
let mut arr = [1, 2, 3];
let mut arr_ptr = &arr;

let change_first = fn(arr_ptr_arg) {
  (*arr_ptr_arg)[0] = 42;
};

change_first(arr_ptr);

print(arr); // [42, 2, 3]
print(*arr_ptr); // [42, 2, 3]

print("--------------------------------");

