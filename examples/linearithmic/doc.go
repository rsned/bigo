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

// Package linearithmic contains implementations of algorithms with
// O(n log n) linearithmic time complexity.
//
// Linearithmic algorithms combine linear iteration with logarithmic operations
// per element. They are slower than linear O(n) but much faster than quadratic
// O(n²) complexity, representing the optimal time for comparison-based sorting.
//
// Examples:
//   - n=1: 1×log₂(1) = 0 operations
//   - n=8: 8×log₂(8) = 24 operations
//   - n=64: 64×log₂(64) = 384 operations
//   - n=1024: 1024×log₂(1024) = 10,240 operations
//   - n=65536: 65536×log₂(65536) = 1,048,576 operations
//
// Common use cases include efficient sorting algorithms (merge sort, heap sort,
// quick sort average case), building heaps from arrays, divide-and-conquer
// algorithms, and many graph algorithms.
//
// Benchmarks can test moderate input sizes (n ≤ 10^4) efficiently before
// performance constraints become significant.
package linearithmic
