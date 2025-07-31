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

import (
	"fmt"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Helper function to check if an array is sorted
func isSortedHeapSort(arr []int) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			return false
		}
	}

	return true
}

// Helper function to check if two arrays contain the same elements
func containsSameElementsHeapSort(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	aCopy := make([]int, len(a))
	bCopy := make([]int, len(b))
	copy(aCopy, a)
	copy(bCopy, b)

	sort.Ints(aCopy)
	sort.Ints(bCopy)

	return cmp.Equal(aCopy, bCopy)
}

// TestHeapSort tests the HeapSort function
func TestHeapSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "empty array",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "single element",
			input:    []int{5},
			expected: []int{5},
		},
		{
			name:     "two elements sorted",
			input:    []int{1, 2},
			expected: []int{1, 2},
		},
		{
			name:     "two elements reverse",
			input:    []int{2, 1},
			expected: []int{1, 2},
		},
		{
			name:     "already sorted",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "reverse sorted",
			input:    []int{5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "random order",
			input:    []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5},
			expected: []int{1, 1, 2, 3, 3, 4, 5, 5, 5, 6, 9},
		},
		{
			name:     "with duplicates",
			input:    []int{4, 2, 2, 1, 3, 2, 4},
			expected: []int{1, 2, 2, 2, 3, 4, 4},
		},
		{
			name:     "all same elements",
			input:    []int{7, 7, 7, 7, 7},
			expected: []int{7, 7, 7, 7, 7},
		},
		{
			name:     "negative numbers",
			input:    []int{-3, 1, -5, 4, -1, 0, 2},
			expected: []int{-5, -3, -1, 0, 1, 2, 4},
		},
		{
			name:     "mixed positive and negative",
			input:    []int{10, -5, 3, -2, 7, -1, 0},
			expected: []int{-5, -2, -1, 0, 3, 7, 10},
		},
		{
			name:     "large numbers",
			input:    []int{1000, 500, 1500, 200, 800},
			expected: []int{200, 500, 800, 1000, 1500},
		},
		{
			name:     "zero included",
			input:    []int{3, 0, -1, 2, 0, 1},
			expected: []int{-1, 0, 0, 1, 2, 3},
		},
		{
			name:     "palindromic sequence",
			input:    []int{1, 2, 3, 4, 3, 2, 1},
			expected: []int{1, 1, 2, 2, 3, 3, 4},
		},
		{
			name:     "powers of two",
			input:    []int{64, 1, 32, 4, 16, 2, 8},
			expected: []int{1, 2, 4, 8, 16, 32, 64},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that input is not modified
			originalInput := make([]int, len(tt.input))
			copy(originalInput, tt.input)

			result := HeapSort(tt.input)

			// Verify input wasn't modified
			if !cmp.Equal(tt.input, originalInput) {
				t.Errorf("HeapSort modified input array")
			}

			// Verify result is correct
			if !cmp.Equal(result, tt.expected) {
				t.Errorf("HeapSort(%v) = %v, want %v", tt.input, result, tt.expected)
			}

			// Verify result is sorted
			if !isSortedHeapSort(result) {
				t.Errorf("HeapSort result is not sorted: %v", result)
			}

			// Verify result contains same elements as input
			if !containsSameElementsHeapSort(result, tt.input) {
				t.Errorf("HeapSort changed elements: input %v, result %v", tt.input, result)
			}
		})
	}
}

