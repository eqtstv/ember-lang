let mut x = 42;
let mut ptr_to_x = &x;

print(x);
let value_at_ptr = *ptr_to_x;  // Dereference and store in a variable first
print(value_at_ptr);           // Then print the variable

*ptr_to_x = 9;                 // This dereference for assignment works fine

print(x);
let new_value_at_ptr = *ptr_to_x;  // Dereference again after modification
print(new_value_at_ptr);           // Print the new value

return 0;