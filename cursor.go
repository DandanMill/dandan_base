package main

import (
	"bytes"
	"sort"
)

type cursor struct {
	current *node
}

func (c *cursor) searchNode(seek []byte) {
	for !c.current.isLeaf {
		index := sort.Search(len(c.current.kvs), func(i int) bool {
			return bytes.Compare(c.current.kvs[i].key, seek) == 1
		})
		c.current = c.current.children[index]
	}
}
