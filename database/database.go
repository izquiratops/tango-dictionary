package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"sync"
	"tango/model"
	"tango/util"

	"github.com/blevesearch/bleve/v2"
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
	// Setup version names
	bleveFilename := fmt.Sprintf("jmdict_%v.bleve", dbVersion)
	blevePath := filepath.Join(indexFolder, bleveFilename)
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

func (di *Database) Search(query string) ([]model.JMdictWord, error) {
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

func (di *Database) runBleveQuery(query string) ([]string, error) {
	meaningsQuery := bleve.NewTermQuery(query)
	meaningsQuery.SetField("meanings")

	kanaBooleanQuery := util.NewJapaneseFieldQuery(query, "kana_exact", "kana_char")
	kanjiBooleanQuery := util.NewJapaneseFieldQuery(query, "kanji_exact", "kanji_char")

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
		var entry model.BleveEntry

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

func (di *Database) runMongoFind(ids []string) ([]model.JMdictWord, error) {
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
	var results []model.JMdictWord
	for cursor.Next(ctx) {
		var result model.JMdictWord
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

func (di *Database) importJmdictEntries(entries <-chan model.JMdictWord, errors chan<- error, wg *sync.WaitGroup) {
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
