package database

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultSearchSize = 20
	defaultSearchFrom = 0
)

func NewDatabase(mongoURI string, dbVersion string, batchSize int, clearCollection bool) (*Database, error) {
	// Setup version names
	bleveFilename := fmt.Sprintf("jmdict_%v.bleve", dbVersion)
	// Mongo do not allow collection names with dots
	mongoCollectionName := strings.Replace(dbVersion, ".", "_", 2)

	// Setup Mongo
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	if clearCollection {
		// Drop tables before loading json data
		client.Database(mongoCollectionName).Collection("words").Drop(ctx)
		client.Database(mongoCollectionName).Collection("tags").Drop(ctx)

		// Clear Bleve index before creating a new one
		os.RemoveAll(fmt.Sprintf("./jmdict_%s.bleve", dbVersion))
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
	bleveIndex, err := bleve.New(bleveFilename, indexMapping)
	if err != nil {
		bleveIndex, err = bleve.Open(bleveFilename)
		if err != nil {
			return nil, fmt.Errorf("error creating/opening Bleve index: %v", err)
		}
	}

	return &Database{
		mongoWords: client.Database(mongoCollectionName).Collection("words"),
		mongoTags:  client.Database(mongoCollectionName).Collection("tags"),
		bleveIndex: bleveIndex,
		batchSize:  batchSize,
	}, nil
}
