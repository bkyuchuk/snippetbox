package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type config struct {
	addr      string
	staticDir string
}

type application struct {
	logger *slog.Logger
}

var cfg config

func main() {
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")

	flag.Parse()

	app := &application{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.staticDir))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippets/{id}", app.getSnippet)
	mux.HandleFunc("GET /snippets/create", app.getSnippetForm)
	mux.HandleFunc("POST /snippets", app.createSnippet)
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	app.logger.Info("Starting server on address", slog.String("addr", cfg.addr))

	err := http.ListenAndServe(cfg.addr, mux)

	app.logger.Error(err.Error())

	os.Exit(1)
}
