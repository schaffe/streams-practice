package index

import (
	"container/heap"
	"sort"
	"strings"

	"streams-practice/assignments/autocomplete/corpus"
)

type heapEntry struct {
	term string
	freq int
}

type maxHeap []heapEntry

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].freq < h[j].freq }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x any)        { *h = append(*h, x.(heapEntry)) }
func (h *maxHeap) Pop() any {
	old := *h
	n := len(old)
	e := old[n-1]
	*h = old[:n-1]
	return e
}

type node struct {
	children map[rune]*node
	top      *maxHeap
}

func (n *node) Add(e heapEntry) {
	if n.top == nil {
		n.top = &maxHeap{}
	}
	heap.Push(n.top, e)
	if n.top.Len() > 10 {
		heap.Pop(n.top)
	}
}

func (n *node) drain() []heapEntry {
	if n.top == nil || n.top.Len() == 0 {
		return nil
	}
	result := make([]heapEntry, len(*n.top))
	copy(result, *n.top)
	sort.Slice(result, func(i, j int) bool {
		return result[i].freq > result[j].freq
	})
	return result
}

type Trie struct {
	root *node
}

func New(entries []corpus.Entry) *Trie {
	t := &Trie{root: &node{children: make(map[rune]*node)}}
	for _, entry := range entries {
		term := strings.ToLower(entry.Term)
		he := heapEntry{entry.Term, entry.Frequency}
		t.root.Add(he)
		current := t.root
		for _, r := range term {
			if _, ok := current.children[r]; !ok {
				current.children[r] = &node{children: make(map[rune]*node)}
			}
			current = current.children[r]
			current.Add(he)
		}
	}
	return t
}

func (t *Trie) Search(prefix string) []corpus.Entry {
	current := t.root
	for _, r := range strings.ToLower(prefix) {
		if current.children == nil {
			return nil
		}
		child, ok := current.children[r]
		if !ok {
			return nil
		}
		current = child
	}
	drained := current.drain()
	if drained == nil {
		return nil
	}
	entries := make([]corpus.Entry, len(drained))
	for i, he := range drained {
		entries[i] = corpus.Entry{Term: he.term, Frequency: he.freq}
	}
	return entries
}
