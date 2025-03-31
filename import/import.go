package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/izquiratops/tango/common/database"
	"github.com/izquiratops/tango/common/jmdict"
	"github.com/izquiratops/tango/common/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	numWorkers = 3
	batchSize  = 1000
)

func Import(db *database.Database, config types.ServerConfig) (string, error) {
	jsonFilename := fmt.Sprintf("jmdict-eng-%v.json", config.JmdictVersion)
	jsonPath := filepath.Join("..", "jmdict_source", jsonFilename)

	jsonSource, err := parseJmdictFile(jsonPath)
	if err != nil {
		return "", err
	}

	if err := importTags(db, jsonSource); err != nil {
		return "", fmt.Errorf("error importing tags: %v", err)
	}
	fmt.Printf("Tags import completed.\n")

	startTime := time.Now()
	if err := importWords(db, jsonSource.Words); err != nil {
		return "", err
	}

	fmt.Printf("Dictionary import completed. Processed %d entries in %v\n", len(jsonSource.Words), time.Since(startTime))
	return jsonPath, nil
}

// Reads and parses a JMdict JSON file
func parseJmdictFile(jsonPath string) (*jmdict.JMdict, error) {
	file, err := os.Open(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var jsonSource jmdict.JMdict
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&jsonSource); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return &jsonSource, nil
}

// Handles the batch import of tags into MongoDB
func importTags(db *database.Database, jsonSource *jmdict.JMdict) error {
	fmt.Println("Importing tags...")

	if len(jsonSource.Tags) == 0 {
		fmt.Println("No tags found in the dictionary.")
		return nil
	}

	var tagEntries []mongo.WriteModel
	for tagName, description := range jsonSource.Tags {
		tagEntry := database.Tag{
			Name:        tagName,
			Description: description,
		}

		model := mongo.NewInsertOneModel().SetDocument(tagEntry)
		tagEntries = append(tagEntries, model)
	}

	ctx := context.Background()
	_, err := db.MongoTags.BulkWrite(ctx, tagEntries)
	if err != nil {
		return fmt.Errorf("error importing tags to MongoDB: %v", err)
	}

	return nil
}

// Handles the worker pool setup and coordination for importing word entries
func importWords(db *database.Database, words []jmdict.JMdictWord) error {
	fmt.Println("Importing words...")
	entriesChan := make(chan jmdict.JMdictWord, batchSize)
	errorsChan := make(chan error, numWorkers)
	doneChan := make(chan struct{})

	// Start worker pool
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go bulkImportJmdictEntries(entriesChan, errorsChan, &wg, db)
	}

	// Monitor for errors from workers in a separate goroutine
	var importErr error
	go func() {
		wg.Wait()
		close(doneChan)
	}()

	// Send entries to workers
	for _, entry := range words {
		select {
		case err := <-errorsChan:
			return fmt.Errorf("worker error: %v", err)
		case entriesChan <- entry:
			// Entry sent to worker
		}
	}

	close(entriesChan)

	// Wait for completion or error
	select {
	case err := <-errorsChan:
		if importErr == nil {
			importErr = fmt.Errorf("worker error: %v", err)
		}
	case <-doneChan:
		// All workers completed successfully
	}

	return importErr
}

func bulkImportJmdictEntries(jsonEntries <-chan jmdict.JMdictWord, errors chan<- error, wg *sync.WaitGroup, di *database.Database) {
	defer wg.Done()

	ctx := context.Background()
	mongoBatch := make([]mongo.WriteModel, 0, batchSize)
	bleveBatch := di.BleveIndex.NewBatch()

	for jsonEntry := range jsonEntries {
		// Save it as DatabaseEntry
		dbEntry := ToWord(&jsonEntry)

		// Prepare MongoDB
		bsonData, err := bson.Marshal(dbEntry)
		if err != nil {
			errors <- fmt.Errorf("error marshalling to BSON: %v", err)
			return
		}

		model := mongo.NewInsertOneModel().SetDocument(bsonData)
		mongoBatch = append(mongoBatch, model)

		// Prepare Bleve
		bleveEntry, err := ToWordSearchable(&jsonEntry)
		if err != nil {
			errors <- err
			return
		}

		if err := bleveBatch.Index(jsonEntry.ID, bleveEntry); err != nil {
			errors <- fmt.Errorf("error indexing in Bleve: %v", err)
			return
		}

		if len(mongoBatch) >= batchSize {
			if _, err := di.MongoWords.BulkWrite(ctx, mongoBatch); err != nil {
				errors <- fmt.Errorf("error writing to MongoDB: %v", err)
				return
			}

			if err := di.BleveIndex.Batch(bleveBatch); err != nil {
				errors <- fmt.Errorf("error writing to Bleve: %v", err)
				return
			}

			mongoBatch = mongoBatch[:0]
			bleveBatch = di.BleveIndex.NewBatch()
		}
	}

	// Process last batch. This runs the batch if mongoBatch never reached the batchSize threshold
	if len(mongoBatch) > 0 {
		if _, err := di.MongoWords.BulkWrite(ctx, mongoBatch); err != nil {
			errors <- fmt.Errorf("error writing final batch to MongoDB: %v", err)
			return
		}

		if err := di.BleveIndex.Batch(bleveBatch); err != nil {
			errors <- fmt.Errorf("error writing final batch to Bleve: %v", err)
			return
		}
	}
}
