# Examples Directory

## Overview

This directory tree contains reference implementations for various Big O complexity classes, organized by their time complexity. Each subdirectory demonstrates algorithms and data structures that exhibit specific computational complexity patterns, providing concrete examples for understanding how different Big O notations manifest in real code.

These examples serve as both educational resources and test cases for the bigo library's complexity analysis capabilities. Each implementation includes comprehensive tests as well as simple benchmarks.

## Complexity Classes

### Constant Time: **O(1)** 

Directory: [constant_time/](constant_time)

Constant time operations that execute in the same amount of time regardless of input size. These operations access data structures directly without iteration or recursion.

**Files and Methods:**
- `array_access.go` - Direct array element access by index
- `basic_math.go` - Basic arithmetic operations that don't depend on input size
- `hash_lookup.go` - Hash table/map lookup operations with O(1) average case
- `linked_list_access.go` - Direct access to linked list head/tail nodes
- `queue_operations.go` - Queue enqueue/dequeue operations using slices
- `stack_operations.go` - Stack push/pop operations using slices

### Inverse Ackermann: **O(α(n))**

Directory: [inverse_ackermann/](inverse_ackermann/)

Inverse Ackermann function complexity, an extremely slowly growing function that appears in advanced data structures like Union-Find with path compression.

**Files and Methods:**
- `chazelle_minimum_spanning_tree.go` - Chazelle's linear-time MST algorithm
- `tarjans_lowest_common_ancestor.go` - Tarjan's offline LCA algorithm with Union-Find

### Double Log: **O(log log n)**

Directory: [loglog_time/](loglog_time/)

Double logarithmic time operations, a very slow-growing complexity class seen in specialized data structures and algorithms.

**Files and Methods:**
- `fusion_tree.go` - Fusion tree operations for integer sorting
- `interpolation_search.go` - Interpolation search on uniformly distributed data
- `parallel_algorithm.go` - Parallel algorithms with logarithmic depth
- `radix_optimization.go` - Optimized radix-based operations
- `van_emde_boas.go` - Van Emde Boas tree operations

### Logarithmic Time: **O(log n)**

Directory: [logarithmic_time/](logarithmic_time/)

Logarithmic time operations that repeatedly divide the problem space in half, typically seen in tree-based operations and binary search algorithms.

**Files and Methods:**
- `binary_indexed_tree.go` - Binary Indexed Tree (Fenwick Tree) operations
- `binary_search.go` - Classic binary search on sorted arrays
- `binary_tree_search.go` - Binary search tree lookup operations
- `heap_operations.go` - Min/max heap insertion and deletion
- `tree_height.go` - Binary tree height calculation
- `type_tree_node.go` - Binary tree node definition

### Linear: **O(n)**

Directory: [linear_time/](linear_time/)

Linear time operations that scale proportionally with input size. These algorithms typically involve a single pass through the data.

**Files and Methods:**
- `count_elements.go` - Element counting operations that traverse arrays once
- `find_minmax.go` - `FindMinimum()`, `FindMaximum()`, `FindMinMax()`: Single-pass searches
- `search.go` - Linear search through unsorted arrays
- `single_pass.go` - Various single-pass array processing algorithms
- `traversal.go` - Array and slice traversal patterns
- `type_list_node.go` - Linked list node definition and linear traversal

### Linearithmic: **O(n log n)**

Directory: [linearithmic_time/](linearithmic_time/)

Linearithmic time operations that combine linear and logarithmic components, most commonly seen in efficient sorting algorithms and divide-and-conquer approaches.

**Files and Methods:**
- `boruvka_mst.go` - Borůvka's minimum spanning tree algorithm
- `build_heap.go` - Heap construction using bottom-up approach
- `comparison_sorts.go` - Various O(n log n) sorting algorithm implementations
- `heap_sort.go` - Heap sort implementation
- `kruskal_mst.go` - Kruskal's minimum spanning tree algorithm
- `merge_sort.go` - Merge sort divide-and-conquer implementation
- `quick_sort.go` - Quick sort with average O(n log n) complexity

### Quadratic: **O(n²)**

Directory: [quadratic_time/](quadratic_time/)

Quadratic time operations featuring nested loops or comparisons between all pairs of elements. Performance degrades rapidly with input size.

**Files and Methods:**
- `all_pairs_comparison.go` - Algorithms that compare every pair of elements
- `bubble_sort.go` - Bubble sort with nested comparison loops
- `insertion_sort.go` - Insertion sort with element shifting
- `matrix_multiplication.go` - Naive matrix multiplication algorithm
- `selection_sort.go` - Selection sort with nested selection loops

### Cubic: **O(n³)**

Directory: [cubic_time/](cubic_time/)

Cubic time operations with three nested loops or operations on three-dimensional data structures. Often found in dynamic programming solutions and triple-nested algorithms.

**Files and Methods:**
- `dynamic_programming.go` - DP solutions with cubic time complexity
- `floyd_warshall.go` - Floyd-Warshall all-pairs shortest path algorithm
- `matrix_multiplication.go` - Three-dimensional matrix operations
- `three_sum.go` - Three-sum problem with triple nested loops
- `triple_nested_brute_force.go` - Brute force algorithms with three nested iterations

### Polynomial: **O(nᵏ) where k > 3**

Directory: [polynomial_time/](polynomial_time/)

Polynomial time operations with higher-degree polynomials, often seen in complex dynamic programming solutions and graph algorithms.

**Files and Methods:**
- `all_pairs_shortest_paths.go` - Advanced shortest path algorithms
- `edit_distance.go` - Edit distance calculation using dynamic programming
- `longest_common_subsequence.go` - LCS dynamic programming solution
- `matrix_chain.go` - Matrix chain multiplication optimization
- `maximum_flow.go` - Maximum flow algorithms with polynomial complexity

### Exponential: **O(2ⁿ)**

Directory: [exponential_time/](exponential_time/)

Exponential time operations where the runtime doubles with each additional input element. Often seen in exhaustive search algorithms and naive recursive solutions.

**Files and Methods:**
- `generate_subsets.go` - Generate all possible subsets of a set
- `n_queens.go` - N-Queens problem backtracking solution
- `recursive_fibonacci.go` - `RecursiveFibonacci()`: Naive recursive Fibonacci implementation
- `tower_of_hanoi.go` - Tower of Hanoi recursive solution
- `traveling_salesman.go` - Traveling Salesman Problem brute force approach

### Factorial: **O(n!)**

Directory: [factorial_time/](factorial_time/)

Factorial time operations where runtime grows as the factorial of input size. These represent the most computationally expensive algorithms, typically involving permutations or exhaustive enumeration.

**Files and Methods:**
- `assignment_problem.go` - Assignment problem with all possible assignments
- `generate_permutations.go` - Generate all permutations of a set
- `n_queens_all_arrangements.go` - N-Queens finding all possible solutions
- `scheduling_problems.go` - Exhaustive scheduling optimization
- `traveling_salesman_brute_force.go` - TSP examining all possible routes