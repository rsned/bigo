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

// VanEmdeBoas represents a simplified Van Emde Boas tree structure
type VanEmdeBoas struct {
	universeSize int
	min, max     *int
	summary      *VanEmdeBoas
	clusters     []*VanEmdeBoas
	clusterSize  int
}

// VanEmdeBoasTreePredecessor performs O(log(log n)) predecessor query
func (v *VanEmdeBoas) VanEmdeBoasTreePredecessor(x int) *int {
	if v.universeSize <= 2 {
		if x == 1 && v.min != nil && *v.min == 0 {
			return v.min
		}

		return nil
	}

	if v.max != nil && x > *v.max {
		return v.max
	}

	if v.min != nil && x > *v.min {
		clusterSize := int(math.Sqrt(float64(v.universeSize)))
		high := x / clusterSize
		low := x % clusterSize

		if v.clusters[high] != nil {
			if maxLow := v.clusters[high].VanEmdeBoasTreePredecessor(low); maxLow != nil {
				result := high*clusterSize + *maxLow

				return &result
			}
		}
	}

	return nil
}
