package main

import (
	"net/http"
)

func (srv *server) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", srv.makeHandler(srv.handleHome()))
	mux.HandleFunc("GET /posts", srv.makeHandler(srv.handlePosts()))
	mux.HandleFunc("GET /posts/{slug}", srv.makeHandler(srv.postView()))

	return commenHeaders(mux)
}
