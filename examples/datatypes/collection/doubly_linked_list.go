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

// doublyLinkedNode represents the nodes used in a doubly linked list.
// Made private to improve encapsulation.
type doublyLinkedNode[T comparable] struct {
	value T
	next  *doublyLinkedNode[T]
	prev  *doublyLinkedNode[T]
}

// DoublyLinkedList represents a generic doubly linked list with head and tail pointers.
// The main advantage over singly linked lists is O(1) removal from the end and
// efficient bidirectional traversal.
type DoublyLinkedList[T comparable] struct {
	head *doublyLinkedNode[T] // Pointer to the first node for O(1) access
	tail *doublyLinkedNode[T] // Pointer to the last node for O(1) access
	size int                  // Number of elements for O(1) length operations
}

// NewDoublyLinkedList creates a new empty doubly linked list.
func NewDoublyLinkedList[T comparable]() *DoublyLinkedList[T] {
	return &DoublyLinkedList[T]{
		head: nil,
		tail: nil,
		size: 0,
	}
}

// NewDoublyLinkedListWithValues creates a new doubly linked list with initial values.
func NewDoublyLinkedListWithValues[T comparable](values ...T) *DoublyLinkedList[T] {
	dll := NewDoublyLinkedList[T]()
	for _, value := range values {
		dll.PushBack(value)
	}

	return dll
}

// DoublyFromSlice creates a doubly linked list from a slice.
func DoublyFromSlice[T comparable](slice []T) *DoublyLinkedList[T] {
	return NewDoublyLinkedListWithValues(slice...)
}

// newDoublyLinkedNode creates a new doubly node with the given value (private helper).
func newDoublyLinkedNode[T comparable](value T) *doublyLinkedNode[T] {
	return &doublyLinkedNode[T]{
		value: value,
		next:  nil,
		prev:  nil,
	}
}

// PushFront adds an element to the beginning of the list - O(1).
func (dll *DoublyLinkedList[T]) PushFront(value T) {
	newNode := newDoublyLinkedNode(value)

	if dll.head == nil {
		// Empty list case
		dll.head = newNode
		dll.tail = newNode
	} else {
		newNode.next = dll.head
		dll.head.prev = newNode
		dll.head = newNode
	}

	dll.size++
}

// PushBack adds an element to the end of the list - O(1).
func (dll *DoublyLinkedList[T]) PushBack(value T) {
	newNode := newDoublyLinkedNode(value)

	if dll.tail == nil {
		// Empty list case
		dll.head = newNode
		dll.tail = newNode
	} else {
		newNode.prev = dll.tail
		dll.tail.next = newNode
		dll.tail = newNode
	}

	dll.size++
}

// PopFront removes and returns the first element - O(1).
func (dll *DoublyLinkedList[T]) PopFront() (T, bool) {
	if dll.head == nil {
		var zero T

		return zero, false
	}

	value := dll.head.value

	if dll.head == dll.tail {
		// Single element case
		dll.head = nil
		dll.tail = nil
	} else {
		dll.head = dll.head.next
		dll.head.prev = nil
	}

	dll.size--

	return value, true
}

// PopBack removes and returns the last element - O(1).
// This is a key advantage of doubly linked lists over singly linked lists.
func (dll *DoublyLinkedList[T]) PopBack() (T, bool) {
	if dll.tail == nil {
		var zero T

		return zero, false
	}

	value := dll.tail.value

	if dll.head == dll.tail {
		// Single element case
		dll.head = nil
		dll.tail = nil
	} else {
		dll.tail = dll.tail.prev
		dll.tail.next = nil
	}

	dll.size--

	return value, true
}

// Front returns the first element without removing it - O(1).
func (dll *DoublyLinkedList[T]) Front() (T, bool) {
	if dll.head == nil {
		var zero T

		return zero, false
	}

	return dll.head.value, true
}

// Back returns the last element without removing it - O(1).
func (dll *DoublyLinkedList[T]) Back() (T, bool) {
	if dll.tail == nil {
		var zero T

		return zero, false
	}

	return dll.tail.value, true
}

