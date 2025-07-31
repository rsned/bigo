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

// InterpolationSearch reports the position of the target value in the
// given slice or -1 if the target was not found.
//
// It performs O(log(log n)) search in a uniformly distributed sorted array.
func InterpolationSearch(arr []int, target int) int {
	if len(arr) == 0 {
		return -1
	}

	low := 0
	high := len(arr) - 1

	for low <= high && target >= arr[low] && target <= arr[high] {
		if low == high {
			if arr[low] == target {
				return low
			}

			return -1
		}

		// Interpolation formula
		// Handle case where all elements in range are same
		if arr[high] == arr[low] {
			if arr[low] == target {
				return low
			}

			break
		}

		pos := low + ((target-arr[low])*(high-low))/(arr[high]-arr[low])

		if pos < low || pos > high {
			break
		}

		if arr[pos] == target {
			return pos
		}

		if arr[pos] < target {
			low = pos + 1
		} else {
			high = pos - 1
		}
	}

	return -1
}
