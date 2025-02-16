package server

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
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

	tmpl, err := template.ParseFiles("./server/template/index.html")
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

	tmpl, err := template.ParseFiles("./server/template/greet.html")
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
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", indexHandler)
	mux.HandleFunc("GET /greet", greetHandler)

	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", mux)
}
