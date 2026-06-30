package main

import (
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"snippetbox.bogdandev.de/internal/models"
)

type config struct {
	addr      string
	staticDir string
	dbDriver  string
	dbConnStr string
}

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
	cache    map[string]*template.Template
}

var cfg config

func main() {
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.dbDriver, "db-driver", "turso", "Name of the database driver")
	flag.StringVar(&cfg.dbConnStr, "conn-str", "./internal/database/app.db", "Database connection string")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	templateCache, err := newTemplateCache()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: initDb(cfg.dbDriver, cfg.dbConnStr)},
		cache:    templateCache,
	}

	app.logger.Info("Starting server on address", slog.String("addr", cfg.addr))

	err = http.ListenAndServe(cfg.addr, app.routes())

	app.logger.Error(err.Error())

	os.Exit(1)
}
