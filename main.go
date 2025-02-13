package main

import (
	"fmt"
	"log"
)

func main() {
	importer, err := NewDictionaryImporter("mongodb://localhost:27017", 1000)
	if err != nil {
		log.Fatal(err)
	}

	if err := importer.ImportFromJSON("./source/test-entry.json"); err != nil {
		log.Fatal(err)
	}

	results, err := importer.Search("プレ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v", results)
}
