package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
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

	results, err := db.Search(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := SearchData{
		Query:   query,
		Results: results,
	}

	// TODO: path 😔
	tmpl, err := template.ParseFiles("./server/template/results.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RunServer() {
	var err error
	db, err = database.NewDatabase("mongodb://localhost:27017", 1000)
	if err != nil {
		log.Fatalf("Couldn't connect to mongo database: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", indexHandler)
	mux.HandleFunc("GET /search", searchHandler)

	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", mux)
}
