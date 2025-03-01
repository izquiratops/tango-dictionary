package main

import (
	"fmt"
	"tango/database"
	"testing"
)

const dbVersion = "3.6.1"

func TestNewDatabase(t *testing.T) {
	// Initialize Mongo and Bleve
	db, err := database.NewDatabase("mongodb://localhost:27017", dbVersion, 1000, true)
	if err != nil {
		t.Errorf("NewDatabase() error = %v", err)
	}

	// Load the dictionary
	jsonName := fmt.Sprintf("jmdict-eng-common-%s.json", dbVersion)
	if err := db.ImportFromJSON(jsonName); err != nil {
		t.Errorf("ImportFromJSON() error = %v", err)
	}
}
