package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.bogdandev.de/internal/assert"
)

func TestCommonHeaders(t *testing.T) {
	// Arrange
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Act
	commonHeaders(dummyHandler).ServeHTTP(rr, req)

	// Assert
	res := rr.Result()
	defer func(c io.ReadCloser) {
		err := c.Close()
		if err != nil {
			t.Fatal()
		}
	}(res.Body)

	expected := "default-src 'self'; style-src 'self'; fonts.googleapis.com; font-src 'fonts.gstatic.com'"
	assert.Equal(t, res.Header.Get("Content-Security-Policy"), expected)

	expected = "origin-when-cross-origin"
	assert.Equal(t, res.Header.Get("Referrer-Policy"), expected)

	expected = "nosniff"
	assert.Equal(t, res.Header.Get("X-Content-Type-Options"), expected)

	expected = "deny"
	assert.Equal(t, res.Header.Get("X-Frame-Options"), expected)

	expected = "0"
	assert.Equal(t, res.Header.Get("X-XSS-Protection"), expected)

	expected = "Go"
	assert.Equal(t, res.Header.Get("Server"), expected)

	assert.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
