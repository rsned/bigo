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
	"sort"

	"github.com/rsned/bigmath"
	"github.com/rsned/stats/correlation"
)

// BigO is used to describe the complexity class of an algorithm, which
// represents how the runtime or space requirements of an algorithm grow as the
// input size grows.
type BigO struct {
	active      bool   // active is used to indicate if this instance should be used.
	rank        int    // rank is used for sorting and ordering.
	label       string // The English language O(...) label for this entry.
	description string // A brief English description of what this Big O means.

	// scalingCutoff is the value of an 'N' which would trigger a shift to
	// pre-scaling the inputs so that high order BigO values don't mushroom
	// out of control. (e.g., Factorial when the N is in the 10's or 100's of
	// thousands). For most Big O, it's not an issue.
	scalingCutoff int

	// floatCutoffMin is the smallest float64/int input value that will allow
	// this function to return a valid float64 value. e.g., Log functions
	// on values <1 will return NaN, 0 will return -Inf, etc.
	floatCutoffMin float64

	// floatCutoffMax is the largest float64/int input value that will allow
	// this function to return a valid float64 value. e.g. 2^n will exceed
	// the limit of a float64 for n > 1024. Values beyond this will need to
	// use variations of functions that can handle big.Float.
	floatCutoffMax float64

	// When reading in the inputs, values are assumed to be strings that need to
	// be converted to numeric values. As the number is read in, it is compared
	// the cutoffs and math.MaxFloat64 to figure out where it falls in our
	// ability to generate a value to correlate it with.
	//
	// These three functions are used for the forward computation of values
	// depending on where the float64 overflow comes into to play. Any value
	// below floatCutoffMax will use the funcFloatFloat helper to generate
	// values to be correlated with. If the value lands above floatCutoffMax,
	// it will end up in the big.Float realm so we would end up using the
	// funcFloatBig helper. And any value beyond math.MaxFloat64 will start as
	// a big.Float and (likely) end up as a big.Float so we would use the
	// funcBigBig. (We say likely because BigO types less than Linear could
	// actually transform the value back into float64 range. e.g. DoubleLog
	// would reduce a big.Float input considerably), but at that point we stay
	// in the big.Float realm.
	funcFloatFloat floatFloatFunc
	funcFloatBig   floatBigFunc
	funcBigBig     bigBigFunc

	/*
		TODO(rsned): When we want to attempt a reverse transform to move from
		other number range back to linear, use these functions.

		// Similarly with the three functions above, we need a way to convert
		// the inputs into a more 'linear' space to make the regression analysis
		// simpler, so we use these functions to 'flatten' the values.
		funcTranformFloatFloat floatFloatFunc
		funcTranformFloatBig   floatBigFunc
		funcTranformBigBig     bigBigFunc
	*/

	// TODO(rsned): Add value range measures to help gauge how will the set
	// of values will be affected by the BigO. e.g., For Linearithmic, it
	// really needs 3-4 order of magnitude change in N to see the difference
	// between O(n), O(log* n), and O(n log n) because within the same
	// magnitude log* and log n are basically constants which makes them
	// O(k*n) which reduces to O(n) and you can't tell the apart.
}

// Label returns the label of this BigO.
func (o *BigO) Label() string {
	return o.label
}

// Description returns the description of this BigO.
func (o *BigO) Description() string {
	return o.description
}

// String returns this instances label in string form.
func (o BigO) String() string {
	return o.label
}

// filterPositiveData filters out non-positive input sizes and their corresponding values.
// Returns filtered slices and the number of filtered entries.
func filterPositiveData(ns []int, vals []float64) ([]int, []float64, int) {
	var filteredNs []int
	var filteredVals []float64
	filtered := 0

	for i, n := range ns {
		if n > 0 {
			filteredNs = append(filteredNs, n)
			filteredVals = append(filteredVals, vals[i])
		} else {
			filtered++
		}
	}

	return filteredNs, filteredVals, filtered
}

// filterPositiveDataBig filters out non-positive input sizes and their corresponding values.
// Returns filtered slices and the number of filtered entries.
func filterPositiveDataBig(ns []int, vals []*big.Float) ([]int, []*big.Float, int) {
	var filteredNs []int
	var filteredVals []*big.Float
	filtered := 0

	for i, n := range ns {
		if n > 0 {
			filteredNs = append(filteredNs, n)
			filteredVals = append(filteredVals, vals[i])
		} else {
			filtered++
		}
	}

	return filteredNs, filteredVals, filtered
}

