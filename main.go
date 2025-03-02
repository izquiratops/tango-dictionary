package main

import (
	"fmt"
	"log"
	"os"
	"tango/server"
	"tango/utils"
)

func main() {
	rebuildDatabase := utils.ResolveBooleanFromEnv("TANGO_REBUILD")
	databaseVersion := os.Getenv("TANGO_VERSION")

	if databaseVersion == "" {
		log.Fatalf("You must set a JMDict version")
	}

	fmt.Printf("\nInitializing server...\n")
	fmt.Printf("JMDict Version: %s\n", databaseVersion)
	fmt.Printf("Database Rebuild: %v\n", rebuildDatabase)

	fmt.Printf("\nStarting server...\n")
	if err := server.RunServer(databaseVersion, rebuildDatabase); err != nil {
		fmt.Fprintf(os.Stderr, "Error Details: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nServer shutting down...\n")
}
