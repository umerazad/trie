package trie

import (
	"fmt"
)

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

type Trie struct {
	root *trieNode
	size int
}

func New() *Trie {
	return &Trie{}
}

func (t *Trie) Size() int {
	return t.size
}

func (t *Trie) IsEmpty() bool {
	return t.Size() == 0
}

func (t *Trie) Depth() int {
	return depth(t.root)
}

func depth(n *trieNode) int {
	max := 0

	if n != nil {
		for k := range n.slots {
			d := 1 + depth(n.slots[k])
			if d > max {
				max += d
			}
		}
	}

	return max
}

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

func keyNotFoundError(k string) error {
	return fmt.Errorf("Key not found: %s", k)
}

func (t *Trie) Get(k string) (Value, error) {
	if t.IsEmpty() {
		return nil, keyNotFoundError(k)
	}

	n := t.root
	for _, ch := range k {
		if n.slots[ch] != nil {
			n = n.slots[ch]
		} else {
			return nil, keyNotFoundError(k)
		}
	}

	if n.v != nil {
		return n.v, nil
	}

	return nil, keyNotFoundError(k)
}

func (t *Trie) Contains(k string) bool {
	_, err := t.Get(k)
	return err == nil
}

func (t *Trie) Keys() []string {
	return nil
}

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
	if n.v != nil {
		keys = append(keys, string(prefix))
	}

	for k, _ := range n.slots {
		prefix = append(prefix, k)
		keys = append(keys, collectKeys(n.slots[k], prefix, nil)...)
		prefix = prefix[:len(prefix)-1]
	}

	return keys
}

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
