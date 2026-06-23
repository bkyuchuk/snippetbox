package main

import (
	"flag"
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
}

var cfg config

func main() {
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.dbDriver, "db-driver", "turso", "Name of the database driver")
	flag.StringVar(&cfg.dbConnStr, "conn-str", "./internal/database/app.db", "Database connection string")

	flag.Parse()

	app := &application{
		logger:   slog.New(slog.NewTextHandler(os.Stdout, nil)),
		snippets: &models.SnippetModel{DB: initDb(cfg.dbDriver, cfg.dbConnStr)},
	}

	app.logger.Info("Starting server on address", slog.String("addr", cfg.addr))

	err := http.ListenAndServe(cfg.addr, app.routes())

	app.logger.Error(err.Error())

	os.Exit(1)
}
