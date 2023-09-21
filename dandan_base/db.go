package dandan_base

import (
	"log"
	"os"
	"unsafe"
)

type DB struct {
	Bucket *Tree
}

func (db *DB) Write() {
	file, err := os.OpenFile("db.db", os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		return
	}
	p := &page{}
	db.Bucket.root.write(p)

	buf := UnsafeByteSlice(unsafe.Pointer(p), 0, 0, 4096)

	file.Write(buf)

}

func (db *DB) Read() *page {
	file, err := os.OpenFile("db.db", os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		log.Fatal("Couldn't open db")
	}

	buf := make([]byte, 4096)
	file.Read(buf)

	return (*page)(unsafe.Pointer(&buf))
}
