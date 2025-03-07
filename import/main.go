package main

import (
	"fmt"
	"os"

	"github.com/izquiratops/tango/common/database"
	"github.com/izquiratops/tango/common/types"
)

func main() {
	config := types.ServerConfig{
		JmdictVersion: "3.6.1",
		MongoURI:      "mongodb://localhost:27017",
	}

	db, err := database.NewDatabase(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Details: %v\n", err)
		os.Exit(1)
	}

	if err := Import("../jmdict_source/jmdict-eng-3.6.1.json", db); err != nil {
		fmt.Fprintf(os.Stderr, "Error Details: %v\n", err)
		os.Exit(1)
	}
}
