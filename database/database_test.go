package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const dbVersion = "3.6.1"

func TestImportFromJSON(t *testing.T) {
	// Remove Bleve index to load a new one
	os.RemoveAll(fmt.Sprintf("./jmdict_%s.bleve", dbVersion))

	// Initialize Mongo and Bleve
	database, err := NewDatabase("mongodb://localhost:27017", ".", dbVersion, 1000)
	if err != nil {
		t.Errorf("NewDatabase() error = %v", err)
	}

	// Clear collection to insert new data
	database.mongoCollection.Drop(context.Background())

	// Load the dictionary
	jsonName := fmt.Sprintf("jmdict-eng-common-%s.json", dbVersion)
	jsonPath := filepath.Join("..", "jmdict", jsonName)
	if err := database.ImportFromJSON(jsonPath); err != nil {
		t.Errorf("ImportFromJSON() error = %v", err)
	}
}

func TestSearch(t *testing.T) {
	database, err := NewDatabase("mongodb://localhost:27017", ".", dbVersion, 1000)
	if err != nil {
		t.Errorf("NewDatabase() error = %v", err)
	}

	results, err := database.Search("いや")
	if err != nil {
		t.Errorf("Search() error = %v", err)
	}

	t.Logf("Done! %v", results)
}
