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

// HashTableLookup performs O(1) hash table lookup (average case).
// Hash tables provide constant time access by computing a hash function
// to directly locate the bucket containing the key-value pair.
// Note: Worst case can be O(n) if all keys hash to the same bucket.
func HashTableLookup(hashMap map[int]int, key int) (int, bool) {
	// Go's map implementation uses a hash table internally
	// The hash function computes an index directly to the bucket
	// This avoids the need to search through all elements
	value, exists := hashMap[key]

	return value, exists
}
