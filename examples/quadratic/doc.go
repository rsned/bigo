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

// Package quadratic contains implementations of algorithms with
// O(n²) quadratic time complexity.
//
// Quadratic algorithms have nested loops with two dimensions of computation.
// They are significantly slower than linearithmic O(n log n) but much faster
// than cubic O(n³) complexity. Growth accelerates rapidly with input size.
//
// Examples:
//   - n=1: 1² = 1 operation
//   - n=10: 10² = 100 operations
//   - n=100: 100² = 10,000 operations
//   - n=1000: 1000² = 1,000,000 operations
//   - n=10000: 10000² = 100,000,000 operations
//
// Common use cases include simple sorting algorithms (bubble sort, insertion
// sort, selection sort), matrix operations, all-pairs comparisons, and
// nested loop algorithms that examine every pair of elements.
//
// Benchmarks should limit input sizes to n ≤ 1000 to avoid excessive
// execution times during performance testing.
package quadratic
