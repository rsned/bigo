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
	"sort"
)

// KruskalEdge represents a weighted edge in Kruskal's algorithm
type KruskalEdge struct {
	U, V   int     // vertices
	Weight float64 // edge weight
}

// KruskalGraph represents a weighted undirected graph for Kruskal's algorithm
type KruskalGraph struct {
	Vertices int
	Edges    []KruskalEdge
}

// KruskalUnionFind implements union-find with path compression and union by rank
type KruskalUnionFind struct {
	parent []int
	rank   []int
}

// NewKruskalUnionFind creates a new union-find structure
func NewKruskalUnionFind(n int) *KruskalUnionFind {
	uf := &KruskalUnionFind{
		parent: make([]int, n),
		rank:   make([]int, n),
	}

	for i := 0; i < n; i++ {
		uf.parent[i] = i
	}

	return uf
}

// Find returns the root of the set containing x with path compression
func (uf *KruskalUnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x])
	}

	return uf.parent[x]
}

// Union merges the sets containing x and y using union by rank
// Returns true if the union was successful (vertices were in different sets)
func (uf *KruskalUnionFind) Union(x, y int) bool {
	rootX := uf.Find(x)
	rootY := uf.Find(y)

	if rootX == rootY {
		return false // already in same component
	}

	switch {
	case uf.rank[rootX] < uf.rank[rootY]:
		uf.parent[rootX] = rootY
	case uf.rank[rootX] > uf.rank[rootY]:
		uf.parent[rootY] = rootX
	default:
		uf.parent[rootY] = rootX
		uf.rank[rootX]++
	}

	return true
}

// KruskalMST implements Kruskal's minimum spanning tree algorithm
// Time complexity: O(m log m) where m is the number of edges
// Space complexity: O(m + n) where n is the number of vertices
func KruskalMST(graph *KruskalGraph) ([]KruskalEdge, float64) {
	if graph.Vertices <= 1 {
		return []KruskalEdge{}, 0
	}

	// Copy and sort edges by weight
	edges := make([]KruskalEdge, len(graph.Edges))
	copy(edges, graph.Edges)

	sort.Slice(edges, func(i, j int) bool {
		return edges[i].Weight < edges[j].Weight
	})

	uf := NewKruskalUnionFind(graph.Vertices)
	mst := make([]KruskalEdge, 0, graph.Vertices-1)
	totalWeight := 0.0

	// Process edges in order of increasing weight
	for _, edge := range edges {
		if uf.Union(edge.U, edge.V) {
			mst = append(mst, edge)
			totalWeight += edge.Weight

			// Early termination when we have n-1 edges
			if len(mst) == graph.Vertices-1 {
				break
			}
		}
	}

	return mst, totalWeight
}

// KruskalMSTWithCustomSort implements Kruskal's algorithm with a custom sorting method
// This version uses a simple bubble sort for educational purposes
func KruskalMSTWithCustomSort(graph *KruskalGraph) ([]KruskalEdge, float64) {
	if graph.Vertices <= 1 {
		return []KruskalEdge{}, 0
	}

	// Copy edges
	edges := make([]KruskalEdge, len(graph.Edges))
	copy(edges, graph.Edges)

	// Bubble sort by weight (O(m^2) sorting for demonstration)
	n := len(edges)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if edges[j].Weight > edges[j+1].Weight {
				edges[j], edges[j+1] = edges[j+1], edges[j]
			}
		}
	}

	uf := NewKruskalUnionFind(graph.Vertices)
	mst := make([]KruskalEdge, 0, graph.Vertices-1)
	totalWeight := 0.0

	for _, edge := range edges {
		if uf.Union(edge.U, edge.V) {
			mst = append(mst, edge)
			totalWeight += edge.Weight

			if len(mst) == graph.Vertices-1 {
				break
			}
		}
	}

	return mst, totalWeight
}

// BuildKruskalSampleGraph creates a sample graph for testing Kruskal's algorithm
func BuildKruskalSampleGraph() *KruskalGraph {
	return &KruskalGraph{
		Vertices: 6,
		Edges: []KruskalEdge{
			{0, 1, 4.0},
			{0, 2, 2.0},
			{1, 2, 1.0},
			{1, 3, 5.0},
			{2, 3, 8.0},
			{2, 4, 10.0},
			{3, 4, 2.0},
			{3, 5, 6.0},
			{4, 5, 3.0},
		},
	}
}

