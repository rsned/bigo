// Copyright 2025 Robert Snedegar
//
//   Licensed under the Apache License, Version 2.0 (the License);
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an AS IS BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

/*
Package bigo is a library to generate Big O estimates for real world timings of
algorithms and code. Given a collection of counts and timings collected
when running code or benchmarks, the library will attempt to characterize
which Big O notations most likely match the values.

The library has support for the following Big O classifications:

  - O(1) - Constant Time
  - O(log (log n)) - Double Log Time
  - O(n (log n)) - Log Time
  - O((log n)^c) - Polylogarithmic Time
  - O(n) - Linear Time
  - O(n log* n) - N Log* Time
  - O(n log n) - Linearithmic Time
  - O(n^2) - Quadratic Time
  - O(n^3) - Cubic Time
  - O(n^c) - Polynomial Time
  - O(2^n) - Exponential Time
  - O(n!) - Factorial Time
  - O(n^n) - Hyper-Exponential Time

The characterization can be applied to anyset of values whether they are times,
bytes, allocations, or any other type.

Underneath the examples directory are subdirectories for each Big O category
that contain a variety of example methods and algorithms that are generally
considered to fall in a given level along with basic benchmarks that can
be used to generate real timing results that can be used.
*/
package bigo
