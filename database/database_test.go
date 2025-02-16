package database

import (
	"fmt"
	"testing"
)

func TestImportFromJSON(t *testing.T) {
	database, err := NewDatabase("mongodb://localhost:27017", 1000)
	if err != nil {
		t.Errorf("NewDatabase() error = %v", err)
	}

	if err := database.ImportFromJSON("../jmdict/jmdict-eng-common-3.6.1.json"); err != nil {
		t.Errorf("ImportFromJSON() error = %v", err)
	}
}

func TestSearch(t *testing.T) {
	database, err := NewDatabase("mongodb://localhost:27017", 1000)
	if err != nil {
		t.Errorf("NewDatabase() error = %v", err)
	}

	results, err := database.Search("いや")
	if err != nil {
		t.Errorf("Search() error = %v", err)
	}

	fmt.Printf("Done! %v", results)
}
