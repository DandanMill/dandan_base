package main

import (
	"bytes"
	"sort"
)

const MAX = 3

type tree struct {
	root *node
}

func (t *tree) put(key, value []byte) {
	if t.root == nil {
		t.root = &node{
			isLeaf: true,
		}
	}
	cursor := &cursor{current: t.root}
	cursor.searchNode(key)

	current := cursor.current
	current.put(key, value)

	for len(current.kvs) == MAX {
		key := current.kvs[len(current.kvs)>>1].key

		next := current.splitNode()
		if t.root == current {
			newRoot := &node{isLeaf: false}
			newRoot.put(key, nil)

			newRoot.children = append(newRoot.children, []*node{current, next}...)
			current.parent = newRoot
			next.parent = newRoot

			t.root = newRoot
		} else {
			parent := current.parent
			parent.internalPut(next)
			current = parent
		}
	}
}

func (t *tree) get(key []byte) []byte {
	cursor := &cursor{current: t.root}
	cursor.searchNode(key)

	node := cursor.current

	index := sort.Search(len(node.kvs), func(i int) bool {
		return bytes.Compare(node.kvs[i].key, key) != -1
	})

	if bytes.Equal(cursor.current.kvs[index].key, key) {
		return cursor.current.kvs[index].value
	}

	return nil
}
