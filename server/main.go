package server

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserData struct {
	Name string
	Time string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := UserData{
		Name: "Guest",
		Time: time.Now().String(),
	}

	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func greetHandler(w http.ResponseWriter, r *http.Request) {
	data := UserData{
		Name: r.URL.Query().Get("name"),
		Time: time.Now().String(),
	}

	tmpl, err := template.ParseFiles("./templates/greet.html")
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
	ctx := context.Background()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// mongoClient.Database("dictionary").Collection("entries").InsertOne(ctx, bson.D{{Key: "name", Value: "pi"}, {Key: "value", Value: 3.14159}})

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", indexHandler)
	mux.HandleFunc("GET /greet", greetHandler)

	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", mux)
}
