package dandan_base

import (
	"bytes"
	"sort"
)

type node struct {
	isLeaf   bool
	parent   *node
	sibling  *node
	kvs      kvs
	children nodes
}

func (n *node) put(key, value []byte, pgid pgid) {
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
	kv.pgid = pgid
}

func (n *node) internalPut(child *node, key []byte) {
	n.put(key, nil, 0)

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
		newNode.sibling = n.sibling
		n.sibling = newNode
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

func (n *node) write(p *page) {

}

type nodes []*node

type kv struct {
	key   []byte
	value []byte
	pgid  pgid
}

type kvs []kv
