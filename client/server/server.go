package server

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/izquiratops/tango/common/database"
	"github.com/izquiratops/tango/common/types"
	"github.com/izquiratops/tango/common/utils"
)

type Server struct {
	db     *database.Database
	config types.ServerConfig
}

type SearchData struct {
	Query   string
	Results []database.Word
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	templatePath, _ := utils.GetAbsolutePath("template/index.html")
	http.ServeFile(w, r, templatePath)
}

func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	query := r.URL.Query().Get("query")
	results, err := s.search(query)

	// Choose template based on results
	var templatePath string
	if err != nil {
		if err.Error() == "EMPTY_LIST" {
			templatePath, _ = utils.GetAbsolutePath("template/not_found.html")
		} else {
			http.Error(w, fmt.Sprintf("Search error: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		templatePath, _ = utils.GetAbsolutePath("template/results.html")
	}

	// Parse template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template parsing error: %v", err), http.StatusInternalServerError)
		return
	}

	// Render template
	data := SearchData{
		Query:   query,
		Results: results,
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("Template rendering error: %v", err), http.StatusInternalServerError)
		return
	}

	duration := time.Since(startTime)
	fmt.Printf("Served search '%s' in %v\n", query, duration)
}

func (s *Server) SetupRoutes() *http.ServeMux {
	fmt.Printf("Setting up routes...\n")

	staticSystem := http.Dir("static")
	staticServer := http.FileServer(staticSystem)
	staticPrefix := http.StripPrefix("/static", staticServer)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.indexHandler)
	mux.HandleFunc("GET /search", s.searchHandler)
	mux.Handle("GET /static/", staticPrefix)

	return mux
}

func NewServer(config types.ServerConfig) (*Server, error) {
	db, err := database.NewDatabase(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &Server{
		db:     db,
		config: config,
	}, nil
}
