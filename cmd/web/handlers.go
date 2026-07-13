package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.bogdandev.de/internal/models"
	"snippetbox.bogdandev.de/internal/validator"
)

type snippetForm struct {
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`
	validator.Validator
}

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
	data := newTemplateData()
	data.Form = snippetForm{
		Expires: 1,
	}

	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	form := &snippetForm{}

	err := app.decodeForm(r, form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(form.IsNotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(form.IsNotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(form.IsUnderMaxChars(form.Title, 100), "title", "This field must have less than 100 characters")
	form.CheckField(validator.HasPermittedValue(form.Expires, 1, 7, 365), "expires", "This field must be either 1, 7 or 365")

	if !form.IsValid() {
		data := newTemplateData()
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippets/%d", id), http.StatusSeeOther)
}
