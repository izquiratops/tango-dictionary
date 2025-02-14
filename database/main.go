package main

import (
	"fmt"
	"log"
)

func main() {
	importer, err := NewDatabase("mongodb://localhost:27017", 1000)
	if err != nil {
		log.Fatal(err)
	}

	if err := importer.ImportFromJSON("./jmdict/test.json"); err != nil {
		log.Fatal(err)
	}

	results, err := importer.Search("いや")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v", results)
}
