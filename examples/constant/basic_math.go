// Copyright 2025 Robert Snedegar
//
// Licensed under the Apache License, Version 2.0 (the License);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an AS IS BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package constant

// Add two integers.
// This operation demonstrates O(1) constant time complexity because
// integer addition is a single CPU instruction regardless of the values.
func Add(a, b int) int {
	// Direct addition - single CPU instruction (ADD)
	return a + b
}

// Subtract two integers.
// This operation demonstrates O(1) constant time complexity because
// integer subtraction is a single CPU instruction regardless of the values.
func Subtract(a, b int) int {
	// Direct subtraction - single CPU instruction (SUB)
	return a - b
}

// Multiply two integers.
// This operation demonstrates O(1) constant time complexity because
// integer multiplication is a single CPU instruction regardless of the values.
func Multiply(a, b int) int {
	// Direct multiplication - single CPU instruction (MUL)
	return a * b
}

// Divide two integers.
// This operation demonstrates O(1) constant time complexity because
// integer division is a single CPU instruction regardless of the values.
func Divide(a, b int) int {
	if b == 0 {
		panic("division by zero")
	}
	// Direct division - single CPU instruction (DIV)
	return a / b
}
