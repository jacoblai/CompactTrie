package test

import (
	"CompactTrie"
	"fmt"
	"testing"
)

func TestCompactTrie(t *testing.T) {
	keys := []string{
		"hello",
		"world",
		"hello world",
		"hello go",
		"golang",
		"go programming",
		"programming in go",
		"data structures",
		"algorithms",
	}

	trie := CompactTrie.NewCompactTrie(keys)

	// 功能完整性测试
	testCases := []struct {
		key   string
		exist bool
	}{
		{"hello", true},
		{"world", true},
		{"hello world", true},
		{"hello go", true},
		{"golang", true},
		{"go programming", true},
		{"programming in go", true},
		{"data structures", true},
		{"algorithms", true},
		{"hello world!", false},
		{"go", false},
		{"algorithm", false},
		{"data", false},
	}

	for _, tc := range testCases {
		if trie.Lookup(tc.key) != tc.exist {
			t.Errorf("Lookup(%q) = %v, want %v", tc.key, !tc.exist, tc.exist)
		}
	}
}

func BenchmarkCompactTrieBuild(b *testing.B) {
	keys := generateKeys(100000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = CompactTrie.NewCompactTrie(keys)
	}
}

func BenchmarkCompactTrieLookup(b *testing.B) {
	keys := generateKeys(100000)
	trie := CompactTrie.NewCompactTrie(keys)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		trie.Lookup(keys[i%len(keys)])
	}
}

func generateKeys(n int) []string {
	keys := make([]string, n)

	for i := 0; i < n; i++ {
		keys[i] = fmt.Sprintf("key%d", i)
	}

	return keys
}
