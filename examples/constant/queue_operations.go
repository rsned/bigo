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

// Queue represents a simple queue data structure.
// This implementation uses a slice with a front index to avoid
// expensive slice operations during dequeue.
type Queue struct {
	items []int // Slice to store queue elements
	front int   // Index of the front element (for O(1) dequeue)
}

// NewQueue creates a new queue.
// Initializes an empty queue with zero capacity and front index at 0.
func NewQueue() *Queue {
	return &Queue{
		items: make([]int, 0), // Start with empty slice
		front: 0,              // Front index starts at 0
	}
}

// CopyQueue creates a new queue with the same items.
// This utility function helps create test queues with predefined values.
func CopyQueue(vals []int) *Queue {
	return &Queue{
		items: append([]int{}, vals...), // Copy slice to avoid sharing
		front: 0,                        // Start with front at beginning
	}
}

// Enqueue performs O(1) enqueue operation.
// Adds an item to the back of the queue. This is constant time
// because we simply append to the end of the slice.
func (q *Queue) Enqueue(item int) {
	// Append to the end of the slice - O(1) amortized time
	q.items = append(q.items, item)
}

// Dequeue performs O(1) dequeue operation.
// Removes and returns the front item from the queue. This is constant time
// because we use a front index instead of shifting all elements.
func (q *Queue) Dequeue() (int, bool) {
	// Check if queue is empty (front index has reached the end)
	if q.front >= len(q.items) {
		return 0, false
	}

	// Get the front item without removing it from the slice
	item := q.items[q.front]
	// Move the front index forward - this is the O(1) trick
	q.front++

	return item, true
}
