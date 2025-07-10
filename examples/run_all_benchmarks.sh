#!/bin/bash

# Script to run all benchmarks in examples/benchmark_test.go with their associated parameters
# Uses the run_benchmark.sh script to execute each benchmark and generate CSV files
#
# Usage: ./run_all_benchmarks.sh [--count=N] [--complexity=CLASS] [--help]
# Example: ./run_all_benchmarks.sh --count=3
# Example: ./run_all_benchmarks.sh --complexity=Quadratic --count=2
# Example: ./run_all_benchmarks.sh --help

# Parse command line arguments
COUNT_FLAG="--count=1"
COMPLEXITY_FILTER=""
for arg in "$@"; do
    case $arg in
        --count=*)
            COUNT_FLAG="$arg"
            ;;
        --complexity=*)
            # Normalize complexity class name (case insensitive)
            raw_complexity="${arg#*=}"
            # Convert to lowercase first
            lower_complexity="$(echo "${raw_complexity}" | tr '[:upper:]' '[:lower:]')"
            # Handle special cases and normalize
            case "$lower_complexity" in
                "constant"|"constanttime"|"constant_time")
                    COMPLEXITY_FILTER="Constant"
                    ;;
                "loglog"|"log-log"|"loglogtime"|"loglog_time")
                    COMPLEXITY_FILTER="LogLog"
                    ;;
                "logarithmic"|"logarithmictime"|"log"|"logarithmic_time")
                    COMPLEXITY_FILTER="Logarithmic"
                    ;;
                "linear"|"lineartime"|"linear_time")
                    COMPLEXITY_FILTER="Linear"
                    ;;
                "linearithmic"|"linearithmictime"|"nlogn")
                    COMPLEXITY_FILTER="Linearithmic"
                    ;;
                "quadratic"|"quadratictime"|"quadratic_time")
                    COMPLEXITY_FILTER="Quadratic"
                    ;;
                "cubic"|"cubictime"|"cubic_time")
                    COMPLEXITY_FILTER="Cubic"
                    ;;
                "exponential"|"exponentialtime"|"exp"|"exponential_time")
                    COMPLEXITY_FILTER="Exponential"
                    ;;
                "polynomial"|"polynomialtime"|"poly"|"polynomial_time")
                    COMPLEXITY_FILTER="Polynomial"
                    ;;
                "factorial"|"factorialtime"|"factorial_time")
                    COMPLEXITY_FILTER="Factorial"
                    ;;
                *)
                    echo "Error: Unknown complexity class '$raw_complexity'" >&2
                    echo "Valid complexity classes (case insensitive):" >&2
                    echo "  Constant, LogLog, Logarithmic, Linear, Linearithmic," >&2
                    echo "  Quadratic, Cubic, Exponential, Polynomial, Factorial" >&2
                    echo "Alternative names also accepted (e.g., nlogn for Linearithmic)" >&2
                    exit 1
                    ;;
            esac
            ;;
        --help|-h)
            echo "Usage: $0 [--count=N] [--complexity=CLASS] [--help]"
            echo ""
            echo "Options:"
            echo "  --count=N           Number of times to run each benchmark (default: 1)"
            echo "  --complexity=CLASS  Run only benchmarks for specific Big O class (case insensitive):"
            echo "                      Constant, LogLog, Logarithmic, Linear, Linearithmic,"
            echo "                      Quadratic, Cubic, Exponential, Polynomial, Factorial"
            echo "  --help, -h          Show this help message"
            echo ""
            echo "Supported benchmark methods:"
            echo ""
            echo "Constant Time (O(1)):"
            echo "  Constant_ArrayAccessByIndex, Constant_HashTableLookup, Constant_LinkedListAccessFirst,"
            echo "  Constant_LinkedListAccessLast, Constant_StackPush, Constant_StackPop,"
            echo "  Constant_QueueEnqueue, Constant_QueueDequeue"
            echo ""
            echo "Log-Log Time (O(log log n)):"
            echo "  LogLog_InterpolationSearch, LogLog_RadixSortOptimization, LogLog_ParallelDivideConquer"
            echo ""
            echo "Logarithmic Time (O(log n)):"
            echo "  Logarithmic_BinarySearch, Logarithmic_BinaryTreeSearch, Logarithmic_FindTreeHeight"
            echo ""
            echo "Linear Time (O(n)):"
            echo "  Linear_Search, Linear_ArrayTraversal, Linear_CountElements, Linear_FindMinimum,"
            echo "  Linear_FindMaximum, Linear_CalculateSum"
            echo ""
            echo "Linearithmic Time (O(n log n)):"
            echo "  Linearithmic_MergeSort, Linearithmic_BuildHeapFromArray, Linearithmic_HeapifyArray,"
            echo "  Linearithmic_IntroSort, Linearithmic_HeapSort, Linearithmic_QuickSort"
            echo ""
            echo "Quadratic Time (O(n²)):"
            echo "  Quadratic_BubbleSort, Quadratic_InsertionSort, Quadratic_SelectionSort,"
            echo "  Quadratic_AllPairsComparison, Quadratic_FindDuplicatePairs, Quadratic_CountInversions,"
            echo "  Quadratic_TwoSum, Quadratic_NaiveMatrixMultiplication, Quadratic_MatrixTranspose"
            echo ""
            echo "Cubic Time (O(n³)):"
            echo "  Cubic_FloydWarshall, Cubic_ThreeSum, Cubic_MatrixChainMultiplication,"
            echo "  Cubic_OptimalBinarySearchTree, Cubic_StandardMatrixMultiplication,"
            echo "  Cubic_TripleNestedProductSum, Cubic_FindTripletsWithSum, Cubic_CountTripletsWithProperty,"
            echo "  Cubic_Generate3DCombinations"
            echo ""
            echo "Exponential Time (O(2ⁿ)):"
            echo "  Exponential_RecursiveFibonacci, Exponential_TowerOfHanoiCount, Exponential_GenerateAllSubsets,"
            echo "  Exponential_TravelingSalesmanBruteForce, Exponential_TSPBitMask"
            echo ""
            echo "Polynomial Time (O(nᵏ)):"
            echo "  Polynomial_EditDistance, Polynomial_JohnsonAlgorithm, Polynomial_LongestCommonSubsequence,"
            echo "  Polynomial_LCSWithSequence, Polynomial_MatrixChainOrder, Polynomial_FordFulkerson,"
            echo "  Polynomial_EdmondsKarp"
            echo ""
            echo "Factorial Time (O(n!)):"
            echo "  Factorial_GenerateAllPermutations, Factorial_AssignmentProblemBruteForce,"
            echo "  Factorial_OptimalMatchingBruteForce, Factorial_NQueensAllArrangements,"
            echo "  Factorial_NQueensCountAllArrangements, Factorial_GenerateAllSchedules,"
            echo "  Factorial_OptimalScheduleBruteForce, Factorial_JobShopSchedulingBruteForce,"
            echo "  Factorial_TSPBruteForceAllRoutes"
            echo ""
            echo "Examples:"
            echo "  $0 --count=3                         # Run all benchmarks 3 times"
            echo "  $0 --complexity=Quadratic --count=2  # Run only quadratic benchmarks 2 times"
            echo "  $0 --complexity=linear               # Run only linear benchmarks (case insensitive)"
            exit 0
            ;;
        *)
            echo "Unknown argument: $arg"
            echo "Usage: $0 [--count=N] [--complexity=CLASS] [--help]"
            echo "  --count=N           Number of times to run each benchmark (default: 1)"
            echo "  --complexity=CLASS  Run only benchmarks for specific Big O class (case insensitive):"
            echo "                      Constant, LogLog, Logarithmic, Linear, Linearithmic,"
            echo "                      Quadratic, Cubic, Exponential, Polynomial, Factorial"
            echo "  --help, -h          Show this help message"
            exit 1
            ;;
    esac
