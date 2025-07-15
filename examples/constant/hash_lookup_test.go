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

import "testing"

func TestHashTableLookup(t *testing.T) {
	tests := []struct {
		name       string
		hashMap    map[int]int
		key        int
		wantValue  int
		wantExists bool
	}{
		{"existing key", map[int]int{1: 10, 2: 20, 3: 30}, 2, 20, true},
		{"non-existing key", map[int]int{1: 10, 2: 20}, 3, 0, false},
		{"empty map", map[int]int{}, 1, 0, false},
		{"zero value", map[int]int{1: 0, 2: 10}, 1, 0, true},
		{"negative key", map[int]int{-1: 100, 0: 200}, -1, 100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotExists := HashTableLookup(tt.hashMap, tt.key)
			if gotValue != tt.wantValue || gotExists != tt.wantExists {
				t.Errorf("HashTableLookup(%v, %d) = (%d, %t), want (%d, %t)",
					tt.hashMap, tt.key, gotValue, gotExists, tt.wantValue, tt.wantExists)
			}
		})
	}
}

// Benchmark functions for hash table operations
func BenchmarkHashTableLookup(b *testing.B) {
	hashMap := map[int]int{1: 10, 2: 20, 3: 30, 4: 40, 5: 50}
	key := 3
	b.ResetTimer()
	for b.Loop() {
		_, _ = HashTableLookup(hashMap, key)
	}
}
