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

// Package constant contains implementations of algorithms with
// O(1) constant time complexity.
//
// Constant time algorithms perform the same number of operations regardless
// of input size. This is the most efficient complexity class, better than
// logarithmic O(log n) complexity. Operations execute in fixed time
// regardless of data size.
//
// Examples:
//   - n=1: 1 operation
//   - n=1000: 1 operation
//   - n=1000000: 1 operation
//   - n=1000000000: 1 operation
//
// Common use cases include array/slice indexing, hash table lookups,
// stack push/pop, queue enqueue/dequeue, and basic arithmetic operations.
//
// Benchmarks can safely test very large input sizes (n ≤ 10^9) since
// performance remains constant regardless of input magnitude.
package constant
