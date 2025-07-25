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
	"math/big"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	// sampleFiles is a set of timing data generated offline for use in testing.
	sampleFiles = []string{
		"testdata/example_constant_time.csv",
		"testdata/example_cubic_time.csv",
		"testdata/example_exponential_time.csv",
		"testdata/example_factorial_time.csv",
		"testdata/example_hyper_exponential_time.csv",
		"testdata/example_linear_time.csv",
		"testdata/example_linearithmic_time.csv",
		"testdata/example_logarithmic_time.csv",
		"testdata/example_log_log_time.csv",
		"testdata/example_n_log_star_n_time.csv",
		"testdata/example_polylogarithmic_time.csv",
		"testdata/example_polynomial_time.csv",
		"testdata/example_quadratic_time.csv",
	}

	// sampleTimings is a map of unique name string to the set of
	// N and corresponding values for that N.
	sampleTimings = map[string][]struct {
		n    int
		vals []float64
	}{}
)

func init() {
	for _, file := range sampleFiles {
		// Extract filename without path and extension to use as key
		base := filepath.Base(file)
		key := strings.TrimSuffix(base, filepath.Ext(base))

		// Read the CSV file (assuming no header and comma delimiter)
		ns, vals, err := readCSV(file, true, ',')
		if err != nil {
			// Skip files that can't be read (they may not exist in test environment)
			continue
		}

		// Convert to the expected struct format
		var timings []struct {
			n    int
			vals []float64
		}

		for i, n := range ns {
			timings = append(timings, struct {
				n    int
				vals []float64
			}{
				n:    n,
				vals: []float64{vals[i]},
			})
		}

		sampleTimings[key] = timings
	}
}

func TestReadCSV(t *testing.T) {
	tests := []struct {
		name      string
		file      string
		header    bool
		delimiter rune
		wantErr   bool
		wantNs    []int
		wantVals  []float64
	}{
		{
			name:      "valid file with header",
			file:      "testdata/valid_with_header.csv",
			header:    true,
			delimiter: ',',
			wantErr:   false,
			wantNs:    []int{40, 100, 200, 300, 400, 500},
			wantVals:  []float64{0.5, 1.5, 30, 45, 60.3, 76.1},
		},
		{
			name:      "valid file with header and extra columns",
			file:      "testdata/valid_with_header_extra_columns.csv",
			header:    true,
			delimiter: ',',
			wantErr:   false,
			wantNs:    []int{100, 200, 300, 400, 500},
			wantVals:  []float64{1.5, 3, 4.5, 6, 7.5},
		},
		{
			name:      "valid file without header",
			file:      "testdata/valid_no_header.csv",
			header:    false,
			delimiter: ',',
			wantErr:   false,
			wantNs:    []int{100, 200, 300, 400, 500},
			wantVals:  []float64{1.5, 3.0, 4.5, 6.0, 7.5},
		},
		{
			name:      "trailing blank lines",
			file:      "testdata/valid_trailing_blank_lines.csv",
			header:    false,
			delimiter: ',',
			wantErr:   false,
			wantNs:    []int{100, 200, 300},
			wantVals:  []float64{432, 578, 924},
		},
		{
			name:      "blank lines in middle",
			file:      "testdata/valid_blank_lines_in_middle.csv",
			header:    false,
			delimiter: ',',
			wantErr:   false,
			wantNs:    []int{100, 200, 300},
			wantVals:  []float64{432, 578, 924},
		},
		// Error cases
		{
			name:      "file not found",
			file:      "testdata/nonexistent.csv",
			header:    false,
			delimiter: ',',
			wantErr:   true,
			wantNs:    nil,
			wantVals:  nil,
		},
		{
			name:      "empty file",
			file:      "testdata/empty.csv",
			header:    false,
			delimiter: ',',
			wantErr:   true,
			wantNs:    nil,
			wantVals:  nil,
		},
		{
			name:      "insufficient columns",
			file:      "testdata/error_insufficient_columns.csv",
			header:    false,
			delimiter: ',',
			wantErr:   true,
			wantNs:    nil,
			wantVals:  nil,
		},
		{
			name:      "invalid number format in column 1",
			file:      "testdata/error_invalid_column_1.csv",
			header:    false,
			delimiter: ',',
			wantErr:   true,
			wantNs:    nil,
			wantVals:  nil,
		},
		{
			name:      "invalid number format in column 2",
			file:      "testdata/error_invalid_column_2.csv",
			header:    false,
			delimiter: ',',
			wantErr:   true,
			wantNs:    nil,
			wantVals:  nil,
		},
		{
			name:      "negative input sizes (filtered out)",
			file:      "testdata/negative_input_sizes.csv",
			header:    true,
			delimiter: ',',
			wantErr:   false,
			wantNs:    []int{},
			wantVals:  []float64{},
		},
		{
			name:      "negative timing values",
			file:      "testdata/negative_timing_values.csv",
			header:    true,
			delimiter: ',',
			wantErr:   false,
			wantNs:    []int{100, 200, 300, 400, 500},
			wantVals:  []float64{-500, -1000, -1500, -2000, -2500},
		},
		{
			name:      "mixed negative values (only positive n kept)",
			file:      "testdata/mixed_negative_values.csv",
			header:    true,
			delimiter: ',',
			wantErr:   false,
			wantNs:    []int{50, 100},
			wantVals:  []float64{-250, -400},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns, vals, err := readCSV(tt.file, tt.header, tt.delimiter)

			if tt.wantErr && err == nil {
				t.Errorf("expected error but got none")

				return
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error = %v", err)

				return
			}

			if !tt.wantErr {
				if len(ns) != len(vals) {
					t.Errorf("returned slices of different lengths: ns=%d, vals=%d", len(ns), len(vals))
				}

				// Only check for empty slices if we expect non-empty results
				expectEmpty := (len(tt.wantNs) == 0) && (len(tt.wantVals) == 0)
				if !expectEmpty && (len(ns) == 0 || len(vals) == 0) {
					t.Errorf("returned empty slices for valid file")
				}

				// Verify loaded data matches expected values if provided
				if tt.wantNs != nil && !cmp.Equal(ns, tt.wantNs, cmpopts.EquateEmpty()) {
					t.Errorf("Ns = %v, want %v, diff = %v", ns, tt.wantNs, cmp.Diff(ns, tt.wantNs, cmpopts.EquateEmpty()))
				}

				if tt.wantVals != nil && !cmp.Equal(vals, tt.wantVals, cmpopts.EquateEmpty()) {
					t.Errorf("Vals = %v, want %v, diff = %v", vals, tt.wantVals, cmp.Diff(vals, tt.wantVals, cmpopts.EquateEmpty()))
				}
			}
		})
	}
}

