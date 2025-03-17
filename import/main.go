package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/izquiratops/tango/common/config"
	"github.com/izquiratops/tango/common/database"
	"github.com/izquiratops/tango/common/types"
)

var mongoDomainMap = map[config.EnvironmentType]string{
	config.LocalEnv:  "localhost",
	config.ServerEnv: "izquiratops.dev", // Import runs externally, so it connects to the server
}

func main() {
	config, err := config.LoadEnvironment(mongoDomainMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Remove Bleve index before database connection
	if err := removeBleveIndex(&config); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing Bleve index: %v\n", err)
		os.Exit(1)
	}

	db, err := database.NewDatabase(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Details: %v\n", err)
		os.Exit(1)
	}

	// Drop MongoDB collections after connection
	if err := dropMongoCollections(db); err != nil {
		fmt.Fprintf(os.Stderr, "Error dropping collections: %v\n", err)
		os.Exit(1)
	}

	jsonPath, err := Import(db, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Details: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n==============================================\n")
	fmt.Printf("✅ IMPORT COMPLETED SUCCESSFULLY!\n")
	fmt.Printf("==============================================\n\n")
	fmt.Printf("The Bleve index was created in: %s\n", jsonPath)

	if !config.MongoRunsLocal {
		fmt.Printf("\nNext step: Upload the index to your server using SCP:\n")
		fmt.Printf("scp -r ./jmdict_source/jmdict_X.Y.Z.bleve user@example.com:jmdict_source\n\n")
	}
}

func removeBleveIndex(config *types.ServerConfig) error {
	bleveFilename := fmt.Sprintf("jmdict_%v.bleve", config.JmdictVersion)
	blevePath := filepath.Join("..", "jmdict_source", bleveFilename)

	return os.RemoveAll(blevePath)
}

func dropMongoCollections(db *database.Database) error {
	ctx := context.Background()

	if err := db.MongoWords.Drop(ctx); err != nil {
		return fmt.Errorf("error dropping Words collection: %w", err)
	}
	if err := db.MongoTags.Drop(ctx); err != nil {
		return fmt.Errorf("error dropping Tags collection: %w", err)
	}

	return nil
}
