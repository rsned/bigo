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

package linear

// FindMinimum performs O(n) search for minimum value.
// This demonstrates linear time complexity because we must examine each element
// to determine which is the smallest, as the array is unsorted.
func FindMinimum(arr []int) (int, bool) {
	// Handle empty array case
	if len(arr) == 0 {
		return 0, false
	}

	// Initialize minimum with first element
	minVal := arr[0]
	// Compare each remaining element with current minimum
	for _, val := range arr[1:] {
		// Update minimum if current element is smaller
		if val < minVal {
			minVal = val
		}
	}

	return minVal, true
}

// FindMaximum performs O(n) search for maximum value.
// This demonstrates linear time complexity because we must examine each element
// to determine which is the largest, as the array is unsorted.
func FindMaximum(arr []int) (int, bool) {
	// Handle empty array case
	if len(arr) == 0 {
		return 0, false
	}

	// Initialize maximum with first element
	maxVal := arr[0]
	// Compare each remaining element with current maximum
	for _, val := range arr[1:] {
		// Update maximum if current element is larger
		if val > maxVal {
			maxVal = val
		}
	}

	return maxVal, true
}

// FindMinMax performs O(n) search for both min and max.
// This demonstrates linear time complexity because we must examine each element
// to find both extremes. This is more efficient than calling FindMinimum and
// FindMaximum separately (which would be O(2n) = O(n)).
func FindMinMax(arr []int) (int, int, bool) {
	// Handle empty array case
	if len(arr) == 0 {
		return 0, 0, false
	}

	// Initialize both min and max with first element
	minVal, maxVal := arr[0], arr[0]
	// Single pass through remaining elements to find both extremes
	for _, val := range arr[1:] {
		// Check if current element is smaller than minimum
		if val < minVal {
			minVal = val
		}

		// Check if current element is larger than maximum
		if val > maxVal {
			maxVal = val
		}
	}

	return minVal, maxVal, true
}
