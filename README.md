# BigO - Big O Complexity Analysis Library

A Go library for analyzing algorithm complexity by characterizing real-world timing data and determining which Big O notation best fits the performance patterns.

## Overview

The `bigo` package analyzes timing measurements collected from algorithms and benchmarks to determine their time complexity characteristics. Given input size (N) and corresponding execution times (or any other measurement such as memory usage), it correlates the data against known Big O patterns to identify the most likely complexity class.

**Supported Complexity Classes:**
- O(1) - Constant Time
- O(log log n) - Double Logarithmic Time
- O(log n) - Logarithmic Time
- O((log n)^c) - Polylogarithmic Time
- O(n) - Linear Time
- O(n log* n) - n Log Star n Time
- O(n log n) - Linearithmic Time  
- O(n²) - Quadratic Time
- O(n³) - Cubic Time
- O(n^c) - Polynomial Time
- O(2^n) - Exponential Time
- O(n!) - Factorial Time
- O(n^n) - Hyper-Exponential Time

## Basic Usage

### Analyzing Timing Data

```go
    // Create a new Classifier analyzer
    b := &bigo.Classifier{}
    
    // Add timing data points (input_size, execution_time_in_ns)
    b.AddDataPoint(100, 1250.5)
    b.AddDataPoint(200, 2501.2) 
    b.AddDataPoint(400, 5002.8)
    b.AddDataPoint(800, 10008.1)
    
    // Characterize the complexity
    rank, err := b.Characterize()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Best fit: %s\n", rank)
    // Also print out the verbose summary.
    fmt.Print(b.Summary())
}
```

### Loading Data from CSV

If your data comes from an external source, **LoadCSV** is a good starting point. It expects data in a two column delimited format. Rows do not need to be unique. Each row is considered a distinct measurement, and all measurements for a given N are averaged together before the characterization is performed.  The method takes a filename, a boolean flag indicating if there is a header row in the file, and the delimiter character used in the file.  The columns are expected to be:

*   Column 1 - A numerical value representing a given N.
*   Column 2 - A numerical measurement value associated with column 1's N.

Extra columns are ignored.

```go
// Load timing data from CSV file
b := &bigo.Classifier{}
if err := b.LoadCSV("timings.csv", true, ','); err != nil {
    panic(err)
}

rank, err := b.Characterize()
if err != nil {
    panic(err)
}

fmt.Printf("BigO: %s\n", rank)
```

LoadCSV can be called multiple times to add more data.  Data is not cleared between calls to LoadCSV.


### Working with Multiple Data Points

Sometimes there are multiple runs for a given size of input. The AddDataPoint method is variadic and can take multiple measurement values for the given size. Similarly, there is a helper method AddDataPoints which takes slices of sizes and corresponding slices of slices of measurements paired up with those sizes.

```go
// Add multiple measurements for each input size
b := &bigo.Classifier{}

// Multiple runs for n=100
b.AddDataPoint(100, 1250.5, 1248.2, 1252.1)

// Batch add data points
inputSizes := []int{100, 200, 400, 800}
timings := [][]float64{
    {1250.5, 1248.2},
    {2501.2, 2498.9}, 
    {5002.8, 4999.1},
    {10008.1, 10011.2, 10007.5, 10009.3},
}

// AddDataPoints can error for things like the the two slices aren't the same length.
err := b.AddDataPoints(inputSizes, timings)
if err != nil {
    panic(err)
}
```

## Example Algorithm Implementations

The `examples/` directory contains comprehensive reference implementations for each Big O complexity class, organized by their time complexity. These implementations serve as both educational resources and test cases for the bigo library's analysis capabilities.

Each complexity class includes multiple algorithm implementations with comprehensive tests and benchmarks that generate timing data for analysis. The examples range from simple operations like array access (O(1)) to complex algorithms like the Traveling Salesman Problem (O(n!)).

For detailed information about the available algorithms and their implementations, see the [examples README](examples/README.md).

## Testing

### Running Tests

I've added comprehensive tests for both the main package and example algorithms:

```bash
# Run all tests
go test ./...
```

## Benchmarking

Because this is about analyzing data, there are a lot of benchmarks here and especially amongst the examples. Most of this is managed though two systems: 
-  Individual benchmark functions
-  A unified example benchmark running method system that can be used to generate CSV timing data for complexity analysis.

### Individual Benchmark Functions

Every example method or algorithm includes unit tests and benchmarks. 

