package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/blevesearch/bleve/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDictionaryImporter(mongoURI string, batchSize int) (*DictionaryImporter, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	// Configurar Bleve
	mapping := bleve.NewIndexMapping()

	documentMapping := bleve.NewDocumentMapping()

	// TODO: Apply index for Japanese too! "プレ" doesn't work
	documentMapping.AddFieldMappingsAt("kanji", bleve.NewTextFieldMapping())
	documentMapping.AddFieldMappingsAt("kana", bleve.NewTextFieldMapping())
	documentMapping.AddFieldMappingsAt("meanings", bleve.NewTextFieldMapping())

	mapping.AddDocumentMapping("_default", documentMapping)

	// Crear o abrir índice Bleve
	bleveIndex, err := bleve.New("dictionary.bleve", mapping)

	if err != nil {
		bleveIndex, err = bleve.Open("dictionary.bleve")
		if err != nil {
			return nil, fmt.Errorf("error creating/opening Bleve index: %v", err)
		}
	}

	return &DictionaryImporter{
		mongoClient: client,
		collection:  client.Database("dictionary").Collection("entries"),
		bleveIndex:  bleveIndex,
		batchSize:   batchSize,
	}, nil
}

func (di *DictionaryImporter) ImportFromJSON(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var source JMdict
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&source); err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	entriesChan := make(chan JMdictWord, di.batchSize)
	errorsChan := make(chan error, 1)
	var wg sync.WaitGroup

	// Iniciar workers
	numWorkers := 3
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go di.ProcessEntries(entriesChan, errorsChan, &wg)
	}

	// Procesar entradas
	count := 0
	startTime := time.Now()

	for _, entry := range source.Words {
		select {
		case err := <-errorsChan:
			close(entriesChan)
			return fmt.Errorf("worker error: %v", err)
		default:
			entriesChan <- entry
			count++

			if count%1000 == 0 {
				elapsed := time.Since(startTime)
				rate := float64(count) / elapsed.Seconds()
				log.Printf("processed %d entries (%.2f entries/sec)", count, rate)
			}
		}
	}

	close(entriesChan)
	wg.Wait()

	log.Printf("import completed. Processed %d entries in %v", count, time.Since(startTime))
	return nil
}

func (di *DictionaryImporter) ProcessEntries(entries <-chan JMdictWord, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx := context.Background()
	mongoBatch := make([]mongo.WriteModel, 0, di.batchSize)
	bleveBatch := di.bleveIndex.NewBatch()

	for entry := range entries {
		// TODO: Avoid save text fields here
		// Prepare MongoDB
		model := mongo.NewInsertOneModel().SetDocument(entry)
		mongoBatch = append(mongoBatch, model)

		// Prepare Bleve
		searchableEntry, err := entry.ToSearchable()
		if err != nil {
			errors <- err
			return
		}

		if err := bleveBatch.Index(entry.ID, searchableEntry); err != nil {
			errors <- fmt.Errorf("error indexing in Bleve: %v", err)
			return
		}

		if len(mongoBatch) >= di.batchSize {
			// MongoDB bulk write
			if _, err := di.collection.BulkWrite(ctx, mongoBatch); err != nil {
				errors <- fmt.Errorf("error writing to MongoDB: %v", err)
				return
			}

			// Bleve batch write
			if err := di.bleveIndex.Batch(bleveBatch); err != nil {
				errors <- fmt.Errorf("error writing to Bleve: %v", err)
				return
			}

			mongoBatch = mongoBatch[:0]
			bleveBatch = di.bleveIndex.NewBatch()
		}
	}

	// Procesar último batch
	if len(mongoBatch) > 0 {
		if _, err := di.collection.BulkWrite(ctx, mongoBatch); err != nil {
			errors <- fmt.Errorf("error writing final batch to MongoDB: %v", err)
			return
		}

		if err := di.bleveIndex.Batch(bleveBatch); err != nil {
			errors <- fmt.Errorf("error writing final batch to Bleve: %v", err)
			return
		}
	}
}

func (di *DictionaryImporter) Search(query string) ([]SearchableEntry, error) {
	q := bleve.NewMatchQuery(query)

	searchRequest := bleve.NewSearchRequest(q)
	searchRequest.Fields = []string{"id", "kana", "kanji", "meanings"}
	searchRequest.Size = 20
	searchRequest.From = 0

	searchResults, err := di.bleveIndex.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	var results []SearchableEntry
	for _, hit := range searchResults.Hits {
		var entry SearchableEntry

		// Serialize the map to a JSON byte slice
		jsonBytes, err := json.Marshal(hit.Fields)
		if err != nil {
			fmt.Println("error marshalling fields:", err)
			continue
		}

		// Unmarshal the JSON byte slice into the SearchableEntry struct using a custom unmarshaler
		if err := json.Unmarshal(jsonBytes, &entry); err != nil {
			fmt.Println("error unmarshalling entry:", err)
			continue
		}

		results = append(results, entry)
	}

	return results, nil
}
