package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Header().Set("Content-Type", "application/json")
	_, err := fmt.Fprint(w, `{"message": "Hello from Snippetbox!"}`)

	if err != nil {
		http.NotFound(w, r)
		return
	}
}

func getSnippet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	_, err = fmt.Fprintf(w, `{"message": "Viewing snippet with id %v"}`, id)

	if err != nil {
		panic("could not write response")
	}
}

func getSnippetForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := fmt.Fprint(w, `{"message": "Display form for creating a snippet."}`)

	if err != nil {
		panic("could not write response")
	}
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, err := fmt.Fprint(w, `{"message": "Creating a snippet."}`)

	if err != nil {
		panic("could not write response")
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippets/{id}", getSnippet)
	mux.HandleFunc("GET /snippets/create", getSnippetForm)
	mux.HandleFunc("POST /snippets", createSnippet)

	log.Print("Starting server on :4000")

	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)
}
