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
	"math/rand/v2"
	"sort"
)

// To cut out some of the timing variability of benchmark functions, pre-create
// and sort a large set of random values.
var (
	testIntVals       []int
	testIntValsSorted []int
)

const (
	// limit is the maximum amount of random values to pre-generate for
	// the benchmarks to run against.
	limit = 10000000
)

func init() {
	testIntVals = make([]int, limit)
	testIntValsSorted = make([]int, limit)

	for i := range limit {
		testIntVals[i] = rand.Int()
	}

	copy(testIntValsSorted, testIntVals)

	sort.Ints(testIntValsSorted)
}
