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

package examples

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/rsned/bigo"
	"github.com/rsned/bigo/examples/constant"
	"github.com/rsned/bigo/examples/datatypes/collection"
	"github.com/rsned/bigo/examples/datatypes/tree"
	"github.com/rsned/bigo/examples/linear"
	"github.com/rsned/bigo/examples/loglog"
)

// These values are small enough for everything below exponential to run if not overridden.
const (
	benchmarkStartDefault = 1
	benchmarkEndDefault   = 100
	benchmarkStepDefault  = 10
)

// benchmarkExampleName is the name of the example method to benchmark.
var benchmarkExampleName = flag.String("benchmark_example_name", "", "Name of the example method to benchmark. See the code for the available options.")

// benchmarkStart is the starting value for benchmark iterations
var benchmarkStart = flag.Int("benchmark_start", benchmarkStartDefault, "Starting value for benchmark iterations")

// benchmarkEnd is the ending value for benchmark iterations
var benchmarkEnd = flag.Int("benchmark_end", benchmarkEndDefault, "Ending value for benchmark iterations")

// benchmarkStep is the step size for benchmark iterations
var benchmarkStep = flag.Int("benchmark_step", benchmarkStepDefault, "Step size for benchmark iterations")

// To cut out some of the timing variability of benchmark functions, pre-create
// and sort a large set of random values.
var (
	testIntVals       []int
	testIntValsSorted []int
	uniformIntVals    []int
)

// variables used by the benchmark methods.
var (
	// Constant benchmark variables
	bmConstantHashTableLookupMap map[int]int
	bmConstantLinkedList         *collection.LinkedList[int]
	bmConstantStack              *constant.DynamicStack
	bmConstantQueue              *constant.Queue

	/*
		// Log-log benchmark variables
		bmLogLogRadixSort []int

		// Logarthmic benchmark variables
		bmLogarithmicBST *logarithmic.TreeNode

		// Polylogarithmic benchmark variables
		bmPolylogarithmicRangeTree *polylogarithmic.RangeTree2D
	*/
	// Linear benchmark variables
	bmLinearBST *tree.BSTNode
	/*
		// NLog*N benchmark variables

		// Linearithmic benchmark variables
		bmLinearithmicBuildHeapFromArray []int
		bmLinearithmicHeapifyArray       []int
		bmLinearithmicMergeSort          []int
		bmLinearithmicIntroSort          []int
		bmLinearithmicHeapSort           []int
		bmLinearithmicQuickSort          []int

		// Quadratic benchmark variables
		bmQuadraticBubbleSort      []int
		bmQuadraticInsertionSort   []int
		bmQuadraticSelectionSort   []int
		bmQuadraticNaiveMatrixA    [][]int
		bmQuadraticNaiveMatrixB    [][]int
		bmQuadraticMatrixTranspose [][]int

		// Cubic benchmark variables
		bmCubicFloydWarshallGraph            [][]int
		bmCubicMatrixChainMultiplication     []int
		bmCubicOptimalBSTKeys                []int
		bmCubicOptimalBSTFreq                []int
		bmCubicStandardMatrixMultiplicationA [][]int
		bmCubicStandardMatrixMultiplicationB [][]int

		// Polynomial benchmark variables
		bmPolynomialEditDistanceS1             string
		bmPolynomialEditDistanceS2             string
		bmPolynomialLongestCommonSubsequenceS1 string
		bmPolynomialLongestCommonSubsequenceS2 string
		bmPolynomialLCSWithSequenceS1          string
		bmPolynomialLCSWithSequenceS2          string
		bmPolynomialMatrixChainOrderDimensions []int

		// Exponential benchmark variables
		bmExponentialTSPBitMaskDistances [][]int

		// Factorial benchmark variables.
		bmFactorialAssignmentMatrix [][]int

		// Hyper-exponential benchmark variables.
	*/
)

const (
	// limit is the maximum amount of random values to pre-generate for
	// the benchmarks to run against.
	limit = 10000000
)

func init() {
	testIntVals = make([]int, limit)
	testIntValsSorted = make([]int, limit)
	uniformIntVals = make([]int, limit)

	for i := range limit {
		testIntVals[i] = rand.Int()
		uniformIntVals[i] = i * 3
	}

	copy(testIntValsSorted, testIntVals)

	sort.Ints(testIntValsSorted)
}

// bigOFuncRunner is a helper type to make benchmarking the various example
// functions easier. Since the varuious functions have different call signatures
// (and sometimes some pre-calculations), this lets us wrap them up nicely.
type bigOFuncRunner func(n int, vals []int)

// BenchmarkSettings contains the settings for running one of the
// various example methods as a Benchmark.
type BenchmarkSettings struct {
	ExpectedBigO *bigo.BigO
	Sorted       bool
	Runner       bigOFuncRunner
	Start        int
	End          int
	Step         int

	// Setup is a function that is called before the main benchmark loops
	// starts running. Use this to build trees, load hash maps, reset or init
	// any other data structures that need to be created before the benchmark.
	//
	// This function does NOT need to include calls to
	// StopTimer/StartTimer/ResetTimer as the benchmark runner handles that.
	Setup func(b *testing.B, n int, vals []int)

	// Cleanup is a function that is called after the main benchmark loops
	// has finished running. Use this to clean up any data structures that
	// were created before the benchmark.
	//
	// This function does NOT need to include calls to
	// StopTimer/StartTimer/ResetTimer as the benchmark runner handles that.
	Cleanup func(b *testing.B)

	// TODO(rsned): Add Best/Avg/Worst case identifiers to enable testing the
	// given benchmark with different levels.  E.g., Searching a Slice for a
	// non-existent value will be Worst case because it must do a complete scan
	// of the entire slice of N elements. Binary Searching for the median value
	// in the input set is going to end up closer to O(1) because it's getting
	// lucky and finding the value on the first lookup.
}

