// Copyright 2025 Robert Snedegar
//
// Licensed under the Apache License, Version 2.0 (the License);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an AS IS BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package linearithmic

// IntroSort performs O(n log n) introspective sort (hybrid of quicksort, heapsort, and insertion sort).
// This demonstrates linearithmic time complexity by combining the best aspects of multiple algorithms:
// - Quicksort for average case O(n log n)
// - Heapsort for worst case O(n log n) when recursion depth exceeds 2*log n
// - Insertion sort for small arrays (< 16 elements) where it's faster due to low overhead
func IntroSort(arr []int) []int {
	// Create a copy to avoid modifying the original array
	result := make([]int, len(arr))
	copy(result, arr)

	// Calculate maximum recursion depth to prevent O(n²) worst case
	// If depth exceeds 2*log₂(n), switch to heapsort
	maxDepth := 2 * logBase2(len(result))
	introSortHelper(result, 0, len(result)-1, maxDepth)

	return result
}

// introSortHelper implements the hybrid sorting strategy with depth limiting.
// This ensures O(n log n) worst case by switching algorithms based on context.
func introSortHelper(arr []int, low, high, depthLimit int) {
	// Use insertion sort for small arrays (< 16 elements)
	// Insertion sort has low overhead and is faster for small datasets
	if high-low < 16 {
		insertionSortRange(arr, low, high)

		return
	}

	// If recursion depth exceeds limit, switch to heapsort
	// This prevents quicksort's O(n²) worst case
	if depthLimit == 0 {
		heapSortRange(arr, low, high)

		return
	}

	// Use quicksort for the main sorting work
	pivotIndex := partition(arr, low, high)
	// Recursively sort both partitions with decremented depth limit
	introSortHelper(arr, low, pivotIndex-1, depthLimit-1)
	introSortHelper(arr, pivotIndex+1, high, depthLimit-1)
}

func insertionSortRange(arr []int, low, high int) {
	for i := low + 1; i <= high; i++ {
		key := arr[i]
		j := i - 1

		for j >= low && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}

		arr[j+1] = key
	}
}

func heapSortRange(arr []int, low, high int) {
	// Build heap
	for i := (high-low)/2 - 1; i >= 0; i-- {
		heapifyRangeDown(arr, low, high, low+i)
	}

	// Extract elements
	for i := high; i > low; i-- {
		arr[low], arr[i] = arr[i], arr[low]
		heapifyRangeDown(arr, low, i-1, low)
	}
}

func heapifyRangeDown(arr []int, low, high, index int) {
	largest := index
	leftChild := 2*(index-low) + 1 + low
	rightChild := 2*(index-low) + 2 + low

	if leftChild <= high && arr[leftChild] > arr[largest] {
		largest = leftChild
	}

	if rightChild <= high && arr[rightChild] > arr[largest] {
		largest = rightChild
	}

	if largest != index {
		arr[index], arr[largest] = arr[largest], arr[index]
		heapifyRangeDown(arr, low, high, largest)
	}
}

// logBase2 calculates the floor of log₂(n).
// This is used to determine the maximum recursion depth for introsort.
func logBase2(n int) int {
	result := 0
	// Count how many times we can divide n by 2
	for n > 1 {
		n /= 2   // Divide by 2
		result++ // Increment counter
	}

	return result
}
