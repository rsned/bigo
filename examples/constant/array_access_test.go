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
	"fmt"
	"testing"
)

func TestArrayAccessByIndex(t *testing.T) {
	tests := []struct {
		name  string
		arr   []int
		index int
		want  int
	}{
		{"valid index", []int{1, 2, 3, 4, 5}, 2, 3},
		{"first index", []int{10, 20, 30}, 0, 10},
		{"last index", []int{10, 20, 30}, 2, 30},
		{"negative index", []int{1, 2, 3}, -1, 0},
		{"index out of bounds", []int{1, 2, 3}, 5, 0},
		{"empty array", []int{}, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ArrayAccessByIndex(tt.arr, tt.index)
			if got != tt.want {
				t.Errorf("ArrayAccessByIndex(%v, %d) = %d, want %d", tt.arr, tt.index, got, tt.want)
			}
		})
	}
}

// Benchmark functions for array access operations

func BenchmarkArrayAccessByIndex(b *testing.B) {
	index := 5
	sizes := []int{10, 100, 1000, 10000, 100000}

	for _, n := range sizes {
		vals := make([]int, n)
		for i := range vals {
			vals[i] = i + 1
		}

		b.Run(fmt.Sprintf("ArrayAccessByIndex-%d", n), func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				_ = ArrayAccessByIndex(vals, index)
			}
		})
	}
}
