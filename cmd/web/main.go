package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/axbrunn/http_web/internals/models"

	_ "github.com/go-sql-driver/mysql"
)

type server struct {
	router        *http.ServeMux
	logger        *slog.Logger
	posts         *models.PostModel
	templateCache map[string]*template.Template
}

type config struct {
	addr      string
	staticDir string
	dsn       string
}

func main() {
	if err := run(); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "%s\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}

func run() error {
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.dsn, "dsn", "web:pass@/personalblog?parseTime=true", "MySQL data source name")

	flag.Parse()

	// Move db setup here
	db, err := openDB(cfg.dsn)
	if err != nil {
		return err
	}
	defer db.Close() // This will now close *after* the server exits

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	templateCache, err := newTemplateCache()
	if err != nil {
		return err
	}

	srv := &server{
		logger:        logger,
		templateCache: templateCache,
		posts:         &models.PostModel{DB: db},
	}

	logger.Info("Started server", slog.String("addr", cfg.addr))
	return http.ListenAndServe(cfg.addr, srv.routes())
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
