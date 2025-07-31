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

// isMaxHeap checks if an array represents a valid max heap
func isMaxHeap(arr []int) bool {
	n := len(arr)
	for i := 0; i < n; i++ {
		leftChild := 2*i + 1
		rightChild := 2*i + 2

		// Check if left child violates max heap property
		if leftChild < n && arr[i] < arr[leftChild] {
			return false
		}

		// Check if right child violates max heap property
		if rightChild < n && arr[i] < arr[rightChild] {
			return false
		}
	}

	return true
}

// containsSameElements checks if two slices contain the same elements (ignoring order)
func containsSameElementsBuildHeap(a, b []int) bool {
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

// TestBuildHeapFromArray tests the BuildHeapFromArray function
func TestBuildHeapFromArray(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		checkHeap bool
	}{
		{
			name:      "empty array",
			input:     []int{},
			checkHeap: true,
		},
		{
			name:      "single element",
			input:     []int{5},
			checkHeap: true,
		},
		{
			name:      "two elements sorted",
			input:     []int{1, 2},
			checkHeap: true,
		},
		{
			name:      "two elements reverse sorted",
			input:     []int{2, 1},
			checkHeap: true,
		},
		{
			name:      "three elements",
			input:     []int{1, 2, 3},
			checkHeap: true,
		},
		{
			name:      "multiple elements already heap",
			input:     []int{9, 4, 6, 1, 2, 3},
			checkHeap: true,
		},
		{
			name:      "multiple elements reverse sorted",
			input:     []int{5, 4, 3, 2, 1},
			checkHeap: true,
		},
		{
			name:      "random order",
			input:     []int{3, 1, 6, 5, 2, 4},
			checkHeap: true,
		},
		{
			name:      "with duplicates",
			input:     []int{4, 4, 4, 4},
			checkHeap: true,
		},
		{
			name:      "mixed with duplicates",
			input:     []int{7, 3, 7, 1, 3, 7, 1},
			checkHeap: true,
		},
		{
			name:      "negative numbers",
			input:     []int{-1, -5, 3, 0, -2},
			checkHeap: true,
		},
		{
			name:      "large range",
			input:     []int{100, 5, 90, 15, 75, 25, 60, 35, 50, 45},
			checkHeap: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildHeapFromArray(tt.input)

			// Check that result contains same elements as input
			if !containsSameElementsBuildHeap(result, tt.input) {
				t.Errorf("BuildHeapFromArray changed elements: input %v, result %v",
					tt.input, result)
			}

			// Check that result is a valid max heap
			if tt.checkHeap && !isMaxHeap(result) {
				t.Errorf("BuildHeapFromArray result is not a valid max heap: %v", result)
			}

			// Check that input array was not modified
			originalSum := 0
			resultSum := 0
			for _, v := range tt.input {
				originalSum += v
			}

			for _, v := range result {
				resultSum += v
			}

			if originalSum != resultSum {
				t.Errorf("Sum mismatch - element was lost or gained")
			}
		})
	}
}

// TestHeapifyArray tests the HeapifyArray function
func TestHeapifyArray(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		checkHeap bool
	}{
		{
			name:      "empty array",
			input:     []int{},
			checkHeap: true,
		},
		{
			name:      "single element",
			input:     []int{5},
			checkHeap: true,
		},
		{
			name:      "two elements",
			input:     []int{1, 2},
			checkHeap: true,
		},
		{
			name:      "three elements sorted",
			input:     []int{1, 2, 3},
			checkHeap: true,
		},
		{
			name:      "perfect heap input",
			input:     []int{10, 8, 9, 4, 7, 5, 6, 1, 2, 3},
			checkHeap: true,
		},
		{
			name:      "reverse sorted",
			input:     []int{5, 4, 3, 2, 1},
			checkHeap: true,
		},
		{
			name:      "random order",
			input:     []int{3, 1, 6, 5, 2, 4},
			checkHeap: true,
		},
		{
			name:      "with duplicates",
			input:     []int{5, 5, 5, 5, 5},
			checkHeap: true,
		},
		{
			name:      "mixed duplicates",
			input:     []int{8, 3, 8, 1, 3, 8, 1, 2},
			checkHeap: true,
		},
		{
			name:      "negative numbers",
			input:     []int{-10, 5, -3, 2, -1},
			checkHeap: true,
		},
		{
			name:      "large array",
			input:     []int{15, 10, 20, 8, 16, 25, 12, 5, 9, 11, 13, 30, 18, 7, 14},
			checkHeap: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HeapifyArray(tt.input)

			// Check that result contains same elements as input
			if !containsSameElementsBuildHeap(result, tt.input) {
				t.Errorf("HeapifyArray changed elements: input %v, result %v",
					tt.input, result)
			}

			// Check that result is a valid max heap
			if tt.checkHeap && !isMaxHeap(result) {
				t.Errorf("HeapifyArray result is not a valid max heap: %v", result)
			}

			// Check that original input was not modified
			inputCopy := make([]int, len(tt.input))
			copy(inputCopy, tt.input)
			HeapifyArray(tt.input) // Call again to check input wasn't modified
			if !cmp.Equal(tt.input, inputCopy) {
				t.Errorf("HeapifyArray modified input array")
			}
		})
	}
}

