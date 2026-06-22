package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.staticDir))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippets/{id}", app.getSnippet)
	mux.HandleFunc("GET /snippets/create", app.getSnippetForm)
	mux.HandleFunc("POST /snippets", app.createSnippet)
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	return mux
}
