package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.bogdandev.de/internal/assert"
)

func TestPing(t *testing.T) {
	// Arrange
	app := &application{logger: slog.New(slog.DiscardHandler)}

	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/healthCheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Assert
	assert.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
