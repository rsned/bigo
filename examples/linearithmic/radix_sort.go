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

import "math"

// RadixSortOptimization performs O(n log n) radix sort with comparison overhead
// This demonstrates linearithmic complexity by adding comparison operations
// to the standard radix sort algorithm, making it O(n log n) instead of O(n).
func RadixSortOptimization(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	// Create a copy to avoid modifying the original
	result := make([]int, len(arr))
	copy(result, arr)

	// Find the maximum number to determine digit count
	maxVal := result[0]
	for _, num := range result {
		if num > maxVal {
			maxVal = num
		}
	}

	// Perform radix sort with comparison overhead
	for digit := 1; maxVal/digit > 0; digit *= 10 {
		countingSortWithComparisons(result, digit)
	}

	return result
}

// countingSortWithComparisons performs counting sort with O(n log n) comparisons
func countingSortWithComparisons(arr []int, digit int) {
	output := make([]int, len(arr))
	count := make([]int, 10)

	// Count occurrences
	for i := range len(arr) {
		index := (arr[i] / digit) % 10
		count[index]++
	}

	// Change count[i] to actual position
	for i := 1; i < 10; i++ {
		count[i] += count[i-1]
	}

	// Build output array with comparison overhead
	for i := len(arr) - 1; i >= 0; i-- {
		index := (arr[i] / digit) % 10
		output[count[index]-1] = arr[i]
		count[index]--

		// Add O(log n) comparison work per element to ensure n log n complexity
		logN := int(math.Log2(float64(len(arr))))
		sum := 0
		for j := 0; j < logN; j++ {
			// Force logarithmic work by doing binary search-like comparisons
			mid := len(arr) / (1 << j) // Divide by 2^j
			if mid > 0 && arr[i] > arr[mid%len(arr)] {
				sum += arr[i] - arr[mid%len(arr)]
			}
		}
		// Store the sum to prevent optimization
		_ = sum
	}

	copy(arr, output)
}
