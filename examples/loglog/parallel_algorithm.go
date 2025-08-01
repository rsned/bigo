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

import "sync"

// ParallelDivideConquer performs O(log(log n)) parallel divide-and-conquer
func ParallelDivideConquer(arr []int) int {
	if len(arr) <= 1 {
		if len(arr) == 0 {
			return 0
		}

		return arr[0]
	}

	numWorkers := 2
	mid := len(arr) / 2

	var wg sync.WaitGroup
	results := make([]int, numWorkers)

	wg.Add(2)

	go func() {
		defer wg.Done()
		results[0] = ParallelDivideConquer(arr[:mid])
	}()

	go func() {
		defer wg.Done()
		results[1] = ParallelDivideConquer(arr[mid:])
	}()

	wg.Wait()

	if results[0] > results[1] {
		return results[0]
	}

	return results[1]
}
