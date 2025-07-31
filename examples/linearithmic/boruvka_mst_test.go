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

package linearithmic

import (
	"fmt"
	"math"
	"testing"
)

func TestBoruvkaMST(t *testing.T) {
	testCases := []struct {
		name           string
		graph          *BoruvkaGraph
		expectedWeight float64
		expectedEdges  int
	}{
		{
			name:           "Sample graph",
			graph:          BuildBoruvkaSampleGraph(),
			expectedWeight: 39.0, // Expected MST weight for the sample graph
			expectedEdges:  6,
		},
		{
			name: "Triangle graph",
			graph: &BoruvkaGraph{
				Vertices: 3,
				Edges: []BoruvkaEdge{
					{0, 1, 1.0},
					{1, 2, 2.0},
					{0, 2, 3.0},
				},
			},
			expectedWeight: 3.0, // 1.0 + 2.0 = 3.0
			expectedEdges:  2,
		},
		{
			name: "Single edge",
			graph: &BoruvkaGraph{
				Vertices: 2,
				Edges: []BoruvkaEdge{
					{0, 1, 5.0},
				},
			},
			expectedWeight: 5.0,
			expectedEdges:  1,
		},
		{
			name: "Square graph",
			graph: &BoruvkaGraph{
				Vertices: 4,
				Edges: []BoruvkaEdge{
					{0, 1, 1.0},
					{1, 2, 2.0},
					{2, 3, 3.0},
					{3, 0, 4.0},
					{0, 2, 5.0}, // diagonal
					{1, 3, 6.0}, // diagonal
				},
			},
			expectedWeight: 6.0, // 1.0 + 2.0 + 3.0 = 6.0
			expectedEdges:  3,
		},
		{
			name: "Star graph",
			graph: &BoruvkaGraph{
				Vertices: 5,
				Edges: []BoruvkaEdge{
					{0, 1, 1.0},
					{0, 2, 2.0},
					{0, 3, 3.0},
					{0, 4, 4.0},
					{1, 2, 10.0},
					{2, 3, 11.0},
					{3, 4, 12.0},
				},
			},
			expectedWeight: 10.0, // 1.0 + 2.0 + 3.0 + 4.0 = 10.0
			expectedEdges:  4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mst, totalWeight := BoruvkaMST(tc.graph)

			if len(mst) != tc.expectedEdges {
				t.Errorf("Expected %d edges in MST, got %d", tc.expectedEdges, len(mst))
			}

			if math.Abs(totalWeight-tc.expectedWeight) > 1e-9 {
				t.Errorf("Expected total weight %.2f, got %.2f", tc.expectedWeight, totalWeight)
			}

			// Validate MST properties
			if err := ValidateMST(tc.graph, mst); err != nil {
				t.Errorf("MST validation failed: %v", err)
			}
		})
	}
}

func TestBoruvkaUnionFind(t *testing.T) {
	uf := NewBoruvkaUnionFind(5)

	// Initially, all elements should be their own parent
	if uf.ComponentCount() != 5 {
		t.Errorf("Expected 5 components initially, got %d", uf.ComponentCount())
	}

	for i := 0; i < 5; i++ {
		if uf.Find(i) != i {
			t.Errorf("Expected Find(%d) = %d, got %d", i, i, uf.Find(i))
		}
	}

	// Union operations
	if !uf.Union(0, 1) {
		t.Error("Union(0, 1) should return true")
	}
	if uf.ComponentCount() != 4 {
		t.Errorf("Expected 4 components after Union(0, 1), got %d", uf.ComponentCount())
	}

	if !uf.Union(2, 3) {
		t.Error("Union(2, 3) should return true")
	}
	if uf.ComponentCount() != 3 {
		t.Errorf("Expected 3 components after Union(2, 3), got %d", uf.ComponentCount())
	}

	// Check that 0 and 1 are in the same set
	if uf.Find(0) != uf.Find(1) {
		t.Error("Expected 0 and 1 to be in the same set")
	}

	// Check that 2 and 3 are in the same set
	if uf.Find(2) != uf.Find(3) {
		t.Error("Expected 2 and 3 to be in the same set")
	}

	// Check that 0 and 2 are in different sets
	if uf.Find(0) == uf.Find(2) {
		t.Error("Expected 0 and 2 to be in different sets")
	}

	// Union the two sets
	if !uf.Union(1, 3) {
		t.Error("Union(1, 3) should return true")
	}
	if uf.ComponentCount() != 2 {
		t.Errorf("Expected 2 components after Union(1, 3), got %d", uf.ComponentCount())
	}

	// Now 0, 1, 2, 3 should be in the same set
	root := uf.Find(0)
	for i := 1; i < 4; i++ {
		if uf.Find(i) != root {
			t.Errorf("Expected all elements 0-3 to be in same set, but Find(%d) != Find(0)", i)
		}
	}

	// 4 should still be separate
	if uf.Find(4) == root {
		t.Error("Expected element 4 to be in a separate set")
	}

	// Try to union elements already in the same set
	if uf.Union(0, 2) {
		t.Error("Union(0, 2) should return false (already in same set)")
	}
	if uf.ComponentCount() != 2 {
		t.Errorf("Component count should remain 2, got %d", uf.ComponentCount())
	}
}

