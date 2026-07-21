package main

import (
	"net/http"
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

	res := ts.get(t, "/users/signup")

	t.Logf("CSRF token is: %q", extractCSRFToken(t, res.body))
	t.Logf("cookies are: %q", res.cookies)

}
