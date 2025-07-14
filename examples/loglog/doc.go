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

// Package loglog contains implementations of algorithms with
// O(log log n) double logarithmic time complexity.
//
// Double logarithmic algorithms achieve sub-logarithmic performance through
// advanced data structures. They are faster than logarithmic O(log n) but
// slower than constant O(1), representing extremely efficient complexity.
//
// Examples:
//   - n=2: log₂(log₂(2)) = 0 operations  (Essentially this bottoms out at 1 in reality)
//   - n=16: log₂(log₂(16)) = 2 operations
//   - n=65536: log₂(log₂(65536)) = 4 operations
//   - n=4294967296: log₂(log₂(4294967296)) = 5 operations
//   - n=2⁶⁴: log₂(log₂(2⁶⁴)) = 6 operations
//
// Common use cases include van Emde Boas trees for successor/predecessor queries,
// interpolation search on uniformly distributed data, fusion trees for integer
// sorting, and Y-fast tries for dynamic predecessor problems.
//
// Benchmarks can efficiently test very large input sizes (n ≤ 10^9) due to
// the extremely slow growth of the double logarithm function.
package loglog
