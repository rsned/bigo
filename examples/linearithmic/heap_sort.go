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

// HeapSort performs O(n log n) heap sort.
// This demonstrates linearithmic time complexity through two phases:
// 1. Build heap: O(n) time
// 2. Extract n elements, each taking O(log n) time: n × log n = O(n log n)
func HeapSort(arr []int) []int {
	// Create a copy to avoid modifying the original array
	result := make([]int, len(arr))
	copy(result, arr)

	// Phase 1: Build max heap - O(n) time
	// Start from last non-leaf node and heapify all internal nodes
	for i := len(result)/2 - 1; i >= 0; i-- {
		heapifyForSort(result, len(result), i)
	}

	// Phase 2: Extract elements from heap one by one - O(n log n) time
	for i := len(result) - 1; i > 0; i-- {
		// Move current root (maximum) to end
		result[0], result[i] = result[i], result[0]
		// Restore heap property for reduced heap - O(log n) operation
		heapifyForSort(result, i, 0)
	}

	return result
}

// heapifyForSort maintains the max heap property by moving an element down.
// This operation is O(log n) because it may traverse the height of the heap.
func heapifyForSort(arr []int, heapSize, rootIndex int) {
	// Assume root is the largest
	largest := rootIndex
	// Calculate children indices
	leftChild := 2*rootIndex + 1
	rightChild := 2*rootIndex + 2

	// Check if left child exists and is larger than root
	if leftChild < heapSize && arr[leftChild] > arr[largest] {
		largest = leftChild
	}

	// Check if right child exists and is larger than current largest
	if rightChild < heapSize && arr[rightChild] > arr[largest] {
		largest = rightChild
	}

	// If largest is not root, swap and continue heapifying
	if largest != rootIndex {
		// Swap root with largest child
		arr[rootIndex], arr[largest] = arr[largest], arr[rootIndex]
		// Recursively heapify the affected subtree
		heapifyForSort(arr, heapSize, largest)
	}
}
