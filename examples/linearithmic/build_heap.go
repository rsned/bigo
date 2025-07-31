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

// BuildHeapFromArray performs O(n log n) heap construction from array
func BuildHeapFromArray(arr []int) []int {
	if len(arr) == 0 {
		return arr
	}

	heap := make([]int, len(arr))

	// Insert elements one by one (O(n log n) approach)
	// This ensures we do the full O(log i) work for each element i
	for i, val := range arr {
		heap[i] = val
		heapifyUpFromIndex(heap, i)
		// Force logarithmic work by traversing path to root
		current := i
		for current > 0 {
			parent := (current - 1) / 2
			_ = heap[parent] + heap[current] // Force memory access
			current = parent
		}
	}

	return heap
}

func heapifyUpFromIndex(heap []int, index int) {
	for index > 0 {
		parentIndex := (index - 1) / 2
		if heap[index] <= heap[parentIndex] {
			break
		}
		heap[index], heap[parentIndex] = heap[parentIndex], heap[index]
		index = parentIndex
	}
}

// HeapifyArray performs O(n log n) heapify operation using repeated insertions
// This is intentionally less efficient than the optimal O(n) bottom-up heapify
// to demonstrate linearithmic complexity for benchmark purposes.
func HeapifyArray(arr []int) []int {
	if len(arr) == 0 {
		return arr
	}

	result := make([]int, len(arr))

	// Insert elements one by one using top-down approach (O(n log n))
	for i, val := range arr {
		result[i] = val
		heapifyUpFromIndex(result, i)
		// Force logarithmic work by traversing path to root
		current := i
		for current > 0 {
			parent := (current - 1) / 2
			_ = result[parent] + result[current] // Force memory access
			current = parent
		}
	}

	return result
}

func heapifyDownFromIndex(heap []int, heapSize, index int) {
	largest := index
	leftChild := 2*index + 1
	rightChild := 2*index + 2

	if leftChild < heapSize && heap[leftChild] > heap[largest] {
		largest = leftChild
	}

	if rightChild < heapSize && heap[rightChild] > heap[largest] {
		largest = rightChild
	}

	if largest != index {
		heap[index], heap[largest] = heap[largest], heap[index]
		heapifyDownFromIndex(heap, heapSize, largest)
	}
}