// TestHeapSortProperties tests specific properties of heap sort
func TestHeapSortProperties(t *testing.T) {
	t.Run("stability_not_required", func(t *testing.T) {
		// Heap sort is not stable, but it should still sort correctly
		// This test verifies correct sorting regardless of stability

		// Convert to simple test that checks sorting works
		input := []int{3, 1, 4, 1, 5, 9, 2, 6}
		result := HeapSort(input)
		expected := []int{1, 1, 2, 3, 4, 5, 6, 9}

		if !cmp.Equal(result, expected) {
			t.Errorf("HeapSort stability test failed: got %v, want %v", result, expected)
		}
	})

	t.Run("in_place_property", func(t *testing.T) {
		// Test that HeapSort returns a new array and doesn't modify input
		input := []int{5, 2, 8, 1, 9, 3}
		inputBackup := make([]int, len(input))
		copy(inputBackup, input)

		result := HeapSort(input)

		// Input should be unchanged
		if !cmp.Equal(input, inputBackup) {
			t.Errorf("HeapSort modified input array: before %v, after %v", inputBackup, input)
		}

		// Result should be sorted
		if !isSortedHeapSort(result) {
			t.Errorf("HeapSort result not sorted: %v", result)
		}
	})

	t.Run("heap_property_during_sort", func(t *testing.T) {
		// Test with a specific case that exercises heap building and extraction
		input := []int{10, 7, 8, 4, 6, 2, 5, 1, 3, 9}
		result := HeapSort(input)

		if !isSortedHeapSort(result) {
			t.Errorf("Heap property test failed: result not sorted %v", result)
		}

		if !containsSameElementsHeapSort(result, input) {
			t.Errorf("Elements changed during heap sort")
		}
	})
}

// TestHeapSortEdgeCases tests edge cases and boundary conditions
func TestHeapSortEdgeCases(t *testing.T) {
	t.Run("large_array", func(t *testing.T) {
		// Test with a larger array
		input := make([]int, 1000)
		for i := 0; i < 1000; i++ {
			input[i] = (i * 37) % 100 // Creates a pseudo-random pattern
		}

		result := HeapSort(input)

		if !isSortedHeapSort(result) {
			t.Errorf("Large array not sorted correctly")
		}

		if !containsSameElementsHeapSort(result, input) {
			t.Errorf("Large array elements changed")
		}
	})

	t.Run("many_duplicates", func(t *testing.T) {
		input := make([]int, 100)
		for i := range input {
			input[i] = i % 5 // Many duplicates
		}

		result := HeapSort(input)

		if !isSortedHeapSort(result) {
			t.Errorf("Array with many duplicates not sorted correctly")
		}

		if !containsSameElementsHeapSort(result, input) {
			t.Errorf("Array with duplicates elements changed")
		}
	})

	t.Run("already_heap", func(t *testing.T) {
		// Input that's already a valid max heap
		input := []int{16, 14, 10, 8, 7, 9, 3, 2, 4, 1}
		result := HeapSort(input)
		expected := []int{1, 2, 3, 4, 7, 8, 9, 10, 14, 16}

		if !cmp.Equal(result, expected) {
			t.Errorf("Already heap test failed: got %v, want %v", result, expected)
		}
	})

	t.Run("reverse_heap", func(t *testing.T) {
		// Input that's a min heap (reverse of max heap)
		input := []int{1, 2, 3, 4, 7, 8, 9, 10, 14, 16}
		result := HeapSort(input)

		if !isSortedHeapSort(result) {
			t.Errorf("Reverse heap not sorted correctly: %v", result)
		}

		if !cmp.Equal(result, input) {
			t.Errorf("Already sorted array changed: got %v, want %v", result, input)
		}
	})
}

// TestHeapSortComparison compares HeapSort with standard library sort
func TestHeapSortComparison(t *testing.T) {
	testCases := [][]int{
		{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5},
		{100, 99, 98, 97, 96, 95, 94, 93, 92, 91},
		{1, 1, 1, 1, 1},
		{-5, -3, -1, 1, 3, 5},
		{},
		{42},
		{10, 7, 8, 4, 6, 2, 5, 1, 3, 9},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("comparison_test_%d", i), func(t *testing.T) {
			// Copy for standard library sort
			stdInput := make([]int, len(testCase))
			copy(stdInput, testCase)
			sort.Ints(stdInput)

			// Our HeapSort
			heapResult := HeapSort(testCase)

			if !cmp.Equal(heapResult, stdInput) {
				t.Errorf("HeapSort result differs from standard library: got %v, want %v",
					heapResult, stdInput)
			}
		})
	}
}

