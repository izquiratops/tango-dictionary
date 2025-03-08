package main

import (
	"fmt"
	"os"

	"github.com/izquiratops/tango/common/database"
	"github.com/izquiratops/tango/common/utils"
)

var mongoDomainMap = map[bool]string{
	true:  "localhost",
	false: "izquiratops.dev",
}

func main() {
	config, err := utils.LoadEnvironmentConfig(mongoDomainMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	db, err := database.NewDatabase(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Details: %v\n", err)
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
	fmt.Printf("\nNext step: Upload the index to your server using SCP:\n")
	fmt.Printf("scp -r ./jmdict_source/jmdict_X.Y.Z.bleve user@example.com:/root/jmdict_source\n\n")
}
