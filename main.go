package main

import "fmt"

func main() {

	tree := tree{}

	tree.put([]byte("3"), []byte("4"))
	tree.put([]byte("5"), []byte("3"))
	tree.put([]byte("4"), []byte("1"))
	tree.put([]byte("8"), []byte("2"))
	tree.put([]byte("2"), []byte("9"))
	tree.put([]byte("9"), []byte("8"))

	fmt.Printf("%s", tree.get([]byte("8")))
}
