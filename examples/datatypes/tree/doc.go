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

// Package tree provides advanced tree data structure implementations for
// efficient range queries, updates, and specialized operations.
//
// This package contains several tree-based data structures that demonstrate
// different complexity classes and optimization techniques:
//
// # Binary Search Tree (BST)
//
// BSTNode implements a standard binary search tree with O(log n) average-case
// operations for search, insertion, and deletion. Includes balanced construction
// from sorted arrays and inorder traversal.
//
//	// Create a balanced BST from sorted values
//	values := []int{1, 2, 3, 4, 5, 6, 7}
//	root := BuildBST(values)
//
//	// Insert new values
//	root = root.InsertBST(8)
//
//	// Get sorted values via inorder traversal
//	sorted := root.InorderTraversal()
//
// # Use Cases
//
// These data structures are commonly used in:
//
//   - Computational geometry (range trees)
//   - Dynamic programming optimizations (segment trees)
//   - Competitive programming (Fenwick trees)
//   - Database indexing (B-trees, range trees)
//   - Real-time systems requiring fast queries (fusion trees)
//
// All implementations include comprehensive benchmarks demonstrating their
// time complexity characteristics for performance analysis and comparison.
package tree
