package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"tango/model"
	"tango/util"

	"github.com/blevesearch/bleve/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (di *Database) Search(query string) ([]model.JMdictWord, error) {
	ids, err := performBleveQuery(query, di)
	if err != nil {
		log.Printf("Failed to run Bleve query: %v", err)
		return nil, err
	}

	if len(ids) == 0 {
		// Define a specific error for empty results
		emptyResultsErr := errors.New("no results found")
		return nil, emptyResultsErr
	}

	results, err := fetchWordsByIDs(ids, di)
	if err != nil {
		log.Printf("Failed to run MongoDB find: %v", err)
		return nil, err
	}

	return results, nil
}

// Code related to Bleve
func performBleveQuery(query string, di *Database) ([]string, error) {
	meaningsQuery := bleve.NewTermQuery(query)
	meaningsQuery.SetField("meanings")

	kanaBooleanQuery := util.NewJapaneseFieldQuery(query, "kana_exact", "kana_char")
	kanjiBooleanQuery := util.NewJapaneseFieldQuery(query, "kanji_exact", "kanji_char")

	booleanQuery := bleve.NewBooleanQuery()
	booleanQuery.AddShould(meaningsQuery)
	booleanQuery.AddShould(kanaBooleanQuery)
	booleanQuery.AddShould(kanjiBooleanQuery)

	searchRequest := bleve.NewSearchRequest(booleanQuery)
	searchRequest.Fields = []string{
		"id",
		"kana_exact",
		"kana_char",
		"kanji_exact",
		"kanji_char",
		"meanings",
	}
	searchRequest.Size = defaultSearchSize
	searchRequest.From = defaultSearchFrom

	searchResults, err := di.bleveIndex.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to search Bleve index: %w", err)
	}

	ids := extractBleveResult(searchResults)

	return ids, nil
}

func extractBleveResult(searchResults *bleve.SearchResult) []string {
	var ids []string // List of Ids for every query hit

	for _, hit := range searchResults.Hits {
		var entry model.SearchableEntry

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

	return ids
}

// Code related to MongoDB
func fetchWordsByIDs(ids []string, di *Database) ([]model.JMdictWord, error) {
	ctx := context.Background()

	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}

	cursor, err := di.mongoDict.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents in MongoDB: %w", err)
	}
	defer cursor.Close(ctx)

	results, err := extractCursorResult(cursor, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over cursor: %w", err)
	}

	// Sorting what mongo returns to match the order of Bleve IDs
	sortedResults := sortWords(results, ids)

	return sortedResults, nil
}

func extractCursorResult(cursor *mongo.Cursor, ctx context.Context) ([]model.JMdictWord, error) {
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

	return results, nil
}

func sortWords(results []model.JMdictWord, targetOrder []string) []model.JMdictWord {
	sort.SliceStable(results, func(i, j int) bool {
		for _, id := range targetOrder {
			if results[i].ID == id {
				return true
			}
			if results[j].ID == id {
				return false
			}
		}
		return false // This should never be reached if all IDs are found
	})

	return results
}
