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

// QuickSort performs O(n log n) quick sort (average case).
// This demonstrates linearithmic time complexity through divide-and-conquer:
// we partition around a pivot (O(n)) and recursively sort two sub-arrays,
// creating on average log n levels of recursion. Total: O(n log n) average case.
func QuickSort(arr []int) []int {
	// Create a copy to avoid modifying the original array
	result := make([]int, len(arr))
	copy(result, arr)
	// Perform in-place sorting on the copy
	quickSortHelper(result, 0, len(result)-1)

	return result
}

// quickSortHelper recursively sorts array sections using divide-and-conquer.
// The recursion depth is on average O(log n), and each level does O(n) work.
func quickSortHelper(arr []int, low, high int) {
	// Base case: if section has 1 or 0 elements, it's already sorted
	if low < high {
		// Partition: rearrange elements around pivot - O(n) operation
		pivotIndex := partition(arr, low, high)
		// Conquer: recursively sort elements before pivot
		quickSortHelper(arr, low, pivotIndex-1)
		// Conquer: recursively sort elements after pivot
		quickSortHelper(arr, pivotIndex+1, high)
	}
}

// partition rearranges elements so all elements ≤ pivot come before it.
// This is the key O(n) operation that processes all elements in the range.
func partition(arr []int, low, high int) int {
	// Choose the last element as pivot
	pivot := arr[high]
	// Index of smaller element (indicates proper position of pivot)
	i := low - 1

	// Move all elements smaller than or equal to pivot to the left
	for j := low; j < high; j++ {
		if arr[j] <= pivot {
			i++                             // Increment index of smaller element
			arr[i], arr[j] = arr[j], arr[i] // Swap elements
		}
	}

	// Place pivot in its correct position
	arr[i+1], arr[high] = arr[high], arr[i+1]

	// Return the partition index
	return i + 1
}
