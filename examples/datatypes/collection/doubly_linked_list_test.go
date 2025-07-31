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
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewDoublyLinkedList(t *testing.T) {
	t.Run("int list", func(t *testing.T) {
		dll := NewDoublyLinkedList[int]()
		if dll.Len() != 0 || !dll.IsEmpty() {
			t.Errorf("NewDoublyLinkedList[int]() should create empty list, got len=%d, empty=%v", dll.Len(), dll.IsEmpty())
		}
	})

	t.Run("string list", func(t *testing.T) {
		dll := NewDoublyLinkedList[string]()
		if dll.Len() != 0 || !dll.IsEmpty() {
			t.Errorf("NewDoublyLinkedList[string]() should create empty list, got len=%d, empty=%v", dll.Len(), dll.IsEmpty())
		}
	})
}

func TestNewDoublyLinkedListWithValues(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{"empty", []int{}, []int{}},
		{"single element", []int{42}, []int{42}},
		{"multiple elements", []int{1, 2, 3, 4}, []int{1, 2, 3, 4}},
		{"negative elements", []int{-1, -2, -3}, []int{-1, -2, -3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dll := NewDoublyLinkedListWithValues(tt.values...)
			got := dll.ToSlice()
			if !cmp.Equal(got, tt.want) {
				t.Errorf("NewDoublyLinkedListWithValues(%v) = %v, want %v", tt.values, got, tt.want)
			}
			if dll.Len() != len(tt.want) {
				t.Errorf("NewDoublyLinkedListWithValues(%v) length = %d, want %d", tt.values, dll.Len(), len(tt.want))
			}
		})
	}
}

func TestDoublyFromSlice(t *testing.T) {
	tests := []struct {
		name  string
		slice []string
		want  []string
	}{
		{"empty", []string{}, []string{}},
		{"single element", []string{"hello"}, []string{"hello"}},
		{"multiple elements", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dll := DoublyFromSlice(tt.slice)
			got := dll.ToSlice()
			if !cmp.Equal(got, tt.want) {
				t.Errorf("DoublyFromSlice(%v) = %v, want %v", tt.slice, got, tt.want)
			}
		})
	}
}

func TestDoublyLinkedListPushFront(t *testing.T) {
	dll := NewDoublyLinkedList[int]()

	// Test pushing to empty list
	dll.PushFront(1)
	if got := dll.ToSlice(); !cmp.Equal(got, []int{1}) {
		t.Errorf("After PushFront(1) on empty list: got %v, want [1]", got)
	}

	// Test pushing to non-empty list
	dll.PushFront(2)
	if got := dll.ToSlice(); !cmp.Equal(got, []int{2, 1}) {
		t.Errorf("After PushFront(2): got %v, want [2, 1]", got)
	}

	// Test length tracking
	if dll.Len() != 2 {
		t.Errorf("Length should be 2, got %d", dll.Len())
	}

	// Test bidirectional links
	if got := dll.ToSliceReverse(); !cmp.Equal(got, []int{1, 2}) {
		t.Errorf("Reverse traversal: got %v, want [1, 2]", got)
	}
}

func TestDoublyLinkedListPushBack(t *testing.T) {
	dll := NewDoublyLinkedList[int]()

	// Test pushing to empty list
	dll.PushBack(1)
	if got := dll.ToSlice(); !cmp.Equal(got, []int{1}) {
		t.Errorf("After PushBack(1) on empty list: got %v, want [1]", got)
	}

	// Test pushing to non-empty list
	dll.PushBack(2)
	if got := dll.ToSlice(); !cmp.Equal(got, []int{1, 2}) {
		t.Errorf("After PushBack(2): got %v, want [1, 2]", got)
	}

	// Test length tracking
	if dll.Len() != 2 {
		t.Errorf("Length should be 2, got %d", dll.Len())
	}

	// Test bidirectional links
	if got := dll.ToSliceReverse(); !cmp.Equal(got, []int{2, 1}) {
		t.Errorf("Reverse traversal: got %v, want [2, 1]", got)
	}
}

