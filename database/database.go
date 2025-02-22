package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/blevesearch/bleve/v2"
	bleveSearch "github.com/blevesearch/bleve/v2/search/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultSearchSize = 20
	defaultSearchFrom = 0
)

type Database struct {
	mongoClient     *mongo.Client
	mongoCollection *mongo.Collection
	bleveIndex      bleve.Index
	batchSize       int
}

// 'indexFolder' makes it easier to run this method from Server and also from Unit test, where paths are different.
func NewDatabase(mongoURI string, indexFolder string, dbVersion string, batchSize int) (*Database, error) {
	// Set Bleve path
	bleveFilename := fmt.Sprintf("jmdict_%v.bleve", dbVersion)
	blevePath := filepath.Join(indexFolder, bleveFilename)

	// Set Mongo collection name
	mongoCollectionName := fmt.Sprintf("jmdict_%v", dbVersion)

	// Setup Mongo
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	// Setup Bleve
	indexMapping := bleve.NewIndexMapping()
	documentMapping := bleve.NewDocumentMapping()

	// Kana
	kanaCharMapping := bleve.NewTextFieldMapping()

	kanaExactMapping := bleve.NewTextFieldMapping()
	kanaExactMapping.Analyzer = "keyword" // This will treat the entire field as a single token

	// Kanji
	kanjiCharMapping := bleve.NewTextFieldMapping()

	kanjiExactMapping := bleve.NewTextFieldMapping()
	kanjiExactMapping.Analyzer = "keyword" // This will treat the entire field as a single token

	// Regular text analyzer for meanings
	meaningsMapping := bleve.NewTextFieldMapping()

	documentMapping.AddFieldMappingsAt("kana_exact", kanaExactMapping)
	documentMapping.AddFieldMappingsAt("kana_char", kanaCharMapping)
	documentMapping.AddFieldMappingsAt("kanji_exact", kanjiExactMapping)
	documentMapping.AddFieldMappingsAt("kanji_char", kanjiCharMapping)
	documentMapping.AddFieldMappingsAt("meanings", meaningsMapping)

	indexMapping.AddDocumentMapping("_default", documentMapping)

	// Try to open index, create one if doesn't exist
	bleveIndex, err := bleve.New(blevePath, indexMapping)
	if err != nil {
		bleveIndex, err = bleve.Open(blevePath)
		if err != nil {
			return nil, fmt.Errorf("error creating/opening Bleve index: %v", err)
		}
	}

	return &Database{
		mongoClient:     client,
		mongoCollection: client.Database("dictionary").Collection(mongoCollectionName),
		bleveIndex:      bleveIndex,
		batchSize:       batchSize,
	}, nil
}

func (di *Database) ImportFromJSON(path string) error {
	file, err := os.Open(path)
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

	numWorkers := 3
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go di.processEntries(entriesChan, errorsChan, &wg)
	}

	startTime := time.Now()
	for _, entry := range source.Words {
		select {
		case err := <-errorsChan:
			close(entriesChan)
			return fmt.Errorf("worker error: %v", err)
		default:
			entriesChan <- entry
		}
	}

	close(entriesChan)
	wg.Wait()

	log.Printf("Import completed. Processed %d entries in %v", len(source.Words), time.Since(startTime))
	return nil
}

func (di *Database) Search(query string) ([]JMdictWord, error) {
	ids, err := di.runBleveQuery(query)
	if err != nil {
		log.Printf("Failed to run Bleve query: %v", err)
		return nil, err
	}

	if len(ids) == 0 {
		// Define a specific error for empty results
		emptyResultsErr := errors.New("no results found")
		return nil, emptyResultsErr
	}

	results, err := di.runMongoFind(ids)
	if err != nil {
		log.Printf("Failed to run MongoDB find: %v", err)
		return nil, err
	}

	return results, nil
}

func NewTermQueryWithBoost(field string, term string, boost float64) *bleveSearch.TermQuery {
	query := bleve.NewTermQuery(term)
	query.SetField(field)
	query.SetBoost(boost)
	return query
}

func NewMatchQueryWithBoost(field string, term string, boost float64) *bleveSearch.MatchQuery {
	query := bleve.NewMatchQuery(term)
	query.SetField(field)
	query.SetBoost(boost)
	return query
}

