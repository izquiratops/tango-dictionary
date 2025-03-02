package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"tango/server"
)

func resolveVersion() string {
	var version string
	flag.StringVar(&version, "version", "", "Set the version value")
	flag.Parse()

	// Set the env var TANGO_VERSION as fallback value
	if version == "" {
		version = os.Getenv("TANGO_VERSION")
	}

	return version
}

func main() {
	var rebuildDatabase bool
	flag.BoolVar(&rebuildDatabase, "rebuild-database", false, "Rebuilds Database for the prompted version")

	dbVersion := resolveVersion()
	if dbVersion == "" {
		log.Fatalf("You must set a JMDict version")
	}

	fmt.Printf("Running server with version: %s\n", dbVersion)

	if err := server.RunServer(dbVersion, rebuildDatabase); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
