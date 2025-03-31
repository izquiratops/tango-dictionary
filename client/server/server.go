package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/izquiratops/tango/common/database"
	"github.com/izquiratops/tango/common/types"
	"go.mongodb.org/mongo-driver/bson"
)

type Server struct {
	db           *database.Database
	config       types.ServerConfig
	staticPrefix http.Handler
	tagsCache    map[string]string
}

type SearchData struct {
	Query   string
	Results []database.Word
}

func NewServer(config types.ServerConfig) (*Server, error) {
	db, err := database.NewDatabase(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	tagsCache, err := setupTagsCache(db)
	if err != nil {
		return nil, fmt.Errorf("failed to setup tags cache: %w", err)
	}

	fmt.Printf("Tags cache initialized successfully\n")

	return &Server{
		db:        db,
		config:    config,
		tagsCache: tagsCache,
	}, nil
}

func (s *Server) SetupRoutes() *http.ServeMux {
	staticSystem := http.Dir("static")
	staticServer := http.FileServer(staticSystem)
	s.staticPrefix = http.StripPrefix("/static", staticServer)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.indexHandler)
	mux.HandleFunc("GET /search", s.searchHandler)
	mux.HandleFunc("GET /static/", s.staticFileHandler)

	return mux
}

func setupTagsCache(db *database.Database) (map[string]string, error) {
	ctx := context.Background()
	cursor, err := db.MongoTags.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error fetching tags: %v", err)
	}
	defer cursor.Close(ctx)

	tags := make(map[string]string)

	for cursor.Next(ctx) {
		var tag database.Tag
		if err := cursor.Decode(&tag); err != nil {
			return nil, fmt.Errorf("error decoding tag: %v", err)
		}

		tags[tag.Name] = tag.Description
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over tags: %v", err)
	}

	return tags, nil
}
