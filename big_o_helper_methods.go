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

package bigo

import (
	"math"
	"math/big"

	"github.com/rsned/bigmath"
)

// factorial is a function that returns the factorial of a given integer as a float.
func factorial(x int) float64 {
	if x <= 1 {
		return 1
	}

	return factorial(x-1) * float64(x)
}

// inverseAckermann computes an approximation of the inverse Ackermann function α(n).
// For practical purposes, this function returns values that grow extremely slowly:
// α(n) = 1 for n ≤ 2
// α(n) = 2 for 3 ≤ n ≤ 7
// α(n) = 3 for 8 ≤ n ≤ 2047
// α(n) = 4 for 2048 ≤ n ≤ 2^2047
// α(n) = 5 for larger n (but this is astronomically large)
func inverseAckermann(n int) float64 {
	switch {
	case n <= 2:
		return 1.0
	case n <= 7:
		return 2.0
	case n <= 2047:
		return 3.0
	case n < math.MaxInt64:
		return 4.0
	default:
		// TODO(rsned): Handle the case where n > 2^2047 which isn't reachable by int64

		return 5.0
	}
}

// inverseAckermannBig computes an approximation of the inverse Ackermann function α(n) for big.Float inputs.
// For practical purposes, this function returns values that grow extremely slowly:
// α(n) = 1 for n ≤ 2
// α(n) = 2 for 3 ≤ n ≤ 7
// α(n) = 3 for 8 ≤ n ≤ 2047
// α(n) = 4 for 2048 ≤ n ≤ 2^2047
// α(n) = 5 for larger n (but this is astronomically large)
func inverseAckermannBig(n *big.Float) *big.Float {
	switch {
	case n.Cmp(big.NewFloat(2)) <= 0:
		return big.NewFloat(1)
	case n.Cmp(big.NewFloat(7)) <= 0:
		return big.NewFloat(2)
	case n.Cmp(big.NewFloat(2047)) <= 0:
		return big.NewFloat(3)
	case n.Cmp(big.NewFloat(math.MaxInt64)) < 0:
		return big.NewFloat(4)
	default:
		return big.NewFloat(5)
	}
}

// logStarFloat computes the iterated logarithm function log*(n) for float64 inputs.
// The iterated logarithm is the number of times the natural logarithm function
// must be applied before the result is less than or equal to 1.
//
// log*(n) = 0 if n ≤ 1
// log*(n) = 1 + log*(log(n)) if n > 1
//
// Examples:
// log*(1) = 0
// log*(2) = 1 (since log(2) ≈ 0.69 ≤ 1)
// log*(16) = 3 (since log(16) = 4, log(4) ≈ 1.39, log(1.39) ≈ 0.33 ≤ 1)
// log*(65536) = 4
// log*(2^65536) = 5
//
// For practical purposes, log*(n) ≤ 5 for any realistically computable n.
// This implementation includes a safety limit to prevent infinite loops.
func logStarFloat(x float64) int {
	// Handle special cases first
	if math.IsInf(x, 1) {
		return 1 // log(+Inf) = +Inf, but we return 1 for consistency
	}

	if math.IsInf(x, -1) || math.IsNaN(x) {
		return 0
	}

	if x <= 1 {
		return 0
	}

	count := 0
	maxIterations := 10 // Safety limit - log*(n) is ≤ 5 for any practical n

	for x > 1 && count < maxIterations {
		x = math.Log(x)
		count++

		// Additional safety check for NaN or Inf results
		if math.IsNaN(x) || math.IsInf(x, -1) {
			break
		}
	}

	return count
}

// logStarBig computes the iterated logarithm function log*(n) for big.Float inputs.
// The iterated logarithm is the number of times the natural logarithm function
// must be applied before the result is less than or equal to 1.
//
// log*(n) = 0 if n ≤ 1
// log*(n) = 1 + log*(log(n)) if n > 1
//
// This version uses big.Float for high-precision arithmetic and can handle
// very large numbers that would overflow float64.
//
// For practical purposes, log*(n) ≤ 5 for any realistically computable n.
// This implementation includes a safety limit to prevent infinite loops.
func logStarBig(x *big.Float) int {
	if x == nil {
		return 0
	}

	// Handle special cases first
	if x.IsInf() {
		if x.Sign() > 0 {
			return 1 // log(+Inf) = +Inf, but we return 1 for consistency
		}

		return -1 // Use -1 to indicate NaN case for negative infinity
	}

	one := big.NewFloat(1.0)
	if x.Cmp(one) <= 0 {
		return 0
	}

	count := 0
	maxIterations := 10 // Safety limit - log*(n) is ≤ 5 for any practical n

	// Make a copy to avoid modifying the original
	current := new(big.Float).Copy(x)
	epsilon := big.NewFloat(1e-9) // A small tolerance

	for current.Cmp(one) > 0 && count < maxIterations {
		current = bigmath.Log(current)
		count++

		// If current is very close to 1, treat it as 1 to avoid precision issues
		diff := new(big.Float).Sub(current, one)
		if diff.Abs(diff).Cmp(epsilon) < 0 {
			current.Set(one)
		}

		// Additional safety check for infinity results
		if current.IsInf() {
			break
		}
	}

	return count
}
