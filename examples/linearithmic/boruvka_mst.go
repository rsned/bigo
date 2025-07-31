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
)

// BoruvkaEdge represents a weighted edge in Borůvka's algorithm
type BoruvkaEdge struct {
	U, V   int     // vertices
	Weight float64 // edge weight
}

// BoruvkaGraph represents a weighted undirected graph for Borůvka's algorithm
type BoruvkaGraph struct {
	Vertices int
	Edges    []BoruvkaEdge
}

// BoruvkaUnionFind implements union-find with path compression and union by rank
type BoruvkaUnionFind struct {
	parent []int
	rank   []int
	count  int // number of components
}

// NewBoruvkaUnionFind creates a new union-find structure
func NewBoruvkaUnionFind(n int) *BoruvkaUnionFind {
	uf := &BoruvkaUnionFind{
		parent: make([]int, n),
		rank:   make([]int, n),
		count:  n,
	}

	for i := 0; i < n; i++ {
		uf.parent[i] = i
	}

	return uf
}

// Find returns the root of the set containing x with path compression
func (uf *BoruvkaUnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x])
	}

	return uf.parent[x]
}

// Union merges the sets containing x and y using union by rank
func (uf *BoruvkaUnionFind) Union(x, y int) bool {
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

	uf.count--

	return true
}

// ComponentCount returns the number of connected components
func (uf *BoruvkaUnionFind) ComponentCount() int {
	return uf.count
}

// BoruvkaMST implements Borůvka's minimum spanning tree algorithm
// Time complexity: O(m log n) where m is edges and n is vertices
// Space complexity: O(n + m)
func BoruvkaMST(graph *BoruvkaGraph) ([]BoruvkaEdge, float64) {
	if graph.Vertices <= 1 {
		return []BoruvkaEdge{}, 0
	}

	uf := NewBoruvkaUnionFind(graph.Vertices)
	mst := make([]BoruvkaEdge, 0, graph.Vertices-1)
	totalWeight := 0.0

	// Continue until we have a single component (spanning tree)
	for uf.ComponentCount() > 1 {
		// Find cheapest edge for each component
		cheapest := make([]int, graph.Vertices)
		for i := range cheapest {
			cheapest[i] = -1
		}

		// For each edge, check if it's the cheapest for its components
		for edgeIdx, edge := range graph.Edges {
			compU := uf.Find(edge.U)
			compV := uf.Find(edge.V)

			// Skip if both vertices are in the same component
			if compU == compV {
				continue
			}

			// Check if this is the cheapest edge for component U
			if cheapest[compU] == -1 || edge.Weight < graph.Edges[cheapest[compU]].Weight {
				cheapest[compU] = edgeIdx
			}

			// Check if this is the cheapest edge for component V
			if cheapest[compV] == -1 || edge.Weight < graph.Edges[cheapest[compV]].Weight {
				cheapest[compV] = edgeIdx
			}
		}

		// Add all cheapest edges to MST (avoiding duplicates)
		addedEdges := make(map[int]bool)

		for i := 0; i < graph.Vertices; i++ {
			if cheapest[i] != -1 && !addedEdges[cheapest[i]] {
				edge := graph.Edges[cheapest[i]]

				// Try to add this edge to MST
				if uf.Union(edge.U, edge.V) {
					mst = append(mst, edge)
					totalWeight += edge.Weight
					addedEdges[cheapest[i]] = true
				}
			}
		}

		// Safety check to prevent infinite loops
		if len(addedEdges) == 0 {
			break
		}
	}

	return mst, totalWeight
}

// BuildBoruvkaSampleGraph creates a sample graph for testing Borůvka's algorithm
func BuildBoruvkaSampleGraph() *BoruvkaGraph {
	return &BoruvkaGraph{
		Vertices: 7,
		Edges: []BoruvkaEdge{
			{0, 1, 7.0},
			{0, 3, 5.0},
			{1, 2, 8.0},
			{1, 3, 9.0},
			{1, 4, 7.0},
			{2, 4, 5.0},
			{3, 4, 15.0},
			{3, 5, 6.0},
			{4, 5, 8.0},
			{4, 6, 9.0},
			{5, 6, 11.0},
		},
	}
}

