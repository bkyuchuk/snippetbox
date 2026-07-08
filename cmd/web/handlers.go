package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.bogdandev.de/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	templateData := newTemplateData()
	templateData.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tmpl", templateData)
}

func (app *application) getSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(int64(id))

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	templateData := newTemplateData()
	templateData.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl", templateData)
}

func (app *application) getSnippetForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := fmt.Fprint(w, `{"message": "Display form for creating a snippet."}`)

	if err != nil {
		panic("could not write response")
	}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Dummy data;
	title := "An old silent pond"
	content := "An old silent pond...\\nA frog jumps into the pond,\\nsplash! Silence again.\\n\\n– Matsuo Bashō"
	expires := 1

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Redirect to relevant snippet page
	http.Redirect(w, r, fmt.Sprintf("/snippets/%d", id), http.StatusSeeOther)
}
