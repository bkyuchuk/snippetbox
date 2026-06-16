package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox!"))
}

func viewSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Viewing a single snippet."))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Creating a snippet."))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("/view", viewSnippet)
	mux.HandleFunc("/create", createSnippet)

	log.Print("Starting server on :4000")

	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)
}
