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

// Package logarithmic contains implementations of algorithms with
// O(log n) logarithmic time complexity.
//
// Logarithmic algorithms repeatedly divide the problem space in half or
// use tree-like structures. They are highly efficient, slower than constant
// O(1) but much faster than linear O(n) time complexity.
//
// Examples:
//   - n=1: log₂(1) = 0 operations (Not really 0 in reality, some work actually has to happen even at n=1)
//   - n=8: log₂(8) = 3 operations
//   - n=1024: log₂(1024) = 10 operations
//   - n=1048576: log₂(1048576) = 20 operations
//   - n=1073741824: log₂(1073741824) = 30 operations
//
// Common use cases include binary search in sorted arrays, binary tree
// traversal, heap operations (insert/delete), and divide-and-conquer
// algorithms that split the problem space in half with each iteration.
//
// Benchmarks can test large input sizes (n ≤ 10^6) efficiently due to
// the slow growth rate of logarithmic functions.
package logarithmic
