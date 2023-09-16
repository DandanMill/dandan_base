package main

import (
	"fmt"
)

func main() {
	t := tree{}

	t.put([]byte("3"), []byte("3"))
	t.put([]byte("5"), []byte("5"))
	t.put([]byte("1"), []byte("1"))
	t.put([]byte("2"), []byte("2"))
	fmt.Println(t.root.kvs)
	for _, n := range t.root.children {
		fmt.Println(n)
	}

}
