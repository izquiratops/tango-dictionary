package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/izquiratops/tango/client/server"
	"github.com/izquiratops/tango/common/utils"
)

var mongoDomainMap = map[bool]string{
	true:  "localhost",
	false: "mongo",
}

func main() {
	config, err := utils.LoadEnvironmentConfig(mongoDomainMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nInitializing server...\n")

	server, err := server.NewServer(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't initialize the server: %v\n", err)
		os.Exit(1)
	}

	mux := server.SetupRoutes()

	fmt.Printf("Server listening at 0.0.0.0:8080\n")
	if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
		fmt.Fprintf(os.Stderr, "Error Details: %v\n", err)
		os.Exit(1)
	}
}
