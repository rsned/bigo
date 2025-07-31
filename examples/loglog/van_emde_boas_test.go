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

import (
	"fmt"
	"testing"
)

func TestVanEmdeBoas_VanEmdeBoasTreePredecessor(t *testing.T) {
	tests := []struct {
		name         string
		universeSize int
		min          *int
		max          *int
		x            int
		want         *int
	}{
		{
			"small universe size 2, x=1, min=0",
			2,
			intPtr(0),
			intPtr(1),
			1,
			intPtr(0),
		},
		{
			"small universe size 2, x=0, no predecessor",
			2,
			intPtr(0),
			intPtr(1),
			0,
			nil,
		},
		{
			"small universe size 1, no predecessor",
			1,
			nil,
			nil,
			0,
			nil,
		},
		{
			"x greater than max",
			16,
			intPtr(2),
			intPtr(10),
			15,
			intPtr(10),
		},
		{
			"x less than min",
			16,
			intPtr(5),
			intPtr(10),
			3,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			veb := &VanEmdeBoas{
				universeSize: tt.universeSize,
				min:          tt.min,
				max:          tt.max,
				summary:      nil,
				clusters:     nil,
				clusterSize:  0,
			}

			// Initialize clusters for larger universe sizes
			if tt.universeSize > 2 {
				veb.clusters = make([]*VanEmdeBoas, tt.universeSize)
			}

			got := veb.VanEmdeBoasTreePredecessor(tt.x)

			if (got == nil && tt.want != nil) || (got != nil && tt.want == nil) {
				t.Errorf("VanEmdeBoasTreePredecessor(%d) = %v, want %v", tt.x, got, tt.want)

				return
			}

			if got != nil && tt.want != nil && *got != *tt.want {
				t.Errorf("VanEmdeBoasTreePredecessor(%d) = %d, want %d", tt.x, *got, *tt.want)
			}
		})
	}
}

func TestVanEmdeBoas_EmptyTree(t *testing.T) {
	veb := &VanEmdeBoas{
		universeSize: 16,
		min:          nil,
		max:          nil,
		summary:      nil,
		clusters:     make([]*VanEmdeBoas, 16),
		clusterSize:  0,
	}

	got := veb.VanEmdeBoasTreePredecessor(8)
	if got != nil {
		t.Errorf("VanEmdeBoasTreePredecessor on empty tree should return nil, got %v", got)
	}
}

func TestVanEmdeBoas_SingleElement(t *testing.T) {
	veb := &VanEmdeBoas{
		universeSize: 16,
		min:          intPtr(5),
		max:          intPtr(5),
		summary:      nil,
		clusters:     make([]*VanEmdeBoas, 16),
		clusterSize:  0,
	}

	tests := []struct {
		x    int
		want *int
	}{
		{6, intPtr(5)}, // x > min, should return min as predecessor
		{5, nil},       // x == min, no predecessor
		{4, nil},       // x < min, no predecessor
	}

	for _, tt := range tests {
		got := veb.VanEmdeBoasTreePredecessor(tt.x)
		if (got == nil && tt.want != nil) || (got != nil && tt.want == nil) {
			t.Errorf("VanEmdeBoasTreePredecessor(%d) = %v, want %v", tt.x, got, tt.want)
		} else if got != nil && tt.want != nil && *got != *tt.want {
			t.Errorf("VanEmdeBoasTreePredecessor(%d) = %d, want %d", tt.x, *got, *tt.want)
		}
	}
}

func TestVanEmdeBoas_LargeUniverse(t *testing.T) {
	// Test with larger universe size to exercise the cluster logic
	veb := &VanEmdeBoas{
		universeSize: 256,
		min:          intPtr(10),
		max:          intPtr(200),
		summary:      nil,
		clusters:     make([]*VanEmdeBoas, 256),
		clusterSize:  0,
	}

	// Test basic functionality
	got := veb.VanEmdeBoasTreePredecessor(250)
	if got == nil || *got != 200 {
		t.Errorf("VanEmdeBoasTreePredecessor(250) should return max=200, got %v", got)
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}

// Benchmark functions for Van Emde Boas Tree (log(log n) complexity)

func BenchmarkVanEmdeBoasTreePredecessor(b *testing.B) {
	// Log-log complexity scales very well - can use larger universe sizes
	universeSizes := []int{16, 256, 4096, 65536}

	for _, universeSize := range universeSizes {
		veb := &VanEmdeBoas{
			universeSize: universeSize,
			min:          intPtr(1),
			max:          intPtr(universeSize - 1),
			summary:      nil,
			clusters:     make([]*VanEmdeBoas, universeSize),
			clusterSize:  0,
		}

		// Test queries at various positions

		b.Run(fmt.Sprintf("universe-size-%d", universeSize), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				x := (i % (universeSize - 1)) + 1
				veb.VanEmdeBoasTreePredecessor(x)
			}
		})
	}
}

func BenchmarkVanEmdeBoasTreePredecessorSmall(b *testing.B) {
	// Test small universe sizes (base cases)
	universeSizes := []int{2, 4, 8, 16}

	for _, universeSize := range universeSizes {
		veb := &VanEmdeBoas{
			universeSize: universeSize,
			min:          intPtr(0),
			max:          intPtr(universeSize - 1),
			summary:      nil,
			clusters:     nil,
			clusterSize:  0,
		}

		b.Run(fmt.Sprintf("small-universe-size-%d", universeSize), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				x := i % universeSize
				veb.VanEmdeBoasTreePredecessor(x)
			}
		})
	}
}

func BenchmarkVanEmdeBoasTreePredecessorDense(b *testing.B) {
	// Test with a more populated tree
	universeSizes := []int{64, 256, 1024}

	for _, universeSize := range universeSizes {
		veb := &VanEmdeBoas{
			universeSize: universeSize,
			min:          intPtr(1),
			max:          intPtr(universeSize - 2),
			summary:      nil,
			clusters:     make([]*VanEmdeBoas, universeSize),
			clusterSize:  0,
		}

		// Initialize some clusters to simulate populated tree
		for i := 0; i < len(veb.clusters); i += 4 {
			veb.clusters[i] = &VanEmdeBoas{
				universeSize: universeSize / 4,
				min:          intPtr(0),
				max:          intPtr(universeSize/4 - 1),
				summary:      nil,
				clusters:     nil,
				clusterSize:  0,
			}
		}

		b.Run(fmt.Sprintf("dense-universe-size-%d", universeSize), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				x := (i % (universeSize - 1)) + 1
				veb.VanEmdeBoasTreePredecessor(x)
			}
		})
	}
}