done

# Check if run_benchmark.sh exists
if [ ! -f "./run_benchmark.sh" ]; then
    echo "Error: run_benchmark.sh not found in current directory"
    exit 1
fi

# Make sure run_benchmark.sh is executable
chmod +x ./run_benchmark.sh

# Function to filter benchmarks by complexity class
filter_by_complexity() {
    local params="$1"
    local filter="$2"
    
    if [ -z "$filter" ]; then
        echo "$params"
        return
    fi
    
    # Define complexity class groups
    case "$filter" in
        "Constant")
            echo "$params" | grep -E "^(Constant_ArrayAccessByIndex|Constant_HashTableLookup|Constant_LinkedListAccessFirst|Constant_LinkedListAccessLast|Constant_StackPush|Constant_StackPop|Constant_QueueEnqueue|Constant_QueueDequeue):"
            ;;
        "LogLog")
            echo "$params" | grep -E "^(LogLog_InterpolationSearch|LogLog_RadixSortOptimization|LogLog_ParallelDivideConquer):"
            ;;
        "Logarithmic")
            echo "$params" | grep -E "^(Logarithmic_BinarySearch|Logarithmic_BinaryTreeSearch|Logarithmic_FindTreeHeight):"
            ;;
        "Linear")
            echo "$params" | grep -E "^(Linear_Search|Linear_ArrayTraversal|Linear_CountElements|Linear_FindMinimum|Linear_FindMaximum|Linear_CalculateSum):"
            ;;
        "Linearithmic")
            echo "$params" | grep -E "^(Linearithmic_MergeSort|Linearithmic_BuildHeapFromArray|Linearithmic_HeapifyArray|Linearithmic_IntroSort|Linearithmic_HeapSort|Linearithmic_QuickSort):"
            ;;
        "Quadratic")
            echo "$params" | grep -E "^(Quadratic_BubbleSort|Quadratic_InsertionSort|Quadratic_SelectionSort|Quadratic_AllPairsComparison|Quadratic_FindDuplicatePairs|Quadratic_CountInversions|Quadratic_TwoSum|Quadratic_NaiveMatrixMultiplication|Quadratic_MatrixTranspose):"
            ;;
        "Cubic")
            echo "$params" | grep -E "^(Cubic_FloydWarshall|Cubic_ThreeSum|Cubic_MatrixChainMultiplication|Cubic_OptimalBinarySearchTree|Cubic_StandardMatrixMultiplication|Cubic_TripleNestedProductSum|Cubic_FindTripletsWithSum|Cubic_CountTripletsWithProperty|Cubic_Generate3DCombinations):"
            ;;
        "Exponential")
            echo "$params" | grep -E "^(Exponential_RecursiveFibonacci|Exponential_TowerOfHanoiCount|Exponential_GenerateAllSubsets|Exponential_TravelingSalesmanBruteForce|Exponential_TSPBitMask):"
            ;;
        "Polynomial")
            echo "$params" | grep -E "^(Polynomial_EditDistance|Polynomial_JohnsonAlgorithm|Polynomial_LongestCommonSubsequence|Polynomial_LCSWithSequence|Polynomial_MatrixChainOrder|Polynomial_FordFulkerson|Polynomial_EdmondsKarp):"
            ;;
        "Factorial")
            echo "$params" | grep -E "^(Factorial_GenerateAllPermutations|Factorial_AssignmentProblemBruteForce|Factorial_OptimalMatchingBruteForce|Factorial_NQueensAllArrangements|Factorial_NQueensCountAllArrangements|Factorial_GenerateAllSchedules|Factorial_OptimalScheduleBruteForce|Factorial_JobShopSchedulingBruteForce|Factorial_TSPBruteForceAllRoutes):"
            ;;
        *)
            echo "Unknown complexity class: $filter" >&2
            echo "$params"
            ;;
    esac
}

