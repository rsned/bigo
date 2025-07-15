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

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewLinkedList(t *testing.T) {
	tests := []struct {
		name string
		val  int
		want *LinkedList
	}{
		{"positive value", 5, &LinkedList{Val: 5, Next: nil}},
		{"zero value", 0, &LinkedList{Val: 0, Next: nil}},
		{"negative value", -3, &LinkedList{Val: -3, Next: nil}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLinkedList(tt.val)
			if got.Val != tt.want.Val || got.Next != nil {
				t.Errorf("NewListNode(%d) = %v, want %v", tt.val, got, tt.want)
			}
		})
	}
}

func TestBuildLinkedList(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{"empty slice", []int{}, nil},
		{"single element", []int{1}, []int{1}},
		{"multiple elements", []int{1, 2, 3, 4}, []int{1, 2, 3, 4}},
		{"negative elements", []int{-1, -2, -3}, []int{-1, -2, -3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildLinkedList(tt.values)
			if got == nil && tt.want == nil {
				return
			}

			if got == nil && tt.want != nil {
				t.Errorf("BuildLinkedList(%v) = nil, want %v", tt.values, tt.want)

				return
			}

			result := got.ToSlice()
			if !cmp.Equal(result, tt.want) {
				t.Errorf("BuildLinkedList(%v) = %v, want %v", tt.values, result, tt.want)
			}
		})
	}
}

func TestListNode_ToSlice(t *testing.T) {
	tests := []struct {
		name string
		list *LinkedList
		want []int
	}{
		{"nil list", nil, []int{}},
		{"single node", NewLinkedList(5), []int{5}},
		{"multiple nodes", BuildLinkedList([]int{1, 2, 3}), []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []int
			if tt.list != nil {
				got = tt.list.ToSlice()
			} else {
				got = []int{}
			}

			if !cmp.Equal(got, tt.want) {
				t.Errorf("ListNode.ToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Benchmark functions for ListNode operations

func BenchmarkNewListNode(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewLinkedList(42)
	}
}

func BenchmarkBuildLinkedList(b *testing.B) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for b.Loop() {
		_ = BuildLinkedList(values)
	}
}

func BenchmarkLinkedListToSlice(b *testing.B) {
	list := BuildLinkedList([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	b.ResetTimer()
	for b.Loop() {
		_ = list.ToSlice()
	}
}

func TestNewDoublyLinkedList(t *testing.T) {
	tests := []struct {
		name string
		val  int
		want *DoublyLinkedList
	}{
		{"positive value", 5, &DoublyLinkedList{Val: 5, Next: nil, Prev: nil}},
		{"zero value", 0, &DoublyLinkedList{Val: 0, Next: nil, Prev: nil}},
		{"negative value", -3, &DoublyLinkedList{Val: -3, Next: nil, Prev: nil}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDoublyLinkedList(tt.val)
			if got.Val != tt.want.Val || got.Next != nil || got.Prev != nil {
				t.Errorf("NewDoublyLinkedList(%d) = %v, want %v", tt.val, got, tt.want)
			}
		})
	}
}

func TestBuildDoublyLinkedList(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{"empty slice", []int{}, nil},
		{"single element", []int{1}, []int{1}},
		{"multiple elements", []int{1, 2, 3, 4}, []int{1, 2, 3, 4}},
		{"negative elements", []int{-1, -2, -3}, []int{-1, -2, -3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildDoublyLinkedList(tt.values)
			if got == nil && tt.want == nil {
				return
			}

			if got == nil && tt.want != nil {
				t.Errorf("BuildDoublyLinkedList(%v) = nil, want %v", tt.values, tt.want)

				return
			}

			result := got.ToSlice()
			if !cmp.Equal(result, tt.want) {
				t.Errorf("BuildDoublyLinkedList(%v) = %v, want %v", tt.values, result, tt.want)
			}

			// Test backward links
			if len(tt.values) > 1 {
				current := got
				for current.Next != nil {
					current = current.Next
				}
				// Verify we can traverse backward
				reverseResult := got.ToSliceReverse()
				expectedReverse := make([]int, len(tt.want))
				for i, v := range tt.want {
					expectedReverse[len(tt.want)-1-i] = v
				}
				if !cmp.Equal(reverseResult, expectedReverse) {
					t.Errorf("BuildDoublyLinkedList(%v) reverse traversal = %v, want %v", tt.values, reverseResult, expectedReverse)
				}
			}
		})
	}
}

func TestDoublyLinkedList_ToSlice(t *testing.T) {
	tests := []struct {
		name string
		list *DoublyLinkedList
		want []int
	}{
		{"nil list", nil, []int{}},
		{"single node", NewDoublyLinkedList(5), []int{5}},
		{"multiple nodes", BuildDoublyLinkedList([]int{1, 2, 3}), []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []int
			if tt.list != nil {
				got = tt.list.ToSlice()
			} else {
				got = []int{}
			}

			if !cmp.Equal(got, tt.want) {
				t.Errorf("DoublyLinkedList.ToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoublyLinkedList_ToSliceReverse(t *testing.T) {
	tests := []struct {
		name string
		list *DoublyLinkedList
		want []int
	}{
		{"nil list", nil, []int{}},
		{"single node", NewDoublyLinkedList(5), []int{5}},
		{"multiple nodes", BuildDoublyLinkedList([]int{1, 2, 3}), []int{3, 2, 1}},
		{"five nodes", BuildDoublyLinkedList([]int{1, 2, 3, 4, 5}), []int{5, 4, 3, 2, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []int
			if tt.list != nil {
				got = tt.list.ToSliceReverse()
			} else {
				got = []int{}
			}

			if !cmp.Equal(got, tt.want) {
				t.Errorf("DoublyLinkedList.ToSliceReverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkNewDoublyLinkedList(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewDoublyLinkedList(42)
	}
}

func BenchmarkBuildDoublyLinkedList(b *testing.B) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for b.Loop() {
		_ = BuildDoublyLinkedList(values)
	}
}

func BenchmarkDoublyLinkedListToSlice(b *testing.B) {
	list := BuildDoublyLinkedList([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	b.ResetTimer()
	for b.Loop() {
		_ = list.ToSlice()
	}
}

func BenchmarkDoublyLinkedListToSliceReverse(b *testing.B) {
	list := BuildDoublyLinkedList([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	b.ResetTimer()
	for b.Loop() {
		_ = list.ToSliceReverse()
	}
}
