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
	fmt.Printf("📤 Served index page to %s\n", r.RemoteAddr)
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
		fmt.Printf("❌ Template parsing error: %v\n", parseErr)
		http.Error(w, parseErr.Error(), http.StatusInternalServerError)
		return
	}

	data := SearchData{
		Query:   query,
		Results: results,
	}

	executeErr := tmpl.Execute(w, data)
	if executeErr != nil {
		fmt.Printf("❌ Template execution error: %v\n", executeErr)
		http.Error(w, executeErr.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(startTime)
	fmt.Printf("📤 Served search '%s' in %s to %s\n", query, duration, r.RemoteAddr)
}

func RunServer(databaseVersion string, rebuildDatabase bool) error {
	fmt.Printf("🔄 Connecting to MongoDB...\n")

	var err error
	db, err = database.NewDatabase("mongodb://localhost:27017", databaseVersion, 1000, rebuildDatabase)
	if err != nil {
		log.Fatalf("⛔ Couldn't setup database: %v", err)
	}
	fmt.Printf("✅ Database setted up successfully\n")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", indexHandler)
	fmt.Printf("📌 Registered route: GET /\n")

	mux.HandleFunc("GET /search", searchHandler)
	fmt.Printf("📌 Registered route: GET /search\n")

	fileSystem := http.Dir("./server/static")
	fileServer := http.FileServer(fileSystem)
	fileHandler := http.StripPrefix("/static", fileServer)
	mux.Handle("GET /static", fileHandler)
	fmt.Printf("📌 Registered route: GET /static\n")

	fmt.Printf("\n🚀 Starting server on localhost:8080\n")
	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		log.Fatalf("⛔ Server failed to start: %v", err)
	}

	return nil
}