func TestLoadCSV(t *testing.T) {
	// Most of this is tested in TestReadCSV. We only need to do simple checks
	// here that it loads something in when readCSV returns data.
	tests := []struct {
		name       string
		file       string
		header     bool
		delimiter  rune
		wantErr    bool
		numRecords int
	}{
		{
			name:       "valid file with header",
			file:       "testdata/valid_with_header.csv",
			header:     true,
			delimiter:  ',',
			wantErr:    false,
			numRecords: 6,
		},
		{
			name:       "valid file with header and extra columns",
			file:       "testdata/valid_with_header_extra_columns.csv",
			header:     true,
			delimiter:  ',',
			wantErr:    false,
			numRecords: 5,
		},
		// Error cases
		{
			name:       "file not found",
			file:       "testdata/nonexistent.csv",
			header:     false,
			delimiter:  ',',
			wantErr:    true,
			numRecords: 0,
		},
		{
			name:       "empty file",
			file:       "testdata/empty.csv",
			header:     false,
			delimiter:  ',',
			wantErr:    true,
			numRecords: 0,
		},
		{
			name:       "insufficient columns",
			file:       "testdata/error_insufficient_columns.csv",
			header:     false,
			delimiter:  ',',
			wantErr:    true,
			numRecords: 0,
		},
		{
			name:       "negative input sizes (filtered out)",
			file:       "testdata/negative_input_sizes.csv",
			header:     true,
			delimiter:  ',',
			wantErr:    false,
			numRecords: 0,
		},
		{
			name:       "negative timing values",
			file:       "testdata/negative_timing_values.csv",
			header:     true,
			delimiter:  ',',
			wantErr:    false,
			numRecords: 5,
		},
		{
			name:       "mixed negative values (only positive n kept)",
			file:       "testdata/mixed_negative_values.csv",
			header:     true,
			delimiter:  ',',
			wantErr:    false,
			numRecords: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewClassifier()
			err := b.LoadCSV(tt.file, tt.header, tt.delimiter)
			if tt.wantErr && err == nil {
				t.Errorf("expected error but got none")

				return
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error = %v", err)

				return
			}

			if !tt.wantErr && len(b.data) != tt.numRecords {
				t.Errorf("did not populate the expected data: %d records != %d", len(b.data), tt.numRecords)
			}
		})
	}
}

