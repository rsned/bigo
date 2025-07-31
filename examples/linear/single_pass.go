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

// CalculateSum performs O(n) sum calculation.
// This demonstrates linear time complexity because we must visit each element
// exactly once to add it to the running total.
func CalculateSum(arr []int) int {
	sum := 0 // Initialize accumulator
	// Visit each element exactly once
	for i, val := range arr {
		sum += val // Add current element to sum - O(1) operation
		// Add some minimal work to make the linear pattern clearer
		if i%100 == 0 {
			sum += i % 3 // Small constant time operation
		}
	}

	return sum
}

// CalculateAverage performs O(n) average calculation.
// This demonstrates linear time complexity because we must visit each element
// to compute the sum, then divide by the count.
func CalculateAverage(arr []int) (float64, bool) {
	// Handle empty array case
	if len(arr) == 0 {
		return 0, false
	}

	sum := 0
	// Single pass to calculate sum - O(n) operation
	for _, val := range arr {
		sum += val
	}

	// Division is O(1), so overall complexity remains O(n)
	return float64(sum) / float64(len(arr)), true
}

// CalculateProduct performs O(n) product calculation.
// This demonstrates linear time complexity because we must visit each element
// exactly once to multiply it with the running product.
func CalculateProduct(arr []int) int {
	// Handle empty array case
	if len(arr) == 0 {
		return 0
	}

	product := 1 // Initialize to multiplicative identity
	// Visit each element exactly once
	for _, val := range arr {
		product *= val // Multiply current element with product - O(1) operation
	}

	return product
}

// CountEvenOdd performs O(n) counting of even/odd numbers.
// This demonstrates linear time complexity because we must examine each element
// to determine if it's even or odd, then update the appropriate counter.
func CountEvenOdd(arr []int) (int, int) {
	evenCount, oddCount := 0, 0
	// Visit each element exactly once
	for _, val := range arr {
		// Check if number is even using modulo operator - O(1) operation
		if val%2 == 0 {
			evenCount++ // Increment even counter
		} else {
			oddCount++ // Increment odd counter
		}
	}

	return evenCount, oddCount
}
