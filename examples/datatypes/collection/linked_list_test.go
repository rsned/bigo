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

func TestNewLinkedList(t *testing.T) {
	t.Run("int list", func(t *testing.T) {
		ll := NewLinkedList[int]()
		if ll.Len() != 0 || !ll.IsEmpty() {
			t.Errorf("NewLinkedList[int]() should create empty list, got len=%d, empty=%v", ll.Len(), ll.IsEmpty())
		}
	})

	t.Run("string list", func(t *testing.T) {
		ll := NewLinkedList[string]()
		if ll.Len() != 0 || !ll.IsEmpty() {
			t.Errorf("NewLinkedList[string]() should create empty list, got len=%d, empty=%v", ll.Len(), ll.IsEmpty())
		}
	})
}

func TestNewLinkedListWithValues(t *testing.T) {
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
			ll := NewLinkedListWithValues(tt.values...)
			got := ll.ToSlice()
			if !cmp.Equal(got, tt.want) {
				t.Errorf("NewLinkedListWithValues(%v) = %v, want %v", tt.values, got, tt.want)
			}
			if ll.Len() != len(tt.want) {
				t.Errorf("NewLinkedListWithValues(%v) length = %d, want %d", tt.values, ll.Len(), len(tt.want))
			}
		})
	}
}

func TestFromSlice(t *testing.T) {
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
			ll := FromSlice(tt.slice)
			got := ll.ToSlice()
			if !cmp.Equal(got, tt.want) {
				t.Errorf("FromSlice(%v) = %v, want %v", tt.slice, got, tt.want)
			}
		})
	}
}

func TestLinkedListPushFront(t *testing.T) {
	ll := NewLinkedList[int]()

	// Test pushing to empty list
	ll.PushFront(1)
	if got := ll.ToSlice(); !cmp.Equal(got, []int{1}) {
		t.Errorf("After PushFront(1) on empty list: got %v, want [1]", got)
	}

	// Test pushing to non-empty list
	ll.PushFront(2)
	if got := ll.ToSlice(); !cmp.Equal(got, []int{2, 1}) {
		t.Errorf("After PushFront(2): got %v, want [2, 1]", got)
	}

	// Test length tracking
	if ll.Len() != 2 {
		t.Errorf("Length should be 2, got %d", ll.Len())
	}
}

func TestLinkedListPushBack(t *testing.T) {
	ll := NewLinkedList[int]()

	// Test pushing to empty list
	ll.PushBack(1)
	if got := ll.ToSlice(); !cmp.Equal(got, []int{1}) {
		t.Errorf("After PushBack(1) on empty list: got %v, want [1]", got)
	}

	// Test pushing to non-empty list
	ll.PushBack(2)
	if got := ll.ToSlice(); !cmp.Equal(got, []int{1, 2}) {
		t.Errorf("After PushBack(2): got %v, want [1, 2]", got)
	}

	// Test length tracking
	if ll.Len() != 2 {
		t.Errorf("Length should be 2, got %d", ll.Len())
	}
}

func TestLinkedListPopFront(t *testing.T) {
	ll := FromSlice([]int{1, 2, 3})

	// Pop from non-empty list
	val, ok := ll.PopFront()
	if !ok || val != 1 {
		t.Errorf("PopFront() = (%d, %v), want (1, true)", val, ok)
	}
	if got := ll.ToSlice(); !cmp.Equal(got, []int{2, 3}) {
		t.Errorf("After PopFront(): got %v, want [2, 3]", got)
	}

	// Pop until empty
	ll.PopFront()
	ll.PopFront()

	// Pop from empty list
	val, ok = ll.PopFront()
	if ok || val != 0 {
		t.Errorf("PopFront() on empty list = (%d, %v), want (0, false)", val, ok)
	}
}

func TestLinkedListPopBack(t *testing.T) {
	ll := FromSlice([]int{1, 2, 3})

	// Pop from non-empty list
	val, ok := ll.PopBack()
	if !ok || val != 3 {
		t.Errorf("PopBack() = (%d, %v), want (3, true)", val, ok)
	}
	if got := ll.ToSlice(); !cmp.Equal(got, []int{1, 2}) {
		t.Errorf("After PopBack(): got %v, want [1, 2]", got)
	}

	// Pop until empty
	ll.PopBack()
	ll.PopBack()

	// Pop from empty list
	val, ok = ll.PopBack()
	if ok || val != 0 {
		t.Errorf("PopBack() on empty list = (%d, %v), want (0, false)", val, ok)
	}
}

