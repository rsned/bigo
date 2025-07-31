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

import "github.com/rsned/bigo/examples/datatypes/tree"

// FindTreeHeight performs O(n) height calculation for a BST.
func FindTreeHeight(root *tree.BSTNode) int {
	if root == nil {
		return 0
	}

	leftHeight := FindTreeHeight(root.Left)
	rightHeight := FindTreeHeight(root.Right)

	if leftHeight > rightHeight {
		return leftHeight + 1
	}

	return rightHeight + 1
}

// IsBalanced checks if tree is balanced in O(log n) for a BST.
func IsBalanced(root *tree.BSTNode) bool {
	_, balanced := checkHeight(root)

	return balanced
}

func checkHeight(root *tree.BSTNode) (int, bool) {
	if root == nil {
		return 0, true
	}

	leftHeight, leftBalanced := checkHeight(root.Left)
	if !leftBalanced {
		return 0, false
	}

	rightHeight, rightBalanced := checkHeight(root.Right)
	if !rightBalanced {
		return 0, false
	}

	diff := leftHeight - rightHeight
	if diff < 0 {
		diff = -diff
	}

	if diff > 1 {
		return 0, false
	}

	maxHeight := leftHeight
	if rightHeight > leftHeight {
		maxHeight = rightHeight
	}

	return maxHeight + 1, true
}
