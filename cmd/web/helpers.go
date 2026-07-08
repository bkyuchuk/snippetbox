package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	_ "turso.tech/database/tursogo"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		stack  = string(debug.Stack())
	)

	app.logger.Error(
		err.Error(),
		slog.String("method", method),
		slog.String("uri", uri),
		slog.String("stack", stack),
	)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func initDb(dbDriver, connStr string) *sql.DB {
	DB, err := sql.Open(dbDriver, connStr)

	if err != nil {
		panic("unable to open database")
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS snippets(
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	title TEXT NOT NULL,
    	content TEXT NOT NULL,
    	created DATETIME NOT NULL,
    	expires DATETIME NOT NULL
	)`)

	if err != nil {
		panic("could not create table")
	}

	DB.SetMaxOpenConns(20)
	DB.SetMaxIdleConns(10)

	return DB
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data *templateData) {
	ts, ok := app.cache[page]
	if !ok {
		err := fmt.Errorf("template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	_, err = buf.WriteTo(w)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}