func TestLinkedListFront(t *testing.T) {
	// Empty list
	ll := NewLinkedList[int]()
	val, ok := ll.Front()
	if ok || val != 0 {
		t.Errorf("Front() on empty list = (%d, %v), want (0, false)", val, ok)
	}

	// Non-empty list
	ll = FromSlice([]int{1, 2, 3})
	val, ok = ll.Front()
	if !ok || val != 1 {
		t.Errorf("Front() = (%d, %v), want (1, true)", val, ok)
	}
}

func TestLinkedListBack(t *testing.T) {
	// Empty list
	ll := NewLinkedList[int]()
	val, ok := ll.Back()
	if ok || val != 0 {
		t.Errorf("Back() on empty list = (%d, %v), want (0, false)", val, ok)
	}

	// Non-empty list
	ll = FromSlice([]int{1, 2, 3})
	val, ok = ll.Back()
	if !ok || val != 3 {
		t.Errorf("Back() = (%d, %v), want (3, true)", val, ok)
	}
}

func TestLinkedListAt(t *testing.T) {
	ll := FromSlice([]int{10, 20, 30, 40})

	tests := []struct {
		index   int
		wantVal int
		wantOk  bool
	}{
		{0, 10, true},
		{1, 20, true},
		{2, 30, true},
		{3, 40, true},
		{4, 0, false},  // out of bounds
		{-1, 0, false}, // negative index
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			gotVal, gotOk := ll.At(tt.index)
			if gotVal != tt.wantVal || gotOk != tt.wantOk {
				t.Errorf("At(%d) = (%d, %v), want (%d, %v)", tt.index, gotVal, gotOk, tt.wantVal, tt.wantOk)
			}
		})
	}
}

func TestLinkedListInsert(t *testing.T) {
	t.Run("valid insertions", func(t *testing.T) {
		ll := FromSlice([]int{1, 3, 5})

		// Insert at beginning
		err := ll.Insert(0, 0)
		if err != nil {
			t.Errorf("Insert(0, 0) error = %v", err)
		}
		if got := ll.ToSlice(); !cmp.Equal(got, []int{0, 1, 3, 5}) {
			t.Errorf("After Insert(0, 0): got %v, want [0, 1, 3, 5]", got)
		}

		// Insert in middle
		err = ll.Insert(2, 2)
		if err != nil {
			t.Errorf("Insert(2, 2) error = %v", err)
		}
		if got := ll.ToSlice(); !cmp.Equal(got, []int{0, 1, 2, 3, 5}) {
			t.Errorf("After Insert(2, 2): got %v, want [0, 1, 2, 3, 5]", got)
		}

		// Insert at end
		err = ll.Insert(5, 6)
		if err != nil {
			t.Errorf("Insert(5, 6) error = %v", err)
		}
		if got := ll.ToSlice(); !cmp.Equal(got, []int{0, 1, 2, 3, 5, 6}) {
			t.Errorf("After Insert(5, 6): got %v, want [0, 1, 2, 3, 5, 6]", got)
		}
	})

	t.Run("invalid insertions", func(t *testing.T) {
		ll := FromSlice([]int{1, 2, 3})

		if err := ll.Insert(-1, 0); err == nil {
			t.Error("Insert(-1, 0) should return error")
		}
		if err := ll.Insert(4, 0); err == nil {
			t.Error("Insert(4, 0) should return error")
		}
	})
}

