let fibonacci = fn(n) {
    if (n <= 1) {
        return n;
    }
    return fibonacci(n - 1) + fibonacci(n - 2);
};

let result = fibonacci(10);

return result;
