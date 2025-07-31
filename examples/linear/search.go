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

// Search performs O(n) linear search on unsorted array.
// This demonstrates linear time complexity because in the worst case,
// we must examine every element to find the target or determine it's
// not present. For unsorted data, this is the best we can do.
func Search(arr []int, target int) int {
	// Check each element sequentially until target is found
	for i, val := range arr {
		// Compare current element to target
		if val == target {
			return i // Return index when found
		}
	}

	// Return -1 if target not found after checking all elements
	return -1
}
