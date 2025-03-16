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
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	MongoWords *mongo.Collection
	MongoTags  *mongo.Collection
	BleveIndex bleve.Index
}

func NewDatabase(config *types.ServerConfig) (*Database, error) {
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

	if err := indexMapping.AddCustomAnalyzer("custom_english", map[string]interface{}{
		"type":      custom.Name,
		"tokenizer": unicode.Name,
		"token_filters": []string{
			lowercase.Name,
		},
	}); err != nil {
		return nil, err
	}

	documentMapping := bleve.NewDocumentMapping()

	// English indexes
	meaningsMapping := bleve.NewTextFieldMapping()
	meaningsMapping.Analyzer = "custom_english"
	documentMapping.AddFieldMappingsAt("meanings", meaningsMapping)

	// Kana indexes
	kanaExactMapping := bleve.NewTextFieldMapping()
	kanaExactMapping.Analyzer = keyword.Name
	documentMapping.AddFieldMappingsAt("kana_exact", kanaExactMapping)

	kanaCharMapping := bleve.NewTextFieldMapping()
	kanaCharMapping.Analyzer = cjk.AnalyzerName
	documentMapping.AddFieldMappingsAt("kana_char", kanaCharMapping)

	// Kanji indexes
	kanjiExactMapping := bleve.NewTextFieldMapping()
	kanjiExactMapping.Analyzer = keyword.Name
	documentMapping.AddFieldMappingsAt("kanji_exact", kanjiExactMapping)

	kanjiCharMapping := bleve.NewTextFieldMapping()
	kanjiCharMapping.Analyzer = cjk.AnalyzerName
	documentMapping.AddFieldMappingsAt("kanji_char", kanjiCharMapping)

	// Default mapping
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
