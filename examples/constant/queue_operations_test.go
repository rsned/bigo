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

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestQueueEnqueue(t *testing.T) {
	tests := []struct {
		name        string
		initial     []int
		enqueueItem int
		want        []int
	}{
		{"enqueue to empty queue", []int{}, 5, []int{5}},
		{"enqueue to existing queue", []int{1, 2}, 3, []int{1, 2, 3}},
		{"enqueue negative value", []int{1}, -5, []int{1, -5}},
		{"enqueue zero", []int{10}, 0, []int{10, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := &Queue{items: make([]int, len(tt.initial)), front: 0}
			copy(queue.items, tt.initial)

			queue.Enqueue(tt.enqueueItem)

			if !cmp.Equal(queue.items, tt.want) {
				t.Errorf("Enqueue(%d) resulted in %v, want %v", tt.enqueueItem, queue.items, tt.want)
			}
		})
	}
}

func TestQueueDequeue(t *testing.T) {
	tests := []struct {
		name        string
		initial     []int
		front       int
		wantValue   int
		wantSuccess bool
		wantFront   int
	}{
		{"dequeue from single item", []int{5}, 0, 5, true, 1},
		{"dequeue from multiple items", []int{1, 2, 3}, 0, 1, true, 1},
		{"dequeue when front at end", []int{1, 2, 3}, 3, 0, false, 3},
		{"dequeue after partial consumption", []int{1, 2, 3}, 1, 2, true, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := CopyQueue(tt.initial)
			queue.front = tt.front

			gotValue, gotSuccess := queue.Dequeue()

			if gotValue != tt.wantValue || gotSuccess != tt.wantSuccess {
				t.Errorf("Dequeue() = (%d, %t), want (%d, %t)",
					gotValue, gotSuccess, tt.wantValue, tt.wantSuccess)
			}

			if queue.front != tt.wantFront {
				t.Errorf("After Dequeue(), front = %d, want %d", queue.front, tt.wantFront)
			}
		})
	}
}

func TestQueueOperations_Integration(t *testing.T) {
	queue := NewQueue()

	// Test enqueue operations
	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)

	if len(queue.items) != 3 {
		t.Errorf("After 3 enqueues, queue size = %d, want 3", len(queue.items))
	}

	// Test dequeue operations (FIFO)
	val, ok := queue.Dequeue()
	if !ok || val != 1 {
		t.Errorf("First dequeue = (%d, %t), want (1, true)", val, ok)
	}

	val, ok = queue.Dequeue()
	if !ok || val != 2 {
		t.Errorf("Second dequeue = (%d, %t), want (2, true)", val, ok)
	}

	val, ok = queue.Dequeue()
	if !ok || val != 3 {
		t.Errorf("Third dequeue = (%d, %t), want (3, true)", val, ok)
	}

	val, ok = queue.Dequeue()
	if ok || val != 0 {
		t.Errorf("Dequeue from empty queue = (%d, %t), want (0, false)", val, ok)
	}
}

// Benchmark functions for queue operations

func BenchmarkQueue_Enqueue(b *testing.B) {
	queue := &Queue{
		items: nil,
		front: 0,
	}
	b.ResetTimer()
	for b.Loop() {
		queue.Enqueue(42)
	}
}

func BenchmarkQueue_Dequeue(b *testing.B) {
	queue := &Queue{
		items: nil,
		front: 0,
	}
	queue.Enqueue(42)
	queue.Enqueue(43)
	queue.Enqueue(44)

	b.ResetTimer()
	for b.Loop() {
		if queue.front >= len(queue.items) {
			queue.front = 0
			queue.items = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100}
		}
		_, _ = queue.Dequeue()
	}
}
