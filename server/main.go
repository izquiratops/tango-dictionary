package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"tango/database"
	"time"
)

var db *database.Database

type SearchData struct {
	Query   string
	Results []database.EntryDatabase
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./server/template/index.html")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	query := r.URL.Query().Get("query")

	var tmpl *template.Template
	var parseErr error

	results, searchErr := db.Search(query)
	if searchErr == nil {
		tmpl, parseErr = template.ParseFiles("./server/template/results.html")
	} else {
		switch {
		case strings.Contains(searchErr.Error(), "no results found"):
			tmpl, parseErr = template.ParseFiles("./server/template/not_found.html")
		default:
			http.Error(w, searchErr.Error(), http.StatusInternalServerError)
			return
		}
	}

	if parseErr != nil {
		fmt.Printf("Template parsing error: %v\n", parseErr)
		http.Error(w, parseErr.Error(), http.StatusInternalServerError)
		return
	}

	data := SearchData{
		Query:   query,
		Results: results,
	}

	executeErr := tmpl.Execute(w, data)
	if executeErr != nil {
		fmt.Printf("Template execution error: %v\n", executeErr)
		http.Error(w, executeErr.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(startTime)
	fmt.Printf("Served search '%s' in %s\n", query, duration)
}

func RunServer(databaseVersion string, rebuildDatabase bool) error {
	var err error
	db, err = database.NewDatabase("mongodb://mongo:27017", databaseVersion, 1000, rebuildDatabase)
	if err != nil {
		log.Fatalf("Couldn't setup database: %v", err)
	}
	fmt.Printf("Database setted up successfully\n")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", indexHandler)
	mux.HandleFunc("GET /search", searchHandler)

	fileSystem := http.Dir("./server/static")
	fileServer := http.FileServer(fileSystem)
	fileHandler := http.StripPrefix("/static", fileServer)
	mux.Handle("GET /static", fileHandler)

	addr := "0.0.0.0:8080"
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

	return nil
}
