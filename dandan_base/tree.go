package dandan_base

import (
	"bytes"
	"sort"
)

const MaxAllocSize = 0xFFFFFFF

const MAX = 100

type Tree struct {
	root *node
}

func (t *Tree) Put(key, value []byte) {
	if t.root == nil {
		t.root = &node{
			isLeaf: true,
		}
	}
	cursor := &cursor{current: t.root}
	cursor.searchNode(key)

	current := cursor.current
	current.put(key, value, 0)

	stackPointer := len(cursor.stack)
	iterationCount := 1

	for len(current.kvs) == MAX {

		middleKey := make([]byte, len(current.kvs[len(current.kvs)>>1].key))
		copy(middleKey, current.kvs[len(current.kvs)>>1].key)

		next := current.splitNode()

		if t.root == current {
			newRoot := &node{isLeaf: false}
			newRoot.put(middleKey, nil, 0)

			newRoot.children = append(newRoot.children, []*node{current, next}...)
			current.parent = newRoot
			next.parent = newRoot

			t.root = newRoot
		} else {
			parent := cursor.stack[stackPointer-iterationCount]
			parent.internalPut(next, middleKey)
			current = parent
			iterationCount++
		}
	}
}

func (t *Tree) Get(key []byte) []byte {
	cursor := &cursor{current: t.root}
	cursor.searchNode(key)

	node := cursor.current

	index := sort.Search(len(node.kvs), func(i int) bool {
		return bytes.Compare(node.kvs[i].key, key) != -1
	})

	if index > len(node.kvs)-1 || index < 0 {
		return []byte("No such key")
	}

	if bytes.Equal(node.kvs[index].key, key) {
		return node.kvs[index].value
	}

	return []byte("No such key")
}
