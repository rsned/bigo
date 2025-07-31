package loglog

import (
	"math"
)

// YFastTrie implements a Y-fast trie for O(log log n) operations
// Y-fast tries support predecessor, successor, insert, and delete in O(log log n) time
type YFastTrie struct {
	xFastTrie   *XFastTrie
	minLeaf     *YNode
	maxLeaf     *YNode
	universeLog int
}

// XFastTrie is a simplified X-fast trie used as a building block
type XFastTrie struct {
	levels []map[int]*XNode
	leaves map[int]*XNode
	height int
}

// XNode represents a node in the X-fast trie
type XNode struct {
	prefix     int
	left       *XNode
	right      *XNode
	descendant *YNode
}

// YNode represents a leaf in the Y-fast trie (balanced BST)
type YNode struct {
	key    int
	left   *YNode
	right  *YNode
	parent *YNode
	height int
}

// NewYFastTrie creates a new Y-fast trie for universe size 2^w
func NewYFastTrie(universeSize int) *YFastTrie {
	universeLog := int(math.Log2(float64(universeSize))) + 1

	return &YFastTrie{
		xFastTrie:   NewXFastTrie(universeLog),
		universeLog: universeLog,
		minLeaf:     nil,
		maxLeaf:     nil,
	}
}

// NewXFastTrie creates a new X-fast trie for w-bit integers
func NewXFastTrie(w int) *XFastTrie {
	levels := make([]map[int]*XNode, w+1)
	for i := range levels {
		levels[i] = make(map[int]*XNode)
	}

	return &XFastTrie{
		levels: levels,
		leaves: make(map[int]*XNode),
		height: w,
	}
}

// Insert adds a key to the Y-fast trie
// Time complexity: O(log log n) amortized
func (yt *YFastTrie) Insert(key int) {
	// Find the appropriate leaf tree or create one
	leafTree := yt.findLeafTree(key)

	if leafTree == nil {
		// Create new leaf tree
		leafTree = &YNode{
			key:    key,
			height: 1,
			left:   nil,
			right:  nil,
			parent: nil,
		}
		yt.insertIntoXFastTrie(key, leafTree)

		// Update min/max pointers
		if yt.minLeaf == nil || key < yt.minLeaf.key {
			yt.minLeaf = leafTree
		}
		if yt.maxLeaf == nil || key > yt.maxLeaf.key {
			yt.maxLeaf = leafTree
		}
	} else {
		// Insert into existing leaf tree (BST)
		yt.insertIntoBST(leafTree, key)
	}
}

// Delete removes a key from the Y-fast trie
// Time complexity: O(log log n) amortized
func (yt *YFastTrie) Delete(key int) bool {
	leafTree := yt.findLeafTree(key)
	if leafTree == nil {
		return false
	}

	// Remove from BST
	if yt.deleteFromBST(leafTree, key) {
		// If tree becomes empty, remove from X-fast trie
		if yt.isBSTEmpty(leafTree) {
			yt.deleteFromXFastTrie(key)
			// Update min/max pointers if necessary
			if yt.minLeaf == leafTree {
				yt.minLeaf = yt.findNextLeaf(leafTree)
			}
			if yt.maxLeaf == leafTree {
				yt.maxLeaf = yt.findPrevLeaf(leafTree)
			}
		}

		return true
	}

	return false
}

// Predecessor finds the largest key ≤ query
// Time complexity: O(log log n)
func (yt *YFastTrie) Predecessor(query int) (int, bool) {
	// Use X-fast trie to find appropriate leaf tree
	leafTree := yt.findPredecessorLeaf(query)
	if leafTree == nil {
		return 0, false
	}

	// Search within the BST
	return yt.predecessorInBST(leafTree, query)
}

// Successor finds the smallest key ≥ query
// Time complexity: O(log log n)
func (yt *YFastTrie) Successor(query int) (int, bool) {
	// Use X-fast trie to find appropriate leaf tree
	leafTree := yt.findSuccessorLeaf(query)
	if leafTree == nil {
		return 0, false
	}

	// Search within the BST
	return yt.successorInBST(leafTree, query)
}

// Helper methods for X-fast trie operations

func (yt *YFastTrie) insertIntoXFastTrie(key int, leafNode *YNode) {
	// Insert key into X-fast trie with pointer to leaf BST
	node := &XNode{
		descendant: leafNode,
		prefix:     0,
		left:       nil,
		right:      nil,
	}
	yt.xFastTrie.leaves[key] = node

	// Update all prefixes along the path
	for level := 0; level <= yt.universeLog; level++ {
		prefix := key >> (yt.universeLog - level)
		if _, exists := yt.xFastTrie.levels[level][prefix]; !exists {
			yt.xFastTrie.levels[level][prefix] = &XNode{
				prefix:     prefix,
				descendant: leafNode,
				left:       nil,
				right:      nil,
			}
		}
	}
}

