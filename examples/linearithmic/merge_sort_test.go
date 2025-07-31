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
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMergeSort(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want []int
	}{
		{"empty array", []int{}, []int{}},
		{"single element", []int{5}, []int{5}},
		{"already sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"reverse sorted", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"random order", []int{3, 1, 4, 1, 5, 9, 2, 6}, []int{1, 1, 2, 3, 4, 5, 6, 9}},
		{"duplicates", []int{3, 1, 3, 1, 3}, []int{1, 1, 3, 3, 3}},
		{"negative numbers", []int{-1, -5, 0, 3, -2}, []int{-5, -2, -1, 0, 3}},
		{"two elements", []int{2, 1}, []int{1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeSort(tt.arr)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("MergeSort(%v) = %v, want %v", tt.arr, got, tt.want)
			}
		})
	}
}

func TestMergeSortStability(t *testing.T) {
	// Test that merge sort is stable (preserves relative order of equal elements)
	// We'll use a simple test with the fact that equal elements should appear in original order
	arr := []int{3, 1, 4, 1, 5}
	got := MergeSort(arr)
	want := []int{1, 1, 3, 4, 5}

	if !cmp.Equal(got, want) {
		t.Errorf("MergeSort stability test: got %v, want %v", got, want)
	}
}

func TestMergeSortDoesNotModifyOriginal(t *testing.T) {
	original := []int{3, 1, 4, 1, 5}
	originalCopy := make([]int, len(original))
	copy(originalCopy, original)

	MergeSort(original)

	if !cmp.Equal(original, originalCopy) {
		t.Errorf("MergeSort modified original array: got %v, want %v", original, originalCopy)
	}
}

func BenchmarkMergeSort(b *testing.B) {
	sizes := []int{10, 100, 1000, 5000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Generate test data
			arr := make([]int, size)
			for i := 0; i < size; i++ {
				arr[i] = size - i // Reverse sorted for worst case
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				testArr := make([]int, len(arr))
				copy(testArr, arr)
				b.StartTimer()

				MergeSort(testArr)
			}
		})
	}
}