var (
	// constantTimeBenchmarks contains O(1) constant time benchmarks
	constantTimeBenchmarks = map[string]BenchmarkSettings{
		"ArrayAccessByIndex": {
			ExpectedBigO: bigo.Constant,
			Sorted:       false,
			Runner:       func(n int, vals []int) { _ = constant.ArrayAccessByIndex(vals, n/2) },
			Start:        100000,
			End:          1000000,
			Step:         100000,
			Setup:        nil,
			Cleanup:      nil,
		},
		"HashTableLookup": {
			ExpectedBigO: bigo.Constant,
			Sorted:       false,
			Runner:       func(n int, _ []int) { _, _ = constant.HashTableLookup(bmConstantHashTableLookupMap, n/2) },
			Start:        100000,
			End:          1000000,
			Step:         100000,
			Setup: func(b *testing.B, n int, vals []int) {
				b.Helper()
				// Fill the hash table with the values from the array.
				bmConstantHashTableLookupMap = map[int]int{}
				for _, v := range vals[:n] {
					bmConstantHashTableLookupMap[v] = v
				}
			},
			Cleanup: func(_ *testing.B) {
				bmConstantHashTableLookupMap = nil
			},
		},
		"BasicMathAdd": {
			ExpectedBigO: bigo.Constant,
			Sorted:       false,
			Runner:       func(_ int, vals []int) { _ = constant.Add(vals[0], vals[1]) },
			Start:        100000,
			End:          1000000,
			Step:         100000,
			Setup:        nil,
			Cleanup:      nil,
		},
		"BasicMathSubtract": {
			ExpectedBigO: bigo.Constant,
			Sorted:       false,
			Runner:       func(_ int, vals []int) { _ = constant.Subtract(vals[0], vals[1]) },
			Start:        100000,
			End:          1000000,
			Step:         100000,
			Setup:        nil,
			Cleanup:      nil,
		},
		"BasicMathMultiply": {
			ExpectedBigO: bigo.Constant,
			Sorted:       false,
			Runner:       func(_ int, vals []int) { _ = constant.Multiply(vals[0], vals[1]) },
			Start:        100000,
			End:          1000000,
			Step:         100000,
			Setup:        nil,
			Cleanup:      nil,
		},
		"BasicMathDivide": {
			ExpectedBigO: bigo.Constant,
			Sorted:       false,
			Runner:       func(_ int, vals []int) { _ = constant.Divide(vals[0], vals[1]) },
			Start:        100000,
			End:          1000000,
			Step:         100000,
			Setup:        nil,
			Cleanup:      nil,
		},
		"LinkedListAccessFirst": {
			ExpectedBigO: bigo.Constant,
			Sorted:       false,
			Runner:       func(_ int, _ []int) { _, _ = bmConstantLinkedList.Front() },
			Start:        100000,
			End:          1000000,
			Step:         100000,
			Setup: func(b *testing.B, n int, vals []int) {
				b.Helper()
				bmConstantLinkedList = collection.FromSlice(vals[:n])
			},
			Cleanup: func(_ *testing.B) {
				bmConstantLinkedList = nil
			},
		},
		"LinkedListAccessLast": {
			ExpectedBigO: bigo.Constant,
			Sorted:       false,
			Runner:       func(_ int, _ []int) { _, _ = bmConstantLinkedList.Back() },
			Start:        100000,
			End:          1000000,
			Step:         100000,
			Setup: func(b *testing.B, n int, vals []int) {
				b.Helper()
				bmConstantLinkedList = collection.FromSlice(vals[:n])
			},
			Cleanup: func(_ *testing.B) {
				bmConstantLinkedList = nil
			},
		},
		"StackPush": {
			ExpectedBigO: bigo.Constant,
			Sorted:       false,
			Runner:       func(n int, vals []int) { bmConstantStack.Push(vals[n/2]) },
			Start:        100000,
			End:          1000000,
			Step:         100000,
			Setup: func(b *testing.B, _ int, _ []int) {
				b.Helper()
				bmConstantStack = &constant.DynamicStack{}
			},
			Cleanup: func(_ *testing.B) {
				bmConstantStack = nil
			},
		},
		"StackPop": {
			ExpectedBigO: bigo.Constant,
			Sorted:       false,
			Runner:       func(_ int, _ []int) { bmConstantStack.Pop() },
			Start:        100000,
			End:          1000000,
			Step:         100000,
			Setup: func(b *testing.B, n int, vals []int) {
				b.Helper()
				bmConstantStack = &constant.DynamicStack{}
				for _, v := range vals[:n] {
					bmConstantStack.Push(v)
				}
			},
			Cleanup: func(_ *testing.B) {
				bmConstantStack = nil
			},
		},
		"QueueEnqueue": {
			ExpectedBigO: bigo.Constant,
			Sorted:       false,
			Runner:       func(n int, vals []int) { bmConstantQueue.Enqueue(vals[n/2]) },
			Start:        100000,
			End:          1000000,
			Step:         100000,
			Setup: func(b *testing.B, _ int, _ []int) {
				b.Helper()
				bmConstantQueue = &constant.Queue{}
			},
			Cleanup: func(_ *testing.B) {
				bmConstantQueue = nil
			},
		},
		"QueueDequeue": {
			ExpectedBigO: bigo.Constant,
			Sorted:       false,
			Runner:       func(_ int, _ []int) { bmConstantQueue.Dequeue() },
			Start:        100000,
			End:          1000000,
			Step:         100000,
			Setup: func(b *testing.B, n int, vals []int) {
				b.Helper()
				bmConstantQueue = &constant.Queue{}
				for _, v := range vals[:n] {
					bmConstantQueue.Enqueue(v)
				}

				b.ResetTimer()
			},
			Cleanup: func(_ *testing.B) {
				bmConstantQueue = nil
			},
		},
	}

	// loglogTimeBenchmarks contains O(log(log n)) benchmarks
	loglogTimeBenchmarks = map[string]BenchmarkSettings{
		/*
			"InterpolationSearch": {
				ExpectedBigO: bigo.LogLog,
				Sorted:       true,
				Runner: func(n int, _ []int) {
					// Create a non-uniform sorted array to prevent O(1) interpolation guesses
					// Use a quadratic spacing to make interpolation search work harder
					nonUniformArr := make([]int, n)
					for i := range n {
						// Quadratic distribution: values grow as i^2
						nonUniformArr[i] = i * i
					}

					// Perform multiple searches with targets that don't align perfectly
					logLogN := max(int(math.Log2(math.Log2(float64(n))))+1, 2)

					for i := range logLogN {
						// Search for values between the quadratic points to force iterations
						targetIndex := (n / (logLogN + 1)) * (i + 1)
						target := targetIndex*targetIndex + i // Offset from perfect square
						_ = loglog.InterpolationSearch(nonUniformArr, target)
					}
				},
				Start:   100000,
				End:     1000000,
				Step:    100000,
				Setup:   nil,
				Cleanup: nil,
			},
			"YFastTrieOperations": {
				ExpectedBigO: bigo.LogLog,
				Sorted:       false,
				Runner: func(n int, _ []int) {
					_ = loglog.YFastTrieOperations(n)
				},
				Start:   100,
				End:     10000,
				Step:    1000,
				Setup:   nil,
				Cleanup: nil,
			},
		*/
	}

	// logarithmicTimeBenchmarks contains O(log n) benchmarks
	logarithmicTimeBenchmarks = map[string]BenchmarkSettings{
		/*
			"BinarySearch": {
				ExpectedBigO: bigo.Log,
				Sorted:       true,
				Runner: func(n int, vals []int) {
					// Search for value larger than any in array to force
					// worst-case O(log n).
					target := vals[len(vals)-1] + 1
					_ = logarithmic.BinarySearch(vals[:n], target)
				},
				Start:   10000,
				End:     1000000,
				Step:    100000,
				Setup:   nil,
				Cleanup: nil,
			},
			"BinaryTreeSearch": {
				ExpectedBigO: bigo.Log,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					// Perform multiple searches to ensure we get consistent O(log n) behavior
					// Search for various targets to traverse different paths in the tree
					searchCount := int(math.Log2(float64(n))) + 1
					for i := 0; i < searchCount; i++ {
						// Search for values that force traversal of the tree
						target := vals[(i*n/searchCount)%len(vals)] + i + 1
						_ = logarithmic.BinaryTreeSearch(bmLogarithmicBST, target)
					}
				},
				Start: 250,
				End:   10000,
				Step:  250,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					bmLogarithmicBST = logarithmic.BuildBST(vals[:n])
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmLogarithmicBST = nil
				},
			},
			"SegmentTreeOperations": {
				ExpectedBigO: bigo.Log,
				Sorted:       false,
				Runner: func(n int, _ []int) {
					_ = logarithmic.SegmentTreeOperations(n)
				},
				Start: 100,
				End:   10000,
				Step:  1000,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					// Pre-create segment tree to avoid O(n) construction overhead
					arr := make([]int, n)
					for i := 0; i < n; i++ {
						arr[i] = vals[i%len(vals)]
					}
					logarithmic.GlobalSegmentTree = logarithmic.NewSegmentTree(arr)
				},
				Cleanup: func(_ *testing.B) {
					logarithmic.GlobalSegmentTree = nil
				},
			},
		*/
	}

	// polylogarithmicTimeBenchmarks contains O((log n)^c) benchmarks
	polylogarithmicTimeBenchmarks = map[string]BenchmarkSettings{
		/*
			"RangeTree2D_Build": {
				ExpectedBigO: bigo.Polylogarithmic,
				Sorted:       false,
				Runner: func(n int, _ []int) {
					// Generate points for 2D range tree
					points := make([]polylogarithmic.Point2D, n)
					for i := 0; i < n; i++ {
						points[i] = polylogarithmic.Point2D{X: i, Y: i * 2}
					}
					_ = polylogarithmic.BuildRangeTree2D(points)
				},
				Start:   100,
				End:     10000,
				Step:    1000,
				Setup:   nil,
				Cleanup: nil,
			},
			"RangeTree2D_Query": {
				ExpectedBigO: bigo.Polylogarithmic,
				Sorted:       false,
				Runner: func(n int, _ []int) {
					// Use global tree that should be pre-built in Setup
					if globalRangeTree != nil {
						// Perform multiple queries to amplify the O((log n)^2) behavior
						numQueries := int(math.Log2(float64(n))) * int(math.Log2(float64(n)))
						for i := 0; i < numQueries; i++ {
							querySize := n / (4 + i%3) // Vary query ranges
							_ = globalRangeTree.Query2D(querySize, 3*querySize, querySize*2, 3*querySize*2)
						}
					}
				},
				Start: 1000,
				End:   100000,
				Step:  10000,
				Setup: func(b *testing.B, n int, _ []int) {
					b.Helper()
					b.StopTimer()
					// Pre-build tree once for all iterations
					points := make([]polylogarithmic.Point2D, n)
					for i := 0; i < n; i++ {
						points[i] = polylogarithmic.Point2D{X: i, Y: i * 2}
					}
					globalRangeTree = polylogarithmic.BuildRangeTree2D(points)
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					globalRangeTree = nil
				},
			},
			"FractionalCascadingSearch": {
				ExpectedBigO: bigo.Polylogarithmic,
				Sorted:       false,
				Runner: func(n int, _ []int) {
					// Create log(n) sorted lists to demonstrate O((log n)^2) complexity
					logN := int(math.Log2(float64(n))) + 1
					if logN < 2 {
						logN = 2
					}
					listSize := n / logN
					if listSize < 1 {
						listSize = 1
					}
					sortedLists := make([][]int, logN)
					for i := 0; i < logN; i++ {
						sortedLists[i] = make([]int, listSize)
						for j := 0; j < listSize; j++ {
							sortedLists[i][j] = i*listSize + j
						}
					}

					// Perform log(n) searches to get O((log n)^2) total complexity
					for searchCount := 0; searchCount < logN; searchCount++ {
						target := (n / 2) + searchCount
						_ = polylogarithmic.FractionalCascadingSearch(sortedLists, target)
					}
				},
				Start:   1000,
				End:     100000,
				Step:    10000,
				Setup:   nil,
				Cleanup: nil,
			},
		*/
	}

	// linearTimeBenchmarks contains O(n) benchmarks
	linearTimeBenchmarks = map[string]BenchmarkSettings{
		"Search": {
			ExpectedBigO: bigo.Linear,
			Sorted:       false,
			// Search for hopefully non-existent value to force worst-case O(n)
			Runner:  func(n int, vals []int) { _ = linear.Search(vals[:n], math.MaxInt) },
			Start:   10000,
			End:     100000,
			Step:    10000,
			Setup:   nil,
			Cleanup: nil,
		},
		"ArrayTraversal": {
			ExpectedBigO: bigo.Linear,
			Sorted:       false,
			Runner: func(n int, vals []int) {
				// Use in-place traversal to avoid memory allocation overhead
				for i, val := range vals[:n] {
					vals[i] = val * 2
				}
			},
			Start:   10000,
			End:     100000,
			Step:    10000,
			Setup:   nil,
			Cleanup: nil,
		},
		"CountElements": {
			ExpectedBigO: bigo.Linear,
			Sorted:       false,
			Runner:       func(n int, vals []int) { _ = linear.CountElements(vals[:n]) },
			Start:        10000,
			End:          100000,
			Step:         10000,
			Setup:        nil,
			Cleanup:      nil,
		},
		"FindMinimum": {
			ExpectedBigO: bigo.Linear,
			Sorted:       false,
			Runner:       func(n int, vals []int) { _, _ = linear.FindMinimum(vals[:n]) },
			Start:        10000,
			End:          100000,
			Step:         10000,
			Setup:        nil,
			Cleanup:      nil,
		},
		"FindMaximum": {
			ExpectedBigO: bigo.Linear,
			Sorted:       false,
			Runner:       func(n int, vals []int) { _, _ = linear.FindMaximum(vals[:n]) },
			Start:        10000,
			End:          100000,
			Step:         10000,
			Setup:        nil,
			Cleanup:      nil,
		},
		"CalculateSum": {
			ExpectedBigO: bigo.Linear,
			Sorted:       false,
			Runner:       func(n int, vals []int) { _ = linear.CalculateSum(vals[:n]) },
			Start:        10000,
			End:          100000,
			Step:         10000,
			Setup:        nil,
			Cleanup:      nil,
		},
		"FindTreeHeight": {
			// FindTreeHeight visits all nodes, so it's O(n) not O(log n)
			ExpectedBigO: bigo.Linear,
			Sorted:       false,
			Runner: func(_ int, _ []int) {
				// TODO(rsned): Move from logarithmic to linear package.
				_ = linear.FindTreeHeight(bmLinearBST)
			},
			Start: 1000,
			End:   10000,
			Step:  1000,
			Setup: func(b *testing.B, n int, vals []int) {
				b.Helper()
				b.StopTimer()
				bmLinearBST = tree.BuildBST(vals[:n])
				b.StartTimer()
			},
			Cleanup: func(_ *testing.B) {
				bmLinearBST = nil
			},
		},
		"ParallelDivideConquer": {
			ExpectedBigO: bigo.Linear,
			Sorted:       false,
			Runner:       func(n int, vals []int) { _ = loglog.ParallelDivideConquer(vals[:n]) },
			Start:        10000,
			End:          100000,
			Step:         10000,
			Setup:        nil,
			Cleanup:      nil,
		},
	}

	// nLogStarNTimeBenchmarks contains O(n log*(n)) benchmarks
	nLogStarNTimeBenchmarks = map[string]BenchmarkSettings{
		/*
			"UnionFindOperations": {
				ExpectedBigO: bigo.NLogStarN,
				Sorted:       false,
				Runner: func(n int, _ []int) {
					_ = nlogstar.PerformUnionFindOperations(n)
				},
				Start:   100,
				End:     1000000,
				Step:    50000,
				Setup:   nil,
				Cleanup: nil,
			},
			"KruskalMST": {
				ExpectedBigO: bigo.NLogStarN,
				Sorted:       false,
				Runner: func(n int, _ []int) {
					// Create linear number of edges O(n) to focus on Union-Find complexity
					var edges []nlogstar.Edge
					weight := 1

					// Create a connected graph with exactly 2n-1 edges (linear in n)
					// This ensures sorting is O(n log n) and Union-Find dominates at O(n log* n)
					for i := 0; i < n-1; i++ {
						// Chain edges to ensure connectivity
						edges = append(edges, nlogstar.Edge{From: i, To: i + 1, Weight: weight})
						weight++
					}

					// Add n additional random edges to exercise Union-Find more
					for i := 0; i < n && len(edges) < 2*n-1; i++ {
						to := (i + n/2) % n
						if to != i {
							edges = append(edges, nlogstar.Edge{From: i, To: to, Weight: weight})
							weight++
						}
					}

					_ = nlogstar.KruskalMST(edges, n)
				},
				Start:   100,
				End:     100000,
				Step:    10000,
				Setup:   nil,
				Cleanup: nil,
			},
			"NetworkConnectivity": {
				ExpectedBigO: bigo.NLogStarN,
				Sorted:       false,
				Runner: func(n int, _ []int) {
					_ = nlogstar.SimulateNetworkConnectivity(n)
				},
				Start:   100,
				End:     1000000,
				Step:    50000,
				Setup:   nil,
				Cleanup: nil,
			},
		*/
	}

	// linearithmicTimeBenchmarks contains O(n log n) benchmarks
	linearithmicTimeBenchmarks = map[string]BenchmarkSettings{
		/*
			"MergeSort": {
				ExpectedBigO: bigo.Linearithmic,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					// TODO(rsned): Because this is a sort, do we need to re-copy
					// the slice inside the inner benchmark loop? After the first
					// iteration, the slice is sorted and the benchmark will run
					// more times on the already sorted data?
					_ = linearithmic.MergeSort(bmLinearithmicMergeSort)
				},
				Start: 100,
				End:   1000000,
				Step:  50000,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					// Copy the original slice because we are going to sort the
					// slice we pass in and don't want to ruin anyone elses tests.
					bmLinearithmicMergeSort = make([]int, n)
					copy(bmLinearithmicMergeSort, vals[:n])
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmLinearithmicMergeSort = nil
				},
			},
			"BuildHeapFromArray": {
				ExpectedBigO: bigo.Linearithmic,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = linearithmic.BuildHeapFromArray(bmLinearithmicBuildHeapFromArray)
				},
				Start: 100,
				End:   1000000,
				Step:  50000,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					bmLinearithmicBuildHeapFromArray = make([]int, n)
					copy(bmLinearithmicBuildHeapFromArray, vals[:n])
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmLinearithmicBuildHeapFromArray = nil
				},
			},
			"HeapifyArray": {
				ExpectedBigO: bigo.Linearithmic,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = linearithmic.HeapifyArray(bmLinearithmicHeapifyArray)
				},
				Start: 100,
				End:   1000000,
				Step:  50000,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					bmLinearithmicHeapifyArray = make([]int, n)
					copy(bmLinearithmicHeapifyArray, vals[:n])
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmLinearithmicHeapifyArray = nil
				},
			},
			"IntroSort": {
				ExpectedBigO: bigo.Linearithmic,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = linearithmic.IntroSort(bmLinearithmicIntroSort)
				},
				Start: 100,
				End:   1000000,
				Step:  50000,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					bmLinearithmicIntroSort = make([]int, n)
					copy(bmLinearithmicIntroSort, vals[:n])
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmLinearithmicIntroSort = nil
				},
			},
			"HeapSort": {
				ExpectedBigO: bigo.Linearithmic,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = linearithmic.HeapSort(bmLinearithmicHeapSort)
				},
				Start: 100,
				End:   1000000,
				Step:  50000,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					bmLinearithmicHeapSort = make([]int, n)
					copy(bmLinearithmicHeapSort, vals[:n])
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmLinearithmicHeapSort = nil
				},
			},
			"QuickSort": {
				ExpectedBigO: bigo.Linearithmic,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = linearithmic.QuickSort(bmLinearithmicQuickSort)
				},
				Start: 100,
				End:   1000000,
				Step:  50000,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					bmLinearithmicQuickSort = make([]int, n)
					copy(bmLinearithmicQuickSort, vals[:n])
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmLinearithmicQuickSort = nil
				},
			},
			"RadixSortOptimization": {
				ExpectedBigO: bigo.Linearithmic, // Moved from LogLog - actually O(n log n)
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = linearithmic.RadixSortOptimization(bmLogLogRadixSort)
				},
				Start: 100,
				End:   1000000,
				Step:  50000,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					bmLogLogRadixSort = make([]int, n)
					copy(bmLogLogRadixSort, vals[:n])
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmLogLogRadixSort = nil
				},
			},
		*/
	}

	// quadraticTimeBenchmarks contains O(n²) benchmarks
	quadraticTimeBenchmarks = map[string]BenchmarkSettings{
		/*
			"BubbleSort": {
				ExpectedBigO: bigo.Quadratic,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = quadratic.BubbleSort(bmQuadraticBubbleSort)
				},
				Start: 100,
				End:   5000,
				Step:  200,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					bmQuadraticBubbleSort = make([]int, n)
					copy(bmQuadraticBubbleSort, vals[:n])
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmQuadraticBubbleSort = nil
				},
			},
			"InsertionSort": {
				ExpectedBigO: bigo.Quadratic,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = quadratic.InsertionSort(bmQuadraticInsertionSort)
				},
				Start: 100,
				End:   5000,
				Step:  200,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					bmQuadraticInsertionSort = make([]int, n)
					copy(bmQuadraticInsertionSort, vals[:n])
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmQuadraticInsertionSort = nil
				},
			},
			"SelectionSort": {
				ExpectedBigO: bigo.Quadratic,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = quadratic.SelectionSort(bmQuadraticSelectionSort)
				},
				Start: 100,
				End:   5000,
				Step:  200,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					bmQuadraticSelectionSort = make([]int, n)
					copy(bmQuadraticSelectionSort, vals[:n])
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmQuadraticSelectionSort = nil
				},
			},
			"AllPairsComparison": {
				ExpectedBigO: bigo.Quadratic,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					_ = quadratic.AllPairsComparison(vals[:n])
				},
				Start:   500,
				End:     5000,
				Step:    500,
				Setup:   nil,
				Cleanup: nil,
			},
			"FindDuplicatePairs": {
				ExpectedBigO: bigo.Quadratic,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					_ = quadratic.FindDuplicatePairs(vals[:n])
				},
				Start:   500,
				End:     10000,
				Step:    500,
				Setup:   nil,
				Cleanup: nil,
			},
			"CountInversions": {
				ExpectedBigO: bigo.Quadratic,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					_ = quadratic.CountInversions(vals[:n])
				},
				Start:   1000,
				End:     10000,
				Step:    1000,
				Setup:   nil,
				Cleanup: nil,
			},
			"TwoSum": {
				ExpectedBigO: bigo.Quadratic,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					target := vals[0] + vals[n/2]
					_ = quadratic.TwoSum(vals[:n], target)
				},
				Start:   1000,
				End:     10000,
				Step:    1000,
				Setup:   nil,
				Cleanup: nil,
			},
			"MatrixTranspose": {
				ExpectedBigO: bigo.Quadratic,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = quadratic.MatrixTranspose(bmQuadraticMatrixTranspose)
				},
				Start: 500,
				End:   5000,
				Step:  500,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					bmQuadraticMatrixTranspose = make([][]int, n)
					for i := range bmQuadraticMatrixTranspose {
						bmQuadraticMatrixTranspose[i] = make([]int, n)
						for j := range bmQuadraticMatrixTranspose[i] {
							bmQuadraticMatrixTranspose[i][j] = vals[(i*n+j)%len(vals)] % 100
						}
					}
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmQuadraticMatrixTranspose = nil
				},
			},
		*/
	}

	// cubicTimeBenchmarks contains O(n³) benchmarks
	cubicTimeBenchmarks = map[string]BenchmarkSettings{
		/*
			"FloydWarshall": {
				ExpectedBigO: bigo.Cubic,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = cubic.FloydWarshall(bmCubicFloydWarshallGraph)
				},
				Start: 10,
				End:   100,
				Step:  10,
				Setup: func(b *testing.B, n int, _ []int) {
					b.Helper()
					b.StopTimer()
					bmCubicFloydWarshallGraph = make([][]int, n)
					for i := range bmCubicFloydWarshallGraph {
						bmCubicFloydWarshallGraph[i] = make([]int, n)
						for j := range bmCubicFloydWarshallGraph[i] {
							if i == j {
								bmCubicFloydWarshallGraph[i][j] = 0
							} else {
								bmCubicFloydWarshallGraph[i][j] = 99999
							}
						}
					}
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmCubicFloydWarshallGraph = nil
				},
			},
			"ThreeSum": {
				ExpectedBigO: bigo.Cubic,
				Sorted:       false,
				Runner: func(_ int, vals []int) {
					target := vals[0] + vals[1] + vals[2] // Use first three values as target
					_ = cubic.ThreeSum(vals, target)
				},
				Start:   100,
				End:     1000,
				Step:    100,
				Setup:   nil,
				Cleanup: nil,
			},
			"NaiveMatrixMultiplication": {
				ExpectedBigO: bigo.Cubic,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = quadratic.NaiveMatrixMultiplication(bmQuadraticNaiveMatrixA, bmQuadraticNaiveMatrixB)
				},
				Start: 3,
				End:   20,
				Step:  1,
				Setup: func(b *testing.B, n int, _ []int) {
					b.Helper()
					b.StopTimer()
					bmQuadraticNaiveMatrixA = make([][]int, n)
					bmQuadraticNaiveMatrixB = make([][]int, n)
					for i := range bmQuadraticNaiveMatrixA {
						bmQuadraticNaiveMatrixA[i] = make([]int, n)
						bmQuadraticNaiveMatrixB[i] = make([]int, n)
						for j := range bmQuadraticNaiveMatrixA[i] {
							bmQuadraticNaiveMatrixA[i][j] = i*n + j + 1 // vals[rand.Intn((i+1)*n)%n] % 1000
							bmQuadraticNaiveMatrixB[i][j] = j*n + i + 1 // vals[rand.Intn((j+1)*n)%n] % 1000
						}
					}
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmQuadraticNaiveMatrixA = nil
					bmQuadraticNaiveMatrixB = nil
				},
			},
			"MatrixChainMultiplication": {
				ExpectedBigO: bigo.Cubic,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = cubic.MatrixChainMultiplication(bmCubicMatrixChainMultiplication)
				},
				Start: 50,
				End:   750,
				Step:  50,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					bmCubicMatrixChainMultiplication = make([]int, n+1)
					for i := range bmCubicMatrixChainMultiplication {
						bmCubicMatrixChainMultiplication[i] = (vals[i%len(vals)] % 50) + 1 // Positive dimensions
					}
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmCubicMatrixChainMultiplication = nil
				},
			},
			"OptimalBinarySearchTree": {
				ExpectedBigO: bigo.Cubic,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = cubic.OptimalBinarySearchTree(bmCubicOptimalBSTKeys, bmCubicOptimalBSTFreq)
				},
				Start: 50,
				End:   500,
				Step:  50,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					bmCubicOptimalBSTKeys = make([]int, n)
					bmCubicOptimalBSTFreq = make([]int, n)
					for i := range bmCubicOptimalBSTKeys {
						bmCubicOptimalBSTKeys[i] = vals[i%len(vals)]
						bmCubicOptimalBSTFreq[i] = (vals[(i+n)%len(vals)] % 10) + 1 // Positive frequencies
					}
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmCubicOptimalBSTKeys = nil
					bmCubicOptimalBSTFreq = nil
				},
			},
			"StandardMatrixMultiplication": {
				ExpectedBigO: bigo.Cubic,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = cubic.StandardMatrixMultiplication(bmCubicStandardMatrixMultiplicationA, bmCubicStandardMatrixMultiplicationB)
				},
				Start: 100,
				End:   1000,
				Step:  100,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					bmCubicStandardMatrixMultiplicationA = make([][]int, n)
					bmCubicStandardMatrixMultiplicationB = make([][]int, n)
					for i := range bmCubicStandardMatrixMultiplicationA {
						bmCubicStandardMatrixMultiplicationA[i] = make([]int, n)
						bmCubicStandardMatrixMultiplicationB[i] = make([]int, n)
						for j := range bmCubicStandardMatrixMultiplicationA[i] {
							bmCubicStandardMatrixMultiplicationA[i][j] = vals[(i*n+j)%len(vals)]
							bmCubicStandardMatrixMultiplicationB[i][j] = vals[(i*n+j+n*n)%len(vals)]
						}
					}
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmCubicStandardMatrixMultiplicationA = nil
					bmCubicStandardMatrixMultiplicationB = nil
				},
			},
			"TripleNestedProductSum": {
				ExpectedBigO: bigo.Cubic,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					_ = cubic.TripleNestedProductSum(vals[:n])
				},
				Start:   250,
				End:     2500,
				Step:    250,
				Setup:   nil,
				Cleanup: nil,
			},
			"FindTripletsWithSum": {
				ExpectedBigO: bigo.Cubic,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					if n < 3 {
						panic("n must be at least 3")
					}
					targetSum := vals[0] + vals[1] + vals[2] // Use first three as target
					_ = cubic.FindTripletsWithSum(vals[:n], targetSum)
				},
				Start:   250,
				End:     2500,
				Step:    250,
				Setup:   nil,
				Cleanup: nil,
			},
			"CountTripletsWithProperty": {
				ExpectedBigO: bigo.Cubic,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					_ = cubic.CountTripletsWithProperty(vals[:n])
				},
				Start:   250,
				End:     2500,
				Step:    250,
				Setup:   nil,
				Cleanup: nil,
			},
			"Generate3DCombinations": {
				ExpectedBigO: bigo.Cubic,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					_ = cubic.Generate3DCombinations(vals[:n])
				},
				Start:   5,
				End:     50,
				Step:    5,
				Setup:   nil,
				Cleanup: nil,
			},
		*/
	}

	// exponentialTimeBenchmarks contains O(2^n) benchmarks
	exponentialTimeBenchmarks = map[string]BenchmarkSettings{
		/*
			"RecursiveFibonacci": {
				ExpectedBigO: bigo.Exponential,
				Sorted:       false,
				Runner: func(n int, _ []int) {
					_ = exponential.RecursiveFibonacci(n)
				},
				Start:   20,
				End:     30,
				Step:    1,
				Setup:   nil,
				Cleanup: nil,
			},
			"TowerOfHanoi": {
				ExpectedBigO: bigo.Exponential,
				Sorted:       false,
				Runner: func(n int, _ []int) {
					// Use TowerOfHanoi instead of TowerOfHanoiCount to actually perform exponential work
					_ = exponential.TowerOfHanoi(n, "A", "B", "C")
				},
				Start:   5,
				End:     20,
				Step:    1,
				Setup:   nil,
				Cleanup: nil,
			},
			"GenerateAllSubsets": {
				ExpectedBigO: bigo.Exponential,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					// Add cap to prevent excessive runtime (2^22 = 4M operations)?
					_ = exponential.GenerateAllSubsets(vals[:n])
				},
				Start:   5,
				End:     21,
				Step:    1,
				Setup:   nil,
				Cleanup: nil,
			},
			"TravelingSalesmanBruteForce": {
				ExpectedBigO: bigo.Exponential,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					// Scale city count directly with input for exponential growth
					cityCount := min(n, 10) // Cap to prevent excessive runtime
					distances := make([][]int, cityCount)
					for i := range distances {
						distances[i] = make([]int, cityCount)
						for j := range distances[i] {
							if i == j {
								distances[i][j] = 0
							} else {
								distances[i][j] = vals[(i*cityCount+j)%len(vals)]%100 + 1
							}
						}
					}
					_, _ = exponential.TravelingSalesmanBruteForce(distances)
				},
				Start:   3,
				End:     10,
				Step:    1,
				Setup:   nil,
				Cleanup: nil,
			},
			"TSPBitMask": {
				ExpectedBigO: bigo.Exponential,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = exponential.TSPBitMask(bmExponentialTSPBitMaskDistances)
				},
				Start: 3,
				End:   10,
				Step:  1,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					bmExponentialTSPBitMaskDistances = make([][]int, n)
					for i := range bmExponentialTSPBitMaskDistances {
						bmExponentialTSPBitMaskDistances[i] = make([]int, n)
						for j := range bmExponentialTSPBitMaskDistances[i] {
							if i == j {
								bmExponentialTSPBitMaskDistances[i][j] = 0
							} else {
								bmExponentialTSPBitMaskDistances[i][j] = vals[(i*n+j)%len(vals)]%100 + 1
							}
						}
					}
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmExponentialTSPBitMaskDistances = nil
				},
			},
		*/
	}

	// polynomialTimeBenchmarks contains O(n^c) and other polynomial benchmarks
	polynomialTimeBenchmarks = map[string]BenchmarkSettings{
		/*
			"EditDistance": {
				ExpectedBigO: bigo.Polynomial,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = polynomial.EditDistance(bmPolynomialEditDistanceS1, bmPolynomialEditDistanceS2)
				},
				Start: 20,
				End:   200,
				Step:  20,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					bmPolynomialEditDistanceS1 = string(rune('a'+vals[0]%26)) + string(rune('a'+vals[1]%26))
					bmPolynomialEditDistanceS2 = string(rune('a'+vals[2]%26)) + string(rune('a'+vals[3]%26))
					for i := 4; i < n && i < len(vals); i++ {
						bmPolynomialEditDistanceS1 += string(rune('a' + vals[i]%26))
						if i+1 < len(vals) {
							bmPolynomialEditDistanceS2 += string(rune('a' + vals[i+1]%26))
						}
					}
					b.StartTimer()
				},
				Cleanup: nil,
			},
			"JohnsonAlgorithm": {
				ExpectedBigO: bigo.Polynomial,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					graph := make([][]int, n)
					for i := range graph {
						graph[i] = make([]int, n)
						for j := range graph[i] {
							if i == j {
								graph[i][j] = 0
							} else {
								graph[i][j] = vals[(i*n+j)%len(vals)]%100 + 1
							}
						}
					}
					_ = polynomial.JohnsonAlgorithm(graph)
				},
				Start:   5,
				End:     50,
				Step:    5,
				Setup:   nil,
				Cleanup: nil,
			},
			"LongestCommonSubsequence": {
				ExpectedBigO: bigo.Polynomial,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = polynomial.LongestCommonSubsequence(bmPolynomialLongestCommonSubsequenceS1, bmPolynomialLongestCommonSubsequenceS2)
				},
				Start: 20,
				End:   200,
				Step:  20,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					bmPolynomialLongestCommonSubsequenceS1 = ""
					bmPolynomialLongestCommonSubsequenceS2 = ""
					for i := 4; i < n && i < len(vals); i++ {
						bmPolynomialLongestCommonSubsequenceS1 += string(rune('a' + vals[i]%26))
						if i+1 < len(vals) {
							bmPolynomialLongestCommonSubsequenceS2 += string(rune('a' + vals[i+1]%26))
						}
					}
					b.StartTimer()
				},
				Cleanup: nil,
			},
			"LCSWithSequence": {
				ExpectedBigO: bigo.Polynomial,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_ = polynomial.LCSWithSequence(bmPolynomialLCSWithSequenceS1, bmPolynomialLCSWithSequenceS2)
				},
				Start: 20,
				End:   200,
				Step:  20,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					for i := 0; i < n && i < len(vals); i++ {
						bmPolynomialLCSWithSequenceS1 += string(rune('a' + vals[i]%26))
						if i+1 < len(vals) {
							bmPolynomialLCSWithSequenceS2 += string(rune('a' + vals[i+1]%26))
						}
					}
					b.StartTimer()
				},
				Cleanup: func(_ *testing.B) {
					bmPolynomialLCSWithSequenceS1 = ""
					bmPolynomialLCSWithSequenceS2 = ""
				},
			},
			"MatrixChainOrder": {
				ExpectedBigO: bigo.Polynomial,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_, _ = polynomial.MatrixChainOrder(bmPolynomialMatrixChainOrderDimensions)
				},
				Start: 5,
				End:   50,
				Step:  5,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					b.StopTimer()
					bmPolynomialMatrixChainOrderDimensions = make([]int, n+1)
					for i := range bmPolynomialMatrixChainOrderDimensions {
						bmPolynomialMatrixChainOrderDimensions[i] = vals[i%len(vals)]%100 + 10
					}
					b.StartTimer()
				},
				Cleanup: nil,
			},
			"FordFulkerson": {
				ExpectedBigO: bigo.Polynomial,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					// Create a denser graph that requires more iterations
					// to demonstrate O(V * E * f) behavior where f is max flow
					capacity := make([][]int, n)
					for i := range capacity {
						capacity[i] = make([]int, n)
						for j := range capacity[i] {
							if i != j {
								// Create higher capacity values to force more iterations
								// Dense connectivity with significant flow capacity
								capacity[i][j] = vals[(i*n+j)%len(vals)]%1000 + 100
							}
						}
					}

					// Run multiple flow computations to amplify the polynomial behavior
					for source := 0; source < n/2; source++ {
						for sink := n / 2; sink < n; sink++ {
							_ = polynomial.FordFulkerson(capacity, source, sink)
						}
					}
				},
				Start:   5,
				End:     50,
				Step:    5,
				Setup:   nil,
				Cleanup: nil,
			},
			"EdmondsKarp": {
				ExpectedBigO: bigo.Polynomial,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					// Create a denser graph to demonstrate O(V * E^2) behavior
					capacity := make([][]int, n)
					for i := range capacity {
						capacity[i] = make([]int, n)
						for j := range capacity[i] {
							if i != j {
								// Dense connectivity with higher capacity values
								capacity[i][j] = vals[(i*n+j)%len(vals)]%1000 + 100
							}
						}
					}

					// Run multiple flow computations to amplify the O(V * E^2) behavior
					for source := 0; source < n/2; source++ {
						for sink := n / 2; sink < n; sink++ {
							_ = polynomial.EdmondsKarp(capacity, source, sink)
						}
					}
				},
				Start:   5,
				End:     50,
				Step:    5,
				Setup:   nil,
				Cleanup: nil,
			},
		*/
	}

	// factorialTimeBenchmarks contains O(n!) benchmarks
	factorialTimeBenchmarks = map[string]BenchmarkSettings{
		/*
			"GenerateAllPermutations": {
				ExpectedBigO: bigo.Factorial,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					_ = factorial.GenerateAllPermutations(vals[:n])
				},
				Start:   1,
				End:     10,
				Step:    1,
				Setup:   nil,
				Cleanup: nil,
			},
			"AssignmentProblemBruteForce": {
				ExpectedBigO: bigo.Factorial,
				Sorted:       false,
				Runner: func(_ int, _ []int) {
					_, _ = factorial.AssignmentProblemBruteForce(bmFactorialAssignmentMatrix)
				},
				Start: 2,
				End:   10,
				Step:  1,
				Setup: func(b *testing.B, n int, vals []int) {
					b.Helper()
					size := min(1+n, 25) // Sizes 2 - n
					bmFactorialAssignmentMatrix = make([][]int, size)
					for i := range bmFactorialAssignmentMatrix {
						bmFactorialAssignmentMatrix[i] = make([]int, size)
						for j := range bmFactorialAssignmentMatrix[i] {
							bmFactorialAssignmentMatrix[i][j] = vals[(i*size+j)%len(vals)]%100 + 1
						}
					}
					b.ResetTimer()
				},
				Cleanup: func(b *testing.B) {
					b.Helper()
					bmFactorialAssignmentMatrix = nil
					b.StopTimer()
				},
			},
			"OptimalMatchingBruteForce": {
				ExpectedBigO: bigo.Factorial,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					size := 2 + n // Sizes 2-5
					if size > 25 {
						size = 25
					}
					weights := make([][]int, size)
					for i := range weights {
						weights[i] = make([]int, size)
						for j := range weights[i] {
							weights[i][j] = vals[(i*size+j)%len(vals)]%100 + 1
						}
					}
					_, _ = factorial.OptimalMatchingBruteForce(weights)
				},
				Start:   1,
				End:     8,
				Step:    1,
				Setup:   nil,
				Cleanup: nil,
			},
			"NQueensAllArrangements": {
				ExpectedBigO: bigo.Factorial,
				Sorted:       false,
				Runner: func(n int, _ []int) {
					queenN := min(n, 8) // N-Queens sizes directly match n
					_ = factorial.NQueensAllArrangements(queenN)
				},
				Start:   1,
				End:     6, // Going beyond this puts really puts the hurt on your machine.
				Step:    1,
				Setup:   nil,
				Cleanup: nil,
			},
			"NQueensCountAllArrangements": {
				ExpectedBigO: bigo.Factorial,
				Sorted:       false,
				Runner: func(n int, _ []int) {
					queenN := min(n, 8) // N-Queens sizes directly match n
					_ = factorial.NQueensCountAllArrangements(queenN)
				},
				Start:   1,
				End:     6, // Going beyond this puts really puts the hurt on your machine.
				Step:    1,
				Setup:   nil,
				Cleanup: nil,
			},
			"GenerateAllSchedules": {
				ExpectedBigO: bigo.Factorial,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					taskCount := n // Task counts directly match n
					if taskCount > 8 {
						taskCount = 8
					}
					tasks := make([]factorial.Task, taskCount)
					for i := range tasks {
						tasks[i] = factorial.Task{ID: i + 1, Duration: vals[i%len(vals)]%10 + 1, Priority: vals[(i+1)%len(vals)]%5 + 1}
					}
					_ = factorial.GenerateAllSchedules(tasks)
				},
				Start:   1,
				End:     8,
				Step:    1,
				Setup:   nil,
				Cleanup: nil,
			},
			"OptimalScheduleBruteForce": {
				ExpectedBigO: bigo.Factorial,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					taskCount := n // Task counts directly match n
					if taskCount > 8 {
						taskCount = 8
					}
					tasks := make([]factorial.Task, taskCount)
					for i := range tasks {
						tasks[i] = factorial.Task{ID: i + 1, Duration: vals[i%len(vals)]%10 + 1, Priority: vals[(i+1)%len(vals)]%5 + 1}
					}
					_, _ = factorial.OptimalScheduleBruteForce(tasks)
				},
				Start:   1,
				End:     8,
				Step:    1,
				Setup:   nil,
				Cleanup: nil,
			},
			"JobShopSchedulingBruteForce": {
				ExpectedBigO: bigo.Factorial,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					jobCount := n // Job counts directly match n
					if jobCount > 8 {
						jobCount = 8
					}
					jobs := make([][]int, jobCount)
					for i := range jobs {
						jobs[i] = []int{vals[i%len(vals)]%10 + 1, vals[(i+1)%len(vals)]%10 + 1}
					}
					_, _ = factorial.JobShopSchedulingBruteForce(jobs)
				},
				Start:   1,
				End:     8,
				Step:    1,
				Setup:   nil,
				Cleanup: nil,
			},
			"TSPBruteForceAllRoutes": {
				ExpectedBigO: bigo.Factorial,
				Sorted:       false,
				Runner: func(n int, vals []int) {
					if n < 1 {
						n = 1
					}
					cityCount := 2 + n // City counts 3 - n
					if cityCount > 20 {
						cityCount = 20
					}
					distances := make([][]int, cityCount)
					for i := range distances {
						distances[i] = make([]int, cityCount)
						for j := range distances[i] {
							if i == j {
								distances[i][j] = 0
							} else {
								distances[i][j] = vals[(i*cityCount+j)%len(vals)]%100 + 1
							}
						}
					}
					_, _ = factorial.TSPBruteForceAllRoutes(distances)
				},
				Start:   1,
				End:     10,
				Step:    1,
				Setup:   nil,
				Cleanup: nil,
			},
		*/
	}

	// hyperExponentialTimeBenchmarks contains O(n^n) benchmarks
	hyperExponentialTimeBenchmarks = map[string]BenchmarkSettings{
		/*
			"GenerateAllAssignments": {
				ExpectedBigO: bigo.HyperExponential,
				Sorted:       false,
				Runner:       func(n int, _ []int) { _ = hyperexponential.GenerateAllAssignments(n) },
				Start:        1,
				End:          8,
				Step:         1,
				Setup:        nil,
				Cleanup:      nil,
			},
			"CompleteGraphColoring": {
				ExpectedBigO: bigo.HyperExponential,
				Sorted:       false,
				Runner:       func(n int, _ []int) { _ = hyperexponential.CompleteGraphColoring(n) },
				Start:        1,
				End:          8,
				Step:         1,
				Setup:        nil,
				Cleanup:      nil,
			},
			"GenerateAllPasswords": {
				ExpectedBigO: bigo.HyperExponential,
				Sorted:       false,
				Runner:       func(n int, _ []int) { _ = hyperexponential.GenerateAllPasswords(n) },
				Start:        1,
				End:          8,
				Step:         1,
				Setup:        nil,
				Cleanup:      nil,
			},
			"WorkSimulation": {
				ExpectedBigO: bigo.HyperExponential,
				Sorted:       false,
				Runner:       func(n int, _ []int) { _ = hyperexponential.WorkSimulation(n) },
				Start:        1,
				End:          8,
				Step:         1,
				Setup:        nil,
				Cleanup:      nil,
			},
		*/
	}
)
var (
	// exampleMethodsBenchmarkSettings is a mapping between benchmark
	// method name prefixed with the BigO category and the method name to
	// the set of reasonable values to run it that will cover enough distinct
	// values to ideally let us determine the BigO complexity of the
	// function. The Setup() and Cleanup() fields allow setting up any
	// pre-work such as allocating any variables required and populating
	// them with the values to run the benchmark against.
	exampleMethodsBenchmarkSettings = map[string]BenchmarkSettings{}
)