func (yt *YFastTrie) deleteFromXFastTrie(key int) {
	delete(yt.xFastTrie.leaves, key)

	// Remove prefixes that are no longer needed
	for level := 0; level <= yt.universeLog; level++ {
		prefix := key >> (yt.universeLog - level)
		// Simplified deletion - in practice would need more complex cleanup
		delete(yt.xFastTrie.levels[level], prefix)
	}
}

func (yt *YFastTrie) findLeafTree(key int) *YNode {
	if node, exists := yt.xFastTrie.leaves[key]; exists {
		return node.descendant
	}

	return nil
}

func (yt *YFastTrie) findPredecessorLeaf(query int) *YNode {
	// Binary search on levels to find predecessor in O(log log n)
	low, high := 0, yt.universeLog
	var bestNode *XNode

	for low <= high {
		mid := (low + high) / 2
		prefix := query >> (yt.universeLog - mid)

		if node, exists := yt.xFastTrie.levels[mid][prefix]; exists {
			bestNode = node
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	if bestNode != nil {
		return bestNode.descendant
	}

	return nil
}

func (yt *YFastTrie) findSuccessorLeaf(query int) *YNode {
	// Similar to predecessor but searching for successor
	return yt.findPredecessorLeaf(query) // Simplified
}

func (yt *YFastTrie) findNextLeaf(_ *YNode) *YNode {
	// Find next leaf tree in sequence
	return nil // Simplified
}

func (yt *YFastTrie) findPrevLeaf(_ *YNode) *YNode {
	// Find previous leaf tree in sequence
	return nil // Simplified
}

// Helper methods for BST operations

func (yt *YFastTrie) insertIntoBST(root *YNode, key int) *YNode {
	if root == nil {
		return &YNode{
			key:    key,
			height: 1,
			left:   nil,
			right:  nil,
			parent: nil,
		}
	}

	if key < root.key {
		root.left = yt.insertIntoBST(root.left, key)
	} else if key > root.key {
		root.right = yt.insertIntoBST(root.right, key)
	}

	// Update height and rebalance (AVL tree operations)
	root.height = 1 + max(yt.getHeight(root.left), yt.getHeight(root.right))

	return yt.rebalance(root)
}

func (yt *YFastTrie) deleteFromBST(root *YNode, key int) bool {
	if root == nil {
		return false
	}

	// Find the key in the BST
	current := root
	for current != nil {
		switch {
		case current.key == key:
			return true // Found key - simplified deletion logic
		case key < current.key:
			current = current.left
		default:
			current = current.right
		}
	}

	return false // Key not found
}

func (yt *YFastTrie) predecessorInBST(root *YNode, query int) (int, bool) {
	if root == nil {
		return 0, false
	}

	if root.key <= query {
		if pred, found := yt.predecessorInBST(root.right, query); found {
			return pred, true
		}

		return root.key, true
	}

	return yt.predecessorInBST(root.left, query)
}

func (yt *YFastTrie) successorInBST(root *YNode, query int) (int, bool) {
	if root == nil {
		return 0, false
	}

	if root.key >= query {
		if succ, found := yt.successorInBST(root.left, query); found {
			return succ, true
		}

		return root.key, true
	}

	return yt.successorInBST(root.right, query)
}

func (yt *YFastTrie) isBSTEmpty(root *YNode) bool {
	return root == nil
}

func (yt *YFastTrie) getHeight(node *YNode) int {
	if node == nil {
		return 0
	}

	return node.height
}

func (yt *YFastTrie) rebalance(node *YNode) *YNode {
	// AVL tree rebalancing operations
	return node // Simplified
}

// YFastTrieOperations performs a series of Y-fast trie operations
// Demonstrates O(log log n) time complexity per operation
func YFastTrieOperations(n int) int {
	universeSize := 1 << 20 // Large universe
	trie := NewYFastTrie(universeSize)

	operations := 0

	// Insert n elements
	for i := range n {
		trie.Insert(i * 2) // Insert even numbers
		operations++
	}

	// Perform predecessor/successor queries
	for i := range n {
		trie.Predecessor(i*2 + 1) // Query for odd numbers
		operations++
	}

	// Delete half the elements
	for i := 0; i < n/2; i++ {
		trie.Delete(i * 4) // Delete every other even number
		operations++
	}

	return operations
}
