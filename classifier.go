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
	"bytes"
	"encoding/csv"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// Classifier is used to Classify a set of data points to find the rating they
// are most likely to correspond to.
type Classifier struct {
	// Interim storage of values until we are ready to Classify the data.
	// This is a mapping of the 'N' to the values recorded for that N.
	data    map[int][]float64
	dataBig map[int][]*big.Float

	// classified is set to true once the data has been Classified at least once.
	classified bool

	// Rating is the top rating that the data has been Classified as.
	rating *Rating

	// ratings is the set of all ratings for this data.
	ratings []*Rating
}

// NewClassifier creates a new Classifier.
func NewClassifier() *Classifier {
	return &Classifier{
		data:       make(map[int][]float64),
		dataBig:    make(map[int][]*big.Float),
		classified: false,
		rating:     defaultRating,
		ratings:    make([]*Rating, 0),
	}
}

// AddDataPoint adds the given values to the data.
// Non-positive input sizes (n <= 0) are ignored and not added to the dataset.
func (o *Classifier) AddDataPoint(n int, values ...float64) error {
	// Ignore non-positive input sizes
	if n <= 0 {
		return nil
	}

	if o.data == nil {
		o.data = make(map[int][]float64)
	}

	o.data[n] = append(o.data[n], values...)

	return nil
}

// AddDataPointBig adds the given values to the data.
// Non-positive input sizes (n <= 0) are ignored and not added to the dataset.
func (o *Classifier) AddDataPointBig(n int, values ...*big.Float) error {
	// Ignore non-positive input sizes
	if n <= 0 {
		return nil
	}

	if o.dataBig == nil {
		o.dataBig = make(map[int][]*big.Float)
	}

	o.dataBig[n] = append(o.dataBig[n], values...)

	return nil
}