func TestClassify(t *testing.T) {
	// For this test we use exact values to ensure we trigger the right rating.
	tests := []struct {
		name       string
		setupData  func(*Classifier) error
		wantErr    bool
		wantRating bool
		expectBigO *BigO
	}{
		{
			name: "insufficient data points",
			setupData: func(c *Classifier) error {
				_ = c.AddDataPoint(10, 1.0)
				_ = c.AddDataPoint(20, 2.0)

				return nil
			},
			wantErr:    true,
			wantRating: false,
			expectBigO: Unrated,
		},
		{
			name: "constant time pattern",
			setupData: func(c *Classifier) error {
				_ = c.AddDataPoint(100, 1.0)
				_ = c.AddDataPoint(200, 1.1)
				_ = c.AddDataPoint(400, 0.9)
				_ = c.AddDataPoint(800, 1.0)
				_ = c.AddDataPoint(1600, 1.2)

				return nil
			},
			wantErr:    false,
			wantRating: true,
			expectBigO: Constant,
		},
		{
			name: "linear time pattern",
			setupData: func(c *Classifier) error {
				_ = c.AddDataPoint(100, 100.0)
				_ = c.AddDataPoint(200, 200.0)
				_ = c.AddDataPoint(400, 400.0)
				_ = c.AddDataPoint(800, 800.0)
				_ = c.AddDataPoint(1600, 1600.0)

				return nil
			},
			wantErr:    false,
			wantRating: true,
			expectBigO: Linear,
		},
		{
			name: "quadratic time pattern",
			setupData: func(c *Classifier) error {
				_ = c.AddDataPoint(10, 100.0)
				_ = c.AddDataPoint(20, 400.0)
				_ = c.AddDataPoint(30, 900.0)
				_ = c.AddDataPoint(40, 1600.0)
				_ = c.AddDataPoint(50, 2500.0)

				return nil
			},
			wantErr:    false,
			wantRating: true,
			expectBigO: Quadratic,
		},
		{
			name: "multiple values per data point",
			setupData: func(c *Classifier) error {
				_ = c.AddDataPoint(100, 98.0, 99.0, 101.0, 102.0)
				_ = c.AddDataPoint(200, 196.0, 199.0, 201.0, 204.0)
				_ = c.AddDataPoint(400, 398.0, 401.0, 399.0, 402.0)
				_ = c.AddDataPoint(800, 796.0, 801.0, 804.0, 799.0)
				_ = c.AddDataPoint(1600, 1598.0, 1601.0, 1604.0, 1597.0)

				return nil
			},
			wantErr:    false,
			wantRating: true,
			expectBigO: Linear,
		},
		{
			name: "with big.Float data (not yet supported)",
			setupData: func(c *Classifier) error {
				_ = c.AddDataPoint(100, 100.0)
				_ = c.AddDataPoint(200, 200.0)
				_ = c.AddDataPoint(400, 400.0)
				_ = c.AddDataPointBig(800, big.NewFloat(800.0))

				return nil
			},
			wantErr:    true,
			wantRating: false,
			expectBigO: Unrated,
		},
		{
			name: "already Classified - can re-Classify",
			setupData: func(c *Classifier) error {
				_ = c.AddDataPoint(100, 100.0)
				_ = c.AddDataPoint(200, 200.0)
				_ = c.AddDataPoint(400, 400.0)
				// First characterization
				_, err := c.Classify()
				if err != nil {
					return err
				}
				// Add more data
				_ = c.AddDataPoint(800, 800.0)
				_ = c.AddDataPoint(1600, 1600.0)

				return nil
			},
			wantErr:    false,
			wantRating: true,
			expectBigO: Linear,
		},
		{
			name: "scaling cutoff exceeded for some BigO types",
			setupData: func(c *Classifier) error {
				// Use large N values that would exceed scaling cutoff for exponential/factorial
				_ = c.AddDataPoint(1000, 1000.0)
				_ = c.AddDataPoint(2000, 2000.0)
				_ = c.AddDataPoint(4000, 4000.0)
				_ = c.AddDataPoint(8000, 8000.0)
				_ = c.AddDataPoint(16000, 16000.0)

				return nil
			},
			wantErr:    false,
			wantRating: true,
			expectBigO: Linear,
		},
		{
			name: "linear pattern - data added in decreasing N order",
			setupData: func(c *Classifier) error {
				// Add data points in reverse order to verify order independence
				_ = c.AddDataPoint(1600, 1600.0)
				_ = c.AddDataPoint(800, 800.0)
				_ = c.AddDataPoint(400, 400.0)
				_ = c.AddDataPoint(200, 200.0)
				_ = c.AddDataPoint(100, 100.0)

				return nil
			},
			wantErr:    false,
			wantRating: true,
			expectBigO: Linear,
		},
		{
			name: "quadratic pattern - data added in random N order",
			setupData: func(c *Classifier) error {
				// Add data points in random order to verify order independence
				_ = c.AddDataPoint(30, 900.0)
				_ = c.AddDataPoint(10, 100.0)
				_ = c.AddDataPoint(50, 2500.0)
				_ = c.AddDataPoint(20, 400.0)
				_ = c.AddDataPoint(40, 1600.0)

				return nil
			},
			wantErr:    false,
			wantRating: true,
			expectBigO: Quadratic,
		},
		{
			name: "constant pattern - data added in mixed N order",
			setupData: func(c *Classifier) error {
				// Add data points in mixed order to verify order independence
				_ = c.AddDataPoint(800, 1.0)
				_ = c.AddDataPoint(100, 1.1)
				_ = c.AddDataPoint(1600, 1.2)
				_ = c.AddDataPoint(200, 0.9)
				_ = c.AddDataPoint(400, 1.0)

				return nil
			},
			wantErr:    false,
			wantRating: true,
			expectBigO: Constant,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClassifier()

			err := tt.setupData(c)
			if err != nil {
				t.Errorf("failed to setup test data: %v", err)

				return
			}

			rating, err := c.Classify()
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")

					return
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error = %v", err)

				return
			}

			if !tt.wantRating {
				t.Errorf("expected no rating but got one")

				return
			}

			if rating == nil {
				t.Errorf("returned nil rating")

				return
			}

			// Verify the classifier is marked as Classified
			if !c.classified {
				t.Errorf("did not mark classifier as classified")
			}

			// Verify ratings were populated
			if len(c.ratings) == 0 {
				t.Errorf("did not populate ratings")
			}

			// Verify the rating matches expectation if specified
			if tt.expectBigO != nil && rating.bigO != tt.expectBigO {
				t.Errorf("rating = %v, want %v", rating.bigO.label, tt.expectBigO.label)
			}

			// Verify score is valid (should be >= 0 for successful
			// characterization with the given test data)
			if rating.score < 0 {
				t.Errorf("Classify() rating score = %v, want >= 0", rating.score)
			}

			// Verify Summary works after characterization
			summary := c.Summary()
			if strings.Contains(summary, "Not Classified yet") {
				t.Errorf("Summary() reports not Classified after successful Classify()")
			}
		})
	}
}

