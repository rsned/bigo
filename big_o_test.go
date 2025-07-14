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
	"math/big"
	"testing"
)

func TestOrderUnsortedData(t *testing.T) {
	tests := []struct {
		name     string
		order    *BigO
		ns       []int
		vals     []float64
		expected string  // Expected complexity label
		minScore float64 // Minimum expected correlation score
	}{
		{
			name:     "linear time unsorted",
			order:    Linear,
			ns:       []int{500, 100, 300, 200, 400},              // Unsorted input sizes
			vals:     []float64{250.0, 50.0, 150.0, 100.0, 200.0}, // Corresponding unsorted timings
			expected: "O(n)",
			minScore: 0.8, // Should have high correlation for linear data
		},
		{
			name:     "quadratic time unsorted",
			order:    Quadratic,
			ns:       []int{20, 5, 15, 10},                 // Unsorted input sizes
			vals:     []float64{400.0, 25.0, 225.0, 100.0}, // n^2 pattern: 20²=400, 5²=25, 15²=225, 10²=100
			expected: "O(n^2)",
			minScore: 0.9, // Should have very high correlation for perfect quadratic data
		},
		{
			name:     "logarithmic time unsorted",
			order:    Log,
			ns:       []int{1000, 100, 10000, 10},   // Unsorted powers of 10
			vals:     []float64{3.0, 2.0, 4.0, 1.0}, // Approximate log base 10 pattern
			expected: "O(log n)",
			minScore: 0.7, // Reasonable correlation for log pattern
		},
		{
			name:     "constant time unsorted - low variance",
			order:    Constant,
			ns:       []int{1000, 100, 500, 200, 300},       // Unsorted input sizes
			vals:     []float64{10.1, 9.9, 10.0, 10.2, 9.8}, // Very consistent timings (~10 with small variance)
			expected: "O(1)",
			minScore: 0.9, // Should score high for constant behavior
		},
		{
			name:     "constant time unsorted - higher variance",
			order:    Constant,
			ns:       []int{50, 200, 100, 300, 150},         // Unsorted input sizes
			vals:     []float64{12.0, 8.0, 10.0, 11.0, 9.0}, // More variance but still relatively constant
			expected: "O(1)",
			minScore: 0.5, // Should score lower due to higher variance
		},
		{
			name:     "linearithmic time unsorted",
			order:    Linearithmic,
			ns:       []int{100, 10, 1000},           // Unsorted input sizes
			vals:     []float64{200.0, 10.0, 3000.0}, // Approximate n*log(n) pattern
			expected: "O(n log n)",
			minScore: 0.7, // Reasonable correlation for n log n pattern
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rating, err := tt.order.Rate(tt.ns, tt.vals)
			if err != nil {
				t.Fatalf("Rate() returned error: %v", err)
			}

			if rating.bigO.label != tt.expected {
				t.Errorf("Expected complexity %s, got %s", tt.expected, rating.bigO.label)
			}

			if rating.score < tt.minScore {
				t.Errorf("Expected score >= %f, got %f", tt.minScore, rating.score)
			}

			t.Logf("Complexity: %s, Score: %f", rating.bigO.label, rating.score)
		})
	}
}

