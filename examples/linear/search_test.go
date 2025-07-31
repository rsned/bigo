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

package linear

import "testing"

func TestSearch(t *testing.T) {
	tests := []struct {
		name   string
		arr    []int
		target int
		want   int
	}{
		{"found at beginning", []int{1, 2, 3, 4, 5}, 1, 0},
		{"found at middle", []int{1, 2, 3, 4, 5}, 3, 2},
		{"found at end", []int{1, 2, 3, 4, 5}, 5, 4},
		{"not found", []int{1, 2, 3, 4, 5}, 6, -1},
		{"empty array", []int{}, 1, -1},
		{"single element found", []int{42}, 42, 0},
		{"single element not found", []int{42}, 5, -1},
		{"duplicate elements", []int{1, 2, 2, 3}, 2, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Search(tt.arr, tt.target)
			if got != tt.want {
				t.Errorf("Search(%v, %d) = %d, want %d", tt.arr, tt.target, got, tt.want)
			}
		})
	}
}

// Benchmark functions for linear search

func BenchmarkSearch(b *testing.B) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for b.Loop() {
		_ = Search(arr, arr[len(arr)-1])
	}
}
