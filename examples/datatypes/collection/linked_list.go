// Copyright 2025 Robert Snedegar
//
// Licensed under the Apache License, Version 2.0 (the License);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collection

import (
	"errors"
	"fmt"
)

// node represents the nodes used in a singly linked list.
type node[T comparable] struct {
	value T
	next  *node[T]
}

// LinkedList represents a generic singly linked list with head and tail pointers.
// By maintaining both head and tail pointers, we can achieve O(1) access
// to both the first and last elements, avoiding O(n) traversal for append operations.
type LinkedList[T comparable] struct {
	head *node[T] // Pointer to the first node for O(1) access
	tail *node[T] // Pointer to the last node for O(1) access
	size int      // Number of elements for O(1) length operations
}

// NewLinkedList creates a new empty singly linked list.
func NewLinkedList[T comparable]() *LinkedList[T] {
	return &LinkedList[T]{
		head: nil,
		tail: nil,
		size: 0,
	}
}

// NewLinkedListWithValues creates a new singly linked list with initial values.
func NewLinkedListWithValues[T comparable](values ...T) *LinkedList[T] {
	ll := NewLinkedList[T]()
	for _, value := range values {
		ll.PushBack(value)
	}

	return ll
}

// FromSlice creates a singly linked list from a slice.
func FromSlice[T comparable](slice []T) *LinkedList[T] {
	return NewLinkedListWithValues(slice...)
}

// PushFront adds an element to the beginning of the list - O(1).
func (ll *LinkedList[T]) PushFront(value T) {
	newNode := &node[T]{value: value, next: ll.head}
	ll.head = newNode

	// If list was empty, update tail pointer
	if ll.tail == nil {
		ll.tail = newNode
	}

	ll.size++
}

// PushBack adds an element to the end of the list - O(1).
func (ll *LinkedList[T]) PushBack(value T) {
	newNode := &node[T]{value: value, next: nil}

	if ll.tail == nil {
		// Empty list case
		ll.head = newNode
		ll.tail = newNode
	} else {
		ll.tail.next = newNode
		ll.tail = newNode
	}

	ll.size++
}

// PopFront removes and returns the first element - O(1).
func (ll *LinkedList[T]) PopFront() (T, bool) {
	if ll.head == nil {
		var zero T

		return zero, false
	}

	value := ll.head.value
	ll.head = ll.head.next

	// If list becomes empty, update tail pointer
	if ll.head == nil {
		ll.tail = nil
	}

	ll.size--

	return value, true
}

// PopBack removes and returns the last element - O(n) for singly linked list.
// Note: This is a limitation of singly linked lists. DoublyLinkedList provides O(1).
func (ll *LinkedList[T]) PopBack() (T, bool) {
	if ll.head == nil {
		var zero T

		return zero, false
	}

	// Single element case
	if ll.head == ll.tail {
		value := ll.head.value
		ll.head = nil
		ll.tail = nil
		ll.size--

		return value, true
	}

	// Find the second-to-last node
	current := ll.head
	for current.next != ll.tail {
		current = current.next
	}

	value := ll.tail.value
	current.next = nil
	ll.tail = current
	ll.size--

	return value, true
}

// Front returns the first element without removing it - O(1).
func (ll *LinkedList[T]) Front() (T, bool) {
	if ll.head == nil {
		var zero T

		return zero, false
	}

	return ll.head.value, true
}

// Back returns the last element without removing it - O(1).
func (ll *LinkedList[T]) Back() (T, bool) {
	if ll.tail == nil {
		var zero T

		return zero, false
	}

	return ll.tail.value, true
}

// At returns the element at the specified index - O(n).
func (ll *LinkedList[T]) At(index int) (T, bool) {
	var zero T

	if index < 0 || index >= ll.size {
		return zero, false
	}

	current := ll.head
	for range index {
		current = current.next
	}

	return current.value, true
}

// Insert adds an element at the specified index - O(n).
func (ll *LinkedList[T]) Insert(index int, value T) error {
	if index < 0 || index > ll.size {
		return errors.New("index out of bounds")
	}

	// Insert at beginning
	if index == 0 {
		ll.PushFront(value)

		return nil
	}

	// Insert at end
	if index == ll.size {
		ll.PushBack(value)

		return nil
	}

	// Insert in middle
	current := ll.head
	for range index - 1 {
		current = current.next
	}

	newNode := &node[T]{value: value, next: current.next}
	current.next = newNode
	ll.size++

	return nil
}

// Remove removes the element at the specified index - O(n).
func (ll *LinkedList[T]) Remove(index int) (T, error) {
	var zero T

	if index < 0 || index >= ll.size {
		return zero, errors.New("index out of bounds")
	}

	// Remove from beginning
	if index == 0 {
		value, _ := ll.PopFront()

		return value, nil
	}

	// Remove from end
	if index == ll.size-1 {
		value, _ := ll.PopBack()

		return value, nil
	}

	// Remove from middle
	current := ll.head
	for range index - 1 {
		current = current.next
	}

	nodeToRemove := current.next
	value := nodeToRemove.value
	current.next = nodeToRemove.next
	ll.size--

	return value, nil
}

// Find returns the index of the first occurrence of the value, or -1 if not found - O(n).
func (ll *LinkedList[T]) Find(value T) int {
	current := ll.head
	for i := 0; current != nil; i++ {
		if current.value == value {
			return i
		}
		current = current.next
	}

	return -1
}

// Contains checks if the list contains the specified value - O(n).
func (ll *LinkedList[T]) Contains(value T) bool {
	return ll.Find(value) != -1
}

// Len returns the number of elements in the list - O(1).
func (ll *LinkedList[T]) Len() int {
	return ll.size
}

// IsEmpty returns true if the list is empty - O(1).
func (ll *LinkedList[T]) IsEmpty() bool {
	return ll.size == 0
}

// Clear removes all elements from the list - O(1).
func (ll *LinkedList[T]) Clear() {
	ll.head = nil
	ll.tail = nil
	ll.size = 0
}

// ToSlice converts the linked list to a slice - O(n).
func (ll *LinkedList[T]) ToSlice() []T {
	result := make([]T, ll.size)
	current := ll.head

	for i := 0; current != nil; i++ {
		result[i] = current.value
		current = current.next
	}

	return result
}

// String returns a string representation of the list.
func (ll *LinkedList[T]) String() string {
	return fmt.Sprintf("LinkedList%v", ll.ToSlice())
}

// Iterator returns a forward iterator for the list.
func (ll *LinkedList[T]) Iterator() Iterator[T] {
	return &linkedListIterator[T]{
		current: ll.head,
	}
}

// linkedListIterator implements Iterator for LinkedList.
type linkedListIterator[T comparable] struct {
	current *node[T]
}

// HasNext returns true if there are more elements to iterate over.
func (it *linkedListIterator[T]) HasNext() bool {
	return it.current != nil
}

// Next returns the next element and advances the iterator.
func (it *linkedListIterator[T]) Next() (T, bool) {
	if it.current == nil {
		var zero T

		return zero, false
	}

	value := it.current.value
	it.current = it.current.next

	return value, true
}

// Value returns the current element without advancing the iterator.
func (it *linkedListIterator[T]) Value() T {
	if it.current == nil {
		var zero T

		return zero
	}

	return it.current.value
}
