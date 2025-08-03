package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type server struct {
	router        *http.ServeMux
	logger        *slog.Logger
	templateCache map[string]*template.Template
}

type config struct {
	addr      string
	staticDir string
	dsn       string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	// Define a new command-line flag for the MySQL DSN string.
	flag.StringVar(&cfg.dsn, "dsn", "web:pass@/personalblog?parseTime=true", "MySQL data source name")

	flag.Parse()

	srv := newServer(cfg)

	srv.logger.Info("Started server", slog.String("addr", cfg.addr))

	return http.ListenAndServe(cfg.addr, srv.routes())
}

func newServer(cfg config) *server {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	srv := &server{
		logger:        logger,
		templateCache: templateCache,
	}

	srv.routes()

	return srv
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
