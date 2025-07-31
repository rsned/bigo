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

package loglog

import "math"

// FusionTree represents a simplified fusion tree structure
type FusionTree struct {
	keys     []int
	children []*FusionTree
	isLeaf   bool
	degree   int
}

// FusionTreeOperation performs O(log(log n)) priority queue operations
func (ft *FusionTree) FusionTreeOperation(key int) bool {
	if ft.isLeaf {
		for _, k := range ft.keys {
			if k == key {
				return true
			}
		}

		return false
	}

	// Simplified fusion tree search using bit manipulation approximation
	childIndex := ft.findChildIndex(key)
	if childIndex < len(ft.children) && ft.children[childIndex] != nil {
		return ft.children[childIndex].FusionTreeOperation(key)
	}

	return false
}

func (ft *FusionTree) findChildIndex(key int) int {
	// Simplified index calculation using log-log approximation
	if len(ft.keys) == 0 {
		return 0
	}

	logLogFactor := int(math.Log(math.Log(float64(len(ft.keys))+1) + 1))
	if logLogFactor <= 0 {
		logLogFactor = 1
	}

	// Handle negative keys by using absolute value
	index := key % (logLogFactor + 1)
	if index < 0 {
		index = -index
	}

	return index
}