type dataPoint struct {
	n   int
	val float64
}

func TestClassifierNegativeIntegerHandling(t *testing.T) {
	tests := []struct {
		name           string
		addDataPoints  []dataPoint
		expectedPoints int
		expectedErr    bool
	}{
		{
			name: "all positive input sizes",
			addDataPoints: []dataPoint{
				{n: 5, val: 10.0},
				{n: 10, val: 20.0},
				{n: 15, val: 30.0},
			},
			expectedPoints: 3,
			expectedErr:    false,
		},
		{
			name: "mixed positive and negative input sizes",
			addDataPoints: []dataPoint{
				{n: 5, val: 10.0},
				{n: -3, val: 5.0},
				{n: 10, val: 20.0},
				{n: 0, val: 0.0},
				{n: 15, val: 30.0},
			},
			expectedPoints: 3, // Only positive values should be added
			expectedErr:    false,
		},
		{
			name: "all negative and zero input sizes",
			addDataPoints: []dataPoint{
				{n: -5, val: 10.0},
				{n: -3, val: 5.0},
				{n: 0, val: 0.0},
			},
			expectedPoints: 0,
			expectedErr:    true,
		},
		{
			name: "insufficient positive data after filtering",
			addDataPoints: []dataPoint{
				{n: 5, val: 10.0},
				{n: -3, val: 5.0},
				{n: 10, val: 20.0},
				{n: 0, val: 0.0},
			},
			expectedPoints: 2,
			expectedErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			classifier := NewClassifier()

			// Add all data points
			for _, dp := range tt.addDataPoints {
				err := classifier.AddDataPoint(dp.n, dp.val)
				if err != nil {
					t.Fatalf("AddDataPoint returned unexpected error: %v", err)
				}
			}

			// Check the number of data points actually stored
			if len(classifier.data) != tt.expectedPoints {
				t.Errorf("Expected %d data points, got %d", tt.expectedPoints, len(classifier.data))
			}

			// Try to classify
			_, err := classifier.Classify()

			if tt.expectedErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

type dataPointBig struct {
	n   int
	val *big.Float
}

func TestClassifierNegativeIntegerHandlingBig(t *testing.T) {
	tests := []struct {
		name           string
		addDataPoints  []dataPointBig
		expectedPoints int
		expectedErr    bool
	}{
		{
			name: "all positive input sizes big",
			addDataPoints: []dataPointBig{
				{n: 5, val: big.NewFloat(10.0)},
				{n: 10, val: big.NewFloat(20.0)},
				{n: 15, val: big.NewFloat(30.0)},
			},
			expectedPoints: 3,
			expectedErr:    false,
		},
		{
			name: "mixed positive and negative input sizes big",
			addDataPoints: []dataPointBig{
				{n: 5, val: big.NewFloat(10.0)},
				{n: -3, val: big.NewFloat(5.0)},
				{n: 10, val: big.NewFloat(20.0)},
				{n: 0, val: big.NewFloat(0.0)},
				{n: 15, val: big.NewFloat(30.0)},
			},
			expectedPoints: 3, // Only positive values should be added
			expectedErr:    false,
		},
		{
			name: "all negative and zero input sizes big",
			addDataPoints: []dataPointBig{
				{n: -5, val: big.NewFloat(10.0)},
				{n: -3, val: big.NewFloat(5.0)},
				{n: 0, val: big.NewFloat(0.0)},
			},
			expectedPoints: 0,
			expectedErr:    false, // Classify() doesn't handle big.Float yet
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			classifier := NewClassifier()

			// Add all data points
			for _, dp := range tt.addDataPoints {
				err := classifier.AddDataPointBig(dp.n, dp.val)
				if err != nil {
					t.Fatalf("AddDataPointBig returned unexpected error: %v", err)
				}
			}

			// Check the number of data points actually stored
			if len(classifier.dataBig) != tt.expectedPoints {
				t.Errorf("Expected %d data points, got %d", tt.expectedPoints, len(classifier.dataBig))
			}
		})
	}
}

func TestAddBenchmarkResult(t *testing.T) {
	tests := []struct {
		name           string
		benchResult    testing.BenchmarkResult
		expectedN      int
		expectedValue  float64
		expectedPoints int
		wantErr        bool
	}{
		{
			name: "valid benchmark result",
			benchResult: testing.BenchmarkResult{
				N: 1000,
				T: 5000 * time.Nanosecond, // 5000ns total, so 5ns per op
			},
			expectedN:      1000,
			expectedValue:  5.0,
			expectedPoints: 1,
			wantErr:        false,
		},
		{
			name: "benchmark with zero iterations",
			benchResult: testing.BenchmarkResult{
				N: 0,
				T: 1000 * time.Nanosecond,
			},
			expectedN:      0,
			expectedValue:  0,
			expectedPoints: 0, // Should be filtered out due to non-positive N
			wantErr:        false,
		},
		{
			name: "benchmark with negative iterations",
			benchResult: testing.BenchmarkResult{
				N: -100,
				T: 1000 * time.Nanosecond,
			},
			expectedN:      -100,
			expectedValue:  0,
			expectedPoints: 0, // Should be filtered out due to negative N
			wantErr:        false,
		},
		{
			name: "benchmark with very fast operation",
			benchResult: testing.BenchmarkResult{
				N: 1000000,
				T: 1 * time.Nanosecond, // 1ns total, so 0.000001ns per op
			},
			expectedN:      1000000,
			expectedValue:  0.000001,
			expectedPoints: 1,
			wantErr:        false,
		},
		{
			name: "benchmark with slow operation",
			benchResult: testing.BenchmarkResult{
				N: 10,
				T: 10 * time.Second, // Very slow: 1 billion ns per op
			},
			expectedN:      10,
			expectedValue:  1000000000.0,
			expectedPoints: 1,
			wantErr:        false,
		},
		{
			name: "benchmark with memory allocations (ignored)",
			benchResult: testing.BenchmarkResult{
				N:         500,
				T:         2500 * time.Nanosecond, // 5ns per op
				MemAllocs: 1000,                   // Should be ignored
				MemBytes:  50000,                  // Should be ignored
			},
			expectedN:      500,
			expectedValue:  5.0,
			expectedPoints: 1,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			classifier := NewClassifier()

			err := classifier.AddBenchmarkResult(tt.benchResult)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Check that the correct number of data points were added
			if len(classifier.data) != tt.expectedPoints {
				t.Errorf("expected %d data points, got %d", tt.expectedPoints, len(classifier.data))
				return
			}

			// If we expect data points, verify the values
			if tt.expectedPoints > 0 {
				values, exists := classifier.data[tt.expectedN]
				if !exists {
					t.Errorf("expected data point with N=%d not found", tt.expectedN)
					return
				}

				if len(values) != 1 {
					t.Errorf("expected 1 value for N=%d, got %d", tt.expectedN, len(values))
					return
				}

				if values[0] != tt.expectedValue {
					t.Errorf("expected value %f for N=%d, got %f", tt.expectedValue, tt.expectedN, values[0])
				}
			}
		})
	}
}

func TestClassifierAddBenchmarkResultMultipleResults(t *testing.T) {
	classifier := NewClassifier()

	// Add multiple benchmark results with slight variation to avoid zero variance
	results := []testing.BenchmarkResult{
		{N: 100, T: 500 * time.Nanosecond},   // 5.0ns per op
		{N: 200, T: 1020 * time.Nanosecond},  // 5.1ns per op
		{N: 400, T: 1960 * time.Nanosecond},  // 4.9ns per op
		{N: 800, T: 4040 * time.Nanosecond},  // 5.05ns per op
		{N: 1600, T: 7920 * time.Nanosecond}, // 4.95ns per op
	}

	for _, result := range results {
		err := classifier.AddBenchmarkResult(result)
		if err != nil {
			t.Errorf("unexpected error adding benchmark result: %v", err)
			return
		}
	}

	// Verify all data points were added
	if len(classifier.data) != 5 {
		t.Errorf("expected 5 data points, got %d", len(classifier.data))
		return
	}

	// Verify the data represents roughly constant time complexity
	expectedData := map[int][]float64{
		100:  {5.0},
		200:  {5.1},
		400:  {4.9},
		800:  {5.05},
		1600: {4.95},
	}

	for n, expectedValues := range expectedData {
		actualValues, exists := classifier.data[n]
		if !exists {
			t.Errorf("expected data point with N=%d not found", n)
			continue
		}

		if !cmp.Equal(actualValues, expectedValues) {
			t.Errorf("data for N=%d: got %v, want %v", n, actualValues, expectedValues)
		}
	}

	// Test that we can classify this as constant time
	rating, err := classifier.Classify()
	if err != nil {
		t.Errorf("classification failed: %v", err)
		return
	}

	if rating.bigO != Constant {
		t.Errorf("expected Constant classification, got %s", rating.bigO.label)
	}
}

func TestClassifierAddBenchmarkResultIntegrationWithRealBenchmark(t *testing.T) {
	classifier := NewClassifier()

	// Create a linear-time function that actually scales with input size
	linearFunc := func(arr []int) int {
		sum := 0
		for _, v := range arr {
			sum += v
		}
		return sum
	}

	// Create arrays of different sizes
	problemSizes := []int{100, 200, 400, 800, 1600}

	for _, size := range problemSizes {
		// Create test data of the specified size
		testData := make([]int, size)
		for i := range testData {
			testData[i] = i + 1
		}

		// Run a benchmark that processes the full array each time
		result := testing.Benchmark(func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = linearFunc(testData)
			}
		})

		// Create a fake benchmark result with the size as N for our test
		// In real usage, you'd run the benchmark with different problem sizes
		// and collect the results appropriately
		fakeResult := testing.BenchmarkResult{
			N: size,     // Use problem size as N
			T: result.T, // Use actual timing
		}

		// Add the benchmark result to our classifier
		err := classifier.AddBenchmarkResult(fakeResult)
		if err != nil {
			t.Errorf("failed to add benchmark result for size %d: %v", size, err)
			continue
		}
	}

	// Verify we collected data
	if len(classifier.data) != len(problemSizes) {
		t.Errorf("expected %d data points, got %d", len(problemSizes), len(classifier.data))
	}

	// The classification might vary due to benchmark timing variability,
	// but we should be able to classify it without error
	rating, err := classifier.Classify()
	if err != nil {
		t.Errorf("classification failed: %v", err)
		return
	}

	if rating == nil {
		t.Errorf("got nil rating")
		return
	}

	// Verify the classifier was marked as classified
	if !classifier.classified {
		t.Errorf("classifier not marked as classified")
	}
}

