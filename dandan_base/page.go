package dandan_base

import "unsafe"

const (
	branchPageFlag = 0x01
	leafPageFlag   = 0x02
)

const LeafPageElementSize = unsafe.Sizeof(leafPageElement{})
const BranchPageElementSize = unsafe.Sizeof(branchPageElement{})

type pgid uint64

type page struct {
	id    pgid
	flags uint16
	count uint16
}

func (p *page) PageElementSize() uintptr {
	if p.flags == leafPageFlag {
		return LeafPageElementSize
	} else {
		return BranchPageElementSize
	}
}

func (p *page) LeafPageELement(idx uint16) *leafPageElement {
	return (*leafPageElement)(UnsafeIndex(unsafe.Pointer(p), unsafe.Sizeof(*p),
		LeafPageElementSize, int(idx)))
}

func (p *page) BranchPageElement(idx uint16) *branchPageElement {
	return (*branchPageElement)(UnsafeIndex(unsafe.Pointer(p), unsafe.Sizeof(*p),
		BranchPageElementSize, int(idx)))

}

type branchPageElement struct {
	pos   uint32
	ksize uint32
	pgid  pgid
}

type leafPageElement struct {
	pos   uint32
	ksize uint32
	vsize uint32
}
