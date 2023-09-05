package btree

import (
	"bytes"
	"encoding/binary"
)

const (
	BNODE_NODE         = 1
	BNODE_LEAF         = 2
	HEADER             = 4
	BTREE_PAGE_SIZE    = 4096
	BTREE_MAX_KEY_SIZE = 1000
	BTREE_MAX_VAL_SIZE = 3000
)

type BTree struct {
	Root uint64

	get func(uint64) BNode
	new func(BNode) uint64
	del func(uint64)
}

func leafInsert(new *BNode, old *BNode, idx uint16, key []byte, val []byte) {
	new.setHeader(BNODE_LEAF, old.nkeys()+1)
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKV(new, idx, 0, key, val)
}

func nodeAppendRange(new *BNode, old *BNode, dstNew uint16, srcOld uint16, n uint16) {
	if n == 0 {
		return
	}

	//pointers
	for i := uint16(0); i < n; i++ {
		new.setPtr(dstNew+i, old.getPtr(srcOld+i))
	}

	dstBegin := new.getOffset(dstNew)
	srcBegin := new.getOffset(srcOld)

	//offsets
	for i := uint16(1); i <= n; i++ { //NOTE: the range is [1,n]
		offset := dstBegin + old.getOffset(srcBegin+i)
		new.setOffset(dstNew+i, offset)
	}

	//kvs
	begin := old.kvPos(srcOld)
	end := old.kvPos(srcOld + n)
	copy(new.data[new.kvPos(dstNew):], old.data[begin:end])
}

func nodeAppendKV(new *BNode, idx uint16, ptr uint64, key []byte, val []byte) {
	//pointers
	new.setPtr(idx, ptr)

	//KVs
	pos := new.kvPos(idx)
	binary.LittleEndian.PutUint16(new.data[pos:], uint16(len(key)))
	binary.LittleEndian.PutUint16(new.data[pos+2:], uint16(len(val)))
	copy(new.data[pos+4:], key)
	copy(new.data[pos+4+uint16(len(key)):], val)

	new.setOffset(idx+1, new.getOffset(idx)+4+uint16(len(key)))
}

func treeInsert(tree *BTree, node *BNode, key []byte, val []byte) BNode {
	new := BNode{data: make([]byte, 2*BTREE_PAGE_SIZE)}

	idx := nodeLookupLE(node, key)

	switch node.btype() {
	case BNODE_LEAF:
		if bytes.Equal(key, node.getKey(idx)) {

			leafInsert(&new, node, idx, key, val)
		} else {
			leafInsert(&new, node, idx+1, key, val)
		}
	case BNODE_NODE:
		nodeInsert(tree, &new, node, idx, key, val)
	default:
		panic("Bad node")
	}
	return new
}

func nodeInsert(tree *BTree, new *BNode, node *BNode, idx uint16, key []byte, val []byte) {

	//get and deallocate the kid node
	kptr := node.getPtr(idx)
	knode := tree.get(kptr)
	tree.del(kptr)

	knode = treeInsert(tree, &knode, key, val)

	nsplit, splited := nodeSplit3(knode)

	nodeReplaceKidN(tree, new, node, idx, splited[:nsplit]...)
}

func nodeSplit2(left BNode, right BNode, old BNode) {
	right.data = old.data[len(old.data)-BTREE_PAGE_SIZE:]
	left.data = old.data[:len(old.data)-BTREE_PAGE_SIZE]
}

func nodeSplit3(old BNode) (uint16, [3]BNode) {
	if old.nbytes() <= BTREE_PAGE_SIZE {
		old.data = old.data[:BTREE_PAGE_SIZE]
		return 1, [3]BNode{old}
	}
	left := BNode{make([]byte, 2*BTREE_PAGE_SIZE)}
	right := BNode{make([]byte, BTREE_PAGE_SIZE)}

	nodeSplit2(left, right, old)

	if left.nbytes() <= BTREE_PAGE_SIZE {
		left.data = left.data[:BTREE_PAGE_SIZE]
		return 2, [3]BNode{left, right}
	}
	leftleft := BNode{make([]byte, BTREE_PAGE_SIZE)}
	middle := BNode{make([]byte, BTREE_PAGE_SIZE)}
	nodeSplit2(leftleft, middle, left)

	return 3, [3]BNode{leftleft, middle, right}
}

func nodeReplaceKidN(tree *BTree, new *BNode, old *BNode, idx uint16, kids ...BNode) {
	inc := uint16(len(kids))
	new.setHeader(BNODE_NODE, old.nkeys()+inc-1)
	nodeAppendRange(new, old, 0, 0, idx)
	for i, node := range kids {
		nodeAppendKV(new, idx+uint16(i), tree.new(node), node.getKey(0), nil)
	}
	nodeAppendRange(new, old, idx+inc, idx+1, old.nkeys()-(idx+1))
}

type C struct {
	tree  BTree
	ref   map[string]string
	pages map[uint64]BNode
}
