package database

import (
	"context"
	"fmt"
	"os"
	"strings"
	"tango/types"

	"github.com/blevesearch/bleve/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDatabase(config types.ServerConfig) (*Database, error) {
	mongoCollectionName := strings.Replace(config.JmdictVersion, ".", "_", -1)

	mongoDB, err := setupMongoDB(config.MongoURI, mongoCollectionName, config.ShouldRebuild)
	if err != nil {
		return nil, fmt.Errorf("error setting up MongoDB: %v", err)
	}

	fmt.Printf("MongoDB initialized successfully")

	bleveIndex, err := setupBleve(config.JmdictVersion, config.ShouldRebuild)
	if err != nil {
		return nil, fmt.Errorf("error setting up Bleve: %v", err)
	}

	fmt.Printf("Bleve initialized successfully")

	// if config.ShouldRebuild {
	// 	TODO: Run import here!
	// }

	return &Database{
		mongoWords: mongoDB.Collection("words"),
		mongoTags:  mongoDB.Collection("tags"),
		bleveIndex: bleveIndex,
		batchSize:  1000,
	}, nil
}

func setupMongoDB(mongoURI string, collectionName string, rebuildDatabase bool) (*mongo.Database, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	if rebuildDatabase {
		client.Database(collectionName).Collection("words").Drop(ctx)
		client.Database(collectionName).Collection("tags").Drop(ctx)
	}

	return client.Database(collectionName), nil
}

func setupBleve(dbVersion string, rebuildDatabase bool) (bleve.Index, error) {
	indexMapping := bleve.NewIndexMapping()
	documentMapping := bleve.NewDocumentMapping()

	// "keyword" will treat the entire field as a single token
	kanaExactMapping := bleve.NewTextFieldMapping()
	kanaExactMapping.Analyzer = "keyword"
	kanjiExactMapping := bleve.NewTextFieldMapping()
	kanjiExactMapping.Analyzer = "keyword"
	// Default mappings
	kanaCharMapping := bleve.NewTextFieldMapping()
	kanjiCharMapping := bleve.NewTextFieldMapping()
	meaningsMapping := bleve.NewTextFieldMapping()

	documentMapping.AddFieldMappingsAt("kana_exact", kanaExactMapping)
	documentMapping.AddFieldMappingsAt("kana_char", kanaCharMapping)
	documentMapping.AddFieldMappingsAt("kanji_exact", kanjiExactMapping)
	documentMapping.AddFieldMappingsAt("kanji_char", kanjiCharMapping)
	documentMapping.AddFieldMappingsAt("meanings", meaningsMapping)

	indexMapping.AddDocumentMapping("_default", documentMapping)

	bleveFilename := fmt.Sprintf("jmdict_%v.bleve", dbVersion)
	if rebuildDatabase {
		os.RemoveAll(bleveFilename)
	}

	bleveIndex, err := bleve.New(bleveFilename, indexMapping)
	if err != nil {
		bleveIndex, err = bleve.Open(bleveFilename)
		if err != nil {
			return nil, fmt.Errorf("error creating/opening Bleve index: %v", err)
		}
	}

	return bleveIndex, nil
}