func CreateBooleanQueryForField(query string, exactField string, charField string) *bleveSearch.BooleanQuery {
	booleanQuery := bleve.NewBooleanQuery()

	exactQuery := NewTermQueryWithBoost(exactField, query, 10.0)
	booleanQuery.AddShould(exactQuery)

	charQuery := NewMatchQueryWithBoost(charField, query, 1.0)
	disjunctionQuery := bleve.NewDisjunctionQuery(charQuery)
	disjunctionQuery.SetMin(1)
	booleanQuery.AddShould(disjunctionQuery)

	return booleanQuery
}

func (di *Database) runBleveQuery(query string) ([]string, error) {
	meaningsQuery := bleve.NewTermQuery(query)
	meaningsQuery.SetField("meanings")

	kanaBooleanQuery := CreateBooleanQueryForField(query, "kana_exact", "kana_char")
	kanjiBooleanQuery := CreateBooleanQueryForField(query, "kanji_exact", "kanji_char")

	booleanQuery := bleve.NewBooleanQuery()
	booleanQuery.AddShould(meaningsQuery)
	booleanQuery.AddShould(kanaBooleanQuery)
	booleanQuery.AddShould(kanjiBooleanQuery)

	searchRequest := bleve.NewSearchRequest(booleanQuery)
	searchRequest.Fields = []string{"id", "kana_exact", "kana_char", "kanji_exact", "kanji_char", "meanings"}
	searchRequest.Size = defaultSearchSize
	searchRequest.From = defaultSearchFrom

	searchResults, err := di.bleveIndex.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to search Bleve index: %w", err)
	}

	var ids []string // List of Ids for every query hit
	for _, hit := range searchResults.Hits {
		var entry BleveEntry

		// Serialize the map to a JSON byte slice
		jsonBytes, err := json.Marshal(hit.Fields)
		if err != nil {
			fmt.Println("error marshalling fields:", err)
			continue
		}

		// Unmarshal the JSON byte slice into the BleveEntry struct using a custom unmarshaler
		if err := json.Unmarshal(jsonBytes, &entry); err != nil {
			fmt.Println("error unmarshalling entry:", err)
			continue
		}

		ids = append(ids, entry.ID)
	}

	return ids, nil
}

func (di *Database) runMongoFind(ids []string) ([]JMdictWord, error) {
	ctx := context.Background()

	f := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}

	cursor, err := di.mongoCollection.Find(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents in MongoDB: %w", err)
	}
	defer cursor.Close(ctx)

	// Iter through the Mongo cursor to fetch the returned Find data
	var results []JMdictWord
	for cursor.Next(ctx) {
		var result JMdictWord
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode document: %w", err)
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	// Sort the results based on the original order of IDs
	sort.SliceStable(results, func(i, j int) bool {
		for _, id := range ids {
			if results[i].ID == id {
				return true
			}
			if results[j].ID == id {
				return false
			}
		}
		return false // This should never be reached if all IDs are found
	})

	return results, nil
}

func (di *Database) processEntries(entries <-chan JMdictWord, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx := context.Background()
	mongoBatch := make([]mongo.WriteModel, 0, di.batchSize)
	bleveBatch := di.bleveIndex.NewBatch()

	for entry := range entries {
		// Prepare MongoDB
		bsonData, err := bson.Marshal(entry)
		if err != nil {
			errors <- fmt.Errorf("error marshalling to BSON: %v", err)
			return
		}

		model := mongo.NewInsertOneModel().SetDocument(bsonData)
		mongoBatch = append(mongoBatch, model)

		// Prepare Bleve
		bleveEntry, err := entry.ToBleveEntry()
		if err != nil {
			errors <- err
			return
		}

		if err := bleveBatch.Index(entry.ID, bleveEntry); err != nil {
			errors <- fmt.Errorf("error indexing in Bleve: %v", err)
			return
		}

		if len(mongoBatch) >= di.batchSize {
			// MongoDB bulk write
			if _, err := di.mongoCollection.BulkWrite(ctx, mongoBatch); err != nil {
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

	// Process last batch. This runs the batch if mongoBatch never reached the batchSize threshold
	if len(mongoBatch) > 0 {
		if _, err := di.mongoCollection.BulkWrite(ctx, mongoBatch); err != nil {
			errors <- fmt.Errorf("error writing final batch to MongoDB: %v", err)
			return
		}

		if err := di.bleveIndex.Batch(bleveBatch); err != nil {
			errors <- fmt.Errorf("error writing final batch to Bleve: %v", err)
			return
		}
	}
}
