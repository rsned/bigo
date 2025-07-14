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
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestLogStarFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected int
	}{
		// Base cases
		{
			name:     "zero",
			input:    0,
			expected: 0,
		},
		{
			name:     "one",
			input:    1,
			expected: 0,
		},
		{
			name:     "negative",
			input:    -1,
			expected: 0,
		},
		// Valid positive cases
		{
			name:     "two",
			input:    2,
			expected: 1,
		},
		{
			name:     "e (natural log base)",
			input:    math.E,
			expected: 1,
		},
		{
			name:     "four",
			input:    4,
			expected: 2,
		},
		{
			name:     "sixteen",
			input:    16,
			expected: 3,
		},
		{
			name:     "256",
			input:    256,
			expected: 3,
		},
		{
			name:     "65536 (2^16)",
			input:    65536,
			expected: 3,
		},
		// Edge cases with special float values
		{
			name:     "very small positive",
			input:    1e-10,
			expected: 0,
		},
		{
			name:     "just above one",
			input:    1.0001,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := logStarFloat(tt.input)
			if result != tt.expected {
				t.Errorf("logStar(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLogStarFloatSpecialFloatValues(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected int
	}{
		{
			name:     "positive infinity",
			input:    math.Inf(1),
			expected: 1, // log(+Inf) = +Inf, so it should iterate once and return 1
		},
		{
			name:     "negative infinity",
			input:    math.Inf(-1),
			expected: 0, // returns 0 for negative infinity
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := logStarFloat(tt.input)
			if result != tt.expected {
				t.Errorf("logStarFloat(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLogStarBig(t *testing.T) {
	tests := []struct {
		name     string
		input    *big.Float
		expected int
	}{
		// Base cases
		{
			name:     "nil input",
			input:    nil,
			expected: 0,
		},
		{
			name:     "zero",
			input:    big.NewFloat(0),
			expected: 0,
		},
		{
			name:     "one",
			input:    big.NewFloat(1),
			expected: 0,
		},
		{
			name:     "negative",
			input:    big.NewFloat(-1),
			expected: 0,
		},
		// Valid positive cases
		{
			name:     "two",
			input:    big.NewFloat(2),
			expected: 1,
		},
		{
			name:     "e (natural log base)",
			input:    big.NewFloat(math.E),
			expected: 1,
		},
		{
			name:     "four",
			input:    big.NewFloat(4),
			expected: 2,
		},
		{
			name:     "sixteen",
			input:    big.NewFloat(16),
			expected: 3,
		},
		{
			name:     "256",
			input:    big.NewFloat(256),
			expected: 3,
		},
		{
			name:     "65536 (2^16)",
			input:    big.NewFloat(65536),
			expected: 3,
		},
		// Very large number that would overflow float64
		{
			name: "very large number",
			input: func() *big.Float {
				f := new(big.Float).SetPrec(256)
				f.SetString("1e308")

				return f
			}(),
			expected: 4,
		},
		// Edge cases
		{
			name:     "very small positive",
			input:    big.NewFloat(1e-10),
			expected: 0,
		},
		{
			name:     "just above one",
			input:    big.NewFloat(1.0001),
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := logStarBig(tt.input)
			if result != tt.expected {
				t.Errorf("logStarBig(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLogStarBigSpecialValues(t *testing.T) {
	tests := []struct {
		name     string
		input    *big.Float
		expected int
	}{
		{
			name:     "positive infinity",
			input:    big.NewFloat(0).SetInf(false),
			expected: 1,
		},
		{
			name:     "negative infinity",
			input:    big.NewFloat(0).SetInf(true),
			expected: -1, // Special return value for NaN case
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := logStarBig(tt.input)
			if result != tt.expected {
				t.Errorf("logStarBig(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLogStarBigConsistencyWithFloat(t *testing.T) {
	// Test that logStarBig produces consistent results with logStarFloat
	// for values that can be represented in both float64 and big.Float
	testValues := []float64{0, 1, 2, math.E, 4, 16, 256, 65536}

	for _, val := range testValues {
		t.Run(fmt.Sprintf("consistency_test_%g", val), func(t *testing.T) {
			floatResult := logStarFloat(val)
			bigResult := logStarBig(big.NewFloat(val))

			if bigResult != floatResult {
				t.Errorf("logStarBig(%g) = %d, but logStarFloat(%g) = %d",
					val, bigResult, val, floatResult)
			}
		})
	}
}

func TestFactorial(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected float64
	}{
		{
			name:     "factorial of 0",
			input:    0,
			expected: 1,
		},
		{
			name:     "factorial of 1",
			input:    1,
			expected: 1,
		},
		{
			name:     "factorial of 2",
			input:    2,
			expected: 2,
		},
		{
			name:     "factorial of 3",
			input:    3,
			expected: 6,
		},
		{
			name:     "factorial of 4",
			input:    4,
			expected: 24,
		},
		{
			name:     "factorial of 5",
			input:    5,
			expected: 120,
		},
		{
			name:     "factorial of 10",
			input:    10,
			expected: 3628800,
		},
		{
			name:     "factorial of negative number",
			input:    -1,
			expected: 1,
		},
		{
			name:     "factorial of negative number",
			input:    -5,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := factorial(tt.input)
			if result != tt.expected {
				t.Errorf("factorial(%d) = %g, want %g", tt.input, result, tt.expected)
			}
		})
	}
}

func TestInverseAckermann(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected float64
	}{
		{
			name:     "n = 0",
			input:    0,
			expected: 1.0,
		},
		{
			name:     "n = 1",
			input:    1,
			expected: 1.0,
		},
		{
			name:     "n = 2",
			input:    2,
			expected: 1.0,
		},
		{
			name:     "n = 3",
			input:    3,
			expected: 2.0,
		},
		{
			name:     "n = 7",
			input:    7,
			expected: 2.0,
		},
		{
			name:     "n = 8",
			input:    8,
			expected: 3.0,
		},
		{
			name:     "n = 100",
			input:    100,
			expected: 3.0,
		},
		{
			name:     "n = 2047",
			input:    2047,
			expected: 3.0,
		},
		{
			name:     "n = 2048",
			input:    2048,
			expected: 4.0,
		},
		{
			name:     "n = 100000",
			input:    100000,
			expected: 4.0,
		},
		{
			name:     "n = math.MaxInt64-1",
			input:    math.MaxInt64 - 1,
			expected: 4.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := inverseAckermann(tt.input)
			if result != tt.expected {
				t.Errorf("inverseAckermann(%d) = %g, want %g", tt.input, result, tt.expected)
			}
		})
	}
}

func TestInverseAckermannBig(t *testing.T) {
	tests := []struct {
		name     string
		input    *big.Float
		expected float64
	}{
		{
			name:     "n = 0",
			input:    big.NewFloat(0),
			expected: 1.0,
		},
		{
			name:     "n = 1",
			input:    big.NewFloat(1),
			expected: 1.0,
		},
		{
			name:     "n = 2",
			input:    big.NewFloat(2),
			expected: 1.0,
		},
		{
			name:     "n = 3",
			input:    big.NewFloat(3),
			expected: 2.0,
		},
		{
			name:     "n = 8",
			input:    big.NewFloat(8),
			expected: 3.0,
		},
		{
			name:     "n = 100",
			input:    big.NewFloat(100),
			expected: 3.0,
		},
		{
			name:     "n = math.MaxInt64",
			input:    new(big.Float).SetInt64(math.MaxInt64),
			expected: 4.0,
		},
		{
			name:     "n = math.MaxInt64+1",
			input:    new(big.Float).Add(new(big.Float).SetInt64(math.MaxInt64), big.NewFloat(1)),
			expected: 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := inverseAckermannBig(tt.input)
			resultFloat, _ := result.Float64()
			if resultFloat != tt.expected {
				t.Errorf("inverseAckermannBig(%v) = %g, want %g", tt.input, resultFloat, tt.expected)
			}
		})
	}
}

func TestInverseAckermannConsistency(t *testing.T) {
	// Test that inverseAckermannBig produces consistent results with
	// inverseAckermann for values that can be represented in both
	// int and big.Float
	testValues := []int{0, 1, 2, 3, 4, 7, 8, 100, 65536, math.MaxInt64}

	for _, val := range testValues {
		t.Run(fmt.Sprintf("consistency_test_%d", val), func(t *testing.T) {
			intResult := inverseAckermann(val)
			bigResult := inverseAckermannBig(big.NewFloat(float64(val)))
			resultFloat, _ := bigResult.Float64()

			if resultFloat != intResult {
				t.Errorf("inverseAckermannBig(%d) = %g, but inverseAckermann(%d) = %g",
					val, resultFloat, val, intResult)
			}
		})
	}
}