// Rate generates the ranking for this particular BigO using its defined
// cutoffs and mapping functions.
//
// To try to make some of the math easier, we scale the incoming values by
// the smallest N in the data so that the comparison may hopefully stay in
// the float64 range.
//
// Non-positive input sizes are filtered out before analysis.
func (o *BigO) Rate(ns []int, vals []float64) (*Rating, error) {
	if len(ns) != len(vals) {
		return defaultRating, fmt.Errorf("the N's and values must be the same length")
	}

	// Filter out non-positive input sizes
	filteredNs, filteredVals, filteredCount := filterPositiveData(ns, vals)
	if filteredCount > 0 {
		// Continue with filtered data but don't return an error
		ns = filteredNs
		vals = filteredVals
	}

	if len(ns) < 3 {
		return defaultRating, fmt.Errorf("there must be at least 3 (but preferably more) data points")
	}

	// Special handling for O(1) constant time detection
	if o == Constant {
		return o.detectConstantTime(vals)
	}

	// If the N crosses the max float cutoff, we will need to switch to big math.
	needsBig := false

	// Store the predicted values for the correlation analysis.
	// Because this function is called few times and not in any critical paths,
	// we will store both the float64 and big.Float versions of the values just
	// in case.
	var predicteds []float64
	var predictedBigs []*big.Float

	valsBig := make([]*big.Float, len(ns))

	// Range over the data, generating a comparison sequence using the
	// appropriate helpers for this BigO.
	for i, n := range ns {
		// If the N is below the limit that can be handled by this.
		// then set the predicted value to 0. (e.g., log(x) for x < 1 goes to
		// -Infinity)
		switch {
		case float64(n) < o.floatCutoffMin:
			predicteds = append(predicteds, 0)
			predictedBigs = append(predictedBigs, big.NewFloat(0))
		case float64(n) <= o.floatCutoffMax:
			// For values that will end up in the valid float64 range, use
			// the helper that generates a float output.
			predicteds = append(predicteds, o.funcFloatFloat(float64(n)))
			predictedBigs = append(predictedBigs, o.funcFloatBig(float64(n)))
		default:
			// We've crossed into needing big math. Update the boolean.
			needsBig = true
			// For the float64 predicteds, we have entered the realm of +Inf.
			predicteds = append(predicteds, math.Inf(1))
			predictedBigs = append(predictedBigs, o.funcFloatBig(float64(n)))
		}

		// We are making the bold assumption that no one uses an N larger
		// than math.MaxInt64.
		valsBig[i] = big.NewFloat(vals[i])
	}

	var corr float64
	var err error

	if !needsBig {
		corr, err = correlation.Correlate(predicteds, vals, correlation.Pearson)
		if err != nil {
			return defaultRating, err
		}
	} else {
		corr, err = correlation.CorrelateBig(predictedBigs, valsBig, correlation.Pearson)
		if err != nil {
			return defaultRating, err
		}
	}

	rating := &Rating{
		bigO:  o,
		score: corr,
	}

	return rating, nil
}

// RateBig is a helper function that rates the data using big.Float values to
// perform the correlation and scoring.
//
// Non-positive input sizes are filtered out before analysis.
func (o *BigO) RateBig(ns []int, vals []*big.Float) (*Rating, error) {
	if len(ns) != len(vals) {
		return defaultRating, fmt.Errorf("the N's and values must be the same length")
	}

	// Filter out non-positive input sizes
	filteredNs, filteredVals, filteredCount := filterPositiveDataBig(ns, vals)
	if filteredCount > 0 {
		// Continue with filtered data but don't return an error
		ns = filteredNs
		vals = filteredVals
	}

	if len(ns) < 3 {
		return defaultRating, fmt.Errorf("there must be at least 3 (but preferably more) data points")
	}

	// Special handling for O(1) constant time detection
	if o == Constant {
		return o.detectConstantTimeBig(vals)
	}

	// Figure out the smallest N and it's corresponding value. We would
	// like to believe that the values coming in here are sorted, but it's
	// not guaranteed.
	var startN = math.MaxInt
	var startVal *big.Float
	for i, n := range ns {
		if n < startN {
			startN = n
			startVal = new(big.Float).Copy(vals[i])
		}
	}

	// Store the predicted values for the correlation analysis.
	var predicteds []*big.Float
	scaledVals := make([]*big.Float, len(ns))

	// Range over the data, generating a comparison sequence using the
	// appropriate helpers for this BigO.
	for i, k := range ns {
		// Scale the N down to a reasonable starting point as a float.
		// TODO(rsned): What should be done if multiple values scale down
		// to the same value?
		scaledN := math.Max(1, float64(k)/float64(startN))

		// If the scaled N is below the limit that can be handled by this.
		// then set the predicted value to 0. (e.g., log(x) for x < 1 goes to
		// -Infinity)
		if scaledN < o.floatCutoffMin {
			predicteds = append(predicteds, big.NewFloat(0))
		} else if scaledN <= o.floatCutoffMax {
			predicteds = append(predicteds, o.funcFloatBig(scaledN))
		}

		scaledVals[i] = new(big.Float).Quo(vals[i], startVal)
	}

	var corr float64
	var err error

	corr, err = correlation.CorrelateBig(predicteds, scaledVals, correlation.Pearson)
	if err != nil {
		return defaultRating, err
	}

	rating := &Rating{
		bigO:  o,
		score: corr,
	}

	return rating, nil
}

