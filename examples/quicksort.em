let partition = fn(arr, low, high) {
    let mut arr = arr;
    let pivot = arr[high];
    
    let mut i = low - 1;
    
    for (let j = low; j < high; j++) {
        if (arr[j] < pivot) {
            i = i + 1;
            
            let temp = arr[i];
            arr[i] = arr[j];
            arr[j] = temp;
        }
    }
    
    let temp = arr[i + 1];
    arr[i + 1] = arr[high];
    arr[high] = temp;
    
    return i + 1;
};

let quicksort = fn(arr, low, high) {
    if (low < high) {
        let pi = partition(arr, low, high);
        
        quicksort(arr, low, pi - 1);
        quicksort(arr, pi + 1, high);
    }
    
    return arr;
};

let array = [10, 7, 8, 9, 1, 5, 3, 2, 6, 4];

print("Original array:", array);
quicksort(array, 0, len(array) - 1);
print("Sorted array:", array);