// At returns the element at the specified index - O(n) with bidirectional optimization.
// Chooses the more efficient direction based on index position.
func (dll *DoublyLinkedList[T]) At(index int) (T, bool) {
	var zero T

	if index < 0 || index >= dll.size {
		return zero, false
	}

	// Optimize by choosing direction based on index position
	if index < dll.size/2 {
		// Traverse from head
		current := dll.head
		for range index {
			current = current.next
		}

		return current.value, true
	}
	// Traverse from tail
	current := dll.tail
	for i := dll.size - 1; i > index; i-- {
		current = current.prev
	}

	return current.value, true
}

// Insert adds an element at the specified index - O(n).
func (dll *DoublyLinkedList[T]) Insert(index int, value T) error {
	if index < 0 || index > dll.size {
		return errors.New("index out of bounds")
	}

	// Insert at beginning
	if index == 0 {
		dll.PushFront(value)

		return nil
	}

	// Insert at end
	if index == dll.size {
		dll.PushBack(value)

		return nil
	}

	// Insert in middle - optimize direction
	newNode := newDoublyLinkedNode(value)

	if index < dll.size/2 {
		// Traverse from head
		current := dll.head
		for i := 0; i < index; i++ {
			current = current.next
		}

		newNode.next = current
		newNode.prev = current.prev
		current.prev.next = newNode
		current.prev = newNode
	} else {
		// Traverse from tail
		current := dll.tail
		for i := dll.size - 1; i > index; i-- {
			current = current.prev
		}

		newNode.next = current
		newNode.prev = current.prev
		current.prev.next = newNode
		current.prev = newNode
	}

	dll.size++

	return nil
}

// Remove removes the element at the specified index - O(n).
func (dll *DoublyLinkedList[T]) Remove(index int) (T, error) {
	var zero T

	if index < 0 || index >= dll.size {
		return zero, errors.New("index out of bounds")
	}

	// Remove from beginning
	if index == 0 {
		value, _ := dll.PopFront()

		return value, nil
	}

	// Remove from end
	if index == dll.size-1 {
		value, _ := dll.PopBack()

		return value, nil
	}

	// Remove from middle - optimize direction
	var nodeToRemove *doublyLinkedNode[T]

	if index < dll.size/2 {
		// Traverse from head
		nodeToRemove = dll.head
		for i := 0; i < index; i++ {
			nodeToRemove = nodeToRemove.next
		}
	} else {
		// Traverse from tail
		nodeToRemove = dll.tail
		for i := dll.size - 1; i > index; i-- {
			nodeToRemove = nodeToRemove.prev
		}
	}

	value := nodeToRemove.value
	nodeToRemove.prev.next = nodeToRemove.next
	nodeToRemove.next.prev = nodeToRemove.prev
	dll.size--

	return value, nil
}

// Find returns the index of the first occurrence of the value, or -1 if not found - O(n).
func (dll *DoublyLinkedList[T]) Find(value T) int {
	current := dll.head
	for i := 0; current != nil; i++ {
		if current.value == value {
			return i
		}
		current = current.next
	}

	return -1
}

// Contains checks if the list contains the specified value - O(n).
func (dll *DoublyLinkedList[T]) Contains(value T) bool {
	return dll.Find(value) != -1
}

// Len returns the number of elements in the list - O(1).
func (dll *DoublyLinkedList[T]) Len() int {
	return dll.size
}

// IsEmpty returns true if the list is empty - O(1).
func (dll *DoublyLinkedList[T]) IsEmpty() bool {
	return dll.size == 0
}

// Clear removes all elements from the list - O(1).
func (dll *DoublyLinkedList[T]) Clear() {
	dll.head = nil
	dll.tail = nil
	dll.size = 0
}

// ToSlice converts the doubly linked list to a slice - O(n).
func (dll *DoublyLinkedList[T]) ToSlice() []T {
	result := make([]T, dll.size)
	current := dll.head

	for i := 0; current != nil; i++ {
		result[i] = current.value
		current = current.next
	}

	return result
}

// ToSliceReverse converts the doubly linked list to a slice in reverse order - O(n).
// This is efficient as it traverses backward from tail without needing to find the tail first.
func (dll *DoublyLinkedList[T]) ToSliceReverse() []T {
	result := make([]T, dll.size)
	current := dll.tail

	for i := 0; current != nil; i++ {
		result[i] = current.value
		current = current.prev
	}

	return result
}

