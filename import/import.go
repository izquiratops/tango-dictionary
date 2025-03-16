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

	// Open JMdict json file
	file, err := os.Open(jsonPath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var jsonSource jmdict.JMdict
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&jsonSource); err != nil {
		return "", fmt.Errorf("error decoding JSON: %v", err)
	}

	entriesChan := make(chan jmdict.JMdictWord, batchSize)
	errorsChan := make(chan error, 1)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go bulkImportJmdictEntries(entriesChan, errorsChan, &wg, db)
	}

	startTime := time.Now()
	for _, entry := range jsonSource.Words {
		select {
		case err := <-errorsChan:
			close(entriesChan)
			return "", fmt.Errorf("worker error: %v", err)
		default:
			entriesChan <- entry
		}
	}

	close(entriesChan)
	wg.Wait()

	fmt.Printf("Dictionary import completed. Processed %d entries in %v\n", len(jsonSource.Words), time.Since(startTime))
	return jsonPath, nil
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