// detectConstantTime implements special detection logic for O(1) complexity.
// Since constant time algorithms should have minimal variance in timing,
// we use coefficient of variation (CV = stddev/mean) instead of correlation.
func (o *BigO) detectConstantTime(vals []float64) (*Rating, error) {
	if len(vals) < 3 {
		return defaultRating, fmt.Errorf("not enough data points for constant time detection")
	}

	// Calculate mean
	mean := 0.0
	for _, v := range vals {
		mean += v
	}

	mean /= float64(len(vals))

	// Calculate standard deviation
	variance := 0.0
	for _, v := range vals {
		diff := v - mean
		variance += diff * diff
	}

	variance /= float64(len(vals))
	stddev := math.Sqrt(variance)

	// Calculate coefficient of variation (CV = stddev/mean)
	// Lower CV indicates more constant behavior
	// Handle edge case where mean is zero
	var cv float64
	if mean == 0 {
		// If mean is zero and all values are zero, this is perfectly constant
		if stddev == 0 {
			cv = 0
		} else {
			// If mean is zero but there's variance, this is not constant
			cv = math.Inf(1)
		}
	} else {
		cv = stddev / mean
	}
	score := cvToScore(cv)
	rating := &Rating{
		bigO:  o,
		score: score,
	}

	return rating, nil
}

// cvToScore converts a coefficient of variation to a score.
func cvToScore(cv float64) float64 {
	// Convert CV to a correlation-like score (higher is better)
	// CV < 0.1 (10% variation) = very constant = high score
	// CV > 0.5 (50% variation) = not constant = low score
	var score float64
	switch {
	case cv == 0 || math.IsNaN(cv):
		score = 1.0 // Perfectly constant
	case math.IsInf(cv, 1):
		score = 0.0 // Infinite variation, not constant at all
	case cv < 0.05:
		score = 1.0 // Essentially constant
	case cv < 0.1:
		score = 0.9 // Pretty close to constant
	case cv < 0.2:
		score = 0.7 // Moderately constant
	case cv < 0.3:
		score = 0.5 // Somewhat constant
	case cv < 0.5:
		score = 0.3 // Not very constant
	default:
		score = 0.1 // Not constant
	}

	return score
}

// detectConstantTimeBig implements special detection logic for O(1) complexity
// using *big.Float inputs instead of float64.
// Since constant time algorithms should have minimal variance in timing,
// we use coefficient of variation (CV = stddev/mean) instead of correlation.
func (o *BigO) detectConstantTimeBig(vals []*big.Float) (*Rating, error) {
	if len(vals) < 3 {
		return defaultRating, fmt.Errorf("not enough data points for constant time detection")
	}

	// Calculate mean
	mean := big.NewFloat(0.0)
	for _, v := range vals {
		mean.Add(mean, v)
	}

	mean.Quo(mean, big.NewFloat(float64(len(vals))))

	// Calculate standard deviation
	variance := big.NewFloat(0.0)
	for _, v := range vals {
		diff := new(big.Float).Sub(v, mean)
		diffSquared := new(big.Float).Mul(diff, diff)
		variance.Add(variance, diffSquared)
	}

	variance.Quo(variance, big.NewFloat(float64(len(vals))))

	// Convert to float64 for math.Sqrt - this is acceptable since we're dealing with variance
	varianceFloat, _ := variance.Float64()
	stddev := big.NewFloat(math.Sqrt(varianceFloat))

	// Calculate coefficient of variation (CV = stddev/mean)
	// Lower CV indicates more constant behavior
	// Handle edge case where mean is zero
	var cvFloat float64
	meanFloat, _ := mean.Float64()
	stddevFloat, _ := stddev.Float64()

	if meanFloat == 0 {
		// If mean is zero and all values are zero, this is perfectly constant
		if stddevFloat == 0 {
			cvFloat = 0
		} else {
			// If mean is zero but there's variance, this is not constant
			cvFloat = math.Inf(1)
		}
	} else {
		cv := new(big.Float).Quo(stddev, mean)
		cvFloat, _ = cv.Float64()
	}
	score := cvToScore(cvFloat)
	rating := &Rating{
		bigO:  o,
		score: score,
	}

	return rating, nil
}

