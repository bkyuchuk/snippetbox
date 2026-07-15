package main

import (
	"crypto/tls"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"snippetbox.bogdandev.de/internal/models"
)

type config struct {
	addr      string
	staticDir string
	dbDriver  string
	dbConnStr string
}

type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	cache          map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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
	DB := initDb(cfg.dbDriver, cfg.dbConnStr)
	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(DB)
	sessionManager.Lifetime = 12 * time.Hour

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: DB},
		users:          &models.UserModel{DB: DB},
		cache:          templateCache,
		formDecoder:    form.NewDecoder(),
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:           cfg.addr,
		Handler:        app.routes(),
		ErrorLog:       slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:      tlsConfig,
		IdleTimeout:    1 * time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 524288,
	}

	app.logger.Info("Starting server on address", slog.String("addr", srv.Addr))

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	app.logger.Error(err.Error())
	os.Exit(1)
}