func TestDoublyLinkedListPopFront(t *testing.T) {
	dll := DoublyFromSlice([]int{1, 2, 3})

	// Pop from non-empty list
	val, ok := dll.PopFront()
	if !ok || val != 1 {
		t.Errorf("PopFront() = (%d, %v), want (1, true)", val, ok)
	}
	if got := dll.ToSlice(); !cmp.Equal(got, []int{2, 3}) {
		t.Errorf("After PopFront(): got %v, want [2, 3]", got)
	}

	// Test bidirectional links maintained
	if got := dll.ToSliceReverse(); !cmp.Equal(got, []int{3, 2}) {
		t.Errorf("Reverse after PopFront(): got %v, want [3, 2]", got)
	}

	// Pop until empty
	dll.PopFront()
	dll.PopFront()

	// Pop from empty list
	val, ok = dll.PopFront()
	if ok || val != 0 {
		t.Errorf("PopFront() on empty list = (%d, %v), want (0, false)", val, ok)
	}
}

func TestDoublyLinkedListPopBack(t *testing.T) {
	dll := DoublyFromSlice([]int{1, 2, 3})

	// Pop from non-empty list - this is O(1) for doubly linked list!
	val, ok := dll.PopBack()
	if !ok || val != 3 {
		t.Errorf("PopBack() = (%d, %v), want (3, true)", val, ok)
	}
	if got := dll.ToSlice(); !cmp.Equal(got, []int{1, 2}) {
		t.Errorf("After PopBack(): got %v, want [1, 2]", got)
	}

	// Test bidirectional links maintained
	if got := dll.ToSliceReverse(); !cmp.Equal(got, []int{2, 1}) {
		t.Errorf("Reverse after PopBack(): got %v, want [2, 1]", got)
	}

	// Pop until empty
	dll.PopBack()
	dll.PopBack()

	// Pop from empty list
	val, ok = dll.PopBack()
	if ok || val != 0 {
		t.Errorf("PopBack() on empty list = (%d, %v), want (0, false)", val, ok)
	}
}

func TestDoublyLinkedListFront(t *testing.T) {
	// Empty list
	dll := NewDoublyLinkedList[int]()
	val, ok := dll.Front()
	if ok || val != 0 {
		t.Errorf("Front() on empty list = (%d, %v), want (0, false)", val, ok)
	}

	// Non-empty list
	dll = DoublyFromSlice([]int{1, 2, 3})
	val, ok = dll.Front()
	if !ok || val != 1 {
		t.Errorf("Front() = (%d, %v), want (1, true)", val, ok)
	}
}

func TestDoublyLinkedListBack(t *testing.T) {
	// Empty list
	dll := NewDoublyLinkedList[int]()
	val, ok := dll.Back()
	if ok || val != 0 {
		t.Errorf("Back() on empty list = (%d, %v), want (0, false)", val, ok)
	}

	// Non-empty list
	dll = DoublyFromSlice([]int{1, 2, 3})
	val, ok = dll.Back()
	if !ok || val != 3 {
		t.Errorf("Back() = (%d, %v), want (3, true)", val, ok)
	}
}

func TestDoublyLinkedListAt(t *testing.T) {
	dll := DoublyFromSlice([]int{10, 20, 30, 40, 50, 60})

	tests := []struct {
		index   int
		wantVal int
		wantOk  bool
	}{
		{0, 10, true},
		{1, 20, true},
		{2, 30, true},
		{3, 40, true}, // This should use forward traversal (index < size/2)
		{4, 50, true}, // This should use backward traversal (index >= size/2)
		{5, 60, true},
		{6, 0, false},  // out of bounds
		{-1, 0, false}, // negative index
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			gotVal, gotOk := dll.At(tt.index)
			if gotVal != tt.wantVal || gotOk != tt.wantOk {
				t.Errorf("At(%d) = (%d, %v), want (%d, %v)", tt.index, gotVal, gotOk, tt.wantVal, tt.wantOk)
			}
		})
	}
}

