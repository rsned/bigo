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

// Package nlogstar contains implementations of algorithms with
// O(n log*(n)) complexity, where log*(n) is the iterated logarithm.
//
// These algorithms use the iterated logarithm log*(n), which grows extremely
// slowly. They are slower than linear O(n) but much faster than linearithmic
// O(n log n), representing near-linear performance for practical purposes.
//
// Examples (log*(n) is the number of times log must be applied to reach ≤ 1):
//   - n=2: log*(2) = 1, so 2×1 = 2 operations
//   - n=16: log*(16) = 3, so 16×3 = 48 operations
//   - n=65536: log*(65536) = 4, so 65536×4 = 262,144 operations
//   - n=2⁶⁵⁵³⁶: log*(2⁶⁵⁵³⁶) = 5, so 2⁶⁵⁵³⁶×5 operations
//   - For all practical n: log*(n) ≤ 5
//
// Common use cases include union-find data structures with path compression
// and link-by-rank, which achieve amortized O(n log*(n)) for sequences of
// union and find operations.
//
// Benchmarks can handle very large input sizes (n ≤ 10^6) efficiently since
// log*(n) is effectively constant for all practical values.
package nlogstar