// Benchmark functions
func BenchmarkHeapSort(b *testing.B) {
	sizes := []int{100, 1000, 5000, 10000, 50000}

	for _, size := range sizes {
		// Random data
		randomData := make([]int, size)
		for i := range randomData {
			randomData[i] = (i * 71) % size
		}

		b.Run(fmt.Sprintf("random_size_%d", size), func(b *testing.B) {
			for b.Loop() {
				HeapSort(randomData)
			}
		})

		// Sorted data (best case for some algorithms, but not heap sort)
		sortedData := make([]int, size)
		for i := range sortedData {
			sortedData[i] = i
		}

		b.Run(fmt.Sprintf("sorted_size_%d", size), func(b *testing.B) {
			for b.Loop() {
				HeapSort(sortedData)
			}
		})

		// Reverse sorted data (worst case for some algorithms)
		reverseSortedData := make([]int, size)
		for i := range reverseSortedData {
			reverseSortedData[i] = size - i
		}

		b.Run(fmt.Sprintf("reverse_size_%d", size), func(b *testing.B) {
			for b.Loop() {
				HeapSort(reverseSortedData)
			}
		})

		// Data with many duplicates
		duplicateData := make([]int, size)
		for i := range duplicateData {
			duplicateData[i] = i % 10
		}

		b.Run(fmt.Sprintf("duplicates_size_%d", size), func(b *testing.B) {
			for b.Loop() {
				HeapSort(duplicateData)
			}
		})
	}
}

// BenchmarkHeapSortVsStdSort compares performance with standard library
func BenchmarkHeapSortVsStdSort(b *testing.B) {
	sizes := []int{1000, 10000}

	for _, size := range sizes {
		data := make([]int, size)
		for i := range data {
			data[i] = (i * 73) % size
		}

		b.Run(fmt.Sprintf("HeapSort_size_%d", size), func(b *testing.B) {
			for b.Loop() {
				b.StopTimer()
				dataCopy := make([]int, len(data))
				copy(dataCopy, data)
				b.StartTimer()
				HeapSort(dataCopy)
			}
		})

		b.Run(fmt.Sprintf("StdSort_size_%d", size), func(b *testing.B) {
			for b.Loop() {
				b.StopTimer()
				dataCopy := make([]int, len(data))
				copy(dataCopy, data)
				b.StartTimer()
				sort.Ints(dataCopy)
			}
		})
	}
}

// TestHeapSortConsistency tests that HeapSort produces consistent results
func TestHeapSortConsistency(t *testing.T) {
	t.Run("multiple_runs_same_result", func(t *testing.T) {
		input := []int{5, 2, 8, 1, 9, 3, 7, 4, 6}

		result1 := HeapSort(input)
		result2 := HeapSort(input)
		result3 := HeapSort(input)

		if !cmp.Equal(result1, result2) {
			t.Errorf("HeapSort gave different results: %v vs %v", result1, result2)
		}

		if !cmp.Equal(result2, result3) {
			t.Errorf("HeapSort gave different results: %v vs %v", result2, result3)
		}
	})

	t.Run("deterministic_output", func(t *testing.T) {
		// Same input should always produce same output
		input := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
		expected := []int{1, 1, 2, 3, 3, 4, 5, 5, 6, 9}

		for i := 0; i < 5; i++ {
			result := HeapSort(input)
			if !cmp.Equal(result, expected) {
				t.Errorf("Run %d: HeapSort gave unexpected result: got %v, want %v",
					i, result, expected)
			}
		}
	})
}

// TestHeapSortExtremeValues tests with extreme values
func TestHeapSortExtremeValues(t *testing.T) {
	t.Run("max_min_int", func(t *testing.T) {
		input := []int{2147483647, -2147483648, 0, 1, -1}
		result := HeapSort(input)
		expected := []int{-2147483648, -1, 0, 1, 2147483647}

		if !cmp.Equal(result, expected) {
			t.Errorf("Extreme values test failed: got %v, want %v", result, expected)
		}
	})

	t.Run("single_repeated_value", func(t *testing.T) {
		input := make([]int, 1000)
		for i := range input {
			input[i] = 42
		}

		result := HeapSort(input)

		if len(result) != 1000 {
			t.Errorf("Length changed for repeated values: got %d, want 1000", len(result))
		}

		for i, val := range result {
			if val != 42 {
				t.Errorf("Value changed at index %d: got %d, want 42", i, val)
			}
		}
	})
}
