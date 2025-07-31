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

// CountElements performs O(n) counting of elements.  Yes, this is the dumb
// way of calling len() on a slice.
//
// This demonstrates linear time complexity because we must visit each element
// exactly once to count the total number of elements in the array.
func CountElements(arr []int) int {
	count := 0
	// Visit each element exactly once - this loop runs n times
	for range arr {
		count++ // O(1) operation performed n times = O(n)
	}

	return count
}

// CountOccurrences performs O(n) counting of specific value occurrences.
// This demonstrates linear time complexity because we must examine each element
// to determine if it matches the target value.
func CountOccurrences(arr []int, target int) int {
	count := 0
	// Must check every element - no way to skip elements
	for _, val := range arr {
		// Compare each element to target - O(1) operation
		if val == target {
			count++ // Increment counter when match found
		}
	}

	return count
}

// CountLinkedListElements performs O(n) counting of linked list elements using manual traversal.
// This demonstrates linear time complexity because we must traverse the entire
// linked list to count all nodes, visiting each node exactly once.
// Note: This is deliberately inefficient to demonstrate O(n) traversal - use list.Len() for O(1) counting.
func CountLinkedListElements(list *collection.LinkedList[int]) int {
	count := 0

	// Use iterator to traverse each element - O(n) operation
	iter := list.Iterator()
	for iter.HasNext() {
		_, ok := iter.Next()
		if ok {
			count++ // Count this node - O(1) operation
		}
	}

	return count
}

// CountLinkedListElementsModern performs O(1)* counting using the modern LinkedList API.
// *Note: The modern LinkedList maintains a size counter, making this O(1) instead of O(n).
// This demonstrates how proper data structure design can improve algorithmic complexity.
func CountLinkedListElementsModern(list *collection.LinkedList[int]) int {
	// The modern LinkedList tracks size internally - O(1) operation!
	return list.Len()
}

// CountLinkedListElementsWithCondition demonstrates O(n) conditional counting.
// This shows linear traversal with a condition, which still requires visiting
// each element to apply the predicate function.
func CountLinkedListElementsWithCondition(list *collection.LinkedList[int], predicate func(int) bool) int {
	count := 0

	// Use iterator to traverse and count elements matching condition
	iter := list.Iterator()
	for iter.HasNext() {
		value, ok := iter.Next()
		if ok && predicate(value) {
			count++
		}
	}

	return count
}