func TestAddBenchmarkResult_ZeroDurationHandling(t *testing.T) {
	classifier := NewClassifier()

	// Test with zero duration (could happen with very fast operations)
	result := testing.BenchmarkResult{
		N: 1000,
		T: 0, // Zero duration
	}

	err := classifier.AddBenchmarkResult(result)
	if err != nil {
		t.Errorf("unexpected error with zero duration: %v", err)
		return
	}

	// Should add data point with 0.0 ns per operation
	values, exists := classifier.data[1000]
	if !exists {
		t.Errorf("expected data point not found")
		return
	}

	if len(values) != 1 || values[0] != 0.0 {
		t.Errorf("expected [0.0], got %v", values)
	}
}

func TestClassifierGetAllRatings(t *testing.T) {
	tests := []struct {
		name             string
		setupData        func(*Classifier) error
		expectNil        bool
		expectNumRatings int
		expectSorted     bool
	}{
		{
			name: "not classified yet",
			setupData: func(c *Classifier) error {
				_ = c.AddDataPoint(100, 100.0)
				_ = c.AddDataPoint(200, 200.0)
				_ = c.AddDataPoint(400, 400.0)
				// Don't call Classify()
				return nil
			},
			expectNil:        true,
			expectNumRatings: 0,
			expectSorted:     false,
		},
		{
			name: "after classification - linear pattern",
			setupData: func(c *Classifier) error {
				_ = c.AddDataPoint(100, 100.0)
				_ = c.AddDataPoint(200, 200.0)
				_ = c.AddDataPoint(400, 400.0)
				_ = c.AddDataPoint(800, 800.0)
				_ = c.AddDataPoint(1600, 1600.0)
				_, err := c.Classify()
				return err
			},
			expectNil:        false,
			expectNumRatings: 11, // Some BigO types get filtered out due to scaling cutoffs
			expectSorted:     true,
		},
		{
			name: "after classification - constant pattern",
			setupData: func(c *Classifier) error {
				_ = c.AddDataPoint(100, 1.0)
				_ = c.AddDataPoint(200, 1.1)
				_ = c.AddDataPoint(400, 0.9)
				_ = c.AddDataPoint(800, 1.0)
				_ = c.AddDataPoint(1600, 1.2)
				_, err := c.Classify()
				return err
			},
			expectNil:        false,
			expectNumRatings: 11, // Some BigO types get filtered out due to scaling cutoffs
			expectSorted:     true,
		},
		{
			name: "after classification - quadratic pattern",
			setupData: func(c *Classifier) error {
				_ = c.AddDataPoint(10, 100.0)
				_ = c.AddDataPoint(20, 400.0)
				_ = c.AddDataPoint(30, 900.0)
				_ = c.AddDataPoint(40, 1600.0)
				_ = c.AddDataPoint(50, 2500.0)
				_, err := c.Classify()
				return err
			},
			expectNil:        false,
			expectNumRatings: 13, // All BigO types should be available for small input sizes
			expectSorted:     true,
		},
		{
			name: "re-classification replaces ratings (fixed behavior)",
			setupData: func(c *Classifier) error {
				// First classification with limited data
				_ = c.AddDataPoint(100, 100.0)
				_ = c.AddDataPoint(200, 200.0)
				_ = c.AddDataPoint(400, 400.0)
				_, err := c.Classify()
				if err != nil {
					return err
				}

				// Add more data and re-classify
				_ = c.AddDataPoint(800, 800.0)
				_ = c.AddDataPoint(1600, 1600.0)
				_, err = c.Classify()
				return err
			},
			expectNil:        false,
			expectNumRatings: 11, // Should have fresh ratings (not appended), some filtered due to scaling cutoffs
			expectSorted:     true,
		},
		{
			name: "large input sizes filter some BigO types",
			setupData: func(c *Classifier) error {
				// Use large N values that would exceed scaling cutoff for some BigO types
				_ = c.AddDataPoint(1000, 1000.0)
				_ = c.AddDataPoint(2000, 2000.0)
				_ = c.AddDataPoint(4000, 4000.0)
				_ = c.AddDataPoint(8000, 8000.0)
				_ = c.AddDataPoint(16000, 16000.0)
				_, err := c.Classify()
				return err
			},
			expectNil:        false,
			expectNumRatings: -1, // Variable number due to scaling cutoffs
			expectSorted:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClassifier()

			err := tt.setupData(c)
			if err != nil {
				t.Errorf("failed to setup test data: %v", err)
				return
			}

			ratings := c.GetAllRatings()

			if tt.expectNil {
				if ratings != nil {
					t.Errorf("expected nil ratings but got %d ratings", len(ratings))
				}
				return
			}

			if ratings == nil {
				t.Errorf("expected ratings but got nil")
				return
			}

			if tt.expectNumRatings >= 0 && len(ratings) != tt.expectNumRatings {
				t.Errorf("expected %d ratings, got %d", tt.expectNumRatings, len(ratings))
			}

			if len(ratings) == 0 {
				t.Errorf("expected at least some ratings but got empty slice")
				return
			}

			// Verify that all ratings have valid BigO types and scores
			for i, rating := range ratings {
				if rating == nil {
					t.Errorf("rating %d is nil", i)
					continue
				}
				if rating.bigO == nil {
					t.Errorf("rating %d has nil BigO", i)
				}
				// Score can be negative for poor fits, so no validation on score value
			}

			// Verify ratings are sorted by BigO rank if expected
			if tt.expectSorted && len(ratings) > 1 {
				for i := 1; i < len(ratings); i++ {
					if ratings[i-1].bigO.rank > ratings[i].bigO.rank {
						t.Errorf("ratings not sorted by rank: rating %d (%s, rank %d) > rating %d (%s, rank %d)",
							i-1, ratings[i-1].bigO.label, ratings[i-1].bigO.rank,
							i, ratings[i].bigO.label, ratings[i].bigO.rank)
						break
					}
				}
			}

			// Verify that the ratings slice is a copy and safe to modify
			// The returned slice should NOT affect internal state when modified
			originalLen := len(ratings)
			if originalLen > 0 {
				// Save original value for comparison
				originalFirstRating := ratings[0]
				// Modify the returned slice
				ratings[0] = nil
				// Get ratings again and verify it was NOT affected (fixed behavior)
				newRatings := c.GetAllRatings()
				if newRatings[0] == nil {
					t.Errorf("modifying returned ratings slice affected internal state - should return a copy")
				}
				if newRatings[0] != originalFirstRating {
					t.Errorf("expected internal state to be unchanged after modifying returned slice")
				}
			}
		})
	}
}

