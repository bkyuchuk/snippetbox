package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestApplication() *application {
	return &application{logger: slog.New(slog.DiscardHandler)}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	return &testServer{ts}
}

type testResponse struct {
	status  int
	headers *http.Header
	body    string
	cookies []*http.Cookie
}

func (ts *testServer) get(t *testing.T, urlPath string) *testResponse {
	req, err := http.NewRequest("GET", ts.URL+urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return &testResponse{
		status:  res.StatusCode,
		headers: &res.Header,
		body:    string(bytes.TrimSpace(body)),
		cookies: res.Cookies(),
	}
}
