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

func TestDynamicStackPush(t *testing.T) {
	tests := []struct {
		name     string
		initial  []int
		pushItem int
		want     []int
	}{
		{
			name:     "push to empty stack",
			initial:  []int{},
			pushItem: 5,
			want:     []int{5},
		},
		{
			name:     "push to existing stack",
			initial:  []int{1, 2},
			pushItem: 3,
			want:     []int{1, 2, 3},
		},
		{
			name:     "push negative value",
			initial:  []int{1},
			pushItem: -5,
			want:     []int{1, -5},
		},
		{
			name:     "push zero",
			initial:  []int{10},
			pushItem: 0,
			want:     []int{10, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := CopyDynamicStack(tt.initial)

			stack.Push(tt.pushItem)

			if !cmp.Equal(stack.items, tt.want) {
				t.Errorf("StackPush(%d) resulted in %v, want %v", tt.pushItem, stack.items, tt.want)
			}
		})
	}
}

func TestDynamicStackPop(t *testing.T) {
	tests := []struct {
		name          string
		initial       []int
		wantValue     int
		wantSuccess   bool
		wantRemaining []int
	}{
		{
			name:          "pop from single item",
			initial:       []int{5},
			wantValue:     5,
			wantSuccess:   true,
			wantRemaining: []int{},
		},
		{
			name:          "pop from multiple items",
			initial:       []int{1, 2, 3},
			wantValue:     3,
			wantSuccess:   true,
			wantRemaining: []int{1, 2},
		},
		{
			name:          "pop from empty stack",
			initial:       []int{},
			wantValue:     0,
			wantSuccess:   false,
			wantRemaining: []int{},
		},
		{
			name:          "pop negative value",
			initial:       []int{1, -5},
			wantValue:     -5,
			wantSuccess:   true,
			wantRemaining: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := CopyDynamicStack(tt.initial)

			gotValue, gotSuccess := stack.Pop()

			if gotValue != tt.wantValue || gotSuccess != tt.wantSuccess {
				t.Errorf("StackPop() = (%d, %t), want (%d, %t)",
					gotValue, gotSuccess, tt.wantValue, tt.wantSuccess)
			}

			if !cmp.Equal(stack.items, tt.wantRemaining) {
				t.Errorf("After StackPop(), stack contains %v, want %v", stack.items, tt.wantRemaining)
			}
		})
	}
}

func TestStackOperations_Integration(t *testing.T) {
	stack := NewDynamicStack()

	// Test push operations
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	if len(stack.items) != 3 {
		t.Errorf("After 3 pushes, stack size = %d, want 3", len(stack.items))
	}

	// Test pop operations
	val, ok := stack.Pop()
	if !ok || val != 3 {
		t.Errorf("First pop = (%d, %t), want (3, true)", val, ok)
	}

	val, ok = stack.Pop()
	if !ok || val != 2 {
		t.Errorf("Second pop = (%d, %t), want (2, true)", val, ok)
	}

	val, ok = stack.Pop()
	if !ok || val != 1 {
		t.Errorf("Third pop = (%d, %t), want (1, true)", val, ok)
	}

	val, ok = stack.Pop()
	if ok || val != 0 {
		t.Errorf("Pop from empty stack = (%d, %t), want (0, false)", val, ok)
	}
}

// Benchmark functions for stack operations

func BenchmarkDynamicStackPush(b *testing.B) {
	stack := &DynamicStack{
		items: nil,
	}
	b.ResetTimer()
	for b.Loop() {
		stack.Push(42)
	}
}

func BenchmarkDynamicStackPop(b *testing.B) {
	stack := &DynamicStack{
		items: nil,
	}
	stack.Push(42)
	stack.Push(43)
	stack.Push(44)

	b.ResetTimer()
	for b.Loop() {
		if len(stack.items) == 0 {
			stack.items = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100}
		}
		_, _ = stack.Pop()
	}
}
