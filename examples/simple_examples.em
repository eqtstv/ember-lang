5 + 5 * 10

let x = 5;
let y = 10;
print(x + y);

let add = fn(a, b) {
    return a + b;
};
print(add(x, y));

let mut i = 0;
while (i < 10) {
    i = i++;
}
print(i);


for (let i = 0; i < 10; i++) {
    let x = i;
}
print(x);


let mut numbers = [1, 2, 3];
numbers[0] = 10;
print(numbers);

let mut mapping = {"a": 1, "b": 2, "c": 3};
mapping["a"] = 10;
print(mapping);