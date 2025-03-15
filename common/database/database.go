package database

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/izquiratops/tango/common/types"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/lang/cjk"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/analysis/token/ngram"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	MongoWords *mongo.Collection
	MongoTags  *mongo.Collection
	BleveIndex bleve.Index
}

func NewDatabase(config types.ServerConfig) (*Database, error) {
	// Collections doesn't allow '.'s on their names
	mongoCollectionName := strings.Replace(config.JmdictVersion, ".", "_", -1)

	mongoDB, err := setupMongoDB(config.MongoURI, mongoCollectionName)
	if err != nil {
		return nil, err
	}

	fmt.Printf("MongoDB initialized successfully\n")

	bleveIndex, err := setupBleve(config.JmdictVersion)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Bleve initialized successfully\n")

	return &Database{
		MongoWords: mongoDB.Collection("words"),
		MongoTags:  mongoDB.Collection("tags"),
		BleveIndex: bleveIndex,
	}, nil
}

func setupMongoDB(mongoURI string, collectionName string) (*mongo.Database, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	return client.Database(collectionName), nil
}

func setupBleve(dbVersion string) (bleve.Index, error) {
	indexMapping := bleve.NewIndexMapping()

	if err := indexMapping.AddCustomAnalyzer("english_no_stop", map[string]interface{}{
		"type":      custom.Name,
		"tokenizer": unicode.Name,
		"token_filters": []string{
			lowercase.Name,
			en.PossessiveName,
			en.StopName,
		},
	}); err != nil {
		return nil, err
	}

	if err := indexMapping.AddCustomAnalyzer("japanese_ngram", map[string]interface{}{
		"type":      custom.Name,
		"tokenizer": unicode.Name,
		"token_filters": []string{
			lowercase.Name,
			cjk.WidthName,
			ngram.Name,
		},
		// https://github.com/blevesearch/bleve/blob/master/analysis/token/ngram/ngram.go
		"token_maps": map[string]interface{}{
			ngram.Name: map[string]interface{}{
				"min": 2,
				"max": 3,
			},
		},
	}); err != nil {
		return nil, err
	}

	documentMapping := bleve.NewDocumentMapping()

	// English indexes
	meaningsMapping := bleve.NewTextFieldMapping()
	meaningsMapping.Analyzer = "english_no_stop"
	documentMapping.AddFieldMappingsAt("meanings", meaningsMapping)

	// Kana indexes
	kanaExactMapping := bleve.NewTextFieldMapping()
	kanaExactMapping.Analyzer = keyword.Name
	documentMapping.AddFieldMappingsAt("kana_exact", kanaExactMapping)

	kanaCharMapping := bleve.NewTextFieldMapping()
	kanaExactMapping.Analyzer = "japanese_ngram"
	documentMapping.AddFieldMappingsAt("kana_char", kanaCharMapping)

	// Kanji indexes
	kanjiExactMapping := bleve.NewTextFieldMapping()
	kanjiExactMapping.Analyzer = keyword.Name
	documentMapping.AddFieldMappingsAt("kanji_exact", kanjiExactMapping)

	kanjiCharMapping := bleve.NewTextFieldMapping()
	kanaExactMapping.Analyzer = "japanese_ngram"
	documentMapping.AddFieldMappingsAt("kanji_char", kanjiCharMapping)

	// TODO: Romaji!

	// Default mapping
	indexMapping.DefaultAnalyzer = "english_no_stop"
	indexMapping.AddDocumentMapping("_default", documentMapping)

	bleveFilename := fmt.Sprintf("jmdict_%v.bleve", dbVersion)
	blevePath := filepath.Join("..", "jmdict_source", bleveFilename)
	bleveIndex, err := bleve.New(blevePath, indexMapping)
	if err != nil {
		bleveIndex, err = bleve.Open(blevePath)
		if err != nil {
			return nil, fmt.Errorf("error creating/opening Bleve index: %v", err)
		}
	}

	return bleveIndex, nil
}
