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

	"github.com/google/go-cmp/cmp"
	"github.com/rsned/bigo/examples/datatypes/collection"
)

func TestArrayTraversal(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want []int
	}{
		{"empty array", []int{}, []int{}},
		{"single element", []int{5}, []int{10}},
		{"multiple elements", []int{1, 2, 3, 4}, []int{2, 4, 6, 8}},
		{"negative numbers", []int{-1, -2, 0, 3}, []int{-2, -4, 0, 6}},
		{"zero", []int{0}, []int{0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ArrayTraversal(tt.arr)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("ArrayTraversal(%v) = %v, want %v", tt.arr, got, tt.want)
			}
		})
	}
}

func TestLinkedListTraversal(t *testing.T) {
	tests := []struct {
		name string
		list *collection.LinkedList[int]
		want []int
	}{
		{"empty list", collection.NewLinkedList[int](), []int{}},
		{"single node", collection.NewLinkedListWithValues(5), []int{5}},
		{"multiple nodes", collection.FromSlice([]int{1, 2, 3, 4}), []int{1, 2, 3, 4}},
		{"negative values", collection.FromSlice([]int{-1, 0, 1}), []int{-1, 0, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LinkedListTraversal(tt.list)
			if len(got) != len(tt.want) {
				t.Errorf("LinkedListTraversal() length = %d, want %d", len(got), len(tt.want))

				return
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("LinkedListTraversal()[%d] = %d, want %d", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestArrayTraversalDoesNotModifyOriginal(t *testing.T) {
	original := []int{1, 2, 3}
	originalCopy := make([]int, len(original))
	copy(originalCopy, original)

	ArrayTraversal(original)

	if !cmp.Equal(original, originalCopy) {
		t.Errorf("ArrayTraversal modified original array: got %v, want %v", original, originalCopy)
	}
}

func BenchmarkArrayTraversal(b *testing.B) {
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
				ArrayTraversal(arr)
			}
		})
	}
}

func TestLinkedListTraversalModern(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{"empty list", []int{}, []int{}},
		{"single element", []int{5}, []int{5}},
		{"multiple elements", []int{1, 2, 3, 4}, []int{1, 2, 3, 4}},
		{"negative values", []int{-1, 0, 1}, []int{-1, 0, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := collection.FromSlice(tt.values)
			got := LinkedListTraversalModern(list)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("LinkedListTraversalModern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkedListTraversalWithIterator(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int // Expected doubled values
	}{
		{"empty list", []int{}, []int{}},
		{"single element", []int{5}, []int{10}},
		{"multiple elements", []int{1, 2, 3, 4}, []int{2, 4, 6, 8}},
		{"negative values", []int{-1, 0, 1}, []int{-2, 0, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := collection.FromSlice(tt.values)
			got := LinkedListTraversalWithIterator(list)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("LinkedListTraversalWithIterator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkLinkedListTraversal(b *testing.B) {
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
				LinkedListTraversal(list)
			}
		})
	}
}

func BenchmarkLinkedListTraversalModern(b *testing.B) {
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
				LinkedListTraversalModern(list)
			}
		})
	}
}

func BenchmarkLinkedListTraversalWithIterator(b *testing.B) {
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
				LinkedListTraversalWithIterator(list)
			}
		})
	}
}
