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

import (
	"fmt"
	"testing"
)

func TestFindMinimum(t *testing.T) {
	tests := []struct {
		name        string
		arr         []int
		wantValue   int
		wantSuccess bool
	}{
		{"empty array", []int{}, 0, false},
		{"single element", []int{5}, 5, true},
		{"positive numbers", []int{3, 1, 4, 1, 5}, 1, true},
		{"negative numbers", []int{-3, -1, -4, -1, -5}, -5, true},
		{"mixed numbers", []int{-3, 0, 4, -1, 5}, -3, true},
		{"all same", []int{7, 7, 7, 7}, 7, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotSuccess := FindMinimum(tt.arr)
			if gotValue != tt.wantValue || gotSuccess != tt.wantSuccess {
				t.Errorf("FindMinimum(%v) = (%d, %t), want (%d, %t)",
					tt.arr, gotValue, gotSuccess, tt.wantValue, tt.wantSuccess)
			}
		})
	}
}

func TestFindMaximum(t *testing.T) {
	tests := []struct {
		name        string
		arr         []int
		wantValue   int
		wantSuccess bool
	}{
		{"empty array", []int{}, 0, false},
		{"single element", []int{5}, 5, true},
		{"positive numbers", []int{3, 1, 4, 1, 5}, 5, true},
		{"negative numbers", []int{-3, -1, -4, -1, -5}, -1, true},
		{"mixed numbers", []int{-3, 0, 4, -1, 5}, 5, true},
		{"all same", []int{7, 7, 7, 7}, 7, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotSuccess := FindMaximum(tt.arr)
			if gotValue != tt.wantValue || gotSuccess != tt.wantSuccess {
				t.Errorf("FindMaximum(%v) = (%d, %t), want (%d, %t)",
					tt.arr, gotValue, gotSuccess, tt.wantValue, tt.wantSuccess)
			}
		})
	}
}

func TestFindMinMax(t *testing.T) {
	tests := []struct {
		name        string
		arr         []int
		wantMin     int
		wantMax     int
		wantSuccess bool
	}{
		{"empty array", []int{}, 0, 0, false},
		{"single element", []int{5}, 5, 5, true},
		{"positive numbers", []int{3, 1, 4, 1, 5}, 1, 5, true},
		{"negative numbers", []int{-3, -1, -4, -1, -5}, -5, -1, true},
		{"mixed numbers", []int{-3, 0, 4, -1, 5}, -3, 5, true},
		{"all same", []int{7, 7, 7, 7}, 7, 7, true},
		{"two elements", []int{10, 3}, 3, 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax, gotSuccess := FindMinMax(tt.arr)
			if gotMin != tt.wantMin || gotMax != tt.wantMax || gotSuccess != tt.wantSuccess {
				t.Errorf("FindMinMax(%v) = (%d, %d, %t), want (%d, %d, %t)",
					tt.arr, gotMin, gotMax, gotSuccess, tt.wantMin, tt.wantMax, tt.wantSuccess)
			}
		})
	}
}

func BenchmarkFindMinimum(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Generate test data with minimum at the end for worst case
			arr := make([]int, size)
			for i := 0; i < size-1; i++ {
				arr[i] = i + 100
			}
			arr[size-1] = 1 // Minimum at the end

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				FindMinimum(arr)
			}
		})
	}
}

func BenchmarkFindMaximum(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Generate test data with maximum at the end for worst case
			arr := make([]int, size)
			for i := 0; i < size-1; i++ {
				arr[i] = i + 1
			}
			arr[size-1] = size * 10 // Maximum at the end

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				FindMaximum(arr)
			}
		})
	}
}

func BenchmarkFindMinMax(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Generate test data with random values
			arr := make([]int, size)
			for i := 0; i < size; i++ {
				arr[i] = (i*13 + 7) % 1000 // Some pseudo-random pattern
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				FindMinMax(arr)
			}
		})
	}
}