func init() {
	for k, v := range constantTimeBenchmarks {
		exampleMethodsBenchmarkSettings[fmt.Sprintf("Constant_%s", k)] = v
	}
	for k, v := range loglogTimeBenchmarks {
		exampleMethodsBenchmarkSettings[fmt.Sprintf("LogLog_%s", k)] = v
	}
	for k, v := range logarithmicTimeBenchmarks {
		exampleMethodsBenchmarkSettings[fmt.Sprintf("Logarithmic_%s", k)] = v
	}
	for k, v := range polylogarithmicTimeBenchmarks {
		exampleMethodsBenchmarkSettings[fmt.Sprintf("Polylogarithmic_%s", k)] = v
	}
	for k, v := range linearTimeBenchmarks {
		exampleMethodsBenchmarkSettings[fmt.Sprintf("Linear_%s", k)] = v
	}
	for k, v := range nLogStarNTimeBenchmarks {
		exampleMethodsBenchmarkSettings[fmt.Sprintf("NLogStarN_%s", k)] = v
	}
	for k, v := range linearithmicTimeBenchmarks {
		exampleMethodsBenchmarkSettings[fmt.Sprintf("Linearithmic_%s", k)] = v
	}
	for k, v := range quadraticTimeBenchmarks {
		exampleMethodsBenchmarkSettings[fmt.Sprintf("Quadratic_%s", k)] = v
	}
	for k, v := range cubicTimeBenchmarks {
		exampleMethodsBenchmarkSettings[fmt.Sprintf("Cubic_%s", k)] = v
	}
	for k, v := range exponentialTimeBenchmarks {
		exampleMethodsBenchmarkSettings[fmt.Sprintf("Exponential_%s", k)] = v
	}
	for k, v := range polynomialTimeBenchmarks {
		exampleMethodsBenchmarkSettings[fmt.Sprintf("Polynomial_%s", k)] = v
	}
	for k, v := range factorialTimeBenchmarks {
		exampleMethodsBenchmarkSettings[fmt.Sprintf("Factorial_%s", k)] = v
	}
	for k, v := range hyperExponentialTimeBenchmarks {
		exampleMethodsBenchmarkSettings[fmt.Sprintf("HyperExponential_%s", k)] = v
	}
}

