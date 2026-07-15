package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := &chain{app.sessionManager.LoadAndSave}

	// Snippets
	mux.Handle("GET /{$}", dynamic.thenFunc(http.HandlerFunc(app.home)))
	mux.Handle("GET /snippets/{id}", dynamic.thenFunc(http.HandlerFunc(app.getSnippet)))
	mux.Handle("GET /snippets/create", dynamic.thenFunc(http.HandlerFunc(app.getSnippetForm)))
	mux.Handle("POST /snippets", dynamic.thenFunc(http.HandlerFunc(app.createSnippet)))

	// Auth
	mux.Handle("GET /users/signup", dynamic.thenFunc(http.HandlerFunc(app.getSignupForm)))
	mux.Handle("POST /users/signup", dynamic.thenFunc(http.HandlerFunc(app.signup)))
	mux.Handle("GET /users/login", dynamic.thenFunc(http.HandlerFunc(app.getLoginForm)))
	mux.Handle("POST /users/login", dynamic.then(http.HandlerFunc(app.login)))
	mux.Handle("POST /users/logout", dynamic.thenFunc(http.HandlerFunc(app.logout)))

	standard := &chain{app.recoverPanic, app.logRequest, commonHeaders}

	return standard.then(mux)
}
