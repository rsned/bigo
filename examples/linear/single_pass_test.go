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
	"math"
	"testing"
)

func TestCalculateSum(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want int
	}{
		{"empty array", []int{}, 0},
		{"single element", []int{5}, 5},
		{"positive numbers", []int{1, 2, 3, 4, 5}, 15},
		{"negative numbers", []int{-1, -2, -3}, -6},
		{"mixed numbers", []int{-2, 0, 2}, 0},
		{"zeros", []int{0, 0, 0}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateSum(tt.arr)
			if got != tt.want {
				t.Errorf("CalculateSum(%v) = %d, want %d", tt.arr, got, tt.want)
			}
		})
	}
}

func TestCalculateAverage(t *testing.T) {
	tests := []struct {
		name        string
		arr         []int
		wantValue   float64
		wantSuccess bool
	}{
		{"empty array", []int{}, 0, false},
		{"single element", []int{6}, 6.0, true},
		{"positive numbers", []int{2, 4, 6}, 4.0, true},
		{"negative numbers", []int{-3, -6, -9}, -6.0, true},
		{"mixed numbers", []int{-1, 0, 1}, 0.0, true},
		{"integer division", []int{1, 2}, 1.5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotSuccess := CalculateAverage(tt.arr)
			if gotSuccess != tt.wantSuccess {
				t.Errorf("CalculateAverage(%v) success = %t, want %t", tt.arr, gotSuccess, tt.wantSuccess)

				return
			}

			if gotSuccess && math.Abs(gotValue-tt.wantValue) > 1e-9 {
				t.Errorf("CalculateAverage(%v) = %f, want %f", tt.arr, gotValue, tt.wantValue)
			}
		})
	}
}

func TestCalculateProduct(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want int
	}{
		{"empty array", []int{}, 0},
		{"single element", []int{5}, 5},
		{"positive numbers", []int{2, 3, 4}, 24},
		{"with zero", []int{2, 0, 4}, 0},
		{"negative numbers", []int{-2, 3}, -6},
		{"two negatives", []int{-2, -3}, 6},
		{"ones", []int{1, 1, 1}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateProduct(tt.arr)
			if got != tt.want {
				t.Errorf("CalculateProduct(%v) = %d, want %d", tt.arr, got, tt.want)
			}
		})
	}
}

func TestCountEvenOdd(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		wantEven int
		wantOdd  int
	}{
		{"empty array", []int{}, 0, 0},
		{"single even", []int{2}, 1, 0},
		{"single odd", []int{3}, 0, 1},
		{"mixed numbers", []int{1, 2, 3, 4, 5}, 2, 3},
		{"all even", []int{2, 4, 6, 8}, 4, 0},
		{"all odd", []int{1, 3, 5, 7}, 0, 4},
		{"with zero", []int{0, 1, 2, 3}, 2, 2},
		{"negative numbers", []int{-2, -1, 0, 1, 2}, 3, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEven, gotOdd := CountEvenOdd(tt.arr)
			if gotEven != tt.wantEven || gotOdd != tt.wantOdd {
				t.Errorf("CountEvenOdd(%v) = (%d, %d), want (%d, %d)",
					tt.arr, gotEven, gotOdd, tt.wantEven, tt.wantOdd)
			}
		})
	}
}

// Benchmark functions for linear single pass operations

func BenchmarkCalculateSum(b *testing.B) {
	// Test with various sizes
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		arr := make([]int, size)
		for i := 0; i < size; i++ {
			arr[i] = i + 1
		}

		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				CalculateSum(arr)
			}
		})
	}
}

func BenchmarkCalculateAverage(b *testing.B) {
	// Test with various sizes
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		arr := make([]int, size)
		for i := 0; i < size; i++ {
			arr[i] = i + 1
		}

		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				CalculateAverage(arr)
			}
		})
	}
}

func BenchmarkCalculateProduct(b *testing.B) {
	// Test with smaller sizes to avoid overflow
	sizes := []int{10, 20, 30}

	for _, size := range sizes {
		arr := make([]int, size)
		for i := 0; i < size; i++ {
			arr[i] = 2 // Small values to avoid overflow
		}

		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				CalculateProduct(arr)
			}
		})
	}
}

func BenchmarkCountEvenOdd(b *testing.B) {
	// Test with various sizes
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		arr := make([]int, size)
		for i := 0; i < size; i++ {
			arr[i] = i + 1
		}

		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				CountEvenOdd(arr)
			}
		})
	}
}
