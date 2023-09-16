package main

import (
	"bytes"
	"sort"
)

type node struct {
	isLeaf   bool
	kvs      kvs
	children nodes
}

func (n *node) put(oldKey, key, value []byte) {

	index := sort.Search(len(n.kvs), func(i int) bool { return bytes.Compare(n.kvs[i].key, oldKey) != -1 })

	exact := len(n.kvs) > 0 && index < len(n.kvs) && bytes.Equal(n.kvs[index].key, oldKey)

	if !exact {
		n.kvs = append(n.kvs, kv{})
		copy(n.kvs[index+1:], n.kvs[index:])
	}

	kv := &n.kvs[index]
	kv.key = key
	kv.value = value
}

func (n *node) insertInternal(children []*node) {

	lastNode := children[len(children)-1]
	key := lastNode.kvs[0].key
	n.put(key, key, nil)

	for _, child := range children {
		n.children = append(n.children, child)
	}
}

func (n *node) splitNode() *node {
	half := len(n.kvs) >> 1
	newNode := &node{
		isLeaf: n.isLeaf,
	}

	newNode.kvs = n.kvs[half:]
	n.kvs = n.kvs[:half]

	half = len(n.children) >> 1
	newNode.children = n.children[half:]
	n.children = n.children[:half]

	return newNode
}

type nodes []*node

type kv struct {
	key   []byte
	value []byte
}

type kvs []kv
