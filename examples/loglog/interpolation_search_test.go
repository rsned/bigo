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

import (
	"fmt"
	"testing"

	"github.com/rsned/bigo/examples/linear"
)

func TestInterpolationSearch(t *testing.T) {
	tests := []struct {
		name   string
		arr    []int
		target int
		want   int
	}{
		{"empty array", []int{}, 5, -1},
		{"single element found", []int{10}, 10, 0},
		{"single element not found", []int{10}, 5, -1},
		{"target at beginning", []int{1, 2, 3, 4, 5}, 1, 0},
		{"target at end", []int{1, 2, 3, 4, 5}, 5, 4},
		{"target in middle", []int{1, 2, 3, 4, 5}, 3, 2},
		{"target not found smaller", []int{2, 4, 6, 8, 10}, 1, -1},
		{"target not found larger", []int{2, 4, 6, 8, 10}, 12, -1},
		{"target not found middle", []int{2, 4, 8, 10}, 6, -1},
		{"uniformly spaced array", []int{10, 20, 30, 40, 50}, 30, 2},
		{"larger uniform array", []int{5, 10, 15, 20, 25, 30, 35, 40}, 25, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InterpolationSearch(tt.arr, tt.target)
			if got != tt.want {
				t.Errorf("InterpolationSearch(%v, %d) = %d, want %d", tt.arr, tt.target, got, tt.want)
			}
		})
	}
}

func TestInterpolationSearch_NonUniformDistribution(t *testing.T) {
	// Test with non-uniformly distributed data where interpolation might overshoot
	arr := []int{1, 2, 3, 100, 101, 102}

	tests := []struct {
		target int
		want   int
	}{
		{1, 0},
		{3, 2},
		{100, 3},
		{102, 5},
		{50, -1}, // Should not be found
	}

	for _, tt := range tests {
		got := InterpolationSearch(arr, tt.target)
		if got != tt.want {
			t.Errorf("InterpolationSearch(non-uniform, %d) = %d, want %d", tt.target, got, tt.want)
		}
	}
}

func TestInterpolationSearch_EdgeCases(t *testing.T) {
	// Test edge cases for interpolation formula
	tests := []struct {
		name   string
		arr    []int
		target int
		want   int
	}{
		{"target equals first element", []int{5, 10, 15}, 5, 0},
		{"target equals last element", []int{5, 10, 15}, 15, 2},
		{"all elements same", []int{7, 7, 7, 7}, 7, 0}, // Should find first occurrence
		{"two elements", []int{1, 2}, 2, 1},
		{"negative numbers", []int{-10, -5, 0, 5, 10}, 0, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InterpolationSearch(tt.arr, tt.target)
			if got != tt.want {
				t.Errorf("InterpolationSearch(%v, %d) = %d, want %d", tt.arr, tt.target, got, tt.want)
			}
		})
	}
}

func TestInterpolationSearch_LargeArray(t *testing.T) {
	// Create a large uniformly distributed array
	size := 1000
	arr := make([]int, size)
	for i := range size {
		arr[i] = i * 10 // 0, 10, 20, 30, ..., 9990
	}

	tests := []struct {
		target int
		want   int
	}{
		{0, 0},
		{500, 50},
		{9990, 999},
		{250, 25},
		{9995, -1}, // Not in array
	}

	for _, tt := range tests {
		got := InterpolationSearch(arr, tt.target)
		if got != tt.want {
			t.Errorf("InterpolationSearch(large array, %d) = %d, want %d", tt.target, got, tt.want)
		}
	}
}

func TestInterpolationSearch_OutOfBounds(t *testing.T) {
	arr := []int{10, 20, 30, 40, 50}

	// Test values outside the range of array values
	tests := []struct {
		target int
		want   int
	}{
		{5, -1},  // Below minimum
		{60, -1}, // Above maximum
	}

	for _, tt := range tests {
		got := InterpolationSearch(arr, tt.target)
		if got != tt.want {
			t.Errorf("InterpolationSearch(%v, %d) = %d, want %d", arr, tt.target, got, tt.want)
		}
	}
}

func TestInterpolationSearch_CompareWithLinearSearch(t *testing.T) {
	// Test that interpolation search produces same results as linear search
	testArrays := [][]int{
		{1, 3, 5, 7, 9, 11, 13, 15},
		{2, 4, 6, 8, 10, 12, 14, 16, 18, 20},
		{100, 200, 300, 400, 500},
	}

	for _, arr := range testArrays {
		for _, target := range arr {
			interpolationResult := InterpolationSearch(arr, target)
			linearResult := linear.Search(arr, target)

			if interpolationResult != linearResult {
				t.Errorf("InterpolationSearch and LinearSearch disagree for %v, target %d: interpolation=%d, linear=%d",
					arr, target, interpolationResult, linearResult)
			}
		}
	}
}

// Benchmark functions for Interpolation Search O(log(log n)) complexity

func BenchmarkInterpolationSearch(b *testing.B) {
	// Log-log complexity works best with uniformly distributed data
	for size := 100; size <= 1000000; size *= 10 {
		// Create uniformly distributed sorted array
		arr := make([]int, size)
		for i := range arr {
			arr[i] = i * 10 // Uniformly spaced
		}

		b.Run(fmt.Sprintf("uniform-size-%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				target := (i % size) * 10
				InterpolationSearch(arr, target)
			}
		})
	}
}

func BenchmarkInterpolationSearchWorstCase(b *testing.B) {
	// Non-uniform distribution can degrade to O(n)
	sizes := []int{1000, 10000, 100000}

	for _, size := range sizes {
		// Create non-uniformly distributed array (exponential growth)
		arr := make([]int, size)
		for i := range arr {
			arr[i] = i * i // Quadratic growth
		}

		b.Run(fmt.Sprintf("non-uniform-size-%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				target := (i % size) * (i % size)
				InterpolationSearch(arr, target)
			}
		})
	}
}

func BenchmarkInterpolationSearchBestCase(b *testing.B) {
	// Best case: target is at the interpolated position
	sizes := []int{1000, 10000, 100000, 1000000}

	for _, size := range sizes {
		arr := make([]int, size)
		for i := range arr {
			arr[i] = i * 5 // Uniformly spaced
		}

		// Target that should be found immediately
		target := (size / 2) * 5

		b.Run(fmt.Sprintf("best-case-size-%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				InterpolationSearch(arr, target)
			}
		})
	}
}

func BenchmarkInterpolationSearchMissing(b *testing.B) {
	// Search for values not in the array
	sizes := []int{1000, 10000, 100000}

	for _, size := range sizes {
		arr := make([]int, size)
		for i := range arr {
			arr[i] = i * 2 // Even numbers only
		}

		b.Run(fmt.Sprintf("missing-size-%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				target := (i%size)*2 + 1 // Odd numbers (not in array)
				InterpolationSearch(arr, target)
			}
		})
	}
}
