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

// MergeSort performs O(n log n) merge sort.
// This demonstrates linearithmic time complexity because we divide the array
// into halves (log n levels) and at each level we merge all elements (n work).
// The divide-and-conquer approach gives us log n levels × n work = O(n log n).
func MergeSort(arr []int) []int {
	// Base case: arrays with 0 or 1 element are already sorted
	if len(arr) <= 1 {
		return arr
	}

	// Divide: split array into two halves
	mid := len(arr) / 2
	// Conquer: recursively sort both halves
	left := MergeSort(arr[:mid])  // Sort left half
	right := MergeSort(arr[mid:]) // Sort right half

	// Combine: merge the sorted halves
	return merge(left, right)
}

// merge combines two sorted arrays into a single sorted array.
// This operation is O(n) where n is the total number of elements,
// and it's the key operation that makes merge sort O(n log n).
func merge(left, right []int) []int {
	// Pre-allocate result slice with exact capacity for efficiency
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0 // Indices for left and right arrays

	// Merge elements in sorted order - O(n) operation
	for i < len(left) && j < len(right) {
		// Take smaller element from either left or right
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++ // Move left pointer
		} else {
			result = append(result, right[j])
			j++ // Move right pointer
		}
	}

	// Append remaining elements from left array (if any)
	result = append(result, left[i:]...)
	// Append remaining elements from right array (if any)
	result = append(result, right[j:]...)

	return result
}
