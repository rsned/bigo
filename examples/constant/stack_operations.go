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

// DynamicStack represents a simple stack data structure.
// This implementation uses a slice that grows dynamically as needed,
// providing O(1) amortized time for push/pop operations.
type DynamicStack struct {
	items []int // Slice to store stack elements (top is at the end)
}

// NewDynamicStack creates a new stack.
// Initializes an empty stack with zero capacity that will grow as needed.
func NewDynamicStack() *DynamicStack {
	return &DynamicStack{
		items: make([]int, 0), // Start with empty slice
	}
}

// CopyDynamicStack creates a new stack with the same items.
// This utility function helps create test stacks with predefined values.
func CopyDynamicStack(vals []int) *DynamicStack {
	return &DynamicStack{
		items: append([]int{}, vals...), // Copy slice to avoid sharing
	}
}

// Push performs O(1) stack push operation.
// Adds an item to the top of the stack. This is constant time (amortized)
// because we append to the end of the slice.
func (s *DynamicStack) Push(item int) {
	// Append to the end of the slice - O(1) amortized time
	// The slice may occasionally need to be resized, but this is amortized O(1)
	s.items = append(s.items, item)
}

// Pop performs O(1) stack pop operation.
// Removes and returns the top item from the stack. This is constant time
// because we simply remove the last element from the slice.
func (s *DynamicStack) Pop() (int, bool) {
	// Check if stack is empty
	if len(s.items) == 0 {
		return 0, false
	}

	// Get the top item (last element in the slice)
	item := s.items[len(s.items)-1]
	// Remove the last element by reslicing - O(1) operation
	s.items = s.items[:len(s.items)-1]

	return item, true
}
