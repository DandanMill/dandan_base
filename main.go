package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	tree := mock_data()
	rand.NewSource(time.Now().Unix())

	for i := 0; i < 30; i++ {
		key := []byte(fmt.Sprintf("%d", rand.Intn(100)))
		fmt.Printf("%s -> %s\n", key, tree.get(key))
	}
}

func mock_data() *tree {
	f, err := os.Open("MOCK_DATA.csv")

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	tree := &tree{}

	for _, str := range data {
		tree.put([]byte(str[0]), []byte(str[1]))
	}
	return tree
}
