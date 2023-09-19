package main

import (
	"bytes"
	"sort"
)

type cursor struct {
	current *node
	stack   []*node
}

func (c *cursor) searchNode(seek []byte) {
	c.stack = c.stack[:0]
	for !c.current.isLeaf {
		c.stack = append(c.stack, c.current)
		index := sort.Search(len(c.current.kvs), func(i int) bool {
			return bytes.Compare(c.current.kvs[i].key, seek) == 1
		})
		c.current = c.current.children[index]
	}
}
