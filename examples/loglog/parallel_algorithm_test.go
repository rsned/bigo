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

	"github.com/google/go-cmp/cmp"
	"github.com/rsned/bigo/examples/linear"
)

func TestParallelDivideConquer(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want int
	}{
		{"empty array", []int{}, 0},
		{"single element", []int{42}, 42},
		{"two elements", []int{10, 20}, 20},
		{"multiple elements", []int{1, 5, 3, 9, 2}, 9},
		{"negative numbers", []int{-5, -1, -10, -3}, -1},
		{"mixed numbers", []int{-2, 0, 5, -8, 3}, 5},
		{"all same", []int{7, 7, 7, 7}, 7},
		{"decreasing order", []int{10, 8, 6, 4, 2}, 10},
		{"increasing order", []int{1, 3, 5, 7, 9}, 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParallelDivideConquer(tt.arr)
			if got != tt.want {
				t.Errorf("ParallelDivideConquer(%v) = %d, want %d", tt.arr, got, tt.want)
			}
		})
	}
}

func TestParallelDivideConquer_CompareWithSequential(t *testing.T) {
	// Test that parallel algorithm produces same results as sequential max finding
	testArrays := [][]int{
		{1, 2, 3, 4, 5},
		{5, 4, 3, 2, 1},
		{3, 1, 4, 1, 5, 9, 2, 6},
		{-1, -5, -2, -8, -3},
		{0},
		{100, 50, 75, 25, 90},
	}

	for _, arr := range testArrays {
		parallelResult := ParallelDivideConquer(arr)
		sequentialResult, _ := linear.FindMaximum(arr)

		if len(arr) > 0 && parallelResult != sequentialResult {
			t.Errorf("ParallelDivideConquer and FindMaximum disagree for %v: parallel=%d, sequential=%d",
				arr, parallelResult, sequentialResult)
		}
	}
}

func TestParallelDivideConquer_LargeArray(t *testing.T) {
	// Test with larger array to verify parallel processing
	size := 100
	arr := make([]int, size)
	maxVal := 0

	// Create array with known maximum at position 73
	for i := 0; i < size; i++ {
		arr[i] = i
		if i == 73 {
			arr[i] = 1000 // This should be the maximum
			maxVal = 1000
		}
	}

	got := ParallelDivideConquer(arr)
	if got != maxVal {
		t.Errorf("ParallelDivideConquer(large array) = %d, want %d", got, maxVal)
	}
}

func TestParallelDivideConquer_PowerOfTwo(t *testing.T) {
	// Test with array sizes that are powers of 2 (optimal for divide-and-conquer)
	tests := []struct {
		size int
		max  int
	}{
		{2, 100},
		{4, 200},
		{8, 300},
		{16, 400},
		{32, 500},
	}

	for _, tt := range tests {
		arr := make([]int, tt.size)
		for i := 0; i < tt.size; i++ {
			arr[i] = i
		}
		arr[tt.size/2] = tt.max // Place max in middle

		got := ParallelDivideConquer(arr)
		if got != tt.max {
			t.Errorf("ParallelDivideConquer(size %d) = %d, want %d", tt.size, got, tt.max)
		}
	}
}

func TestParallelDivideConquer_DoesNotModifyInput(t *testing.T) {
	original := []int{3, 1, 4, 1, 5, 9, 2, 6}
	input := make([]int, len(original))
	copy(input, original)

	ParallelDivideConquer(input)

	if !cmp.Equal(input, original) {
		t.Errorf("ParallelDivideConquer modified input array: got %v, want %v", input, original)
	}
}

func TestParallelDivideConquer_DeterministicResults(t *testing.T) {
	// Test that algorithm produces consistent results across multiple runs
	arr := []int{5, 2, 8, 1, 9, 3}
	expected := ParallelDivideConquer(arr)

	// Run multiple times to ensure consistency
	for i := 0; i < 10; i++ {
		got := ParallelDivideConquer(arr)
		if got != expected {
			t.Errorf("ParallelDivideConquer produced inconsistent result on run %d: got %d, want %d", i+1, got, expected)
		}
	}
}

// Benchmark functions for Parallel Divide and Conquer O(log(log n)) complexity

func BenchmarkParallelDivideConquer(b *testing.B) {
	sizes := []int{1000, 10000, 100000, 1000000}

	for _, size := range sizes {
		arr := make([]int, size)
		for i := 0; i < size; i++ {
			arr[i] = i + 1
		}
		arr[size/2] = size * 2 // Max element in middle

		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ParallelDivideConquer(arr)
			}
		})
	}
}

func BenchmarkParallelDivideConquerWorstCase(b *testing.B) {
	// Test with varying data patterns
	sizes := []int{1000, 10000, 100000}

	for _, size := range sizes {
		arr := make([]int, size)
		for i := 0; i < size; i++ {
			arr[i] = size - i // Descending order
		}

		b.Run(fmt.Sprintf("descending-size-%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ParallelDivideConquer(arr)
			}
		})
	}
}

func BenchmarkParallelDivideConquerSmallArrays(b *testing.B) {
	// Test with small arrays where parallelization overhead matters
	sizes := []int{10, 50, 100, 500}

	for _, size := range sizes {
		arr := make([]int, size)
		for i := 0; i < size; i++ {
			arr[i] = i + 1
		}

		b.Run(fmt.Sprintf("small-size-%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ParallelDivideConquer(arr)
			}
		})
	}
}
