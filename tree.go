package main

import (
	"bytes"
	"sort"
)

const MAX_DEGREE = 3

type tree struct {
	root *node
}

func (t *tree) put(key, value []byte) {
	if t.root == nil {
		t.root = &node{
			isLeaf: true,
		}
	}
	cursor := t.searchNode(key)
	cursor.put(key, key, value)

	if len(cursor.kvs) == MAX_DEGREE {
		next := cursor.splitNode()
		if cursor == t.root {
			newRoot := &node{
				isLeaf: false,
			}

			newRoot.insertInternal([]*node{cursor, next})

			t.root = newRoot
		}
	}
}

func (t *tree) searchNode(key []byte) *node {
	cursor := t.root
	for !cursor.isLeaf {
		var exact bool
		index := sort.Search(len(cursor.kvs), func(i int) bool {
			ret := bytes.Compare(cursor.kvs[i].key, key)
			if bytes.Equal(cursor.kvs[i].key, key) {
				exact = true
			}
			return ret != -1
		})
		if exact {
			index++
		}
		cursor = cursor.children[index]
	}
	return cursor
}
