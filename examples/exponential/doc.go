// Copyright 2025 Robert Snedegar
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package exponential contains implementations of algorithms with
// O(2ⁿ) exponential time complexity.
//
// Exponential algorithms generate or explore all possible combinations or subsets.
// They represent a major complexity boundary, being much slower than polynomial
// O(nᵏ) but faster than factorial O(n!) for moderately small values of n.
//
// Examples:
//   - n=1: 2¹ = 2 operations
//   - n=5: 2⁵ = 32 operations
//   - n=10: 2¹⁰ = 1,024 operations
//   - n=15: 2¹⁵ = 32,768 operations
//   - n=20: 2²⁰ = 1,048,576 operations
//
// Common use cases include recursive Fibonacci, Tower of Hanoi, subset/powerset
// generation, many NP-complete problems solved by brute force, and backtracking
// algorithms that explore all possible solution spaces.
//
// Benchmarks must use very small input sizes (n ≤ 20) due to the explosive
// growth rate that makes larger inputs computationally infeasible.
package exponential
