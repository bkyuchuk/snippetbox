package main

import (
	"log"
	"net/http"
)

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