// AddDataPoints adds all the data points and associated values.
func (o *Classifier) AddDataPoints(n []int, values [][]float64) error {
	if o.data == nil {
		o.data = make(map[int][]float64)
	}

	if len(n) != len(values) {
		return fmt.Errorf("sizes and corresponding values must be the same length")
	}

	for i, n := range n {
		err := o.AddDataPoint(n, values[i]...)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddDataPointsBig adds all the data points and associated values.
func (o *Classifier) AddDataPointsBig(n []int, values [][]*big.Float) error {
	if o.dataBig == nil {
		o.dataBig = make(map[int][]*big.Float)
	}

	if len(n) != len(values) {
		return fmt.Errorf("sizes and corresponding values must be the same length")
	}

	for i, n := range n {
		err := o.AddDataPointBig(n, values[i]...)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddBenchmarkResult adds the result of a benchmark test to the data.
// It extracts the nanoseconds per operation from the BenchmarkResult and adds it
// as a data point using the number of iterations (N) as the input size.
// This method provides seamless integration with Go's built-in benchmarking system.
//
// The timing data is converted from nanoseconds to a float64 value representing
// the average execution time per operation. Non-positive iteration counts are ignored.
//
// Usage example:
//
//	result := testing.Benchmark(func(b *testing.B) {
//	    for i := 0; i < b.N; i++ {
//	        // your algorithm here
//	    }
//	})
//	classifier.AddBenchmarkResult(result)
func (o *Classifier) AddBenchmarkResult(result testing.BenchmarkResult) error {
	// Calculate nanoseconds per operation more precisely
	// to avoid integer division truncation for very fast operations
	var timeValue float64
	if result.N > 0 {
		// Use float division to preserve precision
		timeValue = float64(result.T.Nanoseconds()) / float64(result.N)
	} else {
		timeValue = 0.0
	}

	// Use the iteration count as the input size (N)
	// This assumes the benchmark scales the problem size with b.N
	inputSize := result.N

	// Add the data point using the existing method
	return o.AddDataPoint(inputSize, timeValue)
}

// readCSV reads and parses a 2-column delimiter separated file.
// The header parameter controls whether the first line of the file is a header
// or not and should be skipped.
// Returns two slices: one for the N values and one for the corresponding measurements.
// Non-positive input sizes are filtered out during parsing.
func readCSV(path string, header bool, delimiter rune) ([]int, []float64, error) {
	csvFile, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	defer func() {
		_ = csvFile.Close()
	}()

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = delimiter

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	// Check there was something in there.
	if len(records) == 0 {
		return nil, nil, fmt.Errorf("no records found in file %s", path)
	}

	if header {
		records = records[1:]
	}

	var ns []int
	var vals []float64

	for i, record := range records {
		// We only care that there are at least 2 columns. Any additional
		// columns are ignored.
		if len(record) < 2 {
			return nil, nil, fmt.Errorf("not enough columns (%d) for record %d", len(record), i)
		}

		// We do not expect the N values to ever exceed an int64, if that occurs
		// return an error.
		n, err := strconv.Atoi(strings.TrimSpace(record[0]))
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing field 1 value for record %d: %w", i, err)
		}

		// Skip non-positive input sizes
		if n <= 0 {
			continue
		}

		// It could be possible someone is running extreme algorithms and getting
		// values beyond math.MaxFloat64, or values that would fall in the cracks
		// between recordable values of a float64, so we want to support reading
		// values in as big.Float if that occurs.

		// TODO(rsned): Add support to try big.Float if ParseFloat fails.
		val, err := strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing field 2 value for record %d: %w", i, err)
		}

		ns = append(ns, n)
		vals = append(vals, val)
	}

	return ns, vals, nil
}

// LoadCSV loads the data from a 2-column delimiter separated file and adds the values.
// The header parameter controls whether the first line of the file is a header
// or not and should be skipped. If an error occurred, no data will be loaded.
// Any errors encountered are returned.
// Non-positive input sizes are automatically filtered out during loading.
func (o *Classifier) LoadCSV(path string, header bool, delimiter rune) error {
	ns, vals, err := readCSV(path, header, delimiter)
	if err != nil {
		return err
	}

	for i, n := range ns {
		if err := o.AddDataPoint(n, vals[i]); err != nil {
			return fmt.Errorf("failed to add data point: %w", err)
		}
	}

	return nil
}

// Classify is used to Classify the data so far and determine the most
// Big O fit. Can be run as often as needed when more data are added.
//
// Common errors that can occur here are lack of distinct data points to be able
// to analyzye, or errors in a specific BigO test.
//
// For every distinct N value, the values stored for it are averaged to get a working
// value. The assumption in this package is that the user will do more than a single
// run of their code to get real world timing results, so we average all the run values
// for a given N.
//
// TODO(rsned): Add support for the big.Float values as well.
// TODO(rsned): If there are at least 30-50 values for a given N, then we can run some
// basic stats tests on the data points to test for outliers in the data	.
// TODO(rsned): If the number of values for a given N are large enough,
// should we include the option to discard the outlier in the values?
func (o *Classifier) Classify() (*Rating, error) {
	if len(o.data) < 3 {
		return defaultRating, fmt.Errorf("not enough data points (%d) to Classify", len(o.data))
	}

	// Start with an unset ranking.
	o.rating = &Rating{
		bigO:  Unrated,
		score: -1,
	}

	var Ns []int
	// First pass through, pull out the distinct N values and order them.
	for N := range o.data {
		Ns = append(Ns, N)
	}

	sort.Ints(Ns)

	// Second pass through, compute the average result value for the given N
	// and store it.
	vals := make([]float64, len(Ns))
	for i, N := range Ns {
		val := 0.0
		for _, v := range o.data[N] {
			val += v
		}

		// TODO(rsned): Check for divide by 0.
		vals[i] = val / float64(len(o.data[N]))
	}

	if len(o.dataBig) > 0 {
		return defaultRating, fmt.Errorf("big.Float data not implemented yet")
	}

	// Reset ratings slice for fresh classification
	o.ratings = make([]*Rating, 0)

	// Now for each potential BigO complexity, generate and save its ranking.
	//
	// TODO(rsned): Add support for the big.Float values as well.
	var lastErr error
	for _, b := range BigOOrdered {
		// In some cases, to prevent huge amounts of computation in *big.Float
		// land, it is advisable to pre-scale the values down to smaller ranges.
		// e.g.,
		if Ns[len(Ns)-1] > b.scalingCutoff {
			// TODO(rsned): Scaling down Ns and vals to avoid blowout computation.
			continue
		}

		rating, err := b.Rate(Ns, vals)
		if err != nil {
			fmt.Printf("Error ranking %s: %v\n", b.label, err)
			lastErr = err
		}

		o.ratings = append(o.ratings, rating)

		// If this score is higher than the last best, promote it to the leading result.
		if rating.score > o.rating.score {
			o.rating = rating
		}
	}

	o.classified = true

	sort.Slice(o.ratings, func(i, j int) bool {
		return o.ratings[i].bigO.rank < o.ratings[j].bigO.rank
	})

	return o.rating, lastErr
}

// GetAllRatings returns a copy of all the ratings generated by the most recent Classify() call.
// Returns nil if Classify() has not been called yet.
// The ratings are sorted by BigO rank (lowest rank first).
// The returned slice is a copy and can be safely modified without affecting internal state.
func (o *Classifier) GetAllRatings() []*Rating {
	if !o.classified {
		return nil
	}
	// Return a copy of the ratings slice to prevent external modification of internal state
	ratingsCopy := make([]*Rating, len(o.ratings))
	copy(ratingsCopy, o.ratings)
	return ratingsCopy
}

// Summary returns a longer form view of the results as a formatted text blob.
func (o *Classifier) Summary() string {
	if !o.classified {
		return "Not Classified yet"
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "\nBigO:  %s\n", o.rating.bigO.String())
	fmt.Fprintf(&buf, "Num data points: %d\n", len(o.data))

	// TODO(rsned): Add min/max values for N and Vals to the output.

	const winner = " *"
	const other = ""

	for _, r := range o.ratings {
		addedText := other
		if r.score == o.rating.score {
			addedText = winner
		}

		fmt.Fprintf(&buf, "%15s:   %0.8f%s\n", r.bigO.label, r.score, addedText)
	}

	return buf.String()
}
