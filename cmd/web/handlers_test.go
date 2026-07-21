package main

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"snippetbox.bogdandev.de/internal/assert"
)

func TestHealthCheck(t *testing.T) {
	// Arrange
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Act
	res := ts.get(t, "/healthCheck")

	// Assert
	assert.Equal(t, res.status, http.StatusOK)
	assert.Equal(t, res.body, "OK")
}

func TestGetSnippet(t *testing.T) {
	// Arrange
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name       string
		urlPath    string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Valid Id",
			urlPath:    "/snippets/1",
			wantStatus: http.StatusOK,
			wantBody:   "An old silent pond...",
		},
		{
			name:       "Non-existent Id",
			urlPath:    "/snippets/2",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Negative Id",
			urlPath:    "/snippets/-1",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Decimal Id",
			urlPath:    "/snippets/1.32",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "String Id",
			urlPath:    "/snippets/culture",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Empty Id",
			urlPath:    "/snippets/",
			wantStatus: http.StatusNotFound,
		},
	}

	// Act
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts.resetCookieJar(t)

			res := ts.get(t, tt.urlPath)

			// Assert
			assert.Equal(t, res.status, tt.wantStatus)
			assert.True(t, strings.Contains(res.body, tt.wantBody))
		})
	}
}

func TestUserSignup(t *testing.T) {
	// Arrange
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	const (
		validName     = "Bogdan"
		validPassword = "validPa$$word"
		validEmail    = "bogdan@example.com"
		formTag       = "<form action='/users/signup' method='POST' novalidate>"
	)

	tests := []struct {
		name              string
		userName          string
		userEmail         string
		userPassword      string
		useValidCSRFToken bool
		wantStatus        int
		wantFormTag       string
	}{
		{
			name:              "Valid Submission",
			userName:          validName,
			userEmail:         validEmail,
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusSeeOther,
		},
		{
			name:              "Invalid CSRF Token",
			userName:          validName,
			userEmail:         validEmail,
			userPassword:      validPassword,
			useValidCSRFToken: false,
			wantStatus:        http.StatusBadRequest,
		},
		{
			name:              "Empty Name",
			userName:          "",
			userEmail:         validEmail,
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Empty Email",
			userName:          validName,
			userEmail:         "",
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Empty Password",
			userName:          validName,
			userEmail:         validEmail,
			userPassword:      "",
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Invalid Email",
			userName:          validName,
			userEmail:         "bogdan@example.",
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Short Password",
			userName:          validName,
			userEmail:         validEmail,
			userPassword:      "pa$$",
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Duplicate Email",
			userName:          validName,
			userEmail:         "doe@example.com",
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts.resetCookieJar(t)

			res := ts.get(t, "/users/signup")

			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			if tt.useValidCSRFToken {
				form.Add("csrf_token", extractCSRFToken(t, res.body))
			}

			// Act
			res = ts.postForm(t, "/users/signup", form)

			// Assert
			assert.Equal(t, res.status, tt.wantStatus)
			assert.True(t, strings.Contains(res.body, tt.wantFormTag))
		})
	}
}