func TestDoublyLinkedListInsert(t *testing.T) {
	t.Run("valid insertions", func(t *testing.T) {
		dll := DoublyFromSlice([]int{1, 3, 5})

		// Insert at beginning
		err := dll.Insert(0, 0)
		if err != nil {
			t.Errorf("Insert(0, 0) error = %v", err)
		}
		if got := dll.ToSlice(); !cmp.Equal(got, []int{0, 1, 3, 5}) {
			t.Errorf("After Insert(0, 0): got %v, want [0, 1, 3, 5]", got)
		}

		// Test bidirectional links after insertion
		if got := dll.ToSliceReverse(); !cmp.Equal(got, []int{5, 3, 1, 0}) {
			t.Errorf("Reverse after Insert(0, 0): got %v, want [5, 3, 1, 0]", got)
		}

		// Insert in middle
		err = dll.Insert(2, 2)
		if err != nil {
			t.Errorf("Insert(2, 2) error = %v", err)
		}
		if got := dll.ToSlice(); !cmp.Equal(got, []int{0, 1, 2, 3, 5}) {
			t.Errorf("After Insert(2, 2): got %v, want [0, 1, 2, 3, 5]", got)
		}

		// Insert at end
		err = dll.Insert(5, 6)
		if err != nil {
			t.Errorf("Insert(5, 6) error = %v", err)
		}
		if got := dll.ToSlice(); !cmp.Equal(got, []int{0, 1, 2, 3, 5, 6}) {
			t.Errorf("After Insert(5, 6): got %v, want [0, 1, 2, 3, 5, 6]", got)
		}
	})

	t.Run("invalid insertions", func(t *testing.T) {
		dll := DoublyFromSlice([]int{1, 2, 3})

		if err := dll.Insert(-1, 0); err == nil {
			t.Error("Insert(-1, 0) should return error")
		}
		if err := dll.Insert(4, 0); err == nil {
			t.Error("Insert(4, 0) should return error")
		}
	})
}

func TestDoublyLinkedListRemove(t *testing.T) {
	t.Run("valid removals", func(t *testing.T) {
		dll := DoublyFromSlice([]int{1, 2, 3, 4, 5})

		// Remove from beginning
		val, err := dll.Remove(0)
		if err != nil || val != 1 {
			t.Errorf("Remove(0) = (%d, %v), want (1, nil)", val, err)
		}
		if got := dll.ToSlice(); !cmp.Equal(got, []int{2, 3, 4, 5}) {
			t.Errorf("After Remove(0): got %v, want [2, 3, 4, 5]", got)
		}

		// Test bidirectional links maintained
		if got := dll.ToSliceReverse(); !cmp.Equal(got, []int{5, 4, 3, 2}) {
			t.Errorf("Reverse after Remove(0): got %v, want [5, 4, 3, 2]", got)
		}

		// Remove from middle
		val, err = dll.Remove(1)
		if err != nil || val != 3 {
			t.Errorf("Remove(1) = (%d, %v), want (3, nil)", val, err)
		}
		if got := dll.ToSlice(); !cmp.Equal(got, []int{2, 4, 5}) {
			t.Errorf("After Remove(1): got %v, want [2, 4, 5]", got)
		}

		// Remove from end
		val, err = dll.Remove(2)
		if err != nil || val != 5 {
			t.Errorf("Remove(2) = (%d, %v), want (5, nil)", val, err)
		}
		if got := dll.ToSlice(); !cmp.Equal(got, []int{2, 4}) {
			t.Errorf("After Remove(2): got %v, want [2, 4]", got)
		}
	})

	t.Run("invalid removals", func(t *testing.T) {
		dll := DoublyFromSlice([]int{1, 2, 3})

		if _, err := dll.Remove(-1); err == nil {
			t.Error("Remove(-1) should return error")
		}
		if _, err := dll.Remove(3); err == nil {
			t.Error("Remove(3) should return error")
		}
	})
}

func TestDoublyLinkedListFind(t *testing.T) {
	dll := DoublyFromSlice([]int{10, 20, 30, 20, 40})

	tests := []struct {
		value int
		want  int
	}{
		{10, 0},  // first element
		{20, 1},  // first occurrence
		{30, 2},  // middle element
		{40, 4},  // last element
		{50, -1}, // not found
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := dll.Find(tt.value)
			if got != tt.want {
				t.Errorf("Find(%d) = %d, want %d", tt.value, got, tt.want)
			}
		})
	}
}