echo "Running benchmarks with their associated parameters..."
echo "Using count flag: $COUNT_FLAG"
if [ -n "$COMPLEXITY_FILTER" ]; then
    echo "Filtering by complexity class: $COMPLEXITY_FILTER"
fi
echo "Running full benchmark set (all example methods)"
echo "=========================================================="

# Benchmark parameters: function_name:start:end:step
# Optimized for 120s max runtime with sufficient data points for Big O analysis
BENCHMARK_PARAMS="
Constant_ArrayAccessByIndex:50000:2000000:150000
Constant_HashTableLookup:50000:2000000:150000
Constant_LinkedListAccessFirst:50000:2000000:150000
Constant_LinkedListAccessLast:50000:2000000:150000
Constant_StackPush:50000:2000000:150000
Constant_StackPop:50000:2000000:150000
Constant_QueueEnqueue:50000:2000000:150000
Constant_QueueDequeue:50000:2000000:150000
LogLog_InterpolationSearch:50000:500000:50000
LogLog_RadixSortOptimization:10000:200000:20000
LogLog_ParallelDivideConquer:10000:200000:20000
Logarithmic_BinarySearch:10000:5000000:500000
Logarithmic_BinaryTreeSearch:1000:50000:5000
Logarithmic_FindTreeHeight:1000:50000:5000
Linear_Search:5000:200000:15000
Linear_ArrayTraversal:5000:200000:15000
Linear_CountElements:5000:200000:15000
Linear_FindMinimum:5000:200000:15000
Linear_FindMaximum:5000:200000:15000
Linear_CalculateSum:5000:200000:15000
Linearithmic_MergeSort:1000:200000:15000
Linearithmic_BuildHeapFromArray:1000:200000:15000
Linearithmic_HeapifyArray:1000:200000:15000
Linearithmic_IntroSort:1000:200000:15000
Linearithmic_HeapSort:1000:200000:15000
Linearithmic_QuickSort:1000:200000:15000
Quadratic_BubbleSort:100:8000:300
Quadratic_InsertionSort:100:8000:300
Quadratic_SelectionSort:100:8000:300
Quadratic_AllPairsComparison:50:1000:75
Quadratic_FindDuplicatePairs:200:8000:400
Quadratic_CountInversions:200:8000:400
Quadratic_TwoSum:200:8000:400
Quadratic_NaiveMatrixMultiplication:10:200:15
Quadratic_MatrixTranspose:100:2000:150
Cubic_CountTripletsWithProperty:20:400:30
Cubic_FindTripletsWithSum:20:400:30
Cubic_FloydWarshall:10:200:15
Cubic_Generate3DCombinations:5:30:2
Cubic_MatrixChainMultiplication:20:300:20
Cubic_OptimalBinarySearchTree:10:150:10
Cubic_StandardMatrixMultiplication:20:300:20
Cubic_TripleNestedProductSum:20:400:30
Cubic_ThreeSum:20:400:30
Exponential_RecursiveFibonacci:5:40:3
Exponential_TowerOfHanoiCount:5:25:2
Exponential_GenerateAllSubsets:5:25:2
Exponential_TravelingSalesmanBruteForce:3:15:1
Exponential_TSPBitMask:5:20:1
Polynomial_EditDistance:100:3000:200
Polynomial_JohnsonAlgorithm:100:2000:150
Polynomial_LongestCommonSubsequence:100:3000:200
Polynomial_LCSWithSequence:100:3000:200
Polynomial_MatrixChainOrder:100:2000:150
Polynomial_FordFulkerson:100:2000:150
Polynomial_EdmondsKarp:100:2000:150
Factorial_GenerateAllPermutations:3:9:1
Factorial_AssignmentProblemBruteForce:3:10:1
Factorial_OptimalMatchingBruteForce:3:9:1
Factorial_NQueensAllArrangements:3:8:1
Factorial_NQueensCountAllArrangements:3:8:1
Factorial_GenerateAllSchedules:3:6:1
Factorial_OptimalScheduleBruteForce:3:8:1
Factorial_JobShopSchedulingBruteForce:3:9:1
Factorial_TSPBruteForceAllRoutes:3:9:1
"

