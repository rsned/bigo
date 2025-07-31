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
func isSortedQuickSort(arr []int) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			return false
		}
	}

	return true
}

// Helper function to check if two arrays contain the same elements
func containsSameElementsQuickSort(a, b []int) bool {
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

// TestQuickSort tests the QuickSort function
func TestQuickSort(t *testing.T) {
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
			name:     "three elements",
			input:    []int{3, 1, 2},
			expected: []int{1, 2, 3},
		},
		{
			name:     "pivot at start",
			input:    []int{1, 5, 3, 4, 2},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "pivot at end",
			input:    []int{5, 3, 4, 2, 1},
			expected: []int{1, 2, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that input is not modified
			originalInput := make([]int, len(tt.input))
			copy(originalInput, tt.input)

			result := QuickSort(tt.input)

			// Verify input wasn't modified
			if !cmp.Equal(tt.input, originalInput) {
				t.Errorf("QuickSort modified input array")
			}

			// Verify result is correct
			if !cmp.Equal(result, tt.expected) {
				t.Errorf("QuickSort(%v) = %v, want %v", tt.input, result, tt.expected)
			}

			// Verify result is sorted
			if !isSortedQuickSort(result) {
				t.Errorf("QuickSort result is not sorted: %v", result)
			}

			// Verify result contains same elements as input
			if !containsSameElementsQuickSort(result, tt.input) {
				t.Errorf("QuickSort changed elements: input %v, result %v", tt.input, result)
			}
		})
	}
}

// TestQuickSortWorstCase tests QuickSort worst-case scenarios
func TestQuickSortWorstCase(t *testing.T) {
	t.Run("already_sorted_small", func(t *testing.T) {
		// Already sorted can be worst case for QuickSort
		input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		result := QuickSort(input)

		if !isSortedQuickSort(result) {
			t.Errorf("Already sorted array not handled correctly: %v", result)
		}

		if !cmp.Equal(result, input) {
			t.Errorf("Already sorted array changed: got %v, want %v", result, input)
		}
	})

	t.Run("reverse_sorted_small", func(t *testing.T) {
		// Reverse sorted can also be worst case
		input := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
		expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		result := QuickSort(input)

		if !isSortedQuickSort(result) {
			t.Errorf("Reverse sorted array not handled correctly: %v", result)
		}

		if !cmp.Equal(result, expected) {
			t.Errorf("Reverse sorted array: got %v, want %v", result, expected)
		}
	})

	t.Run("many_duplicates", func(t *testing.T) {
		input := make([]int, 50)
		for i := range input {
			input[i] = i % 3 // Only values 0, 1, 2
		}

		result := QuickSort(input)

		if !isSortedQuickSort(result) {
			t.Errorf("Many duplicates not sorted correctly")
		}

		if !containsSameElementsQuickSort(result, input) {
			t.Errorf("Many duplicates changed elements")
		}
	})
}

// TestQuickSortPartitioning tests the partitioning behavior
func TestQuickSortPartitioning(t *testing.T) {
	t.Run("partition_correctness", func(t *testing.T) {
		// Test that partitioning works correctly
		input := []int{3, 6, 8, 10, 1, 2, 1}
		result := QuickSort(input)
		expected := []int{1, 1, 2, 3, 6, 8, 10}

		if !cmp.Equal(result, expected) {
			t.Errorf("Partitioning test failed: got %v, want %v", result, expected)
		}
	})

	t.Run("pivot_handling", func(t *testing.T) {
		// Test with pivot element at different positions
		testCases := [][]int{
			{5, 1, 3, 9, 8, 2, 4, 7}, // Pivot 5 in middle
			{1, 5, 3, 9, 8, 2, 4, 7}, // Pivot 1 at start
			{9, 1, 3, 5, 8, 2, 4, 7}, // Pivot 9 near end
		}

		for i, testCase := range testCases {
			result := QuickSort(testCase)
			if !isSortedQuickSort(result) {
				t.Errorf("Pivot test case %d failed: result not sorted %v", i, result)
			}
		}
	})
}

// TestQuickSortStability tests QuickSort behavior (note: QuickSort is not stable)
func TestQuickSortStability(t *testing.T) {
	t.Run("not_stable_but_correct", func(t *testing.T) {
		// QuickSort is not stable, but should still sort correctly
		input := []int{3, 1, 4, 1, 5, 9, 2, 6}
		result := QuickSort(input)
		expected := []int{1, 1, 2, 3, 4, 5, 6, 9}

		if !cmp.Equal(result, expected) {
			t.Errorf("Stability test failed (correctness): got %v, want %v", result, expected)
		}
	})

	t.Run("large_array_stability", func(t *testing.T) {
		input := make([]int, 100)
		for i := range input {
			input[i] = (i * 37) % 50 // Creates many duplicates
		}

		result := QuickSort(input)

		if !isSortedQuickSort(result) {
			t.Errorf("Large array with duplicates not sorted correctly")
		}

		if !containsSameElementsQuickSort(result, input) {
			t.Errorf("Large array stability test changed elements")
		}
	})
}

// TestQuickSortComparison compares QuickSort with standard library sort
func TestQuickSortComparison(t *testing.T) {
	testCases := [][]int{
		{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5},
		{100, 99, 98, 97, 96, 95, 94, 93, 92, 91},
		{1, 1, 1, 1, 1},
		{-5, -3, -1, 1, 3, 5},
		{},
		{42},
		{64, 34, 25, 12, 22, 11, 90},
		{5, 2, 8, 6, 1, 9, 4, 0, 3, 7},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("comparison_test_%d", i), func(t *testing.T) {
			// Copy for standard library sort
			stdInput := make([]int, len(testCase))
			copy(stdInput, testCase)
			sort.Ints(stdInput)

			// Our QuickSort
			quickResult := QuickSort(testCase)

			if !cmp.Equal(quickResult, stdInput) {
				t.Errorf("QuickSort result differs from standard library: got %v, want %v",
					quickResult, stdInput)
			}
		})
	}
}

