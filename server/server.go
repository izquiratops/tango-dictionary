package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"tango/database"
	"tango/types"
	"tango/utils"
	"time"
)

type Server struct {
	db     *database.Database
	config types.ServerConfig
}

type SearchData struct {
	Query   string
	Results []database.EntryDatabase
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./server/template/index.html")
}

func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	query := r.URL.Query().Get("query")

	tmpl, results, err := s.handleSearchRequest(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.renderTemplate(w, tmpl, SearchData{Query: query, Results: results}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(startTime)
	fmt.Printf("Served search '%s' in %s\n", query, duration)
}

func (s *Server) handleSearchRequest(query string) (*template.Template, []database.EntryDatabase, error) {
	results, err := s.db.Search(query)
	if err != nil {
		if strings.Contains(err.Error(), "no results found") {
			tmpl, err := template.ParseFiles("./server/template/not_found.html")
			return tmpl, nil, err
		}
		return nil, nil, err
	}

	tmpl, err := template.ParseFiles("./server/template/results.html")
	return tmpl, results, err
}

func (s *Server) renderTemplate(w http.ResponseWriter, tmpl *template.Template, data interface{}) error {
	if err := tmpl.Execute(w, data); err != nil {
		fmt.Printf("Template execution error: %v\n", err)
		return err
	}
	return nil
}

func (s *Server) SetupRoutes() *http.ServeMux {
	fmt.Printf("Setting up routes...\n")
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.indexHandler)
	mux.HandleFunc("GET /search", s.searchHandler)

	fileSystem := http.Dir("./server/static")
	fileServer := http.FileServer(fileSystem)
	fileHandler := http.StripPrefix("/static", fileServer)
	mux.Handle("GET /static/", fileHandler)

	return mux
}

func NewServer() (*Server, error) {
	config, err := loadEnvironmentConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Printf("\nInitializing server...\n")
	fmt.Printf("----------- Config values ----------\n")
	fmt.Printf("JMDict Version: %s\n", config.JmdictVersion)
	fmt.Printf("Rebuild database: %v\n", config.ShouldRebuild)
	fmt.Printf("Mongo connection Uri: %v\n", config.MongoURI)
	fmt.Printf("Local env: %v\n", config.IsLocalEnvironment)
	fmt.Printf("------------------------------------\n")

	db, err := database.NewDatabase(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	if config.ShouldRebuild {
		// Get executable directory
		execPath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		execDir := filepath.Dir(execPath)

		jsonPath := filepath.Join(execDir, "jmdict_source", fmt.Sprintf("jmdict-eng-%s.json", config.JmdictVersion))
		if err := db.ImportFromJSON(jsonPath); err != nil {
			return nil, fmt.Errorf("failed to import from json: %w", err)
		}
	}

	return &Server{
		db:     db,
		config: config,
	}, nil
}

func loadEnvironmentConfig() (types.ServerConfig, error) {
	jmdictVersion := os.Getenv("TANGO_VERSION")
	if jmdictVersion == "" {
		return types.ServerConfig{}, fmt.Errorf("TANGO_VERSION environment variable must be set")
	}

	isLocalEnvironment := utils.ResolveBooleanFromEnv("TANGO_LOCAL")
	mongoURI := map[bool]string{
		true:  "mongodb://localhost:27017",
		false: "mongodb://mongo:27017", // The docker service is currently called "mongo"
	}[isLocalEnvironment]

	return types.ServerConfig{
		IsLocalEnvironment: isLocalEnvironment,
		ShouldRebuild:      utils.ResolveBooleanFromEnv("TANGO_REBUILD"),
		JmdictVersion:      jmdictVersion,
		MongoURI:           mongoURI,
	}, nil
}