// String returns a string representation of the list.
func (dll *DoublyLinkedList[T]) String() string {
	return fmt.Sprintf("DoublyLinkedList%v", dll.ToSlice())
}

// Iterator returns a forward iterator for the list.
func (dll *DoublyLinkedList[T]) Iterator() Iterator[T] {
	return &doublyLinkedListIterator[T]{
		current: dll.head,
	}
}

// ReverseIterator returns a reverse iterator for the list.
// This is a unique advantage of doubly linked lists.
func (dll *DoublyLinkedList[T]) ReverseIterator() ReverseIterator[T] {
	return &doublyLinkedListReverseIterator[T]{
		current: dll.tail,
	}
}

// doublyLinkedListIterator implements Iterator for DoublyLinkedList.
type doublyLinkedListIterator[T comparable] struct {
	current *doublyLinkedNode[T]
}

// HasNext returns true if there are more elements to iterate over.
func (it *doublyLinkedListIterator[T]) HasNext() bool {
	return it.current != nil
}

// Next returns the next element and advances the iterator.
func (it *doublyLinkedListIterator[T]) Next() (T, bool) {
	if it.current == nil {
		var zero T

		return zero, false
	}

	value := it.current.value
	it.current = it.current.next

	return value, true
}

// Value returns the current element without advancing the iterator.
func (it *doublyLinkedListIterator[T]) Value() T {
	if it.current == nil {
		var zero T

		return zero
	}

	return it.current.value
}

// doublyLinkedListReverseIterator implements ReverseIterator for DoublyLinkedList.
type doublyLinkedListReverseIterator[T comparable] struct {
	current *doublyLinkedNode[T]
}

// HasPrev returns true if there are more elements to iterate over in reverse.
func (it *doublyLinkedListReverseIterator[T]) HasPrev() bool {
	return it.current != nil
}

// Prev returns the previous element and moves the iterator backward.
func (it *doublyLinkedListReverseIterator[T]) Prev() (T, bool) {
	if it.current == nil {
		var zero T

		return zero, false
	}

	value := it.current.value
	it.current = it.current.prev

	return value, true
}

// Value returns the current element without moving the iterator.
func (it *doublyLinkedListReverseIterator[T]) Value() T {
	if it.current == nil {
		var zero T

		return zero
	}

	return it.current.value
}

// Advanced node reference operations - unique advantages of doubly linked lists

// insertAfter inserts a new element after the given node - O(1).
// This is only possible with node references and doubly linked structure.
func (dll *DoublyLinkedList[T]) insertAfter(node *doublyLinkedNode[T], value T) *doublyLinkedNode[T] {
	if node == nil {
		return nil
	}

	newNode := newDoublyLinkedNode(value)
	newNode.next = node.next
	newNode.prev = node

	if node.next != nil {
		node.next.prev = newNode
	} else {
		// Inserting after tail
		dll.tail = newNode
	}

	node.next = newNode
	dll.size++

	return newNode
}

// insertBefore inserts a new element before the given node - O(1).
// This is only possible with node references and doubly linked structure.
func (dll *DoublyLinkedList[T]) insertBefore(node *doublyLinkedNode[T], value T) *doublyLinkedNode[T] {
	if node == nil {
		return nil
	}

	newNode := newDoublyLinkedNode(value)
	newNode.next = node
	newNode.prev = node.prev

	if node.prev != nil {
		node.prev.next = newNode
	} else {
		// Inserting before head
		dll.head = newNode
	}

	node.prev = newNode
	dll.size++

	return newNode
}

// removeNode removes the specified node from the list - O(1).
// This is only possible with node references and doubly linked structure.
func (dll *DoublyLinkedList[T]) removeNode(node *doublyLinkedNode[T]) T {
	if node == nil {
		var zero T

		return zero
	}

	value := node.value

	if node.prev != nil {
		node.prev.next = node.next
	} else {
		// Removing head
		dll.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		// Removing tail
		dll.tail = node.prev
	}

	dll.size--

	return value
}
