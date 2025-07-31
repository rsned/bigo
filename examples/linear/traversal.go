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

import "github.com/rsned/bigo/examples/datatypes/collection"

// ArrayTraversal performs O(n) traversal of array.
// This demonstrates linear time complexity because we must visit each element
// exactly once to transform it and store the result.
func ArrayTraversal(arr []int) []int {
	// Pre-allocate result slice with same length - O(1) operation
	result := make([]int, len(arr))
	// Visit each element exactly once
	for i, val := range arr {
		// Transform element and store in result - O(1) operation
		result[i] = val * 2
	}

	return result
}

// LinkedListTraversal performs O(n) traversal of linked list using iterator.
// This demonstrates linear time complexity because we must visit each node
// exactly once to collect its value, using the iterator pattern for safety.
func LinkedListTraversal(list *collection.LinkedList[int]) []int {
	var result []int

	// Use iterator to traverse the list - O(n) operation
	iter := list.Iterator()
	for iter.HasNext() {
		value, ok := iter.Next()
		if ok {
			// Add current node's value to result - O(1) amortized
			result = append(result, value)
		}
	}

	return result
}

// LinkedListTraversalModern performs O(n) traversal using the modern LinkedList API.
// This demonstrates the same linear time complexity but with a cleaner, safer API
// that encapsulates the node structure and provides iterator-based traversal.
func LinkedListTraversalModern(list *collection.LinkedList[int]) []int {
	// Option 1: Use ToSlice() - this is the most direct approach
	return list.ToSlice()
}

// LinkedListTraversalWithIterator demonstrates O(n) traversal using iterators.
// This shows how to traverse the list element by element using the iterator pattern,
// which is useful when you need to process elements individually or apply
// early termination conditions.
func LinkedListTraversalWithIterator(list *collection.LinkedList[int]) []int {
	// Initialize result as empty slice to ensure consistency
	result := []int{}

	// Get an iterator for the list
	iter := list.Iterator()

	// Traverse using iterator - each HasNext() and Next() is O(1)
	for iter.HasNext() {
		value, ok := iter.Next()
		if ok {
			// Transform element during traversal (example: double the value)
			result = append(result, value*2)
		}
	}

	return result
}
