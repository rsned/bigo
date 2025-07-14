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

// Package hyperexponential contains implementations of algorithms with
// O(n^n) hyperexponential time complexity.
//
// Hyperexponential algorithms have complexity that grows as n raised to the
// power of n, which is even worse than factorial O(n!) for large n.
// These algorithms are extremely limited in practical application due to
// their explosive growth rate.
//
// Examples:
//   - n=1: 1^1 = 1 operation
//   - n=2: 2^2 = 4 operations
//   - n=3: 3^3 = 27 operations
//   - n=4: 4^4 = 256 operations
//   - n=5: 5^5 = 3,125 operations
//   - n=6: 6^6 = 46,656 operations
//
// Due to the extreme growth rate, benchmarks are typically limited to n ≤ 5.
package hyperexponential
