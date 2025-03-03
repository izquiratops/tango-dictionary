package main

import (
	"fmt"
	"os"
	"tango/server"
)

func main() {
	fmt.Printf("\nStarting server...\n")
	if err := server.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "Error Details: %v\n", err)
		os.Exit(1)
	}
}
