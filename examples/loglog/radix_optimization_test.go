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
	"math/rand/v2"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRadixSortOptimization(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want []int
	}{
		{"empty array", []int{}, []int{}},
		{"single element", []int{42}, []int{42}},
		{"already sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"reverse sorted", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"random order", []int{170, 45, 75, 90, 2, 802, 24, 66}, []int{2, 24, 45, 66, 75, 90, 170, 802}},
		{"duplicates", []int{3, 1, 3, 1, 3}, []int{1, 1, 3, 3, 3}},
		{"single digit numbers", []int{7, 3, 9, 1, 5}, []int{1, 3, 5, 7, 9}},
		{"two digit numbers", []int{23, 12, 34, 56, 78}, []int{12, 23, 34, 56, 78}},
		{"three digit numbers", []int{123, 456, 789, 101, 202}, []int{101, 123, 202, 456, 789}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RadixSortOptimization(tt.arr)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("RadixSortOptimization(%v) = %v, want %v", tt.arr, got, tt.want)
			}
		})
	}
}

func TestRadixSortOptimization_ZeroValues(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want []int
	}{
		{"with zeros", []int{0, 3, 0, 1, 0}, []int{0, 0, 0, 1, 3}},
		{"all zeros", []int{0, 0, 0}, []int{0, 0, 0}},
		{"single zero", []int{0}, []int{0}},
		{"zeros and positives", []int{5, 0, 3, 0, 1}, []int{0, 0, 1, 3, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RadixSortOptimization(tt.arr)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("RadixSortOptimization(%v) = %v, want %v", tt.arr, got, tt.want)
			}
		})
	}
}

func TestRadixSortOptimization_LargeNumbers(t *testing.T) {
	// Test with larger numbers to verify digit processing
	arr := []int{9999, 1111, 5555, 3333, 7777, 2222, 8888, 4444, 6666}
	want := []int{1111, 2222, 3333, 4444, 5555, 6666, 7777, 8888, 9999}

	got := RadixSortOptimization(arr)
	if !cmp.Equal(got, want) {
		t.Errorf("RadixSortOptimization(large numbers) = %v, want %v", got, want)
	}
}

func TestRadixSortOptimization_DoesNotModifyOriginal(t *testing.T) {
	original := []int{3, 1, 4, 1, 5, 9, 2, 6}
	input := make([]int, len(original))
	copy(input, original)

	RadixSortOptimization(input)

	if !cmp.Equal(input, original) {
		t.Errorf("RadixSortOptimization modified original array: got %v, want %v", input, original)
	}
}

func TestRadixSortOptimization_CompareWithStandardSort(t *testing.T) {
	// Test that radix sort produces same results as standard sort
	testArrays := [][]int{
		{170, 45, 75, 90, 2, 802, 24, 66},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		{5, 5, 5, 5, 5},
		{100, 1, 10, 1000, 10000},
	}

	for _, arr := range testArrays {
		radixResult := RadixSortOptimization(arr)

		standardSorted := make([]int, len(arr))
		copy(standardSorted, arr)
		sort.Ints(standardSorted)

		if !cmp.Equal(radixResult, standardSorted) {
			t.Errorf("RadixSortOptimization and sort.Ints disagree for %v: radix=%v, standard=%v",
				arr, radixResult, standardSorted)
		}
	}
}

func TestRadixSortOptimization_StabilityCheck(t *testing.T) {
	// While this implementation may not be stable, we test that equal elements appear correctly
	arr := []int{21, 11, 31, 41, 21, 11}
	got := RadixSortOptimization(arr)
	want := []int{11, 11, 21, 21, 31, 41}

	if !cmp.Equal(got, want) {
		t.Errorf("RadixSortOptimization stability test: got %v, want %v", got, want)
	}
}

func TestRadixSortOptimization_VaryingDigits(t *testing.T) {
	// Test with numbers having different digit counts
	arr := []int{5, 25, 125, 1, 51, 251, 1251}
	want := []int{1, 5, 25, 51, 125, 251, 1251}

	got := RadixSortOptimization(arr)
	if !cmp.Equal(got, want) {
		t.Errorf("RadixSortOptimization(varying digits) = %v, want %v", got, want)
	}
}

func TestRadixSortOptimization_PowersOfTen(t *testing.T) {
	// Test with powers of ten to verify digit boundary handling
	arr := []int{1000, 100, 10, 1, 10000}
	want := []int{1, 10, 100, 1000, 10000}

	got := RadixSortOptimization(arr)
	if !cmp.Equal(got, want) {
		t.Errorf("RadixSortOptimization(powers of ten) = %v, want %v", got, want)
	}
}

func TestRadixSortOptimization_Performance(t *testing.T) {
	// Test with larger array to verify it completes efficiently
	size := 1000
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = size - i // Reverse order
	}

	got := RadixSortOptimization(arr)

	// Verify it's sorted
	for i := 1; i < len(got); i++ {
		if got[i] < got[i-1] {
			t.Errorf("RadixSortOptimization(large array) not properly sorted at index %d: %d > %d",
				i-1, got[i-1], got[i])

			break
		}
	}

	// Verify first and last elements
	if got[0] != 1 || got[len(got)-1] != size {
		t.Errorf("RadixSortOptimization(large array) boundary error: first=%d, last=%d, want first=1, last=%d",
			got[0], got[len(got)-1], size)
	}
}

// Benchmark functions for Radix Sort Optimization (O(d * n) where d is digits)

func BenchmarkRadixSortOptimizationReverseOrder(b *testing.B) {
	// Radix sort is linear in the number of digits times number of elements
	sizes := []int{1000, 10000, 100000, 1000000}

	for _, size := range sizes {
		arr := make([]int, size)
		for i := range size {
			arr[i] = size - i // Reverse order
		}

		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				RadixSortOptimization(arr)
			}
		})
	}
}

func BenchmarkRadixSortOptimizationLargeNumbers(b *testing.B) {
	// Test with larger numbers (more digits)
	sizes := []int{1000, 10000, 100000}

	for _, size := range sizes {
		arr := make([]int, size)
		for i := range size {
			arr[i] = (size - i) * 100000
		}

		b.Run(fmt.Sprintf("large-numbers-size-%d", size), func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				RadixSortOptimization(arr)
			}
		})
	}
}

func BenchmarkRadixSortOptimizationSmallNumbers(b *testing.B) {
	// Test with small numbers (fewer digits)
	sizes := []int{1000, 10000, 100000, 1000000}

	for _, size := range sizes {
		arr := make([]int, size)
		for i := range size {
			arr[i] = (size - i) % 100 // 1-2 digit numbers
		}

		b.Run(fmt.Sprintf("small-numbers-size-%d", size), func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				RadixSortOptimization(arr)
			}
		})
	}
}

func BenchmarkRadixSortOptimizationRandomOrder(b *testing.B) {
	// Test with random data
	sizes := []int{1000, 10000, 100000}

	for _, size := range sizes {
		arr := make([]int, size)
		rng := rand.New(rand.NewPCG(42, uint64(size))) // Deterministic seed for reproducible benchmarks
		for i := range size {
			arr[i] = rng.IntN(size * 10) // Random numbers in range [0, size*10)
		}

		b.Run(fmt.Sprintf("random-size-%d", size), func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				RadixSortOptimization(arr)
			}
		})
	}
}
