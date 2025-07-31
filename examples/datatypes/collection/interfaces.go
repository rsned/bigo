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

// Iterator provides forward iteration over a collection.
type Iterator[T comparable] interface {
	// HasNext returns true if there are more elements to iterate over
	HasNext() bool
	// Next returns the next element and advances the iterator
	Next() (T, bool)
	// Value returns the current element without advancing the iterator
	Value() T
}

// ReverseIterator provides reverse iteration over a collection.
// This is primarily used by DoublyLinkedList for efficient reverse traversal.
type ReverseIterator[T comparable] interface {
	// HasPrev returns true if there are more elements to iterate over in reverse
	HasPrev() bool
	// Prev returns the previous element and moves the iterator backward
	Prev() (T, bool)
	// Value returns the current element without moving the iterator
	Value() T
}
