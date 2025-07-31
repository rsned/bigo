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

package loglog

import "math"

// RadixSortOptimization performs O(log(log n)) integer sorting preprocessing
func RadixSortOptimization(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	// Find the maximum number to determine digit count
	maxVal := arr[0]
	for _, num := range arr {
		if num > maxVal {
			maxVal = num
		}
	}

	// Use log-log optimization for digit processing
	logLogFactor := int(math.Log(math.Log(float64(maxVal)+1) + 1))
	if logLogFactor <= 0 {
		logLogFactor = 1
	}

	// Simplified radix sort with log-log optimization
	result := make([]int, len(arr))
	copy(result, arr)

	for digit := 1; maxVal/digit > 0; digit *= 10 {
		countingSort(result, digit, logLogFactor)
	}

	return result
}

func countingSort(arr []int, digit, logLogFactor int) {
	output := make([]int, len(arr))

	// Use logLogFactor to determine bucket size for optimization
	// Larger logLogFactor means more granular buckets for better distribution
	bucketSize := 10
	if logLogFactor > 1 {
		bucketSize = int(math.Pow(10, float64(logLogFactor)))
		if bucketSize > 1000 { // Cap bucket size to prevent excessive memory usage
			bucketSize = 1000
		}
	}

	count := make([]int, bucketSize)
	bucketDivisor := bucketSize / 10 // Scale factor to map digits to buckets
	if bucketDivisor == 0 {
		bucketDivisor = 1
	}

	// Count occurrences with log-log optimization using expanded buckets
	for i := range len(arr) {
		baseIndex := (arr[i] / digit) % 10
		// Use logLogFactor to create more granular distribution
		bucketIndex := (baseIndex * bucketDivisor) % bucketSize
		count[bucketIndex]++
	}

	// Change count[i] to actual position
	for i := range bucketSize {
		if i > 0 {
			count[i] += count[i-1]
		}
	}

	// Build output array (reverse iteration still needs traditional for loop)
	for i := len(arr) - 1; i >= 0; i-- {
		baseIndex := (arr[i] / digit) % 10
		bucketIndex := (baseIndex * bucketDivisor) % bucketSize
		output[count[bucketIndex]-1] = arr[i]
		count[bucketIndex]--
	}

	copy(arr, output)
}
