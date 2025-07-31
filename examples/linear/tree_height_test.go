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

import (
	"fmt"
	"testing"

	"github.com/rsned/bigo/examples/datatypes/tree"
)

func TestFindTreeHeight(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   int
	}{
		{"nil tree", nil, 0},
		{"single node", []int{5}, 1},
		{"balanced tree", []int{1, 2, 3, 4, 5, 6, 7}, 3},
		{"larger balanced tree", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var root *tree.BSTNode
			if tt.values != nil {
				root = tree.BuildBST(tt.values)
			}

			got := FindTreeHeight(root)
			if got != tt.want {
				t.Errorf("FindTreeHeight() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestFindTreeHeightUnbalanced(t *testing.T) {
	// Create an unbalanced tree manually: 1 -> 2 -> 3 -> 4
	root := tree.NewBSTNode(1)
	root.Right = tree.NewBSTNode(2)
	root.Right.Right = tree.NewBSTNode(3)
	root.Right.Right.Right = tree.NewBSTNode(4)

	got := FindTreeHeight(root)
	want := 4
	if got != want {
		t.Errorf("FindTreeHeight(unbalanced) = %d, want %d", got, want)
	}
}

func TestIsBalanced(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   bool
	}{
		{"nil tree", nil, true},
		{"single node", []int{5}, true},
		{"balanced tree", []int{1, 2, 3, 4, 5, 6, 7}, true},
		{"larger balanced tree", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var root *tree.BSTNode
			if tt.values != nil {
				root = tree.BuildBST(tt.values)
			}

			got := IsBalanced(root)
			if got != tt.want {
				t.Errorf("IsBalanced() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestIsBalanced_Unbalanced(t *testing.T) {
	// Create an unbalanced tree manually: 1 -> 2 -> 3 -> 4
	root := tree.NewBSTNode(1)
	root.Right = tree.NewBSTNode(2)
	root.Right.Right = tree.NewBSTNode(3)
	root.Right.Right.Right = tree.NewBSTNode(4)

	got := IsBalanced(root)
	if got {
		t.Errorf("IsBalanced(unbalanced tree) = true, want false")
	}
}

func TestIsBalanced_SlightlyUnbalanced(t *testing.T) {
	// Create a tree that's slightly unbalanced but still valid
	root := tree.NewBSTNode(2)
	root.Left = tree.NewBSTNode(1)
	root.Right = tree.NewBSTNode(3)
	root.Right.Right = tree.NewBSTNode(4)
	root.Right.Right.Right = tree.NewBSTNode(5)

	got := IsBalanced(root)
	if got {
		t.Errorf("IsBalanced(slightly unbalanced tree) = true, want false")
	}
}

func TestTreeHeightAndBalance_Integration(t *testing.T) {
	// Test that balanced trees built with BuildBST are indeed balanced
	testSizes := []int{1, 3, 7, 15, 31}

	for _, size := range testSizes {
		values := make([]int, size)
		for i := 0; i < size; i++ {
			values[i] = i + 1
		}

		root := tree.BuildBST(values)

		// Check that it's balanced
		if !IsBalanced(root) {
			t.Errorf("BuildBST(%d elements) should create balanced tree", size)
		}

		// Check height is logarithmic
		height := FindTreeHeight(root)
		expectedMaxHeight := 0
		temp := size
		for temp > 0 {
			expectedMaxHeight++
			temp /= 2
		}

		if height > expectedMaxHeight+1 {
			t.Errorf("BuildBST(%d elements) height = %d, expected <= %d", size, height, expectedMaxHeight+1)
		}
	}
}

// Benchmark functions for tree height operations

func BenchmarkFindTreeHeight(b *testing.B) {
	// Test with various tree sizes
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		values := make([]int, size)
		for i := 0; i < size; i++ {
			values[i] = i + 1
		}
		root := tree.BuildBST(values)

		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				FindTreeHeight(root)
			}
		})
	}
}

func BenchmarkIsBalanced(b *testing.B) {
	// Test with various tree sizes
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		values := make([]int, size)
		for i := 0; i < size; i++ {
			values[i] = i + 1
		}
		root := tree.BuildBST(values)

		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				IsBalanced(root)
			}
		})
	}
}

func BenchmarkBuildBSTForHeight(b *testing.B) {
	// Test BST construction with various sizes
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		values := make([]int, size)
		for i := range size {
			values[i] = i + 1
		}

		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			b.ResetTimer()
			for range b.N {
				tree.BuildBST(values)
			}
		})
	}
}