```bash
# Run individual benchmarks by name (from root directory)
go test -bench=BenchmarkBubbleSort ./examples/quadratic/
go test -bench=BenchmarkRecursiveFibonacci ./examples/exponential/
go test -bench=BenchmarkBinarySearch ./examples/logarithmic/

# Run benchmarks by complexity class
go test -bench=BenchmarkBubbleSort -bench=BenchmarkInsertionSort -bench=BenchmarkSelectionSort ./examples/quadratic/

# For exponential and factorial benchmarks, increase timeout
go test -bench=BenchmarkRecursiveFibonacci -timeout=10m ./examples/exponential/
go test -bench=BenchmarkGenerateAllPermutations -timeout=10m ./examples/factorial/
```

### Parameterized Benchmark System

The unified example method system allows parameterized benchmarking with CSV output generation:

```bash
cd examples

# Run parameterized benchmarks with custom ranges
go test -bench=BenchmarkExampleMethod --benchmark_example_name=Quadratic_BubbleSort --benchmark_start=100 --benchmark_end=5000 --benchmark_step=200 --count=3

# Run all example method benchmarks 5 times using their default parameters
go test -bench=BenchmarkAllBigOExamples --count=5
```

### Shell Script Automation

The project includes two shell scripts for automated benchmarking:

-   **./run_all_benchmarks.sh**
-   **./run_benchmark.sh**


#### run_all_benchmarks.sh
Executes all benchmarks with optimized parameter ranges for each complexity class.

The main features of this shell script are:
- Pre-tuned parameter ranges for each complexity class (e.g., Factorial: N=3-9, Exponential: N=5-25) so they don't crush your machine or run forever.
- Complexity class filtering with flexible naming (supports `quadratic`, `QuadraticTime`, `quadratic_time`, etc.)
- Batch CSV generation with consistent file naming


```bash
# Run all benchmarks with default parameters
./run_all_benchmarks.sh
```

Run one/any/all multiple times with the **--count** flag.

```bash
# Run with custom iteration count
./run_all_benchmarks.sh --count=3
```

The resulting CSV files will look like:

```pre
n,ns_per_op
0000100,5159
0000100,4884
0000100,5003
0000400,132290
0000400,132111
0000400,132738
```

To run just a specific complexity classes example method benchmarks use **--complexity**.

Accepted case insensitive options are: Constant, LogLog, Logarithmic, Linear, Linearithmic, Quadratic, Cubic, Exponential, Polynomial, Factorial

Some alternative spelling/forms are also accepted (e.g., nlogn for Linearithmic)

_Note: This is not a repeated value flag, only one class at a time._

```bash
# Run specific complexity class only
./run_all_benchmarks.sh --complexity=Quadratic --count=2
./run_all_benchmarks.sh --complexity=Linear --count=5
```

#### run_benchmark.sh
Runs individual benchmarks or example methods with configurable parameters and generates CSV output:

The main features of this shell script are:
- Supports both individual benchmark functions and example method names
- Configurable parameters: `--benchmark_start`, `--benchmark_end`, `--benchmark_step`, `--count`
- Automatic CSV generation with standardized naming (e.g., `Quadratic_BubbleSort_start100_end5000_step200_count3.csv`)
- Parses benchmark output and extracts timing data in `n,ns_per_op` format

```bash
# Run individual benchmark function
./run_benchmark.sh BenchmarkBubbleSort

# Run example method with custom parameters
./run_benchmark.sh --benchmark_example_name=Quadratic_BubbleSort --benchmark_start=100 --benchmark_end=5000 --benchmark_step=200 --count=3
```

The run_all_benchmarks.sh script will list out the name of the supported benchmarks when passed the --help flag.


### CSV Output and Analysis

All benchmarks run by the shell scripts generate CSV files in the `timings/` directory with the format:
```csv
n,ns_per_op
100,1250.5
200,2501.2
400,5002.8
```

These CSV files can be analyzed using the bigo library itself or external tools to determine algorithmic complexity patterns.

Under the **analyze/** directory is a standalone tool for running the classifier against a CSV file. 

## High-Precision Support

Although unlikely, for algorithms with extreme, extreme performance characteristics, the library supports Go's `*big.Float` and `*big.Int` types in the **AddDataPointBig/AddDataPointsBig** methods:

```go
// Add high-precision timing data
b := &bigo.Classifier{}
bigTime := big.NewFloat(1.23456789012345e15) 
b.AddDataPointBig(1000000, bigTime)

// Batch add big float data
bigTimings := [][]*big.Float{
    {big.NewFloat(1.5e10)},
    {big.NewFloat(3.2e12)},
}

err := b.AddDataPointsBig(inputSizes, bigTimings)
```

Generally the high-precision value support is used internally for Big O complexity classes greater than Linear where the chance of generating a comparison value that overflows a float64 becomes likely.  For example, any value of N > 170 for factorial will exceed a float64's limit, but it's possible a data set matching a lower Big O will have N values running into the thousands or millions that we are hoping to compare to. Even algorithms in Linear time are likely to be able to generate results with a million or more inputs on moderated hardware in reasonable time.
