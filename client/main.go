package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/izquiratops/tango/client/server"
	"github.com/izquiratops/tango/common/config"
)

var mongoDomainMap = map[config.EnvironmentType]string{
	config.LocalEnv:  "localhost",
	config.ServerEnv: "mongo", // Client runs in a container, so it can connect to the 'mongo' container
}

func main() {
	config, err := config.LoadEnvironment(mongoDomainMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Initializing server...\n")
	server, err := server.NewServer(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't initialize the server: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Setting up routes...\n")
	mux := server.SetupRoutes()

	fmt.Printf("Server listening at 0.0.0.0:8080\n")
	if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
		fmt.Fprintf(os.Stderr, "Error Details: %v\n", err)
		os.Exit(1)
	}
}
