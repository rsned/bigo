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

// Package polynomial contains implementations of algorithms with
// polynomial time complexity O(nᵏ) where k > 3.
//
// Polynomial algorithms solve complex optimization problems with higher-degree
// polynomial complexity. They are slower than cubic O(n³) but remain in
// polynomial time, making them much faster than exponential O(2ⁿ) complexity.
//
// Examples for O(n⁴):
//   - n=1: 1⁴ = 1 operation
//   - n=5: 5⁴ = 625 operations
//   - n=10: 10⁴ = 10,000 operations
//   - n=20: 20⁴ = 160,000 operations
//   - n=50: 50⁴ = 6,250,000 operations
//
// Common use cases include complex dynamic programming solutions, some graph
// algorithms (maximum flow, min-cost flow), matrix operations with higher
// dimensions, and optimization problems with polynomial-time solutions.
//
// Benchmarks should use very small input sizes (n ≤ 50) due to the
// rapidly increasing execution time with higher polynomial degrees.
package polynomial
