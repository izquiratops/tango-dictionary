package database

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestImportFromJSON(t *testing.T) {
	// Remove Bleve index to load new ones
	os.RemoveAll("./jmdict.bleve")

	// Initialize Mongo and Bleve
	database, err := NewDatabase("mongodb://localhost:27017", "./jmdict.bleve", 1000)
	if err != nil {
		t.Errorf("NewDatabase() error = %v", err)
	}

	// Clear collection to insert data again
	database.mongoCollection.Drop(context.Background())

	// Load the dictionary
	if err := database.ImportFromJSON("jmdict-eng-common-3.6.1.json"); err != nil {
		t.Errorf("ImportFromJSON() error = %v", err)
	}
}

func TestSearch(t *testing.T) {
	database, err := NewDatabase("mongodb://localhost:27017", "./jmdict.bleve", 1000)
	if err != nil {
		t.Errorf("NewDatabase() error = %v", err)
	}

	results, err := database.Search("いや")
	if err != nil {
		t.Errorf("Search() error = %v", err)
	}

	fmt.Printf("Done! %v", results)
}
