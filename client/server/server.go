package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/izquiratops/tango/common/database"
	"github.com/izquiratops/tango/common/types"
	"github.com/izquiratops/tango/common/utils"
)

type Server struct {
	db           *database.Database
	config       types.ServerConfig
	staticPrefix http.Handler
}

type SearchData struct {
	Query   string
	Results []database.Word
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	templatePath, _ := utils.GetAbsolutePath("template/index.html")
	http.ServeFile(w, r, templatePath)

	duration := time.Since(startTime)
	s.logRequest(r, http.StatusOK, duration)
}

func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	statusCode := http.StatusOK

	query := r.URL.Query().Get("query")
	results, err := s.search(query)

	var templatePath string
	if err != nil {
		if err.Error() == "EMPTY_LIST" {
			templatePath, _ = utils.GetAbsolutePath("template/not_found.html")
		} else {
			statusCode = http.StatusInternalServerError
			http.Error(w, fmt.Sprintf("Search error: %v", err), statusCode)

			duration := time.Since(startTime)
			s.logRequest(r, statusCode, duration)
			return
		}
	} else {
		templatePath, _ = utils.GetAbsolutePath("template/results.html")
	}

	// Parse template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		statusCode = http.StatusInternalServerError
		http.Error(w, fmt.Sprintf("Template parsing error: %v", err), statusCode)

		duration := time.Since(startTime)
		s.logRequest(r, statusCode, duration)
		return
	}

	// Render template
	data := SearchData{
		Query:   query,
		Results: results,
	}
	if err := tmpl.Execute(w, data); err != nil {
		statusCode = http.StatusInternalServerError
		http.Error(w, fmt.Sprintf("Template rendering error: %v", err), statusCode)

		duration := time.Since(startTime)
		s.logRequest(r, statusCode, duration)
		return
	}

	duration := time.Since(startTime)
	s.logRequest(r, http.StatusOK, duration)
}

func (s *Server) staticFileHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	statusCode := http.StatusOK

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		path := strings.TrimPrefix(r.URL.Path, "/static/")
		gzPath := filepath.Join("static", path+".gz")

		if _, err := os.Stat(gzPath); err == nil {
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Content-Type", getContentType(path))
			http.ServeFile(w, r, gzPath)

			duration := time.Since(startTime)
			s.logRequest(r, statusCode, duration)
			return
		}
	}

	// Fallback to the original file server if no compressed version exists
	// or if the client doesn't accept gzip
	s.staticPrefix.ServeHTTP(w, r)

	duration := time.Since(startTime)
	s.logRequest(r, statusCode, duration)
}

func (s *Server) SetupRoutes() *http.ServeMux {
	fmt.Printf("Setting up routes...\n")

	staticSystem := http.Dir("static")
	staticServer := http.FileServer(staticSystem)
	s.staticPrefix = http.StripPrefix("/static", staticServer)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.indexHandler)
	mux.HandleFunc("GET /search", s.searchHandler)
	mux.HandleFunc("GET /static/", s.staticFileHandler)

	return mux
}

func NewServer(config types.ServerConfig) (*Server, error) {
	db, err := database.NewDatabase(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &Server{
		db:     db,
		config: config,
	}, nil
}
