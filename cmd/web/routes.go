package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := &chain{app.sessionManager.LoadAndSave}

	mux.Handle("GET /{$}", dynamic.thenFunc(http.HandlerFunc(app.home)))
	mux.Handle("GET /snippets/{id}", dynamic.thenFunc(http.HandlerFunc(app.getSnippet)))
	mux.Handle("GET /snippets/create", dynamic.thenFunc(http.HandlerFunc(app.getSnippetForm)))
	mux.Handle("POST /snippets", dynamic.thenFunc(http.HandlerFunc(app.createSnippet)))

	standard := &chain{app.recoverPanic, app.logRequest, commonHeaders}

	return standard.then(mux)
}
