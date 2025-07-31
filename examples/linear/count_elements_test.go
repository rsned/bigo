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

	"github.com/rsned/bigo/examples/datatypes/collection"
)

func TestCountElements(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want int
	}{
		{"empty array", []int{}, 0},
		{"single element", []int{5}, 1},
		{"multiple elements", []int{1, 2, 3, 4, 5}, 5},
		{"ten elements", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountElements(tt.arr)
			if got != tt.want {
				t.Errorf("CountElements(%v) = %d, want %d", tt.arr, got, tt.want)
			}
		})
	}
}

func TestCountOccurrences(t *testing.T) {
	tests := []struct {
		name   string
		arr    []int
		target int
		want   int
	}{
		{"empty array", []int{}, 5, 0},
		{"target not found", []int{1, 2, 3}, 5, 0},
		{"single occurrence", []int{1, 2, 3, 4, 5}, 3, 1},
		{"multiple occurrences", []int{1, 2, 2, 3, 2}, 2, 3},
		{"all elements same", []int{5, 5, 5, 5}, 5, 4},
		{"target zero", []int{0, 1, 0, 2, 0}, 0, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountOccurrences(tt.arr, tt.target)
			if got != tt.want {
				t.Errorf("CountOccurrences(%v, %d) = %d, want %d", tt.arr, tt.target, got, tt.want)
			}
		})
	}
}

func TestCountLinkedListElements(t *testing.T) {
	tests := []struct {
		name string
		list *collection.LinkedList[int]
		want int
	}{
		{"empty list", collection.NewLinkedList[int](), 0},
		{"single node", collection.NewLinkedListWithValues(5), 1},
		{"multiple nodes", collection.FromSlice([]int{1, 2, 3, 4}), 4},
		{"long list", collection.FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}), 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountLinkedListElements(tt.list)
			if got != tt.want {
				t.Errorf("CountLinkedListElements() = %d, want %d", got, tt.want)
			}
		})
	}
}

func BenchmarkCountElements(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Generate test data
			arr := make([]int, size)
			for i := 0; i < size; i++ {
				arr[i] = i + 1
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				CountElements(arr)
			}
		})
	}
}

func BenchmarkCountOccurrences(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Generate test data with repeated target
			arr := make([]int, size)
			target := 42
			for i := 0; i < size; i++ {
				if i%10 == 0 {
					arr[i] = target // Every 10th element is the target
				} else {
					arr[i] = i
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				CountOccurrences(arr, target)
			}
		})
	}
}

func TestCountLinkedListElementsModern(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   int
	}{
		{"empty list", []int{}, 0},
		{"single element", []int{5}, 1},
		{"multiple elements", []int{1, 2, 3, 4}, 4},
		{"long list", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := collection.FromSlice(tt.values)
			got := CountLinkedListElementsModern(list)
			if got != tt.want {
				t.Errorf("CountLinkedListElementsModern() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCountLinkedListElementsWithCondition(t *testing.T) {
	tests := []struct {
		name      string
		values    []int
		predicate func(int) bool
		want      int
	}{
		{"empty list", []int{}, func(x int) bool { return x > 0 }, 0},
		{"count positive", []int{-1, 0, 1, 2, -3}, func(x int) bool { return x > 0 }, 2},
		{"count even", []int{1, 2, 3, 4, 5, 6}, func(x int) bool { return x%2 == 0 }, 3},
		{"count all", []int{1, 2, 3}, func(_ int) bool { return true }, 3},
		{"count none", []int{1, 2, 3}, func(_ int) bool { return false }, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := collection.FromSlice(tt.values)
			got := CountLinkedListElementsWithCondition(list, tt.predicate)
			if got != tt.want {
				t.Errorf("CountLinkedListElementsWithCondition() = %d, want %d", got, tt.want)
			}
		})
	}
}

func BenchmarkCountLinkedListElements(b *testing.B) {
	sizes := []int{100, 1000, 10000, 50000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Generate test data
			values := make([]int, size)
			for i := 0; i < size; i++ {
				values[i] = i + 1
			}
			list := collection.FromSlice(values)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				CountLinkedListElements(list)
			}
		})
	}
}

func BenchmarkCountLinkedListElementsModern(b *testing.B) {
	sizes := []int{100, 1000, 10000, 50000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Generate test data
			values := make([]int, size)
			for i := 0; i < size; i++ {
				values[i] = i + 1
			}
			list := collection.FromSlice(values)

			b.ResetTimer()
			for b.Loop() {
				CountLinkedListElementsModern(list) // This should be O(1)!
			}
		})
	}
}

func BenchmarkCountLinkedListElementsWithCondition(b *testing.B) {
	sizes := []int{100, 1000, 10000, 50000}
	predicate := func(x int) bool { return x%2 == 0 } // Count even numbers

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Generate test data
			values := make([]int, size)
			for i := 0; i < size; i++ {
				values[i] = i + 1
			}
			list := collection.FromSlice(values)

			b.ResetTimer()
			for b.Loop() {
				CountLinkedListElementsWithCondition(list, predicate)
			}
		})
	}
}
