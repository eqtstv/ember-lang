let partition = fn(arr_ptr, low, high) {
    let pivot = (*arr_ptr)[high];
    
    let mut i = low - 1;
    
    for (let j = low; j < high; j++) {
        if ((*arr_ptr)[j] < pivot) {
            i = i + 1;
            
            let temp = (*arr_ptr)[i];
            (*arr_ptr)[i] = (*arr_ptr)[j];
            (*arr_ptr)[j] = temp;
        }
    }
    
    let temp = (*arr_ptr)[i + 1];
    (*arr_ptr)[i + 1] = (*arr_ptr)[high];
    (*arr_ptr)[high] = temp;
    
    return i + 1;
};

let quicksort = fn(arr_ptr, low, high) {
    if (low < high) {
        let pi = partition(arr_ptr, low, high);
        
        quicksort(arr_ptr, low, pi - 1);
        quicksort(arr_ptr, pi + 1, high);
    }
    
    return *arr_ptr;
};

let array = [10, 7, 8, 9, 1, 5, 3, 2, 6, 4];
let ptr = &array;

print("Original array:", array);
quicksort(ptr, 0, len(array) - 1);
print("Sorted array:", array);