func TestClassifierGetAllRatingsConsistencyWithClassify(t *testing.T) {
	// Test that GetAllRatings returns consistent data with what Classify() produces
	c := NewClassifier()

	// Add linear pattern data
	_ = c.AddDataPoint(100, 100.0)
	_ = c.AddDataPoint(200, 200.0)
	_ = c.AddDataPoint(400, 400.0)
	_ = c.AddDataPoint(800, 800.0)
	_ = c.AddDataPoint(1600, 1600.0)

	topRating, err := c.Classify()
	if err != nil {
		t.Errorf("classification failed: %v", err)
		return
	}

	allRatings := c.GetAllRatings()
	if allRatings == nil {
		t.Errorf("GetAllRatings returned nil after successful classification")
		return
	}

	// Verify the top rating from Classify() appears in the all ratings
	foundTopRating := false
	for _, rating := range allRatings {
		if rating.bigO == topRating.bigO && rating.score == topRating.score {
			foundTopRating = true
			break
		}
	}

	if !foundTopRating {
		t.Errorf("top rating from Classify() not found in GetAllRatings() results")
	}

	// Verify that ratings are ordered consistently
	// The best score should be among the ratings (though not necessarily first due to rank sorting)
	bestScore := topRating.score
	foundBestScore := false
	for _, rating := range allRatings {
		if rating.score == bestScore {
			foundBestScore = true
			break
		}
	}

	if !foundBestScore {
		t.Errorf("best score %.8f not found in ratings", bestScore)
	}
}