func TestDoublyLinkedListContains(t *testing.T) {
	dll := DoublyFromSlice([]int{1, 2, 3, 4, 5})

	tests := []struct {
		value int
		want  bool
	}{
		{1, true},
		{3, true},
		{5, true},
		{6, false},
		{0, false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := dll.Contains(tt.value)
			if got != tt.want {
				t.Errorf("Contains(%d) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestDoublyLinkedListClear(t *testing.T) {
	dll := DoublyFromSlice([]int{1, 2, 3, 4, 5})
	dll.Clear()

	if !dll.IsEmpty() || dll.Len() != 0 {
		t.Errorf("After Clear(): IsEmpty()=%v, Len()=%d, want true, 0", dll.IsEmpty(), dll.Len())
	}

	if got := dll.ToSlice(); len(got) != 0 {
		t.Errorf("After Clear(): ToSlice()=%v, want []", got)
	}

	if got := dll.ToSliceReverse(); len(got) != 0 {
		t.Errorf("After Clear(): ToSliceReverse()=%v, want []", got)
	}
}

func TestDoublyLinkedListToSliceReverse(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{"empty", []int{}, []int{}},
		{"single element", []int{42}, []int{42}},
		{"multiple elements", []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{"two elements", []int{10, 20}, []int{20, 10}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dll := DoublyFromSlice(tt.values)
			got := dll.ToSliceReverse()
			if !cmp.Equal(got, tt.want) {
				t.Errorf("ToSliceReverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoublyLinkedListIterator(t *testing.T) {
	dll := DoublyFromSlice([]int{1, 2, 3})
	it := dll.Iterator()

	var result []int
	for it.HasNext() {
		val, ok := it.Next()
		if !ok {
			t.Error("Next() returned false when HasNext() was true")
		}
		result = append(result, val)
	}

	if !cmp.Equal(result, []int{1, 2, 3}) {
		t.Errorf("Iterator produced %v, want [1, 2, 3]", result)
	}

	// Test iteration on empty list
	empty := NewDoublyLinkedList[int]()
	emptyIt := empty.Iterator()
	if emptyIt.HasNext() {
		t.Error("Empty list iterator should not have next")
	}
}

func TestDoublyLinkedListReverseIterator(t *testing.T) {
	dll := DoublyFromSlice([]int{1, 2, 3})
	it := dll.ReverseIterator()

	var result []int
	for it.HasPrev() {
		val, ok := it.Prev()
		if !ok {
			t.Error("Prev() returned false when HasPrev() was true")
		}
		result = append(result, val)
	}

	if !cmp.Equal(result, []int{3, 2, 1}) {
		t.Errorf("ReverseIterator produced %v, want [3, 2, 1]", result)
	}

	// Test reverse iteration on empty list
	empty := NewDoublyLinkedList[int]()
	emptyIt := empty.ReverseIterator()
	if emptyIt.HasPrev() {
		t.Error("Empty list reverse iterator should not have prev")
	}
}

// Tests for performance advantages of doubly linked lists

func TestDoublyLinkedListPopBackPerformance(t *testing.T) {
	// This test documents that PopBack is O(1) for doubly linked lists
	// vs O(n) for singly linked lists
	dll := DoublyFromSlice([]int{1, 2, 3, 4, 5})

	// Multiple PopBack operations should all be O(1)
	for i := 5; i >= 1; i-- {
		val, ok := dll.PopBack()
		if !ok || val != i {
			t.Errorf("PopBack() = (%d, %v), want (%d, true)", val, ok, i)
		}
		if dll.Len() != i-1 {
			t.Errorf("After PopBack(), length = %d, want %d", dll.Len(), i-1)
		}
	}
}

func TestDoublyLinkedListBidirectionalAtOptimization(t *testing.T) {
	// Test that At() chooses the optimal direction
	dll := DoublyFromSlice([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

	// Access elements from both ends to verify optimization works
	val, ok := dll.At(1) // Should traverse from head (forward)
	if !ok || val != 1 {
		t.Errorf("At(1) = (%d, %v), want (1, true)", val, ok)
	}

	val, ok = dll.At(8) // Should traverse from tail (backward)
	if !ok || val != 8 {
		t.Errorf("At(8) = (%d, %v), want (8, true)", val, ok)
	}
}

// Benchmark functions for performance testing and comparison

func BenchmarkNewDoublyLinkedList(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewDoublyLinkedList[int]()
	}
}

func BenchmarkDoublyLinkedListPushFront(b *testing.B) {
	dll := NewDoublyLinkedList[int]()
	b.ResetTimer()
	for b.Loop() {
		dll.PushFront(42)
	}
}

func BenchmarkDoublyLinkedListPushBack(b *testing.B) {
	dll := NewDoublyLinkedList[int]()
	b.ResetTimer()
	for b.Loop() {
		dll.PushBack(42)
	}
}

func BenchmarkDoublyLinkedListPopFront(b *testing.B) {
	dll := NewDoublyLinkedList[int]()
	for i := 0; i < 1000; i++ {
		dll.PushBack(i)
	}
	b.ResetTimer()
	for b.Loop() {
		if dll.IsEmpty() {
			for i := 0; i < 1000; i++ {
				dll.PushBack(i)
			}
		}
		dll.PopFront()
	}
}

func BenchmarkDoublyLinkedListPopBack(b *testing.B) {
	dll := NewDoublyLinkedList[int]()
	for i := 0; i < 1000; i++ {
		dll.PushBack(i)
	}
	b.ResetTimer()
	for b.Loop() {
		if dll.IsEmpty() {
			for i := 0; i < 1000; i++ {
				dll.PushBack(i)
			}
		}
		dll.PopBack() // This is O(1) for doubly linked list!
	}
}

func BenchmarkDoublyLinkedListToSlice(b *testing.B) {
	dll := DoublyFromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	b.ResetTimer()
	for b.Loop() {
		_ = dll.ToSlice()
	}
}

func BenchmarkDoublyLinkedListToSliceReverse(b *testing.B) {
	dll := DoublyFromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	b.ResetTimer()
	for b.Loop() {
		_ = dll.ToSliceReverse()
	}
}

func BenchmarkDoublyLinkedListFind(b *testing.B) {
	dll := NewDoublyLinkedList[int]()
	for i := 0; i < 1000; i++ {
		dll.PushBack(i)
	}
	b.ResetTimer()
	for b.Loop() {
		_ = dll.Find(500)
	}
}

func BenchmarkDoublyLinkedListAtForward(b *testing.B) {
	dll := NewDoublyLinkedList[int]()
	for i := 0; i < 1000; i++ {
		dll.PushBack(i)
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = dll.At(100) // Should use forward traversal
	}
}

func BenchmarkDoublyLinkedListAtBackward(b *testing.B) {
	dll := NewDoublyLinkedList[int]()
	for i := 0; i < 1000; i++ {
		dll.PushBack(i)
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = dll.At(900) // Should use backward traversal
	}
}

func BenchmarkDoublyLinkedListIterator(b *testing.B) {
	dll := DoublyFromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	b.ResetTimer()
	for b.Loop() {
		it := dll.Iterator()
		for it.HasNext() {
			_, _ = it.Next()
		}
	}
}

func BenchmarkDoublyLinkedListReverseIterator(b *testing.B) {
	dll := DoublyFromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	b.ResetTimer()
	for b.Loop() {
		it := dll.ReverseIterator()
		for it.HasPrev() {
			_, _ = it.Prev()
		}
	}
}

// Comparative benchmarks to demonstrate performance differences

func BenchmarkPopBackLinkedListVsDoublyLinkedList(b *testing.B) {
	b.Run("LinkedListPopBackO(n)", func(b *testing.B) {
		ll := NewLinkedList[int]()
		for i := 0; i < 100; i++ {
			ll.PushBack(i)
		}
		b.ResetTimer()
		for b.Loop() {
			if ll.IsEmpty() {
				for i := 0; i < 100; i++ {
					ll.PushBack(i)
				}
			}
			ll.PopBack() // O(n) operation
		}
	})

	b.Run("DoublyLinkedListPopBackO(1)", func(b *testing.B) {
		dll := NewDoublyLinkedList[int]()
		for i := 0; i < 100; i++ {
			dll.PushBack(i)
		}
		b.ResetTimer()
		for b.Loop() {
			if dll.IsEmpty() {
				for i := 0; i < 100; i++ {
					dll.PushBack(i)
				}
			}
			dll.PopBack() // O(1) operation
		}
	})
}
