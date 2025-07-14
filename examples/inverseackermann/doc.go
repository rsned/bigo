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

// Package inverseackermann contains implementations of algorithms with
// O(α(n)) inverse Ackermann time complexity.
//
// Inverse Ackermann algorithms achieve near-constant performance through
// sophisticated data structures. They are slower than constant O(1) but
// faster than double logarithmic O(log log n), representing extremely
// efficient amortized complexity.
//
// Examples (α(n) ≤ 4 for all practical values):
//   - n=1: α(1) = 1 operation
//   - n=1000: α(1000) ≈ 3 operations
//   - n=10^6: α(10^6) ≈ 4 operations
//   - n=10^9: α(10^9) ≈ 4 operations
//   - n=A(4,2): α(A(4,2)) = 4 operations (where A is Ackermann function)
//
// Common use cases include disjoint set union (union-find) with path compression
// and union by rank, Tarjan's offline lowest common ancestor algorithm, and
// Chazelle's minimum spanning tree algorithm.
//
// Benchmarks can handle very large input sizes (n ≤ 10^9) efficiently since
// α(n) is effectively constant (≤ 4) for all realistic inputs.
package inverseackermann