func TestValidateMST(t *testing.T) {
	graph := &BoruvkaGraph{
		Vertices: 4,
		Edges: []BoruvkaEdge{
			{0, 1, 1.0},
			{1, 2, 2.0},
			{2, 3, 3.0},
			{3, 0, 4.0},
		},
	}

	testCases := []struct {
		name        string
		mst         []BoruvkaEdge
		shouldError bool
	}{
		{
			name: "Valid MST",
			mst: []BoruvkaEdge{
				{0, 1, 1.0},
				{1, 2, 2.0},
				{2, 3, 3.0},
			},
			shouldError: false,
		},
		{
			name: "Too many edges",
			mst: []BoruvkaEdge{
				{0, 1, 1.0},
				{1, 2, 2.0},
				{2, 3, 3.0},
				{3, 0, 4.0},
			},
			shouldError: true,
		},
		{
			name: "Too few edges",
			mst: []BoruvkaEdge{
				{0, 1, 1.0},
				{1, 2, 2.0},
			},
			shouldError: true,
		},
		{
			name: "Contains cycle",
			mst: []BoruvkaEdge{
				{0, 1, 1.0},
				{1, 2, 2.0},
				{0, 2, 5.0}, // creates cycle with first two edges
			},
			shouldError: true,
		},
		{
			name: "Invalid vertex",
			mst: []BoruvkaEdge{
				{0, 1, 1.0},
				{1, 2, 2.0},
				{2, 5, 3.0}, // vertex 5 doesn't exist
			},
			shouldError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMST(graph, tc.mst)
			if tc.shouldError && err == nil {
				t.Error("Expected validation to fail, but it passed")
			}
			if !tc.shouldError && err != nil {
				t.Errorf("Expected validation to pass, but got error: %v", err)
			}
		})
	}
}

func TestCompareMSTAlgorithms(t *testing.T) {
	testGraphs := []*BoruvkaGraph{
		BuildBoruvkaSampleGraph(),
		{
			Vertices: 4,
			Edges: []BoruvkaEdge{
				{0, 1, 1.0},
				{1, 2, 2.0},
				{2, 3, 3.0},
				{3, 0, 4.0},
				{0, 2, 5.0},
				{1, 3, 6.0},
			},
		},
		BuildCompleteGraph(5),
	}

	for i, graph := range testGraphs {
		t.Run(fmt.Sprintf("Graph_%d", i), func(t *testing.T) {
			boruvkaWeight, kruskalWeight, equal := CompareMSTAlgorithms(graph)

			if !equal {
				t.Errorf("Borůvka and Kruskal algorithms produced different results: %.2f vs %.2f",
					boruvkaWeight, kruskalWeight)
			}

			if boruvkaWeight < 0 || kruskalWeight < 0 {
				t.Error("MST weights should be non-negative")
			}
		})
	}
}

func TestBuildCompleteGraph(t *testing.T) {
	sizes := []int{3, 4, 5, 10}

	for _, n := range sizes {
		t.Run(fmt.Sprintf("Size_%d", n), func(t *testing.T) {
			graph := BuildCompleteGraph(n)

			if graph.Vertices != n {
				t.Errorf("Expected %d vertices, got %d", n, graph.Vertices)
			}

			expectedEdges := n * (n - 1) / 2
			if len(graph.Edges) != expectedEdges {
				t.Errorf("Expected %d edges, got %d", expectedEdges, len(graph.Edges))
			}

			// Check that all edges have positive weights
			for _, edge := range graph.Edges {
				if edge.Weight <= 0 {
					t.Errorf("Edge (%d, %d) has non-positive weight: %.2f", edge.U, edge.V, edge.Weight)
				}
				if edge.U >= n || edge.V >= n || edge.U < 0 || edge.V < 0 {
					t.Errorf("Edge (%d, %d) has invalid vertices for graph of size %d", edge.U, edge.V, n)
				}
			}

			// Test that we can run Borůvka's algorithm on it
			mst, weight := BoruvkaMST(graph)
			if len(mst) != n-1 {
				t.Errorf("MST should have %d edges, got %d", n-1, len(mst))
			}
			if weight <= 0 {
				t.Errorf("MST weight should be positive, got %.2f", weight)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	testCases := []struct {
		name  string
		graph *BoruvkaGraph
	}{
		{
			name: "Empty graph",
			graph: &BoruvkaGraph{
				Vertices: 0,
				Edges:    []BoruvkaEdge{},
			},
		},
		{
			name: "Single vertex",
			graph: &BoruvkaGraph{
				Vertices: 1,
				Edges:    []BoruvkaEdge{},
			},
		},
		{
			name: "Disconnected graph",
			graph: &BoruvkaGraph{
				Vertices: 4,
				Edges: []BoruvkaEdge{
					{0, 1, 1.0},
					{2, 3, 2.0},
					// No edges connecting {0,1} component to {2,3} component
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mst, totalWeight := BoruvkaMST(tc.graph)

			if tc.graph.Vertices <= 1 {
				if len(mst) != 0 {
					t.Errorf("Expected empty MST for %s, got %d edges", tc.name, len(mst))
				}
				if totalWeight != 0.0 {
					t.Errorf("Expected zero total weight for %s, got %.2f", tc.name, totalWeight)
				}
			} else if totalWeight < 0 {
				// For disconnected graphs, we should get a forest (partial spanning tree)
				// The algorithm should still produce a valid result
				t.Errorf("Total weight should be non-negative, got %.2f", totalWeight)
			}
		})
	}
}

func BenchmarkBoruvkaMST(b *testing.B) {
	sizes := []int{10, 20, 50}

	for _, size := range sizes {
		graph := BuildCompleteGraph(size)
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				BoruvkaMST(graph)
			}
		})
	}
}

func BenchmarkBoruvkaVsKruskal(b *testing.B) {
	graph := BuildCompleteGraph(20)

	b.Run("Boruvka", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			BoruvkaMST(graph)
		}
	})

	b.Run("Kruskal", func(b *testing.B) {
		kruskalGraph := ConvertBoruvkaToKruskal(graph)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			KruskalMST(kruskalGraph)
		}
	})
}