// BuildCompleteGraph creates a complete graph with n vertices for benchmarking
func BuildCompleteGraph(n int) *BoruvkaGraph {
	graph := &BoruvkaGraph{
		Vertices: n,
		Edges:    make([]BoruvkaEdge, 0, n*(n-1)/2),
	}

	// Add all possible edges with random-like weights
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			// Use a deterministic but varied weight calculation
			weight := float64((i*n+j)%100 + 1)
			graph.Edges = append(graph.Edges, BoruvkaEdge{
				U: i, V: j, Weight: weight,
			})
		}
	}

	return graph
}

// ValidateMST checks if the given edges form a valid minimum spanning tree
func ValidateMST(graph *BoruvkaGraph, mst []BoruvkaEdge) error {
	if len(mst) != graph.Vertices-1 {
		return fmt.Errorf("MST should have %d edges, got %d", graph.Vertices-1, len(mst))
	}

	// Check connectivity using union-find
	uf := NewBoruvkaUnionFind(graph.Vertices)

	for _, edge := range mst {
		if edge.U < 0 || edge.U >= graph.Vertices ||
			edge.V < 0 || edge.V >= graph.Vertices {
			return fmt.Errorf("invalid edge vertices: (%d, %d)", edge.U, edge.V)
		}

		if !uf.Union(edge.U, edge.V) {
			return fmt.Errorf("MST contains a cycle with edge (%d, %d)", edge.U, edge.V)
		}
	}

	if uf.ComponentCount() != 1 {
		return fmt.Errorf("MST is not connected, has %d components", uf.ComponentCount())
	}

	return nil
}

// CompareMSTAlgorithms compares Borůvka's algorithm with Kruskal's for verification
func CompareMSTAlgorithms(graph *BoruvkaGraph) (float64, float64, bool) {
	// Run Borůvka's algorithm
	boruvkaMST, boruvkaWeight := BoruvkaMST(graph)

	// Convert to Kruskal format and run Kruskal's algorithm
	kruskalGraph := ConvertBoruvkaToKruskal(graph)
	kruskalMSTResult, kruskalWeight := KruskalMST(kruskalGraph)

	// Convert Kruskal result back to Borůvka format for validation
	kruskalMST := make([]BoruvkaEdge, len(kruskalMSTResult))
	for i, edge := range kruskalMSTResult {
		kruskalMST[i] = BoruvkaEdge(edge)
	}

	// Check if weights are equal (allowing for floating point precision)
	equal := math.Abs(boruvkaWeight-kruskalWeight) < 1e-9

	// Validate both MSTs
	if err := ValidateMST(graph, boruvkaMST); err != nil {
		fmt.Printf("Borůvka MST validation failed: %v\n", err)
		equal = false
	}

	if err := ValidateMST(graph, kruskalMST); err != nil {
		fmt.Printf("Kruskal MST validation failed: %v\n", err)
		equal = false
	}

	return boruvkaWeight, kruskalWeight, equal
}

// ExampleBoruvkaMST demonstrates the usage of Borůvka's MST algorithm
func ExampleBoruvkaMST() {
	graph := BuildBoruvkaSampleGraph()

	mst, totalWeight := BoruvkaMST(graph)

	fmt.Println("Borůvka's MST Results:")
	fmt.Printf("Total weight: %.2f\n", totalWeight)
	fmt.Println("MST edges:")
	for _, edge := range mst {
		fmt.Printf("  (%d, %d) weight: %.2f\n", edge.U, edge.V, edge.Weight)
	}

	// Validate the result
	if err := ValidateMST(graph, mst); err != nil {
		fmt.Printf("Validation error: %v\n", err)
	} else {
		fmt.Println("MST validation: PASSED")
	}
}
