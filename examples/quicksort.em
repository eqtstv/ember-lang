let swap = fn(arr, i, j) {
    let temp = arr[i];
    arr[i] = arr[j];
    arr[j] = temp;
    return arr;
};

let partition = fn(arr, low, high) {
    let pivot = arr[low];
    let i = low;
    let j = high;

    while (i < j) {
        while (i <= high - 1) {
            if (arr[i] > pivot) {
                break;
            }
            i = i + 1;
        }

        while (j >= low + 1) {
            if (arr[j] <= pivot) {
                break;
            }
            j = j - 1;
        }

        if (i < j) {
            arr = swap(arr, i, j);
        }
    }

    arr = swap(arr, low, j);
    return j;
};

let quickSort = fn(arr, low, high) {
    if (low < high) {
        // Get partition index
        let pi = partition(arr, low, high);

        // Recursively sort left and right partitions
        arr = quickSort(arr, low, pi - 1);
        arr = quickSort(arr, pi + 1, high);
    }
    return arr;
};

let numbers = [4, 2, 5, 3, 1];
let sorted = quickSort(numbers, 0, len(numbers) - 1);
print(sorted);  