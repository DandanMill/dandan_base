package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/DandanMill/dandan-base/v1/dandan_base"
)

func main() {
	db := mock_data()

	db.Write()

	p := db.Read()

	fmt.Printf("%+v", p)
}

func mock_data() *dandan_base.DB {
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
	tree := &dandan_base.Tree{}

	for _, str := range data {
		tree.Put([]byte(str[0]), []byte(str[1]))
	}
	return &dandan_base.DB{Bucket: tree}
}
