package main

import (
	"fmt"
	"log"
	"os"
	"tango/server"
	"tango/utils"
	"time"
)

func main() {
	fmt.Printf("🕒 Startup Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	rebuildDatabase := utils.ResolveBooleanFromEnv("TANGO_REBUILD")
	databaseVersion := os.Getenv("TANGO_VERSION")
	if databaseVersion == "" {
		log.Fatalf("❌ You must set a JMDict version")
	}

	fmt.Printf("\n🚀 Initializing server...\n")
	fmt.Printf("📚 JMDict Version: %s\n", databaseVersion)
	fmt.Printf("🔄 Database Rebuild: %v\n", rebuildDatabase)

	fmt.Printf("\n⚡ Starting server...\n")
	if err := server.RunServer(databaseVersion, rebuildDatabase); err != nil {
		fmt.Fprintf(os.Stderr, "⛔ Error Details: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n👋 Server shutting down...\n")
}
