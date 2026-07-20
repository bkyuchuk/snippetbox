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
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type signupForm struct {
	Name string `form:"name"`
	loginForm
}

type loginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	templateData := app.newTemplateData(r)
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

	templateData := app.newTemplateData(r)
	templateData.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl", templateData)
}

func (app *application) getSnippetForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
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
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippets/%d", id), http.StatusSeeOther)
}

func (app *application) getSignupForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = &signupForm{}

	app.render(w, r, http.StatusOK, "signup.tmpl", data)
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	form := &signupForm{}

	err := app.decodeForm(r, form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(form.IsNotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(form.IsNotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(form.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(form.IsNotBlank(form.Password), "password", "This field must not be blank")
	form.CheckField(form.MinChars(form.Password, 8), "password", "This field must be at least 8 characters")
	form.CheckField(form.MaxBytes(form.Password, 72), "password", "This field must be at most 72 bytes")

	if !form.IsValid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) getLoginForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = &loginForm{}

	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	form := &loginForm{}
	err := app.decodeForm(r, form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(form.IsNotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(form.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(form.IsNotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(form.MaxBytes(form.Password, 72), "password", "This field must be at most 72 bytes")

	if !form.IsValid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	userId, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNoFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
			return
		}

		app.serverError(w, r, err)
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserId", userId)

	http.Redirect(w, r, "/snippets/create", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserId")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
