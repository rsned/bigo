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

// Package linear contains implementations of algorithms with
// O(n) linear time complexity.
//
// Linear algorithms traverse data once with each element examined exactly once.
// They grow proportionally with input size, slower than logarithmic O(log n)
// but faster than linearithmic O(n log n) complexity.
//
// Examples:
//   - n=1: 1 operation
//   - n=10: 10 operations
//   - n=100: 100 operations
//   - n=1000: 1,000 operations
//   - n=10000: 10,000 operations
//
// Common use cases include array/slice traversal, finding minimum/maximum
// values, linear search, counting elements, and single-pass algorithms
// that process each element exactly once.
//
// Benchmarks can efficiently test moderately large input sizes (n ≤ 10^5)
// before execution time becomes a significant constraint.
package linear
