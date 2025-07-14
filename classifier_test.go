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
	"strings"
	"testing"

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