// BuildKruskalCompleteGraph creates a complete graph with n vertices
func BuildKruskalCompleteGraph(n int) *KruskalGraph {
	graph := &KruskalGraph{
		Vertices: n,
		Edges:    make([]KruskalEdge, 0, n*(n-1)/2),
	}

	// Add all possible edges with deterministic weights
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			// Use a deterministic but varied weight calculation
			weight := float64((i*n+j)%100 + 1)
			graph.Edges = append(graph.Edges, KruskalEdge{
				U: i, V: j, Weight: weight,
			})
		}
	}

	return graph
}

// ValidateKruskalMST checks if the given edges form a valid minimum spanning tree
func ValidateKruskalMST(graph *KruskalGraph, mst []KruskalEdge) error {
	if len(mst) != graph.Vertices-1 {
		return fmt.Errorf("MST should have %d edges, got %d", graph.Vertices-1, len(mst))
	}

	// Check connectivity using union-find
	uf := NewKruskalUnionFind(graph.Vertices)

	for _, edge := range mst {
		if edge.U < 0 || edge.U >= graph.Vertices ||
			edge.V < 0 || edge.V >= graph.Vertices {
			return fmt.Errorf("invalid edge vertices: (%d, %d)", edge.U, edge.V)
		}

		if !uf.Union(edge.U, edge.V) {
			return fmt.Errorf("MST contains a cycle with edge (%d, %d)", edge.U, edge.V)
		}
	}

	// Check if all vertices are connected (single component)
	root := uf.Find(0)
	for i := 1; i < graph.Vertices; i++ {
		if uf.Find(i) != root {
			return fmt.Errorf("MST is not connected, vertex %d is not reachable from vertex 0", i)
		}
	}

	return nil
}

// CompareKruskalSorts compares the standard and custom sort implementations
func CompareKruskalSorts(graph *KruskalGraph) (float64, float64, bool) {
	mstStandard, standardWeight := KruskalMST(graph)
	mstCustom, customWeight := KruskalMSTWithCustomSort(graph)

	// Weights should be equal for both implementations
	equal := standardWeight == customWeight

	// Validate both MSTs
	if err := ValidateKruskalMST(graph, mstStandard); err != nil {
		fmt.Printf("Standard Kruskal MST validation failed: %v\n", err)
		equal = false
	}

	if err := ValidateKruskalMST(graph, mstCustom); err != nil {
		fmt.Printf("Custom sort Kruskal MST validation failed: %v\n", err)
		equal = false
	}

	return standardWeight, customWeight, equal
}

// ConvertBoruvkaToKruskal converts a Borůvka graph to Kruskal format
func ConvertBoruvkaToKruskal(boruvkaGraph *BoruvkaGraph) *KruskalGraph {
	kruskalGraph := &KruskalGraph{
		Vertices: boruvkaGraph.Vertices,
		Edges:    make([]KruskalEdge, len(boruvkaGraph.Edges)),
	}

	for i, edge := range boruvkaGraph.Edges {
		kruskalGraph.Edges[i] = KruskalEdge(edge)
	}

	return kruskalGraph
}

// ConvertKruskalToBoruvka converts a Kruskal graph to Borůvka format
func ConvertKruskalToBoruvka(kruskalGraph *KruskalGraph) *BoruvkaGraph {
	boruvkaGraph := &BoruvkaGraph{
		Vertices: kruskalGraph.Vertices,
		Edges:    make([]BoruvkaEdge, len(kruskalGraph.Edges)),
	}

	for i, edge := range kruskalGraph.Edges {
		boruvkaGraph.Edges[i] = BoruvkaEdge(edge)
	}

	return boruvkaGraph
}

// ExampleKruskalMST demonstrates the usage of Kruskal's MST algorithm
func ExampleKruskalMST() {
	graph := BuildKruskalSampleGraph()

	mst, totalWeight := KruskalMST(graph)

	fmt.Println("Kruskal's MST Results:")
	fmt.Printf("Total weight: %.2f\n", totalWeight)
	fmt.Println("MST edges:")
	for _, edge := range mst {
		fmt.Printf("  (%d, %d) weight: %.2f\n", edge.U, edge.V, edge.Weight)
	}

	// Validate the result
	if err := ValidateKruskalMST(graph, mst); err != nil {
		fmt.Printf("Validation error: %v\n", err)
	} else {
		fmt.Println("MST validation: PASSED")
	}

	// Compare with custom sort implementation
	fmt.Println("\nComparing standard vs custom sort:")
	standardWeight, customWeight, equal := CompareKruskalSorts(graph)
	fmt.Printf("Standard sort weight: %.2f\n", standardWeight)
	fmt.Printf("Custom sort weight: %.2f\n", customWeight)
	fmt.Printf("Results equal: %v\n", equal)
}