// floatFloatFunc is an interface for functions that take a float64 input
// and return a value in float64 space.
type floatFloatFunc func(x float64) float64

// floatBigFunc is an interface for functions that take a float64 input
// and return a value in big.Float space.
type floatBigFunc func(x float64) *big.Float

// bigBigFunc is an interface for functions that take a big.Float input
// and return a value in big.Float space.
type bigBigFunc func(x *big.Float) *big.Float

var (
	// Unrated is a BigO instance for when we have not yet rated the data.
	Unrated = &BigO{
		active:         false,
		rank:           0,
		label:          "O(?)",
		description:    "Not Yet Rated",
		scalingCutoff:  0,
		floatCutoffMin: 0,
		floatCutoffMax: 0,
		funcFloatFloat: nil,
		funcFloatBig:   nil,
		funcBigBig:     nil,
	}

	// defaultBigO is the default value to use.
	defaultBigO = Unrated

	// Constant is a BigO instance for constant time functions.
	//
	// Uses special variance-based detection instead of correlation since
	// constant values have zero variance.
	Constant = &BigO{
		active:      true,
		rank:        1,
		label:       "O(1)",
		description: "An algorithm with O(1) complexity runs in constant time, regardless of the input size.",

		floatCutoffMin: math.SmallestNonzeroFloat64,
		floatCutoffMax: math.MaxFloat64,

		scalingCutoff: math.MaxInt64,

		funcFloatFloat: func(_ float64) float64 {
			return 1
		},
		funcFloatBig: func(_ float64) *big.Float {
			return big.NewFloat(1)
		},
		funcBigBig: func(_ *big.Float) *big.Float {
			return big.NewFloat(1)
		},
	}

	// InverseAckerman is a BigO instance for the inverse Ackerman function.
	// Almost as flat a curve as Constant, and quite similar to LogLog.
	InverseAckerman = &BigO{
		active:         false,
		rank:           2,
		label:          "O(α(n))",
		description:    "An algorithm with O(α(n)) complexity has a runtime that grows at the rate of the inverse Ackerman function with the input size.",
		floatCutoffMin: 1,
		floatCutoffMax: math.MaxFloat64,

		scalingCutoff: math.MaxInt64,

		funcFloatFloat: func(x float64) float64 {
			return inverseAckermann(int(x))
		},
		funcFloatBig: func(x float64) *big.Float {
			return inverseAckermannBig(big.NewFloat(x))
		},
		funcBigBig: inverseAckermannBig,
	}

	// LogLog is a BigO instance for the double log function.
	LogLog = &BigO{
		active:         true,
		rank:           4,
		label:          "O(log log n)",
		description:    "An algorithm with O(log log n) complexity has a runtime that grows logarithmically twice with the input size.",
		floatCutoffMin: math.Exp(math.Exp(0)),
		floatCutoffMax: math.MaxFloat64,

		scalingCutoff: math.MaxInt64,

		funcFloatFloat: func(x float64) float64 {
			return math.Log(math.Log(x))
		},
		funcFloatBig: func(x float64) *big.Float {
			return bigmath.Log(bigmath.Log(big.NewFloat(x)))
		},
		funcBigBig: func(x *big.Float) *big.Float {
			return bigmath.Log(bigmath.Log(x))
		},
	}

	// Log is a BigO instance for Log complexity functions.
	Log = &BigO{
		active:      true,
		rank:        8,
		label:       "O(log n)",
		description: "An algorithm with O(log n) complexity has a runtime that grows logarithmically as the input size increases. This often involves dividing the problem in half each time.",

		floatCutoffMin: math.Exp(0),
		floatCutoffMax: math.MaxFloat64,

		scalingCutoff: math.MaxInt64,

		funcFloatFloat: math.Log,
		funcFloatBig: func(x float64) *big.Float {
			return big.NewFloat(math.Log(x))
		},
		funcBigBig: bigmath.Log,
	}

	// Polylogarithmic is a BigO instance for the polylogarithmic function.
	Polylogarithmic = &BigO{
		active:         true,
		rank:           16,
		label:          "O((log n)^c)",
		description:    "An algorithm with O((log n)^c) complexity has a runtime that grows polylogarithmically with the input size.",
		floatCutoffMin: 1,
		floatCutoffMax: math.MaxFloat64,

		scalingCutoff: math.MaxInt64,

		funcFloatFloat: func(x float64) float64 {
			return math.Pow(math.Log(x), 4)
		},
		funcFloatBig: func(x float64) *big.Float {
			return bigmath.Pow(bigmath.Log(big.NewFloat(x)), big.NewFloat(4))
		},
		funcBigBig: func(x *big.Float) *big.Float {
			return bigmath.Pow(bigmath.Log(x), big.NewFloat(4))
		},
	}

	// Linear is a BigO instance for the linear function.
	Linear = &BigO{
		active:      true,
		rank:        32,
		label:       "O(n)",
		description: "An algorithm with O(n) complexity has a runtime that grows linearly with the input size.",

		scalingCutoff: math.MaxInt64,

		floatCutoffMin: 1,
		floatCutoffMax: math.MaxFloat64,

		funcFloatFloat: func(x float64) float64 {
			return x
		},
		funcFloatBig: big.NewFloat,
		funcBigBig: func(x *big.Float) *big.Float {
			return x
		},
	}

	// NLogStarN is a BigO instance for the n log star n function.
	NLogStarN = &BigO{
		active:         true,
		rank:           64,
		label:          "O(n log* n)",
		description:    "An algorithm with O(n log* n) complexity has a runtime that grows with n times the iterated logarithm of the input size.",
		floatCutoffMin: 1,
		floatCutoffMax: math.MaxFloat64,

		scalingCutoff: math.MaxInt64,

		funcFloatFloat: func(x float64) float64 {

			return x * float64(logStarFloat(x))
		},
		funcFloatBig: func(x float64) *big.Float {
			logStarVal := logStarBig(big.NewFloat(x))

			return new(big.Float).Mul(big.NewFloat(x), big.NewFloat(float64(logStarVal)))
		},
		funcBigBig: func(x *big.Float) *big.Float {
			logStarVal := logStarBig(x)

			return new(big.Float).Mul(x, big.NewFloat(float64(logStarVal)))
		},
	}

	// Linearithmic is a BigO instance for the linearithmic function.
	Linearithmic = &BigO{
		active:      true,
		rank:        128,
		label:       "O(n log n)",
		description: "An algorithm with O(n log n) complexity typically involves a combination of linear and logarithmic time, often seen in efficient sorting algorithms such as Merge Sort.",

		scalingCutoff: math.MaxInt64,

		floatCutoffMin: 1,
		floatCutoffMax: math.Log(math.MaxFloat64),

		funcFloatFloat: func(x float64) float64 {
			return x * math.Log(x)
		},
		funcFloatBig: func(x float64) *big.Float {
			return bigmath.Log(big.NewFloat(x)).Mul(bigmath.Log(big.NewFloat(x)), big.NewFloat(x))
		},
		funcBigBig: func(x *big.Float) *big.Float {
			return bigmath.Log(x).Mul(bigmath.Log(x), x)
		},
	}

	// Quadratic is a BigO instance for the quadratic function.
	Quadratic = &BigO{
		active:      true,
		rank:        256,
		label:       "O(n^2)",
		description: "An algorithm with O(n^2) complexity has a runtime that grows quadratically with the input size, often involving nested loops such as Bubble Sort.",

		scalingCutoff: math.MaxInt64,

		floatCutoffMin: 1,
		floatCutoffMax: math.Sqrt(math.MaxFloat64),

		funcFloatFloat: func(x float64) float64 {
			return x * x
		},
		funcFloatBig: func(x float64) *big.Float {
			return new(big.Float).Mul(big.NewFloat(x), big.NewFloat(x))
		},
		funcBigBig: func(x *big.Float) *big.Float {
			return new(big.Float).Mul(x, x)
		},
	}

	// Cubic is a BigO instance for the cubic function.
	Cubic = &BigO{
		active:      true,
		rank:        512,
		label:       "O(n^3)",
		description: "An algorithm with O(n^3) complexity has a runtime that grows cubically with the input size, often involving three nested loops.",

		scalingCutoff: math.MaxInt64,

		floatCutoffMin: 1,
		floatCutoffMax: math.Cbrt(math.MaxFloat64),
		funcFloatFloat: func(x float64) float64 {
			return x * x * x
		},
		funcFloatBig: func(x float64) *big.Float {
			c := new(big.Float).Mul(big.NewFloat(x), big.NewFloat(x))

			return c.Mul(c, big.NewFloat(x))
		},
		funcBigBig: func(x *big.Float) *big.Float {
			c := new(big.Float).Mul(x, x)

			return c.Mul(c, x)
		},
	}

	// Polynomial is a BigO instance for the polynomial function.
	Polynomial = &BigO{
		active:      true,
		rank:        1024,
		label:       "O(n^c)",
		description: "An algorithm with O(n^c) complexity has a runtime that grows polynomially with the input size, often involving nested loops.",

		scalingCutoff: 1000000,

		floatCutoffMin: 1,
		floatCutoffMax: math.MaxFloat64,
		funcFloatFloat: func(x float64) float64 {
			return math.Pow(x, 4)
		},
		funcFloatBig: func(x float64) *big.Float {
			return bigmath.PowFloat64(x, 4)
		},
		funcBigBig: func(x *big.Float) *big.Float {
			return bigmath.Pow(x, big.NewFloat(4))
		},
	}

	// Exponential is a BigO instance for the exponential function.
	Exponential = &BigO{
		active:      true,
		rank:        2048,
		label:       "O(2^n)",
		description: "An algorithm with O(2^n) complexity has a runtime that grows exponentially with the input size, often seen in recursive algorithms solving combinatorial problems.",

		scalingCutoff: 1000000,

		floatCutoffMin: 1,
		floatCutoffMax: 1024.0,

		funcFloatFloat: math.Exp2,
		funcFloatBig: func(x float64) *big.Float {
			return bigmath.PowFloat64(2, x)
		},
		funcBigBig: func(x *big.Float) *big.Float {
			return bigmath.Pow(big.NewFloat(2), x)
		},
	}

	// Factorial is a BigO instance for the factorial function.
	Factorial = &BigO{
		active:      true,
		rank:        4096,
		label:       "O(n!)",
		description: "An algorithm with O(n!) complexity has a runtime that grows factorially with the input size, often seen in problems involving permutations.",

		scalingCutoff: 1000,

		floatCutoffMin: 1,
		floatCutoffMax: 170,

		funcFloatFloat: func(x float64) float64 {
			return factorial(int(x))
		},
		funcFloatBig: func(x float64) *big.Float {
			return bigmath.FactorialFloat(big.NewFloat(x))
		},
		funcBigBig: bigmath.FactorialFloat,
	}

	// HyperExponential is a BigO instance for the hyper-exponential function.
	HyperExponential = &BigO{
		active:      true,
		rank:        8192,
		label:       "O(n^n)",
		description: "An algorithm with O(n^n) complexity has a runtime that grows hyper-exponentially with the input size.",

		scalingCutoff: 500,

		floatCutoffMin: 1,
		floatCutoffMax: 141,

		funcFloatFloat: func(x float64) float64 {
			return math.Pow(x, x)
		},
		funcFloatBig: func(x float64) *big.Float {
			return bigmath.PowFloat64(x, x)
		},
		funcBigBig: func(x *big.Float) *big.Float {
			return bigmath.Pow(x, x)
		},
	}

	// TODO(rsned): Add any other mildy or wildly useful BigO values.
)

var (
	// BigOOrdered holds the set of all BigO instances in rank order.
	BigOOrdered = []*BigO{}

	// allBigO defines the set of all BigO instances.
	allBigO = []*BigO{
		Unrated,
		Constant,
		InverseAckerman,
		LogLog,
		Log,
		Polylogarithmic,
		Linear,
		NLogStarN,
		Linearithmic,
		Quadratic,
		Cubic,
		Polynomial,
		Exponential,
		Factorial,
		HyperExponential,
	}
)

// init is used to initialize and sort the BigOOrder slice.
func init() {
	// At start time, go through the set of all Big O ranks we have
	// and add the active ones to use.
	for _, r := range allBigO {
		if r.active {
			BigOOrdered = append(BigOOrdered, r)
		}
	}

	sort.Slice(BigOOrdered, func(i, j int) bool {
		return BigOOrdered[i].rank < BigOOrdered[j].rank
	})
}
