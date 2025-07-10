#!/bin/bash

# Script to run a specific Go benchmark and convert the output to CSV format.
#
# Usage: ./run_benchmark.sh <benchmark_name> [additional_go_test_args...]
#        ./run_benchmark.sh --benchmark_example_name=<example_name> [additional_go_test_args...]
#
# Example: ./run_benchmark.sh BenchmarkLinearTime --benchmark_start=1 \
#      --benchmark_end=10 --benchmark_step=1
#
# Example with example method: ./run_benchmark.sh --benchmark_example_name=LinearTime
#
# Usage examples:
#
# Basic usage with default parameters
#   ./run_benchmark.sh BenchmarkLinearTime
#
# With custom range parameters  
#   ./run_benchmark.sh BenchmarkLinearTime --benchmark_start=1 \
#       --benchmark_end=10 --benchmark_step=1
#
# With example method name
#   ./run_benchmark.sh --benchmark_example_name=LinearTime
#
# With additional Go test flags
#   ./run_benchmark.sh BenchmarkLinearTime --benchmark_start=1 \
#       --benchmark_end=5 --benchmark_step=1 --count=3
#
# Available example methods (from examples/benchmark_test.go):
# All keys from exampleMethodsBenchmarksSettings map, e.g.:
# ConstantTime_ArrayAccessByIndex, ConstantTime_HashTableLookup, Logarithmic_BinarySearch,
# LinearSearch, Linear_ArrayTraversal, Linearithmic_MergeSort, Quadratic_BubbleSort,
# Cubic_FloydWarshall, Exponential_RecursiveFibonacci, Polynomial_EditDistance,
# Factorial_GenerateAllPermutations, etc.
#
# The output filenames will be:
# - BenchmarkLinearTime.csv (default parameters)
# - BenchmarkLinearTime_start1_end10_step1.csv (custom parameters)
# - LinearTime.csv (when using --benchmark_example_name)
#

# Check if no arguments provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <benchmark_name> [additional_go_test_args...]"
    echo "   OR: $0 --benchmark_example_name=<example_name> [additional_go_test_args...]"
    echo ""
    echo "Example: $0 BenchmarkLinearTime --benchmark_start=100 --benchmark_end=100000 --benchmark_step=10000 --count=3"
    echo "Example: $0 --benchmark_example_name=LinearTime"
    echo ""
    echo "Supported benchmark flags:"
    echo "  --benchmark_start=N           Starting value for benchmark iterations"
    echo "  --benchmark_end=N             Ending value for benchmark iterations" 
    echo "  --benchmark_step=N            Step size for benchmark iterations"
    echo "  --benchmark_example_name=NAME Name of example method to benchmark"
    echo "  --count=N                     Number of times to run the benchmark"
    echo ""
    echo "Available example methods (from examples/benchmark_test.go):"
    echo "  All keys from exampleMethodsBenchmarksSettings map, e.g.:"
    echo "  ConstantTime_ArrayAccessByIndex, ConstantTime_HashTableLookup,"
    echo "  Logarithmic_BinarySearch, LinearSearch, Linear_ArrayTraversal,"
    echo "  Linearithmic_MergeSort, Quadratic_BubbleSort, Cubic_FloydWarshall,"
    echo "  Exponential_RecursiveFibonacci, Polynomial_EditDistance,"
    echo "  Factorial_GenerateAllPermutations, etc."
    exit 1
fi

# Check for --benchmark_example_name flag in all arguments
EXAMPLE_NAME=""
USING_EXAMPLE_METHOD=false
REMAINING_ARGS=()

for arg in "$@"; do
    if [[ "$arg" == --benchmark_example_name=* ]]; then
        EXAMPLE_NAME="${arg#*=}"
        USING_EXAMPLE_METHOD=true
        echo "Running example method benchmark: $EXAMPLE_NAME"
    else
        REMAINING_ARGS+=("$arg")
    fi
done

# Set up benchmark execution based on whether using example method
if [ "$USING_EXAMPLE_METHOD" = true ]; then
    BENCHMARK_NAME="BenchmarkExampleMethod"
    BENCHMARK_FLAGS="--benchmark_example_name=$EXAMPLE_NAME"
    OUTPUT_PREFIX="$EXAMPLE_NAME"
