package main

import (
	"net/http"
	"testing"

	"snippetbox.bogdandev.de/internal/assert"
)

func TestPing(t *testing.T) {
	// Arrange
	testApp := newTestApplication()
	ts := newTestServer(testApp.routes())
	defer ts.Close()

	// Act
	res := ts.get(t, "/healthCheck")

	// Assert
	assert.Equal(t, res.status, http.StatusOK)
	assert.Equal(t, res.body, "OK")
}
