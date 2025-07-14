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

// Package polylogarithmic contains implementations of algorithms with
// O((log n)^c) polylogarithmic time complexity.
//
// Polylogarithmic algorithms have complexity that grows as powers of logarithms.
// They are slower than logarithmic O(log n) but much faster than linear O(n)
// complexity, often appearing in advanced data structures and parallel algorithms.
//
// Examples for O((log n)²):
//   - n=2: (log₂(2))² = 1 operation
//   - n=16: (log₂(16))² = 16 operations
//   - n=1024: (log₂(1024))² = 100 operations
//   - n=65536: (log₂(65536))² = 256 operations
//   - n=1048576: (log₂(1048576))² = 400 operations
//
// Common use cases include range trees, segment trees with additional dimensions,
// fractional cascading, and some parallel algorithms that require multiple
// logarithmic factors for coordination.
//
// Benchmarks can handle large input sizes (n ≤ 10^6) efficiently due to
// the slow growth of logarithmic powers.
package polylogarithmic
