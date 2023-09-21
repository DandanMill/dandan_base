package dandan_base

import "unsafe"

func UnsafeAdd(base unsafe.Pointer, offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(base) + offset)
}

func UnsafeIndex(base unsafe.Pointer, offset uintptr, elemz uintptr, n int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(base) + offset + uintptr(n)*elemz)
}

func UnsafeByteSlice(base unsafe.Pointer, offset uintptr, i, j int) []byte {
	return (*[MaxAllocSize]byte)(UnsafeAdd(base, offset))[i:j:j]
}
