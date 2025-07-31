// Copyright 2025 Robert Snedegar
//
// Licensed under the Apache License, Version 2.0 (the License);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an AS IS BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package collection provides generic collection data structures optimized for
Big O complexity analysis and education. The package demonstrates modern Go
practices with generics while maintaining backward compatibility.

# Core Types

The package provides two main generic collection types:

  - LinkedList[T comparable]: A singly-linked list with O(1) front operations
    and size tracking for O(1) length queries.

  - DoublyLinkedList[T comparable]: A doubly-linked list with O(1) operations
    at both ends and optimized bidirectional traversal.

# Usage Examples

Creating and using a generic LinkedList:

	// Create from slice
	list := collection.FromSlice([]int{1, 2, 3})

	// Modern API
	list.PushFront(0)       // O(1)
	val, ok := list.PopBack() // O(n) - demonstrates singly-linked limitation

	// Iterator traversal
	iter := list.Iterator()
	for iter.HasNext() {
		value, _ := iter.Next()
		// process value
	}

Creating and using a DoublyLinkedList:

	// Create with values
	dll := collection.NewDoublyLinkedListWithValues(1, 2, 3)

	// Bidirectional O(1) operations
	dll.PushBack(4)         // O(1)
	val, ok := dll.PopBack() // O(1)

	// Optimized indexed access
	val, ok := dll.At(index) // O(min(index, size-index))

# Performance Characteristics

LinkedList[T]:
  - PushFront/PopFront: O(1)
  - PushBack: O(1) with tail pointer
  - PopBack: O(n) - requires full traversal
  - Find/Contains: O(n)
  - Len(): O(1) with size tracking

DoublyLinkedList[T]:
  - All front/back operations: O(1)
  - At(index): O(min(index, size-index))
  - Insert/Remove: O(min(index, size-index))
  - Find/Contains: O(n)
  - Reverse iteration: O(1) per step
*/
package collection
