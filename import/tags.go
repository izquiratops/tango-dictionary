package main

import (
	"context"
	"fmt"

	"github.com/izquiratops/tango/common/database"
	"github.com/izquiratops/tango/common/jmdict"
	"go.mongodb.org/mongo-driver/mongo"
)

type TagEntry struct {
	Name        string `bson:"_id"`         // El nombre del tag será la clave primaria
	Description string `bson:"description"` // Descripción del tag
}

func ImportTags(db *database.Database, jsonSource *jmdict.JMdict) error {
	fmt.Println("Importing tags to MongoDB...")

	if len(jsonSource.Tags) == 0 {
		fmt.Println("No tags found in the dictionary.")
		return nil
	}

	var tagEntries []mongo.WriteModel
	for tagName, description := range jsonSource.Tags {
		tagEntry := TagEntry{
			Name:        tagName,
			Description: description,
		}

		model := mongo.NewInsertOneModel().SetDocument(tagEntry)
		tagEntries = append(tagEntries, model)
	}

	ctx := context.Background()
	result, err := db.MongoTags.BulkWrite(ctx, tagEntries)
	if err != nil {
		return fmt.Errorf("error importing tags to MongoDB: %v", err)
	}

	fmt.Printf("Imported %d tags to MongoDB\n", result.InsertedCount)
	return nil
}
