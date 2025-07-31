package loglog

import (
	"testing"
)

func TestNewYFastTrie(t *testing.T) {
	trie := NewYFastTrie(1024)
	if trie == nil {
		t.Fatal("Expected non-nil Y-fast trie")
	}

	if trie.universeLog != 11 { // log2(1024) + 1
		t.Errorf("Expected universeLog to be 11, got %d", trie.universeLog)
	}
}

func TestYFastTrieInsert(t *testing.T) {
	trie := NewYFastTrie(1024)

	// Insert some keys
	trie.Insert(10)
	trie.Insert(20)
	trie.Insert(15)

	// Basic sanity check - should not panic
	if trie.minLeaf == nil {
		t.Error("Expected minLeaf to be set after insertions")
	}
}

func TestYFastTriePredecessor(t *testing.T) {
	trie := NewYFastTrie(1024)

	// Insert keys
	keys := []int{10, 20, 30, 40, 50}
	for _, key := range keys {
		trie.Insert(key)
	}

	// Test predecessor queries
	if pred, found := trie.Predecessor(25); found {
		if pred > 25 {
			t.Errorf("Predecessor of 25 should be ≤ 25, got %d", pred)
		}
	}

	// Test predecessor of non-existent key
	if pred, found := trie.Predecessor(5); found {
		t.Errorf("Should not find predecessor of 5, but got %d", pred)
	}
}

func TestYFastTrieSuccessor(t *testing.T) {
	trie := NewYFastTrie(1024)

	// Insert keys
	keys := []int{10, 20, 30, 40, 50}
	for _, key := range keys {
		trie.Insert(key)
	}

	// Test successor queries
	if succ, found := trie.Successor(25); found {
		if succ < 25 {
			t.Errorf("Successor of 25 should be ≥ 25, got %d", succ)
		}
	}

	// Test successor beyond max
	if succ, found := trie.Successor(100); found {
		t.Errorf("Should not find successor of 100, but got %d", succ)
	}
}

func TestYFastTrieDelete(t *testing.T) {
	trie := NewYFastTrie(1024)

	// Insert and delete
	trie.Insert(10)
	trie.Insert(20)

	// Delete existing key
	if !trie.Delete(10) {
		t.Error("Should successfully delete existing key 10")
	}

	// Delete non-existing key
	if trie.Delete(100) {
		t.Error("Should not delete non-existing key 100")
	}
}

func TestXFastTrie(t *testing.T) {
	xTrie := NewXFastTrie(10)
	if xTrie == nil {
		t.Fatal("Expected non-nil X-fast trie")
	}

	if xTrie.height != 10 {
		t.Errorf("Expected height 10, got %d", xTrie.height)
	}
}

func TestYFastTrieOperations(t *testing.T) {
	n := 100
	operations := YFastTrieOperations(n)

	// Should perform approximately 2.5n operations
	expectedMin := 2*n - 10
	expectedMax := 3*n + 10

	if operations < expectedMin || operations > expectedMax {
		t.Errorf("Expected operations between %d and %d, got %d", expectedMin, expectedMax, operations)
	}
}

func TestYFastTrieEmptyOperations(t *testing.T) {
	trie := NewYFastTrie(1024)

	// Operations on empty trie
	if pred, found := trie.Predecessor(10); found {
		t.Errorf("Should not find predecessor in empty trie, got %d", pred)
	}

	if succ, found := trie.Successor(10); found {
		t.Errorf("Should not find successor in empty trie, got %d", succ)
	}

	if trie.Delete(10) {
		t.Error("Should not delete from empty trie")
	}
}

// Benchmark functions for performance testing

func BenchmarkYFastTrieInsert(b *testing.B) {
	trie := NewYFastTrie(1 << 20)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Insert(i)
	}
}

func BenchmarkYFastTriePredecessor(b *testing.B) {
	trie := NewYFastTrie(1 << 20)

	// Pre-populate with data
	for i := 0; i < 10000; i += 2 {
		trie.Insert(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Predecessor(i%20000 + 1) // Query for odd numbers
	}
}

func BenchmarkYFastTrieSuccessor(b *testing.B) {
	trie := NewYFastTrie(1 << 20)

	// Pre-populate with data
	for i := 0; i < 10000; i += 2 {
		trie.Insert(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Successor(i%20000 + 1) // Query for odd numbers
	}
}

func BenchmarkYFastTrieOperations(b *testing.B) {
	for i := 0; i < b.N; i++ {
		YFastTrieOperations(100)
	}
}

func BenchmarkXFastTrieCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewXFastTrie(20)
	}
}
