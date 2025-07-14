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

// Package cubic contains implementations of algorithms with
// O(n³) cubic time complexity.
//
// Cubic algorithms have three nested loops or require three-dimensional
// computations. They are much slower than quadratic O(n²) but significantly
// faster than exponential O(2ⁿ) complexity. Growth becomes very steep.
//
// Examples:
//   - n=1: 1³ = 1 operation
//   - n=10: 10³ = 1,000 operations
//   - n=50: 50³ = 125,000 operations
//   - n=100: 100³ = 1,000,000 operations
//   - n=200: 200³ = 8,000,000 operations
//
// Common use cases include matrix chain multiplication, Floyd-Warshall
// all-pairs shortest paths, triple nested loops for finding combinations,
// and some dynamic programming solutions with three dimensions.
//
// Benchmarks should be limited to small input sizes (n ≤ 100) to prevent
// excessively long execution times during testing.
package cubic
