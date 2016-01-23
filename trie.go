package trie

// Value is a generic type that can be stored against any key in the trie.
type Value interface{}

type trieNode struct {
	v     Value
	slots map[rune]*trieNode
}

func newTrieNode() *trieNode {
	n := &trieNode{}
	n.slots = make(map[rune]*trieNode)

	return n
}

// Trie represents an efficient Trie for unicode characters.
type Trie struct {
	root *trieNode
	size int
}

// NewTrie creates a new instance of Trie.
func NewTrie() *Trie {
	return &Trie{}
}

// Size returns the number of keys currently stored in the trie.
func (t *Trie) Size() int {
	return t.size
}

// IsEmpty returns true if the trie is empty.
func (t *Trie) IsEmpty() bool {
	return t.Size() == 0
}

// Depth returns the number of levels in the trie. Its effectively
// the length of the longest keys stored in the Trie.
func (t *Trie) Depth() int {
	return depth(t.root)
}

func depth(n *trieNode) int {
	max := 0

	if n != nil {
		for k := range n.slots {
			if n.slots[k] != nil {
				d := 1 + depth(n.slots[k])
				if d > max {
					max = d
				}
			}
		}
	}

	return max
}

// Put inserts a key, value pair in the trie.
func (t *Trie) Put(k string, v Value) {
	if t.IsEmpty() {
		t.root = newTrieNode()
	}

	n := t.root
	for _, ch := range k {
		if n.slots[ch] != nil {
			n = n.slots[ch]
		} else {
			n.slots[ch] = newTrieNode()
			n = n.slots[ch]
		}
	}

	n.v = v
	t.size++
}

// Get retries a value against a given key from the trie.
func (t *Trie) Get(k string) (Value, bool) {
	if t.IsEmpty() {
		return nil, false
	}

	n := t.root
	for _, ch := range k {
		if n.slots[ch] != nil {
			n = n.slots[ch]
		} else {
			return nil, false
		}
	}

	if n.v != nil {
		return n.v, true
	}

	return nil, false
}

// Contains returns trie if the specified key exists in the trie.
func (t *Trie) Contains(k string) bool {
	_, ok := t.Get(k)
	return ok
}

// Keys returns all the keys stored in the trie.
func (t *Trie) Keys() []string {
	return t.KeysWithPrefix("")
}

// KeysWithPrefix returns all the keys starting with the specified prefix
// in the trie.
func (t *Trie) KeysWithPrefix(prefix string) []string {
	if t.IsEmpty() {
		return nil
	}

	n := t.root
	for _, ch := range prefix {
		if n.slots[ch] == nil {
			// No key exists with full prefix.
			return nil
		}

		n = n.slots[ch]
	}

	return collectKeys(n, []rune(prefix), nil)
}

func collectKeys(n *trieNode, prefix []rune, keys []string) []string {
	if n != nil {

		if n.v != nil {
			keys = append(keys, string(prefix))
		}

		for k := range n.slots {
			prefix = append(prefix, k)
			keys = append(keys, collectKeys(n.slots[k], prefix, nil)...)
			prefix = prefix[:len(prefix)-1]
		}
	}

	return keys
}

// LongestPrefix returns the longest valid prefix for the given query
// that exists in trie.
func (t *Trie) LongestPrefix(query string) string {
	length := -1

	n := t.root
	for _, ch := range query {
		if n.slots[ch] != nil {
			length++
			n = n.slots[ch]
			continue
		}
	}

	return query[0 : length+1]
}

// Delete removes a key and its corresponding value from the trie.
func (t *Trie) Delete(key string) {
	t.root = t.deleteInternal(t.root, []rune(key), 0)
}

func (t *Trie) deleteInternal(n *trieNode, key []rune, index int) *trieNode {
	// Nothing to delete.
	if n == nil {
		return nil
	}

	// Delete any value if we found the exact key.
	if index == len(key) {
		if n.v != nil {
			n.v = nil
			t.size--
		}
	} else {
		// Try deletion at next level.
		ch := key[index]
		n.slots[ch] = t.deleteInternal(n.slots[ch], key, index+1)
	}

	if n.v != nil {
		return n
	}

	// Remove empty sub-trie.
	for ch := range n.slots {
		if n.slots[ch] != nil {
			return n
		}
	}

	return nil
}

// KeysWithFuzzyMatch returns all the keys that match the given pattern.
// A dot '.' matches any character. For example, 'd..e' will match 'date' and 'done'.
func (t *Trie) KeysWithFuzzyMatch(pattern string) []string {
	if t.IsEmpty() {
		return nil
	}

	return collectKeysWithFuzzyMatch(t.root, nil, []rune(pattern), nil)
}

func collectKeysWithFuzzyMatch(n *trieNode, prefix []rune, pattern []rune, keys []string) []string {
	if n == nil {
		return nil
	}

	if len(prefix) == len(pattern) {
		if n.v != nil {
			keys = append(keys, string(prefix))
		}

		return keys
	}

	ch := pattern[len(prefix)]

	if ch == rune('.') {
		for i := range n.slots {
			if n.slots[i] != nil {
				prefix = append(prefix, i)
				keys = append(keys, collectKeysWithFuzzyMatch(n.slots[i], prefix, pattern, nil)...)
				prefix = prefix[:len(prefix)-1]
			}
		}
	} else {
		prefix = append(prefix, ch)
		keys = append(keys, collectKeysWithFuzzyMatch(n.slots[ch], prefix, pattern, nil)...)
		prefix = prefix[:len(prefix)-1]
	}

	return keys
}
