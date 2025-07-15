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

package constant

// LinkedList represents a linked list with head and tail pointers.
// By maintaining both head and tail pointers, we can achieve O(1) access
// to both the first and last elements, avoiding O(n) traversal.
type LinkedList struct {
	Head *ListNode // Pointer to the first node for O(1) access
	Tail *ListNode // Pointer to the last node for O(1) access
}

// GetFirstElement performs O(1) access to first element.
// This demonstrates constant time complexity because we have a direct
// pointer to the head node, eliminating the need for traversal.
func (ll *LinkedList) GetFirstElement() (int, bool) {
	// Check if the list is empty
	if ll.Head == nil {
		return 0, false
	}

	// Direct access to head node's value - O(1) operation
	return ll.Head.Val, true
}

// GetLastElement performs O(1) access to last element.
// This demonstrates constant time complexity because we maintain a tail pointer,
// avoiding the need to traverse the entire list to find the last element.
func (ll *LinkedList) GetLastElement() (int, bool) {
	// Check if the list is empty
	if ll.Tail == nil {
		return 0, false
	}

	// Direct access to tail node's value - O(1) operation
	return ll.Tail.Val, true
}

// ListNode represents a node in a linked list.
// Each node contains a value and a pointer to the next node.
type ListNode struct {
	Val  int       // The data stored in this node
	Next *ListNode // Pointer to the next node in the list
}

// NewListNode creates a new ListNode with the given value.
// This is a helper function to simplify node creation.
func NewListNode(val int) *ListNode {
	return &ListNode{
		Val:  val, // Set the node's value
		Next: nil, // Initialize next pointer to nil
	}
}

// BuildLinkedList creates a linked list from a slice of integers.
// This utility function helps create test data for linked list operations.
func BuildLinkedList(values []int) *ListNode {
	// Handle empty slice case
	if len(values) == 0 {
		return nil
	}

	// Create the head node with the first value
	head := NewListNode(values[0])
	current := head

	// Link remaining values as subsequent nodes
	for i := 1; i < len(values); i++ {
		current.Next = NewListNode(values[i])
		current = current.Next // Move to the newly created node
	}

	return head
}

// ToSlice converts a linked list to a slice of integers.
// This utility function helps with testing and debugging linked
// list operations.
//
// Note: This operation is O(n) as it must traverse the entire list.
func (ln *ListNode) ToSlice() []int {
	var result []int
	current := ln

	// Traverse the linked list from head to tail
	for current != nil {
		result = append(result, current.Val) // Add current value to result
		current = current.Next               // Move to next node
	}

	return result
}
