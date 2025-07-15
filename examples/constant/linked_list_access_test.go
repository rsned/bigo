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

func TestLinkedListGetFirstElement(t *testing.T) {
	tests := []struct {
		name        string
		head        *ListNode
		tail        *ListNode
		wantValue   int
		wantSuccess bool
	}{
		{"empty list", nil, nil, 0, false},
		{"single element", NewListNode(5), NewListNode(5), 5, true},
		{"multiple elements", BuildLinkedList([]int{1, 2, 3}), nil, 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := &LinkedList{Head: tt.head, Tail: tt.tail}
			if tt.head != nil && tt.tail == nil {
				// For multi-element lists, set tail to last element
				current := tt.head
				for current.Next != nil {
					current = current.Next
				}

				ll.Tail = current
			}

			gotValue, gotSuccess := ll.GetFirstElement()
			if gotValue != tt.wantValue || gotSuccess != tt.wantSuccess {
				t.Errorf("GetFirstElement() = (%d, %t), want (%d, %t)",
					gotValue, gotSuccess, tt.wantValue, tt.wantSuccess)
			}
		})
	}
}

func TestLinkedListGetLastElement(t *testing.T) {
	tests := []struct {
		name        string
		head        *ListNode
		tail        *ListNode
		wantValue   int
		wantSuccess bool
	}{
		{"empty list", nil, nil, 0, false},
		{"single element", NewListNode(5), NewListNode(5), 5, true},
		{"multiple elements with proper tail", BuildLinkedList([]int{1, 2, 3}), nil, 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := &LinkedList{Head: tt.head, Tail: tt.tail}
			if tt.head != nil && tt.tail == nil {
				// For multi-element lists, set tail to last element
				current := tt.head
				for current.Next != nil {
					current = current.Next
				}

				ll.Tail = current
			}

			gotValue, gotSuccess := ll.GetLastElement()
			if gotValue != tt.wantValue || gotSuccess != tt.wantSuccess {
				t.Errorf("GetLastElement() = (%d, %t), want (%d, %t)",
					gotValue, gotSuccess, tt.wantValue, tt.wantSuccess)
			}
		})
	}
}

func TestLinkedListIntegration(t *testing.T) {
	// Test with properly constructed linked list
	head := BuildLinkedList([]int{10, 20, 30})

	// Find tail
	current := head
	for current.Next != nil {
		current = current.Next
	}

	tail := current

	ll := &LinkedList{Head: head, Tail: tail}

	// Test first element
	first, ok := ll.GetFirstElement()
	if !ok || first != 10 {
		t.Errorf("GetFirstElement() = (%d, %t), want (10, true)", first, ok)
	}

	// Test last element
	last, ok := ll.GetLastElement()
	if !ok || last != 30 {
		t.Errorf("GetLastElement() = (%d, %t), want (30, true)", last, ok)
	}
}

func TestLinkedListNewListNode(t *testing.T) {
	tests := []struct {
		name string
		val  int
		want *ListNode
	}{
		{"positive value", 5, &ListNode{Val: 5, Next: nil}},
		{"zero value", 0, &ListNode{Val: 0, Next: nil}},
		{"negative value", -3, &ListNode{Val: -3, Next: nil}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewListNode(tt.val)
			if got.Val != tt.want.Val || got.Next != nil {
				t.Errorf("NewListNode(%d) = %v, want %v", tt.val, got, tt.want)
			}
		})
	}
}

func TestLinkedListBuildLinkedList(t *testing.T) {
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
				t.Errorf("BuildLinkedList(%v) = %v, want %v, diff: %v", tt.values, result, tt.want, cmp.Diff(result, tt.want))
			}
		})
	}
}

func TestLinkedListToSlice(t *testing.T) {
	tests := []struct {
		name string
		list *ListNode
		want []int
	}{
		{"nil list", nil, []int{}},
		{"single node", NewListNode(5), []int{5}},
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
				t.Errorf("ListNode.ToSlice() = %v, want %v, diff: %v", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

// Benchmark functions for ListNode operations

func BenchmarkNewListNode(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewListNode(42)
	}
}

func BenchmarkBuildLinkedList(b *testing.B) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for b.Loop() {
		_ = BuildLinkedList(values)
	}
}

func BenchmarkListNode_ToSlice(b *testing.B) {
	list := BuildLinkedList([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	b.ResetTimer()
	for b.Loop() {
		_ = list.ToSlice()
	}
}

// Benchmark functions for linked list operations

func BenchmarkLinkedList_GetFirstElement(b *testing.B) {
	head := BuildLinkedList([]int{1, 2, 3, 4, 5})
	current := head
	for current.Next != nil {
		current = current.Next
	}

	tail := current
	ll := &LinkedList{Head: head, Tail: tail}

	b.ResetTimer()
	for b.Loop() {
		_, _ = ll.GetFirstElement()
	}
}

func BenchmarkLinkedList_GetLastElement(b *testing.B) {
	head := BuildLinkedList([]int{1, 2, 3, 4, 5})
	current := head
	for current.Next != nil {
		current = current.Next
	}

	tail := current
	ll := &LinkedList{Head: head, Tail: tail}

	b.ResetTimer()
	for b.Loop() {
		_, _ = ll.GetLastElement()
	}
}
