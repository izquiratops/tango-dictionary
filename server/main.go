package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"tango/database"
)

var db *database.Database

type SearchData struct {
	Query   string
	Results []database.JMdictWord
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./server/template/index.html")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	// Parse HTML based on the response of the database search
	var tmpl *template.Template
	var parseErr error

	results, searchErr := db.Search(query)
	if searchErr == nil {
		// 200 OK, list with valid content
		tmpl, parseErr = template.ParseFiles("./server/template/results.html")
	} else {
		switch {
		case strings.Contains(searchErr.Error(), "no results found"):
			// 200 OK, BUT empty list of results
			tmpl, parseErr = template.ParseFiles("./server/template/not_found.html")
		default:
			// 500 NOK, searchErr is any Errorf throw by Bleve or Mongo
			http.Error(w, searchErr.Error(), http.StatusInternalServerError)
			return
		}
	}

	if parseErr != nil {
		http.Error(w, parseErr.Error(), http.StatusInternalServerError)
		return
	}

	data := SearchData{
		Query:   query,
		Results: results,
	}

	executeErr := tmpl.Execute(w, data)
	if executeErr != nil {
		http.Error(w, executeErr.Error(), http.StatusInternalServerError)
	}
}

func RunServer(dbVersion string) error {
	var err error
	db, err = database.NewDatabase(
		"mongodb://localhost:27017",
		"./database",
		dbVersion,
		1000,
	)
	if err != nil {
		log.Fatalf("Couldn't connect to mongo database: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", indexHandler)
	mux.HandleFunc("GET /search", searchHandler)

	// Serve static files
	fsys := http.Dir("./server/static")
	fileServer := http.FileServer(fsys)
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	fmt.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

	return nil
}
