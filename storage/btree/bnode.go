package btree

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type BNode struct {
	data []byte
}

// Header
func (node *BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node.data)
}

func (node *BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node.data[2:4])
}

func (node *BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node.data[0:2], btype)
	binary.LittleEndian.PutUint16(node.data[2:4], nkeys)
}

// Pointers
func (node *BNode) getPtr(idx uint16) uint64 {
	if idx > node.nkeys() {
		panic(fmt.Sprintf("%d exceeds number of keys", idx))
	}
	pos := HEADER + 8*idx
	return binary.LittleEndian.Uint64(node.data[pos:])
}

func (node *BNode) setPtr(idx uint16, val uint64) {
	if idx > node.nkeys() {
		panic(fmt.Sprintf("%d exceeds number of keys", idx))
	}
	pos := HEADER + 8*idx
	binary.LittleEndian.PutUint64(node.data[pos:], val)
}

// Offset
func offsetPos(node *BNode, idx uint16) uint16 {
	if 0 < idx || idx > node.nkeys() {
		panic(fmt.Sprintf("%d exceeds number of keys", idx))
	}
	return HEADER + 8*node.nkeys() + 2*(idx-1)
}

func (node *BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	return binary.LittleEndian.Uint16(node.data[offsetPos(node, idx):])
}

func (node *BNode) setOffset(idx uint16, offset uint16) {
	binary.LittleEndian.PutUint16(node.data[offsetPos(node, idx):], offset)
}

// key-values
func (node *BNode) kvPos(idx uint16) uint16 {
	if idx > node.nkeys() {
		panic(fmt.Sprintf("%d exceeds number of keys", idx))
	}
	return HEADER + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}

func (node *BNode) getKey(idx uint16) []byte {
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node.data[pos:])

	return node.data[pos+4:][:klen]
}

func (node *BNode) getVal(idx uint16) []byte {
	pos := node.kvPos(idx)
	vlen := binary.LittleEndian.Uint16(node.data[pos+2:])

	return node.data[pos+4:][:vlen]
}

func (node *BNode) nbytes() uint16 {
	return node.kvPos(node.nkeys())
}

func nodeLookupLE(node *BNode, key []byte) uint16 {
	// nkeys := node.nkeys()
	found := uint16(0)

	for i := uint16(1); bytes.Compare(node.getKey(i), key) <= 0; i++ {
		found = i
	}

	return found
}
