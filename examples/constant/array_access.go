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

// ArrayAccessByIndex performs O(1) array access by index. This operation
// demonstrates constant time complexity because accessing an element by its
// index requires only a single memory lookup, regardless of array size.
// The time complexity remains O(1) whether the array has 10 or 10,000 elements.
func ArrayAccessByIndex(arr []int, index int) int {
	// Bounds checking to prevent index out of range errors
	// This check is also O(1) as it's a simple comparison
	if index < 0 || index >= len(arr) {
		// Return zero value for invalid indices
		return 0
	}

	// Direct array access using index - this is the core O(1) operation
	// CPU calculates memory address as: base_address + (index * element_size)
	return arr[index]
}