// TestOrderBigUnsortedData tests the BigO.RankBig function with unsorted data.
// Break the assumption that the data we receive is all clean and proper. Odds
// are that many runs in a row may be added, and the smallest 'N' result may not
// be the first one.
//
// For the purpose of this test we pass in exact predicted scores so the correlation
// is simpler to compare.
func TestOrderBigUnsortedData(t *testing.T) {
	tests := []struct {
		name     string
		order    *BigO
		ns       []int
		vals     []*big.Float
		expected string  // Expected complexity label
		minScore float64 // Minimum expected correlation score
	}{
		{
			name:  "linear time unsorted big",
			order: Linear,
			ns:    []int{500, 100, 300, 200, 400}, // Unsorted input sizes
			vals: []*big.Float{
				big.NewFloat(250.0),
				big.NewFloat(50.0),
				big.NewFloat(150.0),
				big.NewFloat(100.0),
				big.NewFloat(200.0),
			}, // Corresponding unsorted timings
			expected: Linear.label,
			minScore: 0.95, // Should have high correlation for linear data
		},
		{
			name:  "quadratic time unsorted big",
			order: Quadratic,
			ns:    []int{20, 5, 15, 10}, // Unsorted input sizes
			vals: []*big.Float{
				big.NewFloat(400.0), // 20²=400
				big.NewFloat(25.0),  // 5²=25
				big.NewFloat(225.0), // 15²=225
				big.NewFloat(100.0), // 10²=100
			},
			expected: Quadratic.label,
			minScore: 0.95, // Should have very high correlation for perfect quadratic data
		},
		{
			name:  "constant time unsorted big - low variance",
			order: Constant,
			ns:    []int{1000, 100, 500, 200, 300}, // Unsorted input sizes
			vals: []*big.Float{
				big.NewFloat(10.1),
				big.NewFloat(9.9),
				big.NewFloat(10.0),
				big.NewFloat(10.2),
				big.NewFloat(9.8),
			}, // Very consistent timings (~10 with small variance)
			expected: Constant.label,
			minScore: 0.95, // Should score high for constant behavior
		},
		{
			name:  "constant time unsorted big - higher variance",
			order: Constant,
			ns:    []int{50, 200, 100, 300, 150}, // Unsorted input sizes
			vals: []*big.Float{
				big.NewFloat(10.1),
				big.NewFloat(8.9),
				big.NewFloat(10.0),
				big.NewFloat(11.0),
				big.NewFloat(9.8),
			}, // More variance but still relatively constant
			expected: Constant.label,
			minScore: 0.9, // Should score lower due to higher variance
		},
		{
			name:  "exponential time unsorted big",
			order: Exponential,
			ns:    []int{10, 5, 15, 8}, // Unsorted input sizes
			vals: []*big.Float{
				big.NewFloat(1024.0),  // 2^10=1024
				big.NewFloat(32.0),    // 2^5=32
				big.NewFloat(32768.0), // 2^15=32768
				big.NewFloat(256.0),   // 2^8=256
			},
			expected: Exponential.label,
			minScore: 0.95, // Should have high correlation for exponential data
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rating, err := tt.order.RateBig(tt.ns, tt.vals)
			if err != nil {
				t.Fatalf("RateBig() returned error: %v", err)
			}

			if rating.bigO.label != tt.expected {
				t.Errorf("Expected complexity %s, got %s", tt.expected, rating.bigO.label)
			}

			if rating.score < tt.minScore {
				t.Errorf("Expected score >= %f, got %f", tt.minScore, rating.score)
			}

			t.Logf("Complexity: %s, Score: %f", rating.bigO.label, rating.score)
		})
	}
}

func TestOrderErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		order   *BigO
		ns      []int
		vals    []float64
		wantErr bool
	}{
		{
			name:    "mismatched lengths",
			order:   Linear,
			ns:      []int{1, 2, 3},
			vals:    []float64{1.0, 2.0}, // One less value than ns
			wantErr: true,
		},
		{
			name:    "insufficient data points",
			order:   Linear,
			ns:      []int{1, 2},
			vals:    []float64{1.0, 2.0}, // Only 2 points, need at least 3
			wantErr: true,
		},
		{
			name:    "empty data",
			order:   Linear,
			ns:      []int{},
			vals:    []float64{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rating, err := tt.order.Rate(tt.ns, tt.vals)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}

				if rating == nil {
					t.Error("Expected non-nil rating even on error")
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestDetectConstantTime(t *testing.T) {
	tests := []struct {
		name      string
		vals      []float64
		wantScore float64
		wantErr   bool
	}{
		{
			name:      "very constant - low CV",
			vals:      []float64{10.0, 10.01, 9.99, 10.02, 9.98}, // CV ~0.02
			wantScore: 1.0,                                       // Very constant
			wantErr:   false,
		},
		{
			name:      "quite constant",
			vals:      []float64{10.0, 10.5, 9.5, 10.8, 9.2}, // CV ~0.08
			wantScore: 0.9,                                   // Quite constant
			wantErr:   false,
		},
		{
			name:      "moderately constant",
			vals:      []float64{10.0, 11.0, 9.0, 12.0, 8.0}, // CV ~0.15
			wantScore: 0.7,                                   // Moderately constant
			wantErr:   false,
		},
		{
			name:      "somewhat constant",
			vals:      []float64{10.0, 12.0, 8.0, 13.0, 7.0}, // CV ~0.25
			wantScore: 0.5,                                   // Somewhat constant
			wantErr:   false,
		},
		{
			name:      "not very constant",
			vals:      []float64{10.0, 14.0, 6.0, 15.0, 5.0}, // CV ~0.4
			wantScore: 0.3,                                   // Not very constant
			wantErr:   false,
		},
		{
			name:      "not constant",
			vals:      []float64{10.0, 20.0, 5.0, 25.0, 2.0}, // CV ~0.8
			wantScore: 0.1,                                   // Not constant
			wantErr:   false,
		},
		{
			name:      "perfect constant",
			vals:      []float64{5.0, 5.0, 5.0, 5.0, 5.0}, // CV = 0
			wantScore: 1.0,                                // Very constant
			wantErr:   false,
		},
		{
			name:      "insufficient data points",
			vals:      []float64{10.0, 10.1}, // Only 2 points
			wantScore: 0.0,
			wantErr:   true,
		},
		{
			name:      "empty data",
			vals:      []float64{},
			wantScore: 0.0,
			wantErr:   true,
		},
		{
			name:      "single outlier",
			vals:      []float64{10.0, 10.1, 9.9, 10.05, 15.0}, // One outlier increases CV
			wantScore: 0.7,                                     // Should be scored lower due to outlier
			wantErr:   false,
		},
		{
			name:      "zero mean edge case",
			vals:      []float64{0.1, -0.1, 0.05, -0.05, 0.0}, // Mean close to zero
			wantScore: 0.0,                                    // Infinite variation due to zero mean with variance
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rating, err := Constant.detectConstantTime(tt.vals)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if rating.score != tt.wantScore {
				t.Errorf("Expected score %f, got %f", tt.wantScore, rating.score)
			}

			if rating.bigO != Constant {
				t.Errorf("Expected bigO to be Constant, got %v", rating.bigO)
			}

			t.Logf("CV-based score: %f", rating.score)
		})
	}
}

func TestDetectConstantTimeBig(t *testing.T) {
	tests := []struct {
		name      string
		vals      []*big.Float
		wantScore float64
		wantErr   bool
	}{
		{
			name: "very constant - low CV",
			vals: []*big.Float{
				big.NewFloat(10.0),
				big.NewFloat(10.01),
				big.NewFloat(9.99),
				big.NewFloat(10.02),
				big.NewFloat(9.98),
			}, // CV ~0.02
			wantScore: 1.0, // Very constant
			wantErr:   false,
		},
		{
			name: "quite constant",
			vals: []*big.Float{
				big.NewFloat(10.0),
				big.NewFloat(10.5),
				big.NewFloat(9.5),
				big.NewFloat(10.8),
				big.NewFloat(9.2),
			}, // CV ~0.08
			wantScore: 0.9, // Quite constant
			wantErr:   false,
		},
		{
			name: "moderately constant",
			vals: []*big.Float{
				big.NewFloat(10.0),
				big.NewFloat(11.0),
				big.NewFloat(9.0),
				big.NewFloat(12.0),
				big.NewFloat(8.0),
			}, // CV ~0.15
			wantScore: 0.7, // Moderately constant
			wantErr:   false,
		},
		{
			name: "somewhat constant",
			vals: []*big.Float{
				big.NewFloat(10.0),
				big.NewFloat(12.0),
				big.NewFloat(8.0),
				big.NewFloat(13.0),
				big.NewFloat(7.0),
			}, // CV ~0.25
			wantScore: 0.5, // Somewhat constant
			wantErr:   false,
		},
		{
			name: "not very constant",
			vals: []*big.Float{
				big.NewFloat(10.0),
				big.NewFloat(14.0),
				big.NewFloat(6.0),
				big.NewFloat(15.0),
				big.NewFloat(5.0),
			}, // CV ~0.4
			wantScore: 0.3, // Not very constant
			wantErr:   false,
		},
		{
			name: "not constant",
			vals: []*big.Float{
				big.NewFloat(10.0),
				big.NewFloat(20.0),
				big.NewFloat(5.0),
				big.NewFloat(25.0),
				big.NewFloat(2.0),
			}, // CV ~0.8
			wantScore: 0.1, // Not constant
			wantErr:   false,
		},
		{
			name: "perfect constant",
			vals: []*big.Float{
				big.NewFloat(5.0),
				big.NewFloat(5.0),
				big.NewFloat(5.0),
				big.NewFloat(5.0),
				big.NewFloat(5.0),
			}, // CV = 0
			wantScore: 1.0, // Very constant
			wantErr:   false,
		},
		{
			name: "insufficient data points",
			vals: []*big.Float{
				big.NewFloat(10.0),
				big.NewFloat(10.1),
			}, // Only 2 points
			wantScore: 0.0,
			wantErr:   true,
		},
		{
			name:      "empty data",
			vals:      []*big.Float{},
			wantScore: 0.0,
			wantErr:   true,
		},
		{
			name: "single outlier",
			vals: []*big.Float{
				big.NewFloat(10.0),
				big.NewFloat(10.1),
				big.NewFloat(9.9),
				big.NewFloat(10.05),
				big.NewFloat(15.0),
			}, // One outlier increases CV
			wantScore: 0.7, // Should be scored lower due to outlier
			wantErr:   false,
		},
		{
			name: "high precision values",
			vals: []*big.Float{
				big.NewFloat(1000000.00001),
				big.NewFloat(1000000.00002),
				big.NewFloat(999999.99999),
				big.NewFloat(1000000.00003),
				big.NewFloat(999999.99998),
			}, // Very small relative variance with big numbers
			wantScore: 1.0, // Very constant
			wantErr:   false,
		},
		{
			name: "zero mean edge case",
			vals: []*big.Float{
				big.NewFloat(0.1),
				big.NewFloat(-0.1),
				big.NewFloat(0.05),
				big.NewFloat(-0.05),
				big.NewFloat(0.0),
			}, // Mean close to zero
			wantScore: 0.0, // Infinite variation due to zero mean with variance
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rating, err := Constant.detectConstantTimeBig(tt.vals)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if rating.score != tt.wantScore {
				t.Errorf("Expected score %f, got %f", tt.wantScore, rating.score)
			}

			if rating.bigO != Constant {
				t.Errorf("Expected bigO to be Constant, got %v", rating.bigO)
			}

			t.Logf("CV-based score: %f", rating.score)
		})
	}
}

func TestDetectConstantTimeConsistency(t *testing.T) {
	// Test that detectConstantTime and detectConstantTimeBig produce consistent results
	// for the same data converted between float64 and big.Float
	testCases := [][]float64{
		{10.0, 10.01, 9.99, 10.02, 9.98}, // Very constant
		{10.0, 10.5, 9.5, 10.8, 9.2},     // Quite constant
		{10.0, 11.0, 9.0, 12.0, 8.0},     // Moderately constant
		{10.0, 14.0, 6.0, 15.0, 5.0},     // Not very constant
		{5.0, 5.0, 5.0, 5.0, 5.0},        // Perfect constant
	}

	for i, vals := range testCases {
		t.Run(fmt.Sprintf("consistency_test_%d", i), func(t *testing.T) {
			// Test with float64
			ratingFloat, err := Constant.detectConstantTime(vals)
			if err != nil {
				t.Fatalf("detectConstantTime returned error: %v", err)
			}

			// Convert to big.Float
			bigVals := make([]*big.Float, len(vals))
			for j, v := range vals {
				bigVals[j] = big.NewFloat(v)
			}

			// Test with big.Float
			ratingBig, err := Constant.detectConstantTimeBig(bigVals)
			if err != nil {
				t.Fatalf("detectConstantTimeBig returned error: %v", err)
			}

			// Compare scores (should be identical)
			if ratingFloat.score != ratingBig.score {
				t.Errorf("Score mismatch: float64 version = %f, big.Float version = %f",
					ratingFloat.score, ratingBig.score)
			}

			t.Logf("Consistent score: %f", ratingFloat.score)
		})
	}
}

func TestNegativeIntegerHandling(t *testing.T) {
	tests := []struct {
		name    string
		order   *BigO
		ns      []int
		vals    []float64
		wantErr bool
		errMsg  string
	}{
		{
			name:    "all negative input sizes - filtered out, insufficient data",
			order:   Linear,
			ns:      []int{-10, -5, -1},
			vals:    []float64{10.0, 5.0, 1.0},
			wantErr: true,
			errMsg:  "there must be at least 3 (but preferably more) data points",
		},
		{
			name:    "mixed positive and negative input sizes - negatives filtered out",
			order:   Linear,
			ns:      []int{10, -5, 15, 20},
			vals:    []float64{10.0, 5.0, 15.0, 20.0},
			wantErr: false,
			errMsg:  "",
		},
		{
			name:    "zero input size - filtered out with insufficient remaining data",
			order:   Linear,
			ns:      []int{0, 5},
			vals:    []float64{0.0, 5.0},
			wantErr: true,
			errMsg:  "there must be at least 3 (but preferably more) data points",
		},
		{
			name:    "zero input size - filtered out with sufficient remaining data",
			order:   Linear,
			ns:      []int{0, 5, 10, 15},
			vals:    []float64{0.0, 5.0, 10.0, 15.0},
			wantErr: false,
			errMsg:  "",
		},
		{
			name:    "negative timing values - no validation error",
			order:   Linear,
			ns:      []int{5, 10, 15},
			vals:    []float64{-5.0, 10.0, 15.0},
			wantErr: false,
			errMsg:  "",
		},
		{
			name:    "all negative timing values - no validation error",
			order:   Linear,
			ns:      []int{5, 10, 15},
			vals:    []float64{-5.0, -10.0, -15.0},
			wantErr: false,
			errMsg:  "",
		},
		{
			name:    "zero timing values allowed",
			order:   Constant,
			ns:      []int{5, 10, 15},
			vals:    []float64{0.0, 0.0, 0.0},
			wantErr: false,
			errMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rating, err := tt.order.Rate(tt.ns, tt.vals)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errMsg, err.Error())
				}

				if rating == nil {
					t.Error("Expected non-nil rating even on error")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if rating == nil {
					t.Error("Expected non-nil rating")
				}
			}
		})
	}
}

