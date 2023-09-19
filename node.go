package main

import (
	"bytes"
	"sort"
)

type node struct {
	isLeaf   bool
	parent   *node
	kvs      kvs
	children nodes
}

func (n *node) put(key, value []byte) {
	index := sort.Search(len(n.kvs), func(i int) bool {
		return bytes.Compare(n.kvs[i].key, key) != -1
	})

	exact := len(n.kvs) > 0 && index < len(n.kvs) && bytes.Equal(n.kvs[index].key, key)
	if !exact {
		n.kvs = append(n.kvs, kv{})
		copy(n.kvs[index+1:], n.kvs[index:])
	}
	kv := &n.kvs[index]
	kv.key = key
	kv.value = value
}

func (n *node) internalPut(child *node, key []byte) {
	n.put(key, nil)

	index := sort.Search(len(n.kvs), func(i int) bool {
		return bytes.Compare(n.kvs[i].key, key) == 1
	})

	n.children = append(n.children, &node{})
	if index < len(n.children)-1 {
		copy(n.children[index+1:], n.children[index:])
	}

	n.children[index] = child
}

func (n *node) splitNode() (newNode *node) {

	half := len(n.kvs) / 2

	newNode = &node{
		isLeaf: n.isLeaf,
	}

	if newNode.isLeaf {
		newNode.kvs = append(newNode.kvs, n.kvs[half:]...)
		n.kvs = n.kvs[:half]
		return

	}
	newNode.kvs = append(newNode.kvs, n.kvs[half+1:]...)
	n.kvs = n.kvs[:half]

	half = len(n.children) / 2

	if MAX%2 == 0 {
		half++
	}
	newNode.children = append(newNode.children, n.children[half:]...)
	n.children = n.children[:half]

	return
}

type nodes []*node

type kv struct {
	key   []byte
	value []byte
}

type kvs []kv
