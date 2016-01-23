package trie

import "testing"

func TestTrieCreation(t *testing.T) {
	tr := NewTrie()
	if tr.Size() != 0 {
		t.Error(`An empty trie should have zero size.`)
	}
}

func TestPutAndGet(t *testing.T) {
	tr := NewTrie()
	tr.Put("Key", "Value")

	if !tr.Contains("Key") {
		t.Error(`Trie must contain "Key"`)
	}

	if tr.Size() != 1 {
		t.Errorf(`Invalid Size reported. Expected :%d, Actual:%d"`, 1, tr.Size())
	}

	v, ok := tr.Get("Key")
	if !ok {
		t.Errorf(`Get("Key") failed to find the value.`)
	}

	if v != "Value" {
		t.Errorf(`Get("Key") failed to return correct value. Expected: %s, Actual: %s`, "Value", v)
	}
}

func TestGetOnNonExistantKey(t *testing.T) {
	tr := NewTrie()

	v, ok := tr.Get("NonExistant")
	if ok {
		t.Errorf("Get() succeeded unexpectedly. Returned: %s", v)
	}
}

func TestMultipleGetAndPutWithCommonPrefix(t *testing.T) {
	cases := []struct {
		key string
		val int
	}{
		{"ABC", 1},
		{"ABCA", 2},
		{"ABCB", 3},
		{"ABCC", 4},
		{"ABCD", 5},
		{"ABCE", 6},
		{"ABCF", 7},
		{"ABCF", 7},
	}

	tr := NewTrie()
	for i := 0; i < len(cases); i++ {
		c := cases[i]
		tr.Put(c.key, c.val)
	}

	if tr.Size() != len(cases) {
		t.Errorf(`tr.Size() != len(cases)`)
	}

	for i := 0; i < len(cases); i++ {
		v, ok := tr.Get(cases[i].key)

		if !ok {
			t.Errorf(`Get(%q) failed.`, cases[i].key)
		}

		if v != cases[i].val {
			t.Errorf("tr.Get(%q) != %d", cases[i].key, cases[i].val)
		}
	}

	expectedDepth := 4
	if tr.Depth() != expectedDepth {
		t.Errorf("tr.Depth() != %d. Actual depth: %d", expectedDepth, tr.Depth())
	}
}

func TestLongestPrefix(t *testing.T) {
	tr := NewTrie()

	tr.Put("A", true)
	tr.Put("AB", true)
	tr.Put("ABC", true)
	tr.Put("ZZZZABC", true)

	cases := []struct {
		input    string
		expected string
	}{
		{"A", "A"},
		{"AB", "AB"},
		{"ABC", "ABC"},
		{"ABCD", "ABC"},
		{"BCD", ""},
		{"ZZZZABD", "ZZZZAB"},
	}

	for i := 0; i < len(cases); i++ {
		actual := tr.LongestPrefix(cases[i].input)
		if actual != cases[i].expected {
			t.Errorf(`tr.LongestPrefix("%s") != "%s" : Actual "%s"`,
				cases[i].input,
				cases[i].expected,
				actual)
		}
	}
}

func TestKeysWithPrefix(t *testing.T) {
	tr := NewTrie()

	tr.Put("A", true)
	tr.Put("AB", true)
	tr.Put("ABC", true)
	tr.Put("ZZZZABC", true)

	cases := []struct {
		input string
		count int
	}{
		{"A", 3},
		{"AB", 2},
		{"ABC", 1},
		{"ABCD", 0},
		{"BCD", 0},
		{"Z", 1},
		{"", 4},
	}

	for i := 0; i < len(cases); i++ {
		actual := tr.KeysWithPrefix(cases[i].input)

		if len(actual) != cases[i].count {
			t.Errorf(`Expected prefix count: %d, Actual count: %d`, cases[i].count, len(actual))
		}
	}
}

func TestDepth(t *testing.T) {
	tr := NewTrie()

	cases := []struct {
		input         string
		expectedDepth int
	}{
		{"A", 1},
		{"ABC", 3},
		{"AB", 3},
		{"ABCD", 4},
		{"TTTTTT", 6},
	}

	for i := 0; i < len(cases); i++ {
		tr.Put(cases[i].input, struct{}{})
		if tr.Depth() != cases[i].expectedDepth {
			t.Errorf("Expected depth: %d, Actual depth: %d", cases[i].expectedDepth, tr.Depth())
		}
	}

	cases = []struct {
		input         string
		expectedDepth int
	}{
		{"TTTTTT", 4},
		{"A", 4},
		{"ABCD", 3},
		{"AB", 3},
		{"ABC", 0},
	}

	for i := 0; i < len(cases); i++ {
		tr.Delete(cases[i].input)

		if tr.Depth() != cases[i].expectedDepth {
			t.Errorf("After deleting: %q, Expected depth: %d, Actual depth: %d",
				cases[i].input,
				cases[i].expectedDepth,
				tr.Depth())
		}
	}
}

func TestDelete(t *testing.T) {
	tr := NewTrie()

	tr.Put("A", 1)
	tr.Put("AB", 1)
	tr.Put("ABC", 1)
	tr.Put("ABCD", 1)
	tr.Put("ABCDEFGH", 1)

	cases := []struct {
		input         string
		expectedDepth int
	}{
		{"ZZZZZZZZZZZZZZZZZZZZZZZZ", 8}, // Failed deletion doesn't change depth.
		{"ABCD", 8},                     // Deletion successful but subkey exists.
		{"ABCDEFGH", 3},                 // Longest key left would be 'ABC'
		{"ABC", 2},
		{" ABCDEF", 2},
		{"A", 2},
		{"AB", 0},
	}

	for i := 0; i < len(cases); i++ {
		size := tr.Size()
		expectedSize := size - 1

		if !tr.Contains(cases[i].input) {
			expectedSize++
		}

		tr.Delete(cases[i].input)

		if tr.Contains(cases[i].input) {
			t.Errorf(`tr.Contains(%q) != false`, cases[i].input)
		}

		if tr.Size() != expectedSize {
			t.Errorf("Expected size: %d, Actual size: %d", expectedSize, tr.Size())
		}
	}

	if !tr.IsEmpty() {
		t.Error("Expected an empty trie by now.")
	}

	if tr.Depth() != 0 {
		t.Error("Depth of an empty trie must be 0.")
	}
}

func TestKeysWithFuzzyMatch(t *testing.T) {
	tr := NewTrie()

	tr.Put("A", true)
	tr.Put("AB", true)
	tr.Put("BC", true)
	tr.Put("ABC", true)
	tr.Put("CAC", true)
	tr.Put("ZZZZABC", true)

	cases := []struct {
		input string
		count int
	}{
		{"A", 1},
		{"..", 2},
		{"AB", 1},
		{".B", 1},
		{"B.", 1},
		{"ABC", 1},
		{"A.C", 1},
		{"..C", 2},
		{"AB..", 0},
		{"Z", 0},
	}

	for i := 0; i < len(cases); i++ {
		actual := tr.KeysWithFuzzyMatch(cases[i].input)

		if len(actual) != cases[i].count {
			t.Errorf(`Expected prefix count: %d, Actual count: %d`, cases[i].count, len(actual))
		}
	}
}
