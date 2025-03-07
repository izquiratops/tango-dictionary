package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/izquiratops/tango/common/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	defaultSearchSize = 20
	defaultSearchFrom = 0
)

func (s *Server) search(query string) ([]database.Word, error) {
	ids, err := performBleveQuery(query, s.db)
	if err != nil {
		log.Printf("Failed to run Bleve query: %v", err)
		return nil, err
	}

	if len(ids) == 0 {
		emptyResultsErr := errors.New("EMPTY_LIST")
		return nil, emptyResultsErr
	}

	results, err := fetchWordsByIDs(ids, s.db)
	if err != nil {
		log.Printf("Failed to run MongoDB find: %v", err)
		return nil, err
	}

	return results, nil
}

// Code related to Bleve
func performBleveQuery(query string, db *database.Database) ([]string, error) {
	meaningsQuery := bleve.NewTermQuery(query)
	meaningsQuery.SetField("meanings")

	kanaBooleanQuery := newJapaneseFieldQuery(query, "kana_exact", "kana_char")
	kanjiBooleanQuery := newJapaneseFieldQuery(query, "kanji_exact", "kanji_char")

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

	searchResults, err := db.BleveIndex.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to search Bleve index: %w", err)
	}

	ids := extractBleveResult(searchResults)

	return ids, nil
}

func newTermQueryWithBoost(field string, term string, boost float64) *query.TermQuery {
	query := bleve.NewTermQuery(term)
	query.SetField(field)
	query.SetBoost(boost)
	return query
}

func newMatchQueryWithBoost(field string, term string, boost float64) *query.MatchQuery {
	query := bleve.NewMatchQuery(term)
	query.SetField(field)
	query.SetBoost(boost)
	return query
}

func newJapaneseFieldQuery(query string, exactField string, charField string) *query.BooleanQuery {
	booleanQuery := bleve.NewBooleanQuery()

	exactQuery := newTermQueryWithBoost(exactField, query, 2.0)
	booleanQuery.AddShould(exactQuery)

	charQuery := newMatchQueryWithBoost(charField, query, 0.5)
	disjunctionQuery := bleve.NewDisjunctionQuery(charQuery)
	disjunctionQuery.SetMin(1)
	booleanQuery.AddShould(disjunctionQuery)

	return booleanQuery
}

func extractBleveResult(searchResults *bleve.SearchResult) []string {
	var ids []string // List of Ids for every query hit

	for _, hit := range searchResults.Hits {
		var entry database.WordSearchable

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
func fetchWordsByIDs(ids []string, db *database.Database) ([]database.Word, error) {
	ctx := context.Background()

	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}

	cursor, err := db.MongoWords.Find(ctx, filter)
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

func extractCursorResult(cursor *mongo.Cursor, ctx context.Context) ([]database.Word, error) {
	var results []database.Word
	for cursor.Next(ctx) {
		var result database.Word
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

func sortWords(results []database.Word, targetOrder []string) []database.Word {
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
