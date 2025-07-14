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

// Package factorial contains implementations of algorithms with
// O(n!) factorial time complexity.
//
// Factorial algorithms generate or examine all possible permutations of n elements.
// They are extremely slow, worse than exponential O(2ⁿ) for n ≥ 4, but better
// than hyperexponential O(n^n) for small values of n.
//
// Examples:
//   - n=1: 1! = 1 operation
//   - n=3: 3! = 6 operations
//   - n=5: 5! = 120 operations
//   - n=7: 7! = 5,040 operations
//   - n=8: 8! = 40,320 operations
//   - n=10: 10! = 3,628,800 operations
//
// Common use cases include permutation generation, traveling salesman brute force,
// assignment problems requiring all possible arrangements, and optimization
// problems that must examine every possible ordering.
//
// Benchmarks are extremely limited (n ≤ 8) due to factorial growth becoming
// computationally prohibitive very quickly.
package factorial
