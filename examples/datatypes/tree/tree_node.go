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

package tree

// BSTNode represents a binary search tree node
type BSTNode struct {
	Val   int
	Left  *BSTNode
	Right *BSTNode
}

// NewBSTNode creates a new BSTNode with the given value
func NewBSTNode(val int) *BSTNode {
	return &BSTNode{Val: val, Left: nil, Right: nil}
}

// BuildBST creates a balanced BST from a sorted slice of integers
func BuildBST(values []int) *BSTNode {
	if len(values) == 0 {
		return nil
	}

	return buildBSTHelper(values, 0, len(values)-1)
}

func buildBSTHelper(values []int, start, end int) *BSTNode {
	if start > end {
		return nil
	}

	mid := start + (end-start)/2
	node := NewBSTNode(values[mid])
	node.Left = buildBSTHelper(values, start, mid-1)
	node.Right = buildBSTHelper(values, mid+1, end)

	return node
}

// InsertBST inserts a value into a BST
func (tn *BSTNode) InsertBST(val int) *BSTNode {
	if tn == nil {
		return NewBSTNode(val)
	}

	if val < tn.Val {
		tn.Left = tn.Left.InsertBST(val)
	} else if val > tn.Val {
		tn.Right = tn.Right.InsertBST(val)
	}

	return tn
}

// InorderTraversal returns values in inorder traversal
func (tn *BSTNode) InorderTraversal() []int {
	if tn == nil {
		return nil
	}

	var result []int
	result = append(result, tn.Left.InorderTraversal()...)
	result = append(result, tn.Val)
	result = append(result, tn.Right.InorderTraversal()...)

	return result
}
