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

package collection

// LinkedList represents the nodes used in a linked list.
type LinkedList struct {
	Val  int
	Next *LinkedList
}

// NewLinkedList creates a new LinkedList with the given value
// and a nil Next pointer.
func NewLinkedList(val int) *LinkedList {
	return &LinkedList{
		Val:  val,
		Next: nil,
	}
}

// BuildLinkedList creates a linked list from a slice of integers.
func BuildLinkedList(values []int) *LinkedList {
	if len(values) == 0 {
		return nil
	}

	head := NewLinkedList(values[0])
	current := head

	for i := 1; i < len(values); i++ {
		current.Next = NewLinkedList(values[i])
		current = current.Next
	}

	return head
}

// ToSlice converts a linked list to a slice of integers
func (l *LinkedList) ToSlice() []int {
	var result []int
	current := l

	for current != nil {
		result = append(result, current.Val)
		current = current.Next
	}

	return result
}

// DoublyLinkedList represents the nodes used in a doubly linked list.
type DoublyLinkedList struct {
	Val  int
	Next *DoublyLinkedList
	Prev *DoublyLinkedList
}

// NewDoublyLinkedList creates a new DoublyLinkedList with the given value
// and nil Next and Prev pointers.
func NewDoublyLinkedList(val int) *DoublyLinkedList {
	return &DoublyLinkedList{
		Val:  val,
		Next: nil,
		Prev: nil,
	}
}

// BuildDoublyLinkedList creates a doubly linked list from a slice of integers.
func BuildDoublyLinkedList(values []int) *DoublyLinkedList {
	if len(values) == 0 {
		return nil
	}

	head := NewDoublyLinkedList(values[0])
	current := head

	for i := 1; i < len(values); i++ {
		newNode := &DoublyLinkedList{
			Val:  values[i],
			Next: nil,
			Prev: current,
		}
		current.Next = newNode
		current = newNode
	}

	return head
}

// ToSlice converts a doubly linked list to a slice of integers.
func (d *DoublyLinkedList) ToSlice() []int {
	var result []int
	current := d

	for current != nil {
		result = append(result, current.Val)
		current = current.Next
	}

	return result
}

// ToSliceReverse converts a doubly linked list to a slice of integers
// in reverse order.
func (d *DoublyLinkedList) ToSliceReverse() []int {
	if d == nil {
		return []int{}
	}

	// Find the tail
	current := d
	for current.Next != nil {
		current = current.Next
	}

	// Traverse backward
	var result []int
	for current != nil {
		result = append(result, current.Val)
		current = current.Prev
	}

	return result
}
