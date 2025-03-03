package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"tango/database"
	"tango/utils"
	"time"
)

const (
	addr = "0.0.0.0:8080"
	importBatchSize = 1000
)

type Server struct {
	db     *database.Database
	config ServerConfig
}

type ServerConfig struct {
	isLocalEnvironment bool
	shouldRebuild      bool
	jmdictVersion      string
	mongoURI           string
}

type SearchData struct {
	Query   string
	Results []database.EntryDatabase
}

func NewServer() (*Server, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	db, err := initializeDatabase(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &Server{
		db:     db,
		config: config,
	}, nil
}

func loadConfig() (ServerConfig, error) {
	jmdictVersion := os.Getenv("TANGO_VERSION")
	if jmdictVersion == "" {
		return ServerConfig{}, fmt.Errorf("TANGO_VERSION environment variable must be set")
	}

	isLocalEnvironment := utils.ResolveBooleanFromEnv("TANGO_LOCAL")
	mongoURI := map[bool]string{
		true:  "mongodb://localhost:27017",
		false: "mongodb://mongo:27017", // The docker service is currently called "mongo"
	}[isLocalEnvironment]

	return ServerConfig{
		isLocalEnvironment: isLocalEnvironment,
		shouldRebuild:      utils.ResolveBooleanFromEnv("TANGO_REBUILD"),
		jmdictVersion:      jmdictVersion,
		mongoURI:           mongoURI,
	}, nil
}

func initializeDatabase(config ServerConfig) (*database.Database, error) {
	db, err := database.NewDatabase(
		config.mongoURI,
		config.jmdictVersion,
		importBatchSize,
		config.shouldRebuild,
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./server/template/index.html")
}

func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	query := r.URL.Query().Get("query")

	tmpl, results, err := s.handleSearch(query)
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

func (s *Server) handleSearch(query string) (*template.Template, []database.EntryDatabase, error) {
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

func (s *Server) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.indexHandler)
	mux.HandleFunc("GET /search", s.searchHandler)

	fileSystem := http.Dir("./server/static")
	fileServer := http.FileServer(fileSystem)
	fileHandler := http.StripPrefix("/static", fileServer)
	mux.Handle("GET /static/", fileHandler)

	return mux
}

func RunServer() error {
	server, err := NewServer()
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	fmt.Printf("\nInitializing server...\n")
	fmt.Printf("JMDict Version: %s\n", server.config.jmdictVersion)
	fmt.Printf("Database Rebuild: %v\n", server.config.shouldRebuild)
	fmt.Printf("Database setted up successfully\n")

	mux := server.setupRoutes()

	if err := http.ListenAndServe(addr, mux); err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}

	return nil
}
