package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/izquiratops/tango/common/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	defaultSearchSize = 20
	defaultSearchFrom = 0
)

func (s *Server) search(searchTerm string) ([]database.Word, error) {
	// Make sure to search using lowercase only
	searchTerm = strings.ToLower(searchTerm)

	ids, err := performBleveQuery(searchTerm, s.db)
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
func performBleveQuery(searchTerm string, db *database.Database) ([]string, error) {
	mainQuery := bleve.NewBooleanQuery()

	searchTermType := DetectSearchTermType(searchTerm)

	switch searchTermType {
	case "romaji":
		meaningsPhraseQuery := bleve.NewMatchQuery(searchTerm)
		meaningsPhraseQuery.SetField("meanings")
		meaningsPhraseQuery.SetBoost(4.0)

		meaningsFuzzyQuery := bleve.NewFuzzyQuery(searchTerm)
		meaningsFuzzyQuery.SetField("meanings")
		meaningsFuzzyQuery.SetFuzziness(1.0)
		meaningsFuzzyQuery.SetBoost(1.0)

		mainQuery.AddShould(
			meaningsPhraseQuery,
			meaningsFuzzyQuery,
		)
	case "kana":
		kanaExactQuery := bleve.NewMatchQuery(searchTerm)
		kanaExactQuery.SetField("kana_exact")
		kanaExactQuery.SetBoost(5.0)

		kanaPrefixQuery := bleve.NewPrefixQuery(searchTerm)
		kanaPrefixQuery.SetField("kana_char")
		kanaPrefixQuery.SetBoost(4.0)

		kanaCharQuery := bleve.NewMatchQuery(searchTerm)
		kanaCharQuery.SetField("kana_char")
		kanaCharQuery.SetBoost(1.0)

		mainQuery.AddShould(
			kanaExactQuery,
			kanaPrefixQuery,
			kanaCharQuery,
		)
	case "kanji":
		kanjiExactQuery := bleve.NewMatchPhraseQuery(searchTerm)
		kanjiExactQuery.SetField("kanji_exact")
		kanjiExactQuery.SetBoost(5.0)

		kanjiPrefixQuery := bleve.NewPrefixQuery(searchTerm)
		kanjiPrefixQuery.SetField("kanji_char")
		kanjiPrefixQuery.SetBoost(4.0)

		kanjiCharQuery := bleve.NewMatchQuery(searchTerm)
		kanjiCharQuery.SetField("kanji_char")
		kanjiCharQuery.SetBoost(1.0)

		mainQuery.AddShould(
			kanjiExactQuery,
			kanjiPrefixQuery,
			kanjiCharQuery,
		)
	}

	searchRequest := bleve.NewSearchRequest(mainQuery)
	searchRequest.Explain = true // Debug
	searchRequest.Size = defaultSearchSize
	searchRequest.From = defaultSearchFrom
	searchRequest.Fields = []string{
		"id",
		"kana_char",
		"kanji_char",
		"meanings",
	}

	searchResults, err := db.BleveIndex.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to search Bleve index: %w", err)
	}

	ids := extractBleveResult(searchResults)

	return ids, nil
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
	// TODO: use ctx from request
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