func TestClassifierGetAllRatingsEmptyAfterReset(t *testing.T) {
	// Test behavior when classifier state is reset
	c := NewClassifier()

	// Add data and classify
	_ = c.AddDataPoint(100, 100.0)
	_ = c.AddDataPoint(200, 200.0)
	_ = c.AddDataPoint(400, 400.0)
	_ = c.AddDataPoint(800, 800.0)
	_ = c.AddDataPoint(1600, 1600.0)

	_, err := c.Classify()
	if err != nil {
		t.Errorf("classification failed: %v", err)
		return
	}

	// Verify we have ratings
	ratingsBeforeReset := c.GetAllRatings()
	if len(ratingsBeforeReset) == 0 {
		t.Errorf("expected ratings before reset")
		return
	}

	// Create a new classifier (simulating reset)
	c2 := NewClassifier()
	ratingsAfterReset := c2.GetAllRatings()

	if ratingsAfterReset != nil {
		t.Errorf("expected nil ratings from new classifier, got %d ratings", len(ratingsAfterReset))
	}
}

func TestClassifierGetAllRatingsReturnsCopy(t *testing.T) {
	// Test that GetAllRatings returns a copy that can be safely modified
	c := NewClassifier()

	// Add data and classify
	_ = c.AddDataPoint(100, 100.0)
	_ = c.AddDataPoint(200, 200.0)
	_ = c.AddDataPoint(400, 400.0)
	_ = c.AddDataPoint(800, 800.0)
	_ = c.AddDataPoint(1600, 1600.0)

	_, err := c.Classify()
	if err != nil {
		t.Errorf("classification failed: %v", err)
		return
	}

	// Get ratings twice
	ratings1 := c.GetAllRatings()
	ratings2 := c.GetAllRatings()

	if len(ratings1) == 0 || len(ratings2) == 0 {
		t.Errorf("expected ratings from both calls")
		return
	}

	// Verify they are different slice instances (different addresses)
	if &ratings1[0] == &ratings2[0] {
		t.Errorf("expected different slice instances, got same slice reference")
	}

	// Save original values
	originalFirst1 := ratings1[0]
	originalFirst2 := ratings2[0]

	// Modify first slice
	ratings1[0] = nil

	// Verify second slice is unaffected
	if ratings2[0] == nil {
		t.Errorf("modifying first returned slice affected second returned slice")
	}

	// Get fresh ratings and verify internal state is unaffected
	ratings3 := c.GetAllRatings()
	if ratings3[0] == nil {
		t.Errorf("modifying returned slice affected internal state")
	}

	// Verify values are consistent
	if ratings3[0] != originalFirst1 || ratings3[0] != originalFirst2 {
		t.Errorf("internal state changed unexpectedly")
	}

	// Verify we can modify the returned slice extensively without issues
	for i := range ratings1 {
		ratings1[i] = nil
	}

	// Internal state should still be intact
	ratings4 := c.GetAllRatings()
	if len(ratings4) != len(ratings2) {
		t.Errorf("extensive modification of returned slice affected internal state length")
	}

	for i, rating := range ratings4 {
		if rating == nil {
			t.Errorf("internal state corrupted at index %d after extensive modification of returned slice", i)
		}
	}
}