func TestNegativeIntegerHandlingBig(t *testing.T) {
	tests := []struct {
		name    string
		order   *BigO
		ns      []int
		vals    []*big.Float
		wantErr bool
		errMsg  string
	}{
		{
			name:  "all negative input sizes big - filtered out, insufficient data",
			order: Linear,
			ns:    []int{-10, -5, -1},
			vals: []*big.Float{
				big.NewFloat(10.0),
				big.NewFloat(5.0),
				big.NewFloat(1.0),
			},
			wantErr: true,
			errMsg:  "there must be at least 3 (but preferably more) data points",
		},
		{
			name:  "mixed positive and negative input sizes big - negatives filtered out",
			order: Linear,
			ns:    []int{10, -5, 15, 20},
			vals: []*big.Float{
				big.NewFloat(10.0),
				big.NewFloat(5.0),
				big.NewFloat(15.0),
				big.NewFloat(20.0),
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name:  "zero input size big - filtered out with insufficient remaining data",
			order: Linear,
			ns:    []int{0, 5},
			vals: []*big.Float{
				big.NewFloat(0.0),
				big.NewFloat(5.0),
			},
			wantErr: true,
			errMsg:  "there must be at least 3 (but preferably more) data points",
		},
		{
			name:  "zero input size big - filtered out with sufficient remaining data",
			order: Linear,
			ns:    []int{0, 5, 10, 15},
			vals: []*big.Float{
				big.NewFloat(0.0),
				big.NewFloat(5.0),
				big.NewFloat(10.0),
				big.NewFloat(15.0),
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name:  "negative timing values big - no validation error",
			order: Linear,
			ns:    []int{5, 10, 15},
			vals: []*big.Float{
				big.NewFloat(-5.0),
				big.NewFloat(10.0),
				big.NewFloat(15.0),
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name:  "all negative timing values big - no validation error",
			order: Linear,
			ns:    []int{5, 10, 15},
			vals: []*big.Float{
				big.NewFloat(-5.0),
				big.NewFloat(-10.0),
				big.NewFloat(-15.0),
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name:  "zero timing values allowed big",
			order: Constant,
			ns:    []int{5, 10, 15},
			vals: []*big.Float{
				big.NewFloat(0.0),
				big.NewFloat(0.0),
				big.NewFloat(0.0),
			},
			wantErr: false,
			errMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rating, err := tt.order.RateBig(tt.ns, tt.vals)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errMsg, err.Error())
				}

				if rating == nil {
					t.Error("Expected non-nil rating even on error")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if rating == nil {
					t.Error("Expected non-nil rating")
				}
			}
		})
	}
}
