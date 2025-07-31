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

func TestKruskalMST(t *testing.T) {
	testCases := []struct {
		name           string
		graph          *KruskalGraph
		expectedWeight float64
		expectedEdges  int
	}{
		{
			name:           "Sample graph",
			graph:          BuildKruskalSampleGraph(),
			expectedWeight: 13.0, // Expected MST weight: 1.0 + 2.0 + 2.0 + 3.0 + 5.0 = 13.0
			expectedEdges:  5,
		},
		{
			name: "Triangle graph",
			graph: &KruskalGraph{
				Vertices: 3,
				Edges: []KruskalEdge{
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
			graph: &KruskalGraph{
				Vertices: 2,
				Edges: []KruskalEdge{
					{0, 1, 5.0},
				},
			},
			expectedWeight: 5.0,
			expectedEdges:  1,
		},
		{
			name: "Square graph",
			graph: &KruskalGraph{
				Vertices: 4,
				Edges: []KruskalEdge{
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
			graph: &KruskalGraph{
				Vertices: 5,
				Edges: []KruskalEdge{
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
		{
			name: "Path graph",
			graph: &KruskalGraph{
				Vertices: 4,
				Edges: []KruskalEdge{
					{0, 1, 2.0},
					{1, 2, 1.0},
					{2, 3, 3.0},
				},
			},
			expectedWeight: 6.0, // 2.0 + 1.0 + 3.0 = 6.0
			expectedEdges:  3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mst, totalWeight := KruskalMST(tc.graph)

			if len(mst) != tc.expectedEdges {
				t.Errorf("Expected %d edges in MST, got %d", tc.expectedEdges, len(mst))
			}

			if math.Abs(totalWeight-tc.expectedWeight) > 1e-9 {
				t.Errorf("Expected total weight %.2f, got %.2f", tc.expectedWeight, totalWeight)
			}

			// Validate MST properties
			if err := ValidateKruskalMST(tc.graph, mst); err != nil {
				t.Errorf("MST validation failed: %v", err)
			}
		})
	}
}

func TestKruskalMSTWithCustomSort(t *testing.T) {
	testCases := []struct {
		name           string
		graph          *KruskalGraph
		expectedWeight float64
		expectedEdges  int
	}{
		{
			name:           "Sample graph with custom sort",
			graph:          BuildKruskalSampleGraph(),
			expectedWeight: 13.0,
			expectedEdges:  5,
		},
		{
			name: "Triangle graph with custom sort",
			graph: &KruskalGraph{
				Vertices: 3,
				Edges: []KruskalEdge{
					{0, 1, 1.0},
					{1, 2, 2.0},
					{0, 2, 3.0},
				},
			},
			expectedWeight: 3.0,
			expectedEdges:  2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mst, totalWeight := KruskalMSTWithCustomSort(tc.graph)

			if len(mst) != tc.expectedEdges {
				t.Errorf("Expected %d edges in MST, got %d", tc.expectedEdges, len(mst))
			}

			if math.Abs(totalWeight-tc.expectedWeight) > 1e-9 {
				t.Errorf("Expected total weight %.2f, got %.2f", tc.expectedWeight, totalWeight)
			}

			// Validate MST properties
			if err := ValidateKruskalMST(tc.graph, mst); err != nil {
				t.Errorf("MST validation failed: %v", err)
			}
		})
	}
}

func TestKruskalUnionFind(t *testing.T) {
	uf := NewKruskalUnionFind(5)

	// Initially, all elements should be their own parent
	for i := 0; i < 5; i++ {
		if uf.Find(i) != i {
			t.Errorf("Expected Find(%d) = %d, got %d", i, i, uf.Find(i))
		}
	}

	// Union operations
	if !uf.Union(0, 1) {
		t.Error("Union(0, 1) should return true")
	}

	if !uf.Union(2, 3) {
		t.Error("Union(2, 3) should return true")
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
}

func TestValidateKruskalMST(t *testing.T) {
	graph := &KruskalGraph{
		Vertices: 4,
		Edges: []KruskalEdge{
			{0, 1, 1.0},
			{1, 2, 2.0},
			{2, 3, 3.0},
			{3, 0, 4.0},
		},
	}

	testCases := []struct {
		name        string
		mst         []KruskalEdge
		shouldError bool
	}{
		{
			name: "Valid MST",
			mst: []KruskalEdge{
				{0, 1, 1.0},
				{1, 2, 2.0},
				{2, 3, 3.0},
			},
			shouldError: false,
		},
		{
			name: "Too many edges",
			mst: []KruskalEdge{
				{0, 1, 1.0},
				{1, 2, 2.0},
				{2, 3, 3.0},
				{3, 0, 4.0},
			},
			shouldError: true,
		},
		{
			name: "Too few edges",
			mst: []KruskalEdge{
				{0, 1, 1.0},
				{1, 2, 2.0},
			},
			shouldError: true,
		},
		{
			name: "Contains cycle",
			mst: []KruskalEdge{
				{0, 1, 1.0},
				{1, 2, 2.0},
				{0, 2, 5.0}, // creates cycle with first two edges
			},
			shouldError: true,
		},
		{
			name: "Invalid vertex",
			mst: []KruskalEdge{
				{0, 1, 1.0},
				{1, 2, 2.0},
				{2, 5, 3.0}, // vertex 5 doesn't exist
			},
			shouldError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateKruskalMST(graph, tc.mst)
			if tc.shouldError && err == nil {
				t.Error("Expected validation to fail, but it passed")
			}
			if !tc.shouldError && err != nil {
				t.Errorf("Expected validation to pass, but got error: %v", err)
			}
		})
	}
}

func TestCompareKruskalSorts(t *testing.T) {
	testGraphs := []*KruskalGraph{
		BuildKruskalSampleGraph(),
		{
			Vertices: 4,
			Edges: []KruskalEdge{
				{0, 1, 1.0},
				{1, 2, 2.0},
				{2, 3, 3.0},
				{3, 0, 4.0},
				{0, 2, 5.0},
				{1, 3, 6.0},
			},
		},
		BuildKruskalCompleteGraph(5),
	}

	for i, graph := range testGraphs {
		t.Run(fmt.Sprintf("Graph_%d", i), func(t *testing.T) {
			standardWeight, customWeight, equal := CompareKruskalSorts(graph)

			if !equal {
				t.Errorf("Standard and custom sort produced different results: %.2f vs %.2f",
					standardWeight, customWeight)
			}

			if standardWeight < 0 || customWeight < 0 {
				t.Error("MST weights should be non-negative")
			}
		})
	}
}

func TestBuildKruskalCompleteGraph(t *testing.T) {
	sizes := []int{3, 4, 5, 10}

	for _, n := range sizes {
		t.Run(fmt.Sprintf("Size_%d", n), func(t *testing.T) {
			graph := BuildKruskalCompleteGraph(n)

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

			// Test that we can run Kruskal's algorithm on it
			mst, weight := KruskalMST(graph)
			if len(mst) != n-1 {
				t.Errorf("MST should have %d edges, got %d", n-1, len(mst))
			}
			if weight <= 0 {
				t.Errorf("MST weight should be positive, got %.2f", weight)
			}
		})
	}
}

func TestConvertBoruvkaToKruskal(t *testing.T) {
	boruvkaGraph := &BoruvkaGraph{
		Vertices: 3,
		Edges: []BoruvkaEdge{
			{0, 1, 1.5},
			{1, 2, 2.5},
			{0, 2, 3.5},
		},
	}

	kruskalGraph := ConvertBoruvkaToKruskal(boruvkaGraph)

	if kruskalGraph.Vertices != boruvkaGraph.Vertices {
		t.Errorf("Expected %d vertices, got %d", boruvkaGraph.Vertices, kruskalGraph.Vertices)
	}

	if len(kruskalGraph.Edges) != len(boruvkaGraph.Edges) {
		t.Errorf("Expected %d edges, got %d", len(boruvkaGraph.Edges), len(kruskalGraph.Edges))
	}

	for i, edge := range kruskalGraph.Edges {
		originalEdge := boruvkaGraph.Edges[i]
		if edge.U != originalEdge.U || edge.V != originalEdge.V || edge.Weight != originalEdge.Weight {
			t.Errorf("Edge %d mismatch: expected (%d, %d, %.2f), got (%d, %d, %.2f)",
				i, originalEdge.U, originalEdge.V, originalEdge.Weight,
				edge.U, edge.V, edge.Weight)
		}
	}
}

func TestConvertKruskalToBoruvka(t *testing.T) {
	kruskalGraph := &KruskalGraph{
		Vertices: 3,
		Edges: []KruskalEdge{
			{0, 1, 1.5},
			{1, 2, 2.5},
			{0, 2, 3.5},
		},
	}

	boruvkaGraph := ConvertKruskalToBoruvka(kruskalGraph)

	if boruvkaGraph.Vertices != kruskalGraph.Vertices {
		t.Errorf("Expected %d vertices, got %d", kruskalGraph.Vertices, boruvkaGraph.Vertices)
	}

	if len(boruvkaGraph.Edges) != len(kruskalGraph.Edges) {
		t.Errorf("Expected %d edges, got %d", len(kruskalGraph.Edges), len(boruvkaGraph.Edges))
	}

	for i, edge := range boruvkaGraph.Edges {
		originalEdge := kruskalGraph.Edges[i]
		if edge.U != originalEdge.U || edge.V != originalEdge.V || edge.Weight != originalEdge.Weight {
			t.Errorf("Edge %d mismatch: expected (%d, %d, %.2f), got (%d, %d, %.2f)",
				i, originalEdge.U, originalEdge.V, originalEdge.Weight,
				edge.U, edge.V, edge.Weight)
		}
	}
}

func TestKruskalEdgeCases(t *testing.T) {
	testCases := []struct {
		name  string
		graph *KruskalGraph
	}{
		{
			name: "Empty graph",
			graph: &KruskalGraph{
				Vertices: 0,
				Edges:    []KruskalEdge{},
			},
		},
		{
			name: "Single vertex",
			graph: &KruskalGraph{
				Vertices: 1,
				Edges:    []KruskalEdge{},
			},
		},
		{
			name: "Disconnected graph",
			graph: &KruskalGraph{
				Vertices: 4,
				Edges: []KruskalEdge{
					{0, 1, 1.0},
					{2, 3, 2.0},
					// No edges connecting {0,1} component to {2,3} component
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mst, totalWeight := KruskalMST(tc.graph)

			if tc.graph.Vertices <= 1 {
				if len(mst) != 0 {
					t.Errorf("Expected empty MST for %s, got %d edges", tc.name, len(mst))
				}
				if totalWeight != 0.0 {
					t.Errorf("Expected zero total weight for %s, got %.2f", tc.name, totalWeight)
				}
			} else if totalWeight < 0 {
				// For disconnected graphs, we should get a forest (partial spanning tree)
				t.Errorf("Total weight should be non-negative, got %.2f", totalWeight)
			}
		})
	}
}

func BenchmarkKruskalMST(b *testing.B) {
	sizes := []int{10, 20, 50}

	for _, size := range sizes {
		graph := BuildKruskalCompleteGraph(size)
		b.Run(fmt.Sprintf("StandardSort_Size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				KruskalMST(graph)
			}
		})
	}
}

func BenchmarkKruskalMSTWithCustomSort(b *testing.B) {
	sizes := []int{10, 20} // Smaller sizes for bubble sort

	for _, size := range sizes {
		graph := BuildKruskalCompleteGraph(size)
		b.Run(fmt.Sprintf("CustomSort_Size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				KruskalMSTWithCustomSort(graph)
			}
		})
	}
}

func BenchmarkKruskalSortComparison(b *testing.B) {
	graph := BuildKruskalCompleteGraph(20)

	b.Run("StandardSort", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			KruskalMST(graph)
		}
	})

	b.Run("CustomBubbleSort", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			KruskalMSTWithCustomSort(graph)
		}
	})
}