// BenchmarkExampleMethod lets us run any of the example methods by itself.
func BenchmarkExampleMethod(b *testing.B) {
	name := *benchmarkExampleName
	settings, ok := exampleMethodsBenchmarkSettings[name]
	if !ok {
		b.Errorf("No benchmark found for %q", name)

		return
	}

	runOneBenchmark(b, name, settings)
}

// BenchmarkAllBigOExamples is a benchmark harness to run all of the example
// methods. Consider increasing the timeout to 15m or more to allow it to
// really every benchmark here with the each ones full range of values.
func BenchmarkAllBigOExamples(b *testing.B) {
	for name, settings := range exampleMethodsBenchmarkSettings {
		runOneBenchmark(b, name, settings)
	}
}

// runOneBenchmark is a helper function to run a single benchmark with the given settings.
func runOneBenchmark(b *testing.B, name string, settings BenchmarkSettings) {
	b.Helper()
	start := settings.Start
	end := settings.End
	step := settings.Step

	if end > limit {
		b.Logf("Warning: Benchmark %q end value %d is greater than limit %d, setting to limit", name, end, limit)
		end = limit
	}

	// Override with flag values if set
	if *benchmarkStart != benchmarkStartDefault {
		start = *benchmarkStart
	}
	if *benchmarkEnd != benchmarkEndDefault {
		end = *benchmarkEnd
	}
	if *benchmarkStep != benchmarkStepDefault {
		step = *benchmarkStep
	}

	// Validate the range
	if start >= end {
		b.Skipf("Invalid benchmark range: start (%d) >= end (%d)", start, end)
	}
	if step <= 0 {
		b.Skipf("Invalid benchmark step: %d", step)
	}

	for n := start; n <= end; n += step {
		var vals []int
		if settings.Sorted {
			vals = testIntValsSorted[:n]
		} else {
			vals = testIntVals[:n]
		}

		b.Run(fmt.Sprintf("%s_n=%d", name, n),
			func(b *testing.B) {
				if settings.Setup != nil {
					b.StopTimer()
					settings.Setup(b, n, vals)
					b.StartTimer()
				}
				for b.Loop() {
					settings.Runner(n, vals)
				}
			})
		if settings.Cleanup != nil {
			b.Cleanup(func() {
				settings.Cleanup(b)
			})
		}
	}
}
