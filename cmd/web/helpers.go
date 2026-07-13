package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/form/v4"
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

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func initDb(dbDriver, connStr string) *sql.DB {
	DB, err := sql.Open(dbDriver, connStr)

	if err != nil {
		panic("unable to open database")
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS snippets (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	title TEXT NOT NULL,
    	content TEXT NOT NULL,
    	created DATETIME NOT NULL,
    	expires DATETIME NOT NULL
	)`)

	_, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS sessions (
    	token TEXT PRIMARY KEY,
    	data BLOB NOT NULL,
    	expiry REAL NOT NULL
	);
	CREATE INDEX sessions_expiry_idx ON sessions(expiry);
	`)

	if err != nil {
		panic("could not create tables")
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

func (app *application) decodeForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		if _, ok := errors.AsType[*form.InvalidDecoderError](err); !ok {
			panic(err)
		}

		return err
	}

	return nil
}