func TestClassifierGetAllRatingsReClassificationReplacesRatings(t *testing.T) {
	// Test that re-classification replaces ratings instead of appending
	c := NewClassifier()

	// First classification with limited data
	_ = c.AddDataPoint(100, 100.0)
	_ = c.AddDataPoint(200, 200.0)
	_ = c.AddDataPoint(400, 400.0)

	_, err := c.Classify()
	if err != nil {
		t.Errorf("first classification failed: %v", err)
		return
	}

	ratingsAfterFirst := c.GetAllRatings()
	firstCount := len(ratingsAfterFirst)
	if firstCount == 0 {
		t.Errorf("expected ratings after first classification")
		return
	}

	// Add more data and re-classify
	_ = c.AddDataPoint(800, 800.0)
	_ = c.AddDataPoint(1600, 1600.0)

	_, err = c.Classify()
	if err != nil {
		t.Errorf("second classification failed: %v", err)
		return
	}

	ratingsAfterSecond := c.GetAllRatings()
	secondCount := len(ratingsAfterSecond)

	// The key test: ratings should be replaced, not appended
	// Both classifications should have similar counts (with possible filtering
	// differences) but definitely NOT double the count.
	if secondCount >= firstCount*2 {
		t.Errorf("ratings appear to be appended instead of replaced: first=%d, second=%d",
			firstCount, secondCount)
	}

	// Verify that ratings are fresh (not containing duplicates from first
	// classification). We can check this by ensuring no duplicate BigO
	// types exist.
	seenBigO := make(map[string]bool)
	for _, rating := range ratingsAfterSecond {
		if seenBigO[rating.bigO.label] {
			t.Errorf("found duplicate BigO type %s, indicating ratings were appended instead of replaced",
				rating.bigO.label)
		}
		seenBigO[rating.bigO.label] = true
	}

	// Additional verification: the count should be reasonable
	// (similar to single classification)

	// Allow for some filtering variations
	expectedRange := []int{10, 11, 12, 13}
	countInRange := slices.Contains(expectedRange, secondCount)

	if !countInRange {
		t.Errorf("rating count %d after re-classification seems unreasonable, expected one of %v",
			secondCount, expectedRange)
	}
}