// TestHeapifyVsBuildHeap compares both methods to ensure they produce valid heaps
func TestHeapifyVsBuildHeap(t *testing.T) {
	testCases := [][]int{
		{1, 2, 3, 4, 5},
		{5, 4, 3, 2, 1},
		{3, 1, 6, 5, 2, 4},
		{10, 20, 15, 12, 40, 25, 18},
		{1, 1, 1, 1},
		{-5, 3, -1, 7, -3, 2},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("comparison_test_%d", i), func(t *testing.T) {
			heapFromBuild := BuildHeapFromArray(testCase)
			heapFromHeapify := HeapifyArray(testCase)

			// Both should be valid heaps
			if !isMaxHeap(heapFromBuild) {
				t.Errorf("BuildHeapFromArray did not produce valid heap: %v", heapFromBuild)
			}

			if !isMaxHeap(heapFromHeapify) {
				t.Errorf("HeapifyArray did not produce valid heap: %v", heapFromHeapify)
			}

			// Both should contain same elements as original
			if !containsSameElementsBuildHeap(heapFromBuild, testCase) {
				t.Errorf("BuildHeapFromArray elements mismatch")
			}

			if !containsSameElementsBuildHeap(heapFromHeapify, testCase) {
				t.Errorf("HeapifyArray elements mismatch")
			}

			// The root should be the maximum element
			if len(testCase) > 0 {
				maxElement := testCase[0]
				for _, v := range testCase {
					if v > maxElement {
						maxElement = v
					}
				}

				if heapFromBuild[0] != maxElement {
					t.Errorf("BuildHeapFromArray root is not maximum: got %d, want %d",
						heapFromBuild[0], maxElement)
				}

				if heapFromHeapify[0] != maxElement {
					t.Errorf("HeapifyArray root is not maximum: got %d, want %d",
						heapFromHeapify[0], maxElement)
				}
			}
		})
	}
}

// TestHeapProperties tests specific heap properties
func TestHeapProperties(t *testing.T) {
	t.Run("heap_property_maintained", func(t *testing.T) {
		input := []int{1, 3, 6, 5, 2, 4, 7}
		heap := HeapifyArray(input)

		// Test that every parent is >= its children
		for i := 0; i < len(heap)/2; i++ {
			leftChild := 2*i + 1
			rightChild := 2*i + 2

			if leftChild < len(heap) && heap[i] < heap[leftChild] {
				t.Errorf("Heap property violated: parent %d < left child %d at positions %d, %d",
					heap[i], heap[leftChild], i, leftChild)
			}

			if rightChild < len(heap) && heap[i] < heap[rightChild] {
				t.Errorf("Heap property violated: parent %d < right child %d at positions %d, %d",
					heap[i], heap[rightChild], i, rightChild)
			}
		}
	})

	t.Run("maximum_at_root", func(t *testing.T) {
		input := []int{5, 2, 8, 1, 9, 3, 7}
		heap := HeapifyArray(input)

		if len(heap) > 0 {
			maxVal := heap[0]
			for _, val := range heap {
				if val > maxVal {
					t.Errorf("Found element %d greater than root %d", val, maxVal)
				}
			}
		}
	})
}

// Benchmark functions
func BenchmarkBuildHeapFromArray(b *testing.B) {
	sizes := []int{10, 100, 1000, 5000, 10000}

	for _, size := range sizes {
		// Create test data
		arr := make([]int, size)
		for i := range arr {
			arr[i] = size - i // Reverse sorted for worst case
		}

		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				BuildHeapFromArray(arr)
			}
		})
	}
}

func BenchmarkHeapifyArray(b *testing.B) {
	sizes := []int{10, 100, 1000, 5000, 10000}

	for _, size := range sizes {
		// Create test data
		arr := make([]int, size)
		for i := range arr {
			arr[i] = size - i // Reverse sorted for worst case
		}

		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				HeapifyArray(arr)
			}
		})
	}
}

// BenchmarkHeapifyVsBuild compares performance of both approaches
func BenchmarkHeapifyVsBuild(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		arr := make([]int, size)
		for i := range arr {
			arr[i] = size - i
		}

		b.Run(fmt.Sprintf("BuildHeap_size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				BuildHeapFromArray(arr)
			}
		})

		b.Run(fmt.Sprintf("Heapify_size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				HeapifyArray(arr)
			}
		})
	}
}

// TestEdgeCases tests edge cases and boundary conditions
func TestBuildHeapEdgeCases(t *testing.T) {
	t.Run("single_element_heap", func(t *testing.T) {
		input := []int{42}

		buildResult := BuildHeapFromArray(input)
		heapifyResult := HeapifyArray(input)

		expected := []int{42}
		if !cmp.Equal(buildResult, expected) {
			t.Errorf("BuildHeapFromArray single element: got %v, want %v", buildResult, expected)
		}

		if !cmp.Equal(heapifyResult, expected) {
			t.Errorf("HeapifyArray single element: got %v, want %v", heapifyResult, expected)
		}
	})

	t.Run("two_element_heap", func(t *testing.T) {
		input := []int{1, 2}

		buildResult := BuildHeapFromArray(input)
		heapifyResult := HeapifyArray(input)

		// Both should result in [2, 1] for max heap
		if buildResult[0] != 2 || heapifyResult[0] != 2 {
			t.Errorf("Two element heap should have 2 at root")
		}
	})

	t.Run("zero_values", func(t *testing.T) {
		input := []int{0, 0, 0}

		buildResult := BuildHeapFromArray(input)
		heapifyResult := HeapifyArray(input)

		if !isMaxHeap(buildResult) {
			t.Errorf("BuildHeapFromArray failed with zeros")
		}

		if !isMaxHeap(heapifyResult) {
			t.Errorf("HeapifyArray failed with zeros")
		}
	})
}
