package main

import (
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

type server struct {
	router        *http.ServeMux
	logger        *slog.Logger
	templateCache map[string]*template.Template
}

type config struct {
	addr      string
	staticDir string
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
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")

	flag.Parse()

	srv := newServer()

	srv.logger.Info("Started server", slog.String("addr", cfg.addr))

	return http.ListenAndServe(cfg.addr, srv.routes())
}

func newServer() *server {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

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