else
    # First remaining argument should be the benchmark name
    if [ ${#REMAINING_ARGS[@]} -eq 0 ]; then
        echo "Error: No benchmark name provided"
        exit 1
    fi
    BENCHMARK_NAME="${REMAINING_ARGS[0]}"
    BENCHMARK_FLAGS=""
    OUTPUT_PREFIX="$BENCHMARK_NAME"
    # Remove benchmark name from remaining args
    REMAINING_ARGS=("${REMAINING_ARGS[@]:1}")
fi

# Parse additional benchmark flags and create output filename suffix
ADDITIONAL_FLAGS=""
START_VAL=""
END_VAL=""
STEP_VAL=""
COUNT_VAL=""

for arg in "${REMAINING_ARGS[@]}"; do
    case $arg in
        --benchmark_start=*)
            START_VAL="${arg#*=}"
            ADDITIONAL_FLAGS="$ADDITIONAL_FLAGS $arg"
            ;;
        --benchmark_end=*)
            END_VAL="${arg#*=}"
            ADDITIONAL_FLAGS="$ADDITIONAL_FLAGS $arg"
            ;;
        --benchmark_step=*)
            STEP_VAL="${arg#*=}"
            ADDITIONAL_FLAGS="$ADDITIONAL_FLAGS $arg"
            ;;
        --count=*)
            COUNT_VAL="${arg#*=}"
            ADDITIONAL_FLAGS="$ADDITIONAL_FLAGS $arg"
            ;;
        *)
            # Other flags (like --run, --timeout, etc.)
            ADDITIONAL_FLAGS="$ADDITIONAL_FLAGS $arg"
            ;;
    esac
done

# Generate output suffix in consistent order: start, end, step, count
OUTPUT_SUFFIX=""
if [ -n "$START_VAL" ]; then
    OUTPUT_SUFFIX="${OUTPUT_SUFFIX}_start${START_VAL}"
fi
if [ -n "$END_VAL" ]; then
    OUTPUT_SUFFIX="${OUTPUT_SUFFIX}_end${END_VAL}"
fi
if [ -n "$STEP_VAL" ]; then
    OUTPUT_SUFFIX="${OUTPUT_SUFFIX}_step${STEP_VAL}"
fi
if [ -n "$COUNT_VAL" ]; then
    OUTPUT_SUFFIX="${OUTPUT_SUFFIX}_count${COUNT_VAL}"
fi

# Ensure output directory exists
mkdir -p timings

# Create output filename
OUTPUT_FILE="timings/${OUTPUT_PREFIX}${OUTPUT_SUFFIX}.csv"

echo "Running benchmark: ${BENCHMARK_NAME}"
if [ -n "$OUTPUT_SUFFIX" ]; then
    echo "Benchmark parameters: start=${START_VAL:-default} end=${END_VAL:-default} step=${STEP_VAL:-default}"
fi
echo "Output will be saved to: ${OUTPUT_FILE}"

# Run the benchmark and process output, save to CSV file
if [ -n "$BENCHMARK_FLAGS" ]; then
    # Using BenchmarkExampleMethod with specific example name from examples package
    go test -timeout=15m -bench="^${BENCHMARK_NAME}$" -run=^$ $BENCHMARK_FLAGS $ADDITIONAL_FLAGS . 2>&1
else
    # Using regular benchmark name
    go test -timeout=15m -bench="^${BENCHMARK_NAME}$" -run=^$ $ADDITIONAL_FLAGS . 2>&1
fi | awk '
BEGIN {
    print "n,ns_per_op"
}
/^Benchmark/ {
    # Split the line by whitespace
    split($0, fields)
    
    # Extract benchmark name (first field)
    benchmark_name = fields[1]
    
    # Extract n value from benchmark name
    # Format: BenchmarkExampleMethod/MethodName-N-8 or BenchmarkName-N-8
    # Look for pattern: -digits-digits at the end
    if (match(benchmark_name, /-([0-9]+)-[0-9]+$/)) {
        # Extract the first set of digits (the N value)
        n = substr(benchmark_name, RSTART+1, RLENGTH-3)
        # Remove the trailing "-digits" part
        gsub(/-[0-9]+$/, "", n)
    } else if (match(benchmark_name, /-([0-9]+)$/)) {
        n = substr(benchmark_name, RSTART+1, RLENGTH-1)
    } else {
        n = "unknown"
    }
    
    # Extract ns/op value (third field, remove " ns/op" suffix)
    ns_per_op = fields[3]
    gsub(/ ns\/op$/, "", ns_per_op)
    
    # Print CSV row
    print n "," ns_per_op
}' > "${OUTPUT_FILE}"

echo "Benchmark complete. Results saved to: ${OUTPUT_FILE}"