// Benchmark functions
func BenchmarkQuickSort(b *testing.B) {
	sizes := []int{100, 1000, 5000, 10000, 50000}

	for _, size := range sizes {
		// Random data (best case for QuickSort)
		randomData := make([]int, size)
		for i := range randomData {
			randomData[i] = (i * 71) % size
		}

		b.Run(fmt.Sprintf("random_size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				QuickSort(randomData)
			}
		})

		// Sorted data (potential worst case for QuickSort)
		sortedData := make([]int, size)
		for i := range sortedData {
			sortedData[i] = i
		}

		b.Run(fmt.Sprintf("sorted_size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				QuickSort(sortedData)
			}
		})

		// Reverse sorted data (potential worst case)
		reverseSortedData := make([]int, size)
		for i := range reverseSortedData {
			reverseSortedData[i] = size - i
		}

		b.Run(fmt.Sprintf("reverse_size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				QuickSort(reverseSortedData)
			}
		})

		// Data with many duplicates
		duplicateData := make([]int, size)
		for i := range duplicateData {
			duplicateData[i] = i % 10
		}

		b.Run(fmt.Sprintf("duplicates_size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				QuickSort(duplicateData)
			}
		})
	}
}

// BenchmarkQuickSortVsStdSort compares performance with standard library
func BenchmarkQuickSortVsStdSort(b *testing.B) {
	sizes := []int{1000, 10000}

	for _, size := range sizes {
		data := make([]int, size)
		for i := range data {
			data[i] = (i * 73) % size
		}

		b.Run(fmt.Sprintf("QuickSort_size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				QuickSort(data)
			}
		})

		b.Run(fmt.Sprintf("StdSort_size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				dataCopy := make([]int, len(data))
				copy(dataCopy, data)
				sort.Ints(dataCopy)
			}
		})
	}
}

// TestQuickSortEdgeCases tests edge cases and boundary conditions
func TestQuickSortEdgeCases(t *testing.T) {
	t.Run("nil_equivalent", func(t *testing.T) {
		var input []int
		result := QuickSort(input)
		if len(result) != 0 {
			t.Errorf("Nil slice should return empty slice, got %v", result)
		}
	})

	t.Run("extreme_values", func(t *testing.T) {
		input := []int{-2147483648, 2147483647, 0, -1, 1}
		result := QuickSort(input)
		expected := []int{-2147483648, -1, 0, 1, 2147483647}

		if !cmp.Equal(result, expected) {
			t.Errorf("Extreme values test failed: got %v, want %v", result, expected)
		}
	})

	t.Run("single_value_repeated", func(t *testing.T) {
		input := make([]int, 100)
		for i := range input {
			input[i] = 42
		}

		result := QuickSort(input)

		if len(result) != 100 {
			t.Errorf("Length changed for repeated values")
		}

		for _, val := range result {
			if val != 42 {
				t.Errorf("Value changed in repeated values test")
			}
		}
	})
}

// TestQuickSortConsistency tests that QuickSort produces consistent results
func TestQuickSortConsistency(t *testing.T) {
	t.Run("multiple_runs_same_result", func(t *testing.T) {
		input := []int{5, 2, 8, 1, 9, 3, 7, 4, 6}

		result1 := QuickSort(input)
		result2 := QuickSort(input)
		result3 := QuickSort(input)

		if !cmp.Equal(result1, result2) {
			t.Errorf("QuickSort gave different results: %v vs %v", result1, result2)
		}

		if !cmp.Equal(result2, result3) {
			t.Errorf("QuickSort gave different results: %v vs %v", result2, result3)
		}
	})

	t.Run("deterministic_output", func(t *testing.T) {
		// Same input should always produce same output
		input := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
		expected := []int{1, 1, 2, 3, 3, 4, 5, 5, 6, 9}

		for i := 0; i < 5; i++ {
			result := QuickSort(input)
			if !cmp.Equal(result, expected) {
				t.Errorf("Run %d: QuickSort gave unexpected result: got %v, want %v",
					i, result, expected)
			}
		}
	})
}

// TestQuickSortRecursion tests recursive behavior
func TestQuickSortRecursion(t *testing.T) {
	t.Run("deep_recursion_small", func(t *testing.T) {
		// Test that doesn't cause stack overflow but tests recursion
		input := make([]int, 30)
		for i := range input {
			input[i] = 30 - i // Reverse sorted to trigger deep recursion
		}

		result := QuickSort(input)

		if !isSortedQuickSort(result) {
			t.Errorf("Deep recursion test failed: not sorted")
		}

		expected := make([]int, 30)
		for i := range expected {
			expected[i] = i + 1
		}

		if !cmp.Equal(result, expected) {
			t.Errorf("Deep recursion test failed: got %v, want %v", result, expected)
		}
	})

	t.Run("balanced_partitioning", func(t *testing.T) {
		// Test case that should result in balanced partitioning
		input := []int{8, 3, 5, 4, 7, 6, 1, 2}
		result := QuickSort(input)
		expected := []int{1, 2, 3, 4, 5, 6, 7, 8}

		if !cmp.Equal(result, expected) {
			t.Errorf("Balanced partitioning test failed: got %v, want %v", result, expected)
		}
	})
}
