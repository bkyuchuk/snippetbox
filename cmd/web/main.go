package main

import (
	"flag"
	"log"
	"net/http"
)

type config struct {
	addr      string
	staticDir string
}

var cfg config

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.staticDir))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippets/{id}", getSnippet)
	mux.HandleFunc("GET /snippets/create", getSnippetForm)
	mux.HandleFunc("POST /snippets", createSnippet)
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")

	flag.Parse()

	log.Printf("Starting server on %s", cfg.addr)

	err := http.ListenAndServe(cfg.addr, mux)

	log.Fatal(err)
}
