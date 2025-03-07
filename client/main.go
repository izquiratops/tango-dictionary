package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/izquiratops/tango/client/server"
)

const (
	addr = "0.0.0.0:8080"
)

func main() {
	server, err := server.NewServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't initialize the server: %v\n", err)
		os.Exit(1)
	}

	mux := server.SetupRoutes()

	fmt.Printf("Server listening at %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "Error Details: %v\n", err)
		os.Exit(1)
	}
}