# Use the full benchmark set
RAW_PARAMS="$BENCHMARK_PARAMS"
echo "Selected: Full benchmark set"

# Apply complexity filter if specified
SELECTED_PARAMS=$(filter_by_complexity "$RAW_PARAMS" "$COMPLEXITY_FILTER")

# Counter for progress tracking
TOTAL_BENCHMARKS=$(echo "$SELECTED_PARAMS" | grep -c "^[A-Z]")
CURRENT=0

# Run each benchmark with its parameters
for param in $(echo "$SELECTED_PARAMS" | grep "^[A-Z]"); do
    CURRENT=$((CURRENT + 1))
    
    # Parse the parameter string
    EXAMPLE_NAME=$(echo "$param" | cut -d: -f1)
    START_VAL=$(echo "$param" | cut -d: -f2)
    END_VAL=$(echo "$param" | cut -d: -f3)
    STEP_VAL=$(echo "$param" | cut -d: -f4)
    
    echo ""
    echo "[$CURRENT/$TOTAL_BENCHMARKS] Running example method: $EXAMPLE_NAME..."
    echo "Parameters: start=$START_VAL end=$END_VAL step=$STEP_VAL"
    
    # Run the benchmark with --benchmark_example_name flag
    time ./run_benchmark.sh --benchmark_example_name="$EXAMPLE_NAME" --benchmark_start="$START_VAL" --benchmark_end="$END_VAL" --benchmark_step="$STEP_VAL" "$COUNT_FLAG"
    
    if [ $? -eq 0 ]; then
        echo "✓ $EXAMPLE_NAME completed successfully"
    else
        echo "✗ $EXAMPLE_NAME failed"
    fi
done

echo ""
echo "=========================================================="
echo "All benchmarks completed!"
echo ""
echo "Generated CSV files:"
ls -1 timings/*.csv 2>/dev/null | sort
echo ""
echo "Total CSV files: $(ls -1 timings/*.csv 2>/dev/null | wc -l)"