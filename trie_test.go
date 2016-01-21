package trie

import "testing"

func TestTrieCreation(t *testing.T) {
	tr := New()
	if tr.Size() != 0 {
		t.Error(`An empty trie should have zero size.`)
	}
}

func TestPutAndGet(t *testing.T) {
	tr := New()
	tr.Put("Key", "Value")

	if !tr.Contains("Key") {
		t.Error(`Trie must contain "Key"`)
	}

	if tr.Size() != 1 {
		t.Errorf(`Invalid Size reported. Expected :%d, Actual:d"`, 1, tr.Size())
	}

	v, err := tr.Get("Key")
	if err != nil {
		t.Errorf(`Get() failed unexpectedly. Error: %v"`, err)
	}

	if v != "Value" {
		t.Errorf(`Get() failed to return correct value. Expected: %s, Actual: %s`, "Value", v)
	}
}

func TestGetOnNonExistantKey(t *testing.T) {
	tr := New()

	v, err := tr.Get("NonExistant")
	if err == nil {
		t.Errorf("Get() succeeded unexpectedly. Returned: %s", v)
	}

	expected := "Key not found: NonExistant"
	if err.Error() != expected {
		t.Errorf("Unexpected error format. Expected: %s, Actual: %s", expected, err.Error())
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

	tr := New()
	for i := 0; i < len(cases); i++ {
		c := cases[i]
		tr.Put(c.key, c.val)
	}

	if tr.Size() != len(cases) {
		t.Errorf(`tr.Size() != len(cases)`)
	}

	for i := 0; i < len(cases); i++ {
		v, err := tr.Get(cases[i].key)

		if err != nil {
			t.Errorf("Unexpcted error: %v", err)
		}

		if v != cases[i].val {
			t.Error(`tr.Get(%q) != %q`, cases[i].key, cases[i].val)
		}
	}

	expected_depth := 4
	if tr.Depth() != expected_depth {
		t.Error("tr.Depth() != %d. Actual depth: %d", expected_depth, tr.Depth())
	}
}

func TestLongestPrefix(t *testing.T) {
	tr := New()

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
	tr := New()

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
