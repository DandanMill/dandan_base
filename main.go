package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	tree := mock_data()

	fmt.Printf("%s", tree.get([]byte("150")))
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