func TestLinkedListRemove(t *testing.T) {
	t.Run("valid removals", func(t *testing.T) {
		ll := FromSlice([]int{1, 2, 3, 4, 5})

		// Remove from beginning
		val, err := ll.Remove(0)
		if err != nil || val != 1 {
			t.Errorf("Remove(0) = (%d, %v), want (1, nil)", val, err)
		}
		if got := ll.ToSlice(); !cmp.Equal(got, []int{2, 3, 4, 5}) {
			t.Errorf("After Remove(0): got %v, want [2, 3, 4, 5]", got)
		}

		// Remove from middle
		val, err = ll.Remove(1)
		if err != nil || val != 3 {
			t.Errorf("Remove(1) = (%d, %v), want (3, nil)", val, err)
		}
		if got := ll.ToSlice(); !cmp.Equal(got, []int{2, 4, 5}) {
			t.Errorf("After Remove(1): got %v, want [2, 4, 5]", got)
		}

		// Remove from end
		val, err = ll.Remove(2)
		if err != nil || val != 5 {
			t.Errorf("Remove(2) = (%d, %v), want (5, nil)", val, err)
		}
		if got := ll.ToSlice(); !cmp.Equal(got, []int{2, 4}) {
			t.Errorf("After Remove(2): got %v, want [2, 4]", got)
		}
	})

	t.Run("invalid removals", func(t *testing.T) {
		ll := FromSlice([]int{1, 2, 3})

		if _, err := ll.Remove(-1); err == nil {
			t.Error("Remove(-1) should return error")
		}
		if _, err := ll.Remove(3); err == nil {
			t.Error("Remove(3) should return error")
		}
	})
}

func TestLinkedListFind(t *testing.T) {
	ll := FromSlice([]int{10, 20, 30, 20, 40})

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
			got := ll.Find(tt.value)
			if got != tt.want {
				t.Errorf("Find(%d) = %d, want %d", tt.value, got, tt.want)
			}
		})
	}
}

func TestLinkedListContains(t *testing.T) {
	ll := FromSlice([]int{1, 2, 3, 4, 5})

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
			got := ll.Contains(tt.value)
			if got != tt.want {
				t.Errorf("Contains(%d) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestLinkedListClear(t *testing.T) {
	ll := FromSlice([]int{1, 2, 3, 4, 5})
	ll.Clear()

	if !ll.IsEmpty() || ll.Len() != 0 {
		t.Errorf("After Clear(): IsEmpty()=%v, Len()=%d, want true, 0", ll.IsEmpty(), ll.Len())
	}

	if got := ll.ToSlice(); len(got) != 0 {
		t.Errorf("After Clear(): ToSlice()=%v, want []", got)
	}
}

func TestLinkedListIterator(t *testing.T) {
	ll := FromSlice([]int{1, 2, 3})
	it := ll.Iterator()

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
	empty := NewLinkedList[int]()
	emptyIt := empty.Iterator()
	if emptyIt.HasNext() {
		t.Error("Empty list iterator should not have next")
	}
}

// Benchmark functions for performance testing

func BenchmarkNewLinkedList(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewLinkedList[int]()
	}
}

func BenchmarkLinkedListPushFront(b *testing.B) {
	ll := NewLinkedList[int]()
	b.ResetTimer()
	for b.Loop() {
		ll.PushFront(42)
	}
}

func BenchmarkLinkedListPushBack(b *testing.B) {
	ll := NewLinkedList[int]()
	b.ResetTimer()
	for b.Loop() {
		ll.PushBack(42)
	}
}

func BenchmarkLinkedListPopFront(b *testing.B) {
	ll := NewLinkedList[int]()
	for i := range 1000 {
		ll.PushBack(i)
	}
	b.ResetTimer()
	for b.Loop() {
		if ll.IsEmpty() {
			for i := range 1000 {
				ll.PushBack(i)
			}
		}
		ll.PopFront()
	}
}

func BenchmarkLinkedListPopBack(b *testing.B) {
	ll := NewLinkedList[int]()
	for i := range 1000 {
		ll.PushBack(i)
	}
	b.ResetTimer()
	for b.Loop() {
		if ll.IsEmpty() {
			for i := range 1000 {
				ll.PushBack(i)
			}
		}
		ll.PopBack()
	}
}

func BenchmarkLinkedListToSlice(b *testing.B) {
	ll := FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	b.ResetTimer()
	for b.Loop() {
		_ = ll.ToSlice()
	}
}

func BenchmarkLinkedListFind(b *testing.B) {
	ll := NewLinkedList[int]()
	for i := range 1000 {
		ll.PushBack(i)
	}
	b.ResetTimer()
	for b.Loop() {
		_ = ll.Find(500)
	}
}

func BenchmarkLinkedListIterator(b *testing.B) {
	ll := FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	b.ResetTimer()
	for b.Loop() {
		it := ll.Iterator()
		for it.HasNext() {
			_, _ = it.Next()
		}
	}
}